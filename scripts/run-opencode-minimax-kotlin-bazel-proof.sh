#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'USAGE'
Usage:
  scripts/run-opencode-minimax-kotlin-bazel-proof.sh \
    --repo <kotlin-bazel-repo> \
    --scope <scope-path> \
    --bazel-target <//scope:target> \
    --bazel-command "<bazel test //scope:target>" \
    --model <opencode-minimax-model-id> \
    --out <.sdp-trace-runs/...> \
    [--timeout-seconds <seconds>]

The runner shells out to external OpenCode and Bazel tools. It does not install
or vendor those tools and it does not make policy verdicts.
USAGE
}

repo=""
scope=""
bazel_target=""
bazel_command=""
model=""
out=""
timeout_seconds="900"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      repo="$2"
      shift 2
      ;;
    --scope)
      scope="$2"
      shift 2
      ;;
    --bazel-target)
      bazel_target="$2"
      shift 2
      ;;
    --bazel-command)
      bazel_command="$2"
      shift 2
      ;;
    --model)
      model="$2"
      shift 2
      ;;
    --out)
      out="$2"
      shift 2
      ;;
    --timeout-seconds)
      timeout_seconds="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

for var_name in repo scope bazel_target bazel_command model out; do
  if [[ -z "${!var_name}" ]]; then
    echo "Missing required option: --${var_name//_/-}" >&2
    usage >&2
    exit 2
  fi
done

if [[ ! "$timeout_seconds" =~ ^[0-9]+$ || "$timeout_seconds" -lt 1 ]]; then
  echo "Timeout must be a positive integer number of seconds: $timeout_seconds" >&2
  exit 2
fi

if [[ ! -d "$repo" ]]; then
  echo "Repository path not found: $repo" >&2
  exit 1
fi

case "$bazel_command" in
  bazel\ *|bazelisk\ *)
    ;;
  *)
    echo "Bazel command must be a single bazel or bazelisk invocation: $bazel_command" >&2
    exit 2
    ;;
esac

