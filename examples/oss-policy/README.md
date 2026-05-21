# OSS Policy Prototype

Status: experimental; verification not_assessed/cannot_verify as listed below
Spec: [017](../../specs/017-oss-replacement-compatibility-and-benchmarks/)

This directory contains a minimal OPA/Rego policy prototype for a
simplified adapter-capture profile. It demonstrates that OPA can express
sdp-trace-like rules, but it does not replace the product verifier.

## Files

| File | Purpose |
|---|---|
| `adapter.rego` | Simplified pass/fail rule with trace_id and provenance bounds |
| `test-fixture.json` | Valid input that should produce `pass = true` |
| `test-fixture-fail.json` | Combined invalid input (both rules broken) |
| `test-fixture-fail-traceid.json` | Invalid trace_id only (number instead of string) |
| `test-fixture-fail-provenance.json` | Invalid provenance only (overclaimed length) |

## Run the Policy

Requires `opa` in `$PATH`. The prototype uses `import rego.v1`, which requires OPA v1.0+ (or a recent v0.x with Rego v1 support). Run from the repository root.

```bash
(
  set -e
  cd examples/oss-policy || exit 1
  RESULT=$(opa eval --data adapter.rego \
    --input test-fixture.json \
    --format raw \
    'data.sdp_trace.adapter.pass')
  if [ "$RESULT" != "true" ]; then
    echo "ERROR: expected true, got: $RESULT"
    exit 1
  fi
)
```

## Test the Failure Fixtures

### Combined failure
```bash
(
  set -e
  cd examples/oss-policy || exit 1
  RESULT=$(opa eval --data adapter.rego \
    --input test-fixture-fail.json \
    --format raw \
    'data.sdp_trace.adapter.pass')
  if [ "$RESULT" != "false" ]; then
    echo "ERROR: expected false, got: $RESULT"
    exit 1
  fi
)
```

Expected: `false` (non-string trace_id and overlong provenance).

### Trace ID rule only
```bash
(
  set -e
  cd examples/oss-policy || exit 1
  RESULT=$(opa eval --data adapter.rego \
    --input test-fixture-fail-traceid.json \
    --format raw \
    'data.sdp_trace.adapter.pass')
  if [ "$RESULT" != "false" ]; then
    echo "ERROR: expected false, got: $RESULT"
    exit 1
  fi
)
```

Expected: `false` (trace_id is a number).

### Provenance rule only
```bash
(
  set -e
  cd examples/oss-policy || exit 1
  RESULT=$(opa eval --data adapter.rego \
    --input test-fixture-fail-provenance.json \
    --format raw \
    'data.sdp_trace.adapter.pass')
  if [ "$RESULT" != "false" ]; then
    echo "ERROR: expected false, got: $RESULT"
    exit 1
  fi
)
```

Expected: `false` (provenance exceeds the bound).

## Substitution Boundary

- **What OPA replaces:** Policy-as-code expressions for simplified profiles.
- **What remains sdp-trace-specific:** Evidence collection, gate verdicts,
  `sdp-trace-claim` tag semantics, and hash-chain validation.
- **Adapter glue required:** JSON translation layer between sdp-trace events
  and OPA input. OPA does not natively understand trace provenance or
  recorder profiles.
