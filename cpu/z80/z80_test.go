package z80

import (
	"testing"

	"github.com/blackchip-org/pac8/util/bits"
	. "github.com/blackchip-org/pac8/util/expect"
)

func TestSetFlags(t *testing.T) {
	b := bits.Parse
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
			WithFormat(t, "%08b").Expect(cpu.F).ToBe(test.register)
		})
	}
}

func TestClearFlags(t *testing.T) {
	b := bits.Parse
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
			WithFormat(t, "%08b").Expect(cpu.F).ToBe(test.register)
		})
	}
}

// Not a real test. Useful for visually looking at the CPU status layout.
func TestString(t *testing.T) {
	cpu := New(nil)
	cpu.A = 0x0a
	cpu.F = 0xff
	cpu.B = 0x0b
	cpu.C = 0x0c
	cpu.D = 0x0d
	cpu.E = 0x0e
	cpu.H = 0xf0
	cpu.L = 0x0f
	cpu.IXH = 0x12
	cpu.IXL = 0x34
	cpu.IYH = 0x56
	cpu.IYL = 0x78
	cpu.SP = 0xabcd
	cpu.I = 0xee
	cpu.R = 0xff

	cpu.A1 = 0xa0
	cpu.F1 = 0x88
	cpu.B1 = 0xb0
	cpu.C1 = 0xc0
	cpu.D1 = 0xd0
	cpu.E1 = 0xe0
	cpu.H1 = 0x0f
	cpu.L1 = 0xf0

	cpu.IFF1 = true
	cpu.IFF2 = true

	//fmt.Println(cpu.String())
	//t.Fail()
}
