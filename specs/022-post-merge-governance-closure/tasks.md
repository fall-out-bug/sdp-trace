---
description: "Task list for Spec 022 post-merge governance closure"
---

# Tasks: Post-Merge Governance Closure

**Input**: Design documents from `specs/022-post-merge-governance-closure/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`,
`quickstart.md`

**Tests**: No behavior, verifier, schema, or product code changes are planned.
This docs-governance slice uses live PR/CI refresh where available plus local
documentation checks.

**Organization**: Tasks are grouped by independently testable governance
outcomes rather than runtime user stories.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it touches different files and does not
  depend on incomplete tasks.
- **[Story]**: Applies only to story phases.
- Every task names exact file paths.

## Phase 1: Setup

**Purpose**: Confirm the planned Spec 022 closure route and available design
artifacts.

- [ ] T001 Confirm `specs/022-post-merge-governance-closure/spec.md`, `specs/022-post-merge-governance-closure/plan.md`, `specs/022-post-merge-governance-closure/research.md`, `specs/022-post-merge-governance-closure/data-model.md`, and `specs/022-post-merge-governance-closure/quickstart.md` are present and describe Spec 022.
- [ ] T002 Confirm `.specify/scripts/bash/setup-tasks.sh --json` branch mismatch is documented in `specs/022-post-merge-governance-closure/plan.md`.
- [ ] T003 Confirm Spec 022 still has no planned product code, schema, command, or `/contracts` changes in `specs/022-post-merge-governance-closure/plan.md`.

## Phase 2: Foundational

**Purpose**: Establish the evidence sources and guardrails that all closure
stories depend on.

- [ ] T004 Map source evidence references from `specs/019-repo-realignment-monitoring-gate-readiness/plan.md`, `specs/019-repo-realignment-monitoring-gate-readiness/tasks.md`, and `specs/019-repo-realignment-monitoring-gate-readiness/post-merge-closure-plan.md`.
- [ ] T005 [P] Map decision state references from `docs/closure-decision-ledger.md` and `docs/spec-reality-ledger.md`.
- [ ] T006 [P] Map navigation state references from `docs/roadmap.md`.
- [ ] T007 Confirm `merge_approval`, `maintainer_approval`, `not_assessed`, and `cannot_verify` remain explicit in `specs/022-post-merge-governance-closure/spec.md`.
- [ ] T008 Confirm no task in `specs/022-post-merge-governance-closure/tasks.md` requires retroactive PR #60 approval or changes to existing commands.

**Checkpoint**: Evidence sources and scope guardrails are ready before closure
surface edits.

## Phase 3: User Story 1 - Governance Evidence Summary (Priority: P1)

**Goal**: A maintainer can inspect Spec 019 residual governance evidence with
exact PR, commit, CI, and review references, while PR #60 approval remains
explicitly `not_assessed` unless new approval evidence exists.

**Independent Test**: Review `docs/spec-reality-ledger.md` and confirm PR #60,
PR #63, Spec 019, and Spec 022 states are stated with exact references and no
retroactive approval claim.

- [ ] T009 [P] [US1] Refresh or record unavailable live state for PR #60 in `docs/spec-reality-ledger.md` using the quickstart commands from `specs/022-post-merge-governance-closure/quickstart.md`.
- [ ] T010 [P] [US1] Refresh or record unavailable live state for PR #63 in `docs/spec-reality-ledger.md` using the quickstart commands from `specs/022-post-merge-governance-closure/quickstart.md`.
- [ ] T011 [US1] Update the Spec 019 row in `docs/spec-reality-ledger.md` with current PR #60, PR #63, CI, review, and missing-approval evidence.
- [ ] T012 [US1] Update the Spec 022 row in `docs/spec-reality-ledger.md` from prepared follow-up state to active closure state with the current residual-governance summary.
- [ ] T013 [US1] Verify `docs/spec-reality-ledger.md` keeps PR #60 merge approval as `not_assessed` unless explicit approval evidence was found.

**Checkpoint**: Governance evidence summary is independently reviewable.

## Phase 4: User Story 2 - Maintainer Decision And Remediation Disposition (Priority: P2)

**Goal**: A maintainer can see that `split_successor` is the current decision
and that residual remediation is either explicitly none or represented by
reviewed successor specs before implementation.

**Independent Test**: Review `docs/closure-decision-ledger.md` and confirm D006
preserves `split_successor`, names any residual remediation state, and does not
infer approval from CI, reviews, or checked task boxes.

