# Quickstart: Reviewing the SpecKit Evidence Package

## 1. Start With the Spec

Read:

```text
specs/001-sdp-trace-time-series-evidence-substrate/spec.md
```

Confirm the feature answers the CTO question as evidence-backed process movement, not as a built-in policy verdict.

## 2. Read the Boundary Contract

Read:

```text
specs/001-sdp-trace-time-series-evidence-substrate/contracts/sdp-trace-sdp-gate-boundary.md
```

Confirm `sdp-trace` owns evidence/provenance/observations/metric streams and `sdp-gate` owns policies, gate decisions, degradation verdicts, readiness, and overrides.

## 3. Inspect the Plan and Tasks

Read:

```text
specs/001-sdp-trace-time-series-evidence-substrate/plan.md
specs/001-sdp-trace-time-series-evidence-substrate/tasks.md
```

Beads issues can mirror execution state, but these SpecKit artifacts are the repo-observable plan.

## 4. Run Current Schema Syntax Check

```bash
jq empty schema/*.json
```

This checks current schema JSON syntax. Full JSON Schema validation is a planned task.

## 5. Inspect Pilot Evidence After Runs

Expected sanitized outputs after pilot execution:

```text
docs/research/
examples/opencode/
examples/superpowers/
examples/jvm-bazel/
```

Raw local outputs may live under:

```text
.sdp-trace-runs/
```

That path is intentionally ignored by git.
