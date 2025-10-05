GO            := go
BIN_NAME      := bip38cli
BUILD_DIR     := ../bin
VERSION       := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME    := $(shell date +%Y-%m-%dT%H:%M:%S%z)
LDFLAGS       := -s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

.DEFAULT_GOAL := help

.PHONY: build clean fmt help install lint test uninstall

build: ## Build bip38cli binary with version metadata
	@echo "Building $(BIN_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@env GOCACHE=$$(mktemp -d) $(GO) build -trimpath -tags netgo -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BIN_NAME) ./cmd/bip38cli

clean: ## Remove build outputs and go caches
	@echo "Cleaning build outputs..."
	@rm -rf $(BUILD_DIR)/$(BIN_NAME)
	@$(GO) clean -cache -testcache >/dev/null 2>&1 || true

fmt: ## Run gofmt across core/
	@echo "Formatting Go sources..."
	@find . -name '*.go' -not -path './vendor/*' -print0 | xargs -0 gofmt -w

help: ## Show go-specific targets
	@echo "BIP38CLI - core targets"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*## "} {printf "  %-15s %s\n", $$1, $$2}'

install: ## Install binary via scripts/install.sh
	@../scripts/install.sh --user

lint: ## Run golangci-lint when available
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not found" && exit 1; }
	@env GOCACHE=$$(mktemp -d) golangci-lint run ./...

test: ## Execute go test ./...
	@echo "Running go test ./..."
	@env GOCACHE=$$(mktemp -d) $(GO) test ./...

uninstall: ## Remove installed binary via scripts/uninstall.sh
	@../scripts/uninstall.sh --user
