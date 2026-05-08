# Block 29: Interaction Trace And Friction Metrics

Status: Draft spec. Implementation is blocked until Socratic review is complete
and the reviewed direction is explicitly approved.

Parent artifacts:

- `docs/process-metric-catalog.md`
- `docs/flight-recorder.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/24-demo-repo-ci-trace-pilot.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/28-repo-observer-install-doctor.md`

## Goal

Make task interaction itself an observable trace surface.

If a human gives a task to an agent, then later corrects scope, plan, evidence,
tooling, or delivery behavior, `sdp-trace` should be able to retain that
correction as part of the task trace. The product value is not "store a note."
The value is to preserve the interaction history that explains why the work
changed and how much friction the task required.

This block introduces a portable interaction trace contract for:

- initial task assignment;
- clarifying questions and answers;
- corrective feedback after task assignment;
- plan approval or rejection;
- implementation pause or resume decisions;
- agent/model/tool drift notes;
- friction metrics derived from those events.

## Product Question

"When an agent worked on this task, what interaction events changed the work,
and how much corrective friction was required before the task reached its
current state?"

This question matters because friction is a product signal:

- high correction count may indicate a weak spec, weak agent, weak model,
  missing setup, or poor operator discipline;
- repeated boundary corrections may indicate unsafe prompt-only workflow;
- plan rejections may indicate poor planning or hidden requirements;
- stale evidence corrections may indicate weak trace literacy;
- tool/model drift corrections may explain why two runs of the same task are not
  comparable.

`sdp-trace` must expose these facts without deciding whether the agent, model,
spec, or employee is "good" or "bad." Interpretation belongs to downstream
review, coaching, governance, or reporting systems.

## Dogfood Finding

The demo dogfood exposed a failure mode:

1. The user gave corrective feedback after the task had started.
2. The feedback changed the expected behavior of the agent and GSD workflow.
3. The repo had no product surface that would make that correction obvious to a
   later reviewer.
4. A manual "feedback event" primitive would still depend on agent goodwill and
   therefore would repeat the same failure in a different shape.

The right product surface is not a prompt instruction and not a manual note.
It is a traceable interaction channel or an explicitly imported transcript whose
source, completeness, and limitations are represented in machine-readable state.

## Non-Goals

- No generic chat app, prompt manager, conversation UI, or agent runtime.
- No hidden dependence on Codex, Claude, OpenCode, GSD, Beads, GitHub, Slack,
  Telegram, Zoom, Fireflies, or any specific connector.
- No background spying on private conversations.
- No claim that imported transcript text is complete unless the source provides
  a verifiable completeness boundary.
- No claim that the agent complied with feedback merely because feedback was
  observed.
- No grading, scoring, ranking, coaching recommendation, employment evaluation,
  or model leaderboard.
- No raw secrets, token-like values, private filesystem paths, authenticated
  URLs, or unnecessary personal data in retained interaction events.
- No Node.js, npm, JavaScript, TypeScript, or `.mjs` active product path.

## Product Boundary

`sdp-trace` may say:

- "this task has an observed assignment event";
- "this task has three corrective feedback events after assignment";
- "this correction was observed through source type
  `observed-control-channel`";
- "this transcript import is partial and therefore `not_assessed` for
  completeness";
- "the agent compliance outcome remains `not_assessed` even when the
  correction is referenced by plan, code, evidence, or review artifacts";
- "the task had N friction events by type during this window."

`sdp-trace` must not say:

- "the agent obeyed the correction";
- "the employee handled the task badly";
- "the model is worse than another model";
- "the spec is bad";
- "the work is approved";
- "the task is done."

## Task Identity

Interaction metrics are only useful if every event belongs to a stable task.
Block 29 therefore introduces a lightweight `task_id` contract.

MVP `task_id` rules:

- `task_id` is a safe identifier supplied by the caller or derived from a
  committed task registry.
- A derived id must be stable across runs and must not be derived from private
  local paths.
- If the source cannot provide or derive a stable task id, interaction trace
  assessment is `not_assessed`.
- If two different task descriptions claim the same `task_id`, the trace is
  `cannot_verify`.

Preferred future artifact:

```text
.sdp-trace/tasks.json
```

The registry is out of scope for the first schema implementation, but the schema
must not permit free-text task labels to masquerade as durable task identity.

## Interaction Sources

The interaction trace contract supports three source classes. Only the first two
can produce observed interaction events in this block.

### Observed Control Channel

An observed control channel is a wrapper or integration point where an
interaction message is retained by `sdp-trace` before it is delivered to the
agent or workflow.

Required properties:

- stable `task_id`;
- stable `interaction_id`;
- source-local ordered event sequence;
- actor type aligned with the common provenance taxonomy: `human`, `model`,
  `tool`, or `system`;
