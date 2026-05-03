#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

mode="complete"
package_dir=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --mode)
      if [[ $# -lt 2 ]]; then
        echo "Usage: $0 [--mode package|complete] <opencode-minimax-kotlin-bazel-package-dir>" >&2
        exit 2
      fi
      mode="$2"
      shift 2
      ;;
    -h|--help)
      echo "Usage: $0 [--mode package|complete] <opencode-minimax-kotlin-bazel-package-dir>" >&2
      exit 0
      ;;
    --*)
      echo "Unsupported argument: $1" >&2
      exit 2
      ;;
    *)
      if [[ -n "$package_dir" ]]; then
        echo "Usage: $0 [--mode package|complete] <opencode-minimax-kotlin-bazel-package-dir>" >&2
        exit 2
      fi
      package_dir="$1"
      shift
      ;;
  esac
done

if [[ "$mode" != "package" && "$mode" != "complete" ]]; then
  echo "Unsupported mode: $mode" >&2
  exit 2
fi

if [[ -z "$package_dir" ]]; then
  echo "Usage: $0 [--mode package|complete] <opencode-minimax-kotlin-bazel-package-dir>" >&2
  exit 2
fi

if [[ ! -d "$package_dir" ]]; then
  echo "E2E pilot package directory not found: $package_dir" >&2
  exit 1
fi

required_files=(
  "README.md"
  "run-report.md"
  "opencode-summary.md"
  "redaction-note.md"
  "evidence/proof-states.json"
  "evidence/evidence-events.json"
  "evidence/provenance-records.json"
  "evidence/observations.json"
  "evidence/metric-stream.json"
  "evidence/trace-snapshot.json"
  "handoff/assessment-input.json"
)

for required_file in "${required_files[@]}"; do
  if [[ ! -f "$package_dir/$required_file" ]]; then
    echo "E2E pilot package missing required file: $required_file" >&2
    exit 1
  fi
done

if find "$package_dir" -path '*/raw/*' -o -name raw | grep -q .; then
  echo "E2E pilot package must not include raw output paths" >&2
  exit 1
fi

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

validate() {
  local schema="$1"
  local data="$2"
  node scripts/validate-json-schema.mjs "$schema" "$data" >/dev/null
}

sha256_file() {
  shasum -a 256 "$1" | awk '{print $1}'
}

jq empty \
  "$package_dir/evidence/proof-states.json" \
  "$package_dir/evidence/evidence-events.json" \
  "$package_dir/evidence/provenance-records.json" \
  "$package_dir/evidence/observations.json" \
  "$package_dir/evidence/metric-stream.json" \
  "$package_dir/evidence/trace-snapshot.json" \
  "$package_dir/handoff/assessment-input.json" >/dev/null

required_states=(
  "opencode_available"
  "minimax_model_listed"
  "minimax_access_verified"
  "kotlin_bazel_target_identified"
  "opencode_minimax_run_completed"
  "bazel_commands_executed"
  "sdp_trace_package_valid"
  "sanitized_report_committed"
)

proof_states="$package_dir/evidence/proof-states.json"
if [[ "$(jq -r '.proof_profile // empty' "$proof_states")" != "opencode-minimax-kotlin-bazel-e2e-v1" ]]; then
  echo "E2E pilot package proof_profile is not opencode-minimax-kotlin-bazel-e2e-v1" >&2
  exit 1
fi
completion_state="$(jq -r '.completion_state // empty' "$proof_states")"
if [[ "$mode" == "complete" && "$completion_state" != "complete" ]]; then
  echo "E2E pilot package completion_state must be complete for complete proof packages" >&2
  exit 1
fi
if [[ "$mode" == "package" && "$completion_state" != "complete" && "$completion_state" != "incomplete" ]]; then
  echo "E2E pilot package completion_state must be complete or incomplete in package mode" >&2
  exit 1
fi

