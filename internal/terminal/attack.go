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
// ships appear only once hit, and sunk ships show their letter.
func RenderAttackBoard(w io.Writer, b *board.Board) error {
	var grid [board.Size][board.Size]string
	for row := range grid {
		for col := range grid[row] {
			switch b.Cell(board.Coord{Col: col, Row: row}) {
			case board.CellHit:
				grid[row][col] = hitSymbol
			case board.CellMiss:
				grid[row][col] = missSymbol
			default:
				grid[row][col] = untouchedSymbol
			}
		}
	}
	return writeBoard(w, &grid, AttackLegend)
}