- target role or runtime when applicable;
- event type;
- timestamp supplied by the observing runtime;
- retention policy;
- digest over retained content or redacted content;
- source identity and source version;
- explicit source trust scope.

An observed control channel can support stronger claims than a manual event
because the message path is:

```text
human message -> sdp-trace observation -> agent/runtime delivery
```

The observation still does not prove compliance.

Observed control channels can be bypassed unless the surrounding runtime or
workflow enforces exclusive use of that channel. If exclusivity is not enforced
or cannot be inspected, the trace may claim `observed_before_delivery=true` for
events that passed through the relay, but it must keep
`channel_exclusivity_state=not_assessed`.

### Transcript Import

A transcript import is a retained conversation export from an external source.
It is useful when `sdp-trace` did not sit in the live message path.

Required properties:

- source type and source reference;
- import time;
- source export time if available;
- completeness state: `complete`, `partial`, `not_assessed`, or
  `cannot_verify`;
- event ordering basis;
- digest over imported content or redacted imported content;
- importer identity when available.

Transcript import is evidence of retained transcript material. It is not proof
that no missing corrective messages exist unless the source can verify
completeness.

### Agent-Reported Interaction

Agent-reported interaction is a claim by an agent, prompt, plan, README, or
summary that a correction occurred or was handled.

Agent-reported interaction is not an accepted source for interaction events in
this block. It may be represented only as `not_assessed` context or as an
unsupported claim in another trace surface. This prevents the product from
depending on the same prompt goodwill that failed during dogfood.

## Source Type Registry

MVP source types are closed:

| Source type | Meaning | Trust boundary |
| --- | --- | --- |
| `observed-control-channel` | `sdp-trace` recorded the event before forwarding it to the next process. | Can assert `observed_before_delivery=true`; cannot assert compliance or exclusivity unless separately evidenced. |
| `preclassified-transcript-import` | A structured transcript import already labels events using the Block 29 event vocabulary. | Post-hoc source; completeness depends on transcript source metadata. |
| `agent-reported` | Agent/prose says an interaction happened. | Rejected for event creation; record only as `not_assessed`/unsupported context. |

Unknown source types are `cannot_verify`.

## Completeness Decision Table

| Scenario | Completeness state | Trace state |
| --- | --- | --- |
| Observed relay wrote event before delivery and write succeeded. | `complete` for that relayed event only. | `observed` |
| Observed relay failed before delivery. | `cannot_verify` | `cannot_verify` |
| Transcript source declares complete export with verifiable source boundary. | `complete` | `observed` |
| Transcript source declares partial export. | `partial` | `observed` with stream-level `partial` |
| Transcript source provides no completeness metadata. | `not_assessed` | `not_assessed` for completeness |
| Transcript ordering is missing, duplicate, or contradictory. | `cannot_verify` | `cannot_verify` |
| Agent summary claims feedback existed but no accepted source exists. | `not_assessed` | no interaction event |

Completeness is about source coverage. It is not compliance proof.

## Event Types

MVP event types:

| Event type | Meaning | Friction class |
| --- | --- | --- |
| `task_assignment` | Initial task/request for the scoped work. | none |
| `clarification_request` | Agent or reviewer asks for missing requirement/context. | clarification |
| `clarification_answer` | Human or source answers clarification. | clarification |
| `plan_proposed` | Agent/workflow proposes plan before implementation. | planning |
| `plan_approved` | Human approves plan or reviewed spec direction. | none |
| `plan_rejected` | Human rejects plan or requires substantial correction. | correction |
| `corrective_feedback` | Human correction after assignment that changes task path, boundary, or evidence expectations. | correction |
| `boundary_violation` | Human or reviewer identifies prompt-only, scope, evidence, or authority violation. | correction |
| `evidence_correction` | Human or reviewer corrects stale, missing, weak, or overclaimed evidence. | evidence |
| `tool_or_model_drift` | Actor records changed model/tool availability or behavior affecting comparability. | drift |
| `pause_requested` | Human pauses or blocks implementation. | coordination |
| `resume_approved` | Human resumes work after pause/block. | coordination |

Event type values are closed in the first implementation. Unknown values are
`cannot_verify`, not silently accepted.

An event carries exactly one `event_type` and one primary `friction_class`.
When a message appears to match multiple event types, the classifier must use
this priority order:

1. `boundary_violation`
2. `evidence_correction`
3. `plan_rejected`
4. `corrective_feedback`
5. coordination events
6. clarification events
7. neutral planning or assignment events

Secondary labels are out of scope for MVP. This keeps metrics reproducible.

## State Model

Interaction trace state is separate from task completion, code quality, and
proof state.

Minimum states:

- `observed`: event exists in an accepted source class and validates against the
  schema;
