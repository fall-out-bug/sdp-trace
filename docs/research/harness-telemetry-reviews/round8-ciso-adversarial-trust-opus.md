```text
VERDICT: ACCEPTABLE_WITH_GAPS

Convergence assessment:
- Can this brief be used to start a v0 implementation? yes
- Are there any critical blockers before implementation? no

Critical blockers:
- none

Major gaps:
- Witness independence is defined structurally (§10: not fork/spawn child,
  records `witness_independence`), but key custody for CI signing is not
  specified. If OIDC/KMS is acceptable v0, say so; otherwise an attacker
  with CI job write access can still sign. Pick one: OIDC-issued token
  bound to job identity, or KMS with non-agent-accessible key.
- `expected_run_absent` predicate (§8) depends on CI artifact store lookup
  by task ref + commit/tree digest. The brief does not require that the
  CI preflight step itself be outside agent control. Without that, the
  agent can suppress the preflight and the gap is never raised. Spec
  should state preflight runs in a CI step the agent cannot skip.
- Adapter identity (§6): `self_claimed` adapters can still register and
  emit events. Brief correctly says trust does not upgrade, but does not
  forbid an adapter from emitting `checkpoint_signed` or
  `verifier_result_observed`. Recommend: restrict signer-bearing event
  types to verified-identity observers; self-claimed adapters limited to
  their declared `allowed_event_types` minus signing-class events.
- Replay defense (§10) lists the witnessed tuple but does not mandate a
  freshness element (nonce or CI-issued challenge) in the tuple. Recorder
  nonce from `recorder_attached` should be explicitly part of the signed
  tuple so an old trace cannot be re-presented against a new CI run.
- Chain corruption handling (§12): `integrity_audit` is written as a
  separate record, good. But its location/signing is unspecified. If it
  lives in the same run dir with the same ephemeral key, the agent can
  delete or forge it. Require CI-side integrity_audit when verification
  runs in CI; local integrity_audit is advisory only.

Accepted V0 limitations:
- Per-run ephemeral key = local structural continuity only, never
  gate-grade. Already stated clearly in §10.
- Pre-expected-run deletion is undetectable without preflight hook.
  Stated in §8.
- Local-only traces cannot reach `ci_witnessed`. Stated in §11, §2.
- Raw prompt/response capture out of scope. Stated in §2, §3.
- Host compromise defeats local chain. Stated in §10.

Minimum viable correction:
- Add a short "CI witness key custody" subsection under §10 naming
  acceptable v0 mechanisms (OIDC job token, KMS, hardware key) and
  explicitly excluding keys readable by the agent's workspace or process.
- In §10, add recorder nonce to the CI-witnessed tuple list.
- In §6, add: "Self-claimed adapters cannot emit `checkpoint_signed` or
  `verifier_result_observed`; these require verified signer identity."
- In §8, add: "CI preflight that raises `expected_run_absent` must run
  in a CI step not controllable by the agent under test."
- In §12, add: "When verification runs in CI, `integrity_audit` is
  emitted and signed by the CI witness, not by the local recorder."

Questions before implementation:
- Is OIDC-per-job-identity the v0 CI signer, or is a KMS key profile
  also in scope? This drives the witness authority policy schema.
- Does the authority policy (referenced by §6, §10) get its own
  artifact, or is it implicit in CI config for v0?
- For `witness_independence: same_job`, is that downgraded to
  `local_observed` or to an intermediate scope? §11 does not name one.

Demo changes required:
- Demo 2: add a variant where the CI signing step is moved into the same
  process as the verifier; verifier must record `same_process` and
  refuse `ci_witnessed`.
- Demo 5: make replay explicitly bind the old recorder nonce to a new
  commit and show the nonce mismatch as the downgrade reason, not just
  workspace/commit mismatch.
- Demo 4: show an `integrity_audit` produced by CI (not local) when the
  local chain is tampered, to prove the audit record itself is not
  agent-forgeable.
- Demo 3: add a case where a `self_claimed` adapter attempts to emit a
  signing-class event and the recorder rejects it.
```

Tie-ins: signing/verification (key custody, tuple freshness, signer
class restriction), observer authority (self-claimed vs verified event
scope), forensic query/replay (nonce binding, CI-signed integrity_audit),
demo credibility (same-process witness, forged integrity record). None
of these prevent starting v0; they are bounded corrections plus four
demo additions.
