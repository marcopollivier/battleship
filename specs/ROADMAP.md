# Roadmap

Backlog ordenado. Cada item vira uma pasta `specs/NNN-slug/` quando for especificado. Itens podem ser reordenados, divididos ou removidos a qualquer momento.

Status: `Backlog` → `Spec` → `Em andamento` → `Concluída`

| # | Feature | Origem | Status |
| - | ------- | ------ | ------ |
| 001 | Coordenadas: converter `A1`…`J10` em posição e vice-versa, rejeitar entradas inválidas | Regras do README | Concluída |
| 002 | Tabuleiro e frota: tipos de navio e formatos (linha e T), quantidades padrão, posicionamento com rotação, validação de limites e sobreposição | Regras do README | Em andamento |
| 003 | Disparo: resultado água/acerto/afundado, tiro repetido, fim de jogo quando toda a frota afunda | Regras do README | Backlog |
| 004 | Renderização no terminal: tabuleiro próprio e tabuleiro de ataque | `cmd/init.go` atual | Backlog |
| 005 | Posicionamento aleatório da frota | TODO do README | Backlog |
| 006 | Modo solo: atacar um tabuleiro gerado aleatoriamente | TODO do README | Backlog |
| 007 | Jogador vs. máquina: posicionar a própria frota e jogar contra jogadas aleatórias, com turno extra após acerto | TODO do README | Backlog |
| 008 | Frota configurável pelo jogador | TODO do README | Backlog |

## Questões em aberto do jogo

Decisões que o README não cobre e que alguma spec vai precisar fechar:

- O jogador é avisado de qual navio afundou ou só de que afundou?
- Na frota configurável, quais limites valem (tamanho máximo, ocupação máxima do tabuleiro)?
