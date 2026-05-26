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

#### Recheck on 2026-05-10 after PR #33 / session command facts merge

Status remains `open`.

Verified:

- `sdp-trace` `main` includes `4fa61ad Merge pull request #33 from
  fall-out-bug/codex/block-31-session-facts`.
- `go test ./...` passes.
- Live `observe session` can execute OpenCode with
  `minimax-coding-plan/MiniMax-M2.5`, collect native JSONL, normalize it, and
  keep the demo repository clean.

Closure failed because the live replay did not emit the claimed model event:

```text
SESSION
"command_digest_state": "pass"
"process_id_state": "pass"
"collection_state": "pass"
"event_count": 3

NORMALIZED EVENTS
harness
interaction
harness

VALIDATION
"validation_state": "not_assessed"
harness: pass, event_count=2
interaction: pass, event_count=1
model: not_assessed, event_count=0
phase/tool/mutation/test: not_assessed
```

This contradicts the expected recheck result recorded above for PR #33, where
`command_model` and a digest-only `model` event were expected from the command
argv. Either the live path is not using the new command-model fact during
normalization, or the example session profile does not activate that path. In
both cases the demo must remain stopped: the product still fails to preserve the
selected model route as verifier-backed harness evidence in this customer-case
run.

#### Implementation response on 2026-05-10

| id | severity | plane | finding | disposition | response |
| --- | --- | --- | --- | --- | --- |
| B31-RETURN-07 | critical | requirements | PR #33 extracted the model only from tokenized argv. The live demo command uses `sh -c`, so `--model minimax-coding-plan/MiniMax-M2.5` lived inside the shell command string and no verifier-backed `model` event was emitted. | accepted_fixed | Session observation now falls back to a shell-aware field scan for strict `sh -c` / `bash -c` wrappers after tokenized argv extraction fails. It still stores only a sanitized model id and command digest; raw command text remains unretained. |
| B31-RETURN-08 | major | tracing/evidence | Review found that a shell-extracted model id could retain whitespace or a line-continuation artifact, weakening the retained model evidence contract. | accepted_fixed | `safeCommandModel` now rejects spaces, tabs, CR, and LF in retained model ids; shell field splitting skips escaped newline continuations. Unit tests cover whitespace and newline rejection. |

Live replay after the fix used the same customer-case command shape:

```text
sdp-trace observe session --profile session-profile.json --out run -- \
  sh -c 'opencode run --format json --model minimax-coding-plan/MiniMax-M2.5 \
  --dir /Users/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd \
  "Respond with OK only." > opencode-gsd-run.jsonl'
```

Observed result:

```text
SESSION
"command_model": "minimax-coding-plan/MiniMax-M2.5"
"command_model_state": "pass"
"collection_state": "pass"

RUN
"event_count": 4
events/raw-000001-harness.json
events/raw-000002-interaction.json
events/raw-000003-harness.json
events/session-command-model.json

VALIDATION
"validation_state": "not_assessed"
harness: pass, event_count=2
interaction: pass, event_count=1
model: pass, event_count=1
phase/tool/mutation/test: not_assessed, event_count=0
```

This closes only the returned model-route finding. The five-feature demo P0
remains open until phase, tool, mutation, and test signals are also emitted as
verifier-backed evidence in the customer-case delivery loop.

Review disposition:

```text
code/correctness: no remaining critical or major findings after accepted fixes;
shell escape and sh -c positional-argv precedence risks accepted and fixed
requirements-vs-implementation: no critical or major findings; minor bash/doc
coverage notes accepted and fixed
tracing/evidence: one major finding accepted and fixed; minor zsh/login-shell
generalization deferred because the customer-case command is strict sh -c and
the current product contract intentionally avoids broad shell inference
focused evidence re-review: pass; no critical or major findings after sanitizer
hardening; reviewer note about --model=value was a false positive against the
full function because extractCommandModelArgs still handles --model= and -m=
```

