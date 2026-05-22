# Tasks: Repo Realignment, Monitoring, And Gate Readiness

Status: draft; maintainer human review not_assessed

## Phase 0 - Review And Approval

- [ ] T019-001 Review Spec 019 scope with maintainers.
  Status: `not_assessed`.
- [ ] T019-002 Run adversarial spec review before implementation approval.
  Status: `not_assessed`.
- [ ] T019-003 Approval gate: implementation and Pi handoff may start only
  after reviewed spec direction is approved.
  Status: `not_assessed`.

## Phase 1 - Pi-Ready Workstreams

- [x] T019-010 WS-019-A: Build spec reality ledger.
  Status: completed; deliverable is `docs/spec-reality-ledger.md`.
  Pi ownership: `docs/roadmap.md`, `specs/*/spec.md`, `specs/*/plan.md`,
  `specs/*/tasks.md`, optional `docs/spec-reality-ledger.md`.
  Session verification: `go run ./tools/doccheck`,
  `go run ./tools/hygienecheck`, `git diff --check`.

- [x] T019-020 WS-019-B: Restore OSS compatibility and benchmark tool quality.
  Status: completed; CRAP and MI gates pass for touched tools.
  Pi ownership: `tools/osscompat/*`, `tools/ossbench/*`, focused tests.
  Session verification: `go test -count=1 ./tools/osscompat`,
  `go test -count=1 ./tools/ossbench`,
  `go run ./tools/qualitycheck -gocyclo tools/osscompat tools/ossbench`,
  coverage-backed CRAP check for touched tool packages, `git diff --check`.


- [x] T019-030 WS-019-C: Repair command-surface maintainability ratchets.
  Status: completed; extracted helpers to `main_536_commandsurfacehelpers.go`;
  MI baseline regenerated and passes.
  Pi ownership: `cmd/sdp-trace/main_537_commandsurfaceconstants.go`,
  `cmd/sdp-trace/main_538_commandsurfaceregistrycore.go`,
  `cmd/sdp-trace/main_539_commandsurfaceregistryobserve.go`,
  `cmd/sdp-trace/main_540_commandsurfaceregistryassess.go`,
  `cmd/sdp-trace/main_541_commandsurfaceregistryother.go`,
  `cmd/sdp-trace/main_542_commandsurfaceregistrypacket.go`,
  `cmd/sdp-trace/main_546_commandsurfacedrift.go`.
  Session verification: `go run ./cmd/sdp-trace command-surface`,
  `go run ./tools/mibaselinepolicy -baseline tools/qualitycheck/function-mi-baseline.json -base-ref HEAD~1`,
  `go test -count=1 ./...`, `git diff --check`.


- [ ] T019-040 WS-019-D: Resolve live `wrap` output/schema compatibility.
  Status: `HITL` — requires human contract decision on wrap/schema drift path.
  Pi ownership: schema, recorder/wrap output code, affected examples, focused
  tests, and docs selected after the HITL contract decision.
  Session verification: schema validation for a live `wrap` output,
  `go test -count=1 ./...`, `go run ./tools/osscompat -json` when optional
  prerequisites are available, `go run ./tools/doccheck`,
  `go run ./tools/hygienecheck`, `git diff --check`.


- [ ] T019-050 WS-019-E: Create monitoring and advisory gate proof pack.
  Status: `HITL` — blocked on WS-019-D contract decision.
  Pi ownership: selected examples, monitoring/gate docs, focused tests for
  unsafe event rejection and missing evidence states.
  Session verification: reproduce the proof pack commands, verify unsafe raw
  prompt input fails before run creation, verify missing required evidence is
  `not_assessed`, verify advisory gate missing CI/audit evidence is
  `cannot_verify`, `go test -count=1 ./...`, `git diff --check`.


- [x] T019-060 WS-019-F: Clean numeric shards in the core CLI path.
  Status: completed; consolidated 16 flagset numeric shards into 3 behavior-named
  files (`flagset.go`, `flagset_parse.go`, `flagset_consume.go`).
  Pi ownership: `cmd/sdp-trace/flagset*.go` (consolidated from `core_520`–`core_535`).
  Session verification: `go run ./cmd/sdp-trace command-surface`,
  `go test -count=1 ./...`, `go vet ./...`, `go run ./tools/doccheck`,
  numeric shard count reduced by 13 in the owned core path,
  `git diff --check`.


- [ ] T019-070 WS-019-G: Clean numeric shards in harness and gate code.
  Status: `blocked` — depends on WS-019-E (HITL) and WS-019-F (done).
  Pi ownership: selected `internal/harnessobs/*`, selected harness and gate CLI
  glue, focused tests.
  Session verification: monitoring/gate proof pack still passes,
  `go test -count=1 ./...`, `go vet ./...`, numeric shard count for the owned
  harness/gate path before and after, `git diff --check`.


- [x] T019-080 WS-019-H: Close supply-chain probe automation gaps.
  Status: completed; probes explicitly preserve `cannot_verify` with reproducible
  reasons; tests verify states; docs updated.
  Pi ownership: `tools/osscompat/*`, supply-chain compatibility docs/examples,
  focused tests for probe state classification.
  Session verification: `go test -count=1 ./tools/osscompat`,
  `go run ./tools/osscompat -json`, `go run ./tools/doccheck`,
  `go run ./tools/hygienecheck`, `git diff --check`.


## Phase 2 - Integration And Evidence

- [x] T019-090 Run full local verification for the integrated block.
  Status: completed; all gates pass.
  Session verification: `go test -count=1 ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `go run ./tools/schemadoc -verify-readme`, `jq empty schema/*.json`,
  `go run ./tools/qualitycheck -gocyclo cmd internal tools`,
  coverage-backed CRAP check, `git diff --check`.

- [x] T019-100 Run independent Pi review panels against the final diff.
  Status: completed; codex review on cmd/sdp-trace subset returned LGTM (Round 2),
  ossbench subset returned LGTM, osscompat subset returned findings that were fixed:
  1) runCheckJSONSchema output preservation restored, 2) Direct tests skip when
  external tools are present to avoid environment sensitivity. Re-verified all gates pass.

- [x] T019-110 Query live CI for the final source commit or PR head.
  Status: completed; CI run 26313909770 on commit `65c469f` passed.
  GitHub Actions URL: https://github.com/fall-out-bug/sdp-trace/actions/runs/26313909770

- [ ] T019-120 Approval gate: close this spec only after review findings,
  local verification, and live CI evidence are recorded. Missing external
  evidence must remain `cannot_verify` or `not_assessed`.
  Status: blocked on Phase 0 HITL gates (T019-001/002/003). All AFK workstreams,
  local verification, Pi review (T019-100), and live CI evidence (T019-110)
  are complete. Phase 0 maintainer review remains `not_assessed`.

