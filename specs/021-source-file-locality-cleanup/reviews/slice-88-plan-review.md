# Slice 88 Plan Review

Date: 2026-06-04

Scope: `internal/prreview` packet construction and packet input reference
shards `prreview_020` through `prreview_037`.

## Plan Under Review

Slice 88 moves packet construction and packet input reference helpers into
cohesive locality files:

- `packet_build.go`: public packet build entrypoint, prepared-directory
  orchestration, packet finalization/write, packet construction, identity,
  provenance, and defaulting.
- `packet_refs.go`: `packetRefs` aggregation and copied diff/metadata/context/
  verification reference collection.
- `packet_unavailable_fields.go`: explicit unavailable-field generation for
  absent optional packet inputs.

The slice intentionally excludes run execution (`prreview_038` onward),
ledger/validation logic, packet option validation (`prreview_087` onward),
generic input-copy utilities, prompt generation, and lower-level IO helpers.

## Behavior Contract

Preserve:

- validation before output directory creation;
- new-output-directory enforcement;
- packet digest and `packet.json` write behavior;
- packet ID format and provenance fields;
- default `CreatedBy` and CI state;
- copied ref kinds/content types;
- context kind inference and verification kind rewriting;
- metadata optionality;
- unavailable-field state/reason strings;
- package boundary, dependency direction, CRAP, and MI baselines.

## Review Lanes

- Lane A scope/behavior review: `LGTM`.
- Lane B locality/boundary review: `LGTM`.
- Lane C tests/evidence review round 1: major finding. T021-6081 did not
  explicitly require focused evidence for new-output-directory enforcement.
  Fixed by extending the focused evidence contract to require rejection of an
  already-existing output directory without overwriting existing contents.

Pending:

- Lane C tests/evidence re-review: `LGTM`.
