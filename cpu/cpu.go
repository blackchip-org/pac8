package cpu

type Get func() uint8
type Put func(uint8)

type Get16 func() uint16
type Put16 func(uint16)

type CPU interface {
	PC() uint16
	SetPC(uint16)
	Next()
	CodeReader() CodeReader
	CodeFormatter() CodeFormatter
	String() string
}
