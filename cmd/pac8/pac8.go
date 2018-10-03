package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/blackchip-org/pac8/cabs"
	"github.com/blackchip-org/pac8/mach"
	"github.com/veandco/go-sdl2/sdl"
)

var monitor bool
var trace bool
var wait bool

func init() {
	flag.BoolVar(&monitor, "m", false, "start monitor")
	flag.BoolVar(&trace, "t", false, "enable tracing on start")
	flag.BoolVar(&wait, "w", false, "wait for go command")
}

func main() {
	flag.Parse()
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "unable to initialize sdl: %v\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()
	sdl.GLSetSwapInterval(1)

	m := cabs.NewPacman().Mach()
	if trace {
		m.Proc.Trace(true)
	}
	if monitor {
		mon := mach.NewMonitor(m)
		go func() {
			err := mon.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to start monitor: %v\n", err)
				os.Exit(1)
			}
		}()
	}
	if !wait {
		m.Start()
	}
	m.Run()
}
