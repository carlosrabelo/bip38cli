# BIP38CLI - BIP38 Bitcoin Key Encryption Tool
# Build automation for Go project

# Variables
GO		= go
BIN		= bip38cli
SRC		= ./cmd/bip38cli
BUILD_DIR	= ./bin
VERSION		= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME	= $(shell date +%Y-%m-%dT%H:%M:%S%z)
LDFLAGS		= -s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

# Default target - show help
.DEFAULT_GOAL := help

.PHONY: help all build clean run test fmt vet lint deps mod-tidy info

help:	## Show this help
	@echo "BIP38CLI - BIP38 Bitcoin Key Encryption Tool"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'

all: clean fmt vet build	## Clean, format, vet and build

build:	## Build the project
	@echo "Building $(BIN)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GO) build -trimpath -tags netgo -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BIN) $(SRC)

clean:	## Remove build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@$(GO) clean -cache -testcache -modcache 2>/dev/null || true

run:	## Run the application
	@$(GO) run $(SRC)

test:	## Run tests
	@echo "Running tests..."
	@$(GO) test -v -race -coverprofile=coverage.out ./...

test-coverage:	## Run tests with coverage report
	@echo "Running tests with coverage..."
	@$(GO) test -v -race -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

fmt:	## Format code
	@echo "Formatting code..."
	@$(GO) fmt ./...

vet:	## Vet code
	@echo "Vetting code..."
	@$(GO) vet ./...

lint:	## Lint code with golangci-lint
	@echo "Linting code..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run

deps:	## Download dependencies
	@echo "Downloading dependencies..."
	@$(GO) mod download

mod-tidy:	## Clean up go.mod and go.sum
	@echo "Tidying modules..."
	@$(GO) mod tidy

install:	## Install the binary
	@echo "Installing $(BIN)..."
	@$(GO) install -ldflags="$(LDFLAGS)" $(SRC)

info:	## Show project information
	@echo "Project: BIP38CLI"
	@echo "Binary: $(BIN)"
	@echo "Source: $(SRC)"
	@echo "Build: $(BUILD_DIR)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go version: $$($(GO) version)"