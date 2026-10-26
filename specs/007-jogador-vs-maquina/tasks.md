# 007 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — Criar `internal/board/machine.go` com `RandomShot` e `ErrNoCellsLeft`. (cobre CA-4)
- [x] **T2** — Refatorar o `soloGame` para dizer o resultado de cada disparo e receber o texto da desistência, sem mudar o modo solo (os testes da 006 passam sem alteração). (preparação)
- [x] **T3** — Criar `internal/terminal/machine.go` com `PlayVsMachine`: tela inteira, turno do jogador com tiro extra depois de acerto, e vitória do jogador. (cobre CA-1, CA-2, CA-5 no terminal)
- [x] **T4** — Turno da máquina com as linhas de disparo, volta da tela inteira e vitória da máquina. (cobre CA-3, CA-6 no terminal)
- [ ] **T5** — Desistência contra a máquina revelando a frota inimiga. (cobre CA-7)
- [ ] **T6** — `cmd`: opção `--modo`, dois tabuleiros com o mesmo `*rand.Rand`, `PlayVsMachine` como padrão; README ("Como jogar") e `CLAUDE.md`. (cobre CA-1, CA-5, CA-6, CA-8, CA-9, CA-10 no cmd)
