# Implementation Ledger: MVP Readiness Hardening

**Date**: 2026-05-10
**Worktree**: `/Users/fall_out_bug/projects/vibe_coding/sdp-trace-mvp-readiness-spec`
**Branch**: `codex/mvp-readiness-spec`
**Intake baseline commit**: `5f6706b398d6d68bb9a37be2dee4e6aec1037df3`
**Clean PR base**: `origin/main` at `a9ded8a246119e90b1f2ed57d71f72202dcf1ab8`

## Scope State

| Area | State | Evidence |
|---|---|---|
| Spec review | pass | Three-axis Socratic review recorded in `socratic-review.md`; focused re-review returned `APPROVE`; user approved implementation. |
| Docs freshness | pass_local | Stale `pr-review packet` flags removed from `docs/agent-entrypoint.md`; global CLI help does not advertise `--context` or `--verification`. |
| Russian docs parity | deferred_scope | Russian command reference routes command contract to English canonical section until full bilingual parity is implemented. |
| Placeholder examples | pass_local | Placeholder/evidence-boundary labels added to first-class example READMEs. |
| Lint | pass_local | `golangci-lint run ./...` exited 0 after authority and telemetry fixes. |
| CI lint enforcement | pass_ci | `.github/workflows/ci.yml` now runs `go test ./... -coverprofile=coverage.out` and `golangci-lint-action@v6` at `v1.62.0`. GitHub CI `verify` passed on PR #37. |
| CRAP < 5 | assessed_gap | Strict CRAP threshold is not satisfied by existing production code; `tools/crapcheck` computes the baseline and exits non-zero at threshold 5. |
| Complexity over 15 | assessed_gap | Existing production functions remain above `gocyclo -over 15`; `internal/harnessobs.normalizeOpenCodeRawLine`, `internal/harnessobs.LoadSessionProfile`, `internal/harnessobs.CollectSession`, and `internal/trace.writeCanonicalJSON` were decomposed below 15. |
| Coverage hardening | pass_partial | MVP-critical zero-coverage packages `contract`, `export`, and `policy` now have focused tests; `harnessobs`, `trace`, and `verifier` were improved. |

## Command Evidence

| Command | Result | Notes |
|---|---|---|
| `go run ./cmd/sdp-trace --help` | pass | Used as canonical command-surface comparison source. |
| `go run ./cmd/sdp-trace pr-review --help` | fail_expected | CLI does not support nested `--help`; global help is the current source of command contracts. |
| `go run ./cmd/sdp-trace pr-review packet --help` | fail_expected | CLI reports `unknown flag --help`; docs were compared against global help. |
| `rg -n -- '--context\|--verification\|This example will show\|controlled-pilot ready\|sidecar trust substrate' README.md docs examples` | pass_absent | Command exits 1 because no matches remain. |
| `go test ./... -coverprofile=/tmp/sdp-trace-continue-cover.out` | pass | Total coverage: 71.5%. |
| `go test ./tools/crapcheck -cover` | pass | Tool coverage: 50.6%. |
| `golangci-lint run ./...` | pass | No findings after fixes. |
| `go vet ./...` | pass | Modern Go suspicious-construct sweep. |
| `/Users/fall_out_bug/go/bin/staticcheck ./...` | pass | Staticcheck 2026.1 / v0.7.0, no findings. |
| `rg --files -g '*.go' \| xargs /Users/fall_out_bug/go/bin/gopls check` | pass | `gopls` v0.21.1 LSP diagnostics across all Go files. |
| `/Users/fall_out_bug/go/bin/govulncheck ./...` | pass | Rebuilt with Go 1.26.3 as govulncheck v1.3.0; vulnerability DB updated 2026-05-07; no vulnerabilities found. |
| `gocyclo -over 15 .` | fail_assessed_gap | Existing production and test functions exceed 15. |
| `gocognit -over 20 .` | fail_assessed_gap | Existing production and test functions exceed 20. |
| `go run ./tools/crapcheck -cover-func /tmp/sdp-trace-continue-cover-func.txt -gocyclo /tmp/sdp-trace-continue-gocyclo.txt -threshold 5` | fail_assessed_gap | 373 functions exceed strict CRAP threshold 5; top high-risk functions improved, but strict repo-wide CRAP remains open. |

