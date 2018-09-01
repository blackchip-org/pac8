package z80

import (
	"testing"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/memory"
)

type fuseTest struct {
	name    string
	af      uint16
	bc      uint16
	de      uint16
	hl      uint16
	af1     uint16
	bc1     uint16
	de1     uint16
	hl1     uint16
	ix      uint16
	iy      uint16
	sp      uint16
	pc      uint16
	i       uint8
	r       uint8
	iff1    int
	iff2    int
	halt    int
	tstates int

	snapshots []memory.Snapshot
}

func TestOps(t *testing.T) {
	for _, test := range fuseTests {
		t.Run(test.name, func(t *testing.T) {
			cpu := load(test)
			cpu.Next()
			expected := load(fuseResults[test.name])
			testRegisters(t, cpu, expected)
			testMemory(t, cpu, fuseResults[test.name])
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

func load(test fuseTest) *CPU {
	mem := memory.NewRAM(0x10000)
	cpu := New(mem)

	cpu.A, cpu.F = bits.Split(test.af)
	cpu.B, cpu.C = bits.Split(test.bc)
	cpu.D, cpu.E = bits.Split(test.de)
	cpu.H, cpu.L = bits.Split(test.hl)
	cpu.IX = test.ix
	cpu.IY = test.iy
	cpu.SP = test.sp
	cpu.PC = test.pc
	cpu.I = test.i
	cpu.R = test.r

	for _, snapshot := range test.snapshots {
		memory.Import(mem, snapshot)
	}

	return cpu
}
