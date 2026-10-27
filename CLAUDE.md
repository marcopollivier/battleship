# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

A terminal Battleship ("Batalha Naval") game in Go (module `github.com/marcopollivier/battleship`, Go 1.24, no external dependencies). Today it plays against the machine, after the player places their own fleet, or a solo game (`--modo solo`) where the player only attacks a random board. The README (in Portuguese) is the game spec.

## Commands

All are Makefile targets (`make help` lists them):

- `make run`: `go run ./main.go`, which starts a game against the machine. Options go in `ARGS`, as in `make run ARGS="--modo solo --frota S=2,P=1 --semente 42 --cores nao"`, or after `go run .`
- `make build`: builds the binary to `bin/battleship`
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

Two layers, and the domain never imports the terminal:

- **`internal/board`: the game domain.** It has no I/O and its names are in English.
  - `coord.go`: `Coord` (zero-based `Col`, `Row`) and `ParseCoord`, which turns `"A1"` into a `Coord`.
  - `ship.go`: `ShipType`, `Orientation`, `DefaultFleet`, `ShipCells` and `ValidOrientations`.
  - `board.go`: `Board` with `Place`, which validates fleet limit, orientation, bounds, overlap and adjacency in that order and leaves the board unchanged on error; `Remaining` (ships still to place, largest first) and `Clear`.
  - `shot.go`: `Fire`, `Cell` and `AllSunk`.
  - `random.go`: `RandomBoard` and `FillRandom` (places only the missing ships), which take an injected `*rand.Rand`.
  - `machine.go`: `RandomShot`, the machine's next shot, drawn among the cells not shot yet.
  - `fleet.go`: `ValidateFleet`, which rejects negative counts, empty fleets and fleets that do not fit on the board.
  - Errors are sentinel values wrapped with `%w`; check them with `errors.Is`.
- **`internal/terminal`: presentation.** It reads the domain and writes to an `io.Writer`. Text shown to the player is in Portuguese.
  - `names.go`: `ShipName` and `OrientationName`.
  - `color.go`: `Style` (`Plain`, the zero value, or `Colored`) and `paint`, which wraps a board symbol in its ANSI color code and a reset when the style is `Colored`.
  - `render.go`: `RenderBoard(w, board, style)`, the player's own board with ships and received shots.
  - `attack.go`: `RenderAttackBoard(w, board, style)`, the opponent's board showing only shots fired, and sunk ships.
  - Both renderers paint each cell while building the grid; `writeBoard` adds the header, row numbers and legend, which stay uncolored.
  - `messages.go`: `ShotMessage`, the result of a shot.
  - `fleet.go`: `ParseFleet`, which reads the `--frota` text (such as `S=2,D=1,P=1`), and `FleetErrorMessage`.
  - `solo.go`: `PlaySolo(in, out, board, style)`, the solo game loop: reads one shot per line until every ship is sunk or the player types `sair` or closes the input. Its `playerAttack` type is the player's turn, shared with `machine.go`, and keeps the `Style` every board of the game is drawn with.
  - `machine.go`: `PlayVsMachine(in, out, player, enemy, rng, style)`, the game against the machine: whoever hits shoots again. If `player` is not complete, it first runs the placing phase.
  - `placement.go`: the placing phase (`casa orientação` lines, `aleatorio`, `recomecar`, `sair`). It runs inside `PlayVsMachine` so one `bufio.Scanner` reads all the input.
- **`cmd/` and `main.go`: wiring.** `cmd.Run(args, stdin, stdout, stderr) int` parses `--modo`, `--cores` (`sim`, the default, or `nao`; colors stay on even when the output is not a terminal), `--frota` and `--semente`, checking the mode, then the colors, then the fleet. One `*rand.Rand` seeded by the seed (a new one each run unless given) draws the machine's board, then the ships the player draws with `aleatorio`, then the machine's shots; `Run` runs `PlayVsMachine` (or `PlaySolo`) and prints the seed. It returns the exit code: 0 at the end of a game or on `-h`, 1 on a fleet error, and 2 on an invalid option. `main.go` passes it to `os.Exit`.

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
