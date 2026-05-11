# Socratic Review - Integration Sequencing & DX Feasibility

**Reviewer role**: Integration sequencing and DX feasibility reviewer
**Package**: `specs/003-agent-supply-chain-roadmap/reviews/2026-05-10-socratic-review-packet.md`
**Date**: 2026-05-10

---

## 1. Verdict

**APPROVE_FOR_USER_REVIEW**, with strong recommendations to resolve the findings listed below before any slice moves into implementation. The roadmap is honest, well-scoped, and properly treats itself as a discovery/planning artifact. The identified sequencing and feasibility gaps are tractable; none are fatal to the overall direction.

---

## 2. Top Socratic Questions

1. **Which concrete CTO packet format will be the **first** deliverable?**
   The spec keeps PR comment, downloadable archive, static HTML, Markdown report, and CLI summary as open options. Without a decision, the adapter teams will build toward different targets. Pick one (e.g. static Markdown report) and align the P0-A exit criteria.

2. **Can we obtain and replay a complete OpenCode + GSD session today, or is there a risk that the toolchain has changed enough to invalidate the assumptions in the research notes?**
   The roadmap hinges on this dogfood path. If the surface has shifted, the P0-C slice will need a different starting point.

3. **How will the wrapper/sidecar evidence mode be implemented and maintained for OpenCode + GSD?**
   The plan lists both import and wrapper, but the tasks only detail mapping of native events. A missing wrapper design task leaves the "wrapper/sidecar" mode undefined.

4. **What is the fallback if neither Hermes nor OpenClaw exposes a stable, inspectable event/session API?**
   The research gap is honest, but the tasks assume a spike will be feasible. Without a fallback, the P1-B slice may stall indefinitely.

5. **What constitutes a "stable" session/export path for pi and GSD2, and who signs off on that determination?**
   The success criterion SC-003 uses the word "stable" but never defines it. Ambiguity here could make the `not_assessed` gate unenforceable.

6. **Should the GitHub change-host adapter (P0-B) be built concurrently with the CTO packet shape (P0-A), or strictly after?**
   They share a contract; sequencing them in serial may be slower. A concurrent design pass with an interface contract could reduce total calendar time.

---

## 3. Findings Table

| id | severity | cited file:line | finding | why it matters | exact fix |
|----|----------|-----------------|---------|----------------|-----------|
| F01 | major | spec.md:270-272 | CTO packet format is an unresolved open question; the downstream slices (P0-A/B/C) imply a uniform output shape. | Adapter outputs will be designed without a target, leading to rework or multiple incompatible packet representations. | Before any P0-A implementation, choose a single format (suggest Markdown report as the zero-dependency MVP) and record the decision in spec.md. |
| F02 | major | plan.md:188-196, tasks.md:44-52 | The plan requires identifying the exact evidence surface for each tool before adapter work, but Phase 3 tasks (T016-T020) jump directly to mapping without a prior step to **verify** the current OpenCode/GSD surface. | If the tool's output format has drifted since the research notes were written, mapping will fail or produce wrong evidence. | Add an explicit task (e.g., T015bis) to run a fresh OpenCode/GSD session and validate the JSONL shape against documented expectations before T016. |
| F03 | major | tasks.md:66-72 (Phase 4), research.md:164-172 | Tasks T021-T025 for pi/GSD2 discovery are listed after the OpenCode + GSD packet work, but a lightweight discovery spike on pi could be done earlier to inform the core packet schema. | The current order may cause the packet model to be over-fitted to OpenCode's event structure; discovering pi's potentially simpler surface later could force a schema redesign. | Allow pi/GSD2 discovery to start in parallel with (or immediately after) P0-A; keep the packet schema extensible. Alternatively, add a note that the packet design must remain provider-neutral beyond the first example. |
| F04 | medium | plan.md:84-94 (integration table), tasks.md (entire) | The "wrapper/sidecar" evidence mode for OpenCode + GSD is mentioned but never materialised in the tasks. | Developers will not know what to build, how to test it, or where it lives. The missing component undermines the P0 adoption path claim (FR-005). | Add a concrete task (e.g., T019bis) to design and implement a minimal wrapper/sidecar (likely a shell script or small Go binary) that captures OpenCode stdout/stderr and enriches it with GSD phase markers. |
| F05 | medium | research.md:88-92, tasks.md:74-81 (Phase 6) | The Hermes/OpenClaw boundary spike is predicated on a stable event/session API, which the research gap lists as unknown. No fallback is provided. | The P1-B slice could be blocked indefinitely, wasting roadmap shelf space. | Define a contingency: if no stable API is found, the spike should be replaced by a design exploration that documents what an ideal general-agent event surface would look like, and the slice stays `not_assessed`. |
| F06 | minor | AGENTS.md:45-46 (Decomposition Rule) | The rule about <100 lines/10 skills is a good self-check, but the roadmap does not estimate expected size of the resulting Go modules. | It is possible that a future adapter module (e.g., GitHub adapter) could breach the rule without the team noticing. | In the plan's exit criteria for P0-B, add a note that the Go package implementing the adapter must stay within the decomposition rule, or it must be split. |
| F07 | minor | spec.md:256-259 (SC-003) | "Stable session/export evidence path" is not defined. | Different stakeholders could interpret "stable" as "works once" or "documented API with versioning", leading to premature promotion of `not_assessed` rows. | Add a short definition in the spec: e.g., "Stable" means the export format is documented, versioned, and replayable from a fresh install without manual intervention. |

