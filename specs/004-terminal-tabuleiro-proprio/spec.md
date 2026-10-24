# 004 — Terminal: tabuleiro próprio

- **Status:** Rascunho
- **Roadmap:** #004 (onda 1)
- **Depende de:** 001 (coordenadas), 002 (tabuleiro e frota)

## Problema

O jogador precisa ver o próprio tabuleiro com os navios posicionados. Hoje o programa só imprime um tabuleiro vazio montado à mão em `cmd/init.go`, sem relação com o domínio de `internal/board`. Além disso, o domínio usa nomes em inglês (`Cruiser`, `Up`...), e o jogador precisa ver navios e orientações em português. Esta spec cria a camada de terminal que as próximas features (#009, #006, #007) vão usar para mostrar o jogo.

## Regras

- O tabuleiro é desenhado como hoje: uma linha de cabeçalho com as colunas `A` a `J` e uma linha por fileira, de `1` a `10`, com o número alinhado à direita em duas posições.
- Uma casa sem navio aparece como `~`, como hoje.
- Uma casa com navio aparece com o símbolo de navio definido pela questão em aberto 1.
- Nenhuma linha termina com espaço sobrando. Hoje cada fileira termina com um espaço a mais.
- Os nomes exibidos ao jogador são em português:
  - **Navios:** Submarino, Destroyer, Cruzador, Encouraçado e Porta-Aviões.
  - **Orientações:** horizontal, vertical, cima, baixo, esquerda e direita.
- O programa (`make run`) continua mostrando um tabuleiro, agora desenhado por esta camada a partir de um tabuleiro do domínio.

## Critérios de aceitação

- **CA-1:** Dado um tabuleiro vazio, quando o jogo o desenha, então a saída é igual à do programa atual, só que sem o espaço no fim das linhas:

  ```
     A B C D E F G H I J
   1 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
   2 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
   …
  10 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
  ```

- **CA-2:** Dado um tabuleiro com um Cruzador em `A1` na horizontal, quando o jogo o desenha, então `A1`, `B1` e `C1` aparecem com o símbolo de navio e as outras casas com `~`.
- **CA-3:** Dado um tabuleiro com o Porta-Aviões em `E5` para `baixo`, quando o jogo o desenha, então as casas `D5`, `E5`, `F5`, `E6` e `E7` aparecem com o símbolo de navio.
- **CA-4:** Dado cada tipo de navio, quando o jogo pede o nome para exibição, então recebe o nome em português da tabela das regras. O mesmo vale para cada orientação.
- **CA-5:** Dado o programa compilado, quando ele é executado, então imprime o tabuleiro desenhado por esta camada, sem diferença em relação ao CA-1.

As decisões das questões em aberto podem acrescentar ou ajustar critérios.

## Fora de escopo

- O tabuleiro de ataque e os disparos (#009).
- Ler comandos do teclado (#006 e #007).
- Mostrar os dois tabuleiros lado a lado. Fica para quando houver uma partida (#006).

## Questões em aberto

- [ ] **1. Símbolo dos navios:** um único símbolo para todos os navios (por exemplo, `■`) ou uma letra por tipo (`S` Submarino, `D` Destroyer, `C` Cruzador, `E` Encouraçado, `P` Porta-Aviões)?
- [ ] **2. Caracteres:** usar Unicode (`■`) ou ficar só em ASCII (`#`), para funcionar em qualquer terminal?
- [ ] **3. Legenda:** mostrar uma legenda com os nomes dos navios abaixo do tabuleiro?
- [ ] **4. Cores:** usar cores no terminal (água azul, navios em destaque) ou deixar para depois?

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
