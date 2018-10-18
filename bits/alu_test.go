package bits

import (
	"fmt"
	"testing"

	. "github.com/blackchip-org/pac8/expect"
)

func ExampleALU_Add() {
	a := uint16(0x1234)
	b := uint16(0x1111)

	alu := ALU{}
	alu.SetCarry(false)
	alu.A = Lo(a)
	alu.Add(Lo(b))
	lo := alu.A

	alu.A = Hi(a)
	alu.Add(Hi(b))
	hi := alu.A

	fmt.Printf("%04x", Join(hi, lo))
	// Ouput: 2345
}

func TestAddWithCarry(t *testing.T) {
	alu := ALU{}
	alu.A = 1
	alu.SetCarry(true)
	alu.Add(1)

	With(t).Expect(alu.A).ToBe(3)
}

func TestNot(t *testing.T) {
	alu := ALU{}
	alu.A = Parse("10101010")
	alu.Not()
	WithFormat(t, "%08b").Expect(alu.A).ToBe(Parse("01010101"))
}

func TestShiftLeft(t *testing.T) {
	alu := ALU{}
	alu.A = Parse("11110000")
	alu.ShiftLeft()
	WithFormat(t, "%08b").Expect(alu.A).ToBe(Parse("11100000"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
	alu.ShiftLeft()
	WithFormat(t, "%08b").Expect(alu.A).ToBe(Parse("11000001"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
}

func TestShiftRight(t *testing.T) {
	alu := ALU{}
	alu.A = Parse("00001111")
	alu.ShiftRight()
	WithFormat(t, "%08b").Expect(alu.A).ToBe(Parse("00000111"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
	alu.ShiftRight()
	WithFormat(t, "%08b").Expect(alu.A).ToBe(Parse("10000011"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
}

func TestRotateLeft(t *testing.T) {
	alu := ALU{}
	alu.A = Parse("11110000")
	alu.RotateLeft()
	WithFormat(t, "%08b").Expect(alu.A).ToBe(Parse("11100001"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
	alu.RotateLeft()
	WithFormat(t, "%08b").Expect(alu.A).ToBe(Parse("11000011"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
}

func TestRotateRight(t *testing.T) {
	alu := ALU{}
	alu.A = Parse("00001111")
	alu.RotateRight()
	WithFormat(t, "%08b").Expect(alu.A).ToBe(Parse("10000111"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
	alu.RotateRight()
	WithFormat(t, "%08b").Expect(alu.A).ToBe(Parse("11000011"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
}
