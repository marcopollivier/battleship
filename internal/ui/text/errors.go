package text

import (
	"errors"
	"strings"

	"github.com/marcopollivier/battleship/internal/core/board"
)

// ShotErrorMessage returns the Portuguese message for a shot at c that was
// refused with err, and false if err is not a mistake of the player.
func ShotErrorMessage(err error, c board.Coord) (string, bool) {
	switch {
	case errors.Is(err, board.ErrInvalidFormat):
		return "Coordenada inválida: use uma letra de A a J e um número de 1 a 10, por exemplo B7.", true
	case errors.Is(err, board.ErrOutOfBounds):
		return "Fora do tabuleiro: use uma letra de A a J e um número de 1 a 10.", true
	case errors.Is(err, board.ErrAlreadyShot):
		return "Você já atirou em " + c.String() + ".", true
	default:
		return "", false
	}
}

// PlacementErrorMessage returns the Portuguese message for a placement of a
// ship of type t that the board refused with err, and false if err is not a
// refusal of the board. How each interface reads a position, and the
// message for a position it cannot read, is up to the interface.
func PlacementErrorMessage(err error, t board.ShipType) (string, bool) {
	switch {
	case errors.Is(err, board.ErrInvalidOrientation):
		return "Orientação inválida para o " + ShipName(t) + ": use " + joinOr(OrientationList(t)) + ".", true
	case errors.Is(err, board.ErrOutOfBounds):
		return "Fora do tabuleiro: o " + ShipName(t) + " não cabe nessa posição.", true
	case errors.Is(err, board.ErrOverlap):
		return "Posição ocupada: o " + ShipName(t) + " ficaria em cima de outro navio.", true
	case errors.Is(err, board.ErrAdjacent):
		return "Muito perto: os navios não podem se encostar, nem pela diagonal.", true
	default:
		return "", false
	}
}

// joinOr lists names with commas and "ou" before the last one, such as
// "cima, baixo, esquerda ou direita".
func joinOr(names []string) string {
	if len(names) < 2 {
		return strings.Join(names, "")
	}
	return strings.Join(names[:len(names)-1], ", ") + " ou " + names[len(names)-1]
}
