package cmd

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/marcopollivier/battleship/internal/board"
	"github.com/marcopollivier/battleship/internal/score"
)

// scoreStart is the score the spec 025 tests start from.
var scoreStart = score.Board{
	"dificil": {Wins: 1, Losses: 2, BestWin: 30},
	"mestre":  {Losses: 1},
}

// scoreCase is a game whose result changes the score: the level it is
// played at, the score it leaves and the end of the output.
type scoreCase struct {
	name       string
	input      string
	args       []string
	level      string
	want       score.Record
	wantSuffix string
}

// runScoreCases plays each case from scoreStart and checks the score file
// and the end of the output.
func runScoreCases(t *testing.T, tests []scoreCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := scoreFile(t, scoreStart)
			code, stdout, stderr := runScore(path, tt.input, append([]string{"--semente", "42", "--cores", "nao"}, tt.args...)...)
			if code != 0 || stderr != "" {
				t.Errorf("Run(%q) code = %d, stderr %q, want 0 and no error", tt.args, code, stderr)
			}
			want := score.Board{}
			for level, r := range scoreStart {
				want[level] = r
			}
			want[tt.level] = tt.want
			if got := loadScore(t, path); !reflect.DeepEqual(got, want) {
				t.Errorf("score after Run(%q) = %+v, want %+v", tt.args, got, want)
			}
			if !strings.HasSuffix(stdout, tt.wantSuffix) {
				t.Errorf("Run(%q) stdout =\n%s\nwant suffix\n%s", tt.args, stdout, tt.wantSuffix)
			}
		})
	}
}

// CA-1 (spec 025): a win with the default fleet adds a win and its shots to
// the score of its level, in lower case, and the score line comes before
// the seed. The player shoots every enemy ship cell, so never misses.
func TestRun_ScoreWin(t *testing.T) {
	input := randomPlacement + machineCells(machineEnemy(t, board.DefaultFleet(), 42), true)
	const won = "Você venceu! Afundou a frota inimiga em 19 disparos.\n"
	runScoreCases(t, []scoreCase{
		{"best win", input, nil, "dificil", score.Record{Wins: 2, Losses: 2, BestWin: 19},
			won + "Placar no nível dificil: 2 vitórias, 2 derrotas. Melhor vitória: 19 disparos.\nSemente: 42\n"},
		{"first win", input, []string{"--nivel", "MESTRE"}, "mestre", score.Record{Wins: 1, Losses: 1, BestWin: 19},
			won + "Placar no nível mestre: 1 vitória, 1 derrota. Melhor vitória: 19 disparos.\nSemente: 42\n"},
	})
}

// CA-2 (spec 025): a loss, or giving up after the first shot, adds a loss
// to the score of its level.
func TestRun_ScoreLoss(t *testing.T) {
	lose := randomPlacement + machineCells(machineEnemy(t, board.DefaultFleet(), 42), false)
	runScoreCases(t, []scoreCase{
		{"machine wins", lose, []string{"--nivel", "facil"}, "facil", score.Record{Losses: 1},
			"Placar no nível facil: 0 vitórias, 1 derrota.\nSemente: 42\n"},
		{"sair after a shot", randomPlacement + "A1\nsair\n", nil, "dificil", score.Record{Wins: 1, Losses: 3, BestWin: 30},
			"Placar no nível dificil: 1 vitória, 3 derrotas. Melhor vitória: 30 disparos.\nSemente: 42\n"},
		{"end of input after a shot", randomPlacement + "A1\n", []string{"--nivel", "mestre"}, "mestre", score.Record{Losses: 2},
			"Placar no nível mestre: 0 vitórias, 2 derrotas.\nSemente: 42\n"},
	})
}

// CA-3 (spec 025): giving up while placing or before the first shot, a solo
// game and a game with --frota leave the score file as it was and show no
// score line.
func TestRun_ScoreNotCounted(t *testing.T) {
	oneSub := []board.FleetEntry{{Type: board.Submarine, Count: 1}}
	tests := []struct {
		name  string
		input string
		args  []string
	}{
		{"sair while placing", "sair\n", nil},
		{"sair before shooting", randomPlacement + "sair\n", nil},
		{"solo", seedShots(t, board.DefaultFleet(), 42), []string{"--modo", "solo"}},
		{"frota win", randomPlacement + machineCells(machineEnemy(t, oneSub, 42), true), []string{"--frota", "S=1"}},
		{"frota loss", randomPlacement + machineCells(machineEnemy(t, oneSub, 42), false), []string{"--frota", "S=1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := scoreFile(t, scoreStart)
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile: %v", err)
			}
			code, stdout, stderr := runScore(path, tt.input, append([]string{"--semente", "42", "--cores", "nao"}, tt.args...)...)
			if code != 0 || stderr != "" {
				t.Errorf("Run(%q) code = %d, stderr %q, want 0 and no error", tt.args, code, stderr)
			}
			if after, _ := os.ReadFile(path); string(after) != string(before) {
				t.Errorf("score file after Run(%q) = %s, want %s", tt.args, after, before)
			}
			if strings.Contains(stdout, "Placar") {
				t.Errorf("Run(%q) stdout =\n%s\nwant no score line", tt.args, stdout)
			}
		})
	}
}
