# 009 — Terminal: tabuleiro de ataque

- **Status:** Concluída
- **Roadmap:** #009 (onda 2)
- **Depende de:** 003 (disparo), 004 (terminal, tabuleiro próprio)

## Problema

Para atacar, o jogador precisa ver o tabuleiro do oponente do jeito que ele conhece: onde já atirou, onde acertou água e onde acertou navio, sem ver os navios que ainda não encontrou. Também precisa ler o resultado de cada disparo em português. O domínio já registra os disparos (003) e o terminal já desenha o tabuleiro próprio (004); falta desenhar o tabuleiro de ataque e as mensagens. Isso é pré-requisito do modo solo (#006).

## Regras

**Tabuleiro de ataque** (o tabuleiro do oponente, visto por quem ataca):

- Tem o mesmo cabeçalho e as mesmas fileiras do tabuleiro próprio (004).
- Nunca mostra um navio que ainda não foi atingido.
- Cada casa aparece de acordo com o que o atacante sabe dela:

  | Casa | Símbolo |
  | ---- | ------- |
  | Ainda não atingida | `.` |
  | Água atingida | `o` |
  | Navio atingido, ainda não afundado | `X` |
  | Navio afundado | letra do tipo: `S`, `D`, `C`, `E`, `P` |

- Abaixo do tabuleiro, depois de uma linha em branco, aparece a legenda:

  ```
  Legenda: . não atingida, o água, X acerto, S D C E P navio afundado
  ```

**Tabuleiro próprio** (004), que agora também mostra os disparos recebidos:

- Uma casa de navio atingida aparece como `X`, no lugar da letra do navio.
- Uma casa de água atingida aparece como `o`, no lugar de `~`.
- A legenda do tabuleiro próprio passa a ser:

  ```
  Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões, X atingido, o tiro na água
  ```

**Mensagens de resultado**, em português:

| Resultado | Mensagem |
| --------- | -------- |
| Água | `Água.` |
| Acerto | `Acerto!` |
| Afundado | `Afundou o <navio>!`, com o nome do navio da 004 (por exemplo, `Afundou o Cruzador!`) |

O tipo do navio só aparece quando ele afunda (decisão da 003).

Os símbolos dos dois tabuleiros são só ASCII, como na 004. Nenhuma linha termina com espaço.

## Critérios de aceitação

Nos critérios abaixo, o "tabuleiro de exemplo" é o mesmo da spec 003: frota padrão completa, com o Porta-Aviões em `A1`, `B1`, `C1`, `B2`, `B3`; o Encouraçado em `E1`–`H1`; o Cruzador em `J1`–`J3`; os Destroyers em `A5`–`B5` e `E3`–`F3`; e os Submarinos em `H5`, `A8` e `E8`.

- **CA-1:** Dado o tabuleiro de exemplo sem disparos, quando o jogo desenha o tabuleiro de ataque, então todas as 100 casas aparecem como `.`, e nenhuma letra de navio aparece.
- **CA-2:** Dado o tabuleiro de exemplo com disparos em `J10` (água) e `E1` (acerto), quando o jogo desenha o tabuleiro de ataque, então `J10` aparece como `o`, `E1` como `X`, e as outras casas como `.`.
- **CA-3:** Dado o tabuleiro de exemplo com o Cruzador afundado (disparos em `J1`, `J2` e `J3`), quando o jogo desenha o tabuleiro de ataque, então essas três casas aparecem como `C`.
- **CA-4:** Dado cada resultado de disparo (água, acerto e afundado de cada um dos 5 tipos), quando o jogo monta a mensagem, então ela é o texto da tabela de mensagens, e só a de afundado cita o navio.
- **CA-5:** Dado qualquer tabuleiro de ataque, quando o jogo o desenha, então a saída é o tabuleiro, uma linha em branco e a legenda do tabuleiro de ataque. Nenhuma linha termina com espaço, e todos os caracteres acima da legenda são ASCII.
- **CA-6:** Dado o tabuleiro de exemplo com disparos em `J10` (água) e `E1` (acerto), quando o jogo desenha o tabuleiro próprio, então `J10` aparece como `o`, `E1` como `X`, as outras casas do Encouraçado continuam como `E`, e a legenda é a nova legenda do tabuleiro próprio.
- **CA-7:** Dado um tabuleiro sem disparos, quando o jogo desenha o tabuleiro próprio, então a saída é a mesma da 004, exceto pela legenda nova.

## Fora de escopo

- Ler o disparo do teclado e o laço da partida (#006).
- Mostrar os dois tabuleiros lado a lado (#006).
- Cores (#010).
- Mostrar no tabuleiro próprio quais navios foram afundados de forma diferente de só atingidos. No tabuleiro próprio, toda casa de navio atingida aparece como `X`.

## Questões em aberto

Todas decididas em 2026-10-24:

- [x] **1. Símbolos:** `.` não atingida, `o` água atingida, `X` acerto (CA-1, CA-2).
- [x] **2. Navio afundado:** mostra a letra do tipo (CA-3).
- [x] **3. Mensagens:** `Água.`, `Acerto!` e `Afundou o <navio>!` (CA-4).
- [x] **4. Disparos recebidos no tabuleiro próprio:** entram agora (CA-6, CA-7).

## Encerramento

Concluída em 2026-10-24, nas tarefas T1–T4. Os 7 critérios de aceitação são cobertos por testes em `internal/terminal/messages_test.go`, `internal/terminal/attack_test.go` e `internal/terminal/render_test.go`, com 100% de cobertura no `internal/terminal`. O domínio não mudou. O `go run .` imprime o tabuleiro vazio da 004 com a legenda nova do tabuleiro próprio.

**Desvios em relação ao plano:**

- Os símbolos de disparo viraram constantes internas em `render.go` (`hitSymbol` para `X` e `missSymbol` para `o`), usadas pelos dois tabuleiros, e o `.` virou `untouchedSymbol` em `attack.go`. O `shipSymbols` não mudou.
- A regra de navio afundado ficou numa função interna, `shipSunk`, em `attack.go`: o navio está afundado quando todas as suas casas estão como `CellHit`, como o plano previa.
- Os testes ganharam, além de `attackExampleBoard` e `attackFire`, os helpers `attackRender` e `attackAssertCells`, que confere todas as 100 casas (as esperadas e `.` no resto). Ficaram em `attack_test.go`, e o `TestRenderBoard_ReceivedShots` de `render_test.go` também os usa.
- O `TestShotMessage` também confere que um `Miss` com `Ship` preenchido não cita o navio. O `TestRenderAttackBoard_Sunk` ganhou dois casos além do Cruzador: um Cruzador só atingido (continua `X`) e um Submarino afundado ao lado de uma água e de um acerto.

**Para o roadmap:**

- O `RenderBoard` (tabuleiro próprio) e o `RenderAttackBoard` calculam a grade inteira consultando `Board.Cell` casa a casa, e cada consulta percorre os navios. Para 100 casas isso é irrelevante, mas, se o domínio ganhar um `Board.Sunk(ship)` ou uma consulta em lote, o `shipSunk` do terminal pode ser trocado por ela.
- O `CLAUDE.md` precisa citar o `RenderAttackBoard`, o `AttackLegend` e o `ShotMessage` no `internal/terminal`.
