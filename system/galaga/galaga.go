package galaga

import (
	"time"

	"github.com/blackchip-org/pac8/app"
	"github.com/blackchip-org/pac8/component"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/proc/z80"
	"github.com/blackchip-org/pac8/machine"
)

type Galaga struct {
	spec *machine.Spec
	regs Registers
}

type Config struct {
	Name string
	M    *memory.BlockMapper
}

type Registers struct {
	InterruptEnable1 uint8 // low bit
	DipSwitches      [8]uint8
}

func New(ctx app.SDLContext, config Config) (machine.System, error) {
	sys := &Galaga{}
	ram := memory.NewRAM(0x2000)
	io := memory.NewIO(0x100)

	config.M.Map(0x6800, io)
	config.M.Map(0x8000, ram)

	mapRegisters(&sys.regs, io)

	mem1 := memory.NewPageMapped(config.M.Blocks)
	cpu1 := z80.New(mem1)

	sys.spec = &machine.Spec{
		Name:        config.Name,
		CharDecoder: GalagaDecoder,
		CPU:         cpu1,
		Mem:         mem1,
		TickCallback: func(m *machine.Mach) {
			if m.Status != machine.Run {
				return
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
	pm.WO(0x20, &r.InterruptEnable1)
}
