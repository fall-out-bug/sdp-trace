# Spec 019: Repo Realignment, Monitoring, And Gate Readiness

Status: split_successor; residual governance moved to Spec 022

## Objective

Turn the repository-wide readiness audit into a reviewed implementation plan
that can be delegated to Pi agents without mixing design decisions into
execution. The work realigns specs, CI quality, schema compatibility, monitoring
evidence, gate inputs, and source-file locality before any production-readiness
or trust-closure claim is made.

## Problem Statement

Current repository behavior is useful as a controlled-pilot evidence substrate,
but the checked-in planning state, quality gates, schema compatibility evidence,
and source layout do not yet support a claim that all specs are complete or that
`sdp-trace` can act as production trust authority.

The intended product role remains narrower:

- collect and preserve evidence for LLM, harness, and command runs;
- normalize monitoring facts into explicit `pass`, `fail`, `cannot_verify`, and
  `not_assessed` states;
- feed gate commands with evidence-backed facts;
- avoid deciding merge, release, readiness, risk, or production trust.

## Audit Inputs

This spec is based on the repository audit performed on 2026-05-22 against
`main` at `45ad723`. The audit state is structural evidence only until replayed
in a clean verification run.

Observed inputs:

- Latest observed `main` CI was failing in the quality gate for
  `tools/osscompat` and `tools/ossbench`.
- Local `go test -count=1 ./...`, `go vet ./...`, `go build ./cmd/sdp-trace`,
  `doccheck`, `hygienecheck`, `schemadoc`, `jq`, and `git diff --check` passed.
- CRAP check (coverage-backed: `go test -count=1 ./... -coverprofile=coverage.out`, `go tool cover -func=coverage.out > coverage-func.txt`, `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`, then `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`) **FAIL**: 21 functions exceed threshold, mostly in `tools/ossbench` and `tools/osscompat`.
- MI baseline check (`go run ./tools/qualitycheck -fail-only -function-mi-under 70
  -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal tools`)
  **FAIL**: 2 regressions in `cmd/sdp-trace/main_540_commandsurfaceregistryassess.go`
  and `main_541_commandsurfaceregistryother.go`; many missing baselines for `tools`
  packages.
- Spec and roadmap status did not support an "all specs implemented" claim.
  See `docs/spec-reality-ledger.md` for the full reconciliation.
- Spec 017 left live `wrap` output versus
  `schema/flight-recorder-run.schema.json` as an open blocker.
- Numeric Go source shards such as `assess_179_runassess.go` remained a broad
  transitional organization artifact, not the target architecture.
- Harness observation rejected unsafe raw prompt fields before writing a run
  directory.
- A zero-event harness validation produced `not_assessed` with explicit missing
  required event-family evidence.
- Gate output stayed advisory: local evidence could pass while CI witness and
  audit-grade evidence remained `cannot_verify`.

## Product Decisions

- Spec and roadmap status must be made truthful before implementation workers
  close tasks or claim readiness.
- CI quality failures are blockers for PR-ready claims, even when unit tests and
  docs checks pass.
- Live recorder/schema drift must be fixed or explicitly versioned before any
  flight-recorder schema compatibility claim.
- Monitoring scope is explicit-export monitoring. `sdp-trace` must not claim to
  detect every unwrapped LLM or harness session.
- Gate scope is evidence preparation and advisory fact production. `sdp-trace`
  must not claim to approve merge, release, production readiness, or risk.
- Numbered one-function Go shards are allowed as historical/transitional code,
  but new cleanup work must converge touched areas toward cohesive,
  behavior-named files.
- File renames must be justified by locality, ownership, and dependency
  direction, not by metric gaming.
- Pi agents may implement approved slices only after this spec handoff is
  reviewed and committed. Workers must not merge, publish, close trust, or
  rewrite task checkboxes outside their owned slice.

## Requirements

- FR-019-001: Produce a spec reality ledger that reconciles spec files,
  roadmap rows, task checkboxes, and live verification states.
- FR-019-002: Restore CI quality gates for OSS compatibility and benchmark
  tooling without changing product command behavior.
- FR-019-003: Restore command-surface maintainability ratchets while preserving
  the current command list.
- FR-019-004: Resolve or version the live `wrap` output/schema compatibility
  blocker from Spec 017.
- FR-019-005: Provide an end-to-end monitoring and gate proof pack that covers
  `wrap` or `run`, harness observation, validation/reporting, and advisory gate
  output.
- FR-019-006: Preserve raw-prompt, raw-response, token, and private-path safety
  boundaries in all monitoring examples and tests.
- FR-019-007: Define and execute source-file locality cleanup slices that reduce
  numeric shard usage in touched areas without broad unrelated rewrites.
- FR-019-008: Close or explicitly retain Spec 017 supply-chain probe gaps with
  reproducible `pass`, `fail`, `cannot_verify`, or `not_assessed` states.
- FR-019-009: Prepare Pi handoff slices with disjoint write ownership,
  verification commands, and review expectations.

## Non-Goals

- No production trust, release-readiness, or merge-approval claim.
- No deletion of command families solely because they are outside the core
  adoption path.
- No repository-wide rename of numbered Go files in one worker.
- No new Node.js, npm, JavaScript, TypeScript, or `.mjs` tooling in the active
  product path.
- No dependency on SDP Operator Mode, Beads, agentloop, or a specific harness
  runtime.
- No automatic OSS replacement decision from benchmark numbers alone.

## Acceptance Criteria

- The roadmap and spec reality ledger describe current repository state without
  claiming completed specs that lack live evidence.
- Latest local quality verification for touched slices is recorded as `pass`,
  `fail`, `cannot_verify`, or `not_assessed`.
- OSS compatibility and benchmark tools pass the repository quality gates or
  their remaining gaps are explicitly blocked with reasons.
- Live `wrap` output either validates against the intended schema or the schema
  contract is versioned and documented with migration impact.
- Monitoring proof pack demonstrates useful LLM/harness evidence collection
  while preserving explicit-export and redaction boundaries.
- Gate proof pack demonstrates evidence-backed advisory facts and keeps missing
  CI/audit-grade evidence as `cannot_verify`.
- Numeric source shard cleanup reduces numeric shards in the owned areas and
  leaves command behavior unchanged.
- Each Pi-ready slice names owned files, dependencies, verification commands,
  and HITL versus AFK status.
