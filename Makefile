CORE_DIR       := core
GO             ?= go

.DEFAULT_GOAL  := help

.PHONY: build clean fmt help install lint test uninstall verify

build: ## Build CLI binary
	@$(MAKE) -C $(CORE_DIR) build

clean: ## Clean build artifacts from core/
	@$(MAKE) -C $(CORE_DIR) clean

fmt: ## Format Go sources with gofmt
	@$(MAKE) -C $(CORE_DIR) fmt

help: ## Show root level targets
	@printf "BIP38CLI - Bitcoin Private Key Encryption Tool\n"
	@printf "==============================================\n\n"
	@printf " Build & Install:\n"
	@printf "   build           Build CLI binary\n"
	@printf "   install         Install CLI system-wide\n"
	@printf "   uninstall       Remove CLI from system\n\n"
	@printf " Quality:\n"
	@printf "   fmt             Format Go sources with gofmt\n"
	@printf "   lint            Run golangci-lint via core/\n\n"
	@printf " Testing:\n"
	@printf "   test            Run go test ./... in core/\n\n"
	@printf " Utilities:\n"
	@printf "   clean           Clean build artifacts from core/\n"
	@printf "   verify          Run lint and test in one shot\n"
	@printf "   help            Show this help\n"

lint: ## Run golangci-lint via core/
	@$(MAKE) -C $(CORE_DIR) lint

test: ## Run go test ./... in core/
	@$(MAKE) -C $(CORE_DIR) test

install: ## Install CLI system-wide
	@$(MAKE) -C $(CORE_DIR) install

uninstall: ## Remove CLI from system
	@$(MAKE) -C $(CORE_DIR) uninstall

verify: ## Run lint and test in one shot
	@$(MAKE) lint
	@$(MAKE) test

%:
	@$(MAKE) -C $(CORE_DIR) $@
