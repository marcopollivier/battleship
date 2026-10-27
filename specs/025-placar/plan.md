# 025 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

1. **Placar (`internal/score`, pacote novo):** os números de cada nível e as regras de registrar uma partida. Ler e gravar o arquivo JSON fica em funções separadas das regras, para que as regras sejam testadas sem disco.
2. **Resultado da partida (`internal/terminal`):** o `PlayVsMachine` passa a devolver como a partida acabou e quantos disparos o jogador fez, além do erro.
3. **Programa (`cmd`):**
   - no fim de uma partida que conta, carrega o placar, registra o resultado, grava e mostra a linha do nível;
   - as opções `--placar` e `--zerar-placar`;
   - o caminho do arquivo é passado por parâmetro a uma função interna, para os testes usarem uma pasta temporária.

## Pacotes e API

`internal/score/score.go`:

```go
// Record is the score of one level.
type Record struct {
	Wins, Losses int
	BestWin      int // fewest player shots in a win, 0 with no win
}

// Board is the score of every level, by level name.
type Board map[string]Record

// AddWin records a win at level in shots; AddLoss records a loss.
func (b Board) AddWin(level string, shots int)
func (b Board) AddLoss(level string)

// Load reads the board at path. A missing file is an empty board.
func Load(path string) (Board, error)

// Save writes b to path, creating its directory, through a temporary file
// renamed into place so a failed write never leaves half a file.
func Save(path string, b Board) error

// DefaultPath is placar.json in the battleship folder of the user's
// configuration directory (os.UserConfigDir).
func DefaultPath() (string, error)
```

- O JSON guarda o mapa direto, por exemplo `{"dificil": {"wins": 3, "losses": 2, "best_win": 41}}`.
- Um arquivo corrompido devolve erro, e o `cmd` mostra o aviso da spec. O placar não é zerado sozinho, para não apagar o que a pessoa tinha.

`internal/terminal/machine.go`:

```go
// Outcome is how a game against the machine ended.
type Outcome int

const (
	GaveUp     Outcome = iota // sair or end of input, while placing or shooting
	PlayerWon
	MachineWon
)

// Result is how a game ended and how many shots the player fired.
type Result struct {
	Outcome Outcome
	Shots   int
}

func PlayVsMachine(...) (Result, error)
```

`internal/terminal/score.go` (novo), os textos:

```go
// ScoreLine is "Placar no nível dificil: 3 vitórias, 2 derrotas. Melhor vitória: 41 disparos."
func ScoreLine(level string, r score.Record) string
```

- Singular e plural: `1 vitória`, `1 derrota`, `1 disparo`. Sem vitórias, a parte da melhor vitória não aparece.

`cmd/init.go`:
- `Run` chama `run(args, stdin, stdout, stderr, scorePath)`, com o caminho do `score.DefaultPath()`. Se não houver pasta de configuração, o caminho fica vazio e todo uso do placar mostra o aviso.
- **Quando conta:** contra a máquina, sem `--frota`, e com um destes resultados:
  - `PlayerWon`: `AddWin(nível, Shots)`;
  - `MachineWon`: `AddLoss`;
  - `GaveUp` com `Shots > 0`: `AddLoss`.

  Nesse caso, o `cmd` carrega, registra, grava e mostra `ScoreLine` antes da semente. Qualquer erro de ler ou gravar mostra `Não foi possível guardar o placar.` no `stderr`, sem a linha do placar, e o código de saída continua o da partida.
- O nível é guardado em minúsculas (`MESTRE` vira `mestre`).
- `--placar`: mostra a `ScoreLine` dos três níveis, na ordem `facil`, `dificil`, `mestre`, e retorna 0.
- `--zerar-placar`: grava um placar vazio, mostra `Placar zerado.` e retorna 0.
- As duas opções vêm logo depois do `Parse` e do `--versao` (024), antes de validar a partida.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| Pacote `internal/score` separado | Guardar no `cmd`; pôr no domínio `board` | O placar não é regra da partida, e o domínio não faz I/O. Um pacote próprio deixa as regras testáveis sem disco. |
| `PlayVsMachine` devolve `Result` | O `cmd` deduzir o resultado pelos tabuleiros | Só o laço sabe se a pessoa desistiu e quantos disparos fez. Deduzir pelos tabuleiros repetiria regras do laço. |
| Caminho por parâmetro numa função interna `run` | Uma variável de ambiente; uma opção escondida | Os testes usam `t.TempDir()` sem mexer no ambiente nem no placar real de quem roda os testes. |
| Gravar num arquivo temporário e renomear | Gravar direto | Uma falha no meio da gravação não deixa um JSON quebrado. |
| Não zerar sozinho um arquivo corrompido | Recomeçar do zero | Não apaga o placar sem a pessoa pedir; o `--zerar-placar` resolve. |

## Estratégia de testes

| Critério | Teste |
| -------- | ----- |
| CA-4 | `TestBoard_AddWin`: duas vitórias, com 50 e 41 disparos, dão a melhor 41; `TestScoreLine` (tabela: plural, singular, sem vitórias) |
| CA-1 | `TestRun_ScoreWin` (cmd): vitória com `--frota` padrão e semente fixa (entrada montada como nos testes da 007); o arquivo tem a vitória e a saída tem a linha antes da semente |
| CA-2 | `TestRun_ScoreLoss`: derrota e desistência depois de um disparo |
| CA-3 | `TestRun_ScoreNotCounted`: desistência no posicionamento, `--modo solo` e `--frota`; o arquivo não muda e a linha não aparece |
| CA-5 | `TestRun_ShowScore`: `--placar` com um arquivo preparado |
| CA-6 | `TestRun_ResetScore`: `--zerar-placar` |
| CA-7 | `TestLoad_Corrupt`, `TestSave_Unwritable` (pacote `score`) e `TestRun_ScoreError`: arquivo corrompido e pasta sem permissão; o aviso no `stderr`, sem a linha do placar, e o código de saída igual |
| CA-8 | `TestRun_SameSeed`: cada execução com uma cópia do mesmo arquivo de placar |
| — | `TestPlayVsMachine_Result` (terminal): os três resultados e o número de disparos |

## Dependências e riscos

- **Conflitos:** a 022 muda o `PlayVsMachine` (`Display`), e a 024 muda o `Run` (`--versao`). Pela ordem combinada (024, 023, 022, 025), esta spec entra por último e traz o `master` antes do merge.
- **Arquivo do placar real nos testes:** todo teste do `cmd` passa por `run` com uma pasta temporária. O `Run` público só é chamado pelo `main.go`.

## Checagem da constituição

- [x] **Domínio separado da interface:** o placar fica num pacote próprio, sem misturar com o domínio da partida. O I/O de arquivo fica em `Load` e `Save`.
- [x] **Aleatoriedade injetada:** não se aplica.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** `encoding/json` e `os` da biblioteca padrão.
- [x] **Erros como valores:** erros de arquivo voltam como `error`, e o `cmd` mostra o aviso.
- [x] **Até ~8 tarefas:** 5.
- [x] **Idiomas:** código em inglês; textos em português.

Exceções e justificativas: nenhuma.
