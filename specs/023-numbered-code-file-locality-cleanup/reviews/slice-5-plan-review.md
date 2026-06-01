# Slice 5 Plan Review: Fixture Validation Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target files:
  - `cmd/sdp-trace/fixture_validation_run.go`
  - `cmd/sdp-trace/fixture_validation_args.go`
  - `cmd/sdp-trace/fixture_expectation_policy.go`
- Source shards:
  - `cmd/sdp-trace/fixture_472_run.go`
  - `cmd/sdp-trace/fixture_473_rootarg.go`
  - `cmd/sdp-trace/fixture_474_validatefixtureruns.go`
  - `cmd/sdp-trace/fixture_475_validatefixturerun.go`
  - `cmd/sdp-trace/fixture_476_expectationfailed.go`
  - `cmd/sdp-trace/fixture_477_expectedresultfailed.go`
  - `cmd/sdp-trace/fixture_478_unexpectedresultfailed.go`

## Decision Gate

- Simpler/Faster: move seven fixture validation shards into behavior-named
  files without changing values, command behavior, package boundaries, or tests.
- Blocking Edge Cases: a single-file grouping measured file MI `65.4`, below
  the absolute threshold, and keeping root-arg parsing in the runner measured
  file MI `68.6`, so runner/iteration, root args, and expectation policy stay
  separate.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- three independent staged-diff reviewer lanes after implementation
