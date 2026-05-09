# Block 31: Non-Interfering Harness Observation

Status: Draft spec. Implementation is blocked until Socratic review is complete
and the reviewed direction is explicitly approved.

Parent artifacts:

- `docs/reviews/demo-jvm-gsd-observation-ledger.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/19-adapter-event-contract-capture-depth.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/28-repo-observer-install-doctor.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/29-interaction-trace-and-friction-metrics.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/30-automated-pr-review-evidence-mechanism.md`

## Goal

Make external harness work observable without turning `sdp-trace` into the
harness, prompt router, or repo mutator.

The immediate dogfood pressure comes from the OpenCode + GSD + MiniMax JVM demo:
the demo cannot honestly claim `sdp-trace` observed the agentic delivery loop
when the product can only diagnose repo observer setup, wrap command execution,
or relay prompts through `sdp-trace`.

The product answer should be:

```text
For this task run, these harness events were imported or observed from an
external harness source, these fields were present, these fields were unavailable,
and these claims remain not_assessed or cannot_verify.
```

## Product Question

"Can `sdp-trace` observe an external agent harness delivery cycle without
changing the harness prompt stream, hand-authoring trace artifacts after the
fact, or depending on one harness runtime?"

The answer must include:

- harness identity and version, when available;
- model route or declared model identity, when available;
- task, phase, review, and stage boundaries;
- prompt or interaction boundary digests without requiring raw prompt retention;
- tool and command boundaries;
- file mutation references;
- test and validation observations;
- PR, review, and merge state references when present;
- explicit `pass`, `fail`, `not_assessed`, or `cannot_verify` for every required
  observation dimension.

## Dogfood Finding

`docs/reviews/demo-jvm-gsd-observation-ledger.md` records P0-001: no
non-interfering OpenCode/GSD observation path exists for the required demo proof.

Current product surfaces are insufficient:

- `doctor --profile github-actions-git-hooks-v1` diagnoses repo observer and CI
  artifact setup, not harness internals.
- `install repo-observer --profile github-actions-git-hooks-v1` can install
  observation files, but that mutates the demo repository outside the
  OpenCode/GSD loop.
- `wrap` and `run` capture command boundaries, not model calls, prompt
  boundaries, GSD phases, tool calls, file mutations, and PR cycles.
- `interaction relay` observes prompt delivery only by routing through
  `sdp-trace`, which changes the harness path.
- `assess preview --profile adapter-capture` can name the missing evidence
  shape, but no OpenCode/GSD adapter event source exists.

The correct product response is not to lower the demo bar. The correct response
is to define a portable harness observation intake that can ingest externally
emitted harness events and validate exactly what was and was not observed.

OpenCode/GSD is the first dogfood exemplar, not the product dependency and not
the canonical harness format. No `sdp-trace` command, schema, Go package, or
acceptance test may require OpenCode, GSD, MiniMax, or an OpenCode/GSD binary to
exist. The first fixture set must include a harness-generic complete export and
may include `opencode-gsd-jsonl-v1` as one tested profile exemplar.

## Non-Goals

- No native dependency on OpenCode, GSD, MiniMax, Codex, Claude, GitHub, Beads,
  Operator Mode, or any specific harness runtime.
- No prompt relay as the required observation path.
- No hidden prompt injection, harness prompt rewriting, or managed harness
  takeover.
- No automatic mutation of the assessed repository.
- No hand-authored trace artifacts as authority for a harness run.
- No raw prompt bodies, model responses, secrets, token-like values, private
  filesystem paths, authenticated URLs, or raw provider logs in committed
  examples or default summaries.
- No claim that missing harness data means no harness was used.
- No production trust, release readiness, merge approval, risk acceptance, health
  score, or buyer-facing trust verdict.
- No Node.js, npm, JavaScript, TypeScript, or `.mjs` active product path.

## Product Boundary

Block 31 adds an observation intake and validation path. It does not operate the
harness.

`sdp-trace` may say:

- "this harness run has imported lifecycle events from source X";
- "model route was observed through event field Y";
- "GSD phase boundaries were observed, but review output was absent";
- "file mutation references are present and bind to head commit Z";
- "prompt body was not retained; prompt digest and boundary metadata were
  retained";
