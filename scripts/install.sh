#!/bin/bash
set -euo pipefail

BIN_NAME="bip38cli"
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_PATH="${PROJECT_ROOT}/bin/${BIN_NAME}"
TARGET_DIR="/usr/local/bin"
MODE="system"

usage() {
  echo "Usage: $0 [--user]" >&2
  exit 1
}

if [[ "${1:-}" == "--user" ]]; then
  TARGET_DIR="${HOME}/.local/bin"
  MODE="user"
  shift
fi

if [[ $# -gt 0 ]]; then
  usage
fi

if [[ ! -f "${BIN_PATH}" ]]; then
  echo "Binary not found at ${BIN_PATH}. Run make build first." >&2
  exit 1
fi

if [[ ! -d "${TARGET_DIR}" ]]; then
  mkdir -p "${TARGET_DIR}"
fi

if [[ ! -w "${TARGET_DIR}" ]]; then
  echo "No write access to ${TARGET_DIR}. Use sudo or --user option." >&2
  exit 1
fi

install -m 0755 "${BIN_PATH}" "${TARGET_DIR}/${BIN_NAME}"

if [[ "${MODE}" == "user" ]] && ! echo "${PATH}" | tr ':' '\n' | grep -qx "${TARGET_DIR}"; then
  echo "Warning: ${TARGET_DIR} not in PATH. Add it to your shell profile." >&2
fi

echo "Installed ${BIN_NAME} to ${TARGET_DIR}" >&2
