package terminal

import "strings"

// Display is how a game is shown: the colors of the boards and messages,
// and whether the terminal is cleared before each screen.
type Display struct {
	Style Style
	Clear bool
}

// clearScreen moves the cursor home and erases the terminal.
const clearScreen = "\x1b[H\x1b[2J"

// boardGap is the number of spaces between two boards side by side.
const boardGap = 5

// Titles of the boards in a game against the machine.
const (
	ownTitle   = "Sua frota"
	enemyTitle = "Frota inimiga"
)

// boardTitles is the line over two boards side by side, each title starting
// at the column of its board.
var boardTitles = ownTitle + strings.Repeat(" ", gridWidth+boardGap-len(ownTitle)) + enemyTitle

// sideBySide joins two board grids line by line, the left one padded to
// its visible width plus a gap, so color codes do not shift the columns.
// Every line of a grid is gridWidth columns wide, so the gap is enough.
func sideBySide(left, right []string) []string {
	lines := make([]string, max(len(left), len(right)))
	gap := strings.Repeat(" ", boardGap)
	for i := range lines {
		l, r := strings.Repeat(" ", gridWidth), ""
		if i < len(left) {
			l = left[i]
		}
		if i < len(right) {
			r = right[i]
		}
		lines[i] = l + gap + r
	}
	return lines
}
