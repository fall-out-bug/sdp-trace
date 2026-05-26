---

description: "Task list template for feature implementation"
---

# Tasks: [FEATURE NAME]

**Input**: Design documents from `/specs/[###-feature-name]/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Behavior changes, verifier/trust changes, schema/contract changes,
and Go tooling changes require focused tests. Docs-only changes may omit tests
only when the task records the resulting verification state as `not_assessed`
or gives a concrete reason.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **CLI/product code**: `cmd/`, `internal/`, and colocated `*_test.go`
- **Repository tools**: `tools/<tool>/` and colocated `*_test.go`
- **Schemas**: `schema/`, with matching examples/fixtures and Go structs
- **Docs/examples/specs**: `docs/`, `examples/`, and `specs/`
- Paths shown below assume Go-first `sdp-trace` structure - adjust based on
  `plan.md`

<!--
  ============================================================================
  IMPORTANT: The tasks below are SAMPLE TASKS for illustration purposes only.

  The /speckit-tasks command MUST replace these with actual tasks based on:
  - User stories from spec.md (with their priorities P1, P2, P3...)
  - Feature requirements from plan.md
  - Entities from data-model.md
  - Endpoints from contracts/

  Tasks MUST be organized by user story so each story can be:
  - Implemented independently
  - Tested independently
  - Delivered as an MVP increment

  DO NOT keep these sample tasks in the generated tasks.md file.
  ============================================================================
-->

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create project structure per implementation plan
- [ ] T002 Confirm affected Go packages, schemas, docs, examples, and specs
- [ ] T003 [P] Add or update focused failing tests for changed behavior

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

Examples of foundational tasks (adjust based on your project):

- [ ] T004 Map requirements to evidence, gate, decision, trace, and provenance
  artifacts
- [ ] T005 [P] Map result state, trust scope, and authority scope boundaries
- [ ] T006 [P] Update JSON schema or contract fixtures if behavior changes
- [ ] T007 [P] Update docs/examples that demonstrate the changed surface
- [ ] T008 Confirm no Node/npm/JS/TS tooling enters the active product path
- [ ] T009 Record `not_assessed` or `cannot_verify` for unavailable evidence
- [ ] T010 Confirm local verification, CI, review, PR readiness, merge approval,
  release approval, and production trust remain distinct

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - [Title] (Priority: P1) 🎯 MVP

**Goal**: [Brief description of what this story delivers]

**Independent Test**: [How to verify this story works on its own]

### Tests for User Story 1 (required for behavior/verifier/trust changes) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T011 [P] [US1] Go test for [behavior] in [package]/[file]_test.go
- [ ] T012 [P] [US1] Fixture or schema validation for [artifact] in [path]

### Implementation for User Story 1

- [ ] T013 [P] [US1] Update Go type or parser in internal/[package]/[file].go
- [ ] T014 [P] [US1] Update CLI/tool surface in cmd/[command]/ or tools/[tool]/
- [ ] T015 [US1] Implement [behavior] in [package]/[file].go (depends on T013, T014)
- [ ] T016 [US1] Update docs/examples/spec evidence in [path]
- [ ] T017 [US1] Preserve explicit evidence states for all unavailable or
  out-of-scope proof
- [ ] T018 [US1] Add validation and error handling

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - [Title] (Priority: P2)

**Goal**: [Brief description of what this story delivers]

**Independent Test**: [How to verify this story works on its own]

### Tests for User Story 2 (required for behavior/verifier/trust changes) ⚠️

- [ ] T019 [P] [US2] Go test for [behavior] in [package]/[file]_test.go
- [ ] T020 [P] [US2] Fixture or schema validation for [artifact] in [path]

### Implementation for User Story 2

- [ ] T021 [P] [US2] Update Go type or parser in internal/[package]/[file].go
- [ ] T022 [US2] Implement [behavior] in [package]/[file].go
- [ ] T023 [US2] Update CLI/tool or docs surface in [path]
- [ ] T024 [US2] Integrate with User Story 1 components (if needed)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - [Title] (Priority: P3)

**Goal**: [Brief description of what this story delivers]

**Independent Test**: [How to verify this story works on its own]

### Tests for User Story 3 (required for behavior/verifier/trust changes) ⚠️

- [ ] T025 [P] [US3] Go test for [behavior] in [package]/[file]_test.go
- [ ] T026 [P] [US3] Fixture or schema validation for [artifact] in [path]

### Implementation for User Story 3

- [ ] T027 [P] [US3] Update Go type or parser in internal/[package]/[file].go
- [ ] T028 [US3] Implement [behavior] in [package]/[file].go
- [ ] T029 [US3] Update CLI/tool or docs surface in [path]

**Checkpoint**: All user stories should now be independently functional

---

[Add more user story phases as needed, following the same pattern]

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] TXXX [P] Documentation updates in docs/
- [ ] TXXX Code cleanup and refactoring
- [ ] TXXX Performance optimization across all stories
- [ ] TXXX [P] Additional Go tests in affected packages
- [ ] TXXX Security hardening
- [ ] TXXX Run quickstart.md validation
- [ ] TXXX Run required verification from `AGENTS.md`
- [ ] TXXX Record `verified`, `not_assessed`, `cannot_verify`, `failed`, or
  `blocked` states for each explicit claim

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - May integrate with US1 but should be independently testable
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - May integrate with US1/US2 but should be independently testable

### Within Each User Story

- Tests for behavior/verifier/trust changes MUST be written and FAIL before
  implementation
- Types/parsers before command surfaces
- Core package behavior before CLI wiring
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Go types, fixture updates, and docs updates within a story can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together (if tests requested):
Task: "Go test for [behavior] in [package]/[file]_test.go"
Task: "Fixture/schema validation for [artifact] in [path]"

# Launch supporting artifacts for User Story 1 together:
Task: "Update Go type or parser in internal/[package]/[file].go"
Task: "Update fixture or example in examples/[name]/"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Demo or report the verified state if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Demo/report (MVP!)
3. Add User Story 2 → Test independently → Demo/report
4. Add User Story 3 → Test independently → Demo/report
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 3
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
