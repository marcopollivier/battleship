# 018 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

1. **Domínio (`internal/board/master.go`, novo):** `MasterShot`, uma terceira `ShotStrategy` ao lado do `RandomShot` e do `SmartShot` da 013.
   - Ela lista as posições possíveis de cada navio que falta afundar, considerando o que o atirador sabe.
   - Para cada casa não atingida, soma o peso das posições que passam por ela.
   - Atira numa das casas de maior soma, sorteando entre as empatadas.
2. **Terminal:** não muda, porque o `PlayVsMachine` já recebe a `ShotStrategy`.
3. **Programa (`cmd/init.go`):** `--nivel mestre` escolhe o `MasterShot`, e a mensagem de nível inválido cita os três níveis.

## Pacotes e API

`internal/board/master.go`:

```go
// MasterShot shoots where the ships still afloat most likely are: it counts,
// for each cell not shot yet, the ways the remaining ships could cover it
// given every shot so far, and picks a cell with the highest count, drawing
// among ties with rng. Like SmartShot, it uses only what the shooter knows.
func MasterShot(b *Board, rng *rand.Rand) (Coord, error)
```

**O que a máquina sabe:**
- a frota (`b.fleet`);
- o estado de cada casa (`Cell`);
- as casas dos navios afundados (`sunkCells`, da 013).

Os navios que faltam são a frota menos os tipos dos navios afundados, um por navio.

**Posição possível** de um navio que falta: as casas do `ShipCells` para uma casa de referência e uma orientação válida. Ela precisa obedecer a três regras:
- todas as casas ficam dentro do tabuleiro;
- nenhuma casa é água (`CellMiss`) nem de navio afundado;
- nenhuma casa encosta, nem pela diagonal, numa casa atingida que não faz parte da posição. Esse "atingida" vale para acerto pendente ou de navio afundado, porque navios não se encostam.

O Submarino tem uma casa só e conta uma vez, não uma por orientação.

**Peso de cada posição:**
- **Procura**, sem acertos pendentes: cada posição vale 1.
- **Cerco**, com acertos pendentes: só contam as posições que passam por pelo menos um acerto pendente. Cada uma vale o número de acertos pendentes que cobre. Assim, depois de dois acertos em linha, as posições que seguem a linha pesam mais.

**Escolha:**
- Só as casas não atingidas recebem soma.
- A máquina junta, em ordem (fileira, depois coluna), as casas com a maior soma e sorteia uma com `rng.IntN`.
- Se nenhuma casa tem soma (um tabuleiro que não bate com a frota, o que não acontece numa partida), ela usa o `SmartShot`.
- Sem casa não atingida, `ErrNoCellsLeft`.

`cmd/init.go`:

- `masterLevel = "mestre"`, no `chooseStrategy`.
- A mensagem passa a ser `Nível inválido: use facil, dificil ou mestre.`, e a ajuda da opção cita o `mestre`.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| Contar posições (densidade) | Xadrez na procura | Decisão da questão 2. A contagem já evita as casas onde nada cabe e dá a direção no cerco. |
| No cerco, peso pelo número de acertos cobertos | Exigir que a posição cubra todos os acertos pendentes; peso 1 | Os acertos pendentes podem ser de dois navios: a contagem às vezes manda atirar numa casa longe do primeiro acerto e acha outro navio. Exigir todos os acertos descartaria a posição certa. O peso faz a linha de acertos valer mais, como uma pessoa faria. |
| Posição não encosta em acerto fora dela | Só não passar por água e por navios afundados | É a regra de navios que não se encostam aplicada ao que a máquina já viu. Tira posições impossíveis e inclui naturalmente a regra da issue #25 (questão 3). |
| Recalcular tudo a cada tiro, sem memória | Guardar as contagens entre tiros | Como o `SmartShot`, o tiro só depende do tabuleiro e do `rng`. São no máximo 5 tipos × 100 casas × 4 orientações por tiro, rápido o bastante. |
| Novo arquivo `master.go` | Tudo no `machine.go` | O `machine.go` já tem o `RandomShot` e o `SmartShot`. A contagem tem várias funções próprias. |

## Estratégia de testes

Testes em `internal/board/master_test.go`, com helpers de prefixo `master`, e tabuleiros montados à mão como os da 013.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestMasterShot_Hunt`: tabuleiro sem tiros; para muitas sementes, o tiro está entre as casas de maior contagem, conferidas contra uma contagem feita no teste |
| CA-2 | `TestMasterShot_NoRoom`: um canto cercado de água onde nenhum navio que falta cabe nunca é escolhido; com só Submarinos afundados, a frota sem eles muda a contagem |
| CA-3 | `TestMasterShot_Target`: um acerto pendente no meio dá uma vizinha; um acerto encostado na borda ou com água de um lado dá uma vizinha do lado que cabe; dois acertos em linha dão uma ponta da linha |
| CA-4 | `TestMasterShot_AfterSunk`: navio afundado sai da contagem, e as casas em volta dele nunca são escolhidas |
| CA-5 | `TestMasterShot_PublicKnowledge`: mesmos tiros e resultados, navios ainda não achados em lugares diferentes, mesmo tiro com a mesma semente |
| CA-6 | `TestMasterShot_Game`: 200 sementes com `RandomBoard`; nunca repete, termina, e a média é menor que a do `SmartShot` nas mesmas sementes. As médias ficam no encerramento |
| CA-7 | `TestRun_Level` (cmd): `mestre` e `MESTRE` jogam como uma partida com `MasterShot`; a mesma semente dá a mesma saída |
| CA-8 | `TestRun_InvalidLevel` (cmd): a mensagem nova e o código 2 |

## Dependências e riscos

- **Média mais baixa que a `dificil` (CA-6):** a literatura de Batalha Naval indica que a contagem ganha com folga do cerco simples. Se a média não ficar abaixo, a implementação para e a spec é revista.
- **Conflito com a 015:** as duas mexem no `cmd/init.go` (`--cores` e `--nivel`) e no README. A ordem combinada é 014, 018, 015, e a 015 resolve a fusão.

## Checagem da constituição

- [x] **Domínio separado da interface:** a estratégia fica no domínio, sem I/O.
- [x] **Aleatoriedade injetada:** `MasterShot` recebe o `*rand.Rand`, usado só nos empates.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** só biblioteca padrão.
- [x] **Erros como valores:** `ErrNoCellsLeft` e o nível inválido como código 2.
- [x] **Até ~8 tarefas:** 4.
- [x] **Idiomas:** código em inglês; a mensagem e a ajuda da opção em português.

Exceções e justificativas: nenhuma.
