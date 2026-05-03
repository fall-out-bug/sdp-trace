#!/usr/bin/env bash
set -uo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

profile="baseline"
output="text"
version="0.1.0"
allow_dirty=false
artifact_role="verifier_output"

usage() {
  cat >&2 <<'USAGE'
Usage: scripts/verify.sh [--profile baseline|source-bound|external-trust] [--json] [--allow-dirty] [--example] [--version]

Exit codes:
  0  selected profile passed
  1  selected profile failed
  2  usage error
  3  required checks could not verify

Dirty checkouts without --allow-dirty exit 3. With --allow-dirty, output is limited to
local structural development and cannot support source-bound or external trust.
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --profile)
      if [[ $# -lt 2 ]]; then
        usage
        exit 2
      fi
      profile="$2"
      shift 2
      ;;
    --json)
      output="json"
      shift
      ;;
    --allow-dirty)
      allow_dirty=true
      shift
      ;;
    --example)
      artifact_role="untrusted_shape_example"
      shift
      ;;
    --version)
      echo "$version"
      exit 0
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unsupported argument: $1" >&2
      usage
      exit 2
      ;;
  esac
done

case "$profile" in
  baseline|repo_baseline)
    profile_id="repo_baseline_structural"
    ;;
  source-bound|source_bound|source_bound_local_release)
    profile_id="source_bound_local_release"
    ;;
  external-trust|external_trust|external_production_trust)
    profile_id="external_production_trust"
    ;;
  *)
    echo "Unsupported profile: $profile" >&2
    exit 2
    ;;
esac

require_tool() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "$1 is required but not found" >&2
    exit 3
  fi
}

require_tool git
require_tool jq
require_tool node
if ! command -v sha256sum >/dev/null 2>&1 && ! command -v shasum >/dev/null 2>&1; then
  echo "sha256sum or shasum is required but not found" >&2
  exit 3
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT
states_file="$tmp_dir/states.jsonl"
: >"$states_file"
zero_sha="0000000000000000000000000000000000000000000000000000000000000000"
zero_commit="0000000000000000000000000000000000000000"

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

add_state() {
  local id="$1"
  local result="$2"
  local required="$3"
  local method="$4"
  local reason="${5:-}"
  local next_required_evidence="${6:-}"

  jq -cn \
    --arg id "$id" \
    --arg result "$result" \
    --argjson required "$required" \
    --arg method "$method" \
    --arg reason "$reason" \
    --arg next_required_evidence "$next_required_evidence" '
      {
        id: $id,
        result: $result,
        required: $required,
        method: $method
      }
      + (if $reason != "" then {reason: $reason} else {} end)
      + (if $next_required_evidence != "" then {next_required_evidence: $next_required_evidence} else {} end)
    ' >>"$states_file"
}

run_state() {
  local id="$1"
  local method="$2"
  shift 2
  local out="$tmp_dir/$id.out"
  set +e
  "$@" >"$out" 2>&1
  local status="$?"
  if [[ "$status" -eq 0 ]]; then
    add_state "$id" "pass" "true" "$method"
  elif [[ "$status" -eq 127 ]]; then
    add_state "$id" "cannot_verify" "true" "$method" "required command was unavailable"
  else
    local first_line
    first_line="$(sed -n '1p' "$out" | tr -d '\r')"
    if [[ -z "$first_line" ]]; then
      first_line="command exited $status"
    fi
    add_state "$id" "fail" "true" "$method" "$first_line"
  fi
}

