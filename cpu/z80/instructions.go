package z80

import (
	"github.com/blackchip-org/pac8/cpu"
)

func ld(cpu *CPU, dest cpu.Out, source cpu.In) {
	dest(source())
}

func ld16(cpu *CPU, dest cpu.Out16, source cpu.In16) {
	dest(source())
}
