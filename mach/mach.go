package mach

import (
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

	start   chan bool
	stop    chan bool
	render  chan bool
	quit    chan bool
	now     time.Time
	devices []Device
}

func New(proc *cpu.Processor) *Mach {
	return &Mach{
		Proc:   proc,
		start:  make(chan bool, 1),
		stop:   make(chan bool, 1),
		render: make(chan bool, 1),
		quit:   make(chan bool, 1),
	}
}

func (m *Mach) Run() {
	for {
		select {
		case <-m.stop:
			for _, d := range m.devices {
				d.Stop()
			}
		case <-m.start:
			for _, d := range m.devices {
				d.Start()
			}
		case <-m.quit:
			return
		case <-m.render:
			m.Display.Render()
		}
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			if _, ok := event.(*sdl.QuitEvent); ok {
				return
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

func (m *Mach) Now() time.Time {
	return m.now
}

func (m *Mach) Render() {
	m.render <- true
}

func (m *Mach) AddDevice(d Device) {
	m.devices = append(m.devices, d)
}
