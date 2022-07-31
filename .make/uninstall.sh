#!/usr/bin/env bash
set -euo pipefail

BINARY_NAME="${BINARY_NAME:-bip38cli}"

# Prefer the invoking user's home when running under sudo.
user_home() {
  if [ "$(id -u)" -eq 0 ] && [ -n "${SUDO_USER:-}" ] && [ "$SUDO_USER" != "root" ]; then
    getent passwd "$SUDO_USER" | cut -d: -f6
    return
  fi
  echo "${HOME}"
}

USER_LOCAL="$(user_home)/.local/bin/${BINARY_NAME}"

targets=()
if [ -n "${PREFIX:-}" ]; then
  targets+=("${PREFIX%/}/bin/${BINARY_NAME}")
else
  targets+=("$USER_LOCAL")
fi

removed=0
for target in "${targets[@]}"; do
  if [ ! -e "$target" ]; then
    continue
  fi
  dir="$(dirname "$target")"
  # Directory write bit is enough to unlink; escalate when the dir is root-owned
  # (e.g. leftover from an older sudo make install).
  if [ -w "$dir" ] 2>/dev/null && rm -f "$target" 2>/dev/null; then
    :
  else
    sudo rm -f "$target"
  fi
  echo "Removed: $target"
  removed=$((removed + 1))
done

if [ "$removed" -eq 0 ]; then
  echo "Not installed (checked: ${targets[*]})"
fi
