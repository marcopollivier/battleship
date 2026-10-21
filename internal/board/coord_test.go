package board

import (
	"errors"
	"strconv"
	"strings"
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

func TestParseCoord_Valid(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Coord
	}{
		{"CA-1 first cell", "A1", Coord{Col: 0, Row: 0}},
		{"CA-2 last cell", "J10", Coord{Col: 9, Row: 9}},
		{"CA-3 middle cell", "C7", Coord{Col: 2, Row: 6}},
		{"CA-8 lowercase first cell", "a1", Coord{Col: 0, Row: 0}},
		{"CA-8 lowercase last cell", "j10", Coord{Col: 9, Row: 9}},
		{"CA-9 surrounding spaces", " A1 ", Coord{Col: 0, Row: 0}},
		{"CA-9 surrounding tab and newline", "\tA1\n", Coord{Col: 0, Row: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCoord(tt.in)
			if err != nil {
				t.Fatalf("ParseCoord(%q) returned error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("ParseCoord(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}

// CA-4: every cell survives String -> ParseCoord unchanged.
func TestCoord_RoundTrip(t *testing.T) {
	for col := range Size {
		for row := range Size {
			c := Coord{Col: col, Row: row}
			got, err := ParseCoord(c.String())
			if err != nil {
				t.Fatalf("ParseCoord(%q) returned error: %v", c.String(), err)
			}
			if got != c {
				t.Errorf("ParseCoord(%q) = %+v, want %+v", c.String(), got, c)
			}
		}
	}
}

func TestParseCoord_InvalidFormat(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"CA-7 empty", ""},
		{"CA-7 only spaces", "   "},
		{"CA-7 letter only", "A"},
		{"CA-7 number only", "1"},
		{"CA-7 two letters", "AA1"},
		{"CA-7 trailing letter", "A1B"},
		{"CA-7 dash", "A-1"},
		{"CA-7 symbol", "@1"},
		{"CA-7 non-ASCII letter", "Ç1"},
		{"CA-10 inner space", "A 1"},
		{"CA-10 reversed", "1A"},
		{"CA-10 leading zero", "A01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCoord(tt.in)
			if !errors.Is(err, ErrInvalidFormat) {
				t.Fatalf("ParseCoord(%q) error = %v, want ErrInvalidFormat", tt.in, err)
			}
			if quoted := strconv.Quote(tt.in); !strings.Contains(err.Error(), quoted) {
				t.Errorf("ParseCoord(%q) error = %q, want it to contain %s", tt.in, err.Error(), quoted)
			}
		})
	}
}
