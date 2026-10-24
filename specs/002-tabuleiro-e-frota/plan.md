# 002 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

Ampliar o pacote `internal/board` da spec 001 com dois arquivos:

- `ship.go`: tipos de navio, orientações, frota padrão e o cálculo das casas que um navio ocupa a partir de uma referência e de uma orientação. São funções puras, sem tabuleiro.
- `board.go`: o tabuleiro de um jogador. Recebe posicionamentos e valida, nesta ordem: quantidade, orientação, limites, sobreposição e encostar. O navio só é gravado se todas as validações passarem, então um erro nunca altera o tabuleiro.

O estado do tabuleiro é só a lista de navios posicionados. Com 19 casas no máximo, as checagens percorrem essa lista em vez de manter uma matriz paralela. A matriz pode entrar na #003, se os disparos precisarem.

Pela estimativa, a spec cabe em 7 tarefas, então não é preciso dividi-la (questão 4).

## Pacotes e API

`internal/board/ship.go`:

```go
// ShipType identifies a kind of ship. Names are in English; the terminal
// layer translates them to Portuguese.
type ShipType int

const (
	Submarine       ShipType = iota // 1 cell
	Destroyer                       // line of 2
	Cruiser                         // line of 3
	Battleship                      // line of 4
	AircraftCarrier                 // T of 5
)

func (t ShipType) Size() int
func (t ShipType) String() string

// Orientation is Horizontal or Vertical for line ships, and the direction
// the stem points (Up, Down, Left, Right) for the T-shaped AircraftCarrier.
type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
	Up
	Down
	Left
	Right
)

func (o Orientation) String() string

// FleetEntry is how many ships of one type a fleet has.
type FleetEntry struct {
	Type  ShipType
	Count int
}

// DefaultFleet returns the fleet from the original paper pad: 8 ships, 19 cells.
func DefaultFleet() []FleetEntry

// ShipCells returns the cells a ship of type t occupies when placed at ref
// with orientation o. It does not check board bounds.
// Returns ErrInvalidOrientation if o does not apply to t.
func ShipCells(t ShipType, ref Coord, o Orientation) ([]Coord, error)
```

Formatos usados por `ShipCells`, com `(dCol, dRow)` relativo a `ref`:

| Tipo | Orientação | Casas |
| ---- | ---------- | ----- |
| Em linha, tamanho `n` | `Horizontal` | `(0,0)` … `(n-1,0)` |
| Em linha, tamanho `n` | `Vertical` | `(0,0)` … `(0,n-1)` |
| Porta-Aviões | `Down` | barra `(-1,0) (0,0) (1,0)`, haste `(0,1) (0,2)` |
| Porta-Aviões | `Up` | barra `(-1,0) (0,0) (1,0)`, haste `(0,-1) (0,-2)` |
| Porta-Aviões | `Right` | barra `(0,-1) (0,0) (0,1)`, haste `(1,0) (2,0)` |
| Porta-Aviões | `Left` | barra `(0,-1) (0,0) (0,1)`, haste `(-1,0) (-2,0)` |

`internal/board/board.go`:

```go
var (
	ErrInvalidOrientation = errors.New("orientation does not apply to ship type")
	ErrOverlap            = errors.New("ship overlaps another ship")
	ErrAdjacent           = errors.New("ship touches another ship")
	ErrFleetLimit         = errors.New("no more ships of this type in the fleet")
)
// A ship that leaves the board reuses ErrOutOfBounds from coord.go.

// PlacedShip is a ship on the board and the cells it occupies.
type PlacedShip struct {
	Type  ShipType
	Cells []Coord
}

// Board is one player's board.
type Board struct { /* fleet []FleetEntry; ships []PlacedShip */ }

func NewBoard(fleet []FleetEntry) *Board

// Place puts a ship on the board. On error the board is left unchanged.
func (b *Board) Place(t ShipType, ref Coord, o Orientation) error

// Ships returns a copy of the placed ships, in placement order.
func (b *Board) Ships() []PlacedShip

// Complete reports whether every ship in the fleet has been placed.
func (b *Board) Complete() bool
```

Ordem das validações em `Place`, onde cada erro embrulha o sentinela com `%w` e cita o tipo, a referência e a orientação:

1. Já existem `Count` navios desse tipo, ou o tipo não está na frota → `ErrFleetLimit`.
2. `ShipCells` retorna erro → `ErrInvalidOrientation`.
3. Alguma casa está fora de `0..Size-1` → `ErrOutOfBounds`.
4. Alguma casa já está ocupada → `ErrOverlap`.
5. Alguma das 8 vizinhas de uma casa nova está ocupada → `ErrAdjacent`.

