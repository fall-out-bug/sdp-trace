# Persona 02 Review: Platform / Harness Owner

## VERDICT: REVISE

Two major findings require resolution before this roadmap can guide deployable work. Several minor findings are attached.

---

## Major Findings

### M1: Interception boundary is architecturally assumed but never committed

The persona's central question — *"Where exactly are events intercepted?"* — is unanswered at the architectural level. The roadmap references wrappers, adapters, process observation, and enrollment markers as deliverables inside Block 16, but never states the physical capture model.

Is the interception point:

- A **process wrapper** (`sdp-trace run -- opencode ...`) that observes child process lifecycle, argv, and exit code?
- A **shell hook or Git hook** that intercepts at the VCS boundary?
- A **filesystem watcher** that observes file mutations in the worktree?
- An **adapter shim** that an agent loads via plugin API?
- Some combination, with different trust grades per layer?

This is not a Block 16 implementation detail. It is an architectural decision that determines:

- what evidence Block 14's gate contract can realistically require;
- what Block 15's checkpoint signs over;
- whether Block 17's adapter events are primary capture or supplementary;
- what "fail-closed" physically means (process refused to start? gate rejects artifact? both?).

**Required:** State the interception architecture — layers, boundaries, and which layer owns which event type — before Block 14 work begins. This can be a short block-0 or a section in Block 13A, but it must precede gate-contract design.

### M2: Block 16 wrapper architecture should inform Block 14 gate contract, not follow it

The recommended order places Block 16 (managed harness enforcement) after Block 15 (signed checkpoint), which itself follows Block 14 (gate contract). But Block 16 is where the physical wrapper model, enrollment, adapter registration, and fail-closed behavior are first specified.

If Block 16 discovers that the wrapper can only reliably capture process lifecycle and argv — and that tool-call, file-mutation, and model-call events require agent cooperation — then Block 14's required-run manifest may have overclaimed what a fail-closed gate can enforce without an adapter.

Conversely, if Block 16 defines a robust enrollment and adapter registration model early, Block 14's gate contract can be designed to require adapter enrollment evidence in managed mode.

**Required:** Extract the wrapper and adapter architectural decisions from Block 16 into a prerequisite design checkpoint — either as part of Block 13A or as a dedicated Block 13B — so that Block 14's gate contract is designed against a known capture boundary rather than an assumed one.

---

## Minor Findings

### m1: Degraded-mode operational behavior is underspecified

The persona asked: *"What works for harnesses with no plugin API?"* Block 16's acceptance criteria state that unmanaged profile reports `missing_telemetry` or `not_integrated` without pretending capture is complete. This is a correct state label, but it is not an operational specification. A platform owner deploying to a team using an unmanaged harness needs to know:

- Does the wrapper still run and emit partial lifecycle events?
- Is the developer informed at runtime, or only at gate time?
- Does the gate treat `not_integrated` as a blocking state in managed mode?
- What does `doctor` report for this case?

Consider adding a small "operational behavior" section to Block 16's deliverables.

### m2: Prompt and model-response events are absent from the adapter contract

Block 17 lists `model_call_observed` but does not include `prompt_sent` or `response_received` as explicit event shapes. These are the highest-value capture points and also the most sensitive (which Block 18's redaction work acknowledges by sequencing after capture depth). Stating explicitly that prompt/response capture is deferred to the retention-profile block — or naming it as a `not_captured` state from day one — prevents false expectations.

### m3: Tool-level wrapping as an interception strategy is not considered

The architectural options discuss agent adapters and CI witness but omit the possibility of wrapping individual tools (`git`, `docker`, `kubectl`, test runners) at the PATH or shell-function level. For harnesses with no plugin API, tool-level wrapping is a pragmatic capture point that does not require agent cooperation. It should appear in the architectural options even if it is ultimately rejected, because it directly addresses the persona's question about harnesses with no plugin API.

### m4: `doctor`-style environment validation should start earlier

Block 16 includes a `doctor` or equivalent environment check. Platform owners will rely on environment validation from the first deployment. Consider making `doctor` a cross-cutting deliverable that first appears in Block 14, so that teams can diagnose gate-contract failures before the managed-harness profile exists.

### m5: No mention of wrapper overhead measurement timeline

The persona's success criteria include operational deployment details. Block 13A's acceptance criteria mention "measured wrapper overhead on real demo work" as a DX gap closure item, but no block explicitly owns overhead measurement. If Block 14 introduces a required wrapper path, overhead should be measured at that point, not deferred to a later DX block.

---

## Responses to Socratic Questions (from persona perspective)

**Where are the actual control points?**
Not yet committed. CI boundary (Block 14) is clear. Process-wrapper boundary is implied but architecturally unspecified. Adapter boundary is deferred to Block 16–17. This is the core of M1.

**What works for harnesses with no plugin API?**
The roadmap correctly identifies degraded mode and `not_integrated` states. The operational behavior of that mode — what still runs, what the developer sees, what the gate does — is not specified. See m1.

**What fails closed, and what only reports degraded posture?**
Managed profile (Block 16) fails closed. Unmanaged profile reports degraded. The boundary between them — enrollment, adapter registration — is a Block 16 deliverable but should be designed earlier. See M2.

**Can teams debug setup failures without reading raw JSON?**
`doctor` appears in Block 16. `gate explain` appears in Block 14. Both are correct answers, but `doctor` arriving only in Block 16 means Blocks 14–15 have no user-facing diagnostic for wrapper or enrollment failures. See m4.

---

## Summary

The roadmap is well-aligned with the persona's rejection criteria (no per-agent changes, no SDK-only model, no manual logging, no git reconstruction). The adapter event contract and managed/unmanaged harness profiles directly address the persona's success criteria. The two major findings are sequencing and commitment problems, not directional disagreements. Resolving M1 and M2 before Block 14 begins would make the roadmap deployable from this persona's perspective.
