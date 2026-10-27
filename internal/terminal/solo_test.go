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
	return soloPlayStyle(t, b, input, Plain)
}

// soloPlayStyle is soloPlay drawing the boards with style.
func soloPlayStyle(t *testing.T, b *board.Board, input string, style Style) string {
	t.Helper()
	var sb strings.Builder
	if err := PlaySolo(strings.NewReader(input), &sb, b, Display{Style: style}); err != nil {
		t.Fatalf("PlaySolo: %v", err)
	}
	return sb.String()
}

// soloStart is the screen at the start of a game on b.
func soloStart(t *testing.T, b *board.Board) string {
	t.Helper()
	return attackRender(t, b) + "\n" + soloPrompt
}

// soloTurn fires at cell on replay and returns the screen PlaySolo should
// write after that shot: the attack board, the result citing the cell, such
// as "Seu disparo em B7: Acerto!", and the prompt.
func soloTurn(t *testing.T, replay *board.Board, cell, message string) string {
	t.Helper()
	attackFire(t, replay, cell)
	return attackRender(t, replay) + "\nSeu disparo em " + cell + ": " + message + "\n" + soloPrompt
}

// soloFinal is the last screen of a solo game on b: the board with its ships
// and the shots received, the short legend, the events and the message.
func soloFinal(t *testing.T, b *board.Board, style Style, message string, events ...string) string {
	t.Helper()
	lines := append(ownGrid(b, style), "", "Legenda: ~ água, S D C E P navio, X atingido, o tiro na água", "")
	lines = append(lines, events...)
	return strings.Join(lines, "\n") + "\n" + message + "\n"
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
		soloTurn(t, replay, "J10", "Água.") +
		soloTurn(t, replay, "E1", "Acerto!") +
		soloTurn(t, replay, "H5", "Afundou o Submarino!")

	got := soloPlay(t, attackExampleBoard(t), "J10\n e1 \nh5\n")
	if !strings.HasPrefix(got, want) {
		t.Errorf("PlaySolo() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-4 (spec 008) and CA-7 (spec 022): sinking the last ship ends the game
// on a last screen with the board, the last result and the number of valid
// shots; input after that is not read.
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
		result := "Seu disparo em " + cell + ": " + ShotMessage(shot)
		if i == len(cells)-1 {
			want.WriteString(soloFinal(t, replay, Plain, "Você afundou toda a frota em 21 disparos!", result))
			break
		}
		want.WriteString(attackRender(t, replay) + "\n" + result + "\n" + soloPrompt)
	}

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
		soloTurn(t, replay, "J10", "Água.") +
		"Coordenada inválida: use uma letra de A a J e um número de 1 a 10, por exemplo B7.\n" + soloPrompt +
		"Coordenada inválida: use uma letra de A a J e um número de 1 a 10, por exemplo B7.\n" + soloPrompt +
		"Fora do tabuleiro: use uma letra de A a J e um número de 1 a 10.\n" + soloPrompt +
		"Fora do tabuleiro: use uma letra de A a J e um número de 1 a 10.\n" + soloPrompt +
		"Você já atirou em J10.\n" + soloPrompt +
		soloTurn(t, replay, "E1", "Acerto!")

	got := soloPlay(t, attackExampleBoard(t), "J10\n7B\n\nK1\nA11\n j10\nE1\n")
	if !strings.HasPrefix(got, want) {
		t.Errorf("PlaySolo() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-7 (spec 022): with Clear, every screen of the solo game starts with the
// code that clears the terminal; invalid input does not clear it.
func TestPlaySolo_Clear(t *testing.T) {
	play := func(clear bool) string {
		var sb strings.Builder
		if err := PlaySolo(strings.NewReader("J10\n7B\nsair\n"), &sb, attackExampleBoard(t), Display{Style: Plain, Clear: clear}); err != nil {
			t.Fatalf("PlaySolo: %v", err)
		}
		return sb.String()
	}
	plain, cleared := play(false), play(true)
	if strings.Contains(plain, "\x1b[2J") {
		t.Errorf("without Clear, the game clears the terminal:\n%q", plain)
	}
	if got := strings.ReplaceAll(cleared, clearScreen, ""); got != plain {
		t.Errorf("with Clear, without the clear codes =\n%s\nwant\n%s", got, plain)
	}
	screens := strings.Split(cleared, clearScreen)
	if len(screens) != 4 || screens[0] != "" {
		t.Fatalf("with Clear, %d parts split by the clear code, want the output to start with it and 3 screens:\n%q", len(screens), cleared)
	}
	for i, screen := range screens[1:] {
		if !strings.HasPrefix(screen, "   A B C D E F G H I J\n") {
			t.Errorf("screen %d does not start with the board:\n%s", i+1, screen)
		}
	}
}

// An error that is not a mistake of the player, such as shooting at an
// incomplete board, ends the game and is returned.
func TestPlaySolo_UnexpectedError(t *testing.T) {
	var sb strings.Builder
	err := PlaySolo(strings.NewReader("A1\n"), &sb, board.NewBoard(board.DefaultFleet()), Display{Style: Plain})
	if !errors.Is(err, board.ErrIncompleteBoard) {
		t.Errorf("PlaySolo() error = %v, want ErrIncompleteBoard", err)
	}
}

// soloReveal is the own board RenderBoard draws for b, revealing its ships.
func soloReveal(t *testing.T, b *board.Board) string {
	t.Helper()
	var sb strings.Builder
	if err := RenderBoard(&sb, b, Plain); err != nil {
		t.Fatalf("RenderBoard: %v", err)
	}
	return sb.String()
}

// CA-5 (spec 008) and CA-7 (spec 022): "sair" in any case, or the end of
// the input, gives up on a last screen that reveals the fleet with the shots
// received so far.
func TestPlaySolo_GiveUp(t *testing.T) {
	afterShot := func(t *testing.T, replay *board.Board) string {
		return soloStart(t, replay) + soloTurn(t, replay, "J10", "Água.")
	}
	tests := []struct {
		name  string
		input string
		want  func(t *testing.T, replay *board.Board) string
	}{
		{"sair after one shot", "J10\nsair\n", func(t *testing.T, replay *board.Board) string {
			return afterShot(t, replay) + soloFinal(t, replay, Plain, "Você desistiu depois de 1 disparo.")
		}},
		{"SAIR with spaces and no shots", " SAIR \n", func(t *testing.T, replay *board.Board) string {
			return soloStart(t, replay) + soloFinal(t, replay, Plain, "Você desistiu depois de 0 disparos.")
		}},
		{"end of input after two shots", "J10\nE1", func(t *testing.T, replay *board.Board) string {
			return afterShot(t, replay) + soloTurn(t, replay, "E1", "Acerto!") +
				"\n" + soloFinal(t, replay, Plain, "Você desistiu depois de 2 disparos.")
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replay := attackExampleBoard(t)
			want := tt.want(t, replay)

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
			if err := PlaySolo(tt.in, tt.out, attackExampleBoard(t), Display{Style: Plain}); !errors.Is(err, tt.want) {
				t.Errorf("PlaySolo() error = %v, want %v", err, tt.want)
			}
		})
	}
}

// CA-7 (spec 010): with Colored, the attack boards and the revealed fleet
// are drawn with colors, and the screen texts are not.
func TestPlaySolo_Colored(t *testing.T) {
	replay := attackExampleBoard(t)
	want := renderColored(t, RenderAttackBoard, replay) + "\n" + soloPrompt
	attackFire(t, replay, "J10")
	want += renderColored(t, RenderAttackBoard, replay) + "\nSeu disparo em J10: " + colored("36", "Água.") + "\n" + soloPrompt +
		soloFinal(t, replay, Colored, "Você desistiu depois de 1 disparo.")

	got := soloPlayStyle(t, attackExampleBoard(t), "J10\nsair\n", Colored)
	if got != want {
		t.Errorf("PlaySolo() =\n%s\nwant\n%s", got, want)
	}
	assertUncolored(t, got, soloPrompt, "Seu disparo em J10: ", "Você desistiu depois de 1 disparo.", "Legenda: ~")
	input := strings.Join(soloShipCells, "\n") + "\n"
	assertBoardsColored(t, soloPlayStyle(t, attackExampleBoard(t), input, Colored), soloPlay(t, attackExampleBoard(t), input))
}

// CA-1 and CA-4 (spec 015): with Colored, a miss is cyan, a hit and a sunk
// ship bold red, the solo win bold green and invalid input yellow; the line
// breaks stay out of the colors.
func TestPlaySolo_ColoredMessages(t *testing.T) {
	cells := append([]string{"J10", "7B", "K1", "j10"}, soloShipCells...)
	got := soloPlayStyle(t, attackExampleBoard(t), strings.Join(cells, "\n")+"\n", Colored)
	for _, want := range []string{
		soloPrompt + colored("33", "Coordenada inválida: use uma letra de A a J e um número de 1 a 10, por exemplo B7.") + "\n" + soloPrompt,
		soloPrompt + colored("33", "Fora do tabuleiro: use uma letra de A a J e um número de 1 a 10.") + "\n" + soloPrompt,
		soloPrompt + colored("33", "Você já atirou em J10.") + "\n" + soloPrompt,
		"\nSeu disparo em J10: " + colored("36", "Água.") + "\n" + soloPrompt,
		"\nSeu disparo em A1: " + colored("1;31", "Acerto!") + "\n" + soloPrompt,
		"\nSeu disparo em H5: " + colored("1;31", "Afundou o Submarino!") + "\n" + soloPrompt,
		"\nSeu disparo em B3: " + colored("1;31", "Afundou o Porta-Aviões!") + "\n" + soloPrompt,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("PlaySolo() =\n%s\nwant it to contain %q", got, want)
		}
	}
	if !strings.HasSuffix(got, "\n"+colored("1;32", "Você afundou toda a frota em 20 disparos!")+"\n") {
		t.Errorf("PlaySolo() does not end with the colored win:\n%s", got)
	}
}
