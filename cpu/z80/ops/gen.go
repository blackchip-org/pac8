package main

// http://www.z80.info/decoding.htm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blackchip-org/pac8/bits"
)

type regtab struct {
	name string
	r    map[int]string
	rp   map[int]string
	rp2  map[int]string
}

// unprefixed registers
var un = &regtab{
	name: "un",
	r: map[int]string{
		0: "B",
		1: "C",
		2: "D",
		3: "E",
		4: "H",
		5: "L",
		6: "IndHL",
		7: "A",
	},
	rp: map[int]string{
		0: "BC",
		1: "DE",
		2: "HL",
		3: "SP",
	},
	rp2: map[int]string{
		0: "BC",
		1: "DE",
		2: "HL",
		3: "AF",
	},
}

// dd-prefixed registers
var dd = &regtab{
	name: "dd",
	r: map[int]string{
		0: "B",
		1: "C",
		2: "D",
		3: "E",
		4: "IXH",
		5: "IXL",
		6: "IndIX",
		7: "A",
	},
	rp: map[int]string{
		0: "BC",
		1: "DE",
		2: "IX",
		3: "SP",
	},
	rp2: map[int]string{
		0: "BC",
		1: "DE",
		2: "IX",
		3: "AF",
	},
}

// fd-prefixed registers
var fd = &regtab{
	name: "fd",
	r: map[int]string{
		0: "B",
		1: "C",
		2: "D",
		3: "E",
		4: "IYH",
		5: "IYL",
		6: "IndIY",
		7: "A",
	},
	rp: map[int]string{
		0: "BC",
		1: "DE",
		2: "IY",
		3: "SP",
	},
	rp2: map[int]string{
		0: "BC",
		1: "DE",
		2: "IY",
		3: "AF",
	},
}

// ddcb-prefixed registers
var ddcb = &regtab{
	name: "ddcb",
	r: map[int]string{
		0: "B",
		1: "C",
		2: "D",
		3: "E",
		4: "IXH",
		5: "IXL",
		6: "IndIX",
		7: "A",
	},
	rp:  nil,
	rp2: nil,
}

// fdcb-prefixed registers
var fdcb = &regtab{
	name: "fdcb",
	r: map[int]string{
		0: "B",
		1: "C",
		2: "D",
		3: "E",
		4: "IYH",
		5: "IYL",
		6: "IndIY",
		7: "A",
	},
	rp:  nil,
	rp2: nil,
}

var cc = map[int]string{
	0: "FlagZ, false",
	1: "FlagZ, true",
	2: "FlagC, false",
	3: "FlagC, true",
	4: "FlagV, false",
	5: "FlagV, true",
	6: "FlagS, false",
	7: "FlagS, true",
}

