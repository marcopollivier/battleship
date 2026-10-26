package terminal

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/marcopollivier/battleship/internal/board"
)

// errInvalidPlacement means the line is not a cell followed by an orientation.
var errInvalidPlacement = errors.New("invalid placement format")

// parsePlacement reads "E5 vertical", "e5 v" or, for a one-cell ship, "E5".
// The orientation is its Portuguese name or first letter, in any case, and
// must apply to t.
func parsePlacement(line string, t board.ShipType) (board.Coord, board.Orientation, error) {
	fields := strings.Fields(line)
	oneCell := t.Size() == 1
	if len(fields) == 0 || len(fields) > 2 || (len(fields) == 1 && !oneCell) {
		return board.Coord{}, 0, fmt.Errorf("%w: %q", errInvalidPlacement, line)
	}
	c, err := board.ParseCoord(fields[0])
	if err != nil {
		return board.Coord{}, 0, err
	}
	if len(fields) == 1 {
		return c, board.Horizontal, nil
	}
	o, ok := orientationByName(fields[1])
	if !ok || !slices.Contains(board.ValidOrientations(t), o) {
		return board.Coord{}, 0, fmt.Errorf("%w: %v %q", board.ErrInvalidOrientation, t, fields[1])
	}
	return c, o, nil
}

// orientationByName finds an orientation by its Portuguese name or first
// letter, in any case.
func orientationByName(name string) (board.Orientation, bool) {
	name = strings.ToLower(name)
	for o, full := range orientationNames {
		if name == full || name == full[:1] {
			return o, true
		}
	}
	return 0, false
}
