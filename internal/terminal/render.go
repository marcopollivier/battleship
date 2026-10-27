package terminal

import (
	"fmt"
	"io"
	"strings"

	"github.com/marcopollivier/battleship/internal/board"
)

// Legend explains the symbols of the player's own board.
const Legend = "Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões, X atingido, # afundado, o tiro na água"

const (
	// water is the symbol of a cell without a ship.
	water = "~"
	// hitSymbol is the symbol of a ship cell that was shot.
	hitSymbol = "X"
	// sunkSymbol is the symbol of a cell of a sunk ship on the own board.
	sunkSymbol = "#"
	// missSymbol is the symbol of a water cell that was shot.
	missSymbol = "o"
)

// shipSymbols are the ASCII letters drawn on the cells of each ship type.
var shipSymbols = map[board.ShipType]string{
	board.Submarine:       "S",
	board.Destroyer:       "D",
	board.Cruiser:         "C",
	board.Battleship:      "E",
	board.AircraftCarrier: "P",
}

// RenderBoard writes the player's own board, including the shots it has
// received, followed by a blank line and the legend. The cells of a sunk
// ship are drawn as #. With Colored, each cell is drawn in the color of its
// symbol.
func RenderBoard(w io.Writer, b *board.Board, style Style) error {
	var grid [board.Size][board.Size]string
	for row := range grid {
		for col := range grid[row] {
			grid[row][col] = style.paint(water, waterColor)
		}
	}
	for _, ship := range b.Ships() {
		for _, c := range ship.Cells {
			grid[c.Row][c.Col] = style.paint(shipSymbol(ship.Type), shipColor)
		}
	}
	for row := range grid {
		for col := range grid[row] {
			switch b.Cell(board.Coord{Col: col, Row: row}) {
			case board.CellHit:
				grid[row][col] = style.paint(hitSymbol, hitColor)
			case board.CellMiss:
				grid[row][col] = style.paint(missSymbol, missColor)
			}
		}
	}
	for _, ship := range b.Ships() {
		if !shipSunk(b, ship) {
			continue
		}
		for _, c := range ship.Cells {
			grid[c.Row][c.Col] = style.paint(sunkSymbol, sunkColor)
		}
	}
	return writeBoard(w, &grid, Legend)
}

// writeBoard writes the header, the rows of grid, a blank line and legend
// as a single Write.
func writeBoard(w io.Writer, grid *[board.Size][board.Size]string, legend string) error {
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

	sb.WriteString("\n" + legend + "\n")

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
