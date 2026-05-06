# Round 2 Forgery Red Team: Classification

**Source text**: `docs/research/harness-telemetry-trust-brief.md`
**Reviewer role**: operational attacker with local user access, no root/OS compromise assumed
**Method**: trace each attack case through the brief's described mechanisms (hash chain, verifier states, signing model, adoption rungs). Classify against what the brief *describes*, not what it should describe.

---

## Results

Of **25 attack cases** evaluated: **0 prevented, 2 detected, 0 downgraded, 23 not covered**.

---

## Detected (2)

| # | Attack | Detection | Limitation |
|---|--------|-----------|------------|
| D1 | Interior event mutation without re-chaining | Hash chain breaks at next event; verifier emits `fail` | Adversary with chain write access simply recomputes hashes → reclassified as not covered |
| D2 | Event reordering without re-chaining | Hash chain's sequential `prev_hash` breaks; verifier emits `fail` | Same escape: competent adversary recomputes. Detection only holds against naive tampering on a frozen chain |

These are the *only* cases where the brief's described mechanisms produce a reliable detection signal. Both assume the adversary modifies an existing chain in place without recomputing dependent hashes, which is the least competent attack profile.

---

## Not Covered (23)

### Bypass / Process Tree

| # | Attack | Why Not Covered |
|---|--------|-----------------|
| N1 | Silent bypass via subagent delegation | Section 3 table claims "subagent boundaries" captured; no subagent event schema exists. Parent harness sees "subagent completed" with zero tool/file/command evidence from the subagent. Verifier cannot link subagent actions to parent run. Section 9 confirms missing: no subagent/session correlation. |
| N2 | Child process bypass (python -c, direct syscalls) | Shell wrapper captures only explicit shell invocations. Tool wrapper captures harness-declared tool calls. Neither intercepts child processes spawned inside a tool's runtime. File watcher sees mutations but cannot attribute to bypass. No process-tree monitoring described. |
| N3 | Alternative interpreter / inline code bypass | Same mechanism as N2. `os.system()`, `subprocess.Popen()`, native extensions all bypass the wrapper boundary. No eBPF, ptrace, or kernel-level interposition mentioned. |
| N4 | Agent bypasses harness plugin entirely | Managed mode presumes plugin is installed and non-bypassable. The brief provides no enforcement mechanism for "cannot run without plugin." If agent launches without plugin, zero harness-level telemetry exists. Tool/shell wrappers may still capture actions but lose task/model context. Managed mode enforcement is undefined. |
| N5 | Selective tool bypass (wrapped tools for clean actions, unwrapped for attacks) | Tool wrapper captures some calls; no completeness contract exists to declare which tools were expected. Verifier sees partial telemetry with no `missing_telemetry` flag because it has no expected-events baseline. Section 9 confirms missing: tool capture implementation incomplete. |

### Local Chain Integrity

| # | Attack | Why Not Covered |
|---|--------|-----------------|
| N6 | Event deletion with re-chaining | Agent removes events, recomputes subsequent `prev_hash` values. Chain is internally consistent. Hash linking only detects tampering on an *existing* chain; it does not prevent an adversary who controls the chain from generating a new valid chain. |
| N7 | Event deletion from chain tail | No chain-end marker or expected-length invariant. Agent deletes last N events; no subsequent event exists to break. Verifier sees shorter chain, no detection. Partial run indistinguishable from short legitimate run. |
| N8 | Full chain rewrite | Agent generates new chain post-hoc with consistent hashes. Hash chain self-verifies. Verifier has no mechanism to distinguish "captured live" from "written after session" for local-only telemetry. |
| N9 | Post-hoc plausible chain (zero live events) | Extension of N8. Agent runs entirely without recording, then fabricates complete plausible chain. No liveness proof: no monotonic timestamp authority, no external time source, no witness check-in during session. |
| N10 | Recorder on-disk store corruption | Recorder's data store lives on same filesystem, same user context as agent. No file permission isolation described. Agent can overwrite, append, or delete recorder data directly. |
| N11 | Recorder process killed mid-run | Agent kills recorder process, continues work, commits. Partial chain lacks end-of-chain marker. Verifier cannot distinguish "recorder crashed" from "session naturally short." No heartbeat or completeness invariant in design. |
| N12 | Agent self-signs as "local recorder" | Signing model says "local recorder-signed" implies local integrity. But no key provenance distinction exists between agent keys and recorder keys. Agent generates key pair, signs chain, declares source="local_recorder." Verifier accepts signature; cannot verify signer identity. |

### Witness / External Trust

