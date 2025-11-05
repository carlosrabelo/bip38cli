SRC_DIR := src

.DEFAULT_GOAL := help

.PHONY: build clean fmt help install lint test uninstall coverage

build: ## Build bip38cli binary with version metadata
	@$(MAKE) --no-print-directory -C $(SRC_DIR) build

clean: ## Remove build outputs and go caches
	@$(MAKE) --no-print-directory -C $(SRC_DIR) clean

fmt: ## Run gofmt across src/
	@$(MAKE) --no-print-directory -C $(SRC_DIR) fmt

help: ## Show available targets
	@echo "BIP38CLI - Main targets (delegates to src/Makefile)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*## "} {printf "  %-15s %s\n", $$1, $$2}'
	@echo ""
	@echo "For more targets, run 'make --no-print-directory -C src help'"

install: build ## Install binary via scripts/install.sh
	@./scripts/install.sh --user

lint: ## Run golangci-lint when available
	@$(MAKE) --no-print-directory -C $(SRC_DIR) lint

test: ## Execute go test ./...
	@$(MAKE) --no-print-directory -C $(SRC_DIR) test

uninstall: ## Remove installed binary via scripts/uninstall.sh
	@./scripts/uninstall.sh --user

coverage: ## Generate test coverage report in project root
	@$(MAKE) --no-print-directory -C $(SRC_DIR) coverage
