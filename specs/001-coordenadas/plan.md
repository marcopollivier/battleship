# 001 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

Criar o pacote de domínio `internal/board`, que nas próximas specs também vai receber o tabuleiro, a frota e os disparos. Nesta spec ele ganha só o tipo `Coord`, com índices de coluna e linha a partir de zero, mais duas operações: `ParseCoord`, que converte texto em casa, e `String`, que faz o caminho inverso. A validação acontece em duas etapas, cada uma com seu próprio erro sentinela: primeiro o formato (regex sobre a entrada já sem espaços nas pontas) e depois os limites do tabuleiro. O `cmd/` não é alterado nesta spec.

## Pacotes e API

`internal/board/coord.go`:

```go
// Size is the number of rows and columns on the board.
const Size = 10

var (
 ErrInvalidFormat = errors.New("invalid coordinate format")
 ErrOutOfBounds   = errors.New("coordinate out of bounds")
)

// Coord is a cell on the board. Col and Row are zero-based: A1 is {0, 0}, J10 is {9, 9}.
type Coord struct {
 Col int
 Row int
}

// ParseCoord converts text such as "A1" or " j10 " into a Coord.
// Errors wrap ErrInvalidFormat or ErrOutOfBounds and include the original input.
func ParseCoord(s string) (Coord, error)

// String returns the canonical form, e.g. "A1", "J10".
func (c Coord) String() string
```

Algoritmo de `ParseCoord`:

1. `strings.TrimSpace(s)`, que remove espaços, tabulações e quebras de linha nas pontas.
2. Casar com `^[A-Za-z](0|[1-9][0-9]*)$`. Se não casar, retornar `ErrInvalidFormat`. Esse passo rejeita `1A`, `A 1`, `A01`, `@1`, `Ç1` e o texto vazio.
3. Converter a letra para maiúscula. Se ela passar de `J`, retornar `ErrOutOfBounds`.
4. Se o número tiver mais de 2 dígitos ou estiver fora de 1–10, retornar `ErrOutOfBounds`. Checar a quantidade de dígitos antes de chamar `strconv.Atoi` evita overflow com entradas enormes.
5. Retornar `Coord{Col: letra - 'A', Row: número - 1}`.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| Pacote `internal/board` | `board` na raiz; `internal/coord` | O pacote vai abrigar o domínio inteiro nas specs 002 e 003. `internal/` deixa claro que não é uma API pública. |
| Índices a partir de zero | A partir de 1, igual ao texto | Indexam direto a matriz `[Size][Size]` que o tabuleiro vai usar. A conversão para 1–10 fica só em `ParseCoord` e `String`. |
| `Coord{Col, Row}` com campos exportados | Campos privados com construtor | As specs 002 e 003 vão precisar montar casas vizinhas (`Col+1`) sem passar por texto. O valor zero é `A1`, uma casa válida. |
| Erros sentinela embrulhados com `%w` e a entrada | Tipos de erro próprios; mensagens em português no domínio | `errors.Is` basta para distinguir os dois casos da spec. Pela constituição, o texto para o jogador fica na camada de terminal (#004), que traduz o sentinela para português. |
| Regex para o formato | Parser manual caractere a caractere | O padrão espelha a regra da spec ("letra + número sem zeros à esquerda") numa linha legível. Desempenho não importa aqui. |
| `String()` de uma `Coord` fora do tabuleiro retorna `Coord(col,row)` | `panic`; string vazia | Respeita o princípio de erros sem `panic` e mantém a depuração legível. Nenhum critério depende disso. |

## Estratégia de testes

Arquivo `internal/board/coord_test.go`, com testes table-driven da biblioteca padrão. Os erros são checados com `errors.Is`.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestParseCoord_Valid`: `"A1"` → `{0, 0}` |
| CA-2 | `TestParseCoord_Valid`: `"J10"` → `{9, 9}` |
| CA-3 | `TestParseCoord_Valid`: `"C7"` → `{2, 6}` |
| CA-4 | `TestCoord_RoundTrip`: percorre as 100 casas e verifica que `ParseCoord(c.String()) == c` e que `c.String()` é igual à letra seguida do número montados no próprio teste |
| CA-5 | `TestParseCoord_OutOfBounds`: `"K1"`, `"Z5"` → `ErrOutOfBounds` |
| CA-6 | `TestParseCoord_OutOfBounds`: `"A0"`, `"A11"`, `"B99"`, `"A99999999999999999999"` → `ErrOutOfBounds` |
| CA-7 | `TestParseCoord_InvalidFormat`: `""`, `"A"`, `"1"`, `"AA1"`, `"A1B"`, `"A-1"`, `"@1"`, `"Ç1"` → `ErrInvalidFormat` |
| CA-8 | `TestParseCoord_Valid`: `"a1"` → `{0, 0}`, `"j10"` → `{9, 9}` |
| CA-9 | `TestParseCoord_Valid`: `" A1 "`, `"\tA1\n"` → `{0, 0}` |
| CA-10 | `TestParseCoord_InvalidFormat`: `"A 1"`, `"1A"`, `"A01"` → `ErrInvalidFormat` |

Também testa que a mensagem de erro inclui a entrada original, para facilitar a depuração.

## Dependências e riscos

- **Lint quebrado (bloqueia a fase 4):** o `.golangci.yml` está no formato da v1, e o golangci-lint instalado (2.15) recusa carregá-lo: `unsupported version of the configuration`. O `make pre-commit` e o hook do pre-commit falham hoje, antes de qualquer mudança. Isso é outro contexto, então a correção vai num PR próprio (`chore/golangci-v2`), mergeado antes da implementação.
- **`goconst` nos testes:** a regra está configurada com `min-occurrences: 2` e vale também para `_test.go`. Pode reclamar de literais repetidos nas tabelas (`"A1"` aparece em mais de um caso). Se isso acontecer, a opção é excluir o `goconst` de `_test.go`, que entra no mesmo PR de migração do lint.

## Checagem da constituição

- [x] **Domínio separado da interface:** `internal/board` não faz I/O, e `cmd/` não muda.
- [x] **Aleatoriedade injetada:** não se aplica, a spec não usa aleatoriedade.
- [x] **Testes antes do código:** cada tarefa começa pelo teste do critério que cobre.
- [x] **Sem dependências externas:** usa só `errors`, `fmt`, `regexp`, `strconv` e `strings`.
- [x] **Erros como valores:** entradas inválidas retornam `error`, e `String()` de uma casa inválida não entra em `panic`.
- [x] **Até ~8 tarefas:** a estimativa é de 4 tarefas.

Exceções e justificativas: nenhuma.
