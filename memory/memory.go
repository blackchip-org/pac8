package memory

import (
	"fmt"

	"github.com/blackchip-org/pac8/util/bits"
)

type Memory interface {
	Load(address uint16) uint8
	Store(address uint16, value uint8)
	Length() int
}

type ram struct {
	bytes []uint8
}

func NewRAM(size int) Memory {
	return ram{bytes: make([]uint8, size, size)}
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

type rom struct {
	bytes []uint8
}

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

type null struct {
	length int
}

func NewNull(length int) Memory {
	return null{length: length}
}

func (n null) Load(address uint16) uint8 {
	return 0
}

func (n null) Store(address uint16, value uint8) {}

func (n null) Length() int {
	return n.length
}

type page struct {
	mem    Memory
	offset uint16
}

type PageMapped struct {
	pages []page
}

func NewPageMapped(blocks []Memory) *PageMapped {
	blockIndex := 0
	block := blocks[blockIndex]
	offset := uint16(0)
	remaining := block.Length()
	mem := &PageMapped{pages: make([]page, 0x100, 0x100)}
	for p := 0; p < 0x100; p++ {
		if remaining < 0 {
			panic(fmt.Sprintf("memory has invalid length: %v", block.Length()))
		}
		if remaining == 0 {
			if blockIndex+1 == len(blocks) {
				mem.pages[p] = page{null{}, 0x100}
				continue
			}
			blockIndex++
			block = blocks[blockIndex]
			offset = 0
			remaining = block.Length()
		}
		mem.pages[p] = page{block, offset}
		offset += 0x100
		remaining -= 0x100
	}
	if remaining != 0 {
		panic(fmt.Sprintf("too many memory blocks"))
	}
	return mem
}

func (m PageMapped) Load(address uint16) uint8 {
	pageN, offset1 := bits.Split(address)
	page := m.pages[pageN]
	return page.mem.Load(page.offset + uint16(offset1))
}

func (m PageMapped) Store(address uint16, value uint8) {
	pageN, offset1 := bits.Split(address)
	page := m.pages[pageN]
	page.mem.Store(page.offset+uint16(offset1), value)
}

func (m PageMapped) Length() int {
	return 0x10000
}

type masked struct {
	mem  Memory
	mask uint16
}

func NewMasked(mem Memory, mask uint16) Memory {
	return masked{
		mem:  mem,
		mask: mask,
	}
}

func (m masked) Load(address uint16) uint8 {
	return m.mem.Load(address & m.mask)
}

func (m masked) Store(address uint16, value uint8) {
	m.mem.Store(address&m.mask, value)
}

func (m masked) Length() int {
	return m.mem.Length()
}
