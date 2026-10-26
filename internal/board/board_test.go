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

// CA-7: a ship cannot share a cell with another ship.
func TestBoard_Place_Overlap(t *testing.T) {
	b := NewBoard(DefaultFleet())
	place(t, b, placement{Cruiser, "A1", Horizontal})

	assertRejected(t, b, placement{Destroyer, "B1", Vertical}, ErrOverlap)
	if n := len(b.Ships()); n != 1 {
		t.Errorf("Ships() has %d ships, want only the cruiser", n)
	}
}

// CA-10: ships cannot touch, not even diagonally.
func TestBoard_Place_Adjacent(t *testing.T) {
	tests := []struct {
		name string
		p    placement
	}{
		{"touching the side", placement{Destroyer, "D1", Horizontal}},
		{"touching below", placement{Destroyer, "A2", Horizontal}},
		{"touching diagonally", placement{Destroyer, "D2", Vertical}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoard(DefaultFleet())
			place(t, b, placement{Cruiser, "A1", Horizontal})
			assertRejected(t, b, tt.p, ErrAdjacent)
		})
	}

	t.Run("one water cell apart", func(t *testing.T) {
		b := NewBoard(DefaultFleet())
		place(t, b, placement{Cruiser, "A1", Horizontal}, placement{Destroyer, "E1", Horizontal})
	})
}

// CA-8: a fleet allows only as many ships of a type as it lists.
func TestBoard_Place_FleetLimit(t *testing.T) {
	t.Run("fourth submarine", func(t *testing.T) {
		b := NewBoard(DefaultFleet())
		place(t, b,
			placement{Submarine, "A1", Horizontal},
			placement{Submarine, "C1", Horizontal},
			placement{Submarine, "E1", Horizontal},
		)
		assertRejected(t, b, placement{Submarine, "G1", Horizontal}, ErrFleetLimit)
	})

	t.Run("type not in fleet", func(t *testing.T) {
		b := NewBoard([]FleetEntry{{Type: Submarine, Count: 1}})
		assertRejected(t, b, placement{Cruiser, "A1", Horizontal}, ErrFleetLimit)
	})
}

// CA-9: the board is complete once all 8 ships of the default fleet are placed.
func TestBoard_Complete(t *testing.T) {
	b := NewBoard(DefaultFleet())
	if b.Complete() {
		t.Fatal("empty board is complete")
	}

	place(t, b,
		placement{AircraftCarrier, "B1", Down}, // A1 B1 C1 B2 B3
		placement{Battleship, "E1", Horizontal},
		placement{Cruiser, "J1", Vertical},
		placement{Destroyer, "A5", Horizontal},
		placement{Destroyer, "E3", Horizontal},
		placement{Submarine, "H5", Horizontal},
		placement{Submarine, "A8", Horizontal},
	)
	if b.Complete() {
		t.Fatal("board with 7 ships is complete")
	}

	place(t, b, placement{Submarine, "E8", Horizontal})
	if !b.Complete() {
		t.Error("board with all 8 ships is not complete")
	}
}

// Remaining lists the ships still to place, largest first.
func TestBoard_Remaining(t *testing.T) {
	b := NewBoard(DefaultFleet())
	want := []ShipType{AircraftCarrier, Battleship, Cruiser, Destroyer, Destroyer, Submarine, Submarine, Submarine}
	if got := b.Remaining(); !slices.Equal(got, want) {
		t.Errorf("Remaining() on an empty board = %v, want %v", got, want)
	}

	place(t, b, placement{Destroyer, "A1", Horizontal}, placement{AircraftCarrier, "E5", Down})
	want = []ShipType{Battleship, Cruiser, Destroyer, Submarine, Submarine, Submarine}
	if got := b.Remaining(); !slices.Equal(got, want) {
		t.Errorf("Remaining() after a destroyer and the carrier = %v, want %v", got, want)
	}
}

// CA-7 (spec 012): Clear removes every ship, so the whole fleet is missing again.
func TestBoard_Clear(t *testing.T) {
	b := NewBoard([]FleetEntry{{Type: Submarine, Count: 1}})
	place(t, b, placement{Submarine, "A1", Horizontal})
	if _, err := b.Fire(Coord{}); err != nil {
		t.Fatalf("Fire(A1): %v", err)
	}

	b.Clear()
	if n := len(b.Ships()); n != 0 {
		t.Errorf("Ships() after Clear has %d ships, want 0", n)
	}
	if got := b.Cell(Coord{}); got != CellUntouched {
		t.Errorf("Cell(A1) after Clear = %v, want CellUntouched", got)
	}
	if got := b.Remaining(); !slices.Equal(got, []ShipType{Submarine}) {
		t.Errorf("Remaining() after Clear = %v, want [Submarine]", got)
	}
	place(t, b, placement{Submarine, "A1", Horizontal})
}
