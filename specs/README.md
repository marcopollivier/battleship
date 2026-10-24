# Spec-Driven Development (SDD)

Toda funcionalidade do jogo nasce como uma spec e só vira código depois que a spec e o plano forem aprovados. O objetivo é iterar em incrementos pequenos, cada um com critérios de aceitação verificáveis por testes.

## Estrutura

```
specs/
├── README.md          # este processo
├── constitution.md    # princípios inegociáveis do projeto
├── ROADMAP.md         # backlog ordenado de features
├── _templates/        # modelos de spec, plan e tasks
└── NNN-slug/          # uma pasta por feature (ex.: 001-coordenadas)
    ├── spec.md        # O QUÊ e POR QUÊ
    ├── plan.md        # COMO
    └── tasks.md       # passos executáveis
```

A numeração `NNN` é sequencial e nunca reaproveitada.

## Fluxo

Cada fase termina num **gate**: a próxima fase só começa depois que a pessoa aprovar explicitamente.

| Fase | Artefato | Conteúdo | Gate |
| ---- | -------- | -------- | ---- |
| 1. Especificar | `spec.md` | Problema, regras, critérios de aceitação (Dado/Quando/Então), fora de escopo, questões em aberto | Nenhuma questão em aberto; status `Aprovada` |
| 2. Planejar | `plan.md` | Pacotes, tipos e funções públicas, decisões e alternativas, estratégia de testes, checagem da constituição | Plano aprovado |
| 3. Quebrar | `tasks.md` | Checklist de tarefas pequenas, cada uma com o critério de aceitação que cobre | Tarefas aprovadas |
| 4. Implementar | código + testes | Uma tarefa por vez: teste falhando → código → `make pre-commit` → commit | Todas as tarefas marcadas |
| 5. Fechar | `spec.md` | Status `Concluída`, registro de desvios em relação ao plano | — |

### Regras do fluxo

- A spec descreve comportamento, não implementação. Nomes de tipos e pacotes ficam no plano.
- Todo critério de aceitação vira pelo menos um teste. As tarefas indicam qual critério cobrem (`CA-1`, `CA-2`...).
- Se a implementação revelar que a spec está errada ou incompleta, pare, atualize a spec e peça nova aprovação antes de seguir.
- Mudanças fora do escopo da spec atual vão para o `ROADMAP.md`, não para o código.
- **Plano e tarefas juntos:** quando a spec aprovada não deixar decisão técnica relevante em aberto, o plano e as tarefas são apresentados juntos e aprovados numa única rodada. Cada artefato continua com seu próprio status. Se a pessoa pedir mudanças no plano, as tarefas são refeitas antes da aprovação.

## Branches e PRs

Cada contexto tem um branch e um PR próprios. Nada é commitado direto no `master`.

| Contexto | Branch | PR contém |
| -------- | ------ | --------- |
| Uma spec | `feat/NNN-slug` | `spec.md`, `plan.md`, `tasks.md` e a implementação, um commit por tarefa |
| Processo, docs e tooling | `docs/<assunto>` ou `chore/<assunto>` | Só a mudança daquele assunto |
| Correção fora de uma spec | `fix/<assunto>` | A correção e o teste que a reproduz |

- O PR de uma spec é aberto como draft na fase 1 e sai de draft quando todas as tarefas estiverem marcadas.
- A descrição do PR linka a spec e lista os critérios de aceitação.
- Algo que surgir fora do contexto do PR vai para outro branch ou para o `ROADMAP.md`.

## Trabalho em paralelo (ondas)

Specs sem dependência entre si são tocadas em paralelo, em **ondas**. O `ROADMAP.md` indica a onda e as dependências de cada feature.

1. **Specs em lote:** as specs da onda são escritas juntas, cada uma no seu branch e PR draft. As questões em aberto de todas são apresentadas numa única rodada de decisões.
2. **Plano e tarefas:** seguem a regra de aprovação conjunta acima, uma spec por vez ou todas juntas.
3. **Implementação:** cada spec é implementada por um agente separado, num git worktree próprio, seguindo `/implement`. Quem coordena revisa cada PR antes de pedir o merge.
4. **Ordem de merge:** os PRs de uma onda são independentes. Se dois tocarem o mesmo arquivo, o segundo faz merge do `master` atualizado antes de ser mergeado.

Para evitar conflitos entre PRs paralelos, **os PRs de spec não editam o `ROADMAP.md`** durante uma onda. Quem coordena atualiza o roadmap num PR `docs/` no início da onda (status `Spec`) e outro no fim (status `Concluída`).

## Status de uma spec

`Rascunho` → `Aprovada` → `Em andamento` → `Concluída`

## Comandos do Claude Code

| Comando | O que faz |
| ------- | --------- |
| `/spec <feature>` | Cria `specs/NNN-slug/spec.md` a partir do template e levanta questões em aberto |
| `/plan NNN` | Escreve o `plan.md` de uma spec aprovada e, quando cabe a aprovação conjunta, também o `tasks.md` |
| `/tasks NNN` | Gera o `tasks.md` a partir do plano aprovado |
| `/implement NNN` | Executa a próxima tarefa pendente com TDD |
