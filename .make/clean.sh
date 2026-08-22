#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Cleaning build outputs..."
rm -rf "${ROOT_DIR}/bin/bip38cli"
go clean -cache -testcache >/dev/null 2>&1 || true
