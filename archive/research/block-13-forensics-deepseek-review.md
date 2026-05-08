# Persona 05 Verdict: REVISE

## Critical Findings

### C1. Test evidence vs. agent claims not distinguished
**Criterion:** *"Which tests were evidence, and which were only claims?"* / Rejects *"The agent said tests passed."*

Block 17 proposes a `test_observed` event but does not model provenance for test observations. A single event type collapses independently-observed test execution (e.g., CI-run tests, harness-witnessed test runs) and agent self-reported test claims into the same bucket. A forensic reviewer a month later cannot determine whether a "pass" came from the agent's own output parsing or from a verifiable test runner.

The rejection criteria explicitly forbid trusting agent-reported test results. The roadmap must require that every test observation carry a provenance field distinguishing at least `ci_executed`, `harness_observed`, `agent_reported`, and `cannot_verify`, and that query output surface this distinction.

---

## Major Findings

### M1. Witness-before-merge temporal binding not specified
**Criterion:** *"Witness before merge"* / *"Can we prove the run existed before merge?"*

Blocks 14–15 cover CI witness and signed checkpoints, but the roadmap never requires that the witness checkpoint be temporally bound to the merge event in a verifiable way. A forensic reviewer needs to answer: was this checkpoint generated *before* the PR merged, or could it have been backfilled? The deliverable list should include a checkpoint→merge binding (e.g., checkpoint hash embedded in merge commit message, or witness attested with a pre-merge CI run ID that the VCS host can confirm).

### M2. PR linkage absent from event model and query surface
**Criterion:** *"Link to commit, PR, and CI"*

Commit and CI are well-covered. PR (pull request / merge request) is mentioned in introductory text but never appears as a first-class concept in:
- The event schema (Block 17)
- The gate contract (Block 14)
- The forensics query pack (Block 19)
- The witness profile (Block 15)

Without first-class PR linkage, the forensic reviewer cannot traverse from an incident timeline to the PR discussion, approval events, review comments, or the merge decision context. PR identity must be a linkable artifact, not a documentation-only concept.

### M3. Approval actor identity and authority not modeled
**Criterion:** *"who approved it, what was checked, and where is the evidence"*

Block 14 covers `policy_override_requested` for sdp-trace-specific overrides, but the broader forensic question — "who approved this change?" — requires linking to the PR approval event (human reviewer identity, timestamp, approval scope). The roadmap treats approvals as an external VCS concept but never specifies how the forensic query surface surfaces or links to them. If PR linkage is added (see M2), the approval trail should be surfaced as part of the timeline or provenance query.

---

## Minor Findings

### m1. Retention enforcement mechanism undefined
Block 18 defines retention *profiles and manifests* but does not specify how retention is *enforced*. The persona needs to know: what prevents early deletion of a run directory? What happens when a retention manifest doesn't match the artifacts actually stored? The verifier behavior for retention violations (degraded state, `cannot_verify`, alert) should be explicit.

### m2. `not_assessed` visibility in query output not mandated
**Criterion:** *"not_assessed gaps visible, not hidden"*

Block 19 specifies `query gaps` and Block 20 mentions `not_assessed` counts in degradation export, but the acceptance criteria for the timeline query do not require `not_assessed` to appear as a visible, always-present dimension on every returned row. A forensic reviewer must be able to see at a glance which facts carry no assessment, not only query for gaps separately.

### m3. Replay / deterministic reconstruction not addressed
**Criterion:** Rejects *"Reports without replay"*

Block 19 provides descriptive timeline query, but the roadmap does not specify whether the recorded event stream is sufficient for deterministic reconstruction of agent actions. "Replay" in a forensic context means the ability to re-execute or step through the recorded sequence. If this is out of scope for the product, it should be explicitly listed as a non-goal rather than left ambiguous.

### m4. Task change attribution underspecified
**Criterion:** *"Where is the original task, and who changed it?"*

Block 17 includes `task_locked` and `task_superseded` events for task drift, but the event schema does not call out *actor attribution* — was the task changed by the agent, by a human, or by an external system? The forensic timeline needs this to distinguish legitimate human scope changes from agent scope creep.
