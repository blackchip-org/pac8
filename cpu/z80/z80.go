package z80

import (
	"fmt"

	"github.com/blackchip-org/pac8/cpu"

	"github.com/blackchip-org/pac8/memory"
	"github.com/blackchip-org/pac8/util/bits"
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
	pc  uint16

	IFF1 bool
	IFF2 bool
	IM   uint8

	Halt  bool
	Ports memory.Memory
	Map   memory.PortMapper

	mem   memory.Memory
	delta uint8
	// address used to load on the last (IX+d) or (IY+d) instruction
	iaddr      uint16
	requestInt chan uint8

	intRequested bool
	intData      uint8
}

func New(m memory.Memory) *CPU {
	io := memory.NewIO(0x100)
	c := &CPU{
		mem:        m,
		Ports:      io,
		Map:        memory.NewPortMapper(io),
		requestInt: make(chan uint8, 1),
	}
	return c
}

func (c *CPU) Next() {
	if !c.Halt {
		opcode := c.fetch()
		execute := ops[opcode]
		c.refreshR()
		execute(c)

		// When an EI instruction is executed, any pending interrupt request
		// is not accepted until after the instruction following EI is
		// executed. This single instruction delay is necessary when the
		// next instruction is a return instruction.
		if opcode == 0xfb {
			return
		}
	}

	select {
	case v := <-c.requestInt:
		c.intRequested = true
		c.intData = v
	default:
	}

	if c.IFF1 && c.intRequested {
		c.intRequested = false
		c.intAck(c.intData)
	}

}

func (c *CPU) PC() uint16 {
	return c.pc
}

func (c *CPU) SetPC(pc uint16) {
	c.pc = pc
}

func (c *CPU) INT(v uint8) {
	c.requestInt <- v
}

func (c *CPU) Ready() bool {
	return c.Halt != true
}

func (c *CPU) Memory() memory.Memory {
	return c.mem
}

func (c *CPU) CycleRate() int {
	// CPU is 3.072 MHz which is one T-State every 325 nanoseconds.
	// Roughly round to 1 instruction every 2 microseconds.
	// 500 instructions per millsecond
	return 500
}

func (c *CPU) Disassembler() *cpu.Disassembler {
	return NewDisassembler(c.mem)
}

func (c *CPU) intAck(v uint8) {
	if c.IM != 2 {
		panic(fmt.Sprintf("unsupported interrupt mode %v", c.IM))
	}
	c.Halt = false
	c.IFF1 = false
	c.IFF2 = false
	c.SP -= 2
	memory.StoreLE(c.mem, c.SP, c.PC())
	vector := bits.Join(c.I, v)
	c.pc = memory.LoadLE(c.mem, vector)
}

func (c *CPU) String() string {
	return fmt.Sprintf(""+
		" pc   af   bc   de   hl   ix   iy   sp   i  r\n"+
		"%04x %04x %04x %04x %04x %04x %04x %04x %02x %02x %v\n"+
		"im %v %04x %04x %04x %04x      %v %v %v %v %v %v %v %v %v\n",
		// line 1
		c.pc,
		bits.Join(c.A, c.F),
		bits.Join(c.B, c.C),
		bits.Join(c.D, c.E),
		bits.Join(c.H, c.L),
		bits.Join(c.IXH, c.IXL),
		bits.Join(c.IYH, c.IYL),
		c.SP,
		c.I,
		c.R,
		bits.FormatB(c.IFF1, "", "iff1"),
		// line 2
		c.IM,
		bits.Join(c.A1, c.F1),
		bits.Join(c.B1, c.C1),
		bits.Join(c.D1, c.E1),
		bits.Join(c.H1, c.L1),
		// flags
		bits.Format(c.F, FlagS, ".", "S"),
		bits.Format(c.F, FlagZ, ".", "Z"),
		bits.Format(c.F, Flag5, ".", "5"),
		bits.Format(c.F, FlagH, ".", "H"),
		bits.Format(c.F, Flag3, ".", "3"),
		bits.Format(c.F, FlagV, ".", "V"),
		bits.Format(c.F, FlagN, ".", "N"),
		bits.Format(c.F, FlagC, ".", "C"),
		bits.FormatB(c.IFF2, "", "iff2"))
}

func (c *CPU) fetch() uint8 {
	c.pc++
	return c.mem.Load(c.pc - 1)
}

func (c *CPU) fetch16() uint16 {
	lo := c.fetch()
	hi := c.fetch()
	return bits.Join(hi, lo)
}

func (c *CPU) fetchd() {
	c.delta = c.fetch()
}

func (c *CPU) refreshR() {
	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	bit7 := c.R & 0x80
	c.R = (c.R+1)&0x7f | bit7
}
