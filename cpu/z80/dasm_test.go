package z80

import (
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
