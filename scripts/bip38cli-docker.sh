#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
IMAGE_NAME="bip38cli:local"
DOCKERFILE_PATH="${REPO_ROOT}/docker/Dockerfile"

usage() {
	cat <<EOF
Usage: $(basename "$0") [--build] [arguments...]

Build (if needed) and run the bip38cli Docker image. All remaining arguments
are passed directly to the CLI inside the container. When no arguments are
provided the binary shows the standard help output.

  --build     Force a rebuild of the Docker image before running
  --help      Show this message
EOF
}

if ! command -v docker >/dev/null 2>&1; then
	echo "docker is required but was not found in PATH" >&2
	exit 1
fi

force_build=false

if [[ "${1:-}" == "--help" ]]; then
	usage
	exit 0
fi

if [[ "${1:-}" == "--build" ]]; then
	force_build=true
	shift
fi

if [[ ! -f "${DOCKERFILE_PATH}" ]]; then
	echo "Dockerfile not found at ${DOCKERFILE_PATH}" >&2
	exit 1
fi

if $force_build || ! docker image inspect "${IMAGE_NAME}" >/dev/null 2>&1; then
	echo "Building ${IMAGE_NAME}..."
	docker build -f "${DOCKERFILE_PATH}" -t "${IMAGE_NAME}" "${REPO_ROOT}"
fi

run_flags=("-i")
if [ -t 1 ]; then
	run_flags+=("-t")
fi

docker run --rm "${run_flags[@]}" "${IMAGE_NAME}" "$@"
