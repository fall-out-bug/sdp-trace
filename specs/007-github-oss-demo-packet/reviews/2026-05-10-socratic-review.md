# Socratic Review Ledger: GitHub OSS Demo Packet

Date: 2026-05-10
Scope: `specs/007-github-oss-demo-packet/`
Status: Prior Socratic review completed for the combined product/demo draft;
split into 006 product core and 007 demo. Focused re-review of the split is
pending.

## Review Runs

| plane | model | status | raw output |
| --- | --- | --- | --- |
| Product proof | `openrouter/qwen/qwen3.6-plus` | usable | `raw/2026-05-10-qwen-product-proof.md` |
| Demo truth and contamination | `openrouter/deepseek/deepseek-v4-pro` | usable | `raw/2026-05-10-deepseek-demo-truth.md` |
| GitHub evidence | `openrouter/xiaomi/mimo-v2.5-pro` | usable | `raw/2026-05-10-xiaomi-github-evidence.md` |
| Theater and scope | `zai/glm-5.1` | usable | `raw/2026-05-10-glm-theater-scope.md` |
| Fix-check: product proof | `openrouter/qwen/qwen3.6-plus` | usable for combined draft | terminal capture, 2026-05-10 |
| Fix-check: GitHub evidence | `openrouter/xiaomi/mimo-v2.5-pro` | usable for combined draft | terminal capture, 2026-05-10 |
| Fix-check: theater and scope | `openrouter/deepseek/deepseek-v4-pro` | usable for combined draft | terminal capture, 2026-05-10 |

`pi` emitted a startup warning that `kimi-coding/k2p6` did not match the local
model registry. Kimi was not assigned a reviewer slot in this cycle.

## Verdict

`REVIEW_REFRESH_REQUIRED_AFTER_SPLIT`

All reviewers agreed on the direction: continue in the existing
`sdp-trace-demo-jvm-gsd` repository as the default. The existing repo's messy
history is the product problem. A clean new repository would be useful later
for a polished public demo, but it should not replace the first proof that
`sdp-trace` can make real agent-delivery evidence legible.

The original combined draft was not approval-ready because it lacked a minimum bar
for the first CTO-visible packet, did not pin the negative theater target, and
left retroactive GitHub artifact/provenance rules too loose.

Focused re-review accepted the combined-draft fixes. After the split, the main
design change is stricter: 007 no longer allows hand-authored CTO-visible demo
packets because 006 now owns the pre-renderer fixture path. This split needs a
focused re-review before asking for approval.

Xiaomi raised three minor clarification items; all three were applied, with the
second item superseded by the 006/007 split:

- excluded contaminated features must not appear on CTO-visible packet/tracker
  surfaces unless the contamination is surfaced;
- if no existing feature meets the first-packet bar, 007 creates a new v2 PR
  only after the 006 renderer/validator exists;
- "exits `not_assessed`" is defined inline as assessed `pass`, `partial`, or
  `fail`.

