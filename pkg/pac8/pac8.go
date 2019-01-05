package pac8

import (
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/machine"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	ROM    *memory.Pack
	Config interface{}
	Init   func(Env, memory.Set) (machine.System, error)
}

type Env struct {
	Renderer  *sdl.Renderer
	AudioSpec sdl.AudioSpec
}
