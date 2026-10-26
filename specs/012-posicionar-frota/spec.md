# 012 — Posicionar a própria frota

- **Status:** Aprovada
- **Roadmap:** #012 (onda 5)
- **Depende de:** 007 (jogador contra a máquina)

## Problema

O README pede que o jogador "arme seu próprio tabuleiro" antes de jogar contra a máquina, como no bloco de papel. Na partida da 007, a frota do jogador é sorteada. O domínio já valida um posicionamento (002: limites, sobreposição e navios encostados) e já sorteia uma frota (005); falta o jogador dizer, pelo terminal, onde fica cada navio. Com esta spec, o item do TODO do README fica completo.

## Regras

**Quando acontece:**

- No modo contra a máquina, antes do primeiro disparo, o jogador posiciona a sua frota. O modo solo não muda.
- O jogador pode sortear a frota inteira, ou o resto dela, a qualquer momento (veja os comandos).

**Ordem e tela:**

- Os navios são posicionados um de cada vez, do maior para o menor: Porta-Aviões, Encouraçado, Cruzador, Destroyers e Submarinos, só os tipos da frota escolhida (`--frota`, 008).
- No início, o programa mostra as instruções:

  ```
  Posicione a sua frota. Para cada navio, digite a casa e a orientação, por exemplo E5 vertical.
  A casa é a primeira do navio (a de cima ou a da esquerda). No Porta-Aviões, é o meio da barra do T.
  Comandos: aleatorio sorteia o resto da frota, recomecar tira todos os navios, sair desiste.
  ```

- Antes de cada navio, o programa mostra o tabuleiro próprio do jogador (004) e pede a posição com o nome do navio, o tamanho e as orientações possíveis (lidas do domínio, nota da 005):

  ```
  Porta-Aviões, 5 casas (cima, baixo, esquerda, direita):
  Destroyer, 2 casas (horizontal, vertical):
  Submarino, 1 casa:
  ```

**Entrada:**

- Uma linha com a casa e a orientação, separadas por espaço: `E5 vertical`, `a1 horizontal`, `C3 baixo`.
- A orientação aceita o nome inteiro ou a primeira letra, em maiúscula ou minúscula: `h`/`horizontal`, `v`/`vertical`, `c`/`cima`, `b`/`baixo`, `e`/`esquerda`, `d`/`direita`.
- O Submarino, que tem uma casa só, dispensa a orientação: `A1`.
- Com um posicionamento válido, o navio entra no tabuleiro e o programa passa para o próximo.
- Com uma entrada inválida, o programa mostra uma mensagem, não posiciona o navio e pede de novo o mesmo navio:

  | Problema | Mensagem |
  | -------- | -------- |
  | Formato (falta a orientação, palavras demais, casa mal escrita) | `Posição inválida: digite a casa e a orientação, por exemplo E5 vertical.` |
  | Orientação desconhecida ou que não vale para o navio | `Orientação inválida para o <navio>: use <orientações>.` |
  | Navio sai do tabuleiro | `Fora do tabuleiro: o <navio> não cabe nessa posição.` |
  | Navio em cima de outro | `Posição ocupada: o <navio> ficaria em cima de outro navio.` |
  | Navio encostado em outro | `Muito perto: os navios não podem se encostar, nem pela diagonal.` |

**Comandos**, aceitos no lugar de uma posição, em maiúscula ou minúscula:

- `aleatorio`: sorteia os navios que faltam, mantendo os que o jogador já posicionou. Se não for possível encaixar o resto, mostra `Não foi possível sortear o resto da frota. Mude a posição dos navios com recomecar.` e continua pedindo o mesmo navio.
- `recomecar`: tira todos os navios do tabuleiro e volta ao primeiro.
- `sair`, ou o fim da entrada: desiste antes da partida. O programa mostra `Você desistiu antes da partida.`, depois `Semente: N`, e termina com código 0.

**Fim do posicionamento:**

- Depois do último navio, o programa mostra o tabuleiro próprio completo e `Frota posicionada!`, e a partida da 007 começa.

