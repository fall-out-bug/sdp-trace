# Tasks: MVP Readiness Hardening

**Input**: Design documents from `/specs/004-mvp-readiness-hardening/`
**Prerequisites**: `spec.md`, `plan.md`, Socratic review completion, explicit user approval before implementation
**Tests**: No implementation before reviewed spec approval. Later implementation must run local docs/quality gates and preserve controlled-pilot trust boundaries.

## Phase 0: Spec Review Gate

- [x] T001 Run Socratic pi-review across documentation promise, MVP scope, code-quality gates, and trust-boundary risks.
- [x] T002 Record reviewer findings and dispositions in `socratic-review.md`.
- [x] T003 Fix or explicitly block every critical/major spec finding.
- [x] T004 Stop for explicit user approval before documentation, CI, or Go implementation changes.

## Phase 1: Documentation Freshness And Language

- [x] T005 Compare `docs/agent-entrypoint.md` command contract against `go run ./cmd/sdp-trace --help`.
- [x] T006 In `docs/agent-entrypoint.md`, remove or correct flags documented for `pr-review packet` that do not appear in `go run ./cmd/sdp-trace --help`; record removed flags and rationale in the implementation ledger.
- [x] T007 Bring Russian command reference to parity for current command families, or explicitly mark routing to the English canonical command contract as temporary `deferred_scope` with a translation follow-up.
- [x] T008 Rewrite first-screen controlled-pilot status in README/adoption docs so it is plain and precise before dense contract terms.
- [x] T009 Verify docs do not upgrade local checks, CI witness, release proof, or examples into external production trust.

## Phase 2: Example Surface Cleanup

- [x] T010 Inventory first-class examples linked from README/docs.
- [x] T011 Label placeholder examples such as `examples/codex`, `examples/claude-code`, `examples/opencode`, and `examples/go-service` as placeholder/retired/not-assessed, or move them out of first-class MVP evidence.
- [x] T012 Ensure each first-class example README states evidence boundary, runnable command if applicable, and current state.
- [x] T013 Re-run changed example validation or record why an example is documentation-only.

## Phase 3: Gate Baseline, Lint, And Small Hygiene Fixes

- [x] T014 Define the initial measurable CRAP, cyclomatic complexity, cognitive complexity, MI ratchet, and coverage gates before decomposition starts.
- [x] T015 Record pre-change per-package coverage baseline for all MVP-critical packages.
- [x] T016 Record pre-change per-function CRAP baseline using function coverage plus cyclomatic complexity, or explicitly mark CRAP computation `assessed_gap` and block CRAP readiness claims.
- [x] T016a Record function/file-level MI baselines and ratchets or explicitly mark absolute MI `>70` as `assessed_gap` and block MI pass claims.
- [x] T017 Fix `internal/authority/authority.go` ineffectual assignment found by `golangci-lint`; if changed package coverage is below the selected floor, add a focused test for the changed path.
- [x] T018 Fix `internal/telemetry/prometheus.go` `gosimple` append-loop findings; if changed package coverage is below the selected floor, add a focused test for the changed path.
- [x] T019 Run `golangci-lint run ./...` and record result.

## Phase 4: Complexity And CRAP Hardening

- [x] T020 Decompose `internal/posture.ValidateExportResult` to the selected maximum complexity/CRAP threshold, or explicitly block MVP readiness until decomposition is approved.
- [x] T021 Decompose `internal/harnessobs.normalizeOpenCodeRawLine` to the selected maximum complexity/CRAP threshold, or explicitly block MVP readiness until decomposition is approved.
- [x] T022 Triage other functions with cyclomatic/cognitive complexity or CRAP above the selected threshold.
- [x] T023 Run `gocyclo`, selected `gocognit`, and CRAP review gates and record result.

## Phase 5: Coverage Hardening

- [x] T024 Add focused tests for `internal/contract` primary happy/error paths, or remove it from MVP claim surface.
- [x] T025 Add focused tests for `internal/policy` primary happy/error paths, or remove it from MVP claim surface.
- [x] T026 Add focused tests for `internal/export` primary happy/error paths, or remove it from MVP claim surface.
- [x] T027 Add focused `internal/trace` tests for core run/event safety and validation paths.
- [x] T028 Add focused tests for `internal/posture` primary happy/error paths for validated export behavior.
- [x] T029 Add focused tests for refactored harness observation behavior.
- [x] T030 Add focused tests for `internal/verifier` MVP-critical happy/error paths, or remove it from MVP claim surface.
- [x] T031 Run `go test -count=1 ./... -coverprofile` and record package-level and function-level coverage deltas.

## Phase 6: CI And Verification

- [x] T032 Add stable CI steps for lint/complexity/coverage gates, or record missing CI enforcement as `assessed_gap` with a concrete follow-up.
- [x] T033 Run `go test -count=1 ./...`.
- [x] T034 Run `jq empty schema/*.json` plus changed JSON examples.
- [x] T035 Run `git diff --check HEAD`.
- [x] T036 Run docs freshness check or manual command-surface comparison and record evidence, including the exact comparison method.
- [x] T037 Verify branch/commit delta against the intake baseline commit and record any stale evidence as `assessed_gap` or `cannot_verify`.
- [x] T038 Run implementation pi-review across code/correctness, tracing/evidence, and requirements-vs-implementation planes.
- [x] T039 Fix valid review findings and rerun focused review.
- [x] T040 Open PR and run PR-level review planes before ready state.
  Evidence: PR #64 and `pr-level-closure-review-2026-05-26.md`.
- [x] T041 Record named reviewer sign-off for MVP bar conditions before ready state.
  Evidence: `pr-level-closure-review-2026-05-26.md`.
- [ ] T042 Stop before merge unless explicit merge approval is present.
  Status: `not_assessed`; explicit merge approval is not represented.
