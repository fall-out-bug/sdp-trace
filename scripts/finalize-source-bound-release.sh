#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

repo="$ROOT_DIR"
manifest="examples/contract-foundation/contract-manifest.example.json"
source_ref="HEAD"
output_file=""

usage() {
  cat >&2 <<'USAGE'
Usage: scripts/finalize-source-bound-release.sh [--repo path] [--manifest path] [--source-ref ref] [--output path]

Verifies that a clean git source reference contains every manifest artifact with
matching SHA-256 digest. This command records source-bound local proof only; it
does not claim external or production release trust.
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      repo="$2"
      shift 2
      ;;
    --manifest)
      manifest="$2"
      shift 2
      ;;
    --source-ref)
      source_ref="$2"
      shift 2
      ;;
    --output)
      output_file="$2"
      shift 2
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      usage
      exit 2
      ;;
  esac
done

if ! git -C "$repo" rev-parse --git-dir >/dev/null 2>&1; then
  echo "Not a git repository: $repo" >&2
  exit 2
fi

dirty_status="$(git -C "$repo" status --porcelain --untracked-files=normal)"
if [[ -n "$dirty_status" ]]; then
  echo "Refusing source-bound finalization from dirty working tree: $repo" >&2
  echo "$dirty_status" >&2
  exit 1
fi

source_commit="$(git -C "$repo" rev-parse --verify "$source_ref^{commit}" 2>/dev/null || true)"
if [[ ! "$source_commit" =~ ^[0-9a-f]{40}$ ]]; then
  echo "Source ref does not resolve to a 40-character commit SHA: $source_ref" >&2
  exit 1
fi

if [[ "$manifest" = /* ]]; then
  manifest_path="$manifest"
  manifest_ref="${manifest#$repo/}"
else
  manifest_path="$repo/$manifest"
  manifest_ref="$manifest"
fi

if [[ ! -f "$manifest_path" ]]; then
  echo "Manifest not found: $manifest" >&2
  exit 1
fi

jq empty "$manifest_path" >/dev/null

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

checked=0
missing=0
mismatched=0
total=0

while IFS=$'\t' read -r artifact_path expected_sha; do
  [[ -z "$artifact_path" ]] && continue
  total=$((total + 1))

  if ! git -C "$repo" cat-file -e "$source_commit:$artifact_path" 2>/dev/null; then
    missing=$((missing + 1))
    continue
  fi

  source_file="$tmpdir/source-artifact-$checked"
  git -C "$repo" show "$source_commit:$artifact_path" >"$source_file"
  actual_sha="$(shasum -a 256 "$source_file" | awk '{print $1}')"
  if [[ "$actual_sha" != "$expected_sha" ]]; then
    mismatched=$((mismatched + 1))
  fi
  checked=$((checked + 1))
done < <(jq -r '.artifacts[] | [.path, .sha256] | @tsv' "$manifest_path")

if [[ "$total" -eq 0 ]]; then
  echo "Manifest contains no artifacts to verify: $manifest" >&2
  exit 1
fi

if [[ "$missing" -gt 0 || "$mismatched" -gt 0 ]]; then
  echo "Source commit does not contain the manifest artifact set with matching digests: $source_commit" >&2
  echo "checked=$checked missing=$missing mismatched=$mismatched" >&2
  exit 1
fi

result="$(jq -n \
  --arg id "source-bound-local-finalization" \
  --arg generated_at "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
  --arg manifest_ref "$manifest_ref" \
  --arg source_commit "$source_commit" \
  --argjson checked "$checked" \
  --argjson missing "$missing" \
  --argjson mismatched "$mismatched" \
  '{
    id: $id,
    schema_version: "0.1.0",
    generated_at: $generated_at,
    manifest_ref: $manifest_ref,
    source_commit: $source_commit,
    source_commit_status: "matched",
    source_commit_artifact_status: "matched",
    source_commit_artifact_counts: {
      checked: $checked,
      missing: $missing,
      mismatched: $mismatched
    },
    proof_states: {
      source_commit_artifacts_verified: {
        state: "assessed",
        value: true,
        reason: "source_commit contains every manifest artifact path with matching digest"
      },
      externally_attested: {
        state: "not_assessed",
        value: null,
        reason: "source-bound local finalization does not include Sigstore/Rekor or accepted customer PKI evidence"
      },
      production_release_verified: {
        state: "not_assessed",
        value: null,
        reason: "production release verification requires external trust evidence in addition to source-bound local proof"
      }
    },
    trusted_contract_release: false
  }')"

if [[ -n "$output_file" ]]; then
  printf '%s\n' "$result" >"$output_file"
else
  printf '%s\n' "$result"
fi
