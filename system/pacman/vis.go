package pacman

import (
	"log"

	"github.com/blackchip-org/pac8/app"
	"github.com/veandco/go-sdl2/sdl"
)

var visPalette = [4][]uint8{
	[]uint8{0, 0, 0, 0},
	[]uint8{128, 128, 128, 255},
	[]uint8{192, 192, 192, 255},
	[]uint8{255, 255, 255, 255},
}

func PacmanTiles(r *sdl.Renderer) *sdl.Texture {
	mem, err := app.LoadROM("pacman/pacman.5e", "06ef227747a440831c9a3a613b76693d52a2f0a9")
	if err != nil {
		log.Fatal(err)
	}
	tex, err := tileSheet(r, mem, visPalette)
	if err != nil {
		log.Fatal(err)
	}
	return tex
}

func PacmanSprites(r *sdl.Renderer) *sdl.Texture {
	mem, err := app.LoadROM("pacman/pacman.5f", "4a937ac02216ea8c96477d4a15522070507fb599")
	if err != nil {
		log.Fatal(err)
	}
	tex, err := spriteSheet(r, mem, visPalette)
	if err != nil {
		log.Fatal(err)
	}
	return tex
}

func MsPacmanTiles(r *sdl.Renderer) *sdl.Texture {
	mem, err := app.LoadROM("mspacman/5e", "5e8b472b615f12efca3fe792410c23619f067845")
	if err != nil {
		log.Fatal(err)
	}
	tex, err := tileSheet(r, mem, visPalette)
	if err != nil {
		log.Fatal(err)
	}
	return tex
}

func MsPacmanSprites(r *sdl.Renderer) *sdl.Texture {
	mem, err := app.LoadROM("mspacman/5f", "fd6a1dde780b39aea76bf1c4befa5882573c2ef4")
	if err != nil {
		log.Fatal(err)
	}
	tex, err := spriteSheet(r, mem, visPalette)
	if err != nil {
		log.Fatal(err)
	}
	return tex
}
