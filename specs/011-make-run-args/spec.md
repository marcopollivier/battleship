# 011 — `make run` com opções

- **Status:** Concluída
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

Concluída em 2026-10-26, nas tarefas T1 e T2. O código Go não mudou: a T1 trocou o alvo `run` do `Makefile` (`go run $(MAIN_PATH) $(ARGS)` e a nova descrição no `make help`), e a T2 atualizou o README ("Como jogar") e o `CLAUDE.md`. Como decidido na questão 4, os critérios foram conferidos à mão antes do commit da T1, com os comandos do plano:

- **CA-1:** `printf 'sair\n' | make run ARGS="--modo solo --semente 42"` e `printf 'sair\n' | go run . --modo solo --semente 42`, cada um gravado num arquivo. Tirando a linha `🚀 Executando...`, o `diff` não mostrou diferença: as duas saídas trazem o tabuleiro vazio, `Você desistiu depois de 0 disparos.`, a frota sorteada e `Semente: 42`. Os dois terminaram com código 0.
- **CA-2:** `printf 'sair\n' | make run` começou com as instruções do posicionamento (`Posicione a sua frota. ...`), pediu o Porta-Aviões, mostrou `Você desistiu antes da partida.` e uma semente nova, com código 0.
- **CA-3:** `make run ARGS="--modo xyz"; echo $?` mostrou:

  ```
  🚀 Executando...
  Modo inválido: use maquina ou solo.
  exit status 2
  make: *** [run] Error 1
  2
  ```

- **CA-4:** `make help | grep run` mostrou `run             Executa o jogo (opções: make run ARGS="--modo solo --semente 42")`.

O `make pre-commit` passou nas duas tarefas, sem mudança nos testes em Go.

**Desvios em relação ao plano:**

- A linha de erro do `make` é `make: *** [run] Error 1`, e não `Error 2` como no exemplo da spec. O `go run` não repassa o código do programa: ele escreve `exit status 2` e termina com 1. O `make` termina com 2 de qualquer jeito, então o CA-3 vale. A linha `exit status 2` também aparece em `go run . --modo xyz`, então as mensagens continuam as mesmas das duas formas.
- No README, a frase das opções ficou "As opções vão em `ARGS`, como em `make run ARGS="..."`, ou depois de `go run .`, como em `go run . ...`", para caber na frase que já existia antes da tabela de opções.

**Para o roadmap:**

- Se o código de saída do jogo importar para quem usa o `make` (por exemplo, distinguir opção inválida de erro na frota), o `run` pode compilar o binário e rodá-lo, no lugar do `go run`, que troca o código 2 por 1.