### P0-002: Real OpenCode/GSD raw stream is rejected before observed run creation

- Status: open
- Severity: P0
- Demo state: stopped during first GSD project-initialization cycle
- Observed on: 2026-05-10
- Disposition: unresolved product blocker. The real OpenCode/GSD run produced
  native JSONL events with model metadata, tool calls, and file mutations, but
  `sdp-trace observe session` rejected the raw stream during collection.

#### Expected product behavior

After bounded setup, `sdp-trace observe session` should run beside the real
OpenCode/GSD command, normalize the native OpenCode JSONL stream, create an
observed run, and preserve unsupported fields as explicit `not_assessed` or
`cannot_verify` states. Safety handling may redact or hash sensitive values, but
it must not make the whole customer-case run unverifiable when the raw stream is
valid harness output.

#### Observed behavior

A real GSD initialization run succeeded under OpenCode with the required model
route:

```text
opencode run --format json --model minimax-coding-plan/MiniMax-M2.5 \
  --dir /Users/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd \
  "/gsd-new-project --auto" \
  --file /private/tmp/sdp-trace-demo-observe-gsd-new-project-input.md \
  --dangerously-skip-permissions
```

The run created GSD planning files in the demo repository:

```text
.planning/PROJECT.md
.planning/REQUIREMENTS.md
.planning/ROADMAP.md
.planning/STATE.md
```

The raw OpenCode stream contained useful evidence:

```text
line 3: tool_use read /private/tmp/sdp-trace-demo-observe-gsd-new-project-input.md
line 9: tool_use glob path /Users/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd
line 12: tool_use bash command "rtk ls -la"
line 15: task metadata model minimax-coding-plan/MiniMax-M2.5
line 18: read .planning/ROADMAP.md
```

However, `sdp-trace observe session` failed during raw collection:

```text
raw source line 9: unsafe_input:part.state.input.path:token_like_value
```

The session metadata retained only command-level facts:

```text
"command_model": "minimax-coding-plan/MiniMax-M2.5"
"command_model_state": "pass"
"process_id_state": "pass"
"collection_state": "cannot_verify"
"collection_reason": "not_collected"
```

No `run/observed` package was produced, so `harness validate` could not assess
the actual GSD activity.

#### Impact

The demo must stop. The product can now see that the selected model route was
requested, but it cannot ingest the real OpenCode/GSD event stream that proves
what happened during the delivery loop. Continuing would require manually
editing or filtering trace artifacts outside the product, which violates the
demo boundary and would be self-forged evidence.

This is buyer-critical: the first non-trivial GSD run already contains exactly
the kind of evidence `sdp-trace` must preserve, but the current safety filter
turns a valid harness field into a fatal collection error instead of retaining a
safe representation or marking only that field unavailable.

#### Required product change

`sdp-trace` must handle real OpenCode/GSD JSONL fields such as
`part.state.input.path` without aborting collection. Acceptable fixes include
schema-aware path classification, safe path hashing/redaction, or per-field
`cannot_verify` retention. The closure test must replay this real run shape and
prove that the observed run contains verifier-backed evidence for at least:

- harness and interaction boundaries;
- selected model route;
- tool calls, including `read`, `glob`, `bash`, and `task`;
- file mutations or generated artifact reads for `.planning/*`;
- explicit `not_assessed` or `cannot_verify` states for missing phase/test/PR
  evidence.

#### Implementation response on 2026-05-10

| id | severity | plane | finding | disposition | response |
| --- | --- | --- | --- | --- | --- |
| B31-RETURN-09 | critical | tracing/evidence | Real OpenCode/GSD `tool_use` events with path-like tool inputs, for example `part.state.input.path`, were rejected as `token_like_value` before an observed run could be created. | accepted_fixed | Raw OpenCode safety scanning now treats known path-like field names as path evidence under digest-only normalization, so private local paths are not retained and do not abort collection. Sensitive keys such as `api_key`, `token`, `authorization`, and authenticated URLs remain fatal unsafe input. |
| B31-RETURN-10 | critical | tracing/evidence | Review challenged whether the path-like exemption could hide secrets in `path` fields. | accepted_fixed | The exemption was narrowed to direct file/directory field names only, and regression tests prove provider-token values and authenticated URLs in `part.state.input.path` still fail before any normalized output is written. |

