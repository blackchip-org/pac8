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

// Logical and
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag set.
func and(cpu *CPU, get cpu.In) {
	a1 := cpu.A
	a2 := get()

	result := a1 & a2

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, true)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagV, bits.Parity(result))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, false)

	cpu.A = result
}

// call, conditional
//
// Pushes the address after the CALL instruction (PC+3) onto the stack and
// jumps to the label. Can also take conditions.
func call(cpu *CPU, flag int, condition bool, get cpu.In16) {
	addr := get()
	if bits.Get(cpu.F, flag) == condition {
		cpu.SP -= 2
		cpu.mem16.Store(cpu.SP, cpu.PC)
		cpu.PC = addr
	}
}

// call, always
func calla(cpu *CPU, get cpu.In16) {
	addr := get()
	cpu.SP -= 2
	cpu.mem16.Store(cpu.SP, cpu.PC)
	cpu.PC = addr
}

func cb(cpu *CPU) {
	opcode := cpu.fetch()
	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	cpu.R = (cpu.R + 1) & 0x7f
	opsCB[opcode](cpu)
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

// CP is a subtraction from A that doesn't update A, only the flags it would
// have set/reset if it really was subtracted.
//
// F5 and F3 are copied from the operand, not the result
func cp(cpu *CPU, get cpu.In) {
	a1 := cpu.A
	a2 := get()
	sub(cpu, func() uint8 { return a2 }, false)
	cpu.A = a1
	bits.Set(&cpu.F, Flag3, bits.Get(a2, 3))
	bits.Set(&cpu.F, Flag5, bits.Get(a2, 5))
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

// TODO: implement
func di() {}

// decrement B and jump if not zero
func djnz(cpu *CPU, get cpu.In) {
	delta := get()
	cpu.B--
	if cpu.B != 0 {
		cpu.PC = bits.Displace(cpu.PC, delta)
	}
}

// TODO: implmenet
func ei() {}

// exchange
func ex(cpu *CPU, geta cpu.In16, puta cpu.Out16, getb cpu.In16, putb cpu.Out16) {
	a := geta()
	b := getb()
	puta(b)
	putb(a)
}

// EXX exchanges BC, DE, and HL with shadow registers with BC', DE', and HL'.
func exx(cpu *CPU) {
	ex(cpu, cpu.loadBC, cpu.storeBC, cpu.loadBC1, cpu.storeBC1)
	ex(cpu, cpu.loadDE, cpu.storeDE, cpu.loadDE1, cpu.storeDE1)
	ex(cpu, cpu.loadHL, cpu.storeHL, cpu.loadHL1, cpu.storeHL1)
}

func halt(cpu *CPU) {
	cpu.Halt = true
}

// TODO: implement
func in(cpu *CPU, put cpu.Out, get cpu.In) {
	get()
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

// jump absolute, conditional
func jp(cpu *CPU, flag int, condition bool, get cpu.In16) {
	addr := get()
	if bits.Get(cpu.F, flag) == condition {
		cpu.PC = addr
	}
}

// jump absolute, always
func jpa(cpu *CPU, get cpu.In16) {
	cpu.PC = get()
}

// jump relative, conditional
func jr(cpu *CPU, flag int, condition bool, get cpu.In) {
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

// Logical or
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag cleared.
func or(cpu *CPU, get cpu.In) {
	a1 := cpu.A
	a2 := get()

	result := a1 | a2

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagV, bits.Parity(result))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, false)

	cpu.A = result
}

// TODO: implement
func out(cpu *CPU) {
	cpu.fetch()
}

// Copies the two bytes from (SP) into the operand, then increases SP by 2.
func pop(cpu *CPU, put cpu.Out16) {
	put(cpu.mem16.Load(cpu.SP))
	cpu.SP += 2
}

// Decrements the SP by 2 then copies the operand into (SP)
func push(cpu *CPU, get cpu.In16) {
	cpu.SP -= 2
	cpu.mem16.Store(cpu.SP, get())
}

// return, conditional
func ret(cpu *CPU, flag int, value bool) {
	if bits.Get(cpu.F, flag) == value {
		reta(cpu)
	}
}

// return, always
func reta(cpu *CPU) {
	cpu.PC = cpu.mem16.Load(cpu.SP)
	cpu.SP += 2
}

// 9-bit rotation to the left.
// The register's bits are shifted left. The carry value is put into 0th bit
// of the register, and the leaving 7th bit is put into the carry. C is
// changed to the leaving 7th bit, H and N are reset, P/V is parity, S and Z
// are modified by definition.
func rl(cpu *CPU, put cpu.Out, get cpu.In) {
	arg := get()
	carry := bits.Get(arg, 7)
	result := arg << 1
	if bits.Get(cpu.F, FlagC) {
		result++
	}

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagV, bits.Parity(result))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, carry)

	put(result)
}

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
//
// 8-bit rotation to the left. The bit leaving on the left is copied into the
// carry, and to bit 0.
// H and N flags are reset, P/V is parity, S and Z are modified by definition.
func rlc(cpu *CPU, put cpu.Out, get cpu.In) {
	arg := get()
	carry := bits.Get(arg, 7)
	result := arg << 1
	if carry {
		result++
	}

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagV, bits.Parity(result))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, carry)

	put(result)
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

// 9-bit rotation to the right. The carry is copied into bit 7, and the bit
// leaving on the right is copied into the carry. Carry becomes the bit
// leaving on the right, H and N flags are reset, P/V is parity, S and Z are
// modified by definition.
func rr(cpu *CPU, put cpu.Out, get cpu.In) {
	arg := get()
	carry := bits.Get(arg, 0)
	result := arg >> 1
	if bits.Get(cpu.F, FlagC) {
		bits.Set(&result, 7, true)
	}

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagV, bits.Parity(result))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, carry)

	put(result)
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

// rotate right with carry
//
// 8-bit rotation to the right. the bit leaving on the right is copied into
// the carry, and into bit 7. The carry becomes the value leaving on the right,
// H and N are reset, P/V is parity, and S and Z are modified by definition.
func rrc(cpu *CPU, put cpu.Out, get cpu.In) {
	arg := get()
	carry := bits.Get(arg, 0)
	result := arg >> 1
	if carry {
		bits.Set(&result, 7, true)
	}

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagV, bits.Parity(result))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, carry)

	put(result)
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

func rst(cpu *CPU, y int) {
	cpu.SP -= 2
	cpu.mem16.Store(cpu.SP, cpu.PC)
	cpu.PC = uint16(y) * 8
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

// Logical exclusive or
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag cleared.
func xor(cpu *CPU, get cpu.In) {
	a1 := cpu.A
	a2 := get()

	result := a1 ^ a2

	bits.Set(&cpu.F, FlagS, bits.Get(result, 7))
	bits.Set(&cpu.F, FlagZ, result == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(result, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagV, bits.Parity(result))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, false)

	cpu.A = result
}
