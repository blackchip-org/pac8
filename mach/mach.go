package mach

import (
	"fmt"
	"os"
	"time"

	"github.com/blackchip-org/pac8/cpu"
	"github.com/blackchip-org/pac8/memory"
	"github.com/veandco/go-sdl2/sdl"
)

type Cab interface {
	Mach() Mach
}

const (
	Init Status = iota
	Halt
	Run
	Break
	Breakpoint
	Trap
)

func (s Status) String() string {
	switch s {
	case Init:
		return "halt"
	case Halt:
		return "halt"
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

type Status int

type Mach struct {
	Mem     memory.Memory
	CPU     cpu.CPU
	Reader  cpu.CodeReader
	Format  cpu.CodeFormatter
	Status  Status
	Tracing bool
	Err     error

	start chan bool
	stop  chan bool
	trace chan bool
}

func New() *Mach {
	return &Mach{
		start: make(chan bool, 1),
		stop:  make(chan bool, 1),
		trace: make(chan bool, 1),
	}
}

func (m *Mach) Run() {
	dasm := m.NewDisassembler()

	lastUpdate := time.Now()
	for {
		if m.Status == Run {
			if m.Tracing {
				dasm.SetPC(m.CPU.PC())
				fmt.Println(m.Format(dasm.Next()))
			}
			m.CPU.Next()
		}
		now := time.Now()
		if now.Sub(lastUpdate) > time.Millisecond {
			lastUpdate = now
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				if _, ok := event.(*sdl.QuitEvent); ok {
					os.Exit(0)
				}
				/*
					for _, input := range m.inputs {
						input.SDLEvent(event)
					}
				*/
			}
		}
		select {
		case <-m.stop:
			m.Status = Halt
		case <-m.start:
			m.Status = Run
		case v := <-m.trace:
			m.Tracing = v
		default:
		}
	}
}

func (m *Mach) NewDisassembler() *cpu.Disassembler {
	return cpu.NewDisassembler(m.Mem, m.Reader)
}

func (m *Mach) Start() {
	m.start <- true
}

func (m *Mach) Stop() {
	m.stop <- true
}

func (m *Mach) Trace(v bool) {
	m.trace <- v
}
