---
description: Quebra um plano aprovado em tarefas (fase 3)
argument-hint: <NNN>
---

Gere as tarefas da spec $ARGUMENTS.

1. Abra `specs/$ARGUMENTS-*/plan.md`. Se o status não for `Aprovado`, pare.
2. Copie `specs/_templates/tasks.md` para a pasta da spec e escreva tarefas pequenas, cada uma resultando em um commit e indicando os critérios de aceitação que cobre.
3. Todo critério de aceitação da spec precisa ser coberto por alguma tarefa. Se passar de ~8 tarefas, sugira dividir a spec.
4. Apresente as tarefas e espere aprovação.
