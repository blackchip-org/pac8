package pacman

import (
	"fmt"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/mach"
	"github.com/blackchip-org/pac8/memory"
	"github.com/veandco/go-sdl2/sdl"
)

// https://www.lomont.org/Software/Games/PacMan/PacmanEmulation.pdf

const (
	w = int32(224)
	h = int32(288)
)

type spriteCoord struct {
	x uint8
	y uint8
}

type Palette [4][]uint8

type Video struct {
	Callback     func()
	r            *sdl.Renderer
	mem          memory.Memory
	tiles        [64]*sdl.Texture
	sprites      [64]*sdl.Texture
	colors       [16][]uint8
	palettes     [64]Palette
	spriteCoords [8]spriteCoord
	frame        mach.RenderFrame
	frameFill    sdl.Rect
	scanLines    *sdl.Texture
}

type VideoROM struct {
	Tiles   memory.Memory
	Sprites memory.Memory
	Color   memory.Memory
	Palette memory.Memory
}

func NewVideo(r *sdl.Renderer, mem memory.Memory, rom VideoROM) (*Video, error) {
	v := &Video{
		r:   r,
		mem: mem,
	}
	if r == nil {
		return v, nil
	}

	winW, winH, err := r.GetOutputSize()
	if err != nil {
		return nil, err
	}
	v.frame = mach.FitInWindow(winW, winH, w, h)
	v.frameFill = sdl.Rect{
		X: v.frame.X,
		Y: v.frame.Y,
		W: v.frame.W,
		H: v.frame.H,
	}
	v.scanLines, err = mach.ScanLines(r, winW, winH, v.frame.Scale)
	if err != nil {
		return nil, err
	}
	v.colorTable(rom.Color)
	v.paletteTable(rom.Palette)

	for pal := 0; pal < 64; pal++ {
		tiles, err := tileSheet(r, rom.Tiles, v.palettes[pal])
		if err != nil {
			return nil, err
		}
		v.tiles[pal] = tiles

		sprites, err := spriteSheet(r, rom.Sprites, v.palettes[pal])
		if err != nil {
			return nil, err
		}
		v.sprites[pal] = sprites
	}

	return v, nil
}

func (v *Video) Render() {
	v.Callback()
	if v.r == nil {
		return
	}

	v.r.SetDrawColorArray(0, 0, 0, 0xff)
	v.r.FillRect(&v.frameFill)

	// Render tiles
	for ty := int32(0); ty < 36; ty++ {
		for tx := int32(0); tx < 28; tx++ {
			var addr int32
			if ty == 0 || ty == 1 {
				addr = 0x43dd + (ty * 0x20) - tx
			} else if ty == 34 || ty == 35 {
				addr = 0x401d + ((ty - 34) * 0x20) - tx
			} else {
				addr = 0x43a0 + (ty - 2) - (tx * 0x20)
			}

			tileN := v.mem.Load(uint16(addr))
			sheetX := (tileN % 16) * 8
			sheetY := (tileN / 16) * 8
			src := sdl.Rect{
				X: int32(sheetX),
				Y: int32(sheetY),
				W: 8,
				H: 8,
			}
			screenX := tx * 8 * v.frame.Scale
			screenY := ty * 8 * v.frame.Scale
			dest := sdl.Rect{
				X: screenX + v.frame.X,
				Y: screenY + v.frame.Y,
				W: 8 * v.frame.Scale,
				H: 8 * v.frame.Scale,
			}

			caddr := addr + 0x0400
			// Only 64 palettes, strip out the higher bits
			pal := v.mem.Load(uint16(caddr)) & 0x3f
			v.r.Copy(v.tiles[pal], &src, &dest)
		}
	}

	// Render sprites, reverse order
	for s := 7; s >= 0; s-- {
		coordX := int32(v.spriteCoords[s].x)
		coordY := int32(v.spriteCoords[s].y)
		info := v.mem.Load(uint16(0x4ff0 + (s * 2)))
		spriteN := info >> 2
		flip := sdl.FLIP_NONE
		if info&0x02 > 0 {
			flip |= sdl.FLIP_HORIZONTAL
		}
		if info&0x01 > 0 {
			flip |= sdl.FLIP_VERTICAL
		}

		// do not render of off screen
		if coordX <= 30 || coordX >= 240 {
			continue
		}
		screenX := (w - coordX + 16) * v.frame.Scale
		screenY := (h - coordY - 16) * v.frame.Scale
		sheetX := (spriteN % 8) * 16
		sheetY := (spriteN / 8) * 16
		src := sdl.Rect{
			X: int32(sheetX),
			Y: int32(sheetY),
			W: 16,
			H: 16,
		}
		dest := sdl.Rect{
			X: screenX + v.frame.X,
			Y: screenY + v.frame.Y,
			W: 16 * v.frame.Scale,
			H: 16 * v.frame.Scale,
		}
		pal := v.mem.Load(uint16(0x4ff1 + (s * 2)))
		v.r.CopyEx(v.sprites[pal], &src, &dest, 0, nil, flip)
	}
	v.r.Copy(v.scanLines, nil, nil)
	v.r.Present()
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

func tileSheet(r *sdl.Renderer, mem memory.Memory, pal Palette) (*sdl.Texture, error) {
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

		r.SetDrawColorArray(pal[pixel4]...)
		r.DrawPoint(int32(x), int32(y+0))
		r.SetDrawColorArray(pal[pixel3]...)
		r.DrawPoint(int32(x), int32(y+1))
		r.SetDrawColorArray(pal[pixel2]...)
		r.DrawPoint(int32(x), int32(y+2))
		r.SetDrawColorArray(pal[pixel1]...)
		r.DrawPoint(int32(x), int32(y+3))
	}

	t.SetBlendMode(sdl.BLENDMODE_BLEND)
	r.SetRenderTarget(nil)

	return t, nil
}

