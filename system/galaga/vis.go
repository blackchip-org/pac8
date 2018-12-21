package galaga

import (
	"github.com/blackchip-org/pac8/app"
	"github.com/blackchip-org/pac8/component/namco"
	"github.com/veandco/go-sdl2/sdl"
)

// http://tech.quarterarcade.com/tech/MAME/src/galaga.c.html.aspx?g=1042

func GalagaTiles(r *sdl.Renderer) (*sdl.Texture, error) {
	mem, err := app.LoadROM("galaga/07m_g08.bin", "62f1279a784ab2f8218c4137c7accda00e6a3490")
	if err != nil {
		return nil, err
	}
	info := namco.SheetInfo{
		W:      16,
		H:      8,
		Colors: namco.VisPalette,
	}
	tex, err := namco.TileSheet(r, mem, info)
	if err != nil {
		return nil, err
	}
	return tex, nil
}
