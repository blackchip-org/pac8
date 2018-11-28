package z80

import (
	"fmt"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/check"
	"github.com/blackchip-org/pac8/component"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/proc"
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
	Halt bool

	Ports memory.IO
	info  proc.Info
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
		requestInt: make(chan uint8, 1),
	}
	c.info = proc.Info{
		// CPU is 3.072 MHz which is one T-State every 325 nanoseconds.
		// Roughly round to 1 instruction every 1 microsecond.
		// 1000 instructions per millsecond
		CycleRate:       1000,
		CodeReader:      ReaderZ80,
		CodeFormatter:   FormatterZ80(),
		NewDisassembler: NewDisassembler,
		Registers:       c.registers(),
	}
	return c
}

func (cpu *CPU) Next() {
	if !cpu.Halt {
		opcode := cpu.fetch()
		execute := ops[opcode]
		cpu.refreshR()
		execute(cpu)

		// When an EI instruction is executed, any pending interrupt request
		// is not accepted until after the instruction following EI is
		// executed. This single instruction delay is necessary when the
		// next instruction is a return instruction.
		if opcode == 0xfb {
			return
		}
	}

	select {
	case v := <-cpu.requestInt:
		cpu.intRequested = true
		cpu.intData = v
	default:
	}

	if cpu.IFF1 && cpu.intRequested {
		cpu.intRequested = false
		cpu.intAck(cpu.intData)
	}

}

func (cpu *CPU) PC() uint16 {
	return cpu.pc
}

func (cpu *CPU) SetPC(pc uint16) {
	cpu.pc = pc
}

func (cpu *CPU) INT(v uint8) {
	cpu.requestInt <- v
}

func (cpu *CPU) Ready() bool {
	return cpu.Halt != true
}

func (cpu *CPU) Info() proc.Info {
	return cpu.info
}

func (cpu *CPU) intAck(v uint8) {
	if cpu.IM == 0 {
		panic(fmt.Sprintf("unsupported interrupt mode %v", cpu.IM))
	}
	cpu.Halt = false
	cpu.IFF1 = false
	cpu.IFF2 = false
	cpu.SP -= 2
	memory.StoreLE(cpu.mem, cpu.SP, cpu.PC())
	if cpu.IM == 2 {
		vector := bits.Join(cpu.I, v)
		cpu.pc = memory.LoadLE(cpu.mem, vector)
	} else {
		cpu.pc = 0x0038
	}
}

