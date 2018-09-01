package memory

type Cursor struct {
	Pos uint16
	mem Memory
}

func NewCursor(mem Memory) *Cursor {
	return &Cursor{mem: mem}
}

func (c *Cursor) Load() uint8 {
	value := c.mem.Load(c.Pos)
	c.Pos++
	return value
}

func (c *Cursor) Store(value uint8) {
	c.mem.Store(c.Pos, value)
	c.Pos++
}

func (c *Cursor) StoreN(values ...uint8) {
	for _, value := range values {
		c.Store(value)
	}
}
