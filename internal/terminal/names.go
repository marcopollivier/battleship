// Package terminal draws the game on a text terminal. It reads the board
// domain and never changes it; text shown to the player is in Portuguese.
package terminal

import "github.com/marcopollivier/battleship/internal/board"

var shipNames = map[board.ShipType]string{
	board.Submarine:       "Submarino",
	board.Destroyer:       "Destroyer",
	board.Cruiser:         "Cruzador",
	board.Battleship:      "Encouraçado",
	board.AircraftCarrier: "Porta-Aviões",
}

var orientationNames = map[board.Orientation]string{
	board.Horizontal: "horizontal",
	board.Vertical:   "vertical",
	board.Up:         "cima",
	board.Down:       "baixo",
	board.Left:       "esquerda",
	board.Right:      "direita",
}

// ShipName returns the Portuguese name of a ship type, e.g. "Cruzador".
// An unknown type falls back to the domain's String().
func ShipName(t board.ShipType) string {
	if name, ok := shipNames[t]; ok {
		return name
	}
	return t.String()
}

// OrientationName returns the Portuguese name of an orientation, e.g. "cima".
// An unknown orientation falls back to the domain's String().
func OrientationName(o board.Orientation) string {
	if name, ok := orientationNames[o]; ok {
		return name
	}
	return o.String()
}
