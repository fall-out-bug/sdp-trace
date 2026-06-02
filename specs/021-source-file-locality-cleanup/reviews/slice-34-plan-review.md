# Slice 34 Plan Review: Raw Unsafe Rule Semantics

Status: passed

## Proposed Slice

- Scope: `internal/harnessobs/harnessobs_223` through
  `internal/harnessobs/harnessobs_244`.
- Intended grouping:
  - `unsafe_value_maps.go` for map traversal, child path rendering, raw-event
    skip rules, and map field reason codes.
  - `unsafe_value_traversal.go` for traversal entrypoints, shared dispatch,
    and slice traversal.
  - `unsafe_value_strings.go` for string reason selection, path checks, and
    token checks.
  - `unsafe_value_urls.go` for authenticated URL checks.
  - Additional cohesive safety-rule files for path, token, URL, raw path-like,
    raw skip, and digest-field exemptions when required by MI gates.
- Explicitly excluded: validation enum helpers (`harnessobs_245` onward),
  session collect option validation, event source resolution, and raw-event
  normalization flow.

## Review Lanes

- lane 1 requirements/scope: `LGTM`
- lane 2 trust/evidence: `LGTM`
- lane 3 maintainability/DX: `LGTM`

## Findings

- none

## Implementation Refinement

- Initial four-file grouping failed file-level MI for dense rule groups. Final
  implementation keeps rule-level cohesion while grouping traversal/slices,
  map traversal, raw-event skip rules, string path/token checks, URL checks,
  raw path-like exemptions, and digest/token exemptions without retaining the
  one/two-helper slice/path/token microfiles.
