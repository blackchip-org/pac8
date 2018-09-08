package z80

// http://z80-heaven.wikidot.com/instructions-set

import (
	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/cpu"
)

//go:generate go run ops/gen.go
//go:generate go fmt ops.go

var alu bits.ALU

// Add
//
// N flag is reset, P/V is interpreted as overflow.
// Rest of the flags is modified by definition.
func add(cpu *CPU, get0 cpu.Get, get1 cpu.Get, withCarry bool) {
	alu.In0 = get0()
	alu.In1 = get1()

	alu.SetCarry(false)
	if withCarry && bits.Get(cpu.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.Add()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Overflow())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	cpu.A = alu.Out
}

// add
// preserve s, z, p/v. h undefined
func add16(cpu *CPU, put cpu.Put16, get0 cpu.Get16, get1 cpu.Get16, withCarry bool) {
	in0 := get0()
	in1 := get1()

	alu.In0 = uint8(in0)
	alu.In1 = uint8(in1)
	alu.SetCarry(false)
	if withCarry && bits.Get(cpu.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.Add()
	lo := alu.Out

	alu.In0 = uint8(in0 >> 8)
	alu.In1 = uint8(in1 >> 8)
	alu.Add()
	hi := alu.Out

	result := uint16(hi)<<8 | uint16(lo)

	if withCarry {
		bits.Set(&cpu.F, FlagS, alu.Sign())
		bits.Set(&cpu.F, FlagZ, alu.Zero())
		bits.Set(&cpu.F, FlagV, alu.Overflow())
	}
	bits.Set(&cpu.F, Flag5, bits.Get(hi, 5))
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(hi, 3))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(result)
}

