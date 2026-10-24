package terminal

import (
	"errors"
	"slices"
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
)

// CA-1, CA-2: a valid fleet text lists only the given types, in type order.
func TestParseFleet(t *testing.T) {
	sp := []board.FleetEntry{
		{Type: board.Submarine, Count: 2},
		{Type: board.AircraftCarrier, Count: 1},
	}
	tests := []struct {
		name string
		in   string
		want []board.FleetEntry
	}{
		{"CA-1 S=2,P=1", "S=2,P=1", sp},
		{"CA-2 lower case and spaces", " s=2 , p=1 ", sp},
		{"out of type order", "P=1,S=2", sp},
		{"every type", "P=1,E=1,C=1,D=2,S=3", board.DefaultFleet()},
		{"zero and negative counts are kept", "S=0,D=-1", []board.FleetEntry{
			{Type: board.Submarine, Count: 0},
			{Type: board.Destroyer, Count: -1},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFleet(tt.in)
			if err != nil {
				t.Fatalf("ParseFleet(%q): %v", tt.in, err)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("ParseFleet(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

// CA-4: a malformed fleet text is rejected with ErrInvalidFleetFormat.
func TestParseFleet_InvalidFormat(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"empty", ""},
		{"only spaces", "   "},
		{"unknown letter", "X=1"},
		{"count is not a number", "S=a"},
		{"missing equals sign", "S2"},
		{"wrong separator", "S=1;D=1"},
		{"repeated letter", "S=1,S=2"},
		{"repeated letter in other case", "S=1,s=2"},
		{"empty pair", "S=1,,D=1"},
		{"trailing comma", "S=1,"},
		{"missing count", "S="},
		{"missing letter", "=1"},
		{"word instead of letter", "SS=1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFleet(tt.in)
			if !errors.Is(err, ErrInvalidFleetFormat) {
				t.Errorf("ParseFleet(%q) error = %v, want ErrInvalidFleetFormat", tt.in, err)
			}
			if got != nil {
				t.Errorf("ParseFleet(%q) = %v with the error, want nil", tt.in, got)
			}
		})
	}
}
