# Slice 2 Plan Review: Observe Command Adapter Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target files:
  - `cmd/sdp-trace/observe_command_adapters.go`
  - `cmd/sdp-trace/observe_exit_policy.go`
- Source shards:
  - `cmd/sdp-trace/observe_012_runcommand.go`
  - `cmd/sdp-trace/observe_013_runharness.go`
  - `cmd/sdp-trace/observe_028_harnessstateexits.go`

## Decision Gate

- Simpler/Faster: move the three observe shards into behavior-named files
  without changing values, command behavior, package boundaries, or tests.
- Blocking Edge Cases: single-file observe grouping fails the absolute file-MI
  threshold, so command adapters and exit policy stay separate.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- three independent staged-diff reviewer lanes after implementation
