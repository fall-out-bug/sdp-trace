# Slice 20 Plan Review: Harnessobs Raw Signals And Key Lookup

Status: pass

## Scope

- Package: `internal/harnessobs`
- Source shards: `harnessobs_065` through `harnessobs_091`
- Target responsibility groups:
  - raw signal dispatch and structure traversal
  - raw map/slice/string/scalar signal extraction
  - exact and prefix signal matching
  - recursive key presence traversal
  - generic key lookup across maps and slices
  - string and numeric key extraction
  - timestamp string, numeric, and fallback conversion

## Decision Gate

- Simpler/Faster: move the existing helper functions into responsibility-named
  Go files in the same package without changing logic.
- Blocking Edge Cases: recursive traversal and timestamp extraction are used by
  OpenCode normalization, so behavior changes would affect evidence family
  detection. The slice therefore excludes token safety, run loading, validation
  loading, and event scanning.
- Existing Open Source: not applicable; this is local Go file locality cleanup
  using existing traversal logic and project quality gates.

## Planned Verification

- focused `internal/harnessobs` test and changed-file MI gate with `70.1`
  threshold
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI baseline gates with `70.1` MI thresholds
- three independent staged-diff reviewer lanes after implementation
