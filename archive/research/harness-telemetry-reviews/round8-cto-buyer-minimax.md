<think>The user wants me to review the agentic-sdlc-evidence-substrate-v4-brief.md through the CTO Buyer persona lens, using the section 15 template. I need to be concise and focus only on convergence - whether this brief can be used to start v0 implementation.

Let me analyze the brief through the CTO Buyer's lens:

CTO Buyer's key concerns:
1. No new harness requirement - want control layer over existing processes
2. Minimal integration contract
3. Telemetry adapter path for any harness
4. Explicit missing_telemetry, not silent pass
5. Query/dashboard surface for task drift, evidence gaps, scope creep

Let me evaluate the brief against these:

1. **No new harness** - The brief addresses this well in Section 3, emphasizing `sdp-trace wrap -- <existing-harness-cli>` as the preferred approach. This is strong.

2. **Minimal integration contract** - Section 3 and Section 6 describe the adapter registration, but is it minimal? The adapter capability contract looks fairly detailed. The Unix domain socket requirement ($SDP_TRACE_SOCKET) is a specific integration point.

3. **Telemetry adapter path** - Section 6 covers this well with capability declaration, allowed event types, etc. V0 Unix socket transport is defined.

4. **Explicit missing_telemetry** - Section 8 covers this comprehensively with MissingEvidenceTable and all observed states.

5. **Query/dashboard surface** - Section 12 lists required queries, but there's no dashboard surface mentioned. This could be a gap.

Let me check for critical blockers:

1. The brief mentions a dashboard/query surface but doesn't detail what it looks like or how it integrates with existing tools. This might be acceptable as v0 local query.

2. The Unix domain socket requirement for V0 adapters - this is a specific platform constraint. Is it a blocker? Probably not - it's just V0 scope.

3. No mention of how adapters work with existing CI systems beyond `ci-verify`. The adapter path seems focused on harness-side telemetry, not CI integration.
   
The brief doesn't clearly address how existing tools like CI/CD pipelines, code review systems, or project management dashboards would consume the evidence substrate. This is a significant gap. The brief also focuses heavily on technical requirements rather than demonstrating business value to a CTO - it emphasizes trust scopes and observability without clearly showing how these translate to ROI or risk mitigation. There's minimal discussion of operational impact or team adoption friction, which are critical for a CTO evaluating this purchase.

For v0, the query surface appears to be command-line only, which addresses basic query needs but doesn't provide the dashboard visibility a CTO would expect. The Unix socket requirement for adapters is platform-specific but manageable for v0. The remaining concerns—dashboard integration, operational overhead, and business value articulation—are real gaps that could be addressed as refinement work rather than blockers, assuming the brief is scoped for a CLI-focused initial release. The adapter contract and Unix socket transport appear workable for V0, though the command-line-only interface limits broader adoption. The absence of integration points with existing CTO tooling like dashboards, Slack, or email is a notable gap that could hinder adoption, but this could be addressed in future iterations. The brief maintains a clear scope that prioritizes technical soundness over immediate market appeal, which is appropriate for an initial implementation. This aligns with an ACCEPTABLE_WITH_GAPS verdict and supports moving forward with v0 development, particularly for local and CI-first scenarios.
</think>

# CTO Buyer Review: Agentic SDLC Evidence Substrate v4 Brief

**VERDICT: ACCEPTABLE_WITH_GAPS**

---

**Convergence Assessment:**
- Can this brief be used to start a v0 implementation? **yes**
- Are there any critical blockers before implementation? **no**
- If yes, list only true blockers that prevent starting v0: *none*

---

**Critical Blockers:** *none*

---

**Major Gaps:**

1. **No adoption surface for CTOs who don't live in CLI.** The brief describes `sdp-trace query`, `sdp-trace explain`, and `sdp-trace export`, but provides no integration path to existing CTO tooling (dashboards, Slack alerts, PR checks, compliance exports). A CTO cannot tell their teams to run CLI commands post-incident.

2. **Adapter deployment reality is underspecified.** Section 6 assumes harness adapters can be modified to emit Unix socket events. For "any harness," this requires an SDK or sidecar pattern. The brief states Unix socket is V0 transport but doesn't address how adapters get deployed alongside existing tools without modification.

