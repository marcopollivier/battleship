# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

A terminal Battleship ("Batalha Naval") game in Go (module `github.com/marcopollivier/battleship`, Go 1.24, no external dependencies). It is at an early stage: today it only renders an empty 10x10 board. The README (in Portuguese) is the game spec.

## Commands

All are Makefile targets (`make help` lists them):

- `make run`: `go run ./main.go`
- `make build`: builds the binary to `bin/battleship`
- `make test`: `go test -v ./...`
- Run a single test: `go test -v -run TestName ./cmd/`
- `make lint`: `golangci-lint run` (config in `.golangci.yml`)
- `make fmt` / `make vet` / `make mod-tidy`
- `make pre-commit`: fmt, vet, lint, sec and test together
- `make setup`: installs pre-commit (via Homebrew) and golangci-lint, then installs the git hooks

The pre-commit hooks run go-fmt, go-mod-tidy, go unit tests and golangci-lint on every commit, so code must pass lint and tests before it can be committed.

## Spec-driven workflow

Features are built through the spec-driven development (SDD) process in `specs/README.md`, and every plan must follow the principles in `specs/constitution.md`. The order is spec, then plan, then tasks, then implementation, with explicit user approval needed before moving from one phase to the next. Don't write production code for a feature that has no approved `specs/NNN-slug/tasks.md`. If implementation contradicts the spec, stop and propose a spec change. Each phase has a slash command: `/spec`, `/plan`, `/tasks` and `/implement`. The feature backlog and the open game-rule questions are in `specs/ROADMAP.md`.

Never commit directly to `master`. Each context gets its own branch and PR: `feat/NNN-slug` for a spec, and `docs/`, `chore/` or `fix/` for anything else. The full convention is under "Branches e PRs" in `specs/README.md`. Work that falls outside the current PR's context goes on a separate branch, or into the roadmap.

## Architecture

- `main.go` is a thin entrypoint that only calls `cmd.Init()`.
- `cmd/init.go` holds the game logic for now: it builds the board as a `[10][10]string` filled with `~` (water) and prints it with columns `A`–`J` and rows `1`–`10`.

## Game rules (from README)

- The board is 10x10: columns A–J, rows 1–10. Shots are written as coordinates such as `A1`, `B3`, `J10`.
- A player who hits a ship shoots again. A miss ("tiro na água") passes the turn to the other player.
- A ship is sunk when all of its cells are hit. The first player to sink every enemy ship wins.
- Ships can be placed horizontally or vertically. The default fleet is:

| Ship (PT)    | Size | Count |
| ------------ | ---- | ----- |
| Rebocador    | 1    | 4     |
| Submarino    | 2    | 3     |
| Destroyer    | 3    | 2     |
| Encouraçado  | 4    | 1     |
| Porta-Aviões | 5    | 1     |

Planned features (README TODO): attacking a randomly generated board, placing your own ships and playing against the computer, and choosing how many ships to play with.

## Conventions

- golangci-lint uses gocyclo (max complexity 15), goconst, dupl, gosec, stylecheck, and misspell with the US locale (so comments and identifiers are checked as English). goimports groups `github.com/marcopollivier/battleship` imports as local.
- The README, Makefile messages and commit messages follow different language conventions: the docs and Makefile output are in Portuguese, while commits use English Conventional Commits (`feat:`, `fix:`).
