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

- [x] T030 Run doccheck and full verification. — doccheck exit=0 (link integrity only); go test ./... pass; jq empty schema/*.json pass; git diff --check pass; GitHub Actions verify job PASS for branch head. Roadmap coverage verified manually, not by doccheck.
- [x] T031 Run multi-LLM review on roadmap accuracy. — GLM-5.1 (architecture doubt) and Qwen-3.6 (wide-context) completed with findings. MiniMax-M2.7 returned 404 (endpoint down). DeepSeek-v4 timed out without output. Review synthesis recorded in `reviews/`.
- [x] T032 Record any historical migration follow-ups separately. — Historical specs (001–014) exempt; no migration required.
- [x] T033 Commit scoped slices and prepare PR evidence. — Multiple scoped commits on branch `015-spec-governance-and-roadmap`.
- [x] T034 Fix GLM+Qwen review findings. — Roadmap statuses aligned with spec.md; `blocks/` overclaim removed; `historical` added to taxonomy; claim tags added; overclaims removed.
- [x] T035 Add claim tags to new/modified files. — Tags added to `docs/roadmap.md` and `final-evidence-map.md`.

## Phase 4 - Status Discipline

- [x] T036 Split roadmap discipline across spec, task, implementation, review,
  merge, and trust axes. — Added `docs/spec-status-discipline.md` and updated
  roadmap navigation so implemented specs are not hidden behind `draft`.

<!-- sdp-trace-claim: claim=task_closed; subject=015-tasks; state=pass; profile=repo_baseline_structural; evidence=command_set:block015-t030 -->
