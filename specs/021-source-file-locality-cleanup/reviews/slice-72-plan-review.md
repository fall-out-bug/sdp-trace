# Slice 72 Plan Review

Status: pass
Date: 2026-06-04

## Scope

Slice 72 is bounded to numbered `cmd/sdp-trace` `pr-review validate`
artifact-validation shards:

- `cmd/sdp-trace/pr_review_115_runvalidate.go`
- `cmd/sdp-trace/pr_review_116_parsevalidateargs.go`
- `cmd/sdp-trace/pr_review_117_validationcliexitcode.go`
- `cmd/sdp-trace/pr_review_118_validationinputs.go`
- `cmd/sdp-trace/pr_review_119_readvalidationinputs.go`
- `cmd/sdp-trace/pr_review_120_readvalidationartifacts.go`

## Decision Gate

- Simpler/Faster: rename/move only; no behavior change, public API change,
  dependency change, or baseline update.
- Blocking Edge Cases: validation joins packet, profile, run-set, and ledger
  artifacts. Every unreadable artifact must remain `cannot_verify`; validation
  verdicts must remain package-owned; stdout must mirror a durable validation
  artifact; and non-zero validation verdicts must map through
  `exitCannotVerify`.
- Existing Open Source: no new CLI parser, validation engine, or JSON writer is
  needed; existing package-local helpers and the internal `prreview` package are
  sufficient.

## Planned File Boundary

- `cmd/sdp-trace/pr_review_validate_command.go`: validate execution, argument
  parsing, output path validation, durable validation write, stdout mirroring,
  and validation-exit mapping.
- `cmd/sdp-trace/pr_review_validation_inputs.go`: packet, profile, run-set, and
  ledger input loading.

Pre-change MI probe for the planned boundary:

```text
pr_review_validate_command.go maintainability_index=72.0
pr_review_validation_inputs.go maintainability_index=73.2
```

## Planned Regression Evidence

- Exact focused test existence:
  - `TestParsePRReviewValidateArgsKeepsUsageBoundaries`
  - `TestReadPRReviewValidationInputsKeepsArtifactBoundaries`
  - `TestRunPRReviewValidateKeepsDurableVerdictAndExitMapping`
- Focused behavior checks: mandatory output path validation, positional
  argument rejection, packet/profile/run-set/ledger read failure
  `cannot_verify` behavior, durable validation write before stdout mirroring,
  non-zero validation verdict mapping through `exitCannotVerify`, and exact
  test list verification.
- Standard verification, conditional `golangci-lint run`, and CRAP/MI gates
  remain required before the implementation review.

## Review Rounds

### Round 1

- scope/correctness: `LGTM`
- trust/evidence: `LGTM`
- maintainability/DX: `LGTM`

Final plan-review status: `LGTM` across scope/correctness, trust/evidence,
and maintainability/DX.
