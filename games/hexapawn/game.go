package hexapawn

import "errors"

/*
	https://en.wikipedia.org/wiki/Hexapawn
 */

var name = "hexapawn"

type Board struct {
	grid [][]int
}

func (b Board) GetName() string {
	return name
}

func (b *Board) Initialize() {
	b.grid = make([][]int, 3)
	for i := 0; i < len(b.grid); i++ {
		b.grid[i] = make([]int, 3)
	}
	for i := 0; i < 3; i++{
		b.grid[0][i] = 2
		b.grid[2][i] = 1
	}
}

func (b Board) GetWinner() int {
	ones, twos := false, false
	for i := 0; i < 3; i++ {
		if b.grid[0][i] == 1 {
			return 1;
		}
		if b.grid[2][i] == 2 {
			return 2;
		}
		for j := 0; j < 3; j++{
			if b.grid[i][j] == 1{
				ones = true
			}
			if b.grid[i][j] == 2{
				ones = true
			}
		}
	}
	if !ones{
		return 2
	}
	if !twos {
		return 1
	}
	return 0
}

func (b Board) checkPoint(role, x1, y1, x2, y2 int) bool {
	if x1 >= len(b.grid) || x2 >= len(b.grid) || y1 >= len(b.grid) || y2 >= len(b.grid) || x1 < 0 || y1 < 0 || x2 < 0 || y2 < 0 {
		return false
	}
	if x2 == x1+1 && y1 == y2 && b.grid[x2][y2] != role && role == 2 {
		return true
	}
	if x2 == x1-1 && y1 == y2 && b.grid[x2][y2] != role && role == 1 {
		return true
	}
	if (x2 == x1+1 && (y2 == y1+1 || y2 == y1 - 1)) && b.grid[x2][y2] != role && b.grid[x2][y2] != 0 && role == 2{
		return true
	}
	if (x2 == x1-1 && (y2 == y1+1 || y2 == y1 - 1)) && b.grid[x2][y2] != role && b.grid[x2][y2] != 0  && role == 1{
		return true
	}
	return false
}

func (b *Board) Turn(params ... int) (bool, error) {
	if len(params) != 5{
		return true, errors.New("not valid point")
	}
	if checker := b.checkPoint(params[0], params[1], params[2], params[3], params[4]); checker {
		b.grid[params[1]][params[2]] = 0
		b.grid[params[3]][params[4]] = params[0]
		return true, nil
	}
	return true, errors.New("not valid point")
}

func (b Board) GetGrid() [][]int {
	return b.grid
}
