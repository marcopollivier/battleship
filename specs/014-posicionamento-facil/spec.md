# 014 — Posicionamento mais fácil

- **Status:** Rascunho
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

- [ ] **1. Nome do comando:** `desfazer`, sem acento como os outros comandos da 012.
- [ ] **2. Desfazer sem navios:** mensagem `Não há navio para desfazer.`. A alternativa é ignorar o comando em silêncio e pedir o mesmo navio.
- [ ] **3. O aviso só avisa:** o jogador continua livre, como proposto. A alternativa é recusar a posição que deixa o resto sem lugar, com uma mensagem no lugar do aviso, como uma posição inválida.
- [ ] **4. Conferência exata:** o programa procura de verdade um jeito de encaixar o resto, e o aviso só aparece quando não há nenhum. Com a frota da cartela, essa busca é rápida. A alternativa é uma conferência aproximada (por exemplo, contar casas livres), que pode avisar sem necessidade ou deixar de avisar.
- [ ] **5. Textos:** a linha nova das instruções, a mensagem de desfazer sem navios e o aviso.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
