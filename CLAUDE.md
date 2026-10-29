# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

A terminal Battleship ("Batalha Naval") game in Go (module `github.com/marcopollivier/battleship`, Go 1.24, no external dependencies). Today it plays against the machine, after the player places their own fleet, or a solo game (`--modo solo`) where the player only attacks a random board. The README (in Portuguese) is the game spec.

## Commands

All are Makefile targets (`make help` lists them):

- `make run`: `go run ./main.go`, which starts a game against the machine. Options go in `ARGS`, as in `make run ARGS="--modo solo --frota S=2,P=1 --semente 42 --cores nao --limpar nao --nivel facil"`, or after `go run .`
- `make build`: builds the binary to `bin/battleship`
- `make release-check`: runs `scripts/build-release.sh dev`, which builds the 5 release binaries (`battleship-<version>-<os>-<arch>`, darwin and linux on amd64 and arm64, and windows/amd64) into `dist/` (git-ignored)
- `make test`: `go test -v ./...`
- Run a single test: `go test -v -run TestName ./cmd/`
- `make lint`: `golangci-lint run` (config in `.golangci.yml`)
- `make fmt` / `make vet` / `make mod-tidy`
- `make pre-commit`: fmt, vet, lint, sec and test together
- `make setup`: installs pre-commit (via Homebrew) and golangci-lint, then installs the git hooks

The pre-commit hooks run go-fmt, go-mod-tidy, go unit tests and golangci-lint on every commit, so code must pass lint and tests before it can be committed.

## Spec-driven workflow

Features are built through the spec-driven development (SDD) process in `specs/README.md`, and every plan must follow the principles in `specs/constitution.md`. The order is spec, then plan, then tasks, then implementation, with explicit user approval needed before moving from one phase to the next. Don't write production code for a feature that has no approved `specs/NNN-slug/tasks.md`. If implementation contradicts the spec, stop and propose a spec change. Each phase has a slash command: `/spec`, `/plan`, `/tasks` and `/implement`. The feature backlog and the open game-rule questions are in `specs/ROADMAP.md`. Independent specs are worked on in parallel waves, each implemented by a separate agent in its own git worktree. During a wave, spec PRs must not edit `specs/ROADMAP.md`; only the coordinator's `docs/` PRs update it. See "Trabalho em paralelo" in `specs/README.md`.

Never commit directly to `master`. Each context gets its own branch and PR: `feat/NNN-slug` for a spec, and `docs/`, `chore/` or `fix/` for anything else. The full convention is under "Branches e PRs" in `specs/README.md`. Work that falls outside the current PR's context goes on a separate branch, or into the roadmap.

## Architecture

The code is split into **core** (`internal/core`: the rules, with no I/O) and **ui** (`internal/ui`: the Portuguese text and each interface). The core never imports the ui (constitution, principle 1). The v2 plan in `specs/ROADMAP.md` adds a game session to the core (spec 028) and a web interface (spec 029):

- **`internal/core/board`: the game domain.** It has no I/O and its names are in English.
  - `coord.go`: `Coord` (zero-based `Col`, `Row`) and `ParseCoord`, which turns `"A1"` into a `Coord`.
  - `ship.go`: `ShipType`, `Orientation`, `DefaultFleet`, `ShipCells` and `ValidOrientations`.
  - `board.go`: `Board` with `Place`, which validates fleet limit, orientation, bounds, overlap and adjacency in that order and leaves the board unchanged on error; `Remaining` (ships still to place, largest first), `RemoveLast` (removes the ship placed last) and `Clear`.
  - `shot.go`: `Fire`, `Cell` and `AllSunk`.
  - `random.go`: `RandomBoard` and `FillRandom` (places only the missing ships), which take an injected `*rand.Rand`.
  - `fit.go`: `Fits`, whether the missing ships can still all be placed. An exact backtracking search without randomness on a copy of the board; it answers true once it passes 1,000,000 attempts, so it only answers false when it is sure.
  - `machine.go`: `ShotStrategy`, the type of a machine's way to choose its next shot. `RandomShot` draws among the cells not shot yet. `SmartShot` hunts at random, then shoots next to the hits on a ship afloat (the ends of a line of hits first) until it sinks; it skips the cells around sunk ships (`aroundSunk`, a provisional rule, issue #25) and only uses what the shooter knows.
  - `master.go`: `MasterShot`, a third `ShotStrategy`. It knows the fleet and which ships have sunk, and counts, for each cell not shot yet, the positions of the ships still afloat that could cover it: on the board, off water and sunk cells, and not touching a hit outside the position. With hits on a ship afloat, only the positions covering them count, each weighing the number of hits it covers. It shoots at a cell with the highest count, drawing among ties with the `rng`.
  - `fleet.go`: `ValidateFleet`, which rejects negative counts, empty fleets and fleets that do not fit on the board.
  - Errors are sentinel values wrapped with `%w`; check them with `errors.Is`.
