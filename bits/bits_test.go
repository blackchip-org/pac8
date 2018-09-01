package bits

import (
	"fmt"
	"testing"
)

func ExampleParse() {
	n := Parse("11110000")
	fmt.Printf("%02x", n)
	// Output: f0
}

func TestParse8Invalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	Parse("100x1000")
}

func ExampleGet() {
	n := Parse("101")
	fmt.Println(Get(n, 2))
	fmt.Println(Get(n, 1))
	fmt.Println(Get(n, 0))
	// Output:
	// true
	// false
	// true
}

func ExampleSet() {
	var n uint8
	Set(&n, 2, true)
	Set(&n, 1, false)
	Set(&n, 0, true)
	fmt.Printf("%03b", n)
	// Output: 101
}

func ExampleHi() {
	n := uint16(0xabcd)
	fmt.Printf("%x", Hi(n))
	// Output: ab
}

func ExampleSetHi() {
	n := uint16(0x00cd)
	SetHi(&n, 0xab)
	fmt.Printf("%x", n)
	// Output: abcd
}

func ExampleLo() {
	n := uint16(0xabcd)
	fmt.Printf("%x", Lo(n))
	// Output: cd
}

func ExampleSetLo() {
	n := uint16(0xab00)
	SetLo(&n, 0xcd)
	fmt.Printf("%x", n)
	// Output: abcd
}

func ExampleJoin() {
	fmt.Printf("%x", Join(0xab, 0xcd))
	// Output: abcd
}

func ExampleSplit() {
	hi, lo := Split(0xabcd)
	fmt.Printf("%x %x", hi, lo)
	// Output: ab cd
}

func TestSlice(t *testing.T) {
	b := Parse
	tests := []struct {
		lo   int
		hi   int
		in   uint8
		out  uint8
		name string
	}{
		{6, 7, b("11000000"), b("011"), "high one"},
		{6, 7, b("00111111"), b("000"), "high zero"},
		{3, 5, b("00111000"), b("111"), "middle one"},
		{3, 5, b("11000111"), b("000"), "middle zero"},
		{0, 2, b("00000111"), b("111"), "low one"},
		{0, 2, b("11111000"), b("000"), "low zero"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			have := Slice(test.in, test.lo, test.hi)
			want := test.out
			if have != want {
				t.Errorf("\n have: %08b \n want: %08b", have, want)
			}
		})
	}
}

func ExampleSlice() {
	value := Parse("00111000")
	fmt.Printf("%03b", Slice(value, 3, 5))
	// Output: 111
}

func ExampleDisplace() {
	plus := Displace(0x8000, 0x01)
	minus := Displace(0x8000, 0xff)
	fmt.Printf("%04x %04x", plus, minus)
	// Output: 8001 7fff
}
