package z80

import (
	"testing"

	"github.com/blackchip-org/pac8/pkg/memory"
	. "github.com/blackchip-org/pac8/pkg/util/expect"
)

var format = FormatterZ80()

func TestReader(t *testing.T) {
	tests := []struct {
		bytes []uint8
		str   string
		name  string
	}{
		{
			[]uint8{0x50},
			"$0000:  50           ld   d,b",
			"opcode 1",
		},
		{
			[]uint8{0x2a, 0x82, 0x4c},
			"$0000:  2a 82 4c     ld   hl,($4c82)",
			"address",
		},
		{
			[]uint8{0x20, 0x02},
			"$0000:  20 02        jr   nz,$0004",
			"relative jump",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mem := memory.NewROM(test.bytes)
			dasm := NewDisassembler(mem)
			result := dasm.Next()
			With(t).Expect(result).ToBe(test.str)
		})
	}
}
