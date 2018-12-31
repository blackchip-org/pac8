package memory

import (
	"fmt"

	"github.com/blackchip-org/pac8/component"
)

type EventType int

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

type EventCallback func(e Event)

type Event struct {
	Type    EventType // Read or Write
	Address uint16    // Address that was accessed
	Value   uint8     // Value read or written
}

func (e Event) String() string {
	return fmt.Sprintf("%v %04x %02x", e.Type, e.Address, e.Value)
}

// Spy watches for when certain memory addresses are accessed.
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

// Callback sets the function that is called when a watched address
// is accessed.
func (s *Spy) Callback(e EventCallback) {
	s.callback = e
}

// WatchR watches addr and invokes the callback when read.
func (s *Spy) WatchR(addr uint16) {
	s.reads[addr] = struct{}{}
}

// WatchW watches addr and invokes the callback when written.
func (s *Spy) WatchW(address uint16) {
	s.writes[address] = struct{}{}
}

// WatchRW watches addr and invokes the callback when read or written.
func (s *Spy) WatchRW(address uint16) {
	s.reads[address] = struct{}{}
	s.writes[address] = struct{}{}
}

// UnwatchR removes a previous read watch on address.
func (s *Spy) UnwatchR(address uint16) {
	delete(s.reads, address)
}

// UnwatchW removes a previous write watch on address.
func (s *Spy) UnwatchW(address uint16) {
	delete(s.writes, address)
}

// UnwatchRW removes previous read and write watches on address.
func (s *Spy) UnwatchRW(address uint16) {
	delete(s.reads, address)
	delete(s.writes, address)
}

func (s Spy) Save(enc component.Encoder) error {
	return s.mem.Save(enc)
}

func (s Spy) Restore(dec component.Decoder) error {
	return s.mem.Restore(dec)
}

// Spy creates a new memory spy on mem.
func NewSpy(mem Memory) *Spy {
	return &Spy{
		mem:      mem,
		callback: func(e Event) {},
		reads:    make(map[uint16]struct{}),
		writes:   make(map[uint16]struct{}),
	}
}
