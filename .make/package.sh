#!/usr/bin/env bash
set -euo pipefail

BINARY_NAME="${BINARY_NAME:-bip38cli}"
PACKAGE_NAME="${PACKAGE_NAME:-$BINARY_NAME}"
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
SRC="$ROOT_DIR/bin/$BINARY_NAME"
DIST_DIR="$ROOT_DIR/dist"
STAGE="$DIST_DIR/deb-root"

VERSION="${VERSION:-$(git -C "$ROOT_DIR" describe --tags --always --dirty 2>/dev/null || echo 0.0.0)}"
# Debian versions cannot have a leading "v".
VERSION="${VERSION#v}"
ARCH="$(dpkg --print-architecture 2>/dev/null || echo amd64)"
MAINTAINER="${MAINTAINER:-$(git -C "$ROOT_DIR" config user.email 2>/dev/null || echo unknown@localhost)}"
MAINTAINER_NAME="${MAINTAINER_NAME:-$(git -C "$ROOT_DIR" config user.name 2>/dev/null || echo Unknown)}"
DESCRIPTION="${DESCRIPTION:-$PACKAGE_NAME command-line tool}"

if [ ! -x "$SRC" ]; then
  echo "error: missing binary: $SRC" >&2
  echo "hint:  run make build first" >&2
  exit 1
fi

if ! command -v dpkg-deb >/dev/null 2>&1; then
  echo "error: dpkg-deb not found (install dpkg)" >&2
  exit 1
fi

rm -rf "$STAGE"
install -D -m 755 "$SRC" "$STAGE/usr/bin/$BINARY_NAME"
mkdir -p "$STAGE/DEBIAN"

cat >"$STAGE/DEBIAN/control" <<EOF
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: ${MAINTAINER_NAME} <${MAINTAINER}>
Description: ${DESCRIPTION}
EOF

OUT="${DIST_DIR}/${PACKAGE_NAME}_${VERSION}_${ARCH}.deb"
mkdir -p "$DIST_DIR"
dpkg-deb --root-owner-group --build "$STAGE" "$OUT"
rm -rf "$STAGE"

echo "Package ready: $OUT"
echo "Install with: sudo dpkg -i $OUT"
