#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

validator="scripts/validate-e2e-pilot-package.sh"
if [[ ! -x "$validator" ]]; then
  echo "Missing executable validator: $validator" >&2
  exit 1
fi

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

make_package() {
  local package_dir="$1"
  mkdir -p "$package_dir/evidence" "$package_dir/handoff"
  printf 'fixture package\n' >"$package_dir/README.md"
  printf 'fixture report\n' >"$package_dir/run-report.md"
  printf 'fixture opencode summary\n' >"$package_dir/opencode-summary.md"
  printf 'fixture redaction note\n' >"$package_dir/redaction-note.md"
  node - "$package_dir" <<'NODE'
const fs = require('fs');
const crypto = require('crypto');
const dir = process.argv[2];
const ts = '2026-05-01T00:00:00Z';
const packageJsonSha = crypto.createHash('sha256').update(fs.readFileSync('package.json')).digest('hex');
const accountability = {
  dri: { identity_ref: 'role:sdp-trace-engineering-dri', actor_type: 'human_role' },
  approver: { identity_ref: 'role:sdp-trace-cto', actor_type: 'human_role' },
  escalation: { identity_ref: 'role:sdp-trace-cto', actor_type: 'human_role' },
  authority_scope: 'evidence',
  accountability_claim: 'recording_only',
  approval_ref: 'specs/001-sdp-trace-time-series-evidence-substrate/blocks/06-opencode-minimax-kotlin-bazel-e2e-proof.md',
  risk_owner: { identity_ref: 'role:sdp-trace-risk-owner', actor_type: 'human_role' },
  line_of_defense: 'first'
};
const event = {
  id: 'evidence-opencode-minimax-run',
  schema_version: '0.1.0',
  source: 'local-command',
  external_ref: 'run-report.md',
  observed_at: ts,
  collected_at: ts,
  actor: { id: 'pilot-operator', actor_type: 'human_user', display_name: 'Pilot operator' },
  event_type: 'command',
  status: 'success',
  summary: 'OpenCode MiniMax Kotlin Bazel proof fixture event.',
  redaction_status: 'redacted',
  integrity_status: 'verified_hash',
  accountability
};
const provenance = {
  id: 'provenance-opencode-minimax-run',
  schema_version: '0.1.0',
  actor_id: 'pilot-operator',
  actor_type: 'human_user',
  harness: 'opencode',
  model_family: 'MiniMax',
  model_version: 'minimax-coding-plan/MiniMax-M2.5',
  tool_name: 'opencode',
  command: 'opencode run --model minimax-coding-plan/MiniMax-M2.5',
  captured_at: ts,
  payload_digest: '0'.repeat(64),
  digest_algorithm: 'sha256',
  chain_scope: 'opencode-minimax-kotlin-bazel'
};
const observation = {
  id: 'observation-opencode-minimax-kotlin-bazel-package',
  schema_version: '0.1.0',
  scope: 'opencode-minimax-kotlin-bazel',
  observed_at: ts,
  statement: 'The E2E pilot package fixture has all required proof states observed.',
  evidence_refs: [event.id],
  provenance_refs: [provenance.id],
  assessment_status: 'assessed'
};
const metric = {
  id: 'metric-opencode-minimax-proof-states',
  schema_version: '0.1.0',
  metric_name: 'observed_required_proof_state_count',
  dimensions: { repository: 'sdp-trace', scope: 'opencode-minimax-kotlin-bazel' },
  samples: [{
    id: 'sample-observed-required-proof-states',
    value: 8,
    unit: 'count',
    window_start: ts,
    window_end: ts,
    dimensions: { repository: 'sdp-trace', scope: 'opencode-minimax-kotlin-bazel' },
    evidence_refs: [event.id],
    provenance_refs: [provenance.id],
    assessment_state: 'assessed'
  }],
  comparisons: [],
  assessment_state: 'assessed',
  created_at: ts,
  updated_at: ts
};
const trace = {
  schema_version: '0.1.0',
  nodes: [
    { id: 'block06-spec', kind: 'spec', label: 'Block 06 spec' },
    { id: event.id, kind: 'evidence', label: 'OpenCode MiniMax run evidence' },
    { id: observation.id, kind: 'observation', label: 'E2E package observation' },
    { id: metric.id, kind: 'metric_stream', label: 'Observed proof states' }
  ],
  edges: [
    { from: event.id, to: observation.id, relation: 'supports' },
    { from: observation.id, to: metric.id, relation: 'supports' }
  ]
};
const risk = {
  schema_version: '0.1.0',
  observed_autonomy_level: 'collaborative',
  observed_impact_level: 'low',
  classification_source: 'human_dri',
  classification_ref: 'run-report.md'
};
const assessment = {
  id: 'assessment-input-opencode-minimax-kotlin-bazel',
  schema_version: '0.1.0',
  scope: 'opencode-minimax-kotlin-bazel',
  trace_snapshot_ref: 'evidence/trace-snapshot.json',
  evidence_events: [event],
  provenance_records: [provenance],
  metric_streams: [metric],
  observations: [observation],
  not_assessed: [],
  generated_at: ts,
  producer: 'sdp-trace e2e pilot package validator fixture',
  accountability: { ...accountability, authority_scope: 'assessment_input' },
  risk_classification: risk,
  contract_release_verification_ref: 'examples/contract-foundation/contract-release-verification.example.json'
};
const proofStates = {
  schema_version: '0.1.0',
  proof_profile: 'opencode-minimax-kotlin-bazel-e2e-v1',
  completion_state: 'complete',
  generated_at: ts,
  tested_on: {
    repository: '.',
    source_ref: `working-tree-scope:${packageJsonSha}`,
    source_commit: 'not_assessed',
    source_commit_artifacts_verified: 'not_assessed',
    source_content_sha256: packageJsonSha,
    source_artifacts: [{ path: 'package.json', sha256: packageJsonSha }],
    scope: 'services/example',
    bazel_target: '//services/example:compile_hello_jar',
    bazel_command: 'bazel build //services/example:compile_hello_jar',
    model: 'minimax-coding-plan/MiniMax-M2.5',
    opencode_version: 'fixture-opencode',
    bazel_version: 'fixture-bazel',
    kotlin_version: 'fixture-kotlin',
    kotlinc_version: 'fixture-kotlinc'
  },
  command_results: {
    opencode_run: {
      started_at: ts,
      ended_at: ts,
      exit_code: '0',
      stdout_sha256: '0'.repeat(64),
      stderr_sha256: '0'.repeat(64)
    },
    bazel_command: {
      started_at: ts,
      ended_at: ts,
      exit_code: '0',
      stdout_sha256: '0'.repeat(64),
      stderr_sha256: '0'.repeat(64)
    },
    opencode_models_sha256: '0'.repeat(64),
    bazel_query_sha256: '0'.repeat(64),
    bazel_target_build_sha256: '0'.repeat(64)
  },
  states: [
    'opencode_available',
    'minimax_model_listed',
    'minimax_access_verified',
    'kotlin_bazel_target_identified',
    'opencode_minimax_run_completed',
    'bazel_commands_executed',
    'sdp_trace_package_valid',
    'sanitized_report_committed'
  ].map((name) => ({
    name,
    state: 'observed',
    evidence_refs: [event.id],
    reason: 'Observed in fixture package.',
    next_required_evidence: null
  }))
};
fs.writeFileSync(`${dir}/evidence/evidence-events.json`, JSON.stringify([event], null, 2) + '\n');
fs.writeFileSync(`${dir}/evidence/provenance-records.json`, JSON.stringify([provenance], null, 2) + '\n');
fs.writeFileSync(`${dir}/evidence/observations.json`, JSON.stringify([observation], null, 2) + '\n');
fs.writeFileSync(`${dir}/evidence/metric-stream.json`, JSON.stringify([metric], null, 2) + '\n');
fs.writeFileSync(`${dir}/evidence/trace-snapshot.json`, JSON.stringify(trace, null, 2) + '\n');
fs.writeFileSync(`${dir}/handoff/assessment-input.json`, JSON.stringify(assessment, null, 2) + '\n');
fs.writeFileSync(`${dir}/evidence/proof-states.json`, JSON.stringify(proofStates, null, 2) + '\n');
NODE
}

