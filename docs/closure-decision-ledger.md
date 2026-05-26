# Closure Decision Ledger

Status: decision ledger, current as of 2026-05-26.

This ledger tracks the remaining non-local decisions blocking full spec closure.
It is not approval, merge authorization, production trust, release readiness, or
external attestation.

For one-row-per-task detail, see
[`docs/open-task-breakdown.md`](open-task-breakdown.md).

## Current State

| Evidence | State |
| --- | --- |
| Task ledger | 595 / 605 checked; 10 open |
| PR surface | PR #64 merged as `e129515e7c4c7a4a9c4b2b53eb4d3694b41eb2bd`; post-merge `main` CI run `26453881873` passed |
| Merge approval | explicit maintainer approval recorded in the 2026-05-26 closure thread |
| Maintainer approval | `not_assessed` unless explicitly recorded per row below |
| External demo repo | `fall-out-bug/sdp-trace-demo-jvm-gsd` is the active demo repository for current closure |
| First-run GSD route | GSD-Redux local replacement works through `opencode run --command`; model/interaction/tool/phase observed; mutation/test remain `not_assessed` for the no-op completed phase |

## Decision Rows

| ID | Open tasks | Decision needed | Current evidence | Acceptable evidence to close | Current state |
| --- | --- | --- | --- | --- | --- |
| D001 | 002 T035, 004 T042 | Approve PR #64 merge, then verify post-merge state. | Explicit maintainer approval recorded in the 2026-05-26 closure thread; PR #64 merged as `e129515e7c4c7a4a9c4b2b53eb4d3694b41eb2bd`; post-merge `main` CI run `26453881873` passed. | Explicit maintainer merge approval, successful merge, final `main` CI evidence, post-merge verification note. | `accepted_closed` |
| D002 | 005 T014 | Decide whether Product Contract v0 historical approval gap is accepted, waived, rejected, or split. | Contract and later packet implementation evidence exist; historical approval evidence is not represented. | Maintainer decision recorded with one of: `accepted_gap`, `waived`, `rejected`, or `split_successor`. | `not_assessed` |
| D003 | 006 T003 | Decide whether Change Evidence Packet Core historical pre-implementation approval gap is accepted, waived, rejected, or split. | Implementation, review, local verification, and PR evidence exist; historical pre-implementation approval evidence is not represented. | Maintainer decision recorded with one of: `accepted_gap`, `waived`, `rejected`, or `split_successor`. | `not_assessed` |
| D004 | 007 T008 | Decide demo-track option, first implementation slice, and demo-repo strategy. | External demo packet evidence, v1 baseline tag, and buyer rehearsal exist; explicit demo-track approval is absent. | Maintainer decision approving, rejecting, or splitting demo-track direction. | `not_assessed` |
| D005 | 018 T018-001, T018-070 | Review core/extension direction and decide whether follow-up specs may be prepared. | Machine review, command stability matrix, ownership map, docs, and local verification exist; maintainer review is absent. | Maintainer review outcome plus follow-up specs if direction is approved. | `not_assessed` |
| D006 | 019 T019-001/T002/T003/T004/T120 | Decide post-merge governance for PR #60/Spec 019: accept partial merge state, reject it, or split remaining debt. | PR #60 merged with missed gates; PR #63 final-head CI passed; post-merge closure plan exists. | Maintainer decision recorded against the closure plan. | `not_assessed` |
| D008 | 001 T226 | Close current OpenCode/GSD-Redux route evidence so first-run delivery-loop evidence can be observed. | GSD-Redux local replacement works through `opencode run --command gsd-plan-phase` and `--command gsd-execute-phase`; `sdp-trace observe session` records setup metadata, command digest, process id, source commit, time bounds, output/normalized digests, and model/interaction/tool/phase evidence. Mutation/test remain `not_assessed` because the replayed phase is already complete and execute-phase performs no mutation or test action. | Closed as observed first-run route with unavailable dimensions retained; no mutation, test success, feature delivery, harness compliance, PR approval, merge approval, or production trust is claimed. | `accepted_closed` |

## Closure Rule

Do not close any row by inference. Green CI, local verification, checked task
boxes, review artifacts, or this ledger can support readiness for a decision,
but they are not the decision itself.

<!-- sdp-trace-claim: claim=profile_passed; subject=closure-decision-ledger; state=pass; profile=decision_surface_recorded; evidence=state:claim_tags_consistent -->
