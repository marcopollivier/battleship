package board

import (
	"errors"
	"maps"
	"testing"
)

// exampleBoard returns the complete example board from spec 003.
func exampleBoard(t *testing.T) *Board {
	t.Helper()
	b := NewBoard(DefaultFleet())
	place(t, b,
		placement{AircraftCarrier, "B1", Down}, // A1 B1 C1 B2 B3
		placement{Battleship, "E1", Horizontal},
		placement{Cruiser, "J1", Vertical},
		placement{Destroyer, "A5", Horizontal},
		placement{Destroyer, "E3", Horizontal},
		placement{Submarine, "H5", Horizontal},
		placement{Submarine, "A8", Horizontal},
		placement{Submarine, "E8", Horizontal},
	)
	return b
}

// fire shoots at the named cells and fails the test if any shot is rejected.
// It returns the last shot.
func fire(t *testing.T, b *Board, cells ...string) Shot {
	t.Helper()
	var shot Shot
	for _, name := range cells {
		var err error
		shot, err = b.Fire(mustCoord(t, name))
		if err != nil {
			t.Fatalf("Fire(%s): %v", name, err)
		}
	}
	return shot
}

// CA-1 and CA-2: a shot on water is a miss; a shot on a ship is a hit that
// does not reveal the ship type.
func TestFire(t *testing.T) {
	tests := []struct {
		cell string
		want Shot
	}{
		{"J10", Shot{Result: Miss}},
		{"E1", Shot{Result: Hit}},
	}
	for _, tt := range tests {
		t.Run(tt.cell, func(t *testing.T) {
			b := exampleBoard(t)
			if got := fire(t, b, tt.cell); got != tt.want {
				t.Errorf("Fire(%s) = %+v, want %+v", tt.cell, got, tt.want)
			}
		})
	}
}

func TestShotResult_String(t *testing.T) {
	tests := []struct {
		r    ShotResult
		want string
	}{
		{Miss, "Miss"},
		{Hit, "Hit"},
		{Sunk, "Sunk"},
		{ShotResult(9), "ShotResult(9)"},
	}
	for _, tt := range tests {
		if got := tt.r.String(); got != tt.want {
			t.Errorf("ShotResult(%d).String() = %q, want %q", int(tt.r), got, tt.want)
		}
	}
}

// CA-3 and CA-4: the shot on the last untouched cell of a ship sinks it and
// reveals the ship type.
func TestFire_Sinks(t *testing.T) {
	t.Run("cruiser after three shots", func(t *testing.T) {
		b := exampleBoard(t)
		for _, cell := range []string{"J1", "J2"} {
			if got := fire(t, b, cell); got != (Shot{Result: Hit}) {
				t.Fatalf("Fire(%s) = %+v, want Hit", cell, got)
			}
		}
		want := Shot{Result: Sunk, Ship: Cruiser}
		if got := fire(t, b, "J3"); got != want {
			t.Errorf("Fire(J3) = %+v, want %+v", got, want)
		}
	})

	t.Run("submarine on the first shot", func(t *testing.T) {
		b := exampleBoard(t)
		want := Shot{Result: Sunk, Ship: Submarine}
		if got := fire(t, b, "H5"); got != want {
			t.Errorf("Fire(H5) = %+v, want %+v", got, want)
		}
	})
}

// cellStates returns the state of every cell on the board, keyed by its name.
func cellStates(b *Board) map[string]CellState {
	states := map[string]CellState{}
	for col := range Size {
		for row := range Size {
			c := Coord{Col: col, Row: row}
			states[c.String()] = b.Cell(c)
		}
	}
	return states
}

// CA-5: the board remembers which cells were shot and whether each shot was a
// miss or a hit.
func TestBoard_Cell(t *testing.T) {
	b := exampleBoard(t)
	fire(t, b, "J10", "E1")

	for name, got := range cellStates(b) {
		want := CellUntouched
		switch name {
		case "J10":
			want = CellMiss
		case "E1":
			want = CellHit
		}
		if got != want {
			t.Errorf("Cell(%s) = %v, want %v", name, got, want)
		}
	}

	for _, c := range []Coord{{Col: Size, Row: 0}, {Col: 0, Row: -1}} {
		if got := b.Cell(c); got != CellUntouched {
			t.Errorf("Cell(%v) = %v, want %v", c, got, CellUntouched)
		}
	}
}

func TestCellState_String(t *testing.T) {
	tests := []struct {
		s    CellState
		want string
	}{
		{CellUntouched, "CellUntouched"},
		{CellMiss, "CellMiss"},
		{CellHit, "CellHit"},
		{CellState(9), "CellState(9)"},
	}
	for _, tt := range tests {
		if got := tt.s.String(); got != tt.want {
			t.Errorf("CellState(%d).String() = %q, want %q", int(tt.s), got, tt.want)
		}
	}
}

// assertFireRejected checks that firing at c fails with want and leaves every
// cell of b unchanged, including c itself.
func assertFireRejected(t *testing.T, b *Board, c Coord, want error) {
	t.Helper()
	before := cellStates(b)
	stateBefore := b.Cell(c)
	shot, err := b.Fire(c)
	if !errors.Is(err, want) {
		t.Fatalf("Fire(%v) = %+v, %v, want error %v", c, shot, err, want)
	}
	if after := cellStates(b); !maps.Equal(before, after) {
		t.Errorf("board changed after rejected shot at %v", c)
	}
	if got := b.Cell(c); got != stateBefore {
		t.Errorf("Cell(%v) = %v after rejected shot, want %v", c, got, stateBefore)
	}
}

// CA-6, CA-8 and CA-9: invalid shots are rejected without changing the board.
func TestFire_Rejected(t *testing.T) {
	t.Run("CA-6 out of bounds", func(t *testing.T) {
		b := exampleBoard(t)
		for _, c := range []Coord{{Col: Size, Row: 0}, {Col: 0, Row: -1}} {
			assertFireRejected(t, b, c, ErrOutOfBounds)
		}
	})

	t.Run("CA-8 repeated shot", func(t *testing.T) {
		b := exampleBoard(t)
		fire(t, b, "J10", "E1")
		for _, name := range []string{"J10", "E1"} {
			assertFireRejected(t, b, mustCoord(t, name), ErrAlreadyShot)
		}
	})

	t.Run("CA-9 incomplete board", func(t *testing.T) {
		b := NewBoard(DefaultFleet())
		place(t, b, placement{Cruiser, "J1", Vertical})
		for _, name := range []string{"J1", "A1"} {
			assertFireRejected(t, b, mustCoord(t, name), ErrIncompleteBoard)
		}
		// The board check comes first, even for a cell outside the board.
		assertFireRejected(t, b, Coord{Col: Size, Row: 0}, ErrIncompleteBoard)
	})
}
