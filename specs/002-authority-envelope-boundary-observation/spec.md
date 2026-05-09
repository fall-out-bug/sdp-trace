# Feature Specification: Authority Envelope Boundary Observation

**Feature Branch**: `002-authority-envelope-boundary-observation`
**Created**: 2026-05-09
**Status**: Draft - revised after initial Socratic review, focused re-review pending
**Input**: User description: "`sdp-trace` should record when an actor acts outside the authority declared for a task. The product must not hard-code demo-specific policy; it should record declared authority, observed actions, evidence sources, and explicit `not_assessed`/`cannot_verify` gaps."

## Product Boundary

`sdp-trace` records authority facts. It does not decide employment discipline, demo contamination, merge blocking, risk acceptance, or organizational blame.

External policy consumers decide what to do with an authority boundary finding. `sdp-trace` must provide enough evidence for those consumers to distinguish:

- an action that is within declared authority;
- an action that is outside declared authority;
- an action whose actor/model/tool attribution is not assessed;
- an action whose required evidence cannot be verified.

## User Scenarios & Testing

### User Story 1 - Reviewer Sees Whether Actor Actions Matched Declared Authority (Priority: P1)

A reviewer can inspect a trace package and see the declared authority for each actor/task, the observed actions, and the evaluation state for each action.

**Why this priority**: Without an explicit authority envelope, a green PR or passing CI can hide that the wrong actor performed the work.

**Independent Test**: A committed fixture shows one actor declared as read-only/review-only, a denied-target mutation observed from evidence, and an evaluation result of `outside_authority` with evidence references.

**Acceptance Scenarios**:

1. **Given** an actor is declared as review-only for a task, **When** a direct mutation event targets a path that matches a denied-events rule, **Then** the evaluation records `outside_authority`.
2. **Given** a mutation is visible only from git diff, **When** no harness/tool-call evidence is available, **Then** mutation existence is recorded but actor/model attribution remains `not_assessed`.
3. **Given** an actor has explicit approval evidence for a mutation, **When** the target path and event type match the approved scope, **Then** the evaluation records `within_authority`.

---

### User Story 2 - Integrator Imports Multiple Observation Sources (Priority: P1)

An integrator can import observations from git, PR/MR metadata, CI artifacts, harness logs, and LLM gateway logs without requiring all sources to exist.

**Why this priority**: Different harnesses expose different evidence. The product must not pretend full attribution exists when only consequences are observable.

**Independent Test**: Fixtures prove that git-only, git-plus-harness, and git-plus-harness-plus-gateway inputs produce different attribution quality states.

**Acceptance Scenarios**:

1. **Given** only git evidence, **When** a file matching a selected envelope target rule changes, **Then** `sdp-trace` records the changed path and source commit but does not claim tool-call or model attribution.
2. **Given** harness evidence with operation ids and affected paths, **When** those paths match git changes, **Then** `sdp-trace` can bind the observed action to the harness actor and operation id.
3. **Given** gateway evidence with request ids but no harness binding, **When** a mutation exists, **Then** model identity remains `not_assessed` for that mutation.

---

### User Story 3 - Policy Consumer Separates Product Facts From Local Policy (Priority: P1)

An external policy consumer can use authority evaluation facts without `sdp-trace` embedding that consumer's specific rules.

**Why this priority**: Hard-coding one team's workflow would make `sdp-trace` a narrow demo harness instead of a portable trust substrate.

**Independent Test**: The same observed action can be evaluated against two different authority envelopes, producing different states without changing the observed action record.

**Acceptance Scenarios**:

1. **Given** one policy allows CI edits for an implementer and another denies them, **When** the same CI mutation is evaluated with an operator-supplied `policy_id`, **Then** the result is policy-specific and cites that selected `policy_id`.
2. **Given** no authority envelope is supplied, **When** observed actions exist, **Then** authority evaluation remains `not_assessed` rather than inferring a default policy.
3. **Given** a downstream gate wants to block on `outside_authority`, **When** it consumes the evaluation, **Then** the block decision is external to `sdp-trace`.

## Edge Cases

- A mutation was performed before recorder attachment: pre-attachment attribution is `not_assessed` or `cannot_verify`.
- A commit author is a bot or shared account: git identity is preserved but human/model attribution is not inferred.
- A harness reports a model name without gateway evidence: model identity is `harness_observed` or `agent_reported`, not gateway-verified.
- A gateway request exists but cannot be linked to a tool call: model attribution for a mutation remains `not_assessed`.
- A selected envelope target pattern is ambiguous or invalid: policy evaluation for matching actions is `cannot_verify`.
- A policy envelope changes mid-task: the original envelope is preserved and the new envelope supersedes it with an explicit effective time or event id.
- A reviewer comment suggests a change but no direct mutation occurs: record a review/feedback event, not a mutation event.
- CI rewrites generated files: record CI as the actor/source when evidence supports that; otherwise keep actor attribution `not_assessed`.
- The only evidence is prose: record it as low-authority external assertion; do not use it to prove mutation attribution.
- A checked-in mutable status file claims authority compliance: it is evidence only if tied to an immutable source and verified by the selected evidence-resolution context.
- Multiple authority envelopes could apply: `sdp-trace` does not choose the canonical policy; the caller must provide the selected `policy_id` or authority evaluation is `not_assessed`.
- An envelope has conflicting allow/deny rules: the envelope is invalid and evaluation is `cannot_verify`, not "deny wins" or "allow wins."
- Evidence sources exist but fail to bind: the binding is `cannot_verify`; attribution from that binding cannot upgrade from `not_assessed`.
- A review or feedback event suggests a change without direct mutation evidence: record the event as review/feedback; do not create a mutation fact from it.

