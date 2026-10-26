package board

import (
	"errors"
	"slices"
	"testing"
)

// machineShootAll fires RandomShot until no cell is left and returns the cells in order.
func machineShootAll(t *testing.T, seed uint64) []Coord {
	t.Helper()
	b := mustRandomBoard(t, 1)
	rng := seeded(seed)
	var shots []Coord
	for range Size * Size {
		c, err := RandomShot(b, rng)
		if err != nil {
			t.Fatalf("RandomShot after %d shots: %v", len(shots), err)
		}
		if _, err := b.Fire(c); err != nil {
			t.Fatalf("Fire(%v): %v", c, err)
		}
		shots = append(shots, c)
	}
	if _, err := RandomShot(b, rng); !errors.Is(err, ErrNoCellsLeft) {
		t.Errorf("RandomShot on a fully shot board error = %v, want ErrNoCellsLeft", err)
	}
	return shots
}

// CA-4: the machine only shoots at cells it has not shot yet, and the same
// seed gives the same sequence.
func TestRandomShot(t *testing.T) {
	first := machineShootAll(t, 7)
	seen := map[Coord]bool{}
	for _, c := range first {
		if seen[c] {
			t.Errorf("RandomShot chose %v twice", c)
		}
		seen[c] = true
	}
	if len(seen) != Size*Size {
		t.Errorf("RandomShot covered %d cells, want %d", len(seen), Size*Size)
	}

	if again := machineShootAll(t, 7); !slices.Equal(first, again) {
		t.Error("the same seed gave a different sequence of shots")
	}
	if other := machineShootAll(t, 8); slices.Equal(first, other) {
		t.Error("seeds 7 and 8 gave the same sequence of shots")
	}
}