func spriteSheet(r *sdl.Renderer, mem memory.Memory, pal Palette) (*sdl.Texture, error) {
	// 64 sprites to be placed in 8x8 matrix each with 16x16 pixels
	w, h := 8*16, 8*16
	t, err := r.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET, int32(w), int32(h))
	if err != nil {
		return nil, fmt.Errorf("unable to create sprite sheet: %v", err)
	}
	r.SetRenderTarget(t)

	for spriteN := 0; spriteN < 64; spriteN++ {
		// In the 8x8 matrix, the cell that contains spriteN
		cellX := spriteN % 8
		cellY := spriteN / 8

		// Upper left pixel on the sprite sheet for this sprite
		sheetULX := cellX * 16
		sheetULY := cellY * 16

		// Each sprite uses 64 bytes of ROM
		baseAddr := spriteN * 64

		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				sheetX := sheetULX + x
				sheetY := sheetULY + y

				// The pixel that is being drawn
				pixelN := spritePixels[y][x]

				// Each byte represents 4 pixels
				byteN := pixelN / 4

				// Not sure why this has to be inverted
				bitOffset := 3 - (pixelN % 4)

				v := mem.Load(uint16(baseAddr + byteN))
				pixelValue := bit2(bits.Get(v, 0+bitOffset), bits.Get(v, 4+bitOffset))
				colors := pal[pixelValue]
				r.SetDrawColorArray(colors...)
				r.DrawPoint(int32(sheetX), int32(sheetY))
			}
		}
	}

	t.SetBlendMode(sdl.BLENDMODE_BLEND)
	r.SetRenderTarget(nil)

	return t, nil
}

func (v *Video) colorTable(mem memory.Memory) {
	for addr := 0; addr < 16; addr++ {
		r, g, b := uint8(0), uint8(0), uint8(0)
		c := mem.Load(uint16(addr))
		for bit := 0; bit < 8; bit++ {
			if bits.Get(c, bit) {
				r += colorWeights[bit][0]
				g += colorWeights[bit][1]
				b += colorWeights[bit][2]
			}
		}
		alpha := uint8(0xff)
		// Color 0 is actually transparent
		if addr == 0 {
			alpha = 0x00
		}
		v.colors[addr] = []uint8{r, g, b, alpha}
	}
}

func (v *Video) paletteTable(mem memory.Memory) {
	for pal := 0; pal < 64; pal++ {
		addr := pal * 4
		var entry [4][]uint8
		entry[0] = v.colors[mem.Load(uint16(addr+0))]
		entry[1] = v.colors[mem.Load(uint16(addr+1))]
		entry[2] = v.colors[mem.Load(uint16(addr+2))]
		entry[3] = v.colors[mem.Load(uint16(addr+3))]
		v.palettes[pal] = entry
	}
}

var spritePixels = [][]int{

	[]int{188, 184, 180, 176, 172, 168, 164, 160, 60, 56, 52, 48, 44, 40, 36, 32},
	[]int{189, 185, 181, 177, 173, 169, 165, 161, 61, 57, 53, 49, 45, 41, 37, 33},
	[]int{190, 186, 182, 178, 174, 170, 166, 162, 62, 58, 54, 50, 46, 42, 38, 34},
	[]int{191, 187, 183, 179, 175, 171, 167, 163, 63, 59, 55, 51, 47, 43, 39, 35},

	[]int{220, 216, 212, 208, 204, 200, 196, 192, 92, 88, 84, 80, 76, 72, 68, 64},
	[]int{221, 217, 213, 209, 205, 201, 197, 193, 93, 89, 85, 81, 77, 73, 69, 65},
	[]int{222, 218, 214, 210, 206, 202, 198, 194, 94, 90, 86, 82, 78, 74, 70, 66},
	[]int{223, 219, 215, 211, 207, 203, 199, 195, 95, 91, 87, 83, 79, 75, 71, 67},

	[]int{252, 248, 244, 240, 236, 232, 228, 224, 124, 120, 116, 112, 108, 104, 100, 96},
	[]int{253, 249, 245, 241, 237, 233, 229, 225, 125, 121, 117, 113, 109, 105, 101, 97},
	[]int{254, 250, 246, 242, 238, 234, 230, 226, 126, 122, 118, 114, 110, 106, 102, 98},
	[]int{255, 251, 247, 243, 239, 235, 231, 227, 127, 123, 119, 115, 111, 107, 103, 99},

	[]int{156, 152, 148, 144, 140, 136, 132, 128, 28, 24, 20, 16, 12, 8, 4, 0},
	[]int{157, 153, 149, 145, 141, 137, 133, 129, 29, 25, 21, 17, 13, 9, 5, 1},
	[]int{158, 154, 150, 146, 142, 138, 134, 130, 30, 26, 22, 18, 14, 10, 6, 2},
	[]int{159, 155, 151, 147, 143, 139, 135, 131, 31, 27, 23, 19, 15, 11, 7, 3},
}

var colorWeights = [][]uint8{
	[]uint8{0x21, 0x00, 0x00},
	[]uint8{0x47, 0x00, 0x00},
	[]uint8{0x97, 0x00, 0x00},
	[]uint8{0x00, 0x21, 0x00},
	[]uint8{0x00, 0x47, 0x00},
	[]uint8{0x00, 0x97, 0x00},
	[]uint8{0x00, 0x00, 0x51},
	[]uint8{0x00, 0x00, 0xae},
}
