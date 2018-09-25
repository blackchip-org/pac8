package z80

// http://z80-heaven.wikidot.com/instructions-set
// https://www.worldofspectrum.org/faq/reference/z80reference.htm
// http://www.z80.info/zip/z80-documented.pdf

import (
	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/cpu"
)

//go:generate go run ops/gen.go
//go:generate go fmt ops.go

var alu bits.ALU

type opsTable map[uint8]func(c *CPU)

// Add
//
// N flag is reset, P/V is interpreted as overflow.
// Rest of the flags is modified by definition.
func add(c *CPU, get0 cpu.Get, get1 cpu.Get, withCarry bool) {
	alu.In0 = get0()
	alu.In1 = get1()

	alu.SetCarry(false)
	if withCarry && bits.Get(c.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.Add()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, alu.Carry4())
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Overflow())
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, FlagC, alu.Carry())

	c.A = alu.Out
}

// add
// preserve s, z, p/v. h undefined
func add16(c *CPU, put cpu.Put16, get0 cpu.Get16, get1 cpu.Get16, withCarry bool) {
	in0 := get0()
	in1 := get1()

	alu.In0 = uint8(in0)
	alu.In1 = uint8(in1)
	alu.SetCarry(false)
	if withCarry && bits.Get(c.F, FlagC) {
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
		bits.Set(&c.F, FlagS, alu.Sign())
		bits.Set(&c.F, FlagZ, alu.Zero())
		bits.Set(&c.F, FlagV, alu.Overflow())
	}
	bits.Set(&c.F, Flag5, bits.Get(hi, 5))
	bits.Set(&c.F, FlagH, alu.Carry4())
	bits.Set(&c.F, Flag3, bits.Get(hi, 3))
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, FlagC, alu.Carry())

	put(result)
}

// Logical and
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag set.
func and(c *CPU, get cpu.Get) {
	alu.In0 = c.A
	alu.In1 = get()
	alu.And()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, true)
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Parity())
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, FlagC, false)

	c.A = alu.Out
}

// Tests if the specified bit is set.
//
// Opposite of the nth bit is written into the Z flag. C is preserved,
// N is reset, H is set, and S and P/V are undefined.
//
// PV as Z, S set only if n=7 and b7 of r set
func bit(c *CPU, n int, get cpu.Get) {
	in0 := get()

	bits.Set(&c.F, FlagS, n == 7 && bits.Get(in0, 7))
	bits.Set(&c.F, FlagZ, !bits.Get(in0, n))
	bits.Set(&c.F, Flag5, bits.Get(in0, 5))
	bits.Set(&c.F, FlagH, true)
	bits.Set(&c.F, Flag3, bits.Get(in0, 3))
	bits.Set(&c.F, FlagV, !bits.Get(in0, n))
	bits.Set(&c.F, FlagN, false)
}

func biti(c *CPU, n int, get cpu.Get) {
	bit(c, n, get)

	// "This is where things start to get strange"
	bits.Set(&c.F, Flag5, bits.Get(bits.Hi(c.iaddr), 5))
	bits.Set(&c.F, Flag3, bits.Get(bits.Hi(c.iaddr), 3))
}

func blockc(c *CPU, hlfn func(*CPU, cpu.Put16, cpu.Get16)) {
	carry := bits.Get(c.F, FlagC)
	in0 := c.A
	sub(c, c.loadIndHL, false)
	out := alu.Out
	c.A = in0

	hlfn(c, c.storeHL, c.loadHL)
	dec16(c, c.storeBC, c.loadBC)

	bits.Set(&c.F, FlagV, c.B != 0 || c.C != 0)
	flagResult := out
	if bits.Get(c.F, FlagH) {
		flagResult--
	}
	bits.Set(&c.F, Flag3, bits.Get(flagResult, 3))
	bits.Set(&c.F, Flag5, bits.Get(flagResult, 1)) // yes, one
	bits.Set(&c.F, FlagC, carry)
}

