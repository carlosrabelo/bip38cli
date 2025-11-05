#!/bin/bash
set -euo pipefail

BIN_NAME="bip38cli"
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

TARGET_PATH="${TARGET_DIR}/${BIN_NAME}"

if [[ ! -e "${TARGET_PATH}" ]]; then
  echo "No ${BIN_NAME} binary found at ${TARGET_PATH}" >&2
  exit 0
fi

if [[ ! -w "${TARGET_DIR}" ]]; then
  echo "No write access to ${TARGET_DIR}. Use sudo or --user option." >&2
  exit 1
fi

rm -f "${TARGET_PATH}"

echo "Removed ${BIN_NAME} from ${TARGET_DIR}" >&2

if [[ "${MODE}" == "user" ]] && ! echo "${PATH}" | tr ':' '\n' | grep -qx "${TARGET_DIR}"; then
  echo "Reminder: ${TARGET_DIR} not in PATH." >&2
fi
