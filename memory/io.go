package memory

// Port represents an input/output port between the CPU and other devices.
// Read points to a value that when reading from the device and Write points
// to a value used when writing to the device. Set the pointers to the
// same value if the port is read/write.
type Port struct {
	Read  *uint8
	Write *uint8
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

// NewIO creates new input/output memory with n ports.
func NewIO(n int) IO {
	return &io{
		ports: make([]Port, n, n),
	}
}

// Port returns the port with port number p.
func (i *io) Port(p int) *Port {
	return &i.ports[p]
}

// Store writes the value to the port mapped at address. If no port is
// mapped at address, this function does nothing.
func (i *io) Store(address uint16, value uint8) {
	p := i.ports[int(address)]
	if p.Write != nil {
		*p.Write = value
	}
}

// Load reads the value from the port mapped at address. If no port is
// mapped at address, zero is returned.
func (i *io) Load(address uint16) uint8 {
	p := i.ports[int(address)]
	if p.Read == nil {
		return 0
	}
	return *p.Read
}

// Length returns the number of ports.
func (i *io) Length() int {
	return len(i.ports)
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

type mockIO struct {
	data map[uint8][]uint8
}

func NewMockIO(snapshots []Snapshot) IO {
	mio := &mockIO{
		data: make(map[uint8][]uint8),
	}
	for _, snapshot := range snapshots {
		addr := uint8(snapshot.Address)
		stack, exists := mio.data[addr]
		if !exists {
			stack = make([]uint8, 0, 0)
		}
		stack = append(stack, snapshot.Values[0])
		mio.data[addr] = stack
	}
	return mio
}

func (m *mockIO) Load(addr uint16) uint8 {
	stack, exists := m.data[uint8(addr)]
	if !exists {
		return 0
	}
	if len(stack) == 0 {
		return 0
	}
	v := stack[0]
	stack = stack[1:]
	m.data[uint8(addr)] = stack
	return v
}

func (m *mockIO) Store(addr uint16, value uint8) {
	stack, exists := m.data[uint8(addr)]
	if !exists {
		stack = make([]uint8, 0, 0)
	}
	stack = append(stack, value)
	m.data[uint8(addr)] = stack
}

func (m *mockIO) Length() int {
	return 0
}

func (m *mockIO) Port(p int) *Port {
	return nil
}
