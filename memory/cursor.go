package memory

import "github.com/blackchip-org/pac8/util/bits"

type Cursor struct {
	Pos uint16
	mem Memory
}

func NewCursor(mem Memory) *Cursor {
	return &Cursor{mem: mem}
}

func (c *Cursor) Fetch() uint8 {
	value := c.mem.Load(c.Pos)
	c.Pos++
	return value
}

func (c *Cursor) Peek() uint8 {
	return c.mem.Load(c.Pos)
}

func (c *Cursor) FetchLE() uint16 {
	lo := c.Fetch()
	hi := c.Fetch()
	return bits.Join(hi, lo)
}

func (c *Cursor) Put(value uint8) {
	c.mem.Store(c.Pos, value)
	c.Pos++
}

func (c *Cursor) PutN(values ...uint8) {
	for _, value := range values {
		c.Put(value)
	}
}
