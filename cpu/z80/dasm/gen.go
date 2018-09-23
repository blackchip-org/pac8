package main

// http://www.z80.info/z80oplist.txt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	lineStart = 7
	lineEnd   = 262
)

var breaks = []int{
	3,
	7,
	18,
	22,
	33,
	38,
	45,
}

func dasm() {
	var out bytes.Buffer

	out.WriteString(`
// Code generated by cpu/z80/dasm/gen.go. DO NOT EDIT.

package z80

import "github.com/blackchip-org/pac8/cpu"

`)

	data, err := ioutil.ReadFile("dasm/z80oplist.txt")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")

	// unprefixed
	out.WriteString("var dasm = map[uint8]func(cpu.Eval){\n")
	for i := lineStart; i <= lineEnd; i++ {
		line := lines[i]
		line = strings.ToLower(line)
		parseTable(&out, line, 0, "")
	}
	out.WriteString("}\n")

	// dd prefix
	out.WriteString("var dasmDD = map[uint8]func(cpu.Eval){\n")
	for i := lineStart; i <= lineEnd; i++ {
		line := lines[i]
		line = strings.ToLower(line)
		parseTable(&out, line, 2, "dd")
	}
	out.WriteString("}\n")

	// fd prefix
	out.WriteString("var dasmFD = map[uint8]func(cpu.Eval){\n")
	for i := lineStart; i <= lineEnd; i++ {
		line := lines[i]
		line = strings.ToLower(line)
		parseTable(&out, line, 2, "fd")
	}
	out.WriteString("}\n")

	// cb prefix
	out.WriteString("var dasmCB = map[uint8]func(cpu.Eval){\n")
	for i := lineStart; i <= lineEnd; i++ {
		line := lines[i]
		line = strings.ToLower(line)
		parseTable(&out, line, 4, "cb")
	}
	out.WriteString("}\n")

	err = ioutil.WriteFile("dasm.go", out.Bytes(), 0644)
	if err != nil {
		fmt.Printf("unable to write file: %v", err)
		os.Exit(1)
	}
}

func parseTable(out *bytes.Buffer, line string, firstBreak int, prefix string) {
	break1 := breaks[firstBreak]
	break2 := breaks[firstBreak+1]
	break3 := breaks[firstBreak+2]

	strOpcode := strings.TrimSpace(line[0:2])
	opcode, _ := strconv.ParseUint(strOpcode, 16, 8)

	switch {
	case prefix == "" && opcode == 0xcb:
		out.WriteString("0xcb: func(e cpu.Eval) { opCB(e) },\n")
		return
	case prefix == "" && opcode == 0xdd:
		out.WriteString("0xdd: func(e cpu.Eval) { opDD(e) },\n")
		return
	case prefix == "" && opcode == 0xfd:
		out.WriteString("0xfd: func(e cpu.Eval) { opFD(e) },\n")
		return
	}

	out.WriteString("0x")
	out.WriteString(fmt.Sprintf("%02x", opcode))
	out.WriteString(": func(e cpu.Eval) { op1(e, ")

	args := make([]string, 1)
	args[0] = `"` + strings.TrimSpace(line[break1:break2]) + `"`

	if args[0] == `"-"` {
		out.WriteString(fmt.Sprintf(`"?%v%02x"`, prefix, opcode))
	} else {
		fields := strings.Split(line[break2:break3], ",")
		for _, field := range fields {
			args = append(args, `"`+strings.TrimSpace(field)+`"`)
		}
		entry := strings.Join(args, ",")
		if prefix == "fd" {
			entry = strings.Replace(entry, "ix", "iy", -1)
		}
		out.WriteString(entry)
	}

	out.WriteString(") },\n")
}

func harston() {
	var out bytes.Buffer

	out.WriteString(`
// Code generated by cpu/z80/dasm/gen.go. DO NOT EDIT.

package z80

type harstonTest struct {
	name string
	op string
	bytes []uint8
}

var harstonTests = []harstonTest{
`)

	data, err := ioutil.ReadFile("dasm/expected.txt")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}
		if line[0] == '=' {
			break
		}
		if line[0] == '#' {
			continue
		}
		data := strings.Split(line, " ")
		strdata := strings.Join(data, " ")
		hexdata := "0x" + strings.Join(data, ", 0x")
		i++
		op := lines[i]
		out.WriteString(fmt.Sprintf(`harstonTest{"%v", "%v", []uint8{%v}},`, strdata, op, hexdata))
		out.WriteString("\n")
	}
	out.WriteString("}\n")

	err = ioutil.WriteFile("harston_test.go", out.Bytes(), 0644)
	if err != nil {
		fmt.Printf("unable to write file: %v", err)
		os.Exit(1)
	}
}

func main() {
	dasm()
	harston()
}
