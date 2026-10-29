package text

import (
	"errors"
	"fmt"
	"testing"

	"github.com/marcopollivier/battleship/internal/core/board"
)

func TestJoinOr(t *testing.T) {
	tests := []struct {
		names []string
		want  string
	}{
		{nil, ""},
		{[]string{"cima"}, "cima"},
		{[]string{"horizontal", "vertical"}, "horizontal ou vertical"},
		{[]string{"cima", "baixo", "esquerda", "direita"}, "cima, baixo, esquerda ou direita"},
	}
	for _, tt := range tests {
		if got := joinOr(tt.names); got != tt.want {
			t.Errorf("joinOr(%q) = %q, want %q", tt.names, got, tt.want)
		}
	}
}

func TestShotErrorMessage(t *testing.T) {
	c := board.Coord{Col: 1, Row: 6}
	tests := []struct {
		err  error
		want string
		ok   bool
	}{
		{board.ErrInvalidFormat, "Coordenada inválida: use uma letra de A a J e um número de 1 a 10, por exemplo B7.", true},
		{board.ErrOutOfBounds, "Fora do tabuleiro: use uma letra de A a J e um número de 1 a 10.", true},
		{fmt.Errorf("context: %w", board.ErrAlreadyShot), "Você já atirou em B7.", true},
		{errors.New("boom"), "", false},
	}
	for _, tt := range tests {
		got, ok := ShotErrorMessage(tt.err, c)
		if got != tt.want || ok != tt.ok {
			t.Errorf("ShotErrorMessage(%v) = %q, %v, want %q, %v", tt.err, got, ok, tt.want, tt.ok)
		}
	}
}

func TestPlacementErrorMessage(t *testing.T) {
	tests := []struct {
		err  error
		want string
		ok   bool
	}{
		{board.ErrInvalidOrientation, "Orientação inválida para o Porta-Aviões: use cima, baixo, esquerda ou direita.", true},
		{board.ErrOutOfBounds, "Fora do tabuleiro: o Porta-Aviões não cabe nessa posição.", true},
		{fmt.Errorf("context: %w", board.ErrOverlap), "Posição ocupada: o Porta-Aviões ficaria em cima de outro navio.", true},
		{board.ErrAdjacent, "Muito perto: os navios não podem se encostar, nem pela diagonal.", true},
		{board.ErrInvalidFormat, "", false},
		{errors.New("boom"), "", false},
	}
	for _, tt := range tests {
		got, ok := PlacementErrorMessage(tt.err, board.AircraftCarrier)
		if got != tt.want || ok != tt.ok {
			t.Errorf("PlacementErrorMessage(%v) = %q, %v, want %q, %v", tt.err, got, ok, tt.want, tt.ok)
		}
	}
}
