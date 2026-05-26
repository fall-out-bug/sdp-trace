# Plan: Repo Realignment, Monitoring, And Gate Readiness

Status: split_successor; residual governance moved to Spec 022

## Post-Merge Closure

PR #60 was merged on 2026-05-25 as commit
`657a343a5f310538def9afd509e6c610c713cab0`, but it did not complete this
spec. The remaining closure work is tracked in
`post-merge-closure-plan.md` and follow-up Spec 022.

The PR #60 merge approval state remains `not_assessed`: GitHub PR metadata
contains no recorded review approval, and the PR body left the review checklist
unchecked. The 2026-05-26 maintainer decision splits residual governance debt
to Spec 022 instead of retroactively approving the missed gate.

## Workstreams

### WS-019-A: Spec Reality Ledger

Risk: high

Depends: none

Mode: HITL

Owned files:

- `docs/roadmap.md`
- `specs/*/spec.md`
- `specs/*/plan.md`
- `specs/*/tasks.md`
- `specs/019-repo-realignment-monitoring-gate-readiness/*`
- optional `docs/spec-reality-ledger.md`

Deliverable:

- Reconcile roadmap, spec status, task checkboxes, and live verification state.
  Remove or qualify any "all specs implemented" implication that is not backed
  by current evidence.

After this:

- A reviewer can identify which specs are draft, in progress, blocked,
  complete, `not_assessed`, or `cannot_verify` without relying on prose memory.
- Deliverable exists at `docs/spec-reality-ledger.md` (created 2026-05-22).

### WS-019-B: OSS Tool Quality Gate Restoration

Risk: high

Depends: none

Mode: AFK

Owned files:

- `tools/osscompat/*`
- `tools/ossbench/*`
- focused tests for those tools

Deliverable:

- Refactor failing compatibility and benchmark tool functions into smaller,
  tested units. Keep probe semantics unchanged unless a change is explicitly
  recorded as a spec decision.

After this:

- `tools/osscompat` and `tools/ossbench` pass local complexity and CRAP gates.

### WS-019-C: Command Surface Maintainability Ratchet Repair

Risk: medium

Depends: none

Mode: AFK

Owned files:

- `cmd/sdp-trace/main_537_commandsurfaceconstants.go`
- `cmd/sdp-trace/main_538_commandsurfaceregistrycore.go`
- `cmd/sdp-trace/main_539_commandsurfaceregistryobserve.go`
- `cmd/sdp-trace/main_540_commandsurfaceregistryassess.go`
- `cmd/sdp-trace/main_541_commandsurfaceregistryother.go`
- `cmd/sdp-trace/main_542_commandsurfaceregistrypacket.go`
- `cmd/sdp-trace/main_546_commandsurfacedrift.go`
- focused command-surface tests or docs when needed

Deliverable:

- Restore command-surface maintainability baseline without changing the emitted
  command registry or command tiers.

After this:

- `go run ./cmd/sdp-trace command-surface` emits the same command set, and
  maintainability ratchets pass for the touched registry files.

### WS-019-D: Live Wrap Schema Compatibility

Risk: high

Depends: WS-019-A

Mode: AFK with maintainer review

Owned files:

- `schema/flight-recorder-run.schema.json`
- recorder and wrap code that produces `run.json`
- affected examples and fixtures
- focused schema and recorder tests
- affected docs that describe the recorder contract

Decision:

- Add a dedicated schema for the current live `run.json` manifest emitted by
  `wrap`, with migration notes that keep the richer flight-recorder schema
  separate from the live recorder manifest contract.

Deliverable:

- Implement the reviewed contract path selected above.

After this:

- The Spec 017 live `wrap` manifest has a schema contract and is no longer an
  ambiguous blocker.

### WS-019-E: Monitoring And Advisory Gate Proof Pack

Risk: high

Depends: WS-019-A, WS-019-D

Mode: HITL

Owned files:

- monitoring/gate examples under `examples/`
- docs that describe LLM and harness monitoring boundaries
- focused tests for unsafe event rejection and missing evidence states
- affected `harness`, `report`, or `gate` docs

Deliverable:

- Create a reproducible proof pack for explicit-export monitoring:
  `wrap` or `run`, harness observation, harness validation/reporting, and
  advisory gate facts.

After this:

- `sdp-trace` can be demonstrated as LLM/harness monitoring input and gate
  food without claiming to be an approval authority.

### WS-019-F: Numeric Shard Cleanup, Core CLI

Risk: medium

Depends: WS-019-C

Mode: AFK

Owned files:

- core CLI files for `wrap`, `run`, `verify`, `explain`, `report`, and
  `query --query missing-evidence`
- focused tests for affected command behavior

Deliverable:

- Replace selected numeric one-function shards in the core CLI path with
  cohesive behavior-named files while preserving command behavior and package
  boundaries.

After this:

- The core adoption path has fewer transitional numeric shards and clearer file
  ownership.

### WS-019-G: Numeric Shard Cleanup, Harness And Gate

Risk: medium

Depends: WS-019-E, WS-019-F

Mode: AFK with review

Owned files:

- selected `internal/harnessobs/*`
- selected harness and gate CLI glue
- focused monitoring and gate tests

Deliverable:

- Improve locality for harness observation and advisory gate code without
  weakening redaction, unsafe-input rejection, or `cannot_verify` boundaries.

After this:

- Harness and gate code has clearer domain locality, and the proof pack still
  passes.

### WS-019-H: Supply-Chain Probe Automation Closure

Risk: medium

Depends: WS-019-B

Mode: AFK

Owned files:

- `tools/osscompat/*`
- supply-chain compatibility docs and examples
- focused tests for probe state classification

Deliverable:

- Close or explicitly preserve the Spec 017 in-toto, Cosign, and SLSA probe
  gaps with reproducible state reporting.

After this:

- Supply-chain probes no longer rely on manual-only ambiguity for their current
  state.

### WS-019-I: Final Evidence And Pi Review Pack

Risk: low

Depends: WS-019-A, WS-019-B, WS-019-C, WS-019-D, WS-019-E, WS-019-F, WS-019-G,
WS-019-H

Mode: HITL

Owned files:

- `docs/roadmap.md`
- `specs/019-repo-realignment-monitoring-gate-readiness/*`
- final evidence or review disposition docs created for this spec

Deliverable:

- Update final status, collect verification evidence, and prepare Pi review
  context packs. Leave missing external evidence as `not_assessed` or
  `cannot_verify`.

After this:

- The block is ready for independent review and explicit approval or rejection.

## Verification

Docs-only changes:

```text
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

Code/tooling changes:

```text
go test -count=1 ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
go run ./tools/schemadoc -verify-readme
jq empty schema/*.json examples/block19-adapter-capture/*.json examples/self-trace/proof-summary.example.json tools/qualitycheck/function-mi-baseline.json tools/qualitycheck/file-mi-baseline.json
go run ./tools/qualitycheck -gocyclo cmd internal tools
go test -count=1 ./... -coverprofile=coverage.out
go tool cover -func=coverage.out -o coverage-func.txt
go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less
git diff --check
```

## Pi Handoff Notes

- Commit the reviewed spec handoff before launching workers.
- Build one context pack per workstream with `AGENTS.md`, relevant project
  skills, this spec, this plan, tasks, and owned files only.
- Launch workers in isolated worktrees with recorded Pi profile/model
  resolution.
- Require `subagent-result/v1`, status, events, logs, structured result, and
  inspectable diff before integration.
- Do not assign overlapping write sets to parallel workers.
- Do not allow Pi workers to merge, publish, close trust, or mark checklist
  items complete outside their accepted slice.
