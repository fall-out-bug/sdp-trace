# Block 14 Implementation Review Disposition

Date: 2026-05-06

Scope:

- `internal/demo/demo.go`
- `internal/demo/demo_test.go`
- `cmd/sdp-trace/main.go`
- `cmd/sdp-trace/main_test.go`
- `internal/trace/store.go`
- `internal/trace/types.go`
- `schema/gate-result.schema.json`
- `examples/block14-gate/`

Review models:

- MiniMax-M2.7, no tools, no context files.
- ZAI GLM-5.1, no tools, no context files.

## MiniMax-M2.7 Findings

| Severity | Finding | Disposition |
|---|---|---|
| Critical | Gate exit code ignored CI/audit dimensions. | Accepted. `runGate` now uses `gateExitCode`, with `fail` before `cannot_verify` before pass. Tests updated for Block 14 semantics. |
| Critical | Override requests were not populated in gate output. | Accepted. `rowFromRun` collects `policy_override_requested` events and `EvaluateGate` exposes them in `override_requests`. |
| Critical | Witness binding did not check artifact digests. | Accepted. `WriteGate` computes current run manifest digests and witness binding fails on mismatched digest. |
| Critical | `AppendRunEvent` was dead or missing. | Rejected as stale. The function exists in `internal/trace/store.go`, and CLI override tests pass against it. |
| Major | Explain/preview secret leakage risk needed negative coverage. | Already covered before review; retained and re-verified. |
| Minor | Missing override-does-not-pass test. | Accepted. Added CLI test proving override visibility does not convert missing evidence to pass. |

Post-fix MiniMax-M2.7 review reported no remaining critical or major findings.

## ZAI GLM-5.1 Findings

| Severity | Finding | Disposition |
|---|---|---|
| Major | Required array fields could serialize as `null`, violating schema. | Accepted. `EvaluateGate` now initializes arrays, and schema/sample fixtures require `witness_bindings` and `override_requests`. |
| Minor | `runs` schema is intentionally loose. | Deferred. Current producer is Go-owned; tighter consumer schema can be added after stable run-row output is finalized. |
| Minor | `witness.status` is not enum-constrained in Go. | Deferred to schema validation and future typed witness result model. |
| Minor | Reason ordering was not severity-perfect. | Partially accepted. Reasons and next actions are sorted deterministically; richer severity ordering can follow if user-facing ordering proves weak. |

## Verification After Review Fixes

- `rtk go test ./...` -> pass, 58 tests across 10 packages.
- `rtk jq empty schema/*.json examples/block14-gate/*.json` -> pass.
- `rtk git diff --check` -> pass.

## Pre-PR Review

Fresh pre-PR pi review:

- MiniMax-M2.7 initially reported critical/major concerns on exit-code
  robustness, preview witness mismatch surfacing, and explain completeness.
  Exit-code robustness and preview mismatch surfacing were accepted and fixed.
  Explain already printed reasons and next actions; that finding was stale
  against the provided code.
- ZAI GLM-5.1 reported no critical or major findings. It listed minor issues
  around preview mismatch surfacing, top-level `missing_telemetry` exit-code
  handling, missing audit evidence display in explain, and override ordering.
  The cheap fixes were applied before PR creation.

PR-level pi review for PR #4:

- MiniMax-M2.7 reported LGTM with no critical findings and minor follow-ups.
- ZAI GLM-5.1 reported one major finding: a CI witness could omit an expected
  run artifact and still pass binding. Accepted and fixed with
  `TestGateCommandCannotVerifyWhenWitnessOmitsRunArtifact`.

## Residual State

- `audit_grade_gate` remains `cannot_verify`.
- External production trust remains `not_assessed`.
- Source-bound release proof regeneration was not performed in this chunk.
