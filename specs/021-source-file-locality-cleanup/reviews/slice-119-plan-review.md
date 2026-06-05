# Slice 119 Plan Review

Date: 2026-06-05T02:07:55Z

Scope:

- Consolidate `internal/authority/authority_event_type.go` into an authority
  event-validation file.
- Preserve standard event membership and `custom:` event acceptance.
- Exclude event-set validation, pre-decision blocker ordering, target-rule
  validation, authority evaluation behavior, schemas, examples, dependencies,
  package boundary, dependency direction, CRAP/MI baselines, and public
  surfaces.

Initial review state: fail

| Reviewer | Harness | Agent ID | Model/provider | Prompt class | Result |
|---|---|---:|---|---|---|
| Beauvoir | Codex subagent | 019e9406-f078-7fd2-b8d0-e22ac17a1e3a | not_assessed | plan/task review | major finding |
| Peirce | Codex subagent | 019e9406-f40c-79f1-904e-54d0f0b73866 | not_assessed | plan/task review | major finding |
| Halley | Codex subagent | 019e9406-f7c2-7f80-80d9-86f7cf7e0c22 | not_assessed | plan/task review | major finding |
| Laplace | Codex subagent | 019e9587-5f30-7923-b764-d723befad229 | not_assessed | plan/task review | major finding |

Findings:

- major: `T021-8280` claimed focused verification would preserve `custom:`
  event acceptance, but the guarded tests only proved unsupported diagnostics
  and standard-event paths. A regression removing the `custom:` prefix rule
  from `validEventType` could pass the planned focused guard.

Fix:

- Updated `T021-8280` to require focused coverage for unsupported diagnostics,
  standard event acceptance, and `custom:` event acceptance through the two
  guarded tests. If current coverage lacks a `custom:` acceptance case, one of
  those two tests must be extended before implementation review; adding a third
  focused test name requires re-review.

Re-review state: pass

| Reviewer | Harness | Agent ID | Model/provider | Prompt class | Result |
|---|---|---:|---|---|---|
| Beauvoir | Codex subagent | 019e9406-f078-7fd2-b8d0-e22ac17a1e3a | not_assessed | plan/task re-review | LGTM |
| Peirce | Codex subagent | 019e9406-f40c-79f1-904e-54d0f0b73866 | not_assessed | plan/task re-review | LGTM |
| Halley | Codex subagent | 019e9406-f7c2-7f80-80d9-86f7cf7e0c22 | not_assessed | plan/task re-review | LGTM |

Post-review MI finding:

- major: consolidating `validEventType` into `authority_vars.go` made
  `authority_vars.go` fail the file-MI gate without an allowed baseline:
  `file MI baseline missing for below-threshold file
  internal/authority/authority_vars.go`.

Fix:

- Updated Slice 119 to consolidate `validEventType` into
  `authority_event_set_validation.go`, co-locating the helper with event-type
  validation callers while keeping `authority_vars.go` as the source of
  `standardEventTypes`.

Second re-review state: pass

| Reviewer | Harness | Agent ID | Model/provider | Prompt class | Result |
|---|---|---:|---|---|---|
| Beauvoir | Codex subagent | 019e9406-f078-7fd2-b8d0-e22ac17a1e3a | not_assessed | plan/task re-review | LGTM |
| Peirce | Codex subagent | 019e9406-f40c-79f1-904e-54d0f0b73866 | not_assessed | plan/task re-review | LGTM |
| Halley | Codex subagent | 019e9406-f7c2-7f80-80d9-86f7cf7e0c22 | not_assessed | plan/task re-review | LGTM |
