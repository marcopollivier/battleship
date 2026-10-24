# 005 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

Criar `internal/board/random.go`, com uma função que monta um tabuleiro completo sorteando posições e reaproveitando o `Board.Place` da 002 para todas as validações. A estratégia é tentar e descartar:

1. Os navios são posicionados do maior para o menor, porque os grandes são os mais difíceis de encaixar.
2. Para cada navio, sorteia casa e orientação válida e tenta `Place`, com um limite de tentativas.
3. Se um navio esgotar as tentativas, o tabuleiro inteiro é descartado e o processo recomeça, também com um limite de recomeços.
4. Esgotados os recomeços, retorna erro.

Para saber quais orientações sortear, `ship.go` ganha `ValidOrientations`, que também vai servir à #007 para mostrar as opções ao jogador.

Na onda 1, esta spec só toca `ship.go` (para acrescentar `ValidOrientations`) e o arquivo novo `random.go`. Não edita `board.go`, que é da 003, nem nada do terminal (004).

## Pacotes e API

`internal/board/ship.go`, que só ganha esta função:

```go
// ValidOrientations returns the orientations that apply to t, or nil for an
// unknown type: Horizontal and Vertical for line ships, Up, Down, Left and
// Right for the AircraftCarrier.
func ValidOrientations(t ShipType) []Orientation
```

`internal/board/random.go`:

```go
// ErrFleetDoesNotFit means no valid board was found for the fleet within the attempt limits.
var ErrFleetDoesNotFit = errors.New("fleet does not fit on the board")

// RandomBoard returns a complete board with the fleet placed at random.
// The same rng state always produces the same board.
func RandomBoard(fleet []FleetEntry, rng *rand.Rand) (*Board, error)
```

Detalhes:

- `rng` é um `*rand.Rand` de `math/rand/v2`. Quem chama decide a semente: nos testes, `rand.New(rand.NewPCG(seed, 0))`; na #006, o relógio ou a semente do jogador.
- Os limites são constantes não exportadas: `maxAttemptsPerShip = 200` e `maxRestarts = 50`. Isso garante que a função termina, com no máximo 10.000 tentativas por navio.
- Um tipo sem orientações válidas (desconhecido) retorna `ErrFleetDoesNotFit` imediatamente, sem gastar tentativas.
- A ordem do maior para o menor é estável (`slices.SortStableFunc` por `Size`), então a mesma semente sempre gera o mesmo tabuleiro.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| Tentar e descartar, reaproveitando `Place` | Listar todas as posições válidas e sortear uma; backtracking | Reaproveita as regras da 002 sem duplicá-las. A spec não exige distribuição uniforme. |
| Maior navio primeiro | Ordem da frota; ordem aleatória | Reduz muito os recomeços: os navios pequenos sempre encontram espaço no fim. |
| Recomeçar o tabuleiro inteiro | Desfazer só o último navio | Não existe "remover navio" no `Board` (decisão da 002). Recomeçar é simples e raro com a frota padrão. |
| `*rand.Rand` de `math/rand/v2` como parâmetro | `rand.Source`; uma interface própria | É o tipo padrão do Go para aleatoriedade com semente, e atende o princípio 2 da constituição. |
| `ValidOrientations` exportada | Função interna em `random.go` | O conhecimento de orientações por tipo é de `ship.go`, e a #007 vai precisar dele. |

## Estratégia de testes

Arquivo `internal/board/random_test.go`. Um helper `seeded(seed)` cria o `*rand.Rand`. Um verificador independente, `assertValidBoard`, confere a partir de `Ships()` as regras da 002 sem passar pelo `Place`: casas dentro do tabuleiro, nenhuma casa repetida, nenhuma casa vizinha de outro navio e quantidade por tipo. Para contar orientações, um helper `orientationOf(ship)` deduz a orientação a partir das casas.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestRandomBoard_Complete`: `Complete()` é verdadeiro e há 8 navios |
| CA-2 | `TestRandomBoard_Valid`: sementes 1 a 1.000, cada tabuleiro passa no `assertValidBoard` |
| CA-3 | `TestRandomBoard_SameSeed`: dois `RandomBoard` com a semente 42 têm os mesmos `Ships()` |
| CA-4 | `TestRandomBoard_DifferentSeeds`: as sementes 1 e 2 geram `Ships()` diferentes |
| CA-5 | `TestRandomBoard_AllOrientations`: nas sementes 1 a 1.000, cada par (tipo, orientação válida) aparece pelo menos uma vez; o Submarino entra como uma orientação só |
| CA-6 | `TestRandomBoard_DoesNotFit`: frota com 30 Encouraçados → `ErrFleetDoesNotFit` |
| — | `TestValidOrientations`: tabela com os 5 tipos e um tipo desconhecido |

## Dependências e riscos

- **`gosec` G404 (gerador fraco):** o `gosec` pode acusar o uso de `math/rand/v2`. Num jogo, aleatoriedade sem uso criptográfico é a escolha certa. Se o lint acusar, a correção é um `//nolint:gosec` na linha, com essa justificativa. Os testes já ignoram o `gosec`.
- **Tempo do CA-6:** com os limites acima, o pior caso são 50 recomeços × 30 navios × 200 tentativas. O teste confirma que isso termina rápido; se ficar lento, os limites são reduzidos.
- **Submarino no CA-5:** o Submarino aceita `Horizontal` e `Vertical`, mas as casas são as mesmas nas duas, então não dá para deduzir a orientação pelas casas. Para ele, o teste só confere que aparece.

## Checagem da constituição

- [x] **Domínio separado da interface:** tudo em `internal/board`, sem I/O.
- [x] **Aleatoriedade injetada:** o `*rand.Rand` vem de quem chama, e os testes usam sementes fixas.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** `math/rand/v2`, `errors`, `fmt` e `slices` são da biblioteca padrão.
- [x] **Erros como valores:** a frota que não cabe retorna `ErrFleetDoesNotFit`.
- [x] **Até ~8 tarefas:** a estimativa é de 4.

Exceções e justificativas: nenhuma.
