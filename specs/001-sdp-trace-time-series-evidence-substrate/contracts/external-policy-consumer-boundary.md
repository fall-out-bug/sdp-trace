# Contract: sdp-trace And External Policy Consumers

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
| External verdict input records | `sdp-trace` records them as external evidence only |
| Evidence-strength or quality assertions | External producer; `sdp-trace` records producer and origin |
| Schema IDs, versions, and compatibility notes | `sdp-trace` |
| Artifact hashes, redaction notes, and integrity status | `sdp-trace` |
| Accountability records for trace artifacts | `sdp-trace` records accountable human-held roles and approval refs |
| Risk and oversight classification metadata | `sdp-trace` records autonomy, impact, oversight, and review-independence metadata |
| Contract manifests and artifact digests | `sdp-trace` |
| Contract release verification records | `sdp-trace` records verification evidence and status |
| Trusted signer identity policy shape | `sdp-trace` |
| Accepted release signer identity policy values | Repository or customer governance process |
| Accepted risk tolerance | external policy consumer, management, or another consuming governance process |
| Policy thresholds | external policy consumer or another consuming policy engine |
| Pass/fail/warn/block decisions | external policy consumer or another consuming policy engine |
| Degradation verdicts | external policy consumer or another consuming policy engine |
| Readiness decisions | external policy consumer or another consuming policy engine |
| Overrides and override reasons | external policy consumer or another consuming policy engine |

## Contract Rules

1. External policy consumers may depend on `sdp-trace` contracts.
2. `sdp-trace` must not depend on an external policy consumer; this means no runtime import, required service call, policy configuration, or gate-specific execution path. A schema designed for downstream consumption is a consumer contract, not a runtime dependency.
3. `sdp-trace` may record an external verdict as evidence, but must identify it as externally produced with producer, origin, verdict kind, source artifact, and policy reference when available.
4. `sdp-trace` must use `not_assessed` when required evidence is missing.
5. `sdp-trace` must not publish opaque health scores.
6. `sdp-trace` must preserve dimensions needed by policy engines: scope, repository, team, harness, model, stack, build system, and time window.
7. Beads may track implementation work, but Beads must not appear as a required product dependency.
8. `sdp-trace` must not assign evidence strength or quality verdicts. If a source supplies such a value, `sdp-trace` records it as an external assertion.
9. `sdp-trace` must keep movement data structural: previous value, current value, delta, units, dimensions, evidence coverage, and `not_assessed` gaps. Interpretation labels are external verdicts.
10. `sdp-trace` contracts use semver schema versions. External policy consumers must declare which schema versions they consume.
11. `sdp-trace` must record human accountability for accountable artifacts. AI actors may produce, review, critique, or judge, but cannot be sole accountable owners or approvers.
12. `sdp-trace` must distinguish schema validity from trusted contract release status. A checkout can be JSON-valid while still failing manifest digest or signature verification.
13. `sdp-trace` records observed risk metadata and externally declared oversight assertions; external policy consumer decides whether that metadata satisfies a policy.
14. `sdp-trace` may define the trusted identity policy schema and record policy values, but it does not decide business acceptance of signer risk.

## Structural Boundary Tests

A trace or assessment input violates this contract if a native `sdp-trace` field contains:

- `pass`, `fail`, `blocked`, `ready`, `not_ready`, or equivalent gate verdict as a native result
- `degrading`, `improving`, `healthy`, or equivalent movement interpretation as a native result
- threshold configuration that tells a policy consumer how to decide
- evidence strength assigned by `sdp-trace` rather than recorded as an external assertion
- an AI actor as the sole accountable owner, approver, risk owner, or escalation owner
- a `trusted_contract_release` claim when manifest digest or signature verification is missing, stale, invalid, or explicitly `not_assessed`
- a `trusted_contract_release` claim when signer identity does not match the trusted identity policy

A trace or assessment input is compliant when policy values appear only inside an external verdict input with explicit producer and origin, and oversight obligations appear only as externally declared assertions with policy reference.

## Schema Compatibility Rule

Every schema introduced by this feature must declare:

- JSON Schema Draft 2020-12
- stable `$id`
- semver `schema_version` in artifacts or package metadata once examples are validated

Additive optional fields are minor-version changes. Required field removals, enum semantic changes, or ownership-boundary changes are major-version changes. `schema/trace.schema.json` remains a compatibility path until a replacement path and migration note are committed.

## Example Flow

```text
spec -> plan -> task -> change -> evidence event -> accountability -> observation -> metric stream -> trace snapshot -> assessment input
contract manifest -> contract release verification -> assessment input
```

An external policy consumer consumes the assessment input:

```text
assessment input -> policy evaluation -> gate decision -> decision record
```

The second line is not owned by `sdp-trace`.
