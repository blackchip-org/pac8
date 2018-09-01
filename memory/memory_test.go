package memory

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
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
