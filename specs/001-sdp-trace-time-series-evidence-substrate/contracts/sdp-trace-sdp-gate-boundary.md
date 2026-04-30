# Contract: sdp-trace and sdp-gate Boundary

## Purpose

This contract prevents product confusion between the trace substrate and the policy engine.

## Ownership

| Concern | Owner |
|---|---|
| Evidence events | `sdp-trace` |
| Provenance records | `sdp-trace` |
| Observations | `sdp-trace` |
| Metric samples and streams | `sdp-trace` |
| Trace snapshots | `sdp-trace` |
| Assessment input packages | `sdp-trace` |
| Policy thresholds | `sdp-gate` or another consuming policy engine |
| Pass/fail/warn/block decisions | `sdp-gate` or another consuming policy engine |
| Degradation verdicts | `sdp-gate` or another consuming policy engine |
| Readiness decisions | `sdp-gate` or another consuming policy engine |
| Overrides and override reasons | `sdp-gate` or another consuming policy engine |

## Contract Rules

1. `sdp-gate` may depend on `sdp-trace` contracts.
2. `sdp-trace` must not depend on `sdp-gate`.
3. `sdp-trace` may record an external verdict as evidence, but must identify it as externally produced.
4. `sdp-trace` must use `not_assessed` when required evidence is missing.
5. `sdp-trace` must not publish opaque health scores.
6. `sdp-trace` must preserve dimensions needed by policy engines: scope, repository, team, harness, model, stack, build system, and time window.
7. Beads may track implementation work, but Beads must not appear as a required product dependency.

## Example Flow

```text
spec -> plan -> task -> change -> evidence event -> observation -> metric stream -> trace snapshot -> assessment input
```

`sdp-gate` consumes the assessment input:

```text
assessment input -> policy evaluation -> gate decision -> decision record
```

The second line is not owned by `sdp-trace`.
