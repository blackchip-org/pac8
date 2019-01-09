package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/blackchip-org/pac8/app"
	"github.com/blackchip-org/pac8/game"
	"github.com/blackchip-org/pac8/pkg/memory"
	"github.com/blackchip-org/pac8/pkg/namco"
	"github.com/blackchip-org/pac8/pkg/video"
	"github.com/blackchip-org/pac8/system/galaga"
	"github.com/blackchip-org/pac8/system/pacman"
	"github.com/veandco/go-sdl2/sdl"
)

type view struct {
	system string
	roms   *memory.Pack
	render func(*sdl.Renderer, memory.Set) (video.Sheet, error)
}

var views = map[string]view{
	"galaga:sprites": view{
		system: "galaga",
		roms:   game.Galaga.ROM,
		render: func(r *sdl.Renderer, roms memory.Set) (video.Sheet, error) {
			return namco.NewSheet(r,
				roms["sprite"],
				galaga.VideoConfig.SpriteLayout,
				namco.ViewerPalette)
		},
	},
	"galaga:tiles": view{
		system: "galaga",
		roms:   game.Galaga.ROM,
		render: func(r *sdl.Renderer, roms memory.Set) (video.Sheet, error) {
			return namco.NewSheet(r,
				roms["tile"],
				galaga.VideoConfig.TileLayout,
				namco.ViewerPalette)
		},
	},
	"mspacman:colors": view{
		system: "mspacman",
		roms:   game.MsPacMan.ROM,
		render: func(r *sdl.Renderer, roms memory.Set) (video.Sheet, error) {
			config := pacman.VideoConfig
			colors := namco.ColorTable(roms["color"], config)
			return video.NewColorSheet(r, []video.Palette{colors})
		},
	},
	"mspacman:palettes": view{
		system: "mspacman",
		roms:   game.MsPacMan.ROM,
		render: func(r *sdl.Renderer, roms memory.Set) (video.Sheet, error) {
			config := pacman.VideoConfig
			colors := namco.ColorTable(roms["color"], config)
			palettes := namco.PaletteTable(roms["palette"], config, colors)
			return video.NewColorSheet(r, palettes)
		},
	},
	"mspacman:sprites": view{
		system: "mspacman",
		roms:   game.MsPacMan.ROM,
		render: func(r *sdl.Renderer, roms memory.Set) (video.Sheet, error) {
			return namco.NewSheet(r,
				roms["sprite"],
				pacman.VideoConfig.SpriteLayout,
				namco.ViewerPalette)
		},
	},
	"mspacman:tile": view{
		system: "mspacman",
		roms:   game.MsPacMan.ROM,
		render: func(r *sdl.Renderer, roms memory.Set) (video.Sheet, error) {
			return namco.NewSheet(r,
				roms["tile"],
				pacman.VideoConfig.TileLayout,
				namco.ViewerPalette)
		},
	},
	"pacman:colors": view{
		system: "pacman",
		roms:   game.PacMan.ROM,
		render: func(r *sdl.Renderer, roms memory.Set) (video.Sheet, error) {
			config := pacman.VideoConfig
			colors := namco.ColorTable(roms["color"], config)
			return video.NewColorSheet(r, []video.Palette{colors})
		},
	},
	"pacman:palettes": view{
		system: "pacman",
		roms:   game.PacMan.ROM,
		render: func(r *sdl.Renderer, roms memory.Set) (video.Sheet, error) {
			config := pacman.VideoConfig
			colors := namco.ColorTable(roms["color"], config)
			palettes := namco.PaletteTable(roms["palette"], config, colors)
			return video.NewColorSheet(r, palettes)
		},
	},
	"pacman:sprites": view{
		system: "pacman",
		roms:   game.PacMan.ROM,
		render: func(r *sdl.Renderer, roms memory.Set) (video.Sheet, error) {
			return namco.NewSheet(r,
				roms["sprite"],
				pacman.VideoConfig.SpriteLayout,
				namco.ViewerPalette)
		},
	},
	"pacman:tile": view{
		system: "pacman",
		roms:   game.PacMan.ROM,
		render: func(r *sdl.Renderer, roms memory.Set) (video.Sheet, error) {
			return namco.NewSheet(r,
				roms["tile"],
				pacman.VideoConfig.TileLayout,
				namco.ViewerPalette)
		},
	},
}

var scale int
var vscan int

func init() {
	flag.IntVar(&scale, "scale", 1, "image `scale`")
	flag.IntVar(&vscan, "vscan", vscan, "add vertical scan line of this `width`")

	flag.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "Usage: pac8-viewer [options] <view>\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(o, "\nAvailable values for <view>:\n\n")
		list := []string{}
		for key := range views {
			list = append(list, key)
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

	v, ok := views[flag.Arg(0)]
	if !ok {
		log.Fatalln("no such view")
	}

	roms, err := v.roms.Load(app.PathFor(app.ROM, v.system))
	if err != nil {
		log.Fatalf("unable to load roms:\n%v\n", err)
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "unable to initialize sdl: %v\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		flag.Arg(0),
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		100, 100,
		sdl.WINDOW_HIDDEN,
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

	sheet, err := v.render(r, roms)
	if err != nil {
		log.Fatalf("unable to create sheet: %v", err)
	}
	_, _, w, h, err := sheet.Texture.Query()
	winX, winY := w*int32(scale), h*int32(scale)
	window.SetSize(winX, winY)
	window.SetPosition(sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED)
	window.Show()

	sheet, _ = v.render(r, roms)
	var scanlines *sdl.Texture
	if vscan > 0 {
		scanlines, err = video.ScanLines(r, winX, winY, int32(vscan))
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			if _, ok := event.(*sdl.QuitEvent); ok {
				os.Exit(0)
			}
		}

		r.SetRenderTarget(nil)
		r.Clear()
		r.Copy(sheet.Texture, nil, nil)
		if scanlines != nil {
			r.Copy(scanlines, nil, nil)
		}
		sdl.Delay(250)
		r.Present()
	}
}
