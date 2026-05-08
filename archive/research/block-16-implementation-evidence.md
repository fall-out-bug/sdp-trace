# Block 16 Implementation Evidence

Date: 2026-05-06

Scope: Protected Gate Enforcement Profile implementation.

## Implemented Surface

- `sdp-trace gate --profile protected`
- `sdp-trace gate preview --profile protected`
- `sdp-trace gate explain` read compatibility for Block 14 and Block 16 gate results
- Block 16 gate-result schema shape
- Protected profile domain evaluator with deterministic condition rows

## Evidence

| Evidence | State | Notes |
|---|---|---|
| TDD red for domain protected evaluator | pass | `go test ./internal/demo` initially failed on missing `EvaluateProtectedGate`, `ProtectedGateInput`, and protected condition types. |
| TDD green for domain protected evaluator | pass | `rtk go test ./internal/demo` passed with 23 tests after implementation and review fixes. |
| TDD red for CLI protected profile | pass | `go test ./cmd/sdp-trace` initially failed on unknown `--profile`. |
| TDD green for CLI protected profile | pass | `rtk go test ./cmd/sdp-trace` passed with 46 tests after implementation and review fixes. |
| Full Go test suite | pass | `rtk go test ./...` passed with 96 tests across 11 packages. |
| Schema syntax | pass | `rtk jq empty schema/*.json` passed. |
| Example syntax | pass | `rtk jq empty examples/block16-protected-gate/*.json` passed. |
| Fixture shape and semantics test | pass | `TestBlock16CommittedFixturesHaveRequiredProtectedRows` verifies 13 committed fixtures, Block 16 schema/profile fields, `gate_mode: protected`, all 10 protected condition rows in order, top-level protected-gate consistency, protected trust-scope dependencies, and local-signed signer-authority non-pass semantics. |
| Diff whitespace | pass | `rtk git diff --check` passed. |

## Task Coverage Mapping

| Task | Evidence state | Evidence |
|---|---|---|
| T126 | implemented_not_closed | Block 16 spec exists and was updated from six-plane spec review. |
| T127 | implemented_not_closed | `schema/gate-result.schema.json` has version-separated Block 14 and Block 16 shapes; `schema/README.md` documents compatibility. |
| T128 | implemented_not_closed | `TestProtectedGateRequiresCheckpointPolicyAndWitnessFlags`, `TestProtectedGateMalformedNamedInputIsUsageError`, `TestProtectedGateRejectsLocalSignedCheckpointCLI`, `TestProtectedGateMissingPolicyCannotVerify`, and `TestProtectedGatePassesWithCISignedCheckpointAndBoundWitnessCLI`. |
| T129 | implemented_not_closed | `TestProtectedGateRejectsLocalSignedCheckpointAndKeepsConditionRows` and `TestProtectedGateRejectsLocalSignedCheckpointCLI`. |
| T130 | implemented_not_closed | `TestProtectedGateMapsAbsentAndStaleWitnessFreshness`, `TestProtectedGateRequiresWitnessRunIDBinding`, `TestGateCommandFailsForWitnessArtifactMismatch`, `TestGateCommandCannotVerifyWhenWitnessOmitsRunArtifact`, and `TestProtectedGatePassesWithCISignedCheckpointAndBoundWitnessCLI`. |
| T131 | implemented_not_closed | `TestProtectedGateRequiresCheckpointPolicyAndWitnessFlags`, `TestProtectedGatePreviewRendersAbsentInputsWithoutWriting`, `TestGateExitCodeUsesProtectedGateWhenSelected`, and protected gate exit-code assertions. |
| T132 | implemented_not_closed | `TestGateExplainRendersProtectedFields`, `TestGateExplainUnsupportedArtifactCannotVerify`, and Block 14 explain compatibility tests, including explicit absence of protected fields for Block 14 artifacts. |
| T133 | implemented_not_closed | `TestReportAndGateArtifactsDoNotLeakSecretLikeCommand`, `TestProtectedGatePreviewRendersAbsentInputsWithoutWriting`, `TestProtectedGateRejectsLocalSignedCheckpointCLI`, and existing Block 14 secret-like output tests. |
| T134 | implemented_not_closed | `examples/block16-protected-gate/` has committed fixtures for missing checkpoint, local-development checkpoint, invalid run binding, missing signer policy, signer mismatch, missing CI witness, absent freshness, stale CI witness, CI source mismatch, CI artifact mismatch, malformed override with trust-scope failure, valid CI-authority protected profile, and override-present protected profile. Fixture rows and semantic consistency are checked by `TestBlock16CommittedFixturesHaveRequiredProtectedRows`. |
| T135 | open | Implementation code review, trace review, requirements-vs-implementation review, PR-level code review, PR-level tracing review, and PR-level requirements-vs-implementation review have no remaining critical or major findings. |

