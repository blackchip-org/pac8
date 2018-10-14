package memory

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/blackchip-org/pac8/util/bits"
	"github.com/blackchip-org/pac8/util/charset"
)

func LoadLE(m Memory, address uint16) uint16 {
	lo := m.Load(address)
	hi := m.Load(address + 1)
	return bits.Join(hi, lo)
}

func StoreLE(m Memory, address uint16, value uint16) {
	hi, lo := bits.Split(value)
	m.Store(address, lo)
	m.Store(address+1, hi)
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

func Dump(m Memory, start uint16, end uint16, decode charset.Decoder) string {
	var buf bytes.Buffer
	var chars bytes.Buffer

	a0 := start / 0x10 * 0x10
	a1 := end / 0x10 * 0x10
	if a1 != end {
		a1 += 0x10
	}
	for addr := a0; addr < a1; addr++ {
		if addr%0x10 == 0 {
			buf.WriteString(fmt.Sprintf("$%04x", addr))
			chars.Reset()
		}
		if addr < start || addr > end {
			buf.WriteString("   ")
			chars.WriteString(" ")
		} else {
			value := m.Load(addr)
			buf.WriteString(fmt.Sprintf(" %02x", value))
			ch, printable := decode(value)
			if printable {
				chars.WriteString(fmt.Sprintf("%c", ch))
			} else {
				chars.WriteString(".")
			}
		}
		if addr%0x10 == 7 {
			buf.WriteString(" ")
		}
		if addr%0x10 == 0x0f {
			buf.WriteString(" " + chars.String())
			if addr < end-1 {
				buf.WriteString("\n")
			}
		}
	}
	return buf.String()
}

func home() string {
	home := os.Getenv("PAC8_HOME")
	if home == "" {
		home = "."
	}
	return home
}

func LoadROM(e *[]error, path string, checksum string) *ROM {
	filename := filepath.Join(home(), path)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		*e = append(*e, err)
		return nil
	}
	rom := NewROM(data)
	romChecksum := rom.Checksum()
	if checksum != romChecksum {
		*e = append(*e, fmt.Errorf("invalid checksum for file: %s\nexpected: %v\nreceived: %v", filename, romChecksum, checksum))
	}
	return rom
}
