# Spec 021 Slice 5 Plan Review

Date: 2026-06-01.

## Scope

Slice 5 is bounded to `cmd/sdp-trace` command-surface core command metadata
shards:

- `cmd/sdp-trace/main_548_commandsurfacecorebasic.go`
- `cmd/sdp-trace/main_549_commandsurfacecorewrap.go`
- `cmd/sdp-trace/main_550_commandsurfacecorerun.go`
- `cmd/sdp-trace/main_551_commandsurfacecorepreview.go`
- `cmd/sdp-trace/main_552_commandsurfacecoredoctor.go`
- `cmd/sdp-trace/main_553_commandsurfacecoreinstall.go`
- `cmd/sdp-trace/main_554_commandsurfacecorefixtures.go`
- `cmd/sdp-trace/main_571_commandsurfacecore.go`

Target file:

- `cmd/sdp-trace/command_surface_core_commands.go`

## Decision Gate

- Simpler/Faster: move core command metadata values and the core command list
  into one behavior-named file.
- Blocking Edge Cases: none found for this bounded data-only grouping; local
  pre-change MI analysis measured the combined file at `100.0`.
- Existing Open Source: not applicable; this is intra-package source locality
  cleanup with existing local metadata values.

## Review

Plan result: proceed.

Rationale:

- The slice removes eight numbered shards.
- The target file has a cohesive core command metadata purpose.
- No package boundary, schema contract, dependency, trust, or baseline change is
  planned.

Not assessed:

- Behavior preservation remains `not_assessed` until local verification and
  three post-implementation reviews are recorded.
