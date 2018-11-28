package machine

import (
	"encoding/gob"
	"log"
	"os"
	"time"

	"github.com/blackchip-org/pac8/component"
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
	Save(component.Encoder) error
	Restore(component.Decoder) error
}

type CmdType int

const (
	RestoreCmd CmdType = iota
	SaveCmd
	StartCmd
	StopCmd
	TraceCmd
	QuitCmd
)

type Cmd struct {
	Type CmdType
	Args []interface{}
}

type EventType int

const (
	StatusEvent EventType = iota
	TraceEvent
)

type Mach struct {
	System        System
	CPU           proc.CPU
	Mem           memory.Memory
	Display       video.Display
	In            input.Input
	Status        Status
	Err           error
	Breakpoints   map[uint16]struct{}
	EventCallback func(EventType, interface{})
	TickCallback  func(*Mach)
	CharDecoder   func(uint8) (rune, bool)
	TickRate      time.Duration
	cyclesPerTick int
	mem           memory.Memory
	dasm          *proc.Disassembler
	cmd           chan Cmd
	tracing       bool
	quit          bool
}

func New(sys System) *Mach {
	spec := sys.Spec()
	return &Mach{
		System:        sys,
		CPU:           spec.CPU,
		Breakpoints:   make(map[uint16]struct{}),
		EventCallback: func(EventType, interface{}) {},
		TickCallback:  spec.TickCallback,
		TickRate:      spec.TickRate,
		Display:       spec.Display,
		CharDecoder:   func(_ uint8) (rune, bool) { return 0, false },
		Mem:           spec.Mem,
		dasm:          spec.CPU.Info().NewDisassembler(spec.Mem),
		cmd:           make(chan Cmd, 10),
	}
}

func (m *Mach) setStatus(s Status) {
	m.Status = s
	m.EventCallback(StatusEvent, s)
}

func (m *Mach) Run() {
	m.quit = false
	ticker := time.NewTicker(m.TickRate)
	m.cyclesPerTick = int(float64(m.TickRate) / float64(time.Millisecond) * float64(m.CPU.Info().CycleRate))
	for {
		select {
		case c := <-m.cmd:
			m.command(c)
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
			if m.tracing && m.CPU.Ready() {
				m.dasm.SetPC(m.CPU.PC())
				m.EventCallback(TraceEvent, m.dasm.Next())
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

func (m *Mach) Send(t CmdType, args ...interface{}) {
	m.cmd <- Cmd{Type: t, Args: args}
}

func (m *Mach) save(path string) {
	out, err := os.Create(path)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	enc := gob.NewEncoder(out)
	if err := m.System.Save(enc); err != nil {
		log.Printf("error: %v\n", err)
		return
	}
}

func (m *Mach) restore(path string) {
	out, err := os.Open(path)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	dec := gob.NewDecoder(out)
	if err := m.System.Restore(dec); err != nil {
		log.Printf("error: %v\n", err)
		return
	}
}

func (m *Mach) command(c Cmd) {
	switch c.Type {
	case RestoreCmd:
		path := c.Args[0].(string)
		m.restore(path)
	case SaveCmd:
		path := c.Args[0].(string)
		m.save(path)
	case StartCmd:
		m.setStatus(Run)
	case StopCmd:
		m.setStatus(Halt)
	case TraceCmd:
		m.tracing = !m.tracing
	case QuitCmd:
		m.quit = true
	default:
		log.Panicf("invalid command: %v", c)
	}
}