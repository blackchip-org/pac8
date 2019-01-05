package galaga

import (
	"github.com/veandco/go-sdl2/sdl"
)

// http://tech.quarterarcade.com/tech/MAME/src/galaga.c.html.aspx?g=1042

func GalagaTiles(r *sdl.Renderer) (*sdl.Texture, error) {
	/*
		mem, err := app.LoadROM("galaga/07m_g08.bin", "62f1279a784ab2f8218c4137c7accda00e6a3490")
		if err != nil {
			return nil, err
		}
		layout := namco.SheetLayout{
			CellW:        8,
			CellH:        8,
			W:            8 * 16,
			H:            8 * 8,
			PixelLayout:  tilePixels,
			PixelReader:  pixelReader,
			BytesPerCell: 16,
		}
		sheet, err := namco.NewSheet(r, mem, layout, namco.VisPalette)
		if err != nil {
			return nil, err
		}
		return sheet.Texture, nil
	*/
	return nil, nil
}

func GalagaSprites(r *sdl.Renderer) (*sdl.Texture, error) {
	/*
		ram0, err := app.LoadROM("galaga/07e_g10.bin", "e697c180178cabd1d32483c5d8889a40633f7857")
		if err != nil {
			return nil, err
		}
		ram1, err := app.LoadROM("galaga/07h_g09.bin", "c340ed8c25e0979629a9a1730edc762bd72d0cff")
		if err != nil {
			return nil, err
		}

		m := memory.NewBlockMapper()
		m.Map(0x0000, ram0)
		m.Map(0x1000, ram1)
		mem := memory.NewPageMapped(m.Blocks)

		info := namco.SheetLayout{
			CellW:        16,
			CellH:        16,
			W:            16 * 16,
			H:            16 * 8,
			PixelLayout:  spritePixels,
			PixelReader:  pixelReader,
			BytesPerCell: 64,
		}
		sheet, err := namco.NewSheet(r, mem, info, namco.VisPalette)
		if err != nil {
			return nil, err
		}
		return sheet.Texture, nil
	*/
	return nil, nil
}
