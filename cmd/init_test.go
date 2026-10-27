package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
	"github.com/marcopollivier/battleship/internal/terminal"
)

// run runs the program with stdin and args and returns the exit code, stdout and stderr.
func run(stdin string, args ...string) (code int, stdout, stderr string) {
	var out, errOut strings.Builder
	code = Run(args, strings.NewReader(stdin), &out, &errOut)
	return code, out.String(), errOut.String()
}

// seedShots returns the cells of every ship on the board the program draws
// for fleet and seed, one per line, ready to be used as input.
func seedShots(t *testing.T, fleet []board.FleetEntry, seed uint64) string {
	t.Helper()
	b, err := board.RandomBoard(fleet, newRNG(seed))
	if err != nil {
		t.Fatalf("RandomBoard: %v", err)
	}
	var cells []string
	for _, ship := range b.Ships() {
		for _, c := range ship.Cells {
			cells = append(cells, c.String())
		}
	}
	return strings.Join(cells, "\n") + "\n"
}

// emptyAttackBoard is the attack board before any shot, drawn with style.
func emptyAttackBoard(t *testing.T, style terminal.Style) string {
	t.Helper()
	var sb strings.Builder
	if err := terminal.RenderAttackBoard(&sb, board.NewBoard(board.DefaultFleet()), style); err != nil {
		t.Fatalf("RenderAttackBoard: %v", err)
	}
	return sb.String()
}

// Solo mode (spec 006): shooting every ship cell wins, and the seed is shown at the end.
func TestRun_Win(t *testing.T) {
	code, stdout, stderr := run(seedShots(t, board.DefaultFleet(), 42), "--modo", "solo", "--semente", "42")
	if code != 0 {
		t.Errorf("Run() code = %d, want 0 (stderr %q)", code, stderr)
	}
	if want := "Você afundou toda a frota em 19 disparos!\nSemente: 42\n"; !strings.HasSuffix(stdout, want) {
		t.Errorf("Run() stdout =\n%s\nwant suffix\n%s", stdout, want)
	}
}

// Solo mode (spec 006): without --semente the game starts with the empty
// attack board, and giving up shows the fleet and the seed that was drawn.
func TestRun_GiveUp(t *testing.T) {
	code, stdout, stderr := run("sair\n", "--modo", "solo", "--cores", "nao")
	if code != 0 {
		t.Errorf("Run() code = %d, want 0 (stderr %q)", code, stderr)
	}
	if want := emptyAttackBoard(t, terminal.Plain) + "\nDisparo: "; !strings.HasPrefix(stdout, want) {
		t.Errorf("Run() stdout =\n%s\nwant prefix\n%s", stdout, want)
	}
	if !strings.Contains(stdout, "Você desistiu depois de 0 disparos. A frota era:") {
		t.Errorf("Run() stdout =\n%s\nwant the give up message", stdout)
	}
	if !regexp.MustCompile(`\nSemente: \d+\n$`).MatchString(stdout) {
		t.Errorf("Run() stdout =\n%s\nwant the seed on the last line", stdout)
	}
}

// CA-9 (spec 012), CA-7 (spec 013), CA-7 (spec 014) and CA-7 (spec 018): with the same seed,
// level and input the game repeats: the machine's board, the ships drawn
// with aleatorio and every machine shot, also after desfazer and the warning
// that the rest of the fleet no longer fits.
func TestRun_SameSeed(t *testing.T) {
	const (
		input = "B1 baixo\nE1 h\ndesfazer\nE3 v\ndesfazer\naleatorio\nA1\nB2\nC3\nsair\n"
		// Six battleships leave only row 10, and the seventh at D10 leaves
		// no room for the cruiser: the warning shows before desfazer.
		warningInput = "A2 h\nF2 h\nA5 h\nF5 h\nA8 h\nF8 h\nD10 h\ndesfazer\naleatorio\nA1\nB2\nC3\nsair\n"
		warning      = "Atenção: o resto da frota não cabe mais no tabuleiro." //nolint:misspell // Portuguese text for the player
	)
	tests := []struct {
		name  string
		input string
		args  []string
	}{
		{"default level", input, nil},
		{"dificil", input, []string{"--nivel", "dificil"}},
		{"facil", input, []string{"--nivel", "facil"}},
		{"mestre", input, []string{"--nivel", "mestre"}},
		{"warning", warningInput, []string{"--frota", "E=7,C=1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, first, _ := run(tt.input, append([]string{"--semente", "42"}, tt.args...)...)
			_, second, _ := run(tt.input, append([]string{"--semente", "42"}, tt.args...)...)
			if first != second {
				t.Errorf("two runs with --semente 42 %q differ:\n%s\n---\n%s", tt.args, first, second)
			}
			if _, other, _ := run(tt.input, append([]string{"--semente", "43"}, tt.args...)...); other == first {
				t.Error("runs with --semente 42 and 43 drew the same board")
			}
			if hasWarning := strings.Contains(first, warning); hasWarning != (tt.input == warningInput) {
				t.Errorf("warning shown = %v, want %v:\n%s", hasWarning, !hasWarning, first)
			}
		})
	}
}

