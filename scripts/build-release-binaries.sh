#!/usr/bin/env bash
set -euo pipefail

# Thin release launcher: Go owns the product binary; this script only repeats
# the cross-compilation matrix used by CI and local release verification.
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="${VERSION:-$(git -C "$ROOT" describe --tags --always --dirty)}"
OUT_DIR="${OUT_DIR:-$ROOT/dist}"

mkdir -p "$OUT_DIR"

build_one() {
  local goos="$1"
  local goarch="$2"
  local ext=""
  if [[ "$goos" == "windows" ]]; then
    ext=".exe"
  fi
  local name="sdp-trace_${VERSION}_${goos}_${goarch}${ext}"
  (
    cd "$ROOT"
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
      go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" \
      -o "$OUT_DIR/$name" ./cmd/sdp-trace
  )
  sha256sum "$OUT_DIR/$name" > "$OUT_DIR/$name.sha256"
}

build_one linux amd64
build_one linux arm64
build_one darwin amd64
build_one darwin arm64
build_one windows amd64
build_one windows arm64

(
  cd "$OUT_DIR"
  find . -maxdepth 1 -type f -name "sdp-trace_${VERSION}_*" ! -name "*.sha256" \
    -printf '%f\n' | sort | xargs sha256sum > SHA256SUMS
)

printf 'built sdp-trace %s into %s\n' "$VERSION" "$OUT_DIR"