tested_model="$(jq -r '.tested_on.model // empty' "$proof_states")"
if [[ ! "$tested_model" =~ [Mm]ini[Mm]ax ]]; then
  echo "E2E pilot package tested_on.model must name a MiniMax model id" >&2
  exit 1
fi
tested_target="$(jq -r '.tested_on.bazel_target // empty' "$proof_states")"
tested_command="$(jq -r '.tested_on.bazel_command // empty' "$proof_states")"
if [[ -z "$tested_target" || "$tested_command" != *"$tested_target"* ]]; then
  echo "E2E pilot package tested_on.bazel_command must include tested_on.bazel_target" >&2
  exit 1
fi
if [[ "$mode" == "complete" && "$(jq -r '.command_results.opencode_run.exit_code // empty' "$proof_states")" != "0" ]]; then
  echo "E2E pilot package OpenCode command result must have exit_code 0" >&2
  exit 1
fi
if [[ "$mode" == "complete" && "$(jq -r '.command_results.bazel_command.exit_code // empty' "$proof_states")" != "0" ]]; then
  echo "E2E pilot package Bazel command result must have exit_code 0" >&2
  exit 1
fi
for digest_expr in \
  '.command_results.opencode_run.stdout_sha256' \
  '.command_results.opencode_run.stderr_sha256' \
  '.command_results.bazel_command.stdout_sha256' \
  '.command_results.bazel_command.stderr_sha256' \
  '.command_results.opencode_models_sha256' \
  '.command_results.bazel_query_sha256' \
  '.command_results.bazel_target_build_sha256'; do
  digest_value="$(jq -r "$digest_expr // empty" "$proof_states")"
  if [[ ! "$digest_value" =~ ^[0-9a-f]{64}$ ]]; then
    echo "E2E pilot package command digest is missing or invalid for $digest_expr" >&2
    exit 1
  fi
done

tested_repo="$(jq -r '.tested_on.repository // empty' "$proof_states")"
if [[ -z "$tested_repo" || ! -d "$tested_repo" ]]; then
  echo "E2E pilot package tested_on.repository must exist in the checkout" >&2
  exit 1
fi
while IFS=$'\t' read -r source_path expected_sha; do
  if [[ -z "$source_path" || -z "$expected_sha" ]]; then
    echo "E2E pilot package source_artifacts must include path and sha256" >&2
    exit 1
  fi
  if [[ ! -f "$tested_repo/$source_path" ]]; then
    echo "E2E pilot package source artifact is missing from checkout: $tested_repo/$source_path" >&2
    exit 1
  fi
  actual_sha="$(sha256_file "$tested_repo/$source_path")"
  if [[ "$actual_sha" != "$expected_sha" ]]; then
    echo "E2E pilot package source artifact digest mismatch: $tested_repo/$source_path" >&2
    exit 1
  fi
done < <(jq -r '.tested_on.source_artifacts[]? | [.path, .sha256] | @tsv' "$proof_states")

for state_name in "${required_states[@]}"; do
  count="$(jq --arg name "$state_name" '[.states[]? | select(.name == $name)] | length' "$proof_states")"
  if [[ "$count" != "1" ]]; then
    echo "E2E pilot package must contain exactly one proof state named $state_name" >&2
    exit 1
  fi
  state_value="$(jq -r --arg name "$state_name" '.states[] | select(.name == $name) | .state' "$proof_states")"
  if [[ "$mode" == "complete" && "$state_value" != "observed" ]]; then
    echo "E2E pilot package proof state $state_name must be observed for complete proof packages" >&2
    exit 1
  fi
  evidence_count="$(jq --arg name "$state_name" '.states[] | select(.name == $name) | (.evidence_refs // []) | length' "$proof_states")"
  if [[ "$state_value" == "observed" && "$evidence_count" -lt 1 ]]; then
    echo "E2E pilot package proof state $state_name must have at least one evidence_ref" >&2
    exit 1
  fi
  if [[ "$state_value" != "observed" ]]; then
    reason="$(jq -r --arg name "$state_name" '.states[] | select(.name == $name) | .reason // empty' "$proof_states")"
    next_required_evidence="$(jq -r --arg name "$state_name" '.states[] | select(.name == $name) | .next_required_evidence // empty' "$proof_states")"
    if [[ -z "$reason" || -z "$next_required_evidence" || "$next_required_evidence" == "null" ]]; then
      echo "E2E pilot package non-observed proof state $state_name must include reason and next_required_evidence" >&2
      exit 1
    fi
  fi
