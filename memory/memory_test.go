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
