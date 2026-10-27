package terminal

import (
	"bufio"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/marcopollivier/battleship/internal/board"
)

// errInvalidPlacement means the line is not a cell followed by an orientation.
var errInvalidPlacement = errors.New("invalid placement format")

// parsePlacement reads "E5 vertical", "e5 v" or, for a one-cell ship, "E5".
// The orientation is its Portuguese name or first letter, in any case, and
// must apply to t.
func parsePlacement(line string, t board.ShipType) (board.Coord, board.Orientation, error) {
	fields := strings.Fields(line)
	oneCell := t.Size() == 1
	if len(fields) == 0 || len(fields) > 2 || (len(fields) == 1 && !oneCell) {
		return board.Coord{}, 0, fmt.Errorf("%w: %q", errInvalidPlacement, line)
	}
	c, err := board.ParseCoord(fields[0])
	if err != nil {
		return board.Coord{}, 0, err
	}
	if len(fields) == 1 {
		return c, board.Horizontal, nil
	}
	o, ok := orientationByName(fields[1])
	if !ok || !slices.Contains(board.ValidOrientations(t), o) {
		return board.Coord{}, 0, fmt.Errorf("%w: %v %q", board.ErrInvalidOrientation, t, fields[1])
	}
	return c, o, nil
}

// orientationByName finds an orientation by its Portuguese name or first
// letter, in any case.
func orientationByName(name string) (board.Orientation, bool) {
	name = strings.ToLower(name)
	for o, full := range orientationNames {
		if name == full || name == full[:1] {
			return o, true
		}
	}
	return 0, false
}

// placementInstructions explain how to place the fleet.
const placementInstructions = "Posicione a sua frota. Para cada navio, digite a casa e a orientação, por exemplo E5 vertical.\n" +
	"A casa é a primeira do navio (a de cima ou a da esquerda). No Porta-Aviões, é o meio da barra do T.\n" +
	"Comandos: aleatorio sorteia o resto da frota, desfazer tira o último navio, recomecar tira todos os navios, sair desiste.\n" //nolint:misspell // Portuguese text for the player

// placementWarning is shown before the board while the ships still missing
// no longer fit on it, in yellow. It does not stop the player.
const placementWarning = "Atenção: o resto da frota não cabe mais no tabuleiro. Use desfazer ou recomecar." //nolint:misspell // Portuguese text for the player

// placeFleet asks for the missing ships of the player's board, largest
// first, and reports false if the player gave up or a write failed.
func (g *machineGame) placeFleet(scanner *bufio.Scanner) bool {
	g.write(placementInstructions + "\n")
	for {
		remaining := g.player.Remaining()
		if len(remaining) == 0 {
			g.write("Frota posicionada!\n\n")
			return g.err == nil
		}
		if len(g.player.Ships()) > 0 && !board.Fits(g.player) {
			g.write(g.style.paint(placementWarning, warnColor) + "\n")
		}
		g.render(RenderBoard, g.player)
		g.write("\n")
		if !g.placeShip(scanner, remaining[0]) {
			if g.err == nil {
				g.write("Você desistiu antes da partida.\n")
			}
			return false
		}
	}
}

// Commands accepted instead of a position, in any letter case.
const (
	randomCommand  = "aleatorio"
	undoCommand    = "desfazer"
	restartCommand = "recomecar"
)

// placeShip asks for the position of one ship of type t until it is placed,
// or until a command fills the board or removes ships from it. Invalid
// positions and a random draw that does not fit show a message in yellow. It
// reports false if the player gave up or a write failed.
func (g *machineGame) placeShip(scanner *bufio.Scanner, t board.ShipType) bool {
	for {
		line, ok := g.prompt(scanner, placementPrompt(t))
		if !ok {
			return false
		}
		switch strings.ToLower(strings.TrimSpace(line)) {
		case randomCommand:
			if err := board.FillRandom(g.player, g.rng); err != nil {
				g.write(g.style.paint("Não foi possível sortear o resto da frota. Mude a posição dos navios com recomecar.", warnColor) + "\n")
				continue
			}
			return true
		case undoCommand:
			// The ships are asked largest first, so the ship removed is the
			// next one placeFleet asks for.
			if g.player.RemoveLast() {
				return true
			}
			g.write(g.style.paint("Não há navio para desfazer.", warnColor) + "\n")
			continue
		case restartCommand:
			g.player.Clear()
			return true
		}
		c, o, err := parsePlacement(line, t)
		if err == nil {
			err = g.player.Place(t, c, o)
		}
		if err == nil {
			return true
		}
		message, ok := placementErrorMessage(err, t)
		if !ok {
			g.err = err
			return false
		}
		g.write(g.style.paint(message, warnColor) + "\n")
	}
}

// placementErrorMessage returns the Portuguese message for a placement of a
// ship of type t refused with err, and false if err is not a mistake of the
// player.
func placementErrorMessage(err error, t board.ShipType) (string, bool) {
	switch {
	case errors.Is(err, errInvalidPlacement), errors.Is(err, board.ErrInvalidFormat):
		return "Posição inválida: digite a casa e a orientação, por exemplo E5 vertical.", true
	case errors.Is(err, board.ErrInvalidOrientation):
		return "Orientação inválida para o " + ShipName(t) + ": use " + joinOr(orientationList(t)) + ".", true
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

// placementPrompt asks for a ship, such as "Destroyer, 2 casas (horizontal,
// vertical): ". One-cell ships list no orientation.
func placementPrompt(t board.ShipType) string {
	if t.Size() == 1 {
		return ShipName(t) + ", 1 casa: "
	}
	return fmt.Sprintf("%s, %d casas (%s): ", ShipName(t), t.Size(), strings.Join(orientationList(t), ", "))
}

// orientationList names the orientations that apply to t, in Portuguese.
func orientationList(t board.ShipType) []string {
	var names []string
	for _, o := range board.ValidOrientations(t) {
		names = append(names, OrientationName(o))
	}
	return names
}

// joinOr lists names with commas and "ou" before the last one, such as
// "cima, baixo, esquerda ou direita".
func joinOr(names []string) string {
	if len(names) < 2 {
		return strings.Join(names, "")
	}
	return strings.Join(names[:len(names)-1], ", ") + " ou " + names[len(names)-1]
}
