# Block 20: Forensics Query Pack

Status: reviewed spec delta and implementation plan; awaiting explicit approval
before implementation.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/09-flight-recorder-trust-kernel.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/18-redaction-retention-forensic-profiles.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/19-adapter-event-contract-capture-depth.md`
- `docs/flight-recorder.md`

## Goal

Make common forensic reviewer questions executable and replayable without
turning `sdp-trace` into an incident, legal, release, or risk decision engine.

The product outcome is a portable query pack that answers "what can I safely
inspect, what is missing, and which claims are over the available evidence?"
from already recorded run, retention, redaction, adapter, witness, and verifier
facts.

## Problem

Blocks 18 and 19 made retention, redaction, and capture depth verifier-visible,
but a reviewer still has to stitch together several JSON outputs by hand. That
creates three practical failures:

- investigators can miss an explicit `not_assessed`, `unsupported`,
  `retention_limited`, or `cannot_verify` gap;
- humans can confuse digest-bound existence with reconstructable evidence;
- teams can create ad hoc reports that imply forensic completeness or policy
  acceptance outside `sdp-trace`.

The weak framing would be "add a nicer forensic report." That is too vague and
too risky. The correct framing is "define versioned, safe, read-only query
packs whose rows are derived from existing verifier facts and whose limitations
are impossible to hide."

## Non-Goals

- No new raw payload capture, transcript capture, prompt capture, model-response
  capture, stdout/stderr body capture, source-snippet capture, or tool
  input/output body capture.
- No native incident severity, audit result, legal hold, release readiness,
  merge readiness, degradation, override approval, or risk-acceptance decision.
- No dependency on GitHub, GitLab, Gerrit, Jira, SIEM, KMS, cloud object store,
  customer audit platform, or any named harness/model provider.
- No natural-language summarizer that can invent facts not present in recorded
  evidence.
- No query-pack pass/fail score, health score, trust score, or forensic
  completeness badge.
- No arbitrary evidence labels outside current schema and claim-tag rules.

## Product Boundary

Block 20 may emit query-derived facts:

- selected query pack id, query pack version, and input artifact digests;
- evidence availability by event family, artifact family, and claim family;
- reconstruction state for evidence that the selected Block 18 policy classifies
  as critical;
- capture-depth state using Block 19 adapter facts;
- witness and chain integrity state using Block 09 facts;
- redaction and retention issues requiring reviewer attention;
- task supersession and unverified task expansion indicators;
- unverified claim rows and the condition ids that block stronger claims;
- chronological incident timeline rows assembled from safe metadata only;
- evidence gap descriptions for `not_assessed` and `cannot_verify` rows.

Block 20 must not decide whether the organization accepts the residual risk.
External incident, audit, compliance, release, and legal consumers may use query
pack output as evidence input.

## Query Pack Model

Query packs are versioned named bundles of read-only queries. A pack is not an
assessment profile and must not produce a top-level pass/fail policy verdict.

Initial query pack:

- `forensics-basic-v1`: reviewer-oriented pack over one selected run. It
  combines safe rows for run summary, chain and witness state, retention and
  redaction issues, capture depth, command and test timeline using opaque
  command/test identifiers only, file mutation evidence, task supersession, and
  unverified claims.

Initial query names:

| Query | Required upstream artifacts | Optional upstream artifacts | Missing required artifact state |
| --- | --- | --- | --- |
| `forensics-summary` | Block 09 run facts | Block 18 forensic-retention result; Block 19 adapter-capture result | `cannot_verify` only when Block 09 run facts cannot be parsed; otherwise missing optional artifacts become referenced `not_assessed` rows |
| `forensics-timeline` | Block 09 run facts | Block 18 forensic-retention result; Block 19 adapter-capture result | `cannot_verify` |
| `forensics-gaps` | Block 09 run facts | Block 18 forensic-retention result; Block 19 adapter-capture result | `cannot_verify` |
| `forensics-redactions` | Block 09 run facts; Block 18 forensic-retention result | none | `cannot_verify` with evidence family `redaction` or `retention` |
| `forensics-capture-depth` | Block 09 run facts; Block 19 adapter-capture result | Block 18 forensic-retention result | `cannot_verify` with evidence family `adapter_capture` |
| `forensics-unverified-claims` | Block 09 run facts | Block 18 forensic-retention result; Block 19 adapter-capture result | `cannot_verify` only when Block 09 run facts cannot be parsed; otherwise missing optional artifacts become referenced `not_assessed` rows |

`forensics-summary` rows reference rows from the other query groups by row id
and must not independently re-derive evidence state.

The CLI surface should prefer an explicit pack command:

```bash
go run ./cmd/sdp-trace query-pack --pack forensics-basic-v1 \
  --run <run-dir> \
  --out <query-pack-result.json>
