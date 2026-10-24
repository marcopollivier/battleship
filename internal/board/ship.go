package board

import "fmt"

// ShipType identifies a kind of ship. Names are in English; the terminal
// layer translates them to Portuguese.
type ShipType int

const (
	Submarine       ShipType = iota // 1 cell
	Destroyer                       // line of 2
	Cruiser                         // line of 3
	Battleship                      // line of 4
	AircraftCarrier                 // T of 5
)

var shipTypeInfo = map[ShipType]struct {
	name string
	size int
}{
	Submarine:       {"Submarine", 1},
	Destroyer:       {"Destroyer", 2},
	Cruiser:         {"Cruiser", 3},
	Battleship:      {"Battleship", 4},
	AircraftCarrier: {"AircraftCarrier", 5},
}

// Size returns how many cells a ship of this type occupies, or 0 for an unknown type.
func (t ShipType) Size() int {
	return shipTypeInfo[t].size
}

func (t ShipType) String() string {
	if info, ok := shipTypeInfo[t]; ok {
		return info.name
	}
	return fmt.Sprintf("ShipType(%d)", int(t))
}

// FleetEntry is how many ships of one type a fleet has.
type FleetEntry struct {
	Type  ShipType
	Count int
}

// DefaultFleet returns the fleet from the original paper pad: 8 ships, 19 cells.
func DefaultFleet() []FleetEntry {
	return []FleetEntry{
		{Type: Submarine, Count: 3},
		{Type: Destroyer, Count: 2},
		{Type: Cruiser, Count: 1},
		{Type: Battleship, Count: 1},
		{Type: AircraftCarrier, Count: 1},
	}
}
