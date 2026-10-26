# 008 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

Três camadas, cada uma com sua parte:

1. **Domínio (`internal/board/fleet.go`, novo):** `ValidateFleet` aplica as regras da frota: nada negativo, pelo menos um navio, cabe no tabuleiro. Para saber se cabe, chama o `RandomBoard` da 005 com uma semente fixa.
2. **Terminal (`internal/terminal/fleet.go`, novo):**
   - `ParseFleet` interpreta o texto `S=2,P=1` usando as mesmas letras do desenho (o mapa `shipSymbols` de `render.go`, só lido).
   - `FleetErrorMessage` traduz cada motivo de recusa para português.
3. **Programa (`cmd/` e `main.go`):** a opção `--frota` é lida com o pacote `flag`. Com a frota aceita, o programa mostra o tabuleiro. Com a frota recusada, mostra a mensagem em `stderr` e termina com código 1.

Na onda 2, a 009 também trabalha em `internal/terminal` (`render.go`, `attack.go`, `messages.go`). Esta spec **não edita** nenhum arquivo existente de `internal/terminal`; só cria `fleet.go` e `fleet_test.go`. Em `internal/board`, só cria `fleet.go` e `fleet_test.go`. Em `cmd/` e `main.go`, é a única spec da onda.

## Pacotes e API

`internal/board/fleet.go`:

```go
var (
	ErrNegativeCount = errors.New("fleet has a negative ship count")
	ErrEmptyFleet    = errors.New("fleet has no ships")
)
// A fleet that cannot be placed reuses ErrFleetDoesNotFit from random.go.

// ValidateFleet reports whether fleet can be used in a game: no negative
// counts, at least one ship, and a board can be built with it.
func ValidateFleet(fleet []FleetEntry) error
```

- A checagem de "cabe" usa `RandomBoard(fleet, rand.New(rand.NewPCG(fitCheckSeed, 0)))`, com `const fitCheckSeed = 1`. A mesma frota sempre tem a mesma resposta.
- Ordem: quantidade negativa, depois frota vazia, depois cabe.

`internal/terminal/fleet.go`:

```go
// ErrInvalidFleetFormat means the fleet text is not a comma-separated list of letter=count pairs.
var ErrInvalidFleetFormat = errors.New("invalid fleet format")

// ParseFleet reads a fleet such as "S=2,P=1" using the board letters
// (S, D, C, E, P). Types not listed get no ships. It does not validate
// counts; use board.ValidateFleet for that.
func ParseFleet(s string) ([]board.FleetEntry, error)

// FleetErrorMessage returns the Portuguese message for an error from
// ParseFleet or board.ValidateFleet.
func FleetErrorMessage(err error) string
```

- `ParseFleet` divide por vírgula, tira os espaços de cada par, divide por `=`, converte a letra para maiúscula e procura o tipo no inverso de `shipSymbols`. A quantidade é lida com `strconv.Atoi`, que aceita `-1` (a recusa de negativo é do domínio). Letra repetida, par sem `=`, letra desconhecida, número inválido e texto vazio → `ErrInvalidFleetFormat`.
- O resultado vem na ordem dos tipos (`Submarine` até `AircraftCarrier`), só com os tipos citados.
- Mensagens de `FleetErrorMessage`:

  | Erro | Mensagem |
  | ---- | -------- |
  | `ErrInvalidFleetFormat` | `Frota inválida: use letra=quantidade separados por vírgula, por exemplo S=2,D=1,P=1.` |
  | `board.ErrNegativeCount` | `Frota inválida: as quantidades não podem ser negativas.` |
  | `board.ErrEmptyFleet` | `Frota inválida: escolha pelo menos um navio.` |
  | `board.ErrFleetDoesNotFit` | `Não foi possível posicionar essa frota no tabuleiro.` |
  | qualquer outro | `Erro: ` seguido do texto do erro |

`cmd/init.go`:

```go
// Run runs the game with the command-line args (without the program name),
// writing the screen to stdout and errors to stderr. It returns the exit code.
func Run(args []string, stdout, stderr io.Writer) int
```

