package memory

// Port represents an input/output port between the CPU and other devices.
// Read points to a value that when reading from the device and Write points
// to a value used when writing to the device. Set the pointers to the
// same value if the port is read/write.
type Port struct {
	Read  *uint8
	Write *uint8
}

type IO interface {
	Memory
	Port(int) *Port
}

type io struct {
	ports []Port
}

func NewIO(n int) IO {
	return &io{
		ports: make([]Port, n, n),
	}
}

func (i *io) Port(p int) *Port {
	return &i.ports[p]
}

func (i *io) Store(address uint16, value uint8) {
	p := i.ports[int(address)]
	if p.Write != nil {
		*p.Write = value
	}
}

func (i *io) Load(address uint16) uint8 {
	p := i.ports[int(address)]
	if p.Read == nil {
		return 0
	}
	return *p.Read
}

func (i *io) Length() int {
	return len(i.ports)
}

type PortMapper struct {
	io IO
}

func NewPortMapper(io IO) PortMapper {
	return PortMapper{io: io}
}

func (m PortMapper) RO(port int, v *uint8) {
	m.io.Port(port).Read = v
}

func (m PortMapper) WO(port int, v *uint8) {
	m.io.Port(port).Write = v
}

func (m PortMapper) RW(port int, v *uint8) {
	m.io.Port(port).Read = v
	m.io.Port(port).Write = v
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
