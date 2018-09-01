package z80

var ops = map[uint8]func(c *CPU){
	0x78: func(c *CPU) { ld(c, c.storeA, c.loadB) },
}
