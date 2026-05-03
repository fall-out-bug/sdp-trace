#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

baseline_json="$tmp_dir/baseline-proof-summary.json"
baseline_text="$tmp_dir/baseline.out"
usage_out="$tmp_dir/usage.out"
dirty_out="$tmp_dir/dirty.out"
unsupported_claim_fixture="$tmp_dir/unsupported-claim.md"
prose_fixture="$tmp_dir/non-authoritative-prose.md"
stale_claim_fixture="$tmp_dir/stale-claim.md"
t070_false_pass_fixture="$tmp_dir/t070-false-pass.md"
unverified_evidence_fixture="$tmp_dir/unverified-evidence.md"
none_evidence_fixture="$tmp_dir/none-evidence.md"
contract_release_verification="examples/contract-foundation/contract-release-verification.example.json"
self_attestation_case="examples/self-trace/self-attestation-verification.json"
block06_report="docs/research/opencode-minimax-kotlin-bazel-proof-report.md"
block06_proof_states="examples/pilot-runs/opencode-minimax-kotlin-bazel/evidence/proof-states.json"
source_bound_json="$tmp_dir/source-bound.json"
external_trust_json="$tmp_dir/external-trust.json"
cross_ref_fixture="$tmp_dir/self-trace-cross-ref-fixture"

printf '%s\n' '<!-- sdp-trace-claim: claim=unsupported; subject=T999; state=pass; profile=repo_baseline; evidence=none -->' >"$unsupported_claim_fixture"
set +e
scripts/validate-claim-consistency.mjs --files "$unsupported_claim_fixture" >"$tmp_dir/unsupported-claim.out" 2>&1
unsupported_claim_status="$?"
set -e
if [[ "$unsupported_claim_status" -eq 0 ]]; then
  echo "Expected unsupported claim tag to fail" >&2
  exit 1
fi
grep -q "unsupported claim value" "$tmp_dir/unsupported-claim.out"

printf '%s\n' 'A reviewer said: "T070 falsely claims validation was verified and complete."' >"$prose_fixture"
scripts/validate-claim-consistency.mjs --files "$prose_fixture" >"$tmp_dir/prose.out"
grep -q "claim tags valid" "$tmp_dir/prose.out"

printf '%s\n' '<!-- sdp-trace-claim: claim=task_closed; subject=T070; state=stale; profile=repo_baseline; evidence=state:claim_tags_consistent -->' >"$stale_claim_fixture"
scripts/validate-claim-consistency.mjs --files "$stale_claim_fixture" >"$tmp_dir/stale.out"
grep -q "claim tags valid" "$tmp_dir/stale.out"

printf '%s\n' '<!-- sdp-trace-claim: claim=task_closed; subject=T070; state=pass; profile=repo_baseline; evidence=command_set:block04-t070 -->' >"$t070_false_pass_fixture"
set +e
scripts/validate-claim-consistency.mjs --files "$t070_false_pass_fixture" >"$tmp_dir/t070-false-pass.out" 2>&1
t070_false_pass_status="$?"
set -e
if [[ "$t070_false_pass_status" -ne 0 ]]; then
  echo "Expected T070 pass claim fixture to pass after source-bound local release repair" >&2
  cat "$tmp_dir/t070-false-pass.out" >&2
  exit 1
fi
grep -q "claim tags valid" "$tmp_dir/t070-false-pass.out"

printf '%s\n' '<!-- sdp-trace-claim: claim=task_closed; subject=T071; state=stale; profile=repo_baseline; evidence=proof:unverified-artifact -->' >"$unverified_evidence_fixture"
set +e
scripts/validate-claim-consistency.mjs --files "$unverified_evidence_fixture" >"$tmp_dir/unverified-evidence.out" 2>&1
unverified_evidence_status="$?"
set -e
if [[ "$unverified_evidence_status" -eq 0 ]]; then
  echo "Expected unverified proof evidence to fail in Slice 1" >&2
  exit 1
