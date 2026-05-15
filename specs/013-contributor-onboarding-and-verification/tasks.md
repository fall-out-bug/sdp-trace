# Tasks: Contributor Onboarding And Verification

**Input**: `spec.md`, `plan.md`
**Tests**: doccheck, smoke replay, full Go tests, `git diff --check`.

## Phase 0 - Review

- [ ] T001 Run PI DX review on the proposed onboarding path.
- [ ] T002 Run PI UX review on command order and failure language.
- [ ] T003 Record review dispositions.

## Phase 1 - Canonical Path

- [ ] T010 Choose canonical smoke command source.
- [ ] T011 Add or revise contributor quick-start documentation.
- [ ] T012 Link it from README and docs map.

## Phase 2 - Drift Control

- [ ] T020 Remove duplicate smoke blocks or make them references.
- [ ] T021 Extend doccheck for remaining canonical snippets if needed.
- [ ] T022 Add expected output/state notes for common failures.

## Phase 3 - Closure

- [ ] T030 Replay smoke path locally.
- [ ] T031 Run full verification.
- [ ] T032 Run cold-reader review.
