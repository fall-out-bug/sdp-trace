# Block 30: Automated PR Review Evidence Mechanism

Status: Reviewed spec. Implementation has started after explicit approval of
the Socratic-reviewed direction.

Parent artifacts:

- `docs/reviewer-entrypoint.md`
- `docs/agent-entrypoint.md`
- `docs/evidence-policy.md`
- `schema/evidence-event.schema.json`
- `schema/provenance-record.schema.json`
- `schema/review-ledger.schema.json`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/19-adapter-event-contract-capture-depth.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/29-interaction-trace-and-friction-metrics.md`

## Goal

Make automated PR review observable, replayable, and honestly bounded.

A repository operator should be able to run a multi-role PR review panel over a
frozen PR packet, retain the review provenance and findings, and inspect which
planes were assessed, which model or harness produced each finding, which
outputs were unusable, and which closure decisions remain external to
`sdp-trace`.

This mechanism contractually prevents AI merge approval. It produces an
immutable review evidence record. It does not output a merge readiness score,
approval signal, policy pass, or release authorization. Any external system
that turns this record into a gate acts outside `sdp-trace` authority.

The mechanism is:

```text
For this PR, these independent review planes assessed this exact frozen change
packet, produced these findings, left these areas not_assessed or
cannot_verify, and recorded this disposition ledger.
```

## Product Question

"What did independent automated reviewers actually assess for this PR, what did
they find, and what evidence lets a later human or policy consumer trust the
review record without trusting the model summary?"

The answer must include:

- immutable source and diff identity;
- review packet digest;
- role and review-plane coverage;
- model, runner, command, prompt, and context provenance;
- reviewer output status, including hung, empty, off-task, and parse failures;
- finding severity, citations, and disposition;
- local verification and CI/check state as separate evidence;
- explicit `not_assessed` or `cannot_verify` gaps;
- raw-output retention and redaction state;
- a clear statement that merge, release, override, and risk acceptance remain
  outside `sdp-trace`.

## Problem

The repository already requires multi-plane review for trust-sensitive work, but
the current execution is operator discipline plus local scratch files. That is
not enough for a repeatable PR mechanism:

- reviewers can receive different context packets without anyone noticing;
- a hung or empty reviewer can be accidentally counted as evidence;
- a model can say "LGTM" without inspected file or diff citations;
- raw review output can contain private source, prompts, provider metadata, or
  unsafe URLs;
- CI absence can be collapsed into green by prose;
- a final synthesis can flatten conflicting findings into one optimistic
  summary;
- PR comments can imply approval even when review coverage is partial.

For `sdp-trace`, this is a product problem. If automated PR review becomes a
trust surface, the review itself needs provenance, evidence, trace, redaction,
and state semantics.

## Mechanism Boundary

`sdp-trace` may provide a Go-first mechanism to:

- build a frozen PR review packet from explicit local/provider inputs;
- run configured external review commands such as `pi` and `opencode` when the
  operator explicitly selects them;
- enforce read-only review mode for repo-aware reviewers;
- record external runner availability, model identity, command line, timeout,
  prompt digest, context digest, packet digest, output digest, parse result, and
  raw-output retention state;
- validate structured reviewer output against a schema;
- synthesize a review ledger without hiding conflicts or unassessed planes;
- render a PR-safe summary that names assessed planes, review-record findings,
  residual gaps, and non-authority boundaries.

`sdp-trace` must not:

- approve, merge, close, label, or mark a PR ready;
- decide risk acceptance, override approval, release readiness, customer
  trust, production trust, or human approval;
- treat model confidence, absence of findings, or a clean synthesized summary as
  policy pass;
- require OpenCode, pi, provider SDKs, GitHub, or any hosted model as product
  dependencies;
- silently fetch live PR state without recording source, identity, and digest;
- commit raw model output, raw prompts, raw customer code, raw provider logs,
  authenticated URLs, private filesystem paths, tokens, or secrets;
- run repo-aware model agents with write/edit permissions;
- use Node.js, npm, JavaScript, TypeScript, or `.mjs` tooling in the active
  product path.

## First Product Surface

The first mechanism uses explicit inputs and optional external runners.

Required command surfaces after implementation:

```text
sdp-trace pr-review packet \
  --out <dir> \
  --repo-id <safe-id> \
  --change-ref <provider-neutral-ref> \
  --base <sha> \
  --head <sha> \
  --diff <file> \
  [--metadata <file>] \
  [--context <file>]... \
  [--verification <file>]...

