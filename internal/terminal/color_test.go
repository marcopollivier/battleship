package terminal

import (
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
)

// colorCode matches an ANSI color code, such as "\x1b[1;31m" or "\x1b[0m".
var colorCode = regexp.MustCompile("\x1b\\[[0-9;]*m")

// stripColors removes the ANSI color codes from s.
func stripColors(s string) string {
	return colorCode.ReplaceAllString(s, "")
}

// colored wraps symbol in the ANSI code and the reset, as the boards draw it.
func colored(code, symbol string) string {
	return "\x1b[" + code + "m" + symbol + "\x1b[0m"
}

// renderColored draws b with draw, RenderBoard or RenderAttackBoard, with colors.
func renderColored(t *testing.T, draw func(io.Writer, *board.Board, Style) error, b *board.Board) string {
	t.Helper()
	var sb strings.Builder
	if err := draw(&sb, b, Colored); err != nil {
		t.Fatalf("draw: %v", err)
	}
	return sb.String()
}

// coloredCellAt returns what is drawn for a cell such as "E5", color codes
// included.
func coloredCellAt(t *testing.T, output, cell string) string {
	t.Helper()
	c, err := board.ParseCoord(cell)
	if err != nil {
		t.Fatalf("ParseCoord(%q): %v", cell, err)
	}
	return strings.Fields(strings.Split(output, "\n")[c.Row+1])[c.Col+1]
}

// assertColoredCells checks that every cell of output is the symbol drawn
// in plain wrapped in the code that codes gives for it.
func assertColoredCells(t *testing.T, output, plain string, codes map[string]string) {
	t.Helper()
	for row := range board.Size {
		for col := range board.Size {
			cell := board.Coord{Col: col, Row: row}.String()
			symbol := cellAt(t, plain, cell)
			code, ok := codes[symbol]
			if !ok {
				t.Fatalf("cell %s has symbol %q with no color", cell, symbol)
			}
			want := colored(code, symbol)
			if got := coloredCellAt(t, output, cell); got != want {
				t.Errorf("cell %s = %q, want %q", cell, got, want)
			}
		}
	}
}

// ownBoardCodes are the colors of the own board's symbols.
var ownBoardCodes = map[string]string{
	"~": "34",
	"S": "1;37", "D": "1;37", "C": "1;37", "E": "1;37", "P": "1;37",
	"X": "1;31",
	"o": "36",
}

// CA-1: on the own board, water is blue, whole ships bold white, hits bold
// red and misses cyan.
func TestRenderBoard_Colored(t *testing.T) {
	b := attackExampleBoard(t)
	attackFire(t, b, "J10", "E1")
	got := renderColored(t, RenderBoard, b)

	want := map[string]string{
		"J10": colored("36", "o"),
		"E1":  colored("1;31", "X"),
		"F1":  colored("1;37", "E"),
		"A1":  colored("1;37", "P"),
		"I10": colored("34", "~"),
	}
	for cell, symbol := range want {
		if s := coloredCellAt(t, got, cell); s != symbol {
			t.Errorf("cell %s = %q, want %q", cell, s, symbol)
		}
	}
	assertColoredCells(t, got, render(t, b), ownBoardCodes)
}

// CA-3: without the color codes the colored board is the plain one, and the
// header, the row numbers and the legend have no color codes.
func TestRender_ColoredMatchesPlain(t *testing.T) {
	b := attackExampleBoard(t)
	attackFire(t, b, "J10", "E1", "A8", "C5")
	tests := []struct {
		name           string
		colored, plain string
	}{
		{"own board", renderColored(t, RenderBoard, b), render(t, b)},
		{"attack board", renderColored(t, RenderAttackBoard, b), attackRender(t, b)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripColors(tt.colored); got != tt.plain {
				t.Errorf("colored board without codes =\n%s\nwant\n%s", got, tt.plain)
			}
			lines := strings.Split(tt.colored, "\n")
			uncolored := append([]string{lines[0]}, lines[board.Size+1:]...)
			for row := 1; row <= board.Size; row++ {
				uncolored = append(uncolored, lines[row][:3])
			}
			for _, text := range uncolored {
				if strings.Contains(text, "\x1b") {
					t.Errorf("%q has a color code", text)
				}
			}
		})
	}
}

// attackBoardCodes are the colors of the attack board's symbols.
var attackBoardCodes = map[string]string{
	".": "90",
	"o": "36",
	"X": "1;31",
	"S": "31", "D": "31", "C": "31", "E": "31", "P": "31",
}

// CA-2: on the attack board, untouched cells are gray, misses cyan, hits
// bold red and sunk ships red.
func TestRenderAttackBoard_Colored(t *testing.T) {
	b := attackExampleBoard(t)
	attackFire(t, b, "J10", "E1", "A8")
	got := renderColored(t, RenderAttackBoard, b)

	want := map[string]string{
		"J10": colored("36", "o"),
		"E1":  colored("1;31", "X"),
		"A8":  colored("31", "S"),
		"F1":  colored("90", "."),
	}
	for cell, symbol := range want {
		if s := coloredCellAt(t, got, cell); s != symbol {
			t.Errorf("cell %s = %q, want %q", cell, s, symbol)
		}
	}
	assertColoredCells(t, got, attackRender(t, b), attackBoardCodes)
}

// boardRow matches a board row, which starts with its number.
var boardRow = regexp.MustCompile(`^( \d|10) `)

// assertBoardsColored checks that output, a game drawn with Colored, is
// plain without its color codes, and that every board row in it is colored.
func assertBoardsColored(t *testing.T, output, plain string) {
	t.Helper()
	if got := stripColors(output); got != plain {
		t.Errorf("colored game without codes =\n%s\nwant\n%s", got, plain)
	}
	rows := 0
	for _, line := range strings.Split(output, "\n") {
		if !boardRow.MatchString(line) {
			continue
		}
		rows++
		if !strings.Contains(line, "\x1b") {
			t.Errorf("board row %q has no color", line)
		}
	}
	if rows == 0 {
		t.Error("the game shows no board")
	}
}

// colorSpan matches a text in color: a color code other than the reset, the
// text and the reset.
var colorSpan = regexp.MustCompile("\x1b\\[[0-9;]*[1-9][0-9;]*m[^\x1b]*\x1b\\[0m")

// assertUncolored checks that every text appears in output, a game drawn
// with Colored, and never inside a color code: prompts, instructions and
// titles have no color.
func assertUncolored(t *testing.T, output string, texts ...string) {
	t.Helper()
	outside := colorSpan.ReplaceAllString(output, "")
	for _, text := range texts {
		if !strings.Contains(output, text) {
			t.Errorf("the game does not show %q", text)
			continue
		}
		if got, want := strings.Count(outside, text), strings.Count(stripColors(output), text); got != want {
			t.Errorf("%q is colored %d of %d times", text, want-got, want)
		}
	}
}

// Spec 015: shotColor gives each shot result the color of its board symbol.
func TestShotColor(t *testing.T) {
	tests := []struct {
		result board.ShotResult
		want   string
	}{
		{board.Miss, "\x1b[36m"},
		{board.Hit, "\x1b[1;31m"},
		{board.Sunk, "\x1b[1;31m"},
	}
	for _, tt := range tests {
		if got := shotColor(tt.result); got != tt.want {
			t.Errorf("shotColor(%v) = %q, want %q", tt.result, got, tt.want)
		}
	}
}
