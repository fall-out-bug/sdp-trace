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

PR-level CI state, merge, and post-merge verification are not completed in this local implementation review.

## PR-Level Review

PR: https://github.com/fall-out-bug/sdp-trace/pull/26

| Plane | Model | Result |
|---|---|---|
| Code/correctness | `zai/glm-5.1` | approved; minor coverage/readability suggestions |
| Tracing/evidence | `minimax/MiniMax-M2.7` | changes requested; focused re-review approved |
| Requirements vs implementation | `openrouter/deepseek/deepseek-v4-pro` | changes requested; focused re-review approved |

### P1: Git-only evidence could verify actor/tool attribution

Severity: major

Reviewer concern: A git-sourced action with non-empty `actor_id` or `operation_id` could produce `actor_attribution: verified` or `tool_attribution: verified`, violating FR-006.

Disposition: accepted/fixed.

Resolution:

- Replaced value-only actor attribution with source-aware `actorAttributionState`.
- `source_type: git` no longer verifies actor attribution even when `actor_id` is present.
- Tool attribution now requires `source_type: harness_log` and `operation_id`.
- Updated git-only regression test so it keeps actor/tool fields present and asserts all attribution dimensions remain `not_assessed`.
- Focused PR-level requirements re-review approved the fix.

### P2: Missing approval evidence branch lacked direct test coverage

Severity: minor

Reviewer concern: The `approval_evidence_missing` branch was implemented but not directly tested.

Disposition: accepted/fixed.

Resolution:

- Added `TestEvaluateMissingApprovalEvidenceIsOutsideAuthority`.

### P3: Git-only fixture reason-code mismatch

Severity: major from reviewer; adjudicated false positive after focused evidence.

Reviewer concern: The git-only fixture expected `target_event_allowed` while implementation allegedly returned `event_allowed`.

Disposition: rejected as false positive with evidence.

Resolution:

- Ran the git-only fixture through the CLI and confirmed `reason_code: target_event_allowed`.
- The same output confirmed actor/tool/model attribution as `not_assessed`.
- Focused PR-level tracing/evidence re-review approved with no remaining findings.
