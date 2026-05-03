#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

envelope="${1:-examples/contract-foundation/contract-release.dsse.json}"
public_key="${2:-examples/contract-foundation/local-dev-signing-public.pem}"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

payload_type="$(jq -r '.payloadType' "$envelope")"
jq -r '.payload' "$envelope" | base64 --decode >"$tmpdir/payload.json"
jq -r '.signatures[0].sig' "$envelope" | base64 --decode >"$tmpdir/signature.bin"

payload_len="$(wc -c <"$tmpdir/payload.json" | tr -d ' ')"
type_len="$(printf '%s' "$payload_type" | wc -c | tr -d ' ')"
{
  printf 'DSSEv1 %s %s %s ' "$type_len" "$payload_type" "$payload_len"
  cat "$tmpdir/payload.json"
} >"$tmpdir/pae.bin"

openssl dgst -sha256 -verify "$public_key" -signature "$tmpdir/signature.bin" "$tmpdir/pae.bin" >/dev/null

manifest_ref="$(jq -r '.subject[0].name' "$tmpdir/payload.json")"
manifest_digest="$(jq -r '.subject[0].digest.sha256' "$tmpdir/payload.json")"
actual="$(shasum -a 256 "$manifest_ref" | awk '{print $1}')"

if [[ "$actual" != "$manifest_digest" ]]; then
  echo "Signed payload manifest digest mismatch" >&2
  echo "  expected: $manifest_digest" >&2
  echo "  actual:   $actual" >&2
  exit 1
fi

