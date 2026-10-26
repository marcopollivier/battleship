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
| 007 | Jogador vs. máquina: jogar com a frota sorteada contra jogadas aleatórias, com turno extra após acerto | TODO do README | 006 | 4 | Spec |
| 012 | Posicionar a própria frota: o jogador escolhe onde fica cada navio antes da partida contra a máquina | Divisão da 007 | 007 | 5 | Backlog |
| 010 | Cores no terminal: água, navios e disparos em cores diferentes | Fora de escopo da 004 | 009 | — | Backlog |
| 011 | `make run` repassar opções ao jogo (`make run ARGS="--semente 42"`) | Encerramento da 006 | 006 | — | Backlog |

## Notas das specs concluídas

Pendências que specs anteriores deixaram para as próximas:

- **#012:** usar `board.ValidOrientations` para mostrar as orientações possíveis de cada navio ao jogador (005).
- **#007:** decidir como o jogador escolhe entre o modo solo e o jogo contra a máquina; hoje o programa começa direto no modo solo (006).

## Questões em aberto do jogo

Decisões que o README não cobre e que alguma spec vai precisar fechar:

Nenhuma no momento. A 009 decidiu que o jogador é avisado de qual navio afundou, e a 008 decidiu que a frota configurável só precisa caber no tabuleiro e ter pelo menos um navio.
