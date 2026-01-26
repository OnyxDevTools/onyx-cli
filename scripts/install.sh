#!/usr/bin/env bash
set -euo pipefail

REPO="OnyxDevTools/onyx-cli"
TAG="${1:-${ONYX_TAG:-latest}}"
BINDIR="${BINDIR:-/usr/local/bin}"

case "$(uname -s)" in
  Darwin) OS="darwin" ;;
  Linux)  OS="linux" ;;
  *) echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported ARCH: $(uname -m)" >&2; exit 1 ;;
esac

if [ "$TAG" = "latest" ]; then
  URL="https://github.com/${REPO}/releases/latest/download/onyx_${OS}_${ARCH}.tar.gz"
else
  URL="https://github.com/${REPO}/releases/download/${TAG}/onyx_${OS}_${ARCH}.tar.gz"
fi

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading $URL ..."
curl -fsSL "$URL" -o "$TMPDIR/onyx.tgz"

tar -C "$TMPDIR" -xzf "$TMPDIR/onyx.tgz"

BIN_SRC="$TMPDIR/onyx"
if [ ! -f "$BIN_SRC" ]; then
  BIN_SRC="$(find "$TMPDIR" -type f -name onyx -perm -u+x | head -n1)"
fi
if [ -z "${BIN_SRC:-}" ] || [ ! -f "$BIN_SRC" ]; then
  echo "Failed to find extracted binary in archive" >&2
  exit 1
fi

if install -m 0755 "$BIN_SRC" "$BINDIR/onyx" 2>/dev/null; then
  TARGET="$BINDIR/onyx"
else
  FALLBACK="$HOME/.local/bin"
  mkdir -p "$FALLBACK"
  install -m 0755 "$BIN_SRC" "$FALLBACK/onyx"
  TARGET="$FALLBACK/onyx"
  echo "Note: could not write to $BINDIR, installed to $FALLBACK"
fi

echo "onyx installed to $TARGET"
echo "Ensure $(dirname "$TARGET") is on your PATH."
