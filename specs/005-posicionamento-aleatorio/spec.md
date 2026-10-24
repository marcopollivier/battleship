# 005 — Posicionamento aleatório

- **Status:** Em andamento
- **Roadmap:** #005 (onda 1)
- **Depende de:** 002 (tabuleiro e frota)

## Problema

O primeiro item do TODO do README é deixar o jogador atacar um tabuleiro montado aleatoriamente. A máquina também vai precisar de uma frota sorteada no jogo contra o jogador (#007). Esta spec gera, sem intervenção humana, um tabuleiro completo que respeite todas as regras de posicionamento da 002.

## Regras

- O jogo gera um tabuleiro completo a partir de uma frota: todos os navios posicionados, todas as regras da 002 respeitadas (dentro do tabuleiro, sem sobreposição, sem encostar, quantidade por tipo).
- O sorteio usa uma fonte de aleatoriedade recebida de fora, como pede a constituição. Com a mesma fonte, inicializada da mesma forma, o tabuleiro gerado é sempre o mesmo.
- Qualquer posição e orientação válida pode aparecer, inclusive as 4 orientações do Porta-Aviões. O sorteio não precisa ser uniforme entre todas as disposições possíveis.
- A geração sempre termina. Se a frota não couber no tabuleiro, o jogo retorna um erro, em vez de ficar tentando para sempre.

## Critérios de aceitação

- **CA-1:** Dada a frota padrão e uma fonte de aleatoriedade, quando o jogo gera um tabuleiro, então ele está completo, com os 8 navios.
- **CA-2:** Dados 1.000 tabuleiros gerados com fontes diferentes, quando o jogo verifica cada um, então nenhum tem navio fora do tabuleiro, sobreposto ou encostado em outro.
- **CA-3:** Dadas duas fontes de aleatoriedade inicializadas com a mesma semente, quando o jogo gera um tabuleiro com cada uma, então os dois tabuleiros são idênticos.
- **CA-4:** Dadas duas sementes diferentes, quando o jogo gera um tabuleiro com cada uma, então os tabuleiros são diferentes.
- **CA-5:** Dados 1.000 tabuleiros gerados, quando o jogo conta as orientações sorteadas, então todas as orientações válidas de cada tipo aparecem pelo menos uma vez: horizontal e vertical para os navios em linha, e as 4 orientações para o Porta-Aviões.
- **CA-6:** Dada uma frota que não cabe no tabuleiro (por exemplo, 30 Encouraçados), quando o jogo tenta gerar um tabuleiro, então ele retorna um erro, em tempo limitado.

## Fora de escopo

- Mostrar o tabuleiro gerado e deixar o jogador atacá-lo (#006).
- Escolher os disparos da máquina (#007).
- Escolher a semente ao executar o programa. Por decisão desta spec, a #006 vai usar o relógio por padrão e permitir que o jogador informe uma semente para repetir uma partida.

## Questões em aberto

Todas decididas em 2026-10-24:

- [x] **1. Distribuição:** não precisa ser uniforme; basta que qualquer posição e orientação válida possa aparecer (CA-5).
- [x] **2. Semente no programa:** relógio por padrão, com opção de o jogador informar a semente. A implementação fica na #006 (fora de escopo).

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
