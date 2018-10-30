package pacman

import (
	"github.com/blackchip-org/pac8/bits"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	JoystickUp    = 0
	JoystickLeft  = 1
	JoystickRight = 2
	JoystickDown  = 3
	Coin1         = 5
	Coin2         = 6
)

const (
	Start1 = 5
	Start2 = 6
)

type Keyboard struct {
	regs *Registers
}

func NewKeyboard(regs *Registers) Keyboard {
	return Keyboard{regs: regs}
}

func (k Keyboard) SDLEvent(event sdl.Event) {
	e, ok := event.(*sdl.KeyboardEvent)
	if !ok {
		return
	}
	if e.Type == sdl.KEYDOWN {
		switch e.Keysym.Sym {
		case sdl.K_1:
			bits.Set(&k.regs.In1, Start1, false)
		case sdl.K_2:
			bits.Set(&k.regs.In1, Start2, false)
		case sdl.K_c:
			bits.Set(&k.regs.In0, Coin1, true)
		case sdl.K_UP:
			bits.Set(&k.regs.In0, JoystickUp, false)
		case sdl.K_DOWN:
			bits.Set(&k.regs.In0, JoystickDown, false)
		case sdl.K_LEFT:
			bits.Set(&k.regs.In0, JoystickLeft, false)
		case sdl.K_RIGHT:
			bits.Set(&k.regs.In0, JoystickRight, false)
		}
	}
	if e.Type == sdl.KEYUP {
		switch e.Keysym.Sym {
		case sdl.K_1:
			bits.Set(&k.regs.In1, Start1, true)
		case sdl.K_2:
			bits.Set(&k.regs.In1, Start2, true)
		case sdl.K_c:
			bits.Set(&k.regs.In0, Coin1, false)
		case sdl.K_UP:
			bits.Set(&k.regs.In0, JoystickUp, true)
		case sdl.K_DOWN:
			bits.Set(&k.regs.In0, JoystickDown, true)
		case sdl.K_LEFT:
			bits.Set(&k.regs.In0, JoystickLeft, true)
		case sdl.K_RIGHT:
			bits.Set(&k.regs.In0, JoystickRight, true)
		}
	}
}
