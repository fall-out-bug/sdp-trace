# Slice 8 Plan Review: Query Command Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target files:
  - `cmd/sdp-trace/query_verify.go`
  - `cmd/sdp-trace/query_verify_args.go`
  - `cmd/sdp-trace/query_verify_exit.go`
  - `cmd/sdp-trace/query_explain.go`
  - `cmd/sdp-trace/query_run.go`
  - `cmd/sdp-trace/query_dispatch.go`
  - `cmd/sdp-trace/query_pack.go`
  - `cmd/sdp-trace/query_pack_build.go`
  - `cmd/sdp-trace/query_pack_explain.go`
  - `cmd/sdp-trace/query_pack_args.go`
  - `cmd/sdp-trace/query_pack_validation.go`
- Source shards:
  - `cmd/sdp-trace/query_430_runverify.go`
  - `cmd/sdp-trace/query_431_parseverifyargs.go`
  - `cmd/sdp-trace/query_432_existingdirectory.go`
  - `cmd/sdp-trace/query_433_verifierresultexitcodes.go`
  - `cmd/sdp-trace/query_434_verifierresultexitcode.go`
  - `cmd/sdp-trace/query_435_runexplain.go`
  - `cmd/sdp-trace/query_436_runquery.go`
  - `cmd/sdp-trace/query_437_runnamedquery.go`
  - `cmd/sdp-trace/query_438_capturedepthquery.go`
  - `cmd/sdp-trace/query_439_missingevidencequery.go`
  - `cmd/sdp-trace/query_440_runquerypack.go`
  - `cmd/sdp-trace/query_441_runquerypackbuild.go`
  - `cmd/sdp-trace/query_442_writequerypackartifact.go`
  - `cmd/sdp-trace/query_443_runquerypackexplain.go`
  - `cmd/sdp-trace/query_444_options.go`
  - `cmd/sdp-trace/query_446_parsequerypackargs.go`
  - `cmd/sdp-trace/query_447_parsequerypackexplainargs.go`
  - `cmd/sdp-trace/query_448_validatequerypackoptions.go`
  - `cmd/sdp-trace/query_449_requirequerypackrequiredinputs.go`
  - `cmd/sdp-trace/query_450_readquerypackresult.go`
  - `cmd/sdp-trace/query_451_validatequerypackexplainresult.go`

## Decision Gate

- Simpler/Faster: move twenty-one query shards into behavior-named files
  without changing verifier output, query payloads, query-pack artifacts,
  package boundaries, or tests.
- Blocking Edge Cases: combined verify, query, and query-pack files measured
  file MI below the absolute threshold, so runner, args, dispatch, and
  validation responsibilities stay separate.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- three independent staged-diff reviewer lanes after implementation
