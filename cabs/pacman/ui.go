package pacman

import (
	"github.com/blackchip-org/pac8/util/bits"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	JoystickUp    = 0
	JoystickLeft  = 1
	JoystickRight = 2
	JoystickDown  = 3
	Coin1         = 5
	Con2          = 6
)

const (
	Start1 = 5
	Start2 = 6
)

type Keyboard struct {
	reg *registers
}

func NewKeyboard(reg *registers) Keyboard {
	return Keyboard{reg: reg}
}

func (k Keyboard) SDLEvent(event sdl.Event) {
	e, ok := event.(*sdl.KeyboardEvent)
	if !ok {
		return
	}
	if e.Type == sdl.KEYDOWN {
		switch e.Keysym.Sym {
		case sdl.K_1:
			bits.Set(&k.reg.in1, Start1, false)
		case sdl.K_5:
			bits.Set(&k.reg.in0, Coin1, true)
		case sdl.K_UP:
			bits.Set(&k.reg.in0, JoystickUp, false)
		case sdl.K_DOWN:
			bits.Set(&k.reg.in0, JoystickDown, false)
		case sdl.K_LEFT:
			bits.Set(&k.reg.in0, JoystickLeft, false)
		case sdl.K_RIGHT:
			bits.Set(&k.reg.in0, JoystickRight, false)
		}
	}
	if e.Type == sdl.KEYUP {
		switch e.Keysym.Sym {
		case sdl.K_1:
			bits.Set(&k.reg.in1, Start1, true)
		case sdl.K_5:
			bits.Set(&k.reg.in0, Coin1, false)
		case sdl.K_UP:
			bits.Set(&k.reg.in0, JoystickUp, true)
		case sdl.K_DOWN:
			bits.Set(&k.reg.in0, JoystickDown, true)
		case sdl.K_LEFT:
			bits.Set(&k.reg.in0, JoystickLeft, true)
		case sdl.K_RIGHT:
			bits.Set(&k.reg.in0, JoystickRight, true)
		}
	}
}
