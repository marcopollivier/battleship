package terminal

import (
	"io"

	"github.com/marcopollivier/battleship/internal/board"
)

// AttackLegend explains the symbols of the attack board.
const AttackLegend = "Legenda: . não atingida, o água, X acerto, S D C E P navio afundado"

// untouchedSymbol is the symbol of a cell the attacker has not shot yet.
const untouchedSymbol = "."

// RenderAttackBoard writes the opponent's board as the attacker sees it:
// ships appear only once hit, and sunk ships show their letter. With
// Colored, each cell is drawn in the color of its symbol.
func RenderAttackBoard(w io.Writer, b *board.Board, style Style) error {
	var grid [board.Size][board.Size]string
	for row := range grid {
		for col := range grid[row] {
			switch b.Cell(board.Coord{Col: col, Row: row}) {
			case board.CellHit:
				grid[row][col] = style.paint(hitSymbol, hitColor)
			case board.CellMiss:
				grid[row][col] = style.paint(missSymbol, missColor)
			default:
				grid[row][col] = style.paint(untouchedSymbol, untouchedColor)
			}
		}
	}
	for _, ship := range b.Ships() {
		if !shipSunk(b, ship) {
			continue
		}
		for _, c := range ship.Cells {
			grid[c.Row][c.Col] = style.paint(shipSymbol(ship.Type), sunkColor)
		}
	}
	return writeBoard(w, &grid, AttackLegend)
}

// shipSunk reports whether every cell of ship has been hit.
func shipSunk(b *board.Board, ship board.PlacedShip) bool {
	for _, c := range ship.Cells {
		if b.Cell(c) != board.CellHit {
			return false
		}
	}
	return true
}
