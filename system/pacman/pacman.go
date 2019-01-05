package pacman

// https://www.lomont.org/Software/Games/PacMan/PacmanEmulation.pdf

import (
	"fmt"
	"time"

	"github.com/blackchip-org/pac8/pkg/machine"
	"github.com/blackchip-org/pac8/pkg/memory"
	"github.com/blackchip-org/pac8/pkg/namco"
	"github.com/blackchip-org/pac8/pkg/pac8"
	"github.com/blackchip-org/pac8/pkg/proc"
	"github.com/blackchip-org/pac8/pkg/util/bits"
	"github.com/blackchip-org/pac8/pkg/util/state"
	"github.com/blackchip-org/pac8/pkg/z80"
	"github.com/veandco/go-sdl2/sdl"
)

type Pacman struct {
	spec      *machine.Spec
	regs      *Registers
	intSelect uint8 // value sent during interrupt to select vector (port 0)
	tiles     *sdl.Texture
}

type Config struct {
	Name string
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
	DipSwitches     uint8
	WatchdogReset   uint8
}

func New(env pac8.Env, config Config, roms memory.Set) (machine.System, error) {
	sys := &Pacman{
		regs: &Registers{},
	}

	ram := memory.NewRAM(0x1000)
	io := memory.NewIO(0x100)

	m := memory.NewBlockMapper()
	m.Map(0x0000, roms["code"])
	m.Map(0x4000, ram)
	m.Map(0x5000, io)
	if code2, ok := roms["code2"]; ok {
		m.Map(0x8000, code2)
	}
	// Pacman is missing address line A15 so an access to $c000 is the
	// same as accessing $4000. Ms. Pacman has additional ROMs in high
	// memory so it has an A15 line but it appears to have the RAM mapped at
	// $c000 as well. Text for HIGH SCORE and CREDIT accesses this high memory
	// when writing to video memory. Copy protection?
	m.Map(0xc000, ram)

	mem := memory.NewPageMapped(m.Blocks)
	cpu := z80.New(mem)

	video, err := NewVideo(env.Renderer, mem, roms)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize video: %v", err)
	}
	audio, err := NewAudio(env.AudioSpec, roms)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize audio: %v", err)
	}
	mapRegisters(sys.regs, io, video, audio)

	// Port 0 gets set with the partial interrupt pointer to be set
	// by the interrupting device
	pm := memory.NewPortMapper(cpu.Ports)
	pm.WO(0, &sys.intSelect)

	// FIXME: this turns the joystick "off", etc.
	// Game does not work unless this is set!
	sys.regs.In0 = 0x3f
	sys.regs.In1 = 0x7f

	bits.Set(&sys.regs.In0, 7, true)          // Service button released
	bits.Set(&sys.regs.In1, 4, true)          // Board test switch disabled
	bits.Set(&sys.regs.In1, 7, true)          // Upright cabinet
	bits.Set(&sys.regs.DipSwitches, 0, true)  // 1 coin per game
	bits.Set(&sys.regs.DipSwitches, 1, false) // ...
	bits.Set(&sys.regs.DipSwitches, 3, true)  // 3 lives
	bits.Set(&sys.regs.DipSwitches, 7, true)  // Normal ghost names

	sys.spec = &machine.Spec{
		Name:        config.Name,
		CharDecoder: PacmanDecoder,
		CPU:         []proc.CPU{cpu},
		Mem:         []memory.Memory{mem},
		Display:     video,
		Audio:       audio,
		TickCallback: func(m *machine.Mach) {
			if m.Status != machine.Run {
				return
			}
			if sys.regs.InterruptEnable != 0 {
				cpu.INT(sys.intSelect)
			}
			sys.handleInput(m)
		},
		TickRate: time.Duration(16670 * time.Microsecond),
	}
	return sys, nil
}

func (p Pacman) Spec() *machine.Spec {
	return p.spec
}

