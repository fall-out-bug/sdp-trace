#!/usr/bin/env bash
set -euo pipefail

slice="all"
if [[ "${1:-}" == "--slice" ]]; then
  slice="${2:-}"
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

pass_count=0

pass() {
  pass_count=$((pass_count + 1))
  printf 'ok - %s\n' "$1"
}

expect_pass() {
  local name="$1"
  shift
  "$@" >/tmp/sdp_trace_fr_pass.json
  pass "$name"
}

expect_fail_contains() {
  local name="$1"
  local expected="$2"
  shift 2
  set +e
  "$@" >/tmp/sdp_trace_fr_fail.json 2>/tmp/sdp_trace_fr_fail.err
  local status=$?
  set -e
  if [[ "$status" -eq 0 ]]; then
    printf 'not ok - %s unexpectedly passed\n' "$name" >&2
    cat /tmp/sdp_trace_fr_fail.json >&2
    exit 1
  fi
  if ! grep -q "$expected" /tmp/sdp_trace_fr_fail.json /tmp/sdp_trace_fr_fail.err; then
    printf 'not ok - %s did not report %s\n' "$name" "$expected" >&2
    cat /tmp/sdp_trace_fr_fail.json >&2 || true
    cat /tmp/sdp_trace_fr_fail.err >&2 || true
    exit 1
  fi
  pass "$name"
}

mutate_json() {
  local file="$1"
  local script="$2"
  node -e "const fs=require('fs'); const file=process.argv[1]; const data=JSON.parse(fs.readFileSync(file,'utf8')); ${script}; fs.writeFileSync(file, JSON.stringify(data, null, 2) + '\n');" "$file"
}

run_chain_tests() {
  expect_pass \
    "local positive chain verifies" \
    scripts/verify-flight-recorder.mjs --profile local examples/flight-recorder/local-positive

  expect_fail_contains \
    "schema-valid mismatched event hash is rejected" \
    "event_hash_mismatch" \
    scripts/verify-flight-recorder.mjs --event-file examples/flight-recorder/negative-mismatched-event-hash/event.json

  local tmp
  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/local-positive "$tmp/mutated"
  mutate_json "$tmp/mutated/events/002-task-locked.json" "data.event_payload.summary='Tampered after recording';"
  expect_fail_contains \
    "mutated event payload is rejected" \
    "payload_digest_mismatch" \
    scripts/verify-flight-recorder.mjs --profile local "$tmp/mutated"

  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/local-positive "$tmp/deleted"
  rm "$tmp/deleted/events/001-source-baseline-recorded.json"
  expect_fail_contains \
    "deleted event is rejected" \
    "event_ref_missing" \
    scripts/verify-flight-recorder.mjs --profile local "$tmp/deleted"

  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/local-positive "$tmp/reordered"
  mutate_json "$tmp/reordered/run.json" "const first=data.event_refs[1].uri; data.event_refs[1].uri=data.event_refs[2].uri; data.event_refs[2].uri=first;"
  expect_fail_contains \
    "reordered event references are rejected" \
    "event_ref_hash_mismatch" \
    scripts/verify-flight-recorder.mjs --profile local "$tmp/reordered"

  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/local-positive "$tmp/time-backward"
  mutate_json "$tmp/time-backward/events/002-task-locked.json" "data.event_time='2026-05-05T09:59:59Z';"
  expect_fail_contains \
    "backward event time is rejected" \
    "event_time_order_mismatch" \
    scripts/verify-flight-recorder.mjs --profile local "$tmp/time-backward"

  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/requirement-supersession "$tmp/run-closed-not-terminal"
  mutate_json "$tmp/run-closed-not-terminal/run.json" "const close=data.event_refs.pop(); data.event_refs.splice(3,0,close);"
  expect_fail_contains \
    "non-terminal run closure is rejected" \
    "run_closed_not_terminal" \
    scripts/verify-flight-recorder.mjs --profile local "$tmp/run-closed-not-terminal"
}

