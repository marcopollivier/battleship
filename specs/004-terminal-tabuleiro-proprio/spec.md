# 004 — Terminal: tabuleiro próprio

- **Status:** Em andamento
- **Roadmap:** #004 (onda 1)
- **Depende de:** 001 (coordenadas), 002 (tabuleiro e frota)

## Problema

O jogador precisa ver o próprio tabuleiro com os navios posicionados. Hoje o programa só imprime um tabuleiro vazio montado à mão em `cmd/init.go`, sem relação com o domínio de `internal/board`. Além disso, o domínio usa nomes em inglês (`Cruiser`, `Up`...), e o jogador precisa ver navios e orientações em português. Esta spec cria a camada de terminal que as próximas features (#009, #006, #007) vão usar para mostrar o jogo.

## Regras

- O tabuleiro é desenhado como hoje: uma linha de cabeçalho com as colunas `A` a `J` e uma linha por fileira, de `1` a `10`, com o número alinhado à direita em duas posições.
- Uma casa sem navio aparece como `~`, como hoje.
- Uma casa com navio aparece com a letra do tipo do navio: `S` Submarino, `D` Destroyer, `C` Cruzador, `E` Encouraçado e `P` Porta-Aviões.
- Os símbolos do tabuleiro (`~` e as letras) são só ASCII, para funcionar em qualquer terminal. Os nomes em português da legenda mantêm os acentos.
- Abaixo do tabuleiro, depois de uma linha em branco, aparece a legenda:

  ```
  Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões
  ```

- A saída não usa cores.
- Nenhuma linha termina com espaço sobrando. Hoje cada fileira termina com um espaço a mais.
- Os nomes exibidos ao jogador são em português:
  - **Navios:** Submarino, Destroyer, Cruzador, Encouraçado e Porta-Aviões.
  - **Orientações:** horizontal, vertical, cima, baixo, esquerda e direita.
- O programa (`make run`) continua mostrando um tabuleiro, agora desenhado por esta camada a partir de um tabuleiro do domínio.

## Critérios de aceitação

- **CA-1:** Dado um tabuleiro vazio, quando o jogo o desenha, então a saída é igual à do programa atual, sem o espaço no fim das linhas e com a legenda no fim:

  ```
     A B C D E F G H I J
   1 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
   2 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
   …
  10 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~

  Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões
  ```

- **CA-2:** Dado um tabuleiro com um Cruzador em `A1` na horizontal, quando o jogo o desenha, então `A1`, `B1` e `C1` aparecem com `C` e as outras casas com `~`.
- **CA-3:** Dado um tabuleiro com o Porta-Aviões em `E5` para `baixo`, quando o jogo o desenha, então as casas `D5`, `E5`, `F5`, `E6` e `E7` aparecem com `P`.
- **CA-4:** Dado cada tipo de navio, quando o jogo pede o nome para exibição, então recebe o nome em português da tabela das regras. O mesmo vale para cada orientação.
- **CA-5:** Dado o programa compilado, quando ele é executado, então imprime o tabuleiro desenhado por esta camada, sem diferença em relação ao CA-1.
- **CA-6:** Dado um tabuleiro completo com a frota padrão, quando o jogo o desenha, então cada navio aparece com a letra do seu tipo, e a saída inteira só tem caracteres ASCII acima da legenda.

## Fora de escopo

- O tabuleiro de ataque e os disparos (#009).
- Ler comandos do teclado (#006 e #007).
- Mostrar os dois tabuleiros lado a lado. Fica para quando houver uma partida (#006).
- Cores no terminal. Fica como item novo no roadmap, a ser incluído por quem coordena no fim da onda 1.

## Questões em aberto

Todas decididas em 2026-10-24:

- [x] **1. Símbolo dos navios:** uma letra por tipo: `S`, `D`, `C`, `E`, `P` (CA-2, CA-3, CA-6).
- [x] **2. Caracteres:** só ASCII nos símbolos do tabuleiro (CA-6).
- [x] **3. Legenda:** sim, numa linha abaixo do tabuleiro (CA-1).
- [x] **4. Cores:** ficam para depois, como item novo no roadmap (fora de escopo).

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
