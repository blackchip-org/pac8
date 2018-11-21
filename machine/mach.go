package machine

import (
	"log"
	"time"

	"github.com/blackchip-org/pac8/component/input"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/proc"
	"github.com/blackchip-org/pac8/component/video"
	"github.com/veandco/go-sdl2/sdl"
)

type Status int

const (
	Halt Status = iota
	Run
	Breakpoint
	Trap
)

func (s Status) String() string {
	switch s {
	case Halt:
		return "halt"
	case Run:
		return "run"
	case Breakpoint:
		return "break"
	case Trap:
		return "trap"
	}
	return "???"
}

type Spec struct {
	CPU          proc.CPU
	Mem          memory.Memory
	Display      video.Display
	TickCallback func(*Mach)
	TickRate     time.Duration
	CharDecoder  func(uint8) (rune, bool)
}

type System interface {
	Spec() *Spec
}

type Cmd interface{}

type StartCmd struct {
	Cmd
}

type StopCmd struct {
	Cmd
}

type QuitCmd struct {
	Cmd
}

type Mach struct {
	System         System
	CPU            proc.CPU
	Mem            memory.Memory
	Display        video.Display
	In             input.Input
	Status         Status
	Err            error
	Tracer         *log.Logger
	Breakpoints    map[uint16]struct{}
	StatusCallback func(Status)
	TickCallback   func(*Mach)
	CharDecoder    func(uint8) (rune, bool)
	TickRate       time.Duration
	cyclesPerTick  int
	mem            memory.Memory
	dasm           *proc.Disassembler
	cmd            chan Cmd
	trace          chan *log.Logger
	quit           bool
}

func New(sys System) *Mach {
	spec := sys.Spec()
	return &Mach{
		System:         sys,
		CPU:            spec.CPU,
		Breakpoints:    make(map[uint16]struct{}),
		StatusCallback: func(Status) {},
		TickCallback:   spec.TickCallback,
		TickRate:       spec.TickRate,
		Display:        spec.Display,
		CharDecoder:    func(_ uint8) (rune, bool) { return 0, false },
		Mem:            spec.Mem,
		dasm:           spec.CPU.Info().NewDisassembler(spec.Mem),
		cmd:            make(chan Cmd, 1),
		trace:          make(chan *log.Logger, 1),
	}
}

func (m *Mach) setStatus(s Status) {
	m.Status = s
	m.StatusCallback(s)
}

func (m *Mach) Run() {
	m.quit = false
	ticker := time.NewTicker(m.TickRate)
	m.cyclesPerTick = int(float64(m.TickRate) / float64(time.Millisecond) * float64(m.CPU.Info().CycleRate))
	for {
		select {
		case c := <-m.cmd:
			m.command(c)
		case v := <-m.trace:
			m.Tracer = v
		case <-ticker.C:
			m.tick()
		}
		if m.quit {
			return
		}
	}
}

func (m *Mach) tick() {
	for i := 0; i < m.cyclesPerTick; i++ {
		if m.Status == Run {
			if m.Tracer != nil && m.CPU.Ready() {
				m.dasm.SetPC(m.CPU.PC())
				m.Tracer.Print(m.dasm.Next())
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
			m.quit = true
		} else if e, ok := event.(*sdl.KeyboardEvent); ok {
			if e.Keysym.Sym == sdl.K_ESCAPE {
				m.quit = true
			}
		}
		handleKeyboard(event, &m.In)
	}
	m.TickCallback(m)
}

func (m *Mach) Trace(l *log.Logger) {
	m.trace <- l
}

func (m *Mach) Send(c Cmd) {
	m.cmd <- c
}

func (m *Mach) command(c Cmd) {
	switch c.(type) {
	case StartCmd:
		m.setStatus(Run)
	case StopCmd:
		m.setStatus(Halt)
	case QuitCmd:
		m.quit = true
	default:
		log.Panicf("invalid command: %v", c)
	}
}