// levelGame returns what the program writes for a game against the machine
// with the default fleet, seed and input, when the machine shoots with
// strategy and the boards drawn with colors, the default: the boards and the
// shots come from one random source, as in play.
func levelGame(t *testing.T, input string, seed uint64, strategy board.ShotStrategy) string {
	t.Helper()
	rng := newRNG(seed)
	enemy, err := board.RandomBoard(board.DefaultFleet(), rng)
	if err != nil {
		t.Fatalf("RandomBoard: %v", err)
	}
	var sb strings.Builder
	if err := terminal.PlayVsMachine(strings.NewReader(input), &sb, board.NewBoard(board.DefaultFleet()), enemy, rng, terminal.Colored, strategy); err != nil {
		t.Fatalf("PlayVsMachine: %v", err)
	}
	return sb.String() + fmt.Sprintf("Semente: %d\n", seed)
}

// CA-8 and the level option (spec 013): dificil, the default, plays with
// SmartShot, and facil with RandomShot, the machine from before spec 013.
// CA-7 (spec 018): mestre plays with MasterShot. The player only hits
// water, so the machine shoots until it wins.
func TestRun_Level(t *testing.T) {
	const seed = 42
	input := randomPlacement + machineCells(machineEnemy(t, board.DefaultFleet(), seed), false)
	smart, random := levelGame(t, input, seed, board.SmartShot), levelGame(t, input, seed, board.RandomShot)
	master := levelGame(t, input, seed, board.MasterShot)
	if smart == random || master == smart {
		t.Fatal("two strategies played the same game; pick another seed")
	}
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"default", nil, smart},
		{"dificil", []string{"--nivel", "dificil"}, smart},
		{"facil", []string{"--nivel", "facil"}, random},
		{"any case", []string{"--nivel", "FACIL"}, random},
		{"mestre", []string{"--nivel", "mestre"}, master},
		{"mestre in any case", []string{"--nivel", "MESTRE"}, master},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := run(input, append([]string{"--semente", "42"}, tt.args...)...)
			if code != 0 {
				t.Errorf("Run(%q) code = %d, want 0 (stderr %q)", tt.args, code, stderr)
			}
			if stdout != tt.want {
				t.Errorf("Run(%q) stdout =\n%s\nwant\n%s", tt.args, stdout, tt.want)
			}
		})
	}
}

// CA-9 (spec 013) and CA-8 (spec 018): an unknown level is an option error, with exit code 2,
// in solo mode too.
func TestRun_InvalidLevel(t *testing.T) {
	for _, args := range [][]string{{"--nivel", "xyz"}, {"--modo", "solo", "--nivel", "xyz"}} {
		code, stdout, stderr := run("sair\n", args...)
		if code != 2 {
			t.Errorf("Run(%q) code = %d, want 2", args, code)
		}
		if stdout != "" {
			t.Errorf("Run(%q) stdout = %q, want empty", args, stdout)
		}
		if want := "Nível inválido: use facil, dificil ou mestre.\n"; stderr != want {
			t.Errorf("Run(%q) stderr = %q, want %q", args, stderr, want)
		}
	}
}