func processMain(tab *regtab, op uint8) string {
	r := tab.r
	rp := tab.rp
	rp2 := tab.rp2

	x := int(bits.Slice(op, 6, 7))
	y := int(bits.Slice(op, 3, 5))
	z := int(bits.Slice(op, 0, 2))
	p := int(bits.Slice(op, 4, 5))
	q := int(bits.Slice(op, 3, 3))

	if x == 0 {
		if z == 0 {
			if y == 0 {
				return "nop()"
			}
			if y == 1 {
				return "ex(c, c.loadAF, c.storeAF, c.loadAF1, c.storeAF1)"
			}
			if y == 2 {
				return "djnz(c, c.loadImm)"
			}
			if y == 3 {
				return "jra(c, c.loadImm)"
			}
			if y >= 4 && y <= 7 {
				return fmt.Sprintf("jr(c, %v, c.loadImm)", cc[y-4])
			}
		}
		if z == 1 {
			if q == 0 {
				return fmt.Sprintf("ld16(c, c.store%v, c.loadImm16)", rp[p])
			}
			if q == 1 {
				return fmt.Sprintf("add16(c, c.store%v, c.load%v, c.load%v, false)", rp2[2], rp[2], rp[p])
			}
		}
		if z == 2 {
			if q == 0 {
				if p == 0 {
					return "ld(c, c.storeIndBC, c.loadA)"
				}
				if p == 1 {
					return "ld(c, c.storeIndDE, c.loadA)"
				}
				if p == 2 {
					return fmt.Sprintf("ld16(c, c.store16IndImm, c.load%v)", rp2[2])
				}
				if p == 3 {
					return "ld(c, c.storeIndImm, c.loadA)"
				}
			}
			if q == 1 {
				if p == 0 {
					return "ld(c, c.storeA, c.loadIndBC)"
				}
				if p == 1 {
					return "ld(c, c.storeA, c.loadIndDE)"
				}
				if p == 2 {
					return fmt.Sprintf("ld16(c, c.store%v, c.load16IndImm)", rp2[2])
				}
				if p == 3 {
					return "ld(c, c.storeA, c.loadIndImm)"
				}
			}
		}
		if z == 3 {
			if q == 0 {
				return fmt.Sprintf("inc16(c, c.store%v, c.load%v)", rp[p], rp[p])
			}
			if q == 1 {
				return fmt.Sprintf("dec16(c, c.store%v, c.load%v)", rp[p], rp[p])
			}
		}
		if z == 4 {
			return fmt.Sprintf("inc(c, c.store%v, c.load%v)", r[y], r[y])
		}
		if z == 5 {
			return fmt.Sprintf("dec(c, c.store%v, c.load%v)", r[y], r[y])
		}
		if z == 6 {
			return fmt.Sprintf("ld(c, c.store%v, c.loadImm)", r[y])
		}
		if z == 7 {
			if y == 0 {
				// rlca
				return "rotla(c)"
			}
			if y == 1 {
				// rrca
				return "rotra(c)"
			}
			if y == 2 {
				// rla
				return "shiftla(c)"
			}
			if y == 3 {
				// rra
				return "shiftra(c)"
			}
			if y == 4 {
				return "daa(c)"
			}
			if y == 5 {
				return "cpl(c)"
			}
			if y == 6 {
				return "scf(c)"
			}
			if y == 7 {
				return "ccf(c)"
			}
		}
	}
	if x == 1 {
		if z == 6 && y == 6 {
			return "halt(c)"
		}
		return fmt.Sprintf("ld(c, c.store%v, c.load%v)", r[y], r[z])
	}
	if x == 2 {
		if y == 0 {
			return fmt.Sprintf("add(c, c.loadA, c.load%v, false)", r[z])
		}
		if y == 1 {
			return fmt.Sprintf("add(c, c.loadA, c.load%v, true)", r[z])
		}
		if y == 2 {
			return fmt.Sprintf("sub(c, c.load%v, false)", r[z])
		}
		if y == 3 {
			return fmt.Sprintf("sub(c, c.load%v, true)", r[z])
		}
		if y == 4 {
			return fmt.Sprintf("and(c, c.load%v)", r[z])
		}
		if y == 5 {
			return fmt.Sprintf("xor(c, c.load%v)", r[z])
		}
		if y == 6 {
			return fmt.Sprintf("or(c, c.load%v)", r[z])
		}
		if y == 7 {
			return fmt.Sprintf("cp(c, c.load%v)", r[z])
		}
	}
	if x == 3 {
		if z == 0 {
			return fmt.Sprintf("ret(c, %v)", cc[y])
		}
		if z == 1 {
			if q == 0 {
				return fmt.Sprintf("pop(c, c.store%v)", rp2[p])
			}
			if q == 1 {
				if p == 0 {
					return "reta(c)"
				}
				if p == 1 {
					return "exx(c)"
				}
				if p == 2 {
					return fmt.Sprintf("jpa(c, c.load%v)", rp2[2])
				}
				if p == 3 {
					return fmt.Sprintf("ld16(c, c.storeSP, c.load%v)", rp2[2])
				}
			}
		}
		if z == 2 {
			return fmt.Sprintf("jp(c, %v, c.loadImm16)", cc[y])
		}
		if z == 3 {
			if y == 0 {
				return "jpa(c, c.loadImm16)"
			}
			if y == 1 && tab.name == "un" {
				return "cb(c)"
			}
			if y == 1 && tab.name != "un" {
				return ""
			}
			if y == 2 {
				return "todo2(c, c.loadImm)"
			}
			if y == 3 {
				return "todo2(c, c.loadImm)"
			}
			if y == 4 {
				return fmt.Sprintf("ex(c, c.load16IndSP, c.store16IndSP, c.load%v, c.store%v)", rp2[2], rp2[2])
			}
			if y == 5 {
				return "ex(c, c.loadDE, c.storeDE, c.loadHL, c.storeHL)"
			}
			if y == 6 {
				return "todo(c)"
			}
			if y == 7 {
				return "todo(c)"
			}
		}
		if z == 4 {
			return fmt.Sprintf("call(c, %v, c.loadImm16)", cc[y])
		}
		if z == 5 {
			if q == 0 {
				return fmt.Sprintf("push(c, c.load%v)", rp2[p])
			}
			if q == 1 {
				if p == 0 {
					return "calla(c, c.loadImm16)"
				}
				if p == 1 && tab.name == "un" {
					return "ddfd(c, opsDD, opsDDCB)"
				}
				if p == 1 && tab.name != "un" {
					return "noni(c)"
				}
				if p == 2 {
					return "ed(c)"
				}
				if p == 3 && tab.name == "un" {
					return "ddfd(c, opsFD, opsFDCB)"
				}
				if p == 3 && tab.name != "un" {
					return "noni(c)"
				}
			}
		}
		if z == 6 {
			if y == 0 {
				return "add(c, c.loadImm, c.loadA, false)"
			}
			if y == 1 {
				return "add(c, c.loadImm, c.loadA, true)"
			}
			if y == 2 {
				return "sub(c, c.loadImm, false)"
			}
			if y == 3 {
				return "sub(c, c.loadImm, true)"
			}
			if y == 4 {
				return "and(c, c.loadImm)"
			}
			if y == 5 {
				return "xor(c, c.loadImm)"
			}
			if y == 6 {
				return "or(c, c.loadImm)"
			}
			if y == 7 {
				return "cp(c, c.loadImm)"
			}
		}
		if z == 7 {
			return fmt.Sprintf("rst(c, %v)", y)
		}
	}
	return ""
}

