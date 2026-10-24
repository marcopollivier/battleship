# 008 — Frota configurável

- **Status:** Rascunho
- **Roadmap:** #008 (onda 2)
- **Depende de:** 002 (tabuleiro e frota), 004 (terminal), 005 (posicionamento aleatório, para saber se uma frota cabe)

## Problema

O terceiro item do TODO do README é deixar o jogador definir quantos navios quer usar. O domínio já aceita qualquer frota (decisão da 002), mas não existe forma de o jogador informar uma frota nem uma validação que impeça frotas sem sentido (vazia, negativa ou que não cabe no tabuleiro). Esta spec define como o jogador descreve a frota e quais frotas são aceitas. Usar a frota numa partida fica para a #006 e a #007.

## Regras

- O jogador só escolhe **quantos** navios de cada um dos 5 tipos existentes quer. Tipos novos e formatos novos ficam fora.
- O jogador descreve a frota num texto, no formato definido pela questão em aberto 1.
- Uma frota é aceita quando:
  - nenhuma quantidade é negativa;
  - há pelo menos um navio no total;
  - respeita os limites da questão em aberto 2;
  - cabe no tabuleiro, respeitando as regras da 002 (sem sobreposição e sem encostar).
- Uma frota recusada gera um erro com o motivo, e o jogador recebe uma mensagem em português.
- Sem configuração, vale a frota padrão da cartela.

## Critérios de aceitação

- **CA-1:** Dado um texto que descreve 2 Submarinos e 1 Porta-Aviões, quando o jogo o interpreta, então a frota tem exatamente esses navios e nenhum outro.
- **CA-2:** Dado um texto vazio (sem configuração), quando o jogo interpreta a frota, então o resultado é a frota padrão.
- **CA-3:** Dado um texto mal formado (tipo desconhecido, quantidade que não é número, separador errado), quando o jogo o interpreta, então retorna um erro de formato.
- **CA-4:** Dada uma frota com quantidade negativa, ou com zero navios no total, quando o jogo a valida, então ela é recusada.
- **CA-5:** Dada uma frota que não cabe no tabuleiro (por exemplo, 30 Encouraçados), quando o jogo a valida, então ela é recusada com um erro de frota que não cabe.
- **CA-6:** Dada uma frota aceita, quando o jogo gera um tabuleiro aleatório com ela (005), então o tabuleiro é completo.
- **CA-7:** Dado cada motivo de recusa, quando o jogo monta a mensagem para o jogador, então ela está em português e explica o motivo.

As decisões das questões em aberto podem acrescentar ou ajustar critérios.

## Fora de escopo

- Usar a frota configurada numa partida (#006 e #007).
- Criar tipos ou formatos de navio novos.
- Mudar o tamanho do tabuleiro.

## Questões em aberto

- [ ] **1. Como o jogador informa a frota:** por uma opção de linha de comando, num texto como `--frota S=2,D=1,C=0,E=1,P=1` (usando as letras da legenda da 004), ou por perguntas no terminal antes da partida ("Quantos Submarinos?")? A opção de linha de comando é mais simples e fácil de testar; as perguntas são mais amigáveis.
- [ ] **2. Limites:** além de "caber no tabuleiro", há um máximo por tipo ou um máximo de casas ocupadas? Por exemplo, no máximo 30 casas ocupadas (30% do tabuleiro), para o jogo não virar só acerto.
- [ ] **3. Tipos omitidos:** na configuração por texto, um tipo que não aparece vale 0 ou vale a quantidade da frota padrão? Por exemplo, `S=5` significa "só 5 Submarinos" ou "a frota padrão, mas com 5 Submarinos"?
- [ ] **4. Como saber se cabe:** é aceitável decidir "cabe" tentando gerar um tabuleiro aleatório (005)? Nesse caso, uma frota no limite do possível pode ser recusada mesmo cabendo, se o sorteio não achar uma disposição nas tentativas permitidas.

## Encerramento

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