**Semente:**

- A semente passa a sortear primeiro o tabuleiro da máquina, depois o que o jogador pedir com `aleatorio`, e depois os disparos da máquina. Com a mesma semente e as mesmas entradas, a partida se repete.
- Isso muda a ordem da 007: uma semente de uma partida anterior a esta spec não repete aquela partida.

## Critérios de aceitação

- **CA-1:** Dado o modo contra a máquina com a frota padrão, quando a partida começa, então o programa mostra as instruções, o tabuleiro vazio e o pedido do Porta-Aviões com as 4 orientações.
- **CA-2:** Dadas as 8 posições válidas da frota padrão, quando o jogador as digita, então os navios são pedidos do maior para o menor, o tabuleiro aparece antes de cada um, e no fim aparece `Frota posicionada!` seguido da tela da partida (007) com os navios onde o jogador os pôs.
- **CA-3:** Dadas as entradas `e5 V`, `E5 vertical` e `A1` para um Submarino, quando o jogo as interpreta, então todas são aceitas; e as orientações aceitam a primeira letra e o nome inteiro, em maiúscula ou minúscula.
- **CA-4:** Dada cada linha da tabela de mensagens (por exemplo `E5`, `E5 vertical sim`, `E5 diagonal`, `C3 horizontal` para o Porta-Aviões, `J9 vertical` para o Encouraçado, uma casa ocupada e uma casa encostada), quando o jogador a digita, então o programa mostra a mensagem da tabela e pede de novo o mesmo navio.
- **CA-5:** Dado um tabuleiro com alguns navios posicionados, quando o jogador digita `aleatorio`, então os navios que já estavam continuam no lugar, os que faltam são sorteados, e o posicionamento termina.
- **CA-6:** Dado um tabuleiro em que o resto da frota não cabe, quando o jogador digita `aleatorio`, então o programa mostra a mensagem de que não foi possível sortear e pede o mesmo navio de novo.
- **CA-7:** Dado um tabuleiro com navios, quando o jogador digita `recomecar`, então o tabuleiro fica vazio e o próximo pedido é o do Porta-Aviões.
- **CA-8:** Dado o posicionamento em andamento, quando o jogador digita `sair` ou a entrada termina, então o programa mostra `Você desistiu antes da partida.` e `Semente: N`, e termina com código 0.
- **CA-9:** Dadas duas execuções com `--semente 42` e as mesmas entradas (com e sem `aleatorio`), quando elas rodam, então as saídas são iguais.
- **CA-10:** Dado o modo solo, quando a partida começa, então não há posicionamento, como na 006.

## Fora de escopo

- Desfazer só o último navio (o jogador usa `recomecar`).
- Escolher a ordem dos navios.
- Avisar, durante o posicionamento, que o resto da frota não cabe mais. O jogador descobre ao tentar e usa `recomecar`.
- Posicionar com setas ou mouse.
- Salvar uma frota para outras partidas.

## Questões em aberto

Todas decididas em 2026-10-26:

- [x] **1. Posicionar é o padrão:** no modo contra a máquina, o posicionamento manual sempre acontece, e quem não quiser digita `aleatorio` no primeiro navio. A alternativa é deixar a frota sorteada como padrão e ligar o posicionamento com uma opção (`--posicionar`).
- [x] **2. Ordem e casa de referência:** do maior para o menor; a casa é a primeira do navio (de cima ou da esquerda), e no Porta-Aviões é o meio da barra do T, como no domínio (002).
- [x] **3. Entrada:** `casa orientação` numa linha, orientação pelo nome ou pela primeira letra, e o Submarino sem orientação.
- [x] **4. Comandos:** `aleatorio` (sorteia só o que falta), `recomecar` (começa de novo) e `sair`. Os comandos são escritos sem acento, para facilitar a digitação.
- [x] **5. Semente:** passa a sortear o tabuleiro da máquina primeiro, mudando a ordem da 007.
- [x] **6. Mensagens:** os textos das instruções, dos pedidos, dos erros e dos comandos.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
