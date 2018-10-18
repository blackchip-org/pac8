package memory

import "github.com/blackchip-org/pac8/bits"

// Cursor points to a location in memory.
type Cursor struct {
	Pos uint16 // Current position.
	mem Memory
}

// NewCursor creates a new Cursor on mem pointing at address 0.
func NewCursor(mem Memory) *Cursor {
	return &Cursor{mem: mem}
}

// Fetch returns the byte at c.Pos as an 8-bit value and advances c.Pos by one.
func (c *Cursor) Fetch() uint8 {
	value := c.mem.Load(c.Pos)
	c.Pos++
	return value
}

// Peek returns the byte at c.Pos as an 8-bit value.
func (c *Cursor) Peek() uint8 {
	return c.mem.Load(c.Pos)
}

// FetchLE returns the next two bytes as a 16-bit value stored in little
// endian format and advances c.Pos by two.
func (c *Cursor) FetchLE() uint16 {
	lo := c.Fetch()
	hi := c.Fetch()
	return bits.Join(hi, lo)
}

// Put sets the value at c.Pos and advances c.Pos by one.
func (c *Cursor) Put(value uint8) {
	c.mem.Store(c.Pos, value)
	c.Pos++
}

// PutN calls c.Put for each value in values.
func (c *Cursor) PutN(values ...uint8) {
	for _, value := range values {
		c.Put(value)
	}
}
