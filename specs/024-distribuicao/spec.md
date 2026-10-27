# 024 — Distribuição

- **Status:** Aprovada
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

Preenchido na fase 5: desvios em relação ao plano e o que ficou para o roadmap.
