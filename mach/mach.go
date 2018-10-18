package mach

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/blackchip-org/pac8/cpu"
	"github.com/blackchip-org/pac8/memory"
	"github.com/blackchip-org/pac8/pac8"
	"github.com/veandco/go-sdl2/sdl"
)

type Status int

const (
	Stop Status = iota
	Run
	Break
	Breakpoint
	Trap
)

func (s Status) String() string {
	switch s {
	case Stop:
		return "stop"
	case Run:
		return "run"
	case Break:
		return "break"
	case Breakpoint:
		return "breakpoint"
	case Trap:
		return "trap"
	}
	return "???"
}

type Display interface {
	Render()
}

type UI interface {
	SDLEvent(sdl.Event)
}

type Mach struct {
	CPU           cpu.CPU
	Mem           memory.Memory
	Display       Display
	UI            UI
	Status        Status
	Err           error
	Tracing       bool
	Breakpoints   map[uint16]struct{}
	Callback      func(Status)
	CharDecoder   CharDecoder
	TickRate      time.Duration
	cyclesPerTick int
	mem           memory.Memory
	dasm          *cpu.Disassembler
	start         chan bool
	stop          chan bool
	trace         chan bool
	quit          chan bool
}

func New(mem memory.Memory, cpu cpu.CPU) *Mach {
	return &Mach{
		CPU:         cpu,
		Breakpoints: make(map[uint16]struct{}),
		Callback:    func(Status) {},
		Mem:         mem,
		dasm:        cpu.Disassembler(),
		start:       make(chan bool, 1),
		stop:        make(chan bool, 1),
		trace:       make(chan bool, 1),
		quit:        make(chan bool, 1),
	}
}

func (m *Mach) setStatus(s Status) {
	m.Status = s
	m.Callback(s)
}

func (m *Mach) Run() {
	ticker := time.NewTicker(m.TickRate)
	m.cyclesPerTick = int(float64(m.TickRate) / float64(time.Millisecond) * float64(m.CPU.CycleRate()))
	for {
		select {
		case <-m.stop:
			m.setStatus(Stop)
		case <-m.start:
			m.setStatus(Run)
		case v := <-m.trace:
			m.Tracing = v
		case <-m.quit:
			return
		case <-ticker.C:
			m.tick()
		}
	}
}

func (m *Mach) tick() {
	for i := 0; i < m.cyclesPerTick; i++ {
		if m.Status == Run {
			if m.Tracing && m.CPU.Ready() {
				m.dasm.SetPC(m.CPU.PC())
				fmt.Println(m.dasm.Next())
			}
			m.CPU.Next()
			if _, exists := m.Breakpoints[m.CPU.PC()]; exists && m.CPU.Ready() {
				m.setStatus(Breakpoint)
				return
			}
		}
	}
	m.Display.Render()
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		if _, ok := event.(*sdl.QuitEvent); ok {
			m.quit <- true
		}
		if m.UI != nil {
			m.UI.SDLEvent(event)
		}
	}
}

func (m *Mach) Start() {
	m.start <- true
}

func (m *Mach) Stop() {
	m.stop <- true
}

func (m *Mach) Quit() {
	m.quit <- true
}

func (m *Mach) Trace(v bool) {
	m.trace <- v
}

func LoadROM(e *[]error, path string, checksum string) memory.Memory {
	romDir := filepath.Join(pac8.Home(), "rom")
	filename := filepath.Join(romDir, path)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		*e = append(*e, err)
		return nil
	}
	rom := memory.NewROM(data)
	romChecksum := fmt.Sprintf("%04x", sha1.Sum(data))
	if checksum != romChecksum {
		*e = append(*e, fmt.Errorf("invalid checksum for file: %s\nexpected: %v\nreceived: %v", filename, checksum, romChecksum))
	}
	return rom
}

func Dump(m memory.Memory, start uint16, end uint16, decode CharDecoder) string {
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
