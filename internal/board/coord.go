// Package board holds the Battleship game domain: coordinates, board, fleet and shots.
package board

import (
	"fmt"
	"strconv"
	"strings"
)

// Size is the number of rows and columns on the board.
const Size = 10

// Coord is a cell on the board. Col and Row are zero-based: A1 is {0, 0}, J10 is {9, 9}.
type Coord struct {
	Col int
	Row int
}

// ParseCoord converts text such as "A1" or " j10 " into a Coord.
func ParseCoord(s string) (Coord, error) {
	text := strings.TrimSpace(s)
	if len(text) < 2 {
		return Coord{}, fmt.Errorf("parse coordinate %q", s)
	}
	row, err := strconv.Atoi(text[1:])
	if err != nil {
		return Coord{}, fmt.Errorf("parse coordinate %q: %w", s, err)
	}
	letter := strings.ToUpper(text[:1])[0]
	return Coord{Col: int(letter - 'A'), Row: row - 1}, nil
}

// String returns the canonical form, e.g. "A1", "J10".
// A Coord outside the board is rendered as "Coord(col,row)".
func (c Coord) String() string {
	if c.Col < 0 || c.Col >= Size || c.Row < 0 || c.Row >= Size {
		return fmt.Sprintf("Coord(%d,%d)", c.Col, c.Row)
	}
	return string(rune('A'+c.Col)) + strconv.Itoa(c.Row+1)
}
