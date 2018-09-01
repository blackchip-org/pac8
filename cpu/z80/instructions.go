package z80

import (
	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/cpu"
)

// load
func ld(cpu *CPU, dest cpu.Out, source cpu.In) {
	dest(source())
}

// load, 16-bit
func ld16(cpu *CPU, dest cpu.Out16, source cpu.In16) {
	dest(source())
}

// decrement B and jump if not zero
func djnz(cpu *CPU, arg cpu.In) {
	delta := arg()
	cpu.B--
	if cpu.B != 0 {
		cpu.PC = bits.Displace(cpu.PC, delta)
	}
}

// exchange
func ex(cpu *CPU, ina cpu.In16, outa cpu.Out16, inb cpu.In16, outb cpu.Out16) {
	a := ina()
	b := inb()
	outa(b)
	outb(a)
}

// jump relative, conditional
func jr(cpu *CPU, arg cpu.In, flag int, condition bool) {
	delta := arg()
	if bits.Get(cpu.F, flag) == condition {
		cpu.PC = bits.Displace(cpu.PC, delta)
	}
}

// jump relative, always
func jra(cpu *CPU, arg cpu.In) {
	delta := arg()
	cpu.PC = bits.Displace(cpu.PC, delta)
}

// no operation
func nop() {}
