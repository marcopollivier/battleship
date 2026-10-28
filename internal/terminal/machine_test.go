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
	err := PlayVsMachine(strings.NewReader(input), &sb, machinePlayerBoard(t), machineEnemyBoard(t), rand.New(rand.NewPCG(seed, 0)), Display{Style: style}, board.RandomShot)
	if err != nil {
		t.Fatalf("PlayVsMachine: %v", err)
	}
	return sb.String()
}

// machineScreen is the screen shown when the turn is the player's, without
// colors and without clearing the terminal: both boards side by side, the
// legends, the events of the turn, one per line, and the prompt.
func machineScreen(t *testing.T, player, enemy *board.Board, events ...string) string {
	t.Helper()
	return machineScreenStyle(t, Plain, player, enemy, events...)
}

// machineScreenStyle is machineScreen with the boards drawn with style.
func machineScreenStyle(t *testing.T, style Style, player, enemy *board.Board, events ...string) string {
	t.Helper()
	lines := append([]string{"Sua frota                  Frota inimiga"}, sideBySide(ownGrid(player, style), attackGrid(enemy, style))...)
	lines = append(lines, "",
		"Sua frota: ~ água, S D C E P navio, X atingido, # afundado, o tiro na água",
		"Frota inimiga: . não atingida, o água, X acerto, S D C E P navio afundado",
		"")
	lines = append(lines, events...)
	return strings.Join(lines, "\n") + "\n" + soloPrompt
}

// machineFinal is the last screen of a game against the machine, without
// colors and without clearing the terminal: both boards side by side with
// their ships, the legend, the events of the turn and the final message.
func machineFinal(t *testing.T, player, enemy *board.Board, message string, events ...string) string {
	t.Helper()
	return machineFinalStyle(t, Plain, player, enemy, message, events...)
}

// machineFinalStyle is machineFinal with the boards drawn with style.
func machineFinalStyle(t *testing.T, style Style, player, enemy *board.Board, message string, events ...string) string {
	t.Helper()
	lines := append([]string{"Sua frota                  Frota inimiga"}, sideBySide(ownGrid(player, style), ownGrid(enemy, style))...)
	lines = append(lines, "", "Legenda: ~ água, S D C E P navio, X atingido, # afundado, o tiro na água", "")
	lines = append(lines, events...)
	return strings.Join(lines, "\n") + "\n" + message + "\n"
}

// machineShot fires at cell on replay and returns the event PlayVsMachine
// shows for it, such as "Seu disparo em E5: Acerto!".
func machineShot(t *testing.T, replay *board.Board, cell string) string {
	t.Helper()
	c, err := board.ParseCoord(cell)
	if err != nil {
		t.Fatalf("ParseCoord(%q): %v", cell, err)
	}
	shot, err := replay.Fire(c)
	if err != nil {
		t.Fatalf("Fire(%s): %v", cell, err)
	}
	return "Seu disparo em " + c.String() + ": " + ShotMessage(shot)
}

// CA-1 (spec 022): the game starts with both boards side by side, their
// titles and legends, and the prompt, no line wider than 80 columns.
func TestPlayVsMachine_Start(t *testing.T) {
	got := machinePlay(t, "", 1)
	want := machineScreen(t, machinePlayerBoard(t), machineEnemyBoard(t))
	if !strings.HasPrefix(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant prefix\n%s", got, want)
	}
	for i, line := range strings.Split(want, "\n") {
		if n := visibleWidth(line); n > maxWidth {
			t.Errorf("line %d is %d columns wide: %q", i+1, n, line)
		}
	}
}

