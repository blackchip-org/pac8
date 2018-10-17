// +build fn

/*
Running the full zexdoc test takes about 7 minutes. This test instead
breaks up each test into an individual run. The HL register is loaded
with the address of the test and the program counter is set to the
beginning of the normal test loop. Execution is stopped when the
program counter returns to the top of the loop. Output is then checked
for "ERROR" to determine if the test passes or fails.
*/

package zex

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/blackchip-org/pac8/cpu/z80"
	"github.com/blackchip-org/pac8/memory"
	"github.com/blackchip-org/pac8/pac8"
	"github.com/blackchip-org/pac8/util/bits"
)

func TestZexdoc(t *testing.T) {
	var tests = []string{
		"adc16",
		"add16",
		"add16x",
		"add16y",
		"alu8i",
		"alu8r",
		"alu8rx",
		"alu8x",
		"bitx",
		"bitz80",
		"cpd1",
		"cpi1",
		"daa",
		"inca",
		"incb",
		"incbc",
		"incc",
		"incd",
		"incde",
		"ince",
		"inch",
		"inchl",
		"incix",
		"inciy",
		"incl",
		"incm",
		"incsp",
		"incx",
		"incxh",
		"incxl",
		"incyh",
		"incyl",
		"ld161",
		"ld162",
		"ld163",
		"ld164",
		"ld165",
		"ld166",
		"ld167",
		"ld168",
		"ld16im",
		"ld16ix",
		"ld8bd",
		"ld8im",
		"ld8imx",
		"ld8ix1",
		"ld8ix2",
		"ld8ix3",
		"ld8ixy",
		"ld8rr",
		"ld8rrx",
		"lda",
		"ldd1",
		"ldd2",
		"ldi1",
		"ldi2",
		"neg",
		"rld",
		"rot8080",
		"rotxy",
		"rotz80",
		"srz80",
		"srzx",
		"st8ix1",
		"st8ix2",
		"st8ix3",
		"stabd",
	}

	zexdocFile := filepath.Join(pac8.Home(), "data", "zex", "zexdoc.com")
	zexdoc, err := ioutil.ReadFile(zexdocFile)
	if err != nil {
		t.Fatalf("unable to read %v: %v", zexdocFile, err)
	}

	testBaseAddr := uint16(0x013a)
	for i, test := range tests {
		addr := testBaseAddr + (uint16(i) * 2)
		t.Run(test, func(t *testing.T) {
			ok := runner(zexdoc, addr)
			if !ok {
				t.Fail()
			}
		})
	}
}

func runner(code []byte, addr uint16) bool {
	var out bytes.Buffer

	// Follow the notes at:
	// https://floooh.github.io/2016/07/12/z80-rust-ms1.html
	mem := memory.NewRAM(0x10000)
	memory.ImportBinary(mem, code, 0x100)

	c := z80.New(mem)
	c.SP = 0xf000
	// HL register is loaded with the address of the test to run
	c.H, c.L = bits.Split(addr)

	loopStart := uint16(0x0122)
	c.SetPC(loopStart)
	for {
		c.Next()
		// Tests complete when back at the start of the loop
		if c.PC() == loopStart {
			break
		}
		// System call that outputs to the screen
		if c.PC() == 0x0005 {
			// Single character out
			if c.C == 0x02 {
				msg := fmt.Sprintf("%c", rune(c.C))
				out.WriteString(msg)
				fmt.Print(msg)
			}
			// String out, terminated by $
			if c.C == 0x09 {
				addr := bits.Join(c.D, c.E)
				for {
					ch := rune(mem.Load(addr))
					if ch == '$' {
						break
					}
					msg := fmt.Sprintf("%c", ch)
					out.WriteString(msg)
					fmt.Print(msg)
					addr++
				}
			}
			// Return from subroutine
			c.SetPC(memory.LoadLE(mem, c.SP))
			c.SP += 2
		}
	}
	return !strings.Contains(out.String(), "ERROR")
}
