```
VERDICT: ACCEPTABLE_WITH_GAPS

Convergence assessment:
- Can this brief be used to start a v0 implementation? yes
- Remaining blockers before implementation:
  1. Local signing key material lifecycle undefined for V0.
  2. Observer role authentication mechanism unspecified.
  3. Missing-run detection absent from V0 boundary.

Critical blockers:
- Key material management is a named concern (section 7) but has no
  V0 resolution. The brief states "Agent-accessible keys cannot sign
  local_observed or stronger scopes" without specifying enforcement.
  In V0, the local wrapper and the agent share a host and likely a
  user context. If the wrapper's signing key is in the workspace or
  user keychain that the agent can read, the authority rule is an
  aspiration, not a constraint. This is not about key rotation or
  PKI complexity; it is about one sentence: where does the V0
  local_recorder key live, and can the agent process read it?
  (signing and verification, observer authority)

- Observer role is asserted, not authenticated. Any process can write
  `"observer_role": "ci"` into an event. The verifier checks schema
  and chain integrity but section 8 does not state how it establishes
  that the signer is actually authorized for the claimed role. For V0
  this is acceptable IF the brief explicitly scopes it as "role is
  claimed, trust_scope is the honest floor, not the claim." Currently
  the brief implies but does not state this.
  (observer authority, signing and verification)

Major gaps:
- Missing-run detection. The brief states "If work happens outside the
  wrapper, there may be no run." This is honest but incomplete. There
  is no mechanism for sdp-gate or sdp-report to know a run should have
  existed. The VCS diff evidence (before/after digests) partially
  addresses this: if source changed without a run, that is detectable.
  But the brief does not connect this evidence to a "no-run gap" finding
  in the verifier. A single addition to the MissingEvidenceTable —
  something like `expected_run_absent` — would close this.
  (v0 capture boundary, CI/gate anchoring)

- Nonce/replay detection mechanism under-specified. Section 4 mentions
  a nonce on `recorder_attached`. Section 7 lists replay protections.
  But Demo 5 promises to show "nonce mismatch" without the brief
  defining what nonce protocol is checked. A nonce on the first event
  alone does not detect mid-chain replay or selective event deletion
  followed by re-signing. The prev_event_hash chain handles deletion
  and reorder, but nonce freshness across runs is not defined.
  (provenance model, demo credibility)

- Late attach boundary is acknowledged but its forensic footprint is
  not. If the wrapper starts after work began, `late_attach_boundary`
  is emitted. But the brief does not state what evidence the verifier
  uses to detect late attach. Is it VCS diff before first event? Git
  commit timestamp vs first event time? Without specifying the
  detection signal, the event is a self-report by the late-attaching
  recorder, which is not trustworthy.
  (trace model, v0 capture boundary)

- Contract provisioning by the agent. The contract file is passed via
  `--contract <contract-file>`. If the agent writes this file, it
  controls the expected evidence bar. The brief states contracts must
  be locked before run_started, but locking a contract you wrote is
  not a meaningful constraint. V0 should note that contract
  provenance is a known gap and that CI-witnessed verification
  requires the contract to come from a trusted source (checked into
  VCS, pinned by digest in CI config).
  (expected evidence contracts, CI/gate anchoring)

False assumptions:
- The brief treats `local_observed` as an honest trust floor. This is
  correct for the trust_scope definition but creates a UX/DX risk:
  Demo 0 outputs `pass + local_observed + partial + partial`. A
  developer or CTO reading "pass" on a local-only trace may treat it
  as sufficient without understanding the trust scope. The word "pass"
  carries weight regardless of qualifiers. Consider whether V0 verifier
  should emit `observed` instead of `pass` for non-CI trust scopes,
  or whether the human-readable output must always lead with trust
  scope.
  (CTO usefulness, adoption and DX)

- Section 7 states "CI must sign verifier result, source commit, chain
  head, expected contract digest, verifier version, policy profile,
  timestamp, and CI identity." This assumes the CI job has access to
  the full local event log. If the log is local-only and CI runs in a
  separate environment, the CI verifier either receives the log as an
  artifact (which the agent could have tampered with before upload) or
  CI re-observes independently. The brief does not specify the CI
  handoff model.
  (CI/gate anchoring, provenance model)

Minimum viable correction:
1. Add one paragraph to section 7 specifying V0 local key material:
   where it lives, who can read it, and what trust scope it honestly
   supports. Even "local_recorder key is a per-run ephemeral key
   generated by the wrapper process, stored only in wrapper memory,
   supporting local_observed at most; host compromise defeats this"
   is sufficient.
2. State explicitly that observer role is self-claimed in V0 and that
   trust_scope is the verifier's honest assessment of what that claim
   is worth, not a confirmation of identity.
3. Add `expected_run_absent` or equivalent to the verifier's gap
   taxonomy, keyed on VCS evidence of untraced source changes.
4. Add a sentence to the contract section noting that contract
   provenance is out of scope for V0 and that CI-witnessed verification
   requires contract digest pinning in CI config.

Questions before implementation:
- Is the V0 local_recorder a long-lived daemon or a per-command
  subprocess? This affects key isolation and late-attach semantics.
- Does the CI verifier receive the event log as a build artifact, or
  does it pull from an external store? This changes the tampering
  model.
- Should the V0 human-readable output for `pass + local_observed`
  lead with the trust scope to prevent over-reliance? Or is that a
  sdp-report/sdp-gate concern?

Demo changes required:
- Demo 0: show the human-readable output format. If it says "pass,"
  explain how the reader is prevented from ignoring the trust scope
  qualifier.
- Demo 5: specify the nonce/timestamp mechanism being demonstrated.
  Without this, the demo cannot be implemented honestly.
- Demo 2: clarify the CI handoff. What artifact does CI verify? How
  does CI obtain the local event log?
```
