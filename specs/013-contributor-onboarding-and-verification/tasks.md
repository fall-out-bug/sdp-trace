# Tasks: Contributor Onboarding And Verification

**Input**: `spec.md`, `plan.md`
**Tests**: doccheck, smoke replay, full Go tests, `git diff --check`.

## Phase 0 - Review

- [x] T001 Run PI DX review on the proposed onboarding path.
  Evidence: `reviews/2026-05-26-closure-review.md`.
- [x] T002 Run PI UX review on command order and failure language.
  Evidence: `reviews/2026-05-26-closure-review.md`.
- [x] T003 Record review dispositions.
  Evidence: `reviews/2026-05-26-closure-review.md`.

## Phase 1 - Canonical Path

- [x] T010 Choose canonical smoke command source.
  Evidence: `docs/contributor-quickstart.md`.
- [x] T011 Add or revise contributor quick-start documentation.
  Evidence: `docs/contributor-quickstart.md`.
- [x] T012 Link it from README and docs map.
  Evidence: `README.md`, `docs/README.md`, and `docs/agent-onboarding.md`.

## Phase 2 - Drift Control

- [x] T020 Remove duplicate smoke blocks or make them references.
  Evidence: `docs/contributor-quickstart.md`, `docs/install.md`,
  `docs/reviewer-entrypoint.md`, and `docs/agent-onboarding.md`.
- [x] T021 Extend doccheck for remaining canonical snippets if needed.
  Evidence: `tools/doccheck/quickstart*.go`.
- [x] T022 Add expected output/state notes for common failures.
  Evidence: `docs/contributor-quickstart.md`.

## Phase 3 - Closure

- [x] T030 Replay smoke path locally.
  Session verification: replayed the `docs/contributor-quickstart.md`
  canonical smoke path from the repository root.
- [x] T031 Run full verification.
  Session verification: `go test -count=1 ./...`, `go run ./tools/doccheck`,
  `go run ./tools/hygienecheck`, `git diff --check`.
- [x] T032 Run cold-reader review.
  Evidence: `reviews/2026-05-26-closure-review.md`.
