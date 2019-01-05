package machine

import (
	"github.com/blackchip-org/pac8/pkg/input"
	"github.com/veandco/go-sdl2/sdl"
)

func handleKeyboard(event sdl.Event, in *input.Input) {
	e, ok := event.(*sdl.KeyboardEvent)
	if !ok {
		return
	}

	var state bool
	switch e.Type {
	case sdl.KEYDOWN:
		state = true
	case sdl.KEYUP:
		state = false
	default:
		return
	}

	switch e.Keysym.Sym {
	case sdl.K_1:
		in.PlayerStart[0].Active = state
	case sdl.K_2:
		in.PlayerStart[0].Active = state
	case sdl.K_c:
		if state {
			in.CoinSlot[0].Deposit()
		}
	case sdl.K_UP:
		in.Joysticks[0].Up = state
	case sdl.K_DOWN:
		in.Joysticks[0].Down = state
	case sdl.K_LEFT:
		in.Joysticks[0].Left = state
	case sdl.K_RIGHT:
		in.Joysticks[0].Right = state
	}
}