func processCB(tab *regtab, op uint8) string {
	r := tab.r
	x := int(bits.Slice(op, 6, 7))
	y := int(bits.Slice(op, 3, 5))
	z := int(bits.Slice(op, 0, 2))

	if x == 0 {
		if y == 0 {
			// rlc
			return fmt.Sprintf("rotl(c, c.store%v, c.load%v)", r[z], r[z])
		}
		if y == 1 {
			// rrc
			return fmt.Sprintf("rotr(c, c.store%v, c.load%v)", r[z], r[z])
		}
		if y == 2 {
			// rl
			return fmt.Sprintf("shiftl(c, c.store%v, c.load%v, true)", r[z], r[z])
		}
		if y == 3 {
			// rr
			return fmt.Sprintf("shiftr(c, c.store%v, c.load%v, true)", r[z], r[z])
		}
		if y == 4 {
			// sla
			return fmt.Sprintf("shiftl(c, c.store%v, c.load%v, false)", r[z], r[z])
		}
		if y == 5 {
			return fmt.Sprintf("sra(c, c.store%v, c.load%v)", r[z], r[z])
		}
		if y == 6 {
			return fmt.Sprintf("sll(c, c.store%v, c.load%v)", r[z], r[z])
		}
		if y == 7 {
			// srl
			return fmt.Sprintf("shiftr(c, c.store%v, c.load%v, false)", r[z], r[z])
		}
	}
	if x == 1 {
		return fmt.Sprintf("bit(c, %v, c.load%v)", y, r[z])
	}
	if x == 2 {
		return fmt.Sprintf("res(c, %v, c.store%v, c.load%v)", y, r[z], r[z])
	}
	if x == 3 {
		return fmt.Sprintf("set(c, %v, c.store%v, c.load%v)", y, r[z], r[z])
	}
	return ""
}

