---
description: Implementa a próxima tarefa pendente de uma spec (fase 4)
argument-hint: <NNN>
---

Implemente a próxima tarefa da spec $ARGUMENTS.

1. Abra `spec.md`, `plan.md` e `tasks.md` em `specs/$ARGUMENTS-*/`. Se `tasks.md` não estiver `Aprovado`, pare. Na primeira tarefa, mude o status da spec para `Em andamento`.
2. Pegue a primeira tarefa não marcada. Faça apenas ela.
3. Escreva o teste primeiro e rode `go test ./...` para confirmar que falha pelo motivo certo.
4. Implemente o mínimo para passar, depois rode `make pre-commit`.
5. Marque a tarefa em `tasks.md` e faça um commit em inglês no padrão Conventional Commits, citando a spec (ex.: `feat(board): parse coordinates (spec 001 T1)`) e dê push no branch.
6. Se a spec ou o plano se mostrarem errados, pare e proponha a correção em vez de improvisar.
7. Se era a última tarefa, preencha "Encerramento" na spec, mude o status para `Concluída` e atualize o roadmap. Dê push e tire o PR de draft com `gh pr ready`.
