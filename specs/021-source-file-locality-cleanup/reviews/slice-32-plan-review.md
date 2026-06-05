# Slice 32 Plan Review: Correct Serialization Microfile Drift

Status: passed

## Proposed Slice

- Scope: corrective consolidation of `internal/harnessobs` microfiles around
  artifact JSON writing, event references, and digest helpers.
- Intended grouping:
  - `artifact_json.go` for JSON data rendering and file writing.
  - `event_refs.go` for event ref rendering, event ref path safety, and
    event artifact path safety.
  - `event_ref_validation.go` for retained event-ref validation checks; this
    stays separate from rendering/path safety because the combined event-ref
    file failed file-level MI while this split remains cohesive.
  - `digest_helpers.go` for SHA-256 hex rendering, source line/file digesting,
    validation digesting, and command digesting.
- Explicitly excluded: command parser, session setup, session collection, raw
  event unsafe rule semantics, and numbered shards `harnessobs_217` onward.

## Review Lanes

- lane 1 requirements/scope: `LGTM`
- lane 2 trust/evidence: `LGTM`
- lane 3 maintainability/DX: `LGTM`

## Findings

- none

## Implementation Refinement

- The reviewed consolidation direction stayed unchanged. During implementation,
  `event_refs.go` plus retained event validation failed file-level MI, so event
  validation remains in one cohesive validation file instead of returning to
  multiple one-helper event-ref microfiles.
