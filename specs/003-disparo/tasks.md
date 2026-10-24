# 003 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — Criar `internal/board/shot.go` com `ShotResult`, `Shot` e `Fire` para água e acerto, e o campo `shots` no `Board`. Criar o helper `exampleBoard` nos testes. (cobre CA-1, CA-2)
- [x] **T2** — Detectar navio afundado em `Fire`, preenchendo `Shot.Ship`. (cobre CA-3, CA-4)
- [x] **T3** — `CellState` e `Board.Cell`. (cobre CA-5)
- [x] **T4** — Rejeições em `Fire`: tabuleiro incompleto, casa fora do tabuleiro e tiro repetido, sem alterar o tabuleiro. (cobre CA-6, CA-8, CA-9)
- [x] **T5** — `Board.AllSunk`. (cobre CA-7)