// The game uses the fleet chosen with --frota; rejected fleets end the
// program before the game, with the reason on stderr and exit code 1.
func TestRun_Fleet(t *testing.T) {
	twoSubs := []board.FleetEntry{{Type: board.Submarine, Count: 2}}
	code, stdout, stderr := run(seedShots(t, twoSubs, 7), "--modo", "solo", "--frota", "S=2", "--semente", "7")
	if code != 0 {
		t.Errorf("Run(--frota S=2) code = %d, want 0 (stderr %q)", code, stderr)
	}
	if want := "Você afundou toda a frota em 2 disparos!\nSemente: 7\n"; !strings.HasSuffix(stdout, want) {
		t.Errorf("Run(--frota S=2) stdout =\n%s\nwant suffix\n%s", stdout, want)
	}

	tests := []struct {
		name       string
		args       []string
		wantStderr string
	}{
		{"does not fit", []string{"--frota", "E=30"}, "Não foi possível posicionar essa frota no tabuleiro.\n"},
		{"unknown letter", []string{"--frota", "X=1"}, "Frota inválida: use letra=quantidade separados por vírgula, por exemplo S=2,D=1,P=1.\n"},
		{"given but empty", []string{"--frota", ""}, "Frota inválida: use letra=quantidade separados por vírgula, por exemplo S=2,D=1,P=1.\n"},
		{"negative count", []string{"--frota", "S=-1"}, "Frota inválida: as quantidades não podem ser negativas.\n"},
		{"no ships", []string{"--frota", "S=0"}, "Frota inválida: escolha pelo menos um navio.\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := run("sair\n", tt.args...)
			if code != 1 {
				t.Errorf("Run(%q) code = %d, want 1", tt.args, code)
			}
			if stdout != "" {
				t.Errorf("Run(%q) stdout = %q, want empty", tt.args, stdout)
			}
			if stderr != tt.wantStderr {
				t.Errorf("Run(%q) stderr = %q, want %q", tt.args, stderr, tt.wantStderr)
			}
		})
	}
}

// Invalid --semente values and unknown options: the flag package rejects them with exit code 2
// and the game does not start.
func TestRun_InvalidOption(t *testing.T) {
	for _, args := range [][]string{{"--semente", "abc"}, {"--semente", "-1"}, {"--nope"}} {
		code, stdout, stderr := run("sair\n", args...)
		if code != 2 {
			t.Errorf("Run(%q) code = %d, want 2", args, code)
		}
		if stdout != "" {
			t.Errorf("Run(%q) stdout = %q, want empty", args, stdout)
		}
		if !strings.Contains(stderr, "-semente") {
			t.Errorf("Run(%q) stderr = %q, want the usage with -semente", args, stderr)
		}
	}
}

// machineEnemy returns the machine's board the program draws in machine
// mode for fleet and seed: it is the first board drawn.
func machineEnemy(t *testing.T, fleet []board.FleetEntry, seed uint64) *board.Board {
	t.Helper()
	enemy, err := board.RandomBoard(fleet, newRNG(seed))
	if err != nil {
		t.Fatalf("RandomBoard: %v", err)
	}
	return enemy
}

// machineCells lists the cells of b with or without a ship, one per line.
func machineCells(b *board.Board, withShip bool) string {
	ships := map[board.Coord]bool{}
	for _, s := range b.Ships() {
		for _, c := range s.Cells {
			ships[c] = true
		}
	}
	var cells []string
	for row := range board.Size {
		for col := range board.Size {
			if c := (board.Coord{Col: col, Row: row}); ships[c] == withShip {
				cells = append(cells, c.String())
			}
		}
	}
	return strings.Join(cells, "\n") + "\n"
}

// Without --modo the game is against the machine and starts by placing the
// fleet (spec 012 CA-1); --modo solo, in any case, is the solo game with no
// placing (spec 012 CA-10).
func TestRun_Mode(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantPrefix string
	}{
		{"default", nil, "Posicione a sua frota."},
		{"maquina", []string{"--modo", "MAQUINA"}, "Posicione a sua frota."},
		{"solo", []string{"--modo", "Solo", "--cores", "nao"}, emptyAttackBoard(t, terminal.Plain) + "\nDisparo: "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := run("sair\n", tt.args...)
			if code != 0 {
				t.Errorf("Run(%q) code = %d, want 0 (stderr %q)", tt.args, code, stderr)
			}
			if !strings.HasPrefix(stdout, tt.wantPrefix) {
				t.Errorf("Run(%q) stdout =\n%s\nwant prefix\n%s", tt.args, stdout, tt.wantPrefix)
			}
		})
	}
}

// CA-5: sinking the machine's fleet wins, with the seed at the end.
func TestRun_PlayerWins(t *testing.T) {
	oneSub := []board.FleetEntry{{Type: board.Submarine, Count: 1}}
	enemy := machineEnemy(t, oneSub, 5)
	code, stdout, stderr := run(randomPlacement+machineCells(enemy, true), "--frota", "S=1", "--semente", "5")
	if code != 0 {
		t.Errorf("Run() code = %d, want 0 (stderr %q)", code, stderr)
	}
	if want := "Você venceu! Afundou a frota inimiga em 1 disparo.\nSemente: 5\n"; !strings.HasSuffix(stdout, want) {
		t.Errorf("Run() stdout =\n%s\nwant suffix\n%s", stdout, want)
	}
}

