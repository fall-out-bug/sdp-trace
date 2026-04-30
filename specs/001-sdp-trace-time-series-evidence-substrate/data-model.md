# Data Model: sdp-trace Time-Series Evidence Substrate

## Entity: Evidence Event

Represents one inspectable proof item.

Fields:

- `id`: stable event identifier
- `source`: producer such as `github-actions`, `opencode`, `manual`, `test-log`, or `file-inspection`
- `external_ref`: source-local reference such as run URL, command ID, file path, or review ID
- `observed_at`: when the source observed the event
- `collected_at`: when `sdp-trace` collected or recorded it
- `actor`: human, agent, system, or tool identity
- `event_type`: command, file, test, ci_run, review, scan, deployment, harness_log, model_output, manual, custom
- `status`: success, failure, warning, skipped, pending, not_assessed
- `summary`: short human-readable description
- `artifact_uri`: optional inspectable artifact reference
- `artifact_hash`: optional hash for artifact integrity
- `strength`: strong, medium, weak, none

Rules:

- Missing evidence must not be converted into success.
- Weak evidence can support observations, but consuming policy decides whether it is enough.

## Entity: Provenance Record

Represents the origin chain for an evidence event, observation, or trace snapshot.

Fields:

- `actor_id`
- `actor_type`
- `harness`
- `model_family`
- `model_version`
- `tool_name`
- `command`
- `prompt_ref`
- `context_refs`
- `captured_at`
- `hash_prev`
- `payload_digest`

Rules:

- Model and tool identity may be `not_assessed` when unavailable.
- Provenance records origin. It does not imply quality.

## Entity: Observation

Represents a dated, evidence-backed statement about process state.

Fields:

- `id`
- `scope`
- `observed_at`
- `statement`
- `evidence_refs`
- `provenance_refs`
- `assessment_status`: assessed, partial, not_assessed

Rules:

- Observations are not policy verdicts.
- Observations may be consumed by `sdp-gate`.

## Entity: Metric Sample

Represents one measured value for a process metric.

Fields:

- `metric_name`
- `value`
- `unit`
- `window_start`
- `window_end`
- `dimensions`
- `evidence_refs`
- `provenance_refs`
- `not_assessed_reason`

Common dimensions:

- repository
- team
- scope
- harness
- model_family
- model_version
- language
- build_system
- stack
- phase

Rules:

- Each sample must have evidence or `not_assessed_reason`.
- Thresholds and traffic-light ratings are external policy.

## Entity: Metric Stream

Represents ordered samples for the same metric and comparable dimensions.

Fields:

- `metric_name`
- `dimensions`
- `samples`
- `created_at`
- `updated_at`

Rules:

- Stream comparison can show movement.
- `sdp-trace` does not label movement as degradation unless an external policy verdict is recorded as evidence.

## Entity: Trace Snapshot

Represents a point-in-time graph across delivery artifacts.

Node kinds:

- spec
- plan
- task
- change
- evidence
- observation
- metric_sample
- external_verdict
- decision
- actor

Relation examples:

- task implements spec
- change fulfills task
- evidence supports observation
- metric_sample derived_from evidence
- external_verdict consumes assessment_input
- decision references external_verdict

## Entity: Assessment Input

Represents the package handed to `sdp-gate` or another policy engine.

Fields:

- `id`
- `scope`
- `trace_snapshot_ref`
- `evidence_bundle_refs`
- `metric_stream_refs`
- `observations`
- `not_assessed`
- `generated_at`

Rules:

- It must be usable without Beads.
- It must not contain policy decisions owned by `sdp-gate`.
