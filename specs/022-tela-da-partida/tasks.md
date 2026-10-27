# 022 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Rascunho

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [ ] **T1** — Grades em linhas no `render.go` e no `attack.go`, sem mudar a saída do `RenderBoard` e do `RenderAttackBoard`. (preparação)
- [ ] **T2** — `screen.go`: `sideBySide`, `clearScreen` e os títulos. (cobre CA-1 na função)
- [ ] **T3** — Tipo `Display` no lugar do `Style` no `PlaySolo` e no `PlayVsMachine`, sem mudar a saída. (preparação)
- [ ] **T4** — Tela da partida contra a máquina: tabuleiros lado a lado, mensagens do turno, `Seu disparo em ...`, erros abaixo do pedido, limpeza. (cobre CA-1 a CA-5)
- [ ] **T5** — Telas finais: vitória, derrota e desistência com os dois tabuleiros revelados. (cobre CA-6)
- [ ] **T6** — Modo solo e posicionamento com a mesma regra de tela e limpeza. (cobre CA-7, CA-8)
- [ ] **T7** — `cmd`: `--limpar auto|sim|nao`, `chooseClear`; README e `CLAUDE.md`. (cobre CA-5, CA-9, CA-10)
