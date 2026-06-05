# Slice 101 Plan Review

Date: 2026-06-05

Scope under review:
- `internal/prreview/prreview_169` through `internal/prreview/prreview_182`.
- Unsafe text, unique string, command digest, context/content type, normalized
  extension, safe ID, and default string helpers only.

Decision gate:
- Simpler/Faster: regroup existing helper predicates and mappings by
  responsibility; no sanitizer, prompt rendering, or schema behavior rewrite.
- Blocking Edge Cases: redaction triggers, safe ID fallback, digest hashing,
  content/context mappings, and extension normalization feed evidence refs and
  must not drift.
- Existing Open Source: no new library is justified; current implementation uses
  small Go string/path/hash helpers.

Plan review lanes:
- Halley the 2nd: `LGTM`.
- Beauvoir the 2nd: major finding that T021-6991 relied on broad tests without
  requiring exact coverage for NUL-separated command digest hashing,
  whitespace-only `defaultString` fallback, normalized extension allow-list and
  `.txt` fallback, and full content-type mapping. Resolution: added T021-6981
  requiring `TestPrreviewSmallHelpersPreserveContracts`, and added that test to
  the T021-6991 exact-count guard.
- Peirce the 2nd: major finding matching the helper-contract coverage gap for
  exact command digest, whitespace `defaultString`, and normalized extension
  allow-list/fallback assertions. Resolution: same T021-6981/T021-6991 update.
