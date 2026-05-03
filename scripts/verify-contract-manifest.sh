#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

manifest="${1:-examples/contract-foundation/contract-manifest.example.json}"

jq empty "$manifest" >/dev/null

status=0
while IFS=$'\t' read -r path digest; do
  if [[ ! -f "$path" ]]; then
    echo "Manifest artifact missing: $path" >&2
    status=1
    continue
  fi

  actual="$(shasum -a 256 "$path" | awk '{print $1}')"
  if [[ "$actual" != "$digest" ]]; then
    echo "Digest mismatch for $path" >&2
    echo "  expected: $digest" >&2
    echo "  actual:   $actual" >&2
    status=1
  fi
done < <(jq -r '.artifacts[] | [.path, .sha256] | @tsv' "$manifest")

exit "$status"

