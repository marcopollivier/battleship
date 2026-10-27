// Package cmd wires the game domain to the terminal.
package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/marcopollivier/battleship/internal/board"
	"github.com/marcopollivier/battleship/internal/score"
	"github.com/marcopollivier/battleship/internal/terminal"
)

// Names of the command-line options.
const (
	fleetFlag = "frota"
	seedFlag  = "semente"
	modeFlag  = "modo"
	levelFlag = "nivel"
	colorFlag = "cores"

	showScoreFlag  = "placar"
	resetScoreFlag = "zerar-placar"
)

// Values of the mode option.
const (
	machineMode = "maquina"
	soloMode    = "solo"
)

// Values of the level option.
const (
	easyLevel   = "facil"
	hardLevel   = "dificil"
	masterLevel = "mestre"
)

// Values of the colors option.
const (
	colorsOn   = "sim"
	colorsOff  = "nao"
	colorsAuto = "auto"
)

// Run runs a game, against the machine or solo, with the command-line args (without the program name),
// reading shots from stdin, writing the game to stdout and errors to stderr.
// It keeps the score in the user's configuration directory. It returns the
// exit code.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	scorePath, err := score.DefaultPath()
	if err != nil {
		// With no configuration directory, every use of the score fails.
		scorePath = ""
	}
	return run(args, stdin, stdout, stderr, scorePath)
}