- `not_assessed`: source was not available or completeness was not assessed;
- `cannot_verify`: source exists but is malformed, unsafe, contradictory, or
  lacks required ordering/content integrity;
- `redacted`: event validates, but content is redacted under a declared policy;
- `referenced`: event has stable references to plan, code, evidence, review,
  PR, or CI artifacts;
- `unreferenced`: event is retained but not linked to downstream work.

`observed` and `referenced` are factual states. `referenced` means only that a
stable id links two artifacts. It is not a causal claim that the interaction
caused the downstream change. Causal impact remains `not_assessed` unless a
future profile explicitly implements impact assessment.

## Ordering

Every event has:

- `source_id`
- `source_sequence`
- `observed_at`
- optional `source_created_at`

`source_sequence` is monotonic within a single source. Duplicate or
non-monotonic sequence values in one source are `cannot_verify`.

Metric computation uses `observed_at` for global ordering across sources. When
post-hoc transcript import introduces events whose source timestamps precede
already-observed events, the metric stream remains usable but must set
`assessment_state=partial` unless the source provides a verifiable complete
ordering boundary.

## Metrics

Block 29 adds interaction-friction metrics to the process metric catalog during
implementation.

Candidate metric names:

| Metric name | Unit | Meaning |
| --- | --- | --- |
| `observed_interaction_event_count` | count | Count of retained interaction events by type and source. |
| `observed_post_assignment_correction_count` | count | Count of `corrective_feedback`, `boundary_violation`, and `evidence_correction` events after `task_assignment`. |
| `observed_plan_rejection_count` | count | Count of rejected plans before approved execution. |
| `observed_clarification_turn_count` | count | Count of clarification request/answer events. |
| `unreferenced_interaction_event_count` | count | Count of observed events not linked to downstream artifacts. |
| `interaction_source_completeness_state` | enum | Whether the interaction source is complete, partial, not assessed, or cannot verify. |

Required dimensions when available:

- `task_id`
- `spec_id`
- `phase_id`
- `agent_runtime`
- `model_family`
- `model_version`
- `source_type`
- `source_version`
- `event_type`
- `friction_class`
- `window_start`
- `window_end`

Metrics must not encode judgement. A high correction count is a signal for
review, not a built-in verdict.

Metric names come from the catalog. `dimensions` provide scope breakdown; they
are not an alternate metric namespace. The implementation must not emit a
separate metric family per event type.

Aggregation guardrail: raw correction counts are comparable only when task
scope, time window, source type, and completeness state are comparable. If those
dimensions differ or are missing, comparisons are `not_assessed`.

## Schema Contract

Implementation should add JSON schema artifacts for:

- `interaction-event.schema.json`
- `interaction-trace.schema.json`
- optional `interaction-metric-stream.schema.json` only if existing
  `metric-stream.schema.json` cannot represent the needed facts.

Minimum interaction event fields:

- `schema_version`
- `interaction_id`
- `task_id`
- `source_id`
- `source_sequence`
- `event_type`
- `friction_class`
- `actor`
- `target`
- `source`
- `content_ref`
- `content_digest`
- `digest_algorithm`
- `retention`
- `state`
- `reference_refs`
- `observed_before_delivery`
- `channel_exclusivity_state`
- `completeness_state`
- `not_retained_reason`
- `redaction`
- `observed_at`
- `created_at`

`digest_algorithm` is `sha256`.

`content_ref` must use a closed reference scheme:

- `evidence:<ref>` for retained content already represented as an evidence ref;
- `sdp://interaction/<task_id>/<interaction_id>` for a retained local
  content-addressed blob owned by this trace package;
- `external:<safe-ref>` for external transcript material that is not copied into
  the package.

Unknown schemes are `cannot_verify`. If content is not retained, `content_ref`
is absent and `not_retained_reason` is required.

## Retention And Redaction

Interaction content can contain secrets and personal data. The implementation
must fail closed rather than persist unsafe content.

MVP retention policy:

- Redaction happens before retained content is written.
- Retained content digest covers the retained post-redaction content.
- The event records whether content was `full`, `redacted`, or `not_retained`.
- `redaction.policy_ref` identifies the policy used.
- `redaction.finding_count` records how many redactions occurred without
  echoing sensitive values.
- Actor identity uses role or safe id. Raw names, emails, private handles, and
  employee ids are not required and must not be emitted unless a source profile
  explicitly declares a safe identity policy.
- If the implementation detects a token-like value, private path, authenticated
  URL, or unsupported personal data category that it cannot redact safely, the
  event is `cannot_verify` and content is not retained.

## CLI Shape

The first useful implementation must include a minimal relay. Import-only MVP
is rejected because it recreates the manual ritual that failed during dogfood.

Preferred MVP commands:

