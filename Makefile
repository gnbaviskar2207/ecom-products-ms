.PHONY: help tools gen-all buf-gen gov-gen lint test test-race build run start-local buf-login buf-dep

.DEFAULT_GOAL := help

# Colors for terminal output
CYAN   := \033[0;36m
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RESET  := \033[0m

GOBIN := $(shell go env GOPATH)/bin

help: ## Show this help message
	@printf "$(CYAN)Usage:$(RESET) make [target]\n\n"
	@printf "$(CYAN)Available Targets:$(RESET)\n"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

tools: ## Install development tools (dlv, goverter, gopls)
	@printf "$(CYAN)==> Installing development tools...$(RESET)\n"
	go install github.com/go-delve/delve/cmd/dlv@v1.24.2
	go install github.com/jmattheis/goverter/cmd/goverter@v1.9.3
	go install golang.org/x/tools/gopls@v0.20.0
	@printf "$(GREEN)✓ Development tools installed.$(RESET)\n"

gen-all: buf-gen gov-gen ## Run all code generators (buf + goverter)
	@printf "$(GREEN)✓ All code generation completed successfully.$(RESET)\n"

buf-gen: remove-stale-buf-gen-files buf-generate ## Clean and generate protobuf/gRPC files

gov-gen: remove-stale-goverter-gen-files goverter-gen ## Clean and generate converter code

buf-generate:
	@printf "$(CYAN)==> Generating protobuf and gRPC code with buf...$(RESET)\n"
	PATH="$(GOBIN):$$PATH" buf generate
	@printf "$(GREEN)✓ Buf generation completed.$(RESET)\n"

goverter-gen:
	@printf "$(CYAN)==> Generating converters with goverter...$(RESET)\n"
	PATH="$(GOBIN):$$PATH" goverter gen ./internal/transform
	@printf "$(GREEN)✓ Goverter generation completed.$(RESET)\n"

remove-stale-buf-gen-files:
	@printf "$(YELLOW)==> Cleaning stale buf generated files...$(RESET)\n"
	@find ./gen -mindepth 1 -delete 2>/dev/null || true
	@printf "$(GREEN)✓ Stale buf files removed.$(RESET)\n"

remove-stale-goverter-gen-files:
	@printf "$(YELLOW)==> Cleaning stale goverter generated files...$(RESET)\n"
	@find ./internal/transform/generated -mindepth 1 -delete 2>/dev/null || true
	@printf "$(GREEN)✓ Stale goverter files removed.$(RESET)\n"

remove-binaries:
	@printf "$(YELLOW)==> Cleaning stale application binaries...$(RESET)\n"
	@find ./bin -mindepth 1 -delete 2>/dev/null || true
	@printf "$(GREEN)✓ Stale application binaries removed.$(RESET)\n"

lint: ## Run linters (buf lint + go vet)
	@printf "$(CYAN)==> Running buf lint...$(RESET)\n"
	buf lint
	@printf "$(CYAN)==> Running go vet...$(RESET)\n"
	go vet ./...
	@printf "$(GREEN)✓ All linters passed.$(RESET)\n"

test: ## Run unit tests with coverage
	@printf "$(CYAN)==> Running tests with coverage...$(RESET)\n"
	go test -cover ./...
	@printf "$(GREEN)✓ Tests completed successfully.$(RESET)\n"

test-race: ## Run unit tests with race detector
	@printf "$(CYAN)==> Running tests with race detector...$(RESET)\n"
	go test -race ./...
	@printf "$(GREEN)✓ Race detector tests completed successfully.$(RESET)\n"

build: ## Build product microservice binary
	@printf "$(CYAN)==> Building binary (bin/products)...$(RESET)\n"
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/products ./cmd/products
	@printf "$(GREEN)✓ Binary built successfully at bin/products$(RESET)\n"

run: ## Run product microservice locally
	@printf "$(CYAN)==> Starting product microservice...$(RESET)\n"
	go run ./cmd/products/main.go --config=./config.yaml

run-binary: ## Run the product service using binary
	@printf "$(CYAN)==> Starting product microservice...$(RESET)\n"
	./bin/products --config=./config.yaml

start-local: ## Start the application locally
	@printf "$(CYAN)==> Starting application locally...$(RESET)\n"
	make remove-binaries
	make gen-all
	make lint
	make test
	make build
	make run-binary
	



















// other command:
buf-login: ## Login to buf to exit the rate limiting
	@printf "$(CYAN)==> Logging in to buf registry...$(RESET)\n"
	buf registry login
	@printf "$(GREEN)✓ Logged in to buf registry successfully.$(RESET)\n"

buf-dep: ## Install buf plugins like validation etc
	@printf "$(CYAN)==> Installing buf plugins...$(RESET)\n"
	buf dep update
	@printf "$(GREEN)✓ Successfully installed buf plugins.$(RESET)\n"