func blockcr(c *CPU, hlfn func(*CPU, cpu.Put16, cpu.Get16)) {
	blockc(c, hlfn)
	for {
		if c.B == 0 && c.C == 0 {
			break
		}
		if c.A == c.mem.Load(bits.Join(c.H, c.L)-1) {
			break
		}
		c.refreshR()
		c.refreshR()
		blockc(c, hlfn)
	}
}

func blockin(c *CPU, increment int) {
	in := c.inIndC()
	alu.SetBorrow(false)
	alu.In0 = c.B
	alu.In1 = 1
	alu.Subtract()

	// https://github.com/mamedev/mame/blob/master/src/devices/cpu/z80/z80.cpp
	// I was unable to figure this out by reading all the conflicting
	// documentation for these "undefined" flags
	t := uint16(c.C+uint8(increment)) + uint16(in)
	p := uint8(t&0x07) ^ alu.Out // parity check
	halfAndCarry := t&0x100 != 0

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, halfAndCarry)
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, bits.Parity(p))
	bits.Set(&c.F, FlagN, bits.Get(in, 7))
	bits.Set(&c.F, FlagC, halfAndCarry)

	c.storeIndHL(in)
	c.B = alu.Out
	c.H, c.L = bits.Split(bits.Join(c.H, c.L) + uint16(increment))
}

func blockinr(c *CPU, increment int) {
	blockin(c, increment)
	for c.B != 0 {
		c.refreshR()
		c.refreshR()
		blockin(c, increment)
	}
}

// Performs a "LD (DE),(HL)", then increments DE and HL, and decrements BC.
//
// P/V is reset in case of overflow (if BC=0 after calling LDI).
func blockl(c *CPU, increment int) {
	source := bits.Join(c.H, c.L)
	target := bits.Join(c.D, c.E)
	v := c.mem.Load(source)
	c.mem.Store(target, v)

	c.H, c.L = bits.Split(source + uint16(increment))
	c.D, c.E = bits.Split(target + uint16(increment))

	counter := bits.Join(c.B, c.C)
	counter--
	bits.Set(&c.F, FlagV, counter != 0)
	c.B, c.C = bits.Split(counter)

	bits.Set(&c.F, Flag5, bits.Get(v+c.A, 1)) // yes, bit one
	bits.Set(&c.F, Flag3, bits.Get(v+c.A, 3))
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, FlagH, false)
}

func blocklr(c *CPU, increment int) {
	blockl(c, increment)
	for c.B != 0 || c.C != 0 {
		c.refreshR()
		c.refreshR()
		blockl(c, increment)
	}
}

func blockout(c *CPU, increment int) {
	in := c.mem.Load(bits.Join(c.H, c.L))
	alu.SetBorrow(false)
	alu.In0 = c.B
	alu.In1 = 1
	alu.Subtract()

	c.B = alu.Out
	c.H, c.L = bits.Split(bits.Join(c.H, c.L) + uint16(increment))

	// https://github.com/mamedev/mame/blob/master/src/devices/cpu/z80/z80.cpp
	// I was unable to figure this out by reading all the conflicting
	// documentation for these "undefined" flags
	t := uint16(c.L) + uint16(in)
	p := uint8(t&0x07) ^ alu.Out // parity check
	halfAndCarry := t&0x100 != 0

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, halfAndCarry)
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, bits.Parity(p))
	bits.Set(&c.F, FlagN, bits.Get(in, 7))
	bits.Set(&c.F, FlagC, halfAndCarry)

	c.Ports.Store(uint16(c.C), in)
}

func blockoutr(c *CPU, increment int) {
	blockout(c, increment)
	for c.B != 0 {
		c.refreshR()
		c.refreshR()
		blockout(c, increment)
	}
}

// call, conditional
//
// Pushes the address after the CALL instruction (PC+3) onto the stack and
// jumps to the label. Can also take conditions.
func call(c *CPU, flag int, condition bool, get cpu.Get16) {
	addr := get()
	if bits.Get(c.F, flag) == condition {
		c.SP -= 2
		c.mem16.Store(c.SP, c.PC)
		c.PC = addr
	}
}

// call, always
func calla(c *CPU, get cpu.Get16) {
	addr := get()
	c.SP -= 2
	c.mem16.Store(c.SP, c.PC)
	c.PC = addr
}

