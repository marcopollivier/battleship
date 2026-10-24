package board

import (
	"errors"
	"fmt"
)

// ErrInvalidOrientation means the orientation does not apply to the ship type.
var ErrInvalidOrientation = errors.New("orientation does not apply to ship type")

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

// Orientation is Horizontal or Vertical for line ships, and the direction
// the stem points (Up, Down, Left, Right) for the T-shaped AircraftCarrier.
type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
	Up
	Down
	Left
	Right
)

var orientationNames = map[Orientation]string{
	Horizontal: "Horizontal",
	Vertical:   "Vertical",
	Up:         "Up",
	Down:       "Down",
	Left:       "Left",
	Right:      "Right",
}

func (o Orientation) String() string {
	if name, ok := orientationNames[o]; ok {
		return name
	}
	return fmt.Sprintf("Orientation(%d)", int(o))
}

// carrierOffsets are the AircraftCarrier cells as (dCol, dRow) from the
// middle of the bar, keyed by the direction the stem points.
var carrierOffsets = map[Orientation][5][2]int{
	Down:  {{-1, 0}, {0, 0}, {1, 0}, {0, 1}, {0, 2}},
	Up:    {{-1, 0}, {0, 0}, {1, 0}, {0, -1}, {0, -2}},
	Right: {{0, -1}, {0, 0}, {0, 1}, {1, 0}, {2, 0}},
	Left:  {{0, -1}, {0, 0}, {0, 1}, {-1, 0}, {-2, 0}},
}

// ShipCells returns the cells a ship of type t occupies when placed at ref
// with orientation o. It does not check board bounds.
// Returns ErrInvalidOrientation if o does not apply to t.
func ShipCells(t ShipType, ref Coord, o Orientation) ([]Coord, error) {
	size := t.Size()
	if size == 0 {
		return nil, fmt.Errorf("unknown ship type %v", t)
	}
	if t == AircraftCarrier {
		return carrierCells(ref, o)
	}

	var dCol, dRow int
	switch o {
	case Horizontal:
		dCol = 1
	case Vertical:
		dRow = 1
	default:
		return nil, fmt.Errorf("%w: %v %v", ErrInvalidOrientation, t, o)
	}

	cells := make([]Coord, size)
	for i := range size {
		cells[i] = Coord{Col: ref.Col + i*dCol, Row: ref.Row + i*dRow}
	}
	return cells, nil
}

func carrierCells(ref Coord, o Orientation) ([]Coord, error) {
	offsets, ok := carrierOffsets[o]
	if !ok {
		return nil, fmt.Errorf("%w: %v %v", ErrInvalidOrientation, AircraftCarrier, o)
	}
	cells := make([]Coord, len(offsets))
	for i, d := range offsets {
		cells[i] = Coord{Col: ref.Col + d[0], Row: ref.Row + d[1]}
	}
	return cells, nil
}