```

`query-pack` is read-only. `--pack` and `--out` are required for forensics query
packs because
stdout is commonly captured by CI logs, terminal scrollback, and shell pipelines.
A future explicit stdout mode may be added only if it prints a warning to stderr
and preserves the same safety guarantees. If the input is malformed or a
referenced verifier artifact cannot be parsed, `query-pack` returns
`cannot_verify` rows and deterministic non-zero exit behavior only when the
command cannot produce a valid query-pack result artifact.

If `--pack` is omitted, the command exits non-zero and writes
`error: ambiguous pack selection; --pack is required` to stderr, even when only
one pack is registered. If `--pack` is unknown, it exits non-zero and writes
`error: unknown pack "<value>"` to stderr. Partial query results that contain
`cannot_verify`, `not_assessed`, `unsupported`, `not_integrated`,
`missing_telemetry`, or `retention_limited` rows still exit zero when a valid
query-pack result artifact is written. Non-zero exit is reserved for usage
errors, unreadable required inputs that prevent artifact generation, or
serialization/write failure.

## Result Contract

The query-pack result should be a stable JSON artifact with:

- schema version;
- query pack id and version;
- selected run id, run nonce, and source baseline when available;
- input artifact SHA-256 digests when files are read; artifact references must
  be omitted or replaced with provider-neutral, credential-free,
  path-redacted identifiers already present in upstream Block 09, Block 18, or
  Block 19 artifacts. The result `input_artifacts` list is the union of every
  upstream file actually read by the pack;
- row arrays grouped by query name;
- row ids in the format `<query-short-name>.<NNNN>`, where the number is a
  zero-padded deterministic counter assigned per query group during assembly.
  Row ids must not use upstream event sequence numbers or shared counters;
- row evidence state using query vocabulary: `present`, `issue_observed`,
  `not_assessed`, `cannot_verify`, `missing_telemetry`, `unsupported`,
  `not_integrated`, and `retention_limited`;
- `evidence_family` using this closed vocabulary: `run_chain`, `witness`,
  `retention`, `redaction`, `adapter_capture`, `task`, `command`,
  `file_mutations`, `test`, `supersession`, `claim`, and `input_artifact`;
- `reconstructable` boolean. Digest-only evidence that exists but cannot
  reconstruct the selected event or artifact must set `reconstructable: false`;
- source condition ids and source condition states when rows derive from Block
  18 or Block 19 assessment facts. Upstream `pass` or `fail` values may appear
  only as source condition states, not as query-pack row verdicts;
- `source_ref` values using one of these formats:
  `block_09.run.<field>`, `block_09.event.<event_family>.<opaque_event_ref>`,
  `block_09.witness.<field>`, `block_18.condition.<condition_id>`, or
  `block_19.condition.<condition_id>`. Block 09 rows that do not have condition
  ids omit source condition fields and still include `source_ref`. The
  `block_09.event.<event_family>.<opaque_event_ref>` format covers task,
  command, file mutation, test, supersession, and claim facts through the row's
  closed `evidence_family` value;
- `evidence_gap` for `not_assessed` and `cannot_verify`, limited to
  coarse-grained evidence families already declared in the selected run,
  assessment, policy, or manifest artifacts;
- `output_safety` section listing the redaction policy id/digest when known and
  the sensitive classes that do not appear in the result.

`output_safety` is a verified absence assertion for the serialized result, not
a statement about whether sensitive material existed upstream. Each listed
sensitive class must come from the Block 20 safety list and must be checked
against the serialized JSON and explain output by deterministic tests.

Safe metadata for timeline assembly consists only of event timestamps, event
family identifiers, artifact family identifiers, opaque command/test ids,
run/chain/witness identifiers, query row ids, row evidence states, and source
condition refs. It excludes payload bodies, argument fields, command names,
executable paths, script paths, test names, response bodies, raw references,
provider URLs, private paths, PII, credentials, tokens, and raw content refs.

Row ids must be query-scoped deterministic identifiers, not internal event
sequence numbers or shared counters. That prevents avoidable cross-query joins
beyond the correlation refs already exposed by safe upstream facts.

Rows must be machine-readable first. Human-facing explain output may render
rows, but it must not add claims that are absent from the JSON result.

Explain output renders field labels, row evidence state, source condition refs,
and safe evidence-gap descriptions from the JSON result. It must not infer
conclusions, synthesize new facts, sort by hidden severity, encode state via
ANSI color, indentation depth, whitespace, or omitted sections, or summarize
rows with language absent from the result. Rows render in stable order by query
name and row id.

## Evidence Derivation Rules

Block 20 must derive facts from existing sources rather than invent new
authority:

- Block 09 run and verifier facts provide chain, witness, task, command, file,
  test, and supersession state.
- Block 18 forensic retention facts provide reconstruction and redaction state.
- Block 19 adapter capture facts provide capture-depth and unsupported observer
  state.
- Missing Block 18 or Block 19 artifacts referenced by the selected query are
  `not_assessed` unless the query table above marks that artifact as required;
  missing required artifacts produce `cannot_verify` rows with reason
  `missing_block_<block_number>_<artifact_kind>`.
- Digest-only evidence is evidence of existence, not reconstructability.
- `retention_limited`, `not_integrated`, `unsupported`, and `not_assessed` must
  remain visible in the final result and explain output.

State propagation is deterministic:

| Upstream fact pattern | Query row evidence state | Required row detail |
| --- | --- | --- |
| Block 09 chain, witness, event, or task fact is present with no issue condition | `present` | `source_ref`; `reconstructable` when applicable |
| Block 09 chain or witness issue is observed | `issue_observed` | `source_ref`; safe issue reason |
| Block 18 or Block 19 source condition state is `pass` | `present` | source condition id and state |
| Block 18 or Block 19 source condition state is `fail` | `issue_observed` | source condition id and state |
| Source condition state is `cannot_verify` | `cannot_verify` | source condition id and state; `evidence_gap` |
| Source condition state is `not_assessed` | `not_assessed` | source condition id and state; `evidence_gap` |
| Source condition state is `missing_telemetry` | `missing_telemetry` | source condition id and state; `evidence_gap` |
| Source condition state is `unsupported` | `unsupported` | source condition id and state; `evidence_gap` |
| Source condition state is `not_integrated` | `not_integrated` | source condition id and state; `evidence_gap` |
| Source condition state is `retention_limited` | `retention_limited` | source condition id and state; `reconstructable: false` when reconstruction is capped |
| Critical evidence is retained as digest-only under Block 18 facts | `retention_limited` | evidence family `retention`; reason `digest_only_not_reconstructable`; `reconstructable: false` |
| Required Block 18 or Block 19 upstream artifact is missing but a result artifact can still be written | `cannot_verify` | query-specific evidence family from the required-artifact table; reason `missing_block_<block_number>_<artifact_kind>` |
| Required upstream artifact is unreadable or malformed but a result artifact can still be written | `cannot_verify` | evidence family `input_artifact`; reason `unreadable_or_malformed_input_artifact` |
| Optional upstream artifact is absent | `not_assessed` | evidence family for the absent optional artifact |

Any unmapped upstream state or artifact shape is `cannot_verify`; the
implementation must not silently choose a default query row state.

The implementation must not recompute legal or release policy. It may report
that a downstream claim is unsupported because required evidence rows are
missing, retention-limited, or unverifiable.

## Safety Requirements

Forensics query output is safety-sensitive because it aggregates evidence from
multiple surfaces.

Query-pack JSON and explain output must not print:

- raw command arguments;
- command names, executable paths, script paths, and test identifiers unless
  they appear in a public non-sensitive catalog;
- stdout/stderr bodies;
- prompts;
- source snippets;
- tool-call input/output bodies;
- adapter configuration;
- gateway evidence refs;
- credentials, tokens, signed URLs, OIDC request tokens, adapter secrets,
  gateway tokens, PR tokens, or authenticated provider URLs;
- raw model request or response payloads;
- raw review bodies;
- raw-reference access notes that could expose storage location secrets;
- key material or key custody secrets.

Provider-neutral refs may appear only after canonical token-free and
credential-free normalization. If a row cannot be rendered safely, the row state
is `cannot_verify` or `retention_limited` with safe reason metadata.

## UX Requirements

The reviewer should not have to read raw JSONL manually for the common forensic
questions. The default result must be dense, deterministic, and scannable:

- group rows by question, not by internal schema type;
- include stable ids for rows and condition refs;
- name the evidence state before the evidence detail;
- make every missing or capped claim visible without scrolling through raw
  event payloads;
- preserve JSON output as the source of truth and keep explain output a
  rendering of that JSON.

The command must reject ambiguous pack selection. No default pack should run if
multiple packs are supported in the future.

Query groups emit flat row arrays only. They must not emit aggregate counts,
ratios, ranked lists, severity ordering, or implicit "critical" ordering. Any
future aggregate query requires a new query-pack version and safety review.

## Acceptance Criteria

- FR and SC deltas are added for versioned forensic query packs.
- `forensics-basic-v1` has a documented result contract and query list.
- `query-pack` is specified as read-only and separate from assessment profiles.
- Query-pack output preserves `not_assessed`, `cannot_verify`,
  `retention_limited`, `unsupported`, and `not_integrated` states.
- Digest-only evidence cannot be rendered as reconstructable evidence.
- Missing Block 18 or Block 19 artifacts have deterministic row states instead
  of disappearing from the result.
- Safety-sensitive output tests cover every sensitive class listed above.
- Safety tests must use synthetic fixture data only: documentation domains such
  as `example.com`, reserved IP space, and reserved token prefixes such as
  `example-token-`.
- Negative leak assertions must not echo candidate sensitive strings in failure
  output; use substring-hash, regex-class, or redacted-marker checks.
- Committed fixtures cover positive mixed evidence, digest-only cap, missing
  forensic-retention assessment, missing adapter-capture assessment,
  unsupported observer, redaction unresolved, task supersession, unverified
  claim, unsafe provider ref, and malformed input.
- A machine-checkable fixture matrix must enumerate every committed Block 20
  scenario, expected query group, expected row evidence state, expected
  `evidence_family`, expected `reconstructable` value when applicable, and
  expected source ref shape. Missing matrix rows are test failures.
- Tests must reject free-text evidence families, unmapped upstream states,
  non-conforming source refs, missing `input_artifacts` entries for files read,
  and explain output that reorders rows outside query-name and row-id order.
- Implementation review repeats code/correctness, tracing/evidence, and
  requirements-vs-implementation planes. PR-level review repeats those planes.

## Implementation Slices After Approval

1. Contract and fixtures: add query-pack schema/result examples and fixture
   matrix without changing CLI behavior.
2. Query derivation: add internal query-pack assembly over existing run,
   forensic, adapter, and verifier facts.
3. CLI and rendering: add `query-pack` command and explain rendering over the
   JSON result.
4. Safety and regression: add negative leak assertions and drift checks against
   Blocks 18 and 19.
5. Review and closure: run separate review planes, record dispositions, verify,
   commit scoped slices, open PR, repeat PR-level review, merge, and verify
   `origin/main`.
