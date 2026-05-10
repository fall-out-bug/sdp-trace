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
| Complexity over 15 | assessed_gap | Existing production functions remain above `gocyclo -over 15`; `cmd/sdp-trace.gateExitCode`, `cmd/sdp-trace.runAssessExplain`, `cmd/sdp-trace.witnessMatchesProtectedInput`, `cmd/sdp-trace.runGateExplain`, `cmd/sdp-trace.runValidateFixtures`, `cmd/sdp-trace.runPRReviewCheck`, `cmd/sdp-trace.runWitness`, `cmd/sdp-trace.(*flagSet).parse`, `cmd/sdp-trace.runQueryPack`, `cmd/sdp-trace.runProtectedGate`, `cmd/sdp-trace.runAssess`, `internal/adaptercapture.runBindingCondition`, `internal/adaptercapture.overclaimCondition`, `internal/authority.evaluateAction`, `internal/authority.validateEnvelope`, `internal/ciartifact.evaluateFamily`, `internal/ciartifact.safeIdentityToken`, `internal/harnessobs.Observe`, `internal/harnessobs.normalizeOpenCodeRawLine`, `internal/harnessobs.LoadSessionProfile`, `internal/harnessobs.CollectSession`, `internal/harnessobs.safeOutDir`, `internal/harnessobs.findUnsafeRawEventAt`, `internal/harnessobs.findUnsafeAt`, `internal/harnessobs.validateEvent`, `internal/harnessobs.normalizeRawEvents`, `internal/harnessobs.shellFields`, `internal/harnessobs.safeParentDir`, `internal/harnessobs.Validate`, `internal/harnessobs.safeExistingDir`, `internal/harnessobs.safeExistingFile`, `internal/harnessobs.validateProfile`, `internal/managed.witnessCondition`, `internal/managed.capabilityCondition`, `internal/posture.Build`, `internal/posture.validateMetricRowShape`, `internal/trace.writeCanonicalJSON`, `internal/interaction.ImportTranscript`, `internal/interaction.ValidateEvent`, `internal/interaction.SummarizeTrace`, `internal/contract.(ExpectedEvidenceContract).Validate`, `internal/prreview.Validate`, `internal/prreview.runRole`, `internal/prreview.BuildPacket`, `internal/prreview.validateProfile`, `internal/prreview.citationResolvable`, `internal/recorder.Run`, `internal/recorder.runCommand`, `internal/repoobserver.writeTarget`, `internal/verifier.VerifyRun`, `internal/verifier.verifyChain`, `internal/witness.BuildCustomerPKI`, `internal/witness.validateCIEnvelope`, `internal/witness.validateCIEnvelopeStates`, `internal/witness.validateCustomerPKIAuthority`, `internal/witness.BuildCIEnvelopeProfile`, `internal/query.safeToken`, `internal/forensic.rawReferenceCondition`, `internal/forensic.policyCondition`, `internal/forensic.prewriteCondition`, `internal/demo.EvaluateGate`, `internal/demo.witnessBindingState`, and `internal/telemetry.unsafeValue` were decomposed below 15. |
| Coverage hardening | pass_partial | MVP-critical zero-coverage packages `contract`, `export`, and `policy` now have focused tests; `cmd/sdp-trace` (71.7%), `authority` (89.0%), `adaptercapture` (84.3%), `ciartifact` (92.9%), `contract` (89.1%), `demo` (74.4%), `forensic` (90.4%), `harnessobs` (79.5%), `interaction` (69.7%), `managed` (88.9%), `posture` (88.7%), `prreview` (77.7%), `query` (92.0%), `repoobserver` (71.7%), `trace` (62.2%), `verifier` (83.4%), and `witness` (77.2%) were improved. |

## Command Evidence

| Command | Result | Notes |
|---|---|---|
| `go run ./cmd/sdp-trace --help` | pass | Used as canonical command-surface comparison source. |
| `go run ./cmd/sdp-trace pr-review --help` | fail_expected | CLI does not support nested `--help`; global help is the current source of command contracts. |
| `go run ./cmd/sdp-trace pr-review packet --help` | fail_expected | CLI reports `unknown flag --help`; docs were compared against global help. |
| `rg -n -- '--context\|--verification\|This example will show\|controlled-pilot ready\|sidecar trust substrate' README.md docs examples` | pass_absent | Command exits 1 because no matches remain. |
| `go test ./... -coverprofile=/tmp/sdp-trace-batch8-full.out` | pass | Total coverage: 77.4%. |
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
| `gocyclo -over 15 .` | fail_assessed_gap | Existing production and test functions exceed 15. |
| `gocognit -over 20 .` | fail_assessed_gap | Existing production and test functions exceed 20. |
| `go run ./tools/crapcheck -cover-func /tmp/sdp-trace-batch8-full-func.txt -gocyclo /tmp/sdp-trace-batch8-full-gocyclo.txt -threshold 5` | fail_assessed_gap | 341 functions exceed strict CRAP threshold 5; `internal/witness.FetchGitHubOIDCToken`, `internal/trace.ValidateRunDirectory`, `internal/repoobserver.writeInstallFiles`, `internal/repoobserver.updateGitignore`, `internal/query.mapSourceState`, and their extracted helpers are now below threshold, but the repo-wide strict target remains open. |

## Coverage Delta

Baseline from intake:

| Package | Baseline | Current |
|---|---:|---:|
| `internal/contract` | 0.0% | 68.1% |
| `internal/export` | 0.0% | 83.3% |
| `internal/policy` | 0.0% | 71.7% |
| `internal/trace` | 2.9% | 62.2% |
| `internal/posture` | 72.4% | 87.8% |
| `internal/harnessobs` | 42.7% | 79.5% |
| `internal/verifier` | 51.1% | 83.4% |
| total | 64.0% | 77.4% |

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
| `internal/witness.loadCustomerPublicKey` | 10 | 52.0% | 21.06 | assessed_gap |
| `internal/prreview.safeID` | 10 | 80.0% | 10.80 | assessed_gap |

## Ratchet

Immediate ratchet now enforced or measurable:

- `golangci-lint run ./...` must stay green locally and in CI.
- `tools/crapcheck` provides reproducible CRAP scoring from `go tool cover -func`
  and `gocyclo`.
- New or materially changed production functions in MVP-critical paths should
  target `CRAP < 5`; exceptions must be recorded here as `assessed_gap`.

Next decomposition candidates before stronger MVP-readiness claim:

1. `cmd/sdp-trace.run`
2. `internal/witness.loadCustomerPublicKey`
3. `internal/prreview.safeID`
4. `internal/prreview.RunReview`
5. `internal/posture.unsafeOutput`
6. `internal/posture.verifyDigestManifest`
7. `internal/harnessobs.findNumberByKeyIn`

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

Unusable attempts:

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
