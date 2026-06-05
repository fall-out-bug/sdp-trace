# Slice 19 Plan Review: Harnessobs OpenCode Normalization

Status: pass

## Scope

- Package: `internal/harnessobs`
- Source shards: `harnessobs_046` through `harnessobs_064`
- Target files:
  - `internal/harnessobs/opencode_raw_line.go`
  - `internal/harnessobs/opencode_events.go`
  - `internal/harnessobs/opencode_event_factory.go`
  - `internal/harnessobs/opencode_family_map.go`
  - `internal/harnessobs/opencode_family_set.go`
  - `internal/harnessobs/opencode_family_core.go`
  - `internal/harnessobs/opencode_family_tools.go`
  - `internal/harnessobs/opencode_family_execution.go`
  - `internal/harnessobs/opencode_family_workflow.go`
  - `internal/harnessobs/opencode_family_order.go`
  - `internal/harnessobs/opencode_observed_at.go`
  - `internal/harnessobs/opencode_actor.go`
  - `internal/harnessobs/session_command_facts.go`
  - `internal/harnessobs/session_command_event.go`
  - `internal/harnessobs/session_command_model.go`
  - `internal/harnessobs/session_command_time.go`
  - `internal/harnessobs/normalized_event_factory.go`

## Decision Gate

- Simpler/Faster: move OpenCode normalization and session command event
  functions into responsibility-named files in the same package without
  changing logic.
- Blocking Edge Cases: family-detection behavior is part of evidence
  normalization, so it stays split by rule group rather than hidden in one broad
  file. Raw signal traversal and token safety are excluded for later focused
  review.
- Existing Open Source: not applicable; this is local Go file locality cleanup
  using existing normalization logic and project quality gates.

## Planned Verification

- focused `internal/harnessobs` test and changed-file MI gate with `70.1`
  threshold
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI baseline gates with `70.1` MI thresholds
- three independent staged-diff reviewer lanes after implementation
