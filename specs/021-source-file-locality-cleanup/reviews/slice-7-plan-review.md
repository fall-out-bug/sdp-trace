# Slice 7 Plan Review: Assess Command Metadata Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target file: `cmd/sdp-trace/command_surface_assess_commands.go`
- Source shards:
  - `cmd/sdp-trace/main_559_commandsurfaceassessassess.go`
  - `cmd/sdp-trace/main_560_commandsurfaceassessreport.go`
  - `cmd/sdp-trace/main_561_commandsurfaceassessgate.go`
  - `cmd/sdp-trace/main_562_commandsurfaceassesscheckpoint.go`
  - `cmd/sdp-trace/main_573_commandsurfaceassess.go`

## Decision Gate

- Simpler/Faster: move the five assess metadata shards into one behavior-named
  file without changing values, command behavior, package boundaries, or tests.
- Blocking Edge Cases: accidental metadata edits, package init/order changes,
  MI/CRAP regression, or stale evidence could invalidate the cleanup claim.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Review

- Requirements fit: pass. The slice is bounded to one command family and avoids
  a repo-wide rename sweep.
- UX impact: not_assessed. No CLI output or behavior change is intended.
- DX impact: pass. The target file removes five numbered shards from the active
  product path.
- Maintainability: pass. The grouped file keeps command metadata and the group
  list together.
- Trust boundary: pass. Production trust, release approval, merge approval, and
  external attestation remain not_assessed.

## Planned Verification

- `gofmt -w cmd/sdp-trace/command_surface_assess_commands.go`
- `go test ./cmd/sdp-trace`
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- three independent reviewer lanes after implementation
