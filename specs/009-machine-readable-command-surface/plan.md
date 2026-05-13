# Implementation Plan: Machine-Readable Command Surface

## Summary

Create a typed command-surface registry that can render the existing help text and emit deterministic JSON for agents and docs checks.

## Technical Context

**Language**: Go
**Dependencies**: Go standard library, existing CLI helpers, existing `tools/doccheck`
**Verification**: `go test -count=1 ./...`, `go run ./tools/doccheck`, `jq empty` for committed JSON examples, `git diff --check`
**Constraints**: Go-only product path; no CLI semantic rewrites in this slice.

## Phase 0 - PI Spec Review

Run PI review before implementation:

- requirements-vs-implementation risk;
- DX/UX clarity;
- code architecture risk;
- trust/compatibility overclaim risk.

Stop for explicit approval before code changes.

## Phase 1 - Registry Skeleton

- Add typed command metadata near the CLI command registry.
- Include schema version, command path, flags, rest behavior, output behavior, and profile values where known.
- Mark unknown or partially modeled sections explicitly.

## Phase 2 - Rendering And Checks

- Render top-level help from the registry, or add a checker that proves frozen help matches the registry.
- Extend `tools/doccheck` to compare docs against the registry instead of help prose alone.
- Add focused tests for metadata output and drift failures.

## Phase 3 - Docs And Review

- Document the JSON surface for agent consumers.
- Run full verification.
- Run PI implementation review.
- Record advisory follow-ups separately from blockers.

## Exit Criteria

- Machine-readable command output is deterministic and parseable.
- Existing help and command behavior remain compatible.
- Docs/help/registry drift is machine-checked.
- No external trust or semantic-completeness overclaim is introduced.
