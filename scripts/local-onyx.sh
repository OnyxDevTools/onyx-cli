#!/usr/bin/env bash
# Build and install the local development CLI as `localonyx`.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# Prefer per-user install to avoid sudo; override with BINDIR if you want.
DEFAULT_BINDIR="$HOME/.local/bin"
BINDIR="${BINDIR:-$DEFAULT_BINDIR}"
FALLBACK="/usr/local/bin"

echo "Building localonyx from $ROOT ..."
cd "$ROOT"
OUT="$(mktemp -d)"
trap 'rm -rf "$OUT"' EXIT

go build -o "$OUT/localonyx" ./cmd/localonyx

install_bin() {
  local src="$1"
  local destdir="$2"
  mkdir -p "$destdir"
  install -m 0755 "$src" "$destdir/localonyx"
  echo "$destdir/localonyx"
}

ensure_path() {
  local dir="$1"
  case ":$PATH:" in
    *":$dir:"*) return 0 ;;
  esac
  local zshrc="$HOME/.zshrc"
  if ! grep -qs "$dir" "$zshrc" 2>/dev/null; then
    echo "export PATH=\"$dir:\$PATH\"" >> "$zshrc"
    echo "Added $dir to PATH in $zshrc"
  fi
  echo "Note: restart your shell or 'source $zshrc' to pick up the new PATH."
}

target=""
if target="$(install_bin "$OUT/localonyx" "$BINDIR" 2>/dev/null)"; then
  :
elif target="$(install_bin "$OUT/localonyx" "$FALLBACK" 2>/dev/null)"; then
  echo "Could not write to $BINDIR, installed to fallback."
else
  echo "Failed to install localonyx (try setting BINDIR to a writable directory)" >&2
  exit 1
fi

echo "localonyx installed to $target"
ensure_path "$(dirname "$target")"
echo "Try: localonyx gen --schema ../examples/api/onyx.schema.json --go"
