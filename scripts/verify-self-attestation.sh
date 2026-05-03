#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

DEFAULT_CASE="examples/self-trace/self-attestation-verification.json"

case_file="$DEFAULT_CASE"
output_file=""
run_negative_suite=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --case)
      case_file="$2"
      shift 2
      ;;
    --output)
      output_file="$2"
      shift 2
      ;;
    --all)
      run_negative_suite=1
      shift
      ;;
    *)
      echo "Usage: $0 [--case path] [--output path] [--all]" >&2
      exit 2
      ;;
  esac
done

json_bool() {
  if [[ "$1" == "true" ]]; then
    printf 'true'
  else
    printf 'false'
  fi
}

validate_schema() {
  local schema="$1"
  local data="$2"
  node scripts/validate-json-schema.mjs "$schema" "$data" >/dev/null
}

verify_source_commit_artifacts() {
  local manifest_path="$1"
  local source_ref="$2"
  local expected_ref="$3"
  local tmpdir="$4"

  if [[ "$source_ref" != "$expected_ref" || ! "$source_ref" =~ ^[0-9a-f]{40}$ ]]; then
    printf 'mismatch|0|0|0|source_commit does not match expected immutable reference or is not a 40-character SHA'
    return
  fi

  if ! git cat-file -e "$source_ref^{commit}" 2>/dev/null; then
    printf 'mismatch|0|0|0|source_commit is not available in this git repository'
    return
  fi

  local checked=0 missing=0 mismatched=0
  while IFS=$'\t' read -r artifact_path expected_sha; do
    [[ -z "$artifact_path" ]] && continue
    if ! git cat-file -e "$source_ref:$artifact_path" 2>/dev/null; then
      missing=$((missing + 1))
      continue
    fi

    local source_file="$tmpdir/source-artifact-$checked"
    git show "$source_ref:$artifact_path" >"$source_file"
    local actual_sha
    actual_sha="$(shasum -a 256 "$source_file" | awk '{print $1}')"
    if [[ "$actual_sha" != "$expected_sha" ]]; then
      mismatched=$((mismatched + 1))
    fi
    checked=$((checked + 1))
  done < <(jq -r '.artifacts[] | [.path, .sha256] | @tsv' "$manifest_path")

  if [[ "$mismatched" -gt 0 || "$missing" -gt 0 ]]; then
    printf 'mismatch|%s|%s|%s|source_commit verification found missing or mismatched manifest artifacts; this proof cannot claim source-content attestation' "$checked" "$missing" "$mismatched"
  elif [[ "$checked" -eq 0 ]]; then
    printf 'not_assessed|%s|%s|%s|manifest contains no source artifacts to compare against source_commit' "$checked" "$missing" "$mismatched"
  else
    printf 'matched|%s|%s|%s|source_commit contains every manifest artifact path with matching digest' "$checked" "$missing" "$mismatched"
  fi
}

