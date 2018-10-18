package cabs

import (
	"errors"

	"github.com/blackchip-org/pac8/cabs/pacman"
	"github.com/blackchip-org/pac8/mach"
	"github.com/veandco/go-sdl2/sdl"
)

var cabs = map[string]func(*sdl.Renderer) *mach.Mach{
	"pacman": pacman.New,
}

func New(name string, r *sdl.Renderer) (*mach.Mach, error) {
	cab, ok := cabs[name]
	if !ok {
		return nil, errors.New("no such cabinet")
	}
	return cab(r), nil
}
