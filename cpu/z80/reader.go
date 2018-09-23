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
		if i == 0 {
			v = fmt.Sprintf("%-4s", part)
		} else if part == "&4546" {
			delta := e.Cursor.Fetch()
			addr := bits.Displace(e.Statement.Address, delta)
			v = fmt.Sprintf("$%04x", addr)
		} else if part == "&0000" {
			addr := e.Cursor.FetchLE()
			v = fmt.Sprintf("$%04x", addr)
		} else if part == "(&0000)" {
			addr := e.Cursor.FetchLE()
			v = fmt.Sprintf("($%04x)", addr)
		} else if part == "&00" {
			arg := e.Cursor.Fetch()
			v = fmt.Sprintf("$%02x", arg)
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
