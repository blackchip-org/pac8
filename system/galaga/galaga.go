package galaga

import (
	"fmt"
	"time"

	"github.com/blackchip-org/pac8/app"
	"github.com/blackchip-org/pac8/component"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/namco"
	"github.com/blackchip-org/pac8/component/proc"
	"github.com/blackchip-org/pac8/component/proc/z80"
	"github.com/blackchip-org/pac8/machine"
)

type Galaga struct {
	spec *machine.Spec
	regs Registers
}

type Config struct {
	Name     string
	ProcROM  [3][]memory.Memory
	VideoROM namco.VideoROM
}

type Registers struct {
	InterruptEnable0 uint8 // low bit
	InterruptEnable1 uint8 // low bit
	InterruptEnable2 uint8 // low bit
	DipSwitches      [8]uint8
}

func New(ctx app.SDLContext, config Config) (machine.System, error) {
	sys := &Galaga{}
	ram := memory.NewRAM(0x2000)
	io := memory.NewIO(0x100)

	mem := make([]memory.Memory, 3, 3)
	for i := 0; i < 3; i++ {
		m := memory.NewBlockMapper()
		m.Map(0x0000, config.ProcROM[i][0])
		m.Map(0x1000, config.ProcROM[i][1])
		m.Map(0x2000, config.ProcROM[i][2])
		m.Map(0x3000, config.ProcROM[i][3])
		m.Map(0x6800, io)
		m.Map(0x8000, ram)
		mem[i] = memory.NewPageMapped(m.Blocks)
	}

	mapRegisters(&sys.regs, io)

	cpu0 := z80.New(mem[0])
	cpu1 := z80.New(mem[1])
	cpu2 := z80.New(mem[2])

	video, err := NewVideo(ctx.Renderer, mem[0], config.VideoROM)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize video: %v", err)
	}

	sys.spec = &machine.Spec{
		Name:        config.Name,
		CharDecoder: GalagaDecoder,
		CPU:         []proc.CPU{cpu0, cpu1, cpu2},
		Mem:         mem,
		Display:     video,
		TickCallback: func(m *machine.Mach) {
			if m.Status != machine.Run {
				return
			}
			if sys.regs.InterruptEnable0 != 0 {
				cpu0.INT(0)
			}
			if sys.regs.InterruptEnable1 != 0 {
				cpu1.INT(0)
			}
		},
		TickRate: time.Duration(16670 * time.Microsecond),
	}
	return sys, nil
}

func (g Galaga) Spec() *machine.Spec {
	return g.spec
}

func (g *Galaga) Save(enc component.Encoder) error {
	return nil
}

func (g *Galaga) Restore(dec component.Decoder) error {
	return nil
}

func mapRegisters(r *Registers, io memory.IO) {
	pm := memory.NewPortMapper(io)
	pm.RW(0x20, &r.InterruptEnable0)
	pm.RW(0x21, &r.InterruptEnable1)
}
