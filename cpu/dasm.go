package cpu

import (
	"github.com/blackchip-org/pac8/memory"
)

type Statement struct {
	Address uint16
	Label   string
	Op      string
	Bytes   []uint8
	Comment string
}

func NewStatement() *Statement {
	return &Statement{Bytes: make([]uint8, 0, 0)}
}

type CodeInfo struct{}

type CodeReader func(Eval) Statement

type Disassembler struct {
	CodeInfo *CodeInfo
	mem      memory.Memory
	cursor   *memory.Cursor
	read     CodeReader
}

type Eval struct {
	Cursor    *memory.Cursor
	CodeInfo  *CodeInfo
	Statement *Statement
}

func NewDisassembler(mem memory.Memory, r CodeReader) *Disassembler {
	return &Disassembler{
		CodeInfo: &CodeInfo{},
		mem:      mem,
		cursor:   memory.NewCursor(mem),
		read:     r,
	}
}

func (d *Disassembler) Next() Statement {
	return d.read(Eval{
		Cursor:    d.cursor,
		CodeInfo:  d.CodeInfo,
		Statement: NewStatement(),
	})
}

func (d *Disassembler) SetPC(addr uint16) {
	d.cursor.Pos = addr
}
