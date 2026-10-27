# 010 — Cores no terminal

- **Status:** Em andamento
- **Roadmap:** #010 (onda 6)
- **Depende de:** 009 (tabuleiro de ataque)

## Problema

Os tabuleiros usam só caracteres (`~`, `o`, `X`, letras dos navios), e num tabuleiro de 100 casas é difícil achar de relance onde estão os acertos, a água e os navios. A 004 deixou as cores para depois. Esta spec pinta os símbolos dos dois tabuleiros sem mudar os caracteres, para que o jogo continue legível sem cor.

## Regras

**O que muda:**

- Cada símbolo dos tabuleiros ganha uma cor; o caractere continua o mesmo. Quem não vê cores, ou usa um terminal sem cores, lê o jogo como hoje.
- Cores (questão 2):

  | Símbolo | Onde | Significado | Cor |
  | ------- | ---- | ----------- | --- |
  | `~` | tabuleiro próprio | água | azul |
  | `S` `D` `C` `E` `P` | tabuleiro próprio | navio inteiro | branco em negrito |
  | `X` | os dois | navio atingido | vermelho em negrito |
  | `o` | os dois | tiro na água | ciano |
  | `.` | tabuleiro de ataque | casa não atingida | cinza (apagado) |
  | `S` `D` `C` `E` `P` | tabuleiro de ataque | navio afundado | vermelho |

- O cabeçalho (letras das colunas e números das fileiras), as legendas e as mensagens continuam sem cor (questão 3).

**Quando há cores (questão 1):**

- As cores estão sempre ligadas, inclusive quando a saída vai para um arquivo (`go run . > partida.txt`), que então guarda os códigos de cor.
- A opção `--cores nao` desliga as cores. `--cores sim` é o padrão. Um valor inválido mostra `Cores inválidas: use sim ou nao.` e termina com código 2, como as outras opções.

**Sem efeito no resto do jogo:**

- Com as cores desligadas, a saída é exatamente a de hoje.
- A `--semente` repete a partida do mesmo jeito, com ou sem cores.

## Critérios de aceitação

- **CA-1:** Dado um tabuleiro próprio com água, navios, acertos e tiros na água, quando ele é desenhado com cores, então cada símbolo aparece com a cor da tabela, e os caracteres são os mesmos de sem cores.
- **CA-2:** Dado um tabuleiro de ataque com casas não atingidas, tiros na água, acertos e um navio afundado, quando ele é desenhado com cores, então cada símbolo aparece com a cor da tabela.
- **CA-3:** Dado qualquer tabuleiro, quando ele é desenhado com cores, então o cabeçalho, os números das fileiras e a legenda não têm códigos de cor, e tirando os códigos de cor o texto é igual ao desenho sem cores.
- **CA-4:** Dado `--cores nao`, quando a partida roda, então a saída não tem nenhum código de cor e é igual à de antes desta spec.
- **CA-5:** Dada a partida sem `--cores`, ou com `--cores sim` ou `--cores SIM`, quando ela roda, então os tabuleiros saem com cores, mesmo com a saída indo para um arquivo.
- **CA-6:** Dado `--cores xyz`, quando o programa começa, então ele mostra `Cores inválidas: use sim ou nao.` e termina com código 2.
- **CA-7:** Dados os modos solo e contra a máquina, incluindo o posicionamento da frota (012), quando as cores estão ligadas, então todos os tabuleiros mostrados aparecem com cores.

## Fora de escopo

- Cores nas mensagens (`Acerto!`, `Água.`, vitória e derrota).
- Escolher a paleta (cores configuráveis).
- Fundo colorido nas casas, ou trocar os caracteres por blocos e emojis.
- Desligar as cores sozinho quando a saída não é um terminal, ou com a variável `NO_COLOR` (questão 1).
- Terminais sem suporte a códigos ANSI (por exemplo o `cmd.exe` antigo do Windows): eles podem usar `--cores nao`.

## Questões em aberto

Todas decididas em 2026-10-26:

- [x] **1. Quando ligar as cores:** sempre ligadas, com `--cores nao` para desligar. As alternativas eram ligar só quando a saída é um terminal e sem `NO_COLOR`, ou deixar desligadas por padrão.
- [x] **2. Paleta:** as cores da tabela acima.
- [x] **3. Mensagens:** só os tabuleiros ficam coloridos.
- [x] **4. Nome da opção:** `--cores` com `sim` e `nao`, sem acento, como os comandos da 012.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
