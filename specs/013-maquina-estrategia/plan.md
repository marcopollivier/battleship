# 013 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

1. **Domínio (`internal/board/machine.go`):** uma função nova, `SmartShot`, ao lado do `RandomShot`, com a mesma assinatura. As duas cabem num tipo de função, `ShotStrategy`, que é o "nível" da máquina. O `SmartShot` lê o tabuleiro só pelo que é público, ou seja, o estado de cada casa (`Cell`) e as casas dos navios afundados, e escolhe entre procura e cerco.
2. **Terminal (`internal/terminal/machine.go`):** o `PlayVsMachine` recebe a `ShotStrategy` e a usa no lugar do `RandomShot` fixo. O `*rand.Rand` continua sendo um só, para o `aleatorio` do posicionamento e para os tiros.
3. **Programa (`cmd/init.go`):** a opção `--nivel` (`dificil`, padrão, ou `facil`) escolhe `SmartShot` ou `RandomShot`.

## Pacotes e API

`internal/board/machine.go`:

```go
// ShotStrategy chooses the machine's next shot at b, drawing with rng when
// there is more than one choice. It does not fire.
type ShotStrategy func(b *Board, rng *rand.Rand) (Coord, error)

// SmartShot hunts at random until it hits a ship, then shoots next to the
// hits until the ship sinks, as a player on paper would. It uses only what
// the shooter knows: the state of each cell and which ships have sunk.
func SmartShot(b *Board, rng *rand.Rand) (Coord, error)
```

- `RandomShot` não muda, então `--nivel facil` repete os tiros de antes (CA-8).
- **Conhecimento público:** um navio afundado foi anunciado ao atirador, e como os navios não se encostam, ele ocupa exatamente as casas atingidas ligadas ao último tiro. Por isso o `SmartShot` pode usar as casas dos navios afundados (`sunkCells`, interna, a partir de `Ships` e do `sunk` que já existe). Ele nunca olha as casas de navios que não afundaram.
- **Acertos pendentes:** casas `CellHit` fora dos navios afundados.
- **Cerco** (há acertos pendentes):
  1. Se há dois ou mais acertos pendentes, todos na mesma fileira ou coluna, as candidatas são as duas pontas da linha: a casa antes do primeiro e a depois do último, dentro do tabuleiro e ainda não atingidas.
  2. Se não há ponta disponível, ou os acertos não estão em linha (barra e haste do T), as candidatas são as vizinhas acima, abaixo, à esquerda e à direita de qualquer acerto pendente, ainda não atingidas.
- **Procura** (sem acertos pendentes): as casas não atingidas que não encostam, nem pela diagonal, num navio afundado. Se todas encostam, o que só acontece com o tabuleiro já perdido, vale qualquer casa não atingida.
- Nos dois casos, as candidatas são juntadas em ordem (fileira, depois coluna) sem repetição, e uma é sorteada com `rng.IntN`. Sem casa não atingida, `ErrNoCellsLeft`, como o `RandomShot`.
- As casas excluídas na procura (em volta de navio afundado) também ficam fora do cerco.

`internal/terminal/machine.go`:

```go
func PlayVsMachine(in io.Reader, out io.Writer, player, enemy *board.Board, rng *rand.Rand, shoot board.ShotStrategy) error
```

`cmd/init.go`:

