# Slice 40 Plan Review

Status: pass

## Scope

Slice 40 is bounded to `internal/harnessobs/harnessobs_310` through
`internal/harnessobs/harnessobs_334`.

Planned consolidation:

- `session_isolation_install.go`: isolation rule installation orchestration,
  target resolution, single-rule install, and installer dispatch.
- `session_isolation_paths.go`: profile-relative isolation path normalization,
  parent containment, and filename validation.
- `session_isolation_lines.go`: line-rule read/append/write helpers.
- `session_isolation_json.go`: JSON read-deny rule creation and mutation.
- `session_isolation_json_object.go`: optional object loading, blank-object
  parsing, and missing-file defaulting.
- `session_isolation_readback.go`: readback verification, pass/cannot-verify
  state handling, digest assignment, and result construction.
- `session_isolation_presence.go`: line/JSON readback presence checks. JSON
  readback intentionally calls existing package-local JSON reader helpers that
  remain numbered for a later slice.

Explicit exclusions:

- loaded session run validation (`harnessobs_335` through `harnessobs_336`)
- shared JSON read/decode helpers (`harnessobs_337` through `harnessobs_339`)
- session run construction (`harnessobs_340` onward) and source commit
  discovery
- event source reading and profile-relative source/output file safety
- raw-event normalization execution

## Decision Gate

- Simpler/Faster: Keep the current one-helper numbered shards. Rejected because
  it preserves decomposition debt in filesystem mutation/readback code and
  blocks the ordered cleanup objective.
- Blocking Edge Cases: Isolation installation mutates local files and then
  derives verification evidence. The slice must preserve path containment,
  filename safety, line idempotence, JSON object preservation, readback
  cannot-verify semantics, digest assignment, and unsupported-kind failures.
- Existing Open Source: No new filesystem, JSON, or policy dependency is
  introduced. Existing package-local path, JSON, digest, and installer helpers
  remain the implementation substrate.

## Reviewer Lanes

- scope reviewer: check slice boundary and exclusions.
- trust/evidence reviewer: check evidence mapping and no overclaiming.
- maintainability/DX reviewer: check grouping avoids numbered microfiles and
  avoids new non-numbered one-helper drift.

## Findings

- initial scope reviewer (`019e87b0-fdbd-74f3-b4f5-744e6e3c72c0`): major
  finding. The initial exclusion wording incorrectly described all
  `harnessobs_335` onward as loaded session run validation even though
  `harnessobs_337` through `harnessobs_339` are shared JSON read/decode helpers
  used by JSON readback.
- initial trust/evidence reviewer (`019e87b1-211e-7503-963b-e72d3eb61dd9`):
  major finding. The initial focused evidence task covered missing/blank line
  reads and line writing but omitted line-read error propagation.
- maintainability/DX reviewer (`019e87b6-01fc-7f71-a3be-7f734bcf8654`): minor
  finding. `T021-2691` originally read as if profile-relative isolation path
  safety was omitted rather than preserved.

## Re-review Results

- scope re-review (`019e87b5-118b-7e13-b27f-96d5f7bf60e4`): `LGTM`.
- trust/evidence re-review (`019e87b5-1669-75b1-8f1d-73b955786249`): `LGTM`.
- maintainability/DX re-review (`019e87b6-01fc-7f71-a3be-7f734bcf8654`):
  `LGTM`.
