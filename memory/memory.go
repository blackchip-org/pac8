package memory

import (
	"fmt"
	"strings"
)

type Memory interface {
	Load(address uint16) uint8
	Store(address uint16, value uint8)
}

type RAM struct {
	bytes []uint8
}

func NewRAM(size int) RAM {
	return RAM{bytes: make([]uint8, size, size)}
}

func (r RAM) Load(address uint16) uint8 {
	return r.bytes[address]
}

func (r RAM) Store(address uint16, value uint8) {
	r.bytes[address] = value
}

type ROM struct {
	bytes []uint8
}

func NewROM(data []uint8) ROM {
	return ROM{bytes: data}
}

func (r ROM) Load(address uint16) uint8 {
	if int(address) >= len(r.bytes) {
		return 0
	}
	return r.bytes[address]
}

func (r ROM) Store(address uint16, value uint8) {}

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
			aval := cursor.Load()
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
