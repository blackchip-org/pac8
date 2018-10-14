package memory

import (
	"testing"

	. "github.com/blackchip-org/pac8/util/expect"
)

func TestIONotMapped(t *testing.T) {
	io := NewIO(0xff)
	io.Store(0x0012, 0x42)
	WithFormat(t, "%02x").Expect(io.Load(0x12)).ToBe(0x00)
}

func TestIOMappedRW(t *testing.T) {
	v := uint8(0x42)
	io := NewIO(0xff)
	pm := NewPortMapper(io)
	pm.RW(0x12, &v)
	WithFormat(t, "%02x").Expect(io.Load(0x12)).ToBe(0x42)
	io.Store(0x0012, 0xff)
	WithFormat(t, "%02x").Expect(v).ToBe(0xff)
}

func TestIOMappedSplit(t *testing.T) {
	read := uint8(0x42)
	write := uint8(0)
	io := NewIO(0xff)
	pm := NewPortMapper(io)
	pm.RO(0x12, &read)
	pm.WO(0x12, &write)
	WithFormat(t, "%02x").Expect(io.Load(0x12)).ToBe(0x42)
	io.Store(0x0012, 0xff)
	WithFormat(t, "%02x").Expect(write).ToBe(0xff)
}

func TestIOMappedMulti(t *testing.T) {
	v := uint8(0x42)
	io := NewIO(0xff)
	pm := NewPortMapper(io)
	pm.RW(0x12, &v)
	pm.RW(0x13, &v)
	WithFormat(t, "%02x").Expect(io.Load(0x12)).ToBe(0x42)
	io.Store(0x0013, 0xff)
	WithFormat(t, "%02x").Expect(io.Load(0x12)).ToBe(0xff)
}
