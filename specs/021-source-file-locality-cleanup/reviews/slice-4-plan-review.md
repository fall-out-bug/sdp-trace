# Spec 021 Slice 4 Plan Review

Date: 2026-06-01.

## Scope

Slice 4 is bounded to `cmd/sdp-trace` command-surface schema and runner shards:

- `cmd/sdp-trace/main_536_commandsurfaceschema.go`
- `cmd/sdp-trace/main_544_commandsurfacejson.go`
- `cmd/sdp-trace/main_545_runcommandsurface.go`

Target files:

- `cmd/sdp-trace/command_surface_schema.go`
- `cmd/sdp-trace/command_surface_runner.go`

## Decision Gate

- Simpler/Faster: move schema types and command-surface runner functions into
  behavior-named files without changing behavior.
- Blocking Edge Cases: broader metadata/registry grouping failed local
  pre-change MI analysis and would risk mixed code/baseline churn.
- Existing Open Source: not applicable; this is intra-package source locality
  cleanup with existing local code.

## Review

Plan result: proceed.

Rationale:

- The slice removes three numbered shards.
- Schema types and command runner behavior are cohesive, stable targets.
- No package boundary, schema contract, dependency, trust, or baseline change is
  planned.

Not assessed:

- Behavior preservation remains `not_assessed` until local verification and
  three post-implementation reviews are recorded.
