package memory

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/blackchip-org/pac8/util/bits"
)

type Memory interface {
	Load(address uint16) uint8
	Store(address uint16, value uint8)
	Length() int
}

type Memory16 interface {
	Load(address uint16) uint16
	Store(address uint16, value uint16)
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

func (r *ROM) Checksum() string {
	return fmt.Sprintf("%040x", sha1.Sum(r.bytes))
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

type Port struct {
	Ready bool
	Read  *uint8
	Write *uint8
}

type IO struct {
	Ports []Port
}

func NewIO(n int) *IO {
	return &IO{
		Ports: make([]Port, n, n),
	}
}

func (i *IO) Store(address uint16, value uint8) {
	p := i.Ports[int(address)]
	if p.Write != nil {
		p.Ready = true
		*p.Write = value
	}
}

func (i *IO) Load(address uint16) uint8 {
	p := i.Ports[int(address)]
	if p.Read == nil {
		return 0
	}
	return *p.Read
}

func (i *IO) Length() int {
	return len(i.Ports)
}

func (i *IO) RO(port int, v *uint8) {
	i.Ports[port].Read = v
}

func (i *IO) WO(port int, v *uint8) {
	i.Ports[port].Write = v
}

func (i *IO) RW(port int, v *uint8) {
	i.Ports[port].Read = v
	i.Ports[port].Write = v
}

type LittleEndian struct {
	mem Memory
}

func NewLittleEndian(mem Memory) LittleEndian {
	return LittleEndian{mem: mem}
}

func (e LittleEndian) Load(address uint16) uint16 {
	lo := e.mem.Load(address)
	hi := e.mem.Load(address + 1)
	return bits.Join(hi, lo)
}

func (e LittleEndian) Store(address uint16, value uint16) {
	hi, lo := bits.Split(value)
	e.mem.Store(address, lo)
	e.mem.Store(address+1, hi)
}

type Snapshot struct {
	Address uint16
	Values  []uint8
}

func Import(m Memory, snapshot Snapshot) {
	for i, value := range snapshot.Values {
		m.Store(snapshot.Address+uint16(i), value)
	}
}

type Diff struct {
	Address uint16
	A       uint8
	B       uint8
}

func (d *Diff) String() string {
	return fmt.Sprintf("%04x: %02x %02x", d.Address, d.A, d.B)
}

type DiffReport []Diff

func (d DiffReport) String() string {
	reports := make([]string, 0, 0)
	for _, diff := range d {
		reports = append(reports, diff.String())
	}
	return strings.Join(reports, "\n")
}

// Compare creates a report of all differences between memory a and
// memory b. Returns true if the memories are identical.
func Compare(a Memory, b Memory) (DiffReport, bool) {
	diff := make([]Diff, 0, 0)
	for addr := 0; addr < 0x10000; addr++ {
		aval := a.Load(uint16(addr))
		bval := b.Load(uint16(addr))
		if aval != bval {
			diff = append(diff, Diff{Address: uint16(addr), A: aval, B: bval})
		}
	}
	return diff, len(diff) == 0
}

func Verify(a Memory, b []Snapshot) (DiffReport, bool) {
	diff := make([]Diff, 0, 0)
	cursor := NewCursor(a)
	for _, snapshot := range b {
		cursor.Pos = snapshot.Address
		for i, bval := range snapshot.Values {
			aval := cursor.Fetch()
			if aval != bval {
				diff = append(diff, Diff{
					Address: snapshot.Address + uint16(i),
					A:       aval,
					B:       bval,
				})
			}
		}
	}
	return diff, len(diff) == 0
}
