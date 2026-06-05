# Slice 14 Plan Review: Core CLI Kernel Cleanup

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Source shards:
  - `cmd/sdp-trace/core_001_const.go`
  - `cmd/sdp-trace/core_002_clistdin.go`
  - `cmd/sdp-trace/core_003_version.go`
  - `cmd/sdp-trace/core_004_commandhandler.go`
  - `cmd/sdp-trace/core_005_subcommandhandler.go`
  - `cmd/sdp-trace/core_006_commandhandlers.go`
  - `cmd/sdp-trace/core_007_runversion.go`
  - `cmd/sdp-trace/core_008_main.go`
  - `cmd/sdp-trace/core_009_run.go`
  - `cmd/sdp-trace/core_010_toplevelhelp.go`
  - `cmd/sdp-trace/core_011_dispatchcommand.go`
  - `cmd/sdp-trace/core_014_runsubcommand.go`
  - `cmd/sdp-trace/core_015_dispatchsubcommand.go`
  - `cmd/sdp-trace/core_016_runoptionalsubcommand.go`
  - `cmd/sdp-trace/core_017_ishelparg.go`
  - `cmd/sdp-trace/core_018_subcommandname.go`
  - `cmd/sdp-trace/core_019_rejectrest.go`
  - `cmd/sdp-trace/core_020_requirestringflag.go`
  - `cmd/sdp-trace/core_021_requiredcliflag.go`
  - `cmd/sdp-trace/core_022_requireonlyflags.go`
  - `cmd/sdp-trace/core_023_requirerequiredflags.go`
  - `cmd/sdp-trace/core_024_writejsonpayload.go`
  - `cmd/sdp-trace/core_025_writejsonpayloadunchecked.go`
  - `cmd/sdp-trace/core_026_requirenamedvalues.go`
  - `cmd/sdp-trace/core_027_stringexitcode.go`
- Target files:
  - `cmd/sdp-trace/cli_handlers.go`
  - `cmd/sdp-trace/cli_dispatch.go`
  - `cmd/sdp-trace/cli_main.go`
  - `cmd/sdp-trace/cli_subcommands.go`
  - `cmd/sdp-trace/cli_subcommand_helpers.go`
  - `cmd/sdp-trace/cli_flag_validation.go`
  - `cmd/sdp-trace/cli_named_values.go`
  - `cmd/sdp-trace/cli_json.go`
  - `cmd/sdp-trace/cli_exit_codes.go`

## Decision Gate

- Simpler/Faster: move only the CLI kernel declarations into responsibility
  files in the same package without changing function bodies, command names,
  help behavior, exit-code values, JSON formatting, or dependencies.
- Blocking Edge Cases: all numbered core files are not one responsibility. A
  single combined core file would mix CLI dispatch, assessment explanation,
  preview helpers, exports, and fixture loading, increasing drift and MI risk.
- Existing Open Source: not applicable; this is local Go file locality cleanup
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate for new CLI kernel files
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- no remaining `cmd/sdp-trace/core_00*_*.go`, `core_01*_*.go`, or
  `core_02[0-7]_*.go` files from this slice scope
- three independent staged-diff reviewer lanes after implementation
