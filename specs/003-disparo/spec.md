# 003 — Disparo

- **Status:** Rascunho
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
- Um disparo fora do tabuleiro é rejeitado com um erro, e o tabuleiro não muda.
- O tabuleiro lembra, para cada casa, se já recebeu disparo e, se recebeu, se foi água ou acerto. É isso que vai permitir desenhar o tabuleiro de ataque (#009).
- O tabuleiro informa se toda a frota já foi afundada.

## Critérios de aceitação

Nos critérios abaixo, "tabuleiro de exemplo" é um tabuleiro com um Cruzador em `A1` na horizontal (`A1`–`C1`) e um Submarino em `E5`.

- **CA-1:** Dado o tabuleiro de exemplo, quando o jogador dispara em `J10`, então o resultado é água.
- **CA-2:** Dado o tabuleiro de exemplo, quando o jogador dispara em `B1`, então o resultado é acerto.
- **CA-3:** Dado o tabuleiro de exemplo com disparos em `A1` e `B1`, quando o jogador dispara em `C1`, então o resultado é afundado.
- **CA-4:** Dado o tabuleiro de exemplo, quando o jogador dispara em `E5`, então o resultado é afundado logo no primeiro disparo, porque o Submarino ocupa uma casa só.
- **CA-5:** Dado o tabuleiro de exemplo com disparos em `J10` (água) e `B1` (acerto), quando o jogo consulta as casas, então `J10` aparece como água atingida, `B1` como navio atingido, e as demais como ainda não atingidas.
- **CA-6:** Dado o tabuleiro de exemplo, quando o jogador dispara numa casa fora do tabuleiro, então o disparo é rejeitado com um erro de fora do tabuleiro, e nenhuma casa muda.
- **CA-7:** Dado o tabuleiro de exemplo com o Cruzador afundado, quando o jogo verifica se a frota inteira afundou, então a resposta é não. Depois do disparo em `E5`, a resposta é sim.

As decisões das questões em aberto podem acrescentar ou ajustar critérios.

## Fora de escopo

- De quem é a vez, turno extra após acerto e alternância entre jogadores. Isso fica para a #006 e a #007.
- Ler o disparo do teclado e desenhar o tabuleiro de ataque (#009).
- Escolher os disparos da máquina (#007).

## Questões em aberto

- [ ] **1. Tiro repetido:** um disparo numa casa que já recebeu disparo deve ser rejeitado com um erro (o jogador atira de novo, sem perder a vez) ou aceito e contado como tiro perdido?
- [ ] **2. Qual navio afundou:** quando um disparo afunda um navio, o resultado informa o tipo do navio ("afundou o Cruzador") ou só que um navio afundou? Esta questão vem do roadmap.
- [ ] **3. Acerto revela o tipo:** num acerto que não afunda, o resultado informa o tipo do navio atingido ou só "acerto"?
- [ ] **4. Tabuleiro incompleto:** um tabuleiro que ainda não tem os 8 navios posicionados pode receber disparos, ou o disparo deve ser rejeitado até ele ficar completo?

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
