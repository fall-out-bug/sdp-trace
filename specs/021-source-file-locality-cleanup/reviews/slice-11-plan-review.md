# Slice 11 Plan Review: Final Numbered File Cleanup

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target files:
  - `cmd/sdp-trace/cli_arg_helpers.go`
  - `cmd/sdp-trace/command_surface_metadata.go`
  - `cmd/sdp-trace/command_surface_registry.go`
- Source shards:
  - `cmd/sdp-trace/main_536_ishelp.go`
  - `cmd/sdp-trace/main_537_isboolliteral.go`
  - `cmd/sdp-trace/main_537_commandsurfaceconstants.go`
  - `cmd/sdp-trace/main_579_commandsurfaceregistryvar.go`

## Decision Gate

- Simpler/Faster: move the remaining numbered declarations into three
  responsibility-named files without changing values, command behavior, package
  boundaries, or tests.
- Blocking Edge Cases: combining all remaining declarations into one file would
  mix unrelated responsibilities and earlier MI analysis showed broader
  metadata/registry grouping below the absolute file-MI threshold.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- zero remaining `cmd/sdp-trace/main_[0-9]*.go` files
- three independent staged-diff reviewer lanes after implementation