done

invalid_state_count="$(jq '[.states[]? | select((.state == "observed" or .state == "not_observed" or .state == "not_assessed") | not)] | length' "$proof_states")"
if [[ "$invalid_state_count" != "0" ]]; then
  echo "E2E pilot package proof states contain invalid state values" >&2
  exit 1
fi

idx=0
while IFS= read -r event; do
  file="$tmp/evidence-event-$idx.json"
  printf '%s\n' "$event" >"$file"
  validate schema/evidence-event.schema.json "$file"
  idx=$((idx + 1))
done < <(jq -c '.[]' "$package_dir/evidence/evidence-events.json")
if [[ "$idx" -eq 0 ]]; then
  echo "E2E pilot package evidence-events.json must contain at least one event" >&2
  exit 1
fi

while IFS= read -r evidence_ref; do
  [[ -z "$evidence_ref" ]] && continue
  ref_count="$(jq --arg ref "$evidence_ref" '[.[] | select(.id == $ref)] | length' "$package_dir/evidence/evidence-events.json")"
  if [[ "$ref_count" != "1" ]]; then
    echo "E2E pilot package proof-state evidence_ref does not resolve to exactly one evidence event: $evidence_ref" >&2
    exit 1
  fi
  ref_status="$(jq -r --arg ref "$evidence_ref" '.[] | select(.id == $ref) | .status' "$package_dir/evidence/evidence-events.json")"
  if [[ "$ref_status" != "success" ]]; then
    echo "E2E pilot package observed proof-state evidence_ref must point to a success event: $evidence_ref" >&2
    exit 1
  fi
done < <(jq -r '.states[]? | (.evidence_refs // [])[]?' "$proof_states")

idx=0
while IFS= read -r record; do
  file="$tmp/provenance-record-$idx.json"
  printf '%s\n' "$record" >"$file"
  validate schema/provenance-record.schema.json "$file"
  idx=$((idx + 1))
done < <(jq -c '.[]' "$package_dir/evidence/provenance-records.json")
if [[ "$idx" -eq 0 ]]; then
  echo "E2E pilot package provenance-records.json must contain at least one record" >&2
  exit 1
fi

idx=0
while IFS= read -r observation; do
  file="$tmp/observation-$idx.json"
  printf '%s\n' "$observation" >"$file"
  validate schema/observation.schema.json "$file"
  idx=$((idx + 1))
done < <(jq -c '.[]' "$package_dir/evidence/observations.json")
if [[ "$idx" -eq 0 ]]; then
  echo "E2E pilot package observations.json must contain at least one observation" >&2
  exit 1
fi

idx=0
while IFS= read -r metric_stream; do
  file="$tmp/metric-stream-$idx.json"
  printf '%s\n' "$metric_stream" >"$file"
  validate schema/metric-stream.schema.json "$file"
  idx=$((idx + 1))
done < <(jq -c '.[]' "$package_dir/evidence/metric-stream.json")
if [[ "$idx" -eq 0 ]]; then
  echo "E2E pilot package metric-stream.json must contain at least one stream" >&2
  exit 1
fi

validate schema/trace.schema.json "$package_dir/evidence/trace-snapshot.json"
validate schema/assessment-input.schema.json "$package_dir/handoff/assessment-input.json"

echo "E2E pilot package valid: $package_dir"
