package z80

import (
	"fmt"
	"strings"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/cpu"
)

//go:generate go run dasm/gen.go
//go:generate go fmt dasm.go
//go:generate go fmt harston_test.go

func ReaderZ80(e cpu.Eval) cpu.Statement {
	e.Statement.Address = e.Cursor.Pos
	opcode := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, opcode)
	dasm[opcode](e)
	return *e.Statement
}

func op1(e cpu.Eval, parts ...string) {
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
			addr := bits.Displace(e.Statement.Address, delta)
			v = fmt.Sprintf("$%04x", addr)
		case part == "&0000":
			addr := e.Cursor.FetchLE()
			v = fmt.Sprintf("$%04x", addr)
		case part == "(&0000)":
			addr := e.Cursor.FetchLE()
			v = fmt.Sprintf("($%04x)", addr)
		case part == "&00":
			arg := e.Cursor.Fetch()
			v = fmt.Sprintf("$%02x", arg)
		case part == "(&00)":
			arg := e.Cursor.Fetch()
			v = fmt.Sprintf("($%02x)", arg)
		case part == "(ix+0)":
			delta := e.Cursor.Fetch()
			v = fmt.Sprintf("(ix+$%02x)", delta)
		case part == "(iy+0)":
			delta := e.Cursor.Fetch()
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

func op2(e cpu.Eval, parts ...string) {
	var out strings.Builder
	for i, part := range parts {
		v := part
		switch {
		case i == 0:
			v = fmt.Sprintf("%-4s", part)
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

func opDD(e cpu.Eval) {
	next := e.Cursor.Peek()
	if next == 0xdd || next == 0xed || next == 0xfd {
		e.Statement.Op = "?dd"
		return
	}
	opcode := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, opcode)
	dasmDD[opcode](e)
}

func opFD(e cpu.Eval) {
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
	dasmFD[opcode](e)
}

func opCB(e cpu.Eval) {
	opcode := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, opcode)
	dasmCB[opcode](e)
}

func opFDCB(e cpu.Eval) {
	delta := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, delta)
	opcode := e.Cursor.Fetch()
	e.Statement.Bytes = append(e.Statement.Bytes, opcode)
	dasmFDCB[opcode](e)
}
