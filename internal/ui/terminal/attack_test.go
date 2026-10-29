package terminal

import (
	"strings"
	"testing"

	"github.com/marcopollivier/battleship/internal/core/board"
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

// attackRender draws the attack board of b.
func attackRender(t *testing.T, b *board.Board) string {
	t.Helper()
	var sb strings.Builder
	if err := RenderAttackBoard(&sb, b, Plain); err != nil {
		t.Fatalf("RenderAttackBoard: %v", err)
	}
	return sb.String()
}

// attackAssertCells checks the symbol of every cell: the ones in want, and
// "." for all the others.
func attackAssertCells(t *testing.T, output string, want map[string]string) {
	t.Helper()
	for row := range board.Size {
		for col := range board.Size {
			cell := board.Coord{Col: col, Row: row}.String()
			symbol, ok := want[cell]
			if !ok {
				symbol = "."
			}
			if got := cellAt(t, output, cell); got != symbol {
				t.Errorf("cell %s = %q, want %q", cell, got, symbol)
			}
		}
	}
}

// CA-1: without shots every cell is untouched and no ship is revealed.
func TestRenderAttackBoard_NoShots(t *testing.T) {
	want := `   A B C D E F G H I J
 1 . . . . . . . . . .
 2 . . . . . . . . . .
 3 . . . . . . . . . .
 4 . . . . . . . . . .
 5 . . . . . . . . . .
 6 . . . . . . . . . .
 7 . . . . . . . . . .
 8 . . . . . . . . . .
 9 . . . . . . . . . .
10 . . . . . . . . . .

Legenda: . não atingida, o água, X acerto, S D C E P navio afundado
`
	got := attackRender(t, attackExampleBoard(t))
	if got != want {
		t.Errorf("RenderAttackBoard() =\n%s\nwant\n%s", got, want)
	}
}

// CA-2: a miss shows o and a hit on a ship that is still afloat shows X.
func TestRenderAttackBoard_Shots(t *testing.T) {
	b := attackExampleBoard(t)
	attackFire(t, b, "J10", "E1")
	attackAssertCells(t, attackRender(t, b), map[string]string{"J10": "o", "E1": "X"})
}

// CA-5: the attack board is the grid, a blank line and the attack legend,
// with no trailing spaces and only ASCII above the legend.
func TestRenderAttackBoard_Format(t *testing.T) {
	b := attackExampleBoard(t)
	attackFire(t, b, "J10", "E1", "A8", "C5")
	got := attackRender(t, b)

	if want := "Legenda: . não atingida, o água, X acerto, S D C E P navio afundado"; AttackLegend != want {
		t.Errorf("AttackLegend = %q, want %q", AttackLegend, want)
	}
	boardPart, found := strings.CutSuffix(got, "\n\n"+AttackLegend+"\n")
	if !found {
		t.Fatalf("output does not end with a blank line and the attack legend:\n%s", got)
	}
	lines := strings.Split(boardPart, "\n")
	if len(lines) != board.Size+1 {
		t.Errorf("board has %d lines, want %d", len(lines), board.Size+1)
	}
	for i, line := range strings.Split(got, "\n") {
		if strings.HasSuffix(line, " ") {
			t.Errorf("line %d has a trailing space: %q", i+1, line)
		}
	}
	for i := 0; i < len(boardPart); i++ {
		if boardPart[i] >= 128 {
			t.Errorf("byte %d above the legend is not ASCII: %#x", i, boardPart[i])
		}
	}
}

// CA-3: every cell of a sunk ship shows the letter of its type, while a ship
// that is only hit still shows X.
func TestRenderAttackBoard_Sunk(t *testing.T) {
	tests := []struct {
		name  string
		shots []string
		want  map[string]string
	}{
		{
			name:  "cruiser sunk",
			shots: []string{"J1", "J2", "J3"},
			want:  map[string]string{"J1": "C", "J2": "C", "J3": "C"},
		},
		{
			name:  "cruiser hit twice",
			shots: []string{"J1", "J2"},
			want:  map[string]string{"J1": "X", "J2": "X"},
		},
		{
			name:  "submarine sunk next to a miss and a hit",
			shots: []string{"A8", "B8", "E1"},
			want:  map[string]string{"A8": "S", "B8": "o", "E1": "X"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := attackExampleBoard(t)
			attackFire(t, b, tt.shots...)
			attackAssertCells(t, attackRender(t, b), tt.want)
		})
	}
}
