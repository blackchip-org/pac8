package bits

import (
	"fmt"
	"testing"

	. "github.com/blackchip-org/pac8/util/expect"
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

	With(t).Expect(alu.Out).ToBe(uint8(3))
}
