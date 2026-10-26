package terminal

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

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

// soloReveal is the own board RenderBoard draws for b, revealing its ships.
func soloReveal(t *testing.T, b *board.Board) string {
	t.Helper()
	var sb strings.Builder
	if err := RenderBoard(&sb, b); err != nil {
		t.Fatalf("RenderBoard: %v", err)
	}
	return sb.String()
}

// CA-5: "sair" in any case, or the end of the input, gives up and reveals
// the fleet with the shots received so far.
func TestPlaySolo_GiveUp(t *testing.T) {
	afterShot := func(t *testing.T, replay *board.Board) string {
		return soloStart(t, replay) + soloTurn(t, replay, "J10", "Água.") + "\n" + soloPrompt
	}
	tests := []struct {
		name  string
		input string
		want  func(t *testing.T, replay *board.Board) string
	}{
		{"sair after one shot", "J10\nsair\n", func(t *testing.T, replay *board.Board) string {
			return afterShot(t, replay) + "Você desistiu depois de 1 disparo. A frota era:\n\n"
		}},
		{"SAIR with spaces and no shots", " SAIR \n", func(t *testing.T, replay *board.Board) string {
			return soloStart(t, replay) + "Você desistiu depois de 0 disparos. A frota era:\n\n"
		}},
		{"end of input after two shots", "J10\nE1", func(t *testing.T, replay *board.Board) string {
			return afterShot(t, replay) + soloTurn(t, replay, "E1", "Acerto!") + "\n" + soloPrompt +
				"\nVocê desistiu depois de 2 disparos. A frota era:\n\n"
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replay := attackExampleBoard(t)
			want := tt.want(t, replay)
			want += soloReveal(t, replay)

			got := soloPlay(t, attackExampleBoard(t), tt.input)
			if got != want {
				t.Errorf("PlaySolo() =\n%s\nwant\n%s", got, want)
			}
		})
	}
}

// A win with a single shot uses the singular.
func TestPlaySolo_WinOneShot(t *testing.T) {
	b := board.NewBoard([]board.FleetEntry{{Type: board.Submarine, Count: 1}})
	if err := b.Place(board.Submarine, board.Coord{}, board.Horizontal); err != nil {
		t.Fatalf("Place: %v", err)
	}
	got := soloPlay(t, b, "A1\n")
	if want := "Você afundou toda a frota em 1 disparo!\n"; !strings.HasSuffix(got, want) {
		t.Errorf("PlaySolo() =\n%s\nwant suffix %q", got, want)
	}
}

// soloFailingWriter fails every write after the first limit bytes.
type soloFailingWriter struct{ limit int }

func (w *soloFailingWriter) Write(p []byte) (int, error) {
	if len(p) > w.limit {
		return 0, errSoloWrite
	}
	w.limit -= len(p)
	return len(p), nil
}

var errSoloWrite = errors.New("write failed")

// Failures to read the input or write the screen end the game with that error.
func TestPlaySolo_IOErrors(t *testing.T) {
	errRead := errors.New("read failed")
	tests := []struct {
		name string
		in   io.Reader
		out  io.Writer
		want error
	}{
		{"read", iotest.ErrReader(errRead), io.Discard, errRead},
		{"write at start", strings.NewReader("J10\n"), &soloFailingWriter{}, errSoloWrite},
		{"write after a shot", strings.NewReader("J10\n"), &soloFailingWriter{limit: 400}, errSoloWrite},
		{"write when giving up", strings.NewReader("sair\n"), &soloFailingWriter{limit: 400}, errSoloWrite},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PlaySolo(tt.in, tt.out, attackExampleBoard(t)); !errors.Is(err, tt.want) {
				t.Errorf("PlaySolo() error = %v, want %v", err, tt.want)
			}
		})
	}
}
