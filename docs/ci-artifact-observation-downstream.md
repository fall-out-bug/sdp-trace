# CI Artifact Observation Downstream Use

`sdp-trace` records CI artifact observation facts. It does not enforce merge,
release, readiness, audit, or risk decisions.

Downstream policy consumers may read
`ci-artifact-observation.schema.json` results and decide how to handle them.
That decision belongs outside `sdp-trace`.

## State Semantics

| State | Downstream meaning |
| --- | --- |
| `pass` | The selected profile observed every required family at the required proof level. A downstream policy may still require more checks. |
| `fail` | Contradictory or unsafe evidence was observed. A downstream gate may block, but `sdp-trace` only records the fact. |
| `cannot_verify` | Required evidence was absent, partial, expired, inaccessible, malformed, unverifiable, or below required producer authority. Treat as no usable proof for that profile. |
| `not_assessed` | The selected profile did not inspect that surface. Do not infer either success or failure. |

## Producer Scopes

Per-family `producer_scope` is more important than the top-level aggregate.
`ci_uploaded` can satisfy a selected CI-uploaded requirement. `checked_in`,
`local_generated`, `agent_reported`, and `harness_observed` remain
lower-authority facts unless a downstream policy explicitly accepts them.

The top-level `producer_scope` is an aggregate for required families:

- one shared required scope is rendered directly;
- multiple required scopes are rendered as `mixed`;
- no required family is `not_assessed`.

## Binding Scopes

`source_binding_state` and `run_binding_state` describe selected source/run
binding. `producer_binding_state` describes whether the selected producer
authority was met. A checked-in artifact can have matched source/run metadata
and still have `producer_binding_state: mismatch` when the profile required
`ci_uploaded`.

## Safe Reasons

Reason strings and next actions are safe templates derived from closed reason
codes. They must not echo raw logs, prompts, model output, token-like material,
private URLs, or private filesystem paths.
