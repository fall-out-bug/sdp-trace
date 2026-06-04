# Slice 86 Plan Review

Date: 2026-06-04

Scope: `internal/packet` Markdown rendering shards `packet_193` through
`packet_200`.

## Plan Under Review

Slice 86 completes the numbered-file cleanup inside `internal/packet` by moving
packet Markdown rendering helpers into cohesive responsibility-named locality
files:

- `packet_render.go`: public render entrypoint, validation-before-render
  projection, and top-level Markdown assembly.
- `render_packet_sections.go`: executive summary, metadata table, metadata
  fields, and required-row rendering.
- `render_theater_helpers.go`: existing theater rendering locality, extended
  with the theater section dispatcher.

The slice intentionally excludes `internal/prreview` numbered files.

## Behavior Contract

Preserve:

- validation-before-render behavior and joined validation errors;
- top-level section order and packet disclaimer wording;
- executive summary rows and metadata field order;
- Markdown escaping through existing `md` helper;
- required-row sorting through `requiredRowIndex`;
- gap fallback from empty reason to `none`;
- clean-theater fallback state/reason from `PC-THEATER`;
- theater finding rendering;
- package boundary, dependency direction, CRAP, and MI baselines.

## Review Lanes

- Lane A scope/behavior review: `LGTM`.
- Lane B locality/boundary review: `LGTM`.
- Lane C tests/evidence review round 1: major finding. T021-5941 did not
  explicitly require focused evidence for joined validation errors, metadata
  field order, required-row sorting, and empty gap fallback. Fixed by adding
  `TestPacketRenderingSectionHelpersPreserveMetadataRowsAndErrors` to the
  exact focused test list and expanding the evidence contract.

Pending:

- Lane C tests/evidence re-review: `LGTM`.
