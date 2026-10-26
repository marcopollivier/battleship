package board

import (
	"errors"
	"math/rand/v2"
)

// ErrNoCellsLeft means every cell of the board has already been shot.
var ErrNoCellsLeft = errors.New("no cells left to shoot")

// RandomShot returns a cell of b that has not been shot yet, chosen with rng.
// It does not fire. The same rng state and board always give the same cell.
func RandomShot(b *Board, rng *rand.Rand) (Coord, error) {
	var free []Coord
	for row := range Size {
		for col := range Size {
			c := Coord{Col: col, Row: row}
			if b.Cell(c) == CellUntouched {
				free = append(free, c)
			}
		}
	}
	if len(free) == 0 {
		return Coord{}, ErrNoCellsLeft
	}
	return free[rng.IntN(len(free))], nil
}
