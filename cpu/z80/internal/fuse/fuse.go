package fuse

import "github.com/blackchip-org/pac8/memory"

//go:generate go run gen.go
//go:generate go fmt in.go
//go:generate go fmt expected.go

type Test struct {
	Name    string
	AF      uint16
	BC      uint16
	DE      uint16
	HL      uint16
	AF1     uint16
	BC1     uint16
	DE1     uint16
	HL1     uint16
	IX      uint16
	IY      uint16
	SP      uint16
	PC      uint16
	I       uint8
	R       uint8
	IFF1    int
	IFF2    int
	IM      int
	Halt    int
	TStates int

	Snapshots  []memory.Snapshot
	PortReads  []memory.Snapshot
	PortWrites []memory.Snapshot
}
