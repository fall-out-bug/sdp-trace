# Wrap Output / Schema Drift Example

Status: `cannot_verify`
Spec: [017](../../../specs/017-oss-replacement-compatibility-and-benchmarks/)

This directory demonstrates that the live `sdp-trace wrap` output does not
conform to `schema/flight-recorder-run.schema.json`.

## Frozen Snapshot

The file `run.json` in this directory is a **frozen historical snapshot** of
`sdp-trace wrap /bin/true` output captured during spec 017 development. It is
not regenerated on every build; it exists as structural evidence that wrap
output once differed from `schema/flight-recorder-run.schema.json`.

```text
run_dir: .sdp-trace-runs/run-3068560305
```

> **Note:** The run ID (`run-3068560305`) is nondeterministic across invocations.
> This fixture is frozen as structural evidence; the exact run ID does not affect
> the drift claim, which depends only on the plain-text format. If `sdp-trace wrap`
> output changes in the future, this snapshot does not automatically update.

## Schema Expectation

`schema/flight-recorder-run.schema.json` requires a JSON object with these
required fields:

- `schema_version` (const `"1.0.0"`)
- `run_id`
- `profile`
- `trust_scope`
- `artifact_role`
- `created_at` (ISO-8601 timestamp)
- `source_summary`
- `task_summary`
- `model_summary`
- `harness_summary`
- `evidence_retention_summary`
- `verifier_states`

## Drift Summary

| Dimension | Live Output | Schema Requirement |
|---|---|---|
| Format | Plain text line | JSON object |
| `schema_version` | Absent | Required `"1.0.0"` |
| `run_id` | Absent | Required non-empty string |
| `created_at` | Absent | Required ISO-8601 timestamp |
| `verifier_states` | Absent | Required array |
| `event_refs` / `event_chain_head` | Absent | At least one required (`anyOf`) |
| `witness_ref` + `event_chain_head` | Absent | Required for witnessed profiles (`allOf`) |

## Blocker Status

This drift is documented as a blocker for flight-recorder schema compatibility.
Resolving it requires either:

1. Updating `sdp-trace wrap` to emit `flight-recorder-run.schema.json`-compliant
   JSON, or
2. Defining a separate "current recorder schema" that matches the actual wrap
   output and versioning it independently.

Until one of these options is implemented and accepted via a spec update,
this example must remain in the repository as structural evidence of the
mismatch. Do not delete this directory or mark the drift resolved without a
source-bound proof commit that changes either the wrap output or the schema.
