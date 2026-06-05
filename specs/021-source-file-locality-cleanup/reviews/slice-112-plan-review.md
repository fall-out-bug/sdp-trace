# Slice 112 Plan Review

Date: 2026-06-05T00:57:31Z

Scope:
- Consolidate `internal/packet/demo_first_state_helpers.go` into `internal/packet/demo_first_closure_requirements.go`.
- Preserve the pass, partial, and fail assessed-state set for the demo-first verification/review closure requirement.
- Preserve cannot_verify, not_assessed, and not_in_scope as not assessed for this closure requirement.
- Exclude demo-first route evidence, row evidence, pass-count logic, closure cap logic, packet schemas, examples, fixtures, dependencies, package boundary, dependency direction, CRAP/MI baselines.

Decision gate:
- Simpler/Faster: Move `rowAssessed` next to its only consumer; no new abstraction, dependency, or public API.
- Blocking Edge Cases: The focused test must cover the state matrix and exact diagnostic, not only one rejection path.
- Existing Open Source: Not applicable; this is local file ownership cleanup using existing Go tests.

Rejected feasibility path:
- Initial Slice 112 candidate attempted movement value/comparison consolidation into `posture_validate_movement_row.go`.
- Result: rejected because the consolidated file failed file-level MI (`maintainability index 65.0`) without a baseline.
- Decision: do not add or adjust MI baselines for this cleanup; choose a smaller cohesive packet slice instead.

Initial packet-scope review:

| Reviewer | Agent ID | Harness | Model/provider | Prompt class | Result |
|---|---|---|---|---|---|
| Beauvoir | `019e9406-f078-7fd2-b8d0-e22ac17a1e3a` | Codex subagent | not_assessed | plan review | finding |
| Peirce | `019e9406-f40c-79f1-904e-54d0f0b73866` | Codex subagent | not_assessed | plan review | finding |
| Halley | `019e9406-f7c2-7f80-80d9-86f7cf7e0c22` | Codex subagent | not_assessed | plan review | LGTM |

Findings:
- major: T021-7790 named `TestValidateFirstPacketDemoRejectsWeakPackets` and `TestValidateFirstPacketDemoAcceptsDemoFixture`, but those tests do not exist in `internal/packet`.
- major: T021-7790 did not require the existing focused test to cover the `rowAssessed` state matrix and exact first-packet gate diagnostic.

Fix:
- Updated T021-7790 to use `TestCheckDemoRequiresVerificationOrReviewAssessed`.
- Required the exact-count guard to expect exactly one focused test.
- Required pass, partial, and fail accepted coverage plus cannot_verify, not_assessed, and not_in_scope rejection coverage.
- Required preserving the exact first-packet gate diagnostic.

Re-review:

| Reviewer | Agent ID | Harness | Model/provider | Prompt class | Result |
|---|---|---|---|---|---|
| Beauvoir | `019e9406-f078-7fd2-b8d0-e22ac17a1e3a` | Codex subagent | not_assessed | plan re-review | LGTM |
| Peirce | `019e9406-f40c-79f1-904e-54d0f0b73866` | Codex subagent | not_assessed | plan re-review | LGTM |
| Halley | `019e9406-f7c2-7f80-80d9-86f7cf7e0c22` | Codex subagent | not_assessed | plan re-review | LGTM |

Review state: pass.
