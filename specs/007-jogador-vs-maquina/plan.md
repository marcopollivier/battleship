# 007 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

Três camadas, como na 006:

1. **Domínio (`internal/board/machine.go`, novo):** `RandomShot` sorteia uma casa ainda não atingida de um tabuleiro, com a fonte de aleatoriedade injetada. É a única regra nova: a máquina não tem estratégia.
2. **Terminal (`internal/terminal/machine.go`, novo):** `PlayVsMachine` é o laço da partida contra a máquina. O turno do jogador reaproveita o que o modo solo já faz (ler a linha, `sair`, mensagens de entrada inválida, mensagem do disparo e tabuleiro de ataque); o turno da máquina chama `RandomShot` e `Fire` no tabuleiro do jogador.
3. **Programa (`cmd/init.go`):** a opção `--modo` escolhe entre `PlayVsMachine` (padrão) e `PlaySolo`. No modo contra a máquina, o mesmo `*rand.Rand` sorteia o tabuleiro do jogador, depois o da máquina, e depois os disparos da máquina.

## Pacotes e API

`internal/board/machine.go`:

```go
// ErrNoCellsLeft means every cell of the board has already been shot.
var ErrNoCellsLeft = errors.New("no cells left to shoot")

// RandomShot returns a cell of b that has not been shot yet, chosen with rng.
// It does not fire. The same rng state and board always give the same cell.
func RandomShot(b *Board, rng *rand.Rand) (Coord, error)
```

- Percorre as 100 casas em ordem, junta as que estão como `CellUntouched` e sorteia uma com `rng.IntN`.

`internal/terminal/machine.go`:

```go
// PlayVsMachine runs a game against the machine: the player shoots at enemy
// with shots read from in, the machine shoots at player with cells drawn
// from rng, and whoever hits shoots again. It writes the game to out and
// returns an error if reading or writing fails.
func PlayVsMachine(in io.Reader, out io.Writer, player, enemy *board.Board, rng *rand.Rand) error
```

- **Reuso do modo solo:** o `soloGame` de `solo.go` passa a dizer o resultado de cada disparo (`shoot` retorna o `board.Shot` e se o disparo foi válido), e a mensagem de desistência recebe o texto final (`A frota era:` no solo, `A frota inimiga era:` contra a máquina). O comportamento do modo solo não muda, e os testes da 006 continuam passando sem alteração.
- **Tela inteira**, no início e cada vez que o turno volta para o jogador:

  ```
  Sua frota:
  <RenderBoard(player)>

  Frota inimiga:
  <RenderAttackBoard(enemy)>

  Disparo:
  ```

- **Disparo do jogador:** como no solo, a mensagem, uma linha em branco e o tabuleiro de ataque. Com acerto, vem o `Disparo: ` de novo. Com água, vem o turno da máquina.
- **Turno da máquina:** uma linha por disparo, `A máquina atirou em C4: <resultado>`, com `Água.`, `Acerto!` ou `Afundou o seu <navio>!`. Depois do disparo na água, uma linha em branco e a tela inteira.
- **Vitória do jogador:** depois do tabuleiro de ataque final, `Você venceu! Afundou a frota inimiga em N disparos.`
- **Vitória da máquina:** depois da linha do último disparo, uma linha em branco, `Sua frota:` com o tabuleiro do jogador, uma linha em branco e `A máquina venceu! Ela afundou a sua frota em N disparos. A frota inimiga era:`, uma linha em branco e `RenderBoard(enemy)`.
- **Desistência:** `Você desistiu depois de N disparos. A frota inimiga era:`, uma linha em branco e `RenderBoard(enemy)`, com a quebra de linha antes no Ctrl+D, como no solo.
- `shotCount` (`1 disparo`) é reaproveitado.

`cmd/init.go`:

