# Slice 70 Plan Review

Status: pass
Date: 2026-06-04

## Scope

Slice 70 is bounded to numbered `cmd/sdp-trace` `pr-review run` execution
shards:

- `cmd/sdp-trace/pr_review_105_runrun.go`
- `cmd/sdp-trace/pr_review_106_executerun.go`
- `cmd/sdp-trace/pr_review_107_writerunoutput.go`
- `cmd/sdp-trace/pr_review_108_parserunargs.go`

## Decision Gate

- Simpler/Faster: rename/move only; no behavior change, public API change,
  dependency change, or baseline update.
- Blocking Edge Cases: a single `pr_review_run_command.go` failed pre-change
  file MI at `66.1`. Runner allow-lists, preview mode, work-dir validation, and
  `cannot_verify` mapping are trust-sensitive and need focused regression
  evidence.
- Existing Open Source: no new CLI parser, runner framework, or output renderer
  is needed; existing package-local helpers and the internal `prreview` package
  are sufficient.

## Planned File Boundary

- `cmd/sdp-trace/pr_review_run_command.go`: run execution and reviewer runner
  handoff.
- `cmd/sdp-trace/pr_review_run_args.go`: run argument parsing and usage
  boundaries.
- `cmd/sdp-trace/pr_review_run_output.go`: preview-vs-run output rendering.

## Planned Regression Evidence

- Exact focused test existence:
  - `TestParsePRReviewRunArgsKeepsUsageBoundaries`
  - `TestExecutePRReviewRunKeepsRunnerOptions`
  - `TestWritePRReviewRunOutputKeepsPreviewBoundary`
- Focused behavior checks: parse defaults, parse errors and positional argument
  rejection, packet/profile read failure `cannot_verify` behavior, work-dir
  directory validation, repeated allowed-runner reconstruction from raw args,
  preview mode and not-assessed reason propagation, and preview output shape not
  implying produced run evidence.
- Standard verification, conditional `golangci-lint run`, and CRAP/MI gates
  remain required before the implementation review.

## Review Rounds

### Round 1

- scope/correctness: `LGTM`
- trust/evidence: `LGTM`
- maintainability/DX: major finding. The plan deletes/replaces
  `pr_review_105` through `pr_review_108` but did not require synchronizing
  `cmd/sdp-trace/FAMILY_INDEX.md`, repeating the stale-index issue from Slice
  69.

### Round 2

- maintainability/DX: `LGTM`

Final plan-review status: `LGTM` across scope/correctness, trust/evidence,
and maintainability/DX after fixes.
