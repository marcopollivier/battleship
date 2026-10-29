package terminal

import "github.com/marcopollivier/battleship/internal/core/board"

// Style says whether the boards and the messages are drawn with ANSI colors.
type Style int

const (
	Plain   Style = iota // no color codes, the output before spec 010
	Colored              // each board symbol and message in its color
)

// ANSI codes of the board symbols and the messages.
const (
	waterColor     = "\x1b[34m"   // blue
	shipColor      = "\x1b[1;37m" // bold white
	hitColor       = "\x1b[1;31m" // bold red
	missColor      = "\x1b[36m"   // cyan
	untouchedColor = "\x1b[90m"   // gray
	sunkColor      = "\x1b[31m"   // red
	winColor       = "\x1b[1;32m" // bold green
	warnColor      = "\x1b[33m"   // yellow
	resetColor     = "\x1b[0m"
)

// paint returns text, a board symbol or a whole message, in color, followed
// by the reset code, or text unchanged when s is Plain.
func (s Style) paint(text, color string) string {
	if s != Colored {
		return text
	}
	return color + text + resetColor
}

// shotColor is the color of the message for a shot result: cyan for a miss,
// as on the boards, and bold red for a hit or a sunk ship.
func shotColor(r board.ShotResult) string {
	if r == board.Miss {
		return missColor
	}
	return hitColor
}
