package cpu

import (
	"fmt"

	"github.com/blackchip-org/pac8/memory"
)

type Status int

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

type Processor struct {
	CPU     CPU
	Status  Status
	Err     error
	Tracing bool
	mem     memory.Memory
	dasm    *Disassembler
	stop    chan bool
	trace   chan bool
}

func NewProcessor(c CPU) *Processor {
	return &Processor{
		CPU:   c,
		mem:   c.Memory(),
		dasm:  c.Disassembler(),
		stop:  make(chan bool, 1),
		trace: make(chan bool, 1),
	}
}

func (p *Processor) Start() {
	go p.run()
}

func (p *Processor) Stop() {
	p.stop <- true
}

func (p *Processor) Trace(value bool) {
	p.trace <- value
}

func (p *Processor) run() {
	p.Status = Run
	for {
		p.tick()
		select {
		case value := <-p.trace:
			p.Tracing = value
		case <-p.stop:
			p.Status = Halt
			return
		default:
		}
		//time.Sleep(p.CPU.Speed())
	}
}

func (p *Processor) tick() {
	if p.Status == Run {
		if p.Tracing && p.CPU.Ready() {
			p.dasm.SetPC(p.CPU.PC())
			fmt.Println(p.dasm.Next())
		}
		p.CPU.Next()
	}
}
