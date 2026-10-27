# 014 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — `Board.RemoveLast`. (cobre CA-1, CA-2 no domínio)
- [x] **T2** — `Fits`: busca exata com poda de navios iguais e limite de tentativas; medir o pior caso com a frota da cartela. (cobre CA-4, CA-5 no domínio)
- [x] **T3** — Comando `desfazer` e a linha nova das instruções. (cobre CA-1, CA-2, CA-3, CA-8)
- [x] **T4** — Aviso de que o resto não cabe, antes de cada tabuleiro do posicionamento. (cobre CA-4, CA-5, CA-6)
- [ ] **T5** — `cmd`: mesma semente com `desfazer`; README e `CLAUDE.md`. (cobre CA-7)