- [ ] T014 [P] [US2] Update D006 in `docs/closure-decision-ledger.md` to cite Spec 022 plan, tasks, and quickstart references from `specs/022-post-merge-governance-closure/`.
- [ ] T015 [P] [US2] Inspect `docs/open-task-breakdown.md` for any Spec 019 or Spec 022 residual task references that contradict `split_successor`.
- [ ] T016 [US2] Record the residual remediation state in `docs/closure-decision-ledger.md` as either no residual remediation remains or successor specs are required.
- [ ] T017 [US2] If successor specs are required, add or update their reviewed triplet references in `docs/closure-decision-ledger.md`; otherwise record no-remediation evidence in `docs/closure-decision-ledger.md`.
- [ ] T018 [US2] Confirm `docs/closure-decision-ledger.md` does not reopen accept/reject/split for Spec 019 unless a new maintainer decision explicitly supersedes `split_successor`.

**Checkpoint**: Maintainer decision and remediation disposition are explicit.

## Phase 5: User Story 3 - Synchronized Closure Navigation (Priority: P3)

**Goal**: Contributors see the same Spec 022 closure state in the decision
ledger, spec reality ledger, and roadmap.

**Independent Test**: Compare `docs/closure-decision-ledger.md`,
`docs/spec-reality-ledger.md`, and `docs/roadmap.md`; all three surfaces report
the same Spec 022 state and next step.

- [ ] T019 [P] [US3] Update Spec 022 status and next-step wording in `docs/roadmap.md`.
- [ ] T020 [P] [US3] Update Spec 019 status and residual-governance wording in `docs/roadmap.md` to point to the current Spec 022 closure state.
- [ ] T021 [US3] Cross-check `docs/closure-decision-ledger.md`, `docs/spec-reality-ledger.md`, and `docs/roadmap.md` for consistent Spec 022 state wording.
- [ ] T022 [US3] Update `specs/022-post-merge-governance-closure/spec.md` status if the implementation changes the lifecycle state.
- [ ] T023 [US3] Update `specs/022-post-merge-governance-closure/tasks.md` checkboxes only after the matching evidence edits and local verification exist.

**Checkpoint**: Closure navigation surfaces are synchronized.

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Validate docs and keep trust-state claims bounded.

- [ ] T024 Run `go run ./tools/doccheck` and record the result in the final response for `specs/022-post-merge-governance-closure/tasks.md`.
- [ ] T025 Run `go run ./tools/hygienecheck` and record the result in the final response for `specs/022-post-merge-governance-closure/tasks.md`.
- [ ] T026 Run `git diff --check` and record the result in the final response for `specs/022-post-merge-governance-closure/tasks.md`.
- [ ] T027 If implementation expands beyond docs, run `go test ./...`, `go vet ./...`, and `go run ./tools/schemadoc`; otherwise record those checks as `not_assessed` in the final response.
- [ ] T028 Review `docs/closure-decision-ledger.md`, `docs/spec-reality-ledger.md`, `docs/roadmap.md`, and `specs/022-post-merge-governance-closure/spec.md` for forbidden claims about retroactive approval, production trust, release approval, or external attestation.

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies.
- **Foundational (Phase 2)**: Depends on Setup completion and blocks all user
  stories.
- **US1 Governance Evidence Summary (Phase 3)**: Depends on Foundational.
- **US2 Maintainer Decision And Remediation Disposition (Phase 4)**: Depends on
  Foundational; can start after T005 but should use US1 evidence when possible.
- **US3 Synchronized Closure Navigation (Phase 5)**: Depends on US1 and US2.
- **Polish (Phase 6)**: Depends on desired story edits.

### User Story Dependencies

- **US1 (P1)**: First independently valuable slice; produces evidence summary.
- **US2 (P2)**: Uses the decision ledger and remediation state; can partially
  proceed in parallel with US1 after evidence sources are mapped.
- **US3 (P3)**: Must run after US1 and US2 so navigation reflects final state.

### Parallel Opportunities

- T005 and T006 can run in parallel after T004.
- T009 and T010 can run in parallel if live PR/CI refresh is available.
- T014 and T015 can run in parallel after Phase 2.
- T019 and T020 can run in parallel after US1 and US2 are complete.

## Parallel Example: User Story 1

```text
Task: "Refresh or record unavailable live state for PR #60 in docs/spec-reality-ledger.md"
Task: "Refresh or record unavailable live state for PR #63 in docs/spec-reality-ledger.md"
```

## Parallel Example: User Story 2

```text
Task: "Update D006 in docs/closure-decision-ledger.md"
Task: "Inspect docs/open-task-breakdown.md for residual task contradictions"
```

## Implementation Strategy

### MVP First

1. Complete Phase 1 and Phase 2.
2. Complete US1 so the current governance evidence is visible and bounded.
3. Stop and verify `docs/spec-reality-ledger.md` before editing closure
   decisions.

### Incremental Delivery

1. US1: Evidence summary and live-state boundaries.
2. US2: Maintainer decision and remediation disposition.
3. US3: Synchronized roadmap/reality/decision surfaces.
4. Polish: local docs checks and forbidden-claim review.

### Completion Rule

Spec 022 can move toward `complete` only when evidence is cited and remaining
work is either explicitly none or filed as reviewed successor specs. Missing
live state remains `not_assessed` or `cannot_verify`.