Regression coverage now replays the returned GSD shape with native OpenCode
`tool_use` events for `read`, `glob`, `bash`, `task` model metadata, and a
`.planning/ROADMAP.md` read. The observed run is created with digest-only
events, and normalized output is asserted not to retain `/private/tmp` or
`/Users/fall_out_bug` paths.

Observed result from the regression:

```text
SESSION
"command_model": "minimax-coding-plan/MiniMax-M2.5"
"command_model_state": "pass"
"collection_state": "pass"

RUN
"event_count": 8
events/raw-000001-harness.json
events/raw-000002-interaction.json
events/raw-000003-tool.json
events/raw-000004-tool.json
events/raw-000005-tool.json
events/raw-000006-model.json
events/raw-000007-tool.json
events/session-command-model.json

VALIDATION
"validation_state": "not_assessed"
harness: pass
interaction: pass
model: pass
tool: pass
phase/mutation/test: not_assessed
```

This closes the collection-abort blocker only. The demo P0 remains open for
phase, mutation, and test evidence until those signals are emitted as
verifier-backed events in the customer-case delivery loop.

Review disposition:

```text
code/correctness: no critical or major findings; minor hard-coded event-count
comment accepted and fixed
tracing/evidence and safety: initial critical concern accepted as proof gap;
allowlist narrowed and token/authenticated-url negative tests added
focused safety re-review: pass; no critical, major, or minor findings
requirements-vs-implementation replacement review: pass; implementation
satisfies P0-002 without claiming phase/mutation/test closure
```

#### Recheck on 2026-05-10 after PR #35 / digest-only OpenCode tool paths

Status remains `open`.

Verified:

- `sdp-trace` `main` includes `9e1027b Allow digest-only OpenCode tool paths
  (#35)`.
- A live customer-case GSD planning run used the required model route:

```text
sdp-trace observe session --profile session-profile.json --out run -- \
  sh -c 'opencode run --format json \
  --model minimax-coding-plan/MiniMax-M2.5 \
  --dir /Users/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd \
  "/gsd-plan-phase 1 --skip-research --text" \
  --dangerously-skip-permissions > opencode-gsd-run.jsonl'
```

The run reached real GSD activity. The raw stream contained a completed
`tool_use` event with `tool:"task"`, `subagent_type:"gsd-planner"`, model
metadata `minimax-coding-plan/MiniMax-M2.5`, and phase-planning output for
Phase 1. The demo repository still had only harness-generated `.opencode/` and
`.planning/` artifacts; no trace artifact was hand-edited.

Closure failed because the native OpenCode/GSD stream is still rejected before
observed run creation:

```text
raw source line 2: unsafe_input:part.state.input.prompt:forbidden_raw_field
```

The session metadata again retained only command-level facts:

```text
"command_model": "minimax-coding-plan/MiniMax-M2.5"
"command_model_state": "pass"
"process_id_state": "pass"
"collection_state": "cannot_verify"
"collection_reason": "not_collected"
```

This proves PR #35 fixed only the first path-field rejection. The broader P0
remains: `sdp-trace` cannot yet ingest a valid real OpenCode/GSD event stream
without aborting on native tool input fields. Manual removal or filtering of
`part.state.input.prompt` would be trace artifact tampering and cannot be used
as demo evidence.

Required closure now needs a live replay of this `gsd-plan-phase` shape where
`part.state.input.prompt` is either safely summarized, digest-retained, or
marked unavailable at field level without aborting the whole collection. The
resulting observed run must then pass `harness validate` for harness,
interaction, model, and tool evidence while keeping phase/mutation/test/PR gaps
explicit as `not_assessed` or `cannot_verify`.

