package app

import (
	"fmt"
	"time"

	"github.com/blackchip-org/pac8/pkg/audio"
	"github.com/blackchip-org/pac8/pkg/machine"
	"github.com/blackchip-org/pac8/pkg/memory"
	"github.com/blackchip-org/pac8/pkg/proc"
	"github.com/blackchip-org/pac8/pkg/util/state"
	"github.com/blackchip-org/pac8/pkg/video"
	"github.com/veandco/go-sdl2/sdl"
)

type fixtureCPU struct {
	A      uint8
	BC     uint16
	mem    memory.Memory
	cursor *memory.Cursor
	info   proc.Info
}

func newFixtureCPU(mem memory.Memory) *fixtureCPU {
	c := &fixtureCPU{
		mem:    mem,
		cursor: memory.NewCursor(mem),
	}
	c.info = proc.Info{
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

func (c *fixtureCPU) Info() proc.Info {
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

func (c *fixtureCPU) registers() map[string]proc.Value {
	return map[string]proc.Value{
		"A": proc.Value{
			Get: func() uint8 { return c.A },
			Put: func(v uint8) { c.A = v },
		},
		"BC": proc.Value{
			Get: func() uint16 { return c.BC },
			Put: func(v uint16) { c.BC = v },
		},
		"PC": proc.Value{Get: c.PC, Put: c.SetPC},
	}
}

func (c *fixtureCPU) Save(*state.Encoder) {}

func (c *fixtureCPU) Restore(*state.Decoder) {}

func fixtureReader(e proc.Eval) proc.Statement {
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

func fixtureFormatter() proc.CodeFormatter {
	options := proc.FormatOptions{
		BytesFormat: "%-8s",
	}
	return func(s proc.Statement) string {
		return proc.Format(s, options)
	}
}

func NewDisassembler(mem memory.Memory) *proc.Disassembler {
	return proc.NewDisassembler(mem, fixtureReader, fixtureFormatter())
}

type fixtureSys struct {
	mem memory.Memory
	cpu proc.CPU
}

func (f fixtureSys) Spec() *machine.Spec {
	callback := func(m *machine.Mach) {
		if m.Status == machine.Run && f.mem.Load(f.cpu.PC()) == 0x00 {
			m.Send(machine.QuitCmd)
		}
		if m.Status == machine.Break {
			m.Send(machine.QuitCmd)
		}
	}
	return &machine.Spec{
		CPU:          []proc.CPU{f.cpu},
		Mem:          []memory.Memory{f.mem},
		Display:      video.NullDisplay{},
		Audio:        audio.NullAudio{},
		TickCallback: callback,
		TickRate:     1 * time.Millisecond,
		CharDecoder:  AsciiDecoder,
	}
}

func (f fixtureSys) Save(*state.Encoder) {}

func (f fixtureSys) Restore(*state.Decoder) {}

func newFixtureCab(renderer *sdl.Renderer) machine.System {
	sys := &fixtureSys{}
	sys.mem = memory.NewRAM(0x1000)
	sys.cpu = newFixtureCPU(sys.mem)
	return sys
}
