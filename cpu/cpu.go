package cpu

import (
	"time"

	"github.com/blackchip-org/pac8/memory"
)

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
	Memory() memory.Memory
	Next()
	String() string
	Ready() bool
	Speed() time.Duration
	Disassembler() *Disassembler
}
