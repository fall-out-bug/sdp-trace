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
| CRAP < 5 | pass_local | `go run ./tools/crapcheck -cover-func /tmp/sdp-trace-pr37-cover-func.txt -gocyclo /tmp/sdp-trace-pr37-gocyclo-prod.txt -threshold 5` exited 0 after repository-wide production decomposition. |
| Complexity over 15 | pass_local | `gocyclo -over 15 $(find cmd internal -name '*.go' ! -name '*_test.go' -print)` and `gocognit -over 15 $(find cmd internal -name '*.go' ! -name '*_test.go' -print)` exited 0 for production code. |
| Coverage hardening | pass_local | MVP-critical zero/low-coverage paths have focused tests; final local package coverage includes `cmd/sdp-trace` (80.3%), `authority` (95.5%), `adaptercapture` (89.2%), `checkpoint` (82.5%), `ciartifact` (96.2%), `contract` (89.1%), `demo` (92.2%), `forensic` (93.2%), `harnessobs` (89.1%), `interaction` (82.9%), `managed` (91.4%), `policy` (86.8%), `posture` (93.2%), `prreview` (86.5%), `query` (95.2%), `repoobserver` (88.6%), `trace` (83.8%), `verifier` (85.7%), and `witness` (85.4%). |

## Command Evidence