- "tool-call coverage remains `not_assessed` because the harness source did not
  emit tool events";
- "this OpenCode/GSD tested profile is supported as an import shape, not as a
  product dependency."

`sdp-trace` must not say:

- "the harness complied with the prompt";
- "the model output was correct";
- "the PR is safe to merge";
- "the feature is delivered";
- "OpenCode/GSD is generally supported" from one tested profile;
- "missing events are green."

## First Product Surface

The first implementation should be a file-based intake and validator:

```text
sdp-trace harness observe \
  --profile <harness-profile.json> \
  --source <harness-events.jsonl> \
  --out <run-dir>

sdp-trace harness validate \
  --profile <harness-profile.json> \
  --run <run-dir> \
  --out <validation.json>

sdp-trace harness summarize \
  --validation <validation.json>
```

`observe` reads explicit local files only. It must not call OpenCode, GSD,
provider APIs, GitHub APIs, or hidden shell commands. External harness tooling can
export events into the agreed event shape before `sdp-trace` ingests them.

`observe` is strict for the initial product surface. If it detects an unsafe
absolute path, parent traversal, URL-like local ref, authenticated URL,
token-like value, forbidden raw prompt/model field, malformed JSONL line, or
source digest mismatch, it exits non-zero and does not write the observed run.
The error output must include a safe event id when available, the source line
number when available, and a closed reason code; it must not echo the unsafe
value.

`validate` checks event shape, required dimensions, source identity, digest
bindings, content state, unavailable fields, and cross-links to existing
adapter, delivery-trace, repo-observer, and PR-review artifacts.

`summarize` renders only validation facts. It must not imply harness compliance,
feature delivery, PR approval, merge approval, or production trust.

## Harness Observation Profile

A profile declares the expected event families and degradation semantics for a
specific harness export shape.

Minimum fields:

| Field | Meaning |
| --- | --- |
| `profile_id` | Stable safe id, for example `opencode-gsd-jsonl-v1`. |
| `harness_family` | Safe family label such as `opencode-gsd`. |
| `schema_version` | Profile schema version. |
| `event_schema_version` | Event schema version accepted by this profile. |
| `required_event_families` | Event families required for a valid observation. |
| `optional_event_families` | Event families that may remain `not_assessed`. |
| `raw_retention_policy` | Whether raw fields are retained, redacted, or digest-only. |
| `unsupported_fields` | Fields known unavailable from this harness source. |
| `degradation_rules` | How missing, malformed, stale, or unsafe events map to states. |

Profile `unsupported_fields` declares fields the harness source does not emit
for any event in that profile. Event `unavailable_fields` records gaps specific
to one event instance. A field listed in profile `unsupported_fields` must not be
repeated as a per-event unavailable field unless the event supplies a narrower
reason for that specific instance.

When a harness export format changes, the profile must change `profile_id` or
`schema_version`, and `validate` must check the profile `event_schema_version`
against imported events. A mismatch produces top-level `cannot_verify` with
reason `schema_version_mismatch`.

OpenCode/GSD may be the first tested exemplar profile, but the contract must
stay generic enough for other harnesses to export equivalent lifecycle events.

### Degradation Rule Grammar

The first implementation uses a closed JSON grammar, not free-form policy prose:

```json
{
  "degradation_rules": {
    "missing_required_family": {
      "state": "not_assessed",
      "reason_code": "required_event_family_absent"
    },
    "missing_optional_family": {
      "state": "not_assessed",
      "reason_code": "optional_event_family_absent"
    },
    "source_unavailable": {
      "state": "cannot_verify",
      "reason_code": "source_unavailable"
    },
    "unsafe_input": {
      "state": "fail",
      "reason_code": "unsafe_input"
    },
    "digest_mismatch": {
      "state": "cannot_verify",
      "reason_code": "source_digest_mismatch"
    },
    "schema_version_mismatch": {
      "state": "cannot_verify",
      "reason_code": "schema_version_mismatch"
    },
    "cross_link_conflict": {
      "state": "cannot_verify",
      "reason_code": "adapter_harness_state_conflict"
    }
  }
}
```

