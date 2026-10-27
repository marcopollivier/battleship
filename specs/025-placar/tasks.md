# 025 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — Pacote `internal/score`: `Record`, `Board`, `AddWin`, `AddLoss`, `Load`, `Save` e `DefaultPath`. (cobre CA-4 e CA-7 no pacote)
- [x] **T2** — `PlayVsMachine` devolve o `Result`. (preparação)
- [x] **T3** — `ScoreLine` e o registro no fim da partida no `cmd`, com o caminho por parâmetro. (cobre CA-1, CA-2, CA-3, CA-8)
- [ ] **T4** — `--placar` e `--zerar-placar`. (cobre CA-5, CA-6)
- [ ] **T5** — Aviso quando o placar não pode ser lido ou gravado; README e `CLAUDE.md`. (cobre CA-7)
