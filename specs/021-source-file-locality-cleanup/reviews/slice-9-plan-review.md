# Slice 9 Plan Review: Other Command Metadata Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target file: `cmd/sdp-trace/command_surface_other_commands.go`
- Source shards:
  - `cmd/sdp-trace/main_565_commandsurfaceotherverify.go`
  - `cmd/sdp-trace/main_566_commandsurfaceotherquery.go`
  - `cmd/sdp-trace/main_567_commandsurfaceotherwitness.go`
  - `cmd/sdp-trace/main_568_commandsurfaceotherrelease.go`
  - `cmd/sdp-trace/main_569_commandsurfaceotheroverride.go`
  - `cmd/sdp-trace/main_570_commandsurfaceotherexport.go`
  - `cmd/sdp-trace/main_574_commandsurfaceother.go`

## Decision Gate

- Simpler/Faster: move the seven other metadata shards into one behavior-named
  file without changing values, command behavior, package boundaries, or tests.
- Blocking Edge Cases: accidental metadata edits, package init/order changes,
  MI/CRAP regression, or stale evidence could invalidate the cleanup claim.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- three independent staged-diff reviewer lanes after implementation
