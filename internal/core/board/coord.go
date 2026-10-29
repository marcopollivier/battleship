// Package board holds the Battleship game domain: coordinates, board, fleet and shots.
package board

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Size is the number of rows and columns on the board.
const Size = 10

var (
	// ErrInvalidFormat means the text is not a letter followed by a number without leading zeros.
	ErrInvalidFormat = errors.New("invalid coordinate format")
	// ErrOutOfBounds means the text is well formed but names a cell outside the board.
	ErrOutOfBounds = errors.New("coordinate out of bounds")
)

// coordPattern is a letter followed by a number without leading zeros.
var coordPattern = regexp.MustCompile(`^[A-Za-z](0|[1-9][0-9]*)$`)

// Coord is a cell on the board. Col and Row are zero-based: A1 is {0, 0}, J10 is {9, 9}.
type Coord struct {
	Col int
	Row int
}

// ParseCoord converts text such as "A1" or " j10 " into a Coord.
// Errors wrap ErrInvalidFormat or ErrOutOfBounds and include the original input.
func ParseCoord(s string) (Coord, error) {
	text := strings.TrimSpace(s)
	if !coordPattern.MatchString(text) {
		return Coord{}, fmt.Errorf("%w: %q", ErrInvalidFormat, s)
	}

	letter := strings.ToUpper(text[:1])[0]
	digits := text[1:]
	// Checking the length first keeps strconv.Atoi from overflowing on huge numbers.
	if letter > 'A'+Size-1 || len(digits) > len(strconv.Itoa(Size)) {
		return Coord{}, fmt.Errorf("%w: %q", ErrOutOfBounds, s)
	}
	row, err := strconv.Atoi(digits)
	if err != nil || row < 1 || row > Size {
		return Coord{}, fmt.Errorf("%w: %q", ErrOutOfBounds, s)
	}
	return Coord{Col: int(letter - 'A'), Row: row - 1}, nil
}

// String returns the canonical form, e.g. "A1", "J10".
// A Coord outside the board is rendered as "Coord(col,row)".
func (c Coord) String() string {
	// The explicit range check lets gosec prove the rune conversion is safe.
	if c.Col < 0 || c.Col >= Size || c.Row < 0 || c.Row >= Size {
		return fmt.Sprintf("Coord(%d,%d)", c.Col, c.Row)
	}
	return string(rune('A'+c.Col)) + strconv.Itoa(c.Row+1)
}

func (c Coord) inBounds() bool {
	return c.Col >= 0 && c.Col < Size && c.Row >= 0 && c.Row < Size
}
