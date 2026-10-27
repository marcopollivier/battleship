package terminal

import (
	"errors"
	"io"
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
	return machinePlayStyle(t, input, seed, Plain)
}

// machinePlayStyle is machinePlay drawing the boards with style.
func machinePlayStyle(t *testing.T, input string, seed uint64, style Style) string {
	t.Helper()
	var sb strings.Builder
	err := PlayVsMachine(strings.NewReader(input), &sb, machinePlayerBoard(t), machineEnemyBoard(t), rand.New(rand.NewPCG(seed, 0)), style)
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

// machineTurn replays a machine turn on player with rng: it fires until a
// miss or until the player's fleet is sunk, and returns the lines
// PlayVsMachine should write and whether the machine won.
func machineTurn(t *testing.T, player *board.Board, rng *rand.Rand) (string, bool) {
	t.Helper()
	var sb strings.Builder
	for {
		c, err := board.RandomShot(player, rng)
		if err != nil {
			t.Fatalf("RandomShot: %v", err)
		}
		shot, err := player.Fire(c)
		if err != nil {
			t.Fatalf("Fire(%v): %v", c, err)
		}
		message := ShotMessage(shot)
		if shot.Result == board.Sunk {
			message = "Afundou o seu " + ShipName(shot.Ship) + "!"
		}
		sb.WriteString("A máquina atirou em " + c.String() + ": " + message + "\n")
		if player.AllSunk() {
			return sb.String(), true
		}
		if shot.Result == board.Miss {
			return sb.String(), false
		}
	}
}

// machineWaterCells lists the cells of b with no ship, in board order.
func machineWaterCells(b *board.Board) []string {
	ships := map[board.Coord]bool{}
	for _, s := range b.Ships() {
		for _, c := range s.Cells {
			ships[c] = true
		}
	}
	var cells []string
	for row := range board.Size {
		for col := range board.Size {
			if c := (board.Coord{Col: col, Row: row}); !ships[c] {
				cells = append(cells, c.String())
			}
		}
	}
	return cells
}

// CA-3: after the player's miss the machine shoots until it misses, one line
// per shot, and then the full screen comes back.
func TestPlayVsMachine_MachineTurn(t *testing.T) {
	const seed = 1
	player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
	rng := rand.New(rand.NewPCG(seed, 0))
	want := machineScreen(t, player, enemy) + soloTurn(t, enemy, "A10", "Água.") + "\n"
	lines, won := machineTurn(t, player, rng)
	if won {
		t.Fatalf("the machine wins in its first turn with seed %d; pick another seed", seed)
	}
	want += lines + "\n" + machineScreen(t, player, enemy)

	if got := machinePlay(t, "A10\n", seed); !strings.HasPrefix(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-6: when the player only hits water, the machine eventually sinks the
// player's fleet, wins, and the enemy fleet is revealed.
func TestPlayVsMachine_MachineWins(t *testing.T) {
	const seed = 3
	player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
	rng := rand.New(rand.NewPCG(seed, 0))
	waters := machineWaterCells(enemy)

	want := machineScreen(t, player, enemy)
	won, machineShots := false, 0
	for _, cell := range waters {
		want += soloTurn(t, enemy, cell, "Água.") + "\n"
		lines, machineWon := machineTurn(t, player, rng)
		want += lines
		machineShots += strings.Count(lines, "\n")
		if machineWon {
			won = true
			want += "\nSua frota:\n" + soloReveal(t, player) +
				"\nA máquina venceu! Ela afundou a sua frota em " + shotCount(machineShots) + ". A frota inimiga era:\n\n" +
				soloReveal(t, enemy)
			break
		}
		want += "\n" + machineScreen(t, player, enemy)
	}
	if !won {
		t.Fatalf("the machine does not win with seed %d before the player runs out of water", seed)
	}

	if got := machinePlay(t, strings.Join(waters, "\n")+"\n", seed); got != want {
		t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, want)
	}
}

// CA-7: "sair" or the end of the input gives up and reveals the enemy fleet.
func TestPlayVsMachine_GiveUp(t *testing.T) {
	t.Run("sair at the start", func(t *testing.T) {
		player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
		want := machineScreen(t, player, enemy) +
			"Você desistiu depois de 0 disparos. A frota inimiga era:\n\n" + soloReveal(t, enemy)
		if got := machinePlay(t, "Sair\n", 1); got != want {
			t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, want)
		}
	})

	t.Run("end of input after a miss", func(t *testing.T) {
		const seed = 1
		player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
		want := machineScreen(t, player, enemy) + soloTurn(t, enemy, "A10", "Água.") + "\n"
		lines, won := machineTurn(t, player, rand.New(rand.NewPCG(seed, 0)))
		if won {
			t.Fatalf("the machine wins in its first turn with seed %d; pick another seed", seed)
		}
		want += lines + "\n" + machineScreen(t, player, enemy) +
			"\nVocê desistiu depois de 1 disparo. A frota inimiga era:\n\n" + soloReveal(t, enemy)
		if got := machinePlay(t, "A10", seed); got != want {
			t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, want)
		}
	})
}

// Errors that are not a mistake of the player end the game and are returned.
func TestPlayVsMachine_Errors(t *testing.T) {
	rng := func() *rand.Rand { return rand.New(rand.NewPCG(1, 0)) }

	t.Run("player shoots at an incomplete enemy board", func(t *testing.T) {
		incomplete := board.NewBoard(board.DefaultFleet())
		err := PlayVsMachine(strings.NewReader("A10\n"), io.Discard, machinePlayerBoard(t), incomplete, rng(), Plain)
		if !errors.Is(err, board.ErrIncompleteBoard) {
			t.Errorf("PlayVsMachine() error = %v, want ErrIncompleteBoard", err)
		}
	})

	t.Run("write fails", func(t *testing.T) {
		err := PlayVsMachine(strings.NewReader("A10\n"), &soloFailingWriter{}, machinePlayerBoard(t), machineEnemyBoard(t), rng(), Plain)
		if !errors.Is(err, errSoloWrite) {
			t.Errorf("PlayVsMachine() error = %v, want %v", err, errSoloWrite)
		}
	})
}

// CA-7 (spec 010): with Colored, both boards of the screen, the boards at
// the end and the revealed enemy fleet are drawn with colors.
func TestPlayVsMachine_Colored(t *testing.T) {
	player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
	want := "Sua frota:\n" + renderColored(t, RenderBoard, player) +
		"\nFrota inimiga:\n" + renderColored(t, RenderAttackBoard, enemy) + "\n" + soloPrompt
	attackFire(t, enemy, "E5")
	want += "Acerto!\n\n" + renderColored(t, RenderAttackBoard, enemy) + "\n" + soloPrompt +
		"Você desistiu depois de 1 disparo. A frota inimiga era:\n\n" + renderColored(t, RenderBoard, enemy)

	if got := machinePlayStyle(t, "E5\nsair\n", 1, Colored); got != want {
		t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, want)
	}

	const seed = 3 // the machine wins, as in TestPlayVsMachine_MachineWins
	input := strings.Join(machineWaterCells(machineEnemyBoard(t)), "\n") + "\n"
	got := machinePlayStyle(t, input, seed, Colored)
	if !strings.Contains(got, "A máquina venceu!") {
		t.Fatalf("the machine does not win with seed %d", seed)
	}
	assertBoardsColored(t, got, machinePlay(t, input, seed))
}
