# Tasks: sdp-trace Time-Series Evidence Substrate

**Input**: Design documents from `/specs/001-sdp-trace-time-series-evidence-substrate/`
**Prerequisites**: `spec.md`, `plan.md`, `research.md`, `data-model.md`, `contracts/sdp-trace-sdp-gate-boundary.md`
**Tests**: Include schema syntax checks now; add JSON Schema validation after validator selection.

**Organization**: Tasks are grouped by user story to preserve independent value and reviewability.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it touches different files or produces an independent artifact.
- **[Story]**: User story from `spec.md`.
- **Beads mirror**: Optional secondary tracking issue. Beads does not replace this task list.

## Phase 1: Setup and Canonical SpecKit Package

**Purpose**: Make SpecKit artifacts the repo-observable planning source.

- [x] T001 [US4] Add root README pointer to `specs/001-sdp-trace-time-series-evidence-substrate/spec.md`
- [x] T002 [US4] Link Beads epic `sdp-trace-cdn` and child issues to this spec with `bd update --spec-id`
- [ ] T003 [US4] Document in `docs/speckit-compatibility.md` that Beads is secondary execution tracking, not the planning source of truth

**Checkpoint**: A repository observer can start from `specs/001-sdp-trace-time-series-evidence-substrate/` without Beads.

## Phase 2: Foundational Boundary and Extraction

**Purpose**: Prevent policy/runtime coupling before schema work starts.

- [x] T004 [US2] Write source-mapped extraction memo in `specs/001-sdp-trace-time-series-evidence-substrate/research.md` from `sdp_lab` sources (Beads mirror: `sdp-trace-cdn.1`)
- [ ] T005 [US2] Finalize `contracts/sdp-trace-sdp-gate-boundary.md` and update `docs/concepts.md` with the same boundary (Beads mirror: `sdp-trace-cdn.2`)
- [ ] T006 [US2] Audit `README.md`, `docs/cto-brief.en.md`, `docs/cto-brief.ru.md`, `docs/team-lead-playbook.en.md`, and `docs/team-lead-playbook.ru.md` for language implying `sdp-trace` owns policy decisions (Beads mirror: `sdp-trace-cdn.11`)
- [ ] T007 [US2] Replace or narrow gate/decision wording so external verdicts are recorded inputs, not `sdp-trace` decisions (Beads mirror: `sdp-trace-cdn.11`)

**Checkpoint**: `sdp-trace` and `sdp-gate` ownership is clear before new schemas are added.

## Phase 3: User Story 1 - CTO Reviews Process Movement (Priority: P1)

**Goal**: Define time-series observations and metric streams without built-in degradation policy.

**Independent Test**: A sample metric stream can show movement across windows while every sample has evidence or `not_assessed`.

### Contract and Schema Tasks

- [ ] T008 [P] [US1] Design `schema/observation.schema.json` for evidence-backed observations (Beads mirror: `sdp-trace-cdn.3`)
- [ ] T009 [P] [US1] Design `schema/metric-stream.schema.json` for metric samples and streams (Beads mirror: `sdp-trace-cdn.3`)
- [ ] T010 [US1] Add examples under `examples/github-speckit/` showing current-window vs previous-window movement without policy verdicts
- [ ] T011 [US1] Define metric catalog in `docs/process-metric-catalog.md` with units, dimensions, evidence source, and `not_assessed` rule (Beads mirror: `sdp-trace-cdn.5`)
- [ ] T012 [US1] Update `schema/trace.schema.json` or document a replacement path so trace snapshots can include observations and metric samples

### Verification Tasks

- [ ] T013 [US1] Run `jq empty schema/*.json`
- [ ] T014 [US1] Record validation output in a sanitized evidence note under `docs/research/`

**Checkpoint**: CTO-facing process movement exists as data, not as a policy verdict.

## Phase 4: User Story 2 - sdp-gate Inherits Trace Contracts (Priority: P1)

**Goal**: Produce the assessment input contract consumed by `sdp-gate`.

**Independent Test**: An assessment input example contains evidence, observations, metric streams, and `not_assessed`, but no pass/fail/degradation decision.

- [ ] T015 [P] [US2] Design `schema/evidence-event.schema.json` from portable evidence-event concepts (Beads mirror: `sdp-trace-cdn.4`)
- [ ] T016 [P] [US2] Design `schema/provenance-record.schema.json` for actor/model/harness/tool provenance (Beads mirror: `sdp-trace-cdn.4`)
- [ ] T017 [US2] Design `schema/assessment-input.schema.json` for policy-engine handoff
- [ ] T018 [US2] Add `examples/go-service/assessment-input.json` or equivalent portable example
- [ ] T019 [US2] Update `schema/README.md` with ownership and validation rules

**Checkpoint**: `sdp-gate` has a clear inherited input contract.

## Phase 5: User Story 3 - Pilot Evaluates Harness, Model, and JVM Stack Slices (Priority: P1)

**Goal**: Create repeatable pilot run-cards and evidence paths for the customer-requested matrix.

**Independent Test**: Each run-card lists prompt, expected artifacts, provenance fields, unsupported claims, validation, and `not_assessed` behavior.

- [ ] T020 [P] [US3] Add OpenCode run-card covering MiniMax, Kimi, and GLM in `docs/research/opencode-model-run-card.md` (Beads mirror: `sdp-trace-cdn.6`)
- [ ] T021 [P] [US3] Add harness run-card for Superpowers, `gsd`, `gsd2`, and Oh My OpenAgent in `docs/research/harness-run-card.md` (Beads mirror: `sdp-trace-cdn.6`)
- [ ] T022 [US3] Add Kotlin+Bazel pilot fixture plan in `docs/research/kotlin-bazel-fixture-plan.md` (Beads mirror: `sdp-trace-cdn.7`)
- [ ] T023 [US3] Update `docs/jvm-bazel-guide.md` with Kotlin+Bazel-specific evidence requirements
- [ ] T024 [US3] Add or update `examples/jvm-bazel/` with a Kotlin+Bazel evidence bundle or fixture placeholder that is explicitly `not_assessed` until run
- [ ] T025 [US3] Update `docs/harness-compatibility-matrix.md` only with evidence-backed status or `TBD`/`not_assessed` (Beads mirror: `sdp-trace-cdn.10`)
- [ ] T026 [US3] Update `docs/model-compatibility.md` only with evidence-backed status or `TBD`/`not_assessed` (Beads mirror: `sdp-trace-cdn.10`)

**Checkpoint**: Pilot scope is executable without unsupported compatibility claims.

## Phase 6: User Story 4 - Repository Observer Finds SpecKit Evidence (Priority: P2)

**Goal**: Make the evidence package self-explanatory from committed files.

**Independent Test**: A reviewer can follow `quickstart.md`, validate schemas, and identify what remains `not_assessed`.

- [ ] T027 [US4] Select and document JSON Schema validator strategy in `schema/README.md` (Beads mirror: `sdp-trace-cdn.8`)
- [ ] T028 [US4] Add pass and `not_assessed` fixtures for new schemas under `examples/`
- [ ] T029 [US4] Add validation command that excludes `.git`, `.beads`, `.sdp-trace-runs`, and `benchmarks/repos/` (Beads mirror: `sdp-trace-cdn.8`)
- [ ] T030 [US4] Build customer pilot evidence package outline in `docs/research/customer-pilot-evidence-package.md` (Beads mirror: `sdp-trace-cdn.9`)
- [ ] T031 [US4] Verify `jq empty schema/*.json`
- [ ] T032 [US4] Verify all committed examples are parseable JSON where applicable

**Checkpoint**: The repository itself explains the plan, proof, gaps, and execution path.

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1**: Start immediately.
- **Phase 2**: Depends on Phase 1.
- **Phase 3**: Depends on boundary clarity from Phase 2.
- **Phase 4**: Depends on Phase 2 and can run partly in parallel with Phase 3.
- **Phase 5**: Depends on evidence/provenance contract direction from Phase 4.
- **Phase 6**: Depends on schema and pilot artifacts from Phases 3-5.

### Parallel Opportunities

- T008 and T009 can run in parallel.
- T015 and T016 can run in parallel.
- T020 and T021 can run in parallel.
- T025 and T026 can run in parallel after pilot evidence exists.

## Implementation Strategy

### MVP First

1. Complete Phase 1.
2. Complete Phase 2.
3. Complete T008, T009, T015, T016, and T017.
4. Add one valid example and one `not_assessed` example.
5. Stop and review the contract before running the full pilot matrix.

### Evidence Discipline

- Do not claim harness/model/stack compatibility without committed evidence.
- Do not add policy thresholds to `sdp-trace`.
- Keep raw pilot outputs ignored until sanitized.
- Every public claim must link to a file, command, example, or `not_assessed` entry.
