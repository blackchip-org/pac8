package memory

import (
	"reflect"
	"testing"
)

func TestLittleTo16(t *testing.T) {
	have := LittleEndian.To16(0x12, 0x34)
	want := uint16(0x3412)
	if have != want {
		t.Errorf("\n have: 0x%x \n want: 0x%x", have, want)
	}
}

func TestBigTo16(t *testing.T) {
	have := BigEndian.To16(0x12, 0x34)
	want := uint16(0x1234)
	if have != want {
		t.Errorf("\n have: 0x%x \n want: 0x%x", have, want)
	}
}

func TestLittleFrom16(t *testing.T) {
	lo, hi := LittleEndian.From16(0x3412)
	have := []uint8{lo, hi}
	want := []uint8{0x12, 0x34}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\n have: 0x%x \n want: 0x%x", have, want)
	}
}

func TestBigFrom16(t *testing.T) {
	lo, hi := BigEndian.From16(0x1234)
	have := []uint8{lo, hi}
	want := []uint8{0x12, 0x34}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\n have: 0x%x \n want: 0x%x", have, want)
	}
}
