# 014 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

1. **Domínio (`internal/board`):** duas funções novas.
   - `RemoveLast` tira o último navio posicionado.
   - `Fits` diz se os navios que faltam ainda cabem no tabuleiro. Ela faz uma busca exata, sem aleatoriedade.
2. **Terminal (`internal/terminal/placement.go`):**
   - o comando `desfazer` e a linha nova nas instruções;
   - antes de desenhar o tabuleiro e pedir o próximo navio, o aviso, quando `Fits` é falso.
3. **Programa:** o `cmd` não muda. README e `CLAUDE.md` citam o comando novo.

## Pacotes e API

`internal/board/board.go`:

```go
// RemoveLast removes the ship placed last and reports false if there is none.
func (b *Board) RemoveLast() bool
```

`internal/board/fit.go` (novo):

```go
// Fits reports whether the ships of the fleet not placed yet can still all
// be placed on b, following the placing rules. It searches every way to
// place them, without randomness, and leaves b unchanged.
func Fits(b *Board) bool
```

- **Busca:** backtracking sobre uma cópia (`clone`).
  - Os navios que faltam vêm do `Remaining`, do maior para o menor.
  - Cada navio tenta as casas em ordem (fileira, depois coluna) e as orientações válidas, chamando `Place`. Se não houver lugar para um navio, desfaz o anterior (`RemoveLast`) e tenta a próxima posição dele.
- **Poda:** navios do mesmo tipo são intercambiáveis. Por isso o segundo navio de um tipo só tenta posições depois da posição do primeiro. Isso corta as repetições da mesma arrumação em outra ordem.
- **Limite de segurança:** a busca conta as tentativas de `Place`. Com a frota da cartela, o pior caso medido deve ficar bem abaixo do limite (por exemplo 1.000.000 de tentativas). Se uma frota grande do `--frota` passar do limite, `Fits` responde `true`. Assim o aviso nunca aparece por engano, só deixa de aparecer nesse caso extremo (veja Decisões).

`internal/terminal/placement.go`:

- `undoCommand = "desfazer"`. No `placeShip`, o caso de `desfazer`:
  - com `RemoveLast` verdadeiro, volta ao laço do `placeFleet`, que desenha o tabuleiro e pede o navio que saiu. Como os navios são pedidos do maior para o menor, o primeiro que falta é sempre o último que saiu;
  - com `RemoveLast` falso, mostra `Não há navio para desfazer.` e pede o mesmo navio.
- No `placeFleet`, antes do `RenderBoard` de cada navio: se há navios posicionados e `!board.Fits(g.player)`, mostra `Atenção: o resto da frota não cabe mais no tabuleiro. Use desfazer ou recomecar.`. O aviso aparece de novo a cada tabuleiro desenhado enquanto o resto não couber, e some quando volta a caber (CA-6).
- `placementInstructions`: a linha de comandos passa a ser `Comandos: aleatorio sorteia o resto da frota, desfazer tira o último navio, recomecar tira todos os navios, sair desiste.`.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| `Fits` no domínio | Conferir no terminal | Encaixar navios é regra do jogo, como o `FillRandom`. |
| Busca exata com backtracking | Conta aproximada de casas livres; tentar o `FillRandom` | Decisão da questão 4. O `FillRandom` usa a semente e pode falhar por azar, o que mudaria a partida e daria avisos falsos. |
| Limite de tentativas, respondendo `true` quando passa | Sem limite | Uma frota grande do `--frota` poderia travar o jogo numa busca longa. Responder `true` mantém a promessa de que o aviso só aparece quando é certo. Com a frota da cartela, o limite nunca é atingido, e um teste confere isso. |
| Conferir antes de cada tabuleiro, não só depois de posicionar | Conferir só depois de um `Place` | Cobre o `desfazer` sem código extra (CA-6) e mantém o aviso visível enquanto o problema existe. |
| `RemoveLast` no domínio | Guardar a ordem no terminal | O `Board` já guarda os navios na ordem em que foram posicionados. |

## Estratégia de testes

| Critério | Teste |
| -------- | ----- |
| CA-1, CA-2 | `TestBoard_RemoveLast` (domínio) e `TestPlayVsMachine_PlacementUndo`: depois de 2 navios, `desfazer` mostra o tabuleiro sem o segundo e pede o Encouraçado. Duas vezes voltam ao Porta-Aviões |
| CA-3 | `TestPlayVsMachine_PlacementUndo`: `desfazer` no tabuleiro vazio mostra a mensagem e o pedido do Porta-Aviões |
| CA-4 | `TestFits` (tabela, domínio): tabuleiros montados à mão em que o resto não cabe. `TestPlayVsMachine_PlacementWarning`: uma posição que bloqueia o resto mostra o aviso antes do tabuleiro |
| CA-5 | `TestFits`: tabuleiro vazio, frota quase completa e um caso apertado que só cabe de um jeito. `TestPlayVsMachine_PlaceFleet` (012) continua sem aviso |
| CA-6 | `TestPlayVsMachine_PlacementWarning`: depois do aviso, `desfazer` e o tabuleiro sem aviso |
| CA-7 | `TestRun_SameSeed` (cmd), com `desfazer` na entrada |
| CA-8 | `TestPlayVsMachine_PlacementStart` (012), atualizado com a linha nova das instruções |
| — | `TestFits_Bounded`: com a frota da cartela, a busca fica abaixo do limite nos piores casos montados no teste, e `Fits` não muda o tabuleiro |

Para um tabuleiro em que o resto não cabe com a frota da cartela: o Porta-Aviões no meio e o Encouraçado e o Cruzador espalhados ainda deixam lugar. O teste usa uma frota pequena do `--frota`, por exemplo 2 Encouraçados num tabuleiro em que o primeiro ocupa uma fileira do meio, e uma frota maior com Submarinos espalhados em xadrez.

## Dependências e riscos

- **Desempenho da busca:** o pior caso é quando o resto quase cabe. A poda dos navios iguais e o limite de tentativas evitam travar. A T2 mede o tempo e o número de tentativas no pior caso com a frota da cartela, e registra o resultado no encerramento.
- **Conflito com a 015:** a 015 colore a mensagem de aviso e a de desfazer, que esta spec cria. A ordem de merge combinada é 014, depois 018, depois 015.

## Checagem da constituição

- [x] **Domínio separado da interface:** `RemoveLast` e `Fits` no domínio, sem I/O.
- [x] **Aleatoriedade injetada:** `Fits` não usa aleatoriedade, e a semente não muda.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** só biblioteca padrão.
- [x] **Erros como valores:** nada de `panic`; `RemoveLast` e `Fits` respondem com `bool`.
- [x] **Até ~8 tarefas:** 5.
- [x] **Idiomas:** código em inglês; mensagens em português.

Exceções e justificativas: nenhuma.
