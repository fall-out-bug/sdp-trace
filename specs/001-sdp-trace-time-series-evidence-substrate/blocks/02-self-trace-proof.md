# Block 02: Self-Trace Proof

Status: implemented; pi review passed
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Beads mirror: `sdp-trace-cdn.12`
Audience: CTO, CIO, CEO, implementation agents, future `sdp-gate` consumers

## Purpose

Self-Trace Proof makes `sdp-trace` its own first real consumer.

The block answers the CTO objection:

> If `sdp-trace` cannot prove its own development, why should anyone trust it for customer work?

## Executive Outcome

The repository must contain a committed self-trace package that describes this feature's own development through the same contracts downstream consumers will use.

This is not a gate verdict. It is a replayable evidence package showing what was specified, planned, changed, checked, reviewed, left `not_assessed`, and who is accountable.

## In Scope

- `examples/self-trace/evidence-events.json`
- `examples/self-trace/provenance-records.json`
- `examples/self-trace/observations.json`
- `examples/self-trace/metric-stream.json`
- `examples/self-trace/trace-snapshot.json`
- `examples/self-trace/assessment-input.json`
- negative self-trace fixture proving native policy fields fail validation
- self-trace validation command
- sanitized summary in retired research artifact

## Required Evidence

The self-trace package must include:

- SpecKit spec, plan, task, and block references
- changed artifact references
- command evidence for validation and safety checks
- Socratic critic, resolution, and judge artifacts as external review evidence
- provenance for Codex, `pi`, model provider, tool, and command where available
- human-held accountability for repository proof ownership
- metric streams for task completion, evidence coverage, `not_assessed` count, schema validity, and review contradiction count
- explicit `not_assessed` entries for missing immutable commit, external attestation, production signing, or unavailable raw logs

## Out of Scope

- customer pilot execution
- `sdp-gate` policy decisions
- production Sigstore/Rekor release proof

## Acceptance

Block 02 passes only when `examples/self-trace/assessment-input.json` validates from a fresh checkout and contains no native pass/fail/readiness/degradation decision.

If Block 02 fails, customer pilot work remains blocked.
