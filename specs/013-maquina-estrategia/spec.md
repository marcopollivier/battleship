# 013 — Máquina com estratégia

- **Status:** Em andamento
- **Roadmap:** #013 (onda 6)
- **Depende de:** 007 (jogador contra a máquina)

## Problema

Na partida da 007, cada tiro da máquina é sorteado, mesmo logo depois de um acerto. Uma pessoa jogando no papel atiraria em volta da casa atingida até afundar o navio. Por isso a máquina é fácil de vencer e não parece um adversário de verdade. Esta spec dá à máquina a estratégia de quem joga no papel: procurar ao acaso e, depois de um acerto, cercar o navio até afundá-lo.

## Regras

**O que a máquina sabe:**

- Só o que um jogador saberia: onde atirou, se foi `Água.`, `Acerto!` ou `Afundou`. Ela nunca olha onde estão os navios que ainda não achou.
- Um navio afundado ocupa as casas atingidas ligadas à casa do último disparo, porque os navios não podem se encostar (002). Assim a máquina sabe quais acertos já são de navios afundados e quais são de um navio ainda boiando.

**Procura:**

- Quando não há acerto em navio ainda boiando, a máquina sorteia uma casa, como na 007.
- A máquina não sorteia casas encostadas, nem pela diagonal, num navio já afundado, porque ali não pode haver navio (questão 2, decisão provisória na issue #25).

**Cerco:**

- Quando há um acerto em navio ainda boiando, a máquina atira numa casa vizinha, acima, abaixo, à esquerda ou à direita, que ainda não atingiu.
- Com dois ou mais acertos em linha nesse navio, ela prefere continuar a linha, nas duas pontas. Se as duas pontas já foram atingidas ou saem do tabuleiro, ela volta a atirar em qualquer vizinha dos acertos. Isso cobre o Porta-Aviões em T, que não é uma linha.
- Quando há mais de uma casa possível, a máquina sorteia entre elas.
- Quando o navio afunda, a máquina volta à procura.

**Nível (questão 1):**

- A opção `--nivel` escolhe a máquina: `dificil` (padrão) usa esta estratégia, e `facil` sorteia todos os tiros, como na 007.
- Um valor inválido mostra `Nível inválido: use facil ou dificil.` e termina com código 2, como as outras opções.
- No modo solo, a máquina não atira, e a opção não muda nada.

**Semente:**

- Com a mesma semente, o mesmo nível e as mesmas entradas, a partida se repete.
- Com o nível `dificil`, uma semente de uma partida anterior a esta spec não repete aquela partida. Com `--nivel facil`, repete.

**Mensagens:** as mesmas da 007. O jogador só percebe a estratégia pelos tiros.

## Critérios de aceitação

- **CA-1:** Dado um tabuleiro sem acertos em navio boiando, quando a máquina escolhe o tiro, então ela sorteia entre as casas ainda não atingidas, sem as casas encostadas em navios afundados.
- **CA-2:** Dado um único acerto em navio boiando, quando a máquina escolhe o tiro, então ele é uma das vizinhas acima, abaixo, à esquerda ou à direita ainda não atingidas.
- **CA-3:** Dados dois acertos em linha num navio boiando, quando a máquina escolhe o tiro, então ele é uma das pontas da linha ainda não atingidas.
- **CA-4:** Dados acertos em linha com as duas pontas já atingidas ou fora do tabuleiro (por exemplo, a barra do T do Porta-Aviões), quando a máquina escolhe o tiro, então ele é uma vizinha de algum dos acertos.
- **CA-5:** Dado um navio que acabou de afundar, quando a máquina escolhe o próximo tiro, então ela volta à procura, e as casas em volta do navio afundado não são sorteadas.
- **CA-6:** Dado qualquer tabuleiro completo da frota padrão, quando a máquina com o nível `dificil` atira até afundar a frota, então ela nunca atira duas vezes na mesma casa e termina; em média, com muitas sementes, ela precisa de menos tiros do que o nível `facil`.
- **CA-7:** Dadas duas execuções com `--semente 42`, o mesmo nível e as mesmas entradas, quando elas rodam, então as saídas são iguais.
- **CA-8:** Dado `--nivel facil`, quando a partida roda, então os tiros da máquina são os mesmos de antes desta spec para a mesma semente.
- **CA-9:** Dado `--nivel xyz`, quando o programa começa, então ele mostra `Nível inválido: use facil ou dificil.` e termina com código 2.

## Fora de escopo

- Estratégias mais fortes, como atirar em xadrez na procura ou calcular onde cada navio ainda cabe.
- Lembrar entre partidas como o jogador costuma posicionar a frota.
- Mostrar ao jogador o que a máquina está pensando.

## Questões em aberto

Todas decididas em 2026-10-26:

- [x] **1. Nível:** a estratégia vira o padrão (`dificil`), e `--nivel facil` mantém a máquina da 007. As alternativas eram substituir a máquina sorteada sem opção, ou manter a sorteada como padrão.
- [x] **2. Casas em volta de navio afundado:** a máquina pula essas casas. A decisão é provisória: o jogador quer testar jogando, e a alternativa (continuar sorteando essas casas) está registrada na issue #25.
- [x] **3. Continuar a linha:** com dois acertos em linha, a máquina prefere as pontas da linha.
- [x] **4. Nome e textos:** `--nivel` com `facil` e `dificil`, e a mensagem `Nível inválido: use facil ou dificil.`.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
