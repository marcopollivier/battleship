# 011 — `make run` com opções

- **Status:** Em andamento
- **Roadmap:** #011 (onda 6)
- **Depende de:** 006 (modo solo)

## Problema

O README manda começar o jogo com `make run`, mas as opções (`--frota`, `--semente`, `--modo`) só funcionam com `go run .`, porque o `make run` não repassa nada ao programa. Quem quer repetir uma partida com a semente que apareceu no fim precisa trocar de comando. A 006 sugeriu aceitar `make run ARGS="--semente 42"`.

## Regras

- `make run ARGS="..."` passa o texto de `ARGS` ao jogo como opções, do mesmo jeito que `go run . ...`. Exemplo: `make run ARGS="--modo solo --semente 42"`.
- Sem `ARGS`, o `make run` continua igual: começa uma partida contra a máquina com a frota da cartela e uma semente nova.
- As opções inválidas mostram as mesmas mensagens de `go run .` (por exemplo `Modo inválido: use maquina ou solo.`).
- O `make help` mostra, na linha do `run`, que ele aceita `ARGS`: `Executa o jogo (opções: make run ARGS="--modo solo --semente 42")`.
- Quando o jogo termina com erro (código 1 ou 2), o `make` mostra a sua linha de erro, como `make: *** [run] Error 2`, e também termina com erro (questão 3).
- O README ("Como jogar") e o `CLAUDE.md` passam a mostrar as opções com `make run ARGS="..."`, mantendo `go run .` como alternativa.

## Critérios de aceitação

- **CA-1:** Dado `make run ARGS="--modo solo --semente 42"` com a entrada `sair`, quando ele roda, então a saída do jogo é igual à de `go run . --modo solo --semente 42` com a mesma entrada.
- **CA-2:** Dado `make run` sem `ARGS`, quando ele roda, então começa a partida contra a máquina, como hoje.
- **CA-3:** Dado `make run ARGS="--modo xyz"`, quando ele roda, então aparece `Modo inválido: use maquina ou solo.`, seguido da linha de erro do `make`, e o `make` termina com erro.
- **CA-4:** Dado `make help`, quando ele roda, então a linha do `run` cita `ARGS`.

## Fora de escopo

- Repassar opções ao `make build` ou a outros alvos.
- Atalhos como `make solo` ou `make run SEMENTE=42` (questão 1).
- Esconder a linha `🚀 Executando...` que o `make run` mostra antes do jogo, ou a linha de erro do `make`.

## Questões em aberto

Todas decididas em 2026-10-26:

- [x] **1. Forma:** uma variável genérica, `make run ARGS="--semente 42"`. A alternativa era uma variável por opção (`make run SEMENTE=42`).
- [x] **2. Ajuda:** a descrição do `run` no `make help` vira `Executa o jogo (opções: make run ARGS="--modo solo --semente 42")`.
- [x] **3. Erro do make:** a linha de erro do `make` continua aparecendo, e os critérios conferem que ela aparece.
- [x] **4. Testes:** os critérios são conferidos rodando os comandos à mão, e o resultado fica no encerramento. É uma exceção ao princípio 3 da constituição (teste primeiro), justificada no plano.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
