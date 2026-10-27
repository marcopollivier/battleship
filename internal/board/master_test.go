package board

import (
	"maps"
	"slices"
	"testing"
)

// masterChoices runs MasterShot on b with smartSeeds seeds, without firing,
// and returns the set of cells it chose.
func masterChoices(t *testing.T, b *Board) map[Coord]bool {
	t.Helper()
	return strategyChoices(t, MasterShot, b)
}

// masterBest counts, on a board with no hits on a ship still afloat, the
// ways ships could cover each cell not shot yet, and returns the cells with
// the highest count.
func masterBest(t *testing.T, b *Board, ships []ShipType) map[Coord]bool {
	t.Helper()
	best, top := map[Coord]bool{}, 0
	for c, n := range masterCounts(t, masterForbidden(b), ships) {
		switch {
		case n > top:
			best, top = map[Coord]bool{c: true}, n
		case n == top:
			best[c] = true
		}
	}
	return best
}

// masterForbidden returns the cells where no ship can be on a board with no
// hits on a ship still afloat: the cells shot, and the cells that touch a
// hit, since every hit is part of a sunk ship.
func masterForbidden(b *Board) map[Coord]bool {
	forbidden := map[Coord]bool{}
	hits := map[Coord]bool{}
	for row := range Size {
		for col := range Size {
			c := Coord{Col: col, Row: row}
			if b.Cell(c) != CellUntouched {
				forbidden[c] = true
			}
			if b.Cell(c) == CellHit {
				hits[c] = true
			}
		}
	}
	maps.Copy(forbidden, smartAround(hits))
	return forbidden
}

// masterCounts counts, for each cell, the positions of ships on the board
// that cover it and avoid the forbidden cells.
func masterCounts(t *testing.T, forbidden map[Coord]bool, ships []ShipType) map[Coord]int {
	t.Helper()
	counts := map[Coord]int{}
	for _, ship := range ships {
		orientations := ValidOrientations(ship)
		if ship == Submarine {
			orientations = orientations[:1]
		}
		for ref := range Size * Size {
			for _, o := range orientations {
				cells, err := ShipCells(ship, Coord{Col: ref % Size, Row: ref / Size}, o)
				if err != nil {
					t.Fatalf("ShipCells: %v", err)
				}
				if slices.ContainsFunc(cells, func(c Coord) bool { return !c.inBounds() || forbidden[c] }) {
					continue
				}
				for _, c := range cells {
					counts[c]++
				}
			}
		}
	}
	return counts
}

// masterFleet lists one entry per ship of fleet, as the ships afloat at the
// start of a game.
func masterFleet(fleet []FleetEntry) []ShipType {
	return shipsBySize(fleet)
}

// masterWithout returns ships without one entry of each of gone.
func masterWithout(ships []ShipType, gone ...ShipType) []ShipType {
	ships = slices.Clone(ships)
	for _, g := range gone {
		if i := slices.Index(ships, g); i >= 0 {
			ships = slices.Delete(ships, i, i+1)
		}
	}
	return ships
}

// CA-1: on a board with no shots, MasterShot draws among the cells most
// ships could cover, which are in the middle of the board.
func TestMasterShot_Hunt(t *testing.T) {
	b := exampleBoard(t)
	want := masterBest(t, b, masterFleet(DefaultFleet()))
	for c := range want {
		if c.Col < 2 || c.Col > 7 || c.Row < 2 || c.Row > 7 {
			t.Errorf("the highest count is at %v, want the middle of the board", c)
		}
	}
	if got := masterChoices(t, b); !maps.Equal(got, want) {
		t.Errorf("MasterShot chose %v, want %v", got, want)
	}
}

// CA-2: MasterShot never shoots at a cell where no ship still afloat fits,
// and the ships that sank no longer count.
func TestMasterShot_NoRoom(t *testing.T) {
	t.Run("a corner closed by water", func(t *testing.T) {
		b := NewBoard([]FleetEntry{{Type: Submarine, Count: 1}, {Type: Destroyer, Count: 1}})
		place(t, b, placement{Destroyer, "C1", Horizontal}, placement{Submarine, "J10", Horizontal})
		smartShootAllBut(t, b, "A1", "C1", "D1")
		if got, want := masterChoices(t, b), smartCells(t, "C1", "D1"); !maps.Equal(got, want) {
			t.Errorf("MasterShot chose %v, want %v: only the Destroyer is afloat", got, want)
		}
	})

	t.Run("every Submarine sunk", func(t *testing.T) {
		b := exampleBoard(t)
		fire(t, b, "H5", "A8", "E8", "I9", "I10", "J9")
		afloat := masterWithout(masterFleet(DefaultFleet()), Submarine, Submarine, Submarine)
		want := masterBest(t, b, afloat)
		if !masterBest(t, b, []ShipType{Submarine})[mustCoord(t, "J10")] {
			t.Fatal("a Submarine does not fit at J10; pick other shots")
		}
		got := masterChoices(t, b)
		if got[mustCoord(t, "J10")] {
			t.Error("MasterShot chose J10, where only a Submarine fits")
		}
		if !maps.Equal(got, want) {
			t.Errorf("MasterShot chose %v, want %v", got, want)
		}
	})
}

// CA-4: a sunk ship no longer counts among the ships afloat, and MasterShot
// never shoots next to it.
func TestMasterShot_AfterSunk(t *testing.T) {
	t.Run("next to a sunk ship", func(t *testing.T) {
		b := exampleBoard(t)
		if shot := fire(t, b, "A5", "B5"); shot.Result != Sunk {
			t.Fatalf("A5 and B5 gave %v, want the Destroyer sunk", shot.Result)
		}
		fire(t, b, "E5", "F6")
		got := masterChoices(t, b)
		for c := range got {
			if smartAround(smartCells(t, "A5", "B5"))[c] {
				t.Errorf("MasterShot chose %v, next to the sunk Destroyer", c)
			}
		}
		if want := masterBest(t, b, masterWithout(masterFleet(DefaultFleet()), Destroyer)); !maps.Equal(got, want) {
			t.Errorf("MasterShot chose %v, want %v", got, want)
		}
	})

	// Only a line of 5 cells and a cross of 5 cells are left. The Cruiser
	// fits best in the middle of the line, and a Destroyer in the middle
	// of the cross.
	t.Run("the sunk ship no longer counts", func(t *testing.T) {
		fleet := []FleetEntry{{Type: Destroyer, Count: 1}, {Type: Cruiser, Count: 1}}
		b := NewBoard(fleet)
		place(t, b, placement{Destroyer, "A1", Horizontal}, placement{Cruiser, "D10", Horizontal})
		smartShootAllBut(t, b, "C10", "D10", "E10", "F10", "G10", "H3", "G4", "H4", "I4", "H5")
		want := smartCells(t, "E10")
		if maps.Equal(want, masterBest(t, b, masterFleet(fleet))) {
			t.Fatal("the sunk Destroyer does not change the highest count")
		}
		if got := masterChoices(t, b); !maps.Equal(got, want) {
			t.Errorf("MasterShot chose %v, want %v", got, want)
		}
	})
}

var _ ShotStrategy = MasterShot
