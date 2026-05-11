# Block 31 Implementation Plan: Non-Interfering Harness Observation

Status: Reopened for T226 first-run customer-case closure. The generic
`harness observe` / `validate` / `summarize` intake is necessary but
insufficient until Block 31 also proves a "set up and forget" OpenCode/GSD path
for the demo customer case.

Spec:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/31-non-interfering-harness-observation.md`

## Slice 0: Socratic Review And Disposition

Write scope:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/31-non-interfering-harness-observation-review-ledger.md`
- Block 31 spec and implementation plan if review findings require changes

Actions:

1. Run Socratic review across product boundary, UX/DX, trace/evidence,
   safety/privacy, and implementation-feasibility planes.
2. Record every valid critical or major finding with disposition.
3. Fix every accepted critical or major finding before implementation starts.
4. Stop for explicit approval of the reviewed direction.

No product code may be written in this slice.

## Slice 1: Contracts And Fixtures

Write scope:

- `schema/harness-observation-profile.schema.json`
- `schema/harness-event.schema.json`
- `schema/harness-observation-run.schema.json`
- `schema/harness-observation-validation.schema.json`
- `examples/harness-observation/`
- `internal/harnessobs/` contract tests if useful

Actions:

1. Define profile, event, run, and validation schemas.
2. Add focused fixtures for:
   - harness-generic complete export;
   - complete OpenCode/GSD export;
   - first-run OpenCode/GSD setup metadata;
   - first-run OpenCode/GSD observed session output;
   - zero-event source;
   - missing model route;
   - missing phase or review events;
   - prompt digest only;
   - tool-event gap;
   - unsafe raw prompt;
   - unsafe source ref;
   - symlink escape;
   - source digest mismatch;
   - schema version mismatch;
   - mutation without source binding;
   - absent PR state;
   - no run supplied.
3. Add fixture matrix with expected top-level and dimension states.
4. Select a Go JSON Schema validator before implementation starts. The default
   candidate is `github.com/santhosh-tekuri/jsonschema/v6` behind a thin local
   wrapper; if the dependency is rejected during review, implement equivalent
   focused Go contract validation before adding CLI behavior.

Expected verification:

- `jq empty schema/*.json examples/harness-observation/*.json`
- focused Go tests for schema/fixture alignment; JSON syntax alone is not
  sufficient

## Slice 2: Observe Command

Write scope:

- `cmd/sdp-trace/main.go`
- `cmd/sdp-trace/*harness*_test.go`
- `internal/harnessobs/`

Actions:

1. Add `sdp-trace harness observe --profile <file> --source <jsonl> --out <dir>`.
2. Read explicit files only; do not invoke OpenCode, GSD, GitHub, provider APIs,
   or hidden shell commands.
3. Process JSONL line by line with bounded memory. Default limits: 1 MiB per
   event line and 100,000 events per source unless a reviewed profile specifies
   a smaller cap.
4. Refuse an existing non-empty `--out` directory unless a reviewed overwrite
   flag is added.
5. Normalize safe event refs, source digests, content states, unavailable
   fields, and profile identity into an observed run directory.
6. Reject unsafe paths, symlink escapes, URL-like refs, token-like values,
   forbidden raw prompt/model fields, malformed JSONL, and digest mismatches
   before writing a run.

Expected verification:

- focused CLI tests for successful observation and unsafe input rejection
- `go test ./cmd/sdp-trace ./internal/harnessobs`

## Slice 3: First-Run Customer Path

Write scope:

- `cmd/sdp-trace/main.go`
- `cmd/sdp-trace/*harness*_test.go`
- `internal/harnessobs/`
- `examples/harness-observation/`
- command docs only after behavior exists

Actions:

1. Add a customer-usable path for the OpenCode/GSD case that allows bounded
   setup before delivery, then observes the real harness workflow without prompt
   relay or extra in-loop operator chores.
2. Add a generic session profile schema or focused Go contract for setup
   actions, safe env names, declared log paths, declared output directories,
   stream capture policy, raw-surface-to-`harness-event-v1` mappings, redaction
   rules, context isolation rules, and external-tool requirements.
3. Ship the OpenCode/GSD mapping as a checked-in profile example under
   `examples/harness-observation/`, not as hidden Go special-case logic.
4. The reviewed CLI shape may change, but profile resolution must use a file path
   by default or an explicit `builtin:` prefix. Target workflow:

   ```text
   sdp-trace observe setup --profile <session-profile.json> --out <run-dir> [--command <harness-command-preview>]
   <normal harness command>
   sdp-trace observe collect --profile <session-profile.json> --run <run-dir>
   ```

   Optional single-command convenience wrapper:

   ```text
   sdp-trace observe session --profile <session-profile.json> --out <run-dir> -- <harness-command>
   ```

