package terminal

import (
	"io"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/marcopollivier/battleship/internal/core/board"
)

// Spec 022 T1: the grids are the header and the 10 rows that RenderBoard and
// RenderAttackBoard write before the blank line and the legend, each line
// gridWidth columns wide without the color codes.
func TestGrids(t *testing.T) {
	b := attackExampleBoard(t)
	attackFire(t, b, "J10", "E1", "A1")
	tests := []struct {
		name   string
		grid   func(*board.Board, Style) []string
		draw   func(io.Writer, *board.Board, Style) error
		legend string
	}{
		{"own board", ownGrid, RenderBoard, Legend},
		{"attack board", attackGrid, RenderAttackBoard, AttackLegend},
	}
	for _, tt := range tests {
		for _, style := range []Style{Plain, Colored} {
			lines := tt.grid(b, style)
			if len(lines) != board.Size+1 {
				t.Fatalf("%s: %d lines, want %d", tt.name, len(lines), board.Size+1)
			}
			for i, line := range lines {
				if n := utf8.RuneCountInString(stripColors(line)); n != gridWidth {
					t.Errorf("%s: line %d is %d columns wide, want %d: %q", tt.name, i+1, n, gridWidth, line)
				}
			}
			var sb strings.Builder
			if err := tt.draw(&sb, b, style); err != nil {
				t.Fatalf("draw: %v", err)
			}
			if got, want := strings.Join(lines, "\n")+"\n\n"+tt.legend+"\n", sb.String(); got != want {
				t.Errorf("%s: grid + legend =\n%s\nwant\n%s", tt.name, got, want)
			}
		}
	}
}
