# 024 — Tarefas

- **Plano:** [plan.md](plan.md)
- **Status:** Aprovado

Cada tarefa é um commit: teste falhando (ou a conferência à mão do plano) → implementação → `make pre-commit` → commit.

- [x] **T1** — `version`, `versionText` e a opção `--versao`. (cobre CA-4)
- [x] **T2** — `scripts/build-release.sh`, `make release-check`, `dist/` no `.gitignore` e o workflow `release.yml`; conferir os 5 binários e a versão com uma tag local de teste. (cobre CA-3 antes da tag, CA-5)
- [ ] **T3** — README ("Instalar" e como publicar) e `CLAUDE.md`; conferir o `go install .` local. (cobre CA-1, CA-2 antes da tag)
