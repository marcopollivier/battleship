# 012 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — `Board.Remaining` e `Board.Clear`. (cobre CA-7 no domínio; prepara CA-2)
- [x] **T2** — `FillRandom`, e o `RandomBoard` passando a usá-lo sem mudar os tabuleiros sorteados. (cobre CA-5 e CA-6 no domínio)
- [x] **T3** — `parsePlacement`: casa, orientação pelo nome ou pela primeira letra, Submarino sem orientação. (cobre CA-3, CA-4 formato e orientação)
- [x] **T4** — Fase de posicionamento no `PlayVsMachine`: instruções, tabuleiro e pedido de cada navio, `Frota posicionada!` e início da partida. (cobre CA-1, CA-2)
- [x] **T5** — Mensagens de posicionamento inválido, pedindo o mesmo navio. (cobre CA-4)
- [x] **T6** — Comandos `aleatorio`, `recomecar` e `sair`/fim da entrada. (cobre CA-5 a CA-8 no terminal)
- [x] **T7** — `cmd`: tabuleiro da máquina sorteado primeiro e tabuleiro vazio para o jogador; testes do `cmd` da 007 com `aleatorio`; README e `CLAUDE.md`. (cobre CA-8 a CA-10 no cmd)
