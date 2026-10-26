package terminal

import (
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
)

// CA-4: each shot result has its Portuguese message, and only a sunk ship
// is named.
func TestShotMessage(t *testing.T) {
	tests := []struct {
		name string
		shot board.Shot
		want string
	}{
		{"miss", board.Shot{Result: board.Miss}, "Água."},
		{"hit", board.Shot{Result: board.Hit}, "Acerto!"},
		{"hit ignores ship", board.Shot{Result: board.Hit, Ship: board.Cruiser}, "Acerto!"},
		{"miss ignores ship", board.Shot{Result: board.Miss, Ship: board.Battleship}, "Água."},
		{"sunk submarine", board.Shot{Result: board.Sunk, Ship: board.Submarine}, "Afundou o Submarino!"},
		{"sunk destroyer", board.Shot{Result: board.Sunk, Ship: board.Destroyer}, "Afundou o Destroyer!"},
		{"sunk cruiser", board.Shot{Result: board.Sunk, Ship: board.Cruiser}, "Afundou o Cruzador!"},
		{"sunk battleship", board.Shot{Result: board.Sunk, Ship: board.Battleship}, "Afundou o Encouraçado!"},
		{"sunk aircraft carrier", board.Shot{Result: board.Sunk, Ship: board.AircraftCarrier}, "Afundou o Porta-Aviões!"},
		{"unknown result", board.Shot{Result: board.ShotResult(99)}, "ShotResult(99)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShotMessage(tt.shot); got != tt.want {
				t.Errorf("ShotMessage(%+v) = %q, want %q", tt.shot, got, tt.want)
			}
		})
	}
}
