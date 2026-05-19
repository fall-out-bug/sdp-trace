# Spec 018: Core/Policy Split And Pi Delivery Plan

Status: draft

## Objective

Define a smaller product shape that separates the stable evidence substrate
from policy, witness, PR-review, telemetry, and forensic extensions. This spec
turns the simplification analysis into work that can be delegated safely to Pi
agents.

## Proposed Product Shape

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

## Non-Goals

- No immediate deletion of non-core commands.
- No behavior change without a reviewed implementation spec.
- No claim that core-only adoption is production trust.

## Acceptance Criteria

- A core/extension command matrix exists.
- Documentation reading order points new adopters to the core path first.
- Extension surfaces have explicit owner, status, and next decision.
- Pi handoff slices are small enough to run in isolated worktrees without
  overlapping write sets.
