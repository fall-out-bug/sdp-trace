# Open Task Breakdown

Status: working breakdown, current as of 2026-05-26.

This file decomposes the remaining 12 open task boxes one by one. It is a
working map for closure, not approval, merge authorization, production trust, or
external attestation.

## Per-Task Breakdown

| Task | What It Requires | Current Evidence | Can Close Locally? | Next Action |
| --- | --- | --- | --- | --- |
| 002 T035 | Merge only after fresh CI, local verification, PR review, and post-merge verification. | PR #64 is open, green, and `CLEAN`; merge approval and post-merge verification are absent. | No. | Requires explicit merge approval, merge, final `main` CI query, and post-merge verification note. |
| 004 T042 | Stop before merge unless explicit merge approval is present. | PR #64 is green/CLEAN; explicit merge approval remains `not_assessed`. | No. | Requires explicit merge approval or an explicit decision to keep the task open. |
| 005 T014 | Stop for explicit user approval of reviewed Product Contract v0. | Contract and later packet implementation evidence exist; historical approval evidence is not represented. | No. | Maintainer decides `accepted_gap`, `waived`, `rejected`, or `split_successor`. |
| 006 T003 | Get explicit user approval before implementation. | Implementation, review, local verification, and PR evidence exist; historical pre-implementation approval is not represented. | No. | Maintainer decides `accepted_gap`, `waived`, `rejected`, or `split_successor`. |
| 007 T008 | Ask for explicit approval of demo-track option, first slice, and demo-repo strategy. | External demo packet evidence, v1 baseline tag, and buyer rehearsal exist; explicit demo-track approval is absent. | No. | Maintainer approves, rejects, or splits demo-track direction. |
| 018 T018-001 | Review core command list with maintainers. | Machine review and command-surface verification exist; maintainer review is `not_assessed`. | No. | Maintainer reviews command list or explicitly keeps review pending. |
| 018 T018-070 | Prepare follow-up implementation specs only after core/extension direction is approved. | Follow-up specs are blocked on T018-001. | No. | After maintainer approval, prepare follow-up specs; otherwise keep open. |
| 019 T019-001 | Review Spec 019 scope with maintainers. | PR #60 merged before this approval was recorded; this is a missed pre-merge gate. | No. | Maintainer decides whether to accept, reject, or split the already-merged state. |
| 019 T019-002 | Run adversarial spec review before implementation approval. | Cross-model review ran after implemented slices; pre-implementation approval evidence is absent. | No. | Treat as `partial_after_merge`; maintainer decides whether to accept or split. |
| 019 T019-003 | Start implementation/Pi handoff only after reviewed direction is approved. | Already missed before merge; cannot be retroactively satisfied as pre-implementation approval. | No. | Maintainer records acceptance/waiver/rejection/split for missed gate. |
| 019 T019-004 | Post-merge approval gate: maintainers approve partial merge state and closure plan or reject/split remaining work. | `post-merge-closure-plan.md` exists; maintainer decision is `not_assessed`. | No. | Maintainer approves, rejects, or splits the closure plan. |
| 019 T019-120 | Close Spec 019 only after review findings, local verification, and live CI evidence are recorded. | Review and CI evidence exist, but this task is blocked on T019-004 maintainer decision. | No. | Complete only after T019-004 has an explicit outcome. |

## Local Closure Candidate

T213 is closed through the active v2 packet/bundle track: demo PR #25 is merged
to `sdp-trace-demo-jvm-gsd` `main` as
`3a9491f734e5214c72014db5d893f125eb254a11`, local Bazel verification passed,
and downloaded artifacts from demo CI run `25724386343` replayed successfully
with the new verifier. T214 is closed through demo PR #26, merged to
`sdp-trace-demo-jvm-gsd` `main` as
`a4d1f755552ba1f411af5edcb7d6caf24a9c39bf`; GitHub run `26447797437`
confirmed the Block 25 negative matrix in `build-and-test`. T215 is closed by
`docs/reviews/block25-jvm-gsd-demo-sanitized-report.md`, which records the
sanitized report and artifact summary without upgrading residual trust states.
T216 is closed by active-demo role reviews in
`specs/001-sdp-trace-time-series-evidence-substrate/blocks/25-compiled-ci-demo-pilot-review-ledger.md`;
one major InfoSec redaction-scan evidence gap was fixed and focused re-review
returned no critical or major findings. T217 is closed by fresh local
`sdp-trace` verification plus PR #64 `verify` and `pr-review-evidence-only`
success at head `ea00be499abe9a211f3fa639be6124863afad36c`. T226 is closed
by the current OpenCode/GSD-Redux first-run observation path:
`opencode run --command gsd-plan-phase` and `--command gsd-execute-phase`
under `sdp-trace observe session` produced setup metadata, command digest,
process id, source commit, time bounds, output and normalized digests, and
model/interaction/tool/phase evidence. `mutation` and `test` remain explicit
`not_assessed` dimensions because the target phase was already complete and the
execute workflow performed no mutation or test action.

The other open tasks need one of:

- explicit maintainer or merge approval;
- a post-merge verification cycle;
- or an explicit successor split / rejection decision.

<!-- sdp-trace-claim: claim=profile_passed; subject=open-task-breakdown; state=pass; profile=open_tasks_classified; evidence=state:claim_tags_consistent -->
