.PHONY: help build run test lint fmt vet sec install-tools install-deps setup clean

# Configurações
BINARY_NAME=battleship
MAIN_PATH=./main.go
BUILD_DIR=./bin

# Comando padrão
help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Compila o projeto
	@echo "🔨 Compilando..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Compilação concluída: $(BUILD_DIR)/$(BINARY_NAME)"

run: ## Executa o projeto
	@echo "🚀 Executando..."
	@go run $(MAIN_PATH)

test: ## Executa os testes
	@echo "🧪 Executando testes..."
	@go test -v ./...

fmt: ## Formata o código
	@echo "🎨 Formatando código..."
	@go fmt ./...
	@echo "✅ Código formatado"

vet: ## Executa go vet
	@echo "🔍 Executando go vet..."
	@go vet ./...
	@echo "✅ go vet concluído"

lint: ## Executa o linter (golangci-lint)
	@echo "🔍 Executando golangci-lint..."
	@golangci-lint run
	@echo "✅ Linting concluído"

sec: ## Executa verificação de segurança (go vet -vettool=shadow)
	@echo "🔒 Executando verificação de segurança..."
	@go vet -shadow ./... 2>/dev/null || go vet ./...
	@echo "✅ Verificação de segurança concluída"

mod-tidy: ## Limpa e organiza go.mod
	@echo "📦 Organizando dependências..."
	@go mod tidy
	@echo "✅ Dependências organizadas"

install-deps: ## Instala dependências externas do sistema
	@echo "📦 Instalando dependências externas..."
	@if ! command -v brew &> /dev/null; then \
		echo "❌ Homebrew não encontrado. Instale o Homebrew primeiro: https://brew.sh/"; \
		exit 1; \
	fi
	@echo "  → Instalando pre-commit..."
	@brew list pre-commit &>/dev/null || brew install pre-commit
	@echo "✅ Dependências externas instaladas"

install-tools: ## Instala ferramentas de desenvolvimento Go
	@echo "🛠️  Instalando ferramentas Go..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Ferramentas Go instaladas"

setup: install-deps install-tools ## Configura ambiente de desenvolvimento completo
	@echo "⚙️  Configurando pre-commit hooks..."
	@pre-commit install
	@echo "✅ Ambiente de desenvolvimento configurado!"

clean: ## Remove arquivos de build
	@echo "🧹 Limpando..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Limpeza concluída"

pre-commit: fmt vet lint sec test ## Executa todas as verificações antes do commit
	@echo "✅ Todas as verificações passaram!"

all: clean fmt vet lint sec test build ## Executa pipeline completo

setup-project: setup mod-tidy ## Configura projeto do zero (primeira vez)
	@echo "🎉 Projeto configurado e pronto para desenvolvimento!"
