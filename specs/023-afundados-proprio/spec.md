# 023 — Navios afundados no próprio tabuleiro

- **Status:** Concluída
- **Roadmap:** #023 (onda 8), issue #33
- **Depende de:** 015 (mais cores)

## Problema

No tabuleiro próprio (004), toda casa de navio atingida aparece como `X`, esteja o navio afundado ou ainda boiando. Para saber quais navios ainda tem, o jogador precisa lembrar das mensagens `Afundou o seu ...!`. No tabuleiro de ataque (009) é diferente: um navio inimigo afundado mostra a sua letra. A 009 deixou essa diferença para depois.

## Regras

- No tabuleiro próprio, as casas de um navio afundado aparecem com o símbolo `#` (questão 1). As casas atingidas de um navio ainda boiando continuam `X`, e as casas não atingidas continuam com a letra do navio.
- Com as cores ligadas (010 e 015), o `#` aparece em vermelho, a mesma cor dos navios afundados no tabuleiro de ataque (questão 2).
- A legenda do tabuleiro próprio ganha o símbolo novo:

  ```
  Legenda: ~ água, S Submarino, D Destroyer, C Cruzador, E Encouraçado, P Porta-Aviões, X atingido, # afundado, o tiro na água
  ```

- Isso vale em todo lugar onde o tabuleiro próprio é desenhado: a tela da partida contra a máquina, a derrota (a sua frota no fim) e a frota inimiga revelada no fim da partida, na desistência ou no modo solo.
- O tabuleiro de ataque não muda.

## Critérios de aceitação

- **CA-1:** Dado um tabuleiro com um navio afundado, um navio atingido e não afundado, e um navio intacto, quando ele é desenhado, então as casas do afundado são `#`, as atingidas do outro são `X`, e as do intacto têm a letra do navio.
- **CA-2:** Dado o mesmo tabuleiro com as cores ligadas, quando ele é desenhado, então o `#` aparece em vermelho, e tirando os códigos de cor o texto é igual ao desenho sem cores.
- **CA-3:** Dado qualquer tabuleiro próprio, quando ele é desenhado, então a legenda cita `# afundado`.
- **CA-4:** Dada uma partida contra a máquina em que a máquina afunda um navio do jogador, quando a tela volta para o jogador, então esse navio aparece com `#` em `Sua frota`.
- **CA-5:** Dada uma desistência ou o fim do modo solo, quando a frota inimiga é revelada, então os navios que o jogador afundou aparecem com `#`.

## Fora de escopo

- Mudar o tabuleiro de ataque.
- Uma lista de navios restantes ao lado do tabuleiro.
- Animação ou destaque do último navio afundado.

## Questões em aberto

Todas decididas em 2026-10-27, como propostas:

- [x] **1. Símbolo:** `#`, que se lê bem sem cor. Alternativas: a letra do navio em minúscula (`s d c e p`), que diz qual navio foi, mas se confunde com as maiúsculas sem cor; ou `*`.
- [x] **2. Cor:** vermelho (o `sunkColor` da 010), igual ao navio afundado do tabuleiro de ataque.
- [x] **3. Legenda:** `# afundado`, logo depois de `X atingido`.

## Encerramento

Concluída em 2026-10-27, nas tarefas T1 e T2. Os 5 critérios de aceitação são cobertos por testes em `internal/terminal`: CA-1 em `TestRenderBoard_Sunk` (`render_test.go`), CA-2 em `TestRenderBoard_SunkColored` (`color_test.go`), CA-3 em `TestRenderBoard_Sunk` e nos testes do `RenderBoard` que já existiam (`TestRenderBoard_Empty`, `TestRenderBoard_FullFleet`, `TestRenderBoard_ReceivedShots`), que passaram a conferir a legenda nova, CA-4 em `TestPlayVsMachine_SunkOwnShip` (`machine_test.go`) e CA-5 em `TestPlaySolo_GiveUpShowsSunk` (`solo_test.go`). A cobertura ficou em 99,7% em `internal/board`, 97,6% em `internal/terminal` e 93,0% em `cmd`.

O binário foi rodado com `--cores nao --semente 42`, `aleatorio` e um tiro em cada casa do tabuleiro: depois de `Afundou o seu Encouraçado!`, a tela seguinte mostra o Encouraçado como `# # # #` em `Sua frota`, com os outros navios ainda com as suas letras; na derrota, a sua frota aparece toda com `#`, e na frota inimiga revelada os navios que o jogador afundou também. Com `--modo solo --semente 7 --cores sim`, três linhas de tiros e `sair`, a frota revelada mostra o Submarino e o Destroyer afundados com `#` em vermelho (`\x1b[31m`) e o Porta-Aviões atingido, mas boiando, com `X` em vermelho negrito.

**Desvios em relação ao plano:**

- Os testes do T2 (CA-4 e CA-5) passaram de cara, porque o T1 já mudou o `RenderBoard`, que todas as telas usam. Para conferir que eles pegam o problema, eles foram rodados com o `sunkSymbol` trocado por `X` e falharam nas linhas esperadas.
- No CA-4, o teste usa uma frota de jogador com dois Submarinos e uma estratégia de tiros fixa (`machineScript`, da 013): a máquina afunda o Submarino em A1, erra em C3 e a tela volta para o jogador. Com um só navio, como nos tabuleiros de teste da 007, afundar o navio já termina a partida, e não haveria "a tela seguinte".
- Os testes do CA-4 e do CA-5 conferem as linhas da grade escritas à mão, e não montadas com o próprio `RenderBoard`, para não passarem com qualquer desenho.
- O README não cita a legenda do tabuleiro próprio, então não mudou. O `CLAUDE.md` cita o `#` na descrição do `render.go`.

**Para o roadmap:**

- A fusão com a 022, que refaz a tela e separa a grade da legenda no `render.go`: ela precisa manter o passo dos navios afundados e a legenda com `# afundado`.
