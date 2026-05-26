# Spec 019 Monitoring And Gate Proof Pack

Status: reproducible local proof pack.

This pack demonstrates `sdp-trace` as explicit-export monitoring input for LLM
and harness activity, and as evidence food for advisory gates. It does not claim
merge approval, release approval, production readiness, or audit-grade trust.

## What This Proves

- A wrapped command produces a local run manifest and replayable event chain.
- `verify`, `query`, `report`, and `gate` can consume the run directory.
- Local gate facts can pass while CI witness and audit-grade evidence remain
  `cannot_verify`.
- Harness observation reads explicit JSONL exports only.
- Missing required harness event families remain `not_assessed`.
- Unsafe raw prompt input is rejected before an observed run is written.
- Telemetry export can render posture facts as Prometheus text for downstream
  monitoring systems.

## Replay

Run from the repository root:

```bash
(
  set -euo pipefail
  TMPDIR=$(mktemp -d)
  trap 'rm -rf "$TMPDIR"' EXIT

  go run ./cmd/sdp-trace wrap --name spec019-proof --output-dir "$TMPDIR/run" -- /bin/echo ok
  go run ./cmd/sdp-trace verify "$TMPDIR/run"
  go run ./cmd/sdp-trace query --query missing-evidence "$TMPDIR/run"
  go run ./cmd/sdp-trace report --out "$TMPDIR/report" "$TMPDIR/run"
  go run ./cmd/sdp-trace gate --out "$TMPDIR/gate-result.json" "$TMPDIR/run" || true
  jq '{local_gate, ci_witness_gate, audit_grade_gate, gate_mode}' "$TMPDIR/gate-result.json"
)
```

Expected gate states:

```json
{
  "local_gate": "pass",
  "ci_witness_gate": "cannot_verify",
  "audit_grade_gate": "cannot_verify",
  "gate_mode": "observation"
}
```

Harness degradation replay:

```bash
(
  set -euo pipefail
  TMPDIR=$(mktemp -d)
  trap 'rm -rf "$TMPDIR"' EXIT

  go run ./cmd/sdp-trace harness observe \
    --profile examples/spec019-monitoring-gate-proof/harness-profile.json \
    --source examples/spec019-monitoring-gate-proof/harness-events-missing-model.jsonl \
    --out "$TMPDIR/harness-run"

  go run ./cmd/sdp-trace harness validate \
    --profile examples/spec019-monitoring-gate-proof/harness-profile.json \
    --run "$TMPDIR/harness-run" \
    --out "$TMPDIR/harness-validation.json"

  jq '{validation_state, reason_code}' "$TMPDIR/harness-validation.json"
)
```

Expected validation state:

```json
{
  "validation_state": "not_assessed",
  "reason_code": "required_event_family_absent"
}
```

Unsafe input replay:

```bash
(
  set -euo pipefail
  TMPDIR=$(mktemp -d)
  trap 'rm -rf "$TMPDIR"' EXIT

  if go run ./cmd/sdp-trace harness observe \
    --profile examples/spec019-monitoring-gate-proof/harness-profile.json \
    --source examples/spec019-monitoring-gate-proof/harness-events-unsafe-raw-prompt.jsonl \
    --out "$TMPDIR/unsafe-run"; then
    echo "ERROR: unsafe raw prompt was accepted"
    exit 1
  fi

  test ! -e "$TMPDIR/unsafe-run"
)
```

Telemetry replay:

```bash
go run ./cmd/sdp-trace export telemetry \
  --profile prometheus-text-v1 \
  --cross-repo-posture examples/block21-cross-repo-posture/valid-movement/cross-repo-posture-export.json \
  --out -
```

Telemetry output is monitoring data only. Alert rules, dashboards, risk
thresholds, and policy decisions remain downstream.
