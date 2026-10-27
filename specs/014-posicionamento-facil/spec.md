# 014 — Posicionamento mais fácil

- **Status:** Concluída
- **Roadmap:** #014 (onda 7)
- **Depende de:** 012 (posicionar a própria frota)

## Problema

No posicionamento da 012, um navio no lugar errado só sai com `recomecar`, que tira a frota inteira. E o jogador pode posicionar navios de um jeito que o resto da frota já não cabe, sem saber; ele só descobre quando nenhuma posição funciona ou quando o `aleatorio` falha. O encerramento da 012 deixou as duas melhorias para depois.

## Regras

**Desfazer:**

- Um comando novo, `desfazer`, aceito no lugar de uma posição, em maiúscula ou minúscula (questão 1).
- Ele tira do tabuleiro o último navio posicionado e pede esse navio de novo, com o tabuleiro atualizado. Pode ser usado várias vezes seguidas, tirando um navio de cada vez, na ordem inversa em que foram posicionados.
- Sem nenhum navio no tabuleiro, o programa mostra `Não há navio para desfazer.` e pede o mesmo navio de novo (questão 2).
- As instruções do início passam a citar o comando:

  ```
  Comandos: aleatorio sorteia o resto da frota, desfazer tira o último navio, recomecar tira todos os navios, sair desiste.
  ```

**Aviso de que o resto não cabe:**

- Depois de cada navio posicionado, o programa confere se os navios que faltam ainda cabem no tabuleiro, respeitando as regras de sempre (dentro do tabuleiro, sem sobrepor e sem encostar).
- Se não couberem de nenhum jeito, o programa mostra, antes do tabuleiro e do pedido do próximo navio:

  ```
  Atenção: o resto da frota não cabe mais no tabuleiro. Use desfazer ou recomecar.
  ```

- O aviso não impede nada: o jogador ainda pode tentar posições, usar `aleatorio`, `desfazer`, `recomecar` ou `sair` (questão 3).
- A conferência é exata: o aviso só aparece quando realmente não existe nenhum jeito de encaixar o resto (questão 4).
- A conferência não usa a semente. Uma partida com a mesma semente e as mesmas entradas continua igual.

## Critérios de aceitação

- **CA-1:** Dados dois navios posicionados, quando o jogador digita `desfazer`, então o segundo navio sai do tabuleiro, o tabuleiro aparece sem ele, e o pedido é o desse navio.
- **CA-2:** Dados dois navios posicionados, quando o jogador digita `desfazer` duas vezes, então os dois saem e o pedido volta ao Porta-Aviões.
- **CA-3:** Dado o tabuleiro vazio, quando o jogador digita `desfazer`, então aparece `Não há navio para desfazer.` e o mesmo pedido.
- **CA-4:** Dada uma posição que deixa o resto da frota sem lugar, quando o jogador a digita, então o navio é posicionado e aparece o aviso antes do pedido do próximo navio.
- **CA-5:** Dadas posições que ainda deixam lugar para o resto da frota, mesmo que apertado, quando o jogador as digita, então o aviso não aparece.
- **CA-6:** Dado o aviso, quando o jogador digita `desfazer`, então o navio sai e, se o resto voltar a caber, o aviso não aparece mais.
- **CA-7:** Dadas duas execuções com `--semente 42` e as mesmas entradas, com `desfazer` e com o aviso, quando elas rodam, então as saídas são iguais.
- **CA-8:** Dado o início do posicionamento, quando o programa mostra as instruções, então elas citam `desfazer`.

## Fora de escopo

- Refazer um navio desfeito.
- Desfazer depois que a partida começou, ou desfazer um `aleatorio` (ele termina o posicionamento).
- Sugerir onde pôr o próximo navio.
- Bloquear posições que deixariam o resto sem lugar.

## Questões em aberto

Todas decididas em 2026-10-27, como propostas:

- [x] **1. Nome do comando:** `desfazer`, sem acento como os outros comandos da 012.
- [x] **2. Desfazer sem navios:** mensagem `Não há navio para desfazer.`. A alternativa é ignorar o comando em silêncio e pedir o mesmo navio.
- [x] **3. O aviso só avisa:** o jogador continua livre, como proposto. A alternativa é recusar a posição que deixa o resto sem lugar, com uma mensagem no lugar do aviso, como uma posição inválida.
- [x] **4. Conferência exata:** o programa procura de verdade um jeito de encaixar o resto, e o aviso só aparece quando não há nenhum. Com a frota da cartela, essa busca é rápida. A alternativa é uma conferência aproximada (por exemplo, contar casas livres), que pode avisar sem necessidade ou deixar de avisar.
- [x] **5. Textos:** a linha nova das instruções, a mensagem de desfazer sem navios e o aviso.

