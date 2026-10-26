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
		shot, ok := g.shoot(line)
		if !ok {
			if g.err != nil {
				return g.err
			}
			continue
		}
		if enemy.AllSunk() {
			g.write("Você venceu! Afundou a frota inimiga em " + shotCount(g.shots) + ".\n")
			return g.err
		}
		if shot.Result != board.Miss {
			continue
		}
		if g.machineTurn() {
			g.machineWins()
			return g.err
		}
		g.write("\n")
		g.showScreen()
	}
}

// machineGame is a game against the machine: the player's attack on the
// enemy board, plus the player's own board that the machine shoots at.
type machineGame struct {
	playerAttack
	player       *board.Board
	rng          *rand.Rand
	machineShots int
}

// machineTurn makes the machine shoot at the player's board until it misses,
// one line per shot, and reports whether it sank the player's whole fleet.
func (g *machineGame) machineTurn() bool {
	for g.err == nil {
		c, err := board.RandomShot(g.player, g.rng)
		if err != nil {
			g.err = err
			return false
		}
		shot, err := g.player.Fire(c)
		if err != nil {
			g.err = err
			return false
		}
		g.machineShots++
		g.write("A máquina atirou em " + c.String() + ": " + machineShotMessage(shot) + "\n")
		if g.player.AllSunk() {
			return true
		}
		if shot.Result == board.Miss {
			return false
		}
	}
	return false
}

// machineWins ends the game showing the player's sunk fleet and revealing
// the enemy fleet.
func (g *machineGame) machineWins() {
	g.write("\nSua frota:\n")
	g.render(RenderBoard, g.player)
	g.write("\nA máquina venceu! Ela afundou a sua frota em " + shotCount(g.machineShots) + ". A frota inimiga era:\n\n")
	g.render(RenderBoard, g.target)
}

// machineShotMessage is the result of a machine shot, naming the player's
// ship when it sinks, such as "Afundou o seu Cruzador!".
func machineShotMessage(s board.Shot) string {
	if s.Result == board.Sunk {
		return "Afundou o seu " + ShipName(s.Ship) + "!"
	}
	return ShotMessage(s)
}

// showScreen writes the player's board and the attack board, each with its title.
func (g *machineGame) showScreen() {
	g.write("Sua frota:\n")
	g.render(RenderBoard, g.player)
	g.write("\nFrota inimiga:\n")
	g.showBoard()
}
