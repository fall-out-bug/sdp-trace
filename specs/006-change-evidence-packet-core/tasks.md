# Tasks: Change Evidence Packet Core

## Phase 1: Spec Review

- [ ] T001 Run Socratic spec review for product proof, evidence/forgery, and
  DX/replayability planes.
  Status: `not_assessed`; the checked-in review ledger exists but records
  pending planes only.
- [ ] T002 Record findings and dispositions in
  `reviews/2026-05-10-socratic-review.md`.
  Status: `partial`; the ledger records pending planes, not findings.
- [ ] T003 Get explicit user approval before implementation.
  Status: `not_assessed`; later implementation exists, but approval evidence is
  not represented in this spec directory.

## Phase 2: Contract

- [x] T004 Add `change-evidence-packet.v0` schema.
  Evidence: `schema/change-evidence-packet.v0.schema.json`.
- [x] T005 Add `evidence-bundle-manifest.v0` schema.
  Evidence: `schema/evidence-bundle-manifest.v0.schema.json`.
- [x] T006 Add minimal GitHub PR evidence input fixture schema or documented
  fixture contract.
  Evidence: `examples/change-evidence-packet/github-input.json` and
  `docs/change-evidence-packet.md`.
- [x] T007 Add valid and invalid schema fixtures.
  Evidence: valid checked-in fixture plus negative validator coverage in
  `cmd/sdp-trace/packet_cli_test.go` and `internal/packet/packet_test.go`.

## Phase 3: Product Code

- [x] T008 Add Go packet/bundle models.
  Evidence: `internal/packet/packet_011_packet.go` through packet model files.
- [x] T009 Add validator for required rows, allowed states, missing reasons,
  evidence refs, resolver entries, expired artifacts, and contradiction rules.
  Evidence: `internal/packet/packet_060_validate.go` and focused packet tests.
- [x] T010 Add Markdown renderer with stable golden output.
  Evidence: `internal/packet/packet_193_rendermarkdown.go` and CLI render
  tests.
- [x] T011 Add CLI validate/render surface.
  Evidence: `sdp-trace packet validate`, `check-demo`, `render`,
  `build-github`, and `build-pr` command surfaces.

## Phase 4: Product Fixtures

- [x] T012 Add happy-path fixture with `PC-THEATER: pass`.
  Evidence: `examples/change-evidence-packet/happy-path.bundle.json`.
- [x] T013 Add missing verification fixture.
  Evidence: generated fixture coverage in `internal/packet/packet_test.go`.
- [x] T014 Add expired artifact fixture.
  Evidence: negative coverage in `cmd/sdp-trace/packet_cli_test.go` and
  `internal/packet/packet_test.go`.
- [x] T015 Add contradictory evidence fixture.
  Evidence: negative and downgraded-state coverage in
  `internal/packet/packet_test.go`.
- [x] T016 Add `agent_claimed_verification` theater fixture.
  Evidence: theater reason coverage in `internal/packet/packet_test.go`.

## Phase 5: Documentation And Trace

- [x] T017 Document packet and bundle authoring contract.
  Evidence: `docs/change-evidence-packet.md` and
  `docs/evidence-bundle-manifest.md`.
- [x] T018 Document canonical artifact vs PR projection rule.
  Evidence: `docs/change-evidence-packet.md` and
  `docs/evidence-bundle-manifest.md`.
- [x] T019 Update relevant product contract references without broad OSS or
  enterprise overclaim.
  Evidence: `docs/agent-entrypoint.md`, `docs/output-location-map.md`, and
  packet docs describe non-approval limits.
- [x] T020 Add trace/evidence notes for implemented behavior.
  Evidence: packet docs and command surface notes preserve `cannot_verify` and
  `not_assessed` boundaries.

## Phase 6: Verification And Review

- [x] T021 Run `go test ./...`.
  Session verification: `go test -count=1 ./...`.
- [x] T022 Run `jq empty schema/*.json`.
  Session verification: `jq empty schema/*.json examples/change-evidence-packet/*.json`.
- [x] T023 Run fixture validator and renderer golden tests.
  Session verification: `go run ./cmd/sdp-trace packet validate --bundle examples/change-evidence-packet/happy-path.bundle.json`;
  `go run ./cmd/sdp-trace packet render --bundle examples/change-evidence-packet/happy-path.bundle.json --out "$tmp"`.
- [x] T024 Run `git diff --check`.
  Session verification: `git diff --check`.
- [ ] T025 Run implementation review planes: code/correctness,
  tracing/evidence, requirements-vs-implementation.
  Status: `not_assessed`; no current implementation-review artifact found in
  this spec directory.
- [ ] T026 Fix accepted critical/major findings and re-review.
  Status: `not_assessed`; depends on T025.
- [ ] T027 Open PR and repeat PR-level review planes before ready.
  Status: `not_assessed`; current branch has not yet opened a PR for this
  closure-route reconciliation.