- Nova opção `--modo` (`fs.String`, padrão `maquina`), comparada com `strings.EqualFold`. Um valor inválido escreve `Modo inválido: use maquina ou solo.` em `stderr` e retorna 2, antes de validar a frota.
- No modo `maquina`: `rng := newRNG(seed)`, `player := RandomBoard(fleet, rng)`, `enemy := RandomBoard(fleet, rng)`, `PlayVsMachine(stdin, stdout, player, enemy, rng)`.
- No modo `solo`: igual à 006. A mesma semente continua gerando o mesmo tabuleiro do solo.
- Depois das duas partidas, `Semente: N`, como hoje.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| `RandomShot` no domínio | Sortear a casa no terminal | É a regra de como a máquina joga, e a 012 ou uma máquina com estratégia vão mexer nela sem tocar no terminal. |
| `RandomShot` só escolhe a casa, sem atirar | Atirar e retornar o `Shot` | Quem atira é o laço, igual para os dois lados; a função fica fácil de testar. |
| Sortear entre as casas livres | Sortear qualquer casa e tentar de novo se já foi atingida | Número fixo de sorteios por disparo, e o fim do tabuleiro não vira um laço longo. |
| Um único `*rand.Rand` para os tabuleiros e os disparos | Um gerador para cada coisa | Uma semente repete a partida inteira, como a spec pede. |
| Reaproveitar o `soloGame` | Copiar o laço do solo | O turno do jogador é igual nos dois modos; copiar duplicaria as mensagens e o lint (`dupl`) reclamaria. |
| `PlayVsMachine` recebe os dois tabuleiros prontos | Receber a frota | Igual ao `PlaySolo`: os testes montam os tabuleiros que quiserem. |
| `--modo` inválido retorna 2 | Retornar 1 | É um erro de opção, como `--semente abc` (008, 006). |

## Estratégia de testes

Arquivos novos: `internal/board/machine_test.go` e `internal/terminal/machine_test.go`. Os helpers novos usam o prefixo `machine`. Para os testes do terminal, o resultado esperado é montado repetindo a partida em tabuleiros de réplica, com um `*rand.Rand` de mesma semente chamando `RandomShot`, como os testes da 006 fazem com `attackFire`. Os tabuleiros são pequenos e montados à mão (por exemplo, só um Submarino de cada lado), para que os testes de vitória da máquina tenham água suficiente para o jogador errar.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestPlayVsMachine_Start`: com a entrada vazia, a saída começa com a tela inteira. `TestRun_MachineStart` (cmd): sem `--modo`, a saída começa com `Sua frota:` |
| CA-2 | `TestPlayVsMachine_HitShootsAgain`: um acerto é seguido de `Disparo: ` sem linha da máquina |
| CA-3 | `TestPlayVsMachine_MachineTurn`: depois de uma água, as linhas da máquina são as da réplica até o primeiro erro, seguidas da tela inteira |
| CA-4 | `TestRandomShot`: 100 chamadas atirando em cada casa sorteada cobrem as 100 casas sem repetir; a mesma semente dá a mesma sequência; tabuleiro todo atingido → `ErrNoCellsLeft` |
| CA-5 | `TestPlayVsMachine_PlayerWins` e `TestRun_PlayerWins` (cmd, com `--frota S=1 --semente N` e as casas lidas dos tabuleiros da semente) |
| CA-6 | `TestPlayVsMachine_MachineWins` e `TestRun_MachineWins` (cmd): o jogador só atira na água, e a máquina afunda a frota do jogador |
| CA-7 | `TestPlayVsMachine_GiveUp`: `sair` e entrada que termina |
| CA-8 | `TestRun_SameSeed`: duas execuções com `--semente 42` e a mesma entrada têm a mesma saída |
| CA-9 | `TestRun_Solo`: `--modo solo` e `--modo SOLO` mostram a tela do solo, sem `Sua frota:`. Os testes atuais do `cmd`, que são do modo solo, passam a usar `--modo solo` |
| CA-10 | `TestRun_InvalidOption`: `--modo xyz` → código 2 e a mensagem em `stderr` |

## Dependências e riscos

- **Mudança no `soloGame`:** a refatoração do turno do jogador mexe no código da 006. Os testes da 006 ficam sem alteração e precisam passar antes e depois.
- **Testes do cmd em modo solo:** como o padrão muda para `maquina`, os testes do solo no `cmd` precisam passar `--modo solo`. É uma mudança nos testes, não no comportamento do solo.
- **Ordem do sorteio:** o tabuleiro do jogador é sorteado antes do da máquina. Os testes do `cmd` repetem essa ordem com `newRNG`.

## Checagem da constituição

- [x] **Domínio separado da interface:** `RandomShot` fica no domínio, sem I/O. O laço só chama o domínio e escreve.
- [x] **Aleatoriedade injetada:** `RandomShot` e `PlayVsMachine` recebem o `*rand.Rand`; os testes usam sementes fixas.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** só biblioteca padrão.
- [x] **Erros como valores:** `ErrNoCellsLeft` e os erros de entrada; nada de `panic`.
- [x] **Até ~8 tarefas:** a estimativa é de 6.
- [x] **Idiomas:** mensagens em português, código em inglês.

Exceções e justificativas: nenhuma.
