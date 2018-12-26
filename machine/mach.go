package machine

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blackchip-org/pac8/component"
	"github.com/blackchip-org/pac8/component/audio"
	"github.com/blackchip-org/pac8/component/input"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/proc"
	"github.com/blackchip-org/pac8/component/video"
	"github.com/veandco/go-sdl2/sdl"
)

type Status int

var measureTime bool

const (
	Halt Status = iota
	Run
	Break
)

func (s Status) String() string {
	switch s {
	case Halt:
		return "halt"
	case Run:
		return "run"
	case Break:
		return "break"
	}
	return "???"
}

type Spec struct {
	Name         string
	CPU          proc.CPU
	Mem          memory.Memory
	Display      video.Display
	Audio        audio.Audio
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
	ErrorEvent
)

type Mach struct {
	System        System
	CPU           proc.CPU
	Mem           memory.Memory
	Display       video.Display
	Audio         audio.Audio
	In            input.Input
	Status        Status
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
		CharDecoder:   spec.CharDecoder,
		Audio:         spec.Audio,
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
				m.setStatus(Break)
				return
			}
		}
	}
	if m.Display != nil {
		m.Display.Render()
	}
	if m.Audio != nil && m.Status == Run {
		if err := m.Audio.Queue(); err != nil {
			log.Panicf("unable to queue audio: %v", err)
		}
	}
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
	if m.TickCallback != nil {
		m.TickCallback(m)
	}
}

func (m *Mach) Send(t CmdType, args ...interface{}) {
	m.cmd <- Cmd{Type: t, Args: args}
}

func (m *Mach) save(path string) {
	out, err := os.Create(path)
	if err != nil {
		m.EventCallback(ErrorEvent, fmt.Sprintf("unable to create snapshot: %v", err))
		return
	}
	enc := gob.NewEncoder(out)
	if err := m.System.Save(enc); err != nil {
		m.EventCallback(ErrorEvent, fmt.Sprintf("unable to save snapshot: %v", err))
		return
	}
}

func (m *Mach) restore(path string) {
	out, err := os.Open(path)
	if err != nil {
		m.EventCallback(ErrorEvent, fmt.Sprintf("unable to open snapshot: %v", err))
		return
	}
	dec := gob.NewDecoder(out)
	if err := m.System.Restore(dec); err != nil {
		m.EventCallback(ErrorEvent, fmt.Sprintf("unable to load snapshot: %v", err))
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
