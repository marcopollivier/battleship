package terminal

import (
	"errors"
	"strings"
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
)

// soloShipCells are the 19 ship cells of the example board from spec 003.
var soloShipCells = []string{
	"A1", "B1", "C1", "B2", "B3",
	"E1", "F1", "G1", "H1",
	"J1", "J2", "J3",
	"A5", "B5", "E3", "F3",
	"H5", "A8", "E8",
}

// soloPlay runs PlaySolo on b with input and returns what it wrote.
func soloPlay(t *testing.T, b *board.Board, input string) string {
	t.Helper()
	var sb strings.Builder
	if err := PlaySolo(strings.NewReader(input), &sb, b); err != nil {
		t.Fatalf("PlaySolo: %v", err)
	}
	return sb.String()
}

// soloStart is the screen at the start of a game on b.
func soloStart(t *testing.T, b *board.Board) string {
	t.Helper()
	return attackRender(t, b) + "\n" + soloPrompt
}

// soloTurn fires at cell on replay and returns the message and the attack
// board PlaySolo should write for that shot.
func soloTurn(t *testing.T, replay *board.Board, cell, message string) string {
	t.Helper()
	attackFire(t, replay, cell)
	return message + "\n\n" + attackRender(t, replay)
}

const soloPrompt = "Disparo: "

// CA-1: the game starts with the untouched attack board and the prompt.
func TestPlaySolo_Start(t *testing.T) {
	b := attackExampleBoard(t)
	got := soloPlay(t, b, "")
	if want := soloStart(t, attackExampleBoard(t)); !strings.HasPrefix(got, want) {
		t.Errorf("PlaySolo() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-2: each valid shot shows its result followed by the updated attack board.
func TestPlaySolo_Shots(t *testing.T) {
	replay := attackExampleBoard(t)
	want := soloStart(t, replay) +
		soloTurn(t, replay, "J10", "Água.") + "\n" + soloPrompt +
		soloTurn(t, replay, "E1", "Acerto!") + "\n" + soloPrompt +
		soloTurn(t, replay, "H5", "Afundou o Submarino!") + "\n" + soloPrompt

	got := soloPlay(t, attackExampleBoard(t), "J10\n e1 \nh5\n")
	if !strings.HasPrefix(got, want) {
		t.Errorf("PlaySolo() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-4: sinking the last ship ends the game with the final board and the
// number of valid shots; input after that is not read.
func TestPlaySolo_Win(t *testing.T) {
	cells := append([]string{"J10", "D7"}, soloShipCells...)
	replay := attackExampleBoard(t)
	var want strings.Builder
	want.WriteString(soloStart(t, replay))
	for i, cell := range cells {
		c, err := board.ParseCoord(cell)
		if err != nil {
			t.Fatalf("ParseCoord(%q): %v", cell, err)
		}
		shot, err := replay.Fire(c)
		if err != nil {
			t.Fatalf("Fire(%s): %v", cell, err)
		}
		want.WriteString(ShotMessage(shot) + "\n\n" + attackRender(t, replay) + "\n")
		if i < len(cells)-1 {
			want.WriteString(soloPrompt)
		}
	}
	want.WriteString("Você afundou toda a frota em 21 disparos!\n")

	got := soloPlay(t, attackExampleBoard(t), strings.Join(cells, "\n")+"\nA10\n")
	if got != want.String() {
		t.Errorf("PlaySolo() =\n%s\nwant\n%s", got, want.String())
	}
}

// CA-3: invalid input shows a message and asks again, without a board and
// without counting as a shot.
func TestPlaySolo_InvalidInput(t *testing.T) {
	replay := attackExampleBoard(t)
	want := soloStart(t, replay) +
		soloTurn(t, replay, "J10", "Água.") + "\n" + soloPrompt +
		"Coordenada inválida: use uma letra de A a J e um número de 1 a 10, por exemplo B7.\n" + soloPrompt +
		"Coordenada inválida: use uma letra de A a J e um número de 1 a 10, por exemplo B7.\n" + soloPrompt +
		"Fora do tabuleiro: use uma letra de A a J e um número de 1 a 10.\n" + soloPrompt +
		"Fora do tabuleiro: use uma letra de A a J e um número de 1 a 10.\n" + soloPrompt +
		"Você já atirou em J10.\n" + soloPrompt +
		soloTurn(t, replay, "E1", "Acerto!") + "\n" + soloPrompt

	got := soloPlay(t, attackExampleBoard(t), "J10\n7B\n\nK1\nA11\n j10\nE1\n")
	if !strings.HasPrefix(got, want) {
		t.Errorf("PlaySolo() =\n%s\nwant prefix\n%s", got, want)
	}
}

// An error that is not a mistake of the player, such as shooting at an
// incomplete board, ends the game and is returned.
func TestPlaySolo_UnexpectedError(t *testing.T) {
	var sb strings.Builder
	err := PlaySolo(strings.NewReader("A1\n"), &sb, board.NewBoard(board.DefaultFleet()))
	if !errors.Is(err, board.ErrIncompleteBoard) {
		t.Errorf("PlaySolo() error = %v, want ErrIncompleteBoard", err)
	}
}
