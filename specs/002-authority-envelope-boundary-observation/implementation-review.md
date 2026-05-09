# Implementation Review: Authority Envelope Boundary Observation

**Date**: 2026-05-09
**Status**: implementation review complete; no unresolved critical or major findings

## Review Planes

| Plane | Model | Result |
|---|---|---|
| Requirements vs implementation | `openrouter/deepseek/deepseek-v4-pro` | approved; no findings |
| Tracing/evidence | `minimax/MiniMax-M2.7` | changes requested; one major fixed; focused re-review approved |
| Code/correctness | `zai/glm-5.1` | approved with minor observations; focused re-review approved after review-packet correction |

## Findings And Dispositions

### I1: Gateway source alone verified model attribution

Severity: major

Reviewer concern: `source_type: llm_gateway` could promote `model_attribution` to `verified` without a verified `same_gateway_request` binding.

Disposition: accepted/fixed.

Resolution:

- Removed the gateway-source shortcut from `hasVerifiedGatewayBinding`.
- `model_attribution` now becomes `verified` only when a `same_gateway_request` binding has state `verified`.
- Added `TestEvaluateGatewaySourcedMutationWithoutBindingDoesNotVerifyModel`.
- Focused tracing/evidence re-review approved the fix.

### I2: Semantically overlapping target rules were only partially handled

Severity: minor

Reviewer concern: Static envelope validation caught only identical-pattern target conflicts. Different glob patterns could both match one target and disagree.

Disposition: accepted/fixed.

Resolution:

- Added runtime conflict detection in `matchDecision` for multiple matching target rules that disagree for the same observed action.
- Added `TestEvaluateOverlappingTargetRulesCannotVerify`.
- Added `examples/authority-envelope-basic/malformed-policy-cannot-verify/authority-package-overlap.json`.
- Added fixture-matrix coverage for `overlapping_target_rules_conflict`.
- Focused code re-review approved the fix.

### I3: Review packet initially omitted a new fixture file

Severity: process/minor

Reviewer concern: The first focused code re-review did not see `authority-package-overlap.json` and reported it missing.

Disposition: accepted/fixed.

Resolution:

- Added the file to intent-to-add so `git diff` included it in the review packet.
- Verified the file exists and is valid JSON.
- Ran `go test ./internal/authority -run TestFixtureMatrixScenarios -count=1`.
- Re-ran focused code review with the corrected packet; verdict approved.

## Verification

Latest verification commands:

- `jq empty schema/*.json examples/authority-envelope-basic/fixture-matrix.json examples/authority-envelope-basic/*/authority-package*.json examples/authority-envelope-basic/policy-specific/authority-package*.json`
- `go test ./...`
- `git diff --check HEAD`

## Remaining State

PR-level review, CI state, merge, and post-merge verification are not completed in this local implementation review.

