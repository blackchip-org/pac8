package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/blackchip-org/pac8/app"
	"github.com/blackchip-org/pac8/machine"
	"github.com/blackchip-org/pac8/pkg/pac8"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	defaultWidth  = 1024
	defaultHeight = 786
)

var (
	gameName      string
	cprof         bool
	monitorEnable bool
	noAudio       bool
	noVideo       bool
	restore       bool
	slowStart     bool
	trace         bool
	wait          bool
)

func init() {
	flag.StringVar(&gameName, "g", "pacman", "use this game")
	flag.BoolVar(&cprof, "cprof", false, "enable cpu profiling")
	flag.BoolVar(&monitorEnable, "m", false, "start monitor")
	flag.BoolVar(&noAudio, "no-audio", false, "disable audio device")
	flag.BoolVar(&noVideo, "no-video", false, "disable video device")
	flag.BoolVar(&restore, "r", false, "restore from previous snapshot")
	flag.BoolVar(&slowStart, "s", false, "slow start -- skip any POST bypass")
	flag.BoolVar(&trace, "t", false, "enable tracing on start")
	flag.BoolVar(&wait, "w", false, "wait for go command")
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

	if noVideo || trace || wait {
		monitorEnable = true
	}

	game, ok := app.Games[gameName]
	if !ok {
		log.Fatalf("no such game: %v", gameName)
	}
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("unable to initialize sdl: %v", err)
	}

	romDir := app.PathFor(app.ROM, gameName)
	roms, err := game.ROM.Load(romDir)
	if err != nil {
		log.Fatalf("unable to load roms\n%v\n", err)
	}

	defer sdl.Quit()

	var env = pac8.Env{}
	if !noVideo {
		fullScreen := uint32(0)
		if !monitorEnable {
			fullScreen = sdl.WINDOW_FULLSCREEN_DESKTOP
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

		r, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
		if err != nil {
			log.Fatalf("unable to initialize renderer: %v", err)
		}
		if err = sdl.GLSetSwapInterval(1); err != nil {
			log.Printf("unable to set swap interval: %v", err)
		}
		env.Renderer = r
	}

	if !noAudio {
		requestSpec := sdl.AudioSpec{
			Freq:     22050,
			Format:   sdl.AUDIO_S16LSB,
			Channels: 2,
			Samples:  367,
		}
		if err := sdl.OpenAudio(&requestSpec, &env.AudioSpec); err != nil {
			log.Fatalf("unable to initialize audio: %v", err)
		}
		sdl.PauseAudio(false)
	}

	runtimeDir := app.PathFor(app.Store, gameName)
	if err := os.MkdirAll(runtimeDir, 0755); err != nil {
		log.Fatalf("unable to create runtime directory %v: %v", runtimeDir, err)
	}

	sys, err := game.Init(env, roms)
	if err != nil {
		log.Fatalf("unable to start game: %v", err)
	}
	m := machine.New(sys)

	if trace {
		m.Send(machine.TraceCmd)
	}

	var mon *app.Monitor
	if monitorEnable {
		mon = app.NewMonitor(m)
		defer func() {
			mon.Close()
		}()
		go func() {
			err := mon.Run()
			if err != nil {
				log.Fatalf("monitor error: %v", err)
			}
		}()
	}
	if !wait {
		m.Send(machine.StartCmd)
	}
	if restore {
		filename := app.PathFor(app.Store, gameName, app.SnapshotFileName)
		m.Send(machine.RestoreCmd, filename)
	} else if !slowStart {
		// If there is a snapshot for bypassing POST, use it
		filename := app.PathFor(app.ROM, gameName, app.InitState)
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			m.Send(machine.RestoreCmd, filename)
		}
	}
	m.Run()
}
