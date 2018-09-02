package z80

// http://z80-heaven.wikidot.com/instructions-set

import (
	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/cpu"
)

//go:generate go run ops/gen.go
//go:generate go fmt ops.go

func add4(arg1 uint8, arg2 uint8, carry uint8) (uint8, uint8) {
	result := arg1 + arg2 + carry
	carry = 0
	if result > 0x0f {
		carry = 1
	}
	return result & 0x0f, carry
}

// add
// preserve s, z, p/v. h undefined
func add16(cpu *CPU, put cpu.Out16, arg1 cpu.In16, arg2 cpu.In16) {
	a1 := arg1()
	a2 := arg2()

	r1, c1 := add4(uint8(a1)&0x0f, uint8(a2)&0x0f, 0)
	r2, c2 := add4(uint8(a1>>4)&0x0f, uint8(a2>>4)&0x0f, c1)
	lo := uint8(r1) + uint8(r2<<4)

	r3, c3 := add4(uint8(a1>>8)&0x0f, uint8(a2>>8)&0x0f, c2)
	r4, c4 := add4(uint8(a1>>12)&0x0f, uint8(a2>>12)&0x0f, c3)
	hi := uint8(r3) + uint8(r4<<4)

	result := uint16(lo) + uint16(hi)<<8

	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, c4 == 1)
	bits.Set(&cpu.F, FlagH, c3 == 1)
	bits.Set(&cpu.F, Flag3, bits.Get(hi, 3))
	bits.Set(&cpu.F, Flag5, bits.Get(hi, 5))

	put(result)
}

// decrement
// C flag preserved, P/V detects overflow and rest modified by definition.
// modified by definition.
func dec(cpu *CPU, put cpu.Out, get cpu.In) {
	arg := get()

	r1, c1 := add4(arg&0x0f, 0xf, 0)
	r2, _ := add4((arg>>4)&0xf, 0xf, c1)
	result := r1 + (r2 << 4)

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, c1 == 0)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagV, bits.Overflow(arg, 0xff, result))
	bits.Set(&cpu.F, FlagN, true)

	put(result)
}

// decrement 16-bit
// No flags altered
func dec16(cpu *CPU, put cpu.Out16, get cpu.In16) {
	arg := get()
	put(arg - 1)
}

// decrement B and jump if not zero
func djnz(cpu *CPU, get cpu.In) {
	delta := get()
	cpu.B--
	if cpu.B != 0 {
		cpu.PC = bits.Displace(cpu.PC, delta)
	}
}

// exchange
func ex(cpu *CPU, geta cpu.In16, puta cpu.Out16, getb cpu.In16, putb cpu.Out16) {
	a := geta()
	b := getb()
	puta(b)
	putb(a)
}

// increment
// Preserves C flag, N flag is reset, P/V detects overflow and rest are
// modified by definition.
func inc(cpu *CPU, put cpu.Out, get cpu.In) {
	arg := get()

	r1, c1 := add4(arg&0x0f, 1, 0)
	r2, _ := add4((arg>>4)&0xf, 0, c1)
	result := r1 + (r2 << 4)

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, c1 == 1)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagV, bits.Overflow(arg, 1, result))
	bits.Set(&cpu.F, FlagN, false)

	put(result)
}

// increment 16-bit
// No flags altered
func inc16(cpu *CPU, put cpu.Out16, get cpu.In16) {
	arg := get()
	put(arg + 1)
}

// jump relative, conditional
func jr(cpu *CPU, get cpu.In, flag int, condition bool) {
	delta := get()
	if bits.Get(cpu.F, flag) == condition {
		cpu.PC = bits.Displace(cpu.PC, delta)
	}
}

// jump relative, always
func jra(cpu *CPU, get cpu.In) {
	delta := get()
	cpu.PC = bits.Displace(cpu.PC, delta)
}

// load
func ld(cpu *CPU, put cpu.Out, get cpu.In) {
	put(get())
}

// load, 16-bit
func ld16(cpu *CPU, put cpu.Out16, get cpu.In16) {
	put(get())
}

// no operation
func nop() {}

// 9-bit rotation to the left
// Performs an RL A, but is much faster and S, Z, and P/V flags are preserved.
// The carry value is put into 0th bit of the register, and the leaving
// 7th bit is put into the carry.
func rla(cpu *CPU) {
	arg := cpu.A
	carry := bits.Get(arg, 7)
	result := arg << 1
	if bits.Get(cpu.F, FlagC) {
		result++
	}

	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, carry)

	cpu.A = result
}

// rotate A left with carry
// S,Z, and P/V are preserved, H and N flags are reset
func rlca(cpu *CPU) {
	arg := cpu.A
	carry := bits.Get(arg, 7)
	result := arg << 1
	if carry {
		result++
	}

	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, carry)

	cpu.A = result
}

// 9-bit rotation to the right.
// The Carry becomes the bit leaving on the right, H, N flags are reset,
// P/V, S, and Z are preserved.
func rra(cpu *CPU) {
	arg := cpu.A
	carry := bits.Get(arg, 0)
	result := arg >> 1
	if bits.Get(cpu.F, FlagC) {
		bits.Set(&result, 7, true)
	}

	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, carry)

	cpu.A = result
}

// rotate A right with carry
// The carry becomes the value leaving on the right, H and N are reset,
// P/V, S, and Z are preserved.
func rrca(cpu *CPU) {
	arg := cpu.A
	carry := bits.Get(arg, 0)
	result := arg >> 1
	if carry {
		bits.Set(&result, 7, true)
	}

	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, carry)

	cpu.A = result
}
