# 015 — Mais cores

- **Status:** Em andamento
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

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
