# 008 — Frota configurável

- **Status:** Concluída
- **Roadmap:** #008 (onda 2)
- **Depende de:** 002 (tabuleiro e frota), 004 (terminal), 005 (posicionamento aleatório, para saber se uma frota cabe)

## Problema

O terceiro item do TODO do README é deixar o jogador definir quantos navios quer usar. O domínio já aceita qualquer frota (decisão da 002), mas não existe forma de o jogador informar uma frota nem uma validação que impeça frotas sem sentido (vazia, negativa ou que não cabe no tabuleiro). Esta spec define como o jogador descreve a frota e quais frotas são aceitas. Usar a frota numa partida fica para a #006 e a #007.

## Regras

- O jogador só escolhe **quantos** navios de cada um dos 5 tipos existentes quer. Tipos novos e formatos novos ficam fora.
- O jogador informa a frota pela opção de linha de comando `--frota`, num texto com pares `letra=quantidade` separados por vírgula, usando as letras da legenda da 004 (`S`, `D`, `C`, `E`, `P`). Por exemplo: `--frota S=2,D=1,P=1`.
  - As letras aceitam maiúscula ou minúscula, e espaços em volta dos pares são ignorados.
  - Um tipo que não aparece no texto vale 0. `--frota S=5` significa "só 5 Submarinos".
  - Uma letra repetida (`S=1,S=2`) é erro de formato.
- Sem a opção `--frota`, vale a frota padrão da cartela.
- Uma frota é aceita quando:
  - nenhuma quantidade é negativa;
  - há pelo menos um navio no total;
  - cabe no tabuleiro, respeitando as regras da 002 (sem sobreposição e sem encostar).
- Não há outros limites (por tipo ou de casas ocupadas).
- "Caber" é decidido tentando gerar um tabuleiro aleatório com a frota (005), com uma semente fixa, para que a mesma frota tenha sempre a mesma resposta. Uma frota no limite do possível pode ser recusada mesmo cabendo, e a mensagem diz isso ("não foi possível posicionar essa frota no tabuleiro").
- Uma frota recusada gera um erro com o motivo. O programa mostra a mensagem em português e termina com erro.
- Com uma frota aceita, o programa mostra o tabuleiro como hoje. Usar a frota numa partida fica para a #006 e a #007.

## Critérios de aceitação

- **CA-1:** Dado o texto `S=2,P=1`, quando o jogo o interpreta, então a frota tem 2 Submarinos e 1 Porta-Aviões, e nenhum outro navio.
- **CA-2:** Dado o texto ` s=2 , p=1 ` (minúsculas e espaços), quando o jogo o interpreta, então o resultado é o mesmo do CA-1.
- **CA-3:** Dado que a opção `--frota` não foi informada, quando o programa monta a frota, então ela é a frota padrão.
- **CA-4:** Dado um texto mal formado, quando o jogo o interpreta, então retorna um erro de formato. Exemplos: vazio, `X=1` (tipo desconhecido), `S=a`, `S2`, `S=1;D=1`, `S=1,S=2`.
- **CA-5:** Dada uma frota com quantidade negativa (`S=-1`), quando o jogo a valida, então ela é recusada com um erro de quantidade negativa.
- **CA-6:** Dada uma frota com zero navios no total (`S=0`), quando o jogo a valida, então ela é recusada com um erro de frota vazia.
- **CA-7:** Dada uma frota que não cabe no tabuleiro (`E=30`), quando o jogo a valida, então ela é recusada com um erro de frota que não cabe.
- **CA-8:** Dada uma frota aceita, quando o jogo gera um tabuleiro aleatório com ela (005), então o tabuleiro é completo.
- **CA-9:** Dado cada motivo de recusa (formato, quantidade negativa, frota vazia e frota que não cabe), quando o jogo monta a mensagem para o jogador, então ela está em português e explica o motivo.
- **CA-10:** Dado o programa executado com `--frota E=30`, quando ele roda, então mostra a mensagem de frota que não cabe e termina com código de saída diferente de 0. Com `--frota S=2,P=1`, mostra o tabuleiro e termina com código 0.

