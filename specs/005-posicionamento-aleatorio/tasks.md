# 005 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — `ValidOrientations` em `ship.go`. (base para as demais; testada por `TestValidOrientations`)
- [x] **T2** — Criar `internal/board/random.go` com `RandomBoard`, posicionando do maior navio para o menor com tentativas e recomeços, e o verificador `assertValidBoard` nos testes. (cobre CA-1, CA-2)
- [x] **T3** — Testes de reprodutibilidade e de orientações. Se a implementação da T2 já satisfizer os critérios, os testes passam de primeira e ficam como proteção contra regressão; isso é registrado no commit. (cobre CA-3, CA-4, CA-5)
- [ ] **T4** — Limites de tentativas e recomeços, com `ErrFleetDoesNotFit`. Se a T2 já tiver os limites, a tarefa fica só com o teste do CA-6. (cobre CA-6)
