# Slice 6 Plan Review: Interaction Command Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target files:
  - `cmd/sdp-trace/interaction_command.go`
  - `cmd/sdp-trace/interaction_relay.go`
  - `cmd/sdp-trace/interaction_relay_args.go`
  - `cmd/sdp-trace/interaction_transcript_import.go`
  - `cmd/sdp-trace/interaction_transcript_import_args.go`
  - `cmd/sdp-trace/interaction_summary.go`
  - `cmd/sdp-trace/cli_flag_requirements.go`
- Source shards:
  - `cmd/sdp-trace/interaction_158_run.go`
  - `cmd/sdp-trace/interaction_159_runrelay.go`
  - `cmd/sdp-trace/interaction_160_relayoptions.go`
  - `cmd/sdp-trace/interaction_161_parserelayargs.go`
  - `cmd/sdp-trace/interaction_162_newrelayflagset.go`
  - `cmd/sdp-trace/interaction_163_relaystringflags.go`
  - `cmd/sdp-trace/interaction_164_requirerest.go`
  - `cmd/sdp-trace/interaction_165_requireonlyflagscode.go`
  - `cmd/sdp-trace/interaction_166_requiredflags.go`
  - `cmd/sdp-trace/interaction_167_runimporttranscript.go`
  - `cmd/sdp-trace/interaction_168_writeimportedtranscript.go`
  - `cmd/sdp-trace/interaction_169_importtranscriptfromoptions.go`
  - `cmd/sdp-trace/interaction_170_parseimporttranscriptargs.go`
  - `cmd/sdp-trace/interaction_171_runsummarize.go`
  - `cmd/sdp-trace/interaction_172_parsesummarizeargs.go`

## Decision Gate

- Simpler/Faster: move fifteen interaction shards into behavior-named files
  without changing command behavior, forwarded command handling, package
  boundaries, or tests.
- Blocking Edge Cases: a combined transcript import file measured file MI
  `69.5`, below the absolute threshold, so import args parsing stays separate.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- three independent staged-diff reviewer lanes after implementation
