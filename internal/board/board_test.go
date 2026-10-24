package board

import (
	"errors"
	"slices"
	"testing"
)

type placement struct {
	shipType    ShipType
	ref         string
	orientation Orientation
}

// place puts ships on b and fails the test if any placement is rejected.
func place(t *testing.T, b *Board, placements ...placement) {
	t.Helper()
	for _, p := range placements {
		if err := b.Place(p.shipType, mustCoord(t, p.ref), p.orientation); err != nil {
			t.Fatalf("Place(%v, %s, %v): %v", p.shipType, p.ref, p.orientation, err)
		}
	}
}

// assertRejected checks that placing p on b fails with want and leaves b unchanged.
func assertRejected(t *testing.T, b *Board, p placement, want error) {
	t.Helper()
	before := b.Ships()
	err := b.Place(p.shipType, mustCoord(t, p.ref), p.orientation)
	if !errors.Is(err, want) {
		t.Fatalf("Place(%v, %s, %v) error = %v, want %v", p.shipType, p.ref, p.orientation, err, want)
	}
	if after := b.Ships(); !slices.EqualFunc(before, after, func(a, b PlacedShip) bool {
		return a.Type == b.Type && slices.Equal(a.Cells, b.Cells)
	}) {
		t.Errorf("board changed after rejected placement: before %v, after %v", before, after)
	}
}

// CA-4 (accepted case): a battleship can reach the last row.
func TestBoard_Place_AtEdge(t *testing.T) {
	b := NewBoard(DefaultFleet())
	place(t, b, placement{Battleship, "J7", Vertical})

	ships := b.Ships()
	if len(ships) != 1 {
		t.Fatalf("Ships() has %d ships, want 1", len(ships))
	}
	want := []string{"J10", "J7", "J8", "J9"}
	if got := cellNames(ships[0].Cells); ships[0].Type != Battleship || !slices.Equal(got, want) {
		t.Errorf("Ships()[0] = %v %v, want Battleship %v", ships[0].Type, got, want)
	}
}

// CA-4 and CA-6: ships that would leave the board are rejected.
func TestBoard_Place_OutOfBounds(t *testing.T) {
	tests := []struct {
		name string
		p    placement
	}{
		{"CA-4 battleship past last row", placement{Battleship, "J8", Vertical}},
		{"CA-4 battleship past last column", placement{Battleship, "H1", Horizontal}},
		{"CA-6 carrier stem past last row", placement{AircraftCarrier, "E9", Down}},
		{"CA-6 carrier bar before column A", placement{AircraftCarrier, "A5", Up}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertRejected(t, NewBoard(DefaultFleet()), tt.p, ErrOutOfBounds)
		})
	}
}

// CA-11: Place rejects orientations that do not apply to the ship type.
func TestBoard_Place_InvalidOrientation(t *testing.T) {
	tests := []struct {
		name string
		p    placement
	}{
		{"cruiser pointing down", placement{Cruiser, "E5", Down}},
		{"carrier horizontal", placement{AircraftCarrier, "E5", Horizontal}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertRejected(t, NewBoard(DefaultFleet()), tt.p, ErrInvalidOrientation)
		})
	}
}

func TestBoard_Ships_ReturnsCopy(t *testing.T) {
	b := NewBoard(DefaultFleet())
	place(t, b, placement{Submarine, "A1", Horizontal})

	ships := b.Ships()
	ships[0].Cells[0] = Coord{Col: 9, Row: 9}

	if got := b.Ships()[0].Cells[0]; got != (Coord{}) {
		t.Errorf("changing the result of Ships() changed the board: cell is %v, want A1", got)
	}
}
