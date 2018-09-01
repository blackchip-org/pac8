package z80

import (
	"github.com/blackchip-org/pac8/cpu"
)

func ld(cpu *CPU, dest cpu.ModePut, source cpu.ModeGet) {
	dest(source())
}
