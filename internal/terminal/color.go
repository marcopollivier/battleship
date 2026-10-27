package terminal

// Style says whether the boards are drawn with ANSI colors.
type Style int

const (
	Plain   Style = iota // no color codes, the output before spec 010
	Colored              // each board symbol in its color
)
