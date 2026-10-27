package board

import "math/rand/v2"

// MasterShot shoots where the ships still afloat most likely are: it counts,
// for each cell not shot yet, the ways the remaining ships could cover it
// given every shot so far, and picks a cell with the highest count, drawing
// among ties with rng. Like SmartShot, it uses only what the shooter knows.
func MasterShot(b *Board, rng *rand.Rand) (Coord, error) {
	counts := shipCounts(newChart(b), afloatShips(b))
	var best []Coord
	top := 0
	for _, c := range freeCells(b, nil) {
		switch n := counts[c]; {
		case n > top:
			best, top = []Coord{c}, n
		case n == top && n > 0:
			best = append(best, c)
		}
	}
	// No count at all only happens on a board that does not match its
	// fleet, or with no cells left.
	if top == 0 {
		return SmartShot(b, rng)
	}
	return draw(best, rng)
}

// sighting is what the shooter knows about a cell.
type sighting int

const (
	unknown    sighting = iota // not shot yet
	water                      // shot, no ship
	pendingHit                 // shot, part of a ship still afloat
	sunkHit                    // shot, part of a sunk ship
)

// chart is what the shooter knows about every cell of a board.
type chart [Size][Size]sighting

// newChart returns what the shooter of b knows: the state of each cell and
// which hits belong to sunk ships.
func newChart(b *Board) *chart {
	sunk := sunkCells(b)
	var ch chart
	for row := range Size {
		for col := range Size {
			c := Coord{Col: col, Row: row}
			switch {
			case sunk[c]:
				ch[row][col] = sunkHit
			case b.Cell(c) == CellHit:
				ch[row][col] = pendingHit
			case b.Cell(c) == CellMiss:
				ch[row][col] = water
			}
		}
	}
	return &ch
}

// at returns what the shooter knows about c, or unknown when c is off the board.
func (ch *chart) at(c Coord) sighting {
	if !c.inBounds() {
		return unknown
	}
	return ch[c.Row][c.Col]
}

// afloatShips lists the ships of b's fleet that have not sunk, one entry per
// ship. The shooter knows them: a ship that sinks is announced.
func afloatShips(b *Board) []ShipType {
	sunk := map[ShipType]int{}
	for i := range b.ships {
		if b.sunk(&b.ships[i]) {
			sunk[b.ships[i].Type]++
		}
	}
	var afloat []ShipType
	for _, t := range shipsBySize(b.fleet) {
		if sunk[t] > 0 {
			sunk[t]--
			continue
		}
		afloat = append(afloat, t)
	}
	return afloat
}

// shipCounts returns, for each cell, how many positions of the ships could
// cover it, given what ch shows.
func shipCounts(ch *chart, ships []ShipType) map[Coord]int {
	counts := map[Coord]int{}
	for _, t := range ships {
		for _, cells := range shipPositions(ch, t) {
			for _, c := range cells {
				counts[c]++
			}
		}
	}
	return counts
}

// shipPositions lists the cells of every position where a ship of type t
// could be, given what ch shows. A Submarine has a single cell, so it is
// counted in one orientation only.
func shipPositions(ch *chart, t ShipType) [][]Coord {
	orientations := ValidOrientations(t)
	if t.Size() == 1 {
		orientations = orientations[:1]
	}
	var positions [][]Coord
	for row := range Size {
		for col := range Size {
			for _, o := range orientations {
				cells, err := ShipCells(t, Coord{Col: col, Row: row}, o)
				if err == nil && ch.fits(cells) {
					positions = append(positions, cells)
				}
			}
		}
	}
	return positions
}

// fits reports whether a ship could be on cells: every cell on the board,
// none of them water or part of a sunk ship, and none touching, diagonals
// included, a hit that is not one of cells, since ships do not touch.
func (ch *chart) fits(cells []Coord) bool {
	own := map[Coord]bool{}
	for _, c := range cells {
		if !c.inBounds() || ch.at(c) == water || ch.at(c) == sunkHit {
			return false
		}
		own[c] = true
	}
	for _, c := range cells {
		for dCol := -1; dCol <= 1; dCol++ {
			for dRow := -1; dRow <= 1; dRow++ {
				n := Coord{Col: c.Col + dCol, Row: c.Row + dRow}
				if s := ch.at(n); !own[n] && (s == pendingHit || s == sunkHit) {
					return false
				}
			}
		}
	}
	return true
}
