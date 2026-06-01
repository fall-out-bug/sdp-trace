---
description: "Task list for Spec 022 post-merge governance closure"
status: "in_progress"
---

# Tasks: Post-Merge Governance Closure

**Input**: Design documents from `specs/022-post-merge-governance-closure/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`,
`quickstart.md`

**Tests**: No behavior, verifier, schema, or product code changes are planned.
This docs-governance slice uses live PR/CI refresh plus local documentation
checks. If GitHub access is unavailable in a future environment, record the
exact access failure as `cannot_verify` or `not_assessed`; in this worktree,
GitHub access is available and PR #60 / PR #63 refresh is required.

**Organization**: Tasks are grouped by independently testable governance
outcomes rather than runtime user stories.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it touches different files and does not
  depend on incomplete tasks.
- **[Story]**: Applies only to story phases.
- Every task names exact file paths.
- Story review tasks are required checkpoints, not optional polish.

## Phase 1: Setup

**Purpose**: Confirm the planned Spec 022 closure route and available design
artifacts.

- [x] T001 Confirm `specs/022-post-merge-governance-closure/spec.md`, `specs/022-post-merge-governance-closure/plan.md`, `specs/022-post-merge-governance-closure/research.md`, `specs/022-post-merge-governance-closure/data-model.md`, and `specs/022-post-merge-governance-closure/quickstart.md` are present and describe Spec 022.
- [x] T002 Confirm `.specify/scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks` resolves `specs/022-post-merge-governance-closure/` on branch `codex/022-post-merge-governance-closure`.
- [x] T003 Confirm Spec 022 still has no planned product code, schema, command, or `/contracts` changes in `specs/022-post-merge-governance-closure/plan.md`.

## Phase 2: Foundational

**Purpose**: Establish the evidence sources and guardrails that all closure
stories depend on.

- [x] T004 Map source evidence references from `specs/019-repo-realignment-monitoring-gate-readiness/plan.md`, `specs/019-repo-realignment-monitoring-gate-readiness/tasks.md`, and `specs/019-repo-realignment-monitoring-gate-readiness/post-merge-closure-plan.md`.
- [x] T005 [P] Map decision state references from `docs/closure-decision-ledger.md` and `docs/spec-reality-ledger.md`.
- [x] T006 [P] Map navigation state references from `docs/roadmap.md`.
- [x] T007 Confirm `merge_approval`, `maintainer_approval`, `not_assessed`, and `cannot_verify` remain explicit in `specs/022-post-merge-governance-closure/spec.md`.
- [x] T008 Confirm no task in `specs/022-post-merge-governance-closure/tasks.md` requires retroactive PR #60 approval or changes to existing commands.
- [x] T009 Confirm FR-022-009 is already represented by this active unchecked `specs/022-post-merge-governance-closure/tasks.md` task list before closure surface edits begin.
- [x] T010 Confirm D006's stale `T120` reference in `docs/closure-decision-ledger.md` is tracked for removal or correction during the US2 D006 update.
- [x] T011 Re-check the post-design constitution gate in `specs/022-post-merge-governance-closure/plan.md` after the active task and review-cadence changes.

## Phase 2A: Pre-Implementation Review Gate

**Purpose**: Review and fix the active spec/plan/tasks before implementation.

- [x] T012 Run plan/task review for `specs/022-post-merge-governance-closure/spec.md`, `specs/022-post-merge-governance-closure/plan.md`, and `specs/022-post-merge-governance-closure/tasks.md` with Kimi for Coding 2.6, Z.AI 5.1, and MiniMax M3 or record exact unavailable model lanes.
- [x] T013 Record retained plan/task review findings and dispositions under `specs/022-post-merge-governance-closure/reviews/`.
- [x] T014 Fix all critical and major plan/task findings before editing `docs/closure-decision-ledger.md`, `docs/spec-reality-ledger.md`, or `docs/roadmap.md`.
- [x] T015 Re-run the affected plan/task review lanes until no retained critical or major findings remain.

**Checkpoint**: Evidence sources and scope guardrails are ready before closure
surface edits.