| Command | Result | Notes |
|---|---|---|
| `go run ./cmd/sdp-trace --help` | pass | Used as canonical command-surface comparison source. |
| `go run ./cmd/sdp-trace pr-review --help` | fail_expected | CLI does not support nested `--help`; global help is the current source of command contracts. |
| `go run ./cmd/sdp-trace pr-review packet --help` | fail_expected | CLI reports `unknown flag --help`; docs were compared against global help. |
| `rg -n -- '--context\|--verification\|This example will show\|controlled-pilot ready\|sidecar trust substrate' README.md docs examples` | pass_absent | Command exits 1 because no matches remain. |
| `go test ./... -coverprofile=/tmp/sdp-trace-pr37-cover.out` | pass | Final local package coverage was recorded before CRAP scoring; `cmd/sdp-trace` reached 80.3% and all packages passed. |
| `go test ./tools/crapcheck -cover` | pass | Tool coverage: 50.6%. |
| `golangci-lint run ./...` | pass | No findings after fixes. |
| `go vet ./...` | pass | Modern Go suspicious-construct sweep. |
| `/Users/fall_out_bug/go/bin/staticcheck ./...` | pass | Staticcheck 2026.1 / v0.7.0, no findings. |
| `rg --files -g '*.go' \| xargs /Users/fall_out_bug/go/bin/gopls check` | pass | `gopls` v0.21.1 LSP diagnostics across all Go files. |
| `/Users/fall_out_bug/go/bin/govulncheck ./...` | pass | Rebuilt with Go 1.26.3 as govulncheck v1.3.0; vulnerability DB updated 2026-05-07; no vulnerabilities found. |
| `gocyclo -over 15 internal/interaction/interaction.go` | pass | `internal/interaction.ValidateEvent` decomposition removed production complexity findings in the interaction file. |
| `gocyclo -over 15 internal/witness/profiles.go` | pass | `internal/witness.BuildCustomerPKI` and `internal/witness.validateCIEnvelope` were decomposed below 15. |
| `gocyclo -over 15 internal/forensic/forensic.go` | pass | `internal/forensic.rawReferenceCondition` was decomposed below 15. |
| `gocyclo -over 14 internal/forensic/forensic.go` | pass | `internal/forensic.policyCondition` was decomposed below 15; no production forensic function exceeds 14. |
| `gocyclo -over 15 internal/prreview/prreview.go` | pass | `internal/prreview.Validate`, `internal/prreview.runRole`, and `internal/prreview.BuildPacket` were decomposed below 15. |
| `gocyclo -over 15 internal/harnessobs/harnessobs.go` | pass | All production harness observation functions are now below 15. |
| `gocyclo -over 14 internal/harnessobs/harnessobs.go` | pass | `internal/harnessobs.Observe` was decomposed below 15; no production harness observation function exceeds 14. |
| `gocyclo -over 15 internal/demo/demo.go` | pass | `internal/demo.witnessBindingState` was decomposed below 15; no production demo function exceeds 15. |
| `gocyclo -over 14 internal/demo/demo.go` | pass | `internal/demo.EvaluateGate` was decomposed below 15; no production demo function exceeds 14. |
| `gocyclo -over 15 internal/posture/posture.go` | pass | `internal/posture.Build` was decomposed below 15; no production posture function exceeds 15. |
| `gocyclo -over 15 internal/adaptercapture/adaptercapture.go` | pass | `internal/adaptercapture.runBindingCondition` and `internal/adaptercapture.overclaimCondition` were decomposed below 15; no production adaptercapture function exceeds 15. |
| `gocyclo -over 15 internal/managed/managed.go` | pass | `internal/managed.witnessCondition` was decomposed below 15; no production managed-harness function exceeds 15. |
| `gocyclo -over 15 internal/authority/authority.go` | pass | `internal/authority.evaluateAction` and `internal/authority.validateEnvelope` were decomposed below 15; no production authority function exceeds 15. |
| `/Users/fall_out_bug/go/bin/gocognit -over 20 internal/authority/authority.go` | pass | No production authority function exceeds 20 cognitive complexity. |
| `gocyclo -over 14 internal/verifier/verify.go` | pass | `internal/verifier.VerifyRun` was decomposed below 15; no production verifier function exceeds 14. |
| `/Users/fall_out_bug/go/bin/gocognit -over 20 internal/verifier/verify.go` | pass | `internal/verifier.VerifyRun` no longer exceeds 20 cognitive complexity. |
| `gocyclo -over 14 internal/repoobserver/repoobserver.go` | pass | `internal/repoobserver.writeTarget` was decomposed below 15; no production repoobserver function exceeds 14. |
| `gocyclo -over 14 internal/posture/posture.go` | pass | `internal/posture.validateMetricRowShape` was decomposed below 15; no production posture function exceeds 14. |
| `gocyclo -over 14 internal/interaction/interaction.go` | pass | `internal/interaction.ImportTranscript` was decomposed below 15; no production interaction function exceeds 14. |
| `gocyclo -over 14 cmd/sdp-trace/main.go` | assessed_gap | `cmd/sdp-trace.gateExitCode` was decomposed below 15; other CLI functions remain above 14. |
| `gocyclo -over 14 cmd/sdp-trace/main.go` | assessed_gap | `cmd/sdp-trace.runAssessExplain` was decomposed below 15; other CLI functions remain above 14. |
| `gocyclo -over 14 cmd/sdp-trace/main.go` | assessed_gap | `cmd/sdp-trace.witnessMatchesProtectedInput` was decomposed below 15; other CLI functions remain above 14. |
| `gocyclo -over 14 cmd/sdp-trace/main.go` | assessed_gap | `cmd/sdp-trace.runGateExplain` was decomposed below 15; other CLI functions remain above 14. |
| `gocyclo -over 14 cmd/sdp-trace/main.go` | assessed_gap | `cmd/sdp-trace.runValidateFixtures` was decomposed below 15; other CLI functions remain above 14. |
| `gocyclo -over 14 cmd/sdp-trace/main.go` | assessed_gap | `cmd/sdp-trace.runPRReviewCheck` was decomposed below 15; other CLI functions remain above 14. |
| `gocyclo -over 14 internal/recorder/recorder.go` | pass | `internal/recorder.Run` was decomposed below 15; no production recorder function exceeds 14. |
| `gocyclo -over 14 internal/ciartifact/ciartifact.go` | pass | `internal/ciartifact.evaluateFamily` was decomposed below 15; no production ciartifact function exceeds 14. |
| `gocyclo -over 15 $(find cmd internal -name '*.go' ! -name '*_test.go' -print)` | pass | No production functions in `cmd` or `internal` exceed cyclomatic complexity 15. |
| `gocognit -over 15 $(find cmd internal -name '*.go' ! -name '*_test.go' -print)` | pass | No production functions in `cmd` or `internal` exceed cognitive complexity 15. |
| `go run ./tools/crapcheck -cover-func /tmp/sdp-trace-pr37-cover-func.txt -gocyclo /tmp/sdp-trace-pr37-gocyclo-prod.txt -threshold 5` | pass | Repository-wide production CRAP threshold 5 passed locally; no functions exceeded the threshold. |

## Coverage Delta

Baseline from intake:

| Package | Baseline | Current |
|---|---:|---:|
| `internal/contract` | 0.0% | 68.1% |
| `internal/export` | 0.0% | 83.3% |
| `internal/policy` | 0.0% | 71.7% |
| `internal/trace` | 2.9% | 83.8% |
| `internal/posture` | 72.4% | 93.2% |
| `internal/harnessobs` | 42.7% | 89.1% |
| `internal/verifier` | 51.1% | 85.7% |
| total | 64.0% | package-level final profile recorded above |

## CRAP Baseline Summary

Formula:

```text
CRAP = complexity^2 * (1 - coverage)^3 + complexity
```

Strict production `CRAP < 5` is now satisfied locally for `cmd` and `internal`
production code by `tools/crapcheck` at threshold 5. The gate is still local
evidence until PR #37 is pushed and fresh GitHub CI completes on the final head.

## Ratchet

Immediate ratchet now enforced or measurable:

- `golangci-lint run ./...` must stay green locally and in CI.
- `tools/crapcheck` provides reproducible CRAP scoring from `go tool cover -func`
  and `gocyclo`.
- New or materially changed production functions in MVP-critical paths should
  target `CRAP < 5`; exceptions must be recorded here as `assessed_gap`.

Next decomposition candidates: none required for the local production
`CRAP < 5`, cyclomatic `< 15`, or cognitive `< 15` gates. Future work should
keep the same gates enforced for changed production paths and move them into CI
if PR-level review accepts the local gate shape.

## External Evidence Boundary

