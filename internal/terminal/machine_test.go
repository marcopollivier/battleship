package terminal

import (
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
)

// machineBoard returns a complete board whose fleet is exactly the placed ships.
func machineBoard(t *testing.T, placements ...placement) *board.Board {
	t.Helper()
	counts := map[board.ShipType]int{}
	for _, p := range placements {
		counts[p.shipType]++
	}
	var fleet []board.FleetEntry
	for shipType, count := range counts {
		fleet = append(fleet, board.FleetEntry{Type: shipType, Count: count})
	}
	b := board.NewBoard(fleet)
	for _, p := range placements {
		ref, err := board.ParseCoord(p.ref)
		if err != nil {
			t.Fatalf("ParseCoord(%q): %v", p.ref, err)
		}
		if err := b.Place(p.shipType, ref, p.orientation); err != nil {
			t.Fatalf("Place(%v, %s, %v): %v", p.shipType, p.ref, p.orientation, err)
		}
	}
	return b
}

// machinePlayerBoard has a single Submarine at A1.
func machinePlayerBoard(t *testing.T) *board.Board {
	t.Helper()
	return machineBoard(t, placement{board.Submarine, "A1", board.Horizontal})
}

// machineEnemyBoard has a Destroyer at E5-F5 and a Submarine at J10.
func machineEnemyBoard(t *testing.T) *board.Board {
	t.Helper()
	return machineBoard(t,
		placement{board.Destroyer, "E5", board.Horizontal},
		placement{board.Submarine, "J10", board.Horizontal},
	)
}

// machinePlay runs PlayVsMachine on fresh boards with input and the machine
// shooting with seed, and returns what it wrote.
func machinePlay(t *testing.T, input string, seed uint64) string {
	t.Helper()
	var sb strings.Builder
	err := PlayVsMachine(strings.NewReader(input), &sb, machinePlayerBoard(t), machineEnemyBoard(t), rand.New(rand.NewPCG(seed, 0)))
	if err != nil {
		t.Fatalf("PlayVsMachine: %v", err)
	}
	return sb.String()
}

// machineScreen is the full screen shown when the turn is the player's.
func machineScreen(t *testing.T, player, enemy *board.Board) string {
	t.Helper()
	return "Sua frota:\n" + soloReveal(t, player) + "\nFrota inimiga:\n" + attackRender(t, enemy) + "\n" + soloPrompt
}

// CA-1: the game starts with both boards and the prompt.
func TestPlayVsMachine_Start(t *testing.T) {
	got := machinePlay(t, "", 1)
	if want := machineScreen(t, machinePlayerBoard(t), machineEnemyBoard(t)); !strings.HasPrefix(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-2 and CA-5: a hit lets the player shoot again, without the machine
// shooting, and sinking the last enemy ship wins.
func TestPlayVsMachine_PlayerWins(t *testing.T) {
	player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
	want := machineScreen(t, player, enemy) +
		soloTurn(t, enemy, "E5", "Acerto!") + "\n" + soloPrompt +
		soloTurn(t, enemy, "F5", "Afundou o Destroyer!") + "\n" + soloPrompt +
		soloTurn(t, enemy, "J10", "Afundou o Submarino!") + "\n" +
		"Você venceu! Afundou a frota inimiga em 3 disparos.\n"

	if got := machinePlay(t, "E5\nF5\nJ10\nA1\n", 1); got != want {
		t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, want)
	}
}