func processED(tab *regtab, op uint8) string {
	rp := tab.rp
	x := int(bits.Slice(op, 6, 7))
	y := int(bits.Slice(op, 3, 5))
	z := int(bits.Slice(op, 0, 2))
	p := int(bits.Slice(op, 4, 5))
	q := int(bits.Slice(op, 3, 3))

	if x == 0 || x == 3 {
		return "invalid()"
	}
	if x == 1 {
		if z == 0 {
			if y != 6 {
				return "todo(c)"
			}
			if y == 6 {
				return "todo(c)"
			}
		}
		if z == 1 {
			if y != 6 {
				return "todo(c)"
			}
			if y == 6 {
				return "todo(c)"
			}
		}
		if z == 2 {
			if q == 0 {
				return fmt.Sprintf("sub16(c, c.storeHL, c.loadHL, c.load%v, true)", rp[p])
			}
			if q == 1 {
				return fmt.Sprintf("add16(c, c.storeHL, c.loadHL, c.load%v, true)", rp[p])
			}
		}
		if z == 3 {
			if q == 0 {
				return fmt.Sprintf("ld16(c, c.store16IndImm, c.load%v)", rp[p])
			}
			if q == 1 {
				return fmt.Sprintf("ld16(c, c.store%v, c.load16IndImm)", rp[p])
			}
		}
		if z == 4 {
			return "neg(c)"
		}
		if z == 5 {
			return "todo(c)"
		}
		if z == 6 {
			return "todo(c)"
		}
		if z == 7 {
			if y == 0 {
				return "ld(c, c.storeI, c.loadA)"
			}
			if y == 1 {
				return "ld(c, c.storeR, c.loadA)"
			}
			if y == 2 {
				// ld a, i
				return "ldair(c, c.loadI)"
			}
			if y == 3 {
				// ld a, r
				return "ldair(c, c.loadR)"
			}
			if y == 4 {
				return "rrd(c)"
			}
			if y == 5 {
				return "rld(c)"
			}
			return "nop()"
		}
	}
	if x == 2 {
		if z == 0 { // b == 0
			if y == 4 {
				// ldi
				return "blockl(c, 1)"
			}
			if y == 5 {
				// ldd
				return "blockl(c, -1)"
			}
			if y == 6 {
				// ldir
				return "blocklr(c, 1)"
			}
			if y == 7 {
				// lddr
				return "blocklr(c, -1)"
			}
		}
		if z == 1 { // b == 1
			if y == 4 {
				// cpi
				return "blockc(c, inc16)"
			}
			if y == 5 {
				// cpd
				return "blockc(c, dec16)"
			}
			if y == 6 {
				// cpir
				return "blockcr(c, inc16)"
			}
			if y == 7 {
				// cpdr
				return "blockcr(c, dec16)"
			}
		}
		if z == 2 && y >= 4 && y <= 7 {
			return "todo(c)"
		}
		if z == 3 && y >= 4 && y <= 7 {
			return "todo(c)"
		}
		return "invalid()"
	}
	return ""
}

func processXCB(tab *regtab, op uint8) string {
	r := tab.r
	x := int(bits.Slice(op, 6, 7))
	y := int(bits.Slice(op, 3, 5))
	z := int(bits.Slice(op, 0, 2))
	// p := int(bits.Slice(op, 4, 5))
	// q := int(bits.Slice(op, 3, 3))

	if x == 0 {
		if z != 6 {
			if y == 0 {
				// rlc
				return fmt.Sprintf("rotl(c, c.store%v, c.load%v); ld(c, c.storeLastInd, c.load%v)", un.r[z], r[6], un.r[z])
			}
			if y == 1 {
				// rrc
				return fmt.Sprintf("rotr(c, c.store%v, c.load%v); ld(c, c.storeLastInd, c.load%v)", un.r[z], r[6], un.r[z])
			}
			if y == 2 {
				// rl
				return fmt.Sprintf("shiftl(c, c.store%v, c.load%v, true); ld(c, c.storeLastInd, c.load%v)", un.r[z], r[6], un.r[z])
			}
			if y == 3 {
				// rr
				return fmt.Sprintf("shiftr(c, c.store%v, c.load%v, true); ld(c, c.storeLastInd, c.load%v)", un.r[z], r[6], un.r[z])
			}
			if y == 4 {
				// sla
				return fmt.Sprintf("shiftl(c, c.store%v, c.load%v, false); ld(c, c.storeLastInd, c.load%v)", un.r[z], r[6], un.r[z])
			}
			if y == 5 {
				return fmt.Sprintf("sra(c, c.store%v, c.load%v); ld(c, c.storeLastInd, c.load%v)", un.r[z], r[6], un.r[z])
			}
			if y == 6 {
				return fmt.Sprintf("sll(c, c.store%v, c.load%v); ld(c, c.storeLastInd, c.load%v)", un.r[z], r[6], un.r[z])
			}
			if y == 7 {
				// srl
				return fmt.Sprintf("shiftr(c, c.store%v, c.load%v, false); ld(c, c.storeLastInd, c.load%v)", un.r[z], r[6], un.r[z])
			}
		}
		if z == 6 {
			if y == 0 {
				// rlc
				return fmt.Sprintf("rotl(c, c.storeLastInd, c.load%v)", r[6])
			}
			if y == 1 {
				// rrc
				return fmt.Sprintf("rotr(c, c.storeLastInd, c.load%v)", r[6])
			}
			if y == 2 {
				// rl
				return fmt.Sprintf("shiftl(c, c.storeLastInd, c.load%v, true)", r[6])
			}
			if y == 3 {
				// rr
				return fmt.Sprintf("shiftr(c, c.storeLastInd, c.load%v, true)", r[6])
			}
			if y == 4 {
				// sla
				return fmt.Sprintf("shiftl(c, c.storeLastInd, c.load%v, false)", r[6])
			}
			if y == 5 {
				return fmt.Sprintf("sra(c, c.storeLastInd, c.load%v)", r[6])
			}
			if y == 6 {
				return fmt.Sprintf("sll(c, c.storeLastInd, c.load%v)", r[6])
			}
			if y == 7 {
				// srl
				return fmt.Sprintf("shiftr(c, c.storeLastInd, c.load%v, false)", r[6])
			}
		}
	}
	if x == 1 {
		return fmt.Sprintf("biti(c, %v, c.load%v)", y, r[6])
	}
	if x == 2 {
		if z != 6 {
			return fmt.Sprintf("res(c, %v, c.store%v, c.load%v); ld(c, c.storeLastInd, c.load%v)", y, un.r[z], r[6], un.r[z])
		}
		if z == 6 {
			return fmt.Sprintf("res(c, %v, c.storeLastInd, c.load%v)", y, r[6])
		}
	}
	if x == 3 {
		if z != 6 {
			return fmt.Sprintf("set(c, %v, c.store%v, c.load%v); ld(c, c.storeLastInd, c.load%v)", y, un.r[z], r[6], un.r[z])
		}
		if z == 6 {
			return fmt.Sprintf("set(c, %v, c.storeLastInd, c.load%v)", y, r[6])
		}
	}
	return ""
}

