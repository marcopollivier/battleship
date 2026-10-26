# 006 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

O domínio não muda: `RandomBoard`, `Fire` e `AllSunk` já fazem tudo o que a partida precisa. O trabalho fica em duas camadas:

1. **Terminal (`internal/terminal/solo.go`, novo):** `PlaySolo` é o laço da partida. Lê linhas de um `io.Reader`, chama `ParseCoord` e `Fire`, escreve as mensagens e o tabuleiro de ataque num `io.Writer`, e para na vitória ou na desistência. O laço é da apresentação: ele só lê a entrada, chama o domínio e mostra o resultado, como pede a constituição.
2. **Programa (`cmd/init.go`):** `Run` ganha a entrada (`stdin`) e a opção `--semente`. Em vez de mostrar o tabuleiro vazio, sorteia o tabuleiro com a frota escolhida, chama `PlaySolo` e, no fim, escreve `Semente: N`. O `main.go` passa `os.Stdin`.

## Pacotes e API

`internal/terminal/solo.go`:

```go
// PlaySolo runs a solo game: it reads shots from in, one per line, and
// writes the attack board and messages to out until every ship of b is sunk
// or the player gives up with "sair" or the end of the input.
// It returns an error only if writing to out fails.
func PlaySolo(in io.Reader, out io.Writer, b *board.Board) error
```

- Lê com `bufio.Scanner`. O fim da entrada conta como desistência.
- Formato da saída, com uma linha em branco antes de cada tabuleiro:

  ```
  <tabuleiro de ataque + legenda>

  Disparo: <entrada do jogador, ecoada pelo terminal, não pelo programa>
  Água.

  <tabuleiro de ataque + legenda>

  Disparo:
  ```

- Entrada inválida: a mensagem e um novo `Disparo: `, sem tabuleiro.
- Vitória: a mensagem do último afundamento, o tabuleiro final e `Você afundou toda a frota em N disparos!`.
- Desistência: `Você desistiu depois de N disparos. A frota era:`, uma linha em branco e `RenderBoard(b)`. Se a entrada terminou (Ctrl+D), o terminal não pula linha, então o programa escreve uma quebra de linha antes da mensagem.
- `1 disparo` no singular, nas duas mensagens de fim.
- Mensagens de entrada inválida numa função interna, `shotErrorMessage(err error, c board.Coord) string`, com `errors.Is` para `ErrInvalidFormat`, `ErrOutOfBounds` e `ErrAlreadyShot`. Um erro inesperado do `Fire`, como `ErrIncompleteBoard`, não pode acontecer com um tabuleiro do `RandomBoard`; se acontecer, `PlaySolo` o retorna.

`cmd/init.go`:

```go
// Run runs the game with the command-line args (without the program name),
// reading shots from stdin, writing the game to stdout and errors to stderr.
// It returns the exit code.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int
```

- `--semente` é um `fs.Uint64`. Um valor que não é inteiro não negativo é recusado pelo pacote `flag`, com código 2.
- Sem `--semente`, a semente vem de `rand.Uint64()` do gerador global do `math/rand/v2`, que o Go já inicializa com uma semente aleatória.
- O tabuleiro vem de `board.RandomBoard(fleet, newRNG(seed))`, com `newRNG(seed) = rand.New(rand.NewPCG(seed, 0))`.
- O `chooseFleet` passa a usar uma função `flagGiven(fs, name)`, que também diz se `--semente` foi informada.
- Depois do `PlaySolo`, `Run` escreve `Semente: N` e retorna 0.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| Laço da partida em `internal/terminal` | Um pacote `internal/game` novo; o laço em `cmd` | O laço é leitura e escrita em volta do domínio. Um pacote novo só teria esse laço, e em `cmd` ele ficaria misturado com as opções de linha de comando. |
| `PlaySolo` recebe o tabuleiro pronto | Receber a frota e o `*rand.Rand` | Os testes montam o tabuleiro que quiserem, e o sorteio fica junto da semente, em `cmd`. |
| `Semente: N` escrito por `cmd` | `PlaySolo` receber a semente | A semente é um detalhe de como o tabuleiro foi sorteado, e `PlaySolo` não sorteia nada. |
| Semente padrão de `rand.Uint64()` | `time.Now().UnixNano()` | Dá uma semente diferente a cada execução sem converter `int64` em `uint64`, que o `gosec` acusa (G115). A spec diz "semente nova a cada execução", sem citar o relógio. |
| Semente `uint64` | `int64` | É o tipo que o `rand.NewPCG` recebe. |
| Ecoar a entrada fica com o terminal | O programa repetir o que leu | No terminal real a entrada já aparece na tela. Nos testes, a saída fica sem a entrada, o que é mais fácil de conferir. |

