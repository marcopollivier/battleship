# 011 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

A mudança é só no `Makefile` e na documentação. O alvo `run` passa a acrescentar `$(ARGS)` ao `go run`. Como o `make` substitui uma variável vazia por nada, `make run` sem `ARGS` continua igual. O código Go não muda.

## Pacotes e API

`Makefile`:

```make
run: ## Executa o jogo (opções: make run ARGS="--modo solo --semente 42")
	@echo "🚀 Executando..."
	@go run $(MAIN_PATH) $(ARGS)
```

- O `make help` lê o texto depois de `##`. O awk do `help` separa os campos em `:.*?## `, e o texto novo não tem `## `, então a linha aparece inteira.
- `ARGS` não ganha valor padrão no `Makefile`, para não aparecer como opção vazia.

Documentação:

- **README, "Como jogar":** a frase das opções passa a citar as duas formas: `make run ARGS="--frota S=2,P=1 --semente 42"` ou `go run . --frota S=2,P=1 --semente 42`.
- **`CLAUDE.md`, "Commands":** a linha do `make run` mostra `ARGS`.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| `$(ARGS)` direto no `go run` | Um script que monta as opções | É uma linha, e o `make` já trata as aspas de `ARGS="..."` como um valor só. |
| Manter o `@echo "🚀 Executando..."` | Tirar a linha | Fora de escopo na spec; os outros alvos também anunciam o que fazem. |
| Conferir à mão, sem teste automático | Teste em Go que roda `make`; um alvo `make test-run` | Decisão da questão 4 da spec. Um teste com `make` dependeria do `make` e do `go` no `PATH` do teste e repetiria o que o teste do `cmd` já confere. |

## Estratégia de testes

Conferência à mão, registrada no encerramento com os comandos e o resultado (exceção ao princípio 3, abaixo):

| Critério | Conferência |
| -------- | ----------- |
| CA-1 | `printf 'sair\n' \| make run ARGS="--modo solo --semente 42"` e `printf 'sair\n' \| go run . --modo solo --semente 42`: a saída do jogo (sem a linha `🚀 Executando...`) é igual, comparada com `diff` |
| CA-2 | `printf 'sair\n' \| make run`: começa com as instruções do posicionamento (012) |
| CA-3 | `make run ARGS="--modo xyz"`: mostra `Modo inválido: use maquina ou solo.` e a linha de erro do `make`, e `echo $?` mostra um código diferente de 0 |
| CA-4 | `make help \| grep run`: a linha cita `ARGS` |

Os testes em Go continuam passando sem mudança (`make pre-commit`).

## Checagem da constituição

- [x] **Domínio separado da interface:** nenhum código Go muda.
- [x] **Aleatoriedade injetada:** não se aplica.
- [ ] **Testes antes do código:** exceção, veja abaixo.
- [x] **Sem dependências externas:** só `make`, que o projeto já usa.
- [x] **Erros como valores:** não se aplica; os erros de opção continuam os do `cmd`.
- [x] **Até ~8 tarefas:** 2.
- [x] **Idiomas:** a descrição do `make help` e a documentação em português.

Exceções e justificativas:

- **Princípio 3 (teste primeiro):** o `Makefile` não tem teste em Go. Os critérios são conferidos à mão antes do commit, com os comandos acima, e o resultado fica no encerramento. Decidido na questão 4 da spec.
