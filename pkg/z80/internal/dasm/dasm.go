package dasm

//go:generate go run gen.go
//go:generate go fmt ../../dasm.go
//go:generate go fmt harston.go

type Test struct {
	Name  string
	Op    string
	Bytes []uint8
}
