package terminal

import (
	"bufio"
	"io"
	"math/rand/v2"
	"strings"

	"github.com/marcopollivier/battleship/internal/board"
)

// PlayVsMachine runs a game against the machine. If player is not complete,
// the player first places the missing ships. Then the player shoots at enemy
// with shots read from in, the machine shoots at player with the cells shoot
// chooses, drawing with rng, and whoever hits shoots again. The player gives
// up with "sair" or the end of the input. It writes the game to out, with the
// boards drawn as display says, and returns an error if reading, writing or
// choosing a machine shot fails.
func PlayVsMachine(in io.Reader, out io.Writer, player, enemy *board.Board, rng *rand.Rand, display Display, shoot board.ShotStrategy) error {
	g := &machineGame{playerAttack: playerAttack{out: out, target: enemy, style: display.Style, clear: display.Clear}, player: player, rng: rng, strategy: shoot}
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
		c, shot, ok := g.shoot(line)
		if !ok {
			if g.err != nil {
				return g.err
			}
			continue
		}
		g.events = append(g.events, g.shotResult(c, shot))
		if enemy.AllSunk() {
			g.writeEvents()
			g.write(g.style.paint("Você venceu! Afundou a frota inimiga em "+shotCount(g.shots)+".", winColor) + "\n")
			return g.err
		}
		if shot.Result == board.Miss && g.machineTurn() {
			g.writeEvents()
			g.machineWins()
			return g.err
		}
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
	// events are the lines of what happened since the last screen, shown
	// below the boards of the next one.
	events []string
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
		g.events = append(g.events, "A máquina atirou em "+c.String()+": "+g.style.paint(machineShotMessage(shot), shotColor(shot.Result)))
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
	g.write("\n" + g.style.paint("A máquina venceu! Ela afundou a sua frota em "+shotCount(g.machineShots)+".", hitColor) + " A frota inimiga era:\n\n")
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

// Legends of the boards on the game screen, which name the ship letters
// together to keep the screen within 80 columns.
const (
	ownScreenLegend   = ownTitle + ": ~ água, S D C E P navio, X atingido, o tiro na água"
	enemyScreenLegend = enemyTitle + ": . não atingida, o água, X acerto, S D C E P navio afundado"
)

// showScreen writes the game screen: the clear code if the terminal is
// cleared, both boards side by side with their titles and legends, and the
// events since the last screen. The prompt comes after it.
func (g *machineGame) showScreen() {
	var sb strings.Builder
	if g.clear {
		sb.WriteString(clearScreen)
	}
	sb.WriteString(boardTitles + "\n")
	for _, line := range sideBySide(ownGrid(g.player, g.style), attackGrid(g.target, g.style)) {
		sb.WriteString(line + "\n")
	}
	sb.WriteString("\n" + ownScreenLegend + "\n" + enemyScreenLegend + "\n\n")
	g.write(sb.String())
	g.writeEvents()
}

// writeEvents writes the events since the last screen, one per line, and
// forgets them.
func (g *machineGame) writeEvents() {
	for _, event := range g.events {
		g.write(event + "\n")
	}
	g.events = nil
}
