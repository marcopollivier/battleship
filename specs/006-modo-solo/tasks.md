# 006 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando → implementação → `make pre-commit` → commit.

- [x] **T1** — Criar `internal/terminal/solo.go` com `PlaySolo`: tabuleiro inicial, pedido de disparo, disparos válidos com mensagem e tabuleiro, e vitória com a contagem. (cobre CA-1, CA-2, CA-4 no terminal)
- [x] **T2** — Mensagens de entrada inválida (`shotErrorMessage`) sem contar o disparo. (cobre CA-3)
- [x] **T3** — Desistência com `sair` e fim da entrada, revelando a frota, e o singular `1 disparo`. (cobre CA-5 no terminal)
- [ ] **T4** — `cmd.Run(args, stdin, stdout, stderr) int`: sortear o tabuleiro, opção `--semente`, chamar `PlaySolo`, escrever `Semente: N`; `main.go` passa `os.Stdin`. (cobre CA-4 e CA-5 no cmd, CA-6, CA-7, CA-8)
- [ ] **T5** — README com uma seção "Como jogar" (`make run`, `--frota`, `--semente`, `sair`) e `CLAUDE.md` com o `PlaySolo` e a nova assinatura do `Run`. (documentação)
