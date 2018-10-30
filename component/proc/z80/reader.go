package z80

import (
	"fmt"
	"strings"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/proc"
)

func ReaderZ80(e proc.Eval) proc.Statement {
	e.Statement.Address = e.Cursor.Pos
	opcode := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, opcode)
	dasmTable[opcode](e)
	return *e.Statement
}

func FormatterZ80() proc.CodeFormatter {
	options := proc.FormatOptions{
		BytesFormat: "%-11s",
	}
	return func(s proc.Statement) string {
		return proc.Format(s, options)
	}
}

func NewDisassembler(mem memory.Memory) *proc.Disassembler {
	return proc.NewDisassembler(mem, ReaderZ80, FormatterZ80())
}

func op1(e proc.Eval, parts ...string) {
	var out strings.Builder
	for i, part := range parts {
		v := part
		switch {
		case i == 0:
			v = fmt.Sprintf("%-4s", part)
		case parts[0] == "rst" && i == 1:
			// Reset statements have the argment encoded in the opcode. Change
			// the hex notation from & to $ in the second part
			v = "$" + v[1:]
		case part == "&4546":
			// This is an address that is a 8-bit displacement from the
			// current program counter
			delta := e.Cursor.Fetch()
			e.Statement.Bytes = append(e.Statement.Bytes, delta)
			addr := bits.Displace(e.Statement.Address+2, delta)
			v = fmt.Sprintf("$%04x", addr)
		case part == "&0000":
			lo := e.Cursor.Fetch()
			e.Statement.Bytes = append(e.Statement.Bytes, lo)
			hi := e.Cursor.Fetch()
			e.Statement.Bytes = append(e.Statement.Bytes, hi)
			addr := bits.Join(hi, lo)
			v = fmt.Sprintf("$%04x", addr)
		case part == "(&0000)":
			lo := e.Cursor.Fetch()
			e.Statement.Bytes = append(e.Statement.Bytes, lo)
			hi := e.Cursor.Fetch()
			e.Statement.Bytes = append(e.Statement.Bytes, hi)
			addr := bits.Join(hi, lo)
			v = fmt.Sprintf("($%04x)", addr)
		case part == "&00":
			arg := e.Cursor.Fetch()
			e.Statement.Bytes = append(e.Statement.Bytes, arg)
			v = fmt.Sprintf("$%02x", arg)
		case part == "(&00)":
			arg := e.Cursor.Fetch()
			e.Statement.Bytes = append(e.Statement.Bytes, arg)
			v = fmt.Sprintf("($%02x)", arg)
		case part == "(ix+0)":
			delta := e.Cursor.Fetch()
			e.Statement.Bytes = append(e.Statement.Bytes, delta)
			v = fmt.Sprintf("(ix+$%02x)", delta)
		case part == "(iy+0)":
			delta := e.Cursor.Fetch()
			e.Statement.Bytes = append(e.Statement.Bytes, delta)
			v = fmt.Sprintf("(iy+$%02x)", delta)
		}

		if i == 1 {
			out.WriteString(" ")
		}
		if i == 2 {
			out.WriteString(",")
		}
		out.WriteString(v)
	}
	e.Statement.Op = strings.TrimSpace(out.String())
}

func op2(e proc.Eval, parts ...string) {
	var out strings.Builder
	for i, part := range parts {
		v := part
		switch {
		case i == 0:
			v = fmt.Sprintf("%-4s", part)
		case part == "(ix+0)":
			delta := e.Statement.Bytes[len(e.Statement.Bytes)-2]
			v = fmt.Sprintf("(ix+$%02x)", delta)
		case part == "(iy+0)":
			delta := e.Statement.Bytes[len(e.Statement.Bytes)-2]
			v = fmt.Sprintf("(iy+$%02x)", delta)
		}

		if i == 1 {
			out.WriteString(" ")
		}
		if i == 2 {
			out.WriteString(",")
		}
		out.WriteString(v)
	}
	e.Statement.Op = strings.TrimSpace(out.String())
}

func opDD(e proc.Eval) {
	next := e.Cursor.Peek()
	if next == 0xdd || next == 0xed || next == 0xfd {
		e.Statement.Op = "?dd"
		return
	}
	opcode := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, opcode)
	if opcode == 0xcb {
		opDDCB(e)
		return
	}
	dasmTableDD[opcode](e)
}

func opFD(e proc.Eval) {
	next := e.Cursor.Peek()
	if next == 0xdd || next == 0xed || next == 0xfd {
		e.Statement.Op = "?fd"
		return
	}
	opcode := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, opcode)
	if opcode == 0xcb {
		opFDCB(e)
		return
	}
	dasmTableFD[opcode](e)
}

func opCB(e proc.Eval) {
	opcode := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, opcode)
	dasmTableCB[opcode](e)
}

func opED(e proc.Eval) {
	opcode := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, opcode)
	dasmTableED[opcode](e)
}

func opFDCB(e proc.Eval) {
	delta := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, delta)
	opcode := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, opcode)
	dasmTableFDCB[opcode](e)
}

func opDDCB(e proc.Eval) {
	delta := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, delta)
	opcode := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, opcode)
	dasmTableDDCB[opcode](e)
}
