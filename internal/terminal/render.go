package terminal

import (
	"fmt"
	"io"
	"strings"

	"github.com/marcopollivier/battleship/internal/board"
)

// Legend explains the board symbols.
const Legend = "Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões"

// water is the symbol of a cell without a ship.
const water = "~"

// RenderBoard writes the player's own board, followed by a blank line and the legend.
func RenderBoard(w io.Writer, b *board.Board) error {
	var sb strings.Builder

	sb.WriteString("  ")
	for col := range board.Size {
		sb.WriteString(" " + columnName(col))
	}
	sb.WriteString("\n")

	for row := range board.Size {
		fmt.Fprintf(&sb, "%2d", row+1)
		for range board.Size {
			sb.WriteString(" " + water)
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n" + Legend + "\n")

	_, err := w.Write([]byte(sb.String()))
	return err
}

// columnName returns the letter of a column, reusing the coordinate naming
// rule from the domain instead of converting the index to a rune here.
func columnName(col int) string {
	return strings.TrimSuffix(board.Coord{Col: col}.String(), "1")
}
