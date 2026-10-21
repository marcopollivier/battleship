# 001 — Coordenadas

- **Status:** Rascunho
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
- Uma entrada que não seja uma coordenada válida é rejeitada com um erro, sem encerrar o jogo.

## Critérios de aceitação

- **CA-1:** Dado o texto `A1`, quando o jogo interpreta a coordenada, então obtém a casa da primeira coluna e primeira linha.
- **CA-2:** Dado o texto `J10`, quando o jogo interpreta a coordenada, então obtém a casa da última coluna e última linha.
- **CA-3:** Dado o texto `C7`, quando o jogo interpreta a coordenada, então obtém a casa da terceira coluna e sétima linha.
- **CA-4:** Dada qualquer uma das 100 casas, quando o jogo a converte para texto e interpreta esse texto de volta, então obtém a mesma casa, e o texto está na forma canônica.
- **CA-5:** Dada uma letra fora de `A`–`J` (ex.: `K1`), quando o jogo interpreta a coordenada, então retorna um erro.
- **CA-6:** Dado um número fora de `1`–`10` (ex.: `A0`, `A11`), quando o jogo interpreta a coordenada, então retorna um erro.
- **CA-7:** Dado um texto fora do formato letra + número (ex.: vazio, `A`, `1`, `AA1`, `A1B`, `A-1`), quando o jogo interpreta a coordenada, então retorna um erro.

As decisões das questões em aberto podem acrescentar critérios.

## Fora de escopo

- Tabuleiros com tamanho diferente de 10 × 10.
- Ler a coordenada do teclado ou exibir mensagens de erro no terminal. Isso fica para a #004 e para os modos de jogo.
- Saber se a casa tem navio ou se já foi atingida. Isso fica para as #002 e #003.
- Coordenadas de intervalo, como `A1-A5` para posicionar navios. Isso fica para a #002, se for necessário.

## Questões em aberto

- [ ] **Maiúsculas e minúsculas:** `a1` deve ser aceito como `A1`?
- [ ] **Espaços:** espaços nas pontas (` A1 `) devem ser ignorados? E um espaço no meio (`A 1`)?
- [ ] **Ordem invertida:** `1A` deve ser aceito ou rejeitado?
- [ ] **Zeros à esquerda:** `A01` deve ser aceito como `A1` ou rejeitado?
- [ ] **Tipos de erro:** o jogo precisa distinguir "formato inválido" de "fora do tabuleiro", para mostrar mensagens diferentes ao jogador, ou basta um erro genérico?
- [ ] **Orientação:** o README chama tanto os números quanto as letras de "colunas" ("Colunas verticais marcadas de 1 a 10", "Colunas horizontais de A a J"). Esta spec segue o tabuleiro já exibido (letras = colunas, números = linhas). Está confirmado?

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
