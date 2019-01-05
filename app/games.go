package app

import "github.com/blackchip-org/pac8/pkg/pac8"
import "github.com/blackchip-org/pac8/game"

var Games = map[string]pac8.Game{
	"galaga":   game.Galaga,
	"mspacman": game.MsPacMan,
	"pacman":   game.PacMan,
}
