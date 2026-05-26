# Closure Decision Ledger

Status: decision ledger, current as of 2026-05-26.

This ledger tracks the remaining non-local decisions blocking full spec closure.
It is not approval, merge authorization, production trust, release readiness, or
external attestation.

## Current State

| Evidence | State |
| --- | --- |
| Task ledger | 585 / 605 checked; 20 open |
| PR surface | PR #64 open, green, and `CLEAN` at `0a20e01ed69aa336b0184c71347c79e6f9590e82` |
| Merge approval | `not_assessed` |
| Maintainer approval | `not_assessed` unless explicitly recorded per row below |
| External demo access | `cannot_verify` for `fall-out-bug/sdp-trace-demo-ci-pilot` with current GitHub and SSH access |
| First-run GSD route | `not_assessed` for tool/phase/mutation/test because current OpenCode/GSD lacks `/gsd-plan-phase` |

## Decision Rows

| ID | Open tasks | Decision needed | Current evidence | Acceptable evidence to close | Current state |
| --- | --- | --- | --- | --- | --- |
| D001 | 002 T035, 004 T042 | Approve PR #64 merge, then verify post-merge state. | PR #64 is open, green, and `CLEAN`; PR review evidence exists; explicit merge approval is absent. | Explicit maintainer merge approval, successful merge, final `main` CI evidence, post-merge verification note. | `not_assessed` |
| D002 | 005 T014 | Decide whether Product Contract v0 historical approval gap is accepted, waived, rejected, or split. | Contract and later packet implementation evidence exist; historical approval evidence is not represented. | Maintainer decision recorded with one of: `accepted_gap`, `waived`, `rejected`, or `split_successor`. | `not_assessed` |
| D003 | 006 T003 | Decide whether Change Evidence Packet Core historical pre-implementation approval gap is accepted, waived, rejected, or split. | Implementation, review, local verification, and PR evidence exist; historical pre-implementation approval evidence is not represented. | Maintainer decision recorded with one of: `accepted_gap`, `waived`, `rejected`, or `split_successor`. | `not_assessed` |
| D004 | 007 T008 | Decide demo-track option, first implementation slice, and demo-repo strategy. | External demo packet evidence, v1 baseline tag, and buyer rehearsal exist; explicit demo-track approval is absent. | Maintainer decision approving, rejecting, or splitting demo-track direction. | `not_assessed` |
| D005 | 018 T018-001, T018-070 | Review core/extension direction and decide whether follow-up specs may be prepared. | Machine review, command stability matrix, ownership map, docs, and local verification exist; maintainer review is absent. | Maintainer review outcome plus follow-up specs if direction is approved. | `not_assessed` |
| D006 | 019 T019-001/T002/T003/T004/T120 | Decide post-merge governance for PR #60/Spec 019: accept partial merge state, reject it, or split remaining debt. | PR #60 merged with missed gates; PR #63 final-head CI passed; post-merge closure plan exists. | Maintainer decision recorded against the closure plan. | `not_assessed` |
| D007 | 001 T211-T217 | Restore demo-pilot repository evidence path or split Block 25 external demo closure. | `fall-out-bug/sdp-trace-demo-ci-pilot` is not resolvable through current GitHub or SSH access. | Accessible demo repo, replayed CI/artifact evidence, role reviews, and local/PR-level verification; or explicit successor split. | `cannot_verify` |
| D008 | 001 T226 | Restore or replace current OpenCode/GSD route so first-run delivery-loop evidence can be observed. | `minimax/MiniMax-M2.5` route observes setup/model/source/digest facts, but `/gsd-plan-phase` is unavailable and tool/phase/mutation/test remain `not_assessed`. | Live first-run observation with setup metadata, command digest, source commit, time bounds, output/normalized digests, and delivery-loop evidence; unavailable dimensions explicitly retained. | `not_assessed` |

## Closure Rule

Do not close any row by inference. Green CI, local verification, checked task
boxes, review artifacts, or this ledger can support readiness for a decision,
but they are not the decision itself.

<!-- sdp-trace-claim: claim=profile_passed; subject=closure-decision-ledger; state=pass; profile=decision_surface_recorded; evidence=state:claim_tags_consistent -->
