package game

import (
	"github.com/blackchip-org/pac8/pkg/machine"
	"github.com/blackchip-org/pac8/pkg/memory"
	"github.com/blackchip-org/pac8/pkg/pac8"
	"github.com/blackchip-org/pac8/system/pacman"
)

var pacManROM = memory.NewPack().
	Add("code    ", "pacman.6e", "e87e059c5be45753f7e9f33dff851f16d6751181").
	Add("code    ", "pacman.6f", "674d3a7f00d8be5e38b1fdc208ebef5a92d38329").
	Add("code    ", "pacman.6h", "8e47e8c2c4d6117d174cdac150392042d3e0a881").
	Add("code    ", "pacman.6j", "d4a70d56bb01d27d094d73db8667ffb00ca69cb9").
	Add("tile    ", "pacman.5e", "06ef227747a440831c9a3a613b76693d52a2f0a9").
	Add("sprite  ", "pacman.5f", "4a937ac02216ea8c96477d4a15522070507fb599").
	Add("color   ", "82s123.7f", "8d0268dee78e47c712202b0ec4f1f51109b1f2a5").
	Add("palette ", "82s126.4a", "19097b5f60d1030f8b82d9f1d3a241f93e5c75d6").
	Add("waveform", "82s126.1m", "bbcec0570aeceb582ff8238a4bc8546a23430081").
	Add("waveform", "82s126.3m", "0c4d0bee858b97632411c440bea6948a74759746")

var PacManConfig = pacman.Config{
	Name: "pacman",
}

var PacMan = pac8.Game{
	ROM: pacManROM,
	Init: func(env pac8.Env, roms memory.Set) (machine.System, error) {
		config := pacman.Config{
			Name: "pacman",
		}
		return pacman.New(env, config, roms)
	},
}
