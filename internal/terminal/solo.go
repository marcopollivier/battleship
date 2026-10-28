package terminal

import (
	"bufio"
	"errors"
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
// writes a screen with the attack board and the result of the last shot to
// out after each shot, until every ship of b is sunk or the player gives up,
// with "sair" or the end of the input. The last screen reveals b. The game
// is shown as display says. It returns an error if writing to out or reading
// from in fails.
func PlaySolo(in io.Reader, out io.Writer, b *board.Board, display Display) error {
	g := &playerAttack{out: out, target: b, style: display.Style, clear: display.Clear}
	g.showBoard()
	scanner := bufio.NewScanner(in)
	for {
		line, ok := g.prompt(scanner, shotPrompt)
		if !ok {
			g.showReveal("Você desistiu depois de " + shotCount(g.shots) + ".")
			return g.err
		}
		c, shot, ok := g.shoot(line)
		if !ok {
			if g.err != nil {
				return g.err
			}
			continue
		}
		g.events = append(g.events, g.shotResult(c, shot))
		if b.AllSunk() {
			g.showReveal(g.style.paint("Você afundou toda a frota em "+shotCount(g.shots)+"!", winColor))
			return g.err
		}
		g.showBoard()
	}
}

// playerAttack holds the player's side of a game: the board they shoot at,
// how many valid shots they made, the style every board is drawn with and
// whether the terminal is cleared before each screen. Writes stop at the
// first error, which is kept in err.
type playerAttack struct {
	out    io.Writer
	target *board.Board
	style  Style
	clear  bool
	shots  int
	err    error
	// events are the lines of what happened since the last screen, shown
	// below the boards of the next one.
	events []string
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

// shoot fires at the cell in line and returns the cell and the shot.
// Invalid input only shows a message, does not count as a shot and returns
// false.
func (g *playerAttack) shoot(line string) (board.Coord, board.Shot, bool) {
	c, err := board.ParseCoord(line)
	if err != nil {
		g.reject(err, c)
		return c, board.Shot{}, false
	}
	shot, err := g.target.Fire(c)
	if err != nil {
		g.reject(err, c)
		return c, board.Shot{}, false
	}
	g.shots++
	return c, shot, true
}

// shotResult is the result of the player's shot at c, citing the cell, such
// as "Seu disparo em B7: Acerto!", with the result in its color.
func (g *playerAttack) shotResult(c board.Coord, shot board.Shot) string {
	return "Seu disparo em " + c.String() + ": " + g.style.paint(ShotMessage(shot), shotColor(shot.Result))
}

// showReveal writes the last screen of a solo game: the target board with
// its ships and the shots it received, the events since the last screen and
// message.
func (g *playerAttack) showReveal(message string) {
	g.clearScreen()
	g.write(strings.Join(ownGrid(g.target, g.style), "\n") + "\n\n" + finalLegend + "\n\n")
	g.writeEvents()
	g.write(message + "\n")
}

// shotCount writes n shots in Portuguese, such as "1 disparo" or "2 disparos".
func shotCount(n int) string {
	return plural(n, "disparo", "disparos")
}

// reject shows the message for an invalid shot at c, in yellow. An error with no
// message for the player is kept in g.err and ends the game.
func (g *playerAttack) reject(err error, c board.Coord) {
	message, ok := shotErrorMessage(err, c)
	if !ok {
		g.err = err
		return
	}
	g.write(g.style.paint(message, warnColor) + "\n")
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

// showBoard writes the screen of a solo game: the attack board, a blank
// line and the events since the last screen. The prompt comes after it.
func (g *playerAttack) showBoard() {
	g.clearScreen()
	g.render(RenderAttackBoard, g.target)
	g.write("\n")
	g.writeEvents()
}

// clearScreen clears the terminal before a screen, if the game clears it.
func (g *playerAttack) clearScreen() {
	if g.clear {
		g.write(clearScreen)
	}
}

// writeEvents writes the events since the last screen, one per line, and
// forgets them.
func (g *playerAttack) writeEvents() {
	for _, event := range g.events {
		g.write(event + "\n")
	}
	g.events = nil
}

// render draws b with draw in the game's style, unless an earlier write failed.
func (g *playerAttack) render(draw func(io.Writer, *board.Board, Style) error, b *board.Board) {
	if g.err == nil {
		g.err = draw(g.out, b, g.style)
	}
}

func (g *playerAttack) write(s string) {
	if g.err != nil {
		return
	}
	_, g.err = io.WriteString(g.out, s)
}
