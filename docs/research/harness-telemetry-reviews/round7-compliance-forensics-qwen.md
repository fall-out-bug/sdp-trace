VERDICT: ACCEPTABLE_WITH_GAPS

Convergence assessment:
- Can this brief be used to start a v0 implementation? yes
- Are there any critical blockers before implementation? yes

Critical blockers:
- None. The hash-chain event model, MissingEvidenceTable, and four-axis verifier output constitute sufficient forensic scaffolding to begin v0.

Major gaps:
- **Trace-to-PR/CI linkage is implied, not explicit.** `run_started` records task ref and parent run refs, but there is no mandatory `pr_ref` or `ci_pipeline_id` field binding the trace to a specific merge gate artifact. Forensics needs this as a first-class join key.
- **No explicit investigation-export surface.** `explain` covers error cases well, but there is no `report` or `export` command that produces a signed, self-contained audit packet (events + chain head + missing table + witness assertions) suitable for incident handoff to legal/compliance.
- **Retention lifecycle is underspecified for cold cases.** "Implementation-defined" local cap is fine for v0, but the brief should mandate a minimum `expiry_time` field on every event so that a forensic query at T+30d can distinguish "expired" from "never recorded."
- **Post-hoc tamper detection relies on demo, not event.** There is no canonical `tamper_detected` or `chain_integrity_failure` event; integrity breakage is only surfaced via `explain` output. For audit, a verifier-level integrity failure should be a first-class event in the chain.

Accepted V0 limitations:
- Local ephemeral keys are host-compromise vulnerable; gate-grade trust requires CI or external witness only.
- No raw prompt/response capture; redaction is digest-only by default.
- Adapter identity may be `self_claimed`; this is acceptable as long as the verifier never upgrades it past the honest floor.
- No reliable `attach <pid>`; late-attach emits explicit boundary events.

Minimum viable correction:
- Add `pr_ref` and `ci_pipeline_run_id` as optional-but-recommended fields on `run_started` and `verifier_result_observed`.
- Add a `chain_integrity_failure` event (or overload `checkpoint_signed` with a failure variant) that the verifier emits before any other output when hash-chain validation fails.
- Add `sdp-trace export <run-dir> --format <audit-bundle|json>` to produce a signed forensic packet including the event chain, MissingEvidenceTable, verifier result, and retention manifest.
- Require `expiry_time` on every retention descriptor; forbid `null` in v0 (derive from profile default).

Questions before implementation:
- Can `recorder_attached` include a `ci_pipeline_id` when running under a known CI environment, or must that come from `run_started`?
- For `chain_integrity_failure`: should it be a verifier-only event appended to an existing chain, or should a new `integrity_audit` record be produced in a separate file if the on-disk chain is already corrupted?
- What is the expected format for the audit export bundle (e.g., JSON Lines + detached signature, or a single signed JSON envelope)?

Demo changes required:
- **Demo 4 (tamper):** After mutating events, use the new `sdp-trace export` to show that the audit bundle carries the `chain_integrity_failure` event and verifier downgrade.
- **Demo 7 (forensic cold query):** Add an `export` invocation showing the audit packet with `expiry_time` fields resolved, distinguishing `expired` vs `missing_telemetry` rows.
- Add a brief **Demo 4.5** (or extend Demo 4) that queries an unrelated run directory and emits `No sdp-trace run found` to validate the `explain` cold-path against `expected_run_absent`.
