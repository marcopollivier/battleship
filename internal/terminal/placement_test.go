package terminal

import (
	"errors"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
)

// CA-3: a cell and an orientation by name or first letter, in any case; a
// one-cell ship needs no orientation.
func TestParsePlacement(t *testing.T) {
	tests := []struct {
		line     string
		shipType board.ShipType
		want     string
		wantO    board.Orientation
	}{
		{"e5 V", board.Battleship, "E5", board.Vertical},
		{"E5 vertical", board.Battleship, "E5", board.Vertical},
		{" a1   HORIZONTAL ", board.Cruiser, "A1", board.Horizontal},
		{"C3 b", board.AircraftCarrier, "C3", board.Down},
		{"C3 BAIXO", board.AircraftCarrier, "C3", board.Down},
		{"C3 cima", board.AircraftCarrier, "C3", board.Up},
		{"C3 e", board.AircraftCarrier, "C3", board.Left},
		{"C3 Direita", board.AircraftCarrier, "C3", board.Right},
		{"A1", board.Submarine, "A1", board.Horizontal},
		{"A1 v", board.Submarine, "A1", board.Vertical},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			c, o, err := parsePlacement(tt.line, tt.shipType)
			if err != nil {
				t.Fatalf("parsePlacement(%q, %v) error = %v", tt.line, tt.shipType, err)
			}
			if c.String() != tt.want || o != tt.wantO {
				t.Errorf("parsePlacement(%q, %v) = %v %v, want %s %v", tt.line, tt.shipType, c, o, tt.want, tt.wantO)
			}
		})
	}
}

// CA-4 (format and orientation): malformed lines and orientations that do
// not exist or do not apply to the ship are rejected.
func TestParsePlacement_Invalid(t *testing.T) {
	tests := []struct {
		line     string
		shipType board.ShipType
		want     error
	}{
		{"", board.Cruiser, errInvalidPlacement},
		{"E5", board.Cruiser, errInvalidPlacement},
		{"E5 vertical sim", board.Cruiser, errInvalidPlacement},
		{"5E vertical", board.Cruiser, board.ErrInvalidFormat},
		{"K1 vertical", board.Cruiser, board.ErrOutOfBounds},
		{"E5 diagonal", board.Cruiser, board.ErrInvalidOrientation},
		{"E5 baixo", board.Cruiser, board.ErrInvalidOrientation},
		{"C3 horizontal", board.AircraftCarrier, board.ErrInvalidOrientation},
		{"C3 x", board.AircraftCarrier, board.ErrInvalidOrientation},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			_, _, err := parsePlacement(tt.line, tt.shipType)
			if !errors.Is(err, tt.want) {
				t.Errorf("parsePlacement(%q, %v) error = %v, want %v", tt.line, tt.shipType, err, tt.want)
			}
		})
	}
}

const placementIntro = "Posicione a sua frota. Para cada navio, digite a casa e a orientação, por exemplo E5 vertical.\n" +
	"A casa é a primeira do navio (a de cima ou a da esquerda). No Porta-Aviões, é o meio da barra do T.\n" +
	"Comandos: aleatorio sorteia o resto da frota, recomecar tira todos os navios, sair desiste.\n" + //nolint:misspell // Portuguese text for the player
	"\n"

// placementExample is the example board from spec 003 as the player types it.
var placementExample = []struct {
	shipType board.ShipType
	line     string
	prompt   string
}{
	{board.AircraftCarrier, "B1 baixo", "Porta-Aviões, 5 casas (cima, baixo, esquerda, direita): "},
	{board.Battleship, "E1 h", "Encouraçado, 4 casas (horizontal, vertical): "},
	{board.Cruiser, "J1 v", "Cruzador, 3 casas (horizontal, vertical): "},
	{board.Destroyer, "A5 horizontal", "Destroyer, 2 casas (horizontal, vertical): "},
	{board.Destroyer, "E3 h", "Destroyer, 2 casas (horizontal, vertical): "},
	{board.Submarine, "H5", "Submarino, 1 casa: "},
	{board.Submarine, "A8", "Submarino, 1 casa: "},
	{board.Submarine, "E8", "Submarino, 1 casa: "},
}

// placementPlay runs PlayVsMachine with an empty default-fleet board for the
// player and returns what it wrote.
func placementPlay(t *testing.T, input string) string {
	t.Helper()
	var sb strings.Builder
	player := board.NewBoard(board.DefaultFleet())
	err := PlayVsMachine(strings.NewReader(input), &sb, player, machineEnemyBoard(t), rand.New(rand.NewPCG(1, 0)))
	if err != nil {
		t.Fatalf("PlayVsMachine: %v", err)
	}
	return sb.String()
}

