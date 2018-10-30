package pacman

// https://www.lomont.org/Software/Games/PacMan/PacmanEmulation.pdf

import (
	"encoding/gob"
	"fmt"
	"time"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/machine"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/proc/z80"
	"github.com/veandco/go-sdl2/sdl"
)

type Pacman struct {
	mem       memory.Memory
	cpu       *z80.CPU
	regs      Registers
	intSelect uint8 // value sent during interrupt to select vector (port 0)
	tiles     *sdl.Texture
}

type Config struct {
	M        *memory.BlockMapper
	VideoROM VideoROM
}

type Registers struct {
	In0             uint8 // joystick and coin slot
	InterruptEnable uint8
	SoundEnable     uint8
	Unknown0        uint8
	FlipScreen      uint8
	Player1Lamp     uint8
	Player2Lamp     uint8
	CoinLockout     uint8
	CoinCounter     uint8
	In1             uint8 // joystick and start buttons
	Voices          [3]Voice
	DipSwitches     uint8
	WatchdogReset   uint8
}

type Voice struct {
	Acc      [5]uint8
	Waveform uint8
	Freq     [5]uint8
	Vol      uint8
}

func New(renderer *sdl.Renderer, config Config) (*machine.Mach, error) {
	sys := &Pacman{}

	ram := memory.NewRAM(0x1000)
	io := memory.NewIO(0x100)

	config.M.Map(0x4000, ram)
	config.M.Map(0x5000, io)
	// Pacman is missing address line A15 so an access to $c000 is the
	// same as accessing $4000. Ms. Pacman has additional ROMs in high
	// memory so it has an A15 line but it appears to have the RAM mapped at
	// $c000 as well. Text for HIGH SCORE and CREDIT accesses this high memory
	// when writing to video memory. Copy protection?
	config.M.Map(0xc000, ram)

	sys.mem = memory.NewPageMapped(config.M.Blocks)
	sys.cpu = z80.New(sys.mem)
	m := machine.New(sys.mem, sys.cpu)

	video, err := NewVideo(renderer, sys.mem, config.VideoROM)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize video: %v", err)
	}
	m.Display = video
	keyboard := NewKeyboard(&sys.regs)
	m.UI = keyboard

	mapRegisters(&sys.regs, io, video)

	// Port 0 gets set with the partial interrupt pointer to be set
	// by the interrupting device
	sys.cpu.Map.WO(0, &sys.intSelect)

	// FIXME: this turns the joystick "off", etc.
	sys.regs.In0 = 0x3f
	sys.regs.In1 = 0x7f

	bits.Set(&sys.regs.In1, 4, true)          // Board test switch disabled
	bits.Set(&sys.regs.In1, 7, true)          // Upright sysinet
	bits.Set(&sys.regs.DipSwitches, 0, true)  // 1 coin per game
	bits.Set(&sys.regs.DipSwitches, 1, false) // ...
	bits.Set(&sys.regs.DipSwitches, 3, true)  // 3 lives
	bits.Set(&sys.regs.DipSwitches, 7, true)  // Normal ghost names

	// 16.67 milliseconds for VBLANK interrupt
	m.TickRate = time.Duration(16670 * time.Microsecond)
	video.Callback = func() {
		if m.Status == machine.Run && sys.regs.InterruptEnable != 0 {
			sys.cpu.INT(sys.intSelect)
		}
	}

	return m, nil
}

func mapRegisters(r *Registers, io memory.IO, v *Video) {
	pm := memory.NewPortMapper(io)
	for i := 0; i <= 0x3f; i++ {
		pm.RO(i, &r.In0)
	}
	pm.WO(0x00, &r.InterruptEnable)
	pm.WO(0x01, &r.SoundEnable)
	pm.WO(0x02, &r.Unknown0)
	pm.RW(0x03, &r.FlipScreen)
	pm.RW(0x04, &r.Player1Lamp)
	pm.RW(0x05, &r.Player2Lamp)
	pm.RW(0x06, &r.CoinLockout)
	pm.RW(0x07, &r.CoinCounter)
	for i := 0x40; i <= 0x7f; i++ {
		pm.RO(i, &r.In1)
	}
	for i, v := 0x40, 0; v < 3; i, v = i+6, v+1 {
		pm.WO(i+0, &r.Voices[v].Acc[0])
		pm.WO(i+1, &r.Voices[v].Acc[1])
		pm.WO(i+2, &r.Voices[v].Acc[2])
		pm.WO(i+3, &r.Voices[v].Acc[3])
		pm.WO(i+4, &r.Voices[v].Acc[4])
		pm.WO(i+5, &r.Voices[v].Waveform)
	}
	for i, v := 0x50, 0; v < 3; i, v = i+6, v+1 {
		pm.WO(i+0, &r.Voices[v].Freq[0])
		pm.WO(i+1, &r.Voices[v].Freq[1])
		pm.WO(i+2, &r.Voices[v].Freq[2])
		pm.WO(i+3, &r.Voices[v].Freq[3])
		pm.WO(i+4, &r.Voices[v].Freq[4])
		pm.WO(i+5, &r.Voices[v].Vol)
	}
	for i, s := 0x60, 0; s < 8; i, s = i+2, s+1 {
		pm.WO(i+0, &v.spriteCoords[s].x)
		pm.WO(i+1, &v.spriteCoords[s].y)
	}
	for i := 0x80; i <= 0xbf; i++ {
		pm.RO(i, &r.DipSwitches)
	}
	for i := 0xc0; i <= 0xff; i++ {
		pm.WO(i, &r.WatchdogReset)
	}
}

func (p *Pacman) Encode(en *gob.Encoder) error {
	if err := en.Encode(p.regs); err != nil {
		return err
	}
	if err := en.Encode(p.cpu); err != nil {
		return err
	}
	return nil
}
