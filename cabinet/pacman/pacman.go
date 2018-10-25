package pacman

// https://www.lomont.org/Software/Games/PacMan/PacmanEmulation.pdf

import (
	"fmt"
	"time"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/cpu/z80"
	"github.com/blackchip-org/pac8/mach"
	"github.com/blackchip-org/pac8/memory"
	"github.com/veandco/go-sdl2/sdl"
)

type Pacman struct {
	mem       memory.Memory
	cpu       *z80.CPU
	regs      registers
	intSelect uint8 // value sent during interrupt to select vector (port 0)
	tiles     *sdl.Texture
}

type registers struct {
	in0             uint8 // joystick and coin slot
	interruptEnable uint8
	soundEnable     uint8
	unknown0        uint8
	flipScreen      uint8
	player1Lamp     uint8
	player2Lamp     uint8
	coinLockout     uint8
	coinCounter     uint8
	in1             uint8 // joystick and start buttons
	voices          [3]voice
	dipSwitches     uint8
	watchdogReset   uint8
}

type voice struct {
	acc      [5]uint8
	waveform uint8
	freq     [5]uint8
	vol      uint8
}

// FIXME: Parameters are ugly
func New(renderer *sdl.Renderer, mem memory.Memory, vroms VideoROM, io memory.IO) (*mach.Mach, error) {
	cab := &Pacman{}
	cab.mem = mem
	cab.cpu = z80.New(cab.mem)
	m := mach.New(cab.mem, cab.cpu)

	video, err := NewVideo(renderer, cab.mem, vroms)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize video: %v", err)
	}
	m.Display = video
	keyboard := NewKeyboard(&cab.regs)
	m.UI = keyboard

	mapRegisters(&cab.regs, io, video)

	// Port 0 gets set with the partial interrupt pointer to be set
	// by the interrupting device
	cab.cpu.Map.WO(0, &cab.intSelect)

	bits.Set(&cab.regs.in1, 4, true)          // Board test switch disabled
	bits.Set(&cab.regs.dipSwitches, 0, true)  // 1 coin per game
	bits.Set(&cab.regs.dipSwitches, 1, false) // ...
	bits.Set(&cab.regs.dipSwitches, 3, true)  // 3 lives
	bits.Set(&cab.regs.dipSwitches, 7, true)  // Normal ghost names

	// Different type of crash when these are set
	cab.regs.in0 = 0x3f
	cab.regs.in1 = 0x7f

	// 16.67 milliseconds for VBLANK interrupt
	m.TickRate = time.Duration(16670 * time.Microsecond)
	video.Callback = func() {
		if m.Status == mach.Run && cab.regs.interruptEnable != 0 {
			cab.cpu.INT(cab.intSelect)
		}
	}

	return m, nil
}

func mapRegisters(r *registers, io memory.IO, v *Video) {
	pm := memory.NewPortMapper(io)
	for i := 0; i <= 0x3f; i++ {
		pm.RO(i, &r.in0)
	}
	pm.WO(0x00, &r.interruptEnable)
	pm.WO(0x01, &r.soundEnable)
	pm.WO(0x02, &r.unknown0)
	pm.RW(0x03, &r.flipScreen)
	pm.RW(0x04, &r.player1Lamp)
	pm.RW(0x05, &r.player2Lamp)
	pm.RW(0x06, &r.coinLockout)
	pm.RW(0x07, &r.coinCounter)
	for i := 0x40; i <= 0x7f; i++ {
		pm.RO(i, &r.in1)
	}
	for i, v := 0x40, 0; v < 3; i, v = i+6, v+1 {
		pm.WO(i+0, &r.voices[v].acc[0])
		pm.WO(i+1, &r.voices[v].acc[1])
		pm.WO(i+2, &r.voices[v].acc[2])
		pm.WO(i+3, &r.voices[v].acc[3])
		pm.WO(i+4, &r.voices[v].acc[4])
		pm.WO(i+5, &r.voices[v].waveform)
	}
	for i, v := 0x50, 0; v < 3; i, v = i+6, v+1 {
		pm.WO(i+0, &r.voices[v].freq[0])
		pm.WO(i+1, &r.voices[v].freq[1])
		pm.WO(i+2, &r.voices[v].freq[2])
		pm.WO(i+3, &r.voices[v].freq[3])
		pm.WO(i+4, &r.voices[v].freq[4])
		pm.WO(i+5, &r.voices[v].vol)
	}
	for i, s := 0x60, 0; s < 8; i, s = i+2, s+1 {
		pm.WO(i+0, &v.spriteCoords[s].x)
		pm.WO(i+1, &v.spriteCoords[s].y)
	}
	for i := 0x80; i <= 0xbf; i++ {
		pm.RO(i, &r.dipSwitches)
	}
	for i := 0xc0; i <= 0xff; i++ {
		pm.WO(i, &r.watchdogReset)
	}
}
