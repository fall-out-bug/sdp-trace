# Self-Trace v0 Summary

Date: 2026-05-01

Scope: Block 02 Self-Trace Proof for `specs/001-sdp-trace-time-series-evidence-substrate/`.

## Purpose

Self-Trace v0 records this repository's own contract-foundation work using the same portable `sdp-trace` contracts that downstream consumers are expected to consume.

This summary is not a gate verdict. It records what is assessed, what remains `not_assessed`, and which commands reproduce the package validation.

## Artifacts

- `examples/self-trace/evidence-events.json`
- `examples/self-trace/provenance-records.json`
- `examples/self-trace/observations.json`
- `examples/self-trace/metric-stream.json`
- `examples/self-trace/trace-snapshot.json`
- `examples/self-trace/assessment-input.json`
- `examples/self-trace/negative-native-policy-field.json`
- `scripts/validate-self-trace.sh`

## Commands

```bash
rtk scripts/validate-self-trace.sh
rtk npm run validate
```

Observed result on 2026-05-01: both commands exited 0 after the Block 01 resume digest fix.

## Assessed

- Evidence events record the metric catalog, contract validation command, and crisis-review critic/judge artifacts.
- Provenance records include Codex, `rtk` command execution, and `pi` reviewer/judge model metadata where available.
- Observations distinguish contract scaffolding from product proof and keep external attestation as `not_assessed`.
- Metric streams cover contract task completion, evidence coverage, `not_assessed` count, schema validation state, and review contradiction count.
- The assessment input contains no native pass/fail/readiness/degradation decision fields.
- The negative self-trace fixture proves a native policy field is rejected by `schema/assessment-input.schema.json`.
- Accountability is human-held through synthetic public roles for DRI, approver, risk owner, and escalation.

## Not Assessed

- Immutable source reference for the contract release remains Block 03 scope.
- External attestation remains `not_assessed`.
- Production release verification remains `not_assessed`.
- Fresh-checkout validation remains `not_assessed`; the current evidence was reproduced in the shared working tree. Block 03 must separate schema validity, digest verification, local attestation, external attestation, and production release verification.
- Customer pilot evidence remains out of scope and blocked until Block 03 records separate proof states.

## Review Inputs

The self-trace package records the crisis Socratic review artifacts as external review evidence:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/01-crisis-self-proof-kimi-critic.json`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/01-crisis-self-proof-glm-critic.json`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/01-crisis-self-proof-judge-result.json`

The artifact references and SHA-256 digests are recorded in `examples/self-trace/evidence-events.json`; model/tool provenance is recorded in `examples/self-trace/provenance-records.json`.

Block 02 is ready for clean-context pi code review only after `scripts/validate-self-trace.sh` and `npm run validate` both exit 0.
