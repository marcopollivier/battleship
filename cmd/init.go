// Package cmd wires the game domain to the terminal.
package cmd

import (
	"io"
	"os"

	"github.com/marcopollivier/battleship/internal/board"
	"github.com/marcopollivier/battleship/internal/terminal"
)

// Run writes the game screen to w.
func Run(w io.Writer) error {
	return terminal.RenderBoard(w, board.NewBoard(board.DefaultFleet()))
}

// Init runs the game on standard output.
func Init() error {
	return Run(os.Stdout)
}
