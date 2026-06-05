# Slice 16 Plan Review: Remaining Core Cleanup

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Source shards:
  - `cmd/sdp-trace/core_152_evaluateandwritereleaseproof.go`
  - `cmd/sdp-trace/core_153_writereleaseproofresult.go`
  - `cmd/sdp-trace/core_383_witnesstargetfromflags.go`
  - `cmd/sdp-trace/core_456_exportcommandis.go`
  - `cmd/sdp-trace/core_457_exportsubcommandis.go`
  - `cmd/sdp-trace/core_458_runtelemetryexport.go`
  - `cmd/sdp-trace/core_459_parsetelemetryexportargs.go`
  - `cmd/sdp-trace/core_460_requiretelemetryexportargs.go`
  - `cmd/sdp-trace/core_461_rendertelemetryexport.go`
  - `cmd/sdp-trace/core_462_writetelemetryexportoutput.go`
  - `cmd/sdp-trace/core_463_requiretelemetryexportinputs.go`
  - `cmd/sdp-trace/core_464_runcrossrepopostureexport.go`
  - `cmd/sdp-trace/core_465_parsecrossrepopostureexportargs.go`
  - `cmd/sdp-trace/core_466_requirecrossrepopostureexportargs.go`
  - `cmd/sdp-trace/core_467_writecrossrepopostureexport.go`
  - `cmd/sdp-trace/core_468_requirecrossrepopostureinputs.go`
  - `cmd/sdp-trace/core_469_runcrossrepopostureexplain.go`
  - `cmd/sdp-trace/core_470_parsecrossrepopostureexplainargs.go`
  - `cmd/sdp-trace/core_471_readcrossrepopostureexplainresult.go`
  - `cmd/sdp-trace/core_517_fixtureexpectation.go`
  - `cmd/sdp-trace/core_518_readfixtureexpectation.go`
  - `cmd/sdp-trace/core_519_readfixtureexpectations.go`
- Target files:
  - `cmd/sdp-trace/release_proof_run.go`
  - `cmd/sdp-trace/release_proof_write.go`
  - `cmd/sdp-trace/witness_required_fields.go`
  - `cmd/sdp-trace/export_command.go`
  - `cmd/sdp-trace/export_telemetry.go`
  - `cmd/sdp-trace/export_telemetry_args.go`
  - `cmd/sdp-trace/export_telemetry_render.go`
  - `cmd/sdp-trace/export_cross_repo_posture.go`
  - `cmd/sdp-trace/export_cross_repo_posture_args.go`
  - `cmd/sdp-trace/export_cross_repo_posture_write.go`
  - `cmd/sdp-trace/export_cross_repo_posture_explain.go`
  - `cmd/sdp-trace/export_cross_repo_posture_explain_args.go`
  - `cmd/sdp-trace/export_cross_repo_posture_explain_read.go`
  - `cmd/sdp-trace/fixture_expectation_policy.go`
  - `cmd/sdp-trace/fixture_expectation_read.go`

## Decision Gate

- Simpler/Faster: move the remaining numbered core declarations into existing
  or adjacent responsibility files in the same package without changing
  function bodies, command routing, output text, exit codes, or dependencies.
- Blocking Edge Cases: broad combined export or fixture/release files fell below
  the absolute file-MI threshold during focused verification, so the final
  grouping is split by command phase and responsibility.
- Existing Open Source: not applicable; this is local Go file locality cleanup
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate for changed core-removal files
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- no remaining `cmd/sdp-trace/core_[0-9][0-9][0-9]_*.go` files
- three independent staged-diff reviewer lanes after implementation
