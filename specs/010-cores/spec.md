# 010 — Cores no terminal

- **Status:** Concluída
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

Concluída em 2026-10-26, nas tarefas T1–T5. Os 7 critérios de aceitação são cobertos por testes em `internal/terminal/color_test.go` (CA-1 em `TestRenderBoard_Colored`, CA-2 em `TestRenderAttackBoard_Colored`, CA-3 em `TestRender_ColoredMatchesPlain`), `solo_test.go`, `machine_test.go` e `placement_test.go` (CA-7 em `TestPlaySolo_Colored`, `TestPlayVsMachine_Colored` e `TestPlayVsMachine_PlacementColored`) e `cmd/init_test.go` (CA-4 e CA-5 em `TestRun_Colors`, CA-6 em `TestRun_InvalidColors`). Os testes que já existiam passaram a usar `Plain` sem mudar nenhum resultado, o que confere a regra "com as cores desligadas, a saída é exatamente a de hoje". A cobertura ficou em 99,5% em `internal/board`, 96,7% em `internal/terminal` e 90,6% em `cmd`.

Executando `go run . --semente 42` com `aleatorio`, `A1` e `sair`, com a saída num arquivo: os tabuleiros saem com os códigos de cor (água azul, navios em branco negrito, casas não atingidas em cinza e o Submarino afundado em A1 em vermelho), e o cabeçalho, os números e as legendas saem sem cor. Com `--cores nao`, a mesma partida não tem nenhum código de cor e é igual à saída com cores depois de tirar os códigos. `--cores xyz` mostra `Cores inválidas: use sim ou nao.` e termina com código 2.

**Desvios em relação ao plano:**

- A leitura do `--cores` ficou numa função interna, `chooseStyle`, em vez de comparar o valor dentro do `Run`, para o `Run` não ficar mais complexo.
- Na T1, o `playerAttack` ainda desenhava sempre com `Plain`; o `Style` guardado no `playerAttack` (e, por ele, no `machineGame`) entrou na T4, junto com a mudança de assinatura do `PlaySolo` e do `PlayVsMachine`.
- Os testes do CA-7 conferem a saída exata de uma partida curta com cores e, numa partida mais longa (vitória no solo, vitória da máquina, posicionamento com `aleatorio`), que a saída sem os códigos é igual à sem cores e que toda fileira de tabuleiro tem cor. O `TestRun_Colors` também confere que a mesma semente joga a mesma partida contra a máquina com e sem cores.

**Para o roadmap:**

- Cores nas mensagens (`Acerto!`, `Água.`, vitória e derrota), que ficaram fora desta spec.
- Desligar as cores sozinho quando a saída não é um terminal, ou com a variável `NO_COLOR`.
