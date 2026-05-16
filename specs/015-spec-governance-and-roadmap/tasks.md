# Tasks: Spec Governance And Roadmap Navigation

**Input**: `spec.md`, `plan.md`
**Tests**: doccheck, roadmap consistency review, `git diff --check`.

## Phase 0 - Review

- [x] T001 Run PI CTO/product review on proposed roadmap shape. — Completed: Socratic review via claim-doubt-cycle.
- [x] T002 Run PI DX review on lifecycle labels. — Completed: 8 findings recorded.
- [x] T003 Record review dispositions. — `review-disposition.md` committed.

## Phase 1 - Roadmap

- [x] T010 Define spec lifecycle taxonomy. — In spec.md US-002; caveat added.
- [x] T011 Add current spec ownership map. — `docs/roadmap.md` created.
- [x] T012 Mark blockers and next steps for active specs. — Populated in roadmap.
- [x] T013 Update spec.md with approved review fixes. — Status, claim verification, historical rule, claim-tag scope, freshness.

## Phase 2 - Governance

- [x] T020 Define task-file expectations for approval gates and blocked specs. — Added to `docs/roadmap.md`.
- [x] T021 Define claim-tag enforcement scope for future authoritative prose. — In spec.md US-004; scope = new/touched files.
- [x] T022 Link roadmap from docs map or README if approved. — Linked from `docs/README.md`.

## Phase 3 - Closure

- [x] T030 Run doccheck and full verification. — doccheck exit=0; go test ./... pass; jq empty schema/*.json pass; git diff --check pass.
- [ ] T031 Run multi-LLM review on roadmap accuracy (GLM/Qwen/DeepSeek). — `cannot_verify`: no API credentials or endpoints available in current harness. Must run before merge if provider access available.
- [x] T032 Record any historical migration follow-ups separately. — Historical specs (001–014) exempt; no migration required.
- [x] T033 Commit scoped slices and prepare PR evidence. — Two scoped commits on branch `015-spec-governance-and-roadmap`.
