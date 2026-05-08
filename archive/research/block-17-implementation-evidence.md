# Block 17 Implementation Evidence

Date: 2026-05-07

Scope: Managed Harness Enforcement Profile implementation slice.

Pull request: https://github.com/fall-out-bug/sdp-trace/pull/7

## Implemented Surface

- `sdp-trace assess --profile managed-harness`
- `sdp-trace assess preview --profile managed-harness`
- `sdp-trace assess explain --assessment-result <file>`
- Block 17 assessment-result schema
- Managed harness policy schema
- Adapter registry schema
- Managed profile domain evaluator with deterministic condition rows

## Evidence

| Evidence | State | Notes |
|---|---|---|
| Spec Socratic review | pass | CTO/product, platform harness owner, and CISO/forensics review found critical/major issues; dispositions are recorded in `archive/research/block-17-spec-review-disposition.md`. |
| Spec re-review | pass | ZAI/GLM-5.1 and MiniMax-M2.7 returned `NO_REMAINING_CRITICAL_OR_MAJOR` after spec fixes. |
| Domain evaluator tests | pass | `rtk go test ./internal/managed` covers managed pass, unmanaged run, late enrollment, post-hoc policy/registry, unauthorized adapter, missing capability, missing harness/tool/file/test telemetry, missing witness, adapter disconnect, agent-reported test evidence, policy-authorized suppression, run-local suppression, empty artifact binding, stale witness, witness mismatch, and override behavior. |
| CLI tests | pass | `rtk go test ./cmd/sdp-trace` covers required managed inputs, valid managed assessment, explain output, post-hoc policy plus witness mismatch, preview no-write behavior, and sensitive marker non-leakage. |
| Full Go test suite | pass | `rtk go test ./...` passed with 119 tests across 12 packages after implementation. |
| Schema and fixture syntax | pass | `rtk jq empty schema/*.json examples/block17-managed-harness/*.json` passed. |
| Diff whitespace | pass | `rtk git diff --check` passed. |
| Deferred-work marker scan | pass | The changed Block 17 code, docs, schema, and fixtures contain no deferred-work markers. |
| GitHub checks | not_assessed | `gh pr checks 7` reported no checks on `codex/block-17-managed-harness-enforcement-profile`; absent checks are recorded as `not_assessed`, not green. |

## Task Coverage Mapping

| Task | Evidence state | Evidence |
|---|---|---|
| T136 | implemented | Block 17 spec and implementation plan exist and were revised after Socratic review. |
| T137 | implemented | `schema/assessment-result.schema.json` defines Block 17 managed assessment output without adding `gate --profile managed-harness`. |
| T138 | implemented | `schema/managed-harness-policy.schema.json` and `schema/adapter-registry.schema.json` define managed policy, adapter authority, capabilities, suppression, provenance, and witness binding shape. |
| T139 | implemented | `internal/managed.Evaluate` and `sdp-trace assess --profile managed-harness` require managed policy, adapter registry, selected run, witness, and pre-run enrollment. |
| T140 | implemented | Go tests prove unmanaged runs, late enrollment, post-hoc policy/registry, self-claimed or unauthorized adapter evidence, and missing managed boundary cannot pass. |
| T141 | implemented | Domain tests cover capability references, capability IDs, and event coverage; missing test telemetry is `missing_telemetry` while agent-reported-only test evidence fails managed profile. |
| T142 | implemented | Run-local or unauthorised suppression fails; policy-authorized suppression can satisfy the selected managed profile only when the pre-run policy rule allows it. |
| T143 | implemented | Managed witness binding checks run id, run nonce, source commit, policy digest, registry digest, enrollment event digest, launch event digest, chain head, event count, freshness, and artifact digests. |
| T144 | implemented | CLI supports `assess --profile managed-harness`, deterministic exit behavior, required input usage errors, and preview rendering. |
| T145 | implemented | CLI supports `assess explain` and `assess preview`; Block 14/16 `gate explain` compatibility remains untouched. |
| T146 | implemented | CLI test asserts managed assessment output does not leak secret-like markers. |
| T147 | implemented | Committed fixtures now cover unmanaged run, late enrollment, post-hoc policy/registry, unauthorized adapter, adapter disconnect, missing capability, missing harness/tool/file/test telemetry, agent-reported test evidence, policy-authorized suppression, suppression without policy, witness missing, stale witness, witness mismatch, override present, override upgrade attempt, and valid managed profile. `TestBlock17CommittedFixturesHaveManagedAssessmentShape` enforces the fixture matrix. |
| T148 | implemented | Implementation code review found two valid Block 17 issues; both fixed. Later trace/evidence review identified capability-reference validation as a useful gap; it was fixed with regression coverage. PR-level code, trace/evidence, and requirements-vs-implementation review returned no remaining critical or major findings. ZAI/GLM-5.1 and Kimi replacement review runs hung and were stopped; they are not counted as clean review evidence. GitHub checks are `not_assessed` because none were reported. |

## Implementation Review Disposition

| Review plane | Reviewer model | State | Disposition |
|---|---|---|---|
| Code review | MiniMax-M2.7 | fixed | Found empty artifact binding returned `fail` instead of `cannot_verify`, and managed preview did not parse directory `run.json`. Both were fixed with regression tests. Other findings were protected-gate legacy concerns, optional-witness requests outside the accepted Block 17 spec, or minor hardening. |
| Narrow re-review | MiniMax-M2.7 | no_remaining_critical_or_major | Re-reviewed empty artifact binding, managed preview parsing, and assess/gate separation after fixes. |
| Trace/evidence review | ZAI/GLM-5.1 | not_assessed | Initial and replacement GLM runs hung and were stopped. |
| Replacement review | Kimi K2P6 low reasoning | not_assessed | Replacement Kimi runs hung and were stopped. |
| Requirements-vs-implementation review | Qwen3 Coder Plus | fixed | Initial bounded review raised a valid capability-reference gap and several invalid or internally contradictory findings. Capability-reference validation was fixed with regression coverage. |
| Requirements-vs-implementation re-review | Qwen3 Coder Plus | no_remaining_critical_or_major | Re-review after suppression, missing-test-telemetry, fixture-matrix, and capability-reference fixes returned `NO_REMAINING_CRITICAL_OR_MAJOR`; residual risks were minor diagnostics and future hardening. |
| PR-level code review | MiniMax-M2.7 | no_remaining_critical_or_major | PR #7 code/correctness review found no critical or major findings. Minor notes were diagnostics and maintainability only; one empty-capability comment was internally inconsistent with the current code path. |
| PR-level trace/evidence review | DeepSeek Chat | no_remaining_critical_or_major | PR #7 trace/evidence review returned `NO_REMAINING_CRITICAL_OR_MAJOR`; residual risks were documentation clarity only. |
| PR-level requirements-vs-implementation review | Qwen3 Coder Plus | no_remaining_critical_or_major | PR #7 requirements review returned no findings and `NO_REMAINING_CRITICAL_OR_MAJOR`. |

## Trust Boundaries

- Managed harness pass is not a native merge, release, readiness, degradation,
  override approval, or risk-acceptance decision.
- Managed mode is opt-in and does not replace observation mode for unmanaged
  harnesses.
- Adapter identity verification does not prove adapter implementation honesty.
- Managed witness binding is not external audit proof.
- `sdp-trace gate --profile managed-harness` is intentionally not implemented;
  Block 17 uses `sdp-trace assess --profile managed-harness`.

## Remaining Work Before Closure

- Merge PR #7, verify the merge commit on `origin/main`, and run the post-merge
  verifier pass before removing the feature branch.