run_witness_tests() {
  expect_pass \
    "witnessed positive chain verifies" \
    scripts/verify-flight-recorder.mjs --profile witnessed examples/flight-recorder/witnessed-positive

  expect_fail_contains \
    "missing witness material is rejected" \
    "witness_ref_missing" \
    scripts/verify-flight-recorder.mjs --profile witnessed examples/flight-recorder/negative-missing-witness-material

  local tmp
  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/witnessed-positive "$tmp/witness-mismatch"
  cp examples/flight-recorder/negative-witness-chain-head-mismatch/witness.json "$tmp/witness-mismatch/witness.json"
  expect_fail_contains \
    "witness chain head mismatch is rejected" \
    "witness_chain_head_mismatch" \
    scripts/verify-flight-recorder.mjs --profile witnessed "$tmp/witness-mismatch"

  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/witnessed-positive "$tmp/witness-task-mismatch"
  mutate_json "$tmp/witness-task-mismatch/witness.json" "data.task_hash='$(printf '1%.0s' {1..64})';"
  expect_fail_contains \
    "witness task hash mismatch is rejected" \
    "witness_task_hash_mismatch" \
    scripts/verify-flight-recorder.mjs --profile witnessed "$tmp/witness-task-mismatch"

  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/witnessed-positive "$tmp/witness-run-mismatch"
  mutate_json "$tmp/witness-run-mismatch/witness.json" "data.run_id='different-run-id';"
  expect_fail_contains \
    "witness run id mismatch is rejected" \
    "witness_run_id_mismatch" \
    scripts/verify-flight-recorder.mjs --profile witnessed "$tmp/witness-run-mismatch"

  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/witnessed-positive "$tmp/witness-scope-mismatch"
  mutate_json "$tmp/witness-scope-mismatch/witness.json" "data.witness_scope='external_witness_extension'; data.profile='externally_witnessed_run'; data.trust_scope='externally_witnessed_run';"
  expect_fail_contains \
    "witness scope mismatch is rejected" \
    "witness_scope_mismatch" \
    scripts/verify-flight-recorder.mjs --profile witnessed "$tmp/witness-scope-mismatch"

  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/witnessed-positive "$tmp/witness-source-mismatch"
  mutate_json "$tmp/witness-source-mismatch/witness.json" "data.source_baseline_hash='$(printf '4%.0s' {1..64})';"
  expect_fail_contains \
    "witness source baseline mismatch is rejected" \
    "witness_source_baseline_mismatch" \
    scripts/verify-flight-recorder.mjs --profile witnessed "$tmp/witness-source-mismatch"

  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/witnessed-positive "$tmp/witness-recorder-version-mismatch"
  mutate_json "$tmp/witness-recorder-version-mismatch/witness.json" "data.recorder_version='9.9.9';"
  expect_fail_contains \
    "witness recorder version mismatch is rejected" \
    "witness_recorder_version_mismatch" \
    scripts/verify-flight-recorder.mjs --profile witnessed "$tmp/witness-recorder-version-mismatch"

  expect_fail_contains \
    "witnessed run cannot be downgraded to local verification" \
    "profile_downgrade" \
    scripts/verify-flight-recorder.mjs --profile local examples/flight-recorder/witnessed-positive
}

run_expectation_tests() {
  expect_pass \
    "late attach fixture records explicit not_assessed boundary" \
    scripts/verify-flight-recorder.mjs --profile local examples/flight-recorder/late-attach

  expect_pass \
    "requirement supersession fixture preserves original task event" \
    scripts/verify-flight-recorder.mjs --profile local examples/flight-recorder/requirement-supersession

  local tmp
  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/local-positive "$tmp/source-summary-mismatch"
  mutate_json "$tmp/source-summary-mismatch/run.json" "data.source_summary.source_baseline_hash='$(printf '2%.0s' {1..64})';"
  expect_fail_contains \
    "source summary mismatch is rejected" \
    "source_summary_mismatch" \
    scripts/verify-flight-recorder.mjs --profile local "$tmp/source-summary-mismatch"

  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/local-positive "$tmp/task-summary-mismatch"
  mutate_json "$tmp/task-summary-mismatch/run.json" "data.task_summary.task_hash='$(printf '3%.0s' {1..64})';"
  expect_fail_contains \
    "task summary mismatch is rejected" \
    "task_summary_mismatch" \
    scripts/verify-flight-recorder.mjs --profile local "$tmp/task-summary-mismatch"

  tmp="$(mktemp -d)"
  cp -R examples/flight-recorder/requirement-supersession "$tmp/task-rewrite"
  mutate_json "$tmp/task-rewrite/events/002-task-locked.json" "data.event_payload.task_ref='rewritten-after-command';"
  expect_fail_contains \
    "task rewrite after command evidence is rejected" \
    "payload_digest_mismatch" \
    scripts/verify-flight-recorder.mjs --profile local "$tmp/task-rewrite"
}

run_redaction_tests() {
  expect_fail_contains \
    "unresolved redaction is rejected" \
    "redaction_unresolved" \
    scripts/verify-flight-recorder.mjs --profile local examples/flight-recorder/redaction-unresolved

  expect_fail_contains \
    "forensic critical digest-only evidence is rejected" \
    "forensic_digest_only_critical" \
    scripts/verify-flight-recorder.mjs --profile forensic examples/flight-recorder/forensic-digest-only-negative
}

case "$slice" in
  chain)
    run_chain_tests
    ;;
  witness)
    run_witness_tests
    ;;
  expectations)
    run_expectation_tests
    ;;
  redaction)
    run_redaction_tests
    ;;
  all)
    run_chain_tests
    run_witness_tests
    run_expectation_tests
    run_redaction_tests
    ;;
  *)
    printf 'Unknown flight-recorder test slice: %s\n' "$slice" >&2
    exit 2
    ;;
esac

printf 'flight-recorder tests passed: %s\n' "$pass_count"
