# 004 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Rascunho

## Abordagem

Criar o pacote `internal/terminal`, a camada de apresentação. Ele lê o domínio (`internal/board`) e escreve texto num `io.Writer`, mas nunca altera o domínio. Ele traz:

- os nomes em português de navios e orientações;
- a função que desenha o tabuleiro próprio com a legenda.

O `cmd/` deixa de montar o tabuleiro à mão: passa a criar um `board.Board` vazio com a frota padrão e chamar `terminal.RenderBoard`.

Esta spec não toca em nenhum arquivo de `internal/board`, então não conflita com a 003 nem com a 005 na onda 1.

## Pacotes e API

`internal/terminal/names.go`:

```go
// ShipName returns the Portuguese name of a ship type, e.g. "Cruzador".
func ShipName(t board.ShipType) string

// OrientationName returns the Portuguese name of an orientation, e.g. "cima".
func OrientationName(o board.Orientation) string
```

Tipos ou orientações desconhecidos caem no `String()` do domínio (por exemplo, `ShipType(99)`), para nunca retornar texto vazio.

`internal/terminal/render.go`:

```go
// Legend explains the board symbols.
const Legend = "Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões"

// RenderBoard writes the player's own board, followed by a blank line and the legend.
func RenderBoard(w io.Writer, b *board.Board) error
```

- O texto é montado num `strings.Builder` e escrito com uma única chamada a `w.Write`, então há um só erro para tratar.
- Letras das colunas: vêm de `board.Coord{Col: c}.String()`, tirando o número, para não duplicar a regra da 001 nem converter `int` para `rune` (o que o `gosec` acusaria).
- Rótulo das fileiras: `fmt.Sprintf("%2d", row+1)`. As casas são separadas por um espaço, sem espaço no fim da linha.
- Navios: um mapa `ShipType → letra` (`S`, `D`, `C`, `E`, `P`) aplicado às casas de cada navio em `b.Ships()`. Um tipo desconhecido aparece como `?`.

`cmd/init.go`:

```go
// Run writes the game screen to w.
func Run(w io.Writer) error

// Init runs the game on standard output.
func Init() error
```

O `main.go` passa a tratar o erro de `cmd.Init()`: imprime em `os.Stderr` e sai com código 1.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| Pacote `internal/terminal` | `internal/ui`; desenhar dentro do `cmd/` | Nome direto e separado do domínio. O `cmd/` fica só com a ligação. |
| Escrever num `io.Writer` | Retornar `string` | Combina com `os.Stdout` no programa e com `strings.Builder` nos testes, e deixa a #009 compor várias telas. |
| Letras das colunas a partir de `Coord.String()` | Constante `"ABCDEFGHIJ"`; conversão `rune('A'+c)` | Uma única regra de nomes de coluna (a da 001), sem conversão que o `gosec` acuse. |
| Nomes em mapas no `terminal` | `String()` em português no domínio | Segue o princípio 7 da constituição: domínio em inglês, textos do jogador em português. |
| `cmd.Run(io.Writer)` | Testar o `main` executando o binário | Testável com a biblioteca padrão, sem processo externo. |
| `Init()` passa a retornar `error` | Ignorar o erro de escrita | O `errcheck` exige tratar o erro, e o `main` é o lugar certo para isso. |

## Estratégia de testes

Arquivos `internal/terminal/names_test.go`, `internal/terminal/render_test.go` e `cmd/init_test.go`. As saídas são comparadas com strings exatas escritas nos próprios testes (golden inline).

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestRenderBoard_Empty`: saída exata do tabuleiro vazio, mais linha em branco e legenda; nenhuma linha termina com espaço |
| CA-2 | `TestRenderBoard_Ships`: Cruzador em `A1` horizontal → fileira 1 começa com `C C C ~` |
| CA-3 | `TestRenderBoard_Ships`: Porta-Aviões em `E5` `Down` → `P` em `D5`, `E5`, `F5`, `E6`, `E7` |
| CA-4 | `TestShipName` e `TestOrientationName`: tabelas com os 5 tipos e as 6 orientações, mais um valor desconhecido |
| CA-5 | `TestRun`: a saída de `cmd.Run` é igual à de `RenderBoard` para um tabuleiro vazio com a frota padrão |
| CA-6 | `TestRenderBoard_FullFleet`: tabuleiro completo (mesma disposição dos testes da 002) com saída exata, e todos os bytes antes da legenda abaixo de 128 |

## Dependências e riscos

- **Mudança na assinatura de `cmd.Init`:** agora retorna `error`. O único chamador é o `main.go`, que muda no mesmo commit.

## Checagem da constituição

- [x] **Domínio separado da interface:** o `terminal` só lê o domínio, e o domínio não conhece o `terminal`.
- [x] **Aleatoriedade injetada:** não se aplica.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** só `fmt`, `io`, `os` e `strings`.
- [x] **Erros como valores:** os erros de escrita sobem até o `main`.
- [x] **Até ~8 tarefas:** a estimativa é de 4.
- [x] **Idiomas:** textos para o jogador em português, código em inglês.

Exceções e justificativas: nenhuma.
