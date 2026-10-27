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
	"Comandos: aleatorio sorteia o resto da frota, desfazer tira o último navio, recomecar tira todos os navios, sair desiste.\n" + //nolint:misspell // Portuguese text for the player
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
	return placementPlayStyle(t, input, Plain)
}

// placementPlayStyle is placementPlay drawing the boards with style.
func placementPlayStyle(t *testing.T, input string, style Style) string {
	t.Helper()
	var sb strings.Builder
	player := board.NewBoard(board.DefaultFleet())
	err := PlayVsMachine(strings.NewReader(input), &sb, player, machineEnemyBoard(t), rand.New(rand.NewPCG(1, 0)), style, board.RandomShot)
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

const (
	placementCarrierPrompt    = "Porta-Aviões, 5 casas (cima, baixo, esquerda, direita): "
	placementBattleshipPrompt = "Encouraçado, 4 casas (horizontal, vertical): "
)

// placementWithCarrier returns an empty default-fleet board and the same
// board with the carrier at B1 pointing down, as "B1 baixo" places it.
func placementWithCarrier(t *testing.T) (empty, withCarrier *board.Board) {
	t.Helper()
	empty = board.NewBoard(board.DefaultFleet())
	withCarrier = board.NewBoard(board.DefaultFleet())
	if err := withCarrier.Place(board.AircraftCarrier, board.Coord{Col: 1}, board.Down); err != nil {
		t.Fatalf("Place: %v", err)
	}
	return empty, withCarrier
}

// CA-5: "aleatorio" keeps the ships already placed, draws the rest with the
// game's random source and starts the game.
func TestPlayVsMachine_PlacementRandom(t *testing.T) {
	empty, replay := placementWithCarrier(t)
	want := placementIntro + placementAsk(t, empty, placementCarrierPrompt) + placementAsk(t, replay, placementBattleshipPrompt)
	if err := board.FillRandom(replay, rand.New(rand.NewPCG(1, 0))); err != nil {
		t.Fatalf("FillRandom: %v", err)
	}
	want += "Frota posicionada!\n\n" + machineScreen(t, replay, machineEnemyBoard(t))

	if got := placementPlay(t, "B1 baixo\nAleatorio\n"); !strings.HasPrefix(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-6: when the rest of the fleet does not fit, "aleatorio" says so and
// asks for the same ship. At most 25 submarines fit on the board.
func TestPlayVsMachine_PlacementRandomDoesNotFit(t *testing.T) {
	const prompt = "Submarino, 1 casa: "
	fleet := []board.FleetEntry{{Type: board.Submarine, Count: 26}}
	want := placementIntro + placementAsk(t, board.NewBoard(fleet), prompt) +
		"Não foi possível sortear o resto da frota. Mude a posição dos navios com recomecar.\n" + prompt

	var sb strings.Builder
	err := PlayVsMachine(strings.NewReader("aleatorio\n"), &sb, board.NewBoard(fleet), machineEnemyBoard(t), rand.New(rand.NewPCG(1, 0)), Plain, board.RandomShot)
	if err != nil {
		t.Fatalf("PlayVsMachine: %v", err)
	}
	if got := sb.String(); !strings.HasPrefix(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-7: "recomecar" removes every ship and asks for the carrier again.
func TestPlayVsMachine_PlacementRestart(t *testing.T) {
	empty, withCarrier := placementWithCarrier(t)
	want := placementIntro + placementAsk(t, empty, placementCarrierPrompt) +
		placementAsk(t, withCarrier, placementBattleshipPrompt) +
		placementAsk(t, empty, placementCarrierPrompt)

	if got := placementPlay(t, "B1 baixo\nRECOMECAR\n"); !strings.HasPrefix(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-1, CA-2 and CA-3 (spec 014): "desfazer" removes the last ship placed
// and asks for it again with the board, one ship at a time; on an empty
// board it says so and asks for the same ship.
func TestPlayVsMachine_PlacementUndo(t *testing.T) {
	empty, withCarrier := placementWithCarrier(t)
	withTwo := board.NewBoard(board.DefaultFleet())
	for _, p := range placementExample[:2] {
		c, o, err := parsePlacement(p.line, p.shipType)
		if err != nil {
			t.Fatalf("parsePlacement(%q): %v", p.line, err)
		}
		if err := withTwo.Place(p.shipType, c, o); err != nil {
			t.Fatalf("Place(%q): %v", p.line, err)
		}
	}
	want := placementIntro + placementAsk(t, empty, placementCarrierPrompt) +
		"Não há navio para desfazer.\n" + placementCarrierPrompt +
		placementAsk(t, withCarrier, placementBattleshipPrompt) +
		placementAsk(t, withTwo, "Cruzador, 3 casas (horizontal, vertical): ") +
		placementAsk(t, withCarrier, placementBattleshipPrompt) +
		placementAsk(t, empty, placementCarrierPrompt) +
		"Você desistiu antes da partida.\n"

	if got := placementPlay(t, "desfazer\nB1 baixo\nE1 h\nDESFAZER\n desfazer \nsair\n"); got != want {
		t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, want)
	}
}

// CA-4, CA-5 and CA-6 (spec 014): with 7 battleships and a cruiser, six
// battleships two per row in rows 2, 5 and 8 leave only row 10. The seventh
// at D10 leaves no room for the cruiser, so the warning comes before the
// board; after desfazer it is gone, and the seventh at C10 leaves exactly
// H10 to J10 for the cruiser, with no warning.
func TestPlayVsMachine_PlacementWarning(t *testing.T) {
	const (
		battleshipPrompt = "Encouraçado, 4 casas (horizontal, vertical): "
		cruiserPrompt    = "Cruzador, 3 casas (horizontal, vertical): "
		warning          = "Atenção: o resto da frota não cabe mais no tabuleiro. Use desfazer ou recomecar.\n" //nolint:misspell // Portuguese text for the player
	)
	fleet := []board.FleetEntry{{Type: board.Battleship, Count: 7}, {Type: board.Cruiser, Count: 1}}
	replay := board.NewBoard(fleet)
	want := placementIntro
	var input []string
	for _, line := range []string{"A2 h", "F2 h", "A5 h", "F5 h", "A8 h", "F8 h"} {
		want += placementAsk(t, replay, battleshipPrompt)
		c, o, err := parsePlacement(line, board.Battleship)
		if err != nil {
			t.Fatalf("parsePlacement(%q): %v", line, err)
		}
		if err := replay.Place(board.Battleship, c, o); err != nil {
			t.Fatalf("Place(%q): %v", line, err)
		}
		input = append(input, line)
	}
	sixShips := placementAsk(t, replay, battleshipPrompt)
	want += sixShips
	if err := replay.Place(board.Battleship, board.Coord{Col: 3, Row: 9}, board.Horizontal); err != nil {
		t.Fatalf("Place(D10): %v", err)
	}
	want += warning + placementAsk(t, replay, cruiserPrompt) + sixShips
	replay.RemoveLast()
	if err := replay.Place(board.Battleship, board.Coord{Col: 2, Row: 9}, board.Horizontal); err != nil {
		t.Fatalf("Place(C10): %v", err)
	}
	want += placementAsk(t, replay, cruiserPrompt) + "Você desistiu antes da partida.\n"
	input = append(input, "D10 h", "desfazer", "C10 h", "sair")

	var sb strings.Builder
	err := PlayVsMachine(strings.NewReader(strings.Join(input, "\n")+"\n"), &sb, board.NewBoard(fleet), machineEnemyBoard(t), rand.New(rand.NewPCG(1, 0)), Plain, board.RandomShot)
	if err != nil {
		t.Fatalf("PlayVsMachine: %v", err)
	}
	if got := sb.String(); got != want {
		t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, want)
	}
}

// CA-8: "sair" or the end of the input gives up before the game.
func TestPlayVsMachine_PlacementGiveUp(t *testing.T) {
	empty, withCarrier := placementWithCarrier(t)
	asked := placementIntro + placementAsk(t, empty, placementCarrierPrompt) + placementAsk(t, withCarrier, placementBattleshipPrompt)
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"sair", "B1 baixo\nsair\n", asked + "Você desistiu antes da partida.\n"},
		{"end of input", "B1 baixo\n", asked + "\nVocê desistiu antes da partida.\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := placementPlay(t, tt.input); got != tt.want {
				t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, tt.want)
			}
		})
	}
}

