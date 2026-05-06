# VERDICT: REVISE

## Critical Findings

**C1. No offline development path defined.**
The entire roadmap assumes network connectivity for CI witnesses, OIDC signing, external artifact refs, and checkpoint verification. Block 15 discusses DSSE/in-toto signing, Block 21 mentions air-gapped *documentation* but not a real air-gapped workflow. A developer on a plane, in a SOC, or during a network outage has no documented local-only mode that degrades gracefully. The verifier must accept `offline_dev` as an explicit state, not treat missing CI witness as `missing_telemetry`.

**C2. No measured latency budget or overhead SLO.**
Block 6 mentions "measured wrapper overhead on real demo work" but the roadmap contains no latency targets. Pre-write redaction, signing checkpoints, artifact digest computation, and adapter event routing *will* add time. Without explicit budgets (e.g., "wrapper adds <200ms per subprocess, <2s per run"), there is no way to evaluate whether the product slows teams down. A staff engineer will not ship a control layer without knowing its cost per command.

**C3. Emergency fix flow is underspecified.**
Block 14 defines `policy_override_requested` as a schema, but the operational path is thin. In an emergency, a developer needs: (a) a one-command or one-CLI-flag bypass, (b) no ceremony (no multi-approval flow in the tool itself), (c) deterministic recording that the *override happened*, not that everything was somehow fine. The roadmap says overrides "never upgrade `audit_grade_gate`" which is correct, but does not guarantee the override path is *fast* or *obvious*. If the tool blocks a hotfix for 30 seconds of checkpointing, it will be bypassed entirely outside the wrapper.

## Major Findings

**M1. `gate explain` exists but deterministic verifier output is not guaranteed.**
Block 14 includes `gate explain` output and Block 6 mentions it, but there is no commitment that running the verifier on identical inputs produces identical output. If gate explanation depends on wall-clock time, nondeterministic ordering, or external services, a developer cannot reproduce a failure locally. Add a determinism requirement: given the same event set and profile, `sdp-trace verify` must produce byte-identical results regardless of execution time.

**M2. No local preview of captured data.**
Success criteria require "local preview of what will be captured or sent." The roadmap has redaction profiles and retention manifests (Block 18) but no `sdp-trace preview` or `--dry-run` that shows a developer what *would* be recorded before actually recording it. Without this, developers cannot audit redaction behavior themselves and must trust the tool — precisely the wrong posture for a skeptical engineering org.

**M3. Signing and checkpoint ceremony risk outweighs capture value early.**
Block 15 introduces DSSE envelopes, authority policies, chain-head verifiers, and negative fixtures for replay/wrong-signer/stale-chain. This is necessary but heavy. A staff engineer's first interaction should not require understanding in-toto statement formats. If Block 15 ships before developers have a simple `sdp-trace run` that "just works" in local mode, adoption friction will push teams to bypass the wrapper entirely.

**M4. Block 16's `doctor`/environment check arrives too late.**
If a wrapper or adapter is misconfigured, the developer experience is a verifier failure with no local guidance. Block 16 places `doctor` at that point, but the roadmap's recommended order puts Block 16 *after* signing (Block 15). Developers need the health-check tool *before* the enforcement profile tightens the gate, or they will get opaque failures.

## Minor Findings

**m1. Demo path is eight blocks deep before showing full value.** A staff engineer evaluating this product wants to experience the control loop by end of Block 14 at latest, not after eight sequential blocks.

**m2. "Unsupported" states are listed in Block 17 but not uniformly typed across earlier blocks.** The gap between `unsupported`, `not_integrated`, `suppressed`, and `missing_telemetry` must be machine-enumerable from Block 14 onward, not introduced at Block 17.

**m3. No explicit developer-facing error taxonomy.** The roadmap describes verifier states (`local_observed`, `ci_witnessed`, `cannot_verify`) but not the CLI output format a developer sees when a gate fails. These must be distinct: a verifier state is not a user message.