validate_schema_examples() {
  node scripts/validate-json-schema.mjs schema/trace.schema.json examples/github-speckit/trace.json || return 1
  node scripts/validate-json-schema.mjs schema/evidence-bundle.schema.json examples/go-service/evidence-bundle.json || return 1
  node scripts/validate-json-schema.mjs schema/gate-verdict.schema.json examples/go-service/gate-verdict.json || return 1
  node scripts/validate-json-schema.mjs schema/evidence-bundle.schema.json examples/jvm-bazel/evidence-bundle.json || return 1
  node scripts/validate-json-schema.mjs schema/assessment-input.schema.json examples/contract-foundation/positive-assessment-input.json || return 1
  node scripts/validate-json-schema.mjs schema/assessment-input.schema.json examples/contract-foundation/not-assessed-assessment-input.json || return 1
  node scripts/validate-json-schema.mjs schema/contract-manifest.schema.json examples/contract-foundation/contract-manifest.example.json || return 1
  node scripts/validate-json-schema.mjs schema/contract-release-verification.schema.json examples/contract-foundation/contract-release-verification.example.json || return 1
  node scripts/validate-json-schema.mjs schema/trusted-identity-policy.schema.json examples/contract-foundation/trusted-identity-policy.example.json || return 1
  node scripts/validate-json-schema.mjs schema/trusted-identity-policy.schema.json examples/contract-foundation/trusted-identity-policy-wrong-profile.example.json || return 1
  node scripts/validate-json-schema.mjs schema/consumer-schema-version-declaration.schema.json examples/contract-foundation/sdp-gate-consumer-declaration.example.json || return 1
  if [[ -f examples/self-trace/proof-summary.example.json ]]; then
    node scripts/validate-json-schema.mjs schema/proof-summary.schema.json examples/self-trace/proof-summary.example.json || return 1
  fi
  node scripts/validate-pilot-matrices.mjs || return 1
}

expect_schema_failure() {
  local schema="$1"
  local data="$2"
  if node scripts/validate-json-schema.mjs "$schema" "$data" >/dev/null 2>&1; then
    echo "Expected validation failure, but passed: $data" >&2
    return 1
  fi
}

validate_negative_fixtures() {
  expect_schema_failure schema/assessment-input.schema.json examples/contract-foundation/negative-native-policy-field.json || return 1
  expect_schema_failure schema/trace.schema.json examples/contract-foundation/negative-native-policy-trace.json || return 1
  expect_schema_failure schema/metric-stream.schema.json examples/contract-foundation/negative-assessed-metric-without-evidence.json || return 1
  expect_schema_failure schema/accountability.schema.json examples/contract-foundation/negative-ai-sole-accountable-owner.json || return 1
  expect_schema_failure schema/contract-release-verification.schema.json examples/contract-foundation/negative-unauthorized-signer.json || return 1
  expect_schema_failure schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-missing-external-evidence.json || return 1
  expect_schema_failure schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-oidc-issuer-mismatch.json || return 1
  expect_schema_failure schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-source-uri-mismatch.json || return 1
  expect_schema_failure schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-protected-ref-mismatch.json || return 1
  expect_schema_failure schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-workflow-identity-mismatch.json || return 1
  expect_schema_failure schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-approval-mismatch.json || return 1
  expect_schema_failure schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-expired-freshness.json || return 1
  expect_schema_failure schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-local-profile.json || return 1
  expect_schema_failure schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-mutable-source-commit.json || return 1
  expect_schema_failure schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-bad-source-counts.json || return 1
  expect_schema_failure schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-sigstore-not-required-transparency.json || return 1
}

spec_path="specs/001-sdp-trace-time-series-evidence-substrate/blocks/07-minimum-trust-kernel.md"
gate_files=(
  "scripts/verify.sh"
  "schema/proof-summary.schema.json"
  "scripts/validate-json-schema.mjs"
  "scripts/validate-claim-consistency.mjs"
  "scripts/validate-self-trace.sh"
  "scripts/validate-cross-references.mjs"
  "scripts/check-artifact-safety.sh"
  "scripts/validate-pilot-matrices.mjs"
  "scripts/verify-self-attestation.sh"
  "schema/contract-release-verification.schema.json"
  "schema/self-attestation-case.schema.json"
)

dirty=false
if [[ -n "$(git status --porcelain)" ]]; then
  dirty=true
  if [[ "$allow_dirty" == "true" && "$profile_id" == "repo_baseline_structural" ]]; then
    add_state "source_checkout_clean" "not_assessed" "false" "git status --porcelain" "dirty checkout allowed by explicit --allow-dirty; this output cannot support source-bound or external trust"
  else
    add_state "source_checkout_clean" "cannot_verify" "true" "git status --porcelain" "working tree is dirty; rerun from a clean checkout or use --allow-dirty for local structural development only"
  fi
