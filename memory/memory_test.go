package memory

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	. "github.com/blackchip-org/pac8/util/expect"
)

func TestImport(t *testing.T) {
	want := []uint8{0xab, 0xcd}
	mem := NewRAM(0x100)
	snapshot := Snapshot{Address: 0x12, Values: want}
	Import(mem, snapshot)
	have := []uint8{
		mem.Load(0x12),
		mem.Load(0x13),
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\n have: %v \n want: %v", hex.EncodeToString(have),
			hex.EncodeToString(want))
	}
}

func ExampleCompare() {
	a := NewROM([]byte{0x12, 0x34, 0x45, 0x67})
	b := NewROM([]byte{0x12, 0xff, 0x45, 0xff})
	diff, equal := Compare(a, b)
	if !equal {
		fmt.Println(diff.String())
	}
	// Output:
	// 0001: 34 ff
	// 0003: 67 ff
}

func ExampleVerify() {
	a := NewROM([]byte{0x12, 0x34, 0x45, 0x67})
	b := []Snapshot{
		Snapshot{Address: 0, Values: []byte{0x12, 0xff}},
		Snapshot{Address: 2, Values: []byte{0x45, 0xff}},
	}
	diff, equal := Verify(a, b)
	if !equal {
		fmt.Println(diff.String())
	}
	// Output:
	// 0001: 34 ff
	// 0003: 67 ff
}

func TestLittleEndianLoad(t *testing.T) {
	mem := NewROM([]uint8{0xcd, 0xab})
	mem16 := NewLittleEndian(mem)
	have := mem16.Load(0)
	want := uint16(0xabcd)
	if have != want {
		t.Errorf("\n have: %02x \n want: %02x", have, want)
	}
}

func TestLittleEndianStore(t *testing.T) {
	mem := NewRAM(0x10)
	mem16 := NewLittleEndian(mem)
	mem16.Store(0, 0xabcd)
	have := []uint8{mem.Load(0), mem.Load(1)}
	want := []uint8{0xcd, 0xab}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\n have: %02x \n want: %02x", have, want)
	}
}

func TestPageMappedStore(t *testing.T) {
	ram1 := NewRAM(0x1000)
	ram2 := NewRAM(0x1000)
	mem := NewPageMapped([]Memory{ram1, ram2})

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
	mem := NewPageMapped([]Memory{ram1, ram2})

	ram1.Store(0x0044, 0x44)
	WithFormat(t, "%02x").Expect(mem.Load(0x0044)).ToBe(uint8(0x44))
	ram2.Store(0x0555, 0x55)
	WithFormat(t, "%02x").Expect(mem.Load(0x1555)).ToBe(uint8(0x55))
}
