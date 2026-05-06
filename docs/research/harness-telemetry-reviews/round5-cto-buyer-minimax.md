<think>The user wants me to act as the CTO buyer persona and run an adversarial Socratic product review of the agentic SDLC evidence substrate brief, using the exact response template from section 12.

Let me review this brief critically from the CTO buyer's perspective, focusing on whether this product would actually be adopted and useful in practice.

Key areas to evaluate:
1. Product layering
2. CTO usefulness
3. Evidence model
4. Provenance model
5. Trace model
6. Observer authority
7. Signing and verification
8. Expected evidence contracts
9. CI/gate anchoring
10. Privacy and retention
11. Adoption and DX


12. Demo credibility

Now I'm looking at the demo plan to assess whether these scenarios will actually demonstrate credibility to a CTO buyer—whether they show real value without requiring deep technical buy-in, and whether they address the practical concerns about adoption and integration that a CTO would care about.</think>

# Adversarial Socratic Review — CTO Buyer Persona

## Critical Blockers

- **Adoption friction (Adoption/DX):** The brief defers multi-harness adapters but treats this as a technical deferral, not a product risk. If the first build slice requires teams to instrument their harness specifically for `sdp-trace`, adoption dies before the CTO report exists. The brief must define a *zero-instrumentation* capture path (process-level observation, shell wrapper, or passive log ingestion) as first-class, not a deferred feature.

- **CTO usefulness (CTO usefulness):** The report shape in Section 10 is a checklist of questions the product *could* answer, not evidence that it *will* answer them for teams who already have no structured harness. The brief does not specify what a CTO sees on day one with zero integration. The demo strategy in Section 8 does not include a "zero-instrumentation day-one" scenario, only scenarios that assume trace infrastructure already exists.

- **Evidence model (Privacy and retention):** The brief mentions raw prompts, responses, argv, stdout, stderr, and file paths as evidence types that "may leak confidential data" and defers this. This is a day-one problem, not a deferrable one. Any CTO with IP-sensitive code (models, internal tooling, proprietary specs) will refuse the product unless retention/redaction is defined before the first demo. The brief must state: what does `sdp-trace` capture by default, what is redacted by default, and what is the opt-in contract?

## Major Gaps

- **CI/gate anchoring:** The brief treats CI as a witness and a signer but never addresses the common case where CI *is also* the agent runtime (e.g., GitHub Actions runs a coding agent). In that topology, the CI witness and the agent-reported events originate from the same execution context. The trust boundary collapses. The brief must define how `ci_witnessed` does not reduce to `agent_reported` in co-located CI agent scenarios.

- **Observer authority (Signing/verification):** The authority rule ("a signer can only support trust for event types it is authorized to observe") is correct but underspecified. Who issues the authorization policy? Who audits the policy-to-signer mapping? If this is implicit in the implementation, the verifier is trusting the tool authors, not the signers. This is a critical gap for CISO evaluation.

- **Adoption/DX:** The persona file says the CTO's rejection criteria include "Agents should follow the envelope." The brief does not address what happens when agents *ignore* `sdp-trace` observation entirely. What does the trace look like when the agent bypasses the harness, wraps itself in a subprocess, or uses an approved tool through an unobserved path? The brief addresses missing telemetry as a signal but does not address *observed telemetry that is known to be incomplete by design*.

- **CTO usefulness:** The degradation indicators in Section 10 (share of PRs with missing telemetry, share of local-only test claims, etc.) are report-layer metrics, but the brief does not specify who owns `sdp-report` or how it connects to `sdp-trace` output. A CTO reading this brief cannot determine whether they buy one product or three, or whether the reports work with the first build slice.

## False Assumptions

- **Evidence model / Signing/verification:** The brief assumes signing provides meaningful trust when the signer identity is derived from a local recorder the developer controls. A local recorder signing a checkpoint is not meaningfully stronger than the agent self-reporting the same event. The brief acknowledges this in Demo 1 but treats it as a "not gate-grade" label rather than a fundamental question: is signed local-only telemetry *ever* useful, and if so, at what trust scope?

- **Product layering:** The brief correctly separates `sdp-gate` from `sdp-trace` but assumes the gate consumer exists and has a policy. In the first build slice, no gate policy exists. The brief should address what `sdp-trace` *does* without a gate — specifically what the verifier output means and who reads it before `sdp-gate` exists.

- **Demo credibility:** Demos 3, 4, and 5 are adversarial cases presented as expected outputs, not as things the product actually catches. Demo 4 (tamper attack) assumes the verifier is running before the tampered trace is presented. Demo 5 (post-hoc fabrication) assumes the system knows when the trace was created relative to the work. Neither is guaranteed without a timestamp authority or external witness. The demos should clarify when these detections are *structurally guaranteed* versus when they rely on freshness checks that are themselves forgeable.

- **CTO usefulness:** The brief assumes the CTO has a governance problem they want to solve. A different plausible buyer position: the CTO wants *adoption of AI-assisted SDLC* and is being asked to buy a control layer before the teams have adopted. `sdp-trace` as framed could be a reason for teams to avoid AI-assisted SDLC rather than a reason to trust it. The brief does not address this adoption-inhibiting framing risk.

## Standards We Should Reuse Instead of Building

- **Observer authority:** The in-toto framework already defines predicate types, attestations, link metadata, and inspection rules for supply chain observers. The brief should reference in-toto layout/inspection/verify as the authority policy model rather than defining a bespoke signer-authorization model. Sigstore and SLSA are correctly identified; in-toto is missing from the reuse list.

