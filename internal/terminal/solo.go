package terminal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/marcopollivier/battleship/internal/board"
)

const (
	// shotPrompt asks the player for the next shot.
	shotPrompt = "Disparo: "
	// quitCommand gives up the game, in any letter case.
	quitCommand = "sair"
)

// PlaySolo runs a solo game: it reads shots from in, one per line, and
// writes the attack board and messages to out until every ship of b is sunk
// or the player gives up, with "sair" or the end of the input. It returns an
// error if writing to out or reading from in fails.
func PlaySolo(in io.Reader, out io.Writer, b *board.Board) error {
	g := &playerAttack{out: out, target: b}
	g.showBoard()
	scanner := bufio.NewScanner(in)
	for {
		line, ok := g.prompt(scanner, shotPrompt)
		if !ok {
			g.giveUp("A frota era:")
			return g.err
		}
		if _, ok := g.shoot(line); !ok {
			if g.err != nil {
				return g.err
			}
			continue
		}
		if b.AllSunk() {
			g.write("Você afundou toda a frota em " + shotCount(g.shots) + "!\n")
			return g.err
		}
	}
}

// playerAttack holds the player's side of a game: the board they shoot at
// and how many valid shots they made. Writes stop at the first error, which
// is kept in err.
type playerAttack struct {
	out    io.Writer
	target *board.Board
	shots  int
	err    error
}

// prompt writes text and returns the line the player typed. It returns
// false when the player gives up, with "sair" or the end of the input, or
// when reading or writing fails.
func (g *playerAttack) prompt(scanner *bufio.Scanner, text string) (string, bool) {
	g.write(text)
	if g.err != nil {
		return "", false
	}
	if !scanner.Scan() {
		if g.err = scanner.Err(); g.err == nil {
			// The terminal does not echo a newline for Ctrl+D.
			g.write("\n")
		}
		return "", false
	}
	line := scanner.Text()
	if strings.EqualFold(strings.TrimSpace(line), quitCommand) {
		return "", false
	}
	return line, true
}

// shoot fires at the cell in line and shows the result followed by the
// attack board. Invalid input only shows a message, does not count as a shot
// and returns false.
func (g *playerAttack) shoot(line string) (board.Shot, bool) {
	c, err := board.ParseCoord(line)
	if err != nil {
		g.reject(err, c)
		return board.Shot{}, false
	}
	shot, err := g.target.Fire(c)
	if err != nil {
		g.reject(err, c)
		return board.Shot{}, false
	}
	g.shots++
	g.write(ShotMessage(shot) + "\n\n")
	g.showBoard()
	return shot, g.err == nil
}

// giveUp ends the game with the number of shots and reveal, then the target
// board with its ships and the shots it received.
func (g *playerAttack) giveUp(reveal string) {
	g.write("Você desistiu depois de " + shotCount(g.shots) + ". " + reveal + "\n\n")
	g.render(RenderBoard, g.target)
}

// shotCount writes n shots in Portuguese, such as "1 disparo" or "2 disparos".
func shotCount(n int) string {
	if n == 1 {
		return "1 disparo"
	}
	return fmt.Sprintf("%d disparos", n)
}

// reject shows the message for an invalid shot at c. An error with no
// message for the player is kept in g.err and ends the game.
func (g *playerAttack) reject(err error, c board.Coord) {
	message, ok := shotErrorMessage(err, c)
	if !ok {
		g.err = err
		return
	}
	g.write(message + "\n")
}

// shotErrorMessage returns the Portuguese message for a shot at c that was
// refused with err, and false if err is not a mistake of the player.
func shotErrorMessage(err error, c board.Coord) (string, bool) {
	switch {
	case errors.Is(err, board.ErrInvalidFormat):
		return "Coordenada inválida: use uma letra de A a J e um número de 1 a 10, por exemplo B7.", true
	case errors.Is(err, board.ErrOutOfBounds):
		return "Fora do tabuleiro: use uma letra de A a J e um número de 1 a 10.", true
	case errors.Is(err, board.ErrAlreadyShot):
		return "Você já atirou em " + c.String() + ".", true
	default:
		return "", false
	}
}

// showBoard writes the attack board followed by a blank line.
func (g *playerAttack) showBoard() {
	g.render(RenderAttackBoard, g.target)
	g.write("\n")
}

// render draws b with draw, unless an earlier write failed.
func (g *playerAttack) render(draw func(io.Writer, *board.Board, Style) error, b *board.Board) {
	if g.err == nil {
		g.err = draw(g.out, b, Plain)
	}
}

func (g *playerAttack) write(s string) {
	if g.err != nil {
		return
	}
	_, g.err = io.WriteString(g.out, s)
}
