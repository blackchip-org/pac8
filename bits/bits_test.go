package bits

import (
	"fmt"
	"testing"

	. "github.com/blackchip-org/pac8/expect"
)

func ExampleParse() {
	n := Parse("11110000")
	fmt.Printf("%02x", n)
	// Output: f0
}

func TestParse8Invalid(t *testing.T) {
	With(t).Expect(func() { Parse("100x1000") }).ToPanic()
}

func ExampleFormat() {
	n := Parse("01001010")
	for i := 7; i >= 0; i-- {
		fmt.Printf(Format(n, i, ".", "*"))
	}
	// Output: .*..*.*.
}

func ExampleFormatB() {
	v := []bool{false, true, false, false, true, false, true, false}
	for i := 0; i <= 7; i++ {
		fmt.Printf(FormatB(v[i], ".", "*"))
	}
	// Output: .*..*.*.
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
			slice := Slice(test.in, test.lo, test.hi)
			WithFormat(t, "%08b").Expect(slice).ToBe(test.out)
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

func TestPlane(t *testing.T) {
	b := Parse
	p := []int{0, 4}
	tests := []struct {
		offset int
		in     uint8
		out    uint8
	}{
		{0, b("00010001"), b("11")},
		{1, b("00100010"), b("11")},
		{2, b("01000100"), b("11")},
		{3, b("10001000"), b("11")},
	}
	for _, test := range tests {
		out := Plane(test.in, p, test.offset)
		WithFormat(t, "%08b").Expect(out).ToBe(test.out)
	}
}
