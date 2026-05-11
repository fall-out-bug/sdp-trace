# Socratic Review: Theater and Scope

**Reviewer role**: Adversarial theater and scope reviewer
**Date**: 2026-05-10

---

```
Verdict: REVISE_BEFORE_USER_APPROVAL
```

The direction is correct. Existing repo, honest history, real evidence gaps - this is the right product problem. But four major findings need resolution before the user is asked to approve, because they affect what the demo will actually claim and how proof vs. hand-authoring is gated.

## Findings

| id | severity | file/section | finding | exact fix |
| --- | --- | --- | --- | --- |
| SR-001 | major | spec.md / Core Demo Claim | "through an OSS agent/harness workflow" is generic. It will be quoted as broad ecosystem support despite the non-claim list. | Change to "through one demonstrated OSS agent/harness route (OpenCode + GSD + MiniMax-M2.5)" to match actual demo scope. |
| SR-002 | major | spec.md / Negative / Theater Demo | "trigger `agent_claimed_verification` or `ci_theater`" is ambiguous for a controlled demo. A CTO asks which one and why. | Pick one: change to "trigger `agent_claimed_verification`" for the first negative example. Note that `ci_theater` may appear in a future negative example. |
| SR-003 | major | spec.md / Success Criteria | No criterion requires any happy-path PR to reach `PC-THEATER: pass`. If all features show theater, the negative PR loses contrast. | Add SC-006: "At least one happy-path feature PR has `PC-THEATER: pass`, demonstrating that the theater assessment produces a clean result when evidence is properly bound." |
| SR-004 | major | plan.md / Workstream A vs Non-Goals | Workstream A describes Go implementation but Non-Goals says "Do not implement packet generation in this spec slice." No prerequisite gate links A completion to B packet quality. Without it, demo could ship hand-authored packets as proof. | Add explicit sequencing note to plan.md: "T009-T015 must complete before T020-T025 can produce generated proof packets. Hand-authored packets are acceptable only for T019 as template validation." Add prerequisite annotations to T020-T025 in tasks.md. |
| SR-005 | minor | demo-repo-plan.md / Decision Ownership | The demo implies enterprise ceremony (risk acceptance, security review) for a Todo app. | Add note: "For the demo's selected profile, risk acceptance and security review owners may be `not_in_scope`." |
| SR-006 | minor | spec.md / Demo Track Policy; demo-repo-plan.md | "Polished public sales demo" is undefined and drifts toward marketing scope. | Replace "sales demo" with "public polished demo." Add to Non-Goals: "This demo is a CTO product proof, not a sales or marketing asset." |
| SR-007 | minor | tasks.md / Phase 2-3 | No dependency edges between T009-T015 and T020-T025. | Add prerequisite annotations to T020-T025: "Requires T009-T015 or documented hand-authoring justification." |

## Existing repo default

Yes. Continuing in `sdp-trace-demo-jvm-gsd` is the right default. The messy history - review contamination, missing PR lifecycle evidence, partial OpenCode/GSD observation - is exactly the product problem `sdp-trace` claims to solve. A clean repo proves only that a clean repo looks clean. The CTO needs to see that real imperfect agent-delivery history becomes legible through the packet. The alternative clean repo is defensible only for a later polished demo after product proof exists.

## Smallest first slice

Tag current demo repo `main` as `demo-v1-observation-baseline`. Select one existing feature PR or history slice (the cleanest of the five). Produce one hand-authored Change Evidence Packet v0 from it with all required rows, honest gaps, and a minimal bundle manifest. Review the packet for CTO readability against the Product Contract v0 template. This validates the packet template against real evidence before any Go implementation begins, and it forces an early answer to whether existing repo evidence is legible enough to organize.

## Claims the demo must still refuse

- The demo must refuse to claim that one generated packet proves the system works for all agent/harness/CI combinations.
- The demo must refuse to claim that the negative PR's theater finding proves the system catches all theater patterns.
- The demo must refuse to present hand-authored packets as generated proof once the product generator exists.
- The demo must refuse to claim the existing repo's observation history is a product strength unless a packet actually organizes it legibly.
- The demo must refuse to claim OSS ecosystem compatibility, tool support, or agent readiness beyond the one demonstrated OpenCode + GSD + MiniMax-M2.5 route.
- The demo must refuse to claim that `PC-DECISION: owner_bound` means the change is approved, ready, or safe to merge.
- The demo must refuse to generalize the Todo app's CI/theater patterns to enterprise delivery contexts.
