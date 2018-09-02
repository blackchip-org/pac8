package z80

import (
	"fmt"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/memory"
)

const (
	FlagS = 7
	FlagZ = 6
	Flag5 = 5
	FlagH = 4
	Flag3 = 3
	FlagV = 2
	FlagN = 1
	FlagC = 0
)

type CPU struct {
	A uint8
	F uint8
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8

	A1 uint8
	F1 uint8
	B1 uint8
	C1 uint8
	D1 uint8
	E1 uint8
	H1 uint8
	L1 uint8

	I  uint8
	R  uint8
	IX uint16
	IY uint16
	SP uint16
	PC uint16

	IFF1 bool
	IFF2 bool

	Halt bool

	mem   memory.Memory
	mem16 memory.Memory16
	skip  bool
}

func New(m memory.Memory) *CPU {
	c := &CPU{
		mem:   m,
		mem16: memory.NewLittleEndian(m),
	}
	return c
}

func (cpu *CPU) Next() {
	opcode := cpu.fetch()
	execute := ops[opcode]
	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	cpu.R = (cpu.R + 1) & 0x7f
	execute(cpu)
}

func (cpu *CPU) String() string {
	return fmt.Sprintf(""+
		" pc   af   bc   de   hl   ix   iy   sp   i  r\n"+
		"%04x %04x %04x %04x %04x %04x %04x %04x %02x %02x %v\n"+
		"     %04x %04x %04x %04x      %v %v %v %v %v %v %v %v %v\n",
		// line 1
		cpu.PC,
		bits.Join(cpu.A, cpu.F),
		bits.Join(cpu.B, cpu.C),
		bits.Join(cpu.D, cpu.E),
		bits.Join(cpu.H, cpu.L),
		cpu.IX,
		cpu.IY,
		cpu.SP,
		cpu.I,
		cpu.R,
		bits.FormatB(cpu.IFF1, "", "iff1"),
		// line 2
		bits.Join(cpu.A1, cpu.F1),
		bits.Join(cpu.B1, cpu.C1),
		bits.Join(cpu.D1, cpu.E1),
		bits.Join(cpu.H1, cpu.L1),
		// flags
		bits.Format(cpu.F, FlagS, ".", "S"),
		bits.Format(cpu.F, FlagZ, ".", "Z"),
		bits.Format(cpu.F, Flag5, ".", "5"),
		bits.Format(cpu.F, FlagH, ".", "H"),
		bits.Format(cpu.F, Flag3, ".", "3"),
		bits.Format(cpu.F, FlagV, ".", "V"),
		bits.Format(cpu.F, FlagN, ".", "N"),
		bits.Format(cpu.F, FlagC, ".", "C"),
		bits.FormatB(cpu.IFF2, "", "iff2"))
}

func (cpu *CPU) fetch() uint8 {
	cpu.PC++
	return cpu.mem.Load(cpu.PC - 1)
}

func (cpu *CPU) fetch16() uint16 {
	lo := cpu.fetch()
	hi := cpu.fetch()
	return bits.Join(hi, lo)
}
