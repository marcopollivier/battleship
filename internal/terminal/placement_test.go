package terminal

import (
	"errors"
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
)

// CA-3: a cell and an orientation by name or first letter, in any case; a
// one-cell ship needs no orientation.
func TestParsePlacement(t *testing.T) {
	tests := []struct {
		line     string
		shipType board.ShipType
		want     string
		wantO    board.Orientation
	}{
		{"e5 V", board.Battleship, "E5", board.Vertical},
		{"E5 vertical", board.Battleship, "E5", board.Vertical},
		{" a1   HORIZONTAL ", board.Cruiser, "A1", board.Horizontal},
		{"C3 b", board.AircraftCarrier, "C3", board.Down},
		{"C3 BAIXO", board.AircraftCarrier, "C3", board.Down},
		{"C3 cima", board.AircraftCarrier, "C3", board.Up},
		{"C3 e", board.AircraftCarrier, "C3", board.Left},
		{"C3 Direita", board.AircraftCarrier, "C3", board.Right},
		{"A1", board.Submarine, "A1", board.Horizontal},
		{"A1 v", board.Submarine, "A1", board.Vertical},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			c, o, err := parsePlacement(tt.line, tt.shipType)
			if err != nil {
				t.Fatalf("parsePlacement(%q, %v) error = %v", tt.line, tt.shipType, err)
			}
			if c.String() != tt.want || o != tt.wantO {
				t.Errorf("parsePlacement(%q, %v) = %v %v, want %s %v", tt.line, tt.shipType, c, o, tt.want, tt.wantO)
			}
		})
	}
}

// CA-4 (format and orientation): malformed lines and orientations that do
// not exist or do not apply to the ship are rejected.
func TestParsePlacement_Invalid(t *testing.T) {
	tests := []struct {
		line     string
		shipType board.ShipType
		want     error
	}{
		{"", board.Cruiser, errInvalidPlacement},
		{"E5", board.Cruiser, errInvalidPlacement},
		{"E5 vertical sim", board.Cruiser, errInvalidPlacement},
		{"5E vertical", board.Cruiser, board.ErrInvalidFormat},
		{"K1 vertical", board.Cruiser, board.ErrOutOfBounds},
		{"E5 diagonal", board.Cruiser, board.ErrInvalidOrientation},
		{"E5 baixo", board.Cruiser, board.ErrInvalidOrientation},
		{"C3 horizontal", board.AircraftCarrier, board.ErrInvalidOrientation},
		{"C3 x", board.AircraftCarrier, board.ErrInvalidOrientation},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			_, _, err := parsePlacement(tt.line, tt.shipType)
			if !errors.Is(err, tt.want) {
				t.Errorf("parsePlacement(%q, %v) error = %v, want %v", tt.line, tt.shipType, err, tt.want)
			}
		})
	}
}
