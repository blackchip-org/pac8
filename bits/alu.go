package bits

type table [256]uint8
type table2 [256 * 256]uint8

var addFlags table2
var adcFlags table2
var subFlags table2
var sbcFlags table2
var szpFlags table

const (
	flagCarry    = uint8(1 << 0)
	flagOverflow = uint8(1 << 2)
	flagParity   = uint8(1 << 3)
	flagCarry4   = uint8(1 << 4)
	flagZero     = uint8(1 << 6)
	flagSign     = uint8(1 << 7)
)

type ALU struct {
	In0   uint8
	In1   uint8
	Out   uint8
	flags uint8
}

func (a *ALU) Add() {
	a.Out = a.In0 + a.In1
	if a.Carry() {
		a.Out++
		a.flags = adcFlags[a.index2()]
	} else {
		a.flags = addFlags[a.index2()]
	}
}

func (a *ALU) Subtract() {
	a.Out = a.In0 - a.In1
	if a.Carry() {
		a.Out--
		a.flags = sbcFlags[a.index2()]
	} else {
		a.flags = subFlags[a.index2()]
	}
}

func (a *ALU) Increment() {
	a.In1 = 1
	a.Add()
}

func (a *ALU) Decrement() {
	a.In1 = 1
	a.Subtract()
}

func (a *ALU) RotateLeft() {
	carryOut := a.In0&0x80 != 0
	a.Out = a.In0 << 1
	if carryOut {
		a.Out++
	}
	a.flags = szpFlags[a.Out]
	if carryOut {
		a.flags |= flagCarry
	}
}

func (a *ALU) ShiftLeft() {
	carryOut := a.In0&0x80 != 0
	a.Out = a.In0 << 1
	if a.flags&flagCarry != 0 {
		a.Out++
	}
	a.flags = szpFlags[a.Out]
	if carryOut {
		a.flags |= flagCarry
	}
}

func (a *ALU) RotateRight() {
	carryOut := a.In0&0x1 != 0
	a.Out = a.In0 >> 1
	if carryOut {
		a.Out |= (1 << 7)
	}
	a.flags = szpFlags[a.Out]
	if carryOut {
		a.flags |= flagCarry
	}
}

func (a *ALU) ShiftRight() {
	carryOut := a.In0&0x01 != 0
	a.Out = a.In0 >> 1
	if a.flags&flagCarry != 0 {
		a.Out |= (1 << 7)
	}
	a.flags = szpFlags[a.Out]
	if carryOut {
		a.flags |= flagCarry
	}
}

func (a *ALU) ShiftRightSigned() {
	sign := a.In0 & (1 << 7)
	carryOut := a.In0&0x01 != 0
	a.Out = a.In0 >> 1
	a.Out |= sign
	a.flags = szpFlags[a.Out]
	if carryOut {
		a.flags |= flagCarry
	}
}

func (a *ALU) And() {
	a.Out = a.In0 & a.In1
	a.flags = szpFlags[a.Out]
}

func (a *ALU) Not() {
	a.Out = ^a.In0
	a.flags = szpFlags[a.Out]
}

func (a *ALU) Or() {
	a.Out = a.In0 | a.In1
	a.flags = szpFlags[a.Out]
}

func (a *ALU) ExclusiveOr() {
	a.Out = a.In0 ^ a.In1
	a.flags = szpFlags[a.Out]
}

func (a ALU) Carry() bool {
	return a.flags&flagCarry != 0
}

func (a ALU) Overflow() bool {
	return a.flags&flagOverflow != 0
}

func (a ALU) Parity() bool {
	return a.flags&flagParity != 0
}

func (a ALU) Carry4() bool {
	return a.flags&flagCarry4 != 0
}

func (a ALU) Zero() bool {
	return a.flags&flagZero != 0
}

func (a ALU) Sign() bool {
	return a.flags&flagSign != 0
}

func (a *ALU) SetCarry(v bool) {
	if v {
		a.flags |= flagCarry
	} else {
		a.flags &^= flagCarry
	}
}

func (a ALU) index2() int {
	return int(a.In0) | int(a.In1)<<8
}

func init() {
	addTable(&addFlags, 0)
	addTable(&adcFlags, 1)
	subTable(&subFlags, 0)
	subTable(&sbcFlags, 1)
	szpTable()
}

func addTable(table *table2, carry int) {
	for i := 0; i < 256*256; i++ {
		in0 := uint8(i)
		in1 := uint8(i >> 8)

		// result of 8 bit addition into 16 bits
		r := uint16(in0) + uint16(in1) + uint16(carry)
		// signed result, 16-bit
		sr := int16(int8(in0)) + int16(int8(in1)) + int16(carry)
		// unsigned result, 8-bit
		ur := uint8(r)
		// result of half add
		hr := in0&0xf + in1&0xf + uint8(carry)

		var flags uint8
		if r > uint16(0xff) {
			flags |= flagCarry
		}
		if sr < MinInt8 || sr > MaxInt8 {
			flags |= flagOverflow
		}
		if Parity(ur) {
			flags |= flagParity
		}
		if hr > 0xf {
			flags |= flagCarry4
		}
		if ur == 0 {
			flags |= flagZero
		}
		if ur&0x80 != 0 {
			flags |= flagSign
		}
		table[i] = flags
	}
}

func subTable(table *table2, carry int) {
	for i := 0; i < 256*256; i++ {
		in0 := uint8(i)
		in1 := uint8(i >> 8)

		// result of 8 bit subtraction into 16 bits
		r := int16(in0) - int16(in1) - int16(carry)
		// signed result, 16-bit
		sr := int16(int8(in0)) - int16(int8(in1)) - int16(carry)
		// unsigned result, 8-bit
		ur := uint8(r)
		// result of half subtraction
		hr := int8(in0)&0xf - int8(in1)&0xf - int8(carry)

		var flags uint8
		if r < 0 {
			flags |= flagCarry
		}
		if sr < MinInt8 || sr > MaxInt8 {
			flags |= flagOverflow
		}
		if Parity(ur) {
			flags |= flagParity
		}
		if hr < 0 {
			flags |= flagCarry4
		}
		if ur == 0 {
			flags |= flagZero
		}
		if ur&0x80 != 0 {
			flags |= flagSign
		}
		table[i] = flags
	}
}

func szpTable() {
	for i := 0; i < 256; i++ {
		in0 := uint8(i)

		var flags uint8
		if in0&0x80 != 0 {
			flags |= flagSign
		}
		if in0 == 0 {
			flags |= flagZero
		}
		if Parity(in0) {
			flags |= flagParity
		}
		szpFlags[i] = flags
	}
}
