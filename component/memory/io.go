package memory

// Port represents an input/output port between the CPU and other devices.
// Read points to a value used when reading from the device and Write points
// to a value used when writing to the device. Set the pointers to the
// same value if the port is read/write.
type Port struct {
	Read  *uint8
	Write *uint8
}

// IO is memory that is mapped to input/output ports. Use PortMapper to
// easily map memory addresses to port values.
type IO struct {
	Ports []Port
}

// NewIO creates a new input/output memory with n ports.
func NewIO(n int) IO {
	return IO{
		Ports: make([]Port, n, n),
	}
}

// Store writes the value to the port mapped at address. If no port is
// mapped at address, this function does nothing.
func (i IO) Store(address uint16, value uint8) {
	p := i.Ports[int(address)]
	if p.Write != nil {
		*p.Write = value
	}
}

// Load reads the value from the port mapped at address. If no port is
// mapped at address, zero is returned.
func (i IO) Load(address uint16) uint8 {
	p := i.Ports[int(address)]
	if p.Read == nil {
		return 0
	}
	return *p.Read
}

// Length returns the number of ports.
func (i IO) Length() int {
	return len(i.Ports)
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
	m.io.Ports[p].Read = v
}

// WO maps the value at v to the port at p only when writing.
func (m PortMapper) WO(p int, v *uint8) {
	m.io.Ports[p].Write = v
}

// RW maps the value at v to the port at p when reading or writing.
func (m PortMapper) RW(p int, v *uint8) {
	m.io.Ports[p].Read = v
	m.io.Ports[p].Write = v
}
