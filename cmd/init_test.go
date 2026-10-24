package cmd

import (
	"strings"
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
	"github.com/marcopollivier/battleship/internal/terminal"
)

// CA-5: the program prints the board drawn by the terminal layer.
func TestRun(t *testing.T) {
	var want strings.Builder
	if err := terminal.RenderBoard(&want, board.NewBoard(board.DefaultFleet())); err != nil {
		t.Fatalf("RenderBoard: %v", err)
	}

	var got strings.Builder
	if err := Run(&got); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got.String() != want.String() {
		t.Errorf("Run() =\n%s\nwant\n%s", got.String(), want.String())
	}
}
