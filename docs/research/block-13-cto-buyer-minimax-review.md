<think>The user wants me to review the roadmap (Block 13 and its roadmap blocks 13A-21) from the perspective of the CTO Buyer persona defined in the persona file. Let me analyze this carefully.

**CTO Buyer's Key Position:**
- Already has their own AI SDLC (OpenCode, Superpowers, GSD, internal agents, etc.)
- Does NOT want to replace this existing system
- Wants a CONTROL LAYER that can be added read-only or as a sidecar
- Then potentially tightened at real control points

**CTO Buyer's Success Criteria:**
1. Minimal integration contract
2. Telemetry adapter path for any harness
3. Explicit `missing_telemetry`, not silent pass
4. Query/dashboard surface showing task drift, evidence gaps, scope creep, failed tests, unverified claims

**CTO Buyer's Rejection Criteria:**
- "Agents should follow the envelope"
- "Developers should remember to log"
- Opaque health scores

**CTO Buyer's Review Bias:**
- Prioritize business usefulness and adoption reality
- Reject product framing that requires replacing the buyer's AI SDLC before value appears

Now let me evaluate the roadmap against these criteria:

**POSITIVE ELEMENTS:**
1. The roadmap is gate-first which aligns with the CTO's desire for a control layer
2. Blocks 13A-21 are sequenced logically
3. Non-goals are clearly stated including "no harness-specific product dependency"
4. Explicit `missing_telemetry` states are addressed throughout
5. No opaque health scores (explicit non-goal)
6. Query/dashboard surface is addressed (Block 19, Block 20)
7. Sidecar/read-only adoption path is mentioned in the introduction

**CRITICAL/MAJOR CONCERNS:**

1. **Block 16 "Managed Harness Enforcement Profile"** - This is the most concerning block. The CTO explicitly does NOT want teams to rewrite their harness. The language "fail-closed mode for missing required adapter lifecycle events" and "fail-closed CI profile" sounds like it REQUIRES harness changes or enrollment. 

The "enrollment marker" and "adapter registration and authorization policy" suggest mandatory integration points that could conflict with the minimal integration contract the CTO requires. Block 14's "required-run manifest" and "required evidence contract" similarly creates friction that works against the sidecar-first, read-only adoption model the CTO needs.

2. **Adoption Path Sequencing Issue** - Block 13A comes first, but that's documentation and planning. The first implementation block (Block 14) jumps straight to "protected gate contract" and "fail-closed CI profile" - this is enforcement-heavy for day-one. The CTO wants read-only or sidecar FIRST, then tightened later. The roadmap should clarify: what's the read-only/sidecar path for day 1 vs. the enforcement path for later?

3. **"Managed harness mode" in G2** - This directly conflicts with the persona's rejection criteria. The roadmap needs to distinguish between optional telemetry gathering for any harness versus forced integration for managed harnesses.

4. **Block 17 "Adapter Event Contract"** - The claim that it exposes provenance that CI and git cannot is a strong overclaim. The adapter can only observe what the harness emits through its plugin API. Without that, it's reading file mutations and command traces—the same gaps that exist with CI and git.

5. **Block 20 "Cross-Repository Degradation Export"** - The metrics for evidence quality across repos are a good direction, but the export structure needs clearer definition around what signals indicate actual degradation versus legitimate variation.

The block explicitly avoids issuing a yes/no degradation verdict, which aligns with the requirement not to produce opaque health scores. However, the "time window, dimensions" framing still risks becoming a score-like abstraction.

6. **Missing: "Telemetry adapter path for any harness"** - The gap map acknowledges the "Capture Depth Gap" and "Adapter Event Contract," but there's no clear architectural statement that ANY harness can integrate without modification. The "stable adapter event contract" in G3 and Block 17 addresses this, but the roadmap should explicitly frame adapter integration as a day-one capability rather than a later feature.

7. **Overclaim in Block 19 "Forensics Query Pack"** - The block claims a reviewer can reconstruct task, actors, harness/model identity state, commands, file mutations, tests, missing evidence, redactions, and overrides without reading raw event files. But this assumes complete adapter coverage, which the current plan doesn't guarantee.

The capture depth gaps mean the query output will have holes. The acceptance criteria should acknowledge that some forensics require direct file inspection when adapter coverage is incomplete.

**MINOR OBSERVATIONS:**

Block 21 mentions GitLab CI, Buildkite, and Jenkins profiles for CI witness support, but no equivalent for issue trackers, documentation systems, or design documents—sources where scope creep and drift actually manifest.

**VERDICT:**

The roadmap has merit and aligns with the CTO's core requirements. However, Block 16's "Managed Harness Enforcement Profile" introduces mandatory enrollment and registration that conflicts with the requirement for minimal integration. The sequencing from Block 14 onward pushes toward fail-closed enforcement before establishing a clear read-only path.

The roadmap doesn't explicitly promise "any harness integration without modification" as a day-one constraint, and Block 17 risks overclaiming on forensics given the capture gaps that persist. These are major issues, though they don't necessarily require scrapping the entire plan—just clarifying that the enforcement path is optional and that forensics has explicit limitations.

