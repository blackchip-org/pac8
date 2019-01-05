package machine

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blackchip-org/pac8/component/audio"
	"github.com/blackchip-org/pac8/component/input"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/proc"
	"github.com/blackchip-org/pac8/component/video"
	"github.com/blackchip-org/pac8/pkg/util/state"
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
	CPU          []proc.CPU
	Mem          []memory.Memory
	Display      video.Display
	Audio        audio.Audio
	TickCallback func(*Mach)
	TickRate     time.Duration
	CharDecoder  func(uint8) (rune, bool)
}

type System interface {
	Spec() *Spec
	Save(*state.Encoder)
	Restore(*state.Decoder)
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
	Display       video.Display
	Audio         audio.Audio
	In            input.Input
	Status        Status
	EventCallback func(EventType, interface{})
	TickCallback  func(*Mach)
	CharDecoder   func(uint8) (rune, bool)
	TickRate      time.Duration
	cyclesPerTick int
	Cores         []Core
	cmd           chan Cmd
	tracing       int
	quit          bool
}

type Core struct {
	CPU         proc.CPU
	Mem         memory.Memory
	Breakpoints map[uint16]struct{}
	Dasm        *proc.Disassembler
}

func New(sys System) *Mach {
	spec := sys.Spec()
	nCores := len(spec.CPU)
	m := &Mach{
		System:        sys,
		EventCallback: func(EventType, interface{}) {},
		TickCallback:  spec.TickCallback,
		TickRate:      spec.TickRate,
		Display:       spec.Display,
		CharDecoder:   spec.CharDecoder,
		Audio:         spec.Audio,
		cmd:           make(chan Cmd, 10),
		Cores:         make([]Core, nCores, nCores),
		tracing:       -1,
	}
	for i := 0; i < len(spec.CPU); i++ {
		core := Core{
			CPU:         spec.CPU[i],
			Mem:         spec.Mem[i],
			Breakpoints: make(map[uint16]struct{}),
			Dasm:        spec.CPU[i].Info().NewDisassembler(spec.Mem[i]),
		}
		m.Cores[i] = core
	}
	return m
}

func (m *Mach) setStatus(s Status) {
	m.Status = s
	m.EventCallback(StatusEvent, s)
}

func (m *Mach) Run() {
	m.quit = false
	ticker := time.NewTicker(m.TickRate)
	// FIXME: This needs to be done better with multi-core
	m.cyclesPerTick = int(float64(m.TickRate) / float64(time.Millisecond) * float64(m.Cores[0].CPU.Info().CycleRate))
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
	if m.Status == Run {
		m.execute()
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

func (m *Mach) execute() {
	for i, core := range m.Cores {
		for t := 0; t < m.cyclesPerTick; t++ {
			if m.tracing == i && core.CPU.Ready() {
				core.Dasm.SetPC(core.CPU.PC())
				m.EventCallback(TraceEvent, core.Dasm.Next())
			}
			core.CPU.Next()
			if _, exists := core.Breakpoints[core.CPU.PC()]; exists && core.CPU.Ready() {
				m.setStatus(Break)
				return
			}
		}
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
	enc := state.NewEncoder(out)
	m.System.Save(enc)
	if enc.Err != nil {
		m.EventCallback(ErrorEvent, fmt.Sprintf("unable to save snapshot: %v", enc.Err))
		return
	}
}

func (m *Mach) restore(path string) {
	out, err := os.Open(path)
	if err != nil {
		m.EventCallback(ErrorEvent, fmt.Sprintf("unable to open snapshot: %v", err))
		return
	}
	dec := state.NewDecoder(out)
	m.System.Restore(dec)
	if dec.Err != nil {
		m.EventCallback(ErrorEvent, fmt.Sprintf("unable to load snapshot: %v", dec.Err))
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
		core := c.Args[0].(int)
		if m.tracing >= 0 {
			m.tracing = -1
		} else {
			m.tracing = core
		}
	case QuitCmd:
		m.quit = true
	default:
		log.Panicf("invalid command: %v", c)
	}
}
