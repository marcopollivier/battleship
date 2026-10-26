package cmd

import (
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

// emptyAttackBoard is the attack board before any shot.
func emptyAttackBoard(t *testing.T) string {
	t.Helper()
	var sb strings.Builder
	if err := terminal.RenderAttackBoard(&sb, board.NewBoard(board.DefaultFleet())); err != nil {
		t.Fatalf("RenderAttackBoard: %v", err)
	}
	return sb.String()
}

// CA-4: shooting every ship cell wins, and the seed is shown at the end.
func TestRun_Win(t *testing.T) {
	code, stdout, stderr := run(seedShots(t, board.DefaultFleet(), 42), "--semente", "42")
	if code != 0 {
		t.Errorf("Run() code = %d, want 0 (stderr %q)", code, stderr)
	}
	if want := "Você afundou toda a frota em 19 disparos!\nSemente: 42\n"; !strings.HasSuffix(stdout, want) {
		t.Errorf("Run() stdout =\n%s\nwant suffix\n%s", stdout, want)
	}
}

// CA-1 and CA-5: without --semente the game starts with the empty attack
// board, and giving up shows the fleet and the seed that was drawn.
func TestRun_GiveUp(t *testing.T) {
	code, stdout, stderr := run("sair\n")
	if code != 0 {
		t.Errorf("Run() code = %d, want 0 (stderr %q)", code, stderr)
	}
	if want := emptyAttackBoard(t) + "\nDisparo: "; !strings.HasPrefix(stdout, want) {
		t.Errorf("Run() stdout =\n%s\nwant prefix\n%s", stdout, want)
	}
	if !strings.Contains(stdout, "Você desistiu depois de 0 disparos. A frota era:") {
		t.Errorf("Run() stdout =\n%s\nwant the give up message", stdout)
	}
	if !regexp.MustCompile(`\nSemente: \d+\n$`).MatchString(stdout) {
		t.Errorf("Run() stdout =\n%s\nwant the seed on the last line", stdout)
	}
}

// CA-6: the same seed draws the same board.
func TestRun_SameSeed(t *testing.T) {
	_, first, _ := run("sair\n", "--semente", "42")
	_, second, _ := run("sair\n", "--semente", "42")
	if first != second {
		t.Errorf("two runs with --semente 42 differ:\n%s\n---\n%s", first, second)
	}
	if _, other, _ := run("sair\n", "--semente", "43"); other == first {
		t.Error("runs with --semente 42 and 43 drew the same board")
	}
}

// CA-7: the game uses the fleet chosen with --frota; rejected fleets end the
// program before the game, with the reason on stderr and exit code 1.
func TestRun_Fleet(t *testing.T) {
	twoSubs := []board.FleetEntry{{Type: board.Submarine, Count: 2}}
	code, stdout, stderr := run(seedShots(t, twoSubs, 7), "--frota", "S=2", "--semente", "7")
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

// CA-8 and unknown options: the flag package rejects them with exit code 2
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
