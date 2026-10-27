package terminal

import (
	"fmt"
	"strings"
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
)

func render(t *testing.T, b *board.Board) string {
	t.Helper()
	var sb strings.Builder
	if err := RenderBoard(&sb, b, Plain); err != nil {
		t.Fatalf("RenderBoard: %v", err)
	}
	return sb.String()
}

// CA-7 (spec 009, was CA-1 of spec 004): an empty board is drawn like the old program, without trailing
// spaces, followed by a blank line and the legend.
func TestRenderBoard_Empty(t *testing.T) {
	want := `   A B C D E F G H I J
 1 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 2 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 3 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 4 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 5 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 6 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 7 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 8 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 9 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
10 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~

Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões, X atingido, o tiro na água
`
	got := render(t, board.NewBoard(board.DefaultFleet()))
	if got != want {
		t.Errorf("RenderBoard() =\n%s\nwant\n%s", got, want)
	}
	for i, line := range strings.Split(got, "\n") {
		if strings.HasSuffix(line, " ") {
			t.Errorf("line %d has a trailing space: %q", i+1, line)
		}
	}
}

type placement struct {
	shipType    board.ShipType
	ref         string
	orientation board.Orientation
}

func boardWith(t *testing.T, placements ...placement) *board.Board {
	t.Helper()
	b := board.NewBoard(board.DefaultFleet())
	for _, p := range placements {
		ref, err := board.ParseCoord(p.ref)
		if err != nil {
			t.Fatalf("ParseCoord(%q): %v", p.ref, err)
		}
		if err := b.Place(p.shipType, ref, p.orientation); err != nil {
			t.Fatalf("Place(%v, %s, %v): %v", p.shipType, p.ref, p.orientation, err)
		}
	}
	return b
}

// cellAt returns the symbol drawn for a cell such as "E5".
func cellAt(t *testing.T, output, cell string) string {
	t.Helper()
	c, err := board.ParseCoord(cell)
	if err != nil {
		t.Fatalf("ParseCoord(%q): %v", cell, err)
	}
	line := strings.Split(output, "\n")[c.Row+1]
	pos := 3 + 2*c.Col
	return line[pos : pos+1]
}

func TestRenderBoard_Ships(t *testing.T) {
	t.Run("CA-2 cruiser at A1 horizontal", func(t *testing.T) {
		got := render(t, boardWith(t, placement{board.Cruiser, "A1", board.Horizontal}))
		lines := strings.Split(got, "\n")
		if want := " 1 C C C ~ ~ ~ ~ ~ ~ ~"; lines[1] != want {
			t.Errorf("row 1 = %q, want %q", lines[1], want)
		}
		for i := 2; i <= board.Size; i++ {
			if want := fmt.Sprintf("%2d", i) + strings.Repeat(" ~", board.Size); lines[i] != want {
				t.Errorf("row %d = %q, want %q", i, lines[i], want)
			}
		}
	})

	t.Run("CA-3 aircraft carrier at E5 down", func(t *testing.T) {
		got := render(t, boardWith(t, placement{board.AircraftCarrier, "E5", board.Down}))
		for _, cell := range []string{"D5", "E5", "F5", "E6", "E7"} {
			if s := cellAt(t, got, cell); s != "P" {
				t.Errorf("cell %s = %q, want %q", cell, s, "P")
			}
		}
		if n := strings.Count(got, "P") - strings.Count(Legend, "P"); n != 5 {
			t.Errorf("board has %d P cells, want 5", n)
		}
	})
}

// CA-7 (spec 009, was CA-6 of spec 004): a complete board shows every ship with its letter, and everything
// above the legend is ASCII.
func TestRenderBoard_FullFleet(t *testing.T) {
	b := boardWith(t,
		placement{board.AircraftCarrier, "B1", board.Down}, // A1 B1 C1 B2 B3
		placement{board.Battleship, "E1", board.Horizontal},
		placement{board.Cruiser, "J1", board.Vertical},
		placement{board.Destroyer, "A5", board.Horizontal},
		placement{board.Destroyer, "E3", board.Horizontal},
		placement{board.Submarine, "H5", board.Horizontal},
		placement{board.Submarine, "A8", board.Horizontal},
		placement{board.Submarine, "E8", board.Horizontal},
	)
	if !b.Complete() {
		t.Fatal("board is not complete")
	}

	want := `   A B C D E F G H I J
 1 P P P ~ E E E E ~ C
 2 ~ P ~ ~ ~ ~ ~ ~ ~ C
 3 ~ P ~ ~ D D ~ ~ ~ C
 4 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 5 D D ~ ~ ~ ~ ~ S ~ ~
 6 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 7 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 8 S ~ ~ ~ S ~ ~ ~ ~ ~
 9 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
10 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~

Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões, X atingido, o tiro na água
`
	got := render(t, b)
	if got != want {
		t.Errorf("RenderBoard() =\n%s\nwant\n%s", got, want)
	}

	boardPart, _, found := strings.Cut(got, Legend)
	if !found {
		t.Fatal("legend not found in output")
	}
	for i := 0; i < len(boardPart); i++ {
		if boardPart[i] >= 128 {
			t.Errorf("byte %d above the legend is not ASCII: %#x", i, boardPart[i])
		}
	}
}

func TestShipSymbol_Unknown(t *testing.T) {
	if got := shipSymbol(board.ShipType(99)); got != "?" {
		t.Errorf("shipSymbol(ShipType(99)) = %q, want %q", got, "?")
	}
}

// CA-6 (spec 009): the own board marks received shots, with X on a hit ship
// cell and o on a missed water cell, and keeps the other ship letters.
func TestRenderBoard_ReceivedShots(t *testing.T) {
	b := attackExampleBoard(t)
	attackFire(t, b, "J10", "E1")
	got := render(t, b)

	want := map[string]string{"J10": "o", "E1": "X", "F1": "E", "G1": "E", "H1": "E", "I10": "~", "A1": "P"}
	for cell, symbol := range want {
		if s := cellAt(t, got, cell); s != symbol {
			t.Errorf("cell %s = %q, want %q", cell, s, symbol)
		}
	}
	if !strings.HasSuffix(got, "\n\n"+Legend+"\n") {
		t.Errorf("output does not end with a blank line and the legend:\n%s", got)
	}
	if want := "Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões, X atingido, o tiro na água"; Legend != want {
		t.Errorf("Legend = %q, want %q", Legend, want)
	}
}
