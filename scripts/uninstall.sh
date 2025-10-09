#!/bin/bash
set -e

# Script to uninstall bip38cli.

# --- Configuration ---
APP_NAME="bip38cli"
INSTALL_DIR_USER="${HOME}/.local/bin"
INSTALL_DIR_ROOT="/usr/local/bin"

# --- Main Logic ---
echo "Starting uninstallation of ${APP_NAME}..."

# Determine potential installation locations
USER_PATH="${INSTALL_DIR_USER}/${APP_NAME}"
ROOT_PATH="${INSTALL_DIR_ROOT}/${APP_NAME}"

REMOVED=false
if [ -f "${ROOT_PATH}" ]; then
    if [ "$(id -u)" -eq 0 ]; then
        echo "Found in ${INSTALL_DIR_ROOT}. Removing..."
        rm -f "${ROOT_PATH}"
        REMOVED=true
        echo "✅ Uninstalled from ${INSTALL_DIR_ROOT}."
    else
        echo "Found in ${INSTALL_DIR_ROOT}, but you are not root. Please run 'sudo make uninstall' to remove." >&2
    fi
fi

if [ -f "${USER_PATH}" ]; then
    echo "Found in ${INSTALL_DIR_USER}. Removing..."
    rm -f "${USER_PATH}"
    REMOVED=true
    echo "✅ Uninstalled from ${INSTALL_DIR_USER}."
fi

if [ "$REMOVED" = false ]; then
    echo "Could not find ${APP_NAME} in standard locations (${INSTALL_DIR_ROOT}, ${INSTALL_DIR_USER}). Nothing to do."
fi
