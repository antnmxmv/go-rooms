package tic_tac_toe

import "errors"

/*
	https://en.wikipedia.org/wiki/Tic-tac-toe

	5x5 field. When player does line of 3 Xs or Os, he wins.
*/


var name = "tic_tac_toe"

type Board struct {
	grid [][]int
}

func (b Board)GetName() string {
	return name
}

func (b *Board) Initialize() {
	b.grid = make([][]int, 5)
	for i := 0; i < len(b.grid); i++ {
		b.grid[i] = make([]int, 5)
	}
}

func (b Board) GetWinner() int {
	for i := 0; i+2 < len(b.grid); i++ {
		for j := 0; j+2 < len(b.grid); j++ {
			for x := 0; x < 3; x++{
				if b.grid[i+x][j+0] == b.grid[i+x][j+1] && b.grid[i+x][j+2] == b.grid[i+x][j+0] && b.grid[i+x][j+1] != 0 {
					return b.grid[i+x][j+0]
				}
				if b.grid[i+0][j+x] == b.grid[i+1][j+x] && b.grid[i+1][j+x] == b.grid[i+2][j+x] && b.grid[i+0][j+x] != 0 {
					return b.grid[i+0][j+x]
				}
			}
			if b.grid[i+0][j+0] == b.grid[i+1][j+1] && b.grid[i+0][j+0] == b.grid[i+2][j+2] && b.grid[i+1][j+1] != 0 {
				return b.grid[i+0][j+0]
			} else if b.grid[i+0][j+2] == b.grid[i+1][j+1] && b.grid[i+0][j+2] == b.grid[i+2][j+0] && b.grid[i+1][j+1] != 0 {
				return b.grid[i+1][j+1]
			}
		}
	}
	return 0
}

func (b Board) checkPoint(x, y int) bool {
	if x >= len(b.grid) || x < 0 || y < 0 || y > len(b.grid) {
		return false
	}
	if b.grid[x][y] != 0 {
		return false
	}
	return true
}

func (b *Board) Turn(params... int) (bool, error) {
	if checker := b.checkPoint(params[1], params[2]); checker {
		b.grid[params[1]][params[2]] = params[0]
		return true, nil
	}
	return true, errors.New("not valid point")
}

func (b Board) GetGrid() [][]int {
	return b.grid
}
