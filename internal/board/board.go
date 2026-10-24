package board

import (
	"errors"
	"fmt"
	"slices"
)

var (
	// ErrOverlap means the ship would share a cell with another ship.
	ErrOverlap = errors.New("ship overlaps another ship")
	// ErrAdjacent means the ship would touch another ship, including diagonally.
	ErrAdjacent = errors.New("ship touches another ship")
	// ErrFleetLimit means every ship of that type in the fleet is already placed.
	ErrFleetLimit = errors.New("no more ships of this type in the fleet")
)

// PlacedShip is a ship on the board and the cells it occupies.
type PlacedShip struct {
	Type  ShipType
	Cells []Coord
}

// Board is one player's board.
type Board struct {
	fleet []FleetEntry
	ships []PlacedShip
}

// NewBoard returns an empty board for the given fleet.
func NewBoard(fleet []FleetEntry) *Board {
	return &Board{fleet: slices.Clone(fleet)}
}

// Place puts a ship of type t on the board at ref with orientation o.
// On error the board is left unchanged.
func (b *Board) Place(t ShipType, ref Coord, o Orientation) error {
	if b.placedCount(t) >= b.fleetCount(t) {
		return fmt.Errorf("%w: %v", ErrFleetLimit, t)
	}
	cells, err := ShipCells(t, ref, o)
	if err != nil {
		return err
	}
	for _, c := range cells {
		if !c.inBounds() {
			return fmt.Errorf("%w: %v at %v %v", ErrOutOfBounds, t, ref, o)
		}
	}
	occupied := b.occupied()
	for _, c := range cells {
		if occupied[c] {
			return fmt.Errorf("%w: %v at %v %v", ErrOverlap, t, ref, o)
		}
	}
	// Overlap is checked first because every overlap also touches.
	for _, c := range cells {
		for dCol := -1; dCol <= 1; dCol++ {
			for dRow := -1; dRow <= 1; dRow++ {
				if occupied[Coord{Col: c.Col + dCol, Row: c.Row + dRow}] {
					return fmt.Errorf("%w: %v at %v %v", ErrAdjacent, t, ref, o)
				}
			}
		}
	}

	b.ships = append(b.ships, PlacedShip{Type: t, Cells: cells})
	return nil
}

// Ships returns a copy of the placed ships, in placement order.
func (b *Board) Ships() []PlacedShip {
	ships := make([]PlacedShip, len(b.ships))
	for i, s := range b.ships {
		ships[i] = PlacedShip{Type: s.Type, Cells: slices.Clone(s.Cells)}
	}
	return ships
}

// Complete reports whether every ship in the fleet has been placed.
func (b *Board) Complete() bool {
	for _, entry := range b.fleet {
		if b.placedCount(entry.Type) < entry.Count {
			return false
		}
	}
	return true
}

func (b *Board) fleetCount(t ShipType) int {
	count := 0
	for _, entry := range b.fleet {
		if entry.Type == t {
			count += entry.Count
		}
	}
	return count
}

func (b *Board) placedCount(t ShipType) int {
	count := 0
	for _, s := range b.ships {
		if s.Type == t {
			count++
		}
	}
	return count
}

// occupied returns the set of cells taken by placed ships.
func (b *Board) occupied() map[Coord]bool {
	cells := map[Coord]bool{}
	for _, s := range b.ships {
		for _, c := range s.Cells {
			cells[c] = true
		}
	}
	return cells
}
