# Tasks: Core/Policy Split And Pi Delivery

Status: reviewed for Spec 018; core/extension direction approved

## Phase 0 - Review

- [x] T018-001 Review core command list with maintainers.
  Evidence: maintainer decision `approve_core_extension_direction` recorded in
  the 2026-05-26 closure thread. Approved core commands are `wrap`, `run`,
  `verify`, `explain`, `report`, and `query --query missing-evidence`; other
  current command families remain extension, experimental, or fixture-only in
  the same binary for this phase.
- [x] T018-001a Machine review of core command list against `sdp-trace command-surface`.
  Session verification: `go run ./cmd/sdp-trace command-surface`,
  `go run ./tools/doccheck`.
- [x] T018-002 Verify docs and implementation slices preserve the spec
  decision that extension surfaces stay in the same binary for this phase.
  Session verification: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`.

## Phase 1 - Pi-Ready Workstreams

- [x] T018-010 WS-018-A: Create command stability matrix.
  Pi ownership: `docs/command-stability-matrix.md`.
  Session verification: `go run ./cmd/sdp-trace command-surface`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.
- [x] T018-020 WS-018-B: Rewrite first-run docs around core-first adoption.
  Pi ownership: `README.md`, `docs/install.md`,
  `docs/contributor-quickstart.md`, `docs/adoption-guide.en.md`,
  `docs/adoption-guide.ru.md`.
  Session verification: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.
- [x] T018-030 WS-018-C: Create package ownership map.
  Pi ownership: `docs/package-ownership-map.md`; optional generated inventory
  under `tools/` if implemented in Go.
  Session verification: `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/...`,
  `go test -count=1 ./...`, `go vet ./...`.
- [x] T018-040 WS-018-D: Create extension/deprecation boundary plan.
  Pi ownership: `docs/extension-boundary-plan.md`.
  Session verification: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.
- [x] T018-045 WS-018-E: Define source file locality policy for numbered Go
  shards.
  Pi ownership: `docs/package-ownership-map.md`,
  `specs/018-core-policy-split-and-pi-delivery/spec.md`.
  Session verification: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.

## Phase 2 - Integration

- [x] T018-050 Verify command docs still match `sdp-trace command-surface`.
  Session verification: `go run ./cmd/sdp-trace command-surface`,
  `go run ./tools/doccheck`, `go test -count=1 ./...`, `go vet ./...`,
  `unformatted=$(gofmt -l cmd/sdp-trace/doctor_515_usagetext.go cmd/sdp-trace/main_540_commandsurfaceregistryassess.go cmd/sdp-trace/main_541_commandsurfaceregistryother.go) && test -z "$unformatted"`,
  `git diff --check`.
- [x] T018-060 Run docs and hygiene checks.
  Session verification: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.
- [x] T018-070 Prepare follow-up implementation specs only after core/extension
  direction is approved.
  Evidence: follow-up Spec 020 `core-query-package-split` and Spec 021
  `source-file-locality-cleanup` are filed as draft implementation specs. They
  prepare future implementation work without opening active closure-route task
  checkboxes or approving command removal, separate binaries, production
  readiness, release approval, or external trust.
