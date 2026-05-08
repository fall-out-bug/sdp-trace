```text
ROUND 2 FORGERY RED TEAM — Verifier/State-Machine Adversary Classification
============================================================================

Methodology: For each attack, I assume the brief's architecture is
implemented as described (layers, hash chain, signing levels, verifier
states). I evaluate what the verifier can actually determine, not what
the brief claims it should determine.

STATE MODEL CLARIFICATION (prerequisite)
-----------------------------------------
The brief's 8 states conflate verdicts and trust scopes. Before classification,
the functional decomposition from round 1 (Mimo) must be stated:

  Verdict:      pass | fail | cannot_verify | not_assessed
  Trust scope:  local_only | witnessed | externally_witnessed
  Completeness: complete | missing_telemetry | partial

The brief does not make this split. All classifications below assume
the split is NOT made (current brief state), making verifier behavior
ambiguous for compound cases. I note where the split would change
the outcome.

ATTACK CLASSIFICATIONS
======================

A1. Post-hoc chain fabrication (agent writes complete plausible chain
    after session ends, before verifier runs)
    ---------------------------------------------------------------
    Status: NOT COVERED

    Local hash chain is internally consistent by construction.
    No monotonic timestamp authority exists. No witness saw the chain
    grow. Verifier has no live-capture proof. It emits `local_only` +
    `pass` — a false positive.

    If verdict×scope split existed: verdict remains `pass` (chain is
    structurally valid), trust scope is `local_only`. Consumer must
    know `local_only`+`pass` ≠ trustworthy. The brief does not
    communicate this distinction.

    Required fix: live-capture proof (e.g., checkpoint hashes sent to
    witness at intervals, or recorder process attestation via separate
    launch identity).

A2. Agent signs its own telemetry as "recorder"
    --------------------------------------------
    Status: NOT COVERED

    No identity attestation binds the signing key to the recorder
    process vs. the agent process. Both run in the same user context.
    Verifier sees a valid signature and emits `local_only` + `pass`.

    If key were hardware-bound or process-attested: detected (signature
    identity mismatch). Neither is specified.

A3. Event deletion from middle of hash chain
    ------------------------------------------
    Status: DETECTED (hash integrity break)

    Hash linking makes deletion detectable: chain is broken at the
    gap. Verifier emits `fail`.

    Caveat: only if the verifier re-checks the full chain from genesis.
    If it checks only the signed checkpoint forward, gaps before the
    checkpoint are invisible. Checkpoint frequency is unspecified.

A4. Event reordering within hash chain
    ------------------------------------
    Status: DETECTED (hash integrity break)

    Same as A3. Reordering breaks hash links.

A5. Event mutation within hash chain
    ----------------------------------
    Status: DETECTED (hash integrity break)

    Same as A3. Mutation breaks forward hash.

A6. Event injection (adding events not from original run)
    -------------------------------------------------------
    Status: DETECTED only if injection occurs between checkpoints
            or after chain head signed by witness

    If injection happens before any witness checkpoint and chain is
    re-hashed: NOT COVERED. The agent controls the chain before the
    first witness anchor.

A7. Recorder process killed mid-run, agent continues
    --------------------------------------------------
    Status: NOT COVERED

    Partial chain is structurally valid up to the kill point.
    Subsequent events are missing. Verifier sees a short chain with
    no explicit "session end" event. It cannot distinguish "short
    legitimate task" from "killed recorder, unrecorded work."

    Classified as `missing_telemetry` only if a completeness contract
    exists. No such contract is specified.

A8. Subagent bypasses tool/shell wrapper
    --------------------------------------
    Status: DOWNGRADED TO cannot_verify

    File watcher sees final mutations but cannot attribute them to a
    tool call or shell command. Harness adapter sees "subagent
    completed" but not what the subagent did. Verifier sees:
      - harness: task done
      - tool/shell: no events for the mutations
      - file watcher: files changed
    This is a layer disagreement with no resolution rule. Current
    brief: no covered state. With state split: verdict `cannot_verify`,
    completeness `partial`.

A9. Two harnesses run concurrently, events interleave
    ---------------------------------------------------
    Status: NOT COVERED

    No `session_id` or `harness_id` field is specified. File watcher
    and shell wrapper cannot attribute events to a session. Verifier
    produces a merged chain with cross-contamination.

    Required fix: mandatory session correlation field in every event.

A10. Gateway telemetry shows model call, no local action telemetry
    ----------------------------------------------------------------
    Status: DOWNGRADED TO cannot_verify

    Brief raises this as a review question but provides no verifier
    behavior. The model was called but nothing happened locally —
    possible reasons: model refused, agent crashed, agent bypassed
    recorder. Verifier cannot determine which.

    With completeness contract: `missing_telemetry` for action layer.
    Without: `cannot_verify`.

A11. Local action telemetry exists, no gateway provenance
    -------------------------------------------------------
    Status: DOWNGRADED TO cannot_verify

    Files changed, commands ran, but no model call is recorded.
    Possible: gateway was down, agent used cached response, agent
    fabricated local actions. Verifier cannot distinguish.

    With completeness contract requiring gateway evidence:
    `missing_telemetry` for provenance layer. Without: `cannot_verify`.

A12. Temporal correlation attack (gateway time T, file changes T+2s,
    no causal proof)
    -------------------------------------------------------------------
    Status: NOT COVERED

    The brief assumes correlation equals causation. It does not.
    Verifier emits `pass` with `local_only` or `witnessed` scope,
    implying provenance linkage that does not exist.

    This is a trust-scope overclaim. The verifier should emit a
    lower trust scope or `cannot_verify` for the specific causal
    claim, but no causal verification logic is specified.

A13. Task rewrite after evidence capture
    --------------------------------------
    Status: NOT COVERED

    Task is an agent-written enrichment field (Section 4). Agent can
    rewrite it at any time. Verifier has no task immutability
    checkpoint. If task changes after file/command evidence, the
    verifier sees the final task text and matching evidence — but
    the evidence was captured against a different task.

    Required fix: task hash bound into chain at session start, before
    any tool/command events. Current brief: no task-binding mechanism.

A14. Scope inflation via task redefinition
    -----------------------------------------
    Status: NOT COVERED (variant of A13)

    Same mechanism. Scope boundary is an agent-reported field.
    Redefining scope post-evidence makes out-of-scope changes appear
    in-scope. No scope hash checkpoint exists.

A15. Agent redacts events then re-links hashes
    ---------------------------------------------
    Status: NOT COVERED (local-only runs)

    If chain is local-only, agent has write access to the chain
    store. Redaction + re-link produces a valid chain. Only detected
    if a witness previously anchored a checkpoint that references
    the redacted events. Otherwise: `pass` + `local_only`.

    If externally witnessed: DETECTED (witness anchor conflicts with
    modified chain).

A16. CI signs chain head without re-verifying
    -------------------------------------------
    Status: DOWNGRADED TO cannot_verify (for CI-witnessed scope)

    CI signature proves "CI saw this hash," not "CI verified this
    chain." The signature is valid but the trust claim is inflated.

    Verifier emits `witnessed` — which the brief defines as CI or
    remote witness. This is a trust-scope overclaim unless CI also
    runs the verifier and signs the verdict.

    Required fix: CI must run verifier and sign verdict + chain head,
    not just chain head.

A17. Slow replay attack (replay witnessed run as new run)
    --------------------------------------------------------
    Status: NOT COVERED

    Witness saw hash H at time T. Agent later produces identical
    content hash H and claims it as a new run. No nonce or run-id
    binding in the witness protocol. Verifier sees valid witness
    record and emits `externally_witnessed`.

    Required fix: witness must bind to a run-scoped nonce, not just
    content hash.

A18. Agent uses unapproved model through impersonating proxy
    -----------------------------------------------------------
    Status: NOT COVERED

    Gateway logs model X with provider P. Proxy presents as model X
    to the gateway. Gateway cannot distinguish approved model from
    proxy. Verifier sees valid gateway record and emits `pass`.

    Requires: approved provider list at gateway, or model response
    attestation from provider. Neither is specified.

A19. Harness adapter lies about tool calls
    -----------------------------------------
    Status: NOT COVERED

    Adapter reports "tool: read_file(path=X)" but agent actually
    executed "tool: write_file(path=Y)". Verifier has no independent
    observation layer unless tool wrapper is also present AND
    cross-correlated. Brief assumes layers agree; no disagreement
    resolution exists.

    If tool wrapper confirms write_file: DETECTED (layer mismatch).
    If no tool wrapper: NOT COVERED.

A20. Backfilled telemetry after failure
    --------------------------------------
    Status: DOWNGRADED TO cannot_verify (without witness)

    Agent fails, then produces clean telemetry. Locally: valid chain.
    Without external time anchor, verifier cannot prove backfill.
    Emits `local_only` + `pass`.

    With external witness that saw absence of checkpoints during the
    original run window: DETECTED. But brief does not define absence-
    of-evidence witnessing.

A21. Partial telemetry with successful summary
    ---------------------------------------------
    Status: NOT COVERED

    Agent produces 3 of 10 expected events, writes a "task completed"
    claim. Verifier has no completeness contract to determine that
    7 events are missing. Emits `pass` unless it has an expected-event
    manifest. No such manifest mechanism exists.

    This is the most dangerous undetected forgery: the CTO sees
    "completed, passed" when the run is incomplete.

A22. Recorder and agent in same process (attestation bypass)
    -----------------------------------------------------------
    Status: NOT COVERED

    Brief never requires process isolation. If both are in-process,
    all "recorder-signed" events are agent-accessible. Verifier sees
    valid signatures from what it believes is an independent recorder.

    Required fix: process attestation (PID, launch-time binding,
    separate binary) enforced by verifier.

A23. File watcher sees atomic swap, misses intermediate mutations
    -----------------------------------------------------------------
    Status: DOWNGRADED TO cannot_verify

    File watcher sees final state only. Intermediate unauthorized
    writes are invisible. Verifier can confirm final state matches
    scope but cannot confirm the path was clean.

    Trust scope should be downgraded, but no downgrade trigger exists.

A24. Gateway prompt hash inverted for short prompts
    --------------------------------------------------
    Status: NOT COVERED (privacy attack, not forgery, but relevant)

    Short prompt hashes can be brute-forced. If gateway logs prompt
    hashes, an adversary who compromises the log can reconstruct
    prompts. This violates the CTO's privacy expectation and makes
    the gateway layer a data exfiltration vector.

    Classification: not a forgery but a trust boundary violation.
    Gateway protocol must specify minimum entropy or salted hashing.


SUMMARY TABLE
=============

| ID  | Attack                                      | Classification      | Layer gap                  |
|-----|---------------------------------------------|---------------------|----------------------------|
| A1  | Post-hoc chain fabrication                  | NOT COVERED         | No live-capture proof      |
| A2  | Agent self-signs as recorder                | NOT COVERED         | No signer attestation      |
| A3  | Mid-chain deletion                          | DETECTED            | Hash linking works         |
| A4  | Event reordering                            | DETECTED            | Hash linking works         |
| A5  | Event mutation                              | DETECTED            | Hash linking works         |
| A6  | Event injection pre-witness                 | NOT COVERED         | No pre-witness anchor      |
| A7  | Recorder killed mid-run                     | NOT COVERED         | No completeness contract   |
| A8  | Subagent bypass                             | DOWNGRADED cannot_v | No session boundary        |
| A9  | Concurrent harness interleaving             | NOT COVERED         | No session correlation     |
| A10 | Gateway only, no local                      | DOWNGRADED cannot_v | No layer completeness      |
| A11 | Local only, no gateway                      | DOWNGRADED cannot_v | No layer completeness      |
| A12 | Temporal = causal assumption                | NOT COVERED         | No causal verification     |
| A13 | Task rewrite post-evidence                  | NOT COVERED         | No task immutability       |
| A14 | Scope inflation                             | NOT COVERED         | No scope checkpoint        |
| A15 | Redaction + re-link                         | NOT COVERED*        | No local append-only       |
| A16 | CI signs without verifying                  | DOWNGRADED cannot_v | CI must sign verdict       |
| A17 | Slow replay                                 | NOT COVERED         | No nonce in witness proto  |
| A18 | Model impersonation via proxy               | NOT COVERED         | No provider attestation    |
| A19 | Harness adapter lies                        | NOT COVERED*        | No cross-layer agreement   |
| A20 | Backfill after failure                      | DOWNGRADED cannot_v | No absence-of-evidence     |
| A21 | Partial telemetry, complete summary         | NOT COVERED         | No completeness contract   |
| A22 | Same-process recorder                       | NOT COVERED         | No process attestation     |
| A23 | Atomic swap hides intermediates             | DOWNGRADED cannot_v | No path audit              |
| A24 | Prompt hash inversion (privacy)             | NOT COVERED         | No hash entropy floor      |

* Conditional: detected only if secondary layer present and correlated.

VERDICT: 24 attacks assessed. 3 detected (A3-A5, all hash-link integrity).
         5 downgraded to cannot_verify. 16 not covered.

STRUCTURAL GAPS (ordered by forgery impact)
============================================

1. COMPLETENESS CONTRACT (blocks A7, A21, and partially A10/A11)
   Without a machine-readable declaration of expected events per run,
   the verifier cannot detect absence. This is the single highest-
   impact gap. Without it, `pass` is meaningless for any non-trivial
   run.

2. LIVE-CAPTURE PROOF (blocks A1, A20, partially A15)
   Hash chain proves ordering, not timing. A chain written in 2
   seconds after a 30-minute session is structurally identical to one
   grown over 30 minutes. Need: checkpoint witnesses at intervals, or
   monotonic external timestamp binding.

3. SIGNER ATTESTATION (blocks A2, A22)
   Signing key must be bound to a process identity, not just a key.
   Without this, "local recorder signed" is indistinguishable from
   "agent signed."

4. SESSION CORRELATION (blocks A8, A9, A19)
   Every event must carry: session_id, harness_id, parent_session_id,
   request_id (for gateway correlation). Without this, multi-session
   and subagent scenarios produce unverifiable merged streams.

5. TASK/SCOPE IMmutABILITY (blocks A13, A14)
   Task and scope must be hashed into the chain at session start,
   before any tool/command events. Post-hoc rewrites must break the
   chain or be visible as a supersession event.

6. WITNESS NONCE BINDING (blocks A17)
   Witness must anchor (run_id, nonce, hash, timestamp), not just
   (hash, timestamp).

7. LAYER AGREEMENT RESOLUTION (blocks A10, A11, A12, A19)
   Verifier needs explicit rules for layer conflicts:
   - gateway + no local → cannot_verify (not pass)
   - local + no gateway → cannot_verify for provenance claims
   - harness says X, tool wrapper says Y → fail unless reconciled

STATE MACHINE PROPOSAL (minimum viable)
========================================

Verdict (per event and per run):
  pass | fail | cannot_verify | not_assessed

Trust scope (per run):
  local_only | ci_witnessed | externally_witnessed

Completeness (per run):
  complete | partial | unknown

Compound state example:
  { verdict: pass, scope: local_only, completeness: complete }
  → "Chain is intact, all expected events present, but trust is
     local-only and cannot resist host-level forgery."

  { verdict: cannot_verify, scope: externally_witnessed, completeness: partial }
  → "Chain head was witnessed externally but expected events are
     missing. Witness proves chain existed, not that it was complete."

This 3-axis model covers every case above with no ambiguity. The
current 1-axis model covers 3 of 24.
```

Save this to a file if you want it persisted in the reviews directory. The main takeaway: **hash linking gives you 3 out of 24 attacks detected. Everything else requires completeness contracts, live-capture proof, and signer attestation — none of which exist in the brief.** The verifier state model must be 3-axis (verdict × scope × completeness) before any of the forgery classifications become reliable.
