# Follow-Up Hardening Handoff Ledger

**Date**: 2026-05-17
**Coordinator**: Codex
**External worker harness**: `codex-subagent` `0.3.0` with Pi runtime

## Coordinator Commits

| Commit | Purpose | Verification |
|---|---|---|
| `01a245b` | Added follow-up hardening SpecKit delta, tasks, Kimi review prompt, Kimi readiness note, and Pi review ledger. | `git diff --check` before commit; Kimi PI review recorded in `followup-hardening-pi-review.md`. |
| `1b936da` | Updated rules to request explicit model selection after default Pi runs stalled. | Later corrected: the stall was not proven to be a model failure; it is an orchestration/profile/context observability finding until logs prove root cause. |
| `8992a57` | Slice 4 high-signal duplication cleanup. | `golangci-lint run --disable-all --enable=dupl ./...`, focused package tests, `go test -count=1 ./...`, and `git diff --check` passed locally. |
| `7522fe1` | Slice 5 gocritic/unparam/prealloc cleanup. | `golangci-lint run --disable-all --enable=gocritic --enable=unparam --enable=prealloc ./...`, focused package tests, `go test -count=1 ./...`, and `git diff --check` passed locally. |
| `a023696` | Slice 6 stale MI claim correction. | Absolute file/function MI replay failed locally; docs now mark absolute MI as open and baseline ratchet as the CI policy. |
| `aeb0429` | Releaseproof commit SHA validator CRAP cleanup. | Fresh strict CRAP replay passed after splitting the validator helper. |

## Subagent Runs

| Run ID | Scope | Runtime/model | State | Disposition |
|---|---|---|---|---|
| `run_7FWrtNDYCb` | Full follow-up package | Pi default model | `cancelled` | `cannot_verify`: run produced only partial Slice 2 format drift and then no status/result/log progress. Not counted as closure. |
| `run_HRz3AfNM82` | Slice 2 format/import | Pi default model | `cancelled` | `cannot_verify`: no diff/log/result progress. Root cause not proven; treat as Pi orchestration/profile/context observability issue, not a default-model defect. |
| `run_XROAfGeIU6` | Slice 2 format/import | `kimi-coding/kimi-for-coding` | `pass` | `accepted_integrated`: produced commit `dfdd470`; focused local verification passed in main. |
| `run_Wauju2IBKg` | Slice 3 releaseproof hardening | `kimi-coding/kimi-for-coding` | `pass` | `accepted_integrated`: produced commit `3718e90`, integrated as main commit `4916b1f`; focused local verification passed in main. |
| `run_9u6G4Ax__O` | Slice 6 docs/MI overclaim closure | `kimi-coding/kimi-for-coding` | `cancelled` | `cannot_verify`: no diff/log/result progress after several minutes. Slice 6 remains open. |
| `run_j_ZbHNgXRr` | Slice 4 duplication cleanup | `kimi-coding/kimi-for-coding` | `cancelled` | `cannot_verify`: process disappeared while run remained marked `running`; no result or worktree diff. Coordinator completed Slice 4 locally. |
| `panel_NRfzTXlC8S` | Final review panel | `kimi-coding/kimi-for-coding` | `cancelled` | `cannot_verify`: three role runs remained running with no result files; direct Pi reviews replaced it. |
| direct Pi Qwen | Final security/trust review | `openrouter/qwen/qwen3-coder` | `pass` | `no actionable findings`; output stored outside source tree for session evidence. |
| direct Pi DeepSeek/Qwen | Final code/correctness review | `openrouter/deepseek/deepseek-chat-v3.1`, fallback `openrouter/qwen/qwen3-coder` | `reviewed` | DeepSeek produced empty output; Qwen fallback raised one artifactBytes validation concern, rejected as false positive after checking `internal/releaseproof/artifacts.go` because the nearest `git show` boundary must retain its own source_commit validation. |
| direct Pi GLM | Final docs/evidence review | `openrouter/z-ai/glm-4.6` | `findings_fixed` | Found stale date and unchecked final-review task state; this commit updates the ledger. |

