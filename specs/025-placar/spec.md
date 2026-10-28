# 025 — Placar entre partidas

- **Status:** Concluída
- **Roadmap:** #025 (onda 8), issue #35
- **Depende de:** 018 (nível `mestre`)

## Problema

O jogo esquece tudo quando a partida acaba. O jogador não sabe quantas vezes já venceu a máquina, nem se está melhorando. Com três níveis (013 e 018), um placar por nível também ajuda a decidir se vale subir de nível, e se o `mestre` deve virar o padrão (issue #38).

## Regras

**O que conta (questões 1 e 2):**

- Só as partidas contra a máquina com a frota da cartela (sem `--frota`). Frotas diferentes mudam a dificuldade e não entram no placar.
- Para cada nível (`facil`, `dificil`, `mestre`), o placar guarda:
  - vitórias;
  - derrotas;
  - a melhor vitória: o menor número de disparos do jogador numa vitória.
- Uma desistência (`sair` ou fim da entrada) depois do primeiro disparo conta como derrota. Desistir no posicionamento, antes de atirar, não conta.
- Partidas com `--semente` contam como as outras.

**Onde fica (questão 3):**

- Num arquivo na pasta de configuração do usuário: `~/Library/Application Support/battleship/placar.json` no macOS, `~/.config/battleship/placar.json` no Linux, e a pasta equivalente no Windows.
- Se o arquivo não existir, o placar começa zerado. Se não der para ler ou gravar o arquivo, o jogo mostra `Não foi possível guardar o placar.` no `stderr` e continua: a partida nunca é interrompida por causa do placar.

**Quando aparece (questão 4):**

- No fim de cada partida que conta, antes da semente, uma linha com o placar do nível jogado:

  ```
  Placar no nível dificil: 3 vitórias, 2 derrotas. Melhor vitória: 41 disparos.
  ```

  Sem vitórias, a parte da melhor vitória não aparece.
- A opção `--placar` mostra o placar dos três níveis e termina com código 0, sem começar partida.

**Zerar (questão 5):**

- A opção `--zerar-placar` apaga o placar, mostra `Placar zerado.` e termina com código 0.

## Critérios de aceitação

- **CA-1:** Dada uma vitória contra a máquina com a frota da cartela, quando a partida termina, então a vitória e o número de disparos entram no placar do nível, e a linha do placar aparece antes da semente.
- **CA-2:** Dada uma derrota, ou uma desistência depois do primeiro disparo, quando a partida termina, então uma derrota entra no placar do nível.
- **CA-3:** Dada uma desistência no posicionamento, uma partida no modo solo, ou uma partida com `--frota`, quando ela termina, então o placar não muda e a linha não aparece.
- **CA-4:** Dadas duas vitórias com 50 e 41 disparos, quando o placar é mostrado, então a melhor vitória é 41.
- **CA-5:** Dado `--placar`, quando o programa roda, então mostra uma linha por nível e termina com código 0, sem partida.
- **CA-6:** Dado `--zerar-placar`, quando o programa roda, então o placar volta a zero, aparece `Placar zerado.` e o código é 0.
- **CA-7:** Dado um arquivo de placar que não pode ser lido ou gravado (por exemplo, corrompido ou numa pasta sem permissão), quando a partida termina, então aparece `Não foi possível guardar o placar.` no `stderr`, e o resto da saída e o código de saída são os de uma partida normal.
- **CA-8:** Dadas duas execuções com a mesma semente e as mesmas entradas, partindo do mesmo placar, quando elas rodam, então as saídas são iguais.

## Fora de escopo

- Placar por frota ou no modo solo.
- Ranking online, ou placar compartilhado entre pessoas.
- Histórico de cada partida (data, semente, jogadas).
- Nomes de jogadores: o placar é de quem usa aquele computador.

## Questões em aberto

Todas decididas em 2026-10-27, como propostas:

- [x] **1. O que conta:** só contra a máquina, com a frota da cartela, separado por nível. A alternativa é contar também frotas do `--frota` e o modo solo, num placar à parte.
- [x] **2. Desistência:** depois do primeiro disparo conta como derrota; no posicionamento, não conta. A alternativa é nunca contar desistências.
- [x] **3. Onde guardar:** um arquivo JSON na pasta de configuração do usuário, achada pela biblioteca padrão do Go. A alternativa é um arquivo na pasta atual, que mudaria conforme de onde o jogo é rodado.
- [x] **4. Quando mostrar:** uma linha no fim de cada partida que conta, e `--placar` para ver os três níveis.
- [x] **5. Zerar:** `--zerar-placar`. A alternativa é a pessoa apagar o arquivo, com o caminho explicado no README.
- [x] **6. Textos:** a linha do placar, `Placar zerado.` e `Não foi possível guardar o placar.`.

## Encerramento

Concluída em 2026-10-27, nas tarefas T1–T5. Os 8 critérios de aceitação são cobertos por testes em `internal/score/score_test.go`, `internal/terminal/score_test.go`, `internal/terminal/machine_test.go`, `cmd/score_test.go` e `cmd/init_test.go`: CA-1 em `TestRun_ScoreWin` e `TestScoreLine`, CA-2 em `TestRun_ScoreLoss` e `TestPlayVsMachine_Result`, CA-3 em `TestRun_ScoreNotCounted`, CA-4 em `TestBoard_AddWin` e `TestScoreLine`, CA-5 em `TestRun_ShowScore`, CA-6 em `TestRun_ResetScore`, CA-7 em `TestLoad_Corrupt`, `TestSave_Unwritable`, `TestRun_ScoreError` e `TestRun_ScoreCommandError`, CA-8 em `TestRun_SameSeed`. `TestSaveLoad` e `TestDefaultPath` conferem o formato do arquivo e o caminho. A cobertura ficou em 86,0% em `internal/score`, 97,9% em `internal/terminal`, 99,7% em `internal/board` e 91,9% em `cmd` (o `Run` público, que só acha o caminho real e chama o `run`, fica de fora dos testes de propósito).

O binário foi rodado com `HOME` e `XDG_CONFIG_HOME` apontando para uma pasta temporária, `--cores nao`, `aleatorio` e um tiro em cada casa, de A1 a J10. Com `--nivel facil` e as sementes 1 a 5, o jogador venceu as 5 partidas (97, 97, 94, 91 e 94 disparos), e a linha do placar foi de `1 vitória, 0 derrotas. Melhor vitória: 97 disparos.` a `5 vitórias, 0 derrotas. Melhor vitória: 91 disparos.`, sempre antes da semente. Com `--nivel dificil` e a semente 7, a máquina venceu duas vezes, e o placar foi a `0 vitórias, 2 derrotas.`. `sair` antes do primeiro disparo e uma partida com `--frota S=1` não mostraram a linha nem mudaram o arquivo. `--placar` mostrou os três níveis com código 0, e o arquivo ficou como `{"dificil":{"wins":0,"losses":2,"best_win":0},"facil":{"wins":5,"losses":0,"best_win":91}}`. Com o arquivo corrompido, `--placar` e uma partida mostraram `Não foi possível guardar o placar.` no `stderr`, com código 0 e a partida normal. `--zerar-placar` mostrou `Placar zerado.` com código 0, e o `--placar` seguinte mostrou tudo zerado. O placar real (`~/Library/Application Support/battleship`) não foi criado.

**Desvios em relação ao plano:**

- O helper dos testes do `cmd` se chamava `run`, o mesmo nome da função interna nova. Ele virou `runCmd(t, ...)`, que usa um placar numa pasta de `t.TempDir()`, e `runScore(caminho, ...)`, para os testes que preparam o arquivo. Assim, todos os testes do `cmd`, inclusive os antigos, passam pelo `run` e nunca pelo placar real.
- As partidas contra a máquina com a frota da cartela agora mostram a linha do placar. Por isso o `levelGame` do `TestRun_Level` recebe o nível e termina com a linha de uma derrota, e o `TestRun_SameSeed` confere que a linha aparece (e que não aparece com `--frota`).
- `Load` e `Save` devolvem `ErrNoPath` para um caminho vazio, que é o caso sem pasta de configuração. Sem isso, o `Load("")` leria um placar vazio, como um arquivo que não existe, e o `Save("")` gravaria na pasta atual.
- O arquivo temporário do `Save` fica na mesma pasta do placar, com o nome `placar.json.*.tmp`, para o `os.Rename` não cruzar sistemas de arquivos.
- Para o `run` não passar do limite do `gocyclo`, a contagem e o registro ficaram em `counts` e `recordScore`, e as duas opções em `scoreCommand`, que devolve o código e se a partida deve começar.
- Com `--placar` e `--zerar-placar` juntos, o placar é zerado e depois mostrado. A spec não falava desse caso.
- Se `--placar` não consegue ler o arquivo, ou `--zerar-placar` não consegue gravar, o aviso aparece e o código continua 0, como numa partida. A spec não dizia o código nesses casos. `--zerar-placar` sobre um arquivo corrompido funciona e grava um placar vazio.
- A linha do placar sai sem cores, mesmo com `--cores sim`.
- Uma partida que termina com erro (código 1) não entra no placar.
- O `shotCount` passou a usar o `plural`, o mesmo helper dos números da linha do placar.
- O teste de cada tarefa falhou primeiro na compilação (tipos e funções novas); o `TestRun_ShowScore` e o `TestRun_ResetScore` falharam também com código 2, pela opção desconhecida, e o `TestRun_ScoreError` e o `TestRun_ScoreCommandError` pela falta do aviso no `stderr`.
- O README ganhou uma seção "Placar", com o que conta e onde fica o arquivo em cada sistema, a linha do placar no exemplo e as duas opções na tabela.
- No merge do `master` com a 022 e a 024, o `PlayVsMachine` passou a receber o `Display` e a devolver o `Result` ao mesmo tempo. A linha do placar sai depois da tela final, antes da semente. O `--versao` é tratado antes das opções do placar. Para o `run` continuar dentro do limite do `gocyclo`, a escrita da versão foi para `printVersion`, em `cmd/version.go`.

**Para o roadmap:**

- A fusão com a 024, que põe o `--versao` logo depois do `Parse`: as opções do placar devem vir depois dele, como diz o plano.
- A fusão com a 022, que troca o `Style` do `PlayVsMachine` por um `Display` e refaz as telas de fim: a linha do placar continua antes da semente, depois da tela de fim.
- Usar o placar para decidir se o `mestre` vira o padrão (issue #38).
- Duas partidas terminando ao mesmo tempo no mesmo computador gravam o placar uma por cima da outra, e um resultado se perde. Se isso importar, uma trava no arquivo resolve.