| # | Attack | Why Not Covered |
|---|--------|-----------------|
| N13 | Pre-witness event rewrite | Agent produces chain, rewrites events, recomputes chain head hash, then submits to witness. Witness sees consistent hash, signs it. Witness confirms "head existed at time T" but not "head matched this event sequence." Content substitution undetectable unless witness receives and stores chain content. Protocol undefined. |
| N14 | Witness chain fork (same head, different content) | Same mechanism as N13. If witness only signs the head hash, content binding is absent. Hash collision or deliberate content swap produces two chains with same head. Witness anchor is non-discriminating. |
| N15 | Witness endpoint spoofing | Agent points to fake witness endpoint. Brief specifies no witness identity verification: no certificate pinning, no known-endpoint registry, no transparency log operator identity. |
| N16 | Slow-replay attack (identical task replay) | Agent captures valid witnessed run for task X. Later replays identical task X (same content hash) and presents old witness as evidence. No nonce, session-id, or freshness binding in witness protocol. Witness binding to task identity is unspecified. |
| N17 | CI signs without re-running verifier | Section 5: "CI identity signs chain head and artifacts." Brief acknowledges "CI signing unverified telemetry" as threat (Section 6). No mitigation specified: no verifier re-run requirement, no chain content validation prerequisite. CI becomes blind trust-stamping machine. |

### Cross-Layer Correlation

| # | Attack | Why Not Covered |
|---|--------|-----------------|
| N18 | Gateway-to-local temporal-ambiguity attack | LLM gateway proves model X called at time T. Local telemetry shows file changes at T+2s. Nothing proves model response *caused* those changes. Correlation is temporal, not causal. No request-id propagation protocol exists (Section 3A lists as "review question"). Agent can claim gateway-backed provenance for changes made for unrelated reasons. |
| N19 | Gateway telemetry exists, local telemetry missing | Brief asks the question (Section 3A) but provides no verifier state. Is this `missing_telemetry` (missing local)? Is it `cannot_verify` (incomplete correlation)? The verifier state for this specific gap is undefined. |
| N20 | Unapproved model via proxy impersonation | Agent uses unapproved model behind proxy that claims approved model identity. Gateway capture logs self-reported model name. No provider certificate enforcement, no model fingerprinting, no approved-provider list enforcement. |
| N21 | Concurrent harnesses interleaving events | Two harnesses running on same repo produce interleaved file mutations and shell commands. No session_id or harness_id boundary in event schema. File watcher and shell wrapper cannot attribute events to source harness. |

### Evidence Tampering & Scope

| # | Attack | Why Not Covered |
|---|--------|-----------------|
| N22 | Test artifact poisoning | Local recorder captures test commands and exit codes at time T1. Agent modifies test suite after capture, before CI runs at T2. CI witnesses different test state. No test-to-evidence binding (e.g., test file hash at capture time). No replay integrity. CI sees passing tests but cannot compare to capture-time test state. |
| N23 | Scope inflation post-evidence | Agent records scope S1, executes within S1, captures evidence. Then rewrites task to S2 (broader scope) before verifier runs. Verifier sees scope S2 and evidence within S2 → passes scope check. No scope snapshots bound to event chain at capture time. Section 6 lists "task rewrite after command evidence" as threat; no mechanism provided. |
| N24 | Atomic file swap hiding intermediate mutations | Agent stages unauthorized writes in temp directory, atomically swaps target file. File watcher sees single final mutation, not intermediate mutations. Tool wrapper may see intermediate writes if intercepted, but atomic swap at OS level can bypass both. No inotify/fsevents granularity described. |
| N25 | Redaction legitimacy ambiguity | Agent redacts failed test run, claims privacy. No redaction policy framework: any redaction is either always suspicious (false positives) or always accepted (exploitable). Brief lists "redaction hiding critical evidence" as threat (Section 6) but provides no distinction between legitimate and adversarial redaction. |

---

## Summary Table

| Category | Prevented | Detected | Downgraded | Not Covered |
|----------|-----------|----------|------------|-------------|
| Bypass / Process Tree | 0 | 0 | 0 | 5 |
| Local Chain Integrity | 0 | 2* | 0 | 5 |
| Witness / External Trust | 0 | 0 | 0 | 5 |
| Cross-Layer Correlation | 0 | 0 | 0 | 4 |
| Evidence Tampering & Scope | 0 | 0 | 0 | 4 |
| **Total** | **0** | **2*** | **0** | **23** |

*D1 and D2 detect only naive tampering; competent adversary with chain write access escapes both.

## Core Finding

The brief's hash chain mechanism detects tampering on a frozen chain against an agent that cannot recompute hashes. This is the *only* operative defense.

Every other attack succeeds because:

1. **No completeness contract**: verifier cannot distinguish "short run" from "truncated run."
2. **No liveness proof**: verifier cannot distinguish "events captured live" from "events fabricated post-hoc" for local-only telemetry.
3. **No recorder attestation**: verifier cannot distinguish agent-signed from recorder-signed.
4. **No witness content binding**: witness signs hash without receiving chain content.
5. **No correlation protocol**: gateway and local events share no trace ID.
6. **No CI verification precondition**: CI can sign unverified telemetry.
7. **No scope/time binding**: task scope and test state are not pinned to event chain at capture time.

The brief's threat list (Section 6) names 13 threat categories correctly. Exactly zero have operative mitigations in the described design.
