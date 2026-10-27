# 015 — Mais cores

- **Status:** Concluída
- **Roadmap:** #015 (onda 7)
- **Depende de:** 010 (cores no terminal)

## Problema

A 010 coloriu só os tabuleiros, e deixou as cores sempre ligadas. Por isso as mensagens que mais importam (`Acerto!`, `Afundou o Cruzador!`, quem venceu) passam despercebidas no meio do texto. E quem grava a partida num arquivo, ou usa um terminal sem cores, recebe os códigos de cor no meio do texto, a não ser que lembre de passar `--cores nao`. O encerramento da 010 deixou as duas melhorias para o roadmap.

## Regras

**Mensagens coloridas (questão 1):**

- Com as cores ligadas, estas mensagens ganham cor, com as mesmas cores dos tabuleiros quando falam da mesma coisa:

  | Mensagem | Cor |
  | -------- | --- |
  | `Água.` | ciano |
  | `Acerto!` | vermelho em negrito |
  | `Afundou o <navio>!` e `Afundou o seu <navio>!` | vermelho em negrito |
  | `Você venceu! ...` | verde em negrito |
  | `A máquina venceu! ...` | vermelho em negrito |
  | mensagens de entrada inválida (coordenada, tiro repetido, posição inválida, aviso da 014) | amarelo |

- Nas linhas de disparo da máquina (`A máquina atirou em C4: Acerto!`), só o resultado depois dos dois-pontos ganha cor.
- Os pedidos (`Disparo: `, `Porta-Aviões, 5 casas ...`), as instruções, os títulos (`Sua frota:`) e a desistência continuam sem cor.
- Com as cores desligadas, as mensagens são exatamente as de hoje.

**Cores desligadas sozinhas (questão 2):**

- A opção `--cores` ganha um terceiro valor, `auto`, que passa a ser o padrão:
  - com `auto`, as cores aparecem quando a saída é um terminal, e não aparecem quando a saída vai para um arquivo ou outro programa (`go run . > partida.txt`, `go run . | less`);
  - com `auto`, a variável de ambiente `NO_COLOR` definida (com qualquer valor que não seja vazio) desliga as cores, como pede a convenção de no-color.org;
  - `--cores sim` liga sempre, e `--cores nao` desliga sempre, como hoje.
- Um valor inválido mostra `Cores inválidas: use sim, nao ou auto.` e termina com código 2.
- Isso muda a decisão da 010 (cores sempre ligadas por padrão). Quem quiser o comportamento da 010 usa `--cores sim`.

## Critérios de aceitação

- **CA-1:** Dadas as cores ligadas, quando o jogador acerta, erra e afunda um navio, então `Acerto!`, `Água.` e `Afundou o <navio>!` aparecem com as cores da tabela.
- **CA-2:** Dadas as cores ligadas, quando a máquina atira, então só o resultado de cada linha `A máquina atirou em ...` tem cor.
- **CA-3:** Dadas as cores ligadas, quando o jogador vence ou a máquina vence, então a mensagem final tem a cor da tabela.
- **CA-4:** Dadas as cores ligadas, quando o jogador digita uma coordenada inválida, um tiro repetido ou uma posição inválida, então a mensagem aparece em amarelo.
- **CA-5:** Dadas as cores ligadas, quando o programa mostra pedidos, instruções e títulos, então eles não têm códigos de cor; e tirando os códigos de cor, a saída é igual à saída sem cores.
- **CA-6:** Dado `--cores auto` (ou nenhuma opção), quando a saída não é um terminal, então não há códigos de cor; quando é um terminal e `NO_COLOR` não está definida, há cores; quando `NO_COLOR` está definida, não há cores.
- **CA-7:** Dado `--cores sim` com a saída num arquivo e `NO_COLOR` definida, quando a partida roda, então há cores; e dado `--cores nao`, nunca há cores.
- **CA-8:** Dado `--cores xyz`, quando o programa começa, então ele mostra `Cores inválidas: use sim, nao ou auto.` e termina com código 2.

## Fora de escopo

- Escolher a paleta.
- Cores nas instruções, nos pedidos e nos títulos.
- Detectar se o terminal suporta cores pelo `TERM` (por exemplo `TERM=dumb`).
- Cores no `stderr` (mensagens de opção inválida).

## Questões em aberto

Todas decididas em 2026-10-27, como propostas:

- [x] **1. Mensagens e cores:** as da tabela, com o resultado da máquina colorido só depois dos dois-pontos. Alternativas: colorir a linha inteira da máquina; ou deixar as mensagens de entrada inválida sem cor.
- [x] **2. Padrão `auto`:** o padrão passa a ser `auto`, desfazendo a decisão da 010 de cores sempre ligadas, como proposto no roadmap. A alternativa é manter `sim` como padrão e só acrescentar `auto` como opção.
- [x] **3. `NO_COLOR` com `sim`:** `--cores sim` ignora `NO_COLOR`, porque quem passa a opção pediu cores de propósito. A alternativa é `NO_COLOR` vencer sempre.
- [x] **4. Texto:** `Cores inválidas: use sim, nao ou auto.`.

