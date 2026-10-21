---
description: Cria uma nova spec SDD (fase 1)
argument-hint: <feature ou número do roadmap>
---

Crie a spec da feature: $ARGUMENTS

1. Leia `specs/README.md`, `specs/constitution.md` e `specs/ROADMAP.md`.
2. Use o número do roadmap se a feature estiver lá; senão, o próximo `NNN` livre (e adicione a linha no roadmap).
3. Crie o branch `feat/NNN-slug` a partir do `master` atualizado (veja "Branches e PRs" em `specs/README.md`).
4. Copie `specs/_templates/spec.md` para `specs/NNN-slug/spec.md` e preencha Problema, Regras, Critérios de aceitação e Fora de escopo, baseando-se nas regras do `README.md` do projeto.
5. Liste em "Questões em aberto" tudo que o README não decide. Não invente respostas.
6. Mantenha o status `Rascunho`, atualize o status no roadmap para `Spec` e apresente as questões em aberto para a pessoa decidir. Não escreva plano nem código. Faça commit da spec, dê push e abra um PR draft com `gh pr create --draft` linkando a spec.