// CA-6: when the player only hits water, the machine wins and the enemy
// fleet is revealed before the seed.
func TestRun_MachineWins(t *testing.T) {
	oneSub := []board.FleetEntry{{Type: board.Submarine, Count: 1}}
	enemy := machineEnemy(t, oneSub, 5)
	code, stdout, stderr := run(randomPlacement+machineCells(enemy, false), "--frota", "S=1", "--semente", "5")
	if code != 0 {
		t.Errorf("Run() code = %d, want 0 (stderr %q)", code, stderr)
	}
	if !strings.Contains(stdout, "A máquina venceu! Ela afundou a sua frota em ") {
		t.Errorf("Run() stdout =\n%s\nwant the machine to win", stdout)
	}
	if !strings.HasSuffix(stdout, "Semente: 5\n") {
		t.Errorf("Run() stdout does not end with the seed:\n%s", stdout)
	}
}

// CA-10: an unknown mode is an option error, with exit code 2.
func TestRun_InvalidMode(t *testing.T) {
	code, stdout, stderr := run("sair\n", "--modo", "xyz")
	if code != 2 {
		t.Errorf("Run(--modo xyz) code = %d, want 2", code)
	}
	if stdout != "" {
		t.Errorf("Run(--modo xyz) stdout = %q, want empty", stdout)
	}
	if want := "Modo inválido: use maquina ou solo.\n"; stderr != want {
		t.Errorf("Run(--modo xyz) stderr = %q, want %q", stderr, want)
	}
}

// CA-4 and CA-5 (spec 010): the boards are colored by default and with
// --cores sim in any case, even with the output going to a buffer instead of
// a terminal; --cores nao leaves no color code in the whole game.
func TestRun_Colors(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		style terminal.Style
	}{
		{"default", nil, terminal.Colored},
		{"sim", []string{"--cores", "sim"}, terminal.Colored},
		{"SIM", []string{"--cores", "SIM"}, terminal.Colored},
		{"nao", []string{"--cores", "nao"}, terminal.Plain},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := append([]string{"--modo", "solo", "--semente", "42"}, tt.args...)
			code, stdout, stderr := run("A1\nsair\n", args...)
			if code != 0 {
				t.Errorf("Run(%q) code = %d, want 0 (stderr %q)", args, code, stderr)
			}
			if want := emptyAttackBoard(t, tt.style) + "\nDisparo: "; !strings.HasPrefix(stdout, want) {
				t.Errorf("Run(%q) stdout =\n%s\nwant prefix\n%s", args, stdout, want)
			}
			if colored := strings.Contains(stdout, "\x1b"); colored != (tt.style == terminal.Colored) {
				t.Errorf("Run(%q) stdout has color codes = %v, want %v", args, colored, !colored)
			}
		})
	}

	// Against the machine, the same seed plays the same game with or without colors.
	const input = "aleatorio\nA1\nB2\nsair\n"
	_, plain, _ := run(input, "--semente", "42", "--cores", "nao")
	if strings.Contains(plain, "\x1b") {
		t.Errorf("Run(--cores nao) against the machine has color codes:\n%s", plain)
	}
	_, colored, _ := run(input, "--semente", "42")
	if got := regexp.MustCompile("\x1b\\[[0-9;]*m").ReplaceAllString(colored, ""); got != plain {
		t.Errorf("colored game without codes =\n%s\nwant\n%s", got, plain)
	}
}

// CA-6 (spec 010): an unknown --cores value is an option error, with exit
// code 2, checked before the fleet.
func TestRun_InvalidColors(t *testing.T) {
	for _, args := range [][]string{{"--cores", "xyz"}, {"--cores", "xyz", "--frota", "X=1"}} {
		code, stdout, stderr := run("sair\n", args...)
		if code != 2 {
			t.Errorf("Run(%q) code = %d, want 2", args, code)
		}
		if stdout != "" {
			t.Errorf("Run(%q) stdout = %q, want empty", args, stdout)
		}
		if want := "Cores inválidas: use sim ou nao.\n"; stderr != want {
			t.Errorf("Run(%q) stderr = %q, want %q", args, stderr, want)
		}
	}
}

// randomPlacement places the whole player's fleet at random.
const randomPlacement = "aleatorio\n"

// CA-8 (spec 012): giving up while placing ends before the game, with the seed.
func TestRun_PlacementGiveUp(t *testing.T) {
	code, stdout, stderr := run("sair\n", "--semente", "7")
	if code != 0 {
		t.Errorf("Run() code = %d, want 0 (stderr %q)", code, stderr)
	}
	if want := "Você desistiu antes da partida.\nSemente: 7\n"; !strings.HasSuffix(stdout, want) {
		t.Errorf("Run() stdout =\n%s\nwant suffix\n%s", stdout, want)
	}
}
