# 003 — Disparo

- **Status:** Em andamento
- **Roadmap:** #003 (onda 1)
- **Depende de:** 001 (coordenadas), 002 (tabuleiro e frota)

## Problema

O jogo acontece em disparos: o jogador escolhe uma casa do tabuleiro do oponente e descobre se acertou a água ou um navio. Para isso, o tabuleiro precisa receber disparos, dizer o resultado de cada um, lembrar onde já houve disparo e saber quando toda a frota afundou. Sem isso, não existe partida: nem o modo solo (#006), nem o jogo contra a máquina (#007), nem o desenho do tabuleiro de ataque (#009).

## Regras

Do README:

- O jogador dispara indicando uma casa do tabuleiro do oponente (ex.: `A1`, `B3`, `J10`).
- Se não houver navio na casa, o disparo é um **tiro na água**.
- Se houver, é um **acerto**.
- Quando todas as casas de um navio são acertadas, o navio está **afundado**.
- Quem afundar todos os navios do oponente ganha o jogo.

Desta spec:

- O resultado de um disparo é um entre três: água, acerto ou afundado. O disparo que acerta a última casa ainda intacta de um navio resulta em afundado, e não em acerto.
- Quando o resultado é afundado, ele informa o tipo do navio afundado ("afundou o Cruzador").
- Quando o resultado é acerto, ele não informa o tipo do navio atingido.
- Só um tabuleiro completo (com todos os navios da frota posicionados) recebe disparos.
- Um disparo é rejeitado com um erro, sem alterar o tabuleiro, quando:
  - a casa está fora do tabuleiro;
  - a casa já recebeu disparo antes, seja água ou acerto. O jogador atira de novo, sem perder a vez (a vez em si é tratada na #006 e na #007);
  - o tabuleiro ainda não está completo.
- O tabuleiro lembra, para cada casa, se já recebeu disparo e, se recebeu, se foi água ou acerto. É isso que vai permitir desenhar o tabuleiro de ataque (#009).
- O tabuleiro informa se toda a frota já foi afundada.

## Critérios de aceitação

Nos critérios abaixo, o "tabuleiro de exemplo" é um tabuleiro completo com a frota padrão, na mesma disposição usada nos testes da 002:

| Navio | Posição | Casas |
| ----- | ------- | ----- |
| Porta-Aviões | `B1`, `baixo` | `A1` `B1` `C1` `B2` `B3` |
| Encouraçado | `E1`, horizontal | `E1` `F1` `G1` `H1` |
| Cruzador | `J1`, vertical | `J1` `J2` `J3` |
| Destroyer | `A5`, horizontal | `A5` `B5` |
| Destroyer | `E3`, horizontal | `E3` `F3` |
| Submarino | `H5` | `H5` |
| Submarino | `A8` | `A8` |
| Submarino | `E8` | `E8` |

- **CA-1:** Dado o tabuleiro de exemplo, quando o jogador dispara em `J10`, então o resultado é água.
- **CA-2:** Dado o tabuleiro de exemplo, quando o jogador dispara em `E1`, então o resultado é acerto, sem informar o tipo do navio.
- **CA-3:** Dado o tabuleiro de exemplo com disparos em `J1` e `J2`, quando o jogador dispara em `J3`, então o resultado é afundado, informando que foi o Cruzador.
- **CA-4:** Dado o tabuleiro de exemplo, quando o jogador dispara em `H5`, então o resultado é afundado logo no primeiro disparo, informando que foi um Submarino.
- **CA-5:** Dado o tabuleiro de exemplo com disparos em `J10` (água) e `E1` (acerto), quando o jogo consulta as casas, então `J10` aparece como água atingida, `E1` como navio atingido, e as demais como ainda não atingidas.
- **CA-6:** Dado o tabuleiro de exemplo, quando o jogador dispara numa casa fora do tabuleiro, então o disparo é rejeitado com um erro de fora do tabuleiro, e nenhuma casa muda.
- **CA-7:** Dado o tabuleiro de exemplo com todas as casas de navio atingidas, exceto `E8`, quando o jogo verifica se a frota inteira afundou, então a resposta é não. Depois do disparo em `E8`, a resposta é sim.
- **CA-8:** Dado o tabuleiro de exemplo com disparos em `J10` e `E1`, quando o jogador dispara de novo em `J10` ou em `E1`, então o disparo é rejeitado com um erro de tiro repetido, e nenhuma casa muda.
- **CA-9:** Dado um tabuleiro com só um Cruzador posicionado, quando o jogador dispara em qualquer casa, então o disparo é rejeitado com um erro de tabuleiro incompleto.

## Fora de escopo

- De quem é a vez, turno extra após acerto e alternância entre jogadores. Isso fica para a #006 e a #007.
- Ler o disparo do teclado e desenhar o tabuleiro de ataque (#009).
- Escolher os disparos da máquina (#007).

## Questões em aberto

Todas decididas em 2026-10-24:

- [x] **1. Tiro repetido:** rejeitado com erro; o jogador atira de novo sem perder a vez (CA-8).
- [x] **2. Qual navio afundou:** o resultado afundado informa o tipo do navio (CA-3, CA-4).
- [x] **3. Acerto revela o tipo:** não, o acerto só diz "acerto" (CA-2).
- [x] **4. Tabuleiro incompleto:** não recebe disparos (CA-9).

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