func cb(c *CPU) {
	opcode := c.fetch()
	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	c.refreshR()
	opsCB[opcode](c)
}

// Inverts the carry flag
//
// Carry flag inverted. Also inverts H and clears N. Rest of the flags are
// preserved.
func ccf(c *CPU) {
	bits.Set(&c.F, FlagC, !bits.Get(c.F, FlagC))
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, Flag3, bits.Get(c.A, 3))
	bits.Set(&c.F, Flag5, bits.Get(c.A, 5))
}

// CP is a subtraction from A that doesn't update A, only the flags it would
// have set/reset if it really was subtracted.
//
// F5 and F3 are copied from the operand, not the result
func cp(c *CPU, get cpu.Get) {
	in0 := c.A
	in1 := get()
	sub(c, func() uint8 { return in1 }, false)
	bits.Set(&c.F, Flag3, bits.Get(in1, 3))
	bits.Set(&c.F, Flag5, bits.Get(in1, 5))
	c.A = in0
}

// inverts all bits of A
//
// Sets H and N, other flags are unmodified.
func cpl(c *CPU) {
	alu.In0 = c.A
	alu.Not()

	bits.Set(&c.F, FlagH, true)
	bits.Set(&c.F, FlagN, true)
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))

	c.A = alu.Out
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
func daa(c *CPU) {
	result := c.A

	half := false
	carry := false
	if bits.Get(c.F, FlagN) {
		if bits.Get(c.F, FlagH) || c.A&0xf > 9 {
			result -= 0x06
			/* TODO: remove if works
			if result < 6 {
				half = true
			}
			*/
		}
		if bits.Get(c.F, FlagC) || c.A > 0x99 {
			result -= 0x60
			carry = true
		}
	} else {
		if bits.Get(c.F, FlagH) || c.A&0xf > 9 {
			result += 0x06
			half = true
		}
		/* TODO: remove if works
		if bits.Get(c.F, FlagC) || c.A > 0x99 {
			result += 0x60
			carry = true
		}
		*/
	}

	bits.Set(&c.F, FlagS, bits.Get(result, 7))
	bits.Set(&c.F, FlagZ, result == 0)
	bits.Set(&c.F, Flag5, bits.Get(result, 5))
	bits.Set(&c.F, FlagH, half)
	bits.Set(&c.F, Flag3, bits.Get(result, 3))
	bits.Set(&c.F, FlagV, bits.Parity(result))
	bits.Set(&c.F, FlagC, carry)

	c.A = result
}

func ddfd(c *CPU, table opsTable, extendedTable opsTable) {
	// Peek at the next opcode, it it doesn't have a function in the table,
	// return now and let it execute as a normal instruction
	opcode := c.mem.Load(c.PC)
	fn := table[opcode]
	if fn == nil {
		return
	}

	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	c.fetch()
	c.refreshR()

	// If the opcode is 0xcb, then this is an extended ddcb or fdcb
	// operation
	if opcode == 0xcb {
		ddfdcb(c, extendedTable)
		return
	}

	fn(c)
}

func ddfdcb(c *CPU, table opsTable) {
	c.fetchd()
	opcode := c.fetch()
	table[opcode](c)
}

// decrement
// C flag preserved, P/V detects overflow and rest modified by definition.
// modified by definition.
func dec(c *CPU, put cpu.Put, get cpu.Get) {
	alu.SetBorrow(false)
	alu.In0 = get()
	alu.In1 = 1
	alu.Subtract()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, alu.Carry4())
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Overflow())
	bits.Set(&c.F, FlagN, true)

	put(alu.Out)
}

// decrement 16-bit
// No flags altered
func dec16(c *CPU, put cpu.Put16, get cpu.Get16) {
	in0 := get()
	put(in0 - 1)
}

func di(c *CPU) {
	c.IFF1 = false
	c.IFF2 = false
}

// decrement B and jump if not zero
func djnz(c *CPU, get cpu.Get) {
	delta := get()
	c.B--
	if c.B != 0 {
		c.PC = bits.Displace(c.PC, delta)
	}
}

