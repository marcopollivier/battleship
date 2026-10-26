# 006 — Modo solo

- **Status:** Rascunho
- **Roadmap:** #006 (onda 3)
- **Depende de:** 003 (disparo), 005 (posicionamento aleatório), 008 (frota configurável), 009 (tabuleiro de ataque)

## Problema

O primeiro item do TODO do README é deixar uma pessoa atacar um tabuleiro armado aleatoriamente. Todas as peças já existem: o domínio sorteia um tabuleiro (005) e registra os disparos (003), o terminal desenha o tabuleiro de ataque e as mensagens (009), e o jogador escolhe a frota com `--frota` (008). Falta juntar as peças numa partida: ler os disparos do teclado até a frota inteira afundar. Esta é a primeira versão jogável do jogo.

## Regras

**Início da partida:**

- Ao rodar o programa, a partida começa direto, sem menu.
- O computador sorteia um tabuleiro com a frota escolhida: a frota padrão, ou a do `--frota` (008).
- O sorteio usa o relógio como semente. Com a opção `--semente N` (um número inteiro), o jogador repete exatamente o mesmo tabuleiro.
- O programa mostra o tabuleiro de ataque (009), todo com `.`, e pede o primeiro disparo.

**Cada disparo:**

- O programa mostra o pedido `Disparo: ` e lê uma linha do teclado.
- A linha é uma coordenada, como `B7` (regras da 001: maiúscula ou minúscula, espaços em volta ignorados).
- Com um disparo válido, o programa mostra a mensagem de resultado (009) e, em seguida, o tabuleiro de ataque atualizado.
- Com uma entrada inválida, o programa mostra uma mensagem de erro, não conta o disparo e pede de novo:

  | Entrada | Mensagem |
  | ------- | -------- |
  | Fora do formato (`7B`, `xyz`, vazio) | `Coordenada inválida: use uma letra de A a J e um número de 1 a 10, por exemplo B7.` |
  | Fora do tabuleiro (`K1`, `A11`) | `Fora do tabuleiro: use uma letra de A a J e um número de 1 a 10.` |
  | Casa já atingida | `Você já atirou em B7.`, com a coordenada na forma canônica |

**Fim da partida:**

- Quando o último navio afunda, o programa mostra a mensagem do afundamento, o tabuleiro de ataque final e `Você afundou toda a frota em N disparos!`, onde N conta só os disparos válidos. O programa termina com código 0.
- O jogador desiste digitando `sair` (maiúscula ou minúscula) ou fechando a entrada (Ctrl+D). O programa mostra `Você desistiu depois de N disparos. A frota era:`, o tabuleiro próprio do computador (004/009, com navios e disparos) e termina com código 0.
- Ao fim da partida, por vitória ou desistência, o programa mostra `Semente: N`, para que o jogador possa repetir o tabuleiro com `--semente N`.

**Opções de linha de comando:**

- `--frota` continua como na 008. Uma frota recusada termina o programa antes da partida, como hoje.
- `--semente` com um valor que não é inteiro termina o programa com código 2, como qualquer opção inválida (008).

## Critérios de aceitação

Nos critérios abaixo, "uma partida com a semente S" é uma partida com a frota padrão e o tabuleiro sorteado com a semente S. Os testes leem as posições dos navios do tabuleiro sorteado, sem depender de uma semente específica.

- **CA-1:** Dado o programa iniciado, quando a partida começa, então a saída mostra o tabuleiro de ataque com as 100 casas como `.` e o pedido `Disparo: `.
- **CA-2:** Dada uma partida com a semente S, quando o jogador dispara em uma água, em uma casa de navio que não o afunda e na última casa de um navio, então o programa mostra `Água.`, `Acerto!` e `Afundou o <navio>!`, cada um seguido do tabuleiro de ataque atualizado.
- **CA-3:** Dada uma partida, quando o jogador digita `7B`, `K1` e uma casa em que já atirou, então o programa mostra a mensagem de cada caso da tabela, pede o disparo de novo, e nenhuma das três entradas conta como disparo.
- **CA-4:** Dada uma partida com a semente S, quando o jogador dispara em todas as casas dos navios, então o programa mostra o tabuleiro final, `Você afundou toda a frota em N disparos!` com N igual ao número de disparos válidos, e `Semente: S`, e termina com código 0.
- **CA-5:** Dada uma partida, quando o jogador digita `sair` ou `SAIR`, ou a entrada termina, então o programa mostra `Você desistiu depois de N disparos. A frota era:`, o tabuleiro próprio do computador e `Semente: S`, e termina com código 0.
- **CA-6:** Dadas duas partidas com `--semente 42`, quando elas começam, então os dois tabuleiros sorteados são iguais.
- **CA-7:** Dado o programa iniciado com `--frota S=2`, quando o jogador afunda os 2 Submarinos, então a partida termina com vitória.
- **CA-8:** Dado o programa iniciado com `--semente abc`, quando ele roda, então termina com código 2 sem começar a partida.

## Fora de escopo

- Turno extra depois de um acerto, que só faz sentido com dois jogadores (#007).
- O jogador posicionar a própria frota e a máquina atirar de volta (#007).
- Limpar a tela entre os disparos: cada tabuleiro é impresso abaixo do anterior.
- Cores (#010).
- Limite de disparos ou de tempo.
- Placar ou histórico entre partidas.

## Questões em aberto

As respostas propostas abaixo estão nas regras e nos critérios. Se alguma mudar, a spec muda junto.

- [ ] **1. Partida direto ao rodar:** o programa começa a partida sem menu, substituindo a tela atual do tabuleiro vazio. Quando a #007 chegar, ela decide como escolher entre os modos.
- [ ] **2. Semente:** relógio por padrão, `--semente N` para repetir, e a semente aparece no fim da partida (nota da 005 no roadmap).
- [ ] **3. Tela:** o tabuleiro de ataque é reimpresso depois de cada disparo válido, sem limpar a tela. Entradas inválidas só mostram a mensagem.
- [ ] **4. Desistência:** `sair` ou Ctrl+D, revelando a frota do computador. Código de saída 0 nos dois casos.
- [ ] **5. Mensagens:** os textos das tabelas de erro e de fim de partida.
- [ ] **6. Contagem:** N conta só os disparos válidos (água, acerto e afundado).

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
