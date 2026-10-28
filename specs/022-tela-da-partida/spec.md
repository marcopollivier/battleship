# 022 — Tela da partida

- **Status:** Concluída
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

  Sua frota: ~ água, S D C E P navio, X atingido, # afundado, o tiro na água
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

Concluída em 2026-10-27, nas tarefas T1–T7. Os 10 critérios de aceitação são cobertos por testes em `internal/terminal` e `cmd/init_test.go`:

- CA-1 em `TestGrids`, `TestSideBySide`, `TestBoardTitles` e `TestPlayVsMachine_Start`;
- CA-2 em `TestPlayVsMachine_HitShootsAgain` e `TestPlayVsMachine_PlayerWins`;
- CA-3 em `TestPlayVsMachine_MachineTurn`;
- CA-4 em `TestPlayVsMachine_InvalidInput`;
- CA-5 em `TestClearScreen`, `TestPlayVsMachine_Clear`, `TestPlaySolo_Clear`, `TestChooseClear` e `TestRun_Clear`;
- CA-6 em `TestPlayVsMachine_PlayerWins`, `_MachineWins`, `_GiveUp` e `_Colored`;
- CA-7 em `TestPlaySolo_Shots`, `_Win`, `_GiveUp`, `_Clear` e `_ColoredMessages`;
- CA-8 em `TestPlayVsMachine_PlacementClear` e nos outros testes do posicionamento;
- CA-9 em `TestRun_InvalidClear`;
- CA-10 em `TestRun_SameSeed` (com `--limpar sim --cores sim`), `TestRun_Clear` e `TestPlayVsMachine_Colored`.

A cobertura ficou em 98,0% em `internal/terminal`, 99,7% em `internal/board` e 93,9% em `cmd`.

O binário foi rodado com `--semente 42`, `aleatorio` e alguns disparos. Com `--cores nao --limpar nao`, as telas ficam uma embaixo da outra; com `--limpar sim`, cada tela (cada navio do posicionamento, cada disparo e a tela final) começa com `\x1b[H\x1b[2J`. Nenhuma linha da tela da partida passa de 49 colunas nos tabuleiros e 73 nas legendas. Uma tela depois de uma água:

```
Sua frota                  Frota inimiga
   A B C D E F G H I J        A B C D E F G H I J
 1 ~ C C C ~ ~ ~ ~ ~ ~      1 S o . . . . . . . .
 2 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~      2 . . . . . . . . . .
 3 D ~ ~ E E E E o ~ ~      3 . . . . . . . . . .
...
10 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~     10 . . . . . . . . . .

Sua frota: ~ água, S D C E P navio, X atingido, # afundado, o tiro na água
Frota inimiga: . não atingida, o água, X acerto, S D C E P navio afundado

Seu disparo em B1: Água.
A máquina atirou em H3: Água.
Disparo:
```

**Desvios em relação ao plano:**

- A legenda da frota própria na tela tinha 114 colunas e contrariava a regra das 80. A implementação parou depois do T3, e a spec mudou durante a implementação (questão 6): a legenda da tela é `Sua frota: ~ água, S D C E P navio, X atingido, o tiro na água`, numa constante própria (`ownScreenLegend`); a `Legend` do `RenderBoard` não mudou. As instruções do posicionamento e a mensagem de coordenada inválida ficaram fora da regra.
- A tela final tem uma legenda só, `Legenda: ~ água, S D C E P navio, X atingido, o tiro na água`, porque os dois tabuleiros mostram os navios. O plano não listava legenda na tela final.
- Com a 023 no `master`, as legendas curtas da tela ganharam `# afundado` (`Sua frota: ~ água, S D C E P navio, X atingido, # afundado, o tiro na água`, 74 colunas), no merge do `master` neste branch.
- No modo solo, a tela final usa a grade do `RenderBoard` com essa legenda curta, e não o `RenderBoard` inteiro, para caber em 80 colunas. As frases `A frota era:` e `A frota inimiga era:` saíram, porque a frota já aparece na tela.
- `Frota posicionada!` virou a primeira mensagem da primeira tela da partida; antes da tela, ela sumiria com a limpeza.
- As instruções do posicionamento aparecem em cada tela de navio, como no plano, também sem limpar a tela.
- O `shoot` do `playerAttack` não escreve mais o resultado: ele devolve a casa e o disparo, e o resultado entra nas `events` do `playerAttack`, que o solo e a partida contra a máquina escrevem abaixo dos tabuleiros.
- No `cmd`, as constantes `colorsOn`, `colorsOff` e `colorsAuto` viraram `optionOn`, `optionOff` e `optionAuto`, usadas por `--cores` e `--limpar`.
- O commit do T6 deixou o `TestRun_GiveUp` do `cmd` esperando a frase antiga; um commit a mais corrigiu o teste.
- Sem terminal (saída num arquivo), o pedido `Disparo: ` fica na mesma linha do começo da tela seguinte, porque a entrada não aparece; foi sempre assim.

**Para o roadmap:**

- A legenda do tabuleiro no posicionamento continua com 110 colunas (o `RenderBoard` com a `Legend`), e as instruções do posicionamento passam de 80 colunas.
- A mensagem de coordenada inválida tem 82 colunas.
- Na fusão com a 023, o `# afundado` entra também nas legendas curtas `ownScreenLegend` e `finalLegend`.