```text
sdp-trace interaction relay \
  --task-id <safe-id> \
  --actor-type human \
  --target <safe-id> \
  --event-type <closed-event-type> \
  --out <trace.json> \
  -- <forward-command> [args...]

sdp-trace interaction import-transcript \
  --source preclassified-transcript-import \
  --source-ref <safe-ref> \
  --task-id <safe-id> \
  --events-jsonl <file> \
  --out <trace.json>

sdp-trace interaction summarize \
  --trace <trace.json> \
  --out <summary.json>
```

`relay` reads message content from stdin, validates and records the event, then
forwards the same content to the supplied command on stdin. If event recording
fails, the message is not forwarded. This gives the product one runtime-agnostic
observed control channel without depending on any specific agent.

`import-transcript` accepts pre-classified JSONL events only. Free-form chat
parsing is out of scope for MVP. A source-specific parser can be added only with
fixtures and explicit `not_assessed` behavior for ambiguous lines.

Both commands are trace-package writers. They do not claim the agent complied
with the message.

## Acceptance Criteria

- Spec defines the difference between observed control channel and transcript
  import, and explicitly rejects agent-reported interaction as event evidence.
- Spec states that prompt cooperation and manual self-reporting are not
  sufficient.
- MVP includes a minimal relay that records before forwarding.
- Schema validates ordered task interaction events with closed event types.
- Schema rejects unknown source types and unknown content reference schemes.
- Transcript import supports explicit `complete`, `partial`, `not_assessed`,
  and `cannot_verify` completeness states.
- Corrective feedback after task assignment is represented as a first-class
  event and counted in friction metrics.
- Friction metrics are factual counts/states, not quality verdicts.
- Metric comparison is `not_assessed` when scope, source, or completeness
  dimensions are not comparable.
- Event content retention and redaction behavior is explicit and safe.
- Existing task/provenance/evidence records can link to interaction events by
  stable ids.
- Focused tests cover relay write-before-forward, correction-after-assignment,
  pre-classified transcript import, transcript partiality, malformed ordering,
  unsafe content refs, redaction refusal, agent-reported rejection, and
  unreferenced events.

## Review Questions

Socratic review must answer:

1. Does the spec create a real product surface, or another manual note ritual?
2. Can a user tell whether a correction was observed before agent delivery,
   imported later, or merely agent-reported?
3. Are friction metrics useful without becoming employee/model scoring?
4. Does the source/completeness model avoid overclaiming transcript coverage?
5. Is privacy/redaction concrete enough for real task conversations?
6. Does the MVP have a narrow implementation path that does not depend on a
   specific agent runtime?
7. What remains `not_assessed` after MVP?

## Socratic Review Ledger

Initial review verdict: `REVISE`.
Focused re-review verdict after fixes: `APPROVE`; no remaining critical or major
findings.

| id | severity | plane | finding | disposition |
| --- | --- | --- | --- | --- |
| S29-PROD-01 | critical | product/UX | Import-only MVP repeats the manual note ritual. | accepted_fixed: MVP now requires `interaction relay` write-before-forward. |
| S29-TRUST-01 | critical | trust boundary | Observed channel did not state bypass/exclusivity limits. | accepted_fixed: added `channel_exclusivity_state=not_assessed` unless externally evidenced. |
| S29-TRUST-02 | critical | trust boundary | `correlated` could be misread as causation. | accepted_fixed: renamed to `referenced` and states causal impact remains `not_assessed`. |
| S29-ARCH-01 | critical | existing contracts | Stable `task_id` did not exist as a contract. | accepted_fixed: added task identity rules and future registry boundary. |
| S29-ARCH-02 | critical | existing contracts | Content storage did not integrate with evidence refs. | accepted_fixed: added closed `content_ref` schemes and `digest_algorithm=sha256`. |
| S29-PRIV-01 | major | privacy | Redaction authority and retention behavior were unspecified. | accepted_fixed: added retention/redaction policy, actor identity constraints, and fail-closed rule. |
| S29-METRIC-01 | major | metrics | Friction counts could become unsupported quality scores. | accepted_fixed: renamed metrics to observed facts and added comparison guardrails. |
| S29-SOURCE-01 | major | source model | Agent-reported interaction was not explicitly rejected. | accepted_fixed: added third source class rejected for event creation. |
| S29-ORDER-01 | major | ordering | Multi-source ordering was undefined. | accepted_fixed: added `source_sequence`, `observed_at`, and partial-state rule. |
| S29-CLASS-01 | major | classification | Transcript classifier authority and overlapping event types were ambiguous. | accepted_fixed: MVP accepts pre-classified JSONL only and defines event priority order. |

## Approval Boundary

This block is a spec only. Implementation must not start until:

- Socratic review is complete;
- critical and major findings are fixed or explicitly blocked;
- the user approves the reviewed direction.
