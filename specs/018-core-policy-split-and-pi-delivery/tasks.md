# Tasks: Core/Policy Split And Pi Delivery

Status: draft

## Phase 0 - Review

- [ ] T018-001 Review core command list with maintainers.
- [ ] T018-002 Decide whether extension surfaces stay in the same binary,
  become sub-binaries, or remain documented as controlled-pilot commands.

## Phase 1 - Pi-Ready Workstreams

- [ ] T018-010 WS-018-A: Create command stability matrix.
  Pi ownership: `docs/command-stability-matrix.md`.
  Verification: `go run ./cmd/sdp-trace command-surface`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.
- [ ] T018-020 WS-018-B: Rewrite first-run docs around core-first adoption.
  Pi ownership: `README.md`, `docs/install.md`,
  `docs/contributor-quickstart.md`, `docs/adoption-guide.en.md`,
  `docs/adoption-guide.ru.md`.
  Verification: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.
- [ ] T018-030 WS-018-C: Create package ownership map.
  Pi ownership: `docs/package-ownership-map.md`; optional generated inventory
  under `tools/` if implemented in Go.
  Verification: `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/...`,
  `go test -count=1 ./...`, `go vet ./...`.
- [ ] T018-040 WS-018-D: Create extension/deprecation boundary plan.
  Pi ownership: `docs/extension-boundary-plan.md`.
  Verification: `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `git diff --check`.

## Phase 2 - Integration

- [ ] T018-050 Verify command docs still match `sdp-trace command-surface`.
- [ ] T018-060 Run docs and hygiene checks.
- [ ] T018-070 Prepare follow-up implementation specs only after core/extension
  direction is approved.
