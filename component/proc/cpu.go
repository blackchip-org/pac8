package proc

import "github.com/blackchip-org/pac8/component/memory"

type Get func() uint8
type Put func(uint8)

type Get16 func() uint16
type Put16 func(uint16)

type PC interface {
	PC() uint16
	SetPC(uint16)
}

type CPU interface {
	PC
	Next()
	String() string
	Ready() bool
	Info() Info
}

type Value struct {
	Get interface{}
	Put interface{}
}

type Info struct {
	CycleRate       int
	CodeReader      CodeReader
	CodeFormatter   CodeFormatter
	NewDisassembler func(memory.Memory) *Disassembler
	Registers       map[string]Value
}
