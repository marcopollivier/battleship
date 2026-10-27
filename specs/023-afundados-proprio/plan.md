# 023 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Rascunho

## Abordagem

Só o desenho do tabuleiro próprio muda (`internal/terminal/render.go`). Depois de marcar os acertos e a água, o `RenderBoard` passa pelos navios afundados e troca as casas deles por `#`, pintado com a cor de navio afundado da 010. A legenda (`Legend`) ganha `# afundado`. Como todos os pontos que desenham o tabuleiro próprio usam o `RenderBoard`, a mudança chega sozinha à tela da partida, à derrota e às revelações. O domínio e o `cmd` não mudam.

## Pacotes e API

`internal/terminal/render.go`:

- Nova constante `sunkSymbol = "#"`.
- No `RenderBoard`, depois do laço dos disparos: para cada navio com `shipSunk(b, ship)` (de `attack.go`), as casas recebem `style.paint(sunkSymbol, sunkColor)`.
- `Legend` passa a ser `Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões, X atingido, # afundado, o tiro na água`.

Nenhuma assinatura muda.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| Reaproveitar o `shipSunk` do tabuleiro de ataque | Perguntar ao domínio | O `shipSunk` já existe no terminal e lê só o estado das casas. |
| Mesma cor (`sunkColor`) do navio afundado do ataque | Uma cor nova | Decisão da questão 2: o vermelho já significa "afundado". |
| Trocar as casas depois de marcar os acertos | Calcular a casa de uma vez | Mantém a ordem atual do `RenderBoard` (água, navios, disparos) e acrescenta só um passo. |

## Estratégia de testes

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestRenderBoard_Sunk`: tabuleiro com um navio afundado, um atingido e um intacto; a grade esperada montada à mão |
| CA-2 | `TestRenderBoard_SunkColored`: o `#` vem com o código vermelho, e `stripColors` dá o desenho sem cores |
| CA-3 | Os testes atuais do `RenderBoard` usam a constante `Legend`, então passam a conferir a legenda nova; um caso confere o texto `# afundado` |
| CA-4 | `TestPlayVsMachine_SunkOwnShip`: a máquina afunda o Submarino do jogador (tabuleiros de teste da 007), e a tela seguinte mostra `#` em `Sua frota` |
| CA-5 | `TestPlaySolo_GiveUpShowsSunk`: depois de afundar um navio e desistir, a frota revelada mostra `#` |

Testes que comparam a saída inteira e têm navios afundados no tabuleiro próprio (derrota da máquina, revelações) mudam o resultado esperado. Eles montam a saída com o próprio `RenderBoard`, então devem continuar passando sem mudança.

## Dependências e riscos

- **Conflito com a 022:** a 022 refaz a tela e separa a grade da legenda no `render.go`. Pela ordem combinada (024, 023, 022, 025), a 023 entra antes, e a 022 traz o `master` com o `#`.

## Checagem da constituição

- [x] **Domínio separado da interface:** só o terminal muda.
- [x] **Aleatoriedade injetada:** não se aplica.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** só biblioteca padrão.
- [x] **Erros como valores:** não se aplica; nenhum erro novo.
- [x] **Até ~8 tarefas:** 2.
- [x] **Idiomas:** código em inglês; legenda em português.

Exceções e justificativas: nenhuma.
