package pacman

import (
	"fmt"

	"github.com/blackchip-org/pac8/memory"
	"github.com/blackchip-org/pac8/util/bits"
	"github.com/veandco/go-sdl2/sdl"
)

// https://www.lomont.org/Software/Games/PacMan/PacmanEmulation.pdf

type Video struct {
	r       *sdl.Renderer
	mem     memory.Memory
	tiles   *sdl.Texture
	sprites *sdl.Texture
	scale   int
}

type VideoROM struct {
	Tiles   memory.Memory
	Sprites memory.Memory
}

func NewVideo(r *sdl.Renderer, mem memory.Memory, rom VideoROM) (*Video, error) {
	v := &Video{
		r:     r,
		scale: 2,
		mem:   mem,
	}
	tiles, err := tileSheet(r, rom.Tiles)
	if err != nil {
		return nil, err
	}
	v.tiles = tiles
	return v, nil
}

func (v *Video) Render() {
	for ty := 0; ty < 36; ty++ {
		for tx := 0; tx < 28; tx++ {
			var addr int
			if ty == 0 || ty == 1 {
				addr = 0x43dd + (ty * 0x20) - tx
			} else if ty == 34 || ty == 35 {
				addr = 0x401d + ((ty - 34) * 0x20) - tx
			} else {
				addr = 0x43a0 + (ty - 2) - (tx * 0x20)
			}
			// fmt.Printf("tx: %v, ty: %v, addr: %02x\n", tx, ty, addr)
			tileN := v.mem.Load(uint16(addr))
			sheetX := (tileN % 16) * 8
			sheetY := (tileN / 16) * 8
			src := sdl.Rect{
				X: int32(sheetX),
				Y: int32(sheetY),
				W: 8,
				H: 8,
			}
			screenX := tx * 8 * v.scale
			screenY := ty * 8 * v.scale
			dest := sdl.Rect{
				X: int32(screenX),
				Y: int32(screenY),
				W: int32(8 * v.scale),
				H: int32(8 * v.scale),
			}
			v.r.Copy(v.tiles, &src, &dest)
		}
	}
	v.r.Present()
}

var palette = [][]uint8{
	[]uint8{0x00, 0x00, 0x00, 0xff},
	[]uint8{0x77, 0x77, 0x77, 0xff},
	[]uint8{0xbb, 0xbb, 0xbb, 0xff},
	[]uint8{0xff, 0xff, 0xff, 0xff},
}

func bit2(b0 bool, b1 bool) int {
	index := 0
	if b0 {
		index |= 0x01
	}
	if b1 {
		index |= 0x02
	}
	return index
}

func tileSheet(r *sdl.Renderer, mem memory.Memory) (*sdl.Texture, error) {
	w := 16 * 8
	h := 16 * 8
	t, err := r.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET, int32(w), int32(h))
	if err != nil {
		return nil, fmt.Errorf("unable to create tile sheet: %v", err)
	}
	r.SetRenderTarget(t)

	for addr := 0; addr < mem.Length(); addr++ {
		tile := addr / 16
		offset := addr % 16
		tileX := (tile % 16) * 8
		tileY := (tile / 16) * 8

		x := tileX + 7 - (offset % 8)
		y := tileY
		if offset < 8 {
			y = tileY + 4
		}
		val := mem.Load(uint16(addr))

		pixel1 := bit2(bits.Get(val, 0), bits.Get(val, 4))
		pixel2 := bit2(bits.Get(val, 1), bits.Get(val, 5))
		pixel3 := bit2(bits.Get(val, 2), bits.Get(val, 6))
		pixel4 := bit2(bits.Get(val, 3), bits.Get(val, 7))

		r.SetDrawColorArray(palette[pixel4]...)
		r.DrawPoint(int32(x), int32(y+0))
		r.SetDrawColorArray(palette[pixel3]...)
		r.DrawPoint(int32(x), int32(y+1))
		r.SetDrawColorArray(palette[pixel2]...)
		r.DrawPoint(int32(x), int32(y+2))
		r.SetDrawColorArray(palette[pixel1]...)
		r.DrawPoint(int32(x), int32(y+3))
	}

	t.SetBlendMode(sdl.BLENDMODE_BLEND)
	r.SetRenderTarget(nil)

	return t, nil
}
