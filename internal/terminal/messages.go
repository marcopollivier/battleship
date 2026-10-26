package terminal

import "github.com/marcopollivier/battleship/internal/board"

// ShotMessage returns the Portuguese message for a shot result, e.g. "Afundou o Cruzador!".
// The ship is named only when it sinks, since Shot.Ship is unset otherwise.
// An unknown result falls back to the domain's String().
func ShotMessage(s board.Shot) string {
	switch s.Result {
	case board.Miss:
		return "Água."
	case board.Hit:
		return "Acerto!"
	case board.Sunk:
		return "Afundou o " + ShipName(s.Ship) + "!"
	default:
		return s.Result.String()
	}
}
