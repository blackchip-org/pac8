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

type Spy struct {
	callback EventCallback
	mem      Memory
	reads    map[uint16]struct{}
	writes   map[uint16]struct{}
}

func (s *Spy) Load(address uint16) uint8 {
	value := s.mem.Load(address)
	if _, exists := s.reads[address]; exists {
		s.callback(Event{ReadEvent, address, value})
	}
	return value
}

func (s *Spy) Store(address uint16, value uint8) {
	s.mem.Store(address, value)
	if _, exists := s.writes[address]; exists {
		s.callback(Event{WriteEvent, address, value})
	}
}

func (s *Spy) Length() int {
	return s.mem.Length()
}

func (s *Spy) Callback(e EventCallback) {
	s.callback = e
}

func (s *Spy) WatchR(address uint16) {
	s.reads[address] = struct{}{}
}

func (s *Spy) WatchW(address uint16) {
	s.writes[address] = struct{}{}
}

func (s *Spy) WatchRW(address uint16) {
	s.reads[address] = struct{}{}
	s.writes[address] = struct{}{}
}

func (s *Spy) UnwatchR(address uint16) {
	delete(s.reads, address)
}

func (s *Spy) UnwatchW(address uint16) {
	delete(s.writes, address)
}

func (s *Spy) UnwatchRW(address uint16) {
	delete(s.reads, address)
	delete(s.writes, address)
}

func NewSpy(mem Memory) *Spy {
	return &Spy{
		mem:      mem,
		callback: func(e Event) {},
		reads:    make(map[uint16]struct{}),
		writes:   make(map[uint16]struct{}),
	}
}
