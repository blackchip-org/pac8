package mach

import (
	"fmt"

	"github.com/blackchip-org/pac8/cpu"
	"github.com/blackchip-org/pac8/memory"
)

type Cab interface {
	Mach() Mach
}

type Mach struct {
	Mem memory.Memory
	CPU cpu.CPU
}

func (m *Mach) Run() {
	dasm := m.NewDisassembler()
	format := m.CPU.CodeFormatter()

	for {
		dasm.SetPC(m.CPU.PC())
		fmt.Println(format(dasm.Next()))
		m.CPU.Next()
	}
}

func (m *Mach) NewDisassembler() *cpu.Disassembler {
	return cpu.NewDisassembler(m.Mem, m.CPU.CodeReader())
}
