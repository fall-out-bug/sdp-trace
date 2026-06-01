# Slice 18 Plan Review: Harnessobs Observe And Session Setup Entrypoints

Status: pass

## Scope

- Package: `internal/harnessobs`
- Source shards: `harnessobs_034` through `harnessobs_045`
- Target files:
  - `internal/harnessobs/observe_entrypoint.go`
  - `internal/harnessobs/observe_options.go`
  - `internal/harnessobs/observe_validation.go`
  - `internal/harnessobs/observe_paths.go`
  - `internal/harnessobs/observation_context.go`
  - `internal/harnessobs/observation_events_writer.go`
  - `internal/harnessobs/observation_prepare.go`
  - `internal/harnessobs/observation_run_factory.go`
  - `internal/harnessobs/observation_source.go`
  - `internal/harnessobs/observation_time.go`
  - `internal/harnessobs/session_setup_entrypoint.go`

## Decision Gate

- Simpler/Faster: move observe and session setup entrypoint functions into
  responsibility-named Go files in the same package without changing logic.
- Blocking Edge Cases: observe path safety, event serialization, source digest
  loading, and session setup validation are user-visible enough to keep focused
  files instead of one broad observe aggregate that fails file MI and weakens
  navigation.
- Existing Open Source: not applicable; this is local file locality cleanup
  using existing Go functions and project quality gates.

## Planned Verification

- focused `internal/harnessobs` test and changed-file MI gate with `70.1`
  threshold
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI baseline gates with `70.1` MI thresholds
- three independent staged-diff reviewer lanes after implementation
