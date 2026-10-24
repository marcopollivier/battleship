# 004 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — Criar `internal/terminal/names.go` com `ShipName` e `OrientationName`. (cobre CA-4)
- [x] **T2** — Criar `internal/terminal/render.go` com `Legend` e `RenderBoard` para o tabuleiro vazio. (cobre CA-1)
- [x] **T3** — Desenhar os navios com a letra de cada tipo em `RenderBoard`. (cobre CA-2, CA-3, CA-6)
- [x] **T4** — `cmd.Run` e `cmd.Init` usando `RenderBoard`, e o `main.go` tratando o erro. Remover o desenho manual antigo. (cobre CA-5)
