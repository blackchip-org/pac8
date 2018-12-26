package pacman

import (
	"fmt"
	"testing"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/component/memory"
	. "github.com/blackchip-org/pac8/expect"
)

func TestPixelReader(t *testing.T) {
	b := bits.Parse
	tests := []struct {
		base     uint16
		pixel    int
		expected uint8
		rom      []uint8
	}{
		{0, 0, b("11"), []uint8{b("00010001")}},
		{0, 2, b("11"), []uint8{b("01000100")}},
		{0, 4, b("11"), []uint8{0, b("00010001")}},
		{1, 0, b("11"), []uint8{0, b("00010001")}},
	}
	for _, test := range tests {
		name := fmt.Sprintf("base %v pixel %v", test.base, test.pixel)
		t.Run(name, func(t *testing.T) {
			mem := memory.NewROM(test.rom)
			out := pixelReader(mem, test.base, test.pixel)
			WithFormat(t, "%08b").Expect(out).ToBe(test.expected)
		})
	}
}
