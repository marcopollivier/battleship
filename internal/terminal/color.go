package terminal

// Style says whether the boards are drawn with ANSI colors.
type Style int

const (
	Plain   Style = iota // no color codes, the output before spec 010
	Colored              // each board symbol in its color
)

// ANSI codes of the board symbols.
const (
	waterColor = "\x1b[34m"   // blue
	shipColor  = "\x1b[1;37m" // bold white
	hitColor   = "\x1b[1;31m" // bold red
	missColor  = "\x1b[36m"   // cyan
	resetColor = "\x1b[0m"
)

// paint returns symbol in color, followed by the reset code, or symbol
// unchanged when s is Plain.
func (s Style) paint(symbol, color string) string {
	if s != Colored {
		return symbol
	}
	return color + symbol + resetColor
}