- **`internal/ui/text`: the Portuguese text shared by every interface.** It reads the core and has no I/O.
  - `names.go`: `ShipName`, `OrientationName`, `ShipSymbol` (the letters `S D C E P`), `OrientationByName` (a Portuguese name or its first letter, in any case) and `OrientationList`.
  - `messages.go`: `ShotMessage`, the result of a shot.
  - `errors.go`: `ShotErrorMessage` and `PlacementErrorMessage`, the messages for a shot or placement refused by the board, with false when the error is not a mistake of the player. Each interface writes its own message for input it cannot read.
  - `fleet.go`: `ParseFleet`, which reads the `--frota` text (such as `S=2,D=1,P=1`), and `FleetErrorMessage`.
  - `score.go`: `ScoreLine(level, record)`, the score line of a level (`Placar no nível dificil: 3 vitórias, 2 derrotas. Melhor vitória: 41 disparos.`), with singular and plural and no best win without a win, and `ShotCount` (`1 disparo`, `2 disparos`).
- **`internal/ui/terminal`: the terminal interface.** It reads the core, takes its text from `ui/text` and writes to an `io.Writer`.
  - `color.go`: `Style` (`Plain`, the zero value, or `Colored`) and `paint`, which wraps a board symbol or a whole message in its ANSI color code and a reset when the style is `Colored`; `shotColor` gives a shot result's message the color of its board symbol (cyan for a miss, bold red for a hit or a sunk ship).
  - Message colors (spec 015): the shot results (the player's, and the machine's only after `A máquina atirou em C4: `), the wins in bold green, the machine's win in bold red, and invalid shots and placements in yellow. Prompts, instructions, titles and the give-up message stay uncolored, and line breaks stay outside the codes. `ShotMessage` itself stays plain; the callers paint it.
  - `render.go`: `RenderBoard(w, board, style)`, the player's own board with ships and received shots, then its legend; the cells of a sunk ship are drawn as `#`. `ownGrid` returns the board's lines (header and 10 rows, each `gridWidth` = 22 columns wide without color codes) without the legend.
  - `attack.go`: `RenderAttackBoard(w, board, style)`, the opponent's board showing only shots fired, and sunk ships; `attackGrid` returns its lines.
  - `screen.go`: `Display` (`Style` and `Clear`, whether the terminal is cleared before each screen), `clearScreen` (the ANSI code that clears it), `sideBySide`, which joins two grids line by line with a 5-space gap, and `boardTitles`.
  - Both renderers paint each cell while building the grid; `writeBoard` adds the header, row numbers and legend, which stay uncolored.
  - `solo.go`: `PlaySolo(in, out, board, display)`, the solo game loop: reads one shot per line until every ship is sunk or the player types `sair` or closes the input. Its `playerAttack` type is the player's turn, shared with `machine.go`: it keeps the `Style`, whether to clear the terminal, and the `events` (such as `Seu disparo em B7: Acerto!`) written below the boards of the next screen. Each shot draws a new screen (clear, attack board, events, prompt); the last screen reveals the board.
  - `machine.go`: `PlayVsMachine(in, out, player, enemy, rng, display, shoot) (Result, error)`, the game against the machine, which shoots with the `ShotStrategy` `shoot`: whoever hits shoots again. If `player` is not complete, it first runs the placing phase. Every player shot draws a new screen: clear (if `display.Clear`), the titles, both boards side by side, the short screen legends (to stay within 80 columns), the events since the last screen (the player's result and the machine's shots) and the prompt. Invalid input is written below the prompt with no new screen. The game ends on a last screen with both fleets revealed and the final message. `Result` holds the `Outcome` (`GaveUp`, while placing or shooting, `PlayerWon` or `MachineWon`) and the player's valid `Shots`.
  - `placement.go`: the placing phase (`casa orientação` lines, `aleatorio`, `desfazer`, `recomecar`, `sair`), each ship on a new screen (clear, instructions, the warning while `Fits` is false, the board, the prompt); `Frota posicionada!` is the first event of the game screen. It runs inside `PlayVsMachine` so one `bufio.Scanner` reads all the input.
