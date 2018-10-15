package mach

import (
	"fmt"
	"time"

	"github.com/blackchip-org/pac8/cpu"
	"github.com/blackchip-org/pac8/memory"
	"github.com/blackchip-org/pac8/util/clock"
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
	CPU         cpu.CPU
	Display     Display
	UI          UI
	Status      Status
	Err         error
	Tracing     bool
	Breakpoints map[uint16]struct{}
	Callback    func(Status)

	mem    memory.Memory
	dasm   *cpu.Disassembler
	events *Cycle
	start  chan bool
	stop   chan bool
	trace  chan bool
	quit   chan bool
}

func New(cpu cpu.CPU) *Mach {
	return &Mach{
		CPU:         cpu,
		Breakpoints: make(map[uint16]struct{}),
		Callback:    func(Status) {},
		mem:         cpu.Memory(),
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
	m.events = NewCycle(1 * time.Millisecond)
	for {
		clock.SetNow(time.Now())
		select {
		case <-m.stop:
			m.setStatus(Stop)
		case <-m.start:
			m.setStatus(Run)
		case v := <-m.trace:
			m.Tracing = v
		case <-m.quit:
			return
		default:
			m.tick()
		}
	}
}

func (m *Mach) tick() {
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
	m.Display.Render()
	if m.events.Next() {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			if _, ok := event.(*sdl.QuitEvent); ok {
				m.quit <- true
			}
			if m.UI != nil {
				m.UI.SDLEvent(event)
			}
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
