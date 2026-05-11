Status: NOT READY — critical trust-gap findings must be resolved before spec can be used as the basis for a CTO-visible demo proof.

---

## Critical Findings

### C1 – Observation log integrity is unverifiable (trust rule: “Machine proof wins over prose”)
The demo claims that `sdp-trace` observed the agent route and produced a tool‑generated packet that separates “observed facts.”  
The packet row `PC-AGENT-ROUTE` points to a recorder observation log that, under the current spec, is an ordinary file committed to the demo repository with no integrity protection (hash chain, signature, or time‑stamp). A reviewer—or the CTO—cannot distinguish a genuine recorder output from one fabricated or edited by Codex after the agent session.  
This violates the trust rule that machine proof wins over prose. The packet’s evidence ref is effectively prose that Codex can rewrite at will.

**Required fix**
Either:
- require the recorder to emit a signed, hashed, or externally time‑stamped observation artifact, **or**
- explicitly downgrade the demo claim: “The demo shows packet *format* and tool pipeline, not the trustworthiness of raw observation data. Recorder integrity is not machine‑verified in this slice.”

### C2 – Passivity of the recorder is asserted, not verified
The spec states `sdp-trace` is a “passive flight recorder” and the developer prompt MUST NOT mention `sdp-trace`.  
There is no mechanism described that would allow a reviewer to prove that the agent never saw `sdp-trace`‑related environment variables, commands, or output, or that the recorder did not indirectly steer the agent. The entire “passive” property rests on Codex’s setup discipline, which is invisible to a CTO or third‑party reviewer.  
This undermines the demo’s core claim that the packet reflects an *untainted* agent workflow.

**Required fix**
- Add a machine‑verifiable proof of passivity (e.g., a recording of the *full* terminal session with a hash of the raw input, or a `tee`‑style capture that can be replayed), **or**
- explicitly state in the demo that passivity is not independently verifiable in this slice, and that the demo illustrates packet structure + tool path, not a complete trust chain.

### C3 – First‑packet minimum bar is not gated at demo‑claim time
The spec defines a minimum bar for the first CTO‑visible packet (at least four rows `pass`/`partial`, `PC‑CHANGE`+`PC‑MUTATION` with resolvable evidence, at least one of `PC‑VERIFICATION`/`PC‑REVIEW`/`PC‑AGENT‑ROUTE` assessed).  
The tasks (T013–T014) describe selecting or creating a feature, but the demo success criteria (SC‑001) only require *a* packet to exist—they do not require a check that the minimum bar was actually met. The demo could produce a packet full of `not_assessed` and still tick the “success” box.  
This would produce a “happy‑path” packet that fails the product proof, leaving the CTO with a misleading picture.

**Required fix**
Insert a gate before the buyer‑demo rehearsal (T021): validate the first packet against the minimum bar using the 006 validator (or a simple row‑state check) and record the pass/fail result in the demo tracker. The rehearsal must be blocked if the bar is not met.

---

## Major Findings

### M1 – Theater assessor is undefined
The negative‑theater example requires triggering `agent_claimed_verification` and a theater finding.  
The 006 slice does not appear to contain a theater assessment module (the current spec for 006 is not in this review, but 007 describes theater rows with reason codes without defining how the assessor works). Without an implementation plan, the negative‑demo packet cannot actually be tool‑generated; it risks becoming a hand‑authored demo.

**Required fix**
Either:
- include a minimal theater rule‑engine in the 006 slice that can detect at least `agent_claimed_verification` from evidence bundle contents, **or**
- mark the theater row of the 007 negative demo as “manual assessment” and record that limitation before the buyer rehearsal.

### M2 – Boundary between “setup‑only” and feature work is not machine‑checkable
Codex may author setup‑only infrastructure, but the spec does not provide a deterministic rule (e.g., allowed file paths) or a CI check to prevent mixing feature code into setup commits. The demo tracker relies on prose labels (`setup_only`) that violate the trust rules.

**Required fix**
Define a file‑scope rule for setup‑only commits (e.g., only `.github/`, `.sdp‑trace/`, `docs/`, and non‑functional build files) and add a CI guard that rejects any functional source or test change in a setup‑only commit. Alternatively, require that such commits be reviewed by an independent reviewer not involved in the demo coordination.

---

## Minor Findings

- **Packet generation provenance**: The spec says the packet must be tool‑generated, but does not require recording the renderer version and input hashes. This makes re‑generation by a reviewer impossible. Add a metadata block in the packet with the tool version, command, and hash of the evidence bundle input. (Enhances reproducible proof.)
- **Git author identity for agent commits**: The spec does not require that commits from the OpenCode/GSD route be distinguishable from Codex commits. A future reviewer might doubt authorship. Defer until the demo identity infrastructure matures.
- **Retroactive task binding is marked `partial`**: The spec correctly caps `PC‑INITIATOR` at `partial` for retroactive bindings, but does not require a reviewer to confirm the binding timestamp against the original PR/commit history. Low risk for the demo; can be noted as a future check.

---

## Required Fixes Summary
1. **C1 & C2**: Strengthen observation integrity and passivity proof, or downgrade demo claims accordingly.
2. **C3**: Gate the first CTO‑visible packet against the minimum bar before claiming demo success.
3. **M1**: Define or defer the theater assessment capability used in the negative demo.
4. **M2**: Machine‑check the setup‑only boundary.

Implementation work on 007 should not proceed until C1, C2, and C3 are addressed. M1 and M2 may be resolved in parallel if the spec is updated to record their deferred status explicitly.
