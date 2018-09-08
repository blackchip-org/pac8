package bits

import "testing"

func TestAdd(t *testing.T) {
	alu := ALU{}
	alu.In0 = 1
	alu.In1 = 1
	alu.SetCarry(false)
	alu.Add()

	have := alu.Out
	want := uint8(2)
	if have != want {
		t.Errorf("\n have: %v \n want: %v", have, want)
	}
}

func TestAddWithCarry(t *testing.T) {
	alu := ALU{}
	alu.In0 = 1
	alu.In1 = 1
	alu.SetCarry(true)
	alu.Add()

	have := alu.Out
	want := uint8(3)
	if have != want {
		t.Errorf("\n have: %v \n want: %v", have, want)
	}
}