// CA-7 (spec 010): with Colored, the boards shown while placing the fleet
// are drawn with colors.
func TestPlayVsMachine_PlacementColored(t *testing.T) {
	empty, withCarrier := placementWithCarrier(t)
	want := placementIntro + renderColored(t, RenderBoard, empty) + "\n" + placementCarrierPrompt +
		renderColored(t, RenderBoard, withCarrier) + "\n" + placementBattleshipPrompt

	const input = "B1 baixo\naleatorio\nA1\nsair\n"
	got := placementPlayStyle(t, input, Colored)
	if !strings.HasPrefix(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant prefix\n%s", got, want)
	}
	assertBoardsColored(t, got, placementPlay(t, input))
	assertUncolored(t, got, placementInstructions, placementCarrierPrompt, placementBattleshipPrompt, "Frota posicionada!")
}

// CA-4 (spec 015): with Colored, an invalid position and a random draw that
// does not fit are yellow, followed by the uncolored prompt.
func TestPlayVsMachine_PlacementColoredErrors(t *testing.T) {
	got := placementPlayStyle(t, "xyz\nB1 baixo\nA1 horizontal\nsair\n", Colored)
	for _, want := range []string{
		colored("33", "Posição inválida: digite a casa e a orientação, por exemplo E5 vertical.") + "\n" + placementCarrierPrompt,
		colored("33", "Posição ocupada: o Encouraçado ficaria em cima de outro navio.") + "\n" + placementBattleshipPrompt,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("PlayVsMachine() =\n%s\nwant it to contain %q", got, want)
		}
	}

	const prompt = "Submarino, 1 casa: "
	fleet := []board.FleetEntry{{Type: board.Submarine, Count: 26}}
	var sb strings.Builder
	err := PlayVsMachine(strings.NewReader("aleatorio\n"), &sb, board.NewBoard(fleet), machineEnemyBoard(t), rand.New(rand.NewPCG(1, 0)), Colored, board.RandomShot)
	if err != nil {
		t.Fatalf("PlayVsMachine: %v", err)
	}
	want := prompt + colored("33", "Não foi possível sortear o resto da frota. Mude a posição dos navios com recomecar.") + "\n" + prompt
	if got := sb.String(); !strings.Contains(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant it to contain %q", got, want)
	}
}

