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

	if b.shipAt(c) == nil {
		return Shot{Result: Miss}, nil
	}
	return Shot{Result: Hit}, nil
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
