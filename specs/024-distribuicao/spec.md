# 024 — Distribuição

- **Status:** Concluída
- **Roadmap:** #024 (onda 8), issue #37
- **Depende de:** nenhuma

## Problema

Para jogar hoje, a pessoa precisa ter Go instalado, clonar o repositório e rodar `make run`. Não há versão publicada nem binário pronto, então quem não programa em Go não consegue jogar. Esta spec torna o jogo fácil de instalar: com Go, um comando só; sem Go, baixar um arquivo da página de releases.

## Regras

**Com Go instalado:**

- O README ganha uma seção "Instalar", antes de "Como jogar", com:

  ```
  go install github.com/marcopollivier/battleship@latest
  battleship
  ```

- O binário instalado se chama `battleship` (questão 1).

**Sem Go (releases):**

- Cada versão nova do jogo é publicada na página de releases do GitHub, com binários prontos para (questão 2):
  - macOS (Intel e Apple Silicon);
  - Linux (amd64 e arm64);
  - Windows (amd64).
- A publicação é automática: quem mantém o projeto cria uma tag de versão (por exemplo `v1.0.0`) e o GitHub monta e publica os binários (questão 3).
- A seção "Instalar" do README explica como baixar o binário do sistema e rodar.
- A primeira versão é `v1.0.0`, porque o jogo do README está completo (questão 4).

**Versão do jogo:**

- A opção `--versao` mostra a versão do jogo (por exemplo `battleship v1.0.0`) e termina com código 0 (questão 5).
- Num binário instalado com `go install ...@v1.0.0` ou baixado da release, a versão é a da tag. Rodando do repositório (`make run`), aparece `battleship (desenvolvimento)`.

## Critérios de aceitação

- **CA-1:** Dado o README, quando a pessoa lê "Instalar", então encontra o comando `go install` e as instruções para baixar o binário da release.
- **CA-2:** Dado `go install github.com/marcopollivier/battleship@<versão>` num ambiente limpo, quando o comando termina, então o binário `battleship` existe e começa uma partida.
- **CA-3:** Dada uma tag `vX.Y.Z` enviada ao GitHub, quando a automação termina, então a release `vX.Y.Z` tem um binário para cada um dos 5 sistemas, e cada binário mostra a versão certa com `--versao`.
- **CA-4:** Dado `--versao`, quando o programa roda do repositório, então mostra `battleship (desenvolvimento)` e termina com código 0, sem começar partida.
- **CA-5:** Dado `make build`, quando o binário é gerado, então ele continua funcionando como hoje.

## Fora de escopo

- Homebrew, Scoop, pacotes `.deb` ou instaladores.
- Assinar os binários (macOS e Windows podem mostrar aviso de programa não verificado; o README explica como liberar).
- Atualização automática do jogo.
- Mudar o fluxo de PRs para gerar versões a cada merge.

## Questões em aberto

Todas decididas em 2026-10-27, como propostas:

- [x] **1. Nome do binário:** `battleship`, que é o que o `go install` gera a partir do nome do módulo. A alternativa, `batalha-naval`, exigiria mudar o módulo ou mover o `main.go` para `cmd/batalha-naval`.
- [x] **2. Sistemas:** os 5 da lista acima.
- [x] **3. Automação:** uma GitHub Action que roda na tag, compila com `go build` para cada sistema e publica com o `gh`, sem ferramentas de fora. A alternativa é o GoReleaser, mais completo, mas uma dependência a mais no CI.
- [x] **4. Primeira versão:** `v1.0.0`. A alternativa é `v0.1.0`, deixando a `1.0` para depois da usabilidade.
- [x] **5. Opção `--versao`:** sim, com o texto acima. A versão vem da tag gravada no binário pela própria ferramenta do Go.
- [x] **6. Quem cria a tag:** você, quando quiser publicar; o PR desta spec não cria a primeira tag, só deixa tudo pronto e explica como fazer.

## Encerramento

