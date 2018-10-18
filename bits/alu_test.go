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
	alu.In0 = Lo(a)
	alu.In1 = Lo(b)
	lo := alu.Out

	alu.In0 = Hi(a)
	alu.In1 = Hi(b)
	hi := alu.Out

	fmt.Printf("%04x", Join(hi, lo))
	// Ouput: 2345
}

func TestAddWithCarry(t *testing.T) {
	alu := ALU{}
	alu.In0 = 1
	alu.In1 = 1
	alu.SetCarry(true)
	alu.Add()

	With(t).Expect(alu.Out).ToBe(3)
}

func TestNot(t *testing.T) {
	alu := ALU{}
	alu.In0 = Parse("10101010")
	alu.Not()
	WithFormat(t, "%08b").Expect(alu.Out).ToBe(Parse("01010101"))
}

func TestShiftLeft(t *testing.T) {
	alu := ALU{}
	alu.In0 = Parse("11110000")
	alu.ShiftLeft()
	WithFormat(t, "%08b").Expect(alu.Out).ToBe(Parse("11100000"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
	alu.In0 = alu.Out
	alu.ShiftLeft()
	WithFormat(t, "%08b").Expect(alu.Out).ToBe(Parse("11000001"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
}

func TestShiftRight(t *testing.T) {
	alu := ALU{}
	alu.In0 = Parse("00001111")
	alu.ShiftRight()
	WithFormat(t, "%08b").Expect(alu.Out).ToBe(Parse("00000111"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
	alu.In0 = alu.Out
	alu.ShiftRight()
	WithFormat(t, "%08b").Expect(alu.Out).ToBe(Parse("10000011"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
}

func TestRotateLeft(t *testing.T) {
	alu := ALU{}
	alu.In0 = Parse("11110000")
	alu.RotateLeft()
	WithFormat(t, "%08b").Expect(alu.Out).ToBe(Parse("11100001"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
	alu.In0 = alu.Out
	alu.RotateLeft()
	WithFormat(t, "%08b").Expect(alu.Out).ToBe(Parse("11000011"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
}

func TestRotateRight(t *testing.T) {
	alu := ALU{}
	alu.In0 = Parse("00001111")
	alu.RotateRight()
	WithFormat(t, "%08b").Expect(alu.Out).ToBe(Parse("10000111"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
	alu.In0 = alu.Out
	alu.RotateRight()
	WithFormat(t, "%08b").Expect(alu.Out).ToBe(Parse("11000011"))
	WithFormat(t, "carry=%v").Expect(alu.Carry()).ToBe(true)
}