if [[ "$bazel_command" =~ [\;\&\|\<\>\`\$] ]]; then
  echo "Bazel command must not contain shell metacharacters: $bazel_command" >&2
  exit 2
fi

read -r -a bazel_command_argv <<< "$bazel_command"
if [[ "${#bazel_command_argv[@]}" -lt 3 ]]; then
  echo "Bazel command must include bazel/bazelisk, build/test, and the target: $bazel_command" >&2
  exit 2
fi
if [[ "${bazel_command_argv[1]}" != "build" && "${bazel_command_argv[1]}" != "test" ]]; then
  echo "Bazel command must use build or test for full proof: $bazel_command" >&2
  exit 2
fi
bazel_target_arg_observed=0
for bazel_arg in "${bazel_command_argv[@]}"; do
  if [[ "$bazel_arg" == "$bazel_target" ]]; then
    bazel_target_arg_observed=1
  fi
done
if [[ "$bazel_target_arg_observed" -ne 1 ]]; then
  echo "Bazel command must include the supplied target as a separate argument: $bazel_target" >&2
  exit 2
fi

case "$out" in
  .sdp-trace-runs/*|*/.sdp-trace-runs/*)
    ;;
  *)
    if ! git check-ignore -q "$out" 2>/dev/null; then
      echo "Output path must be under .sdp-trace-runs or another gitignored path: $out" >&2
      exit 1
    fi
    ;;
esac

mkdir -p "$out/raw" "$out/evidence" "$out/handoff"
out_abs="$(cd "$out" && pwd)"

timestamp="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
opencode_bin="$(command -v opencode || true)"
bazel_bin="$(command -v bazel || command -v bazelisk || true)"
kotlin_bin="$(command -v kotlin || true)"
kotlinc_bin="$(command -v kotlinc || true)"

state_opencode="not_assessed"
state_model_listed="not_assessed"
state_model_access="not_assessed"
state_target="not_assessed"
state_opencode_run="not_assessed"
state_bazel_run="not_assessed"
state_package_valid="not_assessed"
state_committed="not_assessed"

reason_opencode="opencode command was not checked"
reason_model_listed="opencode models was not checked"
reason_model_access="MiniMax access was not checked"
reason_target="Kotlin+Bazel target was not checked"
reason_opencode_run="OpenCode MiniMax run was not checked"
reason_bazel_run="Bazel command was not checked"
reason_package_valid="Package validation is a separate repository step"
reason_committed="Runner output is local under .sdp-trace-runs; committed sanitized report is a separate step"
opencode_run_started_at="not_assessed"
opencode_run_ended_at="not_assessed"
opencode_exit_code="not_assessed"
bazel_run_started_at="not_assessed"
bazel_run_ended_at="not_assessed"
bazel_exit_code="not_assessed"

run_timed() {
  local stdout_file="$1"
  local stderr_file="$2"
  shift 2
  local flag_file="$out/raw/timeout.flag"
  rm -f "$flag_file"
  kill_tree() {
    local target_pid="$1"
    local child
    while IFS= read -r child; do
      [[ -z "$child" ]] && continue
      kill_tree "$child"
    done < <(pgrep -P "$target_pid" 2>/dev/null || true)
    kill "$target_pid" 2>/dev/null || true
  }
  "$@" >"$stdout_file" 2>"$stderr_file" &
  local pid=$!
  (
    sleep "$timeout_seconds"
    if kill -0 "$pid" 2>/dev/null; then
      printf 'timeout\n' >"$flag_file"
      kill_tree "$pid"
    fi
  ) &
  local watcher=$!
  local status=0
  wait "$pid" || status=$?
  kill "$watcher" 2>/dev/null || true
  wait "$watcher" 2>/dev/null || true
  if [[ -f "$flag_file" ]]; then
    return 124
  fi
  return "$status"
}

if [[ -n "$opencode_bin" ]] && "$opencode_bin" --version >"$out/raw/opencode-version.txt" 2>"$out/raw/opencode-version.err"; then
  state_opencode="observed"
  reason_opencode="opencode --version succeeded"
else
  state_opencode="not_observed"
  reason_opencode="opencode --version failed or opencode was not in PATH"
fi

if [[ "$state_opencode" == "observed" ]] && "$opencode_bin" models >"$out/raw/opencode-models.txt" 2>"$out/raw/opencode-models.err"; then
  if grep -Fq "$model" "$out/raw/opencode-models.txt"; then
    state_model_listed="observed"
    reason_model_listed="requested MiniMax model id appears in opencode models output"
  else
    state_model_listed="not_observed"
    reason_model_listed="requested MiniMax model id does not appear in opencode models output"
  fi
else
  state_model_listed="not_assessed"
  reason_model_listed="opencode models could not be run"
fi

scope_path="$repo/$scope"
if [[ ! -e "$scope_path" ]]; then
  state_target="not_observed"
  reason_target="scope path does not exist in repository"
else
  if [[ -z "$bazel_bin" ]]; then
    state_target="not_assessed"
    reason_target="bazel or bazelisk was not found in PATH"
  elif (cd "$repo" && "$bazel_bin" query "$bazel_target" >"$out_abs/raw/bazel-query.txt" 2>"$out_abs/raw/bazel-query.err") &&
       (cd "$repo" && "$bazel_bin" query --output=build "$bazel_target" >"$out_abs/raw/bazel-target-build.txt" 2>"$out_abs/raw/bazel-target-build.err"); then
    if find "$scope_path" -type f \( -name '*.kt' -o -name '*.kts' -o -name 'BUILD' -o -name 'BUILD.bazel' \) | grep -q . &&
       grep -Eq '\.kt|\.kts|kt_jvm|kotlin' "$out_abs/raw/bazel-target-build.txt"; then
      state_target="observed"
      reason_target="bazel query succeeded and target rule output ties Kotlin/Bazel files to the supplied target"
    else
      state_target="not_observed"
      reason_target="bazel query succeeded but target-tied Kotlin/Bazel evidence was not found"
    fi
  else
    state_target="not_observed"
    reason_target="bazel query failed for supplied target"
  fi
fi

if [[ -n "$bazel_bin" ]]; then
  "$bazel_bin" --version >"$out/raw/bazel-version.txt" 2>"$out/raw/bazel-version.err" || true
fi
if [[ -n "$kotlin_bin" ]]; then
  "$kotlin_bin" -version >"$out/raw/kotlin-version.txt" 2>"$out/raw/kotlin-version.err" || true
fi
if [[ -n "$kotlinc_bin" ]]; then
  "$kotlinc_bin" -version >"$out/raw/kotlinc-version.txt" 2>"$out/raw/kotlinc-version.err" || true
fi

before_status=""
after_status=""
source_ref="not_assessed"
if git -C "$repo" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  source_ref="$(git -C "$repo" rev-parse HEAD)"
  before_status="$(git -C "$repo" status --porcelain)"
  printf '%s\n' "$before_status" >"$out/raw/git-status-before.txt"
fi

prompt_file="$out/raw/opencode-prompt.txt"
cat >"$prompt_file" <<PROMPT
Inspect the scoped Kotlin+Bazel target at $scope.
Do not edit files.
Do not make readiness, support, pass/fail, or compatibility claims.
Ignore local Bazel output directories or symlinks named bazel-*; build evidence comes from the supplied runner command.
Report:
1. Kotlin evidence found or missing.
2. Bazel ownership evidence found or missing.
3. The supplied Bazel target $bazel_target and whether the supplied Bazel command evidence was observed.
4. Files inspected.
5. Claims that are not backed by inspected files or command output.
Return a concise summary suitable for a sanitized evidence record.
PROMPT

if [[ "$state_opencode" == "observed" ]]; then
  opencode_run_started_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
  if run_timed "$out/raw/opencode-run.jsonl" "$out/raw/opencode-run.err" \
    "$opencode_bin" run --model "$model" --format json --dir "$repo" "$(cat "$prompt_file")"; then
    opencode_exit_code="0"
    if [[ -s "$out/raw/opencode-run.jsonl" ]]; then
      state_opencode_run="observed"
      state_model_access="observed"
      reason_opencode_run="opencode run completed with requested MiniMax model id and produced captured output"
      reason_model_access="successful opencode run verifies access to requested model id"
    else
      state_opencode_run="not_observed"
      state_model_access="not_observed"
      reason_opencode_run="opencode run exited 0 but produced no captured output"
      reason_model_access="model access was not verified because captured output was empty"
    fi
  else
    opencode_exit_code="$?"
    state_opencode_run="not_observed"
    state_model_access="not_observed"
    reason_opencode_run="opencode run failed or timed out"
    reason_model_access="model access was not verified because opencode run failed or timed out"
  fi
  opencode_run_ended_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
else
  state_opencode_run="not_assessed"
  state_model_access="not_assessed"
  reason_opencode_run="opencode is unavailable"
  reason_model_access="opencode is unavailable"
fi

if git -C "$repo" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  after_status="$(git -C "$repo" status --porcelain)"
  printf '%s\n' "$after_status" >"$out/raw/git-status-after.txt"
  if [[ "$before_status" != "$after_status" ]]; then
    state_opencode_run="not_observed"
    reason_opencode_run="opencode run changed the assessed repository working tree"
  fi
fi

if [[ -n "$bazel_bin" ]]; then
  bazel_run_started_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
  if (cd "$repo" && run_timed "$out_abs/raw/bazel-command.out" "$out_abs/raw/bazel-command.err" "${bazel_command_argv[@]}"); then
    bazel_exit_code="0"
    state_bazel_run="observed"
    reason_bazel_run="operator-approved bazel command completed"
  else
    bazel_exit_code="$?"
    state_bazel_run="not_observed"
    reason_bazel_run="operator-approved bazel command failed or timed out"
  fi
  bazel_run_ended_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
else
  state_bazel_run="not_assessed"
  reason_bazel_run="bazel or bazelisk was not found in PATH"
fi

if git -C "$repo" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  bazel_after_status="$(git -C "$repo" status --porcelain)"
  printf '%s\n' "$bazel_after_status" >"$out/raw/git-status-after-bazel.txt"
  if [[ "$before_status" != "$bazel_after_status" ]]; then
    state_bazel_run="not_observed"
    reason_bazel_run="bazel command changed the assessed repository working tree"
  fi
fi

completion_state="incomplete"
if [[ "$state_opencode" == "observed" &&
      "$state_model_listed" == "observed" &&
      "$state_model_access" == "observed" &&
      "$state_target" == "observed" &&
      "$state_opencode_run" == "observed" &&
      "$state_bazel_run" == "observed" &&
      "$state_package_valid" == "observed" &&
      "$state_committed" == "observed" ]]; then
  completion_state="complete"
fi

first_line_or_not_assessed() {
  local file="$1"
  if [[ -s "$file" ]]; then
    sed -n '1p' "$file"
  else
    printf 'not_assessed'
  fi
}

opencode_version="$(first_line_or_not_assessed "$out/raw/opencode-version.txt")"
bazel_version="$(first_line_or_not_assessed "$out/raw/bazel-version.txt")"
kotlin_version="$(first_line_or_not_assessed "$out/raw/kotlin-version.err")"
if [[ "$kotlin_version" == "not_assessed" ]]; then
  kotlin_version="$(first_line_or_not_assessed "$out/raw/kotlin-version.txt")"
fi
kotlinc_version="$(first_line_or_not_assessed "$out/raw/kotlinc-version.err")"
if [[ "$kotlinc_version" == "not_assessed" ]]; then
  kotlinc_version="$(first_line_or_not_assessed "$out/raw/kotlinc-version.txt")"
fi

node - "$out" "$timestamp" "$repo" "$scope" "$bazel_target" "$bazel_command" "$model" \
  "$source_ref" "$opencode_version" "$bazel_version" "$kotlin_version" "$kotlinc_version" \
  "$opencode_run_started_at" "$opencode_run_ended_at" "$opencode_exit_code" \
  "$bazel_run_started_at" "$bazel_run_ended_at" "$bazel_exit_code" \
  "$state_opencode" "$reason_opencode" \
  "$state_model_listed" "$reason_model_listed" \
  "$state_model_access" "$reason_model_access" \
  "$state_target" "$reason_target" \
  "$state_opencode_run" "$reason_opencode_run" \
  "$state_bazel_run" "$reason_bazel_run" \
  "$state_package_valid" "$reason_package_valid" \
  "$state_committed" "$reason_committed" \
  "$completion_state" <<'NODE'
const fs = require('fs');
const crypto = require('crypto');
const path = require('path');
const { spawnSync } = require('child_process');
const [
  out, timestamp, repo, scope, bazelTarget, bazelCommand, model,
  sourceRef, opencodeVersion, bazelVersion, kotlinVersion, kotlincVersion,
  opencodeRunStartedAt, opencodeRunEndedAt, opencodeExitCode,
  bazelRunStartedAt, bazelRunEndedAt, bazelExitCode,
  stateOpencode, reasonOpencode,
  stateModelListed, reasonModelListed,
  stateModelAccess, reasonModelAccess,
  stateTarget, reasonTarget,
  stateOpencodeRun, reasonOpencodeRun,
  stateBazelRun, reasonBazelRun,
  statePackageValid, reasonPackageValid,
  stateCommitted, reasonCommitted,
  completionState
] = process.argv.slice(2);
function fileSha(path) {
  if (!fs.existsSync(path)) return null;
  return crypto.createHash('sha256').update(fs.readFileSync(path)).digest('hex');
}
function walkFiles(dir) {
  if (!fs.existsSync(dir)) return [];
  return fs.readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const fullPath = path.join(dir, entry.name);
    if (entry.isSymbolicLink() || entry.name.startsWith('bazel-')) return [];
    if (entry.isDirectory()) return walkFiles(fullPath);
    if (entry.isFile()) return [fullPath];
    return [];
  });
}
const sourcePaths = [
  ...walkFiles(path.join(repo, scope)),
  ...['MODULE.bazel', 'WORKSPACE', 'WORKSPACE.bazel']
    .map((name) => path.join(repo, name))
    .filter((candidate) => fs.existsSync(candidate) && fs.statSync(candidate).isFile())
].sort();
const sourceArtifacts = sourcePaths.map((candidate) => ({
  path: path.relative(repo, candidate),
  sha256: fileSha(candidate)
}));
const sourceDigest = crypto
  .createHash('sha256')
  .update(sourceArtifacts.map((artifact) => `${artifact.path}\t${artifact.sha256}`).join('\n'))
  .digest('hex');
let sourceCommitArtifactsVerified = 'not_assessed';
if (/^[0-9a-f]{40}$/.test(sourceRef) && sourceArtifacts.length > 0) {
  sourceCommitArtifactsVerified = sourceArtifacts.every((artifact) => {
    const show = spawnSync('git', ['-C', repo, 'show', `${sourceRef}:${artifact.path}`], { encoding: null });
    if (show.status !== 0) return false;
    const actual = crypto.createHash('sha256').update(show.stdout).digest('hex');
    return actual === artifact.sha256;
  });
}
const effectiveSourceRef = sourceCommitArtifactsVerified === true
  ? `${sourceRef}+scope:${sourceDigest}`
  : `working-tree-scope:${sourceDigest}`;
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
const proofPairs = [
  ['opencode_available', stateOpencode, reasonOpencode],
  ['minimax_model_listed', stateModelListed, reasonModelListed],
  ['minimax_access_verified', stateModelAccess, reasonModelAccess],
  ['kotlin_bazel_target_identified', stateTarget, reasonTarget],
  ['opencode_minimax_run_completed', stateOpencodeRun, reasonOpencodeRun],
  ['bazel_commands_executed', stateBazelRun, reasonBazelRun],
  ['sdp_trace_package_valid', statePackageValid, reasonPackageValid],
  ['sanitized_report_committed', stateCommitted, reasonCommitted]
];
const proofStates = {
  schema_version: '0.1.0',
  proof_profile: 'opencode-minimax-kotlin-bazel-e2e-v1',
  completion_state: completionState,
  generated_at: timestamp,
  tested_on: {
    repository: repo,
    source_ref: effectiveSourceRef,
    source_commit: sourceRef,
    source_commit_artifacts_verified: sourceCommitArtifactsVerified,
    source_content_sha256: sourceDigest,
    source_artifacts: sourceArtifacts,
    scope,
    bazel_target: bazelTarget,
    bazel_command: bazelCommand,
    model,
    opencode_version: opencodeVersion,
    bazel_version: bazelVersion,
    kotlin_version: kotlinVersion,
    kotlinc_version: kotlincVersion
  },
  command_results: {
    opencode_run: {
      started_at: opencodeRunStartedAt,
      ended_at: opencodeRunEndedAt,
      exit_code: opencodeExitCode,
      stdout_sha256: fileSha(`${out}/raw/opencode-run.jsonl`),
      stderr_sha256: fileSha(`${out}/raw/opencode-run.err`)
    },
    bazel_command: {
      started_at: bazelRunStartedAt,
      ended_at: bazelRunEndedAt,
      exit_code: bazelExitCode,
      stdout_sha256: fileSha(`${out}/raw/bazel-command.out`),
      stderr_sha256: fileSha(`${out}/raw/bazel-command.err`)
    },
    opencode_models_sha256: fileSha(`${out}/raw/opencode-models.txt`),
    bazel_query_sha256: fileSha(`${out}/raw/bazel-query.txt`),
    bazel_target_build_sha256: fileSha(`${out}/raw/bazel-target-build.txt`),
    git_status_before_sha256: fileSha(`${out}/raw/git-status-before.txt`),
    git_status_after_sha256: fileSha(`${out}/raw/git-status-after.txt`),
    git_status_after_bazel_sha256: fileSha(`${out}/raw/git-status-after-bazel.txt`)
  },
  states: proofPairs.map(([name, state, reason]) => ({
    name,
    state,
    evidence_refs: state === 'observed' ? [`evidence-${name.replaceAll('_', '-')}`] : [],
    reason,
    next_required_evidence: state === 'observed' ? null : `Provide evidence for ${name}.`
  }))
};
const eventObservedAt = (name) => {
  if (['minimax-access-verified', 'opencode-minimax-run-completed'].includes(name)) return opencodeRunEndedAt;
  if (name === 'bazel-commands-executed') return bazelRunEndedAt;
  if (['sdp-trace-package-valid', 'sanitized-report-committed'].includes(name)) return timestamp;
  return timestamp;
};
const event = (name, status, summary, ref) => ({
  id: `evidence-${name}`,
  schema_version: '0.1.0',
  source: 'local-command',
  external_ref: ref,
  observed_at: eventObservedAt(name),
  collected_at: eventObservedAt(name),
  actor: { id: 'pilot-operator', actor_type: 'human_user', display_name: 'Pilot operator' },
  event_type: 'command',
  status,
  summary,
  redaction_status: 'redacted',
  integrity_status: 'verified_hash',
  accountability
});
const evidenceEvents = proofPairs.map(([name, state, reason]) =>
  event(name.replaceAll('_', '-'), state === 'observed' ? 'success' : 'not_assessed', reason, 'run-report.md')
);
const digestSource = fileSha(`${out}/raw/opencode-run.jsonl`) || '0'.repeat(64);
const provenance = [{
  id: 'provenance-opencode-minimax-kotlin-bazel-run',
  schema_version: '0.1.0',
  actor_id: 'pilot-operator',
  actor_type: 'human_user',
  harness: 'opencode',
  model_family: 'MiniMax',
  model_version: model,
  tool_name: 'opencode',
  command: `opencode run --model ${model} --format json --dir ${repo} <prompt_sha256:${fileSha(`${out}/raw/opencode-prompt.txt`)}>`,
  prompt_ref: `sha256:${fileSha(`${out}/raw/opencode-prompt.txt`)}`,
  captured_at: timestamp,
  payload_digest: digestSource,
  digest_algorithm: 'sha256',
  chain_scope: 'opencode-minimax-kotlin-bazel',
  unavailable_fields: completionState === 'complete' ? [] : [{
    field: 'complete_product_proof',
    state: 'not_assessed',
    reason: 'At least one required proof state was not observed.'
  }]
}];
const observation = {
  id: 'observation-opencode-minimax-kotlin-bazel-proof',
  schema_version: '0.1.0',
  scope: 'opencode-minimax-kotlin-bazel',
  observed_at: timestamp,
  statement: completionState === 'complete'
    ? 'The OpenCode + MiniMax + Kotlin+Bazel proof package has all required proof states observed.'
    : 'The OpenCode + MiniMax + Kotlin+Bazel proof package is incomplete and must not be used as product proof.',
  evidence_refs: evidenceEvents.map((e) => e.id),
  provenance_refs: provenance.map((p) => p.id),
  assessment_status: completionState === 'complete' ? 'assessed' : 'partial'
};
const observedCount = proofStates.states.filter((state) => state.state === 'observed').length;
const metric = {
  id: 'metric-opencode-minimax-kotlin-bazel-proof-states',
  schema_version: '0.1.0',
  metric_name: 'observed_required_proof_state_count',
  dimensions: { repository: 'sdp-trace', scope: 'opencode-minimax-kotlin-bazel', model_family: 'MiniMax', harness: 'OpenCode', stack: 'Kotlin', build_system: 'Bazel' },
  samples: [{
    id: 'sample-opencode-minimax-kotlin-bazel-observed-proof-states',
    value: observedCount,
    unit: 'count',
    window_start: timestamp,
    window_end: timestamp,
    dimensions: { repository: 'sdp-trace', scope: 'opencode-minimax-kotlin-bazel' },
    evidence_refs: evidenceEvents.map((e) => e.id),
    provenance_refs: provenance.map((p) => p.id),
    assessment_state: completionState === 'complete' ? 'assessed' : 'partial'
  }],
  comparisons: [],
  assessment_state: completionState === 'complete' ? 'assessed' : 'partial',
  created_at: timestamp,
  updated_at: timestamp
};
const trace = {
  schema_version: '0.1.0',
  nodes: [
    { id: 'block06-spec', kind: 'spec', label: 'Block 06 spec' },
    ...evidenceEvents.map((e) => ({ id: e.id, kind: 'evidence', label: e.summary })),
    { id: observation.id, kind: 'observation', label: 'Block 06 proof observation' },
    { id: metric.id, kind: 'metric_stream', label: 'Block 06 proof-state metric' }
  ],
  edges: evidenceEvents.map((e) => ({ from: e.id, to: observation.id, relation: 'supports' }))
};
const notAssessed = proofStates.states
  .filter((state) => state.state !== 'observed')
  .map((state) => ({ field: state.name, state: 'not_assessed', reason: state.reason }));
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
  evidence_events: evidenceEvents,
  provenance_records: provenance,
  metric_streams: [metric],
  observations: [observation],
  not_assessed: notAssessed,
  generated_at: timestamp,
  producer: 'sdp-trace Block 06 reference runner',
  accountability: { ...accountability, authority_scope: 'assessment_input' },
  risk_classification: risk,
  contract_release_verification_ref: 'examples/contract-foundation/contract-release-verification.example.json'
};
fs.writeFileSync(`${out}/evidence/proof-states.json`, JSON.stringify(proofStates, null, 2) + '\n');
fs.writeFileSync(`${out}/evidence/evidence-events.json`, JSON.stringify(evidenceEvents, null, 2) + '\n');
fs.writeFileSync(`${out}/evidence/provenance-records.json`, JSON.stringify(provenance, null, 2) + '\n');
fs.writeFileSync(`${out}/evidence/observations.json`, JSON.stringify([observation], null, 2) + '\n');
fs.writeFileSync(`${out}/evidence/metric-stream.json`, JSON.stringify([metric], null, 2) + '\n');
fs.writeFileSync(`${out}/evidence/trace-snapshot.json`, JSON.stringify(trace, null, 2) + '\n');
fs.writeFileSync(`${out}/handoff/assessment-input.json`, JSON.stringify(assessment, null, 2) + '\n');
fs.writeFileSync(`${out}/README.md`, `# OpenCode + MiniMax + Kotlin+Bazel Proof Package\n\nCompletion state: ${completionState}\n\n`);
fs.writeFileSync(`${out}/redaction-note.md`, '# Redaction Note\n\nRaw output remains under the local raw directory and must not be committed.\n');
const sourceArtifactRows = sourceArtifacts.map((artifact) => `| ${artifact.path} | ${artifact.sha256} |`).join('\n');
const proofStateRows = proofStates.states.map((state) =>
  `| ${state.name} | ${state.state} | ${state.reason.replace(/\|/g, '/')} |`
).join('\n');
const runReport = `# OpenCode + MiniMax + Kotlin+Bazel Run Report

Completion state: ${completionState}

## Tested-On Environment

| Field | Value |
|---|---|
| Repository | ${repo} |
| Source ref | ${effectiveSourceRef} |
| Source content sha256 | ${sourceDigest} |
| Scope | ${scope} |
| Bazel target | ${bazelTarget} |
| Bazel command | ${bazelCommand.replace(/\|/g, '/')} |
| Model | ${model} |
| OpenCode version | ${opencodeVersion} |
| Bazel version | ${bazelVersion} |
| Kotlin version | ${kotlinVersion} |
| kotlinc version | ${kotlincVersion} |

## Source Artifacts

| Path | sha256 |
|---|---|
${sourceArtifactRows}

## Command Results

| Command | Started | Ended | Exit code | stdout sha256 | stderr sha256 |
|---|---|---|---|---|---|
| OpenCode run | ${opencodeRunStartedAt} | ${opencodeRunEndedAt} | ${opencodeExitCode} | ${proofStates.command_results.opencode_run.stdout_sha256 || 'not_assessed'} | ${proofStates.command_results.opencode_run.stderr_sha256 || 'not_assessed'} |
| Bazel command | ${bazelRunStartedAt} | ${bazelRunEndedAt} | ${bazelExitCode} | ${proofStates.command_results.bazel_command.stdout_sha256 || 'not_assessed'} | ${proofStates.command_results.bazel_command.stderr_sha256 || 'not_assessed'} |

## Proof States

| Proof state | State | Reason |
|---|---|---|
${proofStateRows}

## Boundary

sdp-trace records observed evidence for this exact slice only. It does not produce a pass/fail, readiness, support, compatibility, or degradation verdict.
`;
fs.writeFileSync(`${out}/run-report.md`, runReport);
NODE

printf 'Proof states:\n'
jq -r '.states[] | "- \(.name): \(.state) (\(.reason))"' "$out/evidence/proof-states.json"
printf 'Completion state: %s\n' "$completion_state"

if [[ "$completion_state" == "complete" ]]; then
  exit 0
fi
exit 1