---

## 4. Missing Evidence / `not_assessed` Areas

The research.md (lines 162-172) explicitly lists the unknowns-pi export format, GSD2 state database, Superpowers artifact stability, Hermes/OpenClaw event API, minimal packet shape, and private signing equivalent. These are appropriate for a discovery roadmap. **No new unknown is hidden**; the team should be commended for making these gaps visible.

However, the tasks do not yet contain a step to **validate the external sources** before acting on them. This gap is partially closed by the plan's review gates, but it would be safer to add a lightweight "verify tool surface" task immediately before each integration phase.

---

## 5. Scope-Control Risks

- **Actor-laundering** (evidence theater taxonomy) remains a conceptual risk until a concrete agent-chain example is traced. The roadmap correctly leaves it as a detection target.
- **"Scope theater"** (overclaiming support from one observation) is addressed by FR-013 and by the explicit "not_assessed" discipline. The main residual risk is that discovery fatigue could pressure the team into marking items "supported" after only a shallow check. The existing review gates are a strong countermeasure.
- **General-purpose agent boundary** (spec.md lines 137-165) is well bounded, but without an API, the boundary may remain theoretical. The roadmap should resist the temptation to build a "full connector" before the spike proves feasibility.

---

## 6. One Strongest Reason to Proceed

The roadmap delivers a **complete, honest, and phased plan** to take `sdp-trace` from a single GitHub PR evidence packet to a multi-tool agent supply-chain evidence layer, while never overclaiming trust, support, or product readiness. The CTO-facing value is crisp, the evidence semantics are explicit, and the decomposition into safe slices respects the repo's strict constraints.

---

## 7. One Strongest Reason Not to Proceed Yet

**The absence of a chosen CTO packet format** and the **lack of a proven stable evidence surface for the first dogfood toolchain (OpenCode/GSD)** create a risk that the first implementation slice will be built on sand. If the packet format selection and tool-surface verification are delayed until after roadmap approval, the team may waste time on a design that must be thrown away. Resolving these two items (via F01 and F02) before the human approval of the first implementation slice would remove the largest single risk.

---

**Bottom line**: Approve the roadmap for user review. Require the packet-format decision and an OpenCode/GSD surface-verification spike as pre-conditions for any implementation task. The other findings can be addressed as the slices are scheduled.
