# 022 — Tela da partida

- **Status:** Em andamento
- **Roadmap:** #022 (onda 8), issue #32
- **Depende de:** 015 (mais cores)

## Problema

Hoje cada tabuleiro é impresso abaixo do anterior. Numa partida contra a máquina, cada turno mostra `Sua frota:` e `Frota inimiga:` um embaixo do outro, cada um com 13 linhas e a legenda, e depois de cada disparo vem mais um tabuleiro de ataque. O texto rola sem parar, e o jogador precisa subir a tela para lembrar o que aconteceu. Esta spec arruma a tela: os dois tabuleiros lado a lado, e a tela limpa a cada jogada, com o que acabou de acontecer logo abaixo.

## Regras

**Tabuleiros lado a lado:**

- Na partida contra a máquina, `Sua frota` fica à esquerda e `Frota inimiga` à direita, com os títulos em cima de cada um e as duas legendas embaixo, uma por linha:

  ```
  Sua frota                  Frota inimiga
     A B C D E F G H I J        A B C D E F G H I J
   1 ~ C C C ~ ~ ~ ~ ~ ~      1 S o . . . . . . . .
   2 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~      2 . . . . . . . . . .
   ...
  10 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~     10 . . . . . . . . . .

  Sua frota: ~ água, S D C E P navio, X atingido, o tiro na água
  Frota inimiga: . não atingida, o água, X acerto, S D C E P navio afundado
  ```

- A tela inteira cabe em 80 colunas e cerca de 20 linhas, contando as mensagens e o pedido. Para isso, a legenda da frota própria nesta tela cita as letras dos navios juntas (`S D C E P navio`), como a legenda inimiga, no lugar dos nomes por extenso (questão 6).
- Ficam fora da regra das 80 colunas as instruções do posicionamento (012), que aparecem uma vez, e a mensagem de coordenada inválida, que só aparece no erro.
- Os dois tabuleiros aparecem lado a lado a cada disparo do jogador, inclusive depois de um acerto (hoje aparece só o de ataque).
- No modo solo, há um tabuleiro só, como hoje.

**Mensagens do turno:**

- Como a tela é limpa, o que o jogador digitou some. Por isso o resultado passa a citar a casa: `Seu disparo em B7: Acerto!` (questão 3).
- Abaixo das legendas, a tela mostra o que aconteceu desde a tela anterior: o resultado do disparo do jogador e, quando foi água, as linhas dos disparos da máquina (`A máquina atirou em C4: Água.`). Depois vem o pedido `Disparo: `.
- Mensagens de entrada inválida (coordenada inválida, tiro repetido) aparecem logo abaixo do pedido, sem limpar a tela, como hoje.

**Limpar a tela (questões 1 e 2):**

- Antes de cada tela, o programa limpa o terminal, para o jogador ver só o estado atual.
- A limpeza segue a regra das cores da 015: com a saída num terminal, limpa; num arquivo ou outro programa, não limpa, e as telas ficam uma abaixo da outra.
- A opção `--limpar` escolhe à mão: `auto` (padrão), `sim` ou `nao`. Um valor inválido mostra `Limpar inválido: use sim, nao ou auto.` e termina com código 2.
- O posicionamento da frota (012) também limpa a tela antes de cada navio, mantendo as instruções no topo.

**Fim da partida:**

- Na vitória, na derrota e na desistência, a última tela mostra os dois tabuleiros lado a lado, com a frota inimiga revelada (como hoje, com os navios), seguidos da mensagem final e da semente.

## Critérios de aceitação

- **CA-1:** Dada uma partida contra a máquina, quando o turno começa, então os dois tabuleiros aparecem lado a lado, com os títulos e as duas legendas embaixo, e nenhuma linha da tela (títulos, tabuleiros, legendas, mensagens do turno e pedido) passa de 80 colunas.
- **CA-2:** Dado um acerto do jogador, quando ele atira de novo, então a tela mostra os dois tabuleiros de novo, com `Seu disparo em B7: Acerto!` abaixo deles.
- **CA-3:** Dada uma água do jogador, quando a máquina atira, então a próxima tela mostra o resultado do jogador e as linhas da máquina abaixo dos tabuleiros, e depois o pedido.
- **CA-4:** Dada uma coordenada inválida ou um tiro repetido, quando o jogador a digita, então a mensagem aparece abaixo do pedido, sem limpar a tela e sem redesenhar os tabuleiros.
- **CA-5:** Dado `--limpar sim` (ou `auto` num terminal), quando uma tela é desenhada, então ela começa com o código de limpar o terminal; dado `--limpar nao` (ou `auto` fora de um terminal), não há esse código.
- **CA-6:** Dado o fim da partida (vitória, derrota ou desistência), quando o programa termina, então a última tela mostra os dois tabuleiros lado a lado, com a frota inimiga revelada, e a mensagem final.
- **CA-7:** Dado o modo solo, quando a partida roda, então há um tabuleiro só, com a mesma regra de limpar a tela e o resultado citando a casa.
- **CA-8:** Dado o posicionamento, quando cada navio é pedido, então a tela é limpa (com a mesma regra) e mostra as instruções, o tabuleiro e o pedido.
- **CA-9:** Dado `--limpar xyz`, quando o programa começa, então ele mostra `Limpar inválido: use sim, nao ou auto.` e termina com código 2.
- **CA-10:** Dadas duas execuções com a mesma semente e as mesmas entradas, quando elas rodam, então as saídas são iguais, e as cores continuam funcionando nos dois tabuleiros.

## Fora de escopo

- Ajustar o desenho à largura do terminal, ou empilhar os tabuleiros em terminais estreitos.
- Uma interface de tela cheia que reage às teclas (setas, cursor).
- Histórico rolável das jogadas anteriores.
- Dois jogadores humanos (issue #34).

## Questões em aberto

Todas decididas em 2026-10-27, como propostas:

- [x] **1. Limpar a tela:** sim, a cada tela, só quando a saída é um terminal, como as cores da 015. A alternativa é nunca limpar e só pôr os tabuleiros lado a lado.
- [x] **2. Opção:** `--limpar auto|sim|nao`, com `auto` como padrão. A alternativa é não ter opção e seguir sempre o `auto`, ou ligar a limpeza junto com `--cores`.
- [x] **3. Resultado com a casa:** `Seu disparo em B7: Acerto!`, no lugar de só `Acerto!`, nos dois modos. As mensagens da máquina já citam a casa.
- [x] **4. Layout:** títulos `Sua frota` e `Frota inimiga` sem os dois-pontos, cinco espaços entre os tabuleiros, as legendas embaixo começando com o nome do tabuleiro. A alternativa é uma legenda só, juntando os símbolos dos dois.
- [x] **5. Posicionamento:** também limpa a tela a cada navio, mantendo as instruções no topo.
- [x] **6. Legenda curta (decidida em 2026-10-27, durante a implementação):** a legenda da frota própria tinha 114 colunas e passava do limite. Ela vira `Sua frota: ~ água, S D C E P navio, X atingido, o tiro na água`. As alternativas eram quebrar a legenda em duas linhas ou limitar a regra das 80 colunas aos tabuleiros.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
