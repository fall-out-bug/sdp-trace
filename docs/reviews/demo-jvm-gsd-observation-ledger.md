# Demo JVM GSD Observation Ledger

This ledger records findings from the attempted `sdp-trace-demo-jvm-gsd`
product proof. The demo stop condition is any P0 finding against `sdp-trace`.

## Scope

- Demo repository: `fall-out-bug/sdp-trace-demo-jvm-gsd`
- Local path: `/Users/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd`
- Harness loop: OpenCode + GSD
- Required model route: `minimax-coding-plan/MiniMax-M2.5`
- Product-under-test: `sdp-trace` as a non-interfering observer of harness work
- Required boundary: `sdp-trace` must not alter harness prompts, manually edit
  demo repository trace artifacts, or require post-hoc hand-authored evidence
  to prove the run.
- Acceptable setup boundary: before the delivery work starts, the customer may
  run one initialization command, select one observation profile/configuration
  file, and install one reviewed wrapper or hook if the profile requires it.
  After that setup, the path must be "set up and forget": no extra prompt steps,
  hand-export steps, adapter authoring, or harness-loop friction during
  delivery.

## P0

### P0-001: No usable first-run OpenCode/GSD observation path exists for the customer case

- Status: open
- Severity: P0
- Demo state: stopped before project initialization
- Observed on: 2026-05-09
- Disposition: unresolved product blocker. Block 31 added a generic harness
  import surface, but it does not yet give a customer-usable first-run path for
  the OpenCode/GSD case that the demo is meant to prove.
- Tracking:
  `specs/001-sdp-trace-time-series-evidence-substrate/blocks/31-non-interfering-harness-observation.md`

#### Expected product behavior

`sdp-trace` should give a customer-usable first-run path for an OpenCode + GSD +
MiniMax-driven delivery cycle. A customer may do a bounded setup step before the
work starts, such as repository initialization, wrapper installation, or profile
configuration. After that, the observation path must be "set up and forget": the
customer should be able to run the real harness workflow without building an
adapter, rewriting the harness prompt stream, hand-exporting trace material, or
hand-authoring proof artifacts in the demo repository during delivery.

For each feature cycle, the product proof needs enough evidence to reconstruct
at least:

- selected model route;
- harness identity and GSD phase/review activity;
- prompts or interaction boundaries without prompt injection into the harness;
- file mutations;
- test and validation evidence;
- PR and merge state;
- explicit `pass`, `fail`, `not_assessed`, or `cannot_verify` states for gaps.

The demo must prove that this path works for the customer case, not merely that
`sdp-trace` can validate a pre-shaped JSONL file after someone else has solved
the observation problem.

#### Observed behavior

The available `sdp-trace` product surfaces are insufficient for this customer
path:

- `doctor --profile github-actions-git-hooks-v1` can diagnose that repo observer
  files and CI artifact surfaces are absent, but it does not observe OpenCode or
  GSD internals.
- `install repo-observer --profile github-actions-git-hooks-v1` can propose
  repository hook and CI observer files, but using it would mutate the demo
  repository outside the OpenCode/GSD harness loop.
- `wrap` / `run` can record a command boundary, stdout/stderr digests, and local
  source state, but that is command wrapping, not a non-interfering OpenCode/GSD
  adapter with model call, prompt, review, tool-call, file-mutation, and PR-cycle
  coverage.
- `interaction relay` would route prompt delivery through `sdp-trace`, which
  violates the demo boundary that `sdp-trace` must not interfere with harness
  prompts.
- `assess preview --profile adapter-capture` reports `run: absent` and only
  names the expected evidence shape; no OpenCode/GSD adapter event source is
  installed or available.
- `harness observe` can ingest an explicit `harness-event-v1` JSONL export, but
  it does not give the customer a supported way to run next to OpenCode/GSD and
  produce that export from the real workflow.

#### Reproduction evidence

From the empty demo repository:

```text
$ git status --short --branch
## main...origin/main

$ git log --oneline --decorate --max-count=1
7bfab75 (HEAD -> main, origin/main, origin/HEAD) Reset demo repository to empty state

$ opencode --version
1.14.41

$ opencode models | rg -i "MiniMax-M2.5|minimax"
minimax-coding-plan/MiniMax-M2.5
...

$ /tmp/sdp-trace-demo-observer doctor --profile github-actions-git-hooks-v1
Install state: fail
Proof state: not_assessed
...
github_actions_artifact_bundle | not_assessed | not_assessed | ci_uploaded | filesystem:.sdp-trace/ci | run CI and inspect uploaded artifact bundle
local_wrapped_commands | not_assessed | not_assessed | not_applicable | sdp_trace_runs:not_inspected | outside selected profile; no action required
agent_prompt | not_assessed | not_assessed | agent_reported | agent_prompt:not_inspected | do not rely on prompt instructions as setup proof

$ /tmp/sdp-trace-demo-observer assess preview --profile adapter-capture
{
  "command": "assess preview",
  "selected_profile": "adapter_capture",
  "inputs": {
    "run": "absent"
  },
  "claim": "preview is read-only and does not emit an adapter capture verdict"
}
```

