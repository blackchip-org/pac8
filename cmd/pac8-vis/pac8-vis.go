package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/blackchip-org/pac8/cabinet/pacman"
	"github.com/blackchip-org/pac8/mach"
	"github.com/veandco/go-sdl2/sdl"
)

var vis = map[string]func(*sdl.Renderer) *sdl.Texture{
	"mspacman:tiles":   pacman.MsPacmanTiles,
	"mspacman:sprites": pacman.MsPacmanSprites,
	"pacman:tiles":     pacman.PacmanTiles,
	"pacman:sprites":   pacman.PacmanSprites,
}

var scale int
var vscan int

func init() {
	flag.IntVar(&scale, "scale", 1, "image `scale`")
	flag.IntVar(&vscan, "vscan", vscan, "add vertical scan line of this `width`")

	flag.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "Usage: pac8-vis [options] <vis>\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(o, "\nAvailable values for <vis>:\n\n")
		list := []string{}
		for name := range vis {
			list = append(list, name)
		}
		sort.Strings(list)
		fmt.Fprintln(o, strings.Join(list, "\n"))
		fmt.Fprintln(o)
	}
}

func main() {
	log.SetFlags(0)

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	v, ok := vis[flag.Arg(0)]
	if !ok {
		log.Fatalln("no such visualization")
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "unable to initialize sdl: %v\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"pac8-viz",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		100, 100,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to initialize window: %v", err)
		os.Exit(1)
	}

	r, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("unable to initialize renderer: %v", err)
	}

	err = sdl.GLSetSwapInterval(1)
	if err != nil {
		fmt.Printf("unable to set swap interval: %v\n", err)
	}

	sheet := v(r)
	_, _, w, h, err := sheet.Query()
	winX, winY := w*int32(scale), h*int32(scale)
	window.SetSize(winX, winY)
	window.SetPosition(sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED)

	var scanlines *sdl.Texture
	if vscan > 0 {
		scanlines, err = mach.ScanLines(r, winX, winY, int32(vscan))
		if err != nil {
			log.Fatal(err)
		}
	}

	r.SetRenderTarget(nil)
	r.Copy(sheet, nil, nil)
	if scanlines != nil {
		r.Copy(scanlines, nil, nil)
	}
	r.Present()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			if _, ok := event.(*sdl.QuitEvent); ok {
				os.Exit(0)
			}
		}
	}
}
