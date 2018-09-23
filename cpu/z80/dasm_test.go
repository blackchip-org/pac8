package z80

import (
	"fmt"
	"testing"

	. "github.com/blackchip-org/pac8/expect"

	"github.com/blackchip-org/pac8/cpu"
	"github.com/blackchip-org/pac8/memory"
)

func TestDasm(t *testing.T) {
	for _, test := range harstonTests {
		t.Run(test.name, func(t *testing.T) {
			mem := memory.NewRAM(0x20)
			c := memory.NewCursor(mem)
			dasm := cpu.NewDisassembler(mem, ReaderZ80)
			dasm.SetPC(0x10)
			c.Pos = 0x10
			c.PutN(test.bytes...)
			c.Pos = 0x10
			s := dasm.Next()
			With(t).Expect(s.Op).ToBe(test.op)
		})
	}
}

func TestInvalidDD(t *testing.T) {
	for i := 0; i <= 0xff; i++ {
		name := fmt.Sprintf("%02x", i)
		t.Run(name, func(t *testing.T) {
			mem := memory.NewRAM(0x20)
			mem.Store(0, 0xdd)
			mem.Store(1, uint8(i))
			dasm := cpu.NewDisassembler(mem, ReaderZ80)
			s := dasm.Next()
			if s.Op[0] == '?' && i != 0xdd && i != 0xed && i != 0xfd {
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
			dasm := cpu.NewDisassembler(mem, ReaderZ80)
			s := dasm.Next()
			if s.Op[0] == '?' && i != 0xdd && i != 0xed && i != 0xfd {
				With(t).Expect(s.Op).ToBe(fmt.Sprintf("?fd%02x", i))
			}
		})
	}
}