- Usa um `flag.FlagSet` próprio, com saída de erros em `stderr`. A opção `--frota` é detectada com `fs.Visit`, para distinguir "não informada" (frota padrão) de "informada vazia" (erro de formato).
- O `Init()` sai. O `main.go` passa a ser `os.Exit(cmd.Run(os.Args[1:], os.Stdout, os.Stderr))`.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| `ValidateFleet` no domínio, `ParseFleet` no terminal | Tudo no terminal | As regras da frota são do jogo. O formato `S=2` usa as letras do desenho, então é da apresentação. |
| Semente fixa na checagem de "cabe" | Semente aleatória; busca exaustiva | Resposta determinística, sem novo algoritmo. A limitação foi aceita na spec. |
| `Run` retorna o código de saída | Retornar `error` e o `main` imprimir | As mensagens precisam sair em português, e a opção inválida do `flag` usa código 2. Retornar o código deixa tudo testável no `cmd`. |
| Reutilizar `shipSymbols` | Um mapa novo de letras | Uma única fonte para as letras: as mesmas do desenho e da legenda. |
| Remover `cmd.Init()` | Manter como atalho | O único chamador era o `main.go`, e manter dois pontos de entrada não acrescenta nada. |

## Estratégia de testes

Arquivos novos: `internal/board/fleet_test.go` e `internal/terminal/fleet_test.go`. O `cmd/init_test.go` é reescrito para a nova assinatura de `Run`.

Para não colidir com os testes da 009 no mesmo pacote, os helpers novos desta spec usam o prefixo `fleet` (por exemplo, `fleetCount`). Os helpers existentes nos testes de `internal/board` e `internal/terminal` podem ser reutilizados, mas não renomeados.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestParseFleet`: `S=2,P=1` → `[{Submarine 2} {AircraftCarrier 1}]` |
| CA-2 | `TestParseFleet`: ` s=2 , p=1 ` → o mesmo |
| CA-3 | `TestRun_DefaultFleet`: sem argumentos, código 0 e saída igual a `RenderBoard` de um tabuleiro com a frota padrão |
| CA-4 | `TestParseFleet_InvalidFormat`: `""`, `X=1`, `S=a`, `S2`, `S=1;D=1`, `S=1,S=2` → `ErrInvalidFleetFormat` |
| CA-5 | `TestValidateFleet`: `{Submarine -1}` → `ErrNegativeCount` |
| CA-6 | `TestValidateFleet`: `{Submarine 0}` e frota vazia → `ErrEmptyFleet` |
| CA-7 | `TestValidateFleet`: `{Battleship 30}` → `ErrFleetDoesNotFit` |
| CA-8 | `TestValidateFleet_Accepted`: a frota padrão e `S=2,P=1` são aceitas, e `RandomBoard` com cada uma (sementes 1 a 100) gera um tabuleiro completo |
| CA-9 | `TestFleetErrorMessage`: tabela com os quatro erros (também embrulhados com `%w`) e um erro qualquer |
| CA-10 | `TestRun_Fleet`: `--frota E=30` → código 1 e mensagem de frota que não cabe em `stderr`; `--frota S=2,P=1` → código 0 e o tabuleiro em `stdout` |

## Dependências e riscos

- **Mudança na assinatura de `cmd.Run`:** os únicos chamadores são o `main.go` e o `cmd/init_test.go`, que mudam nesta spec.
- **Tempo do CA-7:** `ValidateFleet` com `E=30` passa por todas as tentativas do `RandomBoard`. Na 005, isso levou 0,02 s.

## Checagem da constituição

- [x] **Domínio separado da interface:** as regras ficam em `internal/board`, e o texto e as mensagens em `internal/terminal`.
- [x] **Aleatoriedade injetada:** a checagem de "cabe" cria o próprio `*rand.Rand` com semente fixa. É determinística e não depende de estado global.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** `flag`, `strconv`, `strings` e `math/rand/v2` são da biblioteca padrão.
- [x] **Erros como valores:** todas as recusas retornam erros sentinela.
- [x] **Até ~8 tarefas:** a estimativa é de 6.
- [x] **Idiomas:** mensagens em português, código em inglês.

Exceções e justificativas: nenhuma.
