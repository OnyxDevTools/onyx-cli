#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION_FILE="$ROOT/cmd/onyx/cmd/version.go"
FORMULA_FILE="$ROOT/Formula/onyx.rb"
DIST_DIR="$ROOT/dist"

usage() {
  cat <<'EOF'
Usage: scripts/bump-version.sh [major|minor|patch]

Actions (always performed):
  - bump cmd/onyx/cmd/version.go
  - build darwin/linux (amd64, arm64) tarballs into dist/
  - update Formula/onyx.rb with version + shas
  - git add/commit/tag/push
  - create GitHub release with attached tarballs (installs gh CLI if missing)

Environment:
  BUMP=major|minor|patch  Alternative to positional arg
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage; exit 0
fi

bwarn() { printf "WARNING: %s\n" "$*" >&2; }
info() { printf "==> %s\n" "$*" >&2; }

bump="${1:-${BUMP:-}}"
if [[ -z "$bump" ]]; then
  read -rp "Bump type (major/minor/patch) [patch]: " bump
  bump="${bump:-patch}"
fi

# Guardrails
status=$(cd "$ROOT" && git status --porcelain)
if [[ -n "$status" ]]; then
  bwarn "Working tree not clean. Commit or stash first."
  exit 1
fi
ahead=$(cd "$ROOT" && git rev-list --count --left-only @{u}...HEAD 2>/dev/null || echo 0)
behind=$(cd "$ROOT" && git rev-list --count --right-only @{u}...HEAD 2>/dev/null || echo 0)
if [[ "$behind" != "0" ]]; then
  bwarn "Local branch is behind upstream. Pull/merge before releasing."
  exit 1
fi

current=$(grep 'Version = "' "$VERSION_FILE" | sed -E 's/.*"([^"]+)".*/\1/')
if [[ ! "$current" =~ ^([0-9]+)\.([0-9]+)\.([0-9]+) ]]; then
  echo "Unable to parse current version from $VERSION_FILE (got '$current')" >&2
  exit 1
fi
maj=${BASH_REMATCH[1]}
min=${BASH_REMATCH[2]}
pat=${BASH_REMATCH[3]}

case "$bump" in
  major) ((maj++)); min=0; pat=0 ;;
  minor) ((min++)); pat=0 ;;
  patch) ((pat++)) ;;
  *) echo "Unknown bump '$bump' (use major|minor|patch)" >&2; exit 1 ;;
esac

next="${maj}.${min}.${pat}"

tmp="$VERSION_FILE.tmp"
perl -0777 -pe "s/Version = \"[^\"]+\"/Version = \"${next}\"/" "$VERSION_FILE" > "$tmp"
mv "$tmp" "$VERSION_FILE"

echo "Bumped version: $current -> $next"

rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

commit=$(cd "$ROOT" && git rev-parse --short HEAD 2>/dev/null || echo "unknown")
date=$(date -u +%Y-%m-%d)
ldflags="-s -w -X github.com/OnyxDevTools/onyx-cli/cmd/onyx/cmd.Version=${next} -X github.com/OnyxDevTools/onyx-cli/cmd/onyx/cmd.Commit=${commit} -X github.com/OnyxDevTools/onyx-cli/cmd/onyx/cmd.Date=${date}"

oses=("darwin" "linux")
arches=("amd64" "arm64")

for os in "${oses[@]}"; do
  for arch in "${arches[@]}"; do
    outdir="$DIST_DIR/onyx_${os}_${arch}"
    mkdir -p "$outdir"
    echo "Building ${os}/${arch}..."
    (cd "$ROOT" && GOOS="$os" GOARCH="$arch" go build -o "$outdir/onyx" -ldflags "$ldflags" ./cmd/onyx)
    tar -C "$outdir" -czf "$DIST_DIR/onyx_${os}_${arch}.tar.gz" onyx
  done
done

echo "Computing sha256..."
sha_darwin_amd64=$(shasum -a 256 "$DIST_DIR/onyx_darwin_amd64.tar.gz" | awk '{print $1}')
sha_darwin_arm64=$(shasum -a 256 "$DIST_DIR/onyx_darwin_arm64.tar.gz" | awk '{print $1}')

if [[ -f "$FORMULA_FILE" ]]; then
  perl -0777 -pe "s/version \"[^\"]+\"/version \"${next}\"/;
                   s/onyx_darwin_amd64\\.tar\\.gz\"\\n\\s*sha256 \"[^\"]+\"/onyx_darwin_amd64.tar.gz\"\\n      sha256 \"${sha_darwin_amd64}\"/;
                   s/onyx_darwin_arm64\\.tar\\.gz\"\\n\\s*sha256 \"[^\"]+\"/onyx_darwin_arm64.tar.gz\"\\n      sha256 \"${sha_darwin_arm64}\"/" \
    "$FORMULA_FILE" > "$FORMULA_FILE.tmp"
  mv "$FORMULA_FILE.tmp" "$FORMULA_FILE"
  echo "Updated $FORMULA_FILE"
fi

echo
echo "Build artifacts:"
(cd "$DIST_DIR" && ls -1 onyx_*.tar.gz)
echo

info "Committing and tagging..."
git add "$VERSION_FILE" "$FORMULA_FILE" "$DIST_DIR"
git commit -m "release v${next}"
git tag "v${next}"
git push origin HEAD --tags
info "Pushed v${next}"

ensure_gh() {
  if command -v gh >/dev/null 2>&1; then
    return 0
  fi
  info "gh CLI not found; installing locally..."
  os="$(uname -s)"
  arch="$(uname -m)"
  case "$os" in
    Darwin) gh_os="macOS" ;;
    Linux) gh_os="linux" ;;
    *) echo "Unsupported OS for gh install: $os" >&2; exit 1 ;;
  esac
  case "$arch" in
    x86_64|amd64) gh_arch="amd64" ;;
    arm64|aarch64) gh_arch="arm64" ;;
    *) echo "Unsupported arch for gh install: $arch" >&2; exit 1 ;;
  esac
  tag=$(curl -fsSL https://api.github.com/repos/cli/cli/releases/latest | sed -n 's/.*"tag_name": *"\\(.*\\)".*/\\1/p' | head -n1)
  if [[ -z "$tag" ]]; then
    echo "Failed to resolve gh latest tag" >&2; exit 1
  fi
  tmp=$(mktemp -d)
  trap 'rm -rf "$tmp"' EXIT
  tarball="gh_${tag#v}_${gh_os}_${gh_arch}.tar.gz"
  url="https://github.com/cli/cli/releases/download/${tag}/${tarball}"
  info "Downloading $url ..."
  curl -fsSL "$url" -o "$tmp/gh.tgz"
  tar -C "$tmp" -xzf "$tmp/gh.tgz"
  binpath=$(find "$tmp" -type f -path "*/bin/gh" | head -n1)
  if [[ -z "$binpath" ]]; then
    echo "Failed to unpack gh binary" >&2; exit 1
  fi
  target="${HOME}/.local/bin"
  mkdir -p "$target"
  install -m 0755 "$binpath" "$target/gh"
  export PATH="$target:$PATH"
  info "Installed gh to $target/gh (ensure PATH includes $target)"
}

ensure_gh

info "Creating GitHub release v${next}..."
gh release create "v${next}" "$DIST_DIR"/onyx_*.tar.gz --latest --notes "Release v${next}"
info "Done. Version: v${next}"
