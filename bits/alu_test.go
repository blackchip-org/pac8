package bits

import "testing"

func TestAdd(t *testing.T) {
	alu := ALU{}
	alu.O1 = 1
	alu.O2 = 1
	alu.SetCarry(false)
	alu.Add()

	have := alu.Result
	want := uint8(2)
	if have != want {
		t.Errorf("\n have: %v \n want: %v", have, want)
	}
}

func TestAddWithCarry(t *testing.T) {
	alu := ALU{}
	alu.O1 = 1
	alu.O2 = 1
	alu.SetCarry(true)
	alu.Add()

	have := alu.Result
	want := uint8(3)
	if have != want {
		t.Errorf("\n have: %v \n want: %v", have, want)
	}
}