- `level := fs.String("nivel", "dificil", ...)`, comparado com `strings.EqualFold`. Um valor fora de `facil` e `dificil` escreve `Nível inválido: use facil ou dificil.` em `stderr` e retorna 2. É validado no modo solo também, mas só usado contra a máquina.

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| `ShotStrategy` como tipo de função | Uma interface `Machine` com estado; um `enum` de nível passado ao terminal | `RandomShot` já tem a forma certa, e uma função sem estado é fácil de testar. O terminal não precisa saber os nomes dos níveis. |
| `SmartShot` recalcula tudo a partir do tabuleiro | Guardar os acertos numa estrutura da máquina entre os tiros | O tabuleiro já tem todos os tiros. Sem memória, não há estado para manter certo, e o tiro só depende do tabuleiro e do `rng`. |
| Usar as casas dos navios afundados | Reconstruir os navios só pelos acertos ligados | Dá o mesmo resultado, porque os navios não se encostam, e reaproveita o `sunk` do domínio. O teste de "só conhecimento público" confere que navios não afundados não mudam a escolha. |
| Pontas da linha antes das vizinhas | Pesos ou probabilidades por casa | É a regra da questão 3, e cobre o T caindo nas vizinhas quando a linha acaba. |
| Excluir as casas em volta de navio afundado | Deixar a regra para depois | Decisão provisória da questão 2 (issue #25). Fica numa função separada, para sair fácil se a decisão mudar. |

## Estratégia de testes

Testes de domínio em `internal/board/machine_test.go`, com tabuleiros montados à mão e tiros dados com `Fire`. Helpers novos com prefixo `smart`.

| Critério | Teste |
| -------- | ----- |
| CA-1 | `TestSmartShot_Hunt`: sem acertos, muitas sementes só sorteiam casas não atingidas; com um navio afundado, nenhuma semente sorteia as casas em volta dele |
| CA-2 | `TestSmartShot_TargetNeighbors`: um acerto num Destroyer no meio do tabuleiro; para muitas sementes, o tiro é uma das 4 vizinhas; com vizinhas já atingidas ou fora do tabuleiro (acerto no canto), só as que restam |
| CA-3 | `TestSmartShot_TargetLine`: dois acertos num Encouraçado; o tiro é uma das 2 pontas |
| CA-4 | `TestSmartShot_TargetLineBlocked`: barra do T do Porta-Aviões atingida, com as pontas já na água; o tiro é uma vizinha dos acertos (a haste) |
| CA-5 | `TestSmartShot_AfterSunk`: logo depois de afundar um navio, nenhum acerto pendente; o tiro é de procura e não encosta no navio afundado |
| CA-6 | `TestSmartShot_Game`: para 200 sementes, com tabuleiros sorteados (`RandomBoard`), o `SmartShot` atira até `AllSunk`, nunca repete uma casa, e a média de tiros é menor que a do `RandomShot` nas mesmas sementes |
| — | `TestSmartShot_PublicKnowledge`: dois tabuleiros com os mesmos tiros e acertos, mas navios ainda não achados em lugares diferentes, dão o mesmo tiro com a mesma semente |
| CA-7 | `TestRun_SameSeed` (cmd): duas execuções com `--semente 42` e `--nivel dificil` têm a mesma saída |
| CA-8 | Os testes atuais do terminal e do `cmd` contra a máquina, com `RandomShot` e `--nivel facil`, continuam com o mesmo resultado |
| CA-9 | `TestRun_InvalidLevel` (cmd): `--nivel xyz` → código 2 e a mensagem em `stderr` |

No terminal, `TestPlayVsMachine_Strategy`: com uma `ShotStrategy` de teste que sempre escolhe a mesma sequência de casas, a máquina atira nessas casas. Isso confere que o `PlayVsMachine` usa a estratégia recebida.

## Dependências e riscos

- **Conflito com a 010:** as duas specs mudam a assinatura do `PlayVsMachine` e o `cmd/init.go`. Proposta de ordem de merge: 010 antes da 013. A 013 faz merge do `master` e põe a `ShotStrategy` depois do `Style`; a resolução é mecânica.
- **CA-6 pelo valor médio:** a média do `SmartShot` deve ficar bem abaixo da do `RandomShot`. Com 200 sementes fixas, o teste é determinístico.
- **Testes de muitas sementes:** sortear 200 partidas é rápido no domínio (100 casas no máximo por partida).

## Checagem da constituição

- [x] **Domínio separado da interface:** a estratégia fica no domínio, sem I/O; o terminal só recebe a função.
- [x] **Aleatoriedade injetada:** `SmartShot` recebe o `*rand.Rand`, e os testes usam sementes fixas.
- [x] **Testes antes do código:** cada tarefa começa pelo teste.
- [x] **Sem dependências externas:** só biblioteca padrão.
- [x] **Erros como valores:** `ErrNoCellsLeft` e `--nivel` inválido como código 2; nada de `panic`.
- [x] **Até ~8 tarefas:** 6.
- [x] **Idiomas:** código em inglês; a mensagem e a ajuda da opção em português.

Exceções e justificativas: nenhuma.
