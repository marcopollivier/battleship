package board

import "fmt"

// ShotResult is the outcome of a shot: Miss, Hit or Sunk.
type ShotResult int

const (
	Miss ShotResult = iota // no ship on the cell
	Hit                    // a ship was hit but still has untouched cells
	Sunk                   // the shot hit the last untouched cell of a ship
)

var shotResultNames = map[ShotResult]string{
	Miss: "Miss",
	Hit:  "Hit",
	Sunk: "Sunk",
}

func (r ShotResult) String() string {
	if name, ok := shotResultNames[r]; ok {
		return name
	}
	return fmt.Sprintf("ShotResult(%d)", int(r))
}

// Shot is the outcome of firing at a cell. Ship is set only when Result is Sunk,
// so a Hit does not reveal which ship was hit. For any other result Ship holds
// its zero value, so check Result first.
type Shot struct {
	Result ShotResult
	Ship   ShipType
}

// Fire shoots at c. On error the board is left unchanged.
func (b *Board) Fire(c Coord) (Shot, error) {
	b.shots[c] = true

	ship := b.shipAt(c)
	if ship == nil {
		return Shot{Result: Miss}, nil
	}
	if b.sunk(ship) {
		return Shot{Result: Sunk, Ship: ship.Type}, nil
	}
	return Shot{Result: Hit}, nil
}

// sunk reports whether every cell of s has been shot.
func (b *Board) sunk(s *PlacedShip) bool {
	for _, c := range s.Cells {
		if !b.shots[c] {
			return false
		}
	}
	return true
}

// shipAt returns the placed ship that occupies c, or nil if there is none.
func (b *Board) shipAt(c Coord) *PlacedShip {
	for i := range b.ships {
		for _, cell := range b.ships[i].Cells {
			if cell == c {
				return &b.ships[i]
			}
		}
	}
	return nil
}

// CellState is what the shooter knows about a cell.
type CellState int

const (
	CellUntouched CellState = iota // not shot yet
	CellMiss                       // shot, no ship
	CellHit                        // shot, part of a ship
)

var cellStateNames = map[CellState]string{
	CellUntouched: "CellUntouched",
	CellMiss:      "CellMiss",
	CellHit:       "CellHit",
}

func (s CellState) String() string {
	if name, ok := cellStateNames[s]; ok {
		return name
	}
	return fmt.Sprintf("CellState(%d)", int(s))
}

// Cell reports the state of c. Cells outside the board are CellUntouched.
func (b *Board) Cell(c Coord) CellState {
	switch {
	case !b.shots[c]:
		return CellUntouched
	case b.shipAt(c) == nil:
		return CellMiss
	default:
		return CellHit
	}
}
