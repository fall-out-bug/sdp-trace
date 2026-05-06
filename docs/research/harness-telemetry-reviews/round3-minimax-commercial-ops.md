<think>The user wants me to do a Round 3 review of the harness-telemetry-trust-brief.md from a CFO + COO buyer perspective. I need to evaluate it purely on cost, privacy burden, rollout burden, and governance model credibility.

Let me analyze the document through this lens:

**Budget blockers:**
- 4 integration layers (harness plugin + tool/shell + file watcher + LLM gateway) all need to be built and maintained
- Key management infrastructure for signing (unspecified but clearly non-trivial)
- Remote witness service (new infrastructure, ongoing cost)
- CI signing recipe (CI infrastructure changes)


- Verifier state machine expansion (development cost)

**Privacy burden:**
- Gateway captures prompt/response hashes but that's only the starting point—full prompt logging creates enterprise liability
- LLM gateway proxy is itself a man-inhe-middle that may violate terms of service or data processing agreements
- No PII scrubbing defined
- File watcher captures path names (potential project structure disclosure)

**Rollout burden:**
- Requires adoption across multiple harnesses without a clear minimum viable path, creating a chicken-and-egg problem where teams need to adopt before seeing value
- No clear rollback strategy if adoption stalls

**Governance model:**
- No defined ownership model for who controls the system
- No clear escalation path for unapproved harness detection
- No policy engine exists yet
- The "approved harnesses, agents, prompts, or skills" registry is just a placeholder
- Audit trail for governance decisions isn't designed

**CTO value proposition:**
- 12+ queries sound useful but the display layer hasn't been built
- No clear answer on whether "team drifting into unapproved harnesses" actually gets answered
- No ROI model

Now I'm looking at the adoption phases to see if there's a realistic path forward, starting with the observe-only stage and its actual costs.

The minimum viable product isn't clearly defined, and there's no proof of concept for any single layer. The one-week value proof is just a list of 12 queries without explaining how to get from raw JSON to answers. I need to think through the rollout phases more systematically and figure out what can actually be demoed versus what requires infrastructure.
</think>