## Encerramento

Concluída em 2026-10-27, nas tarefas T1–T4. Os 8 critérios de aceitação são cobertos por testes em `internal/terminal` e `cmd`: CA-1 em `TestPlaySolo_ColoredMessages` (e `TestShotColor`), CA-2 em `TestPlayVsMachine_ColoredMachineShots`, CA-3 em `TestPlayVsMachine_ColoredEnd`, CA-4 em `TestPlaySolo_ColoredMessages` (coordenada inválida, fora do tabuleiro, tiro repetido) e `TestPlayVsMachine_PlacementColoredErrors` (posição inválida, posição ocupada, sorteio que não cabe), CA-5 nos testes de partida colorida da 010 (`TestPlaySolo_Colored`, `TestPlayVsMachine_Colored`, `TestPlayVsMachine_PlacementColored`, que agora também conferem com `assertUncolored` que pedidos, instruções e títulos ficam sem cor) e em `TestRun_Colors`, CA-6 e CA-7 em `TestChooseStyle`, `TestIsTerminal` e `TestRun_Colors`, e CA-8 em `TestRun_InvalidColors`. A cobertura ficou em 99,6% em `internal/board`, 97,5% em `internal/terminal` e 92,9% em `cmd` (`chooseStyle`, `isTerminal`, `paint` e `shotColor` com 100%).

**Conferência à mão:** não havia um terminal de verdade no ambiente em que a spec foi implementada, então a conferência usou o binário com pipes, arquivos e o `script -q` do macOS, que roda o programa num pseudoterminal:

- com a saída num pipe ou num arquivo, sem a opção ou com `--cores auto`, não há nenhum código de cor; com `--cores sim` e `NO_COLOR=1`, há;
- no pseudoterminal, sem a opção, há cores; com `NO_COLOR=1`, não há; com `NO_COLOR` vazia, há; com `--cores nao`, não há;
- numa partida contra a máquina no pseudoterminal, `Posição inválida` e `Fora do tabuleiro` saem em amarelo, `Água.` em ciano, `Acerto!` e `Afundou o Submarino!` em vermelho em negrito, e na linha `A máquina atirou em H3: ` só o `Água.` tem cor;
- `--cores xyz` mostra `Cores inválidas: use sim, nao ou auto.` e termina com código 2.

Não foi conferido num terminal de verdade, com os olhos, como as cores novas ficam na tela (por exemplo, se o amarelo se lê bem em fundo claro). Fica para quem revisar o PR: `go run .`, `go run . | cat` e `NO_COLOR=1 go run .`.

**Desvios em relação ao plano:**

- No T1 e no T2, as mensagens passaram a sair coloridas antes de o padrão virar `auto` (T4). Por isso `TestRun_Win`, `TestRun_Fleet` e `TestRun_PlayerWins`, que conferem o fim da partida com o padrão, ganharam `--cores nao` já no T1 e no T2. Depois do T4 a opção continua lá, deixando explícito que esses testes querem a saída sem cores.
- O CA-4 do posicionamento ficou num teste novo, `TestPlayVsMachine_PlacementColoredErrors`, e não no `TestPlayVsMachine_PlacementColored`, porque o sorteio que não cabe precisa de outra frota (26 submarinos). O teste também confere a posição ocupada.
- O `TestPlayVsMachine_ColoredMachineShots` e a derrota do `TestPlayVsMachine_ColoredEnd` usam uma estratégia com as casas fixas (`machineScript`, da 013), para ter acerto, afundou e água numa só vez da máquina.
- Testes a mais: `TestShotColor`, e o `TestIsTerminal` também confere que o `/dev/null` conta como terminal, como o plano já previa nas decisões.
- O texto de ajuda da opção `--cores` mudou para explicar o `auto`.
- Como o plano previa, a vitória do modo solo (`Você afundou toda a frota em ...!`) também fica verde, embora a spec só cite `Você venceu!`.
- Depois do merge da 014, quem coordena trouxe o `master` para esta spec e pintou de amarelo as duas mensagens novas do posicionamento, `Não há navio para desfazer.` e o aviso de que o resto da frota não cabe, com o teste `TestPlayVsMachine_PlacementColoredUndoAndWarning`, como o plano previa.

**Para o roadmap:**

- Conferir num terminal de verdade como as cores ficam, inclusive em fundo claro.
- Os itens que a spec deixou fora de escopo: detectar `TERM=dumb`, cores no `stderr` e escolher a paleta.
