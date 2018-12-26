package namco

import (
	"fmt"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/video"
	"github.com/veandco/go-sdl2/sdl"
)

type Palette [4][]uint8

const (
	w = int32(224)
	h = int32(288)
)

type SpriteCoord struct {
	X uint8
	Y uint8
}

var VisPalette = [4][]uint8{
	[]uint8{0, 0, 0, 0},
	[]uint8{128, 128, 128, 255},
	[]uint8{192, 192, 192, 255},
	[]uint8{255, 255, 255, 255},
}

type SheetLayout struct {
	W            int
	H            int
	CellW        int
	CellH        int
	BytesPerCell int
	PixelLayout  [][]int
	PixelReader  func(memory.Memory, uint16, int) uint8
}

type Sheet struct {
	W       int
	H       int
	CellW   int
	CellH   int
	Texture *sdl.Texture
}

func NewSheet(r *sdl.Renderer, mem memory.Memory, l SheetLayout, pal Palette) (*Sheet, error) {
	t, err := r.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET, int32(l.W), int32(l.H))
	if err != nil {
		return nil, fmt.Errorf("unable to create tile sheet: %v", err)
	}
	r.SetRenderTarget(t)

	rowTiles := l.W / l.CellW
	for i := 0; i < l.W*l.H; i++ {
		targetX := i % l.W
		targetY := i / l.W
		cellX := targetX / l.CellW
		offsetX := targetX % l.CellW
		cellY := targetY / l.CellH
		offsetY := targetY % l.CellH
		cellN := cellX + cellY*rowTiles
		baseAddr := uint16(cellN * l.BytesPerCell)
		pixelN := l.PixelLayout[offsetY][offsetX]
		value := l.PixelReader(mem, baseAddr, pixelN)

		r.SetDrawColorArray(pal[value]...)
		r.DrawPoint(int32(targetX), int32(targetY))
	}
	t.SetBlendMode(sdl.BLENDMODE_BLEND)
	r.SetRenderTarget(nil)

	return &Sheet{
		W:       l.W,
		H:       l.H,
		CellW:   l.CellW,
		CellH:   l.CellH,
		Texture: t,
	}, nil
}

type Video struct {
	Callback     func()
	SpriteCoords [8]SpriteCoord
	r            *sdl.Renderer
	mem          memory.Memory
	layouts      Layouts
	tiles        [64]*Sheet
	sprites      [64]*Sheet
	colors       [16][]uint8
	palettes     [64]Palette
	frame        video.RenderFrame
	frameFill    sdl.Rect
	scanLines    *sdl.Texture
}

type Layouts struct {
	Tile   SheetLayout
	Sprite SheetLayout
}

type VideoROM struct {
	Tiles   memory.Memory
	Sprites memory.Memory
	Color   memory.Memory
	Palette memory.Memory
}

func NewVideo(r *sdl.Renderer, mem memory.Memory, rom VideoROM, layouts Layouts) (*Video, error) {
	v := &Video{
		r:       r,
		mem:     mem,
		layouts: layouts,
	}
	if r == nil {
		return v, nil
	}

	winW, winH, err := r.GetOutputSize()
	if err != nil {
		return nil, err
	}
	v.frame = video.FitInWindow(winW, winH, w, h)
	v.frameFill = sdl.Rect{
		X: v.frame.X,
		Y: v.frame.Y,
		W: v.frame.W,
		H: v.frame.H,
	}
	v.scanLines, err = video.ScanLines(r, winW, winH, v.frame.Scale-1)
	if err != nil {
		return nil, err
	}
	v.colorTable(rom.Color)
	v.paletteTable(rom.Palette)

	for pal := 0; pal < 64; pal++ {
		tiles, err := NewSheet(r, rom.Tiles, layouts.Tile, v.palettes[pal])
		if err != nil {
			return nil, err
		}
		v.tiles[pal] = tiles

		sprites, err := NewSheet(r, rom.Sprites, layouts.Sprite, v.palettes[pal])
		if err != nil {
			return nil, err
		}
		v.sprites[pal] = sprites
	}

	return v, nil
}

func (v *Video) Render() {
	//v.Callback()
	if v.r == nil {
		return
	}
	v.r.SetDrawColorArray(0, 0, 0, 0xff)
	v.r.FillRect(&v.frameFill)
	v.renderTiles()
	v.renderSprites()
	v.r.Copy(v.scanLines, nil, nil)
	v.r.Present()
}

func (v *Video) renderTiles() {
	cellW := v.layouts.Tile.CellW
	rowCells := v.layouts.Tile.W / cellW

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

			tileN := int(v.mem.Load(uint16(addr)))
			sheetX := (tileN % rowCells) * cellW
			sheetY := (tileN / rowCells) * cellW
			src := sdl.Rect{
				X: int32(sheetX),
				Y: int32(sheetY),
				W: int32(v.layouts.Tile.CellW),
				H: int32(v.layouts.Tile.CellH),
			}
			screenX := tx * 8 * v.frame.Scale
			screenY := ty * 8 * v.frame.Scale
			dest := sdl.Rect{
				X: screenX + v.frame.X,
				Y: screenY + v.frame.Y,
				W: int32(v.layouts.Tile.CellW) * v.frame.Scale,
				H: int32(v.layouts.Tile.CellH) * v.frame.Scale,
			}

			caddr := addr + 0x0400
			// Only 64 palettes, strip out the higher bits
			pal := v.mem.Load(uint16(caddr)) & 0x3f
			v.r.Copy(v.tiles[pal].Texture, &src, &dest)
		}
	}
}

func (v *Video) renderSprites() {
	spriteW := int32(v.layouts.Sprite.CellW)
	spriteH := int32(v.layouts.Sprite.CellH)
	rowCells := int32(v.layouts.Sprite.W) / spriteW

	for s := 7; s >= 0; s-- {
		coordX := int32(v.SpriteCoords[s].X)
		coordY := int32(v.SpriteCoords[s].Y)
		info := v.mem.Load(uint16(0x4ff0 + (s * 2)))
		spriteN := int32(info >> 2)
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
		screenX := (w - coordX + spriteW) * v.frame.Scale
		screenY := (h - coordY - spriteH) * v.frame.Scale
		sheetX := (spriteN % rowCells) * spriteW
		sheetY := (spriteN / rowCells) * spriteH
		src := sdl.Rect{
			X: int32(sheetX),
			Y: int32(sheetY),
			W: spriteW,
			H: spriteH,
		}
		dest := sdl.Rect{
			X: screenX + v.frame.X,
			Y: screenY + v.frame.Y,
			W: spriteW * v.frame.Scale,
			H: spriteH * v.frame.Scale,
		}
		pal := v.mem.Load(uint16(0x4ff1 + (s * 2)))
		v.r.CopyEx(v.sprites[pal].Texture, &src, &dest, 0, nil, flip)
	}
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