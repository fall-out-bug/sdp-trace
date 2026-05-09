# Socratic Review: Authority Envelope Boundary Observation

**Status**: Socratic spec review complete; no unresolved critical or major findings; awaiting user approval for implementation
**Date**: 2026-05-09

This file records Socratic pi-review findings and dispositions for the spec package. It is intentionally separate from implementation review.

## Review Inputs

- `spec.md`
- `plan.md`
- `data-model.md`
- `tasks.md`

## Review Plan

Run independent reviewer planes:

1. Product boundary and policy-separation critic.
2. Evidence/source-attribution critic.
3. Schema/data-model consistency critic.

Reviewers must challenge whether the spec:

- accidentally encodes demo-specific policy;
- overclaims actor, tool, or model attribution;
- preserves `not_assessed` and `cannot_verify`;
- gives implementers enough fixture guidance;
- keeps external policy decisions outside `sdp-trace`.

## Reviewer Runs

| Plane | Model | Result |
|---|---|---|
| Product boundary and policy separation | `zai/glm-5.1` | usable review, verdict `REVISE` |
| Evidence/source attribution | `minimax/MiniMax-M2.7` | usable review, verdict `CHANGES_REQUESTED` |
| Schema/data-model consistency | `openrouter/qwen/qwen3.6-plus` | first attempt off-task; retry usable, verdict `CHANGES_REQUESTED` |

## Findings And Dispositions

### S1: Policy selection mechanism undefined

Severity: major

Reviewer concern: The draft required evaluations to cite `policy_id` but did not state who selects that policy when multiple envelopes exist.

Disposition: accepted/fixed.

Resolution:

- `spec.md` now states that authority evaluation uses caller-selected `policy_id`.
- Missing selected policy yields `not_assessed`.
- `plan.md` records caller-supplied `policy_id` as a contract decision.
- `data-model.md` states that `sdp-trace` does not choose among competing envelopes.

### S2: Allow/deny and target-rule conflict semantics undefined

Severity: major

Reviewer concern: Top-level and target-specific allow/deny rules could conflict, leaving implementations to choose allow-wins, deny-wins, or first-match behavior.

Disposition: accepted/fixed.

Resolution:

- `spec.md` edge cases now declare conflicting envelopes invalid.
- `data-model.md` requires conflict-free target rules and makes conflicts `cannot_verify`.
- `plan.md` records "reject ambiguous envelopes" as a contract decision.

### S3: Evidence reference format undefined

Severity: major

Reviewer concern: `evidence_refs` existed without a resolvable format, making fixture validation and referential integrity ambiguous.

Disposition: accepted/fixed.

Resolution:

- `data-model.md` now defines safe URI-style evidence refs: `file:`, `artifact:`, `git:`, and `external:`.
- Unresolvable required references produce `cannot_verify`.
- `plan.md` records evidence ref schemes as a contract decision.

### S4: Binding protocol and failed binding semantics undefined

Severity: critical

Reviewer concern: The draft did not define how bindings across git, harness, and gateway sources succeed, partially succeed, or fail.

Disposition: accepted/fixed.

Resolution:

- `spec.md` now includes binding-state decision rules.
- `data-model.md` adds `matched_fields` and explicit `verified` / `not_assessed` / `cannot_verify` binding criteria.
- `tasks.md` adds failed-binding fixtures and evidence-binding schema work.

### S5: `not_assessed` vs `cannot_verify` decision logic incomplete

Severity: major

Reviewer concern: Evaluation states were listed but not operationalized for absent policy, malformed policy, insufficient attribution, and missing evidence.

Disposition: accepted/fixed.

Resolution:

- `spec.md` now includes an evaluation decision table.
- Action existence can be evaluated while actor/tool/model attribution remains independently `not_assessed`.

### S6: Stale or inaccessible evidence not covered

Severity: major

Reviewer concern: The draft covered policy supersession but not stale, inaccessible, overwritten, or pre-policy evidence.

Disposition: accepted/fixed.

Resolution:

- `spec.md` decision rules now treat stale or inaccessible required evidence as `cannot_verify`.
- `data-model.md` adds event ordering rules.
- `tasks.md` adds stale/inaccessible evidence fixtures.

### S7: `confidence` and `evidence_quality` risked policy creep

Severity: major

Reviewer concern: `confidence` and `evidence_quality` could become implicit policy levers.

Disposition: accepted/fixed.

Resolution:

- `data-model.md` replaces `confidence` with `observation_confidence`, scoped to one observed action.
- `AuthorityEvaluation.evidence_quality` was removed and replaced with descriptive `source_coverage`.
- `source_coverage` explicitly has no native policy meaning.

### S8: `mutation_kind` encoded workflow-specific taxonomy

Severity: major

Reviewer concern: `mutation_kind` values such as `ci_workflow_edit` and `planning_edit` risk encoding local workflow policy in the product model.

Disposition: accepted/fixed.

Resolution:

- `mutation_kind` was removed from `ObservedAction`.
- Target semantics are now expressed through provider-neutral target refs and external policy target rules.

### S9: `AuthorityEnvelope.actor` vs `ActorDeclaration` relationship ambiguous

Severity: major

