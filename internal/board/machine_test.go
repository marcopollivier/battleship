package board

import (
	"errors"
	"maps"
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

// smartSeeds is how many seeds the SmartShot tests try on one board.
const smartSeeds = 200

// smartChoices runs SmartShot on b with smartSeeds seeds, without firing,
// and returns the set of cells it chose. Every choice must be a cell that
// has not been shot yet.
func smartChoices(t *testing.T, b *Board) map[Coord]bool {
	t.Helper()
	chosen := map[Coord]bool{}
	for seed := range uint64(smartSeeds) {
		c, err := SmartShot(b, seeded(seed))
		if err != nil {
			t.Fatalf("SmartShot(seed %d): %v", seed, err)
		}
		if b.Cell(c) != CellUntouched {
			t.Errorf("SmartShot(seed %d) chose %v, which was already shot", seed, c)
		}
		chosen[c] = true
	}
	return chosen
}

// smartCells parses cell names into a set.
func smartCells(t *testing.T, names ...string) map[Coord]bool {
	t.Helper()
	cells := map[Coord]bool{}
	for _, name := range names {
		cells[mustCoord(t, name)] = true
	}
	return cells
}

// smartAround returns the cells on the board that touch any of cells,
// diagonals included, and are not in cells.
func smartAround(cells map[Coord]bool) map[Coord]bool {
	around := map[Coord]bool{}
	for c := range cells {
		for dCol := -1; dCol <= 1; dCol++ {
			for dRow := -1; dRow <= 1; dRow++ {
				n := Coord{Col: c.Col + dCol, Row: c.Row + dRow}
				if n.inBounds() && !cells[n] {
					around[n] = true
				}
			}
		}
	}
	return around
}

// smartHunts fails the test if SmartShot chose a cell next to the sunk
// cells, or chose too few different cells to be a random hunt.
func smartHunts(t *testing.T, chosen, sunk map[Coord]bool) {
	t.Helper()
	around := smartAround(sunk)
	for c := range chosen {
		if around[c] {
			t.Errorf("SmartShot chose %v, next to a sunk ship", c)
		}
	}
	if len(chosen) < smartSeeds/4 {
		t.Errorf("SmartShot chose only %d different cells in %d seeds, want a random hunt", len(chosen), smartSeeds)
	}
}

// smartShootAllBut shoots every cell of b except the ones in keep, in board order.
func smartShootAllBut(t *testing.T, b *Board, keep ...string) {
	t.Helper()
	kept := smartCells(t, keep...)
	for row := range Size {
		for col := range Size {
			if c := (Coord{Col: col, Row: row}); !kept[c] {
				fire(t, b, c.String())
			}
		}
	}
}

// CA-1: with no hit on a ship still afloat, SmartShot draws among the cells
// not shot yet, and never next to a sunk ship.
func TestSmartShot_Hunt(t *testing.T) {
	t.Run("no hits", func(t *testing.T) {
		b := exampleBoard(t)
		fire(t, b, "J10", "A10", "D7")
		smartHunts(t, smartChoices(t, b), nil)
	})

	t.Run("a sunk Submarine", func(t *testing.T) {
		b := exampleBoard(t)
		fire(t, b, "J10", "H5")
		smartHunts(t, smartChoices(t, b), smartCells(t, "H5"))
	})

	t.Run("only cells next to a sunk ship are left", func(t *testing.T) {
		b := NewBoard([]FleetEntry{{Type: Submarine, Count: 1}})
		place(t, b, placement{Submarine, "A1", Horizontal})
		smartShootAllBut(t, b, "B2")
		if c, err := SmartShot(b, seeded(1)); err != nil || c.String() != "B2" {
			t.Errorf("SmartShot() = %v, %v; want B2, the only cell left", c, err)
		}
	})

	t.Run("no cells left", func(t *testing.T) {
		b := NewBoard([]FleetEntry{{Type: Submarine, Count: 1}})
		place(t, b, placement{Submarine, "A1", Horizontal})
		smartShootAllBut(t, b)
		if _, err := SmartShot(b, seeded(1)); !errors.Is(err, ErrNoCellsLeft) {
			t.Errorf("SmartShot on a fully shot board error = %v, want ErrNoCellsLeft", err)
		}
	})
}

// CA-5: right after a ship sinks, its hits are no longer pending: the next
// shot is a hunt shot, away from the sunk ship.
func TestSmartShot_AfterSunk(t *testing.T) {
	b := exampleBoard(t)
	if shot := fire(t, b, "A5", "B5"); shot.Result != Sunk {
		t.Fatalf("A5 and B5 gave %v, want the Destroyer sunk", shot.Result)
	}
	smartHunts(t, smartChoices(t, b), smartCells(t, "A5", "B5"))
}

// smartTargetTest is a board of the example fleet after shots, and the
// cells SmartShot should choose among.
type smartTargetTest struct {
	name  string
	shots []string
	want  []string
}

// smartRunTargets checks, for each test, that SmartShot chooses exactly
// among the wanted cells across the seeds.
func smartRunTargets(t *testing.T, tests []smartTargetTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := exampleBoard(t)
			fire(t, b, tt.shots...)
			if got, want := smartChoices(t, b), smartCells(t, tt.want...); !maps.Equal(got, want) {
				t.Errorf("SmartShot after %v chose %v, want %v", tt.shots, got, want)
			}
		})
	}
}

// CA-2: with a single hit on a ship afloat, SmartShot shoots at one of the
// cells above, below, left or right of it that are not shot yet.
func TestSmartShot_TargetNeighbors(t *testing.T) {
	smartRunTargets(t, []smartTargetTest{
		{"middle of the board", []string{"E3"}, []string{"E2", "E4", "D3", "F3"}},
		{"some neighbors already shot", []string{"D3", "E2", "E3"}, []string{"E4", "F3"}},
		{"corner", []string{"J1"}, []string{"I1", "J2"}},
		{"a neighbor touches a sunk ship", []string{"E1", "F1", "G1", "H1", "E3"}, []string{"D3", "E4", "F3"}},
	})
}

// The strategies fit the ShotStrategy type.
var (
	_ ShotStrategy = RandomShot
	_ ShotStrategy = SmartShot
)
