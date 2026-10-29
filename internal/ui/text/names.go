// Package text holds the Portuguese text shared by every interface: the
// names and letters of ships and orientations, the messages of shots,
// placements and fleets, and the score line. It reads the core and never
// changes it.
package text

import (
	"strings"

	"github.com/marcopollivier/battleship/internal/core/board"
)

var shipNames = map[board.ShipType]string{
	board.Submarine:       "Submarino",
	board.Destroyer:       "Destroyer",
	board.Cruiser:         "Cruzador",
	board.Battleship:      "Encouraçado",
	board.AircraftCarrier: "Porta-Aviões",
}

// shipSymbols are the ASCII letters of each ship type, drawn on its cells
// and used in the fleet text.
var shipSymbols = map[board.ShipType]string{
	board.Submarine:       "S",
	board.Destroyer:       "D",
	board.Cruiser:         "C",
	board.Battleship:      "E",
	board.AircraftCarrier: "P",
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

// ShipSymbol returns the letter of a ship type, or "?" for an unknown type.
func ShipSymbol(t board.ShipType) string {
	if symbol, ok := shipSymbols[t]; ok {
		return symbol
	}
	return "?"
}

// OrientationByName finds an orientation by its Portuguese name or first
// letter, in any case.
func OrientationByName(name string) (board.Orientation, bool) {
	name = strings.ToLower(name)
	for o, full := range orientationNames {
		if name == full || name == full[:1] {
			return o, true
		}
	}
	return 0, false
}

// OrientationList names the orientations that apply to t, in Portuguese.
func OrientationList(t board.ShipType) []string {
	var names []string
	for _, o := range board.ValidOrientations(t) {
		names = append(names, OrientationName(o))
	}
	return names
}
