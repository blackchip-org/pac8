package z80

// http://z80-heaven.wikidot.com/instructions-set
// https://www.worldofspectrum.org/faq/reference/z80reference.htm
// http://www.z80.info/zip/z80-documented.pdf

import (
	"github.com/blackchip-org/pac8/pkg/memory"
	"github.com/blackchip-org/pac8/pkg/proc"
	"github.com/blackchip-org/pac8/pkg/util/bits"
)

var alu bits.ALU

type opsTable map[uint8]func(cpu *CPU)

func add(cpu *CPU, get0 proc.Get, get1 proc.Get, withCarry bool) {
	alu.SetCarry(false)
	if withCarry && bits.Get(cpu.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.A = get0()
	alu.Add(get1())

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Overflow())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	cpu.A = alu.A
}

// add
// preserve s, z, p/v. h undefined
func add16(cpu *CPU, put proc.Put16, get0 proc.Get16, get1 proc.Get16, withCarry bool) {
	in0 := get0()
	in1 := get1()

	alu.SetCarry(false)
	if withCarry && bits.Get(cpu.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.A = bits.Lo(in0)
	alu.Add(bits.Lo(in1))
	zero0 := alu.Zero()
	lo := alu.A

	alu.A = bits.Hi(in0)
	alu.Add(bits.Hi(in1))
	zero1 := alu.Zero()
	hi := alu.A

	result := bits.Join(hi, lo)

	if withCarry {
		bits.Set(&cpu.F, FlagS, alu.Sign())
		bits.Set(&cpu.F, FlagZ, zero0 && zero1)
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
func and(cpu *CPU, get proc.Get) {
	alu.A = cpu.A
	alu.And(get())

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, true)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, false)

	cpu.A = alu.A
}

// Tests if the specified bit is set.
//
// Opposite of the nth bit is written into the Z flag. C is preserved,
// N is reset, H is set, and S and P/V are undefined.
//
// PV as Z, S set only if n=7 and b7 of r set
func bit(cpu *CPU, n int, get proc.Get) {
	in0 := get()

	bits.Set(&cpu.F, FlagS, n == 7 && bits.Get(in0, 7))
	bits.Set(&cpu.F, FlagZ, !bits.Get(in0, n))
	bits.Set(&cpu.F, Flag5, bits.Get(in0, 5))
	bits.Set(&cpu.F, FlagH, true)
	bits.Set(&cpu.F, Flag3, bits.Get(in0, 3))
	bits.Set(&cpu.F, FlagV, !bits.Get(in0, n))
	bits.Set(&cpu.F, FlagN, false)
}

func biti(cpu *CPU, n int, get proc.Get) {
	bit(cpu, n, get)

	// "This is where things start to get strange"
	bits.Set(&cpu.F, Flag5, bits.Get(bits.Hi(cpu.iaddr), 5))
	bits.Set(&cpu.F, Flag3, bits.Get(bits.Hi(cpu.iaddr), 3))
}

func blockc(cpu *CPU, increment int) {
	alu.SetBorrow(false)
	alu.A = cpu.A
	alu.Subtract(cpu.loadIndHL())

	cpu.storeHL(cpu.loadHL() + uint16(increment))
	cpu.storeBC(cpu.loadBC() - uint16(1))

	result := alu.A
	if alu.Carry4() {
		result--
	}

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(result, 1)) // yes, one
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(result, 3))
	bits.Set(&cpu.F, FlagV, cpu.loadBC() != 0)
	// carry unchanged
	bits.Set(&cpu.F, FlagN, true)
}

func blockcr(cpu *CPU, increment int) {
	blockc(cpu, increment)
	for {
		if cpu.B == 0 && cpu.C == 0 {
			break
		}
		if cpu.A == cpu.mem.Load(bits.Join(cpu.H, cpu.L)-1) {
			break
		}
		cpu.refreshR()
		cpu.refreshR()
		blockc(cpu, increment)
	}
}

