# 022 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

Tudo fica no terminal e no `cmd`; o domínio não muda.

1. **Desenho em linhas (`render.go`, `attack.go`):** separar a grade (cabeçalho e 10 fileiras) da legenda. Funções internas devolvem as linhas da grade. `RenderBoard` e `RenderAttackBoard` continuam escrevendo exatamente o que escrevem hoje, montados a partir delas.
2. **Tela (`internal/terminal/screen.go`, novo):**
   - `sideBySide` junta duas grades linha a linha;
   - o código de limpar o terminal;
   - um tipo `Display`, que junta o `Style` da 010 e a escolha de limpar a tela.
3. **Partidas:** o `PlayVsMachine`, o `PlaySolo` e o posicionamento passam a desenhar telas completas. Cada tela tem a limpeza (se ligada), os tabuleiros, as mensagens do turno e o pedido.
4. **Programa (`cmd/init.go`):** a opção `--limpar auto|sim|nao`, com a mesma detecção de terminal da 015.

## Pacotes e API

`internal/terminal/screen.go`:

```go
// Display is how a game is shown: the colors of the boards and messages,
// and whether the terminal is cleared before each screen.
type Display struct {
	Style Style
	Clear bool
}

// clearScreen moves the cursor home and erases the terminal.
const clearScreen = "\x1b[H\x1b[2J"

// sideBySide joins two board grids line by line, the left one padded to
// its visible width plus a gap, so color codes do not shift the columns.
func sideBySide(left, right []string) []string
```

- A largura visível de uma linha de grade é sempre 22 colunas: 2 do número da fileira e 20 das casas. O espaçamento entre as grades é de 5 espaços, então a linha toda tem 49 colunas.
- Os títulos ficam em cima: `Sua frota` com 22 colunas, mais o espaçamento, e depois `Frota inimiga`.

Assinaturas que mudam (o `Style` vira `Display`):

```go
func PlaySolo(in io.Reader, out io.Writer, b *board.Board, display Display) error
func PlayVsMachine(in io.Reader, out io.Writer, player, enemy *board.Board, rng *rand.Rand, display Display, shoot board.ShotStrategy) error
```

`RenderBoard` e `RenderAttackBoard` não mudam.

**Tela da partida contra a máquina:**
1. limpeza, se `Display.Clear`;
2. os títulos e as grades lado a lado;
3. uma linha em branco, a legenda `Sua frota: ...` e a legenda `Frota inimiga: ...`;
4. uma linha em branco;
5. as mensagens do turno, guardadas desde a tela anterior:
   - `Seu disparo em B7: <resultado>`;
   - as linhas da máquina;
6. o pedido `Disparo: `.

- Mensagens de entrada inválida continuam escritas na hora, abaixo do pedido, e o pedido se repete sem tela nova.
- O resultado do jogador é `Seu disparo em <casa>: ` mais a mensagem da 009, pintada como na 015.
- **Fim da partida:** a tela final tem os títulos `Sua frota` e `Frota inimiga`, com o `RenderBoard` (navios à mostra) dos dois lados, as mensagens do turno e a mensagem final. A frase `A frota inimiga era:` sai, porque a frota já aparece na tela.

**Modo solo:**
1. limpeza;
2. o `RenderAttackBoard`;
3. uma linha em branco;
4. `Seu disparo em ...`;
5. o pedido.

No fim, `RenderBoard` do alvo e a mensagem final.

**Posicionamento:**
1. limpeza;
2. as instruções;
3. o aviso da 014, se houver;
4. o `RenderBoard` do jogador;
5. o pedido do navio.

Mensagens de erro continuam abaixo do pedido.

