package board

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

var (
	// ErrNegativeCount means a fleet entry has a count below zero.
	ErrNegativeCount = errors.New("fleet has a negative ship count")
	// ErrEmptyFleet means the fleet has no ships at all.
	ErrEmptyFleet = errors.New("fleet has no ships")
)

// A fleet that cannot be placed reuses ErrFleetDoesNotFit from random.go.

// fitCheckSeed is the fixed seed used to decide whether a fleet fits, so the
// same fleet always gets the same answer.
const fitCheckSeed = 1

// ValidateFleet reports whether fleet can be used in a game: no negative
// counts, at least one ship, and a board can be built with it.
func ValidateFleet(fleet []FleetEntry) error {
	total := 0
	for _, entry := range fleet {
		if entry.Count < 0 {
			return fmt.Errorf("%w: %v %d", ErrNegativeCount, entry.Type, entry.Count)
		}
		total += entry.Count
	}
	if total == 0 {
		return ErrEmptyFleet
	}
	rng := rand.New(rand.NewPCG(fitCheckSeed, 0)) //nolint:gosec // deterministic fit check, not security-sensitive
	if _, err := RandomBoard(fleet, rng); err != nil {
		return err
	}
	return nil
}