func blockin(cpu *CPU, increment int) {
	in := cpu.inIndC()
	alu.SetBorrow(false)
	alu.A = cpu.B
	alu.Subtract(1)

	// https://github.com/mamedev/mame/blob/master/src/devices/device/proc/z80/z80.cpp
	// I was unable to figure this out by reading all the conflicting
	// documentation for these "undefined" flags
	t := uint16(cpu.C+uint8(increment)) + uint16(in)
	p := uint8(t&0x07) ^ alu.A // parity check
	halfAndCarry := t&0x100 != 0

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, halfAndCarry)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, bits.Parity(p))
	bits.Set(&cpu.F, FlagN, bits.Get(in, 7))
	bits.Set(&cpu.F, FlagC, halfAndCarry)

	cpu.storeIndHL(in)
	cpu.B = alu.A
	cpu.H, cpu.L = bits.Split(bits.Join(cpu.H, cpu.L) + uint16(increment))
}

func blockinr(cpu *CPU, increment int) {
	blockin(cpu, increment)
	for cpu.B != 0 {
		cpu.refreshR()
		cpu.refreshR()
		blockin(cpu, increment)
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

func blockout(cpu *CPU, increment int) {
	in := cpu.mem.Load(bits.Join(cpu.H, cpu.L))
	alu.SetBorrow(false)
	alu.A = cpu.B
	alu.Subtract(1)

	cpu.B = alu.A
	cpu.H, cpu.L = bits.Split(bits.Join(cpu.H, cpu.L) + uint16(increment))

	// https://github.com/mamedev/mame/blob/master/src/devices/device/proc/z80/z80.cpp
	// I was unable to figure this out by reading all the conflicting
	// documentation for these "undefined" flags
	t := uint16(cpu.L) + uint16(in)
	p := uint8(t&0x07) ^ alu.A // parity check
	halfAndCarry := t&0x100 != 0

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, halfAndCarry)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, bits.Parity(p))
	bits.Set(&cpu.F, FlagN, bits.Get(in, 7))
	bits.Set(&cpu.F, FlagC, halfAndCarry)

	cpu.Ports.Store(uint16(cpu.C), in)
}

func blockoutr(cpu *CPU, increment int) {
	blockout(cpu, increment)
	for cpu.B != 0 {
		cpu.refreshR()
		cpu.refreshR()
		blockout(cpu, increment)
	}
}

// call, conditional
//
// Pushes the address after the CALL instruction (PC+3) onto the stack and
// jumps to the label. Can also take conditions.
func call(cpu *CPU, flag int, condition bool, get proc.Get16) {
	addr := get()
	if bits.Get(cpu.F, flag) == condition {
		cpu.SP -= 2
		memory.StoreLE(cpu.mem, cpu.SP, cpu.PC())
		cpu.SetPC(addr)
	}
}

// call, always
func calla(cpu *CPU, get proc.Get16) {
	addr := get()
	cpu.SP -= 2
	memory.StoreLE(cpu.mem, cpu.SP, cpu.PC())
	cpu.SetPC(addr)
}

func cb(cpu *CPU) {
	opcode := cpu.fetch()
	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	cpu.refreshR()
	opsCB[opcode](cpu)
}

// inverts the carry flag
func ccf(cpu *CPU) {
	// The H flag was tricky. Correct definition in the Z80 User Manual
	bits.Set(&cpu.F, FlagH, bits.Get(cpu.F, FlagC))
	bits.Set(&cpu.F, FlagC, !bits.Get(cpu.F, FlagC))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, Flag3, bits.Get(cpu.A, 3))
	bits.Set(&cpu.F, Flag5, bits.Get(cpu.A, 5))
}