// run is Run with the score kept in the file at scorePath.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer, scorePath string) int {
	fs := flag.NewFlagSet("battleship", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fleetText := fs.String(fleetFlag, "", "frota no formato letra=quantidade separados por vírgula, por exemplo S=2,D=1,P=1 (padrão: a frota da cartela)")
	seed := fs.Uint64(seedFlag, 0, "número que define os tabuleiros e os disparos da máquina, para repetir uma partida (padrão: um número novo a cada execução)")
	mode := fs.String(modeFlag, machineMode, "maquina para jogar contra a máquina, ou solo para só atacar um tabuleiro sorteado")
	level := fs.String(levelFlag, hardLevel, "dificil para a máquina cercar os navios que acerta, facil para ela sortear todos os tiros, ou mestre para ela atirar onde ainda cabem mais navios")
	colors := fs.String(colorFlag, colorsAuto, "auto para usar cores só quando a saída é um terminal e NO_COLOR não está definida, sim para usar sempre, ou nao para nunca usar")
	showScore := fs.Bool(showScoreFlag, false, "mostra o placar contra a máquina em cada nível, sem começar partida")
	resetScore := fs.Bool(resetScoreFlag, false, "apaga o placar contra a máquina, sem começar partida")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if code, done := scoreCommand(stdout, scorePath, *showScore, *resetScore); done {
		return code
	}
	if !strings.EqualFold(*mode, machineMode) && !strings.EqualFold(*mode, soloMode) {
		fmt.Fprintln(stderr, "Modo inválido: use maquina ou solo.")
		return 2
	}
	style, ok := chooseStyle(*colors, isTerminal(stdout), os.Getenv("NO_COLOR"))
	if !ok {
		fmt.Fprintln(stderr, "Cores inválidas: use sim, nao ou auto.")
		return 2
	}
	strategy, ok := chooseStrategy(*level)
	if !ok {
		fmt.Fprintln(stderr, "Nível inválido: use facil, dificil ou mestre.")
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
	result, err := play(*mode, fleet, strategy, newRNG(*seed), style, stdin, stdout)
	if err != nil {
		fmt.Fprintln(stderr, terminal.FleetErrorMessage(err))
		return 1
	}
	if counts(*mode, flagGiven(fs, fleetFlag), result) {
		recordScore(stdout, scorePath, strings.ToLower(*level), result)
	}
	if _, err := fmt.Fprintf(stdout, "Semente: %d\n", *seed); err != nil {
		return 1
	}
	return 0
}

// play runs a game in mode with fleet. Solo, rng draws the board to attack.
// Against the machine, rng draws the machine's board first, then the ships
// the player asks to draw while placing, then the machine's shots, which
// strategy chooses; the player starts with an empty board and places the fleet.
// The boards are drawn with style. It returns how a game against the
// machine ended; a solo game returns the zero Result.
func play(mode string, fleet []board.FleetEntry, strategy board.ShotStrategy, rng *rand.Rand, style terminal.Style, stdin io.Reader, stdout io.Writer) (terminal.Result, error) {
	target, err := board.RandomBoard(fleet, rng)
	if err != nil {
		return terminal.Result{}, err
	}
	if strings.EqualFold(mode, soloMode) {
		return terminal.Result{}, terminal.PlaySolo(stdin, stdout, target, style)
	}
	return terminal.PlayVsMachine(stdin, stdout, board.NewBoard(fleet), target, rng, style, strategy)
}

// counts reports whether a game in mode that ended with result goes into
// the score: a game against the machine with the default fleet, won, lost
// or given up after the player's first shot.
func counts(mode string, fleetGiven bool, result terminal.Result) bool {
	if !strings.EqualFold(mode, machineMode) || fleetGiven {
		return false
	}
	return result.Outcome != terminal.GaveUp || result.Shots > 0
}

// recordScore adds result to the score of level in the file at scorePath
// and writes the score line of level to stdout.
func recordScore(stdout io.Writer, scorePath, level string, result terminal.Result) {
	b, err := score.Load(scorePath)
	if err != nil {
		return
	}
	if result.Outcome == terminal.PlayerWon {
		b.AddWin(level, result.Shots)
	} else {
		b.AddLoss(level)
	}
	if err := score.Save(scorePath, b); err != nil {
		return
	}
	fmt.Fprintln(stdout, terminal.ScoreLine(level, b[level]))
}

// scoreCommand runs the score options instead of a game: with reset it
// empties the score in the file at scorePath, and with show it writes the
// score line of every level to stdout. It returns the exit code, and false
// when neither option was given and the game should start.
func scoreCommand(stdout io.Writer, scorePath string, show, reset bool) (int, bool) {
	if !show && !reset {
		return 0, false
	}
	if reset {
		if err := score.Save(scorePath, score.Board{}); err != nil {
			return 0, true
		}
		fmt.Fprintln(stdout, "Placar zerado.")
	}
	if show {
		b, err := score.Load(scorePath)
		if err != nil {
			return 0, true
		}
		for _, level := range []string{easyLevel, hardLevel, masterLevel} {
			fmt.Fprintln(stdout, terminal.ScoreLine(level, b[level]))
		}
	}
	return 0, true
}

// chooseStyle returns the board style for the colors option (sim, nao or
// auto, in any case), given whether the output is a terminal and the value
// of NO_COLOR, and false for any other option value. With auto, the game is
// colored only on a terminal with NO_COLOR empty or unset.
func chooseStyle(colors string, tty bool, noColor string) (terminal.Style, bool) {
	switch {
	case strings.EqualFold(colors, colorsOn):
		return terminal.Colored, true
	case strings.EqualFold(colors, colorsOff):
		return terminal.Plain, true
	case strings.EqualFold(colors, colorsAuto):
		if tty && noColor == "" {
			return terminal.Colored, true
		}
		return terminal.Plain, true
	default:
		return terminal.Plain, false
	}
}

// isTerminal reports whether w is a character device, such as a terminal.
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

// chooseStrategy returns the machine's strategy for the level option, in
// any case, and false for an unknown level.
func chooseStrategy(level string) (board.ShotStrategy, bool) {
	switch {
	case strings.EqualFold(level, hardLevel):
		return board.SmartShot, true
	case strings.EqualFold(level, easyLevel):
		return board.RandomShot, true
	case strings.EqualFold(level, masterLevel):
		return board.MasterShot, true
	}
	return nil, false
}

// newRNG returns the random source that draws the boards for seed.
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