// Logical and
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag set.
func and(cpu *CPU, get cpu.Get) {
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

// Tests if the specified bit is set.
//
// Opposite of the nth bit is written into the Z flag. C is preserved,
// N is reset, H is set, and S and P/V are undefined.
//
// PV as Z, S set only if n=7 and b7 of r set
func bit(cpu *CPU, n int, get cpu.Get) {
	arg := get()

	bits.Set(&cpu.F, FlagS, n == 7 && bits.Get(arg, 7))
	bits.Set(&cpu.F, FlagZ, !bits.Get(arg, n))
	bits.Set(&cpu.F, Flag5, bits.Get(arg, 5))
	bits.Set(&cpu.F, FlagH, true)
	bits.Set(&cpu.F, Flag3, bits.Get(arg, 3))
	bits.Set(&cpu.F, FlagV, !bits.Get(arg, n))
	bits.Set(&cpu.F, FlagN, false)
}

// call, conditional
//
// Pushes the address after the CALL instruction (PC+3) onto the stack and
// jumps to the label. Can also take conditions.
func call(cpu *CPU, flag int, condition bool, get cpu.Get16) {
	addr := get()
	if bits.Get(cpu.F, flag) == condition {
		cpu.SP -= 2
		cpu.mem16.Store(cpu.SP, cpu.PC)
		cpu.PC = addr
	}
}

// call, always
func calla(cpu *CPU, get cpu.Get16) {
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
func cp(cpu *CPU, get cpu.Get) {
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
func dec(cpu *CPU, put cpu.Put, get cpu.Get) {
	alu.In0 = get()
	alu.SetCarry(false)
	alu.Decrement()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Overflow())
	bits.Set(&cpu.F, FlagN, true)

	put(alu.Out)
}

// decrement 16-bit
// No flags altered
func dec16(cpu *CPU, put cpu.Put16, get cpu.Get16) {
	arg := get()
	put(arg - 1)
}

// TODO: implement
func di() {}

// decrement B and jump if not zero
func djnz(cpu *CPU, get cpu.Get) {
	delta := get()
	cpu.B--
	if cpu.B != 0 {
		cpu.PC = bits.Displace(cpu.PC, delta)
	}
}

func ed(cpu *CPU) {
	opcode := cpu.fetch()
	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	cpu.R = (cpu.R + 1) & 0x7f
	opsED[opcode](cpu)
}

// TODO: implmenet
func ei() {}

// exchange
func ex(cpu *CPU, geta cpu.Get16, puta cpu.Put16, getb cpu.Get16, putb cpu.Put16) {
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
func in(cpu *CPU, put cpu.Put, get cpu.Get) {
	get()
}

// increment
// Preserves C flag, N flag is reset, P/V detects overflow and rest are
// modified by definition.
func inc(cpu *CPU, put cpu.Put, get cpu.Get) {
	alu.In0 = get()
	alu.SetCarry(false)
	alu.Increment()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Overflow())
	bits.Set(&cpu.F, FlagN, false)

	put(alu.Out)
}

// increment 16-bit
// No flags altered
func inc16(cpu *CPU, put cpu.Put16, get cpu.Get16) {
	arg := get()
	put(arg + 1)
}

func invalid() {}

// jump absolute, conditional
func jp(cpu *CPU, flag int, condition bool, get cpu.Get16) {
	addr := get()
	if bits.Get(cpu.F, flag) == condition {
		cpu.PC = addr
	}
}

// jump absolute, always
func jpa(cpu *CPU, get cpu.Get16) {
	cpu.PC = get()
}

// jump relative, conditional
func jr(cpu *CPU, flag int, condition bool, get cpu.Get) {
	delta := get()
	if bits.Get(cpu.F, flag) == condition {
		cpu.PC = bits.Displace(cpu.PC, delta)
	}
}

// jump relative, always
func jra(cpu *CPU, get cpu.Get) {
	delta := get()
	cpu.PC = bits.Displace(cpu.PC, delta)
}

// load
func ld(cpu *CPU, put cpu.Put, get cpu.Get) {
	put(get())
}

// load, 16-bit
func ld16(cpu *CPU, put cpu.Put16, get cpu.Get16) {
	put(get())
}

// no operation
func nop() {}

// Logical or
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag cleared.
func or(cpu *CPU, get cpu.Get) {
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
func pop(cpu *CPU, put cpu.Put16) {
	put(cpu.mem16.Load(cpu.SP))
	cpu.SP += 2
}

// Decrements the SP by 2 then copies the operand into (SP)
func push(cpu *CPU, get cpu.Get16) {
	cpu.SP -= 2
	cpu.mem16.Store(cpu.SP, get())
}

// Resets the specified byte to zero
func res(cpu *CPU, n int, put cpu.Put, get cpu.Get) {
	arg := get()
	bits.Set(&arg, n, false)
	put(arg)
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

func rotl(cpu *CPU, put cpu.Put, get cpu.Get) {
	alu.In0 = get()
	alu.RotateLeft()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(alu.Out)
}

// rotate A left with carry
// S,Z, and P/V are preserved, H and N flags are reset
func rotla(cpu *CPU) {
	alu.In0 = cpu.A
	alu.RotateLeft()

	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagC, alu.Carry())

	cpu.A = alu.Out
}

func rotr(cpu *CPU, put cpu.Put, get cpu.Get) {
	alu.In0 = get()
	alu.RotateRight()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(alu.Out)
}

func rotra(cpu *CPU) {
	alu.In0 = cpu.A
	alu.RotateRight()

	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagC, alu.Carry())

	cpu.A = alu.Out
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

// Sets the specified byte to one
func set(cpu *CPU, n int, put cpu.Put, get cpu.Get) {
	arg := get()
	bits.Set(&arg, n, true)
	put(arg)
}

func shiftl(cpu *CPU, put cpu.Put, get cpu.Get, withCarry bool) {
	alu.In0 = get()
	alu.SetCarry(false)
	if withCarry && bits.Get(cpu.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.ShiftLeft()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(alu.Out)
}

func shiftla(cpu *CPU) {
	alu.In0 = cpu.A
	alu.SetCarry(bits.Get(cpu.F, FlagC))
	alu.ShiftLeft()

	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagC, alu.Carry())

	cpu.A = alu.Out
}

func shiftr(cpu *CPU, put cpu.Put, get cpu.Get, withCarry bool) {
	alu.In0 = get()
	alu.SetCarry(false)
	if withCarry && bits.Get(cpu.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.ShiftRight()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(alu.Out)
}

func shiftra(cpu *CPU) {
	alu.In0 = cpu.A
	alu.SetCarry(bits.Get(cpu.F, FlagC))
	alu.ShiftRight()

	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagC, alu.Carry())

	cpu.A = alu.Out
}

func sll(cpu *CPU, put cpu.Put, get cpu.Get) {
	alu.In0 = get()
	alu.SetCarry(true)
	alu.ShiftLeft()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(alu.Out)
}

func sra(cpu *CPU, put cpu.Put, get cpu.Get) {
	alu.In0 = get()
	alu.ShiftRightSigned()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(alu.Out)
}

// Subtract
//
// N flag is reset, P/V is interpreted as overflow.
// Rest of the flags is modified by definition.
func sub(cpu *CPU, arg cpu.Get, withCarry bool) {
	alu.In0 = cpu.A
	alu.In1 = arg()

	alu.SetCarry(false)
	if withCarry && bits.Get(cpu.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.Subtract()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Overflow())
	bits.Set(&cpu.F, FlagN, true)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	cpu.A = alu.Out
}

func sub16(cpu *CPU, put cpu.Put16, arg1 cpu.Get16, arg2 cpu.Get16, withCarry bool) {
	n16 := arg1()
	m16 := arg2()

	alu.In0 = uint8(n16)
	alu.In1 = uint8(m16)
	alu.SetCarry(false)
	if withCarry && bits.Get(cpu.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.Subtract()
	lo := alu.Out

	alu.In0 = uint8(n16 >> 8)
	alu.In1 = uint8(m16 >> 8)
	alu.Subtract()
	hi := alu.Out

	result := uint16(hi)<<8 | uint16(lo)

	if withCarry {
		bits.Set(&cpu.F, FlagS, alu.Sign())
		bits.Set(&cpu.F, FlagZ, alu.Zero())
		bits.Set(&cpu.F, FlagV, alu.Overflow())
	}
	bits.Set(&cpu.F, Flag5, bits.Get(hi, 5))
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(hi, 3))
	bits.Set(&cpu.F, FlagN, true)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(result)
}

// Logical exclusive or
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag cleared.
func xor(cpu *CPU, get cpu.Get) {
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
