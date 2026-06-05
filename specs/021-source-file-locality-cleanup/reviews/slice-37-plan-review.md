# Slice 37 Plan Review

Status: pass

## Scope

Slice 37 is bounded to `internal/harnessobs/harnessobs_281` through
`internal/harnessobs/harnessobs_290`.

Planned consolidation:

- `validation_command.go`: `Validate` orchestration, profile loading call site,
  and optional output writing.
- `validation_requirements.go`: required validation options and shared
  non-blank requirement.
- `validation_paths.go`: safe profile/run/out path resolution.
- `validation_evaluation.go`: `LoadRun`-backed run evaluation and
  source-unavailable cannot-verify fallback construction.

Explicit exclusions:

- profile loading implementation (`harnessobs_291` onward)
- session profile validation
- raw-event config validation
- isolation rule validation and installation
- loaded session run validation
- path-safety primitives

## Decision Gate

- Simpler/Faster: Keep the current one-helper numbered shards. Rejected because
  it preserves decomposition debt in a user-visible validation command path.
- Blocking Edge Cases: Validation is trust-sensitive because it may emit
  cannot-verify verdicts and write validation artifacts. The slice must
  preserve required-option messages, path-safety errors, optional output
  semantics, non-authority boundary text, digest generation, and fallback state.
- Existing Open Source: No new parser, workflow engine, JSON writer, or
  dependency is introduced. Existing package-local path, JSON, digest, and
  evaluation helpers remain the implementation substrate.

## Reviewer Lanes

- scope reviewer: check slice boundary and exclusions.
- trust/evidence reviewer: check evidence mapping and no overclaiming.
- maintainability/DX reviewer: check grouping avoids numbered microfiles and
  avoids new non-numbered one-helper drift.

## Findings

- initial scope lane (`019e8788-8f70-7cb2-bc40-bdcb718f7d0c`): found that
  the initial plan wording appeared to include the `LoadProfile`
  implementation, which is outside Slice 37.
- initial trust/evidence lane (`019e8788-fbf2-7bd3-9b10-7f410835ad77`): found
  missing focused evidence mapping for the non-authority boundary.
- initial maintainability/DX lane (`019e8789-17d9-7720-a091-23ccd5febec3`):
  found ambiguous ownership between run evaluation and fallback construction.
- scope re-review (`019e878a-4226-76d3-82c0-22ecea63317e`): LGTM
- trust/evidence re-review (`019e878a-629b-7bb0-ade8-5867a8641499`): LGTM
- maintainability/DX re-review (`019e878a-9b18-7b93-bc79-0a83552cf9b3`):
  LGTM

Final plan review verdict: LGTM across all three lanes. Implementation remains
bounded to `harnessobs_281` through `harnessobs_290`.

Implementation grouping update: local MI rejected a single
`validation_inputs.go`, so validation requirements and path resolution were
split into separate responsibility-named files within the same approved scope.

- maintainability/DX grouping update re-review
  (`019e878f-0f9e-7a83-9690-a95e918f2402`): LGTM
