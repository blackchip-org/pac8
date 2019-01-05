package memory

import (
	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/pkg/util/state"
)

// Memory is a chunk of 8-bit values accessed by a 16-bit address.
type Memory interface {
	// Load returns the value from the address at addr.
	Load(addr uint16) uint8

	// Store puts the value of v at the address at addr.
	Store(addr uint16, v uint8)

	// Length is the number of 8-bit values in this memory.
	Length() int

	Save(*state.Encoder)

	Restore(*state.Decoder)
}

// Block is a chunck of memory found at a specific address.
type Block struct {
	// Mem is the memory for this block
	Mem Memory

	// Addr is the address that Mem.Load(0) represents.
	Addr uint16
}

// NewBlock creates a new Block of memory at the address of addr.
func NewBlock(addr uint16, mem Memory) Block {
	return Block{Mem: mem, Addr: addr}
}

type ram struct {
	bytes []uint8
}

// NewRAM creates a chunk of memory with a length of len that can be
// read and written to.
func NewRAM(len int) Memory {
	return &ram{bytes: make([]uint8, len, len)}
}

func (r ram) Load(address uint16) uint8 {
	return r.bytes[address]
}

func (r ram) Store(address uint16, value uint8) {
	r.bytes[address] = value
}

func (r ram) Length() int {
	return len(r.bytes)
}

func (r ram) Save(enc *state.Encoder) {
	enc.Encode(r.bytes)
}

func (r ram) Restore(dec *state.Decoder) {
	dec.Decode(&r.bytes)
}

type rom struct {
	bytes []uint8
}

// NewROM creates a chunk of read-only memory that accesses data.
func NewROM(data []uint8) Memory {
	return rom{bytes: data}
}

func (r rom) Load(address uint16) uint8 {
	if int(address) >= len(r.bytes) {
		return 0
	}
	return r.bytes[address]
}

func (r rom) Store(address uint16, value uint8) {}

func (r rom) Length() int {
	return len(r.bytes)
}

func (r rom) Save(enc *state.Encoder) {}

func (r rom) Restore(dec *state.Decoder) {}

type null struct {
	length int
}

// NewNull creates a chunk of memory with a length of len that always returns
// zero when read. Writes are ignored.
func NewNull(len int) Memory {
	return null{length: len}
}

func (n null) Load(address uint16) uint8 {
	return 0
}

func (n null) Store(address uint16, value uint8) {}

func (n null) Length() int {
	return n.length
}

func (n null) Save(enc *state.Encoder) {}

func (n null) Restore(dec *state.Decoder) {}

type pageMap struct {
	mem    Memory
	offset uint16
}

type pageMapped struct {
	blocks []Block
	pages  [256]pageMap
}

// NewPageMapped creates a memory that combines multiple memory blocks
// into a single addressable memory mapped at the page level. Each block
// must have a length that is divisible by a page and addressed at a
// page boundary. Unmapped pages return zero when read and are ignored
// when written.
func NewPageMapped(blocks []Block) Memory {
	pm := &pageMapped{}
	for i := 0; i < 256; i++ {
		pm.pages[i] = pageMap{mem: NewNull(0x100), offset: 0}
	}

	for _, block := range blocks {
		if block.Addr%0x100 != 0 {
			panic("memory block must start on page boundary")
		}
		if block.Mem.Length()%0x100 != 0 {
			panic("memory block length must be a multiple of a page")
		}

		for offset := 0; offset < block.Mem.Length(); offset += 256 {
			page := (block.Addr + uint16(offset)) / 256
			pm.pages[page] = pageMap{mem: block.Mem, offset: uint16(offset)}
		}
	}
	pm.blocks = blocks
	return pm
}

func (m pageMapped) Load(address uint16) uint8 {
	pageN, index := bits.Split(address)
	page := m.pages[pageN]
	return page.mem.Load(page.offset + uint16(index))
}

func (m pageMapped) Store(address uint16, value uint8) {
	pageN, index := bits.Split(address)
	page := m.pages[pageN]
	page.mem.Store(page.offset+uint16(index), value)
}

func (m pageMapped) Length() int {
	return 0x10000
}

func (m pageMapped) Save(enc *state.Encoder) {
	for _, b := range m.blocks {
		b.Mem.Save(enc)
	}
}

func (m pageMapped) Restore(dec *state.Decoder) {
	for _, b := range m.blocks {
		b.Mem.Restore(dec)
	}
}

type BlockMapper struct {
	Blocks []Block
}

func NewBlockMapper() *BlockMapper {
	return &BlockMapper{
		Blocks: make([]Block, 0, 0),
	}
}

func (b *BlockMapper) Map(addr uint16, mem Memory) {
	if addr%0x100 != 0 {
		panic("memory block must start on page boundary")
	}
	if mem.Length()%0x100 != 0 {
		panic("memory block length must be a multiple of a page")
	}
	b.Blocks = append(b.Blocks, NewBlock(addr, mem))
}
