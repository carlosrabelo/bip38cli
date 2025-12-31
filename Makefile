MAKEFLAGS += --no-print-directory
CORE_DIR := core

.DEFAULT_GOAL := help

.PHONY: build clean fmt help install lint test uninstall

build: ## Build bip38cli binary with version metadata
	@$(MAKE) -C $(CORE_DIR) build

clean: ## Remove build outputs and go caches
	@$(MAKE) -C $(CORE_DIR) clean

fmt: ## Run gofmt across core/
	@$(MAKE) -C $(CORE_DIR) fmt

help: ## Show available targets
	@echo "BIP38CLI - Main targets (delegates to core/Makefile)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*## "} {printf "  %-15s %s\n", $$1, $$2}'
	@echo ""
	@echo "For more targets, run 'make -C core help'"

install: build ## Install binary via scripts/install.sh
	@./scripts/install.sh --user

lint: ## Run golangci-lint when available
	@$(MAKE) -C $(CORE_DIR) lint

test: ## Execute go test ./...
	@$(MAKE) -C $(CORE_DIR) test

uninstall: ## Remove installed binary via scripts/uninstall.sh
	@./scripts/uninstall.sh --user
