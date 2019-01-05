package game

import (
	"github.com/blackchip-org/pac8/pkg/machine"
	"github.com/blackchip-org/pac8/pkg/memory"
	"github.com/blackchip-org/pac8/pkg/pac8"
	"github.com/blackchip-org/pac8/system/pacman"
)

var msPacManROM = memory.NewPack().
	Add("code    ", "boot1    ", "bc2247ec946b639dd1f00bfc603fa157d0baaa97").
	Add("code    ", "boot2    ", "13ea0c343de072508908be885e6a2a217bbb3047").
	Add("code    ", "boot3    ", "5ea4d907dbb2690698db72c4e0b5be4d3e9a7786").
	Add("code    ", "boot4    ", "3022a408118fa7420060e32a760aeef15b8a96cf").
	Add("code2   ", "boot5    ", "fed6e9a2b210b07e7189a18574f6b8c4ec5bb49b").
	Add("code2   ", "boot6    ", "387010a0c76319a1eab61b54c9bcb5c66c4b67a1").
	Add("tile    ", "5e       ", "5e8b472b615f12efca3fe792410c23619f067845").
	Add("sprite  ", "5f       ", "fd6a1dde780b39aea76bf1c4befa5882573c2ef4").
	Add("color   ", "82s123.7f", "8d0268dee78e47c712202b0ec4f1f51109b1f2a5").
	Add("palette ", "82s126.4a", "19097b5f60d1030f8b82d9f1d3a241f93e5c75d6").
	Add("waveform", "82s126.1m", "bbcec0570aeceb582ff8238a4bc8546a23430081").
	Add("waveform", "82s126.3m", "0c4d0bee858b97632411c440bea6948a74759746")

var MsPacMan = pac8.Game{
	ROM: msPacManROM,
	Init: func(env pac8.Env, roms memory.Set) (machine.System, error) {
		config := pacman.Config{
			Name: "mspacman",
		}
		return pacman.New(env, config, roms)
	},
}