// CA-2 (spec 022): after a hit, the player shoots again on a new screen
// with both boards and the result of the shot, citing the cell.
func TestPlayVsMachine_HitShootsAgain(t *testing.T) {
	player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
	want := machineScreen(t, player, enemy)
	want += machineScreen(t, player, enemy, machineShot(t, enemy, "E5"))
	if !strings.Contains(want, "\nSeu disparo em E5: Acerto!\nDisparo: ") {
		t.Fatalf("the screen does not show the hit:\n%s", want)
	}
	if got := machinePlay(t, "E5\nsair\n", 1); !strings.HasPrefix(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-2 and CA-5 (spec 008), CA-6 (spec 022): a hit lets the player shoot
// again, without the machine shooting, and sinking the last enemy ship wins,
// with a last screen showing both fleets.
func TestPlayVsMachine_PlayerWins(t *testing.T) {
	player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
	want := machineScreen(t, player, enemy)
	want += machineScreen(t, player, enemy, machineShot(t, enemy, "E5"))
	want += machineScreen(t, player, enemy, machineShot(t, enemy, "F5"))
	last := machineShot(t, enemy, "J10")
	want += machineFinal(t, player, enemy, "Você venceu! Afundou a frota inimiga em 3 disparos.", last)

	if got := machinePlay(t, "E5\nF5\nJ10\nA1\n", 1); got != want {
		t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, want)
	}
}

// CA-4 (spec 022): invalid input shows its message below the prompt and
// asks again, without a new screen.
func TestPlayVsMachine_InvalidInput(t *testing.T) {
	player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
	want := machineScreen(t, player, enemy)
	want += machineScreen(t, player, enemy, machineShot(t, enemy, "E5")) +
		"Coordenada inválida: use uma letra de A a J e um número de 1 a 10, por exemplo B7.\n" + soloPrompt +
		"Você já atirou em E5.\n" + soloPrompt +
		machineFinal(t, player, enemy, "Você desistiu depois de 1 disparo.")
	if got := machinePlay(t, "E5\n7B\ne5\nsair\n", 1); got != want {
		t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, want)
	}
}

