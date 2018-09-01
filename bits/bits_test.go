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