#### Impact

The agreed demo cannot honestly start. Starting it would either:

- rely on manual or post-hoc trace artifact creation;
- use command wrapping or prompt relay as a substitute for harness observation;
- mutate the demo repository with observer setup outside the OpenCode/GSD loop;
- require the customer or demo agent to build the missing observation adapter
  before product value is available;
- or reduce the proof to CI/repo facts while missing the actual AI harness work.

All paths would overclaim `sdp-trace` product capability. A demo that hides this
behind prepared artifacts would prove demo choreography, not product
workability.

#### Required product change

Define and implement a customer-usable first-run observation path for the
OpenCode/GSD case. It may require bounded pre-work setup, but it must let a
customer run the real OpenCode/GSD workflow after setup and produce a
verifier-backed observation package for model calls, prompt/interaction
boundaries, GSD phase/review activity, tool calls, file mutations, test
observations, PR state, and merge state without requiring prompt changes, a
customer-built adapter, hand-export chores, or hand-authored demo trace
artifacts during delivery. Unsupported or unavailable fields must remain
explicit `not_assessed` or `cannot_verify` states, but buyer-critical
delivery-loop evidence must be observed rather than replaced by disclaimers.

Allowed observation mechanisms are process-boundary capture, stdout/stderr
digests or retained-safe excerpts, declared log tailing, declared output
directory watching, and filesystem artifact reads. `sdp-trace` must not inject
stdin, rewrite harness arguments after setup, mutate undeclared environment,
hide PATH rewrites, or interpose on provider network calls.

#### Current status

The finding remains an unresolved product blocker. Block 31 merged only an
initial generic harness observation command surface:

- `sdp-trace harness observe --profile <harness-profile.json> --source <harness-events.jsonl> --out <run-dir>`
- `sdp-trace harness validate --profile <harness-profile.json> --run <run-dir> --out <validation.json>`
- `sdp-trace harness summarize --validation <validation.json>`

That surface is not enough for the customer case because it assumes a correctly
shaped harness export already exists. P0-001 can move out of `open` only after
`sdp-trace` has a supported first-run path for the actual OpenCode/GSD workflow,
for example:

```text
sdp-trace observe session --profile opencode-gsd --out <run-dir> -- <harness-command>
```

The exact CLI may change, but the closure contract must not. The path must:

- allow bounded pre-work setup, then run beside the real harness command without
  prompt relay or additional in-loop operator steps;
- produce the observation evidence needed by `harness observe` / `validate` /
  `summarize` without a customer-built adapter;
- bind setup metadata, command digest, process id or unavailable reason,
  start/end time bounds, source commit, and output artifact digests to the run;
- validate the real customer-case workflow, not only a pre-shaped fixture;
- preserve every unsupported field as `not_assessed` or `cannot_verify` without
  weakening the core claim that the delivery loop was observed.

#### Scope clarification on 2026-05-09 after Block 31 merge

Status remains `open`.

The separate demo agent is blocked because this first-run product path does not
exist. The Block 31 merge in `sdp-trace` does not by itself fix the reported
problem; it only creates a generic intake surface. Treating that surface as
sufficient would shift product work onto the customer.

Until `sdp-trace` can run against the actual customer-case workflow and produce
verifier-backed observation evidence, any claim that P0-001 is fixed would be an
overclaim. Current closure state: `unresolved_blocker`.

#### Recheck on 2026-05-09 after PR #29 / first-run observation merge

Status remains `open`.

Verified:

- `sdp-trace` `main` includes `64c3583 Merge pull request #29 from
  fall-out-bug/codex/block-31-first-run-observation`.
- `sdp-trace --help` exposes `observe setup`, `observe collect`, and
  `observe session`.
- Example profiles exist:
  `examples/harness-observation/opencode-gsd-session-profile.example.json` and
  `examples/harness-observation/opencode-gsd-harness-profile.example.json`.
- A live `observe session` run can execute OpenCode with
  `minimax-coding-plan/MiniMax-M2.5` without mutating the demo repository.

Closure failed because the first-run path records only command/process
provenance for native OpenCode execution. It still does not produce or collect
OpenCode/GSD harness events for the real workflow:

```text
$ sdp-trace observe session --profile session-profile.json --out run -- \
    opencode run --format json --model minimax-coding-plan/MiniMax-M2.5 \
    --dir /Users/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd \
    "Respond with OK only."
...
"command_digest_state": "pass",
"process_id_state": "pass",
"source_commit_state": "cannot_verify",
"collection_state": "cannot_verify",
"collection_reason": "source_unavailable"
```

