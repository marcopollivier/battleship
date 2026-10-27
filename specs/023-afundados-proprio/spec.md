# 023 — Navios afundados no próprio tabuleiro

- **Status:** Aprovada
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

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
