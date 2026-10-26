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
	"Comandos: aleatorio sorteia o resto da frota, recomecar tira todos os navios, sair desiste.\n" //nolint:misspell // Portuguese text for the player

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
		g.render(RenderBoard, g.player)
		g.write("\n")
		if !g.placeShip(scanner, remaining[0]) {
			return false
		}
	}
}

// placeShip asks for the position of one ship of type t until it is placed.
func (g *machineGame) placeShip(scanner *bufio.Scanner, t board.ShipType) bool {
	for {
		line, ok := g.prompt(scanner, placementPrompt(t))
		if !ok {
			return false
		}
		c, o, err := parsePlacement(line, t)
		if err == nil {
			err = g.player.Place(t, c, o)
		}
		if err == nil {
			return true
		}
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
