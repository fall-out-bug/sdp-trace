# Slice 26 Plan Review: Harnessobs Output Directory Paths

Status: pass

## Scope

- Package: `internal/harnessobs`
- Source shards: `harnessobs_162` through `harnessobs_173`
- Target responsibility groups:
  - output directory path validation
  - existing output symlink containment
  - missing output parent containment
  - working-directory escape classification
  - empty-or-missing directory validation

## Decision Gate

- Simpler/Faster: move existing functions into responsibility-named Go files in
  the same package without changing logic.
- Blocking Edge Cases: output directories are artifact creation targets, so
  traversal rejection, symlink containment, missing-parent handling, and
  non-empty directory refusal must stay behavior-compatible. The slice
  therefore excludes JSON writing, event refs, and digest helpers.
- Existing Open Source: not applicable; this is local Go file locality cleanup
  using existing output-directory path safety logic and project quality gates.

## Planned Verification

- focused `internal/harnessobs` test and changed-file MI gate with `70.1`
  threshold
- focused output directory regression evidence for traversal rejection,
  existing symlink escape, parent symlink escape, path-escape classification,
  missing/empty directory acceptance, and non-empty directory rejection
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI baseline gates with `70.1` MI thresholds
- three independent staged-diff reviewer lanes after implementation, with every
  finding fixed and the affected lane repeated until it returns exactly `LGTM`
