# 008 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — Criar `internal/board/fleet.go` com `ValidateFleet`, `ErrNegativeCount` e `ErrEmptyFleet`. (cobre CA-5, CA-6)
- [x] **T2** — Checagem de "cabe" em `ValidateFleet` via `RandomBoard` com semente fixa. (cobre CA-7, CA-8)
- [x] **T3** — Criar `internal/terminal/fleet.go` com `ParseFleet` para textos válidos. (cobre CA-1, CA-2)
- [ ] **T4** — Recusas de formato em `ParseFleet` com `ErrInvalidFleetFormat`. (cobre CA-4)
- [ ] **T5** — `FleetErrorMessage`. (cobre CA-9)
- [ ] **T6** — `cmd.Run(args, stdout, stderr) int` com a opção `--frota`, e o `main.go` usando o código de saída. Remover `cmd.Init()`. (cobre CA-3, CA-10)
