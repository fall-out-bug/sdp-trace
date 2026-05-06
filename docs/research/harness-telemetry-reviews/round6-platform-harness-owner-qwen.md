```text
VERDICT: ACCEPTABLE_WITH_GAPS

Convergence assessment:
- Can this brief be used to start a v0 implementation? yes
- Remaining blockers before implementation: 3

Critical blockers:
- Adapter identity/auth enforcement undefined (adapter interface, observer authority). The brief states agent-accessible keys cannot sign `local_observed` or stronger, but does not specify how `harness_adapter` identity is bootstrapped or bound. Without this, any process can self-claim harness role. Tie: observer authority, signing and verification, adapter interface.
- No explicit fail-closed / degraded mode specification (v0 capture boundary). The brief mentions degradation to `cannot_verify` or `partial` but does not define the operational behavior: should the wrapper exit non-zero, allow continuation with degraded evidence, or block? A platform harness owner needs a single configuration flag that determines this. Tie: v0 capture boundary, adoption and DX.

Major gaps:
- No adapter lifecycle events (adapter interface). There is no `adapter_activated`, `adapter_error`, or `adapter_disconnect` event. If a harness adapter crashes mid-run, the verifier cannot distinguish silent absence from partial telemetry. Tie: adapter interface, evidence model.
- Wrapper bypass is detected but not operationally constrained (v0 capture boundary, adoption and DX). The brief correctly states "if work happens outside the wrapper, there may be no run," but provides no Git hook, CI pre-flight, or IDE integration path to make wrapper usage normative. Platform teams need a deployment strategy, not just a gap report. Tie: v0 capture boundary, CI/gate anchoring, adoption and DX.
- Contract locking enforcement mechanism is underspecified (expected evidence contracts). `lock_required_before: run_started` requires a control point that observes the lock before any child process starts. If the wrapper and agent share a process, who prevents the agent from modifying the contract post-attach? Tie: expected evidence contracts, v0 capture boundary.
- Adapter capability negotiation is absent (adapter interface). Different harnesses can observe different things. The brief lists a fixed taxonomy but does not specify how an adapter declares its capability set at registration time so the verifier can correctly distinguish "not observed because adapter cannot" from "not observed because telemetry failed." Tie: adapter interface, evidence model, forensic query/replay.
- Replay/forgery protection for local-only runs is underspecified (signing and verification, demo credibility). The brief mentions nonce mismatch and external anchors as detection signals but does not specify how the initial nonce, process binding, or file-system binding is generated before the wrapper launches. Without filesystem/process binding at attach time, post-hoc chain generation remains trivial. Tie: signing and verification, provenance model.

False assumptions:
- "Local writes are append-only and async" assumes controlled file system access (v0 capture boundary). If the harness or agent shares the same execution context as the recorder, append-only is not enforced at the OS level without a sandbox. The brief should clarify whether the wrapper uses a separate process, a sidecar, or assumes cooperative behavior. Tie: v0 capture boundary.
- "CI must sign verifier result" assumes CI has independent access to the chain head and source commit (CI/gate anchoring). In some runners, the agent produces both the artifact and the verification inputs. The brief addresses this ("downgrade any event whose signer is not independent") but does not specify the topology declaration mechanism. Tie: CI/gate anchoring, observer authority.

Minimum viable correction:
1. Add adapter registration protocol with capability declaration and identity binding.
2. Define explicit fail-closed vs degraded operational mode for the wrapper.
3. Specify bootstrap nonce and process/filesystem binding at attach time to prevent post-hoc forgery without external witness.
4. Add adapter lifecycle events: `adapter_activated`, `adapter_error`, `adapter_disconnect`.
5. Add deployment hook specification (Git hook, CI pre-flight, or equivalent) that makes wrapper bypass detectable before work begins.

Questions before implementation:
- Can the wrapper guarantee append-only writes without a sandbox, or does it assume cooperative harnesses?
- What identity document or attestation binds a `harness_adapter` to a specific harness version?
- Is there a minimum trust scope below which a run is not displayed in `sdp-trace query` at all?
- How does the wrapper handle concurrent runs in the same repository or workspace?

Demo changes required:
- Demo 0: Add a visible wrapper bypass detection path (e.g., run without wrapper, show forensic contrast). Tie: adoption and DX, demo credibility.
- Demo 1: Add adapter lifecycle failure injection (simulated adapter crash mid-run) to show how partial telemetry is reported. Tie: demo credibility, evidence model.
- Demo 4: Show process binding at attach time by demonstrating that a post-hoc chain with correct hashes but wrong process nonce still fails. Tie: signing and verification, demo credibility.
- Demo 8: Add a metric showing the cost of adapter lifecycle events relative to the 5ms p99 target. Tie: v0 capture boundary, demo credibility.
```
