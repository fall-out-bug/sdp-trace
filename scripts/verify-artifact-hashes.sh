#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

status=0

check_hash() {
  local source_file="$1"
  local field_kind="$2"
  local ref="$3"
  local expected="$4"
  local path="${ref%%#*}"

  case "$path" in
    http://*|https://*|role:*|human-review:*|team-lead|pi\ --*)
      return
      ;;
  esac

  if [[ ! -f "$path" ]]; then
    echo "Missing local artifact for verified hash in $source_file: $field_kind $ref" >&2
    status=1
    return
  fi

  local actual
  actual="$(shasum -a 256 "$path" | awk '{print $1}')"
  if [[ "$actual" != "$expected" ]]; then
    echo "Artifact hash mismatch in $source_file: $field_kind $ref" >&2
    echo "  expected: $expected" >&2
    echo "  actual:   $actual" >&2
    status=1
  fi
}

while IFS= read -r -d '' json_file; do
  while IFS=$'\t' read -r kind ref expected; do
    [[ -z "${kind:-}" ]] && continue
    check_hash "$json_file" "$kind" "$ref" "$expected"
  done < <(
    jq -r '
      (
        .. | objects
        | select(.artifact?.integrity_status? == "verified_hash" and .artifact.uri? and .artifact.sha256?)
        | ["artifact", .artifact.uri, .artifact.sha256]
        | @tsv
      ),
      (
        .. | objects
        | select(.integrity_status? == "verified_hash" and .uri? and .sha256?)
        | ["artifact", .uri, .sha256]
        | @tsv
      ),
      (
        .. | objects
        | select(.digest_algorithm? == "sha256" and .payload_digest? and .command? and (.command | test("^[A-Za-z0-9_./-]+$")))
        | ["payload_digest", .command, .payload_digest]
        | @tsv
      )
    ' "$json_file"
  )
done < <(find examples -name '*.json' -print0)

exit "$status"
