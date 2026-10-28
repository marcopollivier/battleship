# 024 — Plano

- **Spec:** [spec.md](spec.md)
- **Status:** Aprovado

## Abordagem

1. **Versão (`cmd`):** a opção `--versao` lê a versão que a ferramenta do Go grava no binário (`runtime/debug.ReadBuildInfo`). O `go install ...@v1.0.0` grava a versão do módulo. A partir do Go 1.24 (o do `go.mod`), o `go build` dentro de um checkout numa tag também grava a versão da tag. Sem versão, `(devel)`, o jogo mostra `battleship (desenvolvimento)`.
2. **Releases (`.github/workflows/release.yml`, novo):** roda quando uma tag `v*` é enviada. Compila os 5 binários com `go build` e publica a release com o `gh`, que já vem no runner do GitHub.
3. **Documentação:** a seção "Instalar" no README, e o `CLAUDE.md` citando a versão e o workflow.

O jogo e o domínio não mudam.

## Pacotes e API

`cmd/version.go` (novo):

```go
// version returns the game version stamped in the binary, such as
// "v1.0.0", or "" when there is none (a development build).
func version(info *debug.BuildInfo, ok bool) string

// versionText is what --versao shows: "battleship v1.0.0" or
// "battleship (desenvolvimento)".
func versionText(v string) string
```

- `version` recebe o resultado do `debug.ReadBuildInfo()` por parâmetro, para os testes passarem um `BuildInfo` montado à mão. `(devel)` e vazio viram `""`.
- No `Run`: `showVersion := fs.Bool("versao", false, "mostra a versão do jogo e termina")`. Logo depois do `Parse`, se `*showVersion`, escreve `versionText(...)` no `stdout` e retorna 0, antes de validar as outras opções.

`.github/workflows/release.yml`:

```yaml
on:
  push:
    tags: ["v*"]
permissions:
  contents: write
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
      - uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
      - name: build
        run: ./scripts/build-release.sh "$GITHUB_REF_NAME"
      - name: release
        env:
          GH_TOKEN: ${{ github.token }}
        run: gh release create "$GITHUB_REF_NAME" dist/* --generate-notes
```

`scripts/build-release.sh` (novo): compila `battleship-<versão>-<sistema>-<arquitetura>` (com `.exe` no Windows) em `dist/`, com `CGO_ENABLED=0`, `-trimpath` e `-ldflags "-s -w"`, para:
- `darwin/amd64` e `darwin/arm64`;
- `linux/amd64` e `linux/arm64`;
- `windows/amd64`.

O mesmo script roda na máquina de quem mantém o projeto, para conferir antes de criar a tag.

`Makefile`: alvo `release-check` (`## Compila os binários de release em dist/`), que chama o script com a versão `dev`. `dist/` entra no `.gitignore`.

README, seção "Instalar", antes de "Como jogar":
- com Go: `go install github.com/marcopollivier/battleship@latest` e `battleship`;
- sem Go: baixar o arquivo do sistema na página de releases, e no macOS e no Linux dar permissão de execução (`chmod +x`);
- o aviso de programa não verificado no macOS e no Windows, e como liberar;
- `battleship --versao`;
- para quem mantém: como publicar uma versão (`git tag v1.0.0 && git push origin v1.0.0`).

## Decisões

| Decisão | Alternativas consideradas | Motivo |
| ------- | ------------------------- | ------ |
| `debug.ReadBuildInfo` | `-ldflags -X main.version=...` | Funciona igual no `go install` e no `go build` na tag, sem flag especial. O `-X` não funcionaria no `go install` de quem instala. |
| Script de build separado do workflow | Tudo no YAML, com uma matriz de jobs | O script roda também localmente (`make release-check`) e permite conferir os binários antes da tag. Um job só evita juntar arquivos de vários jobs. |
| `gh release create --generate-notes` | Escrever as notas à mão | As notas saem dos PRs mergeados, que já descrevem cada mudança em português. |
| `--versao` antes das outras validações | Validar tudo antes | Mostrar a versão não depende das outras opções, como o `-h`. |

## Estratégia de testes

| Critério | Teste ou conferência |
| -------- | -------------------- |
| CA-4 | `TestVersion` (tabela): `BuildInfo` com `v1.0.0`, com `(devel)`, vazio e sem `BuildInfo`. `TestRun_Version`: `--versao` mostra `battleship (desenvolvimento)` (o `go test` não grava versão), código 0 e nada além disso |
| CA-5 | `make build` e `./bin/battleship --versao`, conferidos à mão, e os testes atuais |
| CA-3 | **À mão, antes do merge:** `make release-check` gera os 5 arquivos, e `file dist/*` mostra o sistema e a arquitetura de cada um. **Depois do merge:** quando você criar a tag `v1.0.0`, a release precisa ter os 5 binários. Quem coordena confere e registra no encerramento |
| CA-2 | **À mão:** `GOBIN=$(mktemp -d) go install .` gera o `battleship`. O `go install ...@v1.0.0` real só é possível depois da tag e é conferido junto com o CA-3 |
| CA-1 | Revisão do README no PR |

## Dependências e riscos

- **O workflow só roda de verdade com uma tag.** Um erro no YAML só aparece na primeira publicação. Para reduzir o risco, o script é testado localmente, e o YAML é conferido com `actionlint`, se estiver disponível. Se a primeira release falhar, a correção vai num PR `fix/`, e a tag é recriada.
- **Versão no `go build` da tag:** depende do comportamento do Go 1.24. O teste local é `git tag v0.0.0-teste` num checkout limpo, depois `make release-check` e `--versao`, apagando a tag local em seguida.
- **Conflito com a 022 e a 025:** as três mexem no `Run`. A 024 entra primeiro na ordem combinada.

## Checagem da constituição

- [x] **Domínio separado da interface:** só o `cmd`, o CI e a documentação mudam.
- [x] **Aleatoriedade injetada:** não se aplica.
- [ ] **Testes antes do código:** exceção parcial, veja abaixo.
- [x] **Sem dependências externas:** `runtime/debug` da biblioteca padrão; o workflow usa o `gh` do runner, sem GoReleaser.
- [x] **Erros como valores:** não se aplica.
- [x] **Até ~8 tarefas:** 3.
- [x] **Idiomas:** código em inglês; README e texto da opção em português.

Exceções e justificativas:

- **Princípio 3 (teste primeiro), em parte:** a opção `--versao` tem testes escritos antes. O workflow, o script e o README não têm teste em Go. Eles são conferidos à mão, com os comandos acima, e o CA-3 só é confirmado na primeira tag, depois do merge.