GitHub CI `verify` passed on PR #37 head
`73cb78f30e5edf157ae1e0b1e5bc30d7b0ea95ff`
([run 25640589077/job 75260219377](https://github.com/fall-out-bug/sdp-trace/actions/runs/25640589077/job/75260219377)).
PR-level review planes are recorded below. Merge approval is user-delegated in
the 2026-05-10 task thread.

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
| Harness raw normalization refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Raw normalization equivalence approved: unsafe input rejection, malformed JSONL line numbers, source digest computation, output order, and file close behavior preserved. |
| Harness shell parser refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Controlled shell parser equivalence approved: quoting, backslash handling, newline continuation, separator flush, trailing escape, and safe model gate behavior preserved. |
| Posture Build refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Trust-state/refusal equivalence approved: unsafe labels, malformed paths, stale inputs, digest failures, malformed query packs, malformed signal manifests, trusted grouping, and export fields preserve prior behavior. |
| Adapter run-binding refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Run-binding equivalence approved: run id/nonce mismatch remains `fail`, and missing binding, late same-chain events, missing same-chain hashes, unbound bundles, late bundles, and unsupported binding modes remain `cannot_verify` with unchanged reason codes. |
| Adapter overclaim refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Capture-depth overclaim equivalence approved: reconstructable claims without sufficient evidence and without cap annotation still fail with `capture_depth_overclaimed`; positive cap/sufficient-evidence paths were also covered after review. |
| Managed witness refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Witness-binding equivalence approved. Reviewer noted a possible missing `WitnessID` guard, rejected as false positive after checking the guard remains at function entry and focused tests pass. |
| Authority action refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Focused retry approved behavior equivalence: state/reason precedence, `MatchedRuleRef`, `approval_evidence_missing`, and attribution initialization remain preserved after extraction. |
| Authority envelope refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Reason-code precedence and overlapping-rule equivalence approved; table-driven tests cover each extracted validation path. |
| Verifier run refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Final review approved behavior equivalence for manifest validation, event loading, chain verification, contract digest checks, missing evidence generation, error propagation, audit detail mapping, and corrected shadowing/unused-parameter concerns. |
| Repo observer write-target refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Safety equivalence approved for unsafe target rejection, idempotent existing-file handling, executable chmod, force backup/write, and new-file writes. |
| Posture metric-row shape refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | All 14 malformed metric-row predicates were mapped to the extracted helpers with the same error string; reviewer noted optional sibling cases but no blocking drift. |
| Interaction transcript import refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Final review approved exact call-order and error-semantics equivalence for option validation, JSONL read/empty handling, per-event task/source checks, source mutation before `ValidateEvent`, ordering validation, and `WriteTrace`. |
| Gate exit-code refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Initial review found unknown protected-gate fall-through drift; fixed by map-miss fall-through and explicit tests for component fail/pass. Focused re-review approved fail/missing precedence and cannot-verify handling. |
| Assess explain dispatch refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Generic table dispatch approved as behavior-equivalent: usage errors, envelope read, typed schema reads, same explain functions, unsupported-schema error, and all five schema handlers preserve observable behavior. |
| Protected witness match refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Source/trust/artifact matching equivalence approved. Reviewer low notes accepted with explicit empty/empty artifact coverage and existence-safe path/hash lookup. |
| Recorder run pipeline refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Final review approved execution-order and payload equivalence for output dir setup, contract validation, manifest write, env assignment, event sequence, command timing, closure/finalize, and returned result. |
| Recorder command runner refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Focused retry approved behavior equivalence for command start/wait, stdout/stderr tee hashing, non-zero exit preservation, signal/context semantics, and ledger honesty. Minor signal-driven exit coverage note remains advisory because existing semantics were preserved and focused tests cover output hashes, start failure, and non-zero exit. |
| CI artifact family refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Final review approved access-to-producer-to-binding precedence, default/unknown fallthrough, byte-identical reason/action text, CRAP reduction, coverage bump, and ledger honesty. Prior minor readability notes were accepted before the final review. |
| Gate explain refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Review approved flag parsing exit codes/messages, read/schema handling, block14 protected-absent line, block16 protected fields, rendering order, raw command non-disclosure, CRAP/coverage honesty, and ledger updates. Minor read-helper unit-test note remains advisory because existing end-to-end tests cover happy path, unsupported schema, and secret non-disclosure. |
| Harness observe refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Review approved required-option errors, safe path resolution and error wrapping, profile/source loading order, default `Now` handling, event and `run.json` writes, no partial write on early failures, source digest/event refs preservation, CRAP/coverage honesty, and ledger updates. |
| Forensic policy condition refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Review approved original precedence order, run/event digest mismatch handling, event authority semantics, reason/action preservation, focused precedence tests, CRAP/coverage evidence, and ledger updates. |
| Demo gate evaluation refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Review approved default gate fields, observed evidence and override preservation, local gate aggregation, required-run/evidence effects, gateConditions timing, default reason insertion, sorting, CRAP/coverage evidence, and ledger updates. |
| Validate fixtures refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Focused review returned `VERDICT APPROVE` for fixture-root selection, verifier artifact writing, verify-error reporting, expectation handling, unexpected fail/cannot_verify behavior, final exit aggregation, and ledger honesty. |
| PR review check refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Review approved parse/usage behavior, packet requirements and options, profile/work-dir preparation, preview JSON and no-artifact behavior, RunReview options, artifact write order, validation exit mapping, summary output, focused tests, CRAP/coverage evidence, and ledger updates. |
| Telemetry unsafe-value refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Review approved fail-closed marker equivalence for URL, credential, path, and contact markers, case-normalization behavior, no partial Prometheus output on unsafe labels, expanded tests, CRAP evidence, and ledger updates. |
| PR review profile validator refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Review approved schema/profile/plane/role validation order, runner validation, required-plane coverage, error-string preservation, focused rejection/acceptance tests, CRAP/coverage evidence, and ledger updates. |
| PR review citation resolver refactor follow-up | `openrouter/qwen/qwen3.6-plus` | APPROVE | Review approved empty-citation rejection, diff/diff-alias semantics, context and verification ref semantics, unknown-ref and digest-only fallback preservation, resolver order, characterization tests, CRAP/coverage evidence, and ledger updates. |
| Managed/forensic/harness/verifier/witness batch | `gpt-5.3-codex-spark` implementation subagents | IMPLEMENTED | Spark workers decomposed `internal/forensic.prewriteCondition`, `internal/harnessobs.safeParentDir`, `internal/verifier.verifyChain`, and `cmd/sdp-trace.runWitness`; parent integrated `internal/managed.capabilityCondition`, fixed review/test/CRAP gaps, reran focused package tests, and recorded batch CRAP/coverage evidence. Full PR-level review remains pending for the next larger integration slice. |
| CLI parser and harness validation batch | `gpt-5.3-codex-spark` implementation subagents | IMPLEMENTED | Spark workers decomposed `cmd/sdp-trace.(*flagSet).parse`, `internal/harnessobs.Validate`, `internal/harnessobs.safeExistingDir`, and `internal/harnessobs.safeExistingFile`; parent integrated overlapping harness changes, reran focused cmd/harness tests, and recorded batch CRAP/coverage evidence. Full PR-level review remains pending for the next larger integration slice. |
| Interaction/contract/CI artifact batch | `gpt-5.3-codex-spark` implementation subagents | IMPLEMENTED | Spark workers decomposed `internal/interaction.SummarizeTrace`, `internal/contract.(ExpectedEvidenceContract).Validate`, and `internal/ciartifact.safeIdentityToken`; parent reran focused package tests and recorded batch CRAP/coverage evidence. Full PR-level review remains pending for the next larger integration slice. |
| Harness profile/witness/query batch | `gpt-5.3-codex-spark` implementation subagents | IMPLEMENTED | Spark workers decomposed `internal/harnessobs.validateProfile`, `internal/witness.validateCIEnvelopeStates`, `internal/witness.validateCustomerPKIAuthority`, `internal/witness.BuildCIEnvelopeProfile`, and `internal/query.safeToken`; parent fixed a `strings.ContainsRune` staticcheck finding, reran focused package tests, and recorded batch CRAP/coverage evidence. Full PR-level review remains pending for the next larger integration slice. |
| Setup/sanitize/authority batch | `gpt-5.3-codex-spark` implementation subagents | IMPLEMENTED | Spark workers decomposed `internal/harnessobs.SetupSession`, `internal/ciartifact.sanitizeRun`, and `internal/authority.aggregateState`; parent reran focused package tests and recorded batch CRAP/coverage evidence. Full PR-level review remains pending for the next larger integration slice. |
| CLI query/protected/assess batch | `gpt-5.3-codex-spark` implementation subagents | IMPLEMENTED | Spark workers decomposed `cmd/sdp-trace.runQueryPack`, `cmd/sdp-trace.runProtectedGate`, and `cmd/sdp-trace.runAssess`; parent resolved overlapping `main.go` edits, reran full cmd package tests, and recorded batch CRAP/coverage evidence. Full PR-level review remains pending for the next larger integration slice. |
| Witness/trace/repoobserver/query batch | `gpt-5.3-codex-spark` implementation subagents | IMPLEMENTED | Spark workers decomposed `internal/witness.FetchGitHubOIDCToken`, `internal/trace.ValidateRunDirectory`, `internal/repoobserver.writeInstallFiles`, `internal/repoobserver.updateGitignore`, and `internal/query.mapSourceState`; parent fixed a duplicate hooksPath summary integration bug, reran focused package tests, and recorded batch CRAP/coverage evidence. Full PR-level review remains pending for the next larger integration slice. |
| Witness/prreview/posture/harness batch | `gpt-5.3-codex-spark` implementation subagents | IMPLEMENTED | Spark workers decomposed `internal/witness.loadCustomerPublicKey`, `internal/prreview.safeID`, `internal/posture.unsafeOutput`, `internal/harnessobs.findStringByKeyIn`, and `internal/harnessobs.findNumberByKeyIn`; parent reduced `findByKeyInMap` below the strict CRAP target, reran focused package tests, and recorded batch CRAP/coverage evidence. Full PR-level review remains pending for the next larger integration slice. |
| Batch 7-9 code/correctness review | `zai/glm-5.1` | APPROVE | Review found no blocking findings. Minor mutable `sourceStateMap` concern accepted and fixed by localizing the map inside `mapSourceState`; whitespace-trim note recorded as non-blocking safe behavior. |
| Batch 7-9 trace/evidence review | `kimi-coding/k2p6` | APPROVE | Review approved ledger honesty: repo-wide strict CRAP remains `assessed_gap`, command evidence is cited consistently, and merge approval remains `not_assessed`. |
| Batch 7-9 requirements-vs-implementation review | `openrouter/qwen/qwen3.6-plus` | FALSE_POSITIVE | Reviewer reported a `runQueryPack` fall-through regression. Full-file check shows `runQueryPack` returns `runQueryPackBuild(args, stderr)`, and `go test ./cmd/sdp-trace -run TestQueryPack -count=1` passed; finding rejected. |
| PR review/posture/demo/checkpoint batch | `gpt-5.3-codex-spark` implementation subagents | IMPLEMENTED | Spark workers decomposed `internal/prreview.RunReview`, `internal/posture.verifyDigestManifest`, `internal/demo.overrideRequestsFromEvents`, and `internal/checkpoint.applyPolicy`; parent ported the posture patch into the correct worktree after the worker edited the base checkout, reduced checkpoint/demo extracted helpers below strict CRAP target, reran focused package tests, and recorded batch CRAP/coverage evidence. Full PR-level review remains pending for the next larger integration slice. |
| Posture/forensic/demo/adapter batch | `gpt-5.3-codex-spark` implementation subagents | IMPLEMENTED | Spark workers decomposed `internal/posture.validateMovementRow`, `internal/forensic.validateRawReference`, `internal/demo.evaluateRequiredRuns`, and `internal/adaptercapture.redactionMetadataCondition`; parent reran focused package tests and recorded batch CRAP/coverage evidence. Full PR-level review remains pending for the next larger integration slice. |
| Posture/authority/adapter/demo batch | `gpt-5.3-codex-spark` implementation subagents | IMPLEMENTED | Spark workers decomposed `internal/posture.validateExportCollections`, `internal/posture.validateExportRows`, `internal/authority.evidenceRefsReason`, `internal/adaptercapture.contractCondition`, and `internal/demo.DiscoverRunDirs`; parent reduced `evidenceRefReason`, `adapterEventIsMalformed`, and `validateExportCollections` below the strict CRAP target, reran focused package tests, and recorded batch CRAP/coverage evidence. Full PR-level review remains pending for the next larger integration slice. |
| Batch 10-11 code/correctness review | `zai/glm-5.1` | APPROVE_WITH_MINOR_FIXES | Review found no blocking findings. Minor digest mismatch return-value drift and required-run reason precedence drift were accepted and fixed with focused tests. |
| Batch 10-11 trace/evidence review | `openrouter/qwen/qwen3.6-plus` | APPROVE | Review approved ledger honesty and evidence boundaries; non-blocking notes about focused review rows and package-specific coverage proof were recorded but did not change merge approval state. |
| Batch 10-11 requirements-vs-implementation review | `minimax/MiniMax-M2.7` | APPROVE | Review found no blocking requirement, UX/DX, trust, or user-facing regressions; it confirmed strict CRAP remains `assessed_gap`. |
| PR #37 code/correctness review | `openrouter/deepseek/deepseek-v4-pro` | APPROVE_WITH_MINOR_NOTES | Initial code reviewer received an unexpanded packet and was not counted. Replacement review approved with minor notes: CI was pending at review start, now fixed by final-head CI pass; dependency-coupling automation remains a future enhancement, with current internal import evidence recorded. |
| PR #37 trace/evidence review | `openrouter/qwen/qwen3.6-plus` | APPROVE | Initial evidence review correctly rejected missing packet evidence and was not counted as approval. Focused re-review approved after final-head CI, CRAP/coverage/gocyclo/gocognit outputs, ledger state, and internal import evidence were supplied; prior `main.go` size as domain-leak proof was retracted as unsupported. |
| PR #37 requirements-vs-implementation review | `openrouter/deepseek/deepseek-v4-pro` | APPROVE | Initial MiniMax review raised `main.go` size and missing-evidence concerns from the limited packet. Focused re-review approved after final evidence showed production CRAP <5, cyclomatic/cognitive <15, no internal-to-cmd imports, and ledger honesty. |

Unusable attempts:

- First PR #37 code/correctness review received an unexpanded packet literal;
  not counted as evidence.
- First PR #37 requirements-vs-implementation review hallucinated files outside
  this repo (`internal/metrics/recorder.go`); not counted as evidence.
- First PR #37 trace/evidence review only assessed missing packet evidence;
  useful as a packet-quality finding, not counted as final approval.
- `openrouter/deepseek/deepseek-v4-pro` evidence review returned empty output;
  not counted as evidence.
- First `minimax/MiniMax-M2.7` requirements review did not receive usable diff
  context; not counted as evidence.
- `openrouter/deepseek/deepseek-v4-pro` harness out-dir focused re-review
  returned a tool-call request under `--no-tools`; not counted as evidence.
- First `openrouter/qwen/qwen3.6-plus` harness event validation review returned
  a tool-call request under `--no-tools`; not counted as evidence.
- First `openrouter/qwen/qwen3.6-plus` authority action review returned a
  tool-call request under `--no-tools`; not counted as evidence. Focused retry
  produced usable review evidence.
- First `openrouter/qwen/qwen3.6-plus` verifier run review finished after the
  diff changed; not counted as final evidence. Its cleanup concern was handled,
  and a fresh review approved the final diff.
- First `openrouter/qwen/qwen3.6-plus` recorder command-runner review returned
  a tool-call request under `--no-tools`; not counted as evidence. Focused retry
  produced usable review evidence.
