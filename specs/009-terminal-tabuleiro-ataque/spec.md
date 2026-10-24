# 009 — Terminal: tabuleiro de ataque

- **Status:** Rascunho
- **Roadmap:** #009 (onda 2)
- **Depende de:** 003 (disparo), 004 (terminal, tabuleiro próprio)

## Problema

Para atacar, o jogador precisa ver o tabuleiro do oponente do jeito que ele conhece: onde já atirou, onde acertou água e onde acertou navio, sem ver os navios que ainda não encontrou. Também precisa ler o resultado de cada disparo em português. O domínio já registra os disparos (003) e o terminal já desenha o tabuleiro próprio (004); falta desenhar o tabuleiro de ataque e as mensagens. Isso é pré-requisito do modo solo (#006).

## Regras

- O tabuleiro de ataque tem o mesmo cabeçalho e as mesmas fileiras do tabuleiro próprio (004).
- Ele nunca mostra um navio que ainda não foi atingido.
- Cada casa aparece de acordo com o que o atacante sabe dela, com os símbolos definidos pela questão em aberto 1.
- Abaixo do tabuleiro, depois de uma linha em branco, aparece uma legenda própria do tabuleiro de ataque.
- O resultado de um disparo é exibido em português, com o texto definido pela questão em aberto 3. O tipo do navio só aparece quando ele afunda (decisão da 003).
- Os símbolos do tabuleiro são só ASCII, como na 004.

## Critérios de aceitação

Nos critérios abaixo, o "tabuleiro de exemplo" é o mesmo da spec 003 (frota padrão completa, com o Cruzador em `J1`–`J3`, o Encouraçado em `E1`–`H1` e o Submarino em `H5`, entre outros).

- **CA-1:** Dado o tabuleiro de exemplo sem disparos, quando o jogo desenha o tabuleiro de ataque, então todas as 100 casas aparecem como não atingidas, e nenhuma letra de navio aparece.
- **CA-2:** Dado o tabuleiro de exemplo com disparos em `J10` (água) e `E1` (acerto), quando o jogo desenha o tabuleiro de ataque, então `J10` aparece como água atingida, `E1` como acerto, e as outras casas como não atingidas.
- **CA-3:** Dado o tabuleiro de exemplo com o Cruzador afundado (disparos em `J1`, `J2` e `J3`), quando o jogo desenha o tabuleiro de ataque, então essas três casas aparecem como definido pela questão em aberto 2.
- **CA-4:** Dado cada resultado de disparo (água, acerto e afundado de cada tipo), quando o jogo monta a mensagem, então ela é o texto em português definido pela questão em aberto 3, e só a de afundado cita o navio.
- **CA-5:** Dado qualquer tabuleiro de ataque, quando o jogo o desenha, então nenhuma linha termina com espaço, e todos os caracteres acima da legenda são ASCII.

As decisões das questões em aberto podem acrescentar ou ajustar critérios.

## Fora de escopo

- Ler o disparo do teclado e o laço da partida (#006).
- Mostrar os dois tabuleiros lado a lado (#006).
- Cores (#010).

## Questões em aberto

- [ ] **1. Símbolos:** casa não atingida como `.`, água atingida como `o` e acerto como `X`? O `~` fica reservado para a água conhecida do tabuleiro próprio.
- [ ] **2. Navio afundado:** depois que um navio afunda, as casas dele passam a mostrar a letra do tipo (por exemplo, `C C C` para o Cruzador), ou continuam como `X`? Mostrar a letra combina com a regra de que o afundado revela o tipo.
- [ ] **3. Mensagens:** "Água.", "Acerto!" e "Afundou o Cruzador!" (com o nome de cada navio)? Para o Porta-Aviões e o Encouraçado, o artigo é "o"; para o Submarino e o Destroyer também, então "Afundou o <navio>!" serve para todos.
- [ ] **4. Disparos recebidos no tabuleiro próprio:** o tabuleiro próprio (004) também deve mostrar os disparos que o oponente fez nele (navio atingido e água atingida)? Isso só é usado no jogo contra a máquina (#007). Pode entrar agora, ou ficar para a #007.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
