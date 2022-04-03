#!/usr/bin/env bash
set -euo pipefail

BIN_NAME="bip38cli"
TARGET_DIR="/usr/local/bin"
MODE="system"

usage() {
  echo "Usage: $0 [--user]" >&2
  exit 1
}

resolve_user_target_dir() {
  if [[ -d "${HOME}/.local/go/bin" ]] && echo "${PATH}" | tr ':' '\n' | grep -qx "${HOME}/.local/go/bin"; then
    echo "${HOME}/.local/go/bin"
  else
    echo "${HOME}/.local/bin"
  fi
}

if [[ "${1:-}" == "--user" ]]; then
  TARGET_DIR="$(resolve_user_target_dir)"
  MODE="user"
  shift
fi

if [[ $# -gt 0 ]]; then
  usage
fi

TARGET_PATH="${TARGET_DIR}/${BIN_NAME}"

CANDIDATES=("${TARGET_PATH}")
if [[ "${MODE}" == "user" ]]; then
  for alt in "${HOME}/.local/bin/${BIN_NAME}" "${HOME}/.local/go/bin/${BIN_NAME}"; do
    if [[ "${alt}" != "${TARGET_PATH}" ]]; then
      CANDIDATES+=("${alt}")
    fi
  done
fi

removed=0
for path in "${CANDIDATES[@]}"; do
  dir="$(dirname "${path}")"
  if [[ -e "${path}" ]]; then
    if [[ ! -w "${dir}" ]]; then
      echo "No write access to ${dir}. Use sudo or --user option." >&2
      exit 1
    fi
    rm -f "${path}"
    echo "Removed ${BIN_NAME} from ${dir}" >&2
    removed=1
  fi
done

if [[ "${removed}" -eq 0 ]]; then
  echo "No ${BIN_NAME} binary found at ${TARGET_PATH}" >&2
fi

if [[ "${MODE}" == "user" ]] && ! echo "${PATH}" | tr ':' '\n' | grep -qx "${TARGET_DIR}"; then
  echo "Reminder: ${TARGET_DIR} not in PATH." >&2
fi
