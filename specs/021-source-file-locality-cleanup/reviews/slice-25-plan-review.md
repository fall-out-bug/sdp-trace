# Slice 25 Plan Review: Harnessobs Output File Paths

Status: pass

## Scope

- Package: `internal/harnessobs`
- Source shards: `harnessobs_150` through `harnessobs_161`
- Target responsibility groups:
  - output file path validation
  - output basename validation
  - parent path normalization and rejection
  - missing-parent resolution
  - symlink-aware parent containment
  - working-directory relative path conversion

## Decision Gate

- Simpler/Faster: move existing functions into responsibility-named Go files in
  the same package without changing logic.
- Blocking Edge Cases: output paths create or overwrite local artifacts, so
  traversal rejection, basename validation, symlink resolution, missing-parent
  fallback, and containment checks must stay behavior-compatible. The slice
  therefore excludes output directory materialization helpers.
- Existing Open Source: not applicable; this is local Go file locality cleanup
  using existing path-safety logic and project quality gates.

## Planned Verification

- focused `internal/harnessobs` test and changed-file MI gate with `70.1`
  threshold
- focused output path regression evidence for traversal rejection, basename
  validation, missing-parent resolution, symlink resolution, working-directory
  containment, and successful validation output writing
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI baseline gates with `70.1` MI thresholds
- three independent staged-diff reviewer lanes after implementation, with every
  finding fixed and the affected lane repeated until it returns exactly `LGTM`
