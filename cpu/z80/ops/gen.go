package main

// http://www.z80.info/decoding.htm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/blackchip-org/pac8/bits"
)

var r = map[int]string{
	0: "B",
	1: "C",
	2: "D",
	3: "E",
	4: "H",
	5: "L",
	6: "IndHL",
	7: "A",
}

var rp = map[int]string{
	0: "BC",
	1: "DE",
	2: "HL",
	3: "SP",
}

var rp2 = map[int]string{
	0: "BC",
	1: "DE",
	2: "HL",
	3: "AF",
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

func process(op uint8) string {
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
				return fmt.Sprintf("add16(c, c.storeHL, c.loadHL, c.load%v)", rp[p])
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
					return "ld16(c, c.store16IndImm, c.loadHL)"
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
					return "ld16(c, c.storeHL, c.load16IndImm)"
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
				return "rlca(c)"
			}
			if y == 1 {
				return "rrca(c)"
			}
			if y == 2 {
				return "rla(c)"
			}
			if y == 3 {
				return "rra(c)"
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
			return fmt.Sprintf("add(c, c.loadA, c.load%v, c.F & 0x01 == 1)", r[z])
		}
		if y == 2 {
			return fmt.Sprintf("sub(c, c.load%v, false)", r[z])
		}
		if y == 3 {
			return fmt.Sprintf("sub(c, c.load%v, c.F & 0x01 == 1)", r[z])
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
					return "jpa(c, c.loadHL)"
				}
				if p == 3 {
					return "ld16(c, c.storeSP, c.loadHL)"
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
			if y == 2 {
				return "out(c)"
			}
			if y == 3 {
				return "in(c, c.storeA, c.loadImm)"
			}
			if y == 4 {
				return "ex(c, c.load16IndSP, c.store16IndSP, c.loadHL, c.storeHL)"
			}
			if y == 5 {
				return "ex(c, c.loadDE, c.storeDE, c.loadHL, c.storeHL)"
			}
			if y == 6 {
				return "di()"
			}
			if y == 7 {
				return "ei()"
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
			}
		}
		if z == 6 {
			if y == 0 {
				return "add(c, c.loadImm, c.loadA, false)"
			}
			if y == 1 {
				return "add(c, c.loadImm, c.loadA, c.F & 0x01 == 1)"
			}
			if y == 2 {
				return "sub(c, c.loadImm, false)"
			}
			if y == 3 {
				return "sub(c, c.loadImm, c.F & 0x01 == 1)"
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

func main() {
	var out bytes.Buffer

	out.WriteString(`
// Code generated by cpu/z80/ops/gen.go. DO NOT EDIT.

package z80

var ops = map[uint8]func(c *CPU){
`)
	for i := 0; i < 0x100; i++ {
		fn := process(uint8(i))
		if fn == "" {
			fn = "c.skip = true"
		}
		line := fmt.Sprintf("0x%02x: func(c *CPU){%v},\n", i, fn)
		out.WriteString(line)
	}
	out.WriteString("}\n")

	err := ioutil.WriteFile("ops.go", out.Bytes(), 0644)
	if err != nil {
		fmt.Printf("unable to write file: %v", err)
		os.Exit(1)
	}
}