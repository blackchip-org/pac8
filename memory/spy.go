package memory

import "fmt"

type EventType int
type EventCallback func(e Event)

const (
	ReadEvent EventType = iota
	WriteEvent
)

func (t EventType) String() string {
	switch t {
	case ReadEvent:
		return "r"
	case WriteEvent:
		return "w"
	default:
		return "?"
	}
}

type Event struct {
	Type    EventType
	Address uint16
	Value   uint8
}

func (e Event) String() string {
	return fmt.Sprintf("%v %04x %02x", e.Type, e.Address, e.Value)
}

type spy struct {
	callback EventCallback
	mem      Memory
	reads    map[uint16]struct{}
	writes   map[uint16]struct{}
}

func (s *spy) Load(address uint16) uint8 {
	value := s.mem.Load(address)
	if _, exists := s.reads[address]; exists {
		s.callback(Event{ReadEvent, address, value})
	}
	return value
}

func (s *spy) Store(address uint16, value uint8) {
	s.mem.Store(address, value)
	if _, exists := s.writes[address]; exists {
		s.callback(Event{WriteEvent, address, value})
	}
}

func (s *spy) Length() int {
	return s.mem.Length()
}

func (s *spy) Callback(e EventCallback) {
	s.callback = e
}

func (s *spy) WatchR(address uint16) {
	s.reads[address] = struct{}{}
}

func (s *spy) WatchW(address uint16) {
	s.writes[address] = struct{}{}
}

func (s *spy) WatchRW(address uint16) {
	s.reads[address] = struct{}{}
	s.writes[address] = struct{}{}
}

func (s *spy) UnwatchR(address uint16) {
	delete(s.reads, address)
}

func (s *spy) UnwatchW(address uint16) {
	delete(s.writes, address)
}

func (s *spy) UnwatchRW(address uint16) {
	delete(s.reads, address)
	delete(s.writes, address)
}

type Spy struct {
	spy
}

func NewSpy(mem Memory) *Spy {
	return &Spy{
		spy: spy{
			mem:      mem,
			callback: func(e Event) {},
			reads:    make(map[uint16]struct{}),
			writes:   make(map[uint16]struct{}),
		},
	}
}

type SpyIO struct {
	spy
	io IO
}

func NewSpyIO(io IO) *SpyIO {
	return &SpyIO{
		io: io,
		spy: spy{
			mem:      io,
			callback: func(e Event) {},
			reads:    make(map[uint16]struct{}),
			writes:   make(map[uint16]struct{}),
		},
	}
}

func (s *SpyIO) Port(p int) *Port {
	return s.io.Port(p)
}
