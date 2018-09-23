// Code generated by cpu/z80/dasm/gen.go. DO NOT EDIT.

package z80

type harstonTest struct {
	name  string
	op    string
	bytes []uint8
}

var harstonTests = []harstonTest{
	harstonTest{"00", "nop", []uint8{0x00}},
	harstonTest{"01 34 12", "ld   bc,$1234", []uint8{0x01, 0x34, 0x12}},
	harstonTest{"02", "ld   (bc),a", []uint8{0x02}},
	harstonTest{"03", "inc  bc", []uint8{0x03}},
	harstonTest{"04", "inc  b", []uint8{0x04}},
	harstonTest{"05", "dec  b", []uint8{0x05}},
	harstonTest{"06 12", "ld   b,$12", []uint8{0x06, 0x12}},
	harstonTest{"07", "rlca", []uint8{0x07}},
	harstonTest{"08", "ex   af,af'", []uint8{0x08}},
	harstonTest{"09", "add  hl,bc", []uint8{0x09}},
	harstonTest{"0a", "ld   a,(bc)", []uint8{0x0a}},
	harstonTest{"0b", "dec  bc", []uint8{0x0b}},
	harstonTest{"0c", "inc  c", []uint8{0x0c}},
	harstonTest{"0d", "dec  c", []uint8{0x0d}},
	harstonTest{"0e 12", "ld   c,$12", []uint8{0x0e, 0x12}},
	harstonTest{"0f", "rrca", []uint8{0x0f}},
	harstonTest{"10 10", "djnz $0020", []uint8{0x10, 0x10}},
	harstonTest{"11 34 12", "ld   de,$1234", []uint8{0x11, 0x34, 0x12}},
	harstonTest{"12", "ld   (de),a", []uint8{0x12}},
	harstonTest{"13", "inc  de", []uint8{0x13}},
	harstonTest{"14", "inc  d", []uint8{0x14}},
	harstonTest{"15", "dec  d", []uint8{0x15}},
	harstonTest{"16 12", "ld   d,$12", []uint8{0x16, 0x12}},
	harstonTest{"17", "rla", []uint8{0x17}},
	harstonTest{"18 10", "jr   $0020", []uint8{0x18, 0x10}},
	harstonTest{"19", "add  hl,de", []uint8{0x19}},
	harstonTest{"1a", "ld   a,(de)", []uint8{0x1a}},
	harstonTest{"1b", "dec  de", []uint8{0x1b}},
	harstonTest{"1c", "inc  e", []uint8{0x1c}},
	harstonTest{"1d", "dec  e", []uint8{0x1d}},
	harstonTest{"1e 12", "ld   e,$12", []uint8{0x1e, 0x12}},
	harstonTest{"1f", "rra", []uint8{0x1f}},
	harstonTest{"20 10", "jr   nz,$0020", []uint8{0x20, 0x10}},
	harstonTest{"21 34 12", "ld   hl,$1234", []uint8{0x21, 0x34, 0x12}},
	harstonTest{"22 34 12", "ld   ($1234),hl", []uint8{0x22, 0x34, 0x12}},
	harstonTest{"23", "inc  hl", []uint8{0x23}},
	harstonTest{"24", "inc  h", []uint8{0x24}},
	harstonTest{"25", "dec  h", []uint8{0x25}},
	harstonTest{"26 12", "ld   h,$12", []uint8{0x26, 0x12}},
	harstonTest{"27", "daa", []uint8{0x27}},
	harstonTest{"28 10", "jr   z,$0020", []uint8{0x28, 0x10}},
	harstonTest{"29", "add  hl,hl", []uint8{0x29}},
	harstonTest{"2a 34 12", "ld   hl,($1234)", []uint8{0x2a, 0x34, 0x12}},
	harstonTest{"2b", "dec  hl", []uint8{0x2b}},
	harstonTest{"2c", "inc  l", []uint8{0x2c}},
	harstonTest{"2d", "dec  l", []uint8{0x2d}},
	harstonTest{"2e 12", "ld   l,$12", []uint8{0x2e, 0x12}},
	harstonTest{"2f", "cpl", []uint8{0x2f}},
	harstonTest{"30 10", "jr   nc,$0020", []uint8{0x30, 0x10}},
	harstonTest{"31 34 12", "ld   sp,$1234", []uint8{0x31, 0x34, 0x12}},
	harstonTest{"32 34 12", "ld   ($1234),a", []uint8{0x32, 0x34, 0x12}},
	harstonTest{"33", "inc  sp", []uint8{0x33}},
	harstonTest{"34", "inc  (hl)", []uint8{0x34}},
	harstonTest{"35", "dec  (hl)", []uint8{0x35}},
	harstonTest{"36 12", "ld   (hl),$12", []uint8{0x36, 0x12}},
	harstonTest{"37", "scf", []uint8{0x37}},
	harstonTest{"38 10", "jr   c,$0020", []uint8{0x38, 0x10}},
	harstonTest{"39", "add  hl,sp", []uint8{0x39}},
	harstonTest{"3a 34 12", "ld   a,($1234)", []uint8{0x3a, 0x34, 0x12}},
	harstonTest{"3b", "dec  sp", []uint8{0x3b}},
	harstonTest{"3c", "inc  a", []uint8{0x3c}},
	harstonTest{"3d", "dec  a", []uint8{0x3d}},
	harstonTest{"3e 12", "ld   a,$12", []uint8{0x3e, 0x12}},
	harstonTest{"3f", "ccf", []uint8{0x3f}},
	harstonTest{"40", "ld   b,b", []uint8{0x40}},
	harstonTest{"41", "ld   b,c", []uint8{0x41}},
	harstonTest{"42", "ld   b,d", []uint8{0x42}},
	harstonTest{"43", "ld   b,e", []uint8{0x43}},
	harstonTest{"44", "ld   b,h", []uint8{0x44}},
	harstonTest{"45", "ld   b,l", []uint8{0x45}},
	harstonTest{"46", "ld   b,(hl)", []uint8{0x46}},
	harstonTest{"47", "ld   b,a", []uint8{0x47}},
	harstonTest{"48", "ld   c,b", []uint8{0x48}},
	harstonTest{"49", "ld   c,c", []uint8{0x49}},
	harstonTest{"4a", "ld   c,d", []uint8{0x4a}},
	harstonTest{"4b", "ld   c,e", []uint8{0x4b}},
	harstonTest{"4c", "ld   c,h", []uint8{0x4c}},
	harstonTest{"4d", "ld   c,l", []uint8{0x4d}},
	harstonTest{"4e", "ld   c,(hl)", []uint8{0x4e}},
	harstonTest{"4f", "ld   c,a", []uint8{0x4f}},
	harstonTest{"50", "ld   d,b", []uint8{0x50}},
	harstonTest{"51", "ld   d,c", []uint8{0x51}},
	harstonTest{"52", "ld   d,d", []uint8{0x52}},
	harstonTest{"53", "ld   d,e", []uint8{0x53}},
	harstonTest{"54", "ld   d,h", []uint8{0x54}},
	harstonTest{"55", "ld   d,l", []uint8{0x55}},
	harstonTest{"56", "ld   d,(hl)", []uint8{0x56}},
	harstonTest{"57", "ld   d,a", []uint8{0x57}},
	harstonTest{"58", "ld   e,b", []uint8{0x58}},
	harstonTest{"59", "ld   e,c", []uint8{0x59}},
	harstonTest{"5a", "ld   e,d", []uint8{0x5a}},
	harstonTest{"5b", "ld   e,e", []uint8{0x5b}},
	harstonTest{"5c", "ld   e,h", []uint8{0x5c}},
	harstonTest{"5d", "ld   e,l", []uint8{0x5d}},
	harstonTest{"5e", "ld   e,(hl)", []uint8{0x5e}},
	harstonTest{"5f", "ld   e,a", []uint8{0x5f}},
	harstonTest{"60", "ld   h,b", []uint8{0x60}},
	harstonTest{"61", "ld   h,c", []uint8{0x61}},
	harstonTest{"62", "ld   h,d", []uint8{0x62}},
	harstonTest{"63", "ld   h,e", []uint8{0x63}},
	harstonTest{"64", "ld   h,h", []uint8{0x64}},
	harstonTest{"65", "ld   h,l", []uint8{0x65}},
	harstonTest{"66", "ld   h,(hl)", []uint8{0x66}},
	harstonTest{"67", "ld   h,a", []uint8{0x67}},
	harstonTest{"68", "ld   l,b", []uint8{0x68}},
	harstonTest{"69", "ld   l,c", []uint8{0x69}},
	harstonTest{"6a", "ld   l,d", []uint8{0x6a}},
	harstonTest{"6b", "ld   l,e", []uint8{0x6b}},
	harstonTest{"6c", "ld   l,h", []uint8{0x6c}},
	harstonTest{"6d", "ld   l,l", []uint8{0x6d}},
	harstonTest{"6e", "ld   l,(hl)", []uint8{0x6e}},
	harstonTest{"6f", "ld   l,a", []uint8{0x6f}},
	harstonTest{"70", "ld   (hl),b", []uint8{0x70}},
	harstonTest{"71", "ld   (hl),c", []uint8{0x71}},
	harstonTest{"72", "ld   (hl),d", []uint8{0x72}},
	harstonTest{"73", "ld   (hl),e", []uint8{0x73}},
	harstonTest{"74", "ld   (hl),h", []uint8{0x74}},
	harstonTest{"75", "ld   (hl),l", []uint8{0x75}},
	harstonTest{"76", "halt", []uint8{0x76}},
	harstonTest{"77", "ld   (hl),a", []uint8{0x77}},
	harstonTest{"78", "ld   a,b", []uint8{0x78}},
	harstonTest{"79", "ld   a,c", []uint8{0x79}},
	harstonTest{"7a", "ld   a,d", []uint8{0x7a}},
	harstonTest{"7b", "ld   a,e", []uint8{0x7b}},
	harstonTest{"7c", "ld   a,h", []uint8{0x7c}},
	harstonTest{"7d", "ld   a,l", []uint8{0x7d}},
	harstonTest{"7e", "ld   a,(hl)", []uint8{0x7e}},
	harstonTest{"7f", "ld   a,a", []uint8{0x7f}},
	harstonTest{"80", "add  a,b", []uint8{0x80}},
	harstonTest{"81", "add  a,c", []uint8{0x81}},
	harstonTest{"82", "add  a,d", []uint8{0x82}},
	harstonTest{"83", "add  a,e", []uint8{0x83}},
	harstonTest{"84", "add  a,h", []uint8{0x84}},
	harstonTest{"85", "add  a,l", []uint8{0x85}},
	harstonTest{"86", "add  a,(hl)", []uint8{0x86}},
	harstonTest{"87", "add  a,a", []uint8{0x87}},
	harstonTest{"88", "adc  a,b", []uint8{0x88}},
	harstonTest{"89", "adc  a,c", []uint8{0x89}},
	harstonTest{"8a", "adc  a,d", []uint8{0x8a}},
	harstonTest{"8b", "adc  a,e", []uint8{0x8b}},
	harstonTest{"8c", "adc  a,h", []uint8{0x8c}},
	harstonTest{"8d", "adc  a,l", []uint8{0x8d}},
	harstonTest{"8e", "adc  a,(hl)", []uint8{0x8e}},
	harstonTest{"8f", "adc  a,a", []uint8{0x8f}},
	harstonTest{"90", "sub  a,b", []uint8{0x90}},
	harstonTest{"91", "sub  a,c", []uint8{0x91}},
	harstonTest{"92", "sub  a,d", []uint8{0x92}},
	harstonTest{"93", "sub  a,e", []uint8{0x93}},
	harstonTest{"94", "sub  a,h", []uint8{0x94}},
	harstonTest{"95", "sub  a,l", []uint8{0x95}},
	harstonTest{"96", "sub  a,(hl)", []uint8{0x96}},
	harstonTest{"97", "sub  a,a", []uint8{0x97}},
	harstonTest{"98", "sbc  a,b", []uint8{0x98}},
	harstonTest{"99", "sbc  a,c", []uint8{0x99}},
	harstonTest{"9a", "sbc  a,d", []uint8{0x9a}},
	harstonTest{"9b", "sbc  a,e", []uint8{0x9b}},
	harstonTest{"9c", "sbc  a,h", []uint8{0x9c}},
	harstonTest{"9d", "sbc  a,l", []uint8{0x9d}},
	harstonTest{"9e", "sbc  a,(hl)", []uint8{0x9e}},
	harstonTest{"9f", "sbc  a,a", []uint8{0x9f}},
	harstonTest{"a0", "and  b", []uint8{0xa0}},
	harstonTest{"a1", "and  c", []uint8{0xa1}},
	harstonTest{"a2", "and  d", []uint8{0xa2}},
	harstonTest{"a3", "and  e", []uint8{0xa3}},
	harstonTest{"a4", "and  h", []uint8{0xa4}},
	harstonTest{"a5", "and  l", []uint8{0xa5}},
	harstonTest{"a6", "and  (hl)", []uint8{0xa6}},
	harstonTest{"a7", "and  a", []uint8{0xa7}},
	harstonTest{"a8", "xor  b", []uint8{0xa8}},
	harstonTest{"a9", "xor  c", []uint8{0xa9}},
	harstonTest{"aa", "xor  d", []uint8{0xaa}},
	harstonTest{"ab", "xor  e", []uint8{0xab}},
	harstonTest{"ac", "xor  h", []uint8{0xac}},
	harstonTest{"ad", "xor  l", []uint8{0xad}},
	harstonTest{"ae", "xor  (hl)", []uint8{0xae}},
	harstonTest{"af", "xor  a", []uint8{0xaf}},
	harstonTest{"b0", "or   b", []uint8{0xb0}},
	harstonTest{"b1", "or   c", []uint8{0xb1}},
	harstonTest{"b2", "or   d", []uint8{0xb2}},
	harstonTest{"b3", "or   e", []uint8{0xb3}},
	harstonTest{"b4", "or   h", []uint8{0xb4}},
	harstonTest{"b5", "or   l", []uint8{0xb5}},
	harstonTest{"b6", "or   (hl)", []uint8{0xb6}},
	harstonTest{"b7", "or   a", []uint8{0xb7}},
	harstonTest{"b8", "cp   b", []uint8{0xb8}},
	harstonTest{"b9", "cp   c", []uint8{0xb9}},
	harstonTest{"ba", "cp   d", []uint8{0xba}},
	harstonTest{"bb", "cp   e", []uint8{0xbb}},
	harstonTest{"bc", "cp   h", []uint8{0xbc}},
	harstonTest{"bd", "cp   l", []uint8{0xbd}},
	harstonTest{"be", "cp   (hl)", []uint8{0xbe}},
	harstonTest{"bf", "cp   a", []uint8{0xbf}},
	harstonTest{"c0", "ret  nz", []uint8{0xc0}},
	harstonTest{"c1", "pop  bc", []uint8{0xc1}},
	harstonTest{"c2 34 12", "jp   nz,$1234", []uint8{0xc2, 0x34, 0x12}},
	harstonTest{"c3 34 12", "jp   $1234", []uint8{0xc3, 0x34, 0x12}},
	harstonTest{"c4 34 12", "call nz,$1234", []uint8{0xc4, 0x34, 0x12}},
	harstonTest{"c5", "push bc", []uint8{0xc5}},
	harstonTest{"c6 12", "add  a,$12", []uint8{0xc6, 0x12}},
	harstonTest{"c7 00", "rst  $00", []uint8{0xc7, 0x00}},
	harstonTest{"c8", "ret  z", []uint8{0xc8}},
	harstonTest{"c9", "ret", []uint8{0xc9}},
	harstonTest{"ca 34 12", "jp   z,$1234", []uint8{0xca, 0x34, 0x12}},
	harstonTest{"cc 34 12", "call z,$1234", []uint8{0xcc, 0x34, 0x12}},
	harstonTest{"cd 34 12", "call $1234", []uint8{0xcd, 0x34, 0x12}},
	harstonTest{"ce 12", "adc  a,$12", []uint8{0xce, 0x12}},
	harstonTest{"cf", "rst  $08", []uint8{0xcf}},
	harstonTest{"d0", "ret  nc", []uint8{0xd0}},
	harstonTest{"d1", "pop  de", []uint8{0xd1}},
	harstonTest{"d2 34 12", "jp   nc,$1234", []uint8{0xd2, 0x34, 0x12}},
	harstonTest{"d3 12", "out  ($12),a", []uint8{0xd3, 0x12}},
	harstonTest{"d4 34 12", "call nc,$1234", []uint8{0xd4, 0x34, 0x12}},
	harstonTest{"d5", "push de", []uint8{0xd5}},
	harstonTest{"d6 12", "sub  a,$12", []uint8{0xd6, 0x12}},
	harstonTest{"d7", "rst  $10", []uint8{0xd7}},
	harstonTest{"d8", "ret  c", []uint8{0xd8}},
	harstonTest{"d9", "exx", []uint8{0xd9}},
	harstonTest{"da 34 12", "jp   c,$1234", []uint8{0xda, 0x34, 0x12}},
	harstonTest{"db 12", "in   a,($12)", []uint8{0xdb, 0x12}},
	harstonTest{"dc 34 12", "call c,$1234", []uint8{0xdc, 0x34, 0x12}},
	harstonTest{"de 12", "sbc  a,$12", []uint8{0xde, 0x12}},
	harstonTest{"df", "rst  $18", []uint8{0xdf}},
	harstonTest{"e0", "ret  po", []uint8{0xe0}},
	harstonTest{"e1", "pop  hl", []uint8{0xe1}},
	harstonTest{"e2 34 12", "jp   po,$1234", []uint8{0xe2, 0x34, 0x12}},
	harstonTest{"e3", "ex   (sp),hl", []uint8{0xe3}},
	harstonTest{"e4 34 12", "call po,$1234", []uint8{0xe4, 0x34, 0x12}},
	harstonTest{"e5", "push hl", []uint8{0xe5}},
	harstonTest{"e6 12", "and  $12", []uint8{0xe6, 0x12}},
	harstonTest{"e7", "rst  $20", []uint8{0xe7}},
	harstonTest{"e8", "ret  pe", []uint8{0xe8}},
	harstonTest{"e9", "jp   (hl)", []uint8{0xe9}},
	harstonTest{"ea 34 12", "jp   pe,$1234", []uint8{0xea, 0x34, 0x12}},
	harstonTest{"eb", "ex   de,hl", []uint8{0xeb}},
	harstonTest{"ec 34 12", "call pe,$1234", []uint8{0xec, 0x34, 0x12}},
	harstonTest{"ee 12", "xor  $12", []uint8{0xee, 0x12}},
	harstonTest{"ef", "rst  $28", []uint8{0xef}},
	harstonTest{"f0", "ret  p", []uint8{0xf0}},
	harstonTest{"f1", "pop  af", []uint8{0xf1}},
	harstonTest{"f2 34 12", "jp   p,$1234", []uint8{0xf2, 0x34, 0x12}},
	harstonTest{"f3", "di", []uint8{0xf3}},
	harstonTest{"f4 34 12", "call p,$1234", []uint8{0xf4, 0x34, 0x12}},
	harstonTest{"f5", "push af", []uint8{0xf5}},
	harstonTest{"f6 12", "or   $12", []uint8{0xf6, 0x12}},
	harstonTest{"f7", "rst  $30", []uint8{0xf7}},
	harstonTest{"f8", "ret  m", []uint8{0xf8}},
	harstonTest{"f9", "ld   sp,hl", []uint8{0xf9}},
	harstonTest{"fa 34 12", "jp   m,$1234", []uint8{0xfa, 0x34, 0x12}},
	harstonTest{"fb", "ei", []uint8{0xfb}},
	harstonTest{"fc 34 12", "call m,$1234", []uint8{0xfc, 0x34, 0x12}},
	harstonTest{"fe 12", "cp   $12", []uint8{0xfe, 0x12}},
	harstonTest{"ff", "rst  $38", []uint8{0xff}},
	harstonTest{"dd 00", "?dd00", []uint8{0xdd, 0x00}},
	harstonTest{"dd 01", "?dd01", []uint8{0xdd, 0x01}},
	harstonTest{"dd 02", "?dd02", []uint8{0xdd, 0x02}},
	harstonTest{"dd 03", "?dd03", []uint8{0xdd, 0x03}},
	harstonTest{"dd 04", "?dd04", []uint8{0xdd, 0x04}},
	harstonTest{"dd 05", "?dd05", []uint8{0xdd, 0x05}},
	harstonTest{"dd 06", "?dd06", []uint8{0xdd, 0x06}},
	harstonTest{"dd 07", "?dd07", []uint8{0xdd, 0x07}},
	harstonTest{"dd 08", "?dd08", []uint8{0xdd, 0x08}},
	harstonTest{"dd 09", "add  ix,bc", []uint8{0xdd, 0x09}},
	harstonTest{"dd 0a", "?dd0a", []uint8{0xdd, 0x0a}},
	harstonTest{"dd 0b", "?dd0b", []uint8{0xdd, 0x0b}},
	harstonTest{"dd 0c", "?dd0c", []uint8{0xdd, 0x0c}},
	harstonTest{"dd 0d", "?dd0d", []uint8{0xdd, 0x0d}},
	harstonTest{"dd 0e", "?dd0e", []uint8{0xdd, 0x0e}},
	harstonTest{"dd 0f", "?dd0f", []uint8{0xdd, 0x0f}},
	harstonTest{"dd 10", "?dd10", []uint8{0xdd, 0x10}},
	harstonTest{"dd 11", "?dd11", []uint8{0xdd, 0x11}},
	harstonTest{"dd 12", "?dd12", []uint8{0xdd, 0x12}},
	harstonTest{"dd 13", "?dd13", []uint8{0xdd, 0x13}},
	harstonTest{"dd 14", "?dd14", []uint8{0xdd, 0x14}},
	harstonTest{"dd 15", "?dd15", []uint8{0xdd, 0x15}},
	harstonTest{"dd 16", "?dd16", []uint8{0xdd, 0x16}},
	harstonTest{"dd 17", "?dd17", []uint8{0xdd, 0x17}},
	harstonTest{"dd 18", "?dd18", []uint8{0xdd, 0x18}},
	harstonTest{"dd 19", "add  ix,de", []uint8{0xdd, 0x19}},
	harstonTest{"dd 1a", "?dd1a", []uint8{0xdd, 0x1a}},
	harstonTest{"dd 1b", "?dd1b", []uint8{0xdd, 0x1b}},
	harstonTest{"dd 1c", "?dd1c", []uint8{0xdd, 0x1c}},
	harstonTest{"dd 1d", "?dd1d", []uint8{0xdd, 0x1d}},
	harstonTest{"dd 1e", "?dd1e", []uint8{0xdd, 0x1e}},
	harstonTest{"dd 1f", "?dd1f", []uint8{0xdd, 0x1f}},
	harstonTest{"dd 20", "?dd20", []uint8{0xdd, 0x20}},
	harstonTest{"dd 21 34 12", "ld   ix,$1234", []uint8{0xdd, 0x21, 0x34, 0x12}},
	harstonTest{"dd 22 34 12", "ld   ($1234),ix", []uint8{0xdd, 0x22, 0x34, 0x12}},
	harstonTest{"dd 23", "inc  ix", []uint8{0xdd, 0x23}},
	harstonTest{"dd 24", "inc  ixh", []uint8{0xdd, 0x24}},
	harstonTest{"dd 25", "dec  ixh", []uint8{0xdd, 0x25}},
	harstonTest{"dd 26 12", "ld   ixh,$12", []uint8{0xdd, 0x26, 0x12}},
	harstonTest{"dd 27", "?dd27", []uint8{0xdd, 0x27}},
	harstonTest{"dd 28", "?dd28", []uint8{0xdd, 0x28}},
	harstonTest{"dd 29", "add  ix,ix", []uint8{0xdd, 0x29}},
	harstonTest{"dd 2a 34 12", "ld   ix,($1234)", []uint8{0xdd, 0x2a, 0x34, 0x12}},
	harstonTest{"dd 2b", "dec  ix", []uint8{0xdd, 0x2b}},
	harstonTest{"dd 2c", "inc  ixl", []uint8{0xdd, 0x2c}},
	harstonTest{"dd 2d", "dec  ixl", []uint8{0xdd, 0x2d}},
	harstonTest{"dd 2e 12", "ld   ixl,$12", []uint8{0xdd, 0x2e, 0x12}},
	harstonTest{"dd 2f", "?dd2f", []uint8{0xdd, 0x2f}},
	harstonTest{"dd 30", "?dd30", []uint8{0xdd, 0x30}},
	harstonTest{"dd 31", "?dd31", []uint8{0xdd, 0x31}},
	harstonTest{"dd 32", "?dd32", []uint8{0xdd, 0x32}},
	harstonTest{"dd 33", "?dd33", []uint8{0xdd, 0x33}},
	harstonTest{"dd 34 10", "inc  (ix+$10)", []uint8{0xdd, 0x34, 0x10}},
	harstonTest{"dd 35 10", "dec  (ix+$10)", []uint8{0xdd, 0x35, 0x10}},
	harstonTest{"dd 36 10 20", "ld   (ix+$10),$20", []uint8{0xdd, 0x36, 0x10, 0x20}},
	harstonTest{"dd 37", "?dd37", []uint8{0xdd, 0x37}},
	harstonTest{"dd 38", "?dd38", []uint8{0xdd, 0x38}},
	harstonTest{"dd 39", "add  ix,sp", []uint8{0xdd, 0x39}},
	harstonTest{"dd 3a", "?dd3a", []uint8{0xdd, 0x3a}},
	harstonTest{"dd 3b", "?dd3b", []uint8{0xdd, 0x3b}},
	harstonTest{"dd 3c", "?dd3c", []uint8{0xdd, 0x3c}},
	harstonTest{"dd 3d", "?dd3d", []uint8{0xdd, 0x3d}},
	harstonTest{"dd 3e", "?dd3e", []uint8{0xdd, 0x3e}},
	harstonTest{"dd 3f", "?dd3f", []uint8{0xdd, 0x3f}},
	harstonTest{"dd 40", "?dd40", []uint8{0xdd, 0x40}},
	harstonTest{"dd 41", "?dd41", []uint8{0xdd, 0x41}},
	harstonTest{"dd 42", "?dd42", []uint8{0xdd, 0x42}},
	harstonTest{"dd 43", "?dd43", []uint8{0xdd, 0x43}},
	harstonTest{"dd 44", "ld   b,ixh", []uint8{0xdd, 0x44}},
	harstonTest{"dd 45", "ld   b,ixl", []uint8{0xdd, 0x45}},
	harstonTest{"dd 46 10", "ld   b,(ix+$10)", []uint8{0xdd, 0x46, 0x10}},
	harstonTest{"dd 47", "?dd47", []uint8{0xdd, 0x47}},
	harstonTest{"dd 48", "?dd48", []uint8{0xdd, 0x48}},
	harstonTest{"dd 49", "?dd49", []uint8{0xdd, 0x49}},
	harstonTest{"dd 4a", "?dd4a", []uint8{0xdd, 0x4a}},
	harstonTest{"dd 4b", "?dd4b", []uint8{0xdd, 0x4b}},
	harstonTest{"dd 4c", "ld   c,ixh", []uint8{0xdd, 0x4c}},
	harstonTest{"dd 4d", "ld   c,ixl", []uint8{0xdd, 0x4d}},
	harstonTest{"dd 4e 10", "ld   c,(ix+$10)", []uint8{0xdd, 0x4e, 0x10}},
	harstonTest{"dd 4f", "?dd4f", []uint8{0xdd, 0x4f}},
	harstonTest{"dd 50", "?dd50", []uint8{0xdd, 0x50}},
	harstonTest{"dd 51", "?dd51", []uint8{0xdd, 0x51}},
	harstonTest{"dd 52", "?dd52", []uint8{0xdd, 0x52}},
	harstonTest{"dd 53", "?dd53", []uint8{0xdd, 0x53}},
	harstonTest{"dd 54", "ld   d,ixh", []uint8{0xdd, 0x54}},
	harstonTest{"dd 55", "ld   d,ixl", []uint8{0xdd, 0x55}},
	harstonTest{"dd 56 10", "ld   d,(ix+$10)", []uint8{0xdd, 0x56, 0x10}},
	harstonTest{"dd 57", "?dd57", []uint8{0xdd, 0x57}},
	harstonTest{"dd 58", "?dd58", []uint8{0xdd, 0x58}},
	harstonTest{"dd 59", "?dd59", []uint8{0xdd, 0x59}},
	harstonTest{"dd 5a", "?dd5a", []uint8{0xdd, 0x5a}},
	harstonTest{"dd 5b", "?dd5b", []uint8{0xdd, 0x5b}},
	harstonTest{"dd 5c", "ld   e,ixh", []uint8{0xdd, 0x5c}},
	harstonTest{"dd 5d", "ld   e,ixl", []uint8{0xdd, 0x5d}},
	harstonTest{"dd 5e 10", "ld   e,(ix+$10)", []uint8{0xdd, 0x5e, 0x10}},
	harstonTest{"dd 5f", "?dd5f", []uint8{0xdd, 0x5f}},
	harstonTest{"dd 60", "ld   ixh,b", []uint8{0xdd, 0x60}},
	harstonTest{"dd 61", "ld   ixh,c", []uint8{0xdd, 0x61}},
	harstonTest{"dd 62", "ld   ixh,d", []uint8{0xdd, 0x62}},
	harstonTest{"dd 63", "ld   ixh,e", []uint8{0xdd, 0x63}},
	harstonTest{"dd 64", "ld   ixh,ixh", []uint8{0xdd, 0x64}},
	harstonTest{"dd 65", "ld   ixh,ixl", []uint8{0xdd, 0x65}},
	harstonTest{"dd 66 10", "ld   h,(ix+$10)", []uint8{0xdd, 0x66, 0x10}},
	harstonTest{"dd 67", "ld   ixh,a", []uint8{0xdd, 0x67}},
	harstonTest{"dd 68", "ld   ixl,b", []uint8{0xdd, 0x68}},
	harstonTest{"dd 69", "ld   ixl,c", []uint8{0xdd, 0x69}},
	harstonTest{"dd 6a", "ld   ixl,d", []uint8{0xdd, 0x6a}},
	harstonTest{"dd 6b", "ld   ixl,e", []uint8{0xdd, 0x6b}},
	harstonTest{"dd 6c", "ld   ixl,ixh", []uint8{0xdd, 0x6c}},
	harstonTest{"dd 6d", "ld   ixl,ixl", []uint8{0xdd, 0x6d}},
	harstonTest{"dd 6e 10", "ld   l,(ix+$10)", []uint8{0xdd, 0x6e, 0x10}},
	harstonTest{"dd 6f", "ld   ixl,a", []uint8{0xdd, 0x6f}},
	harstonTest{"dd 70 10", "ld   (ix+$10),b", []uint8{0xdd, 0x70, 0x10}},
	harstonTest{"dd 71 10", "ld   (ix+$10),c", []uint8{0xdd, 0x71, 0x10}},
	harstonTest{"dd 72 10", "ld   (ix+$10),d", []uint8{0xdd, 0x72, 0x10}},
	harstonTest{"dd 73 10", "ld   (ix+$10),e", []uint8{0xdd, 0x73, 0x10}},
	harstonTest{"dd 74 10", "ld   (ix+$10),h", []uint8{0xdd, 0x74, 0x10}},
	harstonTest{"dd 75 10", "ld   (ix+$10),l", []uint8{0xdd, 0x75, 0x10}},
	harstonTest{"dd 76", "?dd76", []uint8{0xdd, 0x76}},
	harstonTest{"dd 77 10", "ld   (ix+$10),a", []uint8{0xdd, 0x77, 0x10}},
	harstonTest{"dd 78", "?dd78", []uint8{0xdd, 0x78}},
	harstonTest{"dd 79", "?dd79", []uint8{0xdd, 0x79}},
	harstonTest{"dd 7a", "?dd7a", []uint8{0xdd, 0x7a}},
	harstonTest{"dd 7b", "?dd7b", []uint8{0xdd, 0x7b}},
	harstonTest{"dd 7c", "ld   a,ixh", []uint8{0xdd, 0x7c}},
	harstonTest{"dd 7d", "ld   a,ixl", []uint8{0xdd, 0x7d}},
	harstonTest{"dd 7e 10", "ld   a,(ix+$10)", []uint8{0xdd, 0x7e, 0x10}},
	harstonTest{"dd 7f", "?dd7f", []uint8{0xdd, 0x7f}},
	harstonTest{"dd 80", "?dd80", []uint8{0xdd, 0x80}},
	harstonTest{"dd 81", "?dd81", []uint8{0xdd, 0x81}},
	harstonTest{"dd 82", "?dd82", []uint8{0xdd, 0x82}},
	harstonTest{"dd 83", "?dd83", []uint8{0xdd, 0x83}},
	harstonTest{"dd 84", "add  a,ixh", []uint8{0xdd, 0x84}},
	harstonTest{"dd 85", "add  a,ixl", []uint8{0xdd, 0x85}},
	harstonTest{"dd 86 10", "add  a,(ix+$10)", []uint8{0xdd, 0x86, 0x10}},
	harstonTest{"dd 87", "?dd87", []uint8{0xdd, 0x87}},
	harstonTest{"dd 88", "?dd88", []uint8{0xdd, 0x88}},
	harstonTest{"dd 89", "?dd89", []uint8{0xdd, 0x89}},
	harstonTest{"dd 8a", "?dd8a", []uint8{0xdd, 0x8a}},
	harstonTest{"dd 8b", "?dd8b", []uint8{0xdd, 0x8b}},
	harstonTest{"dd 8c", "adc  a,ixh", []uint8{0xdd, 0x8c}},
	harstonTest{"dd 8d", "adc  a,ixl", []uint8{0xdd, 0x8d}},
	harstonTest{"dd 8e 10", "adc  a,(ix+$10)", []uint8{0xdd, 0x8e, 0x10}},
	harstonTest{"dd 8f", "?dd8f", []uint8{0xdd, 0x8f}},
	harstonTest{"dd 90", "?dd90", []uint8{0xdd, 0x90}},
	harstonTest{"dd 91", "?dd91", []uint8{0xdd, 0x91}},
	harstonTest{"dd 92", "?dd92", []uint8{0xdd, 0x92}},
	harstonTest{"dd 93", "?dd93", []uint8{0xdd, 0x93}},
	harstonTest{"dd 94", "sub  a,ixh", []uint8{0xdd, 0x94}},
	harstonTest{"dd 95", "sub  a,ixl", []uint8{0xdd, 0x95}},
	harstonTest{"dd 96 10", "sub  a,(ix+$10)", []uint8{0xdd, 0x96, 0x10}},
	harstonTest{"dd 97", "?dd97", []uint8{0xdd, 0x97}},
	harstonTest{"dd 98", "?dd98", []uint8{0xdd, 0x98}},
	harstonTest{"dd 99", "?dd99", []uint8{0xdd, 0x99}},
	harstonTest{"dd 9a", "?dd9a", []uint8{0xdd, 0x9a}},
	harstonTest{"dd 9b", "?dd9b", []uint8{0xdd, 0x9b}},
	harstonTest{"dd 9c", "sbc  a,ixh", []uint8{0xdd, 0x9c}},
	harstonTest{"dd 9d", "sbc  a,ixl", []uint8{0xdd, 0x9d}},
	harstonTest{"dd 9e 10", "sbc  a,(ix+$10)", []uint8{0xdd, 0x9e, 0x10}},
	harstonTest{"dd 9f", "?dd9f", []uint8{0xdd, 0x9f}},
	harstonTest{"dd a0", "?dda0", []uint8{0xdd, 0xa0}},
	harstonTest{"dd a1", "?dda1", []uint8{0xdd, 0xa1}},
	harstonTest{"dd a2", "?dda2", []uint8{0xdd, 0xa2}},
	harstonTest{"dd a3", "?dda3", []uint8{0xdd, 0xa3}},
	harstonTest{"dd a4", "and  ixh", []uint8{0xdd, 0xa4}},
	harstonTest{"dd a5", "and  ixl", []uint8{0xdd, 0xa5}},
	harstonTest{"dd a6 10", "and  (ix+$10)", []uint8{0xdd, 0xa6, 0x10}},
	harstonTest{"dd a7", "?dda7", []uint8{0xdd, 0xa7}},
	harstonTest{"dd a8", "?dda8", []uint8{0xdd, 0xa8}},
	harstonTest{"dd a9", "?dda9", []uint8{0xdd, 0xa9}},
	harstonTest{"dd aa", "?ddaa", []uint8{0xdd, 0xaa}},
	harstonTest{"dd ab", "?ddab", []uint8{0xdd, 0xab}},
	harstonTest{"dd ac", "xor  ixh", []uint8{0xdd, 0xac}},
	harstonTest{"dd ad", "xor  ixl", []uint8{0xdd, 0xad}},
	harstonTest{"dd ae 10", "xor  (ix+$10)", []uint8{0xdd, 0xae, 0x10}},
	harstonTest{"dd af", "?ddaf", []uint8{0xdd, 0xaf}},
	harstonTest{"dd b0", "?ddb0", []uint8{0xdd, 0xb0}},
	harstonTest{"dd b1", "?ddb1", []uint8{0xdd, 0xb1}},
	harstonTest{"dd b2", "?ddb2", []uint8{0xdd, 0xb2}},
	harstonTest{"dd b3", "?ddb3", []uint8{0xdd, 0xb3}},
	harstonTest{"dd b4", "or   ixh", []uint8{0xdd, 0xb4}},
	harstonTest{"dd b5", "or   ixl", []uint8{0xdd, 0xb5}},
	harstonTest{"dd b6 10", "or   (ix+$10)", []uint8{0xdd, 0xb6, 0x10}},
	harstonTest{"dd b7", "?ddb7", []uint8{0xdd, 0xb7}},
	harstonTest{"dd b8", "?ddb8", []uint8{0xdd, 0xb8}},
	harstonTest{"dd b9", "?ddb9", []uint8{0xdd, 0xb9}},
	harstonTest{"dd ba", "?ddba", []uint8{0xdd, 0xba}},
	harstonTest{"dd bb", "?ddbb", []uint8{0xdd, 0xbb}},
	harstonTest{"dd bc", "cp   ixh", []uint8{0xdd, 0xbc}},
	harstonTest{"dd bd", "cp   ixl", []uint8{0xdd, 0xbd}},
	harstonTest{"dd be 10", "cp   (ix+$10)", []uint8{0xdd, 0xbe, 0x10}},
	harstonTest{"dd bf", "?ddbf", []uint8{0xdd, 0xbf}},
	harstonTest{"dd c0", "?ddc0", []uint8{0xdd, 0xc0}},
	harstonTest{"dd c1", "?ddc1", []uint8{0xdd, 0xc1}},
	harstonTest{"dd c2", "?ddc2", []uint8{0xdd, 0xc2}},
	harstonTest{"dd c3", "?ddc3", []uint8{0xdd, 0xc3}},
	harstonTest{"dd c4", "?ddc4", []uint8{0xdd, 0xc4}},
	harstonTest{"dd c5", "?ddc5", []uint8{0xdd, 0xc5}},
	harstonTest{"dd c6", "?ddc6", []uint8{0xdd, 0xc6}},
	harstonTest{"dd c7", "?ddc7", []uint8{0xdd, 0xc7}},
	harstonTest{"dd c8", "?ddc8", []uint8{0xdd, 0xc8}},
	harstonTest{"dd c9", "?ddc9", []uint8{0xdd, 0xc9}},
	harstonTest{"dd ca", "?ddca", []uint8{0xdd, 0xca}},
	harstonTest{"dd cb", "?ddcb", []uint8{0xdd, 0xcb}},
	harstonTest{"dd cc", "?ddcc", []uint8{0xdd, 0xcc}},
	harstonTest{"dd cd", "?ddcd", []uint8{0xdd, 0xcd}},
	harstonTest{"dd ce", "?ddce", []uint8{0xdd, 0xce}},
	harstonTest{"dd cf", "?ddcf", []uint8{0xdd, 0xcf}},
	harstonTest{"dd d0", "?ddd0", []uint8{0xdd, 0xd0}},
	harstonTest{"dd d1", "?ddd1", []uint8{0xdd, 0xd1}},
	harstonTest{"dd d2", "?ddd2", []uint8{0xdd, 0xd2}},
	harstonTest{"dd d3", "?ddd3", []uint8{0xdd, 0xd3}},
	harstonTest{"dd d4", "?ddd4", []uint8{0xdd, 0xd4}},
	harstonTest{"dd d5", "?ddd5", []uint8{0xdd, 0xd5}},
	harstonTest{"dd d6", "?ddd6", []uint8{0xdd, 0xd6}},
	harstonTest{"dd d7", "?ddd7", []uint8{0xdd, 0xd7}},
	harstonTest{"dd d8", "?ddd8", []uint8{0xdd, 0xd8}},
	harstonTest{"dd d9", "?ddd9", []uint8{0xdd, 0xd9}},
	harstonTest{"dd da", "?ddda", []uint8{0xdd, 0xda}},
	harstonTest{"dd db", "?dddb", []uint8{0xdd, 0xdb}},
	harstonTest{"dd dc", "?dddc", []uint8{0xdd, 0xdc}},
	harstonTest{"dd dd", "?dd", []uint8{0xdd, 0xdd}},
	harstonTest{"dd de", "?ddde", []uint8{0xdd, 0xde}},
	harstonTest{"dd df", "?dddf", []uint8{0xdd, 0xdf}},
	harstonTest{"dd e0", "?dde0", []uint8{0xdd, 0xe0}},
	harstonTest{"dd e1", "pop  ix", []uint8{0xdd, 0xe1}},
	harstonTest{"dd e2", "?dde2", []uint8{0xdd, 0xe2}},
	harstonTest{"dd e3", "ex   (sp),ix", []uint8{0xdd, 0xe3}},
	harstonTest{"dd e4", "?dde4", []uint8{0xdd, 0xe4}},
	harstonTest{"dd e5", "push ix", []uint8{0xdd, 0xe5}},
	harstonTest{"dd e6", "?dde6", []uint8{0xdd, 0xe6}},
	harstonTest{"dd e7", "?dde7", []uint8{0xdd, 0xe7}},
	harstonTest{"dd e8", "?dde8", []uint8{0xdd, 0xe8}},
	harstonTest{"dd e9", "jp   (ix)", []uint8{0xdd, 0xe9}},
	harstonTest{"dd ea", "?ddea", []uint8{0xdd, 0xea}},
	harstonTest{"dd eb", "?ddeb", []uint8{0xdd, 0xeb}},
	harstonTest{"dd ec", "?ddec", []uint8{0xdd, 0xec}},
	harstonTest{"dd ed", "?dd", []uint8{0xdd, 0xed}},
	harstonTest{"dd ee", "?ddee", []uint8{0xdd, 0xee}},
	harstonTest{"dd ef", "?ddef", []uint8{0xdd, 0xef}},
	harstonTest{"dd f0", "?ddf0", []uint8{0xdd, 0xf0}},
	harstonTest{"dd f1", "?ddf1", []uint8{0xdd, 0xf1}},
	harstonTest{"dd f2", "?ddf2", []uint8{0xdd, 0xf2}},
	harstonTest{"dd f3", "?ddf3", []uint8{0xdd, 0xf3}},
	harstonTest{"dd f4", "?ddf4", []uint8{0xdd, 0xf4}},
	harstonTest{"dd f5", "?ddf5", []uint8{0xdd, 0xf5}},
	harstonTest{"dd f6", "?ddf6", []uint8{0xdd, 0xf6}},
	harstonTest{"dd f7", "?ddf7", []uint8{0xdd, 0xf7}},
	harstonTest{"dd f8", "?ddf8", []uint8{0xdd, 0xf8}},
	harstonTest{"dd f9", "?ddf9", []uint8{0xdd, 0xf9}},
	harstonTest{"dd fa", "?ddfa", []uint8{0xdd, 0xfa}},
	harstonTest{"dd fb", "?ddfb", []uint8{0xdd, 0xfb}},
	harstonTest{"dd fc", "?ddfc", []uint8{0xdd, 0xfc}},
	harstonTest{"dd fd", "?dd", []uint8{0xdd, 0xfd}},
	harstonTest{"dd fe", "?ddfe", []uint8{0xdd, 0xfe}},
	harstonTest{"dd ff", "?ddff", []uint8{0xdd, 0xff}},
}
