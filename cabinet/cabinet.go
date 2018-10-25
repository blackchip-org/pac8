package cabinet

import (
	"errors"

	"github.com/blackchip-org/pac8/cabinet/pacman"
	"github.com/blackchip-org/pac8/mach"
	"github.com/veandco/go-sdl2/sdl"
)

var cabs = map[string]func(*sdl.Renderer) (*mach.Mach, error){
	"pacman":   pacman.NewPacman,
	"mspacman": pacman.NewMsPacman,
}

func New(name string, r *sdl.Renderer) (*mach.Mach, error) {
	cab, ok := cabs[name]
	if !ok {
		return nil, errors.New("no such cabinet")
	}
	return cab(r)
}
