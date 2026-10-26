package terminal

import (
	"bufio"
	"fmt"
	"io"

	"github.com/marcopollivier/battleship/internal/board"
)

// shotPrompt asks the player for the next shot.
const shotPrompt = "Disparo: "

// PlaySolo runs a solo game: it reads shots from in, one per line, and
// writes the attack board and messages to out until every ship of b is sunk
// or the input ends. It returns an error if writing to out or reading from in
// fails.
func PlaySolo(in io.Reader, out io.Writer, b *board.Board) error {
	g := &soloGame{out: out, board: b}
	g.showBoard()
	scanner := bufio.NewScanner(in)
	for {
		g.write(shotPrompt)
		if g.err != nil {
			return g.err
		}
		if !scanner.Scan() {
			return scanner.Err()
		}
		g.shoot(scanner.Text())
		if g.err != nil {
			return g.err
		}
		if g.board.AllSunk() {
			g.write(fmt.Sprintf("Você afundou toda a frota em %d disparos!\n", g.shots))
			return g.err
		}
	}
}

// soloGame holds the state of a solo game. Writes stop at the first error,
// which is kept in err.
type soloGame struct {
	out   io.Writer
	board *board.Board
	shots int
	err   error
}

// shoot fires at the cell in line and shows the result.
func (g *soloGame) shoot(line string) {
	c, err := board.ParseCoord(line)
	if err != nil {
		g.err = err
		return
	}
	shot, err := g.board.Fire(c)
	if err != nil {
		g.err = err
		return
	}
	g.shots++
	g.write(ShotMessage(shot) + "\n\n")
	g.showBoard()
}

// showBoard writes the attack board followed by a blank line.
func (g *soloGame) showBoard() {
	if g.err != nil {
		return
	}
	if g.err = RenderAttackBoard(g.out, g.board); g.err == nil {
		g.write("\n")
	}
}

func (g *soloGame) write(s string) {
	if g.err != nil {
		return
	}
	_, g.err = io.WriteString(g.out, s)
}