func ed(c *CPU) {
	opcode := c.fetch()
	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	c.refreshR()
	opsED[opcode](c)
}

func ei(c *CPU) {
	c.IFF1 = true
	c.IFF2 = true
}

// exchange
func ex(c *CPU, get0 cpu.Get16, put0 cpu.Put16, get1 cpu.Get16, put1 cpu.Put16) {
	in0 := get0()
	in1 := get1()
	put0(in1)
	put1(in0)
}

// EXX exchanges BC, DE, and HL with shadow registers with BC', DE', and HL'.
func exx(c *CPU) {
	ex(c, c.loadBC, c.storeBC, c.loadBC1, c.storeBC1)
	ex(c, c.loadDE, c.storeDE, c.loadDE1, c.storeDE1)
	ex(c, c.loadHL, c.storeHL, c.loadHL1, c.storeHL1)
}

func halt(c *CPU) {
	c.Halt = true
}

func in(c *CPU, put cpu.Put, get cpu.Get) {
	val := get()

	bits.Set(&c.F, FlagS, bits.Get(val, 7))
	bits.Set(&c.F, FlagZ, val == 0)
	bits.Set(&c.F, Flag5, bits.Get(val, 5))
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, Flag3, bits.Get(val, 3))
	bits.Set(&c.F, FlagV, bits.Parity(val))
	bits.Set(&c.F, FlagN, false)

	put(val)
}

// increment
// Preserves C flag, N flag is reset, P/V detects overflow and rest are
// modified by definition.
func inc(c *CPU, put cpu.Put, get cpu.Get) {
	alu.SetCarry(false)
	alu.In0 = get()
	alu.In1 = 1
	alu.Add()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, alu.Carry4())
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Overflow())
	bits.Set(&c.F, FlagN, false)

	put(alu.Out)
}

// increment 16-bit
// No flags altered
func inc16(c *CPU, put cpu.Put16, get cpu.Get16) {
	in0 := get()
	put(in0 + 1)

}

func invalid() {}

func im(c *CPU, mode int) {
	c.IM = uint8(mode)
}

// jump absolute, conditional
func jp(c *CPU, flag int, condition bool, get cpu.Get16) {
	addr := get()
	if bits.Get(c.F, flag) == condition {
		c.PC = addr
	}
}

// jump absolute, always
func jpa(c *CPU, get cpu.Get16) {
	c.PC = get()
}

// jump relative, conditional
func jr(c *CPU, flag int, condition bool, get cpu.Get) {
	delta := get()
	if bits.Get(c.F, flag) == condition {
		c.PC = bits.Displace(c.PC, delta)
	}
}

// jump relative, always
func jra(c *CPU, get cpu.Get) {
	delta := get()
	c.PC = bits.Displace(c.PC, delta)
}

// load
func ld(c *CPU, put cpu.Put, get cpu.Get) {
	put(get())
}

// load, 16-bit
func ld16(c *CPU, put cpu.Put16, get cpu.Get16) {
	put(get())
}

func ldair(c *CPU, get cpu.Get) {
	in0 := get()

	bits.Set(&c.F, FlagS, bits.Get(in0, 7))
	bits.Set(&c.F, FlagZ, in0 == 0)
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, FlagV, c.IFF2)
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, Flag5, bits.Get(in0, 5))
	bits.Set(&c.F, Flag3, bits.Get(in0, 3))

	c.A = in0
}

func neg(c *CPU) {
	alu.In0 = 0
	alu.In1 = c.A
	alu.SetBorrow(false)
	alu.Subtract()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, alu.Carry4())
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Overflow())
	bits.Set(&c.F, FlagN, true)
	bits.Set(&c.F, FlagC, alu.Borrow())

	c.A = alu.Out
}

// no operation, no interrupt
func noni(c *CPU) {}

// no operation
func nop() {}

// Logical or
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag cleared.
func or(c *CPU, get cpu.Get) {
	alu.In0 = c.A
	alu.In1 = get()
	alu.Or()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Parity())
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, FlagC, false)

	c.A = alu.Out
}

