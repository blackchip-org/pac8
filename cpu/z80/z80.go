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

	I   uint8
	R   uint8
	IXH uint8
	IXL uint8
	IYH uint8
	IYL uint8
	SP  uint16
	PC  uint16

	IFF1 bool
	IFF2 bool
	IM   uint8

	Halt bool

	IOREQ  bool
	IOAddr uint16
	Ports  memory.Memory

	mem   memory.Memory
	mem16 memory.Memory16
	delta uint8
	// address used to load on the last (IX+d) or (IY+d) instruction
	iaddr uint16
}

func New(m memory.Memory) *CPU {
	c := &CPU{
		mem:   m,
		mem16: memory.NewLittleEndian(m),
		Ports: memory.NewRAM(0x100),
	}
	return c
}

func (cpu *CPU) Next() {
	cpu.IOREQ = false
	opcode := cpu.fetch()
	execute := ops[opcode]
	cpu.refreshR()
	execute(cpu)
}

func (cpu *CPU) String() string {
	return fmt.Sprintf(""+
		" pc   af   bc   de   hl   ix   iy   sp   i  r\n"+
		"%04x %04x %04x %04x %04x %04x %04x %04x %02x %02x %v\n"+
		"im %v %04x %04x %04x %04x      %v %v %v %v %v %v %v %v %v\n",
		// line 1
		cpu.PC,
		bits.Join(cpu.A, cpu.F),
		bits.Join(cpu.B, cpu.C),
		bits.Join(cpu.D, cpu.E),
		bits.Join(cpu.H, cpu.L),
		bits.Join(cpu.IXH, cpu.IXL),
		bits.Join(cpu.IYH, cpu.IYL),
		cpu.SP,
		cpu.I,
		cpu.R,
		bits.FormatB(cpu.IFF1, "", "iff1"),
		// line 2
		cpu.IM,
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

func (cpu *CPU) fetchd() {
	cpu.delta = cpu.fetch()
}

func (cpu *CPU) refreshR() {
	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	bit7 := cpu.R & 0x80
	cpu.R = (cpu.R+1)&0x7f | bit7
}
