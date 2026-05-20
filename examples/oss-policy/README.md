# OSS Policy Prototype

Status: locally tested, not externally verified
Spec: [017](../../specs/017-oss-replacement-compatibility-and-benchmarks/)

This directory contains a minimal OPA/Rego policy prototype for a
simplified adapter-capture profile. It demonstrates that OPA can express
sdp-trace-like rules, but it does not replace the product verifier.

## Files

| File | Purpose |
|---|---|
| `adapter.rego` | Simplified pass/fail rule with trace_id and provenance bounds |
| `test-fixture.json` | Valid input that should produce `pass = true` |
| `test-fixture-fail.json` | Invalid input that should produce `pass = false` |

## Run the Policy

Requires `opa` in `$PATH`. Run from the repository root.

```bash
(
  cd examples/oss-policy || exit 1
  opa eval --data adapter.rego \
    --input test-fixture.json \
    'data.sdp_trace.adapter.pass'
)
```

Expected output:

```json
{
  "result": [{
    "expressions": [{
      "value": true,
      "text": "data.sdp_trace.adapter.pass",
      "location": {"row": 1, "col": 1}
    }]
  }]
}
```

## Test the Failure Fixture

```bash
(
  cd examples/oss-policy || exit 1
  opa eval --data adapter.rego \
    --input test-fixture-fail.json \
    --format json \
    'data.sdp_trace.adapter.pass'
)
```

Expected: `false` (the fixture has a non-string trace_id and provenance exceeding the bound).

## Substitution Boundary

- **What OPA replaces:** Policy-as-code expressions for simplified profiles.
- **What remains sdp-trace-specific:** Evidence collection, gate verdicts,
  `sdp-trace-claim` tag semantics, and hash-chain validation.
- **Adapter glue required:** JSON translation layer between sdp-trace events
  and OPA input. OPA does not natively understand trace provenance or
  recorder profiles.
