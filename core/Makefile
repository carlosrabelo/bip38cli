GO            ?= go
BIN_NAME      := bip38cli
BUILD_DIR     := ../bin
VERSION       := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME    := $(shell date +%Y-%m-%dT%H:%M:%S%z)
CGO_ENABLED := 0
LDFLAGS     := -s -w -extldflags "-static" -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)
export GOROOT :=

.DEFAULT_GOAL := help

.PHONY: all build clean fmt help lint test

all: lint test build

build: ## Build bip38cli binary with version metadata (fully static)
	@echo "Building $(BIN_NAME) (static)..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) $(GO) build -trimpath -tags netgo -ldflags='$(LDFLAGS)' -o $(BUILD_DIR)/$(BIN_NAME) ./cmd/bip38cli

clean: ## Remove build outputs and go caches
	@echo "Cleaning build outputs..."
	@rm -rf $(BUILD_DIR)/$(BIN_NAME)
	@$(GO) clean -cache -testcache >/dev/null 2>&1 || true

fmt: ## Run gofmt across all Go sources
	@echo "Formatting Go sources..."
	@go fmt ./...

help: ## Show available targets
	@echo "BIP38CLI - core/ targets"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*## "} {printf "  %-15s %s\n", $$1, $$2}'

lint: ## Run golangci-lint when available
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not found" && exit 1; }
	@golangci-lint run ./...

test: ## Execute go test ./...
	@echo "Running go test ./..."
	@$(GO) test ./...
