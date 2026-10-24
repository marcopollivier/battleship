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