func process(out *bytes.Buffer, getFn func(*regtab, uint8) string, tab *regtab) {
	for i := 0; i < 0x100; i++ {
		fn := getFn(tab, uint8(i))
		if fn == "" {
			fn = "c.skip = true"
		}

		emit := true
		if tab.name == "dd" || tab.name == "fd" {
			// If there is an indirect call, the next byte needs to be
			// fetched for the displacement
			if strings.Contains(fn, "IndIX") || strings.Contains(fn, "IndIY") {
				fn = "c.fetchd(); " + fn
				// "any other instances of H and L will be unaffected".
				// If (HL) was transformed, use H and L instead.
				fn = strings.Replace(fn, "IXH", "H", -1)
				fn = strings.Replace(fn, "IXL", "L", -1)
				fn = strings.Replace(fn, "IYH", "H", -1)
				fn = strings.Replace(fn, "IYL", "L", -1)
			}
			// If the unprefixed version of the op is the same as the
			// prefixed one, just leave the function blank
			unfn := getFn(un, uint8(i))
			if unfn == fn {
				emit = false
			}
		}

		line := fmt.Sprintf("0x%02x: func(c *CPU){%v},\n", i, fn)
		if !emit {
			line = fmt.Sprintf("0x%02x: nil,\n", i)
		}

		out.WriteString(line)
	}
}

func main() {
	var out bytes.Buffer

	out.WriteString(`
// Code generated by cpu/z80/ops/gen.go. DO NOT EDIT.

package z80

`)
	out.WriteString("var ops = map[uint8]func(c *CPU){\n")
	process(&out, processMain, un)
	out.WriteString("}\n")

	out.WriteString("var opsCB = map[uint8]func(c *CPU){\n")
	process(&out, processCB, un)
	out.WriteString("}\n")

	out.WriteString("var opsED = map[uint8]func(c *CPU){\n")
	process(&out, processED, un)
	out.WriteString("}\n")

	out.WriteString("var opsDD = map[uint8]func(c *CPU){\n")
	process(&out, processMain, dd)
	out.WriteString("}\n")

	out.WriteString("var opsFD = map[uint8]func(c *CPU){\n")
	process(&out, processMain, fd)
	out.WriteString("}\n")

	out.WriteString("var opsDDCB = map[uint8]func(c *CPU){\n")
	process(&out, processXCB, ddcb)
	out.WriteString("}\n")

	out.WriteString("var opsFDCB = map[uint8]func(c *CPU){\n")
	process(&out, processXCB, fdcb)
	out.WriteString("}\n")

	err := ioutil.WriteFile("ops.go", out.Bytes(), 0644)
	if err != nil {
		fmt.Printf("unable to write file: %v", err)
		os.Exit(1)
	}
}
