#!/bin/bash
set -e

# Script to install bip38cli for the current user or system-wide.

# --- Configuration ---
APP_NAME="bip38cli"
BUILD_DIR="bin"
BINARY_PATH="${BUILD_DIR}/${APP_NAME}"
INSTALL_DIR_USER="${HOME}/.local/bin"
INSTALL_DIR_ROOT="/usr/local/bin"

# --- Main Logic ---
echo "Starting installation of ${APP_NAME}..."

# 1. Build the binary first
echo "Building the binary..."
if ! make build; then
    echo "Error: Build failed. Aborting installation." >&2
    exit 1
fi

if [ ! -f "${BINARY_PATH}" ]; then
    echo "Error: Binary not found at ${BINARY_PATH} after build. Aborting." >&2
    exit 1
fi

# 2. Determine installation directory and copy the binary
INSTALL_DIR=""
if [ "$(id -u)" -eq 0 ]; then
    # Running as root, install system-wide
    echo "Running as root. Installing to ${INSTALL_DIR_ROOT}..."
    INSTALL_DIR=${INSTALL_DIR_ROOT}
else
    # Running as a regular user, install to user's local bin
    echo "Running as user. Installing to ${INSTALL_DIR_USER}..."
    INSTALL_DIR=${INSTALL_DIR_USER}
fi

# Ensure the destination directory exists
echo "Ensuring destination directory exists: ${INSTALL_DIR}"
mkdir -p "${INSTALL_DIR}"

# Copy the binary
echo "Copying '${BINARY_PATH}' to '${INSTALL_DIR}'..."
cp "${BINARY_PATH}" "${INSTALL_DIR}/"

# 3. Final message
echo ""
echo "✅ ${APP_NAME} installed successfully to ${INSTALL_DIR}/${APP_NAME}"
echo "Please ensure '${INSTALL_DIR}' is in your shell's PATH."
