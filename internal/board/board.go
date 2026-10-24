package board

import (
	"errors"
	"fmt"
	"slices"
)

// ErrOverlap means the ship would share a cell with another ship.
var ErrOverlap = errors.New("ship overlaps another ship")

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
