# 002 — Tabuleiro e frota

- **Status:** Rascunho
- **Roadmap:** #002
- **Depende de:** 001 (coordenadas)

## Problema

Antes de qualquer disparo, cada jogador precisa ter a frota posicionada no tabuleiro. O jogo precisa saber quais navios existem, que formato cada um tem e quais posições são válidas. É essa base que permite registrar disparos (#003), sortear um tabuleiro (#005) e deixar o jogador montar o próprio (#007).

## Regras

- A frota padrão segue a cartela de papel original (`doc/img/batalha-naval.jpg`). São 8 navios, que ocupam 19 casas:

  | Navio | Casas | Quantidade | Formato |
  | ----- | ----- | ---------- | ------- |
  | Submarino | 1 | 3 | 1 casa |
  | Destroyer | 2 | 2 | Linha de 2 casas |
  | Cruzador | 3 | 1 | Linha de 3 casas |
  | Encouraçado | 4 | 1 | Linha de 4 casas |
  | Porta-Aviões | 5 | 1 | T: 3 casas em linha e 2 descendo do meio |

- Os navios em linha podem ser posicionados na horizontal ou na vertical. O Submarino ocupa uma casa só, então a orientação não muda nada para ele.
- O Porta-Aviões pode ser posicionado em qualquer uma das 4 orientações do T:

  ```
  ■■■     ■       ■   ■
   ■      ■     ■■■   ■■■
   ■     ■■■      ■   ■
  ```

- Para posicionar um navio, o jogador indica o tipo do navio, uma casa de referência e uma orientação. A regra que define a casa de referência depende da questão em aberto 2.
- Todas as casas de um navio precisam estar dentro do tabuleiro.
- Dois navios não podem ocupar a mesma casa.
- Não é possível posicionar mais navios de um tipo do que a quantidade prevista na frota.
- Um posicionamento inválido é rejeitado com um erro, e o tabuleiro continua exatamente como estava.
- O tabuleiro está completo quando os 8 navios da frota estiverem posicionados.

## Critérios de aceitação

- **CA-1:** Dada a frota padrão, quando o jogo lista os navios, então são 3 Submarinos, 2 Destroyers, 1 Cruzador, 1 Encouraçado e 1 Porta-Aviões, somando 19 casas.
- **CA-2:** Dado um tabuleiro vazio, quando o jogador posiciona um Cruzador em `A1` na horizontal, então ele ocupa `A1`, `B1` e `C1`. Na vertical, ocupa `A1`, `A2` e `A3`.
- **CA-3:** Dado um tabuleiro vazio, quando o jogador posiciona um Submarino em `E5`, então ele ocupa só `E5`, em qualquer orientação.
- **CA-4:** Dado um tabuleiro vazio, quando o jogador posiciona um Encouraçado em `J7` na vertical, então o posicionamento é aceito e ocupa `J7` a `J10`. Em `J8` na vertical, ou em `H1` na horizontal, o posicionamento é rejeitado com um erro de fora do tabuleiro.
- **CA-5:** Dado um tabuleiro vazio, quando o jogador posiciona o Porta-Aviões em cada uma das 4 orientações, então ele ocupa as 5 casas esperadas. As casas exatas dependem da questão em aberto 2.
- **CA-6:** Dado um tabuleiro vazio, quando o jogador posiciona o Porta-Aviões numa orientação em que alguma das 5 casas sai do tabuleiro, então o posicionamento é rejeitado com um erro de fora do tabuleiro.
- **CA-7:** Dado um Cruzador em `A1` na horizontal, quando o jogador posiciona um Destroyer em `B1` na vertical, então o posicionamento é rejeitado com um erro de sobreposição, e o tabuleiro continua com só o Cruzador.
- **CA-8:** Dados 3 Submarinos posicionados, quando o jogador tenta posicionar um quarto, então o posicionamento é rejeitado com um erro de quantidade excedida.
- **CA-9:** Dado um tabuleiro com 7 navios posicionados, quando o jogo verifica se ele está completo, então a resposta é não. Depois de posicionar o oitavo, a resposta é sim.

As decisões das questões em aberto podem acrescentar ou ajustar critérios.

## Fora de escopo

- Disparos e o estado de casa atingida ou de navio afundado (#003).
- Ler o posicionamento do teclado (por exemplo, `A1 H`) e desenhar o tabuleiro no terminal (#004 e #007).
- Posicionamento aleatório (#005).
- Frotas diferentes da padrão (#008). As regras desta spec devem funcionar para qualquer frota, mas só a padrão é definida e testada.

## Questões em aberto

- [ ] **1. Navios encostados:** dois navios podem ocupar casas vizinhas? E na diagonal? Na versão de papel, o costume é que nenhum navio encoste em outro, nem pela diagonal.
- [ ] **2. Casa de referência e nomes das orientações:**
  - Para os navios em linha, a referência é a primeira casa (a mais à esquerda ou a de cima) e a orientação é `horizontal` ou `vertical`?
  - Para o Porta-Aviões, qual casa é a referência: o meio da barra de 3 casas ou o canto superior esquerdo do quadrado 3 × 3 que envolve o T?
  - As 4 orientações do T são nomeadas pela direção para onde a haste aponta (cima, baixo, esquerda, direita)?
- [ ] **3. Remover navio:** o jogador pode tirar ou reposicionar um navio já posicionado, ou isso fica para a #007, quando ele montar o tabuleiro pelo terminal?
- [ ] **4. Tamanho da spec:** se as decisões acima levarem o plano a mais de ~8 tarefas, dividimos a spec em 002 (frota e formatos) e uma spec nova para o posicionamento no tabuleiro?

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
