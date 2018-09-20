package bits

import (
	"testing"

	. "github.com/blackchip-org/pac8/expect"
)

func TestAdd(t *testing.T) {
	alu := ALU{}
	alu.In0 = 1
	alu.In1 = 1
	alu.SetCarry(false)
	alu.Add()

	With(t).Expect(alu.Out).ToBe(uint8(2))
}

func TestAddWithCarry(t *testing.T) {
	alu := ALU{}
	alu.In0 = 1
	alu.In1 = 1
	alu.SetCarry(true)
	alu.Add()

	With(t).Expect(alu.Out).ToBe(uint8(3))
}
