package board

import (
	"errors"
	"testing"
)

// CA-5, CA-6: a fleet with a negative count or no ships is rejected.
func TestValidateFleet(t *testing.T) {
	tests := []struct {
		name  string
		fleet []FleetEntry
		want  error
	}{
		{"negative count", []FleetEntry{{Type: Submarine, Count: -1}}, ErrNegativeCount},
		{"negative count with other ships", []FleetEntry{{Type: Submarine, Count: 2}, {Type: Destroyer, Count: -1}}, ErrNegativeCount},
		{"zero ships", []FleetEntry{{Type: Submarine, Count: 0}}, ErrEmptyFleet},
		{"empty fleet", nil, ErrEmptyFleet},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateFleet(tt.fleet); !errors.Is(err, tt.want) {
				t.Errorf("ValidateFleet(%v) = %v, want %v", tt.fleet, err, tt.want)
			}
		})
	}
}

// CA-7: a fleet that does not fit on the board is rejected.
func TestValidateFleet_DoesNotFit(t *testing.T) {
	fleet := []FleetEntry{{Type: Battleship, Count: 30}}
	if err := ValidateFleet(fleet); !errors.Is(err, ErrFleetDoesNotFit) {
		t.Errorf("ValidateFleet(%v) = %v, want ErrFleetDoesNotFit", fleet, err)
	}
}

// CA-8: an accepted fleet always produces a complete, valid random board.
func TestValidateFleet_Accepted(t *testing.T) {
	tests := []struct {
		name  string
		fleet []FleetEntry
	}{
		{"default fleet", DefaultFleet()},
		{"S=2,P=1", []FleetEntry{{Type: Submarine, Count: 2}, {Type: AircraftCarrier, Count: 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateFleet(tt.fleet); err != nil {
				t.Fatalf("ValidateFleet(%v) = %v, want nil", tt.fleet, err)
			}
			for seed := uint64(1); seed <= 100; seed++ {
				b, err := RandomBoard(tt.fleet, seeded(seed))
				if err != nil {
					t.Fatalf("RandomBoard(%v, seed %d): %v", tt.fleet, seed, err)
				}
				if !b.Complete() {
					t.Fatalf("seed %d produced an incomplete board: %v", seed, b.Ships())
				}
				assertValidBoard(t, tt.fleet, b)
			}
		})
	}
}
