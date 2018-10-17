package memory

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/blackchip-org/pac8/util/charset"
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
	WithFormat(t, "%04x").Expect(LoadLE(mem, 0)).ToBe(uint16(0xabcd))
}

func TestLittleEndianStore(t *testing.T) {
	mem := NewRAM(0x10)
	StoreLE(mem, 0, 0xabcd)
	WithFormat(t, "02x").Expect(mem.Load(0)).ToBe(uint8(0xcd))
	WithFormat(t, "02x").Expect(mem.Load(1)).ToBe(uint8(0xab))
}

func TestDump(t *testing.T) {
	var dumpTests = []struct {
		name     string
		start    int
		data     func() []int
		showFrom int
		showTo   int
		want     string
	}{
		{
			"one line", 0x10,
			func() []int { return []int{} },
			0x10, 0x20, "" +
				"$0010 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................",
		}, {
			"two lines", 0x10,
			func() []int { return []int{} },
			0x10, 0x30, "" +
				"$0010 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................\n" +
				"$0020 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................",
		}, {
			"jagged top", 0x10,
			func() []int { return []int{} },
			0x14, 0x30, "" +
				"$0010             00 00 00 00  00 00 00 00 00 00 00 00     ............\n" +
				"$0020 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................",
		}, {
			"jagged bottom", 0x10,
			func() []int { return []int{} },
			0x10, 0x2b, "" +
				"$0010 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................\n" +
				"$0020 00 00 00 00 00 00 00 00  00 00 00 00             ............",
		},
		{
			"single value", 0x10,
			func() []int { return []int{0, 0x41} },
			0x11, 0x11, "" +
				"$0010    41                                             A",
		},
		{
			"$40-$5f", 0x10,
			func() []int {
				data := make([]int, 0)
				for i := 0x40; i < 0x60; i++ {
					data = append(data, i)
				}
				return data
			},
			0x10, 0x30, "" +
				"$0010 40 41 42 43 44 45 46 47  48 49 4a 4b 4c 4d 4e 4f @ABCDEFGHIJKLMNO\n" +
				"$0020 50 51 52 53 54 55 56 57  58 59 5a 5b 5c 5d 5e 5f PQRSTUVWXYZ[\\]^_",
		},
	}

	m := NewRAM(0x100)
	for _, test := range dumpTests {
		t.Run(test.name, func(t *testing.T) {
			for i, value := range test.data() {
				m.Store(uint16(test.start+i), uint8(value))
			}
			have := Dump(m, uint16(test.showFrom), uint16(test.showTo),
				charset.AsciiDecoder)
			have = strings.TrimSpace(have)
			With(t).Expect(have).ToBe(test.want)
		})
	}
}
