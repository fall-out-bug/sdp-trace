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
| Complexity over 15 | assessed_gap | Existing production functions remain above `gocyclo -over 15`; `internal/harnessobs.normalizeOpenCodeRawLine`, `internal/harnessobs.LoadSessionProfile`, `internal/harnessobs.CollectSession`, `internal/harnessobs.safeOutDir`, `internal/harnessobs.findUnsafeRawEventAt`, `internal/harnessobs.findUnsafeAt`, `internal/harnessobs.validateEvent`, `internal/trace.writeCanonicalJSON`, `internal/interaction.ValidateEvent`, `internal/prreview.Validate`, `internal/prreview.runRole`, `internal/prreview.BuildPacket`, `internal/witness.BuildCustomerPKI`, `internal/witness.validateCIEnvelope`, `internal/forensic.rawReferenceCondition`, and `internal/demo.witnessBindingState` were decomposed below 15. |
| Coverage hardening | pass_partial | MVP-critical zero-coverage packages `contract`, `export`, and `policy` now have focused tests; `demo` (73.8%), `forensic` (84.0%), `harnessobs` (71.1%), `interaction` (67.3%), `prreview` (73.9%), `trace`, `verifier`, and `witness` (71.2%) were improved. |

## Command Evidence

| Command | Result | Notes |
|---|---|---|
| `go run ./cmd/sdp-trace --help` | pass | Used as canonical command-surface comparison source. |
| `go run ./cmd/sdp-trace pr-review --help` | fail_expected | CLI does not support nested `--help`; global help is the current source of command contracts. |
| `go run ./cmd/sdp-trace pr-review packet --help` | fail_expected | CLI reports `unknown flag --help`; docs were compared against global help. |
| `rg -n -- '--context\|--verification\|This example will show\|controlled-pilot ready\|sidecar trust substrate' README.md docs examples` | pass_absent | Command exits 1 because no matches remain. |
| `go test ./... -coverprofile=/tmp/sdp-trace-harnessobs-validate-full.out` | pass | Total coverage: 72.9%. |
| `go test ./tools/crapcheck -cover` | pass | Tool coverage: 50.6%. |
| `golangci-lint run ./...` | pass | No findings after fixes. |
| `go vet ./...` | pass | Modern Go suspicious-construct sweep. |
| `/Users/fall_out_bug/go/bin/staticcheck ./...` | pass | Staticcheck 2026.1 / v0.7.0, no findings. |
| `rg --files -g '*.go' \| xargs /Users/fall_out_bug/go/bin/gopls check` | pass | `gopls` v0.21.1 LSP diagnostics across all Go files. |
| `/Users/fall_out_bug/go/bin/govulncheck ./...` | pass | Rebuilt with Go 1.26.3 as govulncheck v1.3.0; vulnerability DB updated 2026-05-07; no vulnerabilities found. |
| `gocyclo -over 15 internal/interaction/interaction.go` | pass | `internal/interaction.ValidateEvent` decomposition removed production complexity findings in the interaction file. |
| `gocyclo -over 15 internal/witness/profiles.go` | pass | `internal/witness.BuildCustomerPKI` and `internal/witness.validateCIEnvelope` were decomposed below 15. |
| `gocyclo -over 15 internal/forensic/forensic.go` | pass | `internal/forensic.rawReferenceCondition` was decomposed below 15. |
| `gocyclo -over 15 internal/prreview/prreview.go` | pass | `internal/prreview.Validate`, `internal/prreview.runRole`, and `internal/prreview.BuildPacket` were decomposed below 15. |
| `gocyclo -over 15 internal/harnessobs/harnessobs.go` | fail_assessed_gap | `internal/harnessobs.safeOutDir`, `internal/harnessobs.findUnsafeRawEventAt`, `internal/harnessobs.findUnsafeAt`, and `internal/harnessobs.validateEvent` were decomposed below 15; remaining harness observation findings are `shellFields` and `normalizeRawEvents`. |
| `gocyclo -over 15 internal/demo/demo.go` | pass | `internal/demo.witnessBindingState` was decomposed below 15; no production demo function exceeds 15. |
| `gocyclo -over 15 .` | fail_assessed_gap | Existing production and test functions exceed 15. |
| `gocognit -over 20 .` | fail_assessed_gap | Existing production and test functions exceed 20. |
| `go run ./tools/crapcheck -cover-func /tmp/sdp-trace-harnessobs-validate-full-func.txt -gocyclo /tmp/sdp-trace-harnessobs-validate-full-gocyclo.txt -threshold 5` | fail_assessed_gap | 389 functions exceed strict CRAP threshold 5; `internal/harnessobs.validateEvent` was removed from the over-15 offender list, but the repo-wide strict target remains far open. |

## Coverage Delta

Baseline from intake:

