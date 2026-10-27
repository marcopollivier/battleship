# 013 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — `ShotStrategy` e `SmartShot` na procura: sorteia casas não atingidas fora do entorno dos navios afundados. (cobre CA-1, CA-5)
- [x] **T2** — Cerco pelas vizinhas de um acerto pendente. (cobre CA-2)
- [x] **T3** — Cerco pelas pontas da linha, com volta às vizinhas quando a linha acaba. (cobre CA-3, CA-4)
- [x] **T4** — Testes de partida inteira e de conhecimento público: nunca repete, termina, e precisa de menos tiros que o `RandomShot`. (cobre CA-6)
- [ ] **T5** — `PlayVsMachine` recebe a `ShotStrategy`; os testes atuais passam `RandomShot`. (cobre CA-8 no terminal)
- [ ] **T6** — `cmd`: opção `--nivel`, `dificil` por padrão; README (tabela de opções) e `CLAUDE.md`. (cobre CA-7, CA-8, CA-9)
