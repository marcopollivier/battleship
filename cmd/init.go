// Package cmd wires the game domain to the terminal.
package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/marcopollivier/battleship/internal/board"
	"github.com/marcopollivier/battleship/internal/terminal"
)

// fleetFlag is the name of the command-line option that sets the fleet.
const fleetFlag = "frota"

// Run runs the game with the command-line args (without the program name),
// writing the screen to stdout and errors to stderr. It returns the exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("battleship", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fleetText := fs.String(fleetFlag, "", "frota no formato letra=quantidade separados por vírgula, por exemplo S=2,D=1,P=1 (padrão: a frota da cartela)")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	fleet, err := chooseFleet(fs, *fleetText)
	if err != nil {
		fmt.Fprintln(stderr, terminal.FleetErrorMessage(err))
		return 1
	}
	if err := terminal.RenderBoard(stdout, board.NewBoard(fleet)); err != nil {
		fmt.Fprintln(stderr, terminal.FleetErrorMessage(err))
		return 1
	}
	return 0
}

// chooseFleet returns the default fleet when the fleet option was not given,
// and otherwise the parsed and validated fleet text.
func chooseFleet(fs *flag.FlagSet, text string) ([]board.FleetEntry, error) {
	given := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == fleetFlag {
			given = true
		}
	})
	if !given {
		return board.DefaultFleet(), nil
	}
	fleet, err := terminal.ParseFleet(text)
	if err != nil {
		return nil, err
	}
	if err := board.ValidateFleet(fleet); err != nil {
		return nil, err
	}
	return fleet, nil
}
