# Tasks: Spec Governance And Roadmap Navigation

**Input**: `spec.md`, `plan.md`
**Tests**: doccheck, roadmap consistency review, `git diff --check`.

## Phase 0 - Review

- [x] T001 Run PI CTO/product review on proposed roadmap shape. — Completed: Socratic review via claim-doubt-cycle.
- [x] T002 Run PI DX review on lifecycle labels. — Completed: 8 findings recorded.
- [x] T003 Record review dispositions. — `review-disposition.md` committed.

## Phase 1 - Roadmap

- [x] T010 Define spec lifecycle taxonomy. — In spec.md US-002; caveat added.
- [ ] T011 Add current spec ownership map. → Create `docs/roadmap.md`.
- [ ] T012 Mark blockers and next steps for active specs. → Populate in roadmap.
- [ ] T013 Update spec.md with approved review fixes. → Status, claim verification, historical rule, claim-tag scope, freshness.

## Phase 2 - Governance

- [ ] T020 Define task-file expectations for approval gates and blocked specs. → Add to `docs/roadmap.md` or new `docs/spec-lifecycle.md`.
- [x] T021 Define claim-tag enforcement scope for future authoritative prose. — In spec.md US-004; scope = new/touched files.
- [ ] T022 Link roadmap from docs map or README if approved.

## Phase 3 - Closure

- [ ] T030 Run doccheck and full verification.
- [ ] T031 Run multi-LLM review on roadmap accuracy (GLM/Qwen/DeepSeek).
- [ ] T032 Record any historical migration follow-ups separately.
- [ ] T033 Commit scoped slices and prepare PR evidence.
