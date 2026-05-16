# Review Synthesis: 015-spec-governance-and-roadmap

## Review objective
Verify that the spec governance slice improves roadmap clarity without turning checked-in prose into false authority.

## Review planes

| Plane | Model | Harness | Status |
| --- | --- | --- | --- |
| Architecture doubt | GLM-5.1 | zai | Completed |
| Wide-context code/docs | Qwen-3.6-max-preview | openrouter | Completed |
| Skill adherence / workflow | MiniMax-M2.7 | minimax | 404 (endpoint unavailable) |
| Reasoning | DeepSeek-v4-pro | openrouter | Timeout without output (300s) |

## Findings summary

| # | Severity | Finding | Source | Disposition | Fix |
| --- | --- | --- | --- | --- | --- |
| 1 | P0 | Roadmap status contradicts spec.md for 002–007 (`historical` vs `Draft`) | GLM #1 | accepted_fixed | Changed roadmap to `draft` with note; added "Older Draft Specs" section |
| 2 | P0 | `blocks/` directories don't exist for 002–007 | GLM #2 | accepted_fixed | Removed `blocks/` assertion from roadmap |
| 3 | P1 | `historical` undefined in spec.md taxonomy | GLM #3 | accepted_fixed | Added `historical` to US-002 |
| 4 | P1 | Zero claim tags in new files | GLM #4, Qwen #2 | accepted_fixed | Added tags to `docs/roadmap.md`, `final-evidence-map.md`, `tasks.md` |
| 5 | P1 | `final-evidence-map.md` closure overclaim ("proven") | GLM #5, Qwen #4 | accepted_fixed | Changed to "What the evidence shows"; removed PR #53 reference |
| 6 | P1 | GitHub Actions claim without tag/replay | GLM #6, Qwen #4 | accepted_fixed | Removed PR-specific claim; referenced branch head only |
| 7 | P2 | doccheck overclaimed as coverage verifier | GLM #7, Qwen #3 | accepted_fixed | T030 now states "doccheck verifies link integrity only" |
| 8 | P1 | Task checkboxes closed without proof cycle | GLM #8 | accepted_fixed | Reopened all checkboxes, ran verification, closed again |
| 9 | P2 | Stale "Last updated" date | GLM #9, Qwen #1 | accepted_fixed | Updated to 2026-05-16 |
| 10 | P2 | Blocker notation not applied to roadmap tables | GLM #10, Qwen #5 | accepted_fixed | Applied `→ Blocked on:` format to Blocked specs table |
| 11 | P2 | Capability index hides cross-cutting ownership | GLM #11, Qwen #7 | accepted_fixed | Added caveat: "A capability may be touched by multiple specs" |
| 12 | P3 | docs/README.md duplicate numbering | GLM #12, Qwen #9 | accepted_fixed | Fixed numbering 7→10, 8→11 |
| 13 | P1 | Status transitions undefined makes taxonomy decorative | GLM #13, Qwen #6 | accepted_fixed | Added "Minimal Transitions" section to roadmap |
| 14 | P2 | No existing specs updated to reference taxonomy | GLM #14 | rejected_false_positive | Out of scope per FR-005: "Avoid rewriting historical evidence packages" |
| 15 | P2 | `socratic-judge-result.json` ambiguity | Qwen #8 | rejected_false_positive | Spec 001 artifact; exempt per US-004 scope rule |

## Unresolved states
- **cannot_verify**: MiniMax-M2.7 (404 endpoint); DeepSeek-v4-pro (300s timeout, zero output)
- **not_assessed**: (none)

## What this review does not prove
- That roadmap statuses will stay accurate (depends on manual curation).
- That all specs 001–014 will adopt the taxonomy (out of scope for this slice).
- That the claim-tag grammar is machine-parseable (no parser yet).

## Re-verification after fixes
- `go test ./...`: PASS
- `go run ./tools/doccheck`: PASS (exit=0)
- `git diff --check`: PASS
- `go vet ./...`: PASS
- `jq empty schema/*.json`: PASS

<!-- sdp-trace-claim: claim=profile_passed; subject=015-review-synthesis; state=pass; profile=repo_baseline_structural; evidence=command_set:block015-t030 -->

## Round 2 (post-fix) Review

GLM performed a second-pass review on the v2 artifact (post-fix). Qwen v2 returned no output. MiniMax and DeepSeek remain `cannot_verify`.

| # | Severity | Finding | Disposition | Fix |
| --- | --- | --- | --- | --- |
| F1 | P1 | Spec 001 dual-listed (`blocked` + `historical`) in roadmap | accepted_fixed | Removed from Historical section; kept in Blocked with `blocks/` note |
| F2 | P2 | Unstaged diff in `final-evidence-map.md` | accepted_fixed | Committed PR #54 reference change |
| F3 | P2 | Untracked `glm-review-v2.txt` | accepted_fixed | Committed as review evidence |
| F4 | P3 advisory | `command_set:` references lack machine-readable registry | acknowledged | Out of scope for Slice 1; tracked for later slice |
| F5 | P3 advisory | Transition rules are unenforceable prose | acknowledged | Documented as editorial conventions, not enforced gates |

### Qwen v2 (post-fix) findings

| # | Severity | Finding | Disposition | Fix |
| --- | --- | --- | --- | --- |
| Q1 | P0 | Spec 001 `blocked` vs spec.md `Draft` | accepted_fixed | Changed roadmap status to `draft` (blocked) with blocker note |
| Q2 | P0 | Spec 008 `blocked` vs spec.md `Implemented locally` | accepted_fixed | Changed roadmap status to `in_progress` (blocked) with blocker note |
| Q3 | P1 | Claim tag evidence refs don't match `command_set:block-task` grammar | accepted_fixed | Replaced with `command_set:block015-t030` |
| Q4 | P1 | `not_assessed: (none)` contradicts synthesis gaps | accepted_fixed | Added three `not_assessed` items to final-evidence-map |
| Q5 | P2 | `synthesis.md` verification assertions without claim tag | accepted_fixed | Added claim tag to synthesis.md |
