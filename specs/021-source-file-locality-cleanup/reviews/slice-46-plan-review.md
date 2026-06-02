# Slice 46 Plan Review

Status: pass

## Scope

Slice 46 is bounded to `internal/harnessobs/harnessobs_348` through
`internal/harnessobs/harnessobs_360`.

Planned consolidation:

- `raw_normalization.go`: orchestration, input validation, and zero-time
  fallback.
- `raw_normalization_scan.go`: OpenCode raw JSONL file opening and scanner
  loop.
- `raw_normalization_line.go`: raw line decoding and unsafe raw-event
  rejection.
- `raw_normalization_digest.go`: normalized source digesting.
- `raw_normalization_writer.go`: normalized event output file creation and
  JSONL writing.
- `event_line_parsing.go`: shared blank JSONL line handling used by both normal
  event parsing and raw normalization.

Explicit exclusions:

- generic unsafe raw-value discovery helpers already in non-numbered files
- OpenCode event construction helpers already in non-numbered files
- changes to path policy, raw event schema, output JSONL format, package
  boundary, dependency direction, or MI baselines

## Decision Gate

- Simpler/Faster: Keep the numbered raw-normalization helpers. Rejected because
  they are the next ordered source-file locality debt and currently split one
  behavior across thirteen tiny files.
- Blocking Edge Cases: The slice must preserve supported-format gating,
  raw/output same-file rejection, zero-time fallback, scanner limits, shared
  blank line skipping, malformed JSONL errors, unsafe-input rejection, normalized
  source digest calculation, output parent creation, and JSONL write format.
- Existing Open Source: No new dependency is introduced. The existing stdlib
  `bufio.Scanner`, `encoding/json`, and package-local safety/digest helpers
  remain sufficient.

## Reviewer Lanes

- scope reviewer: `019e8802-b340-76c0-9bda-a82a60138e95`, result `LGTM`.
- trust/evidence reviewer: `019e8802-d261-7f40-b9ed-5df9b15fa92e`,
  result `LGTM`.
- maintainability/DX reviewer: `019e8804-23d6-7e21-8324-a678aab03636`,
  initial result `major`; `019e8805-d8c4-78f3-946e-e1a63b7af3ec`,
  re-review result `LGTM`.

## Findings

- fixed: assigned shared `blankJSONLLine` to neutral `event_line_parsing.go`
  instead of a raw-normalization-specific file.
