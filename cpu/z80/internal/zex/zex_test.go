// +build fn

package main

//go:generate go run gen.go
//go:generate go fmt zexdoc_com_test.go

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/blackchip-org/pac8/cpu/z80"
	"github.com/blackchip-org/pac8/memory"
	"github.com/blackchip-org/pac8/util/bits"
)

const (
	SingleTestPC = uint16(0x0122)
	AllTestsPC   = uint16(0x0100)
)

const (
	TestAdc16 = uint16(0x013a) + (iota * 2)
	TestAdd16
	TestAdd16x
	TestAdd16y
	TestAlu8i
	TestAlu8r
	TestAlu8rx
	TestAlu8x
	TestBitx
	TestBitz80
	TestCpd1
	TestCpi1
	TestDaa
	TestInca
	TestIncb
	TestIncbc
	TestIncc
	TestIncd
	TestIncde
	TestInce
	TestInch
	TestInchl
	TestIncix
	TestInciy
	TestIncl
	TestIncm
	TestIncsp
	TestIncx
	TestIncxh
	TestIncxl
	TestIncyh
	TestIncyl
	TestLd161
	TestLd162
	TestLd163
	TestLd164
	TestLd165
	TestLd166
	TestLd167
	TestLd168
	TestLd16im
	TestLd16ix
	TestLd8bd
	TestLd8im
	TestLd8imx
	TestLd8ix1
	TestLd8ix2
	TestLd8ix3
	TestLd8ixy
	TestLd8rr
	TestLd8rrx
	TestLda
	TestLdd1
	TestLdd2
	TestLdi1
	TestLdi2
	TestNeg
	TestRld
	TestRot8080
	TestRotxy
	TestRotz80
	TestSrz80
	TestSrzx
	TestSt8ix1
	TestSt8ix2
	TestSt8ix3
	TestStabd
)

// To run a single test, change:
// - start to SingleTestPC
// - testN to the "TestXXX" constant of the test to run
func TestZex(t *testing.T) {
	// start := AllTestsPC
	// testN := uint16(0)

	start := SingleTestPC
	testN := TestCpd1

	// Follow the notes at:
	// https://floooh.github.io/2016/07/12/z80-rust-ms1.html
	mem := memory.NewRAM(0x10000)
	c := z80.New(mem)
	memory.ImportBinary(mem, zexdoc, 0x100)
	c.SetPC(start)
	c.H, c.L = bits.Split(testN)
	c.SP = 0xf000
	var out bytes.Buffer
	for {
		c.Next()
		// Tests complete on a jump to $0000
		if c.PC() == 0 {
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
	fmt.Println()
	if strings.Contains(out.String(), "ERROR") {
		t.Errorf("FAILURE")
	} else {
		fmt.Println("SUCCESS")
	}
}
