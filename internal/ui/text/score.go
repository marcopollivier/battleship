package text

import (
	"fmt"

	"github.com/marcopollivier/battleship/internal/core/score"
)

// ScoreLine is the score of level in Portuguese, such as "Placar no nível
// dificil: 3 vitórias, 2 derrotas. Melhor vitória: 41 disparos.". With no
// win, the best win is left out.
func ScoreLine(level string, r score.Record) string {
	line := fmt.Sprintf("Placar no nível %s: %s, %s.", level,
		plural(r.Wins, "vitória", "vitórias"), plural(r.Losses, "derrota", "derrotas"))
	if r.Wins > 0 {
		line += " Melhor vitória: " + ShotCount(r.BestWin) + "."
	}
	return line
}

// plural writes n with the singular word when n is 1, and with the plural
// word otherwise, such as "1 derrota" or "0 derrotas".
func plural(n int, singular, plural string) string {
	if n == 1 {
		return "1 " + singular
	}
	return fmt.Sprintf("%d %s", n, plural)
}

// ShotCount writes n shots in Portuguese, such as "1 disparo" or "2 disparos".
func ShotCount(n int) string {
	return plural(n, "disparo", "disparos")
}
