package games

import (
	"go-rooms/games/hexapawn"
	"go-rooms/games/tic_tac_toe"
	"go-rooms/lib/interfaces"
)

func GetInstance(name string) interfaces.Game {
	switch name {
	case "tic_tac_toe":
		return &tic_tac_toe.Board{}
	case "hexapawn":
		return &hexapawn.Board{}
	default:
		return nil
	}
}
