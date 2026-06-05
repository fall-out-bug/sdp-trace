# Slice 10 Plan Review: Command Surface Catalog Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target file: `cmd/sdp-trace/command_surface_catalog.go`
- Source shards:
  - `cmd/sdp-trace/main_577_commandsurfaceprofiles.go`
  - `cmd/sdp-trace/main_578_commandsurfacestates.go`

## Decision Gate

- Simpler/Faster: move the two catalog metadata shards into one behavior-named
  data file without changing values, command behavior, package boundaries, or
  tests.
- Blocking Edge Cases: broader registry/constants/catalog grouping fails the
  absolute file-MI threshold and would force a mixed code/baseline change.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- three independent staged-diff reviewer lanes after implementation
