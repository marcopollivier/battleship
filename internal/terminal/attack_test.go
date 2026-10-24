package terminal

import (
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
)

// attackExampleBoard returns the complete example board from spec 003.
func attackExampleBoard(t *testing.T) *board.Board {
	t.Helper()
	b := boardWith(t,
		placement{board.AircraftCarrier, "B1", board.Down}, // A1 B1 C1 B2 B3
		placement{board.Battleship, "E1", board.Horizontal},
		placement{board.Cruiser, "J1", board.Vertical},
		placement{board.Destroyer, "A5", board.Horizontal},
		placement{board.Destroyer, "E3", board.Horizontal},
		placement{board.Submarine, "H5", board.Horizontal},
		placement{board.Submarine, "A8", board.Horizontal},
		placement{board.Submarine, "E8", board.Horizontal},
	)
	if !b.Complete() {
		t.Fatal("example board is not complete")
	}
	return b
}

// attackFire shoots at each cell, such as "E1", failing the test on error.
func attackFire(t *testing.T, b *board.Board, cells ...string) {
	t.Helper()
	for _, cell := range cells {
		c, err := board.ParseCoord(cell)
		if err != nil {
			t.Fatalf("ParseCoord(%q): %v", cell, err)
		}
		if _, err := b.Fire(c); err != nil {
			t.Fatalf("Fire(%s): %v", cell, err)
		}
	}
}