Concluída em 2026-10-27, nas tarefas T1–T3. O CA-4 é coberto por testes em `cmd/version_test.go` (`TestVersion`, com `v1.0.0`, `(devel)`, vazio e sem `BuildInfo`; `TestVersionText`; e `TestRun_Version`, que confere `--versao` sozinho e com `--modo xyz`: mostra `battleship (desenvolvimento)`, termina com código 0 e não escreve nada no `stderr`). Os testes falharam primeiro na compilação e depois contra funções provisórias que devolviam `""`, com o `--versao` ainda desconhecido (código 2). A cobertura do `cmd` ficou em 92,7%.

Conferências à mão, antes do merge:

- **CA-3 (binários):** `make release-check` gerou os 5 arquivos, e `file dist/*` mostrou:

  ```
  dist/battleship-dev-darwin-amd64:      Mach-O 64-bit executable x86_64
  dist/battleship-dev-darwin-arm64:      Mach-O 64-bit executable arm64
  dist/battleship-dev-linux-amd64:       ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, ..., stripped
  dist/battleship-dev-linux-arm64:       ELF 64-bit LSB executable, ARM aarch64, version 1 (SYSV), statically linked, ..., stripped
  dist/battleship-dev-windows-amd64.exe: PE32+ executable (console) x86-64, for MS Windows
  ```

  Cada um tem cerca de 2,1 MB. O script sem a versão mostra `uso: ./scripts/build-release.sh <versão>` e termina com código 2, e está no git com permissão de execução (`100755`).
- **CA-3 (versão da tag):** com o T2 commitado e a árvore limpa, `git tag v0.0.0-teste`, `make release-check` e `./dist/battleship-dev-darwin-arm64 --versao` mostraram `battleship v0.0.0-teste` (código 0). A tag foi apagada em seguida com `git tag -d`, e nunca foi enviada ao GitHub.
- **CA-5:** `make build` e `./bin/battleship --versao` mostraram `battleship (desenvolvimento)` (código 0); `./bin/battleship --modo solo --semente 42 --cores nao` com `sair` terminou a partida como antes, com `Semente: 42`. `go run . --versao` e `go run ./main.go --versao` também mostram `battleship (desenvolvimento)`.
- **CA-2 (antes da tag):** `GOBIN=<pasta vazia> go install .` gerou o `battleship`, que começou uma partida e mostrou a semente no fim.
- **Workflow:** o `actionlint` não está instalado nesta máquina, então o `release.yml` não foi conferido por ele. O YAML foi lido sem erro com o Ruby, e o script passou no `bash -n`.

**Só dá para confirmar depois da primeira tag:** que o workflow `Release` roda, que a release `v1.0.0` tem os 5 binários e que cada um mostra `battleship v1.0.0` com `--versao` (CA-3), e que `go install github.com/marcopollivier/battleship@v1.0.0` instala o jogo com essa versão (CA-2). Quem coordena confere e registra aqui, ou num PR `fix/` se algo falhar.

**Desvios em relação ao plano:**

- O plano dizia que, sem tag, a versão é `(devel)`. Isso vale para o `make run`, o `make build` e o `go run`, que não gravam a versão. Mas o `go build .` (o do script) e o `go install .` num checkout sem tag gravam uma pseudo-versão, como `v0.0.0-20261027043025-18a7f9efad5e`, com `+dirty` se houver mudanças não commitadas. Então `make release-check` fora de uma tag gera binários que mostram essa pseudo-versão, e não `(desenvolvimento)`. Isso não contraria a spec, que só fala do `make run`, e ajuda a saber de que commit veio o binário.
- O script compila o pacote (`go build .`), e não o `main.go`: com um arquivo, a ferramenta do Go não grava a versão. Por isso o `make build`, que continua com `./main.go`, mostra `(desenvolvimento)`.
- O README ganhou também a linha do `--versao` na tabela de opções, a pasta onde o `go install` coloca o binário, a tabela de arquivos por sistema e como apagar e recriar a tag se a publicação falhar.
- O `CLAUDE.md` ganhou uma seção "Releases".

**Para o roadmap:**

- Criar a tag `v1.0.0` e conferir a primeira release (CA-2 e CA-3).
- Instalar o `actionlint` no `make setup` (ou rodá-lo no CI), para conferir os workflows antes de uma tag.
- A fusão com a 022, a 023 e a 025, que também mexem no `cmd/init.go`, no README e no `CLAUDE.md`.
