# Slice 100 Plan Review

Date: 2026-06-05

Scope under review:
- `internal/prreview/prreview_157` through `internal/prreview/prreview_168`.
- Citation anchor, resolver dispatch, ref matching, safe ref ID lookup, and
  citation location helpers only.

Decision gate:
- Simpler/Faster: regroup existing citation helpers by responsibility; no
  citation semantics rewrite.
- Blocking Edge Cases: citation resolvability gates finding trust, so empty
  anchors, digest-only citations, unknown-ref digest fallback, and exact
  diff/context/verification location rules must not drift.
- Existing Open Source: no new library is justified; current implementation uses
  small Go predicates over packet refs and citation fields.

Plan review lanes:
- Halley the 2nd: `LGTM`.
- Beauvoir the 2nd: major finding that T021-6921 required unknown-ref
  source-digest fallback coverage but not the paired unknown-ref without source
  digest rejection case. Resolution: T021-6921 now explicitly requires unknown
  ref without source digest rejection.
- Peirce the 2nd: major finding that T021-6921 required exact
  diff/context/verification location rules but did not explicitly require
  context-by-hunk, context-by-source-digest, or verification-by-source-digest
  acceptance. Resolution: T021-6921 now explicitly requires those cases.
