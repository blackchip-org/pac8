package z80

import (
	"fmt"
	"testing"

	"github.com/blackchip-org/pac8/cpu/z80/internal/dasm"
	"github.com/blackchip-org/pac8/memory"
	. "github.com/blackchip-org/pac8/expect"
)

func TestDasm(t *testing.T) {
	for _, test := range dasm.Tests {
		t.Run(test.Name, func(t *testing.T) {
			mem := memory.NewRAM(0x20)
			c := memory.NewCursor(mem)
			dasm := NewDisassembler(mem)
			dasm.SetPC(0x10)
			c.Pos = 0x10
			c.PutN(test.Bytes...)
			c.Pos = 0x10
			s := dasm.NextStatement()
			With(t).Expect(s.Op).ToBe(test.Op)
		})
	}
}

func TestInvalid(t *testing.T) {
	var tests = []struct {
		name   string
		prefix []uint8
	}{
		{"dd", []uint8{0xdd}},
		{"ed", []uint8{0xed}},
		{"fd", []uint8{0xfd}},
		{"ddcb", []uint8{0xdd, 0xcb}},
		{"fdcb", []uint8{0xfd, 0xcb}},
	}

	for _, test := range tests {
		for opcode := 0; opcode < 0x100; opcode++ {
			if test.name == "dd" || test.name == "fd" {
				switch opcode {
				case 0xdd, 0xed, 0xfd, 0xcb:
					continue
				}
			}
			mem := memory.NewRAM(0x20)
			cursor := memory.NewCursor(mem)
			cursor.Put(test.prefix[0])
			if len(test.prefix) > 1 {
				cursor.Put(test.prefix[1])
				cursor.Put(0) // displacement byte
			}
			cursor.Put(uint8(opcode))
			dasm := NewDisassembler(mem)
			s := dasm.NextStatement()
			if s.Op[0] == '?' {
				name := fmt.Sprintf("%v%02x", test.name, opcode)
				t.Run(name, func(t *testing.T) {
					With(t).Expect(s.Op).ToBe(fmt.Sprintf("?%s%02x", test.name, opcode))
				})
			}
		}
	}
}

func TestInvalidDD(t *testing.T) {
	for i := 0; i <= 0xff; i++ {
		name := fmt.Sprintf("%02x", i)
		t.Run(name, func(t *testing.T) {
			mem := memory.NewRAM(0x20)
			mem.Store(0, 0xdd)
			mem.Store(1, uint8(i))
			dasm := NewDisassembler(mem)
			s := dasm.NextStatement()
			if s.Op[0] == '?' && i != 0xdd && i != 0xed && i != 0xfd && i != 0xcb {
				With(t).Expect(s.Op).ToBe(fmt.Sprintf("?dd%02x", i))
			}
		})
	}
}

func TestInvalidFD(t *testing.T) {
	for i := 0; i <= 0xff; i++ {
		name := fmt.Sprintf("%02x", i)
		t.Run(name, func(t *testing.T) {
			mem := memory.NewRAM(0x20)
			mem.Store(0, 0xfd)
			mem.Store(1, uint8(i))
			dasm := NewDisassembler(mem)
			s := dasm.NextStatement()
			if s.Op[0] == '?' && i != 0xdd && i != 0xed && i != 0xfd && i != 0xcb {
				With(t).Expect(s.Op).ToBe(fmt.Sprintf("?fd%02x", i))
			}
		})
	}
}

func TestInvalidFDCB(t *testing.T) {
	for i := 0; i <= 0xff; i++ {
		name := fmt.Sprintf("%02x", i)
		t.Run(name, func(t *testing.T) {
			mem := memory.NewRAM(0x20)
			mem.Store(0, 0xfd)
			mem.Store(1, 0xcb)
			mem.Store(2, 0)
			mem.Store(3, uint8(i))
			dasm := NewDisassembler(mem)
			s := dasm.NextStatement()
			if s.Op[0] == '?' {
				With(t).Expect(s.Op).ToBe(fmt.Sprintf("?fdcb%02x", i))
			}
		})
	}
}
