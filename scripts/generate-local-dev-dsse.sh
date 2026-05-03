#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

manifest="${1:-examples/contract-foundation/contract-manifest.example.json}"
public_key="${2:-examples/contract-foundation/local-dev-signing-public.pem}"
envelope="${3:-examples/contract-foundation/contract-release.dsse.json}"

digest="$(shasum -a 256 "$manifest" | awk '{print $1}')"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out "$tmp/local-dev-private.pem" >/dev/null 2>&1
openssl rsa -in "$tmp/local-dev-private.pem" -pubout -out "$public_key" >/dev/null 2>&1
fingerprint="$(openssl pkey -pubin -in "$public_key" -outform DER | openssl dgst -sha256 | awk '{print $NF}')"

jq -c -n \
  --arg manifest "$manifest" \
  --arg digest "$digest" \
  --arg identity "local-dev-dsse-openssl-v1:$fingerprint" \
  --arg policy "examples/contract-foundation/trusted-identity-policy.example.json" \
  '{
    "_type": "https://in-toto.io/Statement/v1",
    "subject": [
      {
        "name": $manifest,
        "digest": {
          "sha256": $digest
        }
      }
    ],
    "predicateType": "https://schemas.sdp.dev/trace/contract-release-verification/v0.1",
    "predicate": {
      "signing_profile": "sdp-trace-signature/sigstore-dsse-keyless-v1",
      "private_equivalent_profile": "local-dev-dsse-openssl-v1",
      "signer_identity": $identity,
      "identity_policy_ref": $policy
    }
  }' >"$tmp/payload.json"

payload_type="application/vnd.in-toto+json"
payload_len="$(wc -c <"$tmp/payload.json" | tr -d ' ')"
type_len="$(printf '%s' "$payload_type" | wc -c | tr -d ' ')"

{
  printf 'DSSEv1 %s %s %s ' "$type_len" "$payload_type" "$payload_len"
  cat "$tmp/payload.json"
} >"$tmp/pae.bin"

openssl dgst -sha256 -sign "$tmp/local-dev-private.pem" -out "$tmp/signature.bin" "$tmp/pae.bin"

payload_b64="$(base64 <"$tmp/payload.json" | tr -d '\n')"
sig_b64="$(base64 <"$tmp/signature.bin" | tr -d '\n')"

jq -n \
  --arg payloadType "$payload_type" \
  --arg payload "$payload_b64" \
  --arg keyid "$fingerprint" \
  --arg sig "$sig_b64" \
  '{
    "payloadType": $payloadType,
    "payload": $payload,
    "signatures": [
      {
        "keyid": $keyid,
        "sig": $sig
      }
    ]
  }' >"$envelope"

printf '%s %s\n' "$digest" "$fingerprint"
