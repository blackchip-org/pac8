package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/blackchip-org/pac8/app"
	"github.com/blackchip-org/pac8/machine"
	"github.com/blackchip-org/pac8/system/pacman"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	defaultWidth  = 1024
	defaultHeight = 786
)

var gameName string
var cprof bool
var monitor bool
var noVideo bool
var trace bool
var wait bool

func init() {
	flag.StringVar(&gameName, "g", "pacman", "use this game")
	flag.BoolVar(&cprof, "cprof", false, "enable cpu profiling")
	flag.BoolVar(&monitor, "m", false, "start monitor")
	flag.BoolVar(&noVideo, "no-video", false, "do not show video device")
	flag.BoolVar(&trace, "t", false, "enable tracing on start")
	flag.BoolVar(&wait, "w", false, "wait for go command")
}

var games = map[string]func(*sdl.Renderer) (machine.System, error){
	"pacman":   pacman.NewPacman,
	"mspacman": pacman.NewMsPacman,
}

func main() {
	log.SetFlags(0)
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

	newGame, ok := games[gameName]
	if !ok {
		log.Fatalf("no such game: %v", gameName)
	}

	var r *sdl.Renderer
	if !noVideo {
		if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
			log.Fatalf("unable to initialize sdl: %v", err)
		}
		defer sdl.Quit()

		fullScreen := uint32(0)
		if !monitor {
			fullScreen = sdl.WINDOW_FULLSCREEN
		}
		window, err := sdl.CreateWindow(
			"pac8",
			sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			defaultWidth, defaultHeight,
			sdl.WINDOW_SHOWN|fullScreen,
		)
		if err != nil {
			log.Fatalf("unable to initialize window: %v", err)
		}

		renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
		if err != nil {
			log.Fatalf("unable to initialize renderer: %v", err)
		}
		r = renderer

		err = sdl.GLSetSwapInterval(-1)
		if err != nil {
			log.Printf("no adaptive vsync: %v", err)
			err = sdl.GLSetSwapInterval(1)
			if err != nil {
				log.Printf("unable to set swap interval: %v", err)
			}
		}
	}

	sys, err := newGame(r)
	if err != nil {
		log.Fatalf("unable to start game: %v", err)
	}
	mach := machine.New(sys)

	if trace {
		mach.Trace(log.New(os.Stdout, "", 0))
	}
	if monitor {
		mon := app.NewMonitor(mach)
		go func() {
			err := mon.Run()
			if err != nil {
				log.Fatalf("monitor error: %v", err)
			}
		}()
	}
	if !wait {
		mach.Start()
	}
	mach.Run()
}
