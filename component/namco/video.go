package namco

import (
	"fmt"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/veandco/go-sdl2/sdl"
)

type Palette [4][]uint8

var VisPalette = [4][]uint8{
	[]uint8{0, 0, 0, 0},
	[]uint8{128, 128, 128, 255},
	[]uint8{192, 192, 192, 255},
	[]uint8{255, 255, 255, 255},
}

type SheetInfo struct {
	W      int
	H      int
	Colors Palette
}

func TileSheet(r *sdl.Renderer, mem memory.Memory, info SheetInfo) (*sdl.Texture, error) {
	w := info.W * 8
	h := info.H * 8
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

		r.SetDrawColorArray(info.Colors[pixel4]...)
		r.DrawPoint(int32(x), int32(y+0))
		r.SetDrawColorArray(info.Colors[pixel3]...)
		r.DrawPoint(int32(x), int32(y+1))
		r.SetDrawColorArray(info.Colors[pixel2]...)
		r.DrawPoint(int32(x), int32(y+2))
		r.SetDrawColorArray(info.Colors[pixel1]...)
		r.DrawPoint(int32(x), int32(y+3))
	}

	t.SetBlendMode(sdl.BLENDMODE_BLEND)
	r.SetRenderTarget(nil)

	return t, nil
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
