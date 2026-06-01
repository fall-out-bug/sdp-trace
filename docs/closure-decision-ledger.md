# Closure Decision Ledger

Status: decision ledger, current as of 2026-06-01.

This ledger tracks the remaining non-local decisions blocking full spec closure.
It is not approval, merge authorization, production trust, release readiness, or
external attestation.

For one-row-per-task detail, see
[`docs/open-task-breakdown.md`](open-task-breakdown.md).

## Current State

| Evidence | State |
| --- | --- |
| Task ledger | Spec 022 active with 42 / 42 checked after fixed-handoff PR-ready re-review |
| PR surface | PR #64 merged as `e129515e7c4c7a4a9c4b2b53eb4d3694b41eb2bd`; post-merge `main` CI run `26453881873` passed. PR #60 and PR #63 live state refreshed on 2026-06-01 in Spec 022. |
| Merge approval | explicit maintainer approval for PR #64 recorded in the 2026-05-26 closure thread; PR #60 merge approval remains `not_assessed` |
| Maintainer approval | `not_assessed` unless explicitly recorded per row below |
| External demo repo | `fall-out-bug/sdp-trace-demo-jvm-gsd` is the active demo repository for current closure |
| First-run GSD route | GSD-Redux local replacement works through `opencode run --command`; model/interaction/tool/phase observed; mutation/test remain `not_assessed` for the no-op completed phase |

## Decision Rows

| ID | Open tasks | Decision needed | Current evidence | Acceptable evidence to close | Current state |
| --- | --- | --- | --- | --- | --- |
| D001 | 002 T035, 004 T042 | Approve PR #64 merge, then verify post-merge state. | Explicit maintainer approval recorded in the 2026-05-26 closure thread; PR #64 merged as `e129515e7c4c7a4a9c4b2b53eb4d3694b41eb2bd`; post-merge `main` CI run `26453881873` passed. | Explicit maintainer merge approval, successful merge, final `main` CI evidence, post-merge verification note. | `accepted_closed` |
| D002 | 005 T014 | Decide whether Product Contract v0 historical approval gap is accepted, waived, rejected, or split. | Maintainer decision `accept_gap` recorded in the 2026-05-26 closure thread; contract and later packet implementation evidence exist, and Spec 006 represents implementation. | Maintainer decision recorded with one of: `accepted_gap`, `waived`, `rejected`, or `split_successor`. | `accepted_gap` |
| D003 | 006 T003 | Decide whether Change Evidence Packet Core historical pre-implementation approval gap is accepted, waived, rejected, or split. | Maintainer decision `accept_gap` recorded in the 2026-05-26 closure thread; implementation, review, local verification, PR, merge, and post-merge CI evidence exist. | Maintainer decision recorded with one of: `accepted_gap`, `waived`, `rejected`, or `split_successor`. | `accepted_gap` |
| D004 | 007 T008 | Decide demo-track option, first implementation slice, and demo-repo strategy. | Maintainer decision `approve_existing_demo_track` recorded in the 2026-05-26 closure thread; external demo packet evidence, v1 baseline tag, and buyer rehearsal exist. | Maintainer decision approving, rejecting, or splitting demo-track direction. | `approved_existing_demo_track` |
| D005 | 018 T018-001, T018-070 | Review core/extension direction and decide whether follow-up specs may be prepared. | Maintainer decision `approve_core_extension_direction` recorded in the 2026-05-26 closure thread; follow-up Spec 020 and Spec 021 are filed as draft implementation specs. | Maintainer review outcome plus follow-up specs if direction is approved. | `approved_core_extension_direction` |
| D006 | 019 T019-001/T019-002/T019-003/T019-004; 022 PR-ready evidence complete | Decide post-merge governance for PR #60/Spec 019: accept partial merge state, reject it, or split remaining debt. | Maintainer decision `split_successor` recorded in the 2026-05-26 closure thread; PR #60 merged with missed gates; PR #60 live refresh on 2026-06-01 shows `state=MERGED`, merge commit `657a343a5f310538def9afd509e6c610c713cab0`, CI `verify` success, and empty `reviewDecision`; PR #63 live refresh shows `state=MERGED`, merge commit `1ee2c7af53637c7f43bff4e0e7ef9e34d164908e`, CI `verify` success, and `pr-review-evidence-only` success; Spec 022 plan/tasks/review artifacts are active and reviewed through `specs/022-post-merge-governance-closure/reviews/plan-task-review-round-2.md`; fixed-handoff PR-ready re-review returned LGTM across configured lanes; no additional successor spec beyond active Spec 022 is currently identified for this residual governance closure. | Existing `split_successor` decision plus active Spec 022 evidence, live PR/CI refresh, focused story reviews, local docs checks, final PR-ready review, and upcoming live branch PR/CI evidence. | `split_successor` |
| D008 | 001 T226 | Close current OpenCode/GSD-Redux route evidence so first-run delivery-loop evidence can be observed. | GSD-Redux local replacement works through `opencode run --command gsd-plan-phase` and `--command gsd-execute-phase`; `sdp-trace observe session` records setup metadata, command digest, process id, source commit, time bounds, output/normalized digests, and model/interaction/tool/phase evidence. Mutation/test remain `not_assessed` because the replayed phase is already complete and execute-phase performs no mutation or test action. | Closed as observed first-run route with unavailable dimensions retained; no mutation, test success, feature delivery, harness compliance, PR approval, merge approval, or production trust is claimed. | `accepted_closed` |

## Closure Rule

`D007` is intentionally unused in this ledger snapshot; existing decision IDs
are preserved rather than renumbered.

Do not close any row by inference. Green CI, local verification, checked task
boxes, review artifacts, or this ledger can support readiness for a decision,
but they are not the decision itself.

<!-- sdp-trace-claim: claim=profile_passed; subject=closure-decision-ledger; state=pass; profile=decision_surface_recorded; evidence=state:claim_tags_consistent -->
