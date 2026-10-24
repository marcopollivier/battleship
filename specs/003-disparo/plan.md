# 003 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

Os disparos ficam no próprio `Board` de `internal/board`, num arquivo novo, `shot.go`. O `Board` ganha um único campo, o conjunto de casas que já receberam disparo. O resultado de cada casa (água ou acerto) não é guardado: ele é derivado desse conjunto e dos navios posicionados, então há uma única fonte de verdade.

Na onda 1, esta spec é a única que edita `board.go`, e só para acrescentar esse campo ao `Board`. A 005 trabalha em `ship.go` e `random.go`, e a 004 num pacote novo, então os três PRs não devem conflitar.

## Pacotes e API

`internal/board/shot.go`:

```go
var (
	ErrAlreadyShot     = errors.New("cell was already shot")
	ErrIncompleteBoard = errors.New("board is not complete")
)
// A shot outside the board reuses ErrOutOfBounds from coord.go.

// ShotResult is the outcome of a shot: Miss, Hit or Sunk.
type ShotResult int

const (
	Miss ShotResult = iota
	Hit
	Sunk
)

func (r ShotResult) String() string

// Shot is the outcome of firing at a cell. Ship is set only when Result is Sunk,
// so a Hit does not reveal which ship was hit.
type Shot struct {
	Result ShotResult
	Ship   ShipType
}

// Fire shoots at c. On error the board is left unchanged.
func (b *Board) Fire(c Coord) (Shot, error)

// CellState is what the shooter knows about a cell.
type CellState int

const (
	CellUntouched CellState = iota
	CellMiss
	CellHit
)

func (s CellState) String() string

// Cell reports the state of c. Cells outside the board are CellUntouched.
func (b *Board) Cell(c Coord) CellState

// AllSunk reports whether every placed ship has been sunk.
func (b *Board) AllSunk() bool
```

Em `board.go`, só o campo novo: `shots map[Coord]bool`, criado por `NewBoard`.

Ordem das validações em `Fire`, com cada erro embrulhando o sentinela com `%w` e citando a casa:

1. O tabuleiro não está completo → `ErrIncompleteBoard`.
2. A casa está fora do tabuleiro → `ErrOutOfBounds`.
3. A casa já recebeu disparo → `ErrAlreadyShot`.

Depois de registrar o disparo, `Fire` procura o navio que ocupa a casa. Se não houver, o resultado é `Miss`. Se todas as casas desse navio já tiverem recebido disparo, o resultado é `Sunk`, com o tipo do navio. Caso contrário, é `Hit`.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| Disparos no próprio `Board` | Um tipo `Game` separado em volta do `Board` | O estado de cada casa é do tabuleiro. Turnos e jogadores ficam para a #006 e a #007. |
| Guardar só o conjunto de casas atingidas | Matriz de estados por casa | Água ou acerto se deriva dos navios. Uma fonte de verdade só. |
| `Shot.Ship` preenchido só em `Sunk` | Sempre preenchido, deixando o terminal esconder | A regra "acerto não revela o tipo" fica garantida no domínio, e o terminal não precisa lembrar de esconder. |
| Checar tabuleiro completo antes de tudo | Checar a casa primeiro | Num tabuleiro incompleto, nenhum disparo faz sentido, qualquer que seja a casa. |
| `AllSunk` olha os navios posicionados | Olhar a frota | Com a regra do tabuleiro completo, os dois coincidem. Olhar os navios evita depender da frota. |

## Estratégia de testes

Arquivo `internal/board/shot_test.go`. Um helper `exampleBoard(t)` monta o tabuleiro de exemplo da spec, reaproveitando os helpers `place` e `mustCoord` que já existem nos testes do pacote. Erros são checados com `errors.Is`, e os casos de erro também conferem que `Cell` não mudou em nenhuma casa.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestFire`: `J10` → `Miss` |
| CA-2 | `TestFire`: `E1` → `Hit`, com `Ship` vazio (valor zero) |
| CA-3 | `TestFire_Sinks`: `J1`, `J2` → `Hit`; `J3` → `Sunk` com `Ship == Cruiser` |
| CA-4 | `TestFire_Sinks`: `H5` → `Sunk` com `Ship == Submarine` |
| CA-5 | `TestBoard_Cell`: depois de `J10` e `E1`, `Cell(J10) == CellMiss`, `Cell(E1) == CellHit`, e as outras 98 casas são `CellUntouched` |
| CA-6 | `TestFire_Rejected`: `Coord{Col: 10}` e `Coord{Row: -1}` → `ErrOutOfBounds` |
| CA-7 | `TestBoard_AllSunk`: dispara em todas as casas de navio menos `E8` → `false`; depois de `E8` → `true` |
| CA-8 | `TestFire_Rejected`: repetir `J10` e `E1` → `ErrAlreadyShot` |
| CA-9 | `TestFire_Rejected`: tabuleiro só com um Cruzador → `ErrIncompleteBoard` |

## Dependências e riscos

- **Valor zero do `Shot.Ship`:** quando o resultado não é `Sunk`, o campo fica com o valor zero, que é `Submarine`. Quem lê o resultado precisa olhar o `Result` antes. Isso fica documentado no tipo e testado no CA-2.

## Checagem da constituição

- [x] **Domínio separado da interface:** tudo em `internal/board`, sem I/O.
- [x] **Aleatoriedade injetada:** não se aplica.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** só `errors` e `fmt`.
- [x] **Erros como valores:** todas as rejeições retornam `error`, sem alterar o tabuleiro.
- [x] **Até ~8 tarefas:** a estimativa é de 5.

Exceções e justificativas: nenhuma.
