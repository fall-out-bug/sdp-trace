# Block 19: Adapter Event Contract And Capture Depth

Status: spec delta and implementation plan; awaiting explicit approval before
implementation.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/09-flight-recorder-trust-kernel.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13-product-gap-closure-roadmap.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/17-managed-harness-enforcement-profile.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/18-redaction-retention-forensic-profiles.md`
- `docs/flight-recorder.md`

## Goal

Expose agent SDLC provenance that CI and git cannot reconstruct by themselves,
while making capture depth explicit enough that missing adapter coverage is a
visible fact rather than an implied forensic guarantee.

The correct framing is not "capture everything from every harness." That is not
portable, safe, or true. The correct framing is "define the smallest generic
adapter contract that can report high-value lifecycle, task, tool, file, model,
test, and closure facts; then make every unsupported, missing, not-integrated,
or unsafe capture point machine-visible."

## Problem

Current flight-recorder and managed-harness artifacts can evaluate local command
and run evidence, but they still leave important gaps:

- internal harness tool calls can be invisible to CI and git;
- model identity can be self-reported or harness-reported without gateway
  proof;
- file mutation evidence can exist without actor or command correlation;
- task drift and supersession need actor attribution;
- test claims can be confused with executed test evidence;
- PR/MR and review references need a portable shape that is not GitHub-bound;
- missing adapter events must be first-class query output;
- prompt and raw model-response capture is sensitive and must stay unavailable
  by default unless Block 18 policy explicitly permits richer retention.
- adapter event metadata fields can leak secrets if identity, reference, and
  payload fields are not covered by the same redaction rules as run events.

The product risk is overclaim. Block 19 must not say a reviewer can reconstruct
events that were never captured, retained, or bound to the selected run.

## Non-Goals

- No dependency on OpenCode, GSD, Bazel, Kotlin, GitHub, Claude, Codex, or any
  specific harness, model provider, Git host, CI, review, or gateway runtime.
- No default raw prompt, raw model response, stdout/stderr body, command argv,
  source snippet, tool input/output body, adapter configuration, gateway
  evidence body, credential, token, adapter secret, gateway token, PR token, or
  raw review body persistence.
- No native merge, release, readiness, degradation, override approval, audit,
  legal, incident, or risk-acceptance decision.
- No claim that adapter events prove the adapter implementation is honest,
  complete, externally audited, or tamper-proof.
- No replay engine for deterministic re-execution of agent actions. Query output
  reconstructs available evidence, not an executable replay.
- No arbitrary evidence labels outside the current claim-tag and schema rules.

## Product Boundary

Block 19 may emit verifier-derived adapter and capture-depth facts:

- adapter event shape validity;
- adapter producer and identity state;
- event family capture state;
- task lock and task supersession actor attribution;
- tool-call correlation to command, file, test, or model observations;
- model identity provenance state;
- test observation provenance state;
- source, VCS, PR/MR, and review reference binding;
- gateway provenance profile state, including `not_integrated`;
- missing, unsupported, late, or unsafe adapter capture;
- retention and redaction limits inherited from Block 18;
- query-facing rows for task drift, task supersession counts, unverified task
  expansion indicators, unverified claims, missing adapter events, unsupported
  observers, and capture-depth caps.

Block 19 must not decide whether an organization accepts the residual risk.
External CI, incident, compliance, release, and policy consumers may use these
facts.

## Binding Model

Block 19 supports two binding modes. Every adapter event must use exactly one:

1. `same_chain`: the adapter event is written as an ordered flight-recorder
   event in the same Block 09 recorder chain. It must carry monotonic
   `sequence`, `prev_event_hash`, `event_hash`, run id, and run nonce.
2. `adapter_bundle`: the adapter writes a separate canonical event bundle with
   its own ordered event sequence and `adapter_bundle_head_digest`. The
   flight-recorder chain must reference the bundle digest before `run_closed`,
   and `run_closed` must summarize the bundle digest and event count.

For managed-profile use, `same_chain` or a pre-run bundle reference bound by a
managed witness is required. A post-hoc adapter bundle can support local
capture-depth inspection only; it cannot satisfy Block 17 managed enrollment or
required managed event coverage.

`run_binding_established` verifies the selected binding mode. It fails on run
id, run nonce, source, sequence, chain-head, or bundle-head contradiction. It is
`cannot_verify` when a required chain or bundle digest is inaccessible.

An adapter event is late when it appears after the recorder `run_closed`
sequence, when an `adapter_bundle` is first referenced only after closure, or
when event timestamps contradict the chain or bundle order. Late events are
visible capture-depth evidence only and cannot satisfy required managed or
executed-evidence conditions.

## Adapter Event Contract

Block 19 should add a stable adapter event contract for these event families:

- `run_started`;
- `task_locked`;
- `task_superseded`;
- `tool_call`;
- `command_started`;
- `file_mutation`;
- `model_call_observed`;
- `test_observed`;
- `run_closed`.

These are adapter event names. The existing flight-recorder event schema may
either version/extend current event types or bind adapter-event bundles into the
run. The implementation must preserve compatibility with existing
`*_observed` flight-recorder names where they already exist.

Required common fields:

- schema version;
- event id;
- event type;
- event time;
- producer identity;
- adapter identity;
- provenance scope;
- capture state;
- run id;
- run nonce;
- source baseline, source commit, or tree digest when available;
- binding mode;
- monotonic sequence and hash linkage for `same_chain`, or adapter bundle id,
  bundle sequence, and bundle head digest for `adapter_bundle`;
- correlation refs;
- event payload digest;
- redaction policy ref and retention mode where payload classes are sensitive;
- missing, unsupported, or not-assessed reason when capture is incomplete.

Initial provenance scopes:

- `ci_executed`;
- `wrapper_executed`;
- `gateway_observed`;
- `harness_observed`;
- `adapter_observed`;
- `agent_reported`;
- `local_observed`;
- `not_integrated`;
- `unsupported`;
- `not_assessed`;
- `cannot_verify`.

Initial capture states:

- `captured`;
- `missing_telemetry`;
- `not_integrated`;
- `unsupported`;
- `redacted`;
- `retention_limited`;
- `not_assessed`;
- `cannot_verify`.

`captured` only means the adapter event row is present and schema-valid. It
does not mean the adapter was authorized, continuously connected, externally
witnessed, or sufficient for a managed profile. Those remain separate managed
or witness conditions.

Multiple adapters may emit events for the same correlation key. If two events
with the same event family and correlation key contradict source, run, actor,
result, or digest fields, the capture-depth condition for that family is
`cannot_verify` with reason `conflicting_adapter_events`. Non-conflicting
duplicates are retained as separate observations and must not be double-counted
as independent proof.

`provenance_scope` and `capture_state` are independent taxonomies even when they
share values such as `not_integrated`. `provenance_scope` describes who or what
produced the observation. `capture_state` describes whether the expected
evidence family was captured, missing, unsupported, or retention-limited.

## Redaction And Sensitive Metadata

Block 19 adapter events inherit Block 18 forbidden committed-artifact
persistence classes. The following adapter fields are sensitive unless a Block
18 policy explicitly classifies them otherwise:

- tool input and output payloads;
- command argv, stdout, stderr, and working-directory values that could expose
  private paths or source;
- model request, model response, prompt, and gateway evidence bodies;
- adapter configuration;
- producer identity, adapter identity, gateway identity, gateway evidence refs,
  source refs, change refs, and review refs when they contain credentials,
  tokens, authenticated URLs, or raw review bodies.

Committed adapter artifacts must not contain raw sensitive payloads under
`safe_default`. They may contain digests, sanitized excerpts, or Block
18-conforming raw references only when the selected redaction policy permits
that retention mode. A retention-limited reference in Block 19 must conform to
the Block 18 `raw_reference` shape or an explicitly versioned equivalent with
digest binding, access state, access verification time, key custody state where
applicable, retention lifecycle, and unavailable reason.

Provider-neutral refs must be canonical, token-free, and credential-free before
persistence. Embedded credentials, signed URLs, bearer tokens, session ids, and
raw review bodies are forbidden in committed adapter artifacts and in
preview/query/explain output.

## Event Family Semantics

`run_started` records run id, run nonce, selected recorder or adapter profile,
source baseline, producer, and parent run refs when available.

`task_locked` records task/spec/plan digest, actor attribution, authority ref
when available, and lock provenance. Actor attribution is self-reported unless
it is bound to an authenticated provenance or accountability reference.

`task_superseded` records the superseded task lock, replacement task digest,
actor attribution, reason code, and whether supersession happened before or
after dependent command/model/test evidence. Unbound actor attribution must not
be treated as authenticated identity.

`tool_call` records the generic tool family, optional adapter-local tool label,
tool call id, input/output payload digests or Block 18-conforming
retention-limited refs, start/end state when available, and correlation to
command, file, model, or test observations. Stable product behavior must depend
on generic tool families, not harness-specific adapter-local tool labels.
Adapter-local labels may be retained as sanitized metadata but are not stable
contract members and must not appear in product logic or generic fixtures.

`command_started` records command correlation id, command descriptor digest,
working-directory or workspace digest, and safe argv retention metadata. Raw
argv is unavailable by default. Minimum argv retention metadata is retention
mode, argv input digest, redacted payload digest when a sanitized excerpt is
retained, and redaction policy ref.

`file_mutation` records path scope, source baseline, tree or diff digest,
attributed command/tool/task refs, actor attribution state, and source commit or
tree digest correlation. Source baseline and run id correlation are required
for `file_mutation_correlated: pass`; missing source or run binding is
`cannot_verify`.

`model_call_observed` records model identity provenance, harness identity,
gateway evidence ref when present, request/response digest or Block
18-conforming retention-limited refs, and capture state. Committed adapter
artifacts must not contain raw prompt or raw response payloads under
`safe_default`; only digests or Block-18-conforming retention-limited refs are
permitted unless an explicit Block 18 policy permits richer retention.

`test_observed` records test command/tool correlation, result state, execution
environment digest when available, artifact digest, and test provenance. Test
claims from an agent can be recorded, but never as executed test proof.

`run_closed` records closure state, chain head or adapter bundle digest, missing
evidence summary, retention/capture-depth caps, and late/unsupported observer
summary.

## Test Provenance Semantics

Test observation provenance values:

- `ci_executed`: CI or external runner evidence bound to source, run, command,
  artifact digest, and freshness where available.
- `wrapper_executed`: registered process wrapper or tool execution evidence
  with command descriptor, exit state, execution environment digest when
  available, and artifact digest.
- `harness_observed`: harness or adapter saw a test intent or event but lacks
  independent execution evidence.
- `agent_reported`: agent text or self-report claims tests ran or passed.
- `cannot_verify`: required provenance evidence is inaccessible or
  contradictory.

Only `ci_executed` and `wrapper_executed` can satisfy executed-test evidence for
Block 19. `harness_observed` may correlate intent and timing, but cannot become
executed-test evidence unless bound to wrapper or CI execution evidence.
`agent_reported` is visible but non-upgrading.

## Gateway Provenance

Gateway provenance is optional and can remain `not_integrated`.

When gateway evidence exists, it must bind:

- gateway identity;
- model identity;
- request or response digest / retention-limited reference;
- run id and run nonce when available;
- source or task correlation when available;
- gateway event time;
- gateway authority or deployment reference when available;
- redaction/retention profile.

Without gateway evidence, adapter-reported model identity remains
`harness_observed` or `agent_reported`. The verifier must not infer gateway
authority from a model name, harness name, adapter id, file path, environment
variable, or prose.

Adapter and producer identity have two levels:

- `self_asserted`: identity fields are present but not bound to an executable
  digest, deployment manifest, signature, managed policy, or external
  attestation.
- `bound`: identity is bound to an executable digest, deployment manifest,
  signature, managed policy, or external attestation accepted by the selected
  profile.

`adapter_identity_visible` can pass for capture-depth inspection when identity
is visible as `self_asserted`, but managed, forensic, or gateway-observed
claims require `bound` identity where the selected profile requires authority.

## Provider-Neutral References

Source, VCS, PR/MR, and review references must use provider-neutral fields:

- `source_ref`;
- `source_commit` or `source_tree_digest`;
- `change_ref`;
- `review_ref`;
- `review_state`;
- `producer`;
- `artifact_ref` or artifact digest;
- `observed_at`;
- `not_assessed_reason` when any reference is unavailable.

Provider-neutral refs define a minimal abstract shape for source, change, and
review identity. Adapter profiles may add provider-specific extension fields,
but those fields are not stable contract members unless a later block adds them
explicitly. The schema must not require GitHub-specific names or a PR/MR-shaped
workflow. Examples may use generic values only, such as
`provider: generic_git_host`, not real vendor names as product concepts.

## Capture-Depth Assessment

Block 19 should add capture-depth facts, either in a new assessment-result
version or a dedicated query/assessment output that remains read-compatible with
Block 17 and Block 18 result shapes.

Required condition groups:

| Condition | Required behavior |
| --- | --- |
| `adapter_event_contract_valid` | `pass` when adapter events match the selected schema; `fail` for contradictory or malformed fields; `cannot_verify` when required bundle binding is inaccessible. |
| `adapter_identity_visible` | `pass` when producer and adapter identity are present and classified as `self_asserted` or `bound`; `cannot_verify` when identity cannot be resolved; `not_assessed` when no adapter profile was selected. |
| `run_binding_established` | `pass` when adapter events bind to run id and source/run context; `fail` for contradictory run/source binding; `cannot_verify` for missing binding. |
| `task_drift_visible` | `pass` when task locks were assessed and either no supersessions were observed (`no_supersessions_observed`) or supersessions include actor attribution and digest refs; `cannot_verify` when supersession evidence is incomplete; `not_assessed` when task-drift assessment was not selected or implemented. |
| `tool_call_depth_visible` | `pass` when required tool-call families are captured or explicitly unsupported; `missing_telemetry` when expected events are missing; `unsupported` when the adapter declares no capability. |
| `file_mutation_correlated` | `pass` when file mutation evidence correlates to source baseline and run id; `cannot_verify` when source or run binding is unavailable. |
| `model_identity_not_overclaimed` | `pass` when model identity provenance remains within available evidence; `fail` when harness or agent identity is claimed as gateway-observed; `not_integrated` when no gateway exists. |
| `test_provenance_not_overclaimed` | `pass` when executed-test evidence is CI or wrapper executed; `cannot_verify` for inaccessible execution proof; `fail` when agent-reported or harness-observed claims are emitted as executed tests. |
| `provider_refs_portable` | `pass` when source/PR/review refs use provider-neutral fields; `fail` when product code or schemas require a specific Git host. |
| `redaction_metadata_consistent` | `pass` when sensitive adapter fields carry Block 18 redaction policy refs, retention modes, and required raw-reference metadata; `fail` when forbidden raw fields are persisted; `cannot_verify` when required redaction metadata is inaccessible. |
| `capture_depth_not_overclaimed` | `pass` when each event family summary with insufficient evidence emits a visible cap such as `missing_telemetry`, `unsupported`, `not_integrated`, `not_assessed`, `cannot_verify`, `retention_limited`, or `capped_to_retention_mode`; `fail` when any event family emits `pass` or `reconstructable` while its effective evidence is digest-only, missing, unsupported, not integrated, not assessed, or retention-limited without a cap annotation. |

Top-level capture-depth state should use `pass`, `fail`, `cannot_verify`, or
`not_assessed`. `sdp-trace assess --profile adapter-capture` is an assessment
profile: condition rows with `missing_telemetry`, `not_integrated`,
`unsupported`, or `retention_limited` map to top-level `cannot_verify` unless a
dominant `fail` exists. `sdp-trace query capture-depth` is a read-only summary:
it does not emit a top-level pass/fail assessment, does not define exit-code
policy, and must preserve the condition-level states as observed summary rows.

## Block 17 And Block 18 Interaction

Block 19 is independent from Block 17 managed-harness enforcement by default.
`adapter-capture` evaluates capture depth and overclaim boundaries; it does not
authorize an adapter, prove managed enrollment, or satisfy managed witness
binding by itself.

When a later implementation composes Block 19 with Block 17:

- `adapter_identity_visible` may inform diagnostics, but
  `adapter_identity_authorized` remains the Block 17 authority condition;
- `tool_call_depth_visible` may provide evidence rows for
  `required_tool_events_observed`, but only when the binding model satisfies
  Block 17's pre-run enrollment and witness requirements;
- `test_provenance_not_overclaimed` may explain why
  `test_provenance_not_agent_reported` fails or cannot verify, but executed
  test evidence still requires CI or registered wrapper/tool execution proof;
- contradictory Block 17 and Block 19 facts must not be reconciled by prose.
  The stricter failing or `cannot_verify` condition remains visible in each
  profile's own output.

Block 18 owns redaction and retention profile semantics. Block 19 must carry
enough redaction/retention metadata for adapter-sensitive fields so Block 18
can verify them, and `redaction_metadata_consistent` prevents capture-depth pass
from hiding missing safety metadata.

## CLI And UX

Preferred command surfaces:

```text
sdp-trace assess --profile adapter-capture ...
sdp-trace assess preview --profile adapter-capture ...
sdp-trace assess explain <assessment-result.json>
sdp-trace query capture-depth ...
```

`assess --profile adapter-capture` is a new Block 19 assessment surface that
consumes a run artifact plus same-chain adapter events or an adapter bundle
digest. It does not replace Block 14 `gate`, Block 17 `managed-harness`, or
Block 18 `forensic-retention`. `query capture-depth` is a read-only summary
surface for the same facts and is not a policy decision. If Block 20 later
changes the query command family, Block 19 must preserve read compatibility or
update this spec before implementation.

Preview must show expected event families, adapter identity requirements,
capture-depth caps, and missing inputs without printing raw sensitive payloads.
Query output must be useful to a reviewer without raw JSON spelunking, but it
must label the boundary: available evidence is not full replay.

The UX must make these trade-offs visible:

- adapter events increase observation depth but do not prove adapter honesty;
- gateway `not_integrated` keeps model identity below gateway-observed state;
- `harness_observed` test events are correlation evidence, not execution proof;
- safe default retention may cap forensic reconstruction;
- provider-neutral refs preserve portability but may be less ergonomic than
  host-specific deep links.
- task supersession and unverified task expansion are observable conditions;
  "scope creep" remains a human judgment outside native `sdp-trace` output.

## Test And Fixture Expectations

Implementation requires tests and committed fixtures for:

- valid generic adapter event bundle;
- missing required adapter event;
- unsupported observer capability;
- gateway `not_integrated`;
- late adapter attach or missing run binding;
- agent-reported test claim does not become executed evidence;
- harness-observed test correlation without wrapper/CI proof remains
  non-executed;
- CI or wrapper executed test evidence can satisfy executed-test provenance;
- file mutation evidence correlates with source commit/tree digest and run id;
- task supersession includes actor attribution and links to the prior lock;
- provider-neutral PR/MR and review refs validate without a named Git host;
- query output exposes task drift, task supersession counts, unverified task
  expansion indicators, unverified claims, missing adapter events, unsupported
  observers, and retention/capture-depth caps;
- safety-sensitive output does not print raw command args, stdout/stderr
  bodies, prompts, source snippets, tool_call input/output, adapter
  configuration, gateway evidence refs, model request/response payloads,
  credentials, OIDC request tokens, adapter secrets, gateway tokens, model
  responses, PR tokens, or raw review bodies;
- product code and fixtures do not encode demo-specific harness, model, Git
  host, or build-system names as product concepts.

## Review Plan

Spec review must use separate planes:

- requirements-vs-product-boundary review for capture-depth scope,
  provider-neutrality, and no forensic-complete overclaim;
- tracing/evidence review for adapter event binding, model/test provenance,
  task drift, file mutation correlation, and PR/review refs;
- privacy/security review for prompt/response retention, adapter/gateway
  secrets, PR tokens, and safe query/preview output;
- code/correctness review after implementation for schema behavior, verifier
  states, fixtures, CLI/query output, and compatibility with Blocks 17 and 18.

PR-level review repeats the code/correctness, tracing/evidence,
requirements-vs-implementation, and privacy/security planes. Absent CI remains
`not_assessed`.

Review disposition must be recorded in
`specs/001-sdp-trace-time-series-evidence-substrate/blocks/19-adapter-event-contract-capture-depth-review-ledger.md`
with finding id, severity, review plane, finding, disposition, and evidence
reference. Spec-level critical and major findings must be closed or explicitly
blocked before implementation approval.

## Acceptance Criteria

- AC1: Spec, tasks, FRs, and SCs define a provider-neutral adapter event
  contract without making any harness, model provider, Git host, CI, gateway, or
  build system a product dependency.
- AC2: Adapter capture-depth output makes missing, unsupported,
  not-integrated, not-assessed, cannot-verify, late, and retention-limited
  states explicit.
- AC3: Adapter-reported model identity cannot be upgraded to gateway-observed
  without bound gateway evidence.
- AC4: Agent-reported and harness-observed test claims cannot satisfy executed
  test evidence without CI or wrapper/tool execution proof.
- AC5: File mutation evidence can be correlated with source baseline and run id
  without requiring a specific Git host.
- AC6: Query-facing output exposes task drift, task supersession counts,
  unverified task expansion indicators, unverified claims, unsupported
  observers, missing adapter events, and capture-depth caps without claiming
  replay, scope-creep judgment, or forensic completeness beyond retained
  evidence.
- AC7: Default prompt and raw model-response capture remains unavailable unless
  explicit Block 18 retention/redaction policy permits richer retention.
- AC8: Output remains verifier facts only; no native merge, release, readiness,
  degradation, audit, legal, incident, override, or risk-acceptance decision is
  introduced.

## Implementation Slices

1. Schema and contract slice: adapter event schema, provider-neutral refs,
   binding mode fields, Block 18 redaction metadata, capture-depth condition
   rows, and schema parse fixtures.
2. Run binding slice: run id, nonce, chain/bundle digest, source baseline, task
   lock/supersession, and file mutation correlation.
3. Provenance slice: model identity, gateway profile, test provenance semantics,
   and non-upgrading agent/harness claims.
4. Capture-depth evaluator slice: missing/unsupported/not-integrated states,
   retention caps, deterministic top-level state, and exit behavior.
5. Query and safety slice: query capture-depth output, preview/explain text,
   secret-leak negative tests, and provider-neutral fixture docs.

Block 19 must preserve Block 17 managed harness semantics and Block 18
redaction/retention safety. Implementation does not start until this spec
direction is explicitly approved.
