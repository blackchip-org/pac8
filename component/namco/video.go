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

func SpriteSheet(r *sdl.Renderer, mem memory.Memory, info SheetInfo) (*sdl.Texture, error) {
	w := info.W * 8
	h := info.H * 8
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
				colors := info.Colors[pixelValue]
				r.SetDrawColorArray(colors...)
				r.DrawPoint(int32(sheetX), int32(sheetY))
			}
		}
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

// http://tech.quarterarcade.com/tech/MAME/src/galaga.c.html.aspx?g=1042

var comment = `
==== Galaga
   335: static struct GfxLayout charlayout =
   336: {
   337: 	8,8,           /* 8*8 characters */
   338: 	RGN_FRAC(1,1), /* 128 characters */
   339: 	2,             /* 2 bits per pixel */
   340: 	{ 0, 4 },       /* the two bitplanes for 4 pixels are packed into one byte */
   341: 	{ 8*8+0, 8*8+1, 8*8+2, 8*8+3, 0, 1, 2, 3 },   /* bits are packed in groups of four */
   342: 	{ 0*8, 1*8, 2*8, 3*8, 4*8, 5*8, 6*8, 7*8 },   /* characters are rotated 90 degrees */
   343: 	16*8           /* every char takes 16 bytes */
   344: };
   345:
   346: static struct GfxLayout spritelayout =
   347: {
   348: 	16,16,          /* 16*16 sprites */
   349: 	128,            /* 128 sprites */
   350: 	2,              /* 2 bits per pixel */
   351: 	{ 0, 4 },       /* the two bitplanes for 4 pixels are packed into one byte */
   352: 	{ 0, 1, 2, 3, 8*8, 8*8+1, 8*8+2, 8*8+3, 16*8+0, 16*8+1, 16*8+2, 16*8+3,
   353: 			24*8+0, 24*8+1, 24*8+2, 24*8+3 },
   354: 	{ 0*8, 1*8, 2*8, 3*8, 4*8, 5*8, 6*8, 7*8,
   355: 			32*8, 33*8, 34*8, 35*8, 36*8, 37*8, 38*8, 39*8 },
   356: 	64*8    /* every sprite takes 64 bytes */
   357: };

==== PacMan

2030: static struct GfxLayout tilelayout =
2031: {
2032: 	8,8,	/* 8*8 characters */
2033:     256,    /* 256 characters */
2034:     2,  /* 2 bits per pixel */
2035:     { 0, 4 },   /* the two bitplanes for 4 pixels are packed into one byte */
2036:     { 8*8+0, 8*8+1, 8*8+2, 8*8+3, 0, 1, 2, 3 }, /* bits are packed in groups of four */
2037:     { 0*8, 1*8, 2*8, 3*8, 4*8, 5*8, 6*8, 7*8 },
2038:     16*8    /* every char takes 16 bytes */
2039: };
2040:
2041:
2042: static struct GfxLayout spritelayout =
2043: {
2044: 	16,16,	/* 16*16 sprites */
2045: 	64,	/* 64 sprites */
2046: 	2,	/* 2 bits per pixel */
2047: 	{ 0, 4 },	/* the two bitplanes for 4 pixels are packed into one byte */
2048: 	{ 8*8, 8*8+1, 8*8+2, 8*8+3, 16*8+0, 16*8+1, 16*8+2, 16*8+3,
2049: 			24*8+0, 24*8+1, 24*8+2, 24*8+3, 0, 1, 2, 3 },
2050: 	{ 0*8, 1*8, 2*8, 3*8, 4*8, 5*8, 6*8, 7*8,
2051: 			32*8, 33*8, 34*8, 35*8, 36*8, 37*8, 38*8, 39*8 },
2052: 	64*8	/* every sprite takes 64 bytes */
2053: };
`