// CP is a subtraction from A that doesn't update A, only the flags it would
// have set/reset if it really was subtracted.
//
// F5 and F3 are copied from the operand, not the result
func cp(cpu *CPU, get proc.Get) {
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
	alu.A = cpu.A
	alu.Not()

	bits.Set(&cpu.F, FlagH, true)
	bits.Set(&cpu.F, FlagN, true)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))

	cpu.A = alu.A
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
//
// Eventually ported directly from the MAME source code.
func daa(cpu *CPU) {
	a := cpu.A
	half := false
	if bits.Get(cpu.F, FlagN) {
		if bits.Get(cpu.F, FlagH) || cpu.A&0xf > 9 {
			a -= 6
		}
		if bits.Get(cpu.F, FlagC) || cpu.A > 0x99 {
			a -= 0x60
		}
		if bits.Get(cpu.F, FlagH) && cpu.A&0xf <= 0x5 {
			half = true
		}
	} else {
		if bits.Get(cpu.F, FlagH) || cpu.A&0xf > 9 {
			a += 6
		}
		if bits.Get(cpu.F, FlagC) || cpu.A > 0x99 {
			a += 0x60
		}
		if cpu.A&0xf > 0x9 {
			half = true
		}
	}

	if cpu.A > 0x99 {
		bits.Set(&cpu.F, FlagC, true)
	}
	bits.Set(&cpu.F, FlagH, half)
	bits.Set(&cpu.F, FlagS, bits.Get(a, 7))
	bits.Set(&cpu.F, FlagZ, a == 0)
	bits.Set(&cpu.F, FlagV, bits.Parity(a))
	bits.Set(&cpu.F, Flag5, bits.Get(a, 5))
	bits.Set(&cpu.F, Flag3, bits.Get(a, 3))

	cpu.A = a
}

func ddfd(cpu *CPU, table opsTable, extendedTable opsTable) {
	// Peek at the next opcode, it it doesn't have a function in the table,
	// return now and let it execute as a normal instruction
	opcode := cpu.mem.Load(cpu.PC())
	fn := table[opcode]
	if fn == nil {
		return
	}

	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	cpu.fetch()
	cpu.refreshR()

	// If the opcode is 0xcb, then this is an extended ddcb or fdcb
	// operation
	if opcode == 0xcb {
		ddfdcb(cpu, extendedTable)
		return
	}

	fn(cpu)
}

func ddfdcb(cpu *CPU, table opsTable) {
	cpu.fetchd()
	opcode := cpu.fetch()
	table[opcode](cpu)
}

// decrement
// C flag preserved, P/V detects overflow and rest modified by definition.
// modified by definition.
func dec(cpu *CPU, put proc.Put, get proc.Get) {
	alu.SetBorrow(false)
	alu.A = get()
	alu.Subtract(1)

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Overflow())
	bits.Set(&cpu.F, FlagN, true)

	put(alu.A)
}

// decrement 16-bit
// No flags altered
func dec16(cpu *CPU, put proc.Put16, get proc.Get16) {
	in0 := get()
	put(in0 - 1)
}

func di(cpu *CPU) {
	cpu.IFF1 = false
	cpu.IFF2 = false
}

// decrement B and jump if not zero
func djnz(cpu *CPU, get proc.Get) {
	delta := get()
	cpu.B--
	if cpu.B != 0 {
		cpu.SetPC(bits.Displace(cpu.PC(), delta))
	}
}

func ed(cpu *CPU) {
	opcode := cpu.fetch()
	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	cpu.refreshR()
	opsED[opcode](cpu)
}

func ei(cpu *CPU) {
	cpu.IFF1 = true
	cpu.IFF2 = true
}

