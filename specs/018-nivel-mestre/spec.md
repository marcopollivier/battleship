# 018 — Terceiro nível da máquina

- **Status:** Rascunho
- **Roadmap:** #018 (onda 7)
- **Depende de:** 013 (máquina com estratégia)

## Problema

A máquina `dificil` da 013 cerca um navio depois de acertá-lo, mas procura os navios sorteando casas. Um bom jogador de papel procura de um jeito mais esperto: atira onde ainda cabe mais navio, e evita as casas onde nenhum navio que falta caberia. O encerramento da 013 sugeriu esse jogador como um terceiro nível.

## Regras

**O nível (questão 1):**

- `--nivel` ganha um terceiro valor, `mestre`. `dificil` continua o padrão, e `facil` e `dificil` não mudam.
- A mensagem de valor inválido passa a ser `Nível inválido: use facil, dificil ou mestre.`.

**Como a máquina `mestre` escolhe o tiro (questão 2):**

- Ela sabe o mesmo que a `dificil`: onde atirou, o resultado de cada tiro e quais navios afundaram. Ela também sabe a frota da partida (a da cartela ou a do `--frota`), como um jogador que conhece as regras, e por isso sabe quais navios ainda faltam afundar.
- Para cada casa ainda não atingida, ela conta de quantos jeitos os navios que faltam poderiam passar por ela, levando em conta tudo o que sabe:
  - um navio não passa por casas de água, nem por casas de navios afundados, nem encosta num navio afundado;
  - o Porta-Aviões conta nas suas 4 rotações, e os navios em linha na horizontal e na vertical.
- **Procura:** sem acertos em navio boiando, ela atira na casa com a maior contagem.
- **Cerco:** com acertos em navio boiando, ela só conta os jeitos que passam por esses acertos, e atira na casa com a maior contagem. Assim ela segue o navio como a `dificil`, mas já sabe para que lado ele provavelmente continua.
- Quando várias casas empatam na maior contagem, ela sorteia entre elas, com a semente da partida.

**Semente:** com a mesma semente, o mesmo nível e as mesmas entradas, a partida se repete.

**Mensagens:** as mesmas da 007.

## Critérios de aceitação

- **CA-1:** Dado um tabuleiro sem tiros, quando a máquina `mestre` escolhe o tiro, então ele está entre as casas de maior contagem (as do meio do tabuleiro, onde passam mais navios).
- **CA-2:** Dado um canto do tabuleiro cercado de água, onde nenhum navio que falta caberia, quando a máquina escolhe o tiro, então ela nunca atira ali.
- **CA-3:** Dado um acerto em navio boiando, quando a máquina escolhe o tiro, então ela atira numa vizinha do acerto, e, com o navio encostado na borda ou na água de um lado, para o lado em que ele cabe.
- **CA-4:** Dados navios afundados, quando a máquina escolhe o tiro, então ela não conta mais os navios afundados entre os que faltam, nem as casas em volta deles.
- **CA-5:** Dados dois tabuleiros com os mesmos tiros e resultados, mas com navios ainda não achados em lugares diferentes, quando a máquina escolhe o tiro com a mesma semente, então o tiro é o mesmo.
- **CA-6:** Dados 200 tabuleiros sorteados da frota padrão, quando as máquinas `dificil` e `mestre` atiram até afundar a frota, então a `mestre` nunca repete uma casa, sempre termina, e precisa em média de menos tiros que a `dificil`.
- **CA-7:** Dado `--nivel mestre` (ou `MESTRE`), quando a partida roda, então a máquina joga com esta estratégia; e com a mesma semente e as mesmas entradas, as saídas são iguais.
- **CA-8:** Dado `--nivel xyz`, quando o programa começa, então ele mostra `Nível inválido: use facil, dificil ou mestre.` e termina com código 2.

## Fora de escopo

- Mudar o padrão para `mestre`.
- Aprender como o jogador costuma posicionar a frota.
- Uma máquina que posiciona a própria frota de forma esperta (hoje o tabuleiro da máquina é sorteado).

## Questões em aberto

- [ ] **1. Nome e padrão:** `mestre`, com `dificil` continuando o padrão. Alternativas: chamar de `especialista`; ou tornar o novo nível o padrão.
- [ ] **2. Estratégia:** contar de quantos jeitos os navios que faltam passam por cada casa, na procura e no cerco, como proposto. A alternativa, mais simples e mais fraca, é a `dificil` atirando em xadrez na procura (uma casa sim, outra não), porque todo navio ocupa pelo menos duas casas vizinhas, menos o Submarino.
- [ ] **3. Regra da issue #25:** na `mestre`, as casas em volta de navio afundado ficam fora naturalmente, porque nenhum navio pode encostar nele. Se a issue #25 decidir tirar a regra da `dificil`, a `mestre` continua com ela, já que é a regra do jogo usada a favor da máquina.
- [ ] **4. Texto:** `Nível inválido: use facil, dificil ou mestre.`.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
