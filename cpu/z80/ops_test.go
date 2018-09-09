package z80

import (
	"testing"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/memory"
)

//go:generate go run fuse/gen.go
//go:generate go fmt fuse_test.go

// Set a test name here to test a single test
var testSingle = ""

func TestOps(t *testing.T) {
	for _, test := range fuseTests {
		if testSingle != "" && test.name != testSingle {
			continue
		}
		t.Run(test.name, func(t *testing.T) {
			cpu := load(test)
			cpu.testing = true
			i := 0
			for {
				cpu.Next()
				if cpu.skip {
					t.SkipNow()
				}
				if test.name == "dd00" {
					if cpu.PC == 0x0003 {
						break
					}
				} else {
					if cpu.mem.Load(cpu.PC) == 0 && cpu.PC != 0 {
						break
					}
					if test.tstates == 1 {
						break
					}
				}
				if i > 100 {
					t.Fatalf("exceeded execution limit")
				}
				i++
			}
			expected := load(fuseResults[test.name])
			testRegisters(t, cpu, expected)
			testMemory(t, cpu, fuseResults[test.name])
			testHalt(t, cpu, fuseResults[test.name])
		})
	}

}

func testRegisters(t *testing.T, cpu *CPU, expected *CPU) {
	have := cpu.String()
	want := expected.String()
	if have != want {
		t.Fatalf("\n have: \n%v\n want: \n%v\n", have, want)
	}
}

func testMemory(t *testing.T, cpu *CPU, expected fuseTest) {
	diff, equal := memory.Verify(cpu.mem, expected.snapshots)
	if !equal {
		t.Fatalf("\n memory mismatch: \n%v", diff.String())
	}
}

func testHalt(t *testing.T, cpu *CPU, expected fuseTest) {
	have := cpu.Halt
	want := expected.halt != 0
	if have != want {
		t.Fatalf("\n have: %v \n want: %v", have, want)
	}
}

func load(test fuseTest) *CPU {
	mem := memory.NewRAM(0x10000)
	cpu := New(mem)

	cpu.A, cpu.F = bits.Split(test.af)
	cpu.B, cpu.C = bits.Split(test.bc)
	cpu.D, cpu.E = bits.Split(test.de)
	cpu.H, cpu.L = bits.Split(test.hl)

	cpu.A1, cpu.F1 = bits.Split(test.af1)
	cpu.B1, cpu.C1 = bits.Split(test.bc1)
	cpu.D1, cpu.E1 = bits.Split(test.de1)
	cpu.H1, cpu.L1 = bits.Split(test.hl1)

	cpu.IXH, cpu.IXL = bits.Split(test.ix)
	cpu.IYH, cpu.IYL = bits.Split(test.iy)
	cpu.SP = test.sp
	cpu.PC = test.pc
	cpu.I = test.i
	cpu.R = test.r

	for _, snapshot := range test.snapshots {
		memory.Import(mem, snapshot)
	}

	return cpu
}
