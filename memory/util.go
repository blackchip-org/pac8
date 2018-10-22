package memory

import (
	"fmt"
	"strings"

	"github.com/blackchip-org/pac8/bits"
)

// LoadLE loads a 16-bit little endian value from memory m at addr.
func LoadLE(m Memory, addr uint16) uint16 {
	lo := m.Load(addr)
	hi := m.Load(addr + 1)
	return bits.Join(hi, lo)
}

// StoreLE stores a 16-bit little endian value to memory m at addr.
func StoreLE(m Memory, addr uint16, value uint16) {
	hi, lo := bits.Split(value)
	m.Store(addr, lo)
	m.Store(addr+1, hi)
}

// Snapshot represents a series of 8-bit memory Values starting at Address.
type Snapshot struct {
	Address uint16
	Values  []uint8
	Values1 []uint8
}

// Import loads memory m with the values in the snapshot.
func Import(m Memory, snapshot Snapshot) {
	for i, value := range snapshot.Values {
		m.Store(snapshot.Address+uint16(i), value)
	}
}

// ImportBinary loads memory with the data starting at addr.
func ImportBinary(m Memory, data []byte, addr uint16) {
	for i, value := range data {
		m.Store(addr+uint16(i), value)
	}
}

// Diff is the difference between two memory values (A and B) at Address.
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

// Verify checks that the values in snapshot b match up with the values in
// memory a. Returns true if all snapshot values match.
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