- **Retention/privacy (Privacy and retention):** The brief mentions data leakage risk and defers it. Industry standards for SDLC telemetry retention — what to retain, what to redact, how long, and who controls — are not referenced. If an industry standard (e.g., NIST SP 800-53 audit log controls, SOC 2 retention requirements) exists for this class of telemetry, the brief should adopt it rather than deferring a privacy model.

## Minimum Viable Product Correction

The first build slice should be:

1. A passive, zero-instrumentation shell/HTTP/git log observer that captures events from any harness or agent without requiring harness modification.
2. A retention/redaction contract that is explicit and configurable before the first demo, covering prompts, file paths, and test output.
3. A single CTO report for one run showing: observed evidence, trust scope, completeness, missing telemetry, and gate usability — with zero `sdp-gate` dependency.
4. The verifier three-axis output from Demo 1 and Demo 3 as the only trust model for the first build.

Defer the signing chain, checkpoints, and Sigstore integration until the observer capture path and retention model are validated. Signing an incomplete or under-specified evidence substrate makes verification credibility worse, not better.

## Questions Before Implementation

1. What is the zero-instrumentation capture path for the first demo? If it does not exist, what do we show the CTO on day one?
2. Who owns the retention/redaction policy, and is it configurable per-organization before the product is deployed?
3. In the CI co-located agent scenario (CI is both the harness and the witness), what does `ci_witnessed` actually prove that `agent_reported` does not?
4. What does the verifier do when it has a structurally valid local trace with no gate policy? Does it output the three-axis result regardless?
5. Who is the first buyer? A company adopting AI-assisted SDLC, or a company trying to audit an AI-assisted SDLC they already have? The brief tries to serve both but the adoption friction is different.

## Attack Scenarios Not Covered

- **Colluding agent + CI:** The agent generates a trace in CI that looks like `ci_witnessed` but the CI step that runs the verifier has a bypass, a skipped step, or a pre-populated artifact. The brief covers tamper attacks against the trace but not tamper attacks against the verification step itself.

- **Silent evidence suppression:** The agent or harness selectively omits events (e.g., skips logging a tool call that modified a config file) without modifying an existing trace. The missing-telemetry signal requires the *expected contract* to be complete enough to enumerate every event type. The brief does not address how the expected contract is kept complete enough to catch selective suppression.

- **Harness identity spoofing:** A modified harness reports itself as the approved harness. The observer model includes harness identity but does not specify how the harness identity itself is attested, not just signed. Without attested harness identity, a malicious harness can claim to be the approved one.

- **Retroactive contract weakening:** After a run completes, the expected evidence contract is modified to match what was actually observed, making the run appear complete. The brief does not address contract immutability or version-locking.

## Demo Changes Required

- **Add Demo 0:** Zero-instrumentation capture. Run any agent command, capture shell/git events passively, produce a local trace with `local_observed + partial`. No harness modification. This is the adoption demo, and it is missing.

- **Add Demo 6:** CI co-located agent. Run a coding agent inside CI. Show the verifier output and explicitly label whether `ci_witnessed` or `agent_reported` is the honest trust scope. Do not paper over the topology ambiguity.

- **Add Demo 7:** Selective suppression. Run an agent with a known suppressed event (e.g., a modified file not logged). Show whether the missing-telemetry signal fires correctly given the expected contract.

- **Clarify Demo 4:** State explicitly whether tamper detection is guaranteed by hash chain integrity or by checkpoint freshness. If an attacker modifies an event and re-signs the checkpoint, the hash chain is consistent and the tamper is invisible without a prior external anchor. This distinction matters for CISO evaluation.

---

```
VERDICT: CHANGES_REQUIRED

Critical blockers:
- No zero-instrumentation capture path defined; adoption dies on integration requirement
- Retention/redaction model deferred; day-one blocker for IP-sensitive buyers
- Day-one CTO value not specified; brief describes a future state, not a first-demo state

Major gaps:
- CI co-located agent topology collapses ci_witnessed into agent_reported; not addressed
- Signer authorization policy owner and audit path unspecified
- Degradation report connection to sdp-trace output and product boundary undefined
- Agent bypass/ignore paths not addressed beyond missing-telemetry model

False assumptions:
- Local signing meaningfully strengthens agent-authored telemetry
- Gate consumer policy exists before first build slice
- Tamper detection is structurally guaranteed vs. freshness-dependent
- CTO buyer wants control over AI SDLC vs. adoption enablement

Standards we should reuse instead of building:
- in-toto layout/inspection/verify for authority policy model (missing from reuse list)
- NIST SP 800-53 or SOC 2 retention controls for privacy/redaction baseline

Minimum viable product correction:
- Define zero-instrumentation passive capture path as first-class, before signing chain
- Define retention/redaction contract before first demo
- Deliver one CTO report for one run, zero gate dependency
- Defer Sigstore/SLSA until observer capture path and retention model validated

Questions before implementation:
- What is the zero-instrumentation day-one demo path?
- Who owns retention/redaction policy and is it configurable per-org?
- What does ci_witnessed prove when CI is co-located with the agent?
- What does verifier output mean without a gate policy present?
- Is the first buyer adopting AI SDLC or auditing existing AI SDLC?

Attack scenarios not covered:
- CI verifier bypass or skipped step (tamper against verification step, not trace)
- Silent event suppression (requires contract completeness, not just contract existence)
- Harness identity spoofing (attested harness identity not specified)
- Retroactive contract weakening (contract immutability not addressed)

Demo changes required:
- Add Demo 0: zero-instrumentation passive capture
- Add Demo 6: CI co-located agent topology ambiguity
- Add Demo 7: selective event suppression detection
- Clarify Demo 4: hash-chain integrity vs. checkpoint freshness guarantees
```