- **`internal/core/score`: the score across games (spec 025).** `Record` (`Wins`, `Losses`, `BestWin`, the fewest player shots in a win) and `Board`, a `map[string]Record` by level name, with `AddWin` and `AddLoss`, which do no I/O. `Load` reads the JSON file (a missing file is an empty board; a corrupt one is an error and is never reset on its own), `Save` writes it through a temporary file renamed into place, and `DefaultPath` is `battleship/placar.json` in `os.UserConfigDir`. An empty path is `ErrNoPath`.
- **`cmd/` and `main.go`: wiring.** `cmd.Run(args, stdin, stdout, stderr) int` parses `--modo`, `--cores` (`auto`, the default, `sim` or `nao`; `chooseStyle(colors, tty, noColor)` is a pure function, and with `auto` it colors only when `isTerminal(stdout)`, an `*os.File` that is a character device, and `NO_COLOR` is empty; `sim` ignores `NO_COLOR`), `--limpar` (`auto`, the default, `sim` or `nao`; `chooseClear(value, tty)` is pure, and with `auto` it clears only when the output is a terminal), `--nivel` (`dificil`, the default, picks `SmartShot`; `facil` picks `RandomShot`; `mestre` picks `MasterShot`), `--frota`, `--semente`, `--placar` and `--zerar-placar`. `--versao` is handled right after parsing, before the other checks: it prints `versionText(version(debug.ReadBuildInfo()))` (in `version.go`: `battleship v1.0.0`, or `battleship (desenvolvimento)` when the build has no version, `(devel)`) and returns 0. Then it checks the mode, then the colors, then the clearing, then the level, then the fleet. Colors and clearing make the `terminal.Display` passed to the game. One `*rand.Rand` seeded by the seed (a new one each run unless given) draws the machine's board, then the ships the player draws with `aleatorio`, then the machine's shots; `Run` runs `PlayVsMachine` (or `PlaySolo`) and prints the seed. It returns the exit code: 0 at the end of a game, on `-h` or on `--versao`, 1 on a fleet error, and 2 on an invalid option. `main.go` passes it to `os.Exit`.
  - `Run` only finds the score path with `score.DefaultPath` (empty when there is no configuration directory) and calls the internal `run(args, stdin, stdout, stderr, scorePath)`. Tests always call `run` with a path under `t.TempDir()` (the `runCmd` and `runScore` helpers), never `Run`, so they never touch the real score.
  - Right after `--versao`, `--zerar-placar` saves an empty score and prints `Placar zerado.`, and `--placar` prints the `ScoreLine` of `facil`, `dificil` and `mestre`; either one ends with code 0 without a game (`scoreCommand`).
  - A game counts (`counts`) against the machine, without `--frota`, when won, lost, or given up after the first shot. `recordScore` then loads the score, adds the result at the level in lower case, saves it and prints the `ScoreLine` before the seed. Any read or write error prints `Não foi possível guardar o placar.` on stderr instead of the line, and the exit code stays the game's.

## Releases

`.github/workflows/release.yml` runs when a `v*` tag is pushed: it runs `scripts/build-release.sh <tag>` and publishes the binaries with `gh release create --generate-notes`. The version comes from the Go tool, not from `-ldflags`: `go install ...@vX.Y.Z` stamps the module version, and `go build .` in a checkout at a tag stamps the tag (an untagged checkout gets a pseudo-version, with `+dirty` if it has changes). `make run`, `make build` and `go run` stamp nothing, so they show `(desenvolvimento)`. The script must build the package (`.`), not `main.go`, or nothing is stamped. Only the maintainer pushes tags.

## Game rules (from README)

- The board is 10x10: columns A–J, rows 1–10. Shots are written as coordinates such as `A1`, `B3`, `J10`.
- A player who hits a ship shoots again. A miss ("tiro na água") passes the turn to the other player.
- A ship is sunk when all of its cells are hit. The first player to sink every enemy ship wins.
- The default fleet comes from the original paper pad shown in `doc/img/batalha-naval.jpg`. It has 8 ships covering 19 cells:

| Ship (PT)    | Size | Count | Shape |
| ------------ | ---- | ----- | ----- |
| Submarino    | 1    | 3     | single cell |
| Destroyer    | 2    | 2     | line |
| Cruzador     | 3    | 1     | line |
| Encouraçado  | 4    | 1     | line |
| Porta-Aviões | 5    | 1     | T: a row of 3 cells with 2 cells extending from its middle |

- Line ships can be placed horizontally or vertically. The T-shaped Porta-Aviões can be placed in any of its 4 rotations. Ships may not overlap or touch, not even diagonally.


## Conventions

- golangci-lint uses gocyclo (max complexity 15), goconst, dupl, gosec, stylecheck, and misspell with the US locale (so comments and identifiers are checked as English). goimports groups `github.com/marcopollivier/battleship` imports as local.
- The README, Makefile messages and commit messages follow different language conventions: the docs and Makefile output are in Portuguese, while commits use English Conventional Commits (`feat:`, `fix:`).
