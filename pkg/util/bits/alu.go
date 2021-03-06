package bits

const (
	flagCarry    = uint8(1 << 0)
	flagBorrow   = uint8(1 << 1)
	flagOverflow = uint8(1 << 2)
	flagParity   = uint8(1 << 3)
	flagCarry4   = uint8(1 << 4)
	flagZero     = uint8(1 << 6)
	flagSign     = uint8(1 << 7)
)

type table [256]uint8
type table2 [256 * 256]uint8

var addFlags table2
var adcFlags table2
var subFlags table2
var sbcFlags table2
var szpFlags table

// ALU is an 8-bit arithmetic-logic unit.
type ALU struct {
	A     uint8 // Accumulator
	flags uint8 // Flags
}

// Add adds the value of v to alu.A. If the carry is set, increments the
// result by one.
func (a *ALU) Add(v uint8) {
	i := index2(a.A, v)
	a.A += v
	if a.Carry() {
		a.A++
		a.flags = adcFlags[i]
	} else {
		a.flags = addFlags[i]
	}
}

// Subtract subracts the value of alu.A from v. If the borrow is set,
// decrements the result by one.
func (a *ALU) Subtract(v uint8) {
	i := index2(a.A, v)
	a.A -= v
	if a.Borrow() {
		a.A--
		a.flags = sbcFlags[i]
	} else {
		a.flags = subFlags[i]
	}
}

// RotateLeft rotates the bits of alu.A to the left by one. Bit 7 that is
// shifted out becomes the value of bit 0 and the carry.
func (a *ALU) RotateLeft() {
	carryOut := a.A&0x80 != 0
	a.A <<= 1
	if carryOut {
		a.A++
	}
	a.flags = szpFlags[a.A]
	if carryOut {
		a.flags |= flagCarry
	}
}

// ShiftLeft shifts the bits of alu.A to the left by one. Bit 0 becomes the
// value of the carry. Bit 7, that is shifted out, becomes the new value
// of carry.
func (a *ALU) ShiftLeft() {
	carryOut := a.A&0x80 != 0
	a.A <<= 1
	if a.flags&flagCarry != 0 {
		a.A++
	}
	a.flags = szpFlags[a.A]
	if carryOut {
		a.flags |= flagCarry
	}
}

// RotateRight rotates the bits of alu.A to the right by one. Bit 0 that is
// shifted out becomes the value of bit 7 and the carry.
func (a *ALU) RotateRight() {
	carryOut := a.A&0x1 != 0
	a.A >>= 1
	if carryOut {
		a.A |= (1 << 7)
	}
	a.flags = szpFlags[a.A]
	if carryOut {
		a.flags |= flagCarry
	}
}

// ShiftRight shifts the bits of alu.A to the right by one. Bit 0, that is
// shifted out, becomes the new value of carry.
func (a *ALU) ShiftRight() {
	carryOut := a.A&0x01 != 0
	a.A >>= 1
	if a.flags&flagCarry != 0 {
		a.A |= (1 << 7)
	}
	a.flags = szpFlags[a.A]
	if carryOut {
		a.flags |= flagCarry
	}
}

// ShiftRightSigned shifts the lower 7 bits of alu.A to the right by one.
// Bit 6 becomes the value of the carry. Bit 0, that is shifted out, becomes
// the new value of carry. Bit 7 remains unchanged.
func (a *ALU) ShiftRightSigned() {
	sign := a.A & (1 << 7)
	carryOut := a.A&0x01 != 0
	a.A >>= 1
	if a.flags&flagCarry != 0 {
		a.A |= (1 << 6)
	}
	a.A |= sign
	a.flags = szpFlags[a.A]
	if carryOut {
		a.flags |= flagCarry
	}
}

// And performs a logical "and" between alu.A and v.
func (a *ALU) And(v uint8) {
	a.A &= v
	a.flags = szpFlags[a.A]
}

// Not performs a logical "not" on alu.A.
func (a *ALU) Not() {
	a.A ^= 0xff
	a.flags = szpFlags[a.A]
}

// Or performs a logical "or" between alu.A and v.
func (a *ALU) Or(v uint8) {
	a.A |= v
	a.flags = szpFlags[a.A]
}

// ExclusiveOr performs a logical "xor" between alu.A and v.
func (a *ALU) ExclusiveOr(v uint8) {
	a.A ^= v
	a.flags = szpFlags[a.A]
}

// Carry returns true if the last operation was an addttion and the
// result needs to carry over, or the last operation was a shift/rotate
// which moved out a bit. Otherwise returns false.
func (a ALU) Carry() bool {
	return a.flags&flagCarry != 0
}

// Borrow returns true if the last operation was a subtraction and the
// result needs to borrow, otherwise returns false.
func (a ALU) Borrow() bool {
	return a.flags&flagBorrow != 0
}

// Overflow returns true if the last operation was an addition or
// subtraction and the result would have overflowed if In0, In1, and Out
// were treated as signed numbers, otherwise returns false.
func (a ALU) Overflow() bool {
	return a.flags&flagOverflow != 0
}

// Parity returns true if the result of the last operation has an even
// number bits, otherwise retursn false.
func (a ALU) Parity() bool {
	return a.flags&flagParity != 0
}

func (a ALU) Carry4() bool {
	return a.flags&flagCarry4 != 0
}

// Zero returns true if the result of the last operation is zero.
func (a ALU) Zero() bool {
	return a.flags&flagZero != 0
}

// Sign returns true if the result of the last operation set bit 7, otherwise
// returns false.
func (a ALU) Sign() bool {
	return a.flags&flagSign != 0
}

// SetCarry sets the value of the carry used in the next operation.
func (a *ALU) SetCarry(v bool) {
	if v {
		a.flags |= flagCarry
	} else {
		a.flags &^= flagCarry
	}
}

// SetBorrow sets the value of the borrow used in the next operation.
func (a *ALU) SetBorrow(v bool) {
	if v {
		a.flags |= flagBorrow
	} else {
		a.flags &^= flagBorrow
	}
}

func index2(a uint8, v uint8) int {
	return int(a) | int(v)<<8
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

func subTable(table *table2, borrow int) {
	for i := 0; i < 256*256; i++ {
		in0 := uint8(i)
		in1 := uint8(i >> 8)

		// result of 8 bit subtraction into 16 bits
		r := int16(in0) - int16(in1) - int16(borrow)
		// signed result, 16-bit
		sr := int16(int8(in0)) - int16(int8(in1)) - int16(borrow)
		// unsigned result, 8-bit
		ur := uint8(r)
		// result of half subtraction
		hr := int8(in0)&0xf - int8(in1)&0xf - int8(borrow)

		var flags uint8
		if r < 0 {
			flags |= flagBorrow
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
