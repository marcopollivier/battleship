# 010 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Rascunho

## Abordagem

As cores são só apresentação, então tudo fica em `internal/terminal` e no `cmd`; o domínio não muda.

1. **Estilo (`internal/terminal/color.go`, novo):** um tipo `Style` com dois valores, `Plain` e `Colored`, e uma função interna que envolve um símbolo no código ANSI da sua cor quando o estilo é `Colored`.
2. **Desenho:** `RenderBoard` e `RenderAttackBoard` recebem o `Style` e pintam cada casa ao montar a grade. O `writeBoard`, que escreve o cabeçalho, os números e a legenda, não muda; por isso esses textos ficam sem cor.
3. **Partidas:** `PlaySolo` e `PlayVsMachine` recebem o `Style` e o repassam a todos os tabuleiros que desenham, inclusive no posicionamento (012).
4. **Programa (`cmd/init.go`):** opção `--cores` (`sim`, padrão, ou `nao`), validada logo depois do `--modo`.

## Pacotes e API

`internal/terminal/color.go`:

```go
// Style says whether the boards are drawn with ANSI colors.
type Style int

const (
	Plain   Style = iota // no color codes, the output before spec 010
	Colored              // each board symbol in its color
)
```

- Códigos ANSI, cada símbolo envolvido individualmente e seguido de `\x1b[0m`:

  | Uso | Código |
  | --- | ------ |
  | água `~` | `\x1b[34m` (azul) |
  | navio inteiro | `\x1b[1;37m` (branco em negrito) |
  | acerto `X` | `\x1b[1;31m` (vermelho em negrito) |
  | tiro na água `o` | `\x1b[36m` (ciano) |
  | não atingida `.` | `\x1b[90m` (cinza) |
  | navio afundado (ataque) | `\x1b[31m` (vermelho) |

- `func (s Style) paint(symbol, color string) string`: devolve o símbolo sem mudança no `Plain`.

Assinaturas que mudam:

```go
func RenderBoard(w io.Writer, b *board.Board, style Style) error
func RenderAttackBoard(w io.Writer, b *board.Board, style Style) error
func PlaySolo(in io.Reader, out io.Writer, b *board.Board, style Style) error
func PlayVsMachine(in io.Reader, out io.Writer, player, enemy *board.Board, rng *rand.Rand, style Style) error
```

- O `Style` fica guardado no `playerAttack` e no `machineGame`, e cada chamada interna de desenho usa esse valor.

`cmd/init.go`:

- `colors := fs.String("cores", "sim", ...)`, comparado com `strings.EqualFold`. Um valor fora de `sim` e `nao` escreve `Cores inválidas: use sim ou nao.` em `stderr` e retorna 2.
- `play` recebe o `Style` e o repassa.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| `Style` como parâmetro explícito | Variável global de cores; um `io.Writer` que pinta a saída depois de pronta | A constituição pede o domínio e a interface testáveis sem estado global. Pintar a saída pronta exigiria reconhecer os símbolos no meio do texto. |
| Pintar ao montar a grade, antes do `writeBoard` | Pintar no `writeBoard` | O `writeBoard` não sabe se `S` é um navio inteiro ou afundado; a grade sabe. E o cabeçalho e a legenda ficam sem cor sem nenhum cuidado extra. |
| Cada símbolo com o seu código e o seu reset | Um código só para uma sequência de casas da mesma cor | Mais simples e fácil de testar; o tamanho extra da saída não importa num terminal. |
| `Plain` como valor zero | `Colored` como valor zero | Um `Style` esquecido não espalha códigos de cor; o padrão "cores ligadas" fica no `cmd`, onde a spec o define. |
| Ordem de validação: `--modo`, depois `--cores`, depois a frota | Validar as cores por último | Opções inválidas (código 2) vêm antes da frota inválida (código 1), como hoje. |

## Estratégia de testes

Os testes que já existem passam a usar `Plain` e não mudam de resultado: isso confere a regra "com as cores desligadas, a saída é exatamente a de hoje". Arquivo novo: `internal/terminal/color_test.go`, com um helper `stripColors` que tira os códigos ANSI.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestRenderBoard_Colored`: tabuleiro com água, navios, acerto e tiro na água; cada casa conferida com o código da tabela |
| CA-2 | `TestRenderAttackBoard_Colored`: casas não atingidas, água, acerto e navio afundado |
| CA-3 | `TestRender_ColoredMatchesPlain`: nos dois desenhos, `stripColors(colorido) == sem cor`, e o cabeçalho, os números e a legenda não têm `\x1b` |
| CA-4 | `TestRun_Colors` (cmd): com `--cores nao`, a saída não tem `\x1b`; os testes atuais do `cmd`, com `--cores nao`, continuam com o resultado de hoje |
| CA-5 | `TestRun_Colors` (cmd): sem a opção, com `sim` e com `SIM`, a saída tem os tabuleiros com cores, montados com `Colored` |
| CA-6 | `TestRun_InvalidColors` (cmd): `--cores xyz` → código 2 e a mensagem em `stderr` |
| CA-7 | `TestPlaySolo_Colored`, `TestPlayVsMachine_Colored` e `TestPlayVsMachine_PlacementColored`: a saída usa os tabuleiros desenhados com `Colored` |

## Dependências e riscos

- **Conflito com a 013:** as duas specs mudam a assinatura do `PlayVsMachine` e o `cmd/init.go`. Proposta de ordem de merge: 010 antes da 013; a 013 faz merge do `master` e acrescenta o seu parâmetro depois do `Style`. A resolução é mecânica.
- **Muitas chamadas para mudar:** cerca de 50 chamadas nos testes de `internal/terminal` e do `cmd` ganham o `Plain`. A mudança fica numa tarefa só, junto com a mudança de assinatura.

## Checagem da constituição

- [x] **Domínio separado da interface:** o domínio não muda.
- [x] **Aleatoriedade injetada:** as cores não usam aleatoriedade, e a semente repete a partida igual.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** códigos ANSI escritos à mão, só biblioteca padrão.
- [x] **Erros como valores:** `--cores` inválido vira código 2 com mensagem, sem `panic`.
- [x] **Até ~8 tarefas:** 5.
- [x] **Idiomas:** código em inglês; a mensagem e a ajuda da opção em português.

Exceções e justificativas: nenhuma.
