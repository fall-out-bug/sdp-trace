# Block 31 Implementation Plan: Non-Interfering Harness Observation

Status: Draft implementation plan. Implementation is blocked until Block 31
Socratic review is complete and the reviewed direction is explicitly approved.

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

## Slice 3: Validate And Summarize Commands

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

## Slice 4: Docs, Ledger, And Drift Checks

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
3. Keep P0-001 open unless a real OpenCode/GSD export has been observed and
   validated.
4. Run drift checks against Blocks 19, 28, 29, 30, and 31.
5. Run implementation review across code/correctness, trace/evidence, and
   requirements-vs-implementation planes; fix and re-review accepted critical or
   major findings.

Expected verification:

- `go test ./...`
- `jq empty schema/*.json examples/harness-observation/*.json`
- `git diff --check`

## Approval Gate

Do not write CLI code, schemas, fixtures, or tests for Block 31 until:

1. Socratic review has been recorded.
2. All valid critical and major findings have been fixed or explicitly blocked.
3. The user explicitly approves the reviewed direction.

Until then, the demo P0 remains open.