// Copies the two bytes from (SP) into the operand, then increases SP by 2.
func pop(c *CPU, put cpu.Put16) {
	put(c.mem16.Load(c.SP))
	c.SP += 2
}

// Decrements the SP by 2 then copies the operand into (SP)
func push(c *CPU, get cpu.Get16) {
	c.SP -= 2
	c.mem16.Store(c.SP, get())
}

// Resets the specified byte to zero
func res(c *CPU, n int, put cpu.Put, get cpu.Get) {
	arg := get()
	bits.Set(&arg, n, false)
	put(arg)
}

// return, conditional
func ret(c *CPU, flag int, value bool) {
	if bits.Get(c.F, flag) == value {
		reta(c)
	}
}

// return, always
func reta(c *CPU) {
	c.PC = c.mem16.Load(c.SP)
	c.SP += 2
}

func reti(c *CPU) {
	c.PC = c.mem16.Load(c.SP)
	c.SP += 2
}

func retn(c *CPU) {
	c.IFF1 = c.IFF2
	c.PC = c.mem16.Load(c.SP)
	c.SP += 2
}

func rld(c *CPU) {
	addr := bits.Join(c.H, c.L)
	ahi, alo := bits.Split4(c.A)
	memhi, memlo := bits.Split4(c.mem.Load(addr))

	c.A = bits.Join4(ahi, memhi)
	memval := bits.Join4(memlo, alo)
	c.mem.Store(addr, memval)

	bits.Set(&c.F, FlagS, bits.Get(c.A, 7))
	bits.Set(&c.F, FlagZ, c.A == 0)
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, FlagV, bits.Parity(c.A))
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, Flag5, bits.Get(c.A, 5))
	bits.Set(&c.F, Flag3, bits.Get(c.A, 3))

}

func rotl(c *CPU, put cpu.Put, get cpu.Get) {
	alu.In0 = get()
	alu.SetCarry(bits.Get(c.F, FlagC))
	alu.RotateLeft()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Parity())
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, FlagC, alu.Carry())

	put(alu.Out)
}

// rotate A left with carry
// S,Z, and P/V are preserved, H and N flags are reset
func rotla(c *CPU) {
	alu.In0 = c.A
	alu.RotateLeft()

	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagC, alu.Carry())

	c.A = alu.Out
}

func rotr(c *CPU, put cpu.Put, get cpu.Get) {
	alu.In0 = get()
	alu.RotateRight()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Parity())
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, FlagC, alu.Carry())

	put(alu.Out)
}

func rotra(c *CPU) {
	alu.In0 = c.A
	alu.RotateRight()

	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagC, alu.Carry())

	c.A = alu.Out
}

func rrd(c *CPU) {
	addr := bits.Join(c.H, c.L)
	ahi, alo := bits.Split4(c.A)
	memhi, memlo := bits.Split4(c.mem.Load(addr))

	c.A = bits.Join4(ahi, memlo)
	memval := bits.Join4(alo, memhi)
	c.mem.Store(addr, memval)

	bits.Set(&c.F, FlagS, bits.Get(c.A, 7))
	bits.Set(&c.F, FlagZ, c.A == 0)
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, FlagV, bits.Parity(c.A))
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, Flag5, bits.Get(c.A, 5))
	bits.Set(&c.F, Flag3, bits.Get(c.A, 3))

}

func rst(c *CPU, y int) {
	c.SP -= 2
	c.mem16.Store(c.SP, c.PC)
	c.PC = uint16(y) * 8
}

// Set carry flag
//
// Carry flag set, H and N cleared, rest are preserved.
func scf(c *CPU) {
	bits.Set(&c.F, FlagC, true)
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, Flag5, bits.Get(c.A, 5))
	bits.Set(&c.F, Flag3, bits.Get(c.A, 3))
}

// Sets the specified byte to one
func set(c *CPU, n int, put cpu.Put, get cpu.Get) {
	in0 := get()
	bits.Set(&in0, n, true)
	put(in0)
}

