package memory

import (
	"fmt"
	"testing"

	. "github.com/blackchip-org/pac8/expect"
)

func ExampleNewPageMapped() {
	ram1 := NewRAM(0x1000) // 16 pages in length
	ram2 := NewRAM(0x1000) // 16 pages in length

	mem := NewPageMapped([]Block{
		NewBlock(0xa000, ram1),
		NewBlock(0xb000, ram2),
	})

	// Accessing the memory location in ram1 at 0xcd is the same as
	// accessing the mapped memory location at 0xa0cd
	ram1.Store(0xcd, 0x11)
	ram2.Store(0xcd, 0x22)

	fmt.Printf("%02x %02x", mem.Load(0xa0cd), mem.Load(0xb0cd))
	// Output:
	// 11 22
}

func TestPageMappedStore(t *testing.T) {
	ram1 := NewRAM(0x1000)
	ram2 := NewRAM(0x1000)
	mem := NewPageMapped([]Block{
		NewBlock(0x0000, ram1),
		NewBlock(0x1000, ram2),
	})

	mem.Store(0x0044, 0x44)
	WithFormat(t, "%02x").Expect(ram1.Load(0x0044)).ToBe(uint8(0x44))
	mem.Store(0x1555, 0x55)
	WithFormat(t, "%02x").Expect(ram2.Load(0x0555)).ToBe(uint8(0x55))
	mem.Store(0x4000, 0xff)
	WithFormat(t, "%02x").Expect(mem.Load(0x4000)).ToBe(uint8(0))
}

func TestPageMappedLoad(t *testing.T) {
	ram1 := NewRAM(0x1000)
	ram2 := NewRAM(0x1000)
	mem := NewPageMapped([]Block{
		NewBlock(0x0000, ram1),
		NewBlock(0x1000, ram2),
	})

	ram1.Store(0x0044, 0x44)
	WithFormat(t, "%02x").Expect(mem.Load(0x0044)).ToBe(uint8(0x44))
	ram2.Store(0x0555, 0x55)
	WithFormat(t, "%02x").Expect(mem.Load(0x1555)).ToBe(uint8(0x55))
}

func ExampleNewMasked() {
	// Address lines when indicated by A0-A15, this memory has
	// no line A15 connected
	mem := NewMasked(NewRAM(0x10000), 0x7fff)
	mem.Store(0xc000, 0x42)
	fmt.Printf("%02x", mem.Load(0x4000))
	// Output:
	// 42
}
