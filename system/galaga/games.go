package galaga

import (
	"github.com/blackchip-org/pac8/app"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/namco"
	"github.com/blackchip-org/pac8/machine"
)

func NewGalaga(ctx app.SDLContext) (machine.System, error) {
	l := app.NewLoader("galaga")
	rom1a := l.Load("04m_g01.bin", "6907773db7c002ecde5e41853603d53387c5c7cd")
	rom1b := l.Load("04k_g02.bin", "666975aed5ce84f09794c54b550d64d95ab311f0")
	rom1c := l.Load("04j_g03.bin", "481f443aea3ed3504ec2f3a6bfcf3cd47e2f8f81")
	rom1d := l.Load("04h_g04.bin", "366cb0dbd31b787e64f88d182108b670d03b393e")
	rom2 := l.Load("04e_g05.bin", "d29b68d6aab3217fa2106b3507b9273ff3f927bf")
	rom3 := l.Load("04d_g06.bin", "d6cb439de0718826d1a0363c9d77de8740b18ecf")
	tiles := l.Load("07m_g08.bin", "62f1279a784ab2f8218c4137c7accda00e6a3490")
	sprites0 := l.Load("07e_g10.bin", "e697c180178cabd1d32483c5d8889a40633f7857")
	sprites1 := l.Load("07h_g09.bin", "c340ed8c25e0979629a9a1730edc762bd72d0cff")
	empty := memory.NewNull(0x1000)

	sm := memory.NewBlockMapper()
	sm.Map(0x0000, sprites0)
	sm.Map(0x1000, sprites1)
	sprites := memory.NewPageMapped(sm.Blocks)

	vrom := namco.VideoROM{
		Tiles:   tiles,
		Sprites: sprites,
	}

	if err := l.Error(); err != nil {
		return nil, err
	}

	config := Config{
		Name: "galaga",
		ProcROM: [3][]memory.Memory{
			[]memory.Memory{rom1a, rom1b, rom1c, rom1d},
			[]memory.Memory{rom2, empty, empty, empty},
			[]memory.Memory{rom3, empty, empty, empty},
		},
		VideoROM: vrom,
	}
	return New(ctx, config)
}
