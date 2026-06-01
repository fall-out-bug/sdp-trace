# Slice 10 Plan Review: Doctor Repo-Observer Locality Cleanup

Status: accepted for implementation.

## Scope

Slice 10 is bounded to the `cmd/sdp-trace` doctor routing, repo-observer
doctor, and install repo-observer command shards. Local doctor report/check
files remain outside this slice.

## Decision Gate

- Simpler/Faster: grouping all doctor files at once would remove more numbered
  files, but would mix repo-observer install behavior with local report checks.
- Blocking Edge Cases: install and doctor profile paths have different exit
  semantics and persisted JSON surfaces; grouping must preserve those contracts.
- Existing Open Source: not applicable; this is a local file-locality cleanup
  over existing command functions.

## Review

- Spec drift: no drift; the slice excludes local doctor report/check shards.
- Product drift: no product behavior change intended.
- Constitution drift: same package and dependency direction.
- CleanArch hex: no new boundary.
- CleanCode/SOLID/DRY/YAGNI: accepted if grouped by command routing,
  repo-observer doctor output, install execution, install args, and install
  options.
- CRAP/MI: split any broader grouping that falls below file MI 70.
