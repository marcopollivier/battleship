# 001 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — Criar `internal/board/coord.go` com `Size`, `Coord` e `String()`. Testar a forma canônica das 100 casas (`A1`…`J10`) e o retorno `Coord(col,row)` para casas fora do tabuleiro. (cobre a parte de forma canônica do CA-4)
- [x] **T2** — `ParseCoord` para entradas válidas: `A1`, `J10`, `C7`, minúsculas e espaços nas pontas. Incluir o teste de ida e volta das 100 casas. (cobre CA-1, CA-2, CA-3, CA-4, CA-8, CA-9)
- [x] **T3** — Validação de formato com `ErrInvalidFormat`: texto vazio, `A`, `1`, `AA1`, `A1B`, `A-1`, `@1`, `Ç1`, `A 1`, `1A`, `A01`. Testar que a mensagem inclui a entrada original. (cobre CA-7, CA-10)
- [x] **T4** — Validação de limites com `ErrOutOfBounds`: `K1`, `Z5`, `A0`, `A11`, `B99` e um número com 20 dígitos, que não pode causar overflow. (cobre CA-5, CA-6)
