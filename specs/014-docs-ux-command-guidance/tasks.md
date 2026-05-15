# Tasks: Docs UX And Command Guidance

**Input**: `spec.md`, `plan.md`
**Tests**: doccheck, state-term grep, cold-reader review, `git diff --check`.

## Phase 0 - Review

- [x] T001 Run PI UX review on current docs flow.
- [x] T002 Verify current state/profile names against command docs.
- [x] T003 Record accepted/rejected findings.

## Phase 1 - Canonical Guidance

- [x] T010 Restructure `docs/reviewer-entrypoint.md` so a short task path (Quick Reference or task table) appears before any long flat command list.
- [x] T011 Define canonical state contract: result states (`observed`, `pass`, `fail`, `not_assessed`, `cannot_verify`) with exit codes; classify `missing_telemetry` as telemetry label, `warn` as concepts-only sub-state, `coverage_satisfied/partial/unresolved` as pr-review sub-states, `not_integrated`/`unsupported` as integration labels. Add the contract to the canonical doc.
- [x] T012 Add output location map: single table mapping command family → default output path → format → purpose → trust boundary.
- [x] T013 Add profile decision tree mapping trust profile IDs ↔ assessment profiles ↔ witness kinds ↔ authority scopes in one decision aid.
- [x] T014 Create `docs/overclaim-checklist.md` as the canonical overclaim checklist; ensure all other docs link to it.

## Phase 2 - Drift Control

- [x] T020 Remove duplicate overclaim prose from `agent-entrypoint.md`, `agent-onboarding.md`, `adoption-guide.en.md`, replacing with links to `docs/overclaim-checklist.md` plus one-line inline summaries where critical.
- [x] T021 Extend doccheck if command/state claims remain duplicated. (Not required: doccheck passes without extension.)
- [x] T022 Verify docs do not introduce new authority claims.
- [x] T023 Grep audit: scan all `.md` files for state-like tokens; flag any result-state token outside the canonical contract for classification or removal.

## Phase 3 - Closure

- [x] T030 Run doccheck, `go test ./tools/doccheck`, and full verification. Grep audit must show zero orphan result-state tokens.
- [x] T031 Run cold-reader UX review (DeepSeek, MiniMax, Qwen). See `reviews/impl-synthesis.md`.
- [x] T032 Record remaining advisory UX follow-ups separately. (No advisory follow-ups remain; all critical/major findings fixed.)
