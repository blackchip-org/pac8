package z80

// http://z80-heaven.wikidot.com/instructions-set
// https://www.worldofspectrum.org/faq/reference/z80reference.htm

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
	alu.In0 = cpu.A
	alu.In1 = get()
	alu.And()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, true)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, false)

	cpu.A = alu.Out
}

// Tests if the specified bit is set.
//
// Opposite of the nth bit is written into the Z flag. C is preserved,
// N is reset, H is set, and S and P/V are undefined.
//
// PV as Z, S set only if n=7 and b7 of r set
func bit(cpu *CPU, n int, get cpu.Get) {
	in0 := get()

	bits.Set(&cpu.F, FlagS, n == 7 && bits.Get(in0, 7))
	bits.Set(&cpu.F, FlagZ, !bits.Get(in0, n))
	bits.Set(&cpu.F, Flag5, bits.Get(in0, 5))
	bits.Set(&cpu.F, FlagH, true)
	bits.Set(&cpu.F, Flag3, bits.Get(in0, 3))
	bits.Set(&cpu.F, FlagV, !bits.Get(in0, n))
	bits.Set(&cpu.F, FlagN, false)
}

func blockc(cpu *CPU, hlfn func(*CPU, cpu.Put16, cpu.Get16)) {
	carry := bits.Get(cpu.F, FlagC)
	in0 := cpu.A
	sub(cpu, cpu.loadIndHL, false)
	out := alu.Out
	cpu.A = in0

	hlfn(cpu, cpu.storeHL, cpu.loadHL)
	dec16(cpu, cpu.storeBC, cpu.loadBC)

	bits.Set(&cpu.F, FlagV, cpu.B != 0 || cpu.C != 0)
	flagResult := out
	if bits.Get(cpu.F, FlagH) {
		flagResult--
	}
	bits.Set(&cpu.F, Flag3, bits.Get(flagResult, 3))
	bits.Set(&cpu.F, Flag5, bits.Get(flagResult, 1)) // yes, one
	bits.Set(&cpu.F, FlagC, carry)
}

func blockcr(cpu *CPU, hlfn func(*CPU, cpu.Put16, cpu.Get16)) {
	blockc(cpu, hlfn)
	for {
		if cpu.B == 0 && cpu.C == 0 {
			break
		}
		if cpu.A == cpu.mem.Load(bits.Join(cpu.H, cpu.L)-1) {
			break
		}
		cpu.refreshR()
		cpu.refreshR()
		blockc(cpu, hlfn)
	}
}

// Performs a "LD (DE),(HL)", then increments DE and HL, and decrements BC.
//
// P/V is reset in case of overflow (if BC=0 after calling LDI).
func blockl(cpu *CPU, increment int) {
	source := bits.Join(cpu.H, cpu.L)
	target := bits.Join(cpu.D, cpu.E)
	v := cpu.mem.Load(source)
	cpu.mem.Store(target, v)

	cpu.H, cpu.L = bits.Split(source + uint16(increment))
	cpu.D, cpu.E = bits.Split(target + uint16(increment))

	counter := bits.Join(cpu.B, cpu.C)
	counter--
	bits.Set(&cpu.F, FlagV, counter != 0)
	cpu.B, cpu.C = bits.Split(counter)

	bits.Set(&cpu.F, Flag5, bits.Get(v+cpu.A, 1)) // yes, bit one
	bits.Set(&cpu.F, Flag3, bits.Get(v+cpu.A, 3))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagH, false)
}

func blocklr(cpu *CPU, increment int) {
	blockl(cpu, increment)
	for cpu.B != 0 || cpu.C != 0 {
		cpu.refreshR()
		cpu.refreshR()
		blockl(cpu, increment)
	}
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
	cpu.refreshR()
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
	in0 := cpu.A
	in1 := get()
	sub(cpu, func() uint8 { return in1 }, false)
	bits.Set(&cpu.F, Flag3, bits.Get(in1, 3))
	bits.Set(&cpu.F, Flag5, bits.Get(in1, 5))
	cpu.A = in0
}

// inverts all bits of A
//
// Sets H and N, other flags are unmodified.
func cpl(cpu *CPU) {
	alu.In0 = cpu.A
	alu.Not()

	bits.Set(&cpu.F, FlagH, true)
	bits.Set(&cpu.F, FlagN, true)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))

	cpu.A = alu.Out
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

