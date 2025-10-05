CORE_DIR       := core
CORE_MAKEFILE  := $(CORE_DIR)/Makefile

.DEFAULT_GOAL  := help

.PHONY: build clean fmt help install lint test uninstall

build: ## Build CLI binary inside core/
	@$(MAKE) -C $(CORE_DIR) build

clean: ## Clean build artifacts from core/
	@$(MAKE) -C $(CORE_DIR) clean

fmt: ## Format Go sources with gofmt
	@$(MAKE) -C $(CORE_DIR) fmt

help: ## Show root level targets
	@printf "BIP38CLI - Bitcoin Private Key Encryption Tool\n"
	@printf "==============================================\n\n"
	@printf " Build & Install:\n"
	@printf "   build           Build CLI binary inside core/\n"
	@printf "   install         Install CLI via scripts/install.sh\n"
	@printf "   uninstall       Uninstall CLI via scripts/uninstall.sh\n\n"
	@printf " Quality:\n"
	@printf "   fmt             Format Go sources with gofmt\n"
	@printf "   lint            Run golangci-lint via core/\n\n"
	@printf " Testing:\n"
	@printf "   test            Run go test ./... in core/\n\n"
	@printf " Utilities:\n"
	@printf "   clean           Clean build artifacts from core/\n"
	@printf "   help            Show this help\n"

lint: ## Run golangci-lint via core/
	@$(MAKE) -C $(CORE_DIR) lint

test: ## Run go test ./... in core/
	@$(MAKE) -C $(CORE_DIR) test

install: ## Install CLI via core Makefile
	@$(MAKE) -C $(CORE_DIR) install

uninstall: ## Uninstall CLI via core Makefile
	@$(MAKE) -C $(CORE_DIR) uninstall

%:
	@$(MAKE) -C $(CORE_DIR) $@