A sobreposição é checada antes do encostar porque toda sobreposição também encosta, e o CA-7 exige o erro de sobreposição.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| Formatos como deslocamentos `(dCol, dRow)` a partir da referência | Rotacionar uma matriz 3 × 3; tipo `Shape` genérico | Só o Porta-Aviões não é linha. Uma tabela com 4 linhas para ele é mais legível que um algoritmo de rotação e deixa as casas do CA-5 explícitas. |
| `Orientation` única com 6 valores, validada por tipo | Um tipo para linha e outro para o T | Um único `Place(t, ref, o)` atende à spec, e o CA-11 vira uma checagem simples. |
| Submarino aceita só `Horizontal` e `Vertical` | Aceitar as 6 orientações | Segue a regra "navios em linha usam horizontal ou vertical". "Qualquer orientação" no CA-3 é lido como qualquer uma das duas. |
| `ShipCells` separado do `Board` | Calcular dentro de `Place` | Função pura, testável sem tabuleiro, reutilizável no sorteio da #005. |
| Estado do tabuleiro = lista de navios | Matriz `[Size][Size]` com o índice do navio | Uma fonte de verdade só. Com no máximo 19 casas, percorrer a lista é trivial. |
| Reusar `ErrOutOfBounds` | Novo `ErrShipOutOfBounds` | É o mesmo conceito da 001, e o terminal pode tratar os dois casos com uma única mensagem. |
| Frota como `[]FleetEntry` passada a `NewBoard` | Frota fixa dentro do `Board` | Atende "as regras devem funcionar para qualquer frota" e deixa a #008 pronta para receber outra lista. |
| Nomes em inglês no código (`Cruiser`, `AircraftCarrier`, `Up`...) | Nomes em português | Segue o princípio 7 da constituição. A tradução fica na #004. |

## Estratégia de testes

Testes table-driven em `internal/board/ship_test.go` e `internal/board/board_test.go`. Um helper `mustCoord(t, "A1")` usa `ParseCoord` para escrever as casas como na spec. Erros são checados com `errors.Is`.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestDefaultFleet`: contagem por tipo (3, 2, 1, 1, 1) e soma de `Size × Count` igual a 19 |
| CA-2 | `TestShipCells_Line`: Cruzador em `A1`, `Horizontal` → `A1 B1 C1`; `Vertical` → `A1 A2 A3` |
| CA-3 | `TestShipCells_Line`: Submarino em `E5`, `Horizontal` e `Vertical` → só `E5` |
| CA-4 | `TestBoard_Place_OutOfBounds`: Encouraçado em `J7` `Vertical` é aceito e ocupa `J7`–`J10`; `J8` `Vertical` e `H1` `Horizontal` → `ErrOutOfBounds` |
| CA-5 | `TestShipCells_AircraftCarrier`: referência `E5` nas 4 orientações, comparando o conjunto de casas da spec |
| CA-6 | `TestBoard_Place_OutOfBounds`: Porta-Aviões em `E9` `Down` e em `A5` `Up` → `ErrOutOfBounds` |
| CA-7 | `TestBoard_Place_Overlap`: Cruzador em `A1` `Horizontal`, depois Destroyer em `B1` `Vertical` → `ErrOverlap`, e `Ships()` continua com 1 navio |
| CA-8 | `TestBoard_Place_FleetLimit`: 3 Submarinos separados, e o quarto → `ErrFleetLimit` |
| CA-9 | `TestBoard_Complete`: uma disposição fixa dos 8 navios, sem encostar; `false` depois do 7º e `true` depois do 8º |
| CA-10 | `TestBoard_Place_Adjacent`: Cruzador em `A1` `Horizontal`; Destroyer em `D1` `Horizontal`, `A2` `Horizontal` e `D2` `Vertical` → `ErrAdjacent`; em `E1` `Horizontal` é aceito |
| CA-11 | `TestShipCells_InvalidOrientation`: Cruzador `Down` e Porta-Aviões `Horizontal` → `ErrInvalidOrientation`; o mesmo erro sai de `Board.Place` |

Os testes de erro também verificam que `Ships()` não mudou, cobrindo a regra de que um erro deixa o tabuleiro como estava.

## Dependências e riscos

- **`goconst` nos testes:** as tabelas vão repetir literais como `"A1"` e `"E5"`. Na 001 isso não deu problema, mas aqui a repetição é maior. Se o lint reclamar, a correção é excluir o `goconst` de `_test.go` no `.golangci.yml`. Como isso é tooling, vai num PR `chore/` separado, mergeado antes de continuar.
- **Disposição do CA-9:** o teste precisa de uma disposição válida dos 8 navios sem encostar. Ela vai ser montada à mão no teste, e o próprio `Place` confirma que é válida.

## Checagem da constituição

- [x] **Domínio separado da interface:** tudo em `internal/board`, sem I/O. Os nomes em português ficam para a #004.
- [x] **Aleatoriedade injetada:** não se aplica, porque não há aleatoriedade nesta spec.
- [x] **Testes antes do código:** cada tarefa começa pelo teste do critério que cobre.
- [x] **Sem dependências externas:** só `errors` e `fmt`.
- [x] **Erros como valores:** todas as rejeições retornam `error`, e o tabuleiro não muda.
- [x] **Até ~8 tarefas:** a estimativa é de 7, sem necessidade de dividir a spec:
  1. Tipos e frota padrão
  2. Casas dos navios em linha
  3. Casas do Porta-Aviões
  4. `Board.Place` com limites
  5. Sobreposição
  6. Encostar
  7. Quantidade e tabuleiro completo

Exceções e justificativas: nenhuma.