These mappings are implementation evidence, not task closure claims. Task
checkbox closure remains blocked until source-bound final closure policy is
satisfied; this evidence document is not a task-closure authority.

## Implementation Review Disposition

| Review plane | Reviewer model | State | Disposition |
|---|---|---|---|
| Code review | MiniMax-M2.7 | no_remaining_critical_or_major | Initial code review found raw command persistence risk and test gaps. Re-review found no critical or major findings; remaining minors are documented as low-risk hardening. |
| Tracing review | ZAI/GLM-5.1 | no_remaining_critical_or_major | Initial trace review found gate-mode/trust-cap/fixture consistency gaps. Re-review verified the fixes and found no critical or major findings. |
| Requirements-vs-implementation review | Kimi K2P6 low reasoning | no_remaining_critical_or_major | Initial requirements review found top-level state mapping, fixture-row, run-id binding, and exit-code gaps. Re-review major hermetic timestamp finding was fixed; remaining points are minor hardening. |

## PR-Level Review Disposition

| Review plane | Reviewer model | State | Disposition |
|---|---|---|---|
| PR code review | MiniMax-M2.7 | no_remaining_critical_or_major | PR-level code review found no critical or major findings. |
| PR tracing review | ZAI/GLM-5.1 | no_remaining_critical_or_major | Found fixture semantic inconsistencies where static condition rows overclaimed `protected_trust_scope_satisfied`. Fixtures were corrected, `TestBlock16CommittedFixturesHaveRequiredProtectedRows` now checks semantic consistency, and PR tracing re-review found no remaining critical or major findings. |
| PR requirements-vs-implementation review | Kimi K2P6 low reasoning, then ZAI/GLM-5.1 replacement re-review | no_remaining_critical_or_major | Kimi found malformed named protected inputs returned `exitCannotVerify` instead of usage exit `2`; fixed with usage exit behavior and `TestProtectedGateMalformedNamedInputIsUsageError`. Kimi also flagged minor coverage/fixture drift; malformed checkpoint/policy/witness coverage, missing-policy trust-scope handling, and missing-policy fixture text were fixed. A follow-up Kimi run hung after reading fixtures and was replaced. GLM then found two local-signed fixture signer-authority overclaims; fixtures and regression test were fixed, and GLM re-review confirmed no remaining critical or major findings. |

## Trust Boundaries

- Local signed checkpoint evidence is rejected for protected pass.
- CI signed protected pass requires checkpoint signature/binding, signer policy, CI witness binding, and freshness.
- Block 16 does not decide merge, release, readiness, degradation, override approval, or risk acceptance.
- External witness trust remains `not_integrated`.

## Review Requirements

Before closure, run and record:

- implementation code review: no remaining critical or major findings before PR;
- tracing review: no remaining critical or major findings before PR;
- requirements-vs-implementation review: no remaining critical or major findings before PR;
- PR-level code review: no remaining critical or major findings;
- PR-level tracing review: no remaining critical or major findings;
- PR-level requirements-vs-implementation review: no remaining critical or major findings.
