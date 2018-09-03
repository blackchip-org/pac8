package z80

// http://z80-heaven.wikidot.com/instructions-set

import (
	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/cpu"
)

//go:generate go run ops/gen.go
//go:generate go fmt ops.go

// Add
//
// N flag is reset, P/V is interpreted as overflow.
// Rest of the flags is modified by definition.
func add(cpu *CPU, arg1 cpu.In, arg2 cpu.In, carry bool) {
	a1 := arg1()
	a2 := arg2()
	c := uint8(0)

	if carry {
		c = 1
	}

	result16 := uint16(a1) + uint16(a2) + uint16(c)
	result := uint8(result16)
	hresult := a1&0xf + a2&0xf + c

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, bits.Get(hresult, 4))
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagV, bits.Overflow(a1, a2, result))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, bits.Get16(result16, 8))

	cpu.A = result
}

// add
// preserve s, z, p/v. h undefined
func add16(cpu *CPU, put cpu.Out16, arg1 cpu.In16, arg2 cpu.In16) {
	a1 := arg1()
	a2 := arg2()

	result := uint32(a1) + uint32(a2)
	hresult := uint8(bits.Slice16(a1, 8, 11) + bits.Slice16(a2, 8, 11))

	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, bits.Get32(result, 16))
	bits.Set(&cpu.F, FlagH, bits.Get(hresult, 4))
	bits.Set(&cpu.F, Flag3, bits.Get32(result, 11))
	bits.Set(&cpu.F, Flag5, bits.Get32(result, 13))

	put(uint16(result))
}

// Inverts the carry flag
//
// Carry flag inverted. Also inverts H and clears N. Rest of the flags are
// preserved.
func ccf(cpu *CPU) {
	bits.Set(&cpu.F, FlagC, !bits.Get(cpu.F, FlagC))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, Flag3, bits.Get(cpu.A, 3))
	bits.Set(&cpu.F, Flag5, bits.Get(cpu.A, 5))
}

// inverts all bits of A
//
// Sets H and N, other flags are unmodified.
func cpl(cpu *CPU) {
	cpu.A ^= 0xff
	bits.Set(&cpu.F, FlagH, true)
	bits.Set(&cpu.F, FlagN, true)
	bits.Set(&cpu.F, Flag3, bits.Get(cpu.A, 3))
	bits.Set(&cpu.F, Flag5, bits.Get(cpu.A, 5))
}

// decimal adjust in a
//
// When this instruction is executed, the A register is BCD corrected using
// the contents of the flags. The exact process is the following: if the
// least significant four bits of A contain a non-BCD digit (i. e. it is
// greater than 9) or the H flag is set, then $06 is added to the register.
// Then the four most significant bits are checked. If this more significant
// digit also happens to be greater than 9 or the C flag is set, then $60
// is added.
//
// If the second addition was needed, the C flag is set after execution,
// otherwise it is reset. The N flag is preserved, P/V is parity and the
// others are altered by definition.
//
// https://stackoverflow.com/questions/13572638/z80-daa-flags-affected
//
// Note: some documentation omits that the adjustment is negative when the
// N flag is set.
func daa(cpu *CPU) {
	result := cpu.A

	half := false
	carry := false
	if bits.Get(cpu.F, FlagN) {
		if bits.Get(cpu.F, FlagH) || cpu.A&0xf > 9 {
			result -= 0x06
			if result < 6 {
				half = true
			}
		}
		if bits.Get(cpu.F, FlagC) || cpu.A > 0x99 {
			result -= 0x60
			carry = true
		}
	} else {
		if bits.Get(cpu.F, FlagH) || cpu.A&0xf > 9 {
			result += 0x06
			half = true
		}
		if bits.Get(cpu.F, FlagC) || cpu.A > 0x99 {
			result += 0x60
			carry = true
		}
	}

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, half)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagV, bits.Parity(result))
	bits.Set(&cpu.F, FlagC, carry)

	cpu.A = result
}

// decrement
// C flag preserved, P/V detects overflow and rest modified by definition.
// modified by definition.
func dec(cpu *CPU, put cpu.Out, get cpu.In) {
	arg := get()

	result := arg - 1
	hresult := arg&0xf - 1

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, bits.Get(hresult, 4))
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

func halt(cpu *CPU) {
	cpu.Halt = true
}

// increment
// Preserves C flag, N flag is reset, P/V detects overflow and rest are
// modified by definition.
func inc(cpu *CPU, put cpu.Out, get cpu.In) {
	arg := get()

	result := arg + 1
	hresult := arg&0xf + 1

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, bits.Get(hresult, 4))
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
//
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

// Set carry flag
//
// Carry flag set, H and N cleared, rest are preserved.
func scf(cpu *CPU) {
	bits.Set(&cpu.F, FlagC, true)
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, Flag5, bits.Get(cpu.A, 5))
	bits.Set(&cpu.F, Flag3, bits.Get(cpu.A, 3))
}

// Subtract
//
// N flag is reset, P/V is interpreted as overflow.
// Rest of the flags is modified by definition.
func sub(cpu *CPU, arg cpu.In, carry bool) {
	add(cpu, cpu.loadA, func() uint8 { return ^arg() }, !carry)
	bits.Set(&cpu.F, FlagN, true)
	bits.Set(&cpu.F, FlagC, !bits.Get(cpu.F, FlagC))
	bits.Set(&cpu.F, FlagH, !bits.Get(cpu.F, FlagH))
}
