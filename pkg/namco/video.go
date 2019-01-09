package namco

import (
	"fmt"

	"github.com/blackchip-org/pac8/pkg/memory"
	"github.com/blackchip-org/pac8/pkg/util/bits"
	"github.com/blackchip-org/pac8/pkg/video"
	"github.com/veandco/go-sdl2/sdl"
)

type Palette [][]uint8

const (
	w = int32(224)
	h = int32(288)
)

type SpriteCoord struct {
	X uint8
	Y uint8
}

var ViewerPalette = [][]uint8{
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

func NewSheet(r *sdl.Renderer, mem memory.Memory, l SheetLayout, pal Palette) (video.Sheet, error) {
	t, err := r.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET, int32(l.W), int32(l.H))
	if err != nil {
		return video.Sheet{}, fmt.Errorf("unable to create tile sheet: %v", err)
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

	return video.Sheet{
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
	config       Config
	tiles        [64]video.Sheet
	sprites      [64]video.Sheet
	colors       [16][]uint8
	palettes     []Palette
	frame        video.RenderFrame
	frameFill    sdl.Rect
	scanLines    *sdl.Texture
}

type Config struct {
	TileLayout     SheetLayout
	SpriteLayout   SheetLayout
	VideoAddr      uint16
	PaletteEntries int
	PaletteColors  int
	Hack           bool
}

func NewVideo(r *sdl.Renderer, mem memory.Memory, rom memory.Set, config Config) (*Video, error) {
	v := &Video{
		r:        r,
		mem:      mem,
		config:   config,
		palettes: make([]Palette, config.PaletteEntries, config.PaletteEntries),
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
	v.colorTable(rom["color"])
	v.paletteTable(rom["palette"])

	for pal := 0; pal < v.config.PaletteEntries; pal++ {
		tiles, err := NewSheet(r, rom["tile"], config.TileLayout, v.palettes[pal])
		if err != nil {
			return nil, err
		}
		v.tiles[pal] = tiles

		sprites, err := NewSheet(r, rom["sprite"], config.SpriteLayout, v.palettes[pal])
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
	layout := v.config.TileLayout
	cellW := layout.CellW
	rowCells := layout.W / cellW

	// Render tiles
	for ty := uint16(0); ty < 36; ty++ {
		for tx := uint16(0); tx < 28; tx++ {
			var addr uint16
			if ty == 0 || ty == 1 {
				addr = v.addr(0x3dd) + (ty * 0x20) - tx
			} else if ty == 34 || ty == 35 {
				addr = v.addr(0x01d) + ((ty - 34) * 0x20) - tx
			} else {
				addr = v.addr(0x3a0) + (ty - 2) - (tx * 0x20)
			}

			tileN := int(v.mem.Load(addr))
			sheetX := (tileN % rowCells) * cellW
			sheetY := (tileN / rowCells) * cellW
			src := sdl.Rect{
				X: int32(sheetX),
				Y: int32(sheetY),
				W: int32(layout.CellW),
				H: int32(layout.CellH),
			}
			screenX := int32(tx) * 8 * v.frame.Scale
			screenY := int32(ty) * 8 * v.frame.Scale
			dest := sdl.Rect{
				X: screenX + v.frame.X,
				Y: screenY + v.frame.Y,
				W: int32(layout.CellW) * v.frame.Scale,
				H: int32(layout.CellH) * v.frame.Scale,
			}

			caddr := addr + 0x0400
			// Only 64 palettes, strip out the higher bits
			// pal := v.mem.Load(caddr) & 0x3f
			pal := v.mem.Load(caddr) & 0x1f
			v.r.Copy(v.tiles[pal].Texture, &src, &dest)
		}
	}
}

func (v *Video) renderSprites() {
	// FIXME: Galaga testing
	if v.config.Hack {
		return
	}
	layout := v.config.SpriteLayout
	spriteW := int32(layout.CellW)
	spriteH := int32(layout.CellH)
	rowCells := int32(layout.W) / spriteW

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
	// FIXME: Galaga testing
	if v.config.Hack {
		return
	}
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
	for pal := 0; pal < v.config.PaletteEntries; pal++ {
		// FIXME: Galaga testing
		if v.config.Hack {
			v.palettes[pal] = ViewerPalette
			continue
		}
		addr := pal * v.config.PaletteColors
		entry := make([][]uint8, v.config.PaletteColors, v.config.PaletteColors)
		for i := 0; i < v.config.PaletteColors; i++ {
			ref := mem.Load(uint16(addr+i)) & 0x0f
			entry[i] = v.colors[ref]
		}
		v.palettes[pal] = entry
	}
}

func (v *Video) addr(offset int) uint16 {
	return v.config.VideoAddr + uint16(offset)
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