3. **Business value framing is absent.** The brief defines trust scopes and event schemas but never explains why a CTO should care about `ci_witnessed` vs `local_observed` in terms of risk, compliance, or liability. This is product marketing, not technical spec—but a CTO buyer needs it to champion adoption internally.

---

**Accepted V0 Limitations:**

- CLI-only query/export is acceptable for v0; dashboard surface is future profile
- Agent-authored contract downgrade is correctly scoped as not gate-grade
- Local ephemeral keys are explicitly not security-grade
- Retroactive process attach is not supported (correctly scoped as V0 gap)

---

**Minimum Viable Correction:**

Add one paragraph to Section 2 (V0 Product Promise) or Section 12 (Query and Explain):

> *"V0 query output is CLI-native. Integration with dashboards, notification systems, or compliance exports is a post-V0 profile. V0 enables adoption by teams and organizations that require a machine-readable audit surface before committing to a full dashboard investment."*

This tells the CTO buyer that the brief anticipates their need without committing v0 scope to solve it.

---

**Questions Before Implementation:**

1. **Who owns the Unix socket server?** If the recorder creates `$SDP_TRACE_SOCKET`, does the adapter connect to the recorder, or does the recorder poll the adapter? Protocol direction matters for sidecar deployment.

2. **What is the minimum viable adapter sample?** For "any harness" adoption, the brief needs at minimum one reference adapter (e.g., OpenCode) showing the registration → capability → event emission sequence. Without this, "telemetry adapter path for any harness" is a promise, not a plan.

3. **Does `sdp-trace verify` produce a machine-readable gate input?** Section 11 shows human-readable output format but doesn't confirm JSON output for CI/gate consumption. A CTO integrating this into existing CI requires programmatic pass/fail, not text parsing.

---

**Demo Changes Required:**

| Demo | CTO Buyer Gap | Suggested Addition |
|------|---------------|-------------------|
| Demo 0 | No evidence of "what's visible without rewriting harness" | Show adapter-suppressed run with MissingEvidenceTable showing harness rows as `unsupported` |
| Demo 0.5 | No CI integration visibility | Add `sdp-trace export --format audit-bundle` pipeline showing compliance artifact generation |
| Demo 5 | Fabricated traces are a CTO nightmare | Add scenario: local trace replayed with mismatched source digest → `ci_witnessed` downgrade to `vcs_observed` |

---

**CTO-Specific Assessment:**

| Area | v4 Assessment | Verdict |
|------|--------------|---------|
| CTO usefulness | Correctly scopes `sdp-trace` as evidence substrate, not governance. Explicitly defers business pass/fail to `sdp-gate`. | **Adequate** — avoids the "replace your SDLC" trap |
| v0 capture boundary | `wrap` proxies existing commands. No harness rewrite required. | **Strong** — directly addresses CTO rejection criteria |
| Adapter interface | Unix socket transport defined; capability contract defined; registration lifecycle defined. | **Adequate** — missing deployment/sample path |
| Expected evidence contracts | Gate-grade requires VCS-pinned or human-signed contract. Agent-authored contracts cannot self-certify. | **Correct** — prevents false trust claims |
| Observer authority | Unauthorized signers remain `self_claimed`. Topology independence recorded. | **Adequate** |
| Signing and verification | Local ephemeral keys are explicitly not security-grade. CI witness boundary requires independent process. | **Transparent** — CTOs will appreciate explicit trust limits |
| Privacy and retention | No raw prompts/responses by default. Redaction pre-write. | **CTO requirement met** |
| Adoption and DX | `wrap` for existing harness; `ci-verify` for CI pipelines; offline-first. | **Adequate** — missing dashboard path |
| Forensic query/replay | Full query set defined; `explain` covers edge cases; integrity audit on chain corruption. | **Strong** |

---

**Summary:** The v4 brief is implementation-ready for its stated scope. The CTO buyer's primary gap is the absence of a consumption/integration surface beyond CLI. This is a v0 scope boundary issue, not a design flaw. The brief correctly avoids over-promising on dashboard, notification, or external integration paths. Recommend accepting with the minor clarification noted above and proceeding to implementation.
