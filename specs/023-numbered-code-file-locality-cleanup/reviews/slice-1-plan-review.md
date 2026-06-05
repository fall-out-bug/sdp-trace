# Slice 1 Plan Review: Release-Proof Command Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target files:
  - `cmd/sdp-trace/release_proof_run.go`
  - `cmd/sdp-trace/release_proof_args.go`
  - `cmd/sdp-trace/release_proof_policy.go`
- Source shards:
  - `cmd/sdp-trace/release_151_run.go`
  - `cmd/sdp-trace/release_154_parseargs.go`
  - `cmd/sdp-trace/release_155_flagsandexits.go`

## Decision Gate

- Simpler/Faster: move the three release-proof shards into behavior-named files
  without changing values, command behavior, package boundaries, or tests.
- Blocking Edge Cases: single-file release-proof grouping fails the absolute
  file-MI threshold; args+exits grouping also fails the threshold.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- three independent staged-diff reviewer lanes after implementation
