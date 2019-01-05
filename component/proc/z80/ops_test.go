package z80

import (
	"testing"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/proc/z80/internal/fuse"
	. "github.com/blackchip-org/pac8/pkg/util/expect"
	"github.com/blackchip-org/pac8/pkg/util/state"
)

// Set a test name here to test a single test
var testSingle = ""

// TODO: Write single tests for:
// ADC/SBC: Check that both bytes are zero for zero flag when doing 16-bits

func TestOps(t *testing.T) {
	for _, test := range fuse.Tests {
		if testSingle != "" && test.Name != testSingle {
			continue
		}
		t.Run(test.Name, func(t *testing.T) {
			cpu := load(test)
			i := 0
			setupPorts(cpu, fuse.Expected[test.Name])
			for {
				cpu.Next()
				if test.Name == "dd00" {
					if cpu.PC() == 0x0003 {
						break
					}
				} else if test.Name == "ddfd00" {
					if cpu.PC() == 0x0004 {
						break
					}
				} else {
					if cpu.mem.Load(cpu.PC()) == 0 && cpu.PC() != 0 {
						break
					}
					if test.TStates == 1 {
						break
					}
				}
				if i > 100 {
					t.Fatalf("exceeded execution limit")
				}
				i++
			}
			expected := load(fuse.Expected[test.Name])

			testMemory(t, cpu, fuse.Expected[test.Name])
			WithFormat(t, "\n%v").Expect(cpu.String()).ToBe(expected.String())
			testHalt(t, cpu, fuse.Expected[test.Name])
			testPorts(t, cpu, fuse.Expected[test.Name])
		})
	}

}

func testMemory(t *testing.T, cpu *CPU, expected fuse.Test) {
	diff, equal := memory.Verify(cpu.mem, expected.Snapshots)
	if !equal {
		t.Fatalf("\nmemory mismatch (have, want): \n%v", diff.String())
	}
}

func testHalt(t *testing.T, cpu *CPU, expected fuse.Test) {
	WithFormat(t, "halt(%v)").Expect(cpu.Halt).ToBe(expected.Halt != 0)
}

func setupPorts(cpu *CPU, expected fuse.Test) {
	cpu.Ports = newMockIO(expected.PortReads)
}

func testPorts(t *testing.T, cpu *CPU, expected fuse.Test) {
	diff, ok := memory.Verify(cpu.Ports, expected.PortWrites)
	if !ok {
		t.Fatalf("\n write ports mismatch: \n%v", diff.String())
	}
}

func load(test fuse.Test) *CPU {
	mem := memory.NewRAM(0x10000)
	cpu := New(mem)

	cpu.A, cpu.F = bits.Split(test.AF)
	cpu.B, cpu.C = bits.Split(test.BC)
	cpu.D, cpu.E = bits.Split(test.DE)
	cpu.H, cpu.L = bits.Split(test.HL)

	cpu.A1, cpu.F1 = bits.Split(test.AF1)
	cpu.B1, cpu.C1 = bits.Split(test.BC1)
	cpu.D1, cpu.E1 = bits.Split(test.DE1)
	cpu.H1, cpu.L1 = bits.Split(test.HL1)

	cpu.IXH, cpu.IXL = bits.Split(test.IX)
	cpu.IYH, cpu.IYL = bits.Split(test.IY)
	cpu.SP = test.SP
	cpu.SetPC(test.PC)
	cpu.I = test.I
	cpu.R = test.R
	cpu.IFF1 = test.IFF1 != 0
	cpu.IFF2 = test.IFF2 != 0
	cpu.IM = uint8(test.IM)

	for _, snapshot := range test.Snapshots {
		memory.Import(mem, snapshot)
	}

	return cpu
}

type mockIO struct {
	data map[uint8][]uint8
}

func newMockIO(snapshots []memory.Snapshot) memory.IO {
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

func (m *mockIO) Store(addr uint16, value uint8) {
	stack, exists := m.data[uint8(addr)]
	if !exists {
		stack = make([]uint8, 0, 0)
	}
	stack = append(stack, value)
	m.data[uint8(addr)] = stack
}

func (m *mockIO) Length() int {
	return 0
}

func (m *mockIO) Port(n int) *memory.Port {
	return &memory.Port{}
}

func (m *mockIO) Save(_ *state.Encoder) {}

func (m *mockIO) Restore(_ *state.Decoder) {}