func dd(cpu *CPU) {
	opcode := cpu.fetch()
	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	cpu.refreshR()
	opsDD[opcode](cpu)
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
	in0 := get()
	put(in0 - 1)
}

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
	cpu.refreshR()
	opsED[opcode](cpu)
}

// exchange
func ex(cpu *CPU, get0 cpu.Get16, put0 cpu.Put16, get1 cpu.Get16, put1 cpu.Put16) {
	in0 := get0()
	in1 := get1()
	put0(in1)
	put1(in0)
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
	in0 := get()
	put(in0 + 1)
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

func ldair(cpu *CPU, get cpu.Get) {
	in0 := get()

	bits.Set(&cpu.F, FlagS, bits.Get(in0, 7))
	bits.Set(&cpu.F, FlagZ, in0 == 0)
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, FlagV, cpu.IFF2)
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, Flag5, bits.Get(in0, 5))
	bits.Set(&cpu.F, Flag3, bits.Get(in0, 3))

	cpu.A = in0
}

func neg(cpu *CPU) {
	alu.In0 = 0
	alu.In1 = cpu.A
	alu.SetCarry(false)
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

// no operation, no interrupt
func noni(cpu *CPU) {}

// no operation
func nop() {}

// Logical or
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag cleared.
func or(cpu *CPU, get cpu.Get) {
	alu.In0 = cpu.A
	alu.In1 = get()
	alu.Or()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, false)

	cpu.A = alu.Out
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

func rld(cpu *CPU) {
	addr := bits.Join(cpu.H, cpu.L)
	ahi, alo := bits.Split4(cpu.A)
	memhi, memlo := bits.Split4(cpu.mem.Load(addr))

	cpu.A = bits.Join4(ahi, memhi)
	memval := bits.Join4(memlo, alo)
	cpu.mem.Store(addr, memval)

	bits.Set(&cpu.F, FlagS, bits.Get(cpu.A, 7))
	bits.Set(&cpu.F, FlagZ, cpu.A == 0)
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, FlagV, bits.Parity(cpu.A))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, Flag5, bits.Get(cpu.A, 5))
	bits.Set(&cpu.F, Flag3, bits.Get(cpu.A, 3))

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

func rrd(cpu *CPU) {
	addr := bits.Join(cpu.H, cpu.L)
	ahi, alo := bits.Split4(cpu.A)
	memhi, memlo := bits.Split4(cpu.mem.Load(addr))

	cpu.A = bits.Join4(ahi, memlo)
	memval := bits.Join4(alo, memhi)
	cpu.mem.Store(addr, memval)

	bits.Set(&cpu.F, FlagS, bits.Get(cpu.A, 7))
	bits.Set(&cpu.F, FlagZ, cpu.A == 0)
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, FlagV, bits.Parity(cpu.A))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, Flag5, bits.Get(cpu.A, 5))
	bits.Set(&cpu.F, Flag3, bits.Get(cpu.A, 3))

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
	in0 := get()
	bits.Set(&in0, n, true)
	put(in0)
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
func sub(cpu *CPU, get cpu.Get, withCarry bool) {
	alu.In0 = cpu.A
	alu.In1 = get()

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

func sub16(cpu *CPU, put cpu.Put16, get0 cpu.Get16, get1 cpu.Get16, withCarry bool) {
	in0 := get0()
	in1 := get1()

	alu.In0 = uint8(in0)
	alu.In1 = uint8(in1)
	alu.SetCarry(false)
	if withCarry && bits.Get(cpu.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.Subtract()
	lo := alu.Out

	alu.In0 = uint8(in0 >> 8)
	alu.In1 = uint8(in1 >> 8)
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

func todo(cpu *CPU) {
	if cpu.testing {
		cpu.skip = true
		return
	}
	panic("not implemented")
}

func todo2(cpu *CPU, get cpu.Get) {
	get()
	if cpu.testing {
		cpu.skip = true
		return
	}
	panic("not implemented")
}

// Logical exclusive or
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag cleared.
func xor(cpu *CPU, get cpu.Get) {
	alu.In0 = cpu.A
	alu.In1 = get()
	alu.ExclusiveOr()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, false)

	cpu.A = alu.Out
}
