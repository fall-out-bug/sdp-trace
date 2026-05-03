#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

tmp="$(mktemp)"
out="$(mktemp)"
trap 'rm -f "$tmp" "$out"' EXIT

find . -type f \
  -not -path './.git/*' \
  -not -path './.beads/*' \
  -not -path './.sdp-trace-runs/*' \
  -not -path './benchmarks/repos/*' \
  -not -path './node_modules/*' \
  -not -path './scripts/check-artifact-safety.sh' \
  -not -path './.DS_Store' \
  -print0 >"$tmp"

patterns='(AKIA[0-9A-Z]{16}|-----BEGIN [A-Z ]*PRIVATE KEY-----|password[[:space:]]*=|api[_-]?key[[:space:]]*=|secret[[:space:]]*=|raw_customer_data|private_prompt_contents)'

if xargs -0 grep -EIni "$patterns" <"$tmp" >"$out" 2>/dev/null; then
  echo "Artifact safety scan found prohibited secret/customer-data markers:" >&2
  cat "$out" >&2
  exit 1
fi
