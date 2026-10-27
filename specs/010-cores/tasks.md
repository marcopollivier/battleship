# 010 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — `Style`, `Plain` e `Colored` em `color.go`; `RenderBoard` e `RenderAttackBoard` recebem o `Style`, e todas as chamadas atuais passam `Plain`, sem mudar nenhum resultado. (preparação; cobre CA-4 no terminal)
- [x] **T2** — Cores do tabuleiro próprio. (cobre CA-1, CA-3)
- [ ] **T3** — Cores do tabuleiro de ataque. (cobre CA-2, CA-3)
- [ ] **T4** — `PlaySolo` e `PlayVsMachine`, incluindo o posicionamento, recebem o `Style` e desenham todos os tabuleiros com ele. (cobre CA-7)
- [ ] **T5** — `cmd`: opção `--cores`, cores ligadas por padrão; README (tabela de opções) e `CLAUDE.md`. (cobre CA-4, CA-5, CA-6)
