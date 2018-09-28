package cabs

import (
	"fmt"

	"os"

	"github.com/blackchip-org/pac8/cpu/z80"
	"github.com/blackchip-org/pac8/mach"
	"github.com/blackchip-org/pac8/memory"
)

type Pacman struct {
	mem memory.Memory
	cpu *z80.CPU
}

func NewPacman() *Pacman {
	cab := &Pacman{}

	e := []error{}
	rom0 := memory.LoadROM(&e, "pacman/pacman.6e", "e87e059c5be45753f7e9f33dff851f16d6751181")
	rom1 := memory.LoadROM(&e, "pacman/pacman.6f", "674d3a7f00d8be5e38b1fdc208ebef5a92d38329")
	rom2 := memory.LoadROM(&e, "pacman/pacman.6h", "8e47e8c2c4d6117d174cdac150392042d3e0a881")
	rom3 := memory.LoadROM(&e, "pacman/pacman.6j", "d4a70d56bb01d27d094d73db8667ffb00ca69cb9")
	ram := memory.NewRAM(0x1000)
	io := memory.NewRAM(0x100)

	for _, err := range e {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	if len(e) > 0 {
		os.Exit(1)
	}

	cab.mem = memory.NewPageMapped([]memory.Memory{
		rom0, // $0000 - $0fff
		rom1, // $1000 - $1fff
		rom2, // $2000 - $2fff
		rom3, // $3000 - $3fff
		ram,  // $4000 - $4fff
		io,   // $5000 - $50ff
	})
	cab.cpu = z80.New(cab.mem)

	return cab
}

func (c *Pacman) Mach() mach.Mach {
	return mach.Mach{
		Mem: c.mem,
		CPU: c.cpu,
	}
}
