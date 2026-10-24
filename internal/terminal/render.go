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

// shipSymbols are the ASCII letters drawn on the cells of each ship type.
var shipSymbols = map[board.ShipType]string{
	board.Submarine:       "S",
	board.Destroyer:       "D",
	board.Cruiser:         "C",
	board.Battleship:      "E",
	board.AircraftCarrier: "P",
}

// RenderBoard writes the player's own board, followed by a blank line and the legend.
func RenderBoard(w io.Writer, b *board.Board) error {
	var grid [board.Size][board.Size]string
	for row := range grid {
		for col := range grid[row] {
			grid[row][col] = water
		}
	}
	for _, ship := range b.Ships() {
		for _, c := range ship.Cells {
			grid[c.Row][c.Col] = shipSymbol(ship.Type)
		}
	}

	var sb strings.Builder

	sb.WriteString("  ")
	for col := range board.Size {
		sb.WriteString(" " + columnName(col))
	}
	sb.WriteString("\n")

	for row := range board.Size {
		fmt.Fprintf(&sb, "%2d", row+1)
		for col := range board.Size {
			sb.WriteString(" " + grid[row][col])
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

// shipSymbol returns the letter drawn for a ship type, or "?" for an unknown type.
func shipSymbol(t board.ShipType) string {
	if symbol, ok := shipSymbols[t]; ok {
		return symbol
	}
	return "?"
}
