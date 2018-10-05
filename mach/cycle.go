package mach

import (
	"time"

	"github.com/blackchip-org/pac8/util/clock"
)

type Cycle struct {
	t0       time.Time
	interval time.Duration
}

func NewCycle(interval time.Duration) *Cycle {
	return &Cycle{
		t0:       clock.Now(),
		interval: interval,
	}
}

func (c *Cycle) Next() bool {
	now := clock.Now()
	if now.Sub(c.t0) < c.interval {
		return false
	}
	c.t0 = now
	return true
}
