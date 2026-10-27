# 010 — Cores no terminal

- **Status:** Rascunho
- **Roadmap:** #010 (onda 6)
- **Depende de:** 009 (tabuleiro de ataque)

## Problema

Os tabuleiros usam só caracteres (`~`, `o`, `X`, letras dos navios), e num tabuleiro de 100 casas é difícil achar de relance onde estão os acertos, a água e os navios. A 004 deixou as cores para depois. Esta spec pinta os símbolos dos dois tabuleiros, sem mudar os símbolos, para que o jogo continue legível sem cor.

## Regras

**O que muda:**

- Cada símbolo dos tabuleiros ganha uma cor; o caractere continua o mesmo. Quem não vê cores, ou usa um terminal sem cores, lê o jogo como hoje.
- Cores propostas (questão 2):

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

- As cores aparecem quando a saída é um terminal.
- Não aparecem quando a saída vai para um arquivo ou outro programa (por exemplo `go run . > partida.txt`), para não gravar códigos de cor no texto.
- Não aparecem quando a variável de ambiente `NO_COLOR` está definida, uma convenção comum entre programas de terminal (no-color.org).
- A opção `--cores` escolhe à mão: `--cores sim` liga, `--cores nao` desliga, `--cores auto` (padrão) segue as duas regras acima. Um valor inválido mostra `Cores inválidas: use sim, nao ou auto.` e termina com código 2, como as outras opções.

**Sem efeito no resto do jogo:**

- Com as cores desligadas, a saída é exatamente a de hoje.
- A `--semente` repete a partida do mesmo jeito, com ou sem cores.

## Critérios de aceitação

- **CA-1:** Dado um tabuleiro próprio com água, navios, acertos e tiros na água, quando ele é desenhado com cores, então cada símbolo aparece com a cor da tabela, e os caracteres são os mesmos de sem cores.
- **CA-2:** Dado um tabuleiro de ataque com casas não atingidas, tiros na água, acertos e um navio afundado, quando ele é desenhado com cores, então cada símbolo aparece com a cor da tabela.
- **CA-3:** Dado qualquer tabuleiro, quando ele é desenhado com cores, então o cabeçalho, os números das fileiras e a legenda não têm códigos de cor, e tirando os códigos de cor o texto é igual ao desenho sem cores.
- **CA-4:** Dado `--cores nao`, ou `--cores auto` com a saída indo para um arquivo, ou `--cores auto` com `NO_COLOR` definida, quando a partida roda, então a saída não tem nenhum código de cor.
- **CA-5:** Dado `--cores sim`, quando a partida roda com a saída indo para um arquivo, então os tabuleiros saem com cores.
- **CA-6:** Dado `--cores xyz`, quando o programa começa, então ele mostra `Cores inválidas: use sim, nao ou auto.` e termina com código 2.
- **CA-7:** Dados os modos solo e contra a máquina, incluindo o posicionamento da frota (012), quando as cores estão ligadas, então todos os tabuleiros mostrados aparecem com cores.

## Fora de escopo

- Cores nas mensagens (`Acerto!`, `Água.`, vitória e derrota), se a questão 3 decidir assim.
- Escolher a paleta (cores configuráveis).
- Fundo colorido nas casas, ou trocar os caracteres por blocos e emojis.
- Terminais sem suporte a códigos ANSI (por exemplo o `cmd.exe` antigo do Windows): eles podem ligar `--cores nao`.

## Questões em aberto

- [ ] **1. Quando ligar as cores:** automático (terminal sim, arquivo não, `NO_COLOR` desliga), com `--cores sim|nao|auto` para escolher à mão, como proposto. Alternativas: sempre ligadas, com só uma opção para desligar; ou desligadas por padrão, com `--cores` para ligar.
- [ ] **2. Paleta:** as cores da tabela acima.
- [ ] **3. Mensagens:** só os tabuleiros ficam coloridos, como proposto, ou também as mensagens de resultado (por exemplo `Acerto!` em vermelho e `Água.` em azul)?
- [ ] **4. Nome da opção:** `--cores` com `sim`, `nao` e `auto`, sem acento, como os comandos da 012.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
