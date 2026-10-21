package board

import (
	"strconv"
	"testing"
)

func TestCoord_String(t *testing.T) {
	for col := range Size {
		for row := range Size {
			c := Coord{Col: col, Row: row}
			want := string(rune('A'+col)) + strconv.Itoa(row+1)
			if got := c.String(); got != want {
				t.Errorf("Coord{%d, %d}.String() = %q, want %q", col, row, got, want)
			}
		}
	}
}

func TestCoord_String_OutOfBounds(t *testing.T) {
	tests := []struct {
		c    Coord
		want string
	}{
		{Coord{Col: -1, Row: 0}, "Coord(-1,0)"},
		{Coord{Col: 0, Row: -1}, "Coord(0,-1)"},
		{Coord{Col: Size, Row: 0}, "Coord(10,0)"},
		{Coord{Col: 0, Row: Size}, "Coord(0,10)"},
	}
	for _, tt := range tests {
		if got := tt.c.String(); got != tt.want {
			t.Errorf("Coord{%d, %d}.String() = %q, want %q", tt.c.Col, tt.c.Row, got, tt.want)
		}
	}
}