## Dispositions

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| 006-P1 | major | product proof | Product tasks T009-T015 lacked P0 classification fields. | `accepted_fixed` | `tasks.md` now maps T009-T015 to packet rows, evidence surface, start state, transition, buyer effect, and non-goal. |
| 006-P2 | major | product proof | No row-state threshold for demo success. | `accepted_fixed` | `spec.md` now defines first-packet minimum bar: four rows pass/partial, `PC-CHANGE` and `PC-MUTATION` resolvable, one of verification/review/agent-route exits `not_assessed`. |
| 006-P3 | major | product proof | Negative demo reason code was ambiguous. | `accepted_fixed` | `spec.md` and `demo-repo-plan.md` now make `agent_claimed_verification` the primary first negative example. |
| 006-P4 | minor | product proof | Locally-generated packets could become permanent proof. | `accepted_fixed_superseded_by_split` | 006 owns pre-renderer fixture validation; 007 requires CTO-visible demo packets to be renderer-generated and validator-checked. |
| 006-P5 | minor | product proof | Evidence contradiction rule missing. | `accepted_fixed` | `spec.md` now marks contradictory surfaces as `partial` and requires both refs in residual gaps. |
| 006-P6 | minor | product proof | PR comment lifecycle/canonicality missing. | `accepted_fixed` | `plan.md` now says uploaded packet artifact is canonical and PR comment/body link is informational. |
| CT-001 | critical | demo truth | First packet target unspecified. | `accepted_fixed` | `spec.md` now has feature evidence inventory and first-packet selection/minimum-bar rules. |
| CT-002 | critical | demo truth | Feature 4 contamination not handled in spec. | `accepted_fixed` | `spec.md` now requires contamination surfaced in residual gaps/theater or explicit exclusion; Feature 4 cannot be first packet unless surfaced. |
| CT-003 | critical | demo truth | First packet could be honest but useless. | `accepted_fixed` | First-packet minimum bar added. If no existing feature meets it, first packet comes from a new v2 PR. |
| CT-004 | major | demo truth | Reset decision rule lacked "too noisy" threshold. | `accepted_fixed` | `plan.md` reset table switches to fresh root branch only if all existing features fail first-packet rule and minimum bar. |
| CT-005 | major | demo truth | Negative PR could poison first buyer impression. | `accepted_fixed` | `spec.md` and `demo-repo-plan.md` require `DEMO-NEGATIVE:` prefix, `demo-theater` label, separate handling, and showing it after happy-path packet. |
| CT-006 | minor | demo truth | Setup PR could green without proving artifact path. | `accepted_fixed` | `plan.md` requires setup-only artifact or packet validation command; setup is not feature proof. |
| CT-007 | minor | demo truth | Feature evidence inventory missing. | `accepted_fixed` | `spec.md` now includes initial inventory table. |
| G1 | critical | GitHub evidence | Retroactive task binding could backdate provenance. | `accepted_fixed` | `spec.md` adds provenance rules for PR body/task artifact binding and caps retroactive initiator confidence at `partial`. |
| G2 | critical | GitHub evidence | GitHub artifact retention/expiry ignored. | `accepted_fixed` | `spec.md` requires 180-day retention for v2 packet artifacts and `artifact_expired` handling for retroactive packets. |
| G3 | major | GitHub evidence | First packet selection criteria missing. | `accepted_fixed` | First-packet evidence-richness selection rule added. |
| G4 | major | GitHub evidence | Draft PR and retained fixture treated as equivalent negative demo. | `accepted_fixed` | Negative CTO demo must be a GitHub draft PR; fixtures are tests only. |
| G5 | major | GitHub evidence | PR-visible packet pointer optional. | `accepted_fixed` | `plan.md` requires PR comment or PR body link for demo PRs. |
| G6 | major | GitHub evidence | Historical CI citation vs fresh CI was ambiguous. | `accepted_fixed` | `spec.md` now distinguishes retroactive artifact availability and v2 retention. |
| G7 | major | GitHub evidence | Setup/first packet ordering ambiguous. | `accepted_fixed` | `plan.md` clarifies setup PR is infrastructure only; first packet follows selection rule. |
| G8 | minor | GitHub evidence | Packet/bundle path convention unclear. | `accepted_fixed` | `demo-repo-plan.md` defines flat packet files and directory-per-feature bundles. |
| G9 | minor | GitHub evidence | Hand-authored/tool-generated transition unclear. | `accepted_fixed_superseded_by_split` | 006 owns pre-renderer fixtures; 007 depends on 006 renderer/validator. |
| G10 | minor | GitHub evidence | "generate" overclaimed before renderer exists. | `accepted_fixed_superseded_by_split` | 007 starts after renderer completion, so CTO-visible demo packets are tool-generated. |
| SR-001 | major | theater/scope | Core claim implied broad OSS support. | `accepted_fixed` | `spec.md` now says one demonstrated route: OpenCode + GSD + MiniMax-M2.5. |
| SR-002 | major | theater/scope | Negative reason code ambiguous. | `accepted_fixed` | Primary negative reason code fixed to `agent_claimed_verification`. |
| SR-003 | major | theater/scope | No happy-path clean theater criterion. | `accepted_fixed` | Added SC-006 requiring at least one happy-path `PC-THEATER: pass`. |
| SR-004 | major | theater/scope | Product prerequisites not linked to generated proof. | `accepted_fixed_superseded_by_split` | 006 is now a separate prerequisite spec; 007 starts only after 006 renderer/validator behavior exists. |
| SR-005 | minor | theater/scope | Todo app implied enterprise ceremony. | `accepted_fixed` | `demo-repo-plan.md` says risk/security owners may be `not_in_scope`. |
| SR-006 | minor | theater/scope | Sales-demo wording drift. | `accepted_narrower` | Existing wording keeps "polished public sales demo" only as later non-default option; current demo remains CTO product proof. |
| SR-007 | minor | theater/scope | Task dependency edges missing. | `accepted_fixed_superseded_by_split` | 007 tasks now begin by confirming 006 availability and no longer contain product implementation tasks. |

## Remaining Risk

No known critical or major finding remains intentionally open from the combined
draft review. The split itself still needs focused re-review because it changes
sequencing and removes the pre-renderer hand-authored demo path.

## Current Recommended First Slice

1. Complete 006 Change Evidence Packet Core first.
2. Keep the existing `sdp-trace-demo-jvm-gsd` repository.
3. Tag current `main` as `demo-v1-observation-baseline`.
4. Add packetization setup PR with artifact path and tracker only.
5. Select the first packet target by evidence inventory and minimum bar.
6. If no existing feature meets the bar, create one new v2 feature PR under the
   packetization track.
7. Produce one CTO-readable Change Evidence Packet v0 before attempting the
   full five-feature demo.
