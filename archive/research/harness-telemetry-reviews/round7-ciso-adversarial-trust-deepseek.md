VERDICT: ACCEPTABLE_WITH_GAPS

Convergence assessment:
- Can this brief be used to start a v0 implementation? **yes**
- Are there any critical blockers before implementation? **no** (none that prevent writing v0 code, provided the trust model is documented as incomplete by design for `local_observed`)
- If yes, list only true blockers that prevent starting v0. *N/A*

Critical blockers:
- None. The design acknowledges that local-only traces are not trustworthy (`local_observed` does not upgrade without external witness). The per-run ephemeral key destruction after run closure gives the trace **post-run tamper resistance** (you cannot insert signed events later). The verifier’s `expected_run_absent` and CI-witness signing provide a credible escalation path for gate-grade trust. These are sufficient to begin building the recorder, even if some security invariants need further hardening.

Major gaps:
1. **CI witness signing key boundary not defined.** If the CI job’s signing key is available to the agent (same process, shared memory, same token scope), the witness is not functionally independent. The brief says “Events whose signer is not independent … are downgraded” but never defines the independence threshold (process boundary? network-separated service? hardware module?). Without a bright line, implementers may ship a “witness” that is trivially spoofable.
2. **Contract tamper surface in VCS-backed mode.** A malicious agent can modify the contract inside the source tree and commit it before the run. The brief lists “contract committed in VCS before the source commit being verified” as valid provenance, but a commit that also contains the contract change can be the same commit under review. No mechanism (e.g., contract digest pinned in CI config) is required for VCS-sourced contracts; the verifier would accept the attacker-supplied relaxed contract unless an out-of-band digest rule is enforced. The specification should make pinned-digest enforcement mandatory for `ci_witnessed` when contract source is VCS.
3. **Honest‑absence vs. deliberate deletion on local disk.** The system emits `expected_run_absent` only when external observers (VCS, CI) detect a change that should have produced a trace. A developer (or compromised agent) who deletes the local run directory before any external observer notices will leave no `expected_run_absent` and no trace. The forensic story is: “if nobody looked, there is no evidence of absence.” That is an inherent design choice, but the brief does not state this boundary explicitly in the privacy/forensic section. V0 docs must be explicit that local-only absence is unprovable and that deletion before external noticing is indistinguishable from an honest no‑run unless a pre‑commit hook or mandatory telemetry upload is deployed.
4. **Replay of old valid local trace with a new source commit.** The verifier can detect VCS mismatch and timestamp mismatch, but an attacker who replays a previous trace against a different source tree and adjusts the timestamps could still produce a chain that passes structural verification, provided the private key is regenerated. The nonce and run-id would differ, but there is no binding between the run and the specific source commit hash in the `recorder_attached` event beyond a workspace digest. The attack window closes only if the verifier enforces that the workspace digest matches the CI‑observed commit digest. That linkage is implied but not required as a verifier check in the spec. (Demo 5 mentions VCS mismatch, but the test should demonstrate that a crafted trace with a fake workspace digest that matches the current commit would be accepted unless the hash includes ext‑workspace metadata.)

Accepted V0 limitations:
- Local traces are `local_observed` and cannot be trusted for gate decisions; this is by design.
- The ephemeral private key lives in process memory only; host compromise (memory dump, agent attaching to the process) defeats it. Acceptable for an honest local developer scenario.
- No transparency log, no append‑only service; external witness is CI signing only.
- Tool‑call capture depends on adapter cooperation; a harness that simply refuses to emit telemetry yields missing‑telemetry rows, not false trust.

Minimum viable correction:
- Add a **mandatory contract digest pin** rule for `ci_witnessed` whenever the contract source is VCS. If the contract digest is not pinned in CI config (or human‑signed), downgrade the maximum trust scope to `harness_observed` even if all events are present. This eliminates the “attacker commits a relaxed contract” path without blocking implementation.
- Define the **CI witness signing key boundary**: the signing key must reside in a process or service **not** fork/spawn of the agent process, and the verifier must record a signed assertion of independence (e.g., distinct CI job id, separate runner, attestation). Provide a clear pass/fail signal in the verifier output.

Questions before implementation:
- Will the CI witness signing key be a per‑job ephemeral key or a long‑lived pipeline key? If per‑job, how is the public key distributed to downstream verifiers?
- How are clock state / nonce values protected from an attacker who can set the system clock before starting the recorder? (Not critical for v0, but affects timestamp‑based replay detection.)
- Does the recorder have a mechanism to detect if its own binary or config has been modified by the agent? (Out of scope for v0, but a question that will arise during integration.)

Demo changes required:
- Demo 5 (post‑hoc fabrication) must explicitly include the case where an attacker regenerates a trace with a **workspace digest matching the current commit** and no VCS mismatch. Show that the verifier can still downgrade because no external witness signed the combination of run-id and commit-digest (if the verifier enforces that linkage).
- Demo 2 (CI‑witnessed gate) should show a deliberate attempt to use an agent‑authored contract relaxed to skip `test_observed`. Demonstrate that the verifier rejects the upgrade because contract provenance does not match a pinned digest in CI config. This makes the contract‑tamper threat concrete for reviewers.
- Add a demo (or explicit note) for the **CI witness key boundary**: show that when the signing key is in the same process as the agent, the verifier records `dependency_same_job` and downgrades to `local_observed` or `harness_observed`. This prevents false trust from co‑located “witnesses.”
