package text

import (
	"strings"
	"testing"

	"github.com/marcopollivier/battleship/internal/core/board"
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

func TestShipSymbol(t *testing.T) {
	tests := []struct {
		shipType board.ShipType
		want     string
	}{
		{board.Submarine, "S"},
		{board.Destroyer, "D"},
		{board.Cruiser, "C"},
		{board.Battleship, "E"},
		{board.AircraftCarrier, "P"},
		{board.ShipType(99), "?"},
	}
	for _, tt := range tests {
		if got := ShipSymbol(tt.shipType); got != tt.want {
			t.Errorf("ShipSymbol(%v) = %q, want %q", tt.shipType, got, tt.want)
		}
	}
}

func TestOrientationByName(t *testing.T) {
	tests := []struct {
		name string
		want board.Orientation
		ok   bool
	}{
		{"vertical", board.Vertical, true},
		{"V", board.Vertical, true},
		{"Cima", board.Up, true},
		{"e", board.Left, true},
		{"diagonal", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		got, ok := OrientationByName(tt.name)
		if got != tt.want || ok != tt.ok {
			t.Errorf("OrientationByName(%q) = %v, %v, want %v, %v", tt.name, got, ok, tt.want, tt.ok)
		}
	}
}

func TestOrientationList(t *testing.T) {
	tests := []struct {
		shipType board.ShipType
		want     string
	}{
		{board.Cruiser, "horizontal vertical"},
		{board.AircraftCarrier, "cima baixo esquerda direita"},
	}
	for _, tt := range tests {
		if got := strings.Join(OrientationList(tt.shipType), " "); got != tt.want {
			t.Errorf("OrientationList(%v) = %q, want %q", tt.shipType, got, tt.want)
		}
	}
}