#### Implementation response on 2026-05-10

| id | severity | plane | finding | disposition | response |
| --- | --- | --- | --- | --- | --- |
| B31-RETURN-11 | critical | tracing/evidence | Real OpenCode/GSD `tool_use` task events with `part.state.input.prompt` were rejected as `forbidden_raw_field`, so no observed run could be created for the `gsd-plan-phase` customer-case stream. | accepted_fixed | Raw OpenCode scanning now treats string-valued `part.state.input.prompt` as an unretained native tool input body under digest-only normalization. Top-level `prompt` remains forbidden, and sensitive nested keys still remain fatal. |
| B31-RETURN-12 | major | safety | Review found that the `*.input.prompt` allowance needed tighter shape constraints so non-string prompt objects could not bypass nested safety scanning. | accepted_fixed | The allowance now applies only to string-valued `input.prompt` fields with exact path segments; object-valued prompt remains forbidden before normalized output is written. |

Regression coverage now extends the real GSD shape with a native OpenCode
`tool_use` task event containing `input.prompt` and `subagent_type:
gsd-planner`. The observed run is created, and normalized output is asserted
not to retain the prompt body or `gsd-planner` input metadata.

Observed result from the regression:

```text
SESSION
"command_model": "minimax-coding-plan/MiniMax-M2.5"
"command_model_state": "pass"
"collection_state": "pass"

RUN
"event_count": 9
events/raw-000001-harness.json
events/raw-000002-interaction.json
events/raw-000003-tool.json
events/raw-000004-tool.json
events/raw-000005-tool.json
events/raw-000006-model.json
events/raw-000006-tool.json
events/raw-000007-tool.json
events/session-command-model.json

VALIDATION
"validation_state": "not_assessed"
harness: pass
interaction: pass
model: pass
tool: pass
phase/mutation/test: not_assessed
```

This closes only the prompt-field collection-abort blocker. The demo P0 remains
open for phase, mutation, and test evidence until those signals are emitted as
verifier-backed events in the customer-case delivery loop.

Review disposition:

```text
requirements-vs-implementation: pass; native gsd-plan-phase input.prompt
ingests without manual filtering, prompt body is not retained, and top-level
prompt remains forbidden
tracing/evidence: pass; no overclaim beyond collection-abort closure
code/safety: initial critical/major concerns accepted; allowance narrowed to
string-valued part.state.input.prompt only, object prompt remains forbidden
focused safety re-review: pass; no remaining critical or major findings
```

#### Recheck on 2026-05-26 during spec closure route

Status remains `open`.

Verified:

- `opencode --version` returned `1.15.10`.
- `opencode models | rg -i 'minimax|m2'` returned current provider-qualified
  MiniMax routes including `minimax/MiniMax-M2.5`.
- External demo repository `fall-out-bug/sdp-trace-demo-jvm-gsd` was cloned and
  used as the harness target.
- `sdp-trace observe session` created session and observed-run directories with
  setup isolation, command digest, process id, source commit, time bounds, raw
  OpenCode JSONL collection, and normalized digest states.

Attempted customer-case replay:

```text
sdp-trace observe session --profile .tmp-t226-run/session-profile.json \
  --out .tmp-t226-run/session-run -- \
  sh -c 'opencode run --format json \
    --model minimax-coding-plan/MiniMax-M2.5 \
    --dir /tmp/.../sdp-trace-demo-jvm-gsd \
    "/gsd-plan-phase 1 --skip-research --text" \
    --dangerously-skip-permissions > .tmp-t226-run/opencode-gsd-run.jsonl'
```

OpenCode 1.15.10 rejected the historical provider-qualified model id:

```text
Model not found: minimax-coding-plan/MiniMax-M2.5. Did you mean: MiniMax-M2.5, MiniMax-M2.5-highspeed?
```

A retry with `--model MiniMax-M2.5` was also rejected:

```text
Model not found: MiniMax-M2.5/.
```

