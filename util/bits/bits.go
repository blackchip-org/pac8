// Package bits provides utilities for working with bit values.
package bits

import (
	"math/bits"
	"strconv"
)

const (
	// MaxInt8 is the maximum 8-bit signed value
	MaxInt8 = 127

	// MinInt8 is the minimum 8-bit signed value
	MinInt8 = -128
)

// Parse parses the base-2 string value s to a uint8. Panics if s is not
// a valid number. Use strconv.ParseUint for input which may be malformed.
func Parse(s string) uint8 {
	value, err := strconv.ParseUint(s, 2, 8)
	if err != nil {
		panic(err)
	}
	return uint8(value)
}

// Format returns on if the bit in n is one, otherwise returns off.
func Format(n uint8, bit int, off string, on string) string {
	if Get(n, bit) {
		return on
	}
	return off
}

// FormatB returns on if n is true, otherwise returns off.
func FormatB(n bool, off string, on string) string {
	if n {
		return on
	}
	return off
}

// Get returns the bit value from uint8 n as a bool.
func Get(n uint8, bit int) bool {
	return n&(1<<uint8(bit)) != 0
}

// Get16 returns the bit value from uint16 n as a bool.
func Get16(n uint16, bit int) bool {
	return n&(1<<uint16(bit)) != 0
}

// Set changes the bit in uint8 n to value.
func Set(n *uint8, bit int, value bool) {
	if value {
		*n = *n | (1 << uint8(bit))
	} else {
		*n = *n & (1<<uint8(bit) ^ 0xff)
	}
}

// Hi gets the higher byte of n.
func Hi(n uint16) uint8 {
	return uint8(n >> 8)
}

// SetHi sets the higher byte of n to value.
func SetHi(n *uint16, value uint8) {
	*n = (*n & 0x00ff) + uint16(value)<<8
}

// Lo gets the lower byte of n.
func Lo(n uint16) uint8 {
	return uint8(n)
}

// SetLo sets the lower byte of n to value.
func SetLo(n *uint16, value uint8) {
	*n = (*n & 0xff00) + uint16(value)
}

// Join combines a hi byte and a lo byte to create a uint16.
func Join(hi uint8, lo uint8) uint16 {
	return uint16(hi)<<8 + uint16(lo)
}

// Join4 combines a hi nibble and a lo nibble to create a uint8.
func Join4(hi uint8, lo uint8) uint8 {
	return hi<<4 + lo
}

// Split takes a uint16 and splits it out into a hi byte and a lo byte.
func Split(value uint16) (hi uint8, lo uint8) {
	hi = uint8(value >> 8)
	lo = uint8(value)
	return
}

// Split4 takes a uint8 and splits it out into a hi nibble and lo nibble.
func Split4(value uint8) (hi uint8, lo uint8) {
	hi = uint8(value >> 4)
	lo = uint8(value & 0xf)
	return
}

// Slice extracts a sequence of bits in value from bit lo to bit hi,
// inclusive.
func Slice(value uint8, lo int, hi int) uint8 {
	value = value >> uint(lo)
	bits := uint(hi - lo + 1)
	mask := uint8(1)<<bits - 1
	return value & mask
}

// Displace adds delta to value as a signed number.
func Displace(value uint16, delta uint8) uint16 {
	sdelta := int8(delta)
	v := int(value) + int(sdelta)
	return uint16(v)
}

// Parity returns true if an even number of bits are set in value, otherwise
// returns false.
func Parity(value uint8) bool {
	return bits.OnesCount8(value)%2 == 0
}