```text
VERDICT: CHANGES_REQUIRED

Budget blockers:
- B1: Four-layer capture architecture is not one budget item. It is four parallel infrastructure tracks: harness plugin API, tool/shell/file recorder, LLM gateway proxy with privacy controls, and remote witness protocol. Each has independent build, test, and maintenance cost. The brief presents these as a "layered" architecture. A COO sees four engineering quarters, minimum.
- B2: Key management is not scoped. Local recorder signing requires provisioning, rotation, and revocation on heterogeneous developer machines. CI signing requires a separate signing service (the agent cannot have access to the signing key). Remote witness requires a hosted or operated service. None of these are "add a library." Each is a service with its own SLA, security posture, and operational cost. The brief lists all three as missing (Section 9) and presents them as near-term architecture.
- B3: No staffing model exists for the sdp-trace platform itself. Who operates the remote witness? Who rotates keys? Who updates the harness registry? Who triages verifier alerts when the CTO dashboard shows a failed run? This is a platform function with headcount implications that the brief does not address.
- B4: Capture latency and cost are unquantified. A single agent session generates thousands of tool calls per hour. Four interception layers, hash chaining, and store writes per event. At a team of 50 developers running 8-hour sessions, the event volume is non-trivial. Storage, verification time, and network egress for remote witness anchoring are all unmodelled. The brief cannot produce a budget estimate without these numbers.

Privacy burden:
- P1: LLM gateway proxy is a man-in-the-middle. Before the first line of code is written: does every model provider's terms of service permit proxy interception? OpenAI, Anthropic, Google, and self-hosted models all have different API ToS. Enterprise customers have data processing agreements that may prohibit prompt/response interception even in hash form. The brief raises "prompt/response hash" as a feature and "privacy" as a risk in the same sentence without resolving the conflict. This is a legal review item before any engineering, not a follow-up.
- P2: File watcher captures path names. Project structure, file names, and path patterns are disclosed to the recorder. If the recorder is a third-party service or a vendor-operated remote witness, the path data is a significant IP disclosure vector. The brief never scopes where the recorder data flows and who operates it.
- P3: No PII scrubbing model. Agent sessions may include names, emails, internal system paths, and business context in prompts. The brief captures prompt hashes but does not define minimum entropy floors or scrubbing rules for short hashes (a concern raised in the red team review). Enterprise DPO will block this until a documented scrubbing pipeline exists.

Operational blockers:
- O1: The adoption ladder (Section 7) requires instrumenting "existing harnesses" before any trust value is delivered. But the brief has no minimum viable adapter. The first integration effort cannot be scoped because the harness adapter contract does not exist. This is a circular dependency: you need an adapter to demonstrate value, but you cannot scope the adapter without knowing what harness you target first. You need one concrete adapter target before any adoption claim.
- O2: No rollback procedure defined. If teams adopt observe-only and the telemetry shows nothing useful (because the tool/shell layers are blind in many harness contexts), what is the rollback path? The brief does not define an uninstall or detach procedure. A COO who deploys observe-only and sees no value will not trust the product for managed mode.
- O3: No alert response procedure. When the verifier emits `cannot_verify` or `missing_telemetry` for a critical gate, what happens? Who is paged? Is the merge blocked? By whom? At what rung of the adoption ladder does a verifier `fail` actually block a gate? The brief says "managed mode: verifier fails or blocks" but never says who implements the blocking mechanism. If it is not CI, it is not implemented.
- O4: No rollback or data retention policy for the telemetry store. Telemetry data includes: file paths, command history, model prompts (hash form), agent task descriptions, and event chains. GDPR, CCPA, and analogous regulations require data retention limits and deletion procedures. The brief has no data governance section. This is a legal blocker for EU and California customers before the first demo.

Governance requirements:
- G1: No ownership model. Who decides which harnesses are approved? Who approves the model list? Who can declare managed mode? Who sees the verifier output? The CTO dashboard implies a governance layer that has no defined owners or change process. This is not an engineering gap; it is a policy gap that engineering cannot fill.
- G2: The harness registry (Section 9, listed as missing) is a governance artifact, not a technical one. It requires: an approval process for adding new harnesses, a deprecation process for removing approved harnesses, a compliance check for unapproved harness detection, and a policy for what "unapproved harness detected" means operationally (alert? block? report?). None of this is scoped.
- G3: No audit access model. The verifier output and telemetry store contain sensitive operational data: what files were changed, what commands ran, what models were used. Who has read access? Is this SOC 2 scope-creep? The brief does not define an access control model, which means it will default to "whoever can access the repo" — which in many enterprises is not an acceptable audit store.

One-week value proof:
- The brief cannot produce a one-week value proof because no layer is implemented. The only credible one-week proof is: one harness adapter targeting one specific harness (pi or OpenCode), capturing task lifecycle events only, displaying the event log in a structured format, with a worked example showing how one CTO query maps to one event. Nothing else can be demoed in one week at credible quality.
- The brief's one-week claim would be 12+ questions answerable from raw JSON. That is not a product. That is a schema. The CFO needs to see: what does the structured answer look like, not just what questions can be asked.

Rollout phases (credible re-scope):
- Phase 1 (Month 1-2): One harness adapter for one target harness. Shell wrapper for commands only. Baseline verifier for this one layer. No gateway, no remote witness, no file watcher. Deliverable: structured event log for one harness session with worked CTO query examples.
- Phase 2 (Month 3-4): Add tool wrapper and file watcher correlation. Define completeness contract for this layer pair. Add partial verifier state model (verdict × scope × completeness from red team review). Deliverable: cross-layer correlation for tool calls and file mutations within one harness.
- Phase 3 (Month 5-6): Add LLM gateway digest capture (prompt/response hash only, no full prompt). Define gateway-to-local correlation protocol. Add CI signing recipe. Deliverable: two-layer provenance correlation with CI-signed chain.
- Phase 4 (Month 7+): Remote witness protocol, approved harness registry, CTO governance dashboard. Only funded if Phases 1-3 show adoption and measurable CTO value.

Governance requirements for Phase 1 funding approval:
- Define recorder data scope and retention policy before any developer-facing deployment.
- Define who operates the verifier output (internal vs vendor-hosted) before any cloud integration.
- Define the blocking mechanism for managed mode before claiming managed mode is available.
- Define the minimum viable CTO query surface (cut from 12 to 3) before building the dashboard.

Demo scope cuts:
- CUT: LLM gateway capture. Do not demo. No privacy model, no integration, no correlation protocol. Showing a prompt hash in a demo is meaningless without proving correlation.
- CUT: Remote witness. Do not demo. Protocol undefined. Showing a "witnessed" label on a demo run is a trust claim the product cannot support.
- CUT: Multi-harness support. Do not demo. No adapter contract exists for any harness yet. Showing "harness-neutral" with zero adapters is a lie.
- CUT: The full 12-query CTO dashboard. Do not demo. This does not exist and the data model to support it is not specified. Showing a dashboard with placeholder queries is worse than showing nothing.
- DEMO ONLY: One harness adapter (target: whichever harness has the most accessible plugin API — determine this first), shell event capture, structured event log output, and three concrete worked examples mapping event data to a specific CTO question. Show the verifier state for one captured run: pass/fail/cannot_verify/not_assessed with plain-language explanation of what each means for the gate decision.

Attack scenarios from governance/operational failure:
- G-A1: Vendor lock-in via recorder data. If sdp-trace is vendor-operated (remote witness, hosted verifier), the telemetry store becomes a switching cost. Enterprise procurement will flag this. The brief never addresses self-hosted vs vendor-hosted, which means the decision is made by default (whoever ships the first implementation wins).
- G-A2: Verifier alert fatigue. If every developer's debug run emits `cannot_verify` because the tool wrapper is not attached, the verifier signal becomes noise. Teams will ignore it or disable the wrapper. This is the operational death spiral for observe-only mode. The brief assumes adoption will progress through the ladder; it does not model abandonment risk.
- G-A3: CTO dashboard as a false trust signal. If the dashboard shows `pass` for every run because the display layer is not enforcing scope labels, CTO sees a healthy delivery process. The trust degradation the brief promises to detect is invisible because the display layer normalizes away the trust-scope distinction. This is the most commercially dangerous scenario: the product provides false assurance that masks the degradation it was designed to detect.
```
The most commercially dangerous scenario is that the dashboard displays `pass` for every run while masking the trust-scope distinction in the display layer, creating false assurance that obscures the very degradation it was meant to surface. This operational failure mode—where the product promises to detect trust degradation but instead normalizes it away—is the one most likely to damage both the product and the teams relying on it.
```

</file>