else
  add_state "source_checkout_clean" "pass" "true" "git status --porcelain"
fi

if [[ -f "$spec_path" ]]; then
  add_state "spec_subject_present" "pass" "true" "test -f $spec_path"
else
  add_state "spec_subject_present" "cannot_verify" "true" "test -f $spec_path" "required spec subject is missing"
fi

missing_gate_files=()
for gate_file in "${gate_files[@]}"; do
  if [[ ! -f "$gate_file" ]]; then
    missing_gate_files+=("$gate_file")
  fi
done
if [[ "${#missing_gate_files[@]}" -eq 0 ]]; then
  add_state "gate_set_files_present" "pass" "true" "test -f required gate-set files"
else
  missing_gate_files_text="$(printf '%s ' "${missing_gate_files[@]}")"
  add_state "gate_set_files_present" "cannot_verify" "true" "test -f required gate-set files" "missing gate-set files: ${missing_gate_files_text% }"
fi

run_state "schema_json_valid" "jq empty schema/*.json" jq empty schema/*.json
run_state "committed_examples_parse" "jq empty for committed examples" bash -c "find examples -name '*.json' -not -path '*/.git/*' -not -path '*/.beads/*' -not -path '*/.sdp-trace-runs/*' -not -path '*/benchmarks/repos/*' -not -path '*/node_modules/*' -print0 | xargs -0 jq empty"
run_state "portable_schema_examples_valid" "validate portable schema/example pairs" validate_schema_examples
run_state "negative_fixtures_rejected" "execute negative schema fixtures" validate_negative_fixtures
run_state "self_trace_package_valid" "scripts/validate-self-trace.sh" scripts/validate-self-trace.sh
run_state "artifact_safety_valid" "scripts/check-artifact-safety.sh" scripts/check-artifact-safety.sh
run_state "claim_tags_consistent" "scripts/validate-claim-consistency.mjs" scripts/validate-claim-consistency.mjs
run_state "cross_reference_integrity" "scripts/validate-cross-references.mjs" scripts/validate-cross-references.mjs

if [[ "$profile_id" == "source_bound_local_release" || "$profile_id" == "external_production_trust" ]]; then
  self_attestation_json="$tmp_dir/self-attestation.json"
  set +e
  scripts/verify-self-attestation.sh --case examples/self-trace/self-attestation-verification.json >"$self_attestation_json" 2>"$tmp_dir/self-attestation.err"
  self_attestation_status="$?"
  set -e
  if [[ "$self_attestation_status" -ne 0 ]]; then
    add_state "source_bound_self_attestation_ran" "cannot_verify" "true" "scripts/verify-self-attestation.sh --case examples/self-trace/self-attestation-verification.json" "$(sed -n '1p' "$tmp_dir/self-attestation.err")"
  else
    add_state "source_bound_self_attestation_ran" "pass" "true" "scripts/verify-self-attestation.sh --case examples/self-trace/self-attestation-verification.json"
    manifest_digest_status="$(jq -r '.manifest_digest_status' "$self_attestation_json")"
    artifact_digest_status="$(jq -r '.artifact_digest_status' "$self_attestation_json")"
    source_commit_artifact_value="$(jq -r '.proof_states.source_commit_artifacts_verified.value' "$self_attestation_json")"
    locally_attested_value="$(jq -r '.proof_states.locally_attested.value' "$self_attestation_json")"
    source_reason="$(jq -r '.proof_states.source_commit_artifacts_verified.reason // empty' "$self_attestation_json")"
    if [[ "$manifest_digest_status" == "matched" ]]; then
      add_state "manifest_digest_verified" "pass" "true" "verify manifest digest bound by DSSE subject"
    else
      add_state "manifest_digest_verified" "fail" "true" "verify manifest digest bound by DSSE subject" "manifest digest status is $manifest_digest_status"
    fi
    if [[ "$artifact_digest_status" == "matched" ]]; then
      add_state "artifact_subject_verified" "pass" "true" "scripts/verify-contract-manifest.sh examples/contract-foundation/contract-manifest.example.json"
    else
      add_state "artifact_subject_verified" "fail" "true" "scripts/verify-contract-manifest.sh examples/contract-foundation/contract-manifest.example.json" "artifact digest status is $artifact_digest_status"
    fi
    if [[ "$source_commit_artifact_value" == "true" ]]; then
      add_state "source_subject_verified" "pass" "true" "verify manifest artifact set at source commit"
    else
      add_state "source_subject_verified" "fail" "true" "verify manifest artifact set at source commit" "${source_reason:-source commit artifacts were not verified}"
    fi
    if [[ "$locally_attested_value" == "true" ]]; then
      add_state "local_envelope_verified" "pass" "true" "verify local DSSE envelope and identity policy"
    else
      add_state "local_envelope_verified" "fail" "true" "verify local DSSE envelope and identity policy" "local attestation is false because required source-bound states did not pass"
    fi
  fi
fi

if [[ "$profile_id" == "external_production_trust" ]]; then
  add_state "external_trust_profile_selected" "fail" "true" "select supported external trust profile" "no supported external trust profile evidence is committed for this repository self-release"
  add_state "external_attestation_present" "not_assessed" "true" "verify Sigstore/Rekor or customer PKI attestation" "No Sigstore/Rekor bundle or accepted customer PKI evidence is committed" "Commit real external attestation evidence for sigstore-rekor-keyless-v1 or customer-pki-audit-v1"
  add_state "external_identity_policy_matched" "not_assessed" "true" "verify external signer identity policy" "external identity policy cannot be assessed without external attestation evidence" "Commit external signer identity and policy evidence"
  add_state "transparency_or_audit_verified" "not_assessed" "true" "verify transparency log or customer audit evidence" "no transparency or audit evidence is committed" "Commit Rekor inclusion proof or accepted customer audit evidence"
  add_state "production_release_verified" "fail" "true" "derive production release verification from external trust states" "production release verification requires every external production trust state to pass"
fi

cannot_verify_count="$(jq -s '[.[] | select(.required == true and .result == "cannot_verify")] | length' "$states_file")"
failures_count="$(jq -s '[.[] | select(.required == true and (.result == "fail" or .result == "not_assessed" or .result == "incomplete"))] | length' "$states_file")"
if [[ "$cannot_verify_count" -gt 0 ]]; then
  result="cannot_verify"
elif [[ "$failures_count" -gt 0 ]]; then
  result="fail"
else
  result="pass"
fi

commit="$(git rev-parse HEAD)"
trust_scope="clean_checkout_structural"
if [[ "$dirty" == "true" ]]; then
  if [[ "$allow_dirty" == "true" && "$profile_id" == "repo_baseline_structural" ]]; then
    trust_scope="local_dirty_structural_only"
  else
    trust_scope="none_dirty_checkout"
  fi
elif [[ "$profile_id" == "source_bound_local_release" ]]; then
  trust_scope="source_bound_local_release"
elif [[ "$profile_id" == "external_production_trust" ]]; then
  trust_scope="external_production_trust"
fi
source_subject_id="source:${commit}"
if [[ -f "$spec_path" ]]; then
  spec_digest="$(sha256_file "$spec_path")"
else
  spec_digest="$zero_sha"
fi

gate_files_json="$tmp_dir/gate-files.jsonl"
: >"$gate_files_json"
for gate_file in "${gate_files[@]}"; do
  if [[ -f "$gate_file" ]]; then
    jq -cn --arg path "$gate_file" --arg sha "$(sha256_file "$gate_file")" '{path: $path, sha256: $sha}' >>"$gate_files_json"
  fi
done
gate_set_sorted_json="$tmp_dir/gate-set-files.json"
jq -s -c 'sort_by(.path)' "$gate_files_json" >"$gate_set_sorted_json"
gate_set_digest="$(sha256_file "$gate_set_sorted_json")"
generated_at="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

if [[ "$artifact_role" == "untrusted_shape_example" ]]; then
  commit="$zero_commit"
  source_subject_id="source:${commit}"
  trust_scope="untrusted_shape_only"
  spec_digest="$zero_sha"
  gate_set_digest="$zero_sha"
  generated_at="1970-01-01T00:00:00Z"
  : >"$gate_files_json"
  for gate_file in "${gate_files[@]}"; do
    if [[ -f "$gate_file" ]]; then
      jq -cn --arg path "$gate_file" --arg sha "$zero_sha" '{path: $path, sha256: $sha}' >>"$gate_files_json"
    fi
  done
fi

invocation_parts=("scripts/verify.sh" "--profile" "$profile")
if [[ "$output" == "json" ]]; then
  invocation_parts+=("--json")
fi
if [[ "$allow_dirty" == "true" ]]; then
  invocation_parts+=("--allow-dirty")
fi
if [[ "$artifact_role" == "untrusted_shape_example" ]]; then
  invocation_parts+=("--example")
fi
invocation="${invocation_parts[*]}"

summary_json="$(
  jq -n \
    --slurpfile states "$states_file" \
    --slurpfile gate_files "$gate_files_json" \
    --arg id "proof-summary-repo-baseline" \
    --arg schema_version "0.1.0" \
    --arg profile "$profile_id" \
    --arg result "$result" \
    --arg trust_scope "$trust_scope" \
    --arg artifact_role "$artifact_role" \
    --arg generated_at "$generated_at" \
    --arg version "$version" \
    --arg invocation "$invocation" \
    --arg source_subject_id "$source_subject_id" \
    --arg gate_set_digest "$gate_set_digest" \
    --arg commit "$commit" \
    --argjson dirty "$dirty" \
    --argjson dirty_allowed "$allow_dirty" \
    --arg spec_path "$spec_path" \
    --arg spec_digest "$spec_digest" '
      {
        id: $id,
        schema_version: $schema_version,
        artifact_role: $artifact_role,
        profile: $profile,
        result: $result,
        trust_scope: $trust_scope,
        generated_at: $generated_at,
        generated_by: {
          tool: "scripts/verify.sh",
          version: $version,
          invocation: $invocation,
          source_subject_ref: $source_subject_id,
          gate_set_digest: $gate_set_digest
        },
        source_subject: {
          id: $source_subject_id,
          type: "git_commit_v1",
          commit: $commit,
          dirty: $dirty,
          dirty_allowed: $dirty_allowed
        },
        spec_subject: {
          id: "spec:block-07-minimum-trust-kernel",
          path: $spec_path,
          digest_algorithm: "sha256",
          digest: $spec_digest
        },
        gate_set_subject: {
          id: "gate-set:repo-baseline-v0",
          digest_algorithm: "sha256",
          digest: $gate_set_digest,
          files: ($gate_files | sort_by(.path))
        },
        claims: [
          (if $profile == "repo_baseline_structural" then "repo_baseline structural checks executed only" else "repo_baseline structural checks included" end),
          "external production trust not claimed",
          "authoritative claim tags checked by Slice 1 validator",
          "cross-reference integrity checked by Slice 5 validator"
        ],
        states: $states,
        failures: ($states | map(select(.required == true and (.result == "fail" or .result == "cannot_verify" or .result == "not_assessed" or .result == "incomplete")))),
        not_assessed: ($states | map(select(.result == "not_assessed")))
      }
    '
)"

if [[ "$output" == "json" ]]; then
  printf '%s\n' "$summary_json"
else
  echo "sdp-trace verifier"
  echo "profile: $profile_id"
  echo "result: $result"
  echo
  echo "states:"
  jq -r '.states[] | "  \(.result) \(.id)"' <<<"$summary_json"
fi

case "$result" in
  pass)
    exit 0
    ;;
  cannot_verify)
    exit 3
    ;;
  *)
    exit 1
    ;;
esac
