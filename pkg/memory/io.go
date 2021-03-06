package memory

import (
	"fmt"

	"github.com/blackchip-org/pac8/pkg/util/state"
)

// Port represents an input/output port between the CPU and other devices.
// Read points to a value used when reading from the device and Write points
// to a value used when writing to the device. Set the pointers to the
// same value if the port is read/write.
type Port struct {
	Read  *uint8
	Write *uint8
}

func (p Port) String() string {
	r := uint8(0)
	if p.Read != nil {
		r = *p.Read
	}
	w := uint8(0)
	if p.Write != nil {
		w = *p.Write
	}
	return fmt.Sprintf("r(%02x) w(%02x)", r, w)
}

// IO is memory that is mapped to input/output ports. Use PortMapper to
// easily map memory addresses to port values.
type IO interface {
	Memory
	Port(int) *Port
}

type io struct {
	ports []Port
}

// NewIO creates a new input/output memory with n ports.
func NewIO(n int) IO {
	return io{
		ports: make([]Port, n, n),
	}
}

// Store writes the value to the port mapped at address. If no port is
// mapped at address, this function does nothing.
func (i io) Store(address uint16, value uint8) {
	p := i.ports[int(address)]
	if p.Write != nil {
		*p.Write = value
	}
}

// Load reads the value from the port mapped at address. If no port is
// mapped at address, zero is returned.
func (i io) Load(address uint16) uint8 {
	p := i.ports[int(address)]
	if p.Read == nil {
		return 0
	}
	return *p.Read
}

// Length returns the number of ports.
func (i io) Length() int {
	return len(i.ports)
}

func (i io) Save(enc *state.Encoder) {
	for _, p := range i.ports {
		if p.Read != nil {
			enc.Encode(*p.Read)
		}
		if p.Write != nil {
			enc.Encode(*p.Write)
		}
	}
}

func (i io) Restore(dec *state.Decoder) {
	for _, p := range i.ports {
		if p.Read != nil {
			dec.Decode(p.Read)
		}
		if p.Write != nil {
			dec.Decode(p.Write)
		}
	}
}

// Port returns the read/write values for port n.
func (i io) Port(n int) *Port {
	return &i.ports[n]
}

// PortMapper maps memory addresses to port values.
type PortMapper struct {
	io IO
}

// NewPortMapper creates a new mapper for memory io.
func NewPortMapper(io IO) PortMapper {
	return PortMapper{io: io}
}

// RO maps the value at v to the port at p only when reading.
func (m PortMapper) RO(p int, v *uint8) {
	m.io.Port(p).Read = v
}

// WO maps the value at v to the port at p only when writing.
func (m PortMapper) WO(p int, v *uint8) {
	m.io.Port(p).Write = v
}

// RW maps the value at v to the port at p when reading or writing.
func (m PortMapper) RW(p int, v *uint8) {
	m.io.Port(p).Read = v
	m.io.Port(p).Write = v
}