## Phase 3: User Story 1 - Governance Evidence Summary (Priority: P1)

**Goal**: A maintainer can inspect Spec 019 residual governance evidence with
exact PR, commit, CI, and review references, while PR #60 approval remains
explicitly `not_assessed` unless new approval evidence exists.

**Independent Test**: Review `docs/spec-reality-ledger.md` and confirm PR #60,
PR #63, Spec 019, and Spec 022 states are stated with exact references and no
retroactive approval claim.

- [x] T016 [P] [US1] Refresh live state for PR #60 in `docs/spec-reality-ledger.md` using the quickstart commands from `specs/022-post-merge-governance-closure/quickstart.md`.
- [x] T017 [P] [US1] Refresh live state for PR #63 in `docs/spec-reality-ledger.md` using the quickstart commands from `specs/022-post-merge-governance-closure/quickstart.md`.
- [x] T018 [US1] Update the Spec 019 row in `docs/spec-reality-ledger.md` with current PR #60, PR #63, CI, review, and missing-approval evidence.
- [x] T019 [US1] Update the Spec 022 row in `docs/spec-reality-ledger.md` from prepared follow-up state to active closure state with the current residual-governance summary.
- [x] T020 [US1] Verify `docs/spec-reality-ledger.md` keeps PR #60 merge approval as `not_assessed`.
- [x] T021 [US1] Run a focused US1 review of `docs/spec-reality-ledger.md` and record the review result under `specs/022-post-merge-governance-closure/reviews/`.

**Checkpoint**: Governance evidence summary is independently reviewable.

## Phase 4: User Story 2 - Maintainer Decision And Remediation Disposition (Priority: P2)

**Goal**: A maintainer can see that `split_successor` is the current decision
and that residual remediation is either explicitly none or represented by
reviewed successor specs before implementation.

**Independent Test**: Review `docs/closure-decision-ledger.md` and confirm D006
preserves `split_successor`, names any residual remediation state, and does not
infer approval from CI, reviews, or checked task boxes.

- [x] T022 [P] [US2] Update D006 in `docs/closure-decision-ledger.md` to cite Spec 022 plan, tasks, and quickstart references from `specs/022-post-merge-governance-closure/`, and remove or correct the stale non-existent Spec 019 `T120` reference.
- [x] T023 [P] [US2] Inspect `docs/open-task-breakdown.md` for any Spec 019 or Spec 022 residual task references that contradict `split_successor`.
- [x] T024 [US2] Record the residual remediation state in `docs/closure-decision-ledger.md` as either no residual remediation remains or successor specs are required.
- [x] T025 [US2] If successor specs are required, add or update their reviewed triplet and review-artifact references in `docs/closure-decision-ledger.md`; otherwise record no-remediation evidence in `docs/closure-decision-ledger.md`.
- [x] T026 [US2] Confirm `docs/closure-decision-ledger.md` does not reopen accept/reject/split for Spec 019 unless a new maintainer decision explicitly supersedes `split_successor`.
- [x] T027 [US2] Run a focused US2 review of `docs/closure-decision-ledger.md` and `docs/open-task-breakdown.md` and record the review result under `specs/022-post-merge-governance-closure/reviews/`.

**Checkpoint**: Maintainer decision and remediation disposition are explicit.

## Phase 5: User Story 3 - Synchronized Closure Navigation (Priority: P3)

**Goal**: Contributors see the same Spec 022 closure state in the decision
ledger, spec reality ledger, and roadmap.

**Independent Test**: Compare `docs/closure-decision-ledger.md`,
`docs/spec-reality-ledger.md`, and `docs/roadmap.md`; all three surfaces report
the same Spec 022 state and next step.

