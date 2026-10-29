package text

import (
	"testing"

	"github.com/marcopollivier/battleship/internal/core/score"
)

// CA-1 and CA-4 (spec 025): the score line of a level, in singular or
// plural, with the best win only when there is a win.
func TestScoreLine(t *testing.T) {
	tests := []struct {
		level string
		r     score.Record
		want  string
	}{
		{"dificil", score.Record{Wins: 3, Losses: 2, BestWin: 41}, "Placar no nível dificil: 3 vitórias, 2 derrotas. Melhor vitória: 41 disparos."},
		{"facil", score.Record{Wins: 1, Losses: 1, BestWin: 1}, "Placar no nível facil: 1 vitória, 1 derrota. Melhor vitória: 1 disparo."},
		{"mestre", score.Record{Losses: 4}, "Placar no nível mestre: 0 vitórias, 4 derrotas."},
		{"mestre", score.Record{}, "Placar no nível mestre: 0 vitórias, 0 derrotas."},
	}
	for _, tt := range tests {
		if got := ScoreLine(tt.level, tt.r); got != tt.want {
			t.Errorf("ScoreLine(%q, %+v) = %q, want %q", tt.level, tt.r, got, tt.want)
		}
	}
}
