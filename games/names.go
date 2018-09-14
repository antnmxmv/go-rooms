package games

import (
	"go-rooms/games/chatroulette"
	"go-rooms/games/tic_tac_toe"
	"go-rooms/lib/interfaces"
)

func GetInstance(name string) interfaces.Game {
	switch name {
	case "tic_tac_toe":
		return &tic_tac_toe.Board{}
	case "chatroulette":
		return &chatroulette.Chat{}
	default:
		return nil
	}
}