## Integrated Implementation Commits

| Commit | Slice | Evidence |
|---|---|---|
| `dfdd470` | Slice 2 - Format/import target files | Worker result recorded `gofmt`, `go test -count=1 ./...`, and target goimports/gofmt pass. Coordinator replayed `gofmt -l` on Slice 2 target files and focused tests: `go test -count=1 ./cmd/sdp-trace ./internal/demo ./tools/doccheck ./tools/schemadoc`. |
| `4916b1f` | Slice 3 - Releaseproof `source_commit` hardening | Coordinator replayed `go test -count=1 ./internal/releaseproof ./cmd/sdp-trace`, `jq empty schema/*.json`, `golangci-lint run --disable-all --enable=gosec ./internal/releaseproof/...`, and `git diff --check`. |
| `8992a57` | Slice 4 - High-signal duplication cleanup | Coordinator replayed `golangci-lint run --disable-all --enable=dupl ./...`, focused package tests, `go test -count=1 ./...`, and `git diff --check`. |
| `7522fe1` | Slice 5 - Gocritic/unparam/prealloc cleanup | Coordinator replayed `golangci-lint run --disable-all --enable=gocritic --enable=unparam --enable=prealloc ./...`, focused package tests, `go test -count=1 ./...`, and `git diff --check`. |
| `a023696` | Slice 6 - Maintainability docs overclaim closure | Coordinator replayed absolute file/function MI commands, observed `fail_local`, and `go run ./tools/doccheck` passed after doc correction. |
| `aeb0429` | Final CRAP cleanup | Coordinator replayed `go test -count=1 ./... -coverprofile=/tmp/sdp-trace-coverage.out`, `go tool cover`, `go run ./tools/qualitycheck -gocyclo`, and strict `tools/crapcheck`; all passed. |

## Current Open State

- Slice 4 duplication cleanup: `pass_local` after local coordinator implementation; Pi handoff attempt remains `cannot_verify`.
- Slice 5 gocritic/unparam/prealloc cleanup: `pass_local`.
- Slice 6 docs/MI overclaim closure: `pass_local` for stale-claim correction; Pi handoff attempt remains `cannot_verify`, and absolute MI itself remains `fail_local`.
- External GitHub CI: `not_assessed`.
- PR creation: `not_assessed`; ready to open after pushing the branch.

## Final Local Gate

Final gate replayed on 2026-05-17 at head `aeb0429`:

- `go test -count=1 ./...`: `pass`
- `go vet ./...`: `pass`
- `golangci-lint run ./...`: `pass`
- `jq empty schema/*.json`: `pass`
- `git diff --check`: `pass`
- `go run ./tools/hygienecheck`: `pass`
- `go run ./tools/doccheck`: `pass`
- `go run ./tools/schemadoc`: `pass`
- `go run ./tools/schemadoc -verify-readme`: `pass`
- `govulncheck ./...`: `pass`
- strict CRAP replay from fresh `/tmp/sdp-trace-coverage.out`: `pass`

External GitHub CI remains `not_assessed` until queried live for the pushed PR
head SHA.

## Dirty Checkout Note

The coordinator observed pre-existing local modifications in these files and did not revert them:

- `internal/demo/demo_gate_required_runs.go`
- `internal/demo/demo_gate_result.go`
- `internal/demo/demo_gate_rows.go`
- `internal/demo/demo_override_evidence.go`
- `internal/demo/demo_override_field_state.go`
- `internal/demo/demo_override_required_run.go`
- `internal/demo/demo_report_build.go`
- `internal/demo/demo_required_runs_eval.go`
- `internal/demo/demo_required_runs_match.go`
- `internal/demo/demo_rows_verified.go`

These files are local structural state only and are not external trust evidence.
