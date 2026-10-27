# 025 — Placar entre partidas

- **Status:** Rascunho
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

- [ ] **1. O que conta:** só contra a máquina, com a frota da cartela, separado por nível. A alternativa é contar também frotas do `--frota` e o modo solo, num placar à parte.
- [ ] **2. Desistência:** depois do primeiro disparo conta como derrota; no posicionamento, não conta. A alternativa é nunca contar desistências.
- [ ] **3. Onde guardar:** um arquivo JSON na pasta de configuração do usuário, achada pela biblioteca padrão do Go. A alternativa é um arquivo na pasta atual, que mudaria conforme de onde o jogo é rodado.
- [ ] **4. Quando mostrar:** uma linha no fim de cada partida que conta, e `--placar` para ver os três níveis.
- [ ] **5. Zerar:** `--zerar-placar`. A alternativa é a pessoa apagar o arquivo, com o caminho explicado no README.
- [ ] **6. Textos:** a linha do placar, `Placar zerado.` e `Não foi possível guardar o placar.`.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
