# Roadmap

Backlog ordenado. Cada item vira uma pasta `specs/NNN-slug/` quando for especificado. Itens podem ser reordenados, divididos ou removidos a qualquer momento.

Status: `Backlog` → `Spec` → `Em andamento` → `Concluída`

Features da mesma onda não dependem umas das outras e são tocadas em paralelo (veja "Trabalho em paralelo" em `specs/README.md`). Durante uma onda, esta tabela só é atualizada por PRs `docs/` de quem coordena.

| # | Feature | Origem | Depende de | Onda | Status |
| - | ------- | ------ | ---------- | ---- | ------ |
| 001 | Coordenadas: converter `A1`…`J10` em posição e vice-versa, rejeitar entradas inválidas | Regras do README | — | — | Concluída |
| 002 | Tabuleiro e frota: tipos de navio e formatos (linha e T), quantidades padrão, posicionamento com rotação, validação de limites e sobreposição | Regras do README | 001 | — | Concluída |
| 003 | Disparo: resultado água/acerto/afundado, tiro repetido, fim de jogo quando toda a frota afunda | Regras do README | 002 | 1 | Concluída |
| 004 | Terminal, tabuleiro próprio: desenhar o tabuleiro com os navios do jogador e os nomes em português de navios e orientações | `cmd/init.go` atual | 002 | 1 | Concluída |
| 005 | Posicionamento aleatório da frota | TODO do README | 002 | 1 | Concluída |
| 009 | Terminal, tabuleiro de ataque: desenhar os disparos feitos no tabuleiro do oponente e as mensagens de resultado em português | Divisão da 004 | 003, 004 | 2 | Concluída |
| 008 | Frota configurável pelo jogador | TODO do README | 002, 004 | 2 | Concluída |
| 006 | Modo solo: atacar um tabuleiro gerado aleatoriamente | TODO do README | 003, 005, 009 | 3 | Concluída |
| 007 | Jogador vs. máquina: jogar com a frota sorteada contra jogadas aleatórias, com turno extra após acerto | TODO do README | 006 | 4 | Concluída |
| 012 | Posicionar a própria frota: o jogador escolhe onde fica cada navio antes da partida contra a máquina | Divisão da 007 | 007 | 5 | Concluída |
| 010 | Cores no terminal: água, navios e disparos em cores diferentes | Fora de escopo da 004 | 009 | 6 | Concluída |
| 011 | `make run` repassar opções ao jogo (`make run ARGS="--semente 42"`) | Encerramento da 006 | 006 | 6 | Concluída |
| 013 | Máquina com estratégia: atirar em volta de um acerto até afundar o navio | Encerramento da 007 | 007 | 6 | Concluída |
| 014 | Posicionamento mais fácil: desfazer só o último navio, e avisar quando o resto da frota não cabe mais | Encerramento da 012 | 012 | 7 | Concluída |
| 015 | Mais cores: mensagens coloridas (`Acerto!`, `Água.`, vitória e derrota), e cores desligadas sozinhas quando a saída não é um terminal ou com a variável `NO_COLOR` | Encerramento da 010 | 010 | 7 | Concluída |
| 017 | `make run` devolver o código de saída do jogo (compilar e rodar o binário, no lugar do `go run`) | Encerramento da 011 | 011 | — | Backlog |
| 018 | Terceiro nível da máquina, com procura mais forte (xadrez ou onde cada navio ainda cabe) | Encerramento da 013 | 013 | 7 | Concluída |
| 019 | Busca do aviso de posicionamento mais rápida: o Submarino tenta uma orientação só, e frotas grandes do `--frota` não esperam cerca de 1 s antes de cada tabuleiro | Encerramento da 014 | 014 | — | Backlog |
| 020 | Nível `mestre` mais forte: contar as posições de todos os navios juntos, e não de cada navio sozinho | Encerramento da 018 | 018 | — | Backlog |
| 021 | Ajustes de cores: sem cores com `TERM=dumb`, e cores nas mensagens de erro do `stderr` | Encerramento da 015 | 015 | — | Backlog |

## Notas das specs concluídas

Pendências que specs anteriores deixaram para as próximas:

- **Issue #25 (013):** a máquina `dificil` não atira nas casas em volta de um navio afundado. A decisão é provisória: o jogador quer testar jogando. A regra tira cerca de 24 tiros da média (59,7 com a regra, 83,8 sem), e fica numa função separada (`aroundSunk`) para sair fácil.

Decisões para quem joga, sem spec ainda:

- **Aviso de posicionamento (014):** com a frota padrão ele nunca aparece, porque os navios que faltam são sempre os menores e sempre cabem. Só aparece com frotas do `--frota`. Decidir se continua.
- **Nível padrão (018):** decidir, jogando contra a `mestre`, se ela vira o padrão no lugar da `dificil`. Médias em 200 partidas: `mestre` 54,7 tiros, `dificil` 59,7, `facil` 95,6.
- **Cores num terminal de verdade (015):** a conferência foi feita num terminal simulado. Falta ver como ficam, inclusive o amarelo em fundo claro.

## Questões em aberto do jogo

Decisões que o README não cobre e que alguma spec vai precisar fechar:

Nenhuma no momento. A 009 decidiu que o jogador é avisado de qual navio afundou, e a 008 decidiu que a frota configurável só precisa caber no tabuleiro e ter pelo menos um navio.
