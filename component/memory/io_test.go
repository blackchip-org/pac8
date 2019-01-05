package memory

import (
	"bytes"
	"encoding/gob"
	"testing"

	. "github.com/blackchip-org/pac8/pkg/util/expect"
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

func TestIOSaveRestore(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	var ports1, ports2 [3]uint8
	newMemory := func(ports *[3]uint8) Memory {
		io := NewIO(2)
		pm := NewPortMapper(io)
		pm.RO(0x00, &ports[0])
		pm.WO(0x00, &ports[1])
		pm.RW(0x01, &ports[2])
		return io
	}
	mem1 := newMemory(&ports1)
	ports1[0] = 12
	mem1.Store(0x00, 34)
	mem1.Store(0x01, 56)
	mem1.Save(enc)

	mem2 := newMemory(&ports2)
	mem2.Restore(dec)
	With(t).Expect(ports2).ToBe(ports1)
}
