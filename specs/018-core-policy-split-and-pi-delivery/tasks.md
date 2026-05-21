# Tasks: Core/Policy Split And Pi Delivery

Status: accepted for implementation; machine review complete, maintainer human review not_assessed

## Phase 0 - Review

- [ ] T018-001 Review core command list with maintainers.
  Status: `not_assessed` for this controlled-pilot block.
- [x] T018-001a Machine review of core command list against `sdp-trace command-surface`.
  State: `pass`. Verified: `go run ./cmd/sdp-trace command-surface`,
  `go run ./tools/doccheck`.
- [x] T018-002 Verify docs and implementation slices preserve the spec
  decision that extension surfaces stay in the same binary for this phase.
  State: `pass`. Verified: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`.

## Phase 1 - Pi-Ready Workstreams

- [x] T018-010 WS-018-A: Create command stability matrix.
  Pi ownership: `docs/command-stability-matrix.md`.
  State: `pass`. Verified: `go run ./cmd/sdp-trace command-surface`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.
- [x] T018-020 WS-018-B: Rewrite first-run docs around core-first adoption.
  Pi ownership: `README.md`, `docs/install.md`,
  `docs/contributor-quickstart.md`, `docs/adoption-guide.en.md`,
  `docs/adoption-guide.ru.md`.
  State: `pass`. Verified: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.
- [x] T018-030 WS-018-C: Create package ownership map.
  Pi ownership: `docs/package-ownership-map.md`; optional generated inventory
  under `tools/` if implemented in Go.
  State: `pass`. Verified: `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/...`,
  `go test -count=1 ./...`, `go vet ./...`.
- [x] T018-040 WS-018-D: Create extension/deprecation boundary plan.
  Pi ownership: `docs/extension-boundary-plan.md`.
  State: `pass`. Verified: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.
- [x] T018-045 WS-018-E: Define source file locality policy for numbered Go
  shards.
  Pi ownership: `docs/package-ownership-map.md`,
  `specs/018-core-policy-split-and-pi-delivery/spec.md`.
  State: `pass`. Verified: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.

## Phase 2 - Integration

- [x] T018-050 Verify command docs still match `sdp-trace command-surface`.
  State: `pass`. Verified: `go run ./cmd/sdp-trace command-surface`,
  `go run ./tools/doccheck`, `go test -count=1 ./...`, `go vet ./...`,
  `test -z "$(gofmt -l cmd/sdp-trace/doctor_515_usagetext.go cmd/sdp-trace/main_540_commandsurfaceregistryassess.go cmd/sdp-trace/main_541_commandsurfaceregistryother.go)"`,
  `git diff --check`.
- [x] T018-060 Run docs and hygiene checks.
  State: `pass`. Verified: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.
- [ ] T018-070 Prepare follow-up implementation specs only after core/extension
  direction is approved. Status: `not_assessed` until follow-up specs are filed.