`cmd/init.go`: `clear := fs.String("limpar", "auto", ...)`. `chooseClear(value string, tty bool) (bool, bool)` é uma função pura, como o `chooseStyle`. Um valor inválido mostra `Limpar inválido: use sim, nao ou auto.` e retorna 2.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| `Display` juntando `Style` e `Clear` | Mais um parâmetro solto no `PlayVsMachine` | O `PlayVsMachine` já tem 7 parâmetros. As próximas opções de exibição entram no mesmo tipo, sem mudar a assinatura de novo. |
| Largura fixa de 22 colunas para alinhar | Medir o texto tirando os códigos de cor | Toda linha de grade tem a mesma largura visível, e o teste confere isso. |
| Guardar as mensagens do turno e escrevê-las depois dos tabuleiros | Escrever na hora e redesenhar | Com a tela limpa, o que foi escrito antes da limpeza some. As mensagens precisam vir depois dos tabuleiros novos. |
| Erros de entrada sem tela nova | Redesenhar tudo a cada erro | Um erro não muda o tabuleiro, e redesenhar faria a mensagem piscar. |
| `RenderBoard` e `RenderAttackBoard` intactos | Mudar o formato deles | São usados nas revelações e nos testes. A grade em linhas é interna. |

## Estratégia de testes

Os testes de partida montam a saída esperada com funções novas de teste: `screenMachine(player, enemy, events)` e `screenSolo(target, events)`. Elas usam o `sideBySide` e as linhas da grade. Os testes atuais de `PlayVsMachine`, `PlaySolo` e do posicionamento mudam o resultado esperado, porque o formato da tela muda.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestSideBySide`: 11 linhas, colunas alinhadas com e sem cores, nenhuma linha acima de 80 colunas visíveis; `TestPlayVsMachine_Start`: a tela inicial |
| CA-2 | `TestPlayVsMachine_HitShootsAgain`: depois de um acerto, nova tela com `Seu disparo em E5: Acerto!` |
| CA-3 | `TestPlayVsMachine_MachineTurn`: depois de uma água, a tela seguinte tem o resultado do jogador e as linhas da máquina |
| CA-4 | `TestPlayVsMachine_InvalidInput`: a mensagem abaixo do pedido, sem tela nova |
| CA-5 | `TestPlayVsMachine_Clear` e `TestPlaySolo_Clear`: com `Clear`, cada tela começa com `clearScreen`; sem, nenhum `\x1b[2J`. `TestChooseClear` (cmd, tabela) |
| CA-6 | `TestPlayVsMachine_PlayerWins`, `_MachineWins` e `_GiveUp`: a tela final com os dois tabuleiros revelados |
| CA-7 | `TestPlaySolo_*`: tabuleiro único, resultado com a casa |
| CA-8 | `TestPlayVsMachine_Placement*`: com `Clear`, cada navio começa com `clearScreen` e as instruções |
| CA-9 | `TestRun_InvalidClear` |
| CA-10 | `TestRun_SameSeed`, com `--limpar sim` e `--cores sim` |

## Dependências e riscos

- **Tamanho:** 7 tarefas, perto do limite. A mudança de formato quebra muitos testes de uma vez. Por isso a T1 e a T3 são preparações que não mudam a saída, e cada tarefa seguinte muda uma tela só.
- **Conflitos:**
  - com a 023, que muda o `RenderBoard` (`#`), pela ordem 024, 023, 022, 025: esta spec traz o `master` depois da 023;
  - com a 025, que precisa do resultado da partida: ela entra depois e se adapta ao `Display`;
  - com a 024, no `Run`.
- **Terminais pequenos:** a tela tem cerca de 20 linhas e 49 colunas, e cabe num terminal padrão de 80×24.

## Checagem da constituição

- [x] **Domínio separado da interface:** só o terminal e o `cmd` mudam.
- [x] **Aleatoriedade injetada:** não se aplica; a semente continua repetindo a partida.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** códigos ANSI escritos à mão.
- [x] **Erros como valores:** opção inválida vira código 2.
- [x] **Até ~8 tarefas:** 7.
- [x] **Idiomas:** código em inglês; textos em português.

Exceções e justificativas: nenhuma.
