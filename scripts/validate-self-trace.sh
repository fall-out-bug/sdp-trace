#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

validate() {
  local schema="$1"
  local data="$2"
  node scripts/validate-json-schema.mjs "$schema" "$data"
}

validate_array_items() {
  local schema="$1"
  local data="$2"
  local name="$3"
  local idx=0

  jq -e 'type == "array" and length > 0' "$data" >/dev/null

  while IFS= read -r item; do
    local item_file="$tmpdir/${name}-${idx}.json"
    printf '%s\n' "$item" >"$item_file"
    validate "$schema" "$item_file"
    idx=$((idx + 1))
  done < <(jq -c '.[]' "$data")
}

expect_fail() {
  local schema="$1"
  local data="$2"
  if validate "$schema" "$data" >"$tmpdir/negative.out" 2>&1; then
    echo "Expected validation failure, but passed: $data" >&2
    cat "$tmpdir/negative.out" >&2
    exit 1
  fi
}

assert_no_native_policy_keys() {
  local data="$1"
  local found

  found="$(jq -r '
    def prohibited:
      [
        "verdict",
        "decision",
        "gate_result",
        "gate_verdict",
        "readiness",
        "readiness_verdict",
        "degradation_status",
        "policy_result",
        "policy_threshold",
        "evidence_strength",
        "quality_score",
        "override_result"
      ];
    def scan($path):
      if type == "object" then
        to_entries[]
        | (.key as $key
          | .value as $value
          | ($path + [$key]) as $next
          | (if (prohibited | index($key)) then ($next | map(tostring) | join(".")) else empty end),
            ($value | scan($next)))
      elif type == "array" then
        to_entries[]
        | (.key as $idx | .value | scan($path + [$idx]))
      else
        empty
      end;
    scan([])
  ' "$data")"

  if [[ -n "$found" ]]; then
    echo "Native policy keys are not allowed in $data:" >&2
    echo "$found" >&2
    exit 1
  fi
}

jq empty examples/self-trace/*.json

validate_array_items schema/evidence-event.schema.json examples/self-trace/evidence-events.json evidence-event
validate_array_items schema/provenance-record.schema.json examples/self-trace/provenance-records.json provenance-record
validate_array_items schema/observation.schema.json examples/self-trace/observations.json observation
validate_array_items schema/metric-stream.schema.json examples/self-trace/metric-stream.json metric-stream
validate schema/trace.schema.json examples/self-trace/trace-snapshot.json
validate schema/assessment-input.schema.json examples/self-trace/assessment-input.json
expect_fail schema/assessment-input.schema.json examples/self-trace/negative-native-policy-field.json
assert_no_native_policy_keys examples/self-trace/assessment-input.json
assert_no_native_policy_keys examples/self-trace/trace-snapshot.json

diff -u <(jq -S . examples/self-trace/evidence-events.json) <(jq -S '.evidence_events' examples/self-trace/assessment-input.json)
diff -u <(jq -S . examples/self-trace/provenance-records.json) <(jq -S '.provenance_records' examples/self-trace/assessment-input.json)
diff -u <(jq -S . examples/self-trace/observations.json) <(jq -S '.observations' examples/self-trace/assessment-input.json)
diff -u <(jq -S . examples/self-trace/metric-stream.json) <(jq -S '.metric_streams' examples/self-trace/assessment-input.json)