// CA-4 (spec 015) for the messages of spec 014: desfazer on an empty board
// and the warning that the rest of the fleet no longer fits are yellow.
func TestPlayVsMachine_PlacementColoredUndoAndWarning(t *testing.T) {
	got := placementPlayStyle(t, "desfazer\nsair\n", Colored)
	if want := colored("33", "Não há navio para desfazer.") + "\n" + placementCarrierPrompt; !strings.Contains(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant it to contain %q", got, want)
	}

	// Six battleships leave only row 10, and the seventh at D10 leaves no
	// room for the cruiser.
	fleet := []board.FleetEntry{{Type: board.Battleship, Count: 7}, {Type: board.Cruiser, Count: 1}}
	input := "A2 h\nF2 h\nA5 h\nF5 h\nA8 h\nF8 h\nD10 h\nsair\n"
	var sb strings.Builder
	err := PlayVsMachine(strings.NewReader(input), &sb, board.NewBoard(fleet), machineEnemyBoard(t), rand.New(rand.NewPCG(1, 0)), Colored, board.RandomShot)
	if err != nil {
		t.Fatalf("PlayVsMachine: %v", err)
	}
	warning := colored("33", "Atenção: o resto da frota não cabe mais no tabuleiro. Use desfazer ou recomecar.") + "\n" //nolint:misspell // Portuguese text for the player
	if got := sb.String(); !strings.Contains(got, warning) {
		t.Errorf("PlayVsMachine() =\n%s\nwant it to contain %q", got, warning)
	}
}
