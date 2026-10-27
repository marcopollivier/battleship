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
	sunk := sunkCells(b)
	skip := aroundSunk(sunk)
	hits := pendingHits(b, sunk)
	ends, targets := lineEnds(hits), targetCells(hits)
	// From the best choice to the last resort: the ends of the line of hits,
	// the cells next to the hits, a hunt away from sunk ships and, only when
	// every ship has sunk, any cell.
	for _, keep := range []func(Coord) bool{
		func(c Coord) bool { return ends[c] && !skip[c] },
		func(c Coord) bool { return targets[c] && !skip[c] },
		func(c Coord) bool { return !skip[c] },
		nil,
	} {
		if cells := freeCells(b, keep); len(cells) > 0 {
			return draw(cells, rng)
		}
	}
	return Coord{}, ErrNoCellsLeft
}

// lineEnds returns the cell before the first hit and the cell after the
// last one when there are two or more hits, all in one row or all in one
// column, and nil otherwise. The hits must be row by row, as pendingHits
// lists them. The ends may be off the board.
func lineEnds(hits []Coord) map[Coord]bool {
	if len(hits) < 2 {
		return nil
	}
	first, last := hits[0], hits[len(hits)-1]
	switch {
	case allOn(hits, func(c Coord) bool { return c.Row == first.Row }):
		return map[Coord]bool{
			{Col: first.Col - 1, Row: first.Row}: true,
			{Col: last.Col + 1, Row: last.Row}:   true,
		}
	case allOn(hits, func(c Coord) bool { return c.Col == first.Col }):
		return map[Coord]bool{
			{Col: first.Col, Row: first.Row - 1}: true,
			{Col: last.Col, Row: last.Row + 1}:   true,
		}
	}
	return nil
}

// allOn reports whether on is true for every cell.
func allOn(cells []Coord, on func(Coord) bool) bool {
	for _, c := range cells {
		if !on(c) {
			return false
		}
	}
	return true
}

// pendingHits lists the hit cells of b that are not part of a sunk ship,
// row by row: the hits on a ship still afloat.
func pendingHits(b *Board, sunk map[Coord]bool) []Coord {
	var hits []Coord
	for row := range Size {
		for col := range Size {
			c := Coord{Col: col, Row: row}
			if b.Cell(c) == CellHit && !sunk[c] {
				hits = append(hits, c)
			}
		}
	}
	return hits
}

// targetCells returns the cells to shoot at next to the hits: above, below,
// left and right of each one. Cells off the board or already shot are left
// for the caller to filter.
func targetCells(hits []Coord) map[Coord]bool {
	targets := map[Coord]bool{}
	for _, h := range hits {
		for _, n := range neighbors(h) {
			targets[n] = true
		}
	}
	return targets
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

// neighbors returns the cells above, below, left and right of c, which may
// be off the board.
func neighbors(c Coord) []Coord {
	return []Coord{
		{Col: c.Col, Row: c.Row - 1},
		{Col: c.Col, Row: c.Row + 1},
		{Col: c.Col - 1, Row: c.Row},
		{Col: c.Col + 1, Row: c.Row},
	}
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

// aroundSunk returns the cells that touch the sunk cells, diagonals included.
// No ship can be there, because ships do not touch. This rule is a
// provisional decision (issue #25) and may be dropped.
func aroundSunk(sunk map[Coord]bool) map[Coord]bool {
	around := map[Coord]bool{}
	for c := range sunk {
		for dCol := -1; dCol <= 1; dCol++ {
			for dRow := -1; dRow <= 1; dRow++ {
				around[Coord{Col: c.Col + dCol, Row: c.Row + dRow}] = true
			}
		}
	}
	return around
}
