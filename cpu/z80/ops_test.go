package z80

import (
	"testing"

	. "github.com/blackchip-org/pac8/expect"

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
			setupPorts(cpu, fuseResults[test.name])
			for {
				cpu.Next()
				if cpu.skip {
					t.SkipNow()
				}
				if test.name == "dd00" {
					if cpu.PC == 0x0003 {
						break
					}
				} else if test.name == "ddfd00" {
					if cpu.PC == 0x0004 {
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

			testMemory(t, cpu, fuseResults[test.name])
			WithFormat(t, "\n%v").Expect(cpu.String()).ToBe(expected.String())
			testHalt(t, cpu, fuseResults[test.name])
			testPorts(t, cpu, fuseResults[test.name])
		})
	}

}

func testMemory(t *testing.T, cpu *CPU, expected fuseTest) {
	diff, equal := memory.Verify(cpu.mem, expected.snapshots)
	if !equal {
		t.Fatalf("\nmemory mismatch (have, want): \n%v", diff.String())
	}
}

func testHalt(t *testing.T, cpu *CPU, expected fuseTest) {
	WithFormat(t, "halt(%v)").Expect(cpu.Halt).ToBe(expected.halt != 0)
}

func setupPorts(cpu *CPU, expected fuseTest) {
	if len(expected.portReads) > 0 {
		cpu.Ports = newMockIO(expected.portReads)
	}
}

func testPorts(t *testing.T, cpu *CPU, expected fuseTest) {
	diff, ok := memory.Verify(cpu.Ports, expected.portWrites)
	if !ok {
		t.Fatalf("\n write ports mismatch: \n%v", diff.String())
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

type mockIO struct {
	data map[uint8][]uint8
}

func newMockIO(snapshots []memory.Snapshot) *mockIO {
	mio := &mockIO{
		data: make(map[uint8][]uint8),
	}
	for _, snapshot := range snapshots {
		addr := uint8(snapshot.Address)
		stack, exists := mio.data[addr]
		if !exists {
			stack = make([]uint8, 0, 0)
		}
		stack = append(stack, snapshot.Values[0])
		mio.data[addr] = stack
	}
	return mio
}

func (m *mockIO) Load(addr uint16) uint8 {
	stack, exists := m.data[uint8(addr)]
	if !exists {
		return 0
	}
	if len(stack) == 0 {
		return 0
	}
	v := stack[0]
	stack = stack[1:]
	m.data[uint8(addr)] = stack
	return v
}

func (m *mockIO) Store(addr uint16, value uint8) {}
