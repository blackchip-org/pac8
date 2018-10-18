package memory

import (
	"testing"

	. "github.com/blackchip-org/pac8/expect"
)

func TestFetch(t *testing.T) {
	mem := NewRAM(0x10)
	mem.Store(0x04, 0xcd)
	mem.Store(0x05, 0xab)
	c := NewCursor(mem)
	c.Pos = 0x04

	WithFormat(t, "%02x").Expect(c.Fetch()).ToBe(0x0cd)
	WithFormat(t, "%02x").Expect(c.Fetch()).ToBe(0x0ab)
	WithFormat(t, "%04x").Expect(c.Pos).ToBe(0x0006)
}

func TestPeek(t *testing.T) {
	mem := NewRAM(0x10)
	mem.Store(0x04, 0xcd)
	c := NewCursor(mem)
	c.Pos = 0x04

	WithFormat(t, "%02x").Expect(c.Peek()).ToBe(0x0cd)
	WithFormat(t, "%04x").Expect(c.Pos).ToBe(0x0004)
}

func TestFetchLE(t *testing.T) {
	mem := NewRAM(0x10)
	mem.Store(0x04, 0xcd)
	mem.Store(0x05, 0xab)
	c := NewCursor(mem)
	c.Pos = 0x04

	WithFormat(t, "%04x").Expect(c.FetchLE()).ToBe(0x0abcd)
}

func TestPutN(t *testing.T) {
	mem := NewRAM(0x10)
	mem.Store(0x04, 0xcd)
	c := NewCursor(mem)
	c.Pos = 0x04
	c.PutN(0xee, 0xff)

	WithFormat(t, "%02x").Expect(mem.Load(0x04)).ToBe(0x0ee)
	WithFormat(t, "%02x").Expect(mem.Load(0x05)).ToBe(0x0ff)
	WithFormat(t, "%04x").Expect(c.Pos).ToBe(0x0006)
}
