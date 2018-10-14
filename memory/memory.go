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

type RAM struct {
	bytes []uint8
}

func NewRAM(size int) *RAM {
	return &RAM{bytes: make([]uint8, size, size)}
}

func (r *RAM) Load(address uint16) uint8 {
	return r.bytes[address]
}

func (r *RAM) Store(address uint16, value uint8) {
	r.bytes[address] = value
}

func (r *RAM) Length() int {
	return len(r.bytes)
}

type ROM struct {
	bytes []uint8
}

func NewROM(data []uint8) *ROM {
	return &ROM{bytes: data}
}

func (r *ROM) Load(address uint16) uint8 {
	if int(address) >= len(r.bytes) {
		return 0
	}
	return r.bytes[address]
}

func (r *ROM) Store(address uint16, value uint8) {}

func (r *ROM) Length() int {
	return len(r.bytes)
}

type Null struct {
	length int
}

func NewNull(length int) Null {
	return Null{length: length}
}

func (n Null) Load(address uint16) uint8 {
	return 0
}

func (n Null) Store(address uint16, value uint8) {}

func (n Null) Length() int {
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
				mem.pages[p] = page{Null{}, 0x100}
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

type AddrMasked struct {
	mem  Memory
	mask uint16
}

func NewAddrMasked(mem Memory, mask uint16) *AddrMasked {
	return &AddrMasked{
		mem:  mem,
		mask: mask,
	}
}

func (m *AddrMasked) Load(address uint16) uint8 {
	return m.mem.Load(address & m.mask)
}

func (m *AddrMasked) Store(address uint16, value uint8) {
	m.mem.Store(address&m.mask, value)
}

func (m *AddrMasked) Length() int {
	return m.mem.Length()
}