expect_fail() {
  local label="$1"
  shift
  if "$@" >"$tmp/$label.out" 2>&1; then
    echo "Expected failure but command passed: $label" >&2
    cat "$tmp/$label.out" >&2
    exit 1
  fi
}

valid_pkg="$tmp/valid"
make_package "$valid_pkg"
"$validator" "$valid_pkg"
"$validator" --mode complete "$valid_pkg"

incomplete_pkg="$tmp/incomplete"
cp -R "$valid_pkg" "$incomplete_pkg"
jq '
  .completion_state = "incomplete" |
  (.states[] | select(.name == "sdp_trace_package_valid") | .state) = "not_assessed" |
  (.states[] | select(.name == "sdp_trace_package_valid") | .reason) = "Package validation not yet run." |
  (.states[] | select(.name == "sdp_trace_package_valid") | .next_required_evidence) = "Run package validation." |
  (.states[] | select(.name == "sanitized_report_committed") | .state) = "not_assessed" |
  (.states[] | select(.name == "sanitized_report_committed") | .reason) = "Sanitized report not yet committed." |
  (.states[] | select(.name == "sanitized_report_committed") | .next_required_evidence) = "Commit sanitized report."
' "$valid_pkg/evidence/proof-states.json" >"$incomplete_pkg/evidence/proof-states.json"
"$validator" --mode package "$incomplete_pkg"
expect_fail incomplete-complete "$validator" --mode complete "$incomplete_pkg"

missing_state_pkg="$tmp/missing-state"
cp -R "$valid_pkg" "$missing_state_pkg"
jq 'del(.states[] | select(.name == "bazel_commands_executed"))' \
  "$valid_pkg/evidence/proof-states.json" >"$missing_state_pkg/evidence/proof-states.json"
expect_fail missing-state "$validator" "$missing_state_pkg"

partial_pkg="$tmp/partial"
cp -R "$valid_pkg" "$partial_pkg"
jq '(.states[] | select(.name == "bazel_commands_executed") | .state) = "not_assessed"' \
  "$valid_pkg/evidence/proof-states.json" >"$partial_pkg/evidence/proof-states.json"
expect_fail partial "$validator" "$partial_pkg"

raw_pkg="$tmp/raw"
cp -R "$valid_pkg" "$raw_pkg"
mkdir -p "$raw_pkg/raw"
printf 'raw output must not be committed\n' >"$raw_pkg/raw/output.log"
expect_fail raw-output "$validator" "$raw_pkg"

dangling_ref_pkg="$tmp/dangling-ref"
cp -R "$valid_pkg" "$dangling_ref_pkg"
jq '(.states[] | select(.name == "opencode_available") | .evidence_refs) = ["evidence-missing"]' \
  "$valid_pkg/evidence/proof-states.json" >"$dangling_ref_pkg/evidence/proof-states.json"
expect_fail dangling-ref "$validator" "$dangling_ref_pkg"

echo "e2e pilot package tests passed"
