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

// gridWidth is the visible width of every line of a board grid, without the
// color codes: 2 columns for the row number and 2 for each cell.
const gridWidth = 2 + 2*board.Size

// RenderBoard writes the player's own board, including the shots it has
// received, followed by a blank line and the legend. The cells of a sunk
// ship are drawn as #. With Colored, each cell is drawn in the color of its
// symbol.
func RenderBoard(w io.Writer, b *board.Board, style Style) error {
	return writeBoard(w, ownGrid(b, style), Legend)
}

// ownGrid returns the lines of the player's own board, the header and one
// line per row, with the ships and the shots received.
func ownGrid(b *board.Board, style Style) []string {
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
	return gridLines(&grid)
}

// gridLines returns the header and one line per row of grid, the cells
// separated by spaces.
func gridLines(grid *[board.Size][board.Size]string) []string {
	var sb strings.Builder
	sb.WriteString("  ")
	for col := range board.Size {
		sb.WriteString(" " + columnName(col))
	}
	lines := []string{sb.String()}

	for row := range board.Size {
		sb.Reset()
		fmt.Fprintf(&sb, "%2d", row+1)
		for col := range board.Size {
			sb.WriteString(" " + grid[row][col])
		}
		lines = append(lines, sb.String())
	}
	return lines
}

// writeBoard writes the lines of a grid, a blank line and legend as a
// single Write.
func writeBoard(w io.Writer, lines []string, legend string) error {
	_, err := io.WriteString(w, strings.Join(lines, "\n")+"\n\n"+legend+"\n")
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
