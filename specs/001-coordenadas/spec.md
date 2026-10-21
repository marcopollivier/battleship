# 001 — Coordenadas

- **Status:** Aprovada
- **Roadmap:** #001
- **Depende de:** —

## Problema

Todo disparo e todo posicionamento de navio é comunicado por uma coordenada como `A1`, `B3` ou `J10`. O jogo precisa entender o que o jogador digitou, saber a qual casa do tabuleiro aquilo corresponde e conseguir exibir qualquer casa de volta no mesmo formato. Sem isso, nenhuma das próximas features (frota, disparo, terminal) consegue receber jogadas.

## Regras

- O tabuleiro tem 10 × 10 casas.
- Uma coordenada é uma letra de `A` a `J` seguida de um número de `1` a `10`, como nos exemplos do README (`A1`, `B3`, `J10`).
- A letra identifica a coluna e o número identifica a linha, como no tabuleiro que o jogo já exibe hoje: letras no topo e números na lateral. `A1` é o canto superior esquerdo e `J10` é o canto inferior direito.
- Cada uma das 100 casas tem exatamente uma coordenada e cada coordenada válida aponta para exatamente uma casa.
- Ao exibir uma casa, o jogo usa a forma canônica: letra maiúscula seguida do número sem zeros à esquerda.
- Letras minúsculas são aceitas e tratadas como maiúsculas (`a1` é `A1`).
- Espaços nas pontas da entrada são ignorados. Espaço entre a letra e o número não é aceito.
- Só existe uma ordem: letra e depois número. `1A` é rejeitado.
- O número não pode ter zeros à esquerda. `A01` é rejeitado.
- Uma entrada inválida é rejeitada com um erro, sem encerrar o jogo. Há dois tipos de erro, para que o jogador receba uma mensagem útil:
  - **Formato inválido:** o texto não é uma letra de `A` a `Z` seguida de um número sem zeros à esquerda (ex.: vazio, `A`, `1A`, `A 1`, `A01`, `@1`).
  - **Fora do tabuleiro:** o formato está certo, mas a letra passa de `J` ou o número está fora de `1`–`10` (ex.: `K1`, `A0`, `A11`).

## Critérios de aceitação

- **CA-1:** Dado o texto `A1`, quando o jogo interpreta a coordenada, então obtém a casa da primeira coluna e primeira linha.
- **CA-2:** Dado o texto `J10`, quando o jogo interpreta a coordenada, então obtém a casa da última coluna e última linha.
- **CA-3:** Dado o texto `C7`, quando o jogo interpreta a coordenada, então obtém a casa da terceira coluna e sétima linha.
- **CA-4:** Dada qualquer uma das 100 casas, quando o jogo a converte para texto e interpreta esse texto de volta, então obtém a mesma casa, e o texto está na forma canônica.
- **CA-5:** Dada uma letra fora de `A`–`J` (ex.: `K1`, `Z5`), quando o jogo interpreta a coordenada, então retorna um erro de fora do tabuleiro.
- **CA-6:** Dado um número fora de `1`–`10` (ex.: `A0`, `A11`, `B99`), quando o jogo interpreta a coordenada, então retorna um erro de fora do tabuleiro.
- **CA-7:** Dado um texto fora do formato letra + número (ex.: vazio, `A`, `1`, `AA1`, `A1B`, `A-1`, `@1`), quando o jogo interpreta a coordenada, então retorna um erro de formato inválido.
- **CA-8:** Dado o texto `a1` ou `j10`, quando o jogo interpreta a coordenada, então obtém a mesma casa que `A1` ou `J10`.
- **CA-9:** Dado o texto ` A1 ` (com espaços ou tabulações nas pontas), quando o jogo interpreta a coordenada, então obtém a mesma casa que `A1`.
- **CA-10:** Dado o texto `A 1`, `1A` ou `A01`, quando o jogo interpreta a coordenada, então retorna um erro de formato inválido.

## Fora de escopo

- Tabuleiros com tamanho diferente de 10 × 10.
- Ler a coordenada do teclado ou exibir mensagens de erro no terminal. Isso fica para a #004 e para os modos de jogo.
- Saber se a casa tem navio ou se já foi atingida. Isso fica para as #002 e #003.
- Coordenadas de intervalo, como `A1-A5` para posicionar navios. Isso fica para a #002, se for necessário.

## Questões em aberto

Todas decididas em 2026-10-21:

- [x] **Maiúsculas e minúsculas:** aceitas. `a1` é `A1` (CA-8).
- [x] **Espaços:** os das pontas são ignorados e o do meio é rejeitado (CA-9, CA-10).
- [x] **Ordem invertida:** `1A` é rejeitado (CA-10).
- [x] **Zeros à esquerda:** `A01` é rejeitado (CA-10).
- [x] **Tipos de erro:** "formato inválido" e "fora do tabuleiro" são erros distintos (CA-5, CA-6, CA-7).
- [x] **Orientação:** letras são colunas e números são linhas, como no tabuleiro já exibido.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
