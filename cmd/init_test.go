package cmd

import (
	"strings"
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
	"github.com/marcopollivier/battleship/internal/terminal"
)

// fleetRender returns the screen RenderBoard draws for a board with fleet.
func fleetRender(t *testing.T, fleet []board.FleetEntry) string {
	t.Helper()
	var sb strings.Builder
	if err := terminal.RenderBoard(&sb, board.NewBoard(fleet)); err != nil {
		t.Fatalf("RenderBoard: %v", err)
	}
	return sb.String()
}

// fleetRun runs the program with args and returns the exit code, stdout and stderr.
func fleetRun(args ...string) (code int, stdout, stderr string) {
	var out, errOut strings.Builder
	code = Run(args, &out, &errOut)
	return code, out.String(), errOut.String()
}

// CA-3: without --frota, the program uses the default fleet and shows the board.
func TestRun_DefaultFleet(t *testing.T) {
	code, stdout, stderr := fleetRun()
	if code != 0 {
		t.Errorf("Run() code = %d, want 0", code)
	}
	if want := fleetRender(t, board.DefaultFleet()); stdout != want {
		t.Errorf("Run() stdout =\n%s\nwant\n%s", stdout, want)
	}
	if stderr != "" {
		t.Errorf("Run() stderr = %q, want empty", stderr)
	}
}

// CA-10: an accepted fleet shows the board; a rejected fleet shows the
// Portuguese reason on stderr and exits with code 1.
func TestRun_Fleet(t *testing.T) {
	accepted := []board.FleetEntry{
		{Type: board.Submarine, Count: 2},
		{Type: board.AircraftCarrier, Count: 1},
	}
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{"accepted", []string{"--frota", "S=2,P=1"}, 0, fleetRender(t, accepted), ""},
		{"accepted with equals sign", []string{"--frota=S=2,P=1"}, 0, fleetRender(t, accepted), ""},
		{"does not fit", []string{"--frota", "E=30"}, 1, "", "Não foi possível posicionar essa frota no tabuleiro.\n"},
		{"unknown letter", []string{"--frota", "X=1"}, 1, "", "Frota inválida: use letra=quantidade separados por vírgula, por exemplo S=2,D=1,P=1.\n"},
		{"given but empty", []string{"--frota", ""}, 1, "", "Frota inválida: use letra=quantidade separados por vírgula, por exemplo S=2,D=1,P=1.\n"},
		{"negative count", []string{"--frota", "S=-1"}, 1, "", "Frota inválida: as quantidades não podem ser negativas.\n"},
		{"no ships", []string{"--frota", "S=0"}, 1, "", "Frota inválida: escolha pelo menos um navio.\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := fleetRun(tt.args...)
			if code != tt.wantCode {
				t.Errorf("Run(%q) code = %d, want %d", tt.args, code, tt.wantCode)
			}
			if stdout != tt.wantStdout {
				t.Errorf("Run(%q) stdout =\n%s\nwant\n%s", tt.args, stdout, tt.wantStdout)
			}
			if stderr != tt.wantStderr {
				t.Errorf("Run(%q) stderr = %q, want %q", tt.args, stderr, tt.wantStderr)
			}
		})
	}
}

// An unknown option is reported by the flag package with exit code 2.
func TestRun_UnknownOption(t *testing.T) {
	code, stdout, stderr := fleetRun("--nope")
	if code != 2 {
		t.Errorf("Run(--nope) code = %d, want 2", code)
	}
	if stdout != "" {
		t.Errorf("Run(--nope) stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "-frota") {
		t.Errorf("Run(--nope) stderr = %q, want the usage with -frota", stderr)
	}
}
