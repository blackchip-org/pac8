package galaga

import (
	"github.com/blackchip-org/pac8/app"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/machine"
)

func NewGalaga(ctx app.SDLContext) (machine.System, error) {
	l := app.NewLoader("galaga")
	rom0 := l.Load("04m_g01.bin", "6907773db7c002ecde5e41853603d53387c5c7cd")
	rom1 := l.Load("04k_g02.bin", "666975aed5ce84f09794c54b550d64d95ab311f0")
	rom2 := l.Load("04j_g03.bin", "481f443aea3ed3504ec2f3a6bfcf3cd47e2f8f81")
	rom3 := l.Load("04h_g04.bin", "366cb0dbd31b787e64f88d182108b670d03b393e")

	if err := l.Error(); err != nil {
		return nil, err
	}

	m := memory.NewBlockMapper()
	m.Map(0x0000, rom0)
	m.Map(0x1000, rom1)
	m.Map(0x2000, rom2)
	m.Map(0x3000, rom3)

	config := Config{
		Name: "galaga",
		M:    m,
	}
	return New(ctx, config)
}