// exchange
func ex(cpu *CPU, get0 proc.Get16, put0 proc.Put16, get1 proc.Get16, put1 proc.Put16) {
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

func in(cpu *CPU, put proc.Put, get proc.Get) {
	val := get()

	bits.Set(&cpu.F, FlagS, bits.Get(val, 7))
	bits.Set(&cpu.F, FlagZ, val == 0)
	bits.Set(&cpu.F, Flag5, bits.Get(val, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(val, 3))
	bits.Set(&cpu.F, FlagV, bits.Parity(val))
	bits.Set(&cpu.F, FlagN, false)

	put(val)
}

// increment
// Preserves C flag, N flag is reset, P/V detects overflow and rest are
// modified by definition.
func inc(cpu *CPU, put proc.Put, get proc.Get) {
	alu.SetCarry(false)
	alu.A = get()
	alu.Add(1)

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Overflow())
	bits.Set(&cpu.F, FlagN, false)

	put(alu.A)
}

// increment 16-bit
// No flags altered
func inc16(cpu *CPU, put proc.Put16, get proc.Get16) {
	in0 := get()
	put(in0 + 1)
}

func invalid() {}

func im(cpu *CPU, mode int) {
	cpu.IM = uint8(mode)
}

// jump absolute, conditional
func jp(cpu *CPU, flag int, condition bool, get proc.Get16) {
	addr := get()
	if bits.Get(cpu.F, flag) == condition {
		cpu.SetPC(addr)
	}
}

// jump absolute, always
func jpa(cpu *CPU, get proc.Get16) {
	cpu.SetPC(get())
}

// jump relative, conditional
func jr(cpu *CPU, flag int, condition bool, get proc.Get) {
	delta := get()
	if bits.Get(cpu.F, flag) == condition {
		cpu.SetPC(bits.Displace(cpu.PC(), delta))
	}
}

// jump relative, always
func jra(cpu *CPU, get proc.Get) {
	delta := get()
	cpu.SetPC(bits.Displace(cpu.PC(), delta))
}

// load
func ld(cpu *CPU, put proc.Put, get proc.Get) {
	put(get())
}

// load, 16-bit
func ld16(cpu *CPU, put proc.Put16, get proc.Get16) {
	put(get())
}

func ldair(cpu *CPU, get proc.Get) {
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
	alu.SetBorrow(false)
	alu.A = 0
	alu.Subtract(cpu.A)

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Overflow())
	bits.Set(&cpu.F, FlagN, true)
	bits.Set(&cpu.F, FlagC, alu.Borrow())

	cpu.A = alu.A
}

// no operation, no interrupt
func noni(cpu *CPU) {}

// no operation
func nop() {}

// Logical or
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag cleared.
func or(cpu *CPU, get proc.Get) {
	alu.A = cpu.A
	alu.Or(get())

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, false)

	cpu.A = alu.A
}

// Copies the two bytes from (SP) into the operand, then increases SP by 2.
func pop(cpu *CPU, put proc.Put16) {
	put(memory.LoadLE(cpu.mem, cpu.SP))
	cpu.SP += 2
}

// Decrements the SP by 2 then copies the operand into (SP)
func push(cpu *CPU, get proc.Get16) {
	cpu.SP -= 2
	memory.StoreLE(cpu.mem, cpu.SP, get())
}

// Resets the specified byte to zero
func res(cpu *CPU, n int, put proc.Put, get proc.Get) {
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
	cpu.SetPC(memory.LoadLE(cpu.mem, cpu.SP))
	cpu.SP += 2
}

func reti(cpu *CPU) {
	cpu.SetPC(memory.LoadLE(cpu.mem, cpu.SP))
	cpu.SP += 2
}

func retn(cpu *CPU) {
	cpu.IFF1 = cpu.IFF2
	cpu.SetPC(memory.LoadLE(cpu.mem, cpu.SP))
	cpu.SP += 2
}

// Rotate A left
func rla(cpu *CPU) {
	alu.SetCarry(bits.Get(cpu.F, FlagC))
	alu.A = cpu.A
	alu.ShiftLeft()

	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	cpu.A = alu.A
}

// Rotate left with carry
func rlc(cpu *CPU, put proc.Put, get proc.Get) {
	alu.SetCarry(bits.Get(cpu.F, FlagC))
	alu.A = get()
	alu.RotateLeft()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(alu.A)
}

