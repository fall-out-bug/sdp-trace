VERDICT: REVISE

Critical: (none)

Major:

1. **Gate enforcement deployed before trust anchor (Blocks 14 before 15) creates a forgery window.**  
   Block 14 proposes a “protected gate” and fail-closed CI profile before Block 15 delivers signed checkpoints or signer isolation. During that window, the gate consumes local trace JSON, override events, and CI witness metadata that a compromised agent or process can fabricate or rewrite inside the workspace. A policy gate making allow/deny decisions on forgeable inputs is a trust overclaim. The roadmap should either merge minimal DSSE/in-toto signing and authority binding into Block 14, or explicitly label any Block-14 gate behavior as advisory/`cannot_verify` until a distinct observer boundary (signed checkpoint) is operational.

2. **Block 15 omits explicit signer-process isolation from the agent workspace.**  
   The success criteria require “signing outside the agent-controlled workspace” and “key separation: agent cannot access the signing key.” Block 15 describes a “CI OIDC signing profile” but does not mandate that the signing step run in an isolated CI job, hardened step, or service separate from the agent’s filesystem and environment. If the observed agent runs in the same CI workflow with access to OIDC tokens or shared artifacts, it can sign a lie. The spec must require physical or job-level isolation between the agent process and the signing material.

Minor:

3. **Monotonic sequence requirement is missing from Block 15 acceptance criteria.**  
   The persona success criteria require a “monotonic sequence, nonce, and run id.” Block 15 mentions nonce and canonical trace head, but the acceptance criteria do not explicitly require monotonic sequence enforcement or verifier-visible gap detection in the checkpoint chain. Add monotonic index and gap detection to the Block 15 acceptance criteria.

4. **Deletion-and-restart detection relies only on absence, not pre-run intent witness.**  
   The persona asks: “Can the whole run be deleted and restarted?” Block 14 uses a required-run manifest to detect missing telemetry, but the roadmap does not state that the expectation-to-run must originate outside the agent’s control (e.g., from the CI scheduler or trigger event). Without an external intent witness, an attacker who deletes the run directory also deletes the manifest expectation if it lives in the same workspace. Clarify that the required-run manifest must be anchored by the CI trigger or another external observer.

5. **Path from `ci_witnessed` to `externally_witnessed` is under-specified.**  
   Block 15 and Block 21 mention optional external witness profiles, but the roadmap does not map the concrete mechanism (e.g., Sigstore/Rekor, customer transparency log, timestamp authority) that elevates a CI-signed checkpoint to `externally_witnessed`. The verifier states should explicitly state that Block 15 alone caps at `ci_witnessed`; `externally_witnessed` requires an additional append-only or transparency step defined in Block 21 or an equivalent profile.