func mapRegisters(r *Registers, io memory.IO, v *namco.Video, a *Audio) {
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

	pm.WO(0x40, &a.Voices[0].Acc[0])
	pm.WO(0x41, &a.Voices[0].Acc[1])
	pm.WO(0x42, &a.Voices[0].Acc[2])
	pm.WO(0x43, &a.Voices[0].Acc[3])
	pm.WO(0x44, &a.Voices[0].Acc[4])
	pm.WO(0x45, &a.Voices[0].Waveform)

	pm.WO(0x46, &a.Voices[1].Acc[0])
	pm.WO(0x47, &a.Voices[1].Acc[1])
	pm.WO(0x48, &a.Voices[1].Acc[2])
	pm.WO(0x49, &a.Voices[1].Acc[3])
	pm.WO(0x4a, &a.Voices[1].Waveform)

	pm.WO(0x4b, &a.Voices[2].Acc[0])
	pm.WO(0x4c, &a.Voices[2].Acc[1])
	pm.WO(0x4d, &a.Voices[2].Acc[2])
	pm.WO(0x4e, &a.Voices[2].Acc[3])
	pm.WO(0x4f, &a.Voices[2].Waveform)

	pm.WO(0x50, &a.Voices[0].Freq[0])
	pm.WO(0x51, &a.Voices[0].Freq[1])
	pm.WO(0x52, &a.Voices[0].Freq[2])
	pm.WO(0x53, &a.Voices[0].Freq[3])
	pm.WO(0x54, &a.Voices[0].Freq[4])
	pm.WO(0x55, &a.Voices[0].Vol)

	pm.WO(0x56, &a.Voices[1].Freq[0])
	pm.WO(0x57, &a.Voices[1].Freq[1])
	pm.WO(0x58, &a.Voices[1].Freq[2])
	pm.WO(0x59, &a.Voices[1].Freq[3])
	pm.WO(0x5a, &a.Voices[1].Vol)

	pm.WO(0x5b, &a.Voices[2].Freq[0])
	pm.WO(0x5c, &a.Voices[2].Freq[1])
	pm.WO(0x5d, &a.Voices[2].Freq[2])
	pm.WO(0x5e, &a.Voices[2].Freq[3])
	pm.WO(0x5f, &a.Voices[2].Vol)

	for i, s := 0x60, 0; s < 8; i, s = i+2, s+1 {
		pm.WO(i+0, &v.SpriteCoords[s].X)
		pm.WO(i+1, &v.SpriteCoords[s].Y)
	}
	for i := 0x80; i <= 0xbf; i++ {
		pm.RO(i, &r.DipSwitches)
	}
	for i := 0xc0; i <= 0xff; i++ {
		pm.WO(i, &r.WatchdogReset)
	}
}

func (p *Pacman) Save(enc *state.Encoder) {
	p.spec.CPU[0].Save(enc)
	p.spec.Mem[0].Save(enc)
	enc.Encode(p.regs)
}

func (p *Pacman) Restore(dec *state.Decoder) {
	p.spec.CPU[0].Restore(dec)
	p.spec.Mem[0].Restore(dec)
	dec.Decode(&p.regs)
}

func (p *Pacman) handleInput(m *machine.Mach) {
	bits.Set(&p.regs.In0, 0, !m.In.Joysticks[0].Up)
	bits.Set(&p.regs.In0, 1, !m.In.Joysticks[0].Left)
	bits.Set(&p.regs.In0, 2, !m.In.Joysticks[0].Right)
	bits.Set(&p.regs.In0, 3, !m.In.Joysticks[0].Down)
	bits.Set(&p.regs.In0, 5, m.In.CoinSlot[0].Active)

	// No second joystick
	bits.Set(&p.regs.In1, 0, true)
	bits.Set(&p.regs.In1, 1, true)
	bits.Set(&p.regs.In1, 2, true)
	bits.Set(&p.regs.In1, 3, true)
	bits.Set(&p.regs.In1, 5, !m.In.PlayerStart[0].Active)
	bits.Set(&p.regs.In1, 6, !m.In.PlayerStart[1].Active)
}