// rotate A left with carry
func rlca(cpu *CPU) {
	alu.A = cpu.A
	alu.RotateLeft()

	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	v := cpu.A << 1
	if cpu.A&0x80 > 0 {
		v++
	}
	cpu.A = alu.A
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

func rotr(cpu *CPU, put proc.Put, get proc.Get) {
	// FIXME: Should this be here?
	alu.SetCarry(false)
	alu.A = get()
	alu.RotateRight()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(alu.A)
}

// Rotate A right
func rra(cpu *CPU) {
	alu.SetCarry(bits.Get(cpu.F, FlagC))
	alu.A = cpu.A
	alu.ShiftRight()

	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	cpu.A = alu.A
}

func rrca(cpu *CPU) {
	alu.A = cpu.A
	alu.RotateRight()

	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	cpu.A = alu.A
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
	memory.StoreLE(cpu.mem, cpu.SP, cpu.PC())
	cpu.SetPC(uint16(y) * 8)
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
func set(cpu *CPU, n int, put proc.Put, get proc.Get) {
	in0 := get()
	bits.Set(&in0, n, true)
	put(in0)
}

func shiftl(cpu *CPU, put proc.Put, get proc.Get, withCarry bool) {
	alu.SetCarry(false)
	if withCarry && bits.Get(cpu.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.A = get()
	alu.ShiftLeft()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(alu.A)
}

func shiftr(cpu *CPU, put proc.Put, get proc.Get, withCarry bool) {
	alu.SetCarry(false)
	if withCarry && bits.Get(cpu.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.A = get()
	alu.ShiftRight()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(alu.A)
}

func sll(cpu *CPU, put proc.Put, get proc.Get) {
	alu.SetCarry(true)
	alu.A = get()
	alu.ShiftLeft()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(alu.A)
}

func sra(cpu *CPU, put proc.Put, get proc.Get) {
	alu.A = get()
	alu.SetCarry(false)
	alu.ShiftRightSigned()

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, alu.Carry())

	put(alu.A)
}

// Subtract
//
// N flag is reset, P/V is interpreted as overflow.
// Rest of the flags is modified by definition.
func sub(cpu *CPU, get proc.Get, withBorrow bool) {
	alu.SetBorrow(false)
	if withBorrow && bits.Get(cpu.F, FlagC) {
		alu.SetBorrow(true)
	}
	alu.A = cpu.A
	alu.Subtract(get())

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Overflow())
	bits.Set(&cpu.F, FlagN, true)
	bits.Set(&cpu.F, FlagC, alu.Borrow())

	cpu.A = alu.A
}

func sub16(cpu *CPU, put proc.Put16, get0 proc.Get16, get1 proc.Get16, withBorrow bool) {
	in0 := get0()
	in1 := get1()

	alu.SetBorrow(false)
	if withBorrow && bits.Get(cpu.F, FlagC) {
		alu.SetBorrow(true)
	}
	alu.A = bits.Lo(in0)
	alu.Subtract(bits.Lo(in1))
	zero0 := alu.Zero()
	lo := alu.A

	alu.A = bits.Hi(in0)
	alu.Subtract(bits.Hi(in1))
	zero1 := alu.Zero()
	hi := alu.A

	result := bits.Join(hi, lo)

	if withBorrow {
		bits.Set(&cpu.F, FlagS, alu.Sign())
		bits.Set(&cpu.F, FlagZ, zero0 && zero1)
		bits.Set(&cpu.F, FlagV, alu.Overflow())
	}
	bits.Set(&cpu.F, Flag5, bits.Get(hi, 5))
	bits.Set(&cpu.F, FlagH, alu.Carry4())
	bits.Set(&cpu.F, Flag3, bits.Get(hi, 3))
	bits.Set(&cpu.F, FlagN, true)
	bits.Set(&cpu.F, FlagC, alu.Borrow())

	put(result)
}

// Logical exclusive or
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag cleared.
func xor(cpu *CPU, get proc.Get) {
	alu.A = cpu.A
	alu.ExclusiveOr(get())

	bits.Set(&cpu.F, FlagS, alu.Sign())
	bits.Set(&cpu.F, FlagZ, alu.Zero())
	bits.Set(&cpu.F, Flag5, bits.Get(alu.A, 5))
	bits.Set(&cpu.F, FlagH, false)
	bits.Set(&cpu.F, Flag3, bits.Get(alu.A, 3))
	bits.Set(&cpu.F, FlagV, alu.Parity())
	bits.Set(&cpu.F, FlagN, false)
	bits.Set(&cpu.F, FlagC, false)

	cpu.A = alu.A
}
