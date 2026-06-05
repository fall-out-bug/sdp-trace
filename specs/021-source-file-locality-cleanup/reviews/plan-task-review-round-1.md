# Spec 021 Plan/Task Review Round 1

Scope: `AGENTS.md`, `specs/021-source-file-locality-cleanup/spec.md`,
`plan.md`, and `tasks.md`.

Date: 2026-06-01.

## Lanes

| Lane | Harness | Model | Result |
| --- | --- | --- | --- |
| Kimi | `opencode run` | `kimi-for-coding/k2p6` | LGTM with one minor observation about possible MI baseline residue. |
| Z.AI | `opencode run` | `opencode-go/glm-5.1` | LGTM with minor observations about target filename, CRAP task coverage, and wording. |
| MiniMax | `opencode run` | `opencode-go/minimax-m3` | Not LGTM; one major finding about pre-checked confirmation tasks and one major finding about premature `in_progress` status. |

## Retained Findings

| ID | Severity | Finding | Disposition |
| --- | --- | --- | --- |
| R1-M1 | major | Confirmation tasks were pre-checked by the author before review evidence existed. | Fixed by reverting T021-001, T021-002, and T021-010 to unchecked before implementation. They were re-checked only after Slice 1 review and evidence artifacts existed. |
| R1-M2 | major | Spec, plan, and tasks used `in_progress` before implementation began. | Fixed by using `ready_for_review` before implementation, then returning to `in_progress` once Slice 1 implementation started. |
| R1-m1 | minor | Target file name `command_surface_registry_helpers.go` was too narrow because some helpers are general command-surface metadata constructors. | Fixed by splitting the target into `cmd/sdp-trace/command_surface_registry_helpers.go` and `cmd/sdp-trace/command_surface_metadata_helpers.go`. |
| R1-m2 | minor | Final evidence task did not name a review artifact path. | Fixed by naming `specs/021-source-file-locality-cleanup/reviews/slice-1-evidence.md`. |
| R1-m3 | minor | CRAP/MI verification was not explicit in tasks. | Fixed by adding an explicit CRAP/MI gate task and plan note. |

## Current Decision

The first slice remains intentionally small: consolidate nine command-surface
helper shards in `cmd/sdp-trace` into one behavior-named helper file. No package
boundary, command behavior, output contract, schema, or dependency-direction
change is planned.
