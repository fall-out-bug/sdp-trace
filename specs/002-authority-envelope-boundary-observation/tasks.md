# Tasks: Authority Envelope Boundary Observation

**Input**: Design documents from `/specs/002-authority-envelope-boundary-observation/`
**Prerequisites**: `spec.md`, `plan.md`, `data-model.md`, Socratic review completion
**Tests**: No implementation before reviewed spec approval. Later implementation must run `go test ./...`, `jq empty schema/*.json`, changed fixture validation, and `git diff --check`.

## Phase 0: Spec Review Gate

- [x] T001 Run Socratic pi-review across product boundary, data model, and evidence-source semantics.
- [x] T002 Record reviewer findings and dispositions in `socratic-review.md`.
- [x] T003 Fix or explicitly block every critical/major spec finding.
- [x] T004 Stop for explicit approval before schema or Go implementation.

## Phase 1: Schema Contracts

- [x] T005 Add `schema/authority-envelope.schema.json`.
- [x] T006 Add `schema/observed-action.schema.json`.
- [x] T007 Add `schema/evidence-binding.schema.json`.
- [x] T008 Add `schema/authority-evaluation.schema.json`.
- [x] T009 Update schema documentation with authority envelope ownership, evidence ref schemes, binding semantics, and versioning.
- [x] T010 Run `jq empty schema/*.json`.

## Phase 2: Examples and Fixtures

- [x] T011 Add fixture matrix under `examples/authority-envelope-basic/`.
- [x] T012 Add valid `within_authority` fixture with explicit approval evidence.
- [x] T013 Add `outside_authority` denied-target mutation fixture.
- [x] T014 Add git-only fixture where mutation is observed but actor/model attribution remains `not_assessed`.
- [x] T015 Add gateway-unbound fixture where model request exists but mutation model attribution remains `not_assessed`.
- [x] T016 Add malformed-policy fixtures that produce `cannot_verify`, covering both JSON/schema-invalid envelopes and semantically invalid envelopes with conflicting allow/deny or overlapping target rules.
- [x] T017 Add failed-binding fixture where source evidence exists but declared binding fields disagree.
- [x] T018 Add stale or inaccessible evidence fixture that produces `cannot_verify`.
- [x] T019 Add feedback/review fixture proving no mutation fact is inferred without mutation evidence.
- [x] T020 Add policy-specific fixture showing the same observed action evaluated under two selected `policy_id` values with different authority states.

## Phase 3: Verifier Behavior

- [x] T021 Add Go parsing for authority envelopes and observed actions.
- [x] T022 Add path/event matching against caller-selected `policy_id`.
- [x] T023 Add envelope conflict detection for ambiguous allow/deny and overlapping target rules.
- [x] T024 Add evidence binding evaluation.
- [x] T025 Add evaluation output with `within_authority`, `outside_authority`, `not_assessed`, and `cannot_verify`.
- [x] T026 Add regression tests for git-only, harness-bound, gateway-unbound, malformed-policy, failed-binding, stale-evidence, feedback-without-mutation, and policy-specific evaluation cases.
- [x] T027 Add safety tests proving raw prompts, model outputs, credentials, and private source snippets are not printed or persisted.

## Phase 4: Documentation

- [x] T028 Add `docs/authority-envelope.md`.
- [x] T029 Update `docs/agent-entrypoint.md` only after command surface exists.
- [x] T030 Update `docs/reviewer-entrypoint.md` with reviewer interpretation rules.
- [x] T031 Document that contamination/block/merge decisions are external policy decisions.

## Phase 5: Review, PR, and Closure

- [x] T032 Run implementation pi-review across code/correctness, tracing/evidence, and requirements-vs-implementation planes.
- [x] T033 Fix valid findings and rerun focused review.
- [ ] T034 Open PR and run PR-level review planes.
- [ ] T035 Merge only after fresh CI, local verification, PR review, and post-merge verification.
