#!/usr/bin/env bash
set -euo pipefail

provider="${1:?provider required}"
model="${2:?model required}"
prompt_file="$(mktemp)"
trap 'rm -f "$prompt_file"' EXIT

cat >"$prompt_file"

pi --provider "$provider" \
  --model "$model" \
  --no-tools \
  --no-context-files \
  --no-session \
  -p @"$prompt_file" \
  2>/dev/null