- [ ] T028 [P] [US3] Update Spec 022 status and next-step wording in `docs/roadmap.md`.
- [ ] T029 [P] [US3] Update Spec 019 status and residual-governance wording in `docs/roadmap.md` to point to the current Spec 022 closure state.
- [ ] T030 [US3] Cross-check `docs/closure-decision-ledger.md`, `docs/spec-reality-ledger.md`, and `docs/roadmap.md` for consistent Spec 022 state wording.
- [ ] T031 [US3] Ensure the three closure surfaces are committed together; do not split `docs/closure-decision-ledger.md`, `docs/spec-reality-ledger.md`, and `docs/roadmap.md` into separate commits for closure-state changes.
- [ ] T032 [US3] Update `specs/022-post-merge-governance-closure/spec.md` status if the implementation changes the lifecycle state.
- [ ] T033 [US3] Run a focused US3 review of `docs/closure-decision-ledger.md`, `docs/spec-reality-ledger.md`, `docs/roadmap.md`, and `specs/022-post-merge-governance-closure/spec.md`, then record the review result under `specs/022-post-merge-governance-closure/reviews/`.

**Checkpoint**: Closure navigation surfaces are synchronized.

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Validate docs and keep trust-state claims bounded.

- [ ] T034 Run `go run ./tools/doccheck` and record the result in `specs/022-post-merge-governance-closure/reviews/final-evidence.md`.
- [ ] T035 Run `go run ./tools/hygienecheck` and record the result in `specs/022-post-merge-governance-closure/reviews/final-evidence.md`.
- [ ] T036 Run `git diff --check` and record the result in `specs/022-post-merge-governance-closure/reviews/final-evidence.md`.
- [ ] T037 If implementation expands beyond docs, run `go test ./...`, `go vet ./...`, and `go run ./tools/schemadoc`; otherwise record those checks as `not_assessed` in `specs/022-post-merge-governance-closure/reviews/final-evidence.md`.
- [ ] T038 Review `docs/closure-decision-ledger.md`, `docs/spec-reality-ledger.md`, `docs/roadmap.md`, and `specs/022-post-merge-governance-closure/spec.md` for forbidden claims about retroactive approval, production trust, release approval, or external attestation.
- [ ] T039 Run PR-ready review for spec drift, constitution drift, product drift, CRAP `< 5`, MI `> 70` / assessed baseline status, Clean Architecture hex boundaries, Clean Code, SOLID, DRY, and YAGNI; record unavailable planes as `not_assessed` or `cannot_verify`.
- [ ] T040 Ensure new closure-state claims use `sdp-trace-claim` tags where authoritative claim syntax is supported by `docs/claim-authoring.md`; otherwise record why the claim remains prose context.
- [ ] T041 Run iterative adversarial full-diff review through the configured PR review lanes until retained findings are resolved or explicitly dispositioned.
- [ ] T042 Update `specs/022-post-merge-governance-closure/tasks.md` checkboxes only after the matching evidence edits, review artifacts, and local verification exist.

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies.
- **Foundational (Phase 2)**: Depends on Setup completion and blocks all user
  stories.
- **Pre-Implementation Review Gate (Phase 2A)**: Depends on Foundational and
  blocks all closure surface edits.
- **US1 Governance Evidence Summary (Phase 3)**: Depends on Foundational and
  the Pre-Implementation Review Gate.
- **US2 Maintainer Decision And Remediation Disposition (Phase 4)**: Depends on
  Foundational and the Pre-Implementation Review Gate; can start after Phase 2A
  and should use US1 evidence when possible.
- **US3 Synchronized Closure Navigation (Phase 5)**: Depends on US1 and US2.
- **Polish (Phase 6)**: Depends on desired story edits and all focused story
  reviews.

### User Story Dependencies

- **US1 (P1)**: First independently valuable slice; produces evidence summary.
- **US2 (P2)**: Uses the decision ledger and remediation state; can partially
  proceed in parallel with US1 only after evidence sources are mapped and the
  Pre-Implementation Review Gate is complete.
- **US3 (P3)**: Must run after US1 and US2 so navigation reflects final state.

### Parallel Opportunities

- T005 and T006 can run in parallel after T004.
- T016 and T017 can run in parallel after the pre-implementation review gate.
- T022 and T023 can run in parallel after Phase 2A.
- T028 and T029 can run in parallel after US1 and US2 are complete.

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
live state remains `not_assessed` or `cannot_verify`; when GitHub access is
available, PR #60 and PR #63 live refresh is required.
