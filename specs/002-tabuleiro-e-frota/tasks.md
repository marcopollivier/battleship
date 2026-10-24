# 002 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — Criar `internal/board/ship.go` com `ShipType` (`Size`, `String`), `FleetEntry` e `DefaultFleet`. Testar a contagem por tipo e o total de 19 casas. (cobre CA-1)
- [x] **T2** — `Orientation` e `ShipCells` para os navios em linha: `Horizontal` e `Vertical`, com `ErrInvalidOrientation` para `Up`/`Down`/`Left`/`Right`. Testar o Cruzador em `A1`, o Submarino em `E5` e o Cruzador para `Down`. (cobre CA-2, CA-3 e a parte de linha do CA-11)
- [x] **T3** — `ShipCells` para o Porta-Aviões: tabela de deslocamentos das 4 orientações, com `ErrInvalidOrientation` para `Horizontal`/`Vertical`. Testar as casas a partir de `E5`. (cobre CA-5 e a parte do T do CA-11)
- [x] **T4** — Criar `internal/board/board.go` com `NewBoard`, `Place` e `Ships`, validando orientação e limites (`ErrOutOfBounds`) sem alterar o tabuleiro em caso de erro. Testar o Encouraçado em `J7`/`J8`/`H1` e o Porta-Aviões em `E9 Down` e `A5 Up`. Testar também que o `ErrInvalidOrientation` chega por `Place`. (cobre CA-4, CA-6 e CA-11 via `Place`)
- [x] **T5** — Rejeitar sobreposição com `ErrOverlap`. Testar o Destroyer em `B1` sobre o Cruzador em `A1`, conferindo que `Ships()` continua com 1 navio. (cobre CA-7)
- [x] **T6** — Rejeitar navios encostados, inclusive na diagonal, com `ErrAdjacent`. Testar `D1`, `A2` e `D2` rejeitados e `E1` aceito. (cobre CA-10)
- [ ] **T7** — Limite por tipo com `ErrFleetLimit` (incluindo tipos fora da frota) e `Complete`. Testar o quarto Submarino e uma disposição fixa dos 8 navios, com `Complete` falso no 7º e verdadeiro no 8º. (cobre CA-8, CA-9)
