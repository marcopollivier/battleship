package board

import (
	"testing"
	"time"
)

// battleshipRows returns a board for fleet with six battleships in rows 2, 5
// and 8, two per row, and a seventh at last in row 10. They leave no room
// outside row 10: with last at C10, only A10 and H10 to J10 are free.
func battleshipRows(t *testing.T, fleet []FleetEntry, last string) *Board {
	t.Helper()
	b := NewBoard(fleet)
	for _, row := range []string{"2", "5", "8"} {
		place(t, b, placement{Battleship, "A" + row, Horizontal}, placement{Battleship, "F" + row, Horizontal})
	}
	place(t, b, placement{Battleship, last, Horizontal})
	return b
}

// fitBoard builds a board for one Fits case.
type fitBoard func(t *testing.T) *Board

// defaultFleetWith returns a default-fleet board with placements.
func defaultFleetWith(placements ...placement) fitBoard {
	return func(t *testing.T) *Board {
		b := NewBoard(DefaultFleet())
		place(t, b, placements...)
		return b
	}
}

// fleetWith returns a board for fleet with the battleship rows of
// battleshipRows, the last one at last.
func fleetWith(last string, fleet ...FleetEntry) fitBoard {
	return func(t *testing.T) *Board {
		return battleshipRows(t, append([]FleetEntry{{Type: Battleship, Count: 7}}, fleet...), last)
	}
}

// gridSubmarines returns a board for 26 submarines with count of them on a
// grid two cells apart, from A1. At most 25 submarines fit on the board.
func gridSubmarines(count int) fitBoard {
	return func(t *testing.T) *Board {
		b := NewBoard([]FleetEntry{{Type: Submarine, Count: 26}})
		for i := range count {
			if err := b.Place(Submarine, Coord{Col: i % 5 * 2, Row: i / 5 * 2}, Horizontal); err != nil {
				t.Fatalf("Place submarine %d: %v", i, err)
			}
		}
		return b
	}
}

// CA-4 and CA-5 (spec 014): Fits is true only when the missing ships can
// still all be placed, and leaves the board unchanged.
func TestFits(t *testing.T) {
	tests := []struct {
		name  string
		board fitBoard
		want  bool
	}{
		{"CA-5 empty board", defaultFleetWith(), true},
		{"CA-5 carrier in the middle, battleship and cruiser spread", defaultFleetWith(
			placement{AircraftCarrier, "E4", Down},
			placement{Battleship, "A10", Horizontal},
			placement{Cruiser, "J2", Vertical},
		), true},
		{"CA-5 every ship but one submarine", defaultFleetWith(
			placement{AircraftCarrier, "B1", Down},
			placement{Battleship, "E1", Horizontal},
			placement{Cruiser, "J1", Vertical},
			placement{Destroyer, "A5", Horizontal},
			placement{Destroyer, "E3", Horizontal},
			placement{Submarine, "H5", Horizontal},
			placement{Submarine, "A8", Horizontal},
		), true},
		{"complete board", defaultFleetWith(
			placement{AircraftCarrier, "B1", Down},
			placement{Battleship, "E1", Horizontal},
			placement{Cruiser, "J1", Vertical},
			placement{Destroyer, "A5", Horizontal},
			placement{Destroyer, "E3", Horizontal},
			placement{Submarine, "H5", Horizontal},
			placement{Submarine, "A8", Horizontal},
			placement{Submarine, "E8", Horizontal},
		), true},
		{"CA-5 cruiser fits only at H10", fleetWith("C10", FleetEntry{Type: Cruiser, Count: 1}), true},
		{"CA-5 cruiser at H10 and submarine at A10", fleetWith("C10",
			FleetEntry{Type: Cruiser, Count: 1}, FleetEntry{Type: Submarine, Count: 1}), true},
		{"CA-4 no room for the cruiser", fleetWith("D10", FleetEntry{Type: Cruiser, Count: 1}), false},
		{"CA-4 no room for cruiser and destroyer", fleetWith("C10",
			FleetEntry{Type: Cruiser, Count: 1}, FleetEntry{Type: Destroyer, Count: 1}), false},
		{"CA-4 25 submarines on a grid, one missing", gridSubmarines(25), false},
		{"CA-4 20 submarines on a grid, six missing", gridSubmarines(20), false},
		{"unknown ship type", func(*testing.T) *Board {
			return NewBoard([]FleetEntry{{Type: ShipType(99), Count: 1}})
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.board(t)
			before := b.Ships()
			if got := Fits(b); got != tt.want {
				t.Errorf("Fits() = %v, want %v", got, tt.want)
			}
			if after := b.Ships(); !sameShips(before, after) {
				t.Errorf("Fits changed the board: before %v, after %v", before, after)
			}
		})
	}
}

// Fits stops after maxFitAttempts and answers true, so a large fleet does
// not hang the game and the warning is never shown by mistake.
func TestFits_AttemptLimit(t *testing.T) {
	b := gridSubmarines(1)(t)
	start := time.Now()
	fits, attempts := fitSearch(b)
	if attempts <= maxFitAttempts {
		t.Fatalf("fitSearch() used %d attempts, want the limit %d to be passed", attempts, maxFitAttempts)
	}
	if fits {
		t.Error("fitSearch() = true for 26 submarines, want false")
	}
	if !Fits(b) {
		t.Error("Fits() = false after the attempt limit, want true")
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("Fits took %v, want under 5s", elapsed)
	}
}

// With the default fleet, the search stays well below the attempt limit in
// the worst boards found by searching for them (see the spec's Encerramento):
// 283 attempts for a board placed largest first, as the player places it, and
// 2,699 for a board with only some of the smaller ships.
func TestFits_Bounded(t *testing.T) {
	const limit = maxFitAttempts / 100
	tests := []struct {
		name  string
		board fitBoard
	}{
		{"empty board", defaultFleetWith()},
		{"worst board placed largest first", defaultFleetWith(
			placement{AircraftCarrier, "B2", Down},
			placement{Battleship, "F2", Horizontal},
		)},
		{"worst board with smaller ships", defaultFleetWith(
			placement{Cruiser, "E7", Horizontal},
			placement{Destroyer, "I8", Vertical},
			placement{Destroyer, "B10", Horizontal},
			placement{Submarine, "G5", Horizontal},
			placement{Submarine, "B4", Horizontal},
			placement{Submarine, "I2", Horizontal},
		)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, attempts := fitSearch(tt.board(t)); attempts > limit {
				t.Errorf("fitSearch() used %d attempts, want at most %d", attempts, limit)
			}
		})
	}
}
