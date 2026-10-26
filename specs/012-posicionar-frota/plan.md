# 012 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

1. **Domínio (`internal/board`):** três funções pequenas que o posicionamento precisa e que o domínio ainda não tem:
   - `Remaining`, que diz quais navios ainda faltam, do maior para o menor;
   - `Clear`, que tira todos os navios;
   - `FillRandom`, que sorteia só os navios que faltam, mantendo os já posicionados.

   A validação de cada posição continua sendo o `Place` da 002.
2. **Terminal (`internal/terminal/placement.go`, novo):** a fase de posicionamento. Ela lê `casa orientação`, traduz a orientação do português, chama `Place` e mostra as mensagens de erro. Ela roda dentro do `PlayVsMachine`, antes do primeiro disparo, quando o tabuleiro do jogador ainda não está completo. Assim a partida usa um único leitor da entrada, e as linhas não se perdem entre o posicionamento e os disparos.
3. **Programa (`cmd/init.go`):** no modo contra a máquina, sorteia primeiro o tabuleiro da máquina e passa ao `PlayVsMachine` um tabuleiro vazio para o jogador, com a frota escolhida. O mesmo `*rand.Rand` atende depois o `aleatorio` e os disparos da máquina.

## Pacotes e API

`internal/board/board.go`:

```go
// Remaining returns the ships of the fleet that are not placed yet, one
// entry per ship, largest first.
func (b *Board) Remaining() []ShipType

// Clear removes every ship and every shot from the board.
func (b *Board) Clear()
```

- `Remaining` reaproveita a ordenação do `shipsBySize` da 005 (maior primeiro, estável) e desconta os navios já posicionados de cada tipo.

`internal/board/random.go`:

```go
// FillRandom places the ships that are still missing at random, keeping the
// ships already on b. On ErrFleetDoesNotFit, b is left unchanged.
func FillRandom(b *Board, rng *rand.Rand) error
```

- Usa as mesmas tentativas do `RandomBoard` (`maxRestarts` vezes, `placeRandomly` para cada navio) numa cópia de `b` com os navios já posicionados. Só copia o resultado para `b` quando todos os navios couberem.
- O `RandomBoard` passa a ser `NewBoard` mais `FillRandom`, sem mudar os tabuleiros que ele sorteia para uma mesma semente.

`internal/terminal/placement.go`:

```go
// parsePlacement reads "E5 vertical", "e5 v" or, for a one-cell ship, "E5".
func parsePlacement(line string, t board.ShipType) (board.Coord, board.Orientation, error)
```

- A orientação vem do inverso do `orientationNames` de `names.go`, pelo nome inteiro ou pela primeira letra, sem diferenciar maiúscula e minúscula. As seis primeiras letras são distintas.
- Erros: `errInvalidPlacement` (formato), os erros do `ParseCoord`, e `board.ErrInvalidOrientation` (orientação desconhecida ou que não vale para o navio).
- O laço é um método do `machineGame` (`placeFleet`), chamado pelo `PlayVsMachine` antes da tela da partida quando `player.Complete()` é falso. Ele retorna `false` se o jogador desistir. O `PlayVsMachine` mantém a assinatura, e os testes da 007, que passam tabuleiros completos, continuam iguais.
- **Saída:**
  - as três linhas de instruções e uma linha em branco;
  - a cada navio, `RenderBoard(player)`, uma linha em branco e o pedido (`Porta-Aviões, 5 casas (cima, baixo, esquerda, direita): ` ou `Submarino, 1 casa: `);
  - numa entrada inválida, a mensagem e o mesmo pedido, sem tabuleiro, como nos disparos;
  - no fim, `Frota posicionada!`, uma linha em branco e a tela da partida da 007.
- Mensagens de erro numa função interna, `placementErrorMessage(err, t)`. A lista de orientações usa vírgulas no pedido e "ou" antes da última na mensagem (`use cima, baixo, esquerda ou direita`).
- **Comandos:**
  - `aleatorio` chama `FillRandom` com o `rng` da partida;
  - `recomecar` chama `Clear`;
  - `sair` e o fim da entrada (com a quebra de linha antes, como no solo) mostram `Você desistiu antes da partida.` e encerram.

`cmd/init.go`:

