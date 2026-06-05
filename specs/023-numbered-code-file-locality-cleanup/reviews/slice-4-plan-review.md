# Slice 4 Plan Review: Export Command Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target file: `cmd/sdp-trace/export_command.go`
- Source shards:
  - `cmd/sdp-trace/export_452_run.go`
  - `cmd/sdp-trace/export_453_telemetryrequested.go`
  - `cmd/sdp-trace/export_454_crossrepopostureexplainrequested.go`
  - `cmd/sdp-trace/export_455_crossrepoposturerequested.go`

## Decision Gate

- Simpler/Faster: move the export dispatcher and three tiny predicate helpers
  into one behavior-named file without changing values, command behavior,
  package boundaries, or tests.
- Blocking Edge Cases: file MI must remain above the absolute threshold after
  grouping; if it fails, split dispatcher and predicates into separate
  behavior-named files.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- three independent staged-diff reviewer lanes after implementation
