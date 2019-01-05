package memory

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	. "github.com/blackchip-org/pac8/pkg/util/expect"
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

func TestSaveRestoreRAM(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	ram1 := NewRAM(0x1000)
	ram1.Store(0x123, 0xab)
	ram1.Store(0x456, 0xbc)
	ram1.Save(enc)

	ram2 := NewRAM(0x1000)
	ram2.Restore(dec)

	report, err := Compare(ram1, ram2)
	if err != nil {
		t.Errorf("%v\n%v", err, report)
	}
}

func TestSaveRestoreMapped(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	newMemory := func() Memory {
		ram1 := NewRAM(0x1000)
		ram2 := NewRAM(0x1000)
		mem := NewPageMapped([]Block{
			NewBlock(0xa000, ram1),
			NewBlock(0xb000, ram2),
		})
		return mem
	}

	mem1 := newMemory()
	mem1.Store(0xa0cd, 0x11)
	mem1.Store(0xb0cd, 0x22)
	mem1.Save(enc)

	mem2 := newMemory()
	mem2.Restore(dec)

	report, err := Compare(mem1, mem2)
	if err != nil {
		t.Errorf("%v\n%v", err, report)
	}
}
