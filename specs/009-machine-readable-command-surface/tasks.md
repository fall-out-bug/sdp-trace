# Tasks: Machine-Readable Command Surface

**Input**: `spec.md`, `plan.md`  
**Tests**: No implementation before PI spec review and explicit approval. Later implementation must run full tests, doccheck, JSON syntax checks for committed examples, quality gates, and `git diff --check`.

## Phase 0 - PI Review

- [ ] T001 Run PI requirements-vs-implementation review.
- [ ] T002 Run PI DX/UX review.
- [ ] T003 Run PI code architecture review.
- [ ] T004 Run PI trust/compatibility overclaim review.
- [ ] T005 Record findings and dispositions before implementation.

## Phase 1 - Source Of Truth

- [ ] T010 Choose registry/checker/generator shape.
- [ ] T011 Model command paths, flags, repeated flags, rest payloads, outputs, profiles, and partial metadata states.
- [ ] T012 Add schema version to machine-readable output.

## Phase 2 - Drift Checks

- [ ] T020 Add tests for deterministic command-surface output.
- [ ] T021 Extend `doccheck` or equivalent Go check to compare docs/help with the command registry.
- [ ] T022 Add or update committed JSON examples only if they are validated.

## Phase 3 - Closure

- [ ] T030 Run full local verification.
- [ ] T031 Run PI implementation review planes.
- [ ] T032 Record advisory follow-ups separately from blockers.
