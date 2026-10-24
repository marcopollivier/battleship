package board

import (
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