fi
grep -q "unsupported evidence value" "$tmp_dir/unverified-evidence.out"

printf '%s\n' '<!-- sdp-trace-claim: claim=task_closed; subject=T071; state=stale; profile=repo_baseline; evidence=none -->' >"$none_evidence_fixture"
set +e
scripts/validate-claim-consistency.mjs --files "$none_evidence_fixture" >"$tmp_dir/none-evidence.out" 2>&1
none_evidence_status="$?"
set -e
if [[ "$none_evidence_status" -eq 0 ]]; then
  echo "Expected none evidence to fail in Slice 1" >&2
  exit 1
fi
grep -q "unsupported evidence value" "$tmp_dir/none-evidence.out"

if [[ "$(jq -r '.artifact_digest_status' "$contract_release_verification")" != "matched" ]]; then
  echo "Contract release verification example must claim matched artifact digests after source-bound local release repair" >&2
  exit 1
fi

if [[ "$(jq -r '.expected_proof_states.digest_verified' "$self_attestation_case")" != "true" ]]; then
  echo "Self-attestation positive case must expect digest_verified=true after source-bound local release repair" >&2
  exit 1
fi

scripts/verify-self-attestation.sh --case "$self_attestation_case" >/dev/null

if [[ "$(jq -r '.completion_state' "$block06_proof_states")" == "incomplete" ]]; then
  if grep -q 'All required Block 06 proof states are `observed`' "$block06_report"; then
    echo "Block 06 report must not claim all required proof states are observed while proof-states.json is incomplete" >&2
    exit 1
  fi
fi

set +e
scripts/verify.sh --profile source-bound --allow-dirty --json >"$source_bound_json" 2>&1
source_bound_status="$?"
set -e
if [[ -n "$(git status --porcelain)" ]]; then
  if [[ "$source_bound_status" -ne 3 ]]; then
    echo "Expected source-bound verifier to fail closed from a dirty checkout even with --allow-dirty" >&2
    cat "$source_bound_json" >&2
    exit 1
  fi
  jq -e '
    .profile == "source_bound_local_release" and
    .result == "cannot_verify" and
    .trust_scope == "none_dirty_checkout" and
    ([.states[] | select(.id == "source_checkout_clean" and .result == "cannot_verify")] | length == 1)
  ' "$source_bound_json" >/dev/null
else
  if [[ "$source_bound_status" -ne 0 ]]; then
    echo "Expected source-bound verifier to pass in a clean checkout after source-bound local release repair" >&2
    cat "$source_bound_json" >&2
    exit 1
  fi
  jq -e '
    .profile == "source_bound_local_release" and
    .result == "pass" and
    .trust_scope == "source_bound_local_release" and
    ([.states[] | select(.id == "source_checkout_clean" and .result == "pass")] | length == 1) and
    ([.states[] | select(.id == "artifact_subject_verified" and .result == "pass")] | length == 1) and
    ([.states[] | select(.id == "source_subject_verified" and .result == "pass")] | length == 1) and
    ([.states[] | select(.id == "local_envelope_verified" and .result == "pass")] | length == 1)
  ' "$source_bound_json" >/dev/null
fi
node scripts/validate-json-schema.mjs schema/proof-summary.schema.json "$source_bound_json" >/dev/null

set +e
scripts/verify.sh --profile external-trust --allow-dirty --json >"$external_trust_json" 2>&1
external_trust_status="$?"
set -e
if [[ -n "$(git status --porcelain)" ]]; then
  if [[ "$external_trust_status" -ne 3 ]]; then
    echo "Expected external-trust verifier to fail closed from a dirty checkout even with --allow-dirty" >&2
    cat "$external_trust_json" >&2
    exit 1
  fi
  jq -e '
    .profile == "external_production_trust" and
    .result == "cannot_verify" and
    .trust_scope == "none_dirty_checkout" and
    ([.states[] | select(.id == "external_attestation_present" and .result == "not_assessed")] | length == 1)
  ' "$external_trust_json" >/dev/null