The command proves that `sdp-trace` can wrap a first-run process and avoid raw
stdout/stderr retention. It does not prove the required OpenCode/GSD observation
path because the declared `event_source_path` was not produced by the product or
the harness. Continuing the five-feature demo would still require either a
pre-shaped `harness-event-v1` file, a customer-built adapter, or hand-authored
events, which is the original blocker.

#### Block 31 implementation response on 2026-05-09

Status remains `open` pending demo-repo recheck against the real customer-case
workflow, but the accepted product gap is now fixed in `sdp-trace`: a session
profile can declare a raw OpenCode JSONL source via `raw_event_source_path` and
`raw_event_format: opencode-jsonl-v1`; `observe collect` and `observe session`
normalize that raw source into digest-only `harness-event-v1` records before
validation.

Review dispositions:

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| B31-RETURN-01 | critical | requirements | First-run OpenCode/GSD path still required a pre-shaped `harness-event-v1`, customer adapter, or hand-authored events after setup. | accepted_fixed | `SessionProfile.raw_event_source_path`, `raw_event_format`, `CollectSession` raw normalization, and `TestObserveSessionNormalizesOpenCodeRawJSONL`; demo recheck still required before closing this ledger item. |
| B31-REVIEW-01 | critical | tracing/evidence | Normalizer fabricated a `harness_observed` fallback for unrecognized raw lines. | accepted_fixed | Removed fallback event synthesis; `TestObserveCollectDoesNotFabricateEventsForUnrecognizedRawJSONL` asserts unrecognized raw input produces no events and validation remains `not_assessed`. |
| B31-REVIEW-02 | major | tracing/evidence | Raw message content could promote free text such as "tests pass" or "used a tool" into event-family evidence. | accepted_fixed | Family classification now uses structured signal keys and selected event-type metadata, not arbitrary content; `TestObserveSessionDoesNotPromoteMessageTextToEvidence` covers the overclaim case. |
| B31-REVIEW-03 | major | code/safety | Raw sensitive fields such as `api_key` were not rejected by field name alone. | accepted_fixed | Added sensitive field-name rejection; `TestObserveCollectRejectsUnsafeRawOpenCodeJSONL` asserts raw normalization fails before writing normalized output. |
| B31-REVIEW-04 | minor | code/safety | A profile could point raw input and normalized output at the same file. | accepted_fixed | `normalizeRawEvents` rejects equal raw and normalized paths before reading/writing. |
| B31-REVIEW-05 | minor | requirements | `file.*` mutation detection was too broad and could count read-only file events as mutation evidence. | accepted_fixed | Mutation detection now accepts explicit mutation event types only; `TestObserveCollectDoesNotTreatFileReadAsMutation` keeps `file.read` at `not_assessed`. |
| B31-REVIEW-06 | minor | tracing/evidence | Negative content-promotion coverage did not include the `phase` family. | accepted_fixed | `TestObserveSessionDoesNotPromoteMessageTextToEvidence` now requires `phase` and asserts message text does not satisfy it. |

Local verification:

```text
go test ./...
jq empty schema/*.json examples/harness-observation/*.json
git diff --check
```

#### Recheck on 2026-05-09 after PR #30 / OpenCode normalizer merge

Status remains `open`.

Verified:

- `sdp-trace` `main` includes `049658b Merge pull request #30 from
  fall-out-bug/codex/block-31-opencode-normalizer`.
- `SessionProfile` supports `raw_event_source_path` and
  `raw_event_format: opencode-jsonl-v1`.
- A live `observe session` run can execute OpenCode with
  `minimax-coding-plan/MiniMax-M2.5`, redirect native `opencode run --format
  json` output to the declared raw event source, normalize that source, and
  keep the demo repository clean.

Closure failed because the actual OpenCode 1.14.41 JSONL stream was normalized
to zero `harness-event-v1` events:

```text
$ sdp-trace observe session --profile session-profile-safe.json --out run -- \
    sh -c 'opencode run --format json --model minimax-coding-plan/MiniMax-M2.5 \
    --dir /Users/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd \
    "Respond with OK only." > opencode_raw.jsonl'
...
"collection_state": "pass",
"collection_reason": "source_collected",
"event_count": 0

$ sdp-trace harness validate --profile opencode-gsd-harness-profile.example.json \
    --run run/observed --out validation.json
"validation_state": "not_assessed",
"reason_code": "required_event_family_absent"
```

Observed raw OpenCode event types were:

```text
{"type":"step_start", ...}
{"type":"text", ...}
{"type":"step_finish", ...}
```

The current normalizer recognizes fixture-style signals such as
`session.started`, `message`, `phase`, `tool.call`, `file.write`, and
`test.finished`, but it did not map the native OpenCode 1.14.41 stream observed
in the customer-case run. The product now has a raw-source mechanism, but it
still does not observe real OpenCode/GSD delivery-loop evidence for the demo.

