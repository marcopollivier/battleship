# 018 — Terceiro nível da máquina

- **Status:** Concluída
- **Roadmap:** #018 (onda 7)
- **Depende de:** 013 (máquina com estratégia)

## Problema

A máquina `dificil` da 013 cerca um navio depois de acertá-lo, mas procura os navios sorteando casas. Um bom jogador de papel procura de um jeito mais esperto: atira onde ainda cabe mais navio, e evita as casas onde nenhum navio que falta caberia. O encerramento da 013 sugeriu esse jogador como um terceiro nível.

## Regras

**O nível (questão 1):**

- `--nivel` ganha um terceiro valor, `mestre`. `dificil` continua o padrão, e `facil` e `dificil` não mudam.
- A mensagem de valor inválido passa a ser `Nível inválido: use facil, dificil ou mestre.`.

**Como a máquina `mestre` escolhe o tiro (questão 2):**

- Ela sabe o mesmo que a `dificil`: onde atirou, o resultado de cada tiro e quais navios afundaram. Ela também sabe a frota da partida (a da cartela ou a do `--frota`), como um jogador que conhece as regras, e por isso sabe quais navios ainda faltam afundar.
- Para cada casa ainda não atingida, ela conta de quantos jeitos os navios que faltam poderiam passar por ela, levando em conta tudo o que sabe:
  - um navio não passa por casas de água, nem por casas de navios afundados, nem encosta num navio afundado;
  - o Porta-Aviões conta nas suas 4 rotações, e os navios em linha na horizontal e na vertical.
- **Procura:** sem acertos em navio boiando, ela atira na casa com a maior contagem.
- **Cerco:** com acertos em navio boiando, ela só conta os jeitos que passam por esses acertos, e atira na casa com a maior contagem. Assim ela segue o navio como a `dificil`, mas já sabe para que lado ele provavelmente continua.
- Quando várias casas empatam na maior contagem, ela sorteia entre elas, com a semente da partida.

**Semente:** com a mesma semente, o mesmo nível e as mesmas entradas, a partida se repete.

**Mensagens:** as mesmas da 007.

## Critérios de aceitação

- **CA-1:** Dado um tabuleiro sem tiros, quando a máquina `mestre` escolhe o tiro, então ele está entre as casas de maior contagem (as do meio do tabuleiro, onde passam mais navios).
- **CA-2:** Dado um canto do tabuleiro cercado de água, onde nenhum navio que falta caberia, quando a máquina escolhe o tiro, então ela nunca atira ali.
- **CA-3:** Dado um acerto em navio boiando, quando a máquina escolhe o tiro, então ela atira numa vizinha do acerto, e, com o navio encostado na borda ou na água de um lado, para o lado em que ele cabe.
- **CA-4:** Dados navios afundados, quando a máquina escolhe o tiro, então ela não conta mais os navios afundados entre os que faltam, nem as casas em volta deles.
- **CA-5:** Dados dois tabuleiros com os mesmos tiros e resultados, mas com navios ainda não achados em lugares diferentes, quando a máquina escolhe o tiro com a mesma semente, então o tiro é o mesmo.
- **CA-6:** Dados 200 tabuleiros sorteados da frota padrão, quando as máquinas `dificil` e `mestre` atiram até afundar a frota, então a `mestre` nunca repete uma casa, sempre termina, e precisa em média de menos tiros que a `dificil`.
- **CA-7:** Dado `--nivel mestre` (ou `MESTRE`), quando a partida roda, então a máquina joga com esta estratégia; e com a mesma semente e as mesmas entradas, as saídas são iguais.
- **CA-8:** Dado `--nivel xyz`, quando o programa começa, então ele mostra `Nível inválido: use facil, dificil ou mestre.` e termina com código 2.

## Fora de escopo

- Mudar o padrão para `mestre`.
- Aprender como o jogador costuma posicionar a frota.
- Uma máquina que posiciona a própria frota de forma esperta (hoje o tabuleiro da máquina é sorteado).

## Questões em aberto

Todas decididas em 2026-10-27, como propostas:

- [x] **1. Nome e padrão:** `mestre`, com `dificil` continuando o padrão. Alternativas: chamar de `especialista`; ou tornar o novo nível o padrão.
- [x] **2. Estratégia:** contar de quantos jeitos os navios que faltam passam por cada casa, na procura e no cerco, como proposto. A alternativa, mais simples e mais fraca, é a `dificil` atirando em xadrez na procura (uma casa sim, outra não), porque todo navio ocupa pelo menos duas casas vizinhas, menos o Submarino.
- [x] **3. Regra da issue #25:** na `mestre`, as casas em volta de navio afundado ficam fora naturalmente, porque nenhum navio pode encostar nele. Se a issue #25 decidir tirar a regra da `dificil`, a `mestre` continua com ela, já que é a regra do jogo usada a favor da máquina.
- [x] **4. Texto:** `Nível inválido: use facil, dificil ou mestre.`.