5. Setup is limited to the spec-defined bounded setup actions: one
   initialization command, one profile/configuration file selection, and one
   reviewed wrapper or hook installation only when the profile requires it. It
   must be explicit, reviewable, and documented as pre-work; it must not require
   prompt rewriting, prompt relay, or manual trace authoring during delivery.
   When a profile declares context isolation, setup must install and verify the
   declared local file rules before the delivery loop and record rule-level
   states and target-file digests in session evidence.
6. During the harness run, use only allowed observation mechanisms: process
   boundary capture, stdout/stderr digests or retained-safe excerpts only when
   declared safe by the profile, declared log tailing, declared output-directory
   watching, and filesystem artifact reads. Do not inject stdin, rewrite harness
   arguments after setup, mutate undeclared environment, hide PATH rewrites, or
   interpose on provider network calls.
7. During the harness run, collect or normalize the evidence needed by
   `harness observe`, `harness validate`, and `harness summarize` for the same
   run: harness identity, model route, interaction boundaries, phase/review
   activity, tool/command/file mutation evidence, test observations, PR state,
   merge state, and explicit unavailable dimensions.
8. Emit `harness-event-v1` JSONL or an observed run directory with equivalent
   normalized event content. Include setup metadata with setup actions, profile
   id, harness command digest, process id or unavailable reason, start/end time
   bounds, source commit, and output artifact digests.
9. Preserve missing or unsupported fields as `not_assessed` or `cannot_verify`
   without weakening the core claim that the delivery loop was observed.
10. Add acceptance coverage where fixture-only validation remains a separate test
   class and does not satisfy P0-001 closure.

Expected verification:

- focused CLI tests for setup-once/run-normal behavior against a live or
  controlled-proxy harness command invocation
- fixture replay proving the first-run output can be consumed by `harness
  observe`, `harness validate`, and `harness summarize`
- setup provenance tests proving the first-run output is bound to command
  invocation, time bounds, source commit, and output artifact digests
- safety tests proving the first-run path does not print or persist raw prompts,
  raw model responses, provider tokens, authenticated URLs, or private paths
- `go test ./cmd/sdp-trace ./internal/harnessobs`

## Slice 4: Validate And Summarize Commands

Write scope:

- `cmd/sdp-trace/main.go`
- `cmd/sdp-trace/*harness*_test.go`
- `internal/harnessobs/`
- `schema/harness-observation-validation.schema.json`

Actions:

1. Add `sdp-trace harness validate --profile <file> --run <dir> --out <file>`.
2. Validate required and optional dimensions against the profile degradation
   rules.
3. Preserve dimension-level `pass`, `fail`, `not_assessed`, and `cannot_verify`.
4. Add `sdp-trace harness summarize --validation <file>` with safe output only.
5. Ensure no summary implies harness compliance, feature delivery, PR approval,
   merge approval, production trust, or buyer-facing trust.

Expected verification:

- focused CLI tests for complete, gap, unsafe, and absent-run fixtures
- safety tests proving summaries do not print synthetic secret markers, raw
  prompts, raw model responses, authenticated URLs, or private paths
- state-composition tests for mixed required/optional dimensions, zero-event
  source, absent source, and adapter/harness conflict

## Slice 5: Docs, Ledger, And Drift Checks

Write scope:

- `docs/agent-entrypoint.md`
- `docs/reviewer-entrypoint.md`
- `docs/reviews/demo-jvm-gsd-observation-ledger.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md`
- Block 31 review ledger

Actions:

1. Update command docs only after commands exist and tests pass.
2. Update the demo observation ledger to reference implemented validation
   evidence.
3. Keep P0-001 open unless the customer-case OpenCode/GSD workflow has been
   observed through the supported first-run path after bounded setup.
4. Run drift checks against Blocks 19, 28, 29, 30, and 31.
5. Run implementation review across code/correctness, trace/evidence, and
   requirements-vs-implementation planes; fix and re-review accepted critical or
   major findings.

Expected verification:

- `go test ./...`
- `jq empty schema/*.json examples/harness-observation/*.json`
- `git diff --check`

## Reopened Approval Gate

The original reviewed direction produced a generic intake but did not close the
customer-case P0. Before writing the first-run customer-path code, update the
Block 31 review ledger with this T226 correction, run a focused Socratic
re-review on the first-run path, fix accepted critical or major findings, and
get explicit approval of the narrowed direction.

Until the first-run path is implemented and verified against the customer-case
workflow, the demo P0 remains open.
