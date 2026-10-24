package board

import "testing"

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