## Encerramento

Concluída em 2026-10-27, nas tarefas T1–T4. Os 8 critérios de aceitação são cobertos por testes em `internal/board/master_test.go` (CA-1 em `TestMasterShot_Hunt`, CA-2 em `TestMasterShot_NoRoom`, CA-3 em `TestMasterShot_Target`, CA-4 em `TestMasterShot_AfterSunk`, CA-5 em `TestMasterShot_PublicKnowledge`, CA-6 em `TestMasterShot_Game`) e `cmd/init_test.go` (CA-7 em `TestRun_Level` e `TestRun_SameSeed`, CA-8 em `TestRun_InvalidLevel`). A cobertura ficou em 99,7% em `internal/board`, 97,5% em `internal/terminal` e 92,1% em `cmd`.

Em 200 sementes da frota padrão, a máquina `mestre` precisa em média de 54,7 tiros para afundar a frota, a `dificil` de 59,7 e a `facil` de 95,6.

O binário foi rodado com `--cores nao`, as sementes 1, 7, 8, 42, 99 e 2026, `aleatorio` e um tiro em cada casa do tabuleiro, com `--nivel dificil` e `--nivel mestre`. A máquina venceu as 12 partidas: a `mestre` em 44 a 67 tiros (média 55,3), a `dificil` em 54 a 78 (média 65,2). Cada partida rodada duas vezes deu a mesma saída. Nos tiros da `mestre`, dá para ver a procura começar pelo meio do tabuleiro e o cerco seguir o navio para o lado em que ele cabe. `--nivel xyz` mostra a mensagem nova e termina com código 2.

**Desvios em relação ao plano:**

- Os testes do T3 (CA-5 e CA-6) passaram de cara, porque conferem o comportamento feito no T1 e no T2. Os do T1 falharam primeiro na compilação e depois contra um `MasterShot` provisório que só chamava o `SmartShot`; os do T2 falharam contra o código do T1, que ainda contava os acertos pendentes como na procura e atirava longe deles.
- Para conferir a contagem, os testes têm uma contagem própria (`masterBest`), mais simples, que só vale para tabuleiros sem acertos pendentes: cada casa atingida e cada casa em volta de um acerto ficam proibidas.
- No CA-2, tirar os Submarinos afundados não muda as casas de maior contagem no meio do tabuleiro, porque um Submarino soma o mesmo em toda casa onde cabe. Por isso o teste usa um canto fechado por água onde só caberia um Submarino: com o Submarino afundado, a máquina nunca atira ali.
- No CA-4, um subteste novo usa uma frota de um Destroyer e um Cruzador, com só uma linha de 5 casas e uma cruz de 5 casas livres: com o Destroyer afundado, a maior contagem é o meio da linha (onde o Cruzador cabe mais), e não o meio da cruz.
- No CA-3, o caso da borda é um acerto em J2: a máquina atira em J3, para o lado em que o navio tem mais espaço, e não em J1 nem em I2.
- A máquina monta, a cada tiro, um retrato do que sabe de cada casa (`chart`: água, acerto pendente ou navio afundado), para não percorrer os navios a cada consulta de casa. O plano não falava disso.
- Os helpers `strategyChoices` e `samePublicShots` saíram dos testes do `SmartShot` em `machine_test.go`, e servem aos dois, para o `dupl` não acusar código repetido.
- Um teste a mais, `TestMasterShot_FleetSunk`, cobre a volta para o `SmartShot` quando não há contagem (com a frota toda afundada) e o `ErrNoCellsLeft`. Numa partida, a contagem só fica vazia com a frota afundada, porque a posição real de um navio boiando sempre conta.
- A ajuda da opção `--nivel` passou a descrever o `mestre`, e o `TestRun_SameSeed` também roda com `--nivel mestre`.

**Para o roadmap:**

- Decidir se o `mestre` vira o padrão, jogando algumas partidas contra ele.
- A issue #25 continua aberta para a `dificil`; a `mestre` mantém a regra das casas em volta de navio afundado, como decidido na questão 3.
- Uma contagem mais forte, que considera as posições de todos os navios juntos (e não de cada navio sozinho), se a `mestre` parecer fraca.
- A fusão com a 015, que também mexe no `cmd/init.go` e no README.