The schema may add closed rule keys later, but the first implementation must not
accept arbitrary state names, arbitrary reason strings, or executable policy.

## Harness Event Contract

The intake event shape must remain portable and digest-first.

Minimum event fields:

| Field | Meaning |
| --- | --- |
| `event_id` | Stable event id within the source export. |
| `event_schema_version` | Event schema version. |
| `event_family` | Closed family: `harness`, `model`, `interaction`, `phase`, `review`, `tool`, `mutation`, `test`, `pr`, `merge`, or `gap`. |
| `event_type` | Profile-specific type constrained by the profile. |
| `observed_at` | Timestamp if available; otherwise unavailable with reason. |
| `source_ref` | Safe source ref to the harness export. |
| `source_digest` | Digest of the source event or retained source segment. |
| `task_ref` | Task/spec/phase reference when available. |
| `operation_ref` | Delivery trace or adapter operation reference when available. |
| `actor_ref` | Safe actor or harness component ref when available. |
| `content_state` | `redacted`, `digest_only`, `retained_safe`, or `not_applicable`. |
| `unavailable_fields` | Explicit unavailable fields with reason. |

Raw prompts and model responses are not required. If retained, they must be
redacted or stored outside committed examples with digest references only.

`content_state` describes data treatment only. Assessment states such as `pass`,
`fail`, `not_assessed`, and `cannot_verify` belong in validation output, not in
the ingested event.

`not_applicable` means the event type has no raw content field covered by the
retention policy.

`source_digest` is SHA-256 of the event JSON after parsing the JSONL line,
setting only `source_digest` itself to an empty string, and re-encoding the event
with Go `encoding/json` deterministic object-key ordering. This avoids a
self-referential digest while still binding every other event field. The
validator recomputes this digest for each source line and must match
`source_digest` before marking a dimension digest-bound.

`observed_at`, when present, must be RFC 3339 with timezone. If the source does
not expose a timestamp, the relevant field must be listed in
`unavailable_fields` with a closed reason code.

A safe ref is a stable, portable identifier constrained by the schema. It must
not contain paths, URLs, tokens, or personal identifiers. Valid safe refs include
profile ids, event ids, SHA-256 source digests, task refs, operation refs, and
actor refs when they satisfy the safe string grammar.

`gap` is an explicit structural gap event emitted by the harness export itself.
If the harness does not emit gap events, `validate` may synthesize validation
dimensions from `degradation_rules`, but `observe` must not hand-author events to
make the export look more complete.

When `operation_ref` points to existing adapter or delivery-trace artifacts, the
ref must use a deterministic closed scheme such as
`adapter-run:<run_id>/event:<sequence>` or
`delivery-trace:<trace_id>/operation:<operation_id>`. If adapter and harness
sources cover the same required dimension and disagree, validation records
`cannot_verify` with reason `adapter_harness_state_conflict`.

## Validation State Mapping

Top-level validation result:

- `pass`: all required profile dimensions are present, safe, and digest-bound.
- `fail`: an unsafe event, invalid source, broken digest, forbidden raw field, or
  contradictory required state is observed.
- `cannot_verify`: required dimensions exist but cannot be linked to source,
  commit, task, or retained export identity.
- `not_assessed`: the profile cannot assess a dimension because the supplied
  harness export does not provide it.

Every dimension must also carry its own state and reason. The top-level state
must not flatten `not_assessed` or `cannot_verify` into `pass`.

Positive composition rule:

1. Required event-family dimensions determine the top-level state.
2. Optional event-family dimensions are reported at dimension level and do not
   promote the top-level state unless they contain `fail` for unsafe input.
3. Required dimensions compose by strictest state:
   `fail` > `cannot_verify` > `not_assessed` > `pass`.
4. No supplied run or unreadable source maps to top-level `cannot_verify` with a
   concrete source reason. A supplied run with zero events maps each missing
   required family to `not_assessed` with reason
   `required_event_family_absent`.
