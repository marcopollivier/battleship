package board

import (
	"errors"
	"slices"
	"testing"
)

// CA-1: the default fleet is the original paper pad fleet.
func TestDefaultFleet(t *testing.T) {
	want := map[ShipType]int{
		Submarine:       3,
		Destroyer:       2,
		Cruiser:         1,
		Battleship:      1,
		AircraftCarrier: 1,
	}

	fleet := DefaultFleet()
	got := map[ShipType]int{}
	totalCells := 0
	for _, entry := range fleet {
		got[entry.Type] += entry.Count
		totalCells += entry.Type.Size() * entry.Count
	}

	if len(got) != len(want) {
		t.Errorf("DefaultFleet() has %d ship types, want %d", len(got), len(want))
	}
	for shipType, count := range want {
		if got[shipType] != count {
			t.Errorf("DefaultFleet() has %d %v, want %d", got[shipType], shipType, count)
		}
	}
	if totalCells != 19 {
		t.Errorf("DefaultFleet() covers %d cells, want 19", totalCells)
	}
}

func TestShipType_Size(t *testing.T) {
	tests := []struct {
		shipType ShipType
		want     int
	}{
		{Submarine, 1},
		{Destroyer, 2},
		{Cruiser, 3},
		{Battleship, 4},
		{AircraftCarrier, 5},
	}
	for _, tt := range tests {
		if got := tt.shipType.Size(); got != tt.want {
			t.Errorf("%v.Size() = %d, want %d", tt.shipType, got, tt.want)
		}
	}
}

func TestShipType_String(t *testing.T) {
	tests := []struct {
		shipType ShipType
		want     string
	}{
		{Submarine, "Submarine"},
		{Destroyer, "Destroyer"},
		{Cruiser, "Cruiser"},
		{Battleship, "Battleship"},
		{AircraftCarrier, "AircraftCarrier"},
		{ShipType(99), "ShipType(99)"},
	}
	for _, tt := range tests {
		if got := tt.shipType.String(); got != tt.want {
			t.Errorf("ShipType(%d).String() = %q, want %q", int(tt.shipType), got, tt.want)
		}
	}
}

// mustCoord parses a coordinate written as in the spec, such as "A1".
func mustCoord(t *testing.T, s string) Coord {
	t.Helper()
	c, err := ParseCoord(s)
	if err != nil {
		t.Fatalf("ParseCoord(%q): %v", s, err)
	}
	return c
}

// cellNames renders cells in canonical form, sorted, so shapes compare as sets.
func cellNames(cells []Coord) []string {
	names := make([]string, len(cells))
	for i, c := range cells {
		names[i] = c.String()
	}
	slices.Sort(names)
	return names
}

func TestShipCells_Line(t *testing.T) {
	tests := []struct {
		name        string
		shipType    ShipType
		ref         string
		orientation Orientation
		want        []string
	}{
		{"CA-2 cruiser horizontal", Cruiser, "A1", Horizontal, []string{"A1", "B1", "C1"}},
		{"CA-2 cruiser vertical", Cruiser, "A1", Vertical, []string{"A1", "A2", "A3"}},
		{"CA-3 submarine horizontal", Submarine, "E5", Horizontal, []string{"E5"}},
		{"CA-3 submarine vertical", Submarine, "E5", Vertical, []string{"E5"}},
		{"destroyer vertical", Destroyer, "C3", Vertical, []string{"C3", "C4"}},
		{"battleship horizontal", Battleship, "B2", Horizontal, []string{"B2", "C2", "D2", "E2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cells, err := ShipCells(tt.shipType, mustCoord(t, tt.ref), tt.orientation)
			if err != nil {
				t.Fatalf("ShipCells returned error: %v", err)
			}
			if got := cellNames(cells); !slices.Equal(got, tt.want) {
				t.Errorf("ShipCells(%v, %s, %v) = %v, want %v", tt.shipType, tt.ref, tt.orientation, got, tt.want)
			}
		})
	}
}

// CA-11 (line ships): only Horizontal and Vertical apply.
func TestShipCells_LineInvalidOrientation(t *testing.T) {
	for _, o := range []Orientation{Up, Down, Left, Right} {
		_, err := ShipCells(Cruiser, Coord{}, o)
		if !errors.Is(err, ErrInvalidOrientation) {
			t.Errorf("ShipCells(Cruiser, A1, %v) error = %v, want ErrInvalidOrientation", o, err)
		}
	}
}

func TestOrientation_String(t *testing.T) {
	tests := []struct {
		o    Orientation
		want string
	}{
		{Horizontal, "Horizontal"},
		{Vertical, "Vertical"},
		{Up, "Up"},
		{Down, "Down"},
		{Left, "Left"},
		{Right, "Right"},
		{Orientation(99), "Orientation(99)"},
	}
	for _, tt := range tests {
		if got := tt.o.String(); got != tt.want {
			t.Errorf("Orientation(%d).String() = %q, want %q", int(tt.o), got, tt.want)
		}
	}
}

// CA-5: the T-shaped carrier, with the reference at the middle of the bar.
func TestShipCells_AircraftCarrier(t *testing.T) {
	tests := []struct {
		orientation Orientation
		want        []string
	}{
		{Down, []string{"D5", "E5", "F5", "E6", "E7"}},
		{Up, []string{"D5", "E5", "F5", "E4", "E3"}},
		{Right, []string{"E4", "E5", "E6", "F5", "G5"}},
		{Left, []string{"E4", "E5", "E6", "D5", "C5"}},
	}
	for _, tt := range tests {
		t.Run(tt.orientation.String(), func(t *testing.T) {
			cells, err := ShipCells(AircraftCarrier, mustCoord(t, "E5"), tt.orientation)
			if err != nil {
				t.Fatalf("ShipCells returned error: %v", err)
			}
			want := slices.Clone(tt.want)
			slices.Sort(want)
			if got := cellNames(cells); !slices.Equal(got, want) {
				t.Errorf("ShipCells(AircraftCarrier, E5, %v) = %v, want %v", tt.orientation, got, want)
			}
		})
	}
}

// CA-11 (carrier): Horizontal and Vertical do not apply to the T.
func TestShipCells_AircraftCarrierInvalidOrientation(t *testing.T) {
	for _, o := range []Orientation{Horizontal, Vertical} {
		_, err := ShipCells(AircraftCarrier, mustCoord(t, "E5"), o)
		if !errors.Is(err, ErrInvalidOrientation) {
			t.Errorf("ShipCells(AircraftCarrier, E5, %v) error = %v, want ErrInvalidOrientation", o, err)
		}
	}
}
