package terminal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/marcopollivier/battleship/internal/board"
)

// ErrInvalidFleetFormat means the fleet text is not a comma-separated list of letter=count pairs.
var ErrInvalidFleetFormat = errors.New("invalid fleet format")

// ParseFleet reads a fleet such as "S=2,P=1" using the board letters
// (S, D, C, E, P). Types not listed get no ships. It does not validate
// counts; use board.ValidateFleet for that.
func ParseFleet(s string) ([]board.FleetEntry, error) {
	counts := map[board.ShipType]int{}
	for _, pair := range strings.Split(s, ",") {
		letter, value, ok := strings.Cut(strings.TrimSpace(pair), "=")
		if !ok {
			return nil, fmt.Errorf("%w: missing '=' in %q", ErrInvalidFleetFormat, pair)
		}
		shipType, ok := shipTypeBySymbol(strings.ToUpper(letter))
		if !ok {
			return nil, fmt.Errorf("%w: unknown ship letter %q", ErrInvalidFleetFormat, letter)
		}
		if _, repeated := counts[shipType]; repeated {
			return nil, fmt.Errorf("%w: repeated ship letter %q", ErrInvalidFleetFormat, letter)
		}
		count, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid count %q", ErrInvalidFleetFormat, value)
		}
		counts[shipType] = count
	}

	var fleet []board.FleetEntry
	for t := board.Submarine; t <= board.AircraftCarrier; t++ {
		if count, ok := counts[t]; ok {
			fleet = append(fleet, board.FleetEntry{Type: t, Count: count})
		}
	}
	return fleet, nil
}

// shipTypeBySymbol returns the ship type drawn with symbol, the inverse of shipSymbols.
func shipTypeBySymbol(symbol string) (board.ShipType, bool) {
	for t, s := range shipSymbols {
		if s == symbol {
			return t, true
		}
	}
	return 0, false
}