else
  if [[ "$external_trust_status" -ne 1 ]]; then
    echo "Expected external-trust verifier to fail from missing external evidence in a clean checkout" >&2
    cat "$external_trust_json" >&2
    exit 1
  fi
  jq -e '
    .profile == "external_production_trust" and
    .result == "fail" and
    .trust_scope == "external_production_trust" and
    ([.states[] | select(.id == "external_attestation_present" and .result == "not_assessed")] | length == 1) and
    ([.states[] | select(.id == "production_release_verified" and .result == "fail")] | length == 1)
  ' "$external_trust_json" >/dev/null
fi
node scripts/validate-json-schema.mjs schema/proof-summary.schema.json "$external_trust_json" >/dev/null

scripts/validate-cross-references.mjs >/dev/null
mkdir -p "$cross_ref_fixture"
cp examples/self-trace/*.json "$cross_ref_fixture/"
jq '.[0].evidence_refs = ["evidence-does-not-exist"]' \
  examples/self-trace/observations.json >"$cross_ref_fixture/observations.json"
set +e
scripts/validate-cross-references.mjs --self-trace-dir "$cross_ref_fixture" >"$tmp_dir/cross-ref-negative.out" 2>&1
cross_ref_negative_status="$?"
set -e
if [[ "$cross_ref_negative_status" -eq 0 ]]; then
  echo "Expected dangling self-trace evidence reference to fail cross-reference validation" >&2
  exit 1
fi
grep -q "evidence-does-not-exist" "$tmp_dir/cross-ref-negative.out"

if [[ -n "$(git status --porcelain)" ]]; then
  set +e
  scripts/verify.sh --profile baseline --json >"$dirty_out" 2>&1
  dirty_status="$?"
  set -e
  if [[ "$dirty_status" -ne 3 ]]; then
    echo "Expected dirty baseline without --allow-dirty to exit 3" >&2
    cat "$dirty_out" >&2
    exit 1
  fi
  jq -e '
    .result == "cannot_verify" and
    .trust_scope == "none_dirty_checkout" and
    ([.states[] | select(.id == "source_checkout_clean" and .result == "cannot_verify")] | length == 1)
  ' "$dirty_out" >/dev/null
  node scripts/validate-json-schema.mjs schema/proof-summary.schema.json "$dirty_out" >/dev/null
fi

scripts/verify.sh --profile baseline --allow-dirty --json >"$baseline_json"

if [[ -n "$(git status --porcelain)" ]]; then
  expected_baseline_trust_scope="local_dirty_structural_only"
  expected_checkout_result="not_assessed"
else
  expected_baseline_trust_scope="clean_checkout_structural"
  expected_checkout_result="pass"
fi

jq -e \
  --arg expected_trust_scope "$expected_baseline_trust_scope" \
  --arg expected_checkout_result "$expected_checkout_result" '
    .artifact_role == "verifier_output" and
    .profile == "repo_baseline_structural" and
    .result == "pass" and
    .trust_scope == $expected_trust_scope and
    (.generated_by.tool == "scripts/verify.sh") and
    (.generated_by.invocation | contains("--allow-dirty")) and
    (.source_subject.type == "git_commit_v1") and
    (.source_subject.dirty_allowed == true) and
    (.spec_subject.path == "specs/001-sdp-trace-time-series-evidence-substrate/blocks/07-minimum-trust-kernel.md") and
    (.gate_set_subject.digest_algorithm == "sha256") and
    ([.states[] | select(.id == "source_checkout_clean" and .result == $expected_checkout_result)] | length == 1) and
    ([.states[] | select(.id == "claim_tags_consistent" and .result == "pass")] | length == 1) and
    ([.states[] | select(.id == "cross_reference_integrity" and .result == "pass")] | length == 1) and
    (.failures | length == 0)
  ' "$baseline_json" >/dev/null

node scripts/validate-json-schema.mjs schema/proof-summary.schema.json "$baseline_json" >/dev/null

scripts/verify.sh --profile repo_baseline --allow-dirty --json >"$tmp_dir/repo-baseline-alias.json"
jq -e --arg expected_trust_scope "$expected_baseline_trust_scope" '
  .profile == "repo_baseline_structural" and
  .result == "pass" and
  .trust_scope == $expected_trust_scope
' "$tmp_dir/repo-baseline-alias.json" >/dev/null
node scripts/validate-json-schema.mjs schema/proof-summary.schema.json "$tmp_dir/repo-baseline-alias.json" >/dev/null

scripts/verify.sh --profile baseline --allow-dirty --example --json >"$tmp_dir/example-generated.json"
jq -e '
  .artifact_role == "untrusted_shape_example" and
  .trust_scope == "untrusted_shape_only" and
  .generated_at == "1970-01-01T00:00:00Z" and
  (.generated_by.invocation | contains("--example")) and
  (.source_subject.commit == "0000000000000000000000000000000000000000") and
  (.generated_by.gate_set_digest == "0000000000000000000000000000000000000000000000000000000000000000") and
  ((.gate_set_subject.files | length) as $file_count |
    ([.gate_set_subject.files[] | select(.sha256 == "0000000000000000000000000000000000000000000000000000000000000000")] | length) == $file_count)
' "$tmp_dir/example-generated.json" >/dev/null
node scripts/validate-json-schema.mjs schema/proof-summary.schema.json "$tmp_dir/example-generated.json" >/dev/null

if [[ -f examples/self-trace/proof-summary.example.json ]]; then
  jq -e '
    .artifact_role == "untrusted_shape_example" and
    .trust_scope == "untrusted_shape_only" and
    .generated_at == "1970-01-01T00:00:00Z" and
    (.generated_by.invocation | contains("--example")) and
    (.source_subject.commit == "0000000000000000000000000000000000000000") and
    (.generated_by.gate_set_digest == "0000000000000000000000000000000000000000000000000000000000000000")
  ' examples/self-trace/proof-summary.example.json >/dev/null
  node scripts/validate-json-schema.mjs schema/proof-summary.schema.json examples/self-trace/proof-summary.example.json >/dev/null
fi

scripts/verify.sh --profile baseline --allow-dirty >"$baseline_text"
grep -q "profile: repo_baseline_structural" "$baseline_text"
grep -q "result: pass" "$baseline_text"
grep -q "$expected_checkout_result source_checkout_clean" "$baseline_text"
grep -q "pass claim_tags_consistent" "$baseline_text"

scripts/verify.sh --help >"$tmp_dir/help.out" 2>&1
grep -q -- "--allow-dirty" "$tmp_dir/help.out"
grep -q "Dirty checkouts without --allow-dirty exit 3" "$tmp_dir/help.out"

if [[ "$(scripts/verify.sh --version)" != "0.1.0" ]]; then
  echo "Expected --version to output 0.1.0" >&2
  exit 1
fi

set +e
scripts/verify.sh --profile does-not-exist >"$usage_out" 2>&1
usage_status="$?"
set -e
if [[ "$usage_status" -eq 0 ]]; then
  echo "Expected unsupported profile to fail" >&2
  exit 1
fi
if [[ "$usage_status" -ne 2 ]]; then
  echo "Expected unsupported profile to exit 2" >&2
  cat "$usage_out" >&2
  exit 1
fi
grep -q "Unsupported profile" "$usage_out"

set +e
scripts/verify.sh --bogus >"$tmp_dir/bogus.out" 2>&1
bogus_status="$?"
set -e
if [[ "$bogus_status" -ne 2 ]]; then
  echo "Expected unknown argument to exit 2" >&2
  cat "$tmp_dir/bogus.out" >&2
  exit 1
fi
grep -q "Unsupported argument" "$tmp_dir/bogus.out"