## Estratégia de testes

Arquivo novo: `internal/terminal/solo_test.go`. Os helpers novos usam o prefixo `solo` (por exemplo, `soloPlay(t, b, input) string`, que roda `PlaySolo` com a entrada dada e retorna a saída). O tabuleiro dos testes do terminal vem de `board.RandomBoard` com uma semente fixa, e os testes leem as casas dos navios com `b.Ships()`. As casas de água são as que nenhum navio ocupa.

O `cmd/init_test.go` é reescrito para a nova assinatura. O `TestRun_DefaultFleet` atual deixa de valer, porque o programa não mostra mais o tabuleiro vazio.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestPlaySolo_Start`: com a entrada vazia, a saída começa com o tabuleiro de ataque sem disparos e `Disparo: ` |
| CA-2 | `TestPlaySolo_Shots`: água, acerto e afundamento (um Submarino), cada mensagem seguida do tabuleiro igual a `RenderAttackBoard` no estado daquele momento |
| CA-3 | `TestPlaySolo_InvalidInput`: `7B`, `K1` e a mesma casa duas vezes geram as três mensagens, sem tabuleiro depois delas, e a desistência no fim diz `1 disparo` |
| CA-4 | `TestPlaySolo_Win`: disparos em todas as casas dos navios mais duas águas, e a saída termina com o tabuleiro final e `Você afundou toda a frota em 21 disparos!`. `TestRun_Win` (cmd): com `--semente 42` e as casas do tabuleiro da semente 42, código 0 e a última linha é `Semente: 42` |
| CA-5 | `TestPlaySolo_GiveUp`: `sair`, `SAIR` e entrada que termina, cada um com a mensagem de desistência e `RenderBoard`. `TestRun_GiveUp` (cmd): código 0 e `Semente: ` no fim |
| CA-6 | `TestRun_SameSeed`: duas execuções com `--semente 42` e `sair` têm a mesma saída, que revela o tabuleiro |
| CA-7 | `TestRun_Fleet`: com `--frota S=2 --semente 7`, disparos nas 2 casas dos Submarinos terminam em vitória com `2 disparos` |
| CA-8 | `TestRun_InvalidSeed`: `--semente abc` e `--semente -1` dão código 2 e nada em `stdout` |

Os casos de frota recusada do `TestRun_Fleet` atual continuam.

## Dependências e riscos

- **Mudança na assinatura de `cmd.Run`:** os únicos chamadores são o `main.go` e o `cmd/init_test.go`.
- **Testes do `cmd` dependem do `newRNG`:** para saber onde estão os navios da semente 42, o teste usa o mesmo `newRNG` do programa (mesmo pacote). Se o jeito de sortear mudar, o teste acompanha sem mudança.
- **`go run .` passa a esperar entrada:** quem rodar `make run` cai numa partida. É o comportamento pedido.

## Checagem da constituição

- [x] **Domínio separado da interface:** o domínio não muda. O laço só chama `ParseCoord`, `Fire` e `AllSunk`.
- [x] **Aleatoriedade injetada:** o `*rand.Rand` é criado em `cmd` a partir da semente e passado ao `RandomBoard`. O gerador global só escolhe a semente quando o jogador não informa uma, e os testes sempre informam.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** `bufio`, `flag` e `math/rand/v2` são da biblioteca padrão.
- [x] **Erros como valores:** entradas inválidas viram mensagens a partir dos erros sentinela; nada de `panic`.
- [x] **Até ~8 tarefas:** a estimativa é de 5.
- [x] **Idiomas:** mensagens em português, código em inglês.

Exceções e justificativas: o `rand.New(rand.NewPCG(...))` em `cmd` deve precisar de um `//nolint:gosec` (G404), como em `board/fleet.go`, porque o sorteio de um jogo não precisa de gerador criptográfico.
