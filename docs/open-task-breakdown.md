# Open Task Breakdown

Status: working breakdown, current as of 2026-06-01.

This file decomposes open governance closure tasks one by one. It is a working
map for closure, not approval, merge authorization, production trust, or
external attestation.

## Per-Task Breakdown

| Task | What It Requires | Current Evidence | Can Close Locally? | Next Action |
| --- | --- | --- | --- | --- |
| Spec 022 T022-T027 | Update D006, inspect this breakdown for contradictions, record residual remediation state, and run focused US2 review. | Complete in `docs: record spec 022 governance decision state`; US2 review is recorded in `specs/022-post-merge-governance-closure/reviews/us2-review.md`. | Closed for docs-governance scope. | No further US2 action; continue with US3 and final checks. |
| Spec 022 T028-T033 | Synchronize roadmap, decision ledger, spec reality ledger, and spec status wording. | US2 is complete; this navigation update is in progress. | Yes. | Complete synchronized navigation update in the same closure-surface change. |
| Spec 022 T034-T042 | Run local checks, forbidden-claim review, drift/quality review, adversarial PR-ready review, and final checkbox update. | Pending US2/US3. | Partially; live PR/CI final-head checks require GitHub after PR. | Run after all story edits and focused reviews. |

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
success at head `537612723f7dad7f6e3a92657336fafa4ce6bbd5`. T226 is closed
by the current OpenCode/GSD-Redux first-run observation path:
`opencode run --command gsd-plan-phase` and `--command gsd-execute-phase`
under `sdp-trace observe session` produced setup metadata, command digest,
process id, source commit, time bounds, output and normalized digests, and
model/interaction/tool/phase evidence. `mutation` and `test` remain explicit
`not_assessed` dimensions because the target phase was already complete and the
execute workflow performed no mutation or test action.

Spec 002 T035 and Spec 004 T042 are now closed by explicit maintainer merge
approval in the 2026-05-26 closure thread, PR #64 merge commit
`e129515e7c4c7a4a9c4b2b53eb4d3694b41eb2bd`, and post-merge `main` CI run
`26453881873`.

Spec 005 T014 is now closed by maintainer decision `accept_gap` in the
2026-05-26 closure thread. The decision accepts the historical approval gap for
the reviewed Product Contract v0 because later implementation is represented by
Spec 006; it does not create production approval, release approval, or external
trust.

Spec 006 T003 is now closed by maintainer decision `accept_gap` in the
2026-05-26 closure thread. The decision accepts the historical
pre-implementation approval gap after implementation, review, merge, and
post-merge CI evidence; it does not retroactively invent approval, approve
production readiness, or create external trust.

Spec 007 T008 is now closed by maintainer decision
`approve_existing_demo_track` in the 2026-05-26 closure thread. The decision
approves the existing `fall-out-bug/sdp-trace-demo-jvm-gsd` route, setup PR #16,
feature PRs #16-#20 as first slices, negative draft PR #21 for theater
explanation, and no separate polished sales repository for this closure. It
does not approve release, production trust, compliance, semantic quality, or
signed external trust.

Spec 018 T018-001 and T018-070 are now closed by maintainer decision
`approve_core_extension_direction` in the 2026-05-26 closure thread plus filed
follow-up Spec 020 `core-query-package-split` and Spec 021
`source-file-locality-cleanup`. The follow-up specs are draft implementation
surfaces and do not approve command removal, separate binaries, production
readiness, release approval, or external trust.

Spec 022 is now active governance closure work. Follow-up specs 020 and 021
remain prepared draft implementation surfaces and are not active closure debt
until explicitly taken into work. No additional successor spec beyond active
Spec 022 is currently identified for the Spec 019 residual governance closure.

<!-- sdp-trace-claim: claim=profile_passed; subject=open-task-breakdown; state=pass; profile=open_tasks_classified; evidence=state:claim_tags_consistent -->