// CA-5 (spec 022): with Clear, every screen starts with the code that clears
// the terminal, and only the screens; without it, the code never appears.
func TestPlayVsMachine_Clear(t *testing.T) {
	play := func(clear bool) string {
		var sb strings.Builder
		err := PlayVsMachine(strings.NewReader("E5\n7B\nA10\nsair\n"), &sb, machinePlayerBoard(t), machineEnemyBoard(t),
			rand.New(rand.NewPCG(1, 0)), Display{Style: Plain, Clear: clear}, machineScript(t, "J10"))
		if err != nil {
			t.Fatalf("PlayVsMachine: %v", err)
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
	if len(screens) != 5 || screens[0] != "" {
		t.Fatalf("with Clear, %d parts split by the clear code, want the output to start with it and 4 screens:\n%q", len(screens), cleared)
	}
	for i, screen := range screens[1:] {
		if !strings.HasPrefix(screen, "Sua frota                  Frota inimiga\n") {
			t.Errorf("screen %d does not start with the titles:\n%s", i+1, screen)
		}
	}
}

// machineTurn replays a machine turn on player with rng: it fires until a
// miss or until the player's fleet is sunk, and returns the lines
// PlayVsMachine should show and whether the machine won.
func machineTurn(t *testing.T, player *board.Board, rng *rand.Rand) ([]string, bool) {
	t.Helper()
	var lines []string
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
		lines = append(lines, "A máquina atirou em "+c.String()+": "+message)
		if player.AllSunk() {
			return lines, true
		}
		if shot.Result == board.Miss {
			return lines, false
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

// CA-3 (spec 022): after the player's miss the machine shoots until it
// misses, and the next screen shows the player's result and one line per
// machine shot below the boards, then the prompt.
func TestPlayVsMachine_MachineTurn(t *testing.T) {
	const seed = 1
	player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
	rng := rand.New(rand.NewPCG(seed, 0))
	want := machineScreen(t, player, enemy)
	events := []string{machineShot(t, enemy, "A10")}
	lines, won := machineTurn(t, player, rng)
	if won {
		t.Fatalf("the machine wins in its first turn with seed %d; pick another seed", seed)
	}
	want += machineScreen(t, player, enemy, append(events, lines...)...)

	if got := machinePlay(t, "A10\n", seed); !strings.HasPrefix(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant prefix\n%s", got, want)
	}
}

// CA-6 (spec 008 and spec 022): when the player only hits water, the machine
// eventually sinks the player's fleet and wins, and the last screen shows
// both boards with the enemy fleet revealed.
func TestPlayVsMachine_MachineWins(t *testing.T) {
	const seed = 3
	player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
	rng := rand.New(rand.NewPCG(seed, 0))
	waters := machineWaterCells(enemy)

	want := machineScreen(t, player, enemy)
	won, machineShots := false, 0
	for _, cell := range waters {
		events := []string{machineShot(t, enemy, cell)}
		lines, machineWon := machineTurn(t, player, rng)
		events = append(events, lines...)
		machineShots += len(lines)
		if machineWon {
			won = true
			want += machineFinal(t, player, enemy, "A máquina venceu! Ela afundou a sua frota em "+shotCount(machineShots)+".", events...)
			break
		}
		want += machineScreen(t, player, enemy, events...)
	}
	if !won {
		t.Fatalf("the machine does not win with seed %d before the player runs out of water", seed)
	}

	if got := machinePlay(t, strings.Join(waters, "\n")+"\n", seed); got != want {
		t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, want)
	}
}

// CA-7 (spec 008) and CA-6 (spec 022): "sair" or the end of the input gives
// up, and the last screen shows both boards with the enemy fleet revealed.
func TestPlayVsMachine_GiveUp(t *testing.T) {
	t.Run("sair at the start", func(t *testing.T) {
		player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
		want := machineScreen(t, player, enemy) +
			machineFinal(t, player, enemy, "Você desistiu depois de 0 disparos.")
		if got := machinePlay(t, "Sair\n", 1); got != want {
			t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, want)
		}
	})

	t.Run("end of input after a miss", func(t *testing.T) {
		const seed = 1
		player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
		want := machineScreen(t, player, enemy)
		events := []string{machineShot(t, enemy, "A10")}
		lines, won := machineTurn(t, player, rand.New(rand.NewPCG(seed, 0)))
		if won {
			t.Fatalf("the machine wins in its first turn with seed %d; pick another seed", seed)
		}
		want += machineScreen(t, player, enemy, append(events, lines...)...) +
			"\n" + machineFinal(t, player, enemy, "Você desistiu depois de 1 disparo.")
		if got := machinePlay(t, "A10", seed); got != want {
			t.Errorf("PlayVsMachine() =\n%s\nwant\n%s", got, want)
		}
	})
}

// machineScript is a ShotStrategy that shoots at cells in order, and
// returns ErrNoCellsLeft when they run out.
func machineScript(t *testing.T, cells ...string) board.ShotStrategy {
	t.Helper()
	var coords []board.Coord
	for _, cell := range cells {
		c, err := board.ParseCoord(cell)
		if err != nil {
			t.Fatalf("ParseCoord(%q): %v", cell, err)
		}
		coords = append(coords, c)
	}
	return func(*board.Board, *rand.Rand) (board.Coord, error) {
		if len(coords) == 0 {
			return board.Coord{}, board.ErrNoCellsLeft
		}
		c := coords[0]
		coords = coords[1:]
		return c, nil
	}
}

// Spec 013: the machine shoots where the ShotStrategy it receives chooses.
func TestPlayVsMachine_Strategy(t *testing.T) {
	var sb strings.Builder
	err := PlayVsMachine(strings.NewReader("A10\nA9\n"), &sb, machinePlayerBoard(t), machineEnemyBoard(t),
		rand.New(rand.NewPCG(1, 0)), Display{Style: Plain}, machineScript(t, "B5", "A1"))
	if err != nil {
		t.Fatalf("PlayVsMachine: %v", err)
	}
	got := sb.String()
	for _, want := range []string{
		"A máquina atirou em B5: Água.\n",
		"A máquina atirou em A1: Afundou o seu Submarino!\n",
		"A máquina venceu! Ela afundou a sua frota em 2 disparos.",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("PlayVsMachine() =\n%s\nwant it to contain %q", got, want)
		}
	}
}

// CA-4 (spec 023): when the machine sinks a ship of the player and the turn
// comes back, the sunk ship is drawn as # under "Sua frota".
func TestPlayVsMachine_SunkOwnShip(t *testing.T) {
	player := machineBoard(t,
		placement{board.Submarine, "A1", board.Horizontal},
		placement{board.Submarine, "J10", board.Horizontal},
	)
	var sb strings.Builder
	err := PlayVsMachine(strings.NewReader("A10\n"), &sb, player, machineEnemyBoard(t),
		rand.New(rand.NewPCG(1, 0)), Display{Style: Plain}, machineScript(t, "A1", "C3"))
	if err != nil {
		t.Fatalf("PlayVsMachine: %v", err)
	}
	before, _, found := strings.Cut(sb.String(), "A máquina atirou em C3: Água.\n")
	if !found {
		t.Fatalf("PlayVsMachine() =\n%s\nwant the screen after the machine's miss at C3", sb.String())
	}
	title := "Sua frota                  Frota inimiga\n"
	screen := before[strings.LastIndex(before, title)+len(title):]
	lines := strings.Split(screen, "\n")
	for row, want := range map[int]string{
		1:  " 1 # ~ ~ ~ ~ ~ ~ ~ ~ ~",
		3:  " 3 ~ ~ o ~ ~ ~ ~ ~ ~ ~",
		10: "10 ~ ~ ~ ~ ~ ~ ~ ~ ~ S",
	} {
		if got := lines[row][:gridWidth]; got != want {
			t.Errorf("Sua frota row %d = %q, want %q", row, got, want)
		}
	}
}

// Errors that are not a mistake of the player end the game and are returned.
func TestPlayVsMachine_Errors(t *testing.T) {
	rng := func() *rand.Rand { return rand.New(rand.NewPCG(1, 0)) }

	t.Run("player shoots at an incomplete enemy board", func(t *testing.T) {
		incomplete := board.NewBoard(board.DefaultFleet())
		err := PlayVsMachine(strings.NewReader("A10\n"), io.Discard, machinePlayerBoard(t), incomplete, rng(), Display{Style: Plain}, board.RandomShot)
		if !errors.Is(err, board.ErrIncompleteBoard) {
			t.Errorf("PlayVsMachine() error = %v, want ErrIncompleteBoard", err)
		}
	})

	t.Run("write fails", func(t *testing.T) {
		err := PlayVsMachine(strings.NewReader("A10\n"), &soloFailingWriter{}, machinePlayerBoard(t), machineEnemyBoard(t), rng(), Display{Style: Plain}, board.RandomShot)
		if !errors.Is(err, errSoloWrite) {
			t.Errorf("PlayVsMachine() error = %v, want %v", err, errSoloWrite)
		}
	})

	t.Run("the strategy fails", func(t *testing.T) {
		err := PlayVsMachine(strings.NewReader("A10\n"), io.Discard, machinePlayerBoard(t), machineEnemyBoard(t), rng(), Display{Style: Plain}, machineScript(t))
		if !errors.Is(err, board.ErrNoCellsLeft) {
			t.Errorf("PlayVsMachine() error = %v, want ErrNoCellsLeft", err)
		}
	})
}

// CA-7 (spec 010): with Colored, both boards of the screen, the boards at
// the end and the revealed enemy fleet are drawn with colors.
func TestPlayVsMachine_Colored(t *testing.T) {
	player, enemy := machinePlayerBoard(t), machineEnemyBoard(t)
	want := machineScreenStyle(t, Colored, player, enemy)
	attackFire(t, enemy, "E5")
	want += machineScreenStyle(t, Colored, player, enemy, "Seu disparo em E5: "+colored("1;31", "Acerto!")) +
		machineFinalStyle(t, Colored, player, enemy, "Você desistiu depois de 1 disparo.")

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
	assertUncolored(t, got, "Sua frota                  Frota inimiga", "Sua frota: ~", "Frota inimiga: .", soloPrompt, "Seu disparo em ", "A máquina atirou em ", "Legenda: ~")
}

// machineScriptPlay runs PlayVsMachine with colors on player and the
// default enemy board, with input and the machine shooting at cells in
// order, and returns what it wrote.
func machineScriptPlay(t *testing.T, input string, player *board.Board, cells ...string) string {
	t.Helper()
	var sb strings.Builder
	err := PlayVsMachine(strings.NewReader(input), &sb, player, machineEnemyBoard(t),
		rand.New(rand.NewPCG(1, 0)), Display{Style: Colored}, machineScript(t, cells...))
	if err != nil {
		t.Fatalf("PlayVsMachine: %v", err)
	}
	return sb.String()
}

// CA-2 (spec 015): with Colored, only the result after "A máquina atirou
// em A1: " has a color, the color of the shot result.
func TestPlayVsMachine_ColoredMachineShots(t *testing.T) {
	player := machineBoard(t,
		placement{board.Destroyer, "A1", board.Horizontal},
		placement{board.Submarine, "J10", board.Horizontal},
	)
	got := machineScriptPlay(t, "A10\n", player, "A1", "B1", "C5")
	want := "A máquina atirou em A1: " + colored("1;31", "Acerto!") + "\n" +
		"A máquina atirou em B1: " + colored("1;31", "Afundou o seu Destroyer!") + "\n" +
		"A máquina atirou em C5: " + colored("36", "Água.") + "\n"
	if !strings.Contains(got, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant it to contain\n%q", got, want)
	}
}

// CA-3 (spec 015): with Colored, the player's win is bold green, and the
// machine's win bold red.
func TestPlayVsMachine_ColoredEnd(t *testing.T) {
	won := machineScriptPlay(t, "E5\nF5\nJ10\n", machinePlayerBoard(t))
	if want := "\n" + colored("1;32", "Você venceu! Afundou a frota inimiga em 3 disparos.") + "\n"; !strings.HasSuffix(won, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant suffix %q", won, want)
	}

	lost := machineScriptPlay(t, "A10\n", machinePlayerBoard(t), "A1")
	want := "\n" + colored("1;31", "A máquina venceu! Ela afundou a sua frota em 1 disparo.") + "\n"
	if !strings.HasSuffix(lost, want) {
		t.Errorf("PlayVsMachine() =\n%s\nwant suffix %q", lost, want)
	}
}
