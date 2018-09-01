package z80

import (
	"testing"

	"github.com/blackchip-org/pac8/bits"
)

var b = bits.Parse

func TestSetFlags(t *testing.T) {
	tests := []struct {
		name     string
		flag     int
		register uint8
	}{
		{name: "S", flag: FlagS, register: b("10000000")},
		{name: "Z", flag: FlagZ, register: b("01000000")},
		{name: "5", flag: Flag5, register: b("00100000")},
		{name: "H", flag: FlagH, register: b("00010000")},
		{name: "3", flag: Flag3, register: b("00001000")},
		{name: "V", flag: FlagV, register: b("00000100")},
		{name: "N", flag: FlagN, register: b("00000010")},
		{name: "C", flag: FlagC, register: b("00000001")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpu := New(nil)
			bits.Set(&cpu.F, test.flag, true)
			have := cpu.F
			want := test.register
			if have != want {
				t.Errorf("\n have: %02x \n want: %02x", have, want)
			}
		})
	}
}

func TestClearFlags(t *testing.T) {
	tests := []struct {
		name     string
		flag     int
		register uint8
	}{
		{name: "S", flag: FlagS, register: b("01111111")},
		{name: "Z", flag: FlagZ, register: b("10111111")},
		{name: "5", flag: Flag5, register: b("11011111")},
		{name: "H", flag: FlagH, register: b("11101111")},
		{name: "3", flag: Flag3, register: b("11110111")},
		{name: "V", flag: FlagV, register: b("11111011")},
		{name: "N", flag: FlagN, register: b("11111101")},
		{name: "C", flag: FlagC, register: b("11111110")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpu := New(nil)
			cpu.F = 0xff
			bits.Set(&cpu.F, test.flag, false)
			have := cpu.F
			want := test.register
			if have != want {
				t.Errorf("\n have: %02x \n want: %02x", have, want)
			}
		})
	}
}
