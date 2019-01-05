package input

import (
	"time"
)

type Joystick struct {
	Up    bool
	Down  bool
	Right bool
	Left  bool
}

type Coin struct {
	Active bool
}

func (c *Coin) Deposit() {
	c.Active = true
	time.AfterFunc(50*time.Millisecond, func() {
		c.Active = false
	})
}

type Button struct {
	Active bool
}

type Input struct {
	Joysticks   [4]Joystick
	CoinSlot    [2]Coin
	PlayerStart [2]Button
}
