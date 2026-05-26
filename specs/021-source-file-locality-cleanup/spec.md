# Spec 021: Source File Locality Cleanup

Status: draft follow-up prepared by Spec 018 closure.

## Objective

Replace transitional numbered one-function Go source shards with cohesive,
behavior-named files by package or command family while preserving command
behavior, tests, and quality gates.

## Background

Spec 018 recorded numbered source shards such as `family_123_name.go` as a
transitional organization artifact, not the target repository shape. This spec
turns that policy into bounded cleanup slices.

## Requirements

- FR-021-001: Split cleanup by package or command family; no repo-wide rename
  sweep.
- FR-021-002: Preserve current public command behavior and output contracts.
- FR-021-003: Group functions into cohesive files named after behavior,
  command family, or domain concept.
- FR-021-004: Keep complexity, MI, and CRAP gates as verification signals, not
  as a reason for metric-gaming file moves.
- FR-021-005: Preserve package boundaries and dependency direction from
  `docs/package-ownership-map.md`.

## Non-Goals

- No command behavior change.
- No package split unless another reviewed spec owns it.
- No broad formatting churn outside touched packages.
- No production trust, release approval, or external attestation claim.

## Acceptance Criteria

- Each cleanup slice lists the exact package or command family it touches.
- Tests and quality gates pass after each slice.
- `docs/package-ownership-map.md` remains accurate for touched packages.

