# Slice 22 Plan Review: Harnessobs Event Scanning And Validation

Status: pass

## Scope

- Package: `internal/harnessobs`
- Source shards: `harnessobs_111` through `harnessobs_132`
- Target responsibility groups:
  - event source file opening and JSONL scanning
  - scan limits and source hashing
  - line parsing and safe raw-event rejection
  - typed event decoding
  - parsed-event digest validation
  - event identity, reference, and content validation
  - unavailable-field validation

## Decision Gate

- Simpler/Faster: move existing functions into responsibility-named Go files in
  the same package without changing logic.
- Blocking Edge Cases: source hashing, line limits, safe-field rejection, and
  digest mismatch errors are trust-sensitive replay boundaries. The slice
  therefore excludes evaluation/dimension composition and path-safety helpers.
- Existing Open Source: not applicable; this is local Go file locality cleanup
  using existing scanning and validation logic.

## Planned Verification

- focused `internal/harnessobs` test and changed-file MI gate with `70.1`
  threshold
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI baseline gates with `70.1` MI thresholds
- three independent staged-diff reviewer lanes after implementation
