package board

import (
	"errors"
	"math/rand/v2"
)

// ErrNoCellsLeft means every cell of the board has already been shot.
var ErrNoCellsLeft = errors.New("no cells left to shoot")

// ShotStrategy chooses the machine's next shot at b, drawing with rng when
// there is more than one choice. It does not fire.
type ShotStrategy func(b *Board, rng *rand.Rand) (Coord, error)

// RandomShot returns a cell of b that has not been shot yet, chosen with rng.
// It does not fire. The same rng state and board always give the same cell.
func RandomShot(b *Board, rng *rand.Rand) (Coord, error) {
	return draw(freeCells(b, nil), rng)
}

// SmartShot hunts at random until it hits a ship, then shoots next to the
// hits until the ship sinks, as a player on paper would. It uses only what
// the shooter knows: the state of each cell and which ships have sunk.
func SmartShot(b *Board, rng *rand.Rand) (Coord, error) {
	skip := aroundSunk(b)
	candidates := freeCells(b, func(c Coord) bool { return !skip[c] })
	if len(candidates) == 0 {
		// Only when every ship has sunk: any cell is as good as another.
		candidates = freeCells(b, nil)
	}
	return draw(candidates, rng)
}

// draw returns one of cells chosen with rng, or ErrNoCellsLeft if there is none.
func draw(cells []Coord, rng *rand.Rand) (Coord, error) {
	if len(cells) == 0 {
		return Coord{}, ErrNoCellsLeft
	}
	return cells[rng.IntN(len(cells))], nil
}

// freeCells lists the cells of b not shot yet for which keep is true, row by
// row. A nil keep keeps every cell.
func freeCells(b *Board, keep func(Coord) bool) []Coord {
	var free []Coord
	for row := range Size {
		for col := range Size {
			c := Coord{Col: col, Row: row}
			if b.Cell(c) == CellUntouched && (keep == nil || keep(c)) {
				free = append(free, c)
			}
		}
	}
	return free
}

// sunkCells returns the cells of the ships that have sunk. The shooter knows
// them: a ship that sinks is announced, and since ships do not touch, it
// takes exactly the hit cells connected to the last shot.
func sunkCells(b *Board) map[Coord]bool {
	cells := map[Coord]bool{}
	for i := range b.ships {
		if b.sunk(&b.ships[i]) {
			for _, c := range b.ships[i].Cells {
				cells[c] = true
			}
		}
	}
	return cells
}

// aroundSunk returns the cells that touch a sunk ship, diagonals included.
// No ship can be there, because ships do not touch. This rule is a
// provisional decision (issue #25) and may be dropped.
func aroundSunk(b *Board) map[Coord]bool {
	around := map[Coord]bool{}
	for c := range sunkCells(b) {
		for dCol := -1; dCol <= 1; dCol++ {
			for dRow := -1; dRow <= 1; dRow++ {
				around[Coord{Col: c.Col + dCol, Row: c.Row + dRow}] = true
			}
		}
	}
	return around
}
