package galaga

import (
	"github.com/blackchip-org/pac8/pkg/memory"
	"github.com/blackchip-org/pac8/pkg/namco"
	"github.com/blackchip-org/pac8/pkg/util/bits"
	"github.com/veandco/go-sdl2/sdl"
)

func NewVideo(r *sdl.Renderer, mem memory.Memory, rom memory.Set) (*namco.Video, error) {
	return namco.NewVideo(r, mem, rom, VideoConfig)
}

var VideoConfig = namco.Config{
	TileLayout: namco.SheetLayout{
		CellW:        8,
		CellH:        8,
		W:            8 * 16,
		H:            8 * 8,
		PixelLayout:  tilePixels,
		PixelReader:  pixelReader,
		BytesPerCell: 16,
	},
	SpriteLayout: namco.SheetLayout{
		CellW:        16,
		CellH:        16,
		W:            16 * 16,
		H:            16 * 8,
		PixelLayout:  spritePixels,
		PixelReader:  pixelReader,
		BytesPerCell: 64,
	},
	PaletteEntries: 32,
	PaletteColors:  8,
	VideoAddr:      0x8000,
	Hack:           true,
}

var tilePixels = [][]int{
	[]int{63, 59, 55, 51, 47, 43, 39, 35},
	[]int{62, 58, 54, 50, 46, 42, 38, 34},
	[]int{61, 57, 53, 49, 45, 41, 37, 33},
	[]int{60, 56, 52, 48, 44, 40, 36, 32},

	[]int{31, 27, 23, 19, 15, 11, 7, 3},
	[]int{30, 26, 22, 18, 14, 10, 6, 2},
	[]int{29, 25, 21, 17, 13, 9, 5, 1},
	[]int{28, 24, 20, 16, 12, 8, 4, 0},
}

var spritePixels = [][]int{
	[]int{159, 155, 151, 147, 143, 139, 135, 131, 31, 27, 23, 19, 15, 11, 7, 3},
	[]int{158, 154, 150, 146, 142, 138, 134, 130, 30, 26, 22, 18, 14, 10, 6, 2},
	[]int{157, 153, 149, 145, 141, 137, 133, 129, 29, 25, 21, 17, 13, 9, 5, 1},
	[]int{156, 152, 148, 144, 140, 136, 132, 128, 28, 24, 20, 16, 12, 8, 4, 0},

	[]int{191, 187, 183, 179, 175, 171, 167, 163, 63, 59, 55, 51, 47, 43, 39, 35},
	[]int{190, 186, 182, 178, 174, 170, 166, 162, 62, 58, 54, 50, 46, 42, 38, 34},
	[]int{189, 185, 181, 177, 173, 169, 165, 161, 61, 57, 53, 49, 45, 41, 37, 33},
	[]int{188, 184, 180, 176, 172, 168, 164, 160, 60, 56, 52, 48, 44, 40, 36, 32},

	[]int{223, 219, 215, 211, 207, 203, 199, 195, 95, 91, 87, 83, 79, 75, 71, 67},
	[]int{222, 218, 214, 210, 206, 202, 198, 194, 94, 90, 86, 82, 78, 74, 70, 66},
	[]int{221, 217, 213, 209, 205, 201, 197, 193, 93, 89, 85, 81, 77, 73, 69, 65},
	[]int{220, 216, 212, 208, 204, 200, 196, 192, 92, 88, 84, 80, 76, 72, 68, 64},

	[]int{255, 251, 247, 243, 239, 235, 231, 227, 127, 123, 119, 115, 111, 107, 103, 99},
	[]int{254, 250, 246, 242, 238, 234, 230, 226, 126, 122, 118, 114, 110, 106, 102, 98},
	[]int{253, 249, 245, 241, 237, 233, 229, 225, 125, 121, 117, 113, 109, 105, 101, 97},
	[]int{252, 248, 244, 240, 236, 232, 228, 224, 124, 120, 116, 112, 108, 104, 100, 96},
}

func pixelReader(mem memory.Memory, base uint16, pixel int) uint8 {
	addr := base + uint16(pixel/4)
	offset := pixel % 4
	return bits.Plane(mem.Load(addr), []int{0, 4}, offset)
}
