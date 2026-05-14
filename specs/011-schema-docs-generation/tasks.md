# Tasks: Schema Documentation Validation

**Input**: `spec.md`, `plan.md`
**Tests**: No implementation before PI spec review and explicit approval. Later implementation must run full tests, JSON syntax checks, schema-doc validation, doccheck, and `git diff --check`.

## Phase 0 - PI Review

- [x] T001 Run PI docs-completeness review. Evidence: `reviews/011-round1-glm-5.1.md`, `reviews/011-round2-glm-5.1.md`.
- [x] T002 Run PI DX review for agents/downstream integrators. Evidence: DX covered in round-1/round-2 reviews above.
- [x] T003 Run PI requirements-vs-implementation review. Evidence: `reviews/011-round1-minimax-2.7.md`, `reviews/011-round2-minimax-2.7.md`.
- [x] T004 Run PI trust/overclaim review. Evidence: round-2 MiniMax approves with no overclaim findings; synthesis in `reviews/011-synthesis.md`.
- [x] T005 Record findings and dispositions before implementation. Evidence: `reviews/011-synthesis.md`.

## Phase 1 - Metadata Design

- [x] T010 Choose metadata location: checked JSON index (`schema/index.json`). Rationale: machine-readable, versioned, testable, no runtime dependency on Markdown parsing.
- [x] T011 Define required fields: schema `name`, `status`, `purpose`, plus optional `example_coverage` and `examples`. Evidence: `schema/index.json`.
- [x] T012 Define how `not_assessed` example coverage is represented. Evidence: `"example_coverage": "not_assessed"` (or absent, interpreted as `not_assessed`).

## Phase 2 - Checker Or Renderer

- [x] T020 Implement the smallest Go checker or renderer. Evidence: `tools/schemadoc/main.go`.
- [x] T021 Fail missing schema entries, extra schema entries, missing status/purpose, and broken example refs. Evidence: `check()` and sub-validators in `tools/schemadoc/main.go`; 23 tests in `tools/schemadoc/main_test.go`.
- [x] T022 Keep generated output deterministic if docs are generated. Evidence: `TestGenerateTableIsDeterministic` double-call equality.

## Phase 3 - CI And Closure

- [x] T030 Wire validation into CI or existing docs checks. Evidence: `.github/workflows/ci.yml` runs `go run ./tools/schemadoc` and `go run ./tools/schemadoc -verify-readme`; `docs/ci-check-policy.md` lists both as required CI evidence.
- [x] T031 Update docs to distinguish syntax, docs freshness, and semantic validation. Evidence: `schema/README.md` Validation section.
- [x] T032 Run full local verification and PI implementation review. Evidence: round-2 GLM 5.1 (`advisory`, no blockers) + MiniMax 2.7 (`approve`); `reviews/011-synthesis.md`; local verification passes below.
