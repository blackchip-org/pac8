package mach

import (
	"os"
	"time"

	"github.com/blackchip-org/pac8/cpu"

	"github.com/veandco/go-sdl2/sdl"
)

type Cab interface {
	Mach() Mach
}

type Device interface {
	Start()
	Stop()
}

type Display interface {
	Render()
}

type Mach struct {
	Proc    *cpu.Processor
	Display Display

	start  chan bool
	stop   chan bool
	render chan bool

	now     time.Time
	devices []Device
}

func New(proc *cpu.Processor) *Mach {
	return &Mach{
		Proc:   proc,
		start:  make(chan bool, 1),
		stop:   make(chan bool, 1),
		render: make(chan bool, 1),
	}
}

func (m *Mach) Run() {
	lastUpdate := m.now
	for {
		m.now = time.Now()
		/*
			if m.Status == Run {
				if m.Tracing && m.CPU.Ready() {
					dasm.SetPC(m.CPU.PC())
					fmt.Println(m.Format(dasm.Next()))
				}
				m.CPU.Next()
			}
		*/
		if m.now.Sub(lastUpdate) > time.Millisecond {
			lastUpdate = m.now
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
			for _, d := range m.devices {
				d.Stop()
			}
		case <-m.start:
			for _, d := range m.devices {
				d.Start()
			}
		case <-m.render:
			m.Display.Render()
		default:
		}
	}
}

func (m *Mach) Start() {
	m.start <- true
}

func (m *Mach) Stop() {
	m.stop <- true
}

func (m *Mach) Now() time.Time {
	return m.now
}

func (m *Mach) Render() {
	m.render <- true
}

func (m *Mach) AddDevice(d Device) {
	m.devices = append(m.devices, d)
}
