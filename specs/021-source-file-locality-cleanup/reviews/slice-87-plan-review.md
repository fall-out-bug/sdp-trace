# Slice 87 Plan Review

Date: 2026-06-04

Scope: `internal/prreview` portable schema/type shards `prreview_001` through
`prreview_019`.

## Plan Under Review

Slice 87 starts `internal/prreview` cleanup by moving schema constants,
package validation vars, options, packet refs, profile/run/result types, and
ledger/validation types into cohesive locality files:

- `constants.go`: schema/state/ref/content/redaction/plane/runner/status/
  severity/disposition/coverage/authority constants plus package regex/error
  vars.
- `options.go`: packet and run option structs.
- `packet_types.go`: packet, safe refs, and unavailable fields.
- `review_types.go`: review profiles/roles, previews/runsets/results,
  findings/citations, ledgers, validations, and plane results.

The slice intentionally excludes packet construction (`prreview_020` onward),
run execution, validation logic, summary/rendering, file IO, prompt generation,
and utility helpers.

## Behavior Contract

Preserve:

- exported names and API surface;
- JSON tags and optional field `omitempty` behavior;
- trust-boundary comments on portable structs;
- schema/state/ref/content/redaction/plane/runner/status/severity/disposition/
  coverage/authority constants;
- package regex patterns and error values;
- package boundary, dependency direction, CRAP, and MI baselines.

## Review Lanes

- Lane A scope/behavior review: `LGTM`.
- Lane B locality/boundary review: `LGTM`.
- Lane C tests/evidence review: `LGTM`.
