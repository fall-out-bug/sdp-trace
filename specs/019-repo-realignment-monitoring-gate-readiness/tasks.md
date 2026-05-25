# Tasks: Repo Realignment, Monitoring, And Gate Readiness

Status: post-merge partial; maintainer human review not_assessed

## Phase 0 - Review And Approval

- [ ] T019-001 Review Spec 019 scope with maintainers.
  Status: `missed_pre_merge_gate`. PR #60 merged before this approval was
  recorded; this cannot be retroactively satisfied as a pre-implementation
  review.
- [ ] T019-002 Run adversarial spec review before implementation approval.
  Status: `partial_after_merge`. Cross-model review ran for implemented slices,
  but GitHub PR #60 has no recorded approval and the PR body checklist was
  unchecked at merge time.
- [ ] T019-003 Approval gate: implementation and Pi handoff may start only
  after reviewed spec direction is approved.
  Status: `missed_pre_merge_gate`. Treat the missed gate as process evidence,
  not as a remaining task that can still be completed before implementation.
- [ ] T019-004 Post-merge approval gate: maintainers must either approve the
  partial merge state and the closure plan or explicitly reject/split the
  remaining work.
  Status: `not_assessed`. Closure plan:
  `specs/019-repo-realignment-monitoring-gate-readiness/post-merge-closure-plan.md`.

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


- [x] T019-040 WS-019-D: Resolve live `wrap` output/schema compatibility.
  Status: completed; live `wrap` output has a dedicated current-run manifest
  schema at `schema/run-manifest.schema.json`. The richer
  `flight-recorder-run.schema.json` remains a separate profile artifact, not
  the live `wrap` manifest contract.
  Pi ownership: schema, recorder/wrap output contract docs, affected examples,
  focused tests, and OSS compatibility probe docs.
  Session verification: schema validation for a live `wrap` output,
  `go test -count=1 ./...`, `go run ./tools/osscompat -json` when optional
  prerequisites are available, `go run ./tools/doccheck`,
  `go run ./tools/hygienecheck`, `git diff --check`.


- [x] T019-050 WS-019-E: Create monitoring and advisory gate proof pack.
  Status: completed; deliverable is
  `examples/spec019-monitoring-gate-proof/` plus focused CLI replay tests in
  `cmd/sdp-trace/spec019_proof_pack_cli_test.go`.
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


- [x] T019-070 WS-019-G: Clean numeric shards in harness and gate code.
  Status: completed; selected gate/report CLI shards and foundational
  `internal/harnessobs` type shards were renamed to behavior-named files.
  Scoped numeric shard count for owned harness/gate paths moved from 463 to
  435.
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
  Status: completed; cross-model adversarial review performed in two rounds.
  **Round 1**: via `pi` CLI with direct providers (non-OpenRouter):
  - **GLM-5.1** (zai/glm-5.1) — Architecture plane: **LGTM** (zero material findings).
  - **Qwen-3.6-Max** (qwen/qwen3.6-max-preview) — Code/correctness plane: **LGTM**
    (2 advisory style findings, no material issues).
  - **GPT-5.5** (openai-codex/gpt-5.5) — Spec alignment plane: **5 material findings**
    (ledger formatting, stale reconciliation state, approval gate contradiction,
    missing durable evidence, harness-specific config concern).
  **Round 2** (OmPi instrumentation, post-fix): via Oh My Pi `task` tool with bundled
  reviewer agent on 4 focused planes:
  - **cmd/sdp-trace/** — **LGTM**
  - **tools/osscompat/** — **LGTM**
  - **tools/ossbench/** — **LGTM**
  - **docs/config/** — **3 findings addressed** (recursive workflow, governance blind
    spot, broken worktree isolation).
  All material findings addressed or explicitly deferred with `not_assessed`/`HITL`
  reasons. Review disposition artifact:
  `specs/019-repo-realignment-monitoring-gate-readiness/reviews/cross-model-review-disposition.md`.

- [x] T019-110 Query live CI for the final source commit or PR head.
  Status: completed for merged PR #60 source commit; CI PASS on `main` at
  `657a343a5f310538def9afd509e6c610c713cab0` was observed after merge.
  Current post-merge closure branch PR #62 CI is live external evidence and
  must be queried from GitHub at final-head handoff. CI evidence is not merge
  approval evidence.
  PR #60: https://github.com/fall-out-bug/sdp-trace/pull/60

- [ ] T019-120 Approval gate: close this spec only after review findings,
  local verification, and live CI evidence are recorded. Missing external
  evidence must remain `cannot_verify` or `not_assessed`.
  Status: blocked on T019-004. Final post-merge Qwen3.6 Plus alternative-LLM
  review is LGTM after findings were addressed, and PR #62 final-head CI must
  be queried live from GitHub at handoff. PR #60 is a partial merge, not Spec
  019 closure.
  `merge_approval` remains `not_assessed`.
