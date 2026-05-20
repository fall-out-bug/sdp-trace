# Spec 018: Core/Policy Split And Pi Delivery Plan

Status: draft

## Objective

Define a smaller product shape that separates the stable evidence substrate
from policy, witness, PR-review, telemetry, and forensic extensions. This spec
turns the simplification analysis into work that can be delegated safely to Pi
agents.

## Product Shape Decision

### Core

- `wrap`
- `run`
- `verify`
- `explain`
- `report`
- `query --query missing-evidence`
- schema validation for the current run/report artifacts

### Extension Or Downstream Areas

- assessment profiles
- witness and signing
- protected gate/checkpoint
- PR review packet pipeline
- cross-repo posture and telemetry
- forensic query packs
- repo observer installation

## Decisions

- The core command set for this simplification phase is exactly `wrap`, `run`,
  `verify`, `explain`, `report`, and `query --query missing-evidence`.
- All other current command families remain in the same binary for this phase
  and are classified as extension, experimental, or fixture-only. Separate
  binaries, plugins, and removals require a later implementation spec.
- `query --query capture-depth` belongs to the adapter-capture diagnostics
  extension. It is not part of the core query surface.
- `query-pack` remains a forensic extension. The remaining query-pack code in
  `internal/query` must be split before any package-level minimal-core claim.
- `release-proof`, `witness`, `gate`, `checkpoint`, `packet`, and `pr-review`
  are useful trust or policy surfaces but are not core adoption commands.
- First implementation work must reduce package coupling and documentation
  complexity without changing command behavior or deleting commands.
- Numbered one-function Go source shards such as `family_123_name.go` are a
  transitional organization artifact, not the target architecture. Cleanup
  work must converge them into cohesive, behavior-named files by command or
  domain while preserving package boundaries and test coverage.
- File renaming must not be used as metric gaming. The target is readable
  locality, stable ownership, and dependency direction; complexity and MI gates
  remain verification signals, not design goals.
- A future behavior-changing simplification must start from a new reviewed
  implementation spec, not from ad hoc execution inside this planning slice.

## Requirements

- FR-018-001: Define command-family stability tiers: core, extension,
  experimental, fixture-only, and deprecated/not_assessed.
- FR-018-002: Update docs so the first user path starts with the core surface,
  not the full controlled-pilot command list.
- FR-018-003: Add a migration/deprecation plan for non-core surfaces without
  breaking current tests.
- FR-018-004: Identify which packages remain in core and which become extension
  candidates.
- FR-018-005: Keep all trust states explicit; simplification must not remove
  `not_assessed` or `cannot_verify` boundaries.
- FR-018-006: Produce Pi-ready implementation slices with disjoint file
  ownership and verification commands.
- FR-018-007: Define source-file locality rules that reject numbered
  one-function shard names as the target state and require cohesive
  behavior-named files for future cleanup work.

## Non-Goals

- No immediate deletion of non-core commands.
- No behavior change without a reviewed implementation spec.
- No claim that core-only adoption is production trust.
- No decision to split into separate binaries or plugins in this phase.
- No broad rename of numbered Go files in this planning PR.

## Acceptance Criteria

- A core/extension command matrix exists.
- Documentation reading order points new adopters to the core path first.
- Extension surfaces have explicit owner, status, and spec decision.
- Source-file locality policy identifies numbered Go shards as transitional
  and requires future cleanup slices to replace them with cohesive named files.
- Pi handoff slices are small enough to run in isolated worktrees without
  overlapping write sets.
