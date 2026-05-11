# Tasks: Change Evidence Packet Core

## Phase 1: Spec Review

- [ ] T001 Run Socratic spec review for product proof, evidence/forgery, and
  DX/replayability planes.
- [ ] T002 Record findings and dispositions in
  `reviews/2026-05-10-socratic-review.md`.
- [ ] T003 Get explicit user approval before implementation.

## Phase 2: Contract

- [ ] T004 Add `change-evidence-packet.v0` schema.
- [ ] T005 Add `evidence-bundle-manifest.v0` schema.
- [ ] T006 Add minimal GitHub PR evidence input fixture schema or documented
  fixture contract.
- [ ] T007 Add valid and invalid schema fixtures.

## Phase 3: Product Code

- [ ] T008 Add Go packet/bundle models.
- [ ] T009 Add validator for required rows, allowed states, missing reasons,
  evidence refs, resolver entries, expired artifacts, and contradiction rules.
- [ ] T010 Add Markdown renderer with stable golden output.
- [ ] T011 Add CLI validate/render surface.

## Phase 4: Product Fixtures

- [ ] T012 Add happy-path fixture with `PC-THEATER: pass`.
- [ ] T013 Add missing verification fixture.
- [ ] T014 Add expired artifact fixture.
- [ ] T015 Add contradictory evidence fixture.
- [ ] T016 Add `agent_claimed_verification` theater fixture.

## Phase 5: Documentation And Trace

- [ ] T017 Document packet and bundle authoring contract.
- [ ] T018 Document canonical artifact vs PR projection rule.
- [ ] T019 Update relevant product contract references without broad OSS or
  enterprise overclaim.
- [ ] T020 Add trace/evidence notes for implemented behavior.

## Phase 6: Verification And Review

- [ ] T021 Run `go test ./...`.
- [ ] T022 Run `jq empty schema/*.json`.
- [ ] T023 Run fixture validator and renderer golden tests.
- [ ] T024 Run `git diff --check`.
- [ ] T025 Run implementation review planes: code/correctness,
  tracing/evidence, requirements-vs-implementation.
- [ ] T026 Fix accepted critical/major findings and re-review.
- [ ] T027 Open PR and repeat PR-level review planes before ready.