Current-route replay:

```text
sdp-trace observe session --profile tmp/t226-live/session-profile.json \
  --out tmp/t226-live/session-run -- \
  sh -c 'opencode run --format json \
    --model minimax/MiniMax-M2.5 \
    --dir /tmp/sdp-trace-demo-jvm-gsd-full \
    "/gsd-plan-phase 1 --skip-research --text" \
    --dangerously-skip-permissions > tmp/t226-live/opencode-gsd-run.jsonl'
```

Observed session result from the current provider-qualified route:

```text
SESSION
"command_model": "minimax/MiniMax-M2.5"
"command_model_state": "pass"
"process_id_state": "pass"
"source_commit": "50b7ed5ca146a32ce289dc9ee29aa34c8919439d"
"source_commit_state": "pass"
"collection_state": "pass"
"output_digest": "aca2e567edf3541fdf697f634477fad67258937f931d1302dec4d0a2da292a96"
"normalized_digest": "af321c50bcd5627b78bcc49a05d3cc26227727235db3d6d868a8252c9d2f9ddb"

RUN
"event_count": 4
harness: pass
interaction: pass
model: pass

VALIDATION
"validation_state": "not_assessed"
tool/phase/mutation/test: not_assessed
```

This recheck proves the current `sdp-trace` session path can collect and
normalize OpenCode 1.15.10 output with the current MiniMax route while retaining
setup metadata, command digest, source commit, time bounds, output digest, and
normalized digest. It does not close T226 because the current OpenCode/GSD
environment did not expose `/gsd-plan-phase`; the observed run stopped at a
text response that no matching skill was available, so no customer-case tool,
phase, mutation, or test evidence was emitted. At that point T226 remained open
until a current GSD route became available and produced a real first-run
delivery-loop observation, with unavailable dimensions preserved as
`not_assessed` or `cannot_verify`.

#### Recheck with GSD-Redux on 2026-05-26

Status superseded by the closure replay below.

The historical GSD route was replaced locally in the active demo checkout with
GSD-Redux:

```text
npx -y @opengsd/get-shit-done-redux@latest --opencode --local --profile=core
```

Observed installer result:

- package: `@opengsd/get-shit-done-redux`;
- version: `1.1.0`;
- runtime: OpenCode;
- install scope: local demo repository `.opencode`;
- installed commands included `/gsd-plan-phase`;
- global OpenCode configuration was not used for the replacement.

Live replay:

```text
sdp-trace observe session --profile tmp/t226-gsd-redux-live/session-profile.json \
  --out tmp/t226-gsd-redux-live/session-run -- \
  sh -c 'opencode run --format json \
    --model minimax/MiniMax-M2.5 \
    --dir /tmp/sdp-trace-demo-jvm-gsd-full \
    "/gsd-plan-phase 1 --skip-research --text --skip-verify" \
    --dangerously-skip-permissions > tmp/t226-gsd-redux-live/opencode-gsd-redux-run.jsonl'
```

Observed session result:

```text
"command_model": "minimax/MiniMax-M2.5"
"command_model_state": "pass"
"process_id_state": "pass"
"source_commit": "a4d1f755552ba1f411af5edcb7d6caf24a9c39bf"
"source_commit_state": "pass"
"collection_state": "pass"
"output_digest": "23490b30622df97d766cae84c06b7e91235896bb9b72f17d140c55feec1e5359"
"normalized_digest": "90ebbb3fc5352e08f658baf8c09bb3d6db77c1a07ded67ffbb29fa74d6b78b35"
```

Validation result:

```text
"event_count": 39
harness: pass, 24 events
tool: pass, 13 events
interaction: pass, 1 event
model: pass, 1 event
phase: not_assessed
mutation: not_assessed
test: not_assessed
```

