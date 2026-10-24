---
description: Escreve o plano técnico de uma spec aprovada (fase 2)
argument-hint: <NNN>
---

Escreva o plano da spec $ARGUMENTS.

1. Abra `specs/$ARGUMENTS-*/spec.md`. Se o status não for `Aprovada` ou houver questões em aberto sem decisão, pare e diga o que falta.
2. Leia `specs/constitution.md` e o código atual que a feature toca.
3. Copie `specs/_templates/plan.md` para a pasta da spec e preencha todas as seções. Todo critério de aceitação precisa aparecer na tabela de testes.
4. Preencha a checagem da constituição com honestidade; qualquer exceção deve ser justificada.
5. Se o plano não deixar decisão técnica relevante para a pessoa, gere também o `tasks.md` (como em `/tasks`) para aprovação conjunta (veja "Plano e tarefas juntos" em `specs/README.md`). Caso contrário, não gere tarefas.
6. Apresente o plano (e as tarefas, se geradas) e espere aprovação. Não escreva código.
