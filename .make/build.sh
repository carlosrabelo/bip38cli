#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_NAME="bip38cli"

VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo dev)"
BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
LDFLAGS="-s -w -extldflags \"-static\" -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

mkdir -p "${ROOT_DIR}/bin"
cd "${ROOT_DIR}"

echo "Building ${BIN_NAME} (static)..."
CGO_ENABLED=0 go build -trimpath -tags netgo -ldflags="${LDFLAGS}" \
  -o "${ROOT_DIR}/bin/${BIN_NAME}" ./bip38cli/cmd/bip38cli

echo "Binary ready at: ${ROOT_DIR}/bin/${BIN_NAME}"
