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
	reg  *Registers
}

type Config struct {
	Name string
	M    *memory.BlockMapper
}

type Registers struct {
	DipSwitches [8]uint8
}

func New(ctx app.SDLContext, config Config) (machine.System, error) {
	sys := &Galaga{}

	ram := memory.NewRAM(0x2000)
	io := memory.NewIO(0x100)

	config.M.Map(0x6800, io)
	config.M.Map(0x8000, ram)

	mem0 := memory.NewPageMapped(config.M.Blocks)
	cpu0 := z80.New(mem0)
	sys.spec = &machine.Spec{
		Name:     config.Name,
		CPU:      cpu0,
		Mem:      mem0,
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