sdp-trace pr-review run \
  --packet <dir> \
  --profile <review-profile.json> \
  --out <dir> \
  [--preview] \
  [--allow-external-runner pi] \
  [--allow-external-runner opencode]

sdp-trace pr-review synthesize \
  --packet <dir> \
  --runs <dir> \
  --out <file>

sdp-trace pr-review validate \
  --packet <dir> \
  --runs <dir> \
  --ledger <file> \
  --out <file>

sdp-trace pr-review summarize \
  --validation <file> \
  --ledger <file>

sdp-trace pr-review check \
  --out <dir> \
  --repo-id <safe-id> \
  --change-ref <provider-neutral-ref> \
  --base <sha> \
  --head <sha> \
  --diff <file> \
  --profile <review-profile.json> \
  [--ci-state <state>] \
  [--context <file>]... \
  [--verification <file>]... \
  [--allow-external-runner <runner>]...
```

The command names are intentionally PR-specific because the packet must bind to
source and change identity. The implementation may add a generic `review`
namespace later only if the contract is generalized without weakening PR
identity binding.

`packet` must not call GitHub or any provider API in the first implementation.
Provider adapters can prepare `--metadata`, `--diff`, and `--verification`
inputs externally. This prevents hidden GitHub coupling and keeps the first
contract replayable from files.

`run` may call external commands only when the operator names an allowed runner.
If the profile references a runner that is not named by
`--allow-external-runner`, the command fails fast with a usage error before any
external runner is invoked. It does not silently skip the role.

`run --preview` emits the planned roles, runner names, requested model ids,
timeouts, prompt-template refs, prompt digests, output directories, and command
digests without invoking external commands or contacting hosted models.

If a selected runner is unavailable, times out, returns empty output, produces
off-task structured output, or fails structured parsing, that run is recorded
with the deterministic state mapping below. It is not counted as usable review
evidence unless the mapping says it is usable.

`synthesize` creates the initial disposition ledger from reviewer results. Every
critical or major finding starts as `unresolved_review_blocker` unless an input
ledger already records a valid disposition with evidence. The operator or an
external tool may then update dispositions and rerun `validate`.

`validate` verifies the packet, run records, raw-output digest bindings,
required review-plane coverage, finding citation shape, ledger dispositions,
and safety metadata. It emits facts only.

`summarize` renders the validation and ledger into a PR-safe human summary. It
must not add claims absent from the JSON validation artifact and must not imply
merge readiness. The summary must render an authority line stating that merge,
release, risk acceptance, and approval are external decisions.

`check` is the common-path wrapper. It chains `packet`, `run`, `synthesize`,
`validate`, and `summarize` using the same contracts. The subcommands remain the
debuggable surfaces for replay and failure isolation.

Every command that writes `--out <dir>` fails if the directory exists and is
non-empty unless an explicit `--force` flag is added during implementation.
Every command that writes `--out <file>` refuses to overwrite an existing file
unless `--force` is supplied. Reusing a directory across packet digests is
`cannot_verify`.

## Packet Input Guide

The first implementation accepts file inputs only. Provider adapters may prepare
those files, but `sdp-trace` does not call provider APIs in this block.

Required input file formats:

| Input | Required format | Notes |
| --- | --- | --- |
| `--diff <file>` | Unified diff text from `git diff` or provider-equivalent unified diff. | The packet records `content_type=unified_diff` and SHA-256 digest. Provider JSON diffs are out of scope. |
| `--metadata <file>` | JSON object with provider-neutral optional fields: `provider`, `change_url_ref`, `title`, `author_ref`, `created_at`, `updated_at`. | Unsafe URLs or identities must be redacted or digest-only. |
| `--context <file>` | Any text or JSON artifact selected by the operator: spec, plan, tasks, docs, schema, or source excerpt. | Each context ref receives id, kind, safe ref, digest, and redaction state. |
| `--verification <file>` | Text or JSON output from local commands, CI/check summaries, or prior verifier results. | Verification inputs are evidence refs. They do not define policy state by themselves. |

`--ci-state` is an explicit packet flag with allowed values `pass`, `fail`,
`pending`, `not_assessed`, and `cannot_verify`. It defaults to `not_assessed`.
The packet command does not infer CI state from verification files.

Safe identifier rules:

- `repo_id` must match `^[a-z0-9][a-z0-9._-]{0,62}[a-z0-9]$`; it must not be
  derived from an absolute filesystem path.
- `change_ref` must match `^(pr|mr|change)-[A-Za-z0-9._-]{1,64}$`.

## Review Packet Contract

The packet is the immutable target each reviewer receives.

Minimum fields:

| Field | Meaning |
| --- | --- |
| `packet_id` | Stable id for this packet. |
| `schema_version` | Review packet schema version. |
| `repo_id` | Operator-supplied safe repository id, never derived from an absolute path. |
| `change_ref` | Provider-neutral PR/MR/change reference. |
| `base_commit` | Source base SHA. |
| `head_commit` | Source head SHA. |
| `diff_ref` | Safe ref to the diff file and SHA-256 digest. |
| `context_refs` | Spec, plan, task, docs, command contract, or code refs included in the packet. |
| `verification_refs` | Local command output, CI/check state, or `not_assessed` records. |
| `ci_state` | `pass`, `fail`, `pending`, `not_assessed`, or `cannot_verify`. |
| `created_at` | Packet creation timestamp. |
| `created_by` | Human or tool actor. |
| `redaction_status` | Packet redaction state. |
| `unavailable_fields` | Explicit missing or unassessed inputs. |

Every reviewer run must bind to `packet_digest`. If the packet changes, prior
review runs do not apply to the new head.

`diff_ref`, `context_refs[]`, `verification_refs[]`, `raw_output_ref`, and
`prompt_ref` use the safe-ref shape:

| Field | Meaning |
| --- | --- |
| `id` | Stable ref id within the packet or run. |
| `kind` | Closed enum: `diff`, `metadata`, `spec`, `plan`, `task`, `doc`, `schema`, `source_excerpt`, `verification`, `prompt`, `raw_output`, `sanitized_output`, or `external`. |
| `ref` | Repo-relative path, run-local path, or digest-only identifier. Absolute paths are forbidden in committed or PR-rendered artifacts. |
| `digest_sha256` | SHA-256 of the referenced content. Required. |
| `content_type` | Media/content type such as `unified_diff`, `markdown`, `json`, or `text`. |
| `redaction_state` | Closed enum: `none`, `redacted`, `digest_only`, `encrypted_ref`, `external_ref`, `withheld`, or `not_assessed`. |

Every finding citation must reference a packet `context_ref_id` or `diff_ref`
and include either a diff hunk id or source digest plus line range. Citations
that cannot be resolved against the frozen packet digest produce
`cannot_verify`.

## Review Profile Contract

The review profile declares which planes and runner roles are required.

Minimum required planes for trust-sensitive `sdp-trace` PRs:

- `code_correctness`
- `trace_evidence_provenance`
- `requirements_vs_implementation`

Recommended additional planes for high-risk PRs:

- `security_forgery_overclaim`
- `dx_replayability`
- `privacy_output_safety`

Each role must declare:

- `role_id`;
- `plane`;
- `runner`: `pi`, `opencode`, or `manual_external`;
- model or model-pattern requested, when applicable;
- whether repo tool access is allowed;
- timeout;
- required output schema;
- fallback behavior;
- whether raw output may be retained, sanitized, encrypted, or digest-only.

The profile is declarative. `sdp-trace` validates the profile shape; it does not
infer model choices from a user's local `pi` settings. This repository may ship
or document a `trust-sensitive-default` example profile that uses model-family
diversity and excludes OpenAI, Anthropic, and Google unless the operator
explicitly permits them for that review cycle.

Profile validation must be available without running reviewers. The
implementation may expose this as `pr-review validate-profile` or as
`pr-review run --preview` plus schema validation.

## External Runner Rules

### pi

`pi` is the first supported optional external runner for independent model
reviewers. It is not a mandatory product dependency. The mechanism remains
usable for packet validation, ledger validation, and summary rendering when no
model runner is available; reviewer execution then remains `not_assessed`.

When selected, the runner should default to:

```text
pi --no-tools --no-context-files --no-session -p <prompt>
```

The prompt must include the packet digest, review plane, required JSON output
shape, severity scale, citation rules, and "no merge approval" boundary.
Prompt construction must be deterministic. The first implementation must either
ship a prompt template with named placeholders such as `{{packet_digest}}` and
`{{review_plane}}`, or require an explicit `prompt_template_ref` in the review
profile. The final prompt digest is recorded before any runner executes.

If the selected model is unavailable, the run is `not_assessed` unless a
configured fallback model completes the same plane against the same packet.
Fallback records must include `requested_model`, `observed_model`,
`fallback_for_model`, and `fallback_reason`. The actual model identity is the
authority for provenance; fallback fields explain why it differs from the
profile intent.

### OpenCode

OpenCode is allowed only for bounded repo-aware read-only review roles.

The runner must:

- avoid `--dangerously-skip-permissions`;
- use an agent or permission profile that denies write, edit, delete, and
  external mutation operations before the model starts;
- record `opencode --version`;
- record selected `--model` and `--agent`;
- pass the frozen packet by file or prompt ref;
- record a safe working-tree baseline before and after the run when a working
  tree is inspected;
- fail or mark `cannot_verify` if the run mutates files or if mutation
  attribution cannot be separated from pre-existing dirty state;
- retain raw output only under ignored local run directories until sanitized.

If OpenCode cannot enforce read-only permissions for the selected role, that
role remains `not_assessed` and the agent is not executed. A before/after git
baseline is a verification step, not the primary safety control.

OpenCode working-tree modes:

- `clean_required`: the default. The run refuses to start if the working tree is
  dirty.
- `dirty_baseline`: allowed only if the profile explicitly permits it. The run
  records pre-run file/diff digests and marks the result `cannot_verify` if
  pre-existing dirty files intersect packet target files or if the post-run
  baseline cannot prove no runner-caused mutation.

Working-tree diagnostics in committed or PR-rendered artifacts must not list raw
filenames from `git status` or `git diff`. They may record safe counts, digests,
and reason codes such as `working_tree_dirty` or `mutation_detected`.

`opencode --version` failure is `not_assessed` with reason
`runner_unavailable`.

## Reviewer Output Contract

Every reviewer must return structured output. Free-form model prose can be
retained as raw evidence but cannot satisfy a required plane until parsed into
the structured contract.

Minimum reviewer result fields:

| Field | Meaning |
| --- | --- |
| `review_run_id` | Stable id for the model/run output. |
| `packet_digest` | SHA-256 digest of the frozen packet. |
| `plane` | Review plane. |
| `role_id` | Role declared by the review profile. |
| `runner` | `pi`, `opencode`, or external/manual. |
| `requested_model` | Model id from the profile, or `not_assessed`. |
| `observed_model` | Model id observed from runner config/output, or `not_assessed`. |
| `model_family` | Required string; observed family or `not_assessed`. |
| `model_version` | Required string; observed version or `not_assessed`. |
| `fallback_for_model` | Original model id when a fallback completed the plane. |
| `fallback_reason` | Safe reason code for fallback. |
| `status` | `findings_reported`, `no_findings`, `not_assessed`, `failed`, `timed_out`, `empty_output`, `off_task`, `parse_failed`, or `cannot_verify`. |
| `findings` | Structured findings with severity, citation, evidence refs, and exact fix or question. |
| `raw_output_ref` | Safe ref plus digest and retention state. |
| `prompt_ref` | Safe ref plus digest. |
| `context_refs` | Packet refs used by the reviewer. |
| `started_at` / `ended_at` | Timing evidence. |

`findings_reported` and `no_findings` are the only statuses that count as
usable review output. `no_findings` is not a policy pass. It only means the role
produced usable structured output with zero findings for its declared plane.

Failure mapping:

| Condition | Reviewer status | Counts as usable? | Required coverage effect for a required plane |
| --- | --- | --- | --- |
| Runner not selected for the profile | `not_assessed` | No | `not_assessed` if no attempt was requested |
| Runner referenced but not allowed by CLI | usage error before execution | No | no result artifact |
| Runner executable unavailable or version probe fails | `not_assessed` | No | `coverage_partial` or `not_assessed` depending on other planes |
| Timeout | `timed_out` | No | `cannot_verify` |
| Empty output | `empty_output` | No | `cannot_verify` |
| Output not parseable as required JSON | `parse_failed` | No | `cannot_verify` |
| Parsed output declares wrong packet, plane, or role | `off_task` | No | `cannot_verify` |
| OpenCode working tree mutates during execution | `cannot_verify` | No | `cannot_verify` |
| Parsed output matches schema and has findings | `findings_reported` | Yes | depends on disposition |
| Parsed output matches schema and has zero findings | `no_findings` | Yes | depends on other planes and findings |

Reviewer finding shape:

| Field | Meaning |
| --- | --- |
| `id` | Stable finding id within the reviewer run. |
| `severity` | Closed enum: `critical`, `major`, `minor`, or `informational`. |
| `citation` | Object containing `context_ref_id` and either `diff_hunk_id` or `source_digest` plus `line_start`/`line_end`. |
| `summary` | Safe one-line summary. |
| `rationale` | Why this matters under the review plane. |
| `suggested_fix` | Optional exact fix or required spec change. |
| `question` | Optional blocking question when no fix is yet known. |
| `evidence_refs` | Packet or run refs used by the reviewer. |

Requested-vs-observed model mismatch is a validation condition. If the observed
model differs from the profile and no fallback fields explain it, validation
emits `cannot_verify`.

## Ledger And Disposition

The synthesis ledger must preserve independent reviewer findings and conflicts.
It must not collapse them into a single optimistic verdict.

Each finding disposition must be one of:

- `accepted_fixed`
- `accepted_review_blocking`
- `accepted_narrower`
- `rejected_false_positive`
- `deferred_not_assessed`
- `unresolved_review_blocker`

The durable review ledger must record:

- source reviewer run;
- plane;
- model/runner provenance;
- finding severity;
- file/line, diff hunk, spec section, or artifact citation;
- disposition;
- evidence used to accept or reject the finding;
- re-review state after fixes;
- residual `not_assessed` or `cannot_verify` gaps.

A re-review after any diff change must create a new packet and a new
`packet_digest`. The ledger may link the original finding and the re-review
finding, but a re-review bound to a stale packet digest cannot close a finding
for the current head and produces `cannot_verify`.

## State Semantics

Review state and policy state remain separate.

Allowed review coverage states:

- `coverage_satisfied`: every required plane has usable structured output for
  the current packet digest, every critical or major finding has a disposition,
  and there are zero unresolved review blockers.
- `coverage_partial`: at least one required plane has usable output, but
  coverage or disposition is incomplete.
- `coverage_unresolved`: every required plane may have usable output, but at
  least one critical or major finding remains unresolved for review-record
  closure.
- `not_assessed`: no required plane was assessed, or the selected profile was
  out of scope.
- `cannot_verify`: required packet, runner, digest, parse, citation, or
  retention evidence cannot be verified.

`coverage_satisfied` is not merge approval. It is review-record completeness
under the selected profile.

`coverage_unresolved` takes precedence over `coverage_satisfied`. If every plane
was assessed but one critical or major finding remains unresolved, the state is
`coverage_unresolved`.

Every validation output must include:

```json
{
  "authority_scope": "review_record_only",
  "merge_decision": "not_authorized_by_sdp_trace",
  "release_decision": "not_authorized_by_sdp_trace",
  "risk_acceptance": "not_authorized_by_sdp_trace"
}
```

Summaries must render these fields structurally, not as optional prose. They
must not say "merge blocked", "safe to merge", "approved", "ready", or "policy
passed". They may say "review record has unresolved findings" or "review
coverage is not satisfied."

## Safety Requirements

Committed or PR-rendered artifacts must not contain:

- raw prompts unless explicitly approved and sanitized;
- raw model responses;
- raw customer code beyond the PR diff already in the packet;
- private absolute filesystem paths;
- authenticated URLs;
- provider request ids if they reveal tenant or token material;
- credentials, tokens, OIDC request tokens, API keys, cookies, or session ids;
- raw stdout/stderr bodies from arbitrary commands;
- unsafe personal identifiers.

Validation and summary output must include negative leak checks using synthetic
markers and must not echo the marker value in failure output.

Synthetic marker classes:

| Marker class | Injection point |
| --- | --- |
| `SYNTHETIC_PROMPT_SECRET_*` | prompt template/input fixture |
| `SYNTHETIC_TOKEN_SECRET_*` | runner environment/config fixture |
| `SYNTHETIC_PRIVATE_PATH_*` | raw command/path fixture |
| `SYNTHETIC_AUTH_URL_*` | metadata or raw output fixture |
| `SYNTHETIC_MODEL_RESPONSE_*` | raw reviewer output fixture |

Tests must inject markers into packet fields, prompt fields, raw output refs,
metadata refs, and runner result fixtures, then assert that `validate`,
`summarize`, and failure paths do not print any marker value.

Command-line provenance in durable artifacts must use safe refs and digests
instead of absolute local paths. Raw runner argv may be retained only as
local-only or digest-only evidence until sanitized.

## Acceptance Criteria

Spec approval requires:

- Socratic review across product boundary, UX/DX, trace/evidence, security, and
  implementation feasibility planes;
- every valid critical or major spec finding fixed or recorded as a blocker;
- an implementation plan that keeps external model runners optional and
  explicitly recorded;
- explicit user approval of the reviewed direction.

Implementation closure requires:

- Go-first schemas or Go structs for review packet, review profile, reviewer
  result, and validation output;
- focused tests for packet digest invalidation, missing planes, hung/empty
  reviewer states, parse failures, stale packet digest, CI `not_assessed`,
  read-only OpenCode mutation detection, and unsafe output redaction;
- a local fixture that simulates `pi` reviewer output without network calls;
- a local fixture that simulates OpenCode mutation or unavailable read-only
  enforcement and records `cannot_verify` or `not_assessed`;
- fake runners implemented as Go test helpers or Go-built test binaries, not
  shell scripts in the product path;
- PR-safe summary rendering that cannot imply merge approval;
- `go test ./...`, `jq empty schema/*.json`, changed-fixture validation, and
  `git diff --check`;
- implementation review across code/correctness, trace/evidence, and
  requirements-vs-implementation planes;
- PR-level review across the same planes before ready.

## Open Questions For Socratic Review

1. Should first implementation include `pr-review run`, or should it stop at
   `packet`, `validate`, and `summarize` while an external script runs models?
   Draft resolution: include `run`, because the accepted product goal is a
   mechanism, not just a validator. The implementation must still keep runners
   optional, explicit, and fake-runner-testable.
2. Is `pr-review` the right command namespace, or should this be an `assess`
   profile over review packages?
3. What minimum review profile should non-trust-sensitive repositories use?
4. How much raw reviewer output may be retained locally before redaction and
   deletion are required?
5. Should PR comment rendering be part of `sdp-trace`, or should it produce a
   markdown artifact for an external GitHub/GitLab adapter to post?
