# Tasks: Schema Documentation Validation

**Input**: `spec.md`, `plan.md`  
**Tests**: No implementation before PI spec review and explicit approval. Later implementation must run full tests, JSON syntax checks, schema-doc validation, doccheck, and `git diff --check`.

## Phase 0 - PI Review

- [ ] T001 Run PI docs-completeness review.
- [ ] T002 Run PI DX review for agents/downstream integrators.
- [ ] T003 Run PI requirements-vs-implementation review.
- [ ] T004 Run PI trust/overclaim review.
- [ ] T005 Record findings and dispositions before implementation.

## Phase 1 - Metadata Design

- [ ] T010 Choose metadata location: generated README section, checked index, or schema annotations.
- [ ] T011 Define required fields: schema name, status, purpose, and example coverage state.
- [ ] T012 Define how `not_assessed` example coverage is represented.

## Phase 2 - Checker Or Renderer

- [ ] T020 Implement the smallest Go checker or renderer.
- [ ] T021 Fail missing schema entries, extra schema entries, missing status/purpose, and broken example refs.
- [ ] T022 Keep generated output deterministic if docs are generated.

## Phase 3 - CI And Closure

- [ ] T030 Wire validation into CI or existing docs checks.
- [ ] T031 Update docs to distinguish syntax, docs freshness, and semantic validation.
- [ ] T032 Run full local verification and PI implementation review.
