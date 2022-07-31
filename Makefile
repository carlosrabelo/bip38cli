MAKEFLAGS += --no-print-directory

.DEFAULT_GOAL := help

.PHONY: all build clean deb fmt help install lint package quality test uninstall version vet

BINARY_NAME := bip38cli
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME  := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

export BINARY_NAME

all: quality test build ## Run quality checks, tests, and build

help: ## Show available targets
	@echo "$(BINARY_NAME) - Available targets"
	@echo ""
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*## "} {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary
	@./.make/build.sh

test: ## Run tests
	@./.make/test.sh

lint: ## Run linter
	@go vet ./...

fmt: ## Format code
	@go fmt ./...

quality: fmt lint ## Run all quality checks

clean: ## Remove build artifacts
	@./.make/clean.sh

install: ## Install to ~/.local/bin (requires prior make build)
	@./.make/install.sh

uninstall: ## Remove from ~/.local/bin (sudo only if needed)
	@./.make/uninstall.sh

package: ## Build a .deb into dist/ (requires prior make build)
	@./.make/package.sh

deb: package ## Alias for package

version: ## Show version
	@echo "$(BINARY_NAME) $(VERSION) ($(BUILD_TIME))"

vet: lint ## Run go vet
