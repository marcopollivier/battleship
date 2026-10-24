package board

import (
	"math/rand/v2"
	"testing"
)

// seeded returns a random source with a fixed seed, so every run is the same.
func seeded(seed uint64) *rand.Rand {
	return rand.New(rand.NewPCG(seed, 0))
}

// assertValidBoard checks the placement rules from spec 002 straight from
// Ships(), without going through Place: every cell on the board, no shared
// cells, no ship touching another (including diagonally) and the number of
// ships of each type matching the fleet.
func assertValidBoard(t *testing.T, fleet []FleetEntry, b *Board) {
	t.Helper()
	ships := b.Ships()
	owner := map[Coord]int{}
	for i, s := range ships {
		if len(s.Cells) != s.Type.Size() {
			t.Errorf("%v has %d cells, want %d", s.Type, len(s.Cells), s.Type.Size())
		}
		for _, c := range s.Cells {
			if c.Col < 0 || c.Col >= Size || c.Row < 0 || c.Row >= Size {
				t.Errorf("%v has cell %v outside the board", s.Type, c)
			}
			if j, taken := owner[c]; taken {
				t.Errorf("%v and %v share cell %v", s.Type, ships[j].Type, c)
			}
			owner[c] = i
		}
	}
	assertNoTouching(t, ships, owner)
	assertFleetCounts(t, fleet, ships)
}

// assertNoTouching checks that no cell of a ship has a neighbor, in any of
// the 8 directions, that belongs to another ship.
func assertNoTouching(t *testing.T, ships []PlacedShip, owner map[Coord]int) {
	t.Helper()
	for i, s := range ships {
		for _, c := range s.Cells {
			for dCol := -1; dCol <= 1; dCol++ {
				for dRow := -1; dRow <= 1; dRow++ {
					n := Coord{Col: c.Col + dCol, Row: c.Row + dRow}
					if j, taken := owner[n]; taken && j != i {
						t.Errorf("%v at %v touches %v at %v", s.Type, c, ships[j].Type, n)
					}
				}
			}
		}
	}
}

// assertFleetCounts checks that the board has exactly as many ships of each
// type as the fleet lists.
func assertFleetCounts(t *testing.T, fleet []FleetEntry, ships []PlacedShip) {
	t.Helper()
	want := map[ShipType]int{}
	for _, entry := range fleet {
		want[entry.Type] += entry.Count
	}
	got := map[ShipType]int{}
	for _, s := range ships {
		got[s.Type]++
	}
	for shipType, count := range want {
		if got[shipType] != count {
			t.Errorf("board has %d %v, want %d", got[shipType], shipType, count)
		}
	}
	for shipType, count := range got {
		if _, ok := want[shipType]; !ok {
			t.Errorf("board has %d %v, which is not in the fleet", count, shipType)
		}
	}
}

// mustRandomBoard generates a board for the default fleet with the given seed.
func mustRandomBoard(t *testing.T, seed uint64) *Board {
	t.Helper()
	b, err := RandomBoard(DefaultFleet(), seeded(seed))
	if err != nil {
		t.Fatalf("RandomBoard(DefaultFleet(), seed %d): %v", seed, err)
	}
	return b
}

// CA-1: a generated board has the whole default fleet.
func TestRandomBoard_Complete(t *testing.T) {
	b := mustRandomBoard(t, 1)
	if !b.Complete() {
		t.Error("RandomBoard returned an incomplete board")
	}
	if n := len(b.Ships()); n != 8 {
		t.Errorf("RandomBoard placed %d ships, want 8", n)
	}
}

// CA-2: 1,000 generated boards all follow the placement rules.
func TestRandomBoard_Valid(t *testing.T) {
	for seed := uint64(1); seed <= 1000; seed++ {
		b := mustRandomBoard(t, seed)
		assertValidBoard(t, DefaultFleet(), b)
		if t.Failed() {
			t.Fatalf("seed %d produced an invalid board: %v", seed, b.Ships())
		}
	}
}
