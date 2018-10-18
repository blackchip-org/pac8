package mach

import (
	"fmt"
	"time"

	"github.com/blackchip-org/pac8/cpu"
	"github.com/blackchip-org/pac8/memory"
	"github.com/veandco/go-sdl2/sdl"
)

type fixtureCPU struct {
	A      uint8
	BC     uint16
	mem    memory.Memory
	cursor *memory.Cursor
	info   cpu.Info
}

func newFixtureCPU(mem memory.Memory) *fixtureCPU {
	c := &fixtureCPU{
		mem:    mem,
		cursor: memory.NewCursor(mem),
	}
	c.info = cpu.Info{
		CycleRate:       1,
		CodeReader:      fixtureReader,
		CodeFormatter:   fixtureFormatter(),
		NewDisassembler: NewDisassembler,
		Registers:       c.registers(),
	}
	return c
}

func (c *fixtureCPU) Next() {
	opcode := c.cursor.Fetch()
	args := opcode >> 4
	if args == 1 {
		c.cursor.Fetch()
	} else if args == 2 {
		c.cursor.FetchLE()
	}
}

func (c *fixtureCPU) PC() uint16 {
	return c.cursor.Pos
}

func (c *fixtureCPU) SetPC(pc uint16) {
	c.cursor.Pos = pc
}

func (c *fixtureCPU) Ready() bool {
	return true
}

func (c *fixtureCPU) Info() cpu.Info {
	return c.info
}

func (c *fixtureCPU) String() string {
	return fmt.Sprintf(""+
		" pc  a  bc\n"+
		"%04x %02x %04x\n",
		c.cursor.Pos,
		c.A,
		c.BC)
}

func (c *fixtureCPU) registers() map[string]cpu.Value {
	return map[string]cpu.Value{
		"A": cpu.Value{
			Get: func() uint8 { return c.A },
			Put: func(v uint8) { c.A = v },
		},
		"BC": cpu.Value{
			Get: func() uint16 { return c.BC },
			Put: func(v uint16) { c.BC = v },
		},
		"PC": cpu.Value{Get: c.PC, Put: c.SetPC},
	}
}

func fixtureReader(e cpu.Eval) cpu.Statement {
	e.Statement.Address = e.Cursor.Pos
	opcode := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, opcode)
	argN := opcode >> 4
	switch argN {
	case 0:
		e.Statement.Op = fmt.Sprintf("i%02x", opcode)
	case 1:
		value := e.Cursor.Fetch()
		e.Statement.Bytes = append(e.Statement.Bytes, value)
		e.Statement.Op = fmt.Sprintf("i%02x $%02x", opcode, value)
	case 2:
		value := e.Cursor.FetchLE()
		e.Statement.Bytes = append(e.Statement.Bytes, uint8(value&0xff))
		e.Statement.Bytes = append(e.Statement.Bytes, uint8(value>>8))
		e.Statement.Op = fmt.Sprintf("i%02x $%04x", opcode, value)
	default:
		e.Statement.Op = fmt.Sprintf("?%02x", opcode)
	}
	return *e.Statement
}

func fixtureFormatter() cpu.CodeFormatter {
	options := cpu.FormatOptions{
		BytesFormat: "%-8s",
	}
	return func(s cpu.Statement) string {
		return cpu.Format(s, options)
	}
}

func NewDisassembler(mem memory.Memory) *cpu.Disassembler {
	return cpu.NewDisassembler(mem, fixtureReader, fixtureFormatter())
}

type fixtureCab struct {
	mem memory.Memory
	cpu cpu.CPU
}

func newFixtureCab(renderer *sdl.Renderer) *Mach {
	cab := &fixtureCab{}
	cab.mem = memory.NewMasked(memory.NewRAM(0x1000), 0x0fff)
	cab.cpu = newFixtureCPU(cab.mem)
	m := New(cab.mem, cab.cpu)
	m.TickRate = time.Duration(10 * time.Millisecond)
	m.TickCallback = func(m *Mach) {
		if m.Status == Run && cab.mem.Load(cab.cpu.PC()) == 0x00 {
			m.Quit()
		}
		if m.Status == Breakpoint {
			m.Quit()
		}
	}
	return m
}
