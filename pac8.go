package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/blackchip-org/pac8/cabs/pacman"
	"github.com/blackchip-org/pac8/mach"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	defaultWidth  = 1024
	defaultHeight = 786
)

var cprof bool
var monitor bool
var noVideo bool
var trace bool
var wait bool

func init() {
	flag.BoolVar(&cprof, "cprof", false, "enable cpu profiling")
	flag.BoolVar(&monitor, "m", false, "start monitor")
	flag.BoolVar(&noVideo, "no-video", false, "do not show video device")
	flag.BoolVar(&trace, "t", false, "enable tracing on start")
	flag.BoolVar(&wait, "w", false, "wait for go command")
}

func main() {
	flag.Parse()

	if cprof {
		f, err := os.Create("./cpu.prof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		fmt.Println("starting profile")
		defer func() {
			pprof.StopCPUProfile()
			fmt.Println("profile saved")
		}()
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "unable to initialize sdl: %v\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()
	sdl.GLSetSwapInterval(1)

	var m *mach.Mach
	if noVideo {
		m = pacman.New(nil)
	} else {
		window, err := sdl.CreateWindow(
			"pac8",
			sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			defaultWidth, defaultHeight,
			sdl.WINDOW_SHOWN,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to initialize window: %v", err)
			os.Exit(1)
		}

		renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
		if err != nil {
			log.Fatalf("unable to initialize renderer: %v", err)
		}
		m = pacman.New(renderer)
	}

	if trace {
		m.Trace(true)
	}
	if monitor {
		mon := mach.NewMonitor(m)
		go func() {
			err := mon.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "monitor error: %v\n", err)
				os.Exit(1)
			}
		}()
	}
	if !wait {
		m.Start()
	}
	m.Run()
}