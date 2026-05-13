# Tasks: Machine-Readable Command Surface

**Input**: `spec.md`, `plan.md`
**Tests**: No implementation before PI spec review and explicit approval. Later implementation must run full tests, doccheck, JSON syntax checks for committed examples, quality gates, and `git diff --check`.

## Phase 0 - PI Review

- [x] T001 Run PI requirements-vs-implementation review.
- [x] T002 Run PI DX/UX review.
- [x] T003 Run PI code architecture review.
- [x] T004 Run PI trust/compatibility overclaim review.
- [x] T005 Record findings and dispositions before implementation.

## Phase 1 - Source Of Truth

- [x] T010 Choose registry/checker/generator shape.
- [x] T011 Model command paths, flags, repeated flags, rest payloads, outputs, profiles, and partial metadata states.
- [x] T012 Add schema version to machine-readable output.

## Phase 2 - Drift Checks

- [x] T020 Add tests for deterministic command-surface output.
- [x] T021 Extend `doccheck` or equivalent Go check to compare docs/help with the command registry.
- [x] T022 Add or update committed JSON examples only if they are validated.

## Phase 3 - Closure

- [x] T030 Run full local verification.
- [x] T031 Run PI implementation review planes.
- [x] T032 Record advisory follow-ups separately from blockers.
