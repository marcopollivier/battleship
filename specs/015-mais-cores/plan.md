# 015 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

1. **Terminal (`internal/terminal`):**
   - o `Style` da 010 passa a pintar também textos inteiros, não só símbolos;
   - os pontos que escrevem as mensagens pintam cada mensagem com a sua cor. Esses pontos são o `playerAttack` (`shoot`, `reject`), o `machineGame` (`machineTurn`, vitória, derrota) e o posicionamento (`placeShip`);
   - o `ShotMessage` público continua devolvendo o texto puro, e a cor é escolhida por uma função interna.
2. **Programa (`cmd/init.go`):**
   - `--cores` aceita `auto`, que passa a ser o padrão;
   - a decisão fica numa função pura, que recebe o valor da opção, se a saída é um terminal e o valor de `NO_COLOR`;
   - o `Run` descobre se o `stdout` é um terminal e lê a variável.

O domínio não muda.

## Pacotes e API

`internal/terminal/color.go`:

- Duas cores novas: `winColor = "\x1b[1;32m"` (verde em negrito) e `warnColor = "\x1b[33m"` (amarelo). As outras são as da 010: `missColor` (ciano) e `hitColor` (vermelho em negrito).
- `paint(text, color)` já serve para textos: envolve o texto inteiro num código e um reset.
- Uma função nova, `shotColor(r board.ShotResult) string`: `Miss` → `missColor`; `Hit` e `Sunk` → `hitColor`.

Onde cada mensagem é pintada:

| Mensagem | Lugar | Cor |
| -------- | ----- | --- |
| `Água.`, `Acerto!`, `Afundou o <navio>!` | `playerAttack.shoot` | `shotColor` |
| resultado da máquina depois de `A máquina atirou em C4: ` | `machineGame.machineTurn` | `shotColor` |
| `Você venceu! ...` e `Você afundou toda a frota em ...!` (solo) | `PlayVsMachine`, `PlaySolo` | `winColor` |
| `A máquina venceu! Ela afundou a sua frota em N disparos.` | `machineWins` | `hitColor`, até o ponto final; ` A frota inimiga era:` fica sem cor, como as outras revelações |
| mensagens de entrada inválida (`shotErrorMessage`, `placementErrorMessage`, sorteio que falhou) | `reject`, `placeShip` | `warnColor` |

- O texto puro das mensagens não muda, e as quebras de linha ficam fora do código de cor.
- A vitória do modo solo também fica verde, porque é a mesma vitória do jogador. A spec só cita `Você venceu!`, e o encerramento registra essa extensão.

`cmd/init.go`:

```go
// chooseStyle returns the board style for the colors option (sim, nao or
// auto, in any case), given whether the output is a terminal and the value
// of NO_COLOR, and false for any other option value.
func chooseStyle(colors string, tty bool, noColor string) (terminal.Style, bool)

// isTerminal reports whether w is a character device, such as a terminal.
func isTerminal(w io.Writer) bool
```

- Com `auto`, `chooseStyle` só devolve `Colored` quando `tty` é verdadeiro e `noColor == ""`.
- `isTerminal` testa se o `w` é um `*os.File` e se o `Stat` tem `os.ModeCharDevice`. O `Run` chama `chooseStyle(*colors, isTerminal(stdout), os.Getenv("NO_COLOR"))`.
- A mensagem de valor inválido passa a ser `Cores inválidas: use sim, nao ou auto.`.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| `chooseStyle` pura, com o terminal e o `NO_COLOR` como parâmetros | Ler `os.Getenv` e o `Stat` dentro dela | A tabela de casos do CA-6 e do CA-7 fica testável sem terminal de verdade e sem mexer no ambiente do teste. |
| Detectar o terminal pelo `os.ModeCharDevice` | `golang.org/x/term` | Só biblioteca padrão (princípio 4). O `/dev/null` também é um dispositivo de caracteres e ligaria as cores, mas ninguém lê essa saída. |
| `ShotMessage` continua sem cor | Receber o `Style` | É usado pelas mensagens da máquina, e a cor é decisão de quem escreve. Mantém a API da 009. |
| Pintar só o resultado nas linhas da máquina | Pintar a linha inteira | Decisão da questão 1. |
| Testes atuais do `cmd` que esperam cores passam `--cores sim` | Fingir um terminal | Com o padrão `auto` e a saída num `strings.Builder`, o padrão vira sem cores. É a regra nova, e o teste deixa explícito que quer cores. |

## Estratégia de testes

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestPlaySolo_ColoredMessages`: água, acerto e afundou com as cores; e a vitória do solo em verde |
| CA-2 | `TestPlayVsMachine_ColoredMachineShots`: só o resultado depois de `: ` tem cor |
| CA-3 | `TestPlayVsMachine_ColoredEnd`: vitória em verde e derrota em vermelho |
| CA-4 | `TestPlaySolo_ColoredMessages` (coordenada inválida, tiro repetido) e `TestPlayVsMachine_PlacementColored` (posição inválida, sorteio que falhou) em amarelo |
| CA-5 | Os testes de partida colorida da 010 passam a conferir também que, tirando as cores, a saída é igual à sem cores, e que pedidos, instruções e títulos não têm `\x1b` |
| CA-6, CA-7 | `TestChooseStyle` (tabela): `auto`, `sim` e `nao`, com terminal ou não, e com `NO_COLOR` vazia ou definida. `TestIsTerminal`: um `strings.Builder` e a ponta de escrita de um `os.Pipe` não são terminais. `TestRun_Colors`: sem a opção, a saída do teste não tem cores |
| CA-8 | `TestRun_InvalidColors`: a mensagem nova e o código 2 |

O caminho com terminal de verdade é conferido à mão: `go run .` no terminal (com cores), `go run . | cat` (sem cores) e `NO_COLOR=1 go run .` (sem cores). O resultado fica no encerramento.

## Dependências e riscos

- **Conflito com a 014:** a 014 cria duas mensagens do posicionamento: `Não há navio para desfazer.` e o aviso de que o resto da frota não cabe. Pela ordem combinada (014, 018, 015), quem coordena traz o `master` para esta spec depois da 014 e pinta as duas mensagens em amarelo, com o teste. Fica registrado no encerramento.
- **Conflito com a 018:** as duas mexem no `cmd/init.go`, em lugares diferentes (`--cores` e `--nivel`). A resolução é mecânica.

## Checagem da constituição

- [x] **Domínio separado da interface:** o domínio não muda.
- [x] **Aleatoriedade injetada:** não se aplica; a semente continua repetindo a partida.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** a detecção de terminal usa só `os`.
- [x] **Erros como valores:** opção inválida vira código 2 com mensagem.
- [x] **Até ~8 tarefas:** 4.
- [x] **Idiomas:** código em inglês; mensagens em português.

Exceções e justificativas: nenhuma.
