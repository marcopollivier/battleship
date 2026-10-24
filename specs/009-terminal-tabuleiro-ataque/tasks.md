# 009 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Rascunho

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [ ] **T1** — Criar `internal/terminal/messages.go` com `ShotMessage`. (cobre CA-4)
- [ ] **T2** — Extrair `writeBoard` e fazer o `RenderBoard` mostrar os disparos recebidos (`X`, `o`) com a legenda nova, atualizando os testes de saída exata da 004. (cobre CA-6, CA-7)
- [ ] **T3** — Criar `internal/terminal/attack.go` com `AttackLegend` e `RenderAttackBoard` para casas não atingidas, água e acerto. (cobre CA-1, CA-2, CA-5)
- [ ] **T4** — Mostrar a letra do tipo nas casas de navios afundados no tabuleiro de ataque. (cobre CA-3)
