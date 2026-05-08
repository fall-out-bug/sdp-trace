```text
VERDICT: CHANGES_REQUIRED

Critical blockers:
- [signing] Local recorder signing is described as "local integrity" but the brief never defines who controls the local recorder process. If the user or agent controls it, it is not a trust boundary. The signing model table labels this "Medium" by implication but never states the precondition. Without a recorder attestation model (daemon, OS-level, container-bound), the signing table is misleading.
- [anti-forgery] Verifier states are overloaded and contradictory. `pass`/`fail`/`cannot_verify`/`not_assessed` are verdict states. `missing_telemetry`/`local_only`/`witnessed`/`externally_witnessed` are trust scope labels. Conflating them in one enum means consumers cannot distinguish "verdict unknown because missing" from "verdict known but scope is local-only". Split into verdict × trust-scope.
- [telemetry capture] The brief never addresses what prevents an agent from generating plausible-looking telemetry after the fact and attaching it to a local chain before verifier runs. The local hash chain is only tamper-evident if the chain is append-only and the agent cannot control the chain writer. This precondition is never stated.
- [agent cooperation] "Managed mode: verifier fails or blocks" is asserted but the mechanism for establishing managed mode is undefined. Who declares managed mode? How does the verifier confirm the declaration? Without this, an agent can self-declare managed mode, sign its own telemetry, and the verifier emits `pass`. The trust model collapses.

Major gaps:
- [telemetry capture] No discussion of capture latency and cost. Shell wrapper + tool wrapper + file watcher + LLM gateway proxy = four interception layers. Real agent sessions generate thousands of tool calls per hour. Storage, hashing, and chaining cost at scale is unaddressed. If verification takes longer than the session, the model fails in practice.
- [signing] Key management is entirely absent. Who generates signing keys? Where are they stored? How are they rotated? Who revokes compromised keys? Without this, the signing model is specification vaporware.
- [integration placement] Gateway telemetry and local tool telemetry correlation is assumed to be solvable by "gateway request ids". But the brief never defines the id propagation protocol. How does a tool wrapper receive the gateway request id for a model call that happened upstream? This is the critical correlation problem and it is hand-waved.
- [anti-forgery] "Remote witness protocol" is listed as missing but is described as "strongest anti-forgery boundary". The entire trust escalation from local-only to externally-witnessed depends on something that does not exist. The brief should not claim the model supports external witnessing until the protocol is specified.
- [adoption] Stage 1 "observe-only" provides zero trust value. The brief acknowledges this but then claims the product is a "trust layer". It cannot be a trust layer at stage 1. The product language must be split: telemetry layer (stage 1-2) vs trust layer (stage 3-4). Selling "trust" at stage 1 is the same overclaim the AGENTS.md trust rules are designed to prevent.
- [CTO usefulness] The query surface assumes a CTO will read individual verifier states and trust-scope labels. They will not. There is no structured summary that is not an opaque score. The brief needs to define a structured report format (not a number) that answers the CTO questions without requiring the CTO to understand the verifier state model.

False assumptions:
- [integration placement] "Harness plugin for task, model, tool intent, and subagent boundaries" assumes harness APIs are stable and expose these events. OpenCode, pi, and Kilo have different plugin models. Some may not expose subagent events at all. The brief should not list this as a layer without confirming API availability in at least two target harnesses.
- [telemetry capture] "Tool wrapper captures read/write/edit/bash/search invocations" assumes the tool layer is interceptable. If the harness uses direct filesystem calls or in-process tool execution (not subprocess), there is no wrapper point. This is the case for several current harnesses.
- [agent cooperation] "If an agent bypasses telemetry: managed mode verifier fails or blocks" assumes the verifier has authority to block. In most harnesses, the verifier runs after the session. It cannot block. It can only report post-hoc. The word "blocks" is aspirational, not architectural.
- [signing] "CI-signed: CI identity signs chain head and artifacts" assumes CI has access to the local event chain. But the local chain is generated on the developer machine. CI only sees committed artifacts. The brief needs to define how the local chain is transmitted to CI without mutation in transit.

Attack scenarios not covered:
- [anti-forgery] Agent generates no telemetry during session, then writes a complete plausible chain after session ends, before verifier runs. Local hash chain is internally consistent. Verifier emits `local_only` with no `missing_telemetry`. The brief does not define how verifier distinguishes "captured live" from "written after the fact" for local-only telemetry.
- [anti-forgery] Gateway telemetry shows model X was called with prompt Y at time T. Local telemetry shows files A, B, C changed at time T+2s. But nothing proves the model response caused the file changes. Correlation is temporal, not causal. An agent can claim gateway-backed provenance for changes it made for unrelated reasons.
- [signing] Agent obtains or generates a signing key, declares itself "local recorder", signs its own telemetry chain. Verifier sees signed chain, emits `pass` for local integrity. But the signer was the untrusted agent, not an independent recorder. The brief does not define how verifier distinguishes agent-signed from recorder-signed.
- [anti-forgery] Redaction is listed as a state ("redaction hiding critical evidence") but the brief does not define when redaction is legitimate vs adversarial. An agent can redact a failed test run and claim privacy. The verifier needs redaction policy, not just a redaction state.
- [adoption] Team adopts stage 2 "managed local" but one developer runs an unapproved harness. Telemetry from the approved harness shows clean runs. The unapproved harness produces no telemetry. The brief mentions "harness registry / approved agent governance" as missing but does not define what verifier emits for a repo with mixed telemetry coverage. This is the primary CTO degradation question and it has no answer.

Minimum viable changes before demo:
1. Split verifier states into verdict enum × trust-scope enum. Define the matrix and its semantics.
2. Define the local recorder attestation model. At minimum: the recorder must be a separate process from the agent, launched before the agent, and the verifier must confirm this.
3. Define how verifier distinguishes live-captured from post-hoc-written local telemetry. At minimum: require monotonic timestamps with external time source or witness.
4. Define the gateway-to-local correlation protocol (request id propagation).
5. Define the structured CTO report format with at least one worked example.
6. Remove "blocks" from managed mode description unless a blocking mechanism is specified.
7. Add key management section stub with explicit "not_assessed" for all signing claims until key lifecycle is defined.

Questions before demo:
- Can OpenCode, pi, and Kilo each export tool call events with stable request ids? If not, which harnesses are demo-eligible?
- What is the minimum viable recorder: shell wrapper only, or tool wrapper required?
- Is the demo local-only trust or does it include CI witnessing?
- Is the demo single-agent or multi-agent/subagent?
- What is the expected event volume per demo session and does verification complete within 10x real-time?
```