- No `play`, modo `maquina`: `enemy := RandomBoard(fleet, rng)`, `player := board.NewBoard(fleet)`, `PlayVsMachine(stdin, stdout, player, enemy, rng)`. O modo `solo` não muda.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| Posicionamento dentro do `PlayVsMachine` | Uma função `PlaceFleet` separada, chamada pelo `cmd` | Um `bufio.Scanner` lê adiante. Dois leitores sobre o mesmo `stdin` perderiam linhas. Passar o leitor entre as funções mudaria a assinatura do `PlayVsMachine` e do `PlaySolo`. |
| Começar pelo posicionamento quando o tabuleiro do jogador não está completo | Uma opção explícita no `PlayVsMachine` | Um tabuleiro incompleto não serve para a partida, então é o sinal natural. Os testes da 007 continuam valendo sem mudança. |
| `FillRandom` no domínio | Sortear no terminal | Sortear navios é regra do jogo, já está no domínio (005), e o `RandomBoard` passa a reaproveitá-lo. |
| `FillRandom` trabalha numa cópia | Posicionar direto em `b` e desfazer na falha | O tabuleiro não muda numa falha, como o `Place` (002). |
| `Clear` também apaga os disparos | Só os navios | Um tabuleiro sem navios com disparos não faz sentido. No posicionamento ainda não há disparos. |
| Tabuleiro da máquina sorteado primeiro | Manter a ordem da 007 | O tabuleiro do jogador agora depende do que ele digita. Sorteando a máquina antes, o tabuleiro da máquina só depende da semente. É a decisão da questão 5 da spec. |

## Estratégia de testes

Arquivos novos:
- `internal/terminal/placement_test.go`, com helpers de prefixo `placement`;
- os testes de domínio entram em `board_test.go` (`Remaining` e `Clear`) e `random_test.go` (`FillRandom`).

Os testes do terminal montam a saída esperada como os da 006 e da 007: com o tabuleiro de réplica e as funções de desenho. A frota completa usa as posições do tabuleiro de exemplo da 003.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestPlayVsMachine_PlacementStart`: com a entrada vazia, a saída começa com as instruções, o tabuleiro vazio e o pedido do Porta-Aviões |
| CA-2 | `TestPlayVsMachine_PlaceFleet`: as 8 posições do tabuleiro de exemplo, um pedido e um tabuleiro antes de cada navio, `Frota posicionada!` e a tela da 007 com os navios do jogador nesses lugares |
| CA-3 | `TestParsePlacement`: tabela com `e5 V`, `E5 vertical`, `C3 b`, `C3 BAIXO`, `A1` para um Submarino e `A1 h` para um Submarino |
| CA-4 | `TestParsePlacement_Invalid` (formato e orientação) e `TestPlayVsMachine_PlacementErrors`: cada linha da tabela de mensagens, seguida do mesmo pedido |
| CA-5 | `TestFillRandom`: os navios já posicionados continuam, os que faltam entram, o tabuleiro fica completo e válido (`assertValidBoard` da 005). `TestPlayVsMachine_PlacementRandom`: `aleatorio` depois de um navio termina o posicionamento |
| CA-6 | `TestFillRandom_DoesNotFit`: o tabuleiro não muda e o erro é `ErrFleetDoesNotFit`. `TestPlayVsMachine_PlacementRandom`: uma frota de 2 Encouraçados com o primeiro bloqueando o segundo mostra a mensagem e pede o mesmo navio |
| CA-7 | `TestBoard_Clear` e `TestPlayVsMachine_PlacementRestart`: `recomecar` mostra o tabuleiro vazio e o pedido do Porta-Aviões |
| CA-8 | `TestPlayVsMachine_PlacementGiveUp` (`sair` e fim da entrada) e `TestRun_PlacementGiveUp` (cmd): `Você desistiu antes da partida.`, `Semente: N` e código 0 |
| CA-9 | `TestRun_SameSeed` (cmd), com entradas que posicionam navios, usam `aleatorio` e atiram |
| CA-10 | `TestRun_Mode` (cmd): `--modo solo` continua começando pelo tabuleiro de ataque |

Os testes do `cmd` da 007 (`TestRun_PlayerWins`, `TestRun_MachineWins`, `TestRun_SameSeed`) passam a começar a entrada com `aleatorio`, e o helper `machineBoards` passa a sortear a máquina primeiro.

## Dependências e riscos

- **CA-6 precisa de uma frota que trave:** com 2 Encouraçados, o primeiro posto no meio do tabuleiro, na horizontal, ainda deixa espaço para o segundo. O teste usa uma frota em que o primeiro navio, bem posicionado, impede o resto, por exemplo 3 Encouraçados com o primeiro em `A5 horizontal` e mais navios grandes. O valor exato sai do teste, e o teste confere antes que a frota trava mesmo.
- **`RandomBoard` sobre `FillRandom`:** precisa sortear exatamente os mesmos tabuleiros de antes para a mesma semente, ou os testes do modo solo da 006 mudariam. Os testes de `RandomBoard` e da 006 conferem isso.
- **Tamanho:** são 7 tarefas, perto do limite de ~8.

## Checagem da constituição

- [x] **Domínio separado da interface:** `Remaining`, `Clear` e `FillRandom` no domínio, sem I/O. O terminal só lê, traduz e chama `Place`.
- [x] **Aleatoriedade injetada:** `FillRandom` recebe o `*rand.Rand`.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** só biblioteca padrão.
- [x] **Erros como valores:** os erros sentinela do `Place` e o `ErrFleetDoesNotFit` viram mensagens; nada de `panic`.
- [x] **Até ~8 tarefas:** 7.
- [x] **Idiomas:** mensagens em português, código em inglês.

Exceções e justificativas: nenhuma.
