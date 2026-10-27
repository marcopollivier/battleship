package terminal

import (
	"bufio"
	"io"
	"math/rand/v2"

	"github.com/marcopollivier/battleship/internal/board"
)

// PlayVsMachine runs a game against the machine. If player is not complete,
// the player first places the missing ships. Then the player shoots at enemy
// with shots read from in, the machine shoots at player with the cells shoot
// chooses, drawing with rng, and whoever hits shoots again. The player gives
// up with "sair" or the end of the input. It writes the game to out and
// returns an error if reading, writing or choosing a machine shot fails.
func PlayVsMachine(in io.Reader, out io.Writer, player, enemy *board.Board, rng *rand.Rand, shoot board.ShotStrategy) error {
	g := &machineGame{playerAttack: playerAttack{out: out, target: enemy}, player: player, rng: rng, strategy: shoot}
	scanner := bufio.NewScanner(in)
	if !player.Complete() && !g.placeFleet(scanner) {
		return g.err
	}
	g.showScreen()
	for {
		line, ok := g.prompt(scanner, shotPrompt)
		if !ok {
			g.giveUp("A frota inimiga era:")
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
	strategy     board.ShotStrategy
	machineShots int
}

// machineTurn makes the machine shoot at the player's board until it misses,
// one line per shot, and reports whether it sank the player's whole fleet.
func (g *machineGame) machineTurn() bool {
	for g.err == nil {
		c, err := g.strategy(g.player, g.rng)
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
