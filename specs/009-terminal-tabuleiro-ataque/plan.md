# 009 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Rascunho

## Abordagem

Tudo fica em `internal/terminal`, sem mudanças no domínio. O domínio já informa o que é preciso:

- `Board.Cell(c)` diz se a casa está intocada, com água atingida ou com navio atingido.
- `Board.Ships()` diz onde estão os navios.

Um navio está afundado quando todas as suas casas estão como `CellHit`, o que o terminal consegue calcular sozinho.

Mudanças:

1. **Grade comum:** a parte do `RenderBoard` que escreve o cabeçalho, as fileiras e a legenda vira uma função interna, `writeBoard(w, grid, legend)`, usada pelos dois tabuleiros.
2. **Tabuleiro próprio:** o `RenderBoard` passa a marcar os disparos recebidos (`X` e `o`) e a usar a legenda nova.
3. **Tabuleiro de ataque:** a função nova `RenderAttackBoard` desenha o tabuleiro do oponente.
4. **Mensagens:** a função nova `ShotMessage` monta o texto do resultado de um disparo.

Na onda 2, a 008 também cria arquivos em `internal/terminal` (`fleet.go`) e lê o mapa `shipSymbols` de `render.go`. Por isso, esta spec **não renomeia nem muda** o `shipSymbols`, e não toca em `cmd/` nem em `main.go`.

## Pacotes e API

`internal/terminal/render.go`, com as mudanças:

```go
// Legend explains the symbols of the player's own board.
const Legend = "Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões, X atingido, o tiro na água"

// RenderBoard writes the player's own board, including the shots it has
// received, followed by a blank line and the legend.
func RenderBoard(w io.Writer, b *board.Board) error

// writeBoard writes the header, the rows of grid, a blank line and legend
// as a single Write.
func writeBoard(w io.Writer, grid *[board.Size][board.Size]string, legend string) error
```

`internal/terminal/attack.go`, novo:

```go
// AttackLegend explains the symbols of the attack board.
const AttackLegend = "Legenda: . não atingida, o água, X acerto, S D C E P navio afundado"

// RenderAttackBoard writes the opponent's board as the attacker sees it:
// ships appear only once hit, and sunk ships show their letter.
func RenderAttackBoard(w io.Writer, b *board.Board) error
```

`internal/terminal/messages.go`, novo:

```go
// ShotMessage returns the Portuguese message for a shot result, e.g. "Afundou o Cruzador!".
func ShotMessage(s board.Shot) string
```

Para um `ShotResult` desconhecido, `ShotMessage` cai no `String()` do domínio, como os outros nomes da 004.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| Calcular "afundado" no terminal a partir de `Cell` e `Ships` | Expor `Board.Sunk(ship)` no domínio | Não muda o domínio nesta onda (a 008 trabalha em `internal/board`), e a regra é uma linha. |
| Extrair `writeBoard` | Duplicar o laço de desenho | Os dois tabuleiros têm o mesmo cabeçalho e as mesmas fileiras. Uma função só garante que continuem iguais. |
| Manter o nome `Legend` para o tabuleiro próprio | Renomear para `OwnLegend` | Evita mexer em quem já usa `Legend` e reduz o diff. |
| Mensagens com o nome de `ShipName` | Textos fixos por tipo | Uma única fonte para os nomes dos navios (004). |

## Estratégia de testes

Arquivos `internal/terminal/attack_test.go` e `internal/terminal/messages_test.go` (novos), e `internal/terminal/render_test.go`, onde os testes de saída exata da 004 são atualizados para a legenda nova.

Para não colidir com os testes da 008 no mesmo pacote, os helpers novos desta spec usam o prefixo `attack` (por exemplo, `attackExampleBoard`, `attackFire`). Os helpers existentes (`render`, `boardWith`, `placement`, `cellAt`) podem ser reutilizados, mas não renomeados.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestRenderAttackBoard_NoShots`: tabuleiro de exemplo, saída exata com 100 `.` e a `AttackLegend` |
| CA-2 | `TestRenderAttackBoard_Shots`: depois de `J10` e `E1`, `cellAt` dá `o` em `J10`, `X` em `E1` e `.` no resto |
| CA-3 | `TestRenderAttackBoard_Sunk`: depois de `J1`, `J2` e `J3`, as três casas são `C` |
| CA-4 | `TestShotMessage`: tabela com `Miss`, `Hit` (com `Ship` de qualquer valor, que não aparece) e `Sunk` para os 5 tipos, mais um resultado desconhecido |
| CA-5 | `TestRenderAttackBoard_Format`: sem espaço no fim das linhas, ASCII acima da legenda, linha em branco antes da legenda |
| CA-6 | `TestRenderBoard_ReceivedShots`: depois de `J10` e `E1`, `cellAt` dá `o` em `J10`, `X` em `E1` e `E` em `F1`, `G1` e `H1`, com a legenda nova |
| CA-7 | `TestRenderBoard_Empty` e `TestRenderBoard_FullFleet` da 004, atualizados só na legenda |

## Dependências e riscos

- **Mudança de saída do tabuleiro próprio:** a legenda nova altera a saída de `go run .`. O teste `TestRun` do `cmd` compara com o próprio `RenderBoard`, então não quebra.
- **Disparos exigem tabuleiro completo (003):** os testes que disparam usam a disposição completa da frota padrão.

## Checagem da constituição

- [x] **Domínio separado da interface:** o terminal só lê o domínio, e o domínio não muda.
- [x] **Aleatoriedade injetada:** não se aplica.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** só a biblioteca padrão.
- [x] **Erros como valores:** os erros de escrita são retornados.
- [x] **Até ~8 tarefas:** a estimativa é de 4.
- [x] **Idiomas:** textos para o jogador em português, código em inglês.

Exceções e justificativas: nenhuma.
