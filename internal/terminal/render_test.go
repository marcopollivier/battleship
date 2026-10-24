package terminal

import (
	"strings"
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
)

func render(t *testing.T, b *board.Board) string {
	t.Helper()
	var sb strings.Builder
	if err := RenderBoard(&sb, b); err != nil {
		t.Fatalf("RenderBoard: %v", err)
	}
	return sb.String()
}

// CA-1: an empty board is drawn like the old program, without trailing
// spaces, followed by a blank line and the legend.
func TestRenderBoard_Empty(t *testing.T) {
	want := `   A B C D E F G H I J
 1 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 2 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 3 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 4 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 5 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 6 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 7 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 8 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 9 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
10 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~

Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões
`
	got := render(t, board.NewBoard(board.DefaultFleet()))
	if got != want {
		t.Errorf("RenderBoard() =\n%s\nwant\n%s", got, want)
	}
	for i, line := range strings.Split(got, "\n") {
		if strings.HasSuffix(line, " ") {
			t.Errorf("line %d has a trailing space: %q", i+1, line)
		}
	}
}