func (cpu *CPU) String() string {
	return fmt.Sprintf(""+
		" pc   af   bc   de   hl   ix   iy   sp   i  r\n"+
		"%04x %04x %04x %04x %04x %04x %04x %04x %02x %02x %v\n"+
		"im %v %04x %04x %04x %04x      %v %v %v %v %v %v %v %v %v\n",
		// line 1
		cpu.pc,
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
	cpu.pc++
	return cpu.mem.Load(cpu.pc - 1)
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

func (cpu *CPU) registers() map[string]proc.Value {
	return map[string]proc.Value{
		"A": proc.Value{Get: cpu.loadA, Put: cpu.storeA},
		"F": proc.Value{Get: cpu.loadF, Put: cpu.storeF},
		"B": proc.Value{Get: cpu.loadB, Put: cpu.storeB},
		"C": proc.Value{Get: cpu.loadC, Put: cpu.storeC},
		"D": proc.Value{Get: cpu.loadD, Put: cpu.storeD},
		"E": proc.Value{Get: cpu.loadE, Put: cpu.storeE},
		"H": proc.Value{Get: cpu.loadH, Put: cpu.storeH},
		"L": proc.Value{Get: cpu.loadL, Put: cpu.storeL},
		"I": proc.Value{Get: cpu.loadI, Put: cpu.storeI},
		"R": proc.Value{Get: cpu.loadR, Put: cpu.storeR},

		"A1": proc.Value{Get: cpu.loadA1, Put: cpu.storeA1},
		"F1": proc.Value{Get: cpu.loadF1, Put: cpu.storeF1},
		"B1": proc.Value{Get: cpu.loadB1, Put: cpu.storeB1},
		"C1": proc.Value{Get: cpu.loadC1, Put: cpu.storeC1},
		"D1": proc.Value{Get: cpu.loadD1, Put: cpu.storeD1},
		"E1": proc.Value{Get: cpu.loadE1, Put: cpu.storeE1},
		"H1": proc.Value{Get: cpu.loadH1, Put: cpu.storeH1},

		"AF": proc.Value{Get: cpu.loadAF, Put: cpu.storeAF},
		"BC": proc.Value{Get: cpu.loadBC, Put: cpu.storeBC},
		"DE": proc.Value{Get: cpu.loadDE, Put: cpu.storeDE},
		"HL": proc.Value{Get: cpu.loadHL, Put: cpu.storeHL},
		"SP": proc.Value{Get: cpu.loadSP, Put: cpu.storeSP},
		"IX": proc.Value{Get: cpu.loadIX, Put: cpu.storeIX},
		"IY": proc.Value{Get: cpu.loadIY, Put: cpu.storeIY},

		"AF1": proc.Value{Get: cpu.loadAF1, Put: cpu.storeAF1},
		"BC1": proc.Value{Get: cpu.loadBC1, Put: cpu.storeBC1},
		"DE1": proc.Value{Get: cpu.loadDE1, Put: cpu.storeDE1},
		"HL1": proc.Value{Get: cpu.loadHL1, Put: cpu.storeHL1},
		"PC":  proc.Value{Get: cpu.PC, Put: cpu.SetPC},
	}
}

func (c *CPU) Save(enc component.Encoder) error {
	e := check.ForError()
	//e.Check(c.Ports.Save(enc))

	e.Check(enc.Encode(c.A))
	e.Check(enc.Encode(c.F))
	e.Check(enc.Encode(c.B))
	e.Check(enc.Encode(c.C))
	e.Check(enc.Encode(c.D))
	e.Check(enc.Encode(c.H))
	e.Check(enc.Encode(c.L))

	e.Check(enc.Encode(c.A1))
	e.Check(enc.Encode(c.F1))
	e.Check(enc.Encode(c.B1))
	e.Check(enc.Encode(c.C1))
	e.Check(enc.Encode(c.D1))
	e.Check(enc.Encode(c.H1))
	e.Check(enc.Encode(c.L1))

	e.Check(enc.Encode(c.I))
	e.Check(enc.Encode(c.R))
	e.Check(enc.Encode(c.IXH))
	e.Check(enc.Encode(c.IXL))
	e.Check(enc.Encode(c.IYH))
	e.Check(enc.Encode(c.IYL))
	e.Check(enc.Encode(c.SP))
	e.Check(enc.Encode(c.pc))

	e.Check(enc.Encode(c.IFF1))
	e.Check(enc.Encode(c.IFF2))
	e.Check(enc.Encode(c.IM))
	e.Check(enc.Encode(c.Halt))

	return e.Error
}

func (c *CPU) Restore(dec component.Decoder) error {
	e := check.ForError()
	//e.Check(c.Ports.Restore(dec))

	e.Check(dec.Decode(&c.A))
	e.Check(dec.Decode(&c.F))
	e.Check(dec.Decode(&c.B))
	e.Check(dec.Decode(&c.C))
	e.Check(dec.Decode(&c.D))
	e.Check(dec.Decode(&c.H))
	e.Check(dec.Decode(&c.L))

	e.Check(dec.Decode(&c.A1))
	e.Check(dec.Decode(&c.F1))
	e.Check(dec.Decode(&c.B1))
	e.Check(dec.Decode(&c.C1))
	e.Check(dec.Decode(&c.D1))
	e.Check(dec.Decode(&c.H1))
	e.Check(dec.Decode(&c.L1))

	e.Check(dec.Decode(&c.I))
	e.Check(dec.Decode(&c.R))
	e.Check(dec.Decode(&c.IXH))
	e.Check(dec.Decode(&c.IXL))
	e.Check(dec.Decode(&c.IYH))
	e.Check(dec.Decode(&c.IYL))
	e.Check(dec.Decode(&c.SP))
	e.Check(dec.Decode(&c.pc))

	e.Check(dec.Decode(&c.IFF1))
	e.Check(dec.Decode(&c.IFF2))
	e.Check(dec.Decode(&c.IM))
	e.Check(dec.Decode(&c.Halt))

	return e.Error
}