The GSD-Redux route is materially better than the previous recheck:
`/gsd-plan-phase` is now available and emits real OpenCode tool activity over
the demo planning files. At that point T226 still could not close because
current normalization did not produce `phase`, `mutation`, or `test` event
families from this run. That made the remaining issue a `sdp-trace`
route/normalizer evidence gap, not a missing `/gsd-plan-phase` command.

#### Closure replay with GSD-Redux on 2026-05-26

Status: `accepted_closed_with_not_assessed_dimensions`.

The previous replay used slash-command text as a message. OpenCode `run` does
not execute local command files that way; it requires `--command`. The working
route is:

```text
opencode run --format json --command gsd-plan-phase \
  --model minimax/MiniMax-M2.5 \
  --dir /tmp/sdp-trace-demo-jvm-gsd-full \
  "1 --skip-research --text --skip-verify" \
  --dangerously-skip-permissions
```

and:

```text
opencode run --format json --command gsd-execute-phase \
  --model minimax/MiniMax-M2.5 \
  --dir /tmp/sdp-trace-demo-jvm-gsd-full \
  "1 --interactive" \
  --dangerously-skip-permissions
```

Two `sdp-trace` normalizer issues were found and fixed during this replay:

- session collection now regenerates normalized events from
  `raw_event_source_path` instead of reusing a stale existing
  `event_source_path`;
- GSD-Redux phase metadata paths such as `phase_dir` and `verification_path`
  now classify as `phase` evidence without retaining raw command bodies.

Plan-phase replay result:

```text
SESSION
"command_model": "minimax/MiniMax-M2.5"
"command_model_state": "pass"
"process_id_state": "pass"
"source_commit": "a4d1f755552ba1f411af5edcb7d6caf24a9c39bf"
"source_commit_state": "pass"
"collection_state": "pass"
"output_digest": "e0feeba3144dd209125ccdf89f97e8c7f119279fb639bb062c22f25e436d493f"
"normalized_digest": "e5bb7d238f79e4d5cba04cebc9c05f981ae8d1f5f780aa60a387206ef0d3846d"

VALIDATION
"validation_state": "not_assessed"
harness: pass, event_count=40
interaction: pass, event_count=1
model: pass, event_count=1
phase: pass, event_count=4
tool: pass, event_count=19
mutation: not_assessed, event_count=0
test: not_assessed, event_count=0
```

Execute-phase replay result:

```text
SESSION
"command_model": "minimax/MiniMax-M2.5"
"command_model_state": "pass"
"process_id_state": "pass"
"source_commit": "a4d1f755552ba1f411af5edcb7d6caf24a9c39bf"
"source_commit_state": "pass"
"collection_state": "pass"
"output_digest": "729f7765b1c387e22c69f05f2b32d044f67bb3d2cf93ac273e78d2cfc83e1618"
"normalized_digest": "57b9f4c633188641c518480a21317c75d606e5bd7cb1175332d255254cb2151b"

VALIDATION
"validation_state": "not_assessed"
harness: pass, event_count=22
interaction: pass, event_count=1
model: pass, event_count=1
phase: pass, event_count=4
tool: pass, event_count=10
mutation: not_assessed, event_count=0
test: not_assessed, event_count=0
```

The missing mutation/test dimensions are not promoted to green. The GSD-Redux
workflow reported the target phase already complete (`incomplete_count: 0`) and
`gsd-execute-phase 1 --interactive` took no file-mutation or test action. That
is a replay fact, not product success evidence.

T226 is closed for the customer-usable first-run observation path because the
delivery loop is now observable without prompt relay, hand-authored events, or a
customer-built adapter, and unavailable dimensions are retained as
`not_assessed`. This does not claim feature delivery, harness compliance, test
success, PR approval, merge approval, production trust, or broad GSD/OpenCode
support.

Focused diff review:

```text
model: openrouter/z-ai/glm-4.7
scope: stale raw normalization, GSD-Redux phase classification, overclaiming,
missing tests
result: LGTM
```

## P1

No P1 findings recorded yet.

## P2

No P2 findings recorded yet.

## P3

No P3 findings recorded yet.
