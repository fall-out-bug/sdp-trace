# Tasks: Repo Hygiene And Artifact Boundary

**Input**: `spec.md`, `plan.md`
**Tests**: Hygiene check, full Go tests, doccheck, `git diff --check`.

## Phase 0 - Review

- [x] T001 Run PI DX review on artifact classification.
- [x] T002 Run PI evidence review on whether moved/removed review artifacts remain discoverable.
- [x] T003 Record accepted/rejected findings before edits.

## Phase 1 - Policy And Guardrails

- [x] T010 Define durable vs local vs PR-only artifact classes.
- [x] T011 Add or extend a Go hygiene check for forbidden checked-in artifacts.
- [x] T012 Update ignore rules for local worktrees and generated snapshots.

## Phase 2 - Cleanup

- [x] T020 Move or remove root `PR_DESCRIPTION.md`, `design-note.md`, and root `reviews/` artifacts.
- [x] T021 Update docs to name the canonical review evidence location.
- [x] T022 Verify no checked-in absolute local paths remain in durable docs.

## Phase 3 - Closure

- [x] T030 Run full verification.
- [x] T031 Run review on cleanup diff.
- [x] T032 Commit as a scoped repo hygiene slice.
