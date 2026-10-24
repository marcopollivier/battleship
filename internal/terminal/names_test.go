package terminal

import (
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
)

// CA-4: every ship type has its Portuguese name.
func TestShipName(t *testing.T) {
	tests := []struct {
		shipType board.ShipType
		want     string
	}{
		{board.Submarine, "Submarino"},
		{board.Destroyer, "Destroyer"},
		{board.Cruiser, "Cruzador"},
		{board.Battleship, "Encouraçado"},
		{board.AircraftCarrier, "Porta-Aviões"},
		{board.ShipType(99), "ShipType(99)"},
	}
	for _, tt := range tests {
		if got := ShipName(tt.shipType); got != tt.want {
			t.Errorf("ShipName(%v) = %q, want %q", tt.shipType, got, tt.want)
		}
	}
}

// CA-4: every orientation has its Portuguese name.
func TestOrientationName(t *testing.T) {
	tests := []struct {
		orientation board.Orientation
		want        string
	}{
		{board.Horizontal, "horizontal"},
		{board.Vertical, "vertical"},
		{board.Up, "cima"},
		{board.Down, "baixo"},
		{board.Left, "esquerda"},
		{board.Right, "direita"},
		{board.Orientation(99), "Orientation(99)"},
	}
	for _, tt := range tests {
		if got := OrientationName(tt.orientation); got != tt.want {
			t.Errorf("OrientationName(%v) = %q, want %q", tt.orientation, got, tt.want)
		}
	}
}