## Coverage Delta

Baseline from intake:

| Package | Baseline | Current |
|---|---:|---:|
| `internal/contract` | 0.0% | 68.1% |
| `internal/export` | 0.0% | 83.3% |
| `internal/policy` | 0.0% | 71.7% |
| `internal/trace` | 2.9% | 61.1% |
| `internal/posture` | 72.4% | 86.1% |
| `internal/harnessobs` | 42.7% | 65.0% |
| `internal/verifier` | 51.1% | 71.0% |
| total | 64.0% | 71.5% |

## CRAP Baseline Summary

Formula:

```text
CRAP = complexity^2 * (1 - coverage)^3 + complexity
```

Strict `CRAP < 5` cannot be claimed for the current repository. Since the
minimum CRAP score equals cyclomatic complexity at 100% coverage, every
function with cyclomatic complexity 5 or higher must be decomposed before the
strict target can pass.

Top current CRAP/complexity findings:

| Function | Cyclo | Coverage | CRAP | State |
|---|---:|---:|---:|---|
| `internal/harnessobs.findUnsafeRawEventAt` | 21 | 0.0% | 462.00 | assessed_gap |
| `internal/interaction.ValidateEvent` | 33 | 47.1% | 194.21 | assessed_gap |
| `internal/harnessobs.safeOutDir` | 23 | 48.8% | 94.00 | assessed_gap |
| `internal/harnessobs.findUnsafeAt` | 20 | 56.0% | 54.07 | assessed_gap |
| `internal/harnessobs.validateEvent` | 20 | 55.6% | 55.01 | assessed_gap |
| `internal/posture.Build` | 19 | 85.3% | 20.15 | assessed_gap |

## Ratchet

Immediate ratchet now enforced or measurable:

- `golangci-lint run ./...` must stay green locally and in CI.
- `tools/crapcheck` provides reproducible CRAP scoring from `go tool cover -func`
  and `gocyclo`.
- New or materially changed production functions in MVP-critical paths should
  target `CRAP < 5`; exceptions must be recorded here as `assessed_gap`.

Next decomposition candidates before stronger MVP-readiness claim:

1. `internal/interaction.ValidateEvent`
2. `internal/harnessobs.safeOutDir`
3. `internal/harnessobs.findUnsafeRawEventAt`
4. `internal/harnessobs.findUnsafeAt`
5. `internal/harnessobs.validateEvent`

## External Evidence Boundary

GitHub CI `verify` passed on PR #37. Merge approval remains `not_assessed`;
this draft PR must not be treated as approved to merge.

## Implementation Pi Review

| Plane | Reviewer | Verdict | Disposition |
|---|---|---|---|
| Code/correctness | `zai/glm-5.1` | APPROVE | Minor dead branch in `signalMetricMatches` accepted and fixed. Minor `pr-review check` concern rejected with global help evidence: the command is present in `go run ./cmd/sdp-trace --help`. |
| Trust/evidence | `openrouter/qwen/qwen3.6-plus` | APPROVE | Low note to verify removed `pr-review packet` flags against live parsing; accepted as already verified by global help and stale-flag search. |
| Requirements-vs-implementation | `minimax/MiniMax-M2.7` | APPROVE | No required fixes. Reviewer confirmed completed tasks are honestly marked and remaining harnessobs/verifier/PR/CI work remains open. |
| Focused PR-level follow-up | `zai/glm-5.1` | APPROVE | Low test-adequacy note accepted and fixed: harnessobs family test now asserts the exact family set. |
| Continued CRAP reduction follow-up | `openrouter/deepseek/deepseek-v4-pro` | APPROVE | Prior findings fixed: trace canonical JSON tests added, unsafe raw hard-error and source-unavailable paths covered, observed run contents asserted, and extra normalization write removed. |

Unusable attempts:

- `openrouter/deepseek/deepseek-v4-pro` evidence review returned empty output;
  not counted as evidence.
- First `minimax/MiniMax-M2.7` requirements review did not receive usable diff
  context; not counted as evidence.