## Functional Requirements

- **FR-001**: `sdp-trace` MUST define an authority envelope artifact that records task id, policy id, actor id, declared role, allowed events, denied events, target rules, and approval requirements.
- **FR-002**: `sdp-trace` MUST define observed action records for `observe`, `review`, `feedback`, `direct_mutation`, `commit`, `merge`, `ci_run`, `harness_tool_call`, `gateway_request`, and extension event types prefixed with `custom:`.
- **FR-003**: Observed actions MUST record source type and evidence reference: git, PR/MR API, CI artifact, harness log, LLM gateway, manual import, or external assertion.
- **FR-004**: Authority evaluation MUST emit one of `within_authority`, `outside_authority`, `not_assessed`, or `cannot_verify`, using the decision rules below.
- **FR-005**: Authority evaluation MUST distinguish action existence from actor attribution, tool attribution, and model attribution.
- **FR-006**: Git-only evidence MAY prove that a path changed, but MUST NOT prove tool-call, harness, model, or human attribution by itself.
- **FR-007**: Harness-log evidence MAY bind actor, role, operation id, tool call, cwd, and affected paths when the imported log contains those fields.
- **FR-008**: Gateway evidence MAY bind provider/model/request id only when it is linked to a harness operation or observed action by an explicit binding.
- **FR-009**: `sdp-trace` MUST NOT infer authority policy from demo names, harness names, actor names, model names, branch names, commit messages, or prose.
- **FR-010**: Missing authority envelopes MUST produce `not_assessed`, not `within_authority`.
- **FR-011**: Invalid, malformed, stale, or inaccessible required evidence MUST produce `cannot_verify` with a reason code.
- **FR-012**: Authority evaluation MUST cite the selected `policy_id` and evidence refs used for each evaluated action.
- **FR-013**: Authority envelopes MUST support target path rules without requiring GitHub, GitLab, or any specific VCS provider.
- **FR-014**: Authority envelope artifacts MUST be portable JSON and safe to commit: no raw prompts, secrets, credentials, private source snippets, or raw model outputs.
- **FR-015**: `sdp-trace` MUST document that contamination, blocking, disciplinary, readiness, or merge decisions belong to external policy consumers.

## Evaluation Decision Rules

Authority evaluation uses an explicit caller-selected `policy_id`. If multiple envelopes exist, `sdp-trace` does not select the most recent, strictest, or broadest policy.

| Condition | Evaluation state |
|---|---|
| No selected `policy_id` or no matching envelope for the selected `policy_id` | `not_assessed` |
| Selected envelope is malformed, ambiguous, inaccessible, or has conflicting rules | `cannot_verify` |
| Required evidence reference is missing, inaccessible, malformed, stale for the selected event boundary, or fails integrity checks | `cannot_verify` |
| Action is observed but no event or target rule applies in the selected envelope | `not_assessed` |
| Action existence is observed but actor/tool/model attribution source is absent | Evaluate action existence against the selected policy; set missing attribution fields to `not_assessed` |
| Action and required attribution are observed and match allowed policy | `within_authority` |
| Action is observed and contradicts selected policy | `outside_authority` |

Binding state uses these rules:

| Binding condition | Binding state |
|---|---|
| Required source evidence is absent | `not_assessed` |
| Source evidence exists but cannot be read, parsed, integrity-checked, or compared | `cannot_verify` |
| Source evidence exists and declared binding fields all match | `verified` |
| Source evidence exists and declared binding fields disagree | `cannot_verify` |

Gateway evidence without a verified harness/action binding can prove only that a model request existed. It cannot prove mutation existence or model causality for a mutation.

## Success Criteria

- **SC-001**: A reviewer can inspect committed examples and explain why one action is `outside_authority` without relying on private chat context.
- **SC-002**: A git-only fixture proves mutation existence while keeping actor/model attribution `not_assessed`.
- **SC-003**: A harness-plus-git fixture proves stronger actor/operation attribution than git alone.
- **SC-004**: A gateway-without-harness-binding fixture keeps mutation model attribution `not_assessed`.
- **SC-005**: Documentation clearly separates generic authority evaluation from workflow-specific contamination policy.
