# Slice 17 Plan Review: Harnessobs Type And Lookup Cleanup

Status: pass

## Scope

- Package: `internal/harnessobs`
- Source shards:
  - `internal/harnessobs/harnessobs_011_observeoptions.go`
  - `internal/harnessobs/harnessobs_012_sessionprofile.go`
  - `internal/harnessobs/harnessobs_013_sessionsetupaction.go`
  - `internal/harnessobs/harnessobs_014_sessionisolationrule.go`
  - `internal/harnessobs/harnessobs_015_sessionrun.go`
  - `internal/harnessobs/harnessobs_016_sessionisolationresult.go`
  - `internal/harnessobs/harnessobs_017_sessionsetupoptions.go`
  - `internal/harnessobs/harnessobs_018_sessioncollectoptions.go`
  - `internal/harnessobs/harnessobs_019_sessionoptions.go`
  - `internal/harnessobs/harnessobs_020_validateoptions.go`
  - `internal/harnessobs/harnessobs_021_observationcontext.go`
  - `internal/harnessobs/harnessobs_022_sessioncollectioncontext.go`
  - `internal/harnessobs/harnessobs_023_observedcommandresult.go`
  - `internal/harnessobs/harnessobs_024_isolationruleinstallers.go`
  - `internal/harnessobs/harnessobs_025_eventrefcheck.go`
  - `internal/harnessobs/harnessobs_026_staterank.go`
  - `internal/harnessobs/harnessobs_027_existingpathspec.go`
  - `internal/harnessobs/harnessobs_028_shellfieldscanner.go`
  - `internal/harnessobs/harnessobs_029_digestfieldnames.go`
  - `internal/harnessobs/harnessobs_030_validfamilies.go`
  - `internal/harnessobs/harnessobs_031_validcontentstates.go`
  - `internal/harnessobs/harnessobs_032_validstates.go`
  - `internal/harnessobs/harnessobs_033_validrulekeys.go`
- Target files:
  - `internal/harnessobs/options_context.go`
  - `internal/harnessobs/session_model.go`
  - `internal/harnessobs/validation_sets.go`
  - `internal/harnessobs/event_ref_checks.go`
  - `internal/harnessobs/existing_path_spec.go`
  - `internal/harnessobs/shell_field_scanner.go`
  - `internal/harnessobs/isolation_rule_installers.go`

## Decision Gate

- Simpler/Faster: move type and lookup declarations into responsibility-named
  files
  in the same package without changing any field, JSON tag, map entry, or
  dependency.
- Blocking Edge Cases: all selected shards are declaration-only or map-only;
  behavior-heavy observation, parsing, path-safety, and event validation shards
  are excluded to avoid broad mixed-responsibility review.
- Existing Open Source: not applicable; this is local Go file locality cleanup
  using existing project verification tools and conventions.

## Planned Verification

- focused `internal/harnessobs` test and file MI gate with `70.1` threshold for
  changed files
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates with `70.1` MI thresholds
- no remaining `internal/harnessobs/harnessobs_0[1-3][0-9]_*.go` files from
  this slice scope except `harnessobs_034` through `harnessobs_039`, which
  begin the next behavior slice
- three independent staged-diff reviewer lanes after implementation