verify_case() {
  local case_path="$1"
  local tmpdir
  tmpdir="$(mktemp -d)"

  local id manifest envelope public_key policy expected_source expected_signer reference_time
  local verification_artifact verification_sha require_external expected_negative negative_reason

  id="$(jq -r '.id' "$case_path")"
  manifest="$(jq -r '.manifest_ref' "$case_path")"
  envelope="$(jq -r '.dsse_envelope_ref' "$case_path")"
  public_key="$(jq -r '.public_key_ref' "$case_path")"
  policy="$(jq -r '.trusted_identity_policy_ref' "$case_path")"
  expected_source="$(jq -r '.expected_source_commit' "$case_path")"
  expected_signer="$(jq -r '.expected_signer_identity' "$case_path")"
  reference_time="$(jq -r '.reference_time' "$case_path")"
  verification_artifact="$(jq -r '.verification_artifact_ref' "$case_path")"
  verification_sha="$(jq -r '.verification_artifact_sha256' "$case_path")"
  require_external="$(jq -r '.require_external_attestation' "$case_path")"
  expected_negative="$(jq -r '.expected_negative // false' "$case_path")"
  negative_reason="$(jq -r '.negative_reason // ""' "$case_path")"

  local schema_valid=true
  validate_schema schema/self-attestation-case.schema.json "$case_path" || schema_valid=false
  validate_schema schema/contract-manifest.schema.json "$manifest" || schema_valid=false
  jq empty "$case_path" "$envelope" "$policy" "$verification_artifact" >/dev/null || schema_valid=false

  local manifest_digest manifest_digest_signed manifest_subject manifest_digest_status artifact_digest_status
  manifest_digest="$(shasum -a 256 "$manifest" | awk '{print $1}')"
  if scripts/verify-contract-manifest.sh "$manifest" >/dev/null 2>&1; then
    artifact_digest_status="matched"
  else
    artifact_digest_status="mismatch"
  fi

  local signature_valid=true
  local payload_type payload_len type_len key_fingerprint signer_identity source_commit valid_until freshness_status
  local predicate_policy_ref predicate_profile policy_profile policy_expected_signer identity_policy_status
  payload_type="$(jq -r '.payloadType' "$envelope")"
  jq -r '.payload' "$envelope" | base64 --decode >"$tmpdir/payload.json"
  jq -r '.signatures[0].sig' "$envelope" | base64 --decode >"$tmpdir/signature.bin"
  payload_len="$(wc -c <"$tmpdir/payload.json" | tr -d ' ')"
  type_len="$(printf '%s' "$payload_type" | wc -c | tr -d ' ')"
  {
    printf 'DSSEv1 %s %s %s ' "$type_len" "$payload_type" "$payload_len"
    cat "$tmpdir/payload.json"
  } >"$tmpdir/pae.bin"
  openssl dgst -sha256 -verify "$public_key" -signature "$tmpdir/signature.bin" "$tmpdir/pae.bin" >/dev/null 2>&1 || signature_valid=false

  manifest_subject="$(jq -r '.subject[0].name' "$tmpdir/payload.json")"
  manifest_digest_signed="$(jq -r '.subject[0].digest.sha256' "$tmpdir/payload.json")"
  if [[ "$manifest_subject" == "$manifest" && "$manifest_digest_signed" == "$manifest_digest" ]]; then
    manifest_digest_status="matched"
  else
    manifest_digest_status="mismatch"
  fi

  key_fingerprint="$(openssl pkey -pubin -in "$public_key" -outform DER | openssl dgst -sha256 | awk '{print $NF}')"
  signer_identity="$(jq -r '.predicate.signer_identity' "$tmpdir/payload.json")"
  predicate_policy_ref="$(jq -r '.predicate.identity_policy_ref // empty' "$tmpdir/payload.json")"
  predicate_profile="$(jq -r '.predicate.private_equivalent_profile // empty' "$tmpdir/payload.json")"
  policy_profile="$(jq -r '.allowed_private_equivalent_profile.profile // empty' "$policy")"
  policy_expected_signer="$policy_profile:$key_fingerprint"
  source_commit="$(jq -r '.source_commit' "$manifest")"
  valid_until="$(jq -r '.valid_until // empty' "$manifest")"

  local source_status signer_status verification_artifact_status signature_status source_check source_checked source_missing source_mismatched source_reason
  source_check="$(verify_source_commit_artifacts "$manifest" "$source_commit" "$expected_source" "$tmpdir")"
  IFS='|' read -r source_status source_checked source_missing source_mismatched source_reason <<<"$source_check"

  if [[ -n "$policy_profile" && "$predicate_policy_ref" == "$policy" && "$predicate_profile" == "$policy_profile" ]]; then
    identity_policy_status="matched"
  else
    identity_policy_status="mismatch"
  fi

  if [[ "$identity_policy_status" == "matched" && "$signer_identity" == "$policy_expected_signer" ]]; then
    signer_status="matched"
  else
    signer_status="mismatch"
  fi

  if [[ "$signature_valid" == "true" ]]; then
    signature_status="valid"
  else
    signature_status="invalid"
  fi

  if [[ -n "$valid_until" && "$valid_until" > "$reference_time" ]]; then
    freshness_status="current"
  else
    freshness_status="expired"
  fi

  local actual_verification_sha
  actual_verification_sha="$(shasum -a 256 "$verification_artifact" | awk '{print $1}')"
  if [[ "$actual_verification_sha" == "$verification_sha" ]]; then
    verification_artifact_status="matched"
  else
    verification_artifact_status="mismatch"
  fi

  local digest_verified locally_attested externally_attested production_verified negative_detected
  if [[ "$manifest_digest_status" == "matched" && "$artifact_digest_status" == "matched" && "$verification_artifact_status" == "matched" ]]; then
    digest_verified=true
  else
    digest_verified=false
  fi

  if [[ "$digest_verified" == "true" && "$signature_status" == "valid" && "$signer_status" == "matched" && "$source_status" == "matched" && "$freshness_status" == "current" ]]; then
    locally_attested=true
  else
    locally_attested=false
  fi

  externally_attested=false
  production_verified=false

  negative_detected=false
  if [[ "$expected_negative" == "true" ]]; then
    case "$negative_reason" in
      wrong_source_commit)
        [[ "$source_status" == "mismatch" ]] && negative_detected=true
        ;;
      wrong_signer)
        [[ "$signature_status" == "invalid" || "$signer_status" == "mismatch" ]] && negative_detected=true
        ;;
      wrong_policy)
        [[ "$identity_policy_status" == "mismatch" || "$signer_status" == "mismatch" ]] && negative_detected=true
        ;;
      stale_manifest)
        [[ "$freshness_status" == "expired" ]] && negative_detected=true
        ;;
      missing_external_attestation)
        [[ "$require_external" == "true" && "$externally_attested" == "false" && "$production_verified" == "false" ]] && negative_detected=true
        ;;
      modified_verification_artifact)
        [[ "$verification_artifact_status" == "mismatch" ]] && negative_detected=true
        ;;
    esac
  fi

  jq -n \
    --arg id "$id" \
    --arg generated_at "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
    --arg manifest_ref "$manifest" \
    --arg manifest_digest "$manifest_digest" \
    --arg manifest_digest_status "$manifest_digest_status" \
    --arg artifact_digest_status "$artifact_digest_status" \
    --arg dsse_envelope_ref "$envelope" \
    --arg signature_status "$signature_status" \
    --arg signer_identity "$signer_identity" \
    --arg expected_signer_identity "$expected_signer" \
    --arg signer_identity_status "$signer_status" \
    --arg policy_expected_signer_identity "$policy_expected_signer" \
    --arg identity_policy_ref "$policy" \
    --arg identity_policy_status "$identity_policy_status" \
    --arg source_commit "$source_commit" \
    --arg expected_source_commit "$expected_source" \
    --arg source_commit_status "$source_status" \
    --arg source_commit_artifact_status "$source_status" \
    --arg source_commit_artifact_reason "$source_reason" \
    --argjson source_commit_checked_count "$source_checked" \
    --argjson source_commit_missing_count "$source_missing" \
    --argjson source_commit_mismatched_count "$source_mismatched" \
    --arg freshness_status "$freshness_status" \
    --arg verification_artifact_ref "$verification_artifact" \
    --arg verification_artifact_status "$verification_artifact_status" \
    --arg external_reason "No Sigstore/Rekor bundle or customer PKI attestation is committed for this local development proof." \
    --arg production_reason "Production release verification requires external attestation in addition to local digest and DSSE checks." \
    --arg negative_reason "$negative_reason" \
    --argjson schema_valid_value "$(json_bool "$schema_valid")" \
    --argjson digest_verified_value "$(json_bool "$digest_verified")" \
    --argjson locally_attested_value "$(json_bool "$locally_attested")" \
    --argjson externally_attested_value false \
    --argjson production_verified_value false \
    --argjson expected_negative_value "$(json_bool "$expected_negative")" \
    --argjson negative_detected_value "$(json_bool "$negative_detected")" \
    '{
      id: $id,
      schema_version: "0.1.0",
      generated_at: $generated_at,
      manifest_ref: $manifest_ref,
      manifest_digest: $manifest_digest,
      manifest_digest_status: $manifest_digest_status,
      artifact_digest_status: $artifact_digest_status,
      dsse_envelope_ref: $dsse_envelope_ref,
      signature_status: $signature_status,
      signer_identity: $signer_identity,
      expected_signer_identity: $expected_signer_identity,
      policy_expected_signer_identity: $policy_expected_signer_identity,
      signer_identity_status: $signer_identity_status,
      identity_policy_ref: $identity_policy_ref,
      identity_policy_status: $identity_policy_status,
      source_commit: $source_commit,
      expected_source_commit: $expected_source_commit,
      source_commit_status: $source_commit_status,
      source_commit_artifact_status: $source_commit_artifact_status,
      source_commit_artifact_counts: {
        checked: $source_commit_checked_count,
        missing: $source_commit_missing_count,
        mismatched: $source_commit_mismatched_count
      },
      external_trust_profile: "not_assessed",
      external_attestation_ref: null,
      transparency_evidence_ref: null,
      source_uri_status: "not_assessed",
      protected_ref_status: "not_assessed",
      workflow_identity_status: "not_assessed",
      approval_status: "not_assessed",
      production_release_verified: {
        state: "not_assessed",
        value: null,
        reason: $production_reason
      },
      freshness_status: $freshness_status,
      verification_artifact_ref: $verification_artifact_ref,
      verification_artifact_status: $verification_artifact_status,
      proof_states: {
        schema_valid: {
          state: "assessed",
          value: $schema_valid_value
        },
        digest_verified: {
          state: "assessed",
          value: $digest_verified_value
        },
        locally_attested: {
          state: "assessed",
          value: $locally_attested_value,
          reason: (if $locally_attested_value then null else $source_commit_artifact_reason end)
        },
        source_commit_artifacts_verified: {
          state: (if $source_commit_artifact_status == "matched" or $source_commit_artifact_status == "mismatch" then "assessed" else "not_assessed" end),
          value: (if $source_commit_artifact_status == "matched" then true elif $source_commit_artifact_status == "mismatch" then false else null end),
          reason: $source_commit_artifact_reason
        },
        externally_attested: {
          state: "not_assessed",
          value: null,
          reason: $external_reason
        },
        production_release_verified: {
          state: "not_assessed",
          value: null,
          reason: $production_reason
        }
      },
      trusted_contract_release: false,
      expected_negative: $expected_negative_value,
      negative_reason: $negative_reason,
      negative_detected: $negative_detected_value
    }'

  rm -rf "$tmpdir"
}

