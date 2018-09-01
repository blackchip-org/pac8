package memory

type ByteOrder interface {
	To16(uint8, uint8) uint16
	From16(uint16) (uint8, uint8)
}

var LittleEndian littleEndian

type littleEndian struct{}

func (littleEndian) To16(lo uint8, hi uint8) uint16 {
	return uint16(hi)<<8 + uint16(lo)
}

func (littleEndian) From16(value uint16) (lo uint8, hi uint8) {
	lo = uint8(value)
	hi = uint8(value >> 8)
	return
}

var BigEndian bigEndian

type bigEndian struct{}

func (bigEndian) To16(lo uint8, hi uint8) uint16 {
	return uint16(hi) + uint16(lo)<<8
}

func (bigEndian) From16(value uint16) (lo uint8, hi uint8) {
	lo = uint8(value >> 8)
	hi = uint8(value)
	return
}