func shiftl(c *CPU, put cpu.Put, get cpu.Get, withCarry bool) {
	alu.In0 = get()
	alu.SetCarry(false)
	if withCarry && bits.Get(c.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.ShiftLeft()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Parity())
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, FlagC, alu.Carry())

	put(alu.Out)
}

func shiftla(c *CPU) {
	alu.In0 = c.A
	alu.SetCarry(bits.Get(c.F, FlagC))
	alu.ShiftLeft()

	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagC, alu.Carry())

	c.A = alu.Out
}

func shiftr(c *CPU, put cpu.Put, get cpu.Get, withCarry bool) {
	alu.In0 = get()
	alu.SetCarry(false)
	if withCarry && bits.Get(c.F, FlagC) {
		alu.SetCarry(true)
	}
	alu.ShiftRight()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Parity())
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, FlagC, alu.Carry())

	put(alu.Out)
}

func shiftra(c *CPU) {
	alu.In0 = c.A
	alu.SetCarry(bits.Get(c.F, FlagC))
	alu.ShiftRight()

	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagC, alu.Carry())

	c.A = alu.Out
}

func sll(c *CPU, put cpu.Put, get cpu.Get) {
	alu.In0 = get()
	alu.SetCarry(true)
	alu.ShiftLeft()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Parity())
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, FlagC, alu.Carry())

	put(alu.Out)
}

func sra(c *CPU, put cpu.Put, get cpu.Get) {
	alu.In0 = get()
	alu.SetCarry(false)
	alu.ShiftRightSigned()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Parity())
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, FlagC, alu.Carry())

	put(alu.Out)
}

// Subtract
//
// N flag is reset, P/V is interpreted as overflow.
// Rest of the flags is modified by definition.
func sub(c *CPU, get cpu.Get, withBorrow bool) {
	alu.In0 = c.A
	alu.In1 = get()

	alu.SetBorrow(false)
	if withBorrow && bits.Get(c.F, FlagC) {
		alu.SetBorrow(true)
	}
	alu.Subtract()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, alu.Carry4())
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Overflow())
	bits.Set(&c.F, FlagN, true)
	bits.Set(&c.F, FlagC, alu.Borrow())

	c.A = alu.Out
}

func sub16(c *CPU, put cpu.Put16, get0 cpu.Get16, get1 cpu.Get16, withBorrow bool) {
	in0 := get0()
	in1 := get1()

	alu.In0 = uint8(in0)
	alu.In1 = uint8(in1)
	alu.SetBorrow(false)
	if withBorrow && bits.Get(c.F, FlagC) {
		alu.SetBorrow(true)
	}
	alu.Subtract()
	lo := alu.Out

	alu.In0 = uint8(in0 >> 8)
	alu.In1 = uint8(in1 >> 8)
	alu.Subtract()
	hi := alu.Out

	result := bits.Join(hi, lo)

	if withBorrow {
		bits.Set(&c.F, FlagS, alu.Sign())
		bits.Set(&c.F, FlagZ, alu.Zero())
		bits.Set(&c.F, FlagV, alu.Overflow())
	}
	bits.Set(&c.F, Flag5, bits.Get(hi, 5))
	bits.Set(&c.F, FlagH, alu.Carry4())
	bits.Set(&c.F, Flag3, bits.Get(hi, 3))
	bits.Set(&c.F, FlagN, true)
	bits.Set(&c.F, FlagC, alu.Borrow())

	put(result)
}

// Logical exclusive or
//
// C and N flags cleared, P/V is parity, rest are altered by definition.
// H flag cleared.
func xor(c *CPU, get cpu.Get) {
	alu.In0 = c.A
	alu.In1 = get()
	alu.ExclusiveOr()

	bits.Set(&c.F, FlagS, alu.Sign())
	bits.Set(&c.F, FlagZ, alu.Zero())
	bits.Set(&c.F, Flag5, bits.Get(alu.Out, 5))
	bits.Set(&c.F, FlagH, false)
	bits.Set(&c.F, Flag3, bits.Get(alu.Out, 3))
	bits.Set(&c.F, FlagV, alu.Parity())
	bits.Set(&c.F, FlagN, false)
	bits.Set(&c.F, FlagC, false)

	c.A = alu.Out
}
