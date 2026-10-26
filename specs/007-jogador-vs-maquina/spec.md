# 007 — Jogador contra a máquina

- **Status:** Aprovada
- **Roadmap:** #007 (onda 4)
- **Depende de:** 006 (modo solo)

## Problema

O segundo item do TODO do README é jogar contra a máquina, como no bloco de papel: cada lado tem a sua frota, os dois atiram alternadamente, e quem acerta atira de novo. Hoje o jogador só ataca um tabuleiro que não atira de volta (006). Esta spec cria a partida contra a máquina. A frota do jogador é sorteada; posicionar a própria frota à mão fica para a #012 (questão 1).

## Regras

**Início da partida:**

- A partida contra a máquina passa a ser o modo padrão. O modo solo da 006 continua disponível com `--modo solo`; `--modo maquina` é o padrão.
- O programa sorteia dois tabuleiros com a frota escolhida (`--frota`, 008): primeiro o do jogador, depois o da máquina. A mesma `--semente` repete os dois tabuleiros e todos os disparos da máquina.
- O jogador atira primeiro.

**Tela:**

- No turno do jogador, o programa mostra, um abaixo do outro:
  - `Sua frota:` e o tabuleiro próprio do jogador (004/009), com os navios e os disparos que a máquina já fez;
  - `Frota inimiga:` e o tabuleiro de ataque da máquina (009);
  - o pedido `Disparo: `.
- A tela inteira aparece no início da partida e cada vez que o turno volta para o jogador. Depois de um acerto do jogador, aparece só a mensagem do disparo, o tabuleiro de ataque atualizado e o pedido de novo.

**Turno do jogador:**

- Igual ao modo solo (006): uma coordenada por linha, as mesmas mensagens de resultado (009) e de entrada inválida, e `sair` ou Ctrl+D para desistir.
- Com `Acerto!` ou `Afundou o <navio>!`, o jogador atira de novo. Com `Água.`, o turno passa para a máquina.

**Turno da máquina:**

- A máquina atira numa casa sorteada entre as que ainda não atingiu no tabuleiro do jogador. Ela não segue nenhuma estratégia: depois de um acerto, o próximo tiro também é sorteado.
- Cada disparo da máquina aparece numa linha: `A máquina atirou em C4: Água.`, `A máquina atirou em C4: Acerto!` ou `A máquina atirou em C4: Afundou o seu Cruzador!`.
- Com acerto, a máquina atira de novo. Com água, o turno volta para o jogador.

**Fim da partida:**

- Quem afundar primeiro a frota inteira do outro vence.
- Vitória do jogador: a mensagem do último afundamento, o tabuleiro de ataque final e `Você venceu! Afundou a frota inimiga em N disparos.`
- Vitória da máquina: o último disparo da máquina, o tabuleiro próprio do jogador, `A máquina venceu! Ela afundou a sua frota em N disparos. A frota inimiga era:` e o tabuleiro próprio da máquina, revelando onde estavam os navios.
- Desistência: como no modo solo, `Você desistiu depois de N disparos. A frota inimiga era:` e o tabuleiro próprio da máquina.
- N conta só os disparos válidos de quem é citado, com o singular `1 disparo` como na 006.
- Em todos os casos, o programa mostra `Semente: N` no fim e termina com código 0.

**Opções de linha de comando:**

- `--modo` aceita `maquina` (padrão) ou `solo`, em maiúscula ou minúscula. Outro valor mostra `Modo inválido: use maquina ou solo.` e termina com código 2.
- `--frota` e `--semente` valem para os dois modos.

## Critérios de aceitação

Nos critérios abaixo, os testes leem as casas dos navios dos dois tabuleiros sorteados com a semente da partida, sem depender de uma semente específica.

- **CA-1:** Dado o programa iniciado sem `--modo`, quando a partida começa, então a saída mostra `Sua frota:` com o tabuleiro do jogador, `Frota inimiga:` com o tabuleiro de ataque todo em `.`, e `Disparo: `.
- **CA-2:** Dada uma partida, quando o jogador acerta um navio, então ele atira de novo sem que a máquina atire.
- **CA-3:** Dada uma partida, quando o jogador acerta a água, então a máquina atira até errar, cada disparo dela aparece na sua linha, e a tela inteira aparece de novo antes do próximo `Disparo: `.
- **CA-4:** Dado um tabuleiro de jogador e uma fonte de aleatoriedade, quando a máquina atira várias vezes, então cada disparo cai numa casa ainda não atingida, e, com a mesma semente, a sequência de disparos é a mesma.
- **CA-5:** Dada uma partida em que o jogador afunda todos os navios da máquina, quando o último afunda, então o programa mostra `Você venceu! Afundou a frota inimiga em N disparos.`, `Semente: N` e termina com código 0.
- **CA-6:** Dada uma partida em que o jogador só atira na água, quando a máquina afunda todos os navios do jogador, então o programa mostra `A máquina venceu! Ela afundou a sua frota em N disparos. A frota inimiga era:`, o tabuleiro da máquina e `Semente: N`.
- **CA-7:** Dada uma partida, quando o jogador digita `sair` ou a entrada termina, então o programa mostra `Você desistiu depois de N disparos. A frota inimiga era:` e o tabuleiro da máquina.
- **CA-8:** Dadas duas partidas com `--semente 42` e a mesma entrada, quando elas rodam, então as saídas são iguais.
- **CA-9:** Dado o programa iniciado com `--modo solo`, quando ele roda, então a partida é a do modo solo (006), sem frota própria nem disparos da máquina.
- **CA-10:** Dado o programa iniciado com `--modo xyz`, quando ele roda, então mostra `Modo inválido: use maquina ou solo.` e termina com código 2.

## Fora de escopo

- O jogador posicionar a própria frota (#012).
- Uma máquina com estratégia, por exemplo atirar em volta de um acerto.
- Os dois tabuleiros lado a lado, e limpar a tela entre os turnos.
- Dois jogadores humanos no mesmo terminal.
- Cores (#010).

## Questões em aberto

Todas decididas em 2026-10-26:

- [x] **1. Dividir a 007:** o roadmap junta "posicionar a própria frota" e "jogar contra a máquina". Juntas, as duas passam de ~8 tarefas. Proposta: esta spec faz a partida contra a máquina, com a frota do jogador sorteada, e uma spec nova, a 012, faz o posicionamento manual. Assim a 007 já entrega uma partida completa. As duas ficam em ondas seguidas, porque a 012 encaixa o posicionamento na partida da 007.
- [x] **2. Modo padrão:** a partida contra a máquina vira o padrão, e o modo solo passa para `--modo solo`.
- [x] **3. Quem começa:** sempre o jogador.
- [x] **4. Máquina sem estratégia:** cada tiro é sorteado entre as casas ainda não atingidas, como diz o README ("jogadas aleatórias").
- [x] **5. Tela:** os dois tabuleiros um abaixo do outro, com a tela inteira só quando o turno volta para o jogador. Os disparos da máquina aparecem como linhas de texto.
- [x] **6. Mensagens:** os textos de disparo da máquina, de vitória, de derrota, de desistência e de modo inválido.
- [x] **7. Derrota revela a frota inimiga:** sim, como na desistência.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
