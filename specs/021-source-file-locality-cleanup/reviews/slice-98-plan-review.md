# Slice 98 Plan Review

Date: 2026-06-05

Scope under review:
- `internal/prreview/prreview_134` through `internal/prreview/prreview_148`.
- Preview, prompt digest/ref, copied input, packet digest, and output directory
  helpers only.

Decision gate:
- Simpler/Faster: regroup existing helpers by responsibility; no preview,
  digest, copy, or directory behavior rewrite.
- Blocking Edge Cases: prompt refs, copied input refs, packet digest replay, and
  output directory rejection are evidence-path behavior and must preserve exact
  shapes and error strings.
- Existing Open Source: no new library is justified; current implementation uses
  Go standard-library hashing, JSON, filesystem, and filepath handling.

Plan review lanes:
- Beauvoir the 2nd: `LGTM`.
- Peirce the 2nd: major finding that T021-6781 referenced missing
  `TestBuildPacketCopiesInputsAndComputesStableDigests` and required copied
  input mode/ref/digest assertions that existing tests did not yet provide.
  Resolution: added T021-6771 requiring that focused regression before the
  exact-count guard.
- Halley the 2nd: `LGTM`.
