package board

import (
	"cmp"
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
)

// ErrFleetDoesNotFit means no valid board was found for the fleet within the attempt limits.
var ErrFleetDoesNotFit = errors.New("fleet does not fit on the board")

const (
	// maxAttemptsPerShip is how many random spots are tried for one ship
	// before the whole board is discarded.
	maxAttemptsPerShip = 200
	// maxRestarts is how many times the board is started over before giving up.
	maxRestarts = 50
)

// RandomBoard returns a complete board with the fleet placed at random.
// The same rng state always produces the same board.
func RandomBoard(fleet []FleetEntry, rng *rand.Rand) (*Board, error) {
	b := NewBoard(fleet)
	if err := FillRandom(b, rng); err != nil {
		return nil, err
	}
	return b, nil
}

// FillRandom places the ships that are still missing at random, keeping the
// ships already on b. On ErrFleetDoesNotFit, b is left unchanged.
func FillRandom(b *Board, rng *rand.Rand) error {
	ships := b.Remaining()
	for _, t := range ships {
		// A type with no orientations can never be placed, so no attempt is spent on it.
		if len(ValidOrientations(t)) == 0 {
			return fmt.Errorf("%w: unknown ship type %v", ErrFleetDoesNotFit, t)
		}
	}
	for range maxRestarts {
		if filled, ok := tryRandomBoard(b, ships, rng); ok {
			b.ships = filled.ships
			return nil
		}
	}
	return ErrFleetDoesNotFit
}

// shipsBySize lists one entry per ship in the fleet, largest first. The sort
// is stable so the order, and therefore the board, depends only on rng.
func shipsBySize(fleet []FleetEntry) []ShipType {
	entries := slices.Clone(fleet)
	slices.SortStableFunc(entries, func(a, b FleetEntry) int {
		return cmp.Compare(b.Type.Size(), a.Type.Size())
	})
	var ships []ShipType
	for _, entry := range entries {
		for range entry.Count {
			ships = append(ships, entry.Type)
		}
	}
	return ships
}

// tryRandomBoard places every ship in order on a copy of base and reports
// false as soon as one of them cannot be placed.
func tryRandomBoard(base *Board, ships []ShipType, rng *rand.Rand) (*Board, bool) {
	b := base.clone()
	for _, t := range ships {
		if !placeRandomly(b, t, rng) {
			return nil, false
		}
	}
	return b, true
}

// placeRandomly tries random spots and orientations for a ship of type t
// until Place accepts one or the attempts run out.
func placeRandomly(b *Board, t ShipType, rng *rand.Rand) bool {
	orientations := ValidOrientations(t)
	for range maxAttemptsPerShip {
		ref := Coord{Col: rng.IntN(Size), Row: rng.IntN(Size)}
		o := orientations[rng.IntN(len(orientations))]
		if b.Place(t, ref, o) == nil {
			return true
		}
	}
	return false
}
