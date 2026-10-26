package terminal

import (
	"bufio"
	"io"
	"math/rand/v2"

	"github.com/marcopollivier/battleship/internal/board"
)

// PlayVsMachine runs a game against the machine: the player shoots at enemy
// with shots read from in, the machine shoots at player with cells drawn
// from rng, and whoever hits shoots again. It writes the game to out and
// returns an error if reading or writing fails.
func PlayVsMachine(in io.Reader, out io.Writer, player, enemy *board.Board, rng *rand.Rand) error {
	g := &machineGame{playerAttack: playerAttack{out: out, target: enemy}, player: player, rng: rng}
	g.showScreen()
	scanner := bufio.NewScanner(in)
	for {
		line, ok := g.prompt(scanner)
		if !ok {
			return g.err
		}
		if _, ok := g.shoot(line); !ok {
			if g.err != nil {
				return g.err
			}
			continue
		}
		if enemy.AllSunk() {
			g.write("Você venceu! Afundou a frota inimiga em " + shotCount(g.shots) + ".\n")
			return g.err
		}
	}
}

// machineGame is a game against the machine: the player's attack on the
// enemy board, plus the player's own board that the machine shoots at.
type machineGame struct {
	playerAttack
	player *board.Board
	rng    *rand.Rand
}

// showScreen writes the player's board and the attack board, each with its title.
func (g *machineGame) showScreen() {
	g.write("Sua frota:\n")
	g.render(RenderBoard, g.player)
	g.write("\nFrota inimiga:\n")
	g.showBoard()
}