The "minimal integration contract" and "telemetry adapter path for any harness" need to be highlighted as architectural constraints from the start, with Block 16 explicitly designed to avoid mandatory rewrite requirements rather than just stating it afterward. The capture depth gaps should be acknowledged in Block 17's acceptance criteria, and the roadmap should differentiate between the day-one read-only/sidecar path and the later enforcement path. I'm also uncertain whether "degraded mode for unmanaged harnesses" in Block 16's deliverables aligns with the "explicit `missing_telemetry`, not silent pass" criteria.
</think>

# CTO Buyer Review: Block 13 Roadmap

## VERDICT: REVISE

The roadmap direction is sound, but it contains **two major findings** that directly threaten the CTO's adoption requirements, plus one minor finding that risks overclaim.

---

## Critical/ major Findings

### MAJOR-1: Block 16 Enforcement Profile vs. "No Rewrite Required" Constraint

**Problem:**
Block 16's deliverables include:
- "managed harness profile"
- "wrapper enrollment marker"
- "adapter registration and authorization policy"
- "fail-closed mode for missing required adapter lifecycle events"

This conflicts with the CTO's rejection criteria: *"Agents should follow the envelope."* and the success criteria: *"Minimal integration contract."*

The roadmap claims *"sdp-trace still does not require teams to rewrite their harness"* as an acceptance criterion, but the deliverables above describe exactly the kind of mandatory enrollment and required lifecycle events that would constitute a harness rewrite for teams without a plugin API.

**Evidence:**
- Block 16 G3 explicitly states "unclear unsupported versus missing telemetry states" as a problem to solve—but unclear states become enforced states under a managed profile
- The phrase "fail-closed mode" in Block 16 deliverables means non-cooperative harnesses stop working, not just report degraded posture
- "Adapter registration and authorization policy" implies a required control-plane dependency before any harness can participate

**CTO concern this triggers:**
*"Why should teams rewrite their harness?"*

**Required resolution:**
Block 16 must clearly separate:
1. **Observer-only mode** for unmanaged/any harnesses: reports `missing_telemetry` or `not_integrated` explicitly, never blocks, never requires enrollment
2. **Managed enforcement mode** for willing teams: fail-closed, enrollment required, opt-in only

The day-one value is in mode 1. Mode 2 is the tightening path, not the adoption path. The acceptance criteria should measure "days to first evidence from an unmanaged harness" not "fail-closed rate."

---

### MAJOR-2: Block 14 "Required-Run Manifest" Blocks Sidecar-First Adoption

**Problem:**
Block 14's deliverables include:
- "required-run manifest or required-run section in the expected evidence contract"

The CTO's stated position: *"Can sdp-trace attach read-only or as a sidecar?"* and *"wants a control layer that can be added read-only or sidecar-first, then tightened at real control points."*

A "required-run manifest" is enforcement-first, not observation-first. This is the wrong starting point for a buyer who does not want to change the existing AI SDLC.

**Evidence:**
The roadmap introduction explicitly states "read-only or sidecar-first, then tightened"—but Block 14's first deliverable immediately introduces required-run contracts, which is tightening without the observation phase.

**CTO concern this triggers:**
*"Why should teams rewrite their harness?"*

**Required resolution:**
Block 14 must provide two clearly separated paths:
1. **Observation path**: wrapper observes, produces `local_observed`, no required-run contract, `missing_telemetry` is explicit and visible
2. **Enforcement path**: required-run manifest, fail-closed at CI boundary, explicit opt-in

The observation path is day-one value. The enforcement path is Block 14+.

---

### MINOR-1: Block 17 forensics overclaim

**Problem:**
Block 19 acceptance criteria states: *"reviewer can reconstruct task, actors, harness/model identity state, commands, file mutations, tests, missing evidence, redactions, witnesses, and overrides without reading raw event files"*

This is only true if capture depth (Block 17) achieves full adapter integration. Block 17 explicitly depends on "stable adapter event contract"—which is future work after Block 16 enforcement, Block 18 redaction, and Block 15 signing.

The forensics query pack (Block 19) comes *before* capture depth is complete, but claims forensic completeness.

**CTO concern this triggers:**
*"If the agent does not write sdp-trace telemetry, what is still visible?"*

**Required resolution:**
Block 19 acceptance criteria must acknowledge capture depth limitations explicitly. Add: *"Query output reflects available evidence from the observation profile active at run time. Full reconstruction is only achievable with managed harness mode and adapter integration."*

---

## Alignment Check vs. CTO Success Criteria

| Success Criterion | Roadmap Status | Finding |
|---|---|---|
| Minimal integration contract | Unclear in Block 16 | MAJOR-1 |
| Telemetry adapter path for any harness | Block 17 addresses contract, but "any harness" not stated as architectural constraint | MAJOR-1 |
| Explicit `missing_telemetry`, not silent pass | Addresses in blocks but "degraded mode" language risks silent degradation | MAJOR-1 |
| Query/dashboard surface for task drift, gaps, scope creep | Block 19 and 20 address this | OK |
| Opaque health score | Explicitly excluded as non-goal | OK |
| Read-only / sidecar adoption | Block 14 introduces required-run contract before observation path exists | MAJOR-2 |

---

## No Critical or Major Findings If

- Block 16 acceptance criteria explicitly measures "days to first evidence from *any* harness without enrollment"
- Block 14 clearly separates observation path from enforcement path with different acceptance criteria
- Block 19 acceptance criteria acknowledges capture depth preconditions for forensic completeness

Otherwise: REVISE with clarified path-to-value for unmanaged/observation-first adoption.
