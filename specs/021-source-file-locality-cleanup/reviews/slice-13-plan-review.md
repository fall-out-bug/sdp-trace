# Slice 13 Plan Review: Assess Command Cleanup

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Source shards: `cmd/sdp-trace/assess_179_runassess.go` through
  `cmd/sdp-trace/assess_233_previewinputexitcode.go`
- Target files:
  - `cmd/sdp-trace/assess_command.go`
  - `cmd/sdp-trace/assess_command_flags.go`
  - `cmd/sdp-trace/assess_command_registry.go`
  - `cmd/sdp-trace/assess_profiles.go`
  - `cmd/sdp-trace/assess_profiles_artifacts.go`
  - `cmd/sdp-trace/assess_requirements.go`
  - `cmd/sdp-trace/assess_writers.go`
  - `cmd/sdp-trace/assess_inputs.go`
  - `cmd/sdp-trace/assess_inputs_managed.go`
  - `cmd/sdp-trace/assess_inputs_managed_json.go`
  - `cmd/sdp-trace/assess_preview_command.go`
  - `cmd/sdp-trace/assess_preview_registry.go`
  - `cmd/sdp-trace/assess_preview_reports.go`
  - `cmd/sdp-trace/assess_preview_adapter.go`
  - `cmd/sdp-trace/assess_preview_managed.go`
  - `cmd/sdp-trace/assess_preview_forensic.go`
  - `cmd/sdp-trace/assess_preview_ci_artifact.go`
  - `cmd/sdp-trace/assess_preview_authority.go`

## Decision Gate

- Simpler/Faster: move declarations into responsibility-named files inside the
  same package without changing function bodies, profile names, JSON fields,
  package boundaries, or dependencies.
- Blocking Edge Cases: one combined assess file would mix command dispatch,
  profile execution, input loading, and preview report metadata, and would risk
  file-MI regression. Splitting each original function into its own file is the
  debt being removed.
- Existing Open Source: not applicable; this is local Go file locality cleanup
  using existing project verification tools and conventions.

## Planned Verification

- focused `cmd/sdp-trace` test and file MI gate for new assess files
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- zero remaining `cmd/sdp-trace/assess_[0-9]*_*.go` files
- three independent staged-diff reviewer lanes after implementation