| Package | Baseline | Current |
|---|---:|---:|
| `internal/contract` | 0.0% | 68.1% |
| `internal/export` | 0.0% | 83.3% |
| `internal/policy` | 0.0% | 71.7% |
| `internal/trace` | 2.9% | 61.1% |
| `internal/posture` | 72.4% | 86.3% |
| `internal/harnessobs` | 42.7% | 71.1% |
| `internal/verifier` | 51.1% | 71.0% |
| total | 64.0% | 72.9% |

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
| `cmd/sdp-trace.run` | 29 | 91.2% | 29.57 | assessed_gap |
| `internal/posture.Build` | 19 | 85.3% | 20.15 | assessed_gap |
| `internal/adaptercapture.runBindingCondition` | 19 | 50.0% | 64.12 | assessed_gap |
| `internal/recorder.Run` | 18 | 53.3% | 51.00 | assessed_gap |

## Ratchet

Immediate ratchet now enforced or measurable:

- `golangci-lint run ./...` must stay green locally and in CI.
- `tools/crapcheck` provides reproducible CRAP scoring from `go tool cover -func`
  and `gocyclo`.
- New or materially changed production functions in MVP-critical paths should
  target `CRAP < 5`; exceptions must be recorded here as `assessed_gap`.

Next decomposition candidates before stronger MVP-readiness claim:

1. `cmd/sdp-trace.run`
2. `internal/posture.Build`
3. `internal/adaptercapture.runBindingCondition`
4. `cmd/sdp-trace.(*flagSet).parse`
5. `internal/recorder.Run`
6. `internal/ciartifact.evaluateFamily`

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
| Interaction validator refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Minor source-validation coverage note accepted and fixed; CRAP count explanation accepted and added. Minor LLM ordering note rejected as false positive: LLM refs remain in `validateEventRefs` after catalog/source/content/timing validation. |
| PR-review validator refactor code subagent | `codex-subagent pi` / `zai/glm-5.1` | APPROVE | Minor semantic-alias note accepted and addressed with an explicit `planeCannotVerify` comment. No behavior drift found. |
| PR-review validator refactor evidence subagent | `codex-subagent pi` / `openrouter/qwen/qwen3.6-plus` | ASSESSED_WITH_GAPS | Quantitative evidence artifact concern accepted narrower: fresh local commands were rerun and recorded above, but raw `/tmp` coverage/CRAP outputs are ephemeral and not committed as durable release artifacts in this MVP-hardening slice. Merge approval remains `not_assessed`. |
| Witness profile refactor code subagent | `codex-subagent pi` / `zai/glm-5.1` | FINDINGS_ONLY | Accepted and fixed test gaps for invalid freshness signature, policy digest mismatch, and `customerPKIFail` unknown-field consistency. Quantitative evidence note remains handled by local command evidence above. |
| Witness profile refactor evidence subagent | `codex-subagent pi` / `openrouter/qwen/qwen3.6-plus` | APPROVED_WITH_GAPS | Accepted and fixed CustomerPKI authority/freshness branch gaps: unsupported profile, empty signer, public key mismatch, policy digest mismatch, run mismatch, invalid signature, and source-binding invariant assertion. External/ephemeral evidence artifact concerns remain `not_assessed` for durable release proof. |
| Forensic raw-reference refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | No behavior drift found; ledger honesty confirmed. Minor pass-by-value observation accepted and fixed by passing `*RawReference` into `validateRawReference`. |
| PR-review runner/packet refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | No behavior drift found for packet refs, runner status/reason mapping, raw output retention, OpenCode read-only/dirty/mutation checks, and ledger updates. Minor OpenCode preflight status comment accepted and added. |
| Harness out-dir safety refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Initial major finding accepted: extracted path-safety helpers needed isolated tests. Fixed with direct tests for `pathEscapesWorkingDirectory`, `relativeSymlinkTarget`, `safeExistingOutDir`, `outParentEscapes`, and `ensureOutDirEmptyOrMissing`; focused re-review approved path-safety equivalence. |
| Demo witness-binding refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | No trust-semantics drift found: missing expected scalar/artifact bindings remain `cannot_verify`, mismatches remain `fail`, and empty expectations remain non-blocking. Minor diff-context note required no code change because focused tests and full demo tests passed. |
| Harness unsafe-field traversal refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Safety-equivalence approved for forbidden raw fields, sensitive fields, authenticated URLs, token-like values, unsafe paths, digest exemptions, raw path-like exemptions, retained tool input skip, and unstructured body skip. Minor empty-string passthrough test gap accepted and fixed. |
| Harness event validation refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Focused retry approved validation-semantics equivalence: exact error strings and short-circuit order preserved across identity, refs, content state, and unavailable-fields validation. |

Unusable attempts:

- `openrouter/deepseek/deepseek-v4-pro` evidence review returned empty output;
  not counted as evidence.
- First `minimax/MiniMax-M2.7` requirements review did not receive usable diff
  context; not counted as evidence.
- `openrouter/deepseek/deepseek-v4-pro` harness out-dir focused re-review
  returned a tool-call request under `--no-tools`; not counted as evidence.
- First `openrouter/qwen/qwen3.6-plus` harness event validation review returned
  a tool-call request under `--no-tools`; not counted as evidence.
