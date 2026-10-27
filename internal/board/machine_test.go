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

// CA-3: with two or more hits in line on a ship afloat, SmartShot shoots at
// the ends of the line that are on the board and not shot yet.
func TestSmartShot_TargetLine(t *testing.T) {
	smartRunTargets(t, []smartTargetTest{
		{"horizontal", []string{"F1", "G1"}, []string{"E1", "H1"}},
		{"vertical", []string{"J2", "J3"}, []string{"J1", "J4"}},
		{"one end already shot", []string{"D1", "E1", "F1"}, []string{"G1"}},
		{"one end off the board", []string{"J1", "J2"}, []string{"J3"}},
	})
}

// CA-4: when both ends of the line are shot or off the board, or the hits
// are not in line, SmartShot shoots next to any of the hits.
func TestSmartShot_TargetLineBlocked(t *testing.T) {
	smartRunTargets(t, []smartTargetTest{
		{"Porta-Aviões bar with both ends blocked", []string{"D1", "A1", "B1", "C1"}, []string{"A2", "B2", "C2"}},
		{"hits not in line", []string{"A1", "B1", "B2"}, []string{"C1", "A2", "C2", "B3"}},
	})
}

// smartGameShots plays a whole game of strategy against the board drawn for
// seed, with the shots also drawn from seed, and returns how many shots it
// took to sink the fleet. A repeated cell fails the test, since Fire
// rejects it.
func smartGameShots(t *testing.T, strategy ShotStrategy, seed uint64) int {
	t.Helper()
	b := mustRandomBoard(t, seed)
	rng := seeded(seed)
	for shots := 1; shots <= Size*Size; shots++ {
		c, err := strategy(b, rng)
		if err != nil {
			t.Fatalf("seed %d, shot %d: %v", seed, shots, err)
		}
		if _, err := b.Fire(c); err != nil {
			t.Fatalf("seed %d, shot %d: Fire(%v): %v", seed, shots, c, err)
		}
		if b.AllSunk() {
			return shots
		}
	}
	t.Fatalf("seed %d: the fleet is still afloat after %d shots", seed, Size*Size)
	return 0
}

// CA-6: on any board of the default fleet, SmartShot never shoots the same
// cell twice and sinks the whole fleet, with fewer shots on average than
// RandomShot.
func TestSmartShot_Game(t *testing.T) {
	const games = 200
	smartTotal, randomTotal := 0, 0
	for seed := range uint64(games) {
		smartTotal += smartGameShots(t, SmartShot, seed)
		randomTotal += smartGameShots(t, RandomShot, seed)
	}
	smartAvg, randomAvg := float64(smartTotal)/games, float64(randomTotal)/games
	t.Logf("average shots to sink the fleet in %d games: SmartShot %.1f, RandomShot %.1f", games, smartAvg, randomAvg)
	if smartAvg >= randomAvg {
		t.Errorf("SmartShot took %.1f shots on average, want fewer than RandomShot's %.1f", smartAvg, randomAvg)
	}
}

// smartOtherBoard is exampleBoard with two Submarines elsewhere.
func smartOtherBoard(t *testing.T) *Board {
	t.Helper()
	b := NewBoard(DefaultFleet())
	place(t, b,
		placement{AircraftCarrier, "B1", Down},
		placement{Battleship, "E1", Horizontal},
		placement{Cruiser, "J1", Vertical},
		placement{Destroyer, "A5", Horizontal},
		placement{Destroyer, "E3", Horizontal},
		placement{Submarine, "H5", Horizontal},
		placement{Submarine, "C10", Horizontal},
		placement{Submarine, "G8", Horizontal},
	)
	return b
}

// SmartShot only uses what the shooter knows: two boards with the same
// shots and results, but ships not found yet in different places, give the
// same shot with the same seed.
func TestSmartShot_PublicKnowledge(t *testing.T) {
	tests := []struct {
		name  string
		shots []string
	}{
		{"hunt after sinking", []string{"J10", "H5", "A5", "B5"}},
		{"one hit", []string{"J10", "E3"}},
		{"hits in line", []string{"E1", "F1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b1, b2 := exampleBoard(t), smartOtherBoard(t)
			fire(t, b1, tt.shots...)
			fire(t, b2, tt.shots...)
			for seed := range uint64(smartSeeds) {
				c1, err1 := SmartShot(b1, seeded(seed))
				c2, err2 := SmartShot(b2, seeded(seed))
				if c1 != c2 || err1 != nil || err2 != nil {
					t.Fatalf("seed %d: SmartShot = %v, %v on one board and %v, %v on the other", seed, c1, err1, c2, err2)
				}
			}
		})
	}
}

// The strategies fit the ShotStrategy type.
var (
	_ ShotStrategy = RandomShot
	_ ShotStrategy = SmartShot
)
