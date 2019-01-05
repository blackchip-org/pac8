package proc

import (
	"fmt"
	"strings"

	"github.com/blackchip-org/pac8/pkg/memory"
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
type CodeFormatter func(Statement) string

type Disassembler struct {
	CodeInfo *CodeInfo
	mem      memory.Memory
	cursor   *memory.Cursor
	read     CodeReader
	format   CodeFormatter
}

type Eval struct {
	Cursor    *memory.Cursor
	CodeInfo  *CodeInfo
	Statement *Statement
}

func NewDisassembler(mem memory.Memory, r CodeReader, f CodeFormatter) *Disassembler {
	return &Disassembler{
		CodeInfo: &CodeInfo{},
		mem:      mem,
		cursor:   memory.NewCursor(mem),
		read:     r,
		format:   f,
	}
}

func (d *Disassembler) NextStatement() Statement {
	return d.read(Eval{
		Cursor:    d.cursor,
		CodeInfo:  d.CodeInfo,
		Statement: NewStatement(),
	})
}

func (d *Disassembler) Next() string {
	return d.format(d.NextStatement())
}

func (d *Disassembler) SetPC(addr uint16) {
	d.cursor.Pos = addr
}

func (d *Disassembler) PC() uint16 {
	return d.cursor.Pos
}

type FormatOptions struct {
	BytesFormat string
}

func Format(s Statement, options FormatOptions) string {
	bytes := []string{}
	for _, b := range s.Bytes {
		bytes = append(bytes, fmt.Sprintf("%02x", b))
	}
	sbytes := fmt.Sprintf(options.BytesFormat, strings.Join(bytes, " "))
	return fmt.Sprintf("$%04x:  %s  %s", s.Address, sbytes, s.Op)
}