## Encerramento

Concluída em 2026-10-27, nas tarefas T1–T5. Os 8 critérios de aceitação são cobertos por testes em `internal/board/board_test.go` (CA-1 e CA-2 no domínio em `TestBoard_RemoveLast`), `internal/board/fit_test.go` (CA-4 e CA-5 no domínio em `TestFits`, o limite de tentativas em `TestFits_AttemptLimit` e o pior caso medido em `TestFits_Bounded`), `internal/terminal/placement_test.go` (CA-1, CA-2 e CA-3 em `TestPlayVsMachine_PlacementUndo`, CA-4, CA-5 e CA-6 em `TestPlayVsMachine_PlacementWarning`, CA-8 em `TestPlayVsMachine_PlacementStart`, e o CA-5 também em `TestPlayVsMachine_PlaceFleet`, que continua sem aviso) e `cmd/init_test.go` (CA-7 em `TestRun_SameSeed`). A cobertura ficou em 99,7% em `internal/board`, 97,5% em `internal/terminal` e 91,9% em `cmd`.

**Pior caso do `Fits` com a frota da cartela** (medido no T2, com testes temporários, sem commit):

- Tabuleiros como o jogador monta, com os navios do maior para o menor: 1.400.000 tabuleiros sorteados (os 7 primeiros navios de 200.000 sementes) e uma busca por subida de encosta (300 partidas de 3.000 trocas de posição) para achar o pior. O pior caso usou 283 tentativas de `Place` (Porta-Aviões em B2 baixo e Encouraçado em F2 horizontal), em cerca de 0,2 ms. Em nenhum deles o resto deixou de caber.
- Tabuleiros com só alguns dos navios menores, que o posicionamento não monta, mas o domínio aceita: o pior caso usou 2.699 tentativas, em cerca de 2 ms.
- O limite de 1.000.000 de tentativas fica bem longe. Ao ser atingido (26 Submarinos com 1 posicionado, um caso que o `--frota` não aceita), a busca leva cerca de 1 s e o `Fits` responde `true`. `TestFits_Bounded` guarda os dois piores tabuleiros, com teto de 10.000 tentativas.

Com a frota da cartela, o aviso, na prática, não aparece: como os navios são pedidos do maior para o menor, os que faltam são sempre os menores, e eles sempre couberam. Uma conta confirma para os dois últimos: 7 navios bloqueiam no máximo 96 casas (as deles e as em volta), e 6 navios deixam pelo menos 13 casas livres, onde sempre há duas que não se encostam. O aviso serve às frotas do `--frota`.

A partida foi conferida à mão com `--semente 42 --cores nao`: `desfazer` no tabuleiro vazio mostrou `Não há navio para desfazer.`; dois navios e dois `desfazer` voltaram ao Porta-Aviões; a frota do exemplo da 003, com um `desfazer` num Submarino, começou a partida. Com `--frota E=7,C=1`, seis Encouraçados nas fileiras 2, 5 e 8 e o sétimo em D10 mostraram o aviso antes do tabuleiro; o Cruzador em H10 foi recusado (`Muito perto`), e depois do `desfazer` o aviso sumiu. Com `aleatorio` depois do `desfazer`, a frota foi sorteada e a partida começou.

**Desvios em relação ao plano:**

- Os casos dos testes em que o resto não cabe usam uma frota de 7 Encouraçados e 1 Cruzador (seis Encouraçados em três fileiras não deixam lugar fora da fileira 10), e não o exemplo de 2 Encouraçados do plano, que não bloqueia o tabuleiro. O caso em xadrez usa 26 Submarinos, com 20 e 25 posicionados.
- A busca fica numa função interna, `fitSearch`, que devolve também o número de tentativas, para os testes do limite e do pior caso. O limite ganhou um teste próprio, `TestFits_AttemptLimit`.
- O teste do T5 passou de cara, porque o `desfazer` e o aviso vieram no T3 e no T4. Ele ganhou um caso com `--frota E=7,C=1` para repetir uma partida com o aviso, já que a frota da cartela não o mostra.
- O aviso aparece antes de cada tabuleiro, como no plano; depois de uma posição recusada, o tabuleiro não é desenhado de novo, e o aviso também não.

**Para o roadmap:**

- Decidir se o aviso deve continuar, já que só aparece com frotas do `--frota`.
- Na busca, o Submarino tenta as duas orientações, que dão a mesma casa; tirar a repetição corta as tentativas pela metade nas frotas com muitos Submarinos.
- Com frotas grandes do `--frota`, a busca pode chegar ao limite e levar cerca de 1 s antes de cada tabuleiro.
