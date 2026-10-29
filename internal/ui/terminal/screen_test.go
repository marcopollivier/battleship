package terminal

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/marcopollivier/battleship/internal/core/board"
)

// maxWidth is the widest a line of the screen may be, in columns.
const maxWidth = 80

// visibleWidth is the number of columns line takes on the terminal.
func visibleWidth(line string) int {
	return utf8.RuneCountInString(stripColors(line))
}

// CA-1 (spec 022): the two grids are joined line by line, the right one
// always starting at the same column, with or without colors, and no line
// is wider than 80 columns.
func TestSideBySide(t *testing.T) {
	player := attackExampleBoard(t)
	attackFire(t, player, "J10", "E1")
	enemy := attackExampleBoard(t)
	attackFire(t, enemy, "A1", "B1", "J10")

	for _, style := range []Style{Plain, Colored} {
		left, right := ownGrid(player, style), attackGrid(enemy, style)
		lines := sideBySide(left, right)
		if len(lines) != board.Size+1 {
			t.Fatalf("sideBySide() has %d lines, want %d", len(lines), board.Size+1)
		}
		for i, line := range lines {
			if want := left[i] + strings.Repeat(" ", boardGap) + right[i]; line != want {
				t.Errorf("line %d = %q, want %q", i+1, line, want)
			}
			if n := visibleWidth(line); n != 2*gridWidth+boardGap || n > maxWidth {
				t.Errorf("line %d is %d columns wide, want %d: %q", i+1, n, 2*gridWidth+boardGap, line)
			}
		}
	}
}

// CA-1 (spec 022): the titles are over their boards, without colons.
func TestBoardTitles(t *testing.T) {
	want := "Sua frota" + strings.Repeat(" ", 18) + "Frota inimiga"
	if boardTitles != want {
		t.Errorf("boardTitles = %q, want %q", boardTitles, want)
	}
	if got := strings.Index(boardTitles, "Frota inimiga"); got != gridWidth+boardGap {
		t.Errorf("the enemy title starts at column %d, want %d", got, gridWidth+boardGap)
	}
}

// CA-5 (spec 022): the clear code moves the cursor home and erases the
// terminal.
func TestClearScreen(t *testing.T) {
	if clearScreen != "\x1b[H\x1b[2J" {
		t.Errorf("clearScreen = %q", clearScreen)
	}
}