run_one() {
  local result
  result="$(verify_case "$1")"
  if [[ -n "$output_file" ]]; then
    printf '%s\n' "$result" >"$output_file"
  else
    printf '%s\n' "$result"
  fi

  local expected_negative negative_detected schema_valid digest_verified signature_status identity_policy_status source_commit_status verification_artifact_status
  expected_negative="$(jq -r '.expected_negative' <<<"$result")"
  negative_detected="$(jq -r '.negative_detected' <<<"$result")"
  schema_valid="$(jq -r '.proof_states.schema_valid.value' <<<"$result")"
  digest_verified="$(jq -r '.proof_states.digest_verified.value' <<<"$result")"
  signature_status="$(jq -r '.signature_status' <<<"$result")"
  identity_policy_status="$(jq -r '.identity_policy_status' <<<"$result")"
  source_commit_status="$(jq -r '.source_commit_status' <<<"$result")"
  verification_artifact_status="$(jq -r '.verification_artifact_status' <<<"$result")"

  if [[ "$expected_negative" == "true" && "$negative_detected" != "true" ]]; then
    echo "Expected negative self-attestation case was not detected: $1" >&2
    return 1
  fi

  if [[ "$expected_negative" != "true" ]]; then
    expected_digest_verified="$(jq -c 'if (.expected_proof_states | has("digest_verified")) then .expected_proof_states.digest_verified else true end' "$1")"
    if [[ "$schema_valid" != "true" || "$signature_status" != "valid" || "$identity_policy_status" != "matched" || "$verification_artifact_status" != "matched" ]]; then
      echo "Positive self-attestation case failed required structural proof state: $1" >&2
      return 1
    fi
    if [[ "$digest_verified" != "$expected_digest_verified" ]]; then
      echo "Self-attestation digest proof state mismatch for $1: expected $expected_digest_verified actual $digest_verified" >&2
      return 1
    fi
    expected_source_artifacts="$(jq -c 'if (.expected_proof_states | has("source_commit_artifacts_verified")) then .expected_proof_states.source_commit_artifacts_verified else empty end' "$1")"
    if [[ "$source_commit_status" == "mismatch" && "$expected_source_artifacts" != "false" ]]; then
      echo "Positive self-attestation case has a mismatched source commit: $1" >&2
      return 1
    fi
  fi

  local expected_state key expected actual
  while IFS= read -r expected_state; do
    key="$(jq -r '.key' <<<"$expected_state")"
    expected="$(jq -c '.value' <<<"$expected_state")"
    actual="$(jq -c --arg key "$key" 'if $key == "trusted_contract_release" then .trusted_contract_release else .proof_states[$key].value end' <<<"$result")"
    if [[ "$actual" != "$expected" ]]; then
      echo "Self-attestation proof state mismatch for $1: $key expected $expected actual $actual" >&2
      return 1
    fi
  done < <(jq -c '.expected_proof_states // {} | to_entries[]' "$1")
}

if [[ "$run_negative_suite" -eq 1 ]]; then
  run_one "$DEFAULT_CASE" >/dev/null
  for negative_case in examples/self-trace/self-attestation-negative-*.json; do
    run_one "$negative_case" >/dev/null
  done
else
  run_one "$case_file"
fi
