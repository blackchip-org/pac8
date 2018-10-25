package pacman

import (
	"github.com/blackchip-org/pac8/app"
	"github.com/blackchip-org/pac8/mach"
	"github.com/blackchip-org/pac8/memory"
	"github.com/veandco/go-sdl2/sdl"
)

func NewPacman(renderer *sdl.Renderer) (*mach.Mach, error) {
	// Load ROMs
	l := app.NewLoader("pacman")
	rom0 := l.Load("pacman.6e", "e87e059c5be45753f7e9f33dff851f16d6751181")
	rom1 := l.Load("pacman.6f", "674d3a7f00d8be5e38b1fdc208ebef5a92d38329")
	rom2 := l.Load("pacman.6h", "8e47e8c2c4d6117d174cdac150392042d3e0a881")
	rom3 := l.Load("pacman.6j", "d4a70d56bb01d27d094d73db8667ffb00ca69cb9")
	vroms := VideoROM{
		Tiles:   l.Load("pacman.5e", "06ef227747a440831c9a3a613b76693d52a2f0a9"),
		Sprites: l.Load("pacman.5f", "4a937ac02216ea8c96477d4a15522070507fb599"),
		Color:   l.Load("82s123.7f", "8d0268dee78e47c712202b0ec4f1f51109b1f2a5"),
		Palette: l.Load("82s126.4a", "19097b5f60d1030f8b82d9f1d3a241f93e5c75d6"),
	}
	if err := l.Error(); err != nil {
		return nil, err
	}

	ram := memory.NewRAM(0x1000)
	io := memory.NewIO(0x100)
	mem := memory.NewPageMapped([]memory.Block{
		memory.NewBlock(0x0000, rom0),
		memory.NewBlock(0x1000, rom1),
		memory.NewBlock(0x2000, rom2),
		memory.NewBlock(0x3000, rom3),
		memory.NewBlock(0x4000, ram),
		memory.NewBlock(0x5000, io),
	})
	// Mask out the bit 15 address line that is missing in Pacman
	mem = memory.NewMasked(mem, 0x7fff)
	m, err := New(renderer, mem, vroms, io)
	return m, err
}

func NewMsPacman(renderer *sdl.Renderer) (*mach.Mach, error) {
	l := app.NewLoader("mspacman")
	rom0 := l.Load("boot1", "bc2247ec946b639dd1f00bfc603fa157d0baaa97")
	rom1 := l.Load("boot2", "13ea0c343de072508908be885e6a2a217bbb3047")
	rom2 := l.Load("boot3", "5ea4d907dbb2690698db72c4e0b5be4d3e9a7786")
	rom3 := l.Load("boot4", "3022a408118fa7420060e32a760aeef15b8a96cf")
	rom4 := l.Load("boot5", "fed6e9a2b210b07e7189a18574f6b8c4ec5bb49b")
	rom5 := l.Load("boot6", "387010a0c76319a1eab61b54c9bcb5c66c4b67a1")
	vroms := VideoROM{
		Tiles:   l.Load("5e", "5e8b472b615f12efca3fe792410c23619f067845"),
		Sprites: l.Load("5f", "fd6a1dde780b39aea76bf1c4befa5882573c2ef4"),
		Color:   l.Load("82s123.7f", "8d0268dee78e47c712202b0ec4f1f51109b1f2a5"),
		Palette: l.Load("82s126.4a", "19097b5f60d1030f8b82d9f1d3a241f93e5c75d6"),
	}
	if err := l.Error(); err != nil {
		return nil, err
	}

	ram := memory.NewRAM(0x1000)
	io := memory.NewIO(0x100)
	mem := memory.NewPageMapped([]memory.Block{
		memory.NewBlock(0x0000, rom0),
		memory.NewBlock(0x1000, rom1),
		memory.NewBlock(0x2000, rom2),
		memory.NewBlock(0x3000, rom3),
		memory.NewBlock(0x4000, ram),
		memory.NewBlock(0x5000, io),
		memory.NewBlock(0x8000, rom4),
		memory.NewBlock(0x9000, rom5),
	})
	m, err := New(renderer, mem, vroms, io)
	return m, err
}
