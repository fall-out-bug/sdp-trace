# Slice 71 Plan Review

Status: pass
Date: 2026-06-04

## Scope

Slice 71 is bounded to numbered `cmd/sdp-trace` `pr-review synthesize` ledger
collation shards:

- `cmd/sdp-trace/pr_review_109_runsynthesize.go`
- `cmd/sdp-trace/pr_review_110_parsesynthesizeargs.go`
- `cmd/sdp-trace/pr_review_111_synthesisinputs.go`
- `cmd/sdp-trace/pr_review_112_readsynthesisinputs.go`
- `cmd/sdp-trace/pr_review_113_readoptionalledger.go`
- `cmd/sdp-trace/pr_review_114_readexistingledger.go`

## Decision Gate

- Simpler/Faster: rename/move only; no behavior change, public API change,
  dependency change, or baseline update.
- Blocking Edge Cases: synthesis is trust-sensitive evidence collation. Output
  must remain durable, packet/run-set read failures must stay
  `cannot_verify`, optional prior ledger handling must not become authority,
  and stdout must keep mirroring the durable ledger artifact.
- Existing Open Source: no new CLI parser, ledger engine, or JSON writer is
  needed; existing package-local helpers and the internal `prreview` package
  are sufficient.

## Planned File Boundary

- `cmd/sdp-trace/pr_review_synthesize_command.go`: synthesize execution,
  argument parsing, output path validation, durable ledger write, and stdout
  mirroring.
- `cmd/sdp-trace/pr_review_synthesis_inputs.go`: packet, run-set, and optional
  existing-ledger input loading.

Pre-change MI probe for the planned boundary:

```text
pr_review_synthesize_command.go maintainability_index=73.5
pr_review_synthesis_inputs.go maintainability_index=72.8
```

## Planned Regression Evidence

- Exact focused test existence:
  - `TestParsePRReviewSynthesizeArgsKeepsUsageBoundaries`
  - `TestReadPRReviewSynthesisInputsKeepsOptionalLedgerBoundary`
  - `TestRunPRReviewSynthesizeKeepsLedgerDurability`
- Focused behavior checks: mandatory output path validation, positional
  argument rejection, packet read failure `cannot_verify` behavior, run-set
  read failure `cannot_verify` behavior, optional existing-ledger empty path,
  existing ledger read failure propagation, durable ledger write before stdout
  mirroring, and exact test list verification.
- Standard verification, conditional `golangci-lint run`, and CRAP/MI gates
  remain required before the implementation review.

## Review Rounds

### Round 1

- scope/correctness: major finding. Planned focused regression evidence did
  not explicitly cover run-set read failure mapping to `cannot_verify`.
- trust/evidence: `LGTM`
- maintainability/DX: major finding. Planned focused regression evidence did
  not explicitly cover run-set read failure mapping to `cannot_verify`.

### Round 2

- scope/correctness: `LGTM`
- maintainability/DX: `LGTM`

Final plan-review status: `LGTM` across scope/correctness, trust/evidence,
and maintainability/DX after fixes.
