# 011 — `make run` com opções

- **Status:** Rascunho
- **Roadmap:** #011 (onda 6)
- **Depende de:** 006 (modo solo)

## Problema

O README manda começar o jogo com `make run`, mas as opções (`--frota`, `--semente`, `--modo`) só funcionam com `go run .`, porque o `make run` não repassa nada ao programa. Quem quer repetir uma partida com a semente que apareceu no fim precisa trocar de comando. A 006 sugeriu aceitar `make run ARGS="--semente 42"`.

## Regras

- `make run ARGS="..."` passa o texto de `ARGS` ao jogo como opções, do mesmo jeito que `go run . ...`. Exemplo: `make run ARGS="--modo solo --semente 42"`.
- Sem `ARGS`, o `make run` continua igual: começa uma partida contra a máquina com a frota da cartela e uma semente nova.
- As opções inválidas mostram as mesmas mensagens de `go run .` (por exemplo `Modo inválido: use maquina ou solo.`).
- O `make help` mostra, na linha do `run`, que ele aceita `ARGS` (questão 2).
- O README ("Como jogar") e o `CLAUDE.md` passam a mostrar as opções com `make run ARGS="..."`, mantendo `go run .` como alternativa.

## Critérios de aceitação

- **CA-1:** Dado `make run ARGS="--modo solo --semente 42"` com a entrada `sair`, quando ele roda, então a saída do jogo é igual à de `go run . --modo solo --semente 42` com a mesma entrada.
- **CA-2:** Dado `make run` sem `ARGS`, quando ele roda, então começa a partida contra a máquina, como hoje.
- **CA-3:** Dado `make run ARGS="--modo xyz"`, quando ele roda, então aparece `Modo inválido: use maquina ou solo.`.
- **CA-4:** Dado `make help`, quando ele roda, então a linha do `run` cita `ARGS`.

## Fora de escopo

- Repassar opções ao `make build` ou a outros alvos.
- Atalhos como `make solo` ou `make run SEMENTE=42` (questão 1).
- Esconder a linha `🚀 Executando...` que o `make run` mostra antes do jogo.

## Questões em aberto

- [ ] **1. Forma:** uma variável genérica, `make run ARGS="--semente 42"`, como proposto. A alternativa é uma variável por opção (`make run SEMENTE=42 MODO=solo FROTA=S=2`), mais curta de digitar, mas que precisa mudar o `Makefile` a cada opção nova.
- [ ] **2. Ajuda:** a descrição do `run` no `make help` vira `Executa o jogo (opções: make run ARGS="--modo solo --semente 42")`.
- [ ] **3. Erro do make:** quando o jogo termina com código 1 ou 2 (frota ou opção inválida), o `make` acrescenta uma linha como `make: *** [run] Error 2`. Proposta: deixar, porque o código de saída avisa scripts que algo deu errado. A alternativa é o `make run` sempre terminar com sucesso, escondendo essa linha.
- [ ] **4. Testes:** o `Makefile` não tem teste em Go. Proposta: conferir os critérios rodando os comandos à mão e registrar o resultado no encerramento, como exceção ao princípio 3 da constituição (teste primeiro). A alternativa é um teste em Go que roda `make`, que seria lento e dependeria do `make` instalado no CI.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
