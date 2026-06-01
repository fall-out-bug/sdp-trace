# Slice 3 Plan Review: Envelope Command Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target files:
  - `cmd/sdp-trace/envelope_summary_run.go`
  - `cmd/sdp-trace/envelope_summary_args.go`
- Source shards:
  - `cmd/sdp-trace/envelope_173_run.go`
  - `cmd/sdp-trace/envelope_174_requiredflags.go`
  - `cmd/sdp-trace/envelope_174_writeoptionaljsonfile.go`
  - `cmd/sdp-trace/envelope_175_parsesummarizeargs.go`

## Decision Gate

- Simpler/Faster: move four tiny envelope summarize shards into behavior-named
  files without changing values, command behavior, package boundaries, or tests.
- Blocking Edge Cases: a single-file grouping measured file MI `69.0`, below
  the absolute threshold, so runner/output and args/policy stay separate.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- three independent staged-diff reviewer lanes after implementation