## Fora de escopo

- Usar a frota configurada numa partida (#006 e #007).
- Criar tipos ou formatos de navio novos.
- Mudar o tamanho do tabuleiro.

## Questões em aberto

Todas decididas em 2026-10-24:

- [x] **1. Como o jogador informa a frota:** opção de linha de comando `--frota`, com as letras da legenda (CA-1 a CA-4, CA-10).
- [x] **2. Limites:** nenhum além de caber no tabuleiro e ter pelo menos um navio (CA-5 a CA-7).
- [x] **3. Tipos omitidos:** valem 0 (CA-1).
- [x] **4. Como saber se cabe:** tentando gerar um tabuleiro aleatório com semente fixa, aceitando que uma frota no limite possa ser recusada (CA-7).

## Encerramento

Concluída em 2026-10-24, nas tarefas T1–T6. Os 10 critérios de aceitação são cobertos por testes em `internal/board/fleet_test.go` (CA-5 a CA-8), `internal/terminal/fleet_test.go` (CA-1, CA-2, CA-4 e CA-9) e `cmd/init_test.go` (CA-3 e CA-10). A cobertura ficou em 99,4% em `internal/board` e 100% em `internal/terminal` (100% em `fleet.go` nos dois pacotes) e 88,9% em `cmd`. No `cmd`, ficaram sem teste só o `-h` e a falha de escrita do `RenderBoard`.

Executando o programa:

- `go run .` e `go run . --frota S=2,P=1` mostram o tabuleiro e terminam com código 0.
- `go run . --frota E=30` mostra `Não foi possível posicionar essa frota no tabuleiro.` em `stderr` e termina com código 1.
- `go run . --frota X=1` mostra a mensagem de formato inválido em `stderr` e termina com código 1.

**Desvios em relação ao plano:**

- O `ValidateFleet` precisou de um `//nolint:gosec` na linha que cria o `*rand.Rand` com a semente fixa. O `gosec` acusou o G404 (gerador fraco) em `fleet.go`, embora não tenha acusado o `random.go` na 005. A semente fixa é a decisão da spec e a checagem não tem relação com segurança.
- Na T3, o `ParseFleet` já recusava o par sem `=`, a letra desconhecida e o número inválido, com erros comuns. A T4 trocou esses erros pelo `ErrInvalidFleetFormat` e acrescentou a recusa da letra repetida. O texto vazio cai no caso do par sem `=`.
- O `ParseFleet` procura a letra percorrendo o `shipSymbols`, num helper `shipTypeBySymbol`, em vez de montar um mapa inverso. As mensagens do `FleetErrorMessage` ficam numa tabela (`fleetErrorMessages`) percorrida com `errors.Is`.
- O `cmd.Run` trata o `-h` com código 0, como o `flag.ExitOnError` faz. O plano só falava do código 2 para opção inválida. Um teste novo (`TestRun_UnknownOption`) confere o código 2 e o uso da opção em `stderr`.
- Uma falha ao escrever o tabuleiro também sai em `stderr`, pelo `FleetErrorMessage` (`Erro: ...`), com código 1.
- Os testes do `cmd` usam os helpers `fleetRender` e `fleetRun`. Além dos casos do plano, cobrem `--frota=S=2,P=1`, `--frota ""`, `S=-1` e `S=0`.

**Para o roadmap:** a #006 e a #007 recebem a frota validada pelo `cmd.Run` (hoje só usada para montar o tabuleiro vazio) e devem passá-la ao `RandomBoard` e ao laço da partida. O `README.md` e o `CLAUDE.md` podem passar a citar a opção `--frota` e a nova assinatura de `cmd.Run` (o `cmd.Init()` saiu).
