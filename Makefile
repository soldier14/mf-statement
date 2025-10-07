# MF Statement Makefile

# Variables
BINARY_NAME=mf-statement
BUILD_DIR=bin

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show available commands
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  test          - Run all tests"
	@echo "  coverage      - Run tests with coverage"
	@echo "  lint          - Run linting (fmt + vet)"
	@echo "  clean         - Clean build artifacts"
	@echo "  run           - Build and run the application"
	@echo "  install-hooks - Install git hooks (pre-commit, pre-push, commit-msg)"
	@echo "  ci-check      - Run CI/CD checks locally"

.PHONY: build
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/statement
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	ginkgo run ./internal/...

.PHONY: coverage
coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@mkdir -p coverage
	go test -cover -coverprofile=coverage/coverage.out ./internal/... -covermode=count
	@grep -v "internal/util/logger.go" coverage/coverage.out > coverage/coverage_filtered.out
	go tool cover -func=coverage/coverage_filtered.out
	go tool cover -html=coverage/coverage_filtered.out -o coverage/coverage.html
	@echo "Coverage report: coverage/coverage.html"

.PHONY: lint
lint: ## Run linting (fmt + vet)
	@echo "Running linting..."
	go fmt ./...
	go vet ./...
	@echo "Linting complete"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	go clean
	rm -rf $(BUILD_DIR)
	rm -rf coverage
	rm -f coverage.out coverage.html

.PHONY: run
run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: install-hooks
install-hooks: ## Install git hooks
	@echo "Installing git hooks..."
	./scripts/install-hooks.sh

.PHONY: ci-check
ci-check: ## Run CI/CD checks locally
	@echo "Running CI/CD checks..."
	./scripts/ci-check.sh