Reviewer concern: The draft did not say whether the actor was inline or referenced.

Disposition: accepted/fixed.

Resolution:

- `AuthorityEnvelope.actor` is now `actor_ref`, a required reference to `ActorDeclaration.actor_id`.

### S10: Required vs optional fields missing

Severity: major

Reviewer concern: Entities did not distinguish required and optional fields, making schema generation under-specified.

Disposition: accepted/fixed.

Resolution:

- `data-model.md` now marks fields as `[required]` or `[optional]`.

### S11: Event ordering for policy supersession under-specified

Severity: moderate

Reviewer concern: Supersession requires ordering across heterogeneous sources.

Disposition: accepted/fixed.

Resolution:

- `data-model.md` now requires explicit event ids and source refs for ordering.
- Wall-clock timestamps cannot be the only ordering evidence for authority evaluation.

### S12: Feedback/review without mutation needs a fixture

Severity: minor

Reviewer concern: The edge case was documented but not planned as a fixture.

Disposition: accepted/fixed.

Resolution:

- `tasks.md` adds a feedback/review fixture proving no mutation fact is inferred without mutation evidence.

### S13: Example path used block-specific naming

Severity: minor

Reviewer concern: `examples/block29-authority-envelope/` looked project-context-specific.

Disposition: accepted/fixed.

Resolution:

- Planned examples path changed to `examples/authority-envelope-basic/`.

## Remaining Review State

The initial Socratic review found valid blockers. The spec package was revised, then focused re-review found two remaining major product-boundary findings and one minor fixture finding. Those residual findings are recorded below and have been revised. Final focused re-review found no unresolved critical or major findings.

## Focused Re-Review Findings And Dispositions

### R1: Top-level vs target-rule interaction undefined

Severity: major

Reviewer concern: The conflict rule covered target-vs-target overlap but not top-level-vs-target disagreement.

Disposition: accepted/fixed.

Resolution:

- `data-model.md` now states that top-level event rules and target rules participate in the same conflict check.
- Top-level-vs-target disagreement makes the envelope invalid and evaluation `cannot_verify`.

### R2: Unmatched-action fallthrough undefined

Severity: major

Reviewer concern: The decision table did not state what happens when an action is observed but no event or target rule applies.

Disposition: accepted/fixed.

Resolution:

- `spec.md` now states that an observed action with no applicable rule in the selected envelope is `not_assessed`.
- `plan.md` records "unmatched actions are unassessed, not denied by default" as a contract decision.
- `data-model.md` mirrors the same rule.

### R3: "Protected path" wording implied built-in policy

Severity: minor

Reviewer concern: Acceptance scenarios used "protected path/file" terminology even though protection must come from selected envelope rules.

Disposition: accepted/fixed.

Resolution:

- `spec.md` now refers to paths matching selected envelope target or denied-events rules.

### R4: Evidence resolution "profile" referenced but undefined

Severity: minor

Reviewer concern: `data-model.md` referenced selected profile without defining a profile entity.

Disposition: accepted/fixed.

Resolution:

- `data-model.md` now uses "selected evidence-resolution context" and explicitly says unresolved required `external:` evidence is `cannot_verify`.

### R5: `custom` event type underspecified

Severity: minor

Reviewer concern: A bare `custom` event type could not match authority rules deterministically.

Disposition: accepted/fixed.

Resolution:

- `data-model.md` now requires extension event types to use a concrete `custom:` prefix.

### R6: Malformed-policy fixture was too vague

Severity: minor

Reviewer concern: The task could cover only JSON-invalid policy and miss semantically invalid allow/deny conflicts.

Disposition: accepted/fixed.

Resolution:

- `tasks.md` now requires malformed-policy fixtures for both JSON/schema-invalid envelopes and semantically invalid conflicting rules.

## Final Focused Re-Review

| Plane | Model | Result |
|---|---|---|
| Product boundary and policy separation | `zai/glm-5.1` | approved; no critical or major findings |
| Evidence/source attribution | `minimax/MiniMax-M2.7` | approved; no critical or major findings |
| Schema/data-model consistency | `openrouter/deepseek/deepseek-v4-pro` | approved; no critical or major findings |

### F1: Fixture name still used protected-path language

Severity: minor

Reviewer concern: Planned fixture name `outside-authority-protected-path/` could imply a built-in product concept.

Disposition: accepted/fixed.

Resolution:

- `plan.md` now uses `outside-authority-denied-target/`.
- `tasks.md` now uses "denied-target mutation fixture."

### F2: FR-002 used prose event names instead of exact enum values

Severity: minor

Reviewer concern: `spec.md` FR-002 listed prose event names while `data-model.md` defined concrete `event_type` strings and the `custom:` extension prefix.

Disposition: accepted/fixed.

Resolution:

- `spec.md` FR-002 now lists exact event type strings and extension event types prefixed with `custom:`.

### F3: Policy-specific fixture coverage was implicit

Severity: minor

Reviewer concern: User Story 3 requires the same observed action to evaluate differently under two selected `policy_id` values, but the task list did not explicitly require that fixture.

Disposition: accepted/fixed.

Resolution:

- `tasks.md` now adds a policy-specific fixture for two selected `policy_id` values with different authority states.
