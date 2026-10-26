// Package cmd wires the game domain to the terminal.
package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"

	"github.com/marcopollivier/battleship/internal/board"
	"github.com/marcopollivier/battleship/internal/terminal"
)

// Names of the command-line options.
const (
	fleetFlag = "frota"
	seedFlag  = "semente"
)

// Run runs a solo game with the command-line args (without the program name),
// reading shots from stdin, writing the game to stdout and errors to stderr.
// It returns the exit code.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("battleship", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fleetText := fs.String(fleetFlag, "", "frota no formato letra=quantidade separados por vírgula, por exemplo S=2,D=1,P=1 (padrão: a frota da cartela)")
	seed := fs.Uint64(seedFlag, 0, "número que define o tabuleiro sorteado, para repetir uma partida (padrão: um número novo a cada execução)")
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
	if !flagGiven(fs, seedFlag) {
		*seed = rand.Uint64() //nolint:gosec // a game seed, not security-sensitive
	}
	b, err := board.RandomBoard(fleet, newRNG(*seed))
	if err != nil {
		fmt.Fprintln(stderr, terminal.FleetErrorMessage(err))
		return 1
	}
	if err := terminal.PlaySolo(stdin, stdout, b); err != nil {
		fmt.Fprintln(stderr, "Erro:", err)
		return 1
	}
	if _, err := fmt.Fprintf(stdout, "Semente: %d\n", *seed); err != nil {
		return 1
	}
	return 0
}

// newRNG returns the random source that draws the board for seed.
func newRNG(seed uint64) *rand.Rand {
	return rand.New(rand.NewPCG(seed, 0)) //nolint:gosec // a game board, not security-sensitive
}

// chooseFleet returns the default fleet when the fleet option was not given,
// and otherwise the parsed and validated fleet text.
func chooseFleet(fs *flag.FlagSet, text string) ([]board.FleetEntry, error) {
	if !flagGiven(fs, fleetFlag) {
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

// flagGiven reports whether the option name was set on the command line.
func flagGiven(fs *flag.FlagSet, name string) bool {
	given := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			given = true
		}
	})
	return given
}