5. Unsafe input anywhere maps top-level state to `fail`.

## OpenCode/GSD Tested Profile Acceptance

The first tested profile is acceptable only if a reviewer can reconstruct:

1. which OpenCode version and model route were declared or observed;
2. which GSD phases or reviews occurred;
3. which prompt or interaction boundaries exist as digests or safe refs;
4. which tool/command/file mutation events are available;
5. which tests or validations were run;
6. which PR and merge states are observed or explicitly unavailable;
7. which parts remain `not_assessed` or `cannot_verify`.

If OpenCode/GSD cannot export some required fields without prompt relay or
manual trace writing, the profile must report those dimensions honestly instead
of lowering the requirement.

## Safety Requirements

- Event ingestion must reject absolute paths, parent traversal, symlink escapes
  after `filepath.EvalSymlinks`, URL-like local refs, authenticated URLs,
  token-like values, and unsafe personal identifiers in fields intended for
  committed examples or summaries.
- Summaries must not print raw prompt bodies, raw model responses, raw command
  bodies, provider tokens, authenticated URLs, private paths, or synthetic secret
  markers.
- Digest fields must use deterministic lowercase SHA-256 hex.
- Profile ids, family ids, event families, event types, actor refs, and source
  refs must be constrained to safe portable strings.
- Missing external evidence keeps the demo block open.

Token-like value detection must be explicit and testable. The first
implementation must reject known provider-token prefixes used in fixtures,
bearer-token shapes, auth-bearing URL query keys such as `token`, `access_token`,
`api_key`, and high-entropy base64/base64url-like strings of at least 32 bytes.
It must allow lowercase SHA-256 hex only in fields declared as digest fields.

Authenticated URL detection must reject userinfo before the host and known
credential query parameters. The error must identify the field and reason code
without echoing the value.

## Summary Output

`summarize` defaults to safe human text for PR and terminal use. A later JSON
summary mode may be added only if it has a schema and the same output-safety
tests.

The text summary must include:

- top-level validation state and reason;
- per-dimension state, reason, and event count;
- required dimensions that are `not_assessed` or `cannot_verify`;
- unsafe-input failures by reason code without unsafe values;
- profile id, event schema version, and validation artifact digest;
- an explicit non-authority statement: no harness compliance, feature delivery,
  PR approval, merge approval, release readiness, or production trust is claimed.

The text summary must not include raw prompts, raw model responses, raw command
bodies, provider tokens, authenticated URLs, private paths, or unsafe personal
identifiers.

## Acceptance Criteria

- A schema exists for harness observation profiles, harness events, observed
  harness runs, and harness validation output.
- Focused fixtures cover a harness-generic complete export, an OpenCode/GSD
  exemplar export, zero-event source, missing model route, missing phase/review
  events, prompt digest only, tool-event gap, unsafe raw prompt, unsafe source
  ref, symlink escape, source digest mismatch, schema version mismatch, file
  mutation without source binding, absent PR state, and no run supplied.
- Go tests validate every fixture against the implemented contracts. JSON syntax
  checks with `jq empty` are necessary but not sufficient.
- `go test ./...`, schema validation, and `git diff --check` pass for the changed
  scope.
- `docs/agent-entrypoint.md` and `docs/reviewer-entrypoint.md` are updated only
  after the command surface exists.
- `docs/reviews/demo-jvm-gsd-observation-ledger.md` references this block and
  keeps P0-001 open until the implemented validator can assess a real
  OpenCode/GSD export.

## Socratic Review Questions

1. Does the file-based intake avoid interfering with OpenCode/GSD while still
   producing enough evidence to assess harness participation?
2. Is `harness observe` too close to an adapter-capture duplicate, or does it
   correctly sit above adapter events as a harness-export intake?
3. Does the first OpenCode/GSD profile accidentally create a product dependency
   on OpenCode/GSD?
4. Are raw prompt and model response boundaries strict enough for committed
   examples and PR summaries?
5. Can a demo operator produce the required source export without hand-authoring
   proof artifacts?
6. Are `not_assessed` and `cannot_verify` preserved at dimension level and
   top-level validation?
