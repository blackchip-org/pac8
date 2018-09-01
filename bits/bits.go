// Package bits provides utilities for working with bit values.
package bits

import (
	"strconv"
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

func Format(n uint8, bit int, off string, on string) string {
	if Get(n, bit) {
		return on
	} else {
		return off
	}
}

// Get returns the bit value from uint8 n as a bool.
func Get(n uint8, bit int) bool {
	return n&(1<<uint8(bit)) != 0
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

// Join combines a hi byte and a lo byte to create a uint16
func Join(hi uint8, lo uint8) uint16 {
	return uint16(hi)<<8 + uint16(lo)
}

// Split takes a uint16 and splits it out int a hi byte and a lo byte.
func Split(value uint16) (hi uint8, lo uint8) {
	hi = uint8(value >> 8)
	lo = uint8(value)
	return
}

// Slice extracts a sequence of bits in value from bit lo to bit high,
// inclusive.
func Slice(value uint8, lo int, hi int) uint8 {
	value = value >> uint(lo)
	bits := uint(hi - lo + 1)
	mask := uint8(1)<<bits - 1
	return value & mask
}