#### Block 31 implementation response on 2026-05-09 for native OpenCode JSONL

Status remains `open` pending demo-repo recheck against the real delivery-loop
workflow, but the native OpenCode stream regression is accepted and fixed in
`sdp-trace`.

Review dispositions:

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| B31-RETURN-02 | critical | requirements | Actual OpenCode 1.14.41 `--format json` emits `step_start`, `text`, and `step_finish`; PR #30 normalized those native events to zero `harness-event-v1` records. | accepted_fixed | Native `step_start`/`step_finish` now map to `harness`; native `text` maps to `interaction`; `TestObserveSessionNormalizesNativeOpenCodeJSONL` covers the observed stream. |
| B31-RETURN-03 | major | DX/replayability | The example session profile used `opencode-gsd-events.normalized.jsonl`, which failed `safeOutFile` because the output stem contained a dot. | accepted_fixed | Example output path changed to `opencode-gsd-events-normalized.jsonl`; live replay now reaches collection instead of `unsafe output filename`. |
| B31-RETURN-04 | major | tracing/evidence | Native OpenCode `tool_use` can include private absolute paths in unretained `state.output`; rejecting those bodies prevents digest-only observation of real OpenCode tool events. | accepted_fixed | Raw normalization uses a narrower raw-event safety scan: sensitive field names and token-like structural values still fail, while unretained text/input/output bodies are not persisted and do not block classification; `TestObserveCollectNormalizesNativeOpenCodeToolUseWithPrivateOutput` covers the case. |

Live replay after the fix against OpenCode 1.14.41 and
`minimax-coding-plan/MiniMax-M2.5`:

```text
RAW TYPES
step_start
text
step_finish

SESSION
"collection_state": "pass",
"event_count": 3

VALIDATION
"validation_state": "not_assessed"
harness: pass, event_count=2
interaction: pass, event_count=1
model/phase/tool/mutation/test: not_assessed
```

This proves native OpenCode JSONL is no longer normalized to zero events. It
does not close P0-001 yet: the simple OK run does not contain structured model,
phase, tool, mutation, or test evidence, so the full
`opencode-gsd-harness-profile.example.json` validation correctly remains
`not_assessed` until the real demo delivery loop produces those signals or the
profile is narrowed with an explicit evidence rationale.

#### Block 31 implementation response on 2026-05-10 for session facts

Status remains `open` pending real delivery-loop evidence for `phase`, `test`,
and mutation-producing work, but the safe session-command fact gap is accepted
and fixed in `sdp-trace`.

Review dispositions:

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| B31-RETURN-05 | critical | requirements | The selected OpenCode model route was visible in the controlled command but lost as verifier-backed harness evidence because only the command digest was retained. | accepted_fixed | `SessionRun.command_model` records a safe model id extracted from `--model`; raw command bodies remain unretained; normalized output emits a digest-only `model` event with `source_ref: session-command`. |
| B31-RETURN-06 | major | tracing/evidence | Native `tool_use` should satisfy `tool`, but read-only tools such as `glob` or `bash` must not satisfy `mutation`. | accepted_fixed | Native `tool_use` maps to `tool`; `mutation` is emitted only for mutation tool names such as `edit`, `write`, `patch`, `apply_patch`, `update`, or `delete`; tests cover both read-only and edit cases. |
| B31-REVIEW-07 | critical | code/correctness | Reviewer suspected `extractCommandModel(opts.Command)` was a compile-time type mismatch. | false_positive | `SessionOptions.Command` is `[]string`; `go test ./...` passes. |
| B31-REVIEW-08 | major | code/correctness | Positive edit-tool test did not explicitly assert the `mutation` family state. | accepted_fixed | `TestObserveCollectTreatsNativeEditToolAsMutation` now asserts both `tool` and `mutation` pass. |

Live replay after the fix against OpenCode 1.14.41 and
`minimax-coding-plan/MiniMax-M2.5`:

```text
SESSION
"command_model": "minimax-coding-plan/MiniMax-M2.5"
"command_model_state": "pass"
"collection_state": "pass"
"event_count": 4

VALIDATION
"validation_state": "not_assessed"
harness: pass, event_count=2
interaction: pass, event_count=1
model: pass, event_count=1
phase/tool/mutation/test: not_assessed
```

This improves the customer-case proof surface without closing P0-001. The
product can now prove the selected model route and native OpenCode base stream
without asking the customer to hand-author events. It still cannot honestly
claim GSD phase, tool, mutation, or test coverage unless the real delivery
workflow emits corresponding native `tool_use`/mutation/test/phase facts.

## P1

No P1 findings recorded yet.

## P2

No P2 findings recorded yet.

## P3

No P3 findings recorded yet.
