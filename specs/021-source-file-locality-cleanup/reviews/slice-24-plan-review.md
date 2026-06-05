# Slice 24 Plan Review: Harnessobs Existing Path Safety

Status: pass

## Scope

- Package: `internal/harnessobs`
- Source shards: `harnessobs_143` through `harnessobs_149`
- Target responsibility groups:
  - existing file/directory entrypoints
  - traversal and URL-like path rejection
  - symlink and absolute-path resolution
  - working-directory containment
  - expected file/directory type validation

## Decision Gate

- Simpler/Faster: move existing functions into responsibility-named Go files in
  the same package without changing logic.
- Blocking Edge Cases: these helpers protect filesystem trust boundaries, so
  traversal rejection, symlink resolution, and type checks must stay byte-for-byte
  behavior-compatible. The slice therefore excludes output creation helpers.
- Existing Open Source: not applicable; this is local Go file locality cleanup
  using existing path-safety logic and project quality gates.

## Planned Verification

- focused `internal/harnessobs` test and changed-file MI gate with `70.1`
  threshold
- focused path-safety regression evidence for traversal rejection, URL-like
  rejection, symlink resolution, working-directory containment, and expected
  file/directory type checks
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI baseline gates with `70.1` MI thresholds
- three independent staged-diff reviewer lanes after implementation, with every
  finding fixed and the affected lane repeated until it returns exactly `LGTM`
