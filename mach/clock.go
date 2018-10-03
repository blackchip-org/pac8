package mach

import (
	"time"
)

type Clock struct {
	duration time.Duration
	ticker   *time.Ticker
	callback func()
	stop     chan bool
}

func NewClock(d time.Duration, callback func()) *Clock {
	return &Clock{
		duration: d,
		callback: callback,
		stop:     make(chan bool, 1),
	}
}

func (c *Clock) Start() {
	c.ticker = time.NewTicker(c.duration)
	go c.run()
}

func (c *Clock) Stop() {
	c.stop <- true
}

func (c *Clock) run() {
	for {
		select {
		case <-c.ticker.C:
			c.callback()
		case <-c.stop:
			return
		}
	}
}
