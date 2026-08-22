MAKEFLAGS += --no-print-directory

.DEFAULT_GOAL := help

.PHONY: all build clean fmt help install lint quality test uninstall version vet

BINARY_NAME := bip38cli
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME  := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

all: quality test build ## Run quality checks, tests, and build

build: ## Build bip38cli binary with version metadata
	@./.make/build.sh

clean: ## Remove build outputs and go caches
	@./.make/clean.sh

fmt: ## Format Go sources with gofmt
	@go fmt ./...

help: ## Show available targets
	@echo "$(BINARY_NAME) - Available targets"
	@echo ""
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*## "} {printf "  %-15s %s\n", $$1, $$2}'

install: build ## Install binary to $HOME/.local/bin
	@./.make/install.sh --user

lint: ## Run golangci-lint when available
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not found" && exit 1; }
	@golangci-lint run ./...

quality: fmt vet lint ## Run all quality checks

test: ## Execute go test ./...
	@./.make/test.sh

uninstall: ## Remove binary from $HOME/.local/bin
	@./.make/uninstall.sh --user

version: ## Show version
	@echo "$(BINARY_NAME) $(VERSION) ($(BUILD_TIME))"

vet: ## Run go vet
	@go vet ./...