// placementAsk is the board followed by the prompt for the next ship.
func placementAsk(t *testing.T, b *board.Board, prompt string) string {
	t.Helper()
	return soloReveal(t, b) + "\n" + prompt
}

// CA-1: placing starts with the instructions, the empty board and the
// aircraft carrier with its four orientations.
func TestPlayVsMachine_PlacementStart(t *testing.T) {
	want := placementIntro + placementAsk(t, board.NewBoard(board.DefaultFleet()), placementExample[0].prompt)
	if got := placementPlay(t, ""); !strings.HasPrefix(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-2: the ships are asked largest first, each after the board, and the
// game starts with the ships where the player put them.
func TestPlayVsMachine_PlaceFleet(t *testing.T) {
	replay := board.NewBoard(board.DefaultFleet())
	want := placementIntro
	var input []string
	for _, p := range placementExample {
		want += placementAsk(t, replay, p.prompt)
		c, o, err := parsePlacement(p.line, p.shipType)
		if err != nil {
			t.Fatalf("parsePlacement(%q): %v", p.line, err)
		}
		if err := replay.Place(p.shipType, c, o); err != nil {
			t.Fatalf("Place(%q): %v", p.line, err)
		}
		input = append(input, p.line)
	}
	want += "Frota posicionada!\n\n" + machineScreen(t, replay, machineEnemyBoard(t))

	if got := placementPlay(t, strings.Join(input, "\n")+"\n"); !strings.HasPrefix(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-4: each invalid placement shows its message and asks for the same ship
// again, without the board.
func TestPlayVsMachine_PlacementErrors(t *testing.T) {
	const (
		carrierPrompt    = "Porta-Aviões, 5 casas (cima, baixo, esquerda, direita): "
		battleshipPrompt = "Encouraçado, 4 casas (horizontal, vertical): "
		badFormat        = "Posição inválida: digite a casa e a orientação, por exemplo E5 vertical.\n"
	)
	replay := board.NewBoard(board.DefaultFleet())
	want := placementIntro + placementAsk(t, replay, carrierPrompt) +
		"Orientação inválida para o Porta-Aviões: use cima, baixo, esquerda ou direita.\n" + carrierPrompt
	if err := replay.Place(board.AircraftCarrier, board.Coord{Col: 1}, board.Down); err != nil {
		t.Fatalf("Place: %v", err)
	}
	want += placementAsk(t, replay, battleshipPrompt) +
		badFormat + battleshipPrompt +
		badFormat + battleshipPrompt +
		badFormat + battleshipPrompt +
		"Orientação inválida para o Encouraçado: use horizontal ou vertical.\n" + battleshipPrompt +
		"Fora do tabuleiro: o Encouraçado não cabe nessa posição.\n" + battleshipPrompt +
		"Fora do tabuleiro: o Encouraçado não cabe nessa posição.\n" + battleshipPrompt +
		"Posição ocupada: o Encouraçado ficaria em cima de outro navio.\n" + battleshipPrompt +
		"Muito perto: os navios não podem se encostar, nem pela diagonal.\n" + battleshipPrompt
	if err := replay.Place(board.Battleship, board.Coord{Col: 4}, board.Horizontal); err != nil {
		t.Fatalf("Place: %v", err)
	}
	want += placementAsk(t, replay, "Cruzador, 3 casas (horizontal, vertical): ")

	input := "C3 horizontal\nB1 baixo\nE5\nE5 vertical sim\n5E vertical\nE5 diagonal\nJ9 vertical\nK1 h\nA2 h\nD1 v\nE1 h\n"
	if got := placementPlay(t, input); !strings.HasPrefix(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant prefix\n%s", got, want)
	}
}

func TestJoinOr(t *testing.T) {
	tests := []struct {
		names []string
		want  string
	}{
		{nil, ""},
		{[]string{"cima"}, "cima"},
		{[]string{"horizontal", "vertical"}, "horizontal ou vertical"},
		{[]string{"cima", "baixo", "esquerda", "direita"}, "cima, baixo, esquerda ou direita"},
	}
	for _, tt := range tests {
		if got := joinOr(tt.names); got != tt.want {
			t.Errorf("joinOr(%q) = %q, want %q", tt.names, got, tt.want)
		}
	}
}
