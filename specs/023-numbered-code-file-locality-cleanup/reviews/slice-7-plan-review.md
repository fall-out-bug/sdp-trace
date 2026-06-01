# Slice 7 Plan Review: Wrap Command Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target files:
  - `cmd/sdp-trace/wrap_legacy.go`
  - `cmd/sdp-trace/wrap_recorder.go`
  - `cmd/sdp-trace/wrap_run.go`
  - `cmd/sdp-trace/wrap_run_args.go`
  - `cmd/sdp-trace/wrap_preview.go`
  - `cmd/sdp-trace/wrap_preview_args.go`
  - `cmd/sdp-trace/wrap_preview_payload.go`
- Source shards:
  - `cmd/sdp-trace/wrap_399_run.go`
  - `cmd/sdp-trace/wrap_400_runlegacyrecorder.go`
  - `cmd/sdp-trace/wrap_401_parseargs.go`
  - `cmd/sdp-trace/wrap_402_command.go`
  - `cmd/sdp-trace/wrap_403_writerunresult.go`
  - `cmd/sdp-trace/wrap_404_runwrappedcommand.go`
  - `cmd/sdp-trace/wrap_405_runtaskrecorder.go`
  - `cmd/sdp-trace/wrap_406_parsewrappedcommandargs.go`
  - `cmd/sdp-trace/wrap_407_requirewrappedcommandargs.go`
  - `cmd/sdp-trace/wrap_408_missingrequiredcontract.go`
  - `cmd/sdp-trace/wrap_409_rundryrun.go`
  - `cmd/sdp-trace/wrap_410_runpreview.go`
  - `cmd/sdp-trace/wrap_411_runpreviewcommand.go`
  - `cmd/sdp-trace/wrap_412_writepreviewcommandpayload.go`
  - `cmd/sdp-trace/wrap_413_previewcommandpayload.go`
  - `cmd/sdp-trace/wrap_414_parsepreviewcommandargs.go`
  - `cmd/sdp-trace/wrap_415_loadpreviewcontract.go`

## Decision Gate

- Simpler/Faster: move seventeen wrap shards into behavior-named files without
  changing recorder behavior, preview output, package boundaries, or tests.
- Blocking Edge Cases: a combined preview file measured file MI `67.9`, below
  the absolute threshold, so preview payload rendering stays separate.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- three independent staged-diff reviewer lanes after implementation
