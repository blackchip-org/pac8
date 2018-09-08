// Code generated by cpu/z80/ops/gen.go. DO NOT EDIT.

package z80

var ops = map[uint8]func(c *CPU){
	0x00: func(c *CPU) { nop() },
	0x01: func(c *CPU) { ld16(c, c.storeBC, c.loadImm16) },
	0x02: func(c *CPU) { ld(c, c.storeIndBC, c.loadA) },
	0x03: func(c *CPU) { inc16(c, c.storeBC, c.loadBC) },
	0x04: func(c *CPU) { inc(c, c.storeB, c.loadB) },
	0x05: func(c *CPU) { dec(c, c.storeB, c.loadB) },
	0x06: func(c *CPU) { ld(c, c.storeB, c.loadImm) },
	0x07: func(c *CPU) { rotla(c) },
	0x08: func(c *CPU) { ex(c, c.loadAF, c.storeAF, c.loadAF1, c.storeAF1) },
	0x09: func(c *CPU) { add16(c, c.storeHL, c.loadHL, c.loadBC, false) },
	0x0a: func(c *CPU) { ld(c, c.storeA, c.loadIndBC) },
	0x0b: func(c *CPU) { dec16(c, c.storeBC, c.loadBC) },
	0x0c: func(c *CPU) { inc(c, c.storeC, c.loadC) },
	0x0d: func(c *CPU) { dec(c, c.storeC, c.loadC) },
	0x0e: func(c *CPU) { ld(c, c.storeC, c.loadImm) },
	0x0f: func(c *CPU) { rotra(c) },
	0x10: func(c *CPU) { djnz(c, c.loadImm) },
	0x11: func(c *CPU) { ld16(c, c.storeDE, c.loadImm16) },
	0x12: func(c *CPU) { ld(c, c.storeIndDE, c.loadA) },
	0x13: func(c *CPU) { inc16(c, c.storeDE, c.loadDE) },
	0x14: func(c *CPU) { inc(c, c.storeD, c.loadD) },
	0x15: func(c *CPU) { dec(c, c.storeD, c.loadD) },
	0x16: func(c *CPU) { ld(c, c.storeD, c.loadImm) },
	0x17: func(c *CPU) { shiftla(c) },
	0x18: func(c *CPU) { jra(c, c.loadImm) },
	0x19: func(c *CPU) { add16(c, c.storeHL, c.loadHL, c.loadDE, false) },
	0x1a: func(c *CPU) { ld(c, c.storeA, c.loadIndDE) },
	0x1b: func(c *CPU) { dec16(c, c.storeDE, c.loadDE) },
	0x1c: func(c *CPU) { inc(c, c.storeE, c.loadE) },
	0x1d: func(c *CPU) { dec(c, c.storeE, c.loadE) },
	0x1e: func(c *CPU) { ld(c, c.storeE, c.loadImm) },
	0x1f: func(c *CPU) { shiftra(c) },
	0x20: func(c *CPU) { jr(c, FlagZ, false, c.loadImm) },
	0x21: func(c *CPU) { ld16(c, c.storeHL, c.loadImm16) },
	0x22: func(c *CPU) { ld16(c, c.store16IndImm, c.loadHL) },
	0x23: func(c *CPU) { inc16(c, c.storeHL, c.loadHL) },
	0x24: func(c *CPU) { inc(c, c.storeH, c.loadH) },
	0x25: func(c *CPU) { dec(c, c.storeH, c.loadH) },
	0x26: func(c *CPU) { ld(c, c.storeH, c.loadImm) },
	0x27: func(c *CPU) { daa(c) },
	0x28: func(c *CPU) { jr(c, FlagZ, true, c.loadImm) },
	0x29: func(c *CPU) { add16(c, c.storeHL, c.loadHL, c.loadHL, false) },
	0x2a: func(c *CPU) { ld16(c, c.storeHL, c.load16IndImm) },
	0x2b: func(c *CPU) { dec16(c, c.storeHL, c.loadHL) },
	0x2c: func(c *CPU) { inc(c, c.storeL, c.loadL) },
	0x2d: func(c *CPU) { dec(c, c.storeL, c.loadL) },
	0x2e: func(c *CPU) { ld(c, c.storeL, c.loadImm) },
	0x2f: func(c *CPU) { cpl(c) },
	0x30: func(c *CPU) { jr(c, FlagC, false, c.loadImm) },
	0x31: func(c *CPU) { ld16(c, c.storeSP, c.loadImm16) },
	0x32: func(c *CPU) { ld(c, c.storeIndImm, c.loadA) },
	0x33: func(c *CPU) { inc16(c, c.storeSP, c.loadSP) },
	0x34: func(c *CPU) { inc(c, c.storeIndHL, c.loadIndHL) },
	0x35: func(c *CPU) { dec(c, c.storeIndHL, c.loadIndHL) },
	0x36: func(c *CPU) { ld(c, c.storeIndHL, c.loadImm) },
	0x37: func(c *CPU) { scf(c) },
	0x38: func(c *CPU) { jr(c, FlagC, true, c.loadImm) },
	0x39: func(c *CPU) { add16(c, c.storeHL, c.loadHL, c.loadSP, false) },
	0x3a: func(c *CPU) { ld(c, c.storeA, c.loadIndImm) },
	0x3b: func(c *CPU) { dec16(c, c.storeSP, c.loadSP) },
	0x3c: func(c *CPU) { inc(c, c.storeA, c.loadA) },
	0x3d: func(c *CPU) { dec(c, c.storeA, c.loadA) },
	0x3e: func(c *CPU) { ld(c, c.storeA, c.loadImm) },
	0x3f: func(c *CPU) { ccf(c) },
	0x40: func(c *CPU) { ld(c, c.storeB, c.loadB) },
	0x41: func(c *CPU) { ld(c, c.storeB, c.loadC) },
	0x42: func(c *CPU) { ld(c, c.storeB, c.loadD) },
	0x43: func(c *CPU) { ld(c, c.storeB, c.loadE) },
	0x44: func(c *CPU) { ld(c, c.storeB, c.loadH) },
	0x45: func(c *CPU) { ld(c, c.storeB, c.loadL) },
	0x46: func(c *CPU) { ld(c, c.storeB, c.loadIndHL) },
	0x47: func(c *CPU) { ld(c, c.storeB, c.loadA) },
	0x48: func(c *CPU) { ld(c, c.storeC, c.loadB) },
	0x49: func(c *CPU) { ld(c, c.storeC, c.loadC) },
	0x4a: func(c *CPU) { ld(c, c.storeC, c.loadD) },
	0x4b: func(c *CPU) { ld(c, c.storeC, c.loadE) },
	0x4c: func(c *CPU) { ld(c, c.storeC, c.loadH) },
	0x4d: func(c *CPU) { ld(c, c.storeC, c.loadL) },
	0x4e: func(c *CPU) { ld(c, c.storeC, c.loadIndHL) },
	0x4f: func(c *CPU) { ld(c, c.storeC, c.loadA) },
	0x50: func(c *CPU) { ld(c, c.storeD, c.loadB) },
	0x51: func(c *CPU) { ld(c, c.storeD, c.loadC) },
	0x52: func(c *CPU) { ld(c, c.storeD, c.loadD) },
	0x53: func(c *CPU) { ld(c, c.storeD, c.loadE) },
	0x54: func(c *CPU) { ld(c, c.storeD, c.loadH) },
	0x55: func(c *CPU) { ld(c, c.storeD, c.loadL) },
	0x56: func(c *CPU) { ld(c, c.storeD, c.loadIndHL) },
	0x57: func(c *CPU) { ld(c, c.storeD, c.loadA) },
	0x58: func(c *CPU) { ld(c, c.storeE, c.loadB) },
	0x59: func(c *CPU) { ld(c, c.storeE, c.loadC) },
	0x5a: func(c *CPU) { ld(c, c.storeE, c.loadD) },
	0x5b: func(c *CPU) { ld(c, c.storeE, c.loadE) },
	0x5c: func(c *CPU) { ld(c, c.storeE, c.loadH) },
	0x5d: func(c *CPU) { ld(c, c.storeE, c.loadL) },
	0x5e: func(c *CPU) { ld(c, c.storeE, c.loadIndHL) },
	0x5f: func(c *CPU) { ld(c, c.storeE, c.loadA) },
	0x60: func(c *CPU) { ld(c, c.storeH, c.loadB) },
	0x61: func(c *CPU) { ld(c, c.storeH, c.loadC) },
	0x62: func(c *CPU) { ld(c, c.storeH, c.loadD) },
	0x63: func(c *CPU) { ld(c, c.storeH, c.loadE) },
	0x64: func(c *CPU) { ld(c, c.storeH, c.loadH) },
	0x65: func(c *CPU) { ld(c, c.storeH, c.loadL) },
	0x66: func(c *CPU) { ld(c, c.storeH, c.loadIndHL) },
	0x67: func(c *CPU) { ld(c, c.storeH, c.loadA) },
	0x68: func(c *CPU) { ld(c, c.storeL, c.loadB) },
	0x69: func(c *CPU) { ld(c, c.storeL, c.loadC) },
	0x6a: func(c *CPU) { ld(c, c.storeL, c.loadD) },
	0x6b: func(c *CPU) { ld(c, c.storeL, c.loadE) },
	0x6c: func(c *CPU) { ld(c, c.storeL, c.loadH) },
	0x6d: func(c *CPU) { ld(c, c.storeL, c.loadL) },
	0x6e: func(c *CPU) { ld(c, c.storeL, c.loadIndHL) },
	0x6f: func(c *CPU) { ld(c, c.storeL, c.loadA) },
	0x70: func(c *CPU) { ld(c, c.storeIndHL, c.loadB) },
	0x71: func(c *CPU) { ld(c, c.storeIndHL, c.loadC) },
	0x72: func(c *CPU) { ld(c, c.storeIndHL, c.loadD) },
	0x73: func(c *CPU) { ld(c, c.storeIndHL, c.loadE) },
	0x74: func(c *CPU) { ld(c, c.storeIndHL, c.loadH) },
	0x75: func(c *CPU) { ld(c, c.storeIndHL, c.loadL) },
	0x76: func(c *CPU) { halt(c) },
	0x77: func(c *CPU) { ld(c, c.storeIndHL, c.loadA) },
	0x78: func(c *CPU) { ld(c, c.storeA, c.loadB) },
	0x79: func(c *CPU) { ld(c, c.storeA, c.loadC) },
	0x7a: func(c *CPU) { ld(c, c.storeA, c.loadD) },
	0x7b: func(c *CPU) { ld(c, c.storeA, c.loadE) },
	0x7c: func(c *CPU) { ld(c, c.storeA, c.loadH) },
	0x7d: func(c *CPU) { ld(c, c.storeA, c.loadL) },
	0x7e: func(c *CPU) { ld(c, c.storeA, c.loadIndHL) },
	0x7f: func(c *CPU) { ld(c, c.storeA, c.loadA) },
	0x80: func(c *CPU) { add(c, c.loadA, c.loadB, false) },
	0x81: func(c *CPU) { add(c, c.loadA, c.loadC, false) },
	0x82: func(c *CPU) { add(c, c.loadA, c.loadD, false) },
	0x83: func(c *CPU) { add(c, c.loadA, c.loadE, false) },
	0x84: func(c *CPU) { add(c, c.loadA, c.loadH, false) },
	0x85: func(c *CPU) { add(c, c.loadA, c.loadL, false) },
	0x86: func(c *CPU) { add(c, c.loadA, c.loadIndHL, false) },
	0x87: func(c *CPU) { add(c, c.loadA, c.loadA, false) },
	0x88: func(c *CPU) { add(c, c.loadA, c.loadB, true) },
	0x89: func(c *CPU) { add(c, c.loadA, c.loadC, true) },
	0x8a: func(c *CPU) { add(c, c.loadA, c.loadD, true) },
	0x8b: func(c *CPU) { add(c, c.loadA, c.loadE, true) },
	0x8c: func(c *CPU) { add(c, c.loadA, c.loadH, true) },
	0x8d: func(c *CPU) { add(c, c.loadA, c.loadL, true) },
	0x8e: func(c *CPU) { add(c, c.loadA, c.loadIndHL, true) },
	0x8f: func(c *CPU) { add(c, c.loadA, c.loadA, true) },
	0x90: func(c *CPU) { sub(c, c.loadB, false) },
	0x91: func(c *CPU) { sub(c, c.loadC, false) },
	0x92: func(c *CPU) { sub(c, c.loadD, false) },
	0x93: func(c *CPU) { sub(c, c.loadE, false) },
	0x94: func(c *CPU) { sub(c, c.loadH, false) },
	0x95: func(c *CPU) { sub(c, c.loadL, false) },
	0x96: func(c *CPU) { sub(c, c.loadIndHL, false) },
	0x97: func(c *CPU) { sub(c, c.loadA, false) },
	0x98: func(c *CPU) { sub(c, c.loadB, true) },
	0x99: func(c *CPU) { sub(c, c.loadC, true) },
	0x9a: func(c *CPU) { sub(c, c.loadD, true) },
	0x9b: func(c *CPU) { sub(c, c.loadE, true) },
	0x9c: func(c *CPU) { sub(c, c.loadH, true) },
	0x9d: func(c *CPU) { sub(c, c.loadL, true) },
	0x9e: func(c *CPU) { sub(c, c.loadIndHL, true) },
	0x9f: func(c *CPU) { sub(c, c.loadA, true) },
	0xa0: func(c *CPU) { and(c, c.loadB) },
	0xa1: func(c *CPU) { and(c, c.loadC) },
	0xa2: func(c *CPU) { and(c, c.loadD) },
	0xa3: func(c *CPU) { and(c, c.loadE) },
	0xa4: func(c *CPU) { and(c, c.loadH) },
	0xa5: func(c *CPU) { and(c, c.loadL) },
	0xa6: func(c *CPU) { and(c, c.loadIndHL) },
	0xa7: func(c *CPU) { and(c, c.loadA) },
	0xa8: func(c *CPU) { xor(c, c.loadB) },
	0xa9: func(c *CPU) { xor(c, c.loadC) },
	0xaa: func(c *CPU) { xor(c, c.loadD) },
	0xab: func(c *CPU) { xor(c, c.loadE) },
	0xac: func(c *CPU) { xor(c, c.loadH) },
	0xad: func(c *CPU) { xor(c, c.loadL) },
	0xae: func(c *CPU) { xor(c, c.loadIndHL) },
	0xaf: func(c *CPU) { xor(c, c.loadA) },
	0xb0: func(c *CPU) { or(c, c.loadB) },
	0xb1: func(c *CPU) { or(c, c.loadC) },
	0xb2: func(c *CPU) { or(c, c.loadD) },
	0xb3: func(c *CPU) { or(c, c.loadE) },
	0xb4: func(c *CPU) { or(c, c.loadH) },
	0xb5: func(c *CPU) { or(c, c.loadL) },
	0xb6: func(c *CPU) { or(c, c.loadIndHL) },
	0xb7: func(c *CPU) { or(c, c.loadA) },
	0xb8: func(c *CPU) { cp(c, c.loadB) },
	0xb9: func(c *CPU) { cp(c, c.loadC) },
	0xba: func(c *CPU) { cp(c, c.loadD) },
	0xbb: func(c *CPU) { cp(c, c.loadE) },
	0xbc: func(c *CPU) { cp(c, c.loadH) },
	0xbd: func(c *CPU) { cp(c, c.loadL) },
	0xbe: func(c *CPU) { cp(c, c.loadIndHL) },
	0xbf: func(c *CPU) { cp(c, c.loadA) },
	0xc0: func(c *CPU) { ret(c, FlagZ, false) },
	0xc1: func(c *CPU) { pop(c, c.storeBC) },
	0xc2: func(c *CPU) { jp(c, FlagZ, false, c.loadImm16) },
	0xc3: func(c *CPU) { jpa(c, c.loadImm16) },
	0xc4: func(c *CPU) { call(c, FlagZ, false, c.loadImm16) },
	0xc5: func(c *CPU) { push(c, c.loadBC) },
	0xc6: func(c *CPU) { add(c, c.loadImm, c.loadA, false) },
	0xc7: func(c *CPU) { rst(c, 0) },
	0xc8: func(c *CPU) { ret(c, FlagZ, true) },
	0xc9: func(c *CPU) { reta(c) },
	0xca: func(c *CPU) { jp(c, FlagZ, true, c.loadImm16) },
	0xcb: func(c *CPU) { cb(c) },
	0xcc: func(c *CPU) { call(c, FlagZ, true, c.loadImm16) },
	0xcd: func(c *CPU) { calla(c, c.loadImm16) },
	0xce: func(c *CPU) { add(c, c.loadImm, c.loadA, true) },
	0xcf: func(c *CPU) { rst(c, 1) },
	0xd0: func(c *CPU) { ret(c, FlagC, false) },
	0xd1: func(c *CPU) { pop(c, c.storeDE) },
	0xd2: func(c *CPU) { jp(c, FlagC, false, c.loadImm16) },
	0xd3: func(c *CPU) { out(c) },
	0xd4: func(c *CPU) { call(c, FlagC, false, c.loadImm16) },
	0xd5: func(c *CPU) { push(c, c.loadDE) },
	0xd6: func(c *CPU) { sub(c, c.loadImm, false) },
	0xd7: func(c *CPU) { rst(c, 2) },
	0xd8: func(c *CPU) { ret(c, FlagC, true) },
	0xd9: func(c *CPU) { exx(c) },
	0xda: func(c *CPU) { jp(c, FlagC, true, c.loadImm16) },
	0xdb: func(c *CPU) { in(c, c.storeA, c.loadImm) },
	0xdc: func(c *CPU) { call(c, FlagC, true, c.loadImm16) },
	0xdd: func(c *CPU) { c.skip = true },
	0xde: func(c *CPU) { sub(c, c.loadImm, true) },
	0xdf: func(c *CPU) { rst(c, 3) },
	0xe0: func(c *CPU) { ret(c, FlagV, false) },
	0xe1: func(c *CPU) { pop(c, c.storeHL) },
	0xe2: func(c *CPU) { jp(c, FlagV, false, c.loadImm16) },
	0xe3: func(c *CPU) { ex(c, c.load16IndSP, c.store16IndSP, c.loadHL, c.storeHL) },
	0xe4: func(c *CPU) { call(c, FlagV, false, c.loadImm16) },
	0xe5: func(c *CPU) { push(c, c.loadHL) },
	0xe6: func(c *CPU) { and(c, c.loadImm) },
	0xe7: func(c *CPU) { rst(c, 4) },
	0xe8: func(c *CPU) { ret(c, FlagV, true) },
	0xe9: func(c *CPU) { jpa(c, c.loadHL) },
	0xea: func(c *CPU) { jp(c, FlagV, true, c.loadImm16) },
	0xeb: func(c *CPU) { ex(c, c.loadDE, c.storeDE, c.loadHL, c.storeHL) },
	0xec: func(c *CPU) { call(c, FlagV, true, c.loadImm16) },
	0xed: func(c *CPU) { ed(c) },
	0xee: func(c *CPU) { xor(c, c.loadImm) },
	0xef: func(c *CPU) { rst(c, 5) },
	0xf0: func(c *CPU) { ret(c, FlagS, false) },
	0xf1: func(c *CPU) { pop(c, c.storeAF) },
	0xf2: func(c *CPU) { jp(c, FlagS, false, c.loadImm16) },
	0xf3: func(c *CPU) { di() },
	0xf4: func(c *CPU) { call(c, FlagS, false, c.loadImm16) },
	0xf5: func(c *CPU) { push(c, c.loadAF) },
	0xf6: func(c *CPU) { or(c, c.loadImm) },
	0xf7: func(c *CPU) { rst(c, 6) },
	0xf8: func(c *CPU) { ret(c, FlagS, true) },
	0xf9: func(c *CPU) { ld16(c, c.storeSP, c.loadHL) },
	0xfa: func(c *CPU) { jp(c, FlagS, true, c.loadImm16) },
	0xfb: func(c *CPU) { ei() },
	0xfc: func(c *CPU) { call(c, FlagS, true, c.loadImm16) },
	0xfd: func(c *CPU) { c.skip = true },
	0xfe: func(c *CPU) { cp(c, c.loadImm) },
	0xff: func(c *CPU) { rst(c, 7) },
}
var opsCB = map[uint8]func(c *CPU){
	0x00: func(c *CPU) { rotl(c, c.storeB, c.loadB) },
	0x01: func(c *CPU) { rotl(c, c.storeC, c.loadC) },
	0x02: func(c *CPU) { rotl(c, c.storeD, c.loadD) },
	0x03: func(c *CPU) { rotl(c, c.storeE, c.loadE) },
	0x04: func(c *CPU) { rotl(c, c.storeH, c.loadH) },
	0x05: func(c *CPU) { rotl(c, c.storeL, c.loadL) },
	0x06: func(c *CPU) { rotl(c, c.storeIndHL, c.loadIndHL) },
	0x07: func(c *CPU) { rotl(c, c.storeA, c.loadA) },
	0x08: func(c *CPU) { rotr(c, c.storeB, c.loadB) },
	0x09: func(c *CPU) { rotr(c, c.storeC, c.loadC) },
	0x0a: func(c *CPU) { rotr(c, c.storeD, c.loadD) },
	0x0b: func(c *CPU) { rotr(c, c.storeE, c.loadE) },
	0x0c: func(c *CPU) { rotr(c, c.storeH, c.loadH) },
	0x0d: func(c *CPU) { rotr(c, c.storeL, c.loadL) },
	0x0e: func(c *CPU) { rotr(c, c.storeIndHL, c.loadIndHL) },
	0x0f: func(c *CPU) { rotr(c, c.storeA, c.loadA) },
	0x10: func(c *CPU) { shiftl(c, c.storeB, c.loadB, true) },
	0x11: func(c *CPU) { shiftl(c, c.storeC, c.loadC, true) },
	0x12: func(c *CPU) { shiftl(c, c.storeD, c.loadD, true) },
	0x13: func(c *CPU) { shiftl(c, c.storeE, c.loadE, true) },
	0x14: func(c *CPU) { shiftl(c, c.storeH, c.loadH, true) },
	0x15: func(c *CPU) { shiftl(c, c.storeL, c.loadL, true) },
	0x16: func(c *CPU) { shiftl(c, c.storeIndHL, c.loadIndHL, true) },
	0x17: func(c *CPU) { shiftl(c, c.storeA, c.loadA, true) },
	0x18: func(c *CPU) { shiftr(c, c.storeB, c.loadB, true) },
	0x19: func(c *CPU) { shiftr(c, c.storeC, c.loadC, true) },
	0x1a: func(c *CPU) { shiftr(c, c.storeD, c.loadD, true) },
	0x1b: func(c *CPU) { shiftr(c, c.storeE, c.loadE, true) },
	0x1c: func(c *CPU) { shiftr(c, c.storeH, c.loadH, true) },
	0x1d: func(c *CPU) { shiftr(c, c.storeL, c.loadL, true) },
	0x1e: func(c *CPU) { shiftr(c, c.storeIndHL, c.loadIndHL, true) },
	0x1f: func(c *CPU) { shiftr(c, c.storeA, c.loadA, true) },
	0x20: func(c *CPU) { shiftl(c, c.storeB, c.loadB, false) },
	0x21: func(c *CPU) { shiftl(c, c.storeC, c.loadC, false) },
	0x22: func(c *CPU) { shiftl(c, c.storeD, c.loadD, false) },
	0x23: func(c *CPU) { shiftl(c, c.storeE, c.loadE, false) },
	0x24: func(c *CPU) { shiftl(c, c.storeH, c.loadH, false) },
	0x25: func(c *CPU) { shiftl(c, c.storeL, c.loadL, false) },
	0x26: func(c *CPU) { shiftl(c, c.storeIndHL, c.loadIndHL, false) },
	0x27: func(c *CPU) { shiftl(c, c.storeA, c.loadA, false) },
	0x28: func(c *CPU) { sra(c, c.storeB, c.loadB) },
	0x29: func(c *CPU) { sra(c, c.storeC, c.loadC) },
	0x2a: func(c *CPU) { sra(c, c.storeD, c.loadD) },
	0x2b: func(c *CPU) { sra(c, c.storeE, c.loadE) },
	0x2c: func(c *CPU) { sra(c, c.storeH, c.loadH) },
	0x2d: func(c *CPU) { sra(c, c.storeL, c.loadL) },
	0x2e: func(c *CPU) { sra(c, c.storeIndHL, c.loadIndHL) },
	0x2f: func(c *CPU) { sra(c, c.storeA, c.loadA) },
	0x30: func(c *CPU) { sll(c, c.storeB, c.loadB) },
	0x31: func(c *CPU) { sll(c, c.storeC, c.loadC) },
	0x32: func(c *CPU) { sll(c, c.storeD, c.loadD) },
	0x33: func(c *CPU) { sll(c, c.storeE, c.loadE) },
	0x34: func(c *CPU) { sll(c, c.storeH, c.loadH) },
	0x35: func(c *CPU) { sll(c, c.storeL, c.loadL) },
	0x36: func(c *CPU) { sll(c, c.storeIndHL, c.loadIndHL) },
	0x37: func(c *CPU) { sll(c, c.storeA, c.loadA) },
	0x38: func(c *CPU) { shiftr(c, c.storeB, c.loadB, false) },
	0x39: func(c *CPU) { shiftr(c, c.storeC, c.loadC, false) },
	0x3a: func(c *CPU) { shiftr(c, c.storeD, c.loadD, false) },
	0x3b: func(c *CPU) { shiftr(c, c.storeE, c.loadE, false) },
	0x3c: func(c *CPU) { shiftr(c, c.storeH, c.loadH, false) },
	0x3d: func(c *CPU) { shiftr(c, c.storeL, c.loadL, false) },
	0x3e: func(c *CPU) { shiftr(c, c.storeIndHL, c.loadIndHL, false) },
	0x3f: func(c *CPU) { shiftr(c, c.storeA, c.loadA, false) },
	0x40: func(c *CPU) { bit(c, 0, c.loadB) },
	0x41: func(c *CPU) { bit(c, 0, c.loadC) },
	0x42: func(c *CPU) { bit(c, 0, c.loadD) },
	0x43: func(c *CPU) { bit(c, 0, c.loadE) },
	0x44: func(c *CPU) { bit(c, 0, c.loadH) },
	0x45: func(c *CPU) { bit(c, 0, c.loadL) },
	0x46: func(c *CPU) { bit(c, 0, c.loadIndHL) },
	0x47: func(c *CPU) { bit(c, 0, c.loadA) },
	0x48: func(c *CPU) { bit(c, 1, c.loadB) },
	0x49: func(c *CPU) { bit(c, 1, c.loadC) },
	0x4a: func(c *CPU) { bit(c, 1, c.loadD) },
	0x4b: func(c *CPU) { bit(c, 1, c.loadE) },
	0x4c: func(c *CPU) { bit(c, 1, c.loadH) },
	0x4d: func(c *CPU) { bit(c, 1, c.loadL) },
	0x4e: func(c *CPU) { bit(c, 1, c.loadIndHL) },
	0x4f: func(c *CPU) { bit(c, 1, c.loadA) },
	0x50: func(c *CPU) { bit(c, 2, c.loadB) },
	0x51: func(c *CPU) { bit(c, 2, c.loadC) },
	0x52: func(c *CPU) { bit(c, 2, c.loadD) },
	0x53: func(c *CPU) { bit(c, 2, c.loadE) },
	0x54: func(c *CPU) { bit(c, 2, c.loadH) },
	0x55: func(c *CPU) { bit(c, 2, c.loadL) },
	0x56: func(c *CPU) { bit(c, 2, c.loadIndHL) },
	0x57: func(c *CPU) { bit(c, 2, c.loadA) },
	0x58: func(c *CPU) { bit(c, 3, c.loadB) },
	0x59: func(c *CPU) { bit(c, 3, c.loadC) },
	0x5a: func(c *CPU) { bit(c, 3, c.loadD) },
	0x5b: func(c *CPU) { bit(c, 3, c.loadE) },
	0x5c: func(c *CPU) { bit(c, 3, c.loadH) },
	0x5d: func(c *CPU) { bit(c, 3, c.loadL) },
	0x5e: func(c *CPU) { bit(c, 3, c.loadIndHL) },
	0x5f: func(c *CPU) { bit(c, 3, c.loadA) },
	0x60: func(c *CPU) { bit(c, 4, c.loadB) },
	0x61: func(c *CPU) { bit(c, 4, c.loadC) },
	0x62: func(c *CPU) { bit(c, 4, c.loadD) },
	0x63: func(c *CPU) { bit(c, 4, c.loadE) },
	0x64: func(c *CPU) { bit(c, 4, c.loadH) },
	0x65: func(c *CPU) { bit(c, 4, c.loadL) },
	0x66: func(c *CPU) { bit(c, 4, c.loadIndHL) },
	0x67: func(c *CPU) { bit(c, 4, c.loadA) },
	0x68: func(c *CPU) { bit(c, 5, c.loadB) },
	0x69: func(c *CPU) { bit(c, 5, c.loadC) },
	0x6a: func(c *CPU) { bit(c, 5, c.loadD) },
	0x6b: func(c *CPU) { bit(c, 5, c.loadE) },
	0x6c: func(c *CPU) { bit(c, 5, c.loadH) },
	0x6d: func(c *CPU) { bit(c, 5, c.loadL) },
	0x6e: func(c *CPU) { bit(c, 5, c.loadIndHL) },
	0x6f: func(c *CPU) { bit(c, 5, c.loadA) },
	0x70: func(c *CPU) { bit(c, 6, c.loadB) },
	0x71: func(c *CPU) { bit(c, 6, c.loadC) },
	0x72: func(c *CPU) { bit(c, 6, c.loadD) },
	0x73: func(c *CPU) { bit(c, 6, c.loadE) },
	0x74: func(c *CPU) { bit(c, 6, c.loadH) },
	0x75: func(c *CPU) { bit(c, 6, c.loadL) },
	0x76: func(c *CPU) { bit(c, 6, c.loadIndHL) },
	0x77: func(c *CPU) { bit(c, 6, c.loadA) },
	0x78: func(c *CPU) { bit(c, 7, c.loadB) },
	0x79: func(c *CPU) { bit(c, 7, c.loadC) },
	0x7a: func(c *CPU) { bit(c, 7, c.loadD) },
	0x7b: func(c *CPU) { bit(c, 7, c.loadE) },
	0x7c: func(c *CPU) { bit(c, 7, c.loadH) },
	0x7d: func(c *CPU) { bit(c, 7, c.loadL) },
	0x7e: func(c *CPU) { bit(c, 7, c.loadIndHL) },
	0x7f: func(c *CPU) { bit(c, 7, c.loadA) },
	0x80: func(c *CPU) { res(c, 0, c.storeB, c.loadB) },
	0x81: func(c *CPU) { res(c, 0, c.storeC, c.loadC) },
	0x82: func(c *CPU) { res(c, 0, c.storeD, c.loadD) },
	0x83: func(c *CPU) { res(c, 0, c.storeE, c.loadE) },
	0x84: func(c *CPU) { res(c, 0, c.storeH, c.loadH) },
	0x85: func(c *CPU) { res(c, 0, c.storeL, c.loadL) },
	0x86: func(c *CPU) { res(c, 0, c.storeIndHL, c.loadIndHL) },
	0x87: func(c *CPU) { res(c, 0, c.storeA, c.loadA) },
	0x88: func(c *CPU) { res(c, 1, c.storeB, c.loadB) },
	0x89: func(c *CPU) { res(c, 1, c.storeC, c.loadC) },
	0x8a: func(c *CPU) { res(c, 1, c.storeD, c.loadD) },
	0x8b: func(c *CPU) { res(c, 1, c.storeE, c.loadE) },
	0x8c: func(c *CPU) { res(c, 1, c.storeH, c.loadH) },
	0x8d: func(c *CPU) { res(c, 1, c.storeL, c.loadL) },
	0x8e: func(c *CPU) { res(c, 1, c.storeIndHL, c.loadIndHL) },
	0x8f: func(c *CPU) { res(c, 1, c.storeA, c.loadA) },
	0x90: func(c *CPU) { res(c, 2, c.storeB, c.loadB) },
	0x91: func(c *CPU) { res(c, 2, c.storeC, c.loadC) },
	0x92: func(c *CPU) { res(c, 2, c.storeD, c.loadD) },
	0x93: func(c *CPU) { res(c, 2, c.storeE, c.loadE) },
	0x94: func(c *CPU) { res(c, 2, c.storeH, c.loadH) },
	0x95: func(c *CPU) { res(c, 2, c.storeL, c.loadL) },
	0x96: func(c *CPU) { res(c, 2, c.storeIndHL, c.loadIndHL) },
	0x97: func(c *CPU) { res(c, 2, c.storeA, c.loadA) },
	0x98: func(c *CPU) { res(c, 3, c.storeB, c.loadB) },
	0x99: func(c *CPU) { res(c, 3, c.storeC, c.loadC) },
	0x9a: func(c *CPU) { res(c, 3, c.storeD, c.loadD) },
	0x9b: func(c *CPU) { res(c, 3, c.storeE, c.loadE) },
	0x9c: func(c *CPU) { res(c, 3, c.storeH, c.loadH) },
	0x9d: func(c *CPU) { res(c, 3, c.storeL, c.loadL) },
	0x9e: func(c *CPU) { res(c, 3, c.storeIndHL, c.loadIndHL) },
	0x9f: func(c *CPU) { res(c, 3, c.storeA, c.loadA) },
	0xa0: func(c *CPU) { res(c, 4, c.storeB, c.loadB) },
	0xa1: func(c *CPU) { res(c, 4, c.storeC, c.loadC) },
	0xa2: func(c *CPU) { res(c, 4, c.storeD, c.loadD) },
	0xa3: func(c *CPU) { res(c, 4, c.storeE, c.loadE) },
	0xa4: func(c *CPU) { res(c, 4, c.storeH, c.loadH) },
	0xa5: func(c *CPU) { res(c, 4, c.storeL, c.loadL) },
	0xa6: func(c *CPU) { res(c, 4, c.storeIndHL, c.loadIndHL) },
	0xa7: func(c *CPU) { res(c, 4, c.storeA, c.loadA) },
	0xa8: func(c *CPU) { res(c, 5, c.storeB, c.loadB) },
	0xa9: func(c *CPU) { res(c, 5, c.storeC, c.loadC) },
	0xaa: func(c *CPU) { res(c, 5, c.storeD, c.loadD) },
	0xab: func(c *CPU) { res(c, 5, c.storeE, c.loadE) },
	0xac: func(c *CPU) { res(c, 5, c.storeH, c.loadH) },
	0xad: func(c *CPU) { res(c, 5, c.storeL, c.loadL) },
	0xae: func(c *CPU) { res(c, 5, c.storeIndHL, c.loadIndHL) },
	0xaf: func(c *CPU) { res(c, 5, c.storeA, c.loadA) },
	0xb0: func(c *CPU) { res(c, 6, c.storeB, c.loadB) },
	0xb1: func(c *CPU) { res(c, 6, c.storeC, c.loadC) },
	0xb2: func(c *CPU) { res(c, 6, c.storeD, c.loadD) },
	0xb3: func(c *CPU) { res(c, 6, c.storeE, c.loadE) },
	0xb4: func(c *CPU) { res(c, 6, c.storeH, c.loadH) },
	0xb5: func(c *CPU) { res(c, 6, c.storeL, c.loadL) },
	0xb6: func(c *CPU) { res(c, 6, c.storeIndHL, c.loadIndHL) },
	0xb7: func(c *CPU) { res(c, 6, c.storeA, c.loadA) },
	0xb8: func(c *CPU) { res(c, 7, c.storeB, c.loadB) },
	0xb9: func(c *CPU) { res(c, 7, c.storeC, c.loadC) },
	0xba: func(c *CPU) { res(c, 7, c.storeD, c.loadD) },
	0xbb: func(c *CPU) { res(c, 7, c.storeE, c.loadE) },
	0xbc: func(c *CPU) { res(c, 7, c.storeH, c.loadH) },
	0xbd: func(c *CPU) { res(c, 7, c.storeL, c.loadL) },
	0xbe: func(c *CPU) { res(c, 7, c.storeIndHL, c.loadIndHL) },
	0xbf: func(c *CPU) { res(c, 7, c.storeA, c.loadA) },
	0xc0: func(c *CPU) { set(c, 0, c.storeB, c.loadB) },
	0xc1: func(c *CPU) { set(c, 0, c.storeC, c.loadC) },
	0xc2: func(c *CPU) { set(c, 0, c.storeD, c.loadD) },
	0xc3: func(c *CPU) { set(c, 0, c.storeE, c.loadE) },
	0xc4: func(c *CPU) { set(c, 0, c.storeH, c.loadH) },
	0xc5: func(c *CPU) { set(c, 0, c.storeL, c.loadL) },
	0xc6: func(c *CPU) { set(c, 0, c.storeIndHL, c.loadIndHL) },
	0xc7: func(c *CPU) { set(c, 0, c.storeA, c.loadA) },
	0xc8: func(c *CPU) { set(c, 1, c.storeB, c.loadB) },
	0xc9: func(c *CPU) { set(c, 1, c.storeC, c.loadC) },
	0xca: func(c *CPU) { set(c, 1, c.storeD, c.loadD) },
	0xcb: func(c *CPU) { set(c, 1, c.storeE, c.loadE) },
	0xcc: func(c *CPU) { set(c, 1, c.storeH, c.loadH) },
	0xcd: func(c *CPU) { set(c, 1, c.storeL, c.loadL) },
	0xce: func(c *CPU) { set(c, 1, c.storeIndHL, c.loadIndHL) },
	0xcf: func(c *CPU) { set(c, 1, c.storeA, c.loadA) },
	0xd0: func(c *CPU) { set(c, 2, c.storeB, c.loadB) },
	0xd1: func(c *CPU) { set(c, 2, c.storeC, c.loadC) },
	0xd2: func(c *CPU) { set(c, 2, c.storeD, c.loadD) },
	0xd3: func(c *CPU) { set(c, 2, c.storeE, c.loadE) },
	0xd4: func(c *CPU) { set(c, 2, c.storeH, c.loadH) },
	0xd5: func(c *CPU) { set(c, 2, c.storeL, c.loadL) },
	0xd6: func(c *CPU) { set(c, 2, c.storeIndHL, c.loadIndHL) },
	0xd7: func(c *CPU) { set(c, 2, c.storeA, c.loadA) },
	0xd8: func(c *CPU) { set(c, 3, c.storeB, c.loadB) },
	0xd9: func(c *CPU) { set(c, 3, c.storeC, c.loadC) },
	0xda: func(c *CPU) { set(c, 3, c.storeD, c.loadD) },
	0xdb: func(c *CPU) { set(c, 3, c.storeE, c.loadE) },
	0xdc: func(c *CPU) { set(c, 3, c.storeH, c.loadH) },
	0xdd: func(c *CPU) { set(c, 3, c.storeL, c.loadL) },
	0xde: func(c *CPU) { set(c, 3, c.storeIndHL, c.loadIndHL) },
	0xdf: func(c *CPU) { set(c, 3, c.storeA, c.loadA) },
	0xe0: func(c *CPU) { set(c, 4, c.storeB, c.loadB) },
	0xe1: func(c *CPU) { set(c, 4, c.storeC, c.loadC) },
	0xe2: func(c *CPU) { set(c, 4, c.storeD, c.loadD) },
	0xe3: func(c *CPU) { set(c, 4, c.storeE, c.loadE) },
	0xe4: func(c *CPU) { set(c, 4, c.storeH, c.loadH) },
	0xe5: func(c *CPU) { set(c, 4, c.storeL, c.loadL) },
	0xe6: func(c *CPU) { set(c, 4, c.storeIndHL, c.loadIndHL) },
	0xe7: func(c *CPU) { set(c, 4, c.storeA, c.loadA) },
	0xe8: func(c *CPU) { set(c, 5, c.storeB, c.loadB) },
	0xe9: func(c *CPU) { set(c, 5, c.storeC, c.loadC) },
	0xea: func(c *CPU) { set(c, 5, c.storeD, c.loadD) },
	0xeb: func(c *CPU) { set(c, 5, c.storeE, c.loadE) },
	0xec: func(c *CPU) { set(c, 5, c.storeH, c.loadH) },
	0xed: func(c *CPU) { set(c, 5, c.storeL, c.loadL) },
	0xee: func(c *CPU) { set(c, 5, c.storeIndHL, c.loadIndHL) },
	0xef: func(c *CPU) { set(c, 5, c.storeA, c.loadA) },
	0xf0: func(c *CPU) { set(c, 6, c.storeB, c.loadB) },
	0xf1: func(c *CPU) { set(c, 6, c.storeC, c.loadC) },
	0xf2: func(c *CPU) { set(c, 6, c.storeD, c.loadD) },
	0xf3: func(c *CPU) { set(c, 6, c.storeE, c.loadE) },
	0xf4: func(c *CPU) { set(c, 6, c.storeH, c.loadH) },
	0xf5: func(c *CPU) { set(c, 6, c.storeL, c.loadL) },
	0xf6: func(c *CPU) { set(c, 6, c.storeIndHL, c.loadIndHL) },
	0xf7: func(c *CPU) { set(c, 6, c.storeA, c.loadA) },
	0xf8: func(c *CPU) { set(c, 7, c.storeB, c.loadB) },
	0xf9: func(c *CPU) { set(c, 7, c.storeC, c.loadC) },
	0xfa: func(c *CPU) { set(c, 7, c.storeD, c.loadD) },
	0xfb: func(c *CPU) { set(c, 7, c.storeE, c.loadE) },
	0xfc: func(c *CPU) { set(c, 7, c.storeH, c.loadH) },
	0xfd: func(c *CPU) { set(c, 7, c.storeL, c.loadL) },
	0xfe: func(c *CPU) { set(c, 7, c.storeIndHL, c.loadIndHL) },
	0xff: func(c *CPU) { set(c, 7, c.storeA, c.loadA) },
}
var opsED = map[uint8]func(c *CPU){
	0x00: func(c *CPU) { invalid() },
	0x01: func(c *CPU) { invalid() },
	0x02: func(c *CPU) { invalid() },
	0x03: func(c *CPU) { invalid() },
	0x04: func(c *CPU) { invalid() },
	0x05: func(c *CPU) { invalid() },
	0x06: func(c *CPU) { invalid() },
	0x07: func(c *CPU) { invalid() },
	0x08: func(c *CPU) { invalid() },
	0x09: func(c *CPU) { invalid() },
	0x0a: func(c *CPU) { invalid() },
	0x0b: func(c *CPU) { invalid() },
	0x0c: func(c *CPU) { invalid() },
	0x0d: func(c *CPU) { invalid() },
	0x0e: func(c *CPU) { invalid() },
	0x0f: func(c *CPU) { invalid() },
	0x10: func(c *CPU) { invalid() },
	0x11: func(c *CPU) { invalid() },
	0x12: func(c *CPU) { invalid() },
	0x13: func(c *CPU) { invalid() },
	0x14: func(c *CPU) { invalid() },
	0x15: func(c *CPU) { invalid() },
	0x16: func(c *CPU) { invalid() },
	0x17: func(c *CPU) { invalid() },
	0x18: func(c *CPU) { invalid() },
	0x19: func(c *CPU) { invalid() },
	0x1a: func(c *CPU) { invalid() },
	0x1b: func(c *CPU) { invalid() },
	0x1c: func(c *CPU) { invalid() },
	0x1d: func(c *CPU) { invalid() },
	0x1e: func(c *CPU) { invalid() },
	0x1f: func(c *CPU) { invalid() },
	0x20: func(c *CPU) { invalid() },
	0x21: func(c *CPU) { invalid() },
	0x22: func(c *CPU) { invalid() },
	0x23: func(c *CPU) { invalid() },
	0x24: func(c *CPU) { invalid() },
	0x25: func(c *CPU) { invalid() },
	0x26: func(c *CPU) { invalid() },
	0x27: func(c *CPU) { invalid() },
	0x28: func(c *CPU) { invalid() },
	0x29: func(c *CPU) { invalid() },
	0x2a: func(c *CPU) { invalid() },
	0x2b: func(c *CPU) { invalid() },
	0x2c: func(c *CPU) { invalid() },
	0x2d: func(c *CPU) { invalid() },
	0x2e: func(c *CPU) { invalid() },
	0x2f: func(c *CPU) { invalid() },
	0x30: func(c *CPU) { invalid() },
	0x31: func(c *CPU) { invalid() },
	0x32: func(c *CPU) { invalid() },
	0x33: func(c *CPU) { invalid() },
	0x34: func(c *CPU) { invalid() },
	0x35: func(c *CPU) { invalid() },
	0x36: func(c *CPU) { invalid() },
	0x37: func(c *CPU) { invalid() },
	0x38: func(c *CPU) { invalid() },
	0x39: func(c *CPU) { invalid() },
	0x3a: func(c *CPU) { invalid() },
	0x3b: func(c *CPU) { invalid() },
	0x3c: func(c *CPU) { invalid() },
	0x3d: func(c *CPU) { invalid() },
	0x3e: func(c *CPU) { invalid() },
	0x3f: func(c *CPU) { invalid() },
	0x40: func(c *CPU) { c.skip = true },
	0x41: func(c *CPU) { c.skip = true },
	0x42: func(c *CPU) { sub16(c, c.storeHL, c.loadHL, c.loadBC, true) },
	0x43: func(c *CPU) { c.skip = true },
	0x44: func(c *CPU) { c.skip = true },
	0x45: func(c *CPU) { c.skip = true },
	0x46: func(c *CPU) { c.skip = true },
	0x47: func(c *CPU) { c.skip = true },
	0x48: func(c *CPU) { c.skip = true },
	0x49: func(c *CPU) { c.skip = true },
	0x4a: func(c *CPU) { add16(c, c.storeHL, c.loadHL, c.loadBC, true) },
	0x4b: func(c *CPU) { c.skip = true },
	0x4c: func(c *CPU) { c.skip = true },
	0x4d: func(c *CPU) { c.skip = true },
	0x4e: func(c *CPU) { c.skip = true },
	0x4f: func(c *CPU) { c.skip = true },
	0x50: func(c *CPU) { c.skip = true },
	0x51: func(c *CPU) { c.skip = true },
	0x52: func(c *CPU) { sub16(c, c.storeHL, c.loadHL, c.loadDE, true) },
	0x53: func(c *CPU) { c.skip = true },
	0x54: func(c *CPU) { c.skip = true },
	0x55: func(c *CPU) { c.skip = true },
	0x56: func(c *CPU) { c.skip = true },
	0x57: func(c *CPU) { c.skip = true },
	0x58: func(c *CPU) { c.skip = true },
	0x59: func(c *CPU) { c.skip = true },
	0x5a: func(c *CPU) { add16(c, c.storeHL, c.loadHL, c.loadDE, true) },
	0x5b: func(c *CPU) { c.skip = true },
	0x5c: func(c *CPU) { c.skip = true },
	0x5d: func(c *CPU) { c.skip = true },
	0x5e: func(c *CPU) { c.skip = true },
	0x5f: func(c *CPU) { c.skip = true },
	0x60: func(c *CPU) { c.skip = true },
	0x61: func(c *CPU) { c.skip = true },
	0x62: func(c *CPU) { sub16(c, c.storeHL, c.loadHL, c.loadHL, true) },
	0x63: func(c *CPU) { c.skip = true },
	0x64: func(c *CPU) { c.skip = true },
	0x65: func(c *CPU) { c.skip = true },
	0x66: func(c *CPU) { c.skip = true },
	0x67: func(c *CPU) { c.skip = true },
	0x68: func(c *CPU) { c.skip = true },
	0x69: func(c *CPU) { c.skip = true },
	0x6a: func(c *CPU) { add16(c, c.storeHL, c.loadHL, c.loadHL, true) },
	0x6b: func(c *CPU) { c.skip = true },
	0x6c: func(c *CPU) { c.skip = true },
	0x6d: func(c *CPU) { c.skip = true },
	0x6e: func(c *CPU) { c.skip = true },
	0x6f: func(c *CPU) { c.skip = true },
	0x70: func(c *CPU) { c.skip = true },
	0x71: func(c *CPU) { c.skip = true },
	0x72: func(c *CPU) { sub16(c, c.storeHL, c.loadHL, c.loadSP, true) },
	0x73: func(c *CPU) { c.skip = true },
	0x74: func(c *CPU) { c.skip = true },
	0x75: func(c *CPU) { c.skip = true },
	0x76: func(c *CPU) { c.skip = true },
	0x77: func(c *CPU) { c.skip = true },
	0x78: func(c *CPU) { c.skip = true },
	0x79: func(c *CPU) { c.skip = true },
	0x7a: func(c *CPU) { add16(c, c.storeHL, c.loadHL, c.loadSP, true) },
	0x7b: func(c *CPU) { c.skip = true },
	0x7c: func(c *CPU) { c.skip = true },
	0x7d: func(c *CPU) { c.skip = true },
	0x7e: func(c *CPU) { c.skip = true },
	0x7f: func(c *CPU) { c.skip = true },
	0x80: func(c *CPU) { c.skip = true },
	0x81: func(c *CPU) { c.skip = true },
	0x82: func(c *CPU) { c.skip = true },
	0x83: func(c *CPU) { c.skip = true },
	0x84: func(c *CPU) { c.skip = true },
	0x85: func(c *CPU) { c.skip = true },
	0x86: func(c *CPU) { c.skip = true },
	0x87: func(c *CPU) { c.skip = true },
	0x88: func(c *CPU) { c.skip = true },
	0x89: func(c *CPU) { c.skip = true },
	0x8a: func(c *CPU) { c.skip = true },
	0x8b: func(c *CPU) { c.skip = true },
	0x8c: func(c *CPU) { c.skip = true },
	0x8d: func(c *CPU) { c.skip = true },
	0x8e: func(c *CPU) { c.skip = true },
	0x8f: func(c *CPU) { c.skip = true },
	0x90: func(c *CPU) { c.skip = true },
	0x91: func(c *CPU) { c.skip = true },
	0x92: func(c *CPU) { c.skip = true },
	0x93: func(c *CPU) { c.skip = true },
	0x94: func(c *CPU) { c.skip = true },
	0x95: func(c *CPU) { c.skip = true },
	0x96: func(c *CPU) { c.skip = true },
	0x97: func(c *CPU) { c.skip = true },
	0x98: func(c *CPU) { c.skip = true },
	0x99: func(c *CPU) { c.skip = true },
	0x9a: func(c *CPU) { c.skip = true },
	0x9b: func(c *CPU) { c.skip = true },
	0x9c: func(c *CPU) { c.skip = true },
	0x9d: func(c *CPU) { c.skip = true },
	0x9e: func(c *CPU) { c.skip = true },
	0x9f: func(c *CPU) { c.skip = true },
	0xa0: func(c *CPU) { c.skip = true },
	0xa1: func(c *CPU) { c.skip = true },
	0xa2: func(c *CPU) { c.skip = true },
	0xa3: func(c *CPU) { c.skip = true },
	0xa4: func(c *CPU) { c.skip = true },
	0xa5: func(c *CPU) { c.skip = true },
	0xa6: func(c *CPU) { c.skip = true },
	0xa7: func(c *CPU) { c.skip = true },
	0xa8: func(c *CPU) { c.skip = true },
	0xa9: func(c *CPU) { c.skip = true },
	0xaa: func(c *CPU) { c.skip = true },
	0xab: func(c *CPU) { c.skip = true },
	0xac: func(c *CPU) { c.skip = true },
	0xad: func(c *CPU) { c.skip = true },
	0xae: func(c *CPU) { c.skip = true },
	0xaf: func(c *CPU) { c.skip = true },
	0xb0: func(c *CPU) { c.skip = true },
	0xb1: func(c *CPU) { c.skip = true },
	0xb2: func(c *CPU) { c.skip = true },
	0xb3: func(c *CPU) { c.skip = true },
	0xb4: func(c *CPU) { c.skip = true },
	0xb5: func(c *CPU) { c.skip = true },
	0xb6: func(c *CPU) { c.skip = true },
	0xb7: func(c *CPU) { c.skip = true },
	0xb8: func(c *CPU) { c.skip = true },
	0xb9: func(c *CPU) { c.skip = true },
	0xba: func(c *CPU) { c.skip = true },
	0xbb: func(c *CPU) { c.skip = true },
	0xbc: func(c *CPU) { c.skip = true },
	0xbd: func(c *CPU) { c.skip = true },
	0xbe: func(c *CPU) { c.skip = true },
	0xbf: func(c *CPU) { c.skip = true },
	0xc0: func(c *CPU) { c.skip = true },
	0xc1: func(c *CPU) { c.skip = true },
	0xc2: func(c *CPU) { c.skip = true },
	0xc3: func(c *CPU) { c.skip = true },
	0xc4: func(c *CPU) { c.skip = true },
	0xc5: func(c *CPU) { c.skip = true },
	0xc6: func(c *CPU) { c.skip = true },
	0xc7: func(c *CPU) { c.skip = true },
	0xc8: func(c *CPU) { c.skip = true },
	0xc9: func(c *CPU) { c.skip = true },
	0xca: func(c *CPU) { c.skip = true },
	0xcb: func(c *CPU) { c.skip = true },
	0xcc: func(c *CPU) { c.skip = true },
	0xcd: func(c *CPU) { c.skip = true },
	0xce: func(c *CPU) { c.skip = true },
	0xcf: func(c *CPU) { c.skip = true },
	0xd0: func(c *CPU) { c.skip = true },
	0xd1: func(c *CPU) { c.skip = true },
	0xd2: func(c *CPU) { c.skip = true },
	0xd3: func(c *CPU) { c.skip = true },
	0xd4: func(c *CPU) { c.skip = true },
	0xd5: func(c *CPU) { c.skip = true },
	0xd6: func(c *CPU) { c.skip = true },
	0xd7: func(c *CPU) { c.skip = true },
	0xd8: func(c *CPU) { c.skip = true },
	0xd9: func(c *CPU) { c.skip = true },
	0xda: func(c *CPU) { c.skip = true },
	0xdb: func(c *CPU) { c.skip = true },
	0xdc: func(c *CPU) { c.skip = true },
	0xdd: func(c *CPU) { c.skip = true },
	0xde: func(c *CPU) { c.skip = true },
	0xdf: func(c *CPU) { c.skip = true },
	0xe0: func(c *CPU) { c.skip = true },
	0xe1: func(c *CPU) { c.skip = true },
	0xe2: func(c *CPU) { c.skip = true },
	0xe3: func(c *CPU) { c.skip = true },
	0xe4: func(c *CPU) { c.skip = true },
	0xe5: func(c *CPU) { c.skip = true },
	0xe6: func(c *CPU) { c.skip = true },
	0xe7: func(c *CPU) { c.skip = true },
	0xe8: func(c *CPU) { c.skip = true },
	0xe9: func(c *CPU) { c.skip = true },
	0xea: func(c *CPU) { c.skip = true },
	0xeb: func(c *CPU) { c.skip = true },
	0xec: func(c *CPU) { c.skip = true },
	0xed: func(c *CPU) { c.skip = true },
	0xee: func(c *CPU) { c.skip = true },
	0xef: func(c *CPU) { c.skip = true },
	0xf0: func(c *CPU) { c.skip = true },
	0xf1: func(c *CPU) { c.skip = true },
	0xf2: func(c *CPU) { c.skip = true },
	0xf3: func(c *CPU) { c.skip = true },
	0xf4: func(c *CPU) { c.skip = true },
	0xf5: func(c *CPU) { c.skip = true },
	0xf6: func(c *CPU) { c.skip = true },
	0xf7: func(c *CPU) { c.skip = true },
	0xf8: func(c *CPU) { c.skip = true },
	0xf9: func(c *CPU) { c.skip = true },
	0xfa: func(c *CPU) { c.skip = true },
	0xfb: func(c *CPU) { c.skip = true },
	0xfc: func(c *CPU) { c.skip = true },
	0xfd: func(c *CPU) { c.skip = true },
	0xfe: func(c *CPU) { c.skip = true },
	0xff: func(c *CPU) { c.skip = true },
}
