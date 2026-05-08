```text
VERDICT: CHANGES_REQUIRED

Critical blockers:
- No retention policy model. The brief mentions retention state as a dimension of Evidence but provides no schema for retention windows, archival triggers, or deletion attestations. A forensics lead cannot reconstruct a timeline a month later if the substrate silently loses artifacts after a week. [privacy and retention; CTO usefulness]
- No query or replay surface defined. The brief declares signed events and checkpoints but never specifies how an investigator will reconstruct an ordered timeline from multiple observers across harness, gateway, VCS, and CI. Raw JSON dumps are not a forensics product. [CTO usefulness; evidence model]
- No redaction audit trail. Section 6 mentions redaction/retention state compatibility checks, but Section 11 admits prompts, argv, stdout, stderr, and file paths may leak confidential data. There is no model for what was redacted, who redacted it, when, or what the original digest covered. A chain-of-custody gap. [privacy and retention; signing and verification]

Major gaps:
- Supersession events are mentioned in Section 4 but lack a concrete schema. If a task/spec/plan changes mid-run, a forensics investigation must show the before/after digests, who authorized the change, and whether checkpoints before and after the change remain valid. [trace model; evidence model]
- No explicit `not_assessed` rendering in the verifier output or report. The brief says "never hide not_assessed gaps in prose" in Section 11 but does not mandate a machine-readable missing-evidence ledger in the report shape. Investigators need an explicit gap table, not inferred absence. [CTO usefulness; expected evidence contracts]
- Late-attach boundaries are required (Section 4) but there is no schema for distinguishing "recorder started late" from "recorder was backfilled." Post-hoc fabrication detection cannot rely only on missing pre-gate checkpoints. An explicit late_attach_reason and recorder_start_time are needed. [trace model; signing and verification]
- The verifier's three-axis output (Section 6) is excellent but missing a fourth dimension: `replayable: true | partial | false`. A pass verdict on a trace where stdout/stderr were not retained is forensically useless for incident investigation. [evidence model; signing and verification]

False assumptions:
- Assumes in-toto Statement + DSSE is sufficient for long-horizon investigation without a separate replay/index layer. in-toto is attestations-first; forensics needs timeline-first queries. The brief should explicitly defer the query layer but acknowledge it as a required companion, not an afterthought. [CTO usefulness; standards we should reuse]
- Assumes a single "gate checkpoint" anchors the run before merge. For long-lived features with multiple PRs or stacked agent runs, a single pre-merge checkpoint is insufficient. Investigable units may span multiple chained runs. [CI/gate anchoring; trace model]

Standards we should reuse instead of building:
- OpenTelemetry resource attributes for service/harness/model identity (already listed, ensure they carry through to the canonical event schema, not just runtime spans).
- CloudEvents or OpenTelemetry Logs API for structured, queryable event shapes rather than a bespoke canonical JSON format.
- Sigstore Fulcio certificate transparency entries as a tamper-evident log of who signed what and when, separate from the chain head.
- W3C Verifiable Credentials or similar for observer authority claims, rather than inventing a new authority policy format.

Minimum viable product correction:
- Add a signed retention_descriptor field to every Evidence record: {retention_policy, expiry_time, archival_location, redaction_manifest_digest}. Without it, Demo 1-5 are academically interesting but forensically hollow.
- Add an explicit MissingEvidenceTable to the verifier output listing each expected-observer role that was absent and the sequence range where it *should* have appeared.
- Add a replay/ timeline reconstruction acceptance test alongside the five demos. The test should take a real trace, drop into a forensic query, and ask: "What ran between sequence N and M? Which files changed? Which test claims have supporting artifacts?"

Questions before implementation:
- What is the minimum retention window for evidence to be considered forensically useful? 30 days? 90 days? This determines whether archival is in-scope or must be a hard dependency on an external store.
- Who is authorized to redact raw evidence? Redaction is a security action and should itself be a signed event with its own observer role.
- How do we handle multi-day, multi-commit runs where source commit changes mid-trace? The verifier currently checks source commit correlation but does not specify behavior when commits advance inside a single trace.

Attack scenarios not covered:
- Clock manipulation: a local recorder with a shifted system clock can backdate events to appear pre-gate. The brief requires monotonic sequence but not a minimum clock-source trust requirement (e.g., CI time authority, NTP attestation, or external timestamp).
- Fork-and-replay: clone a valid local trace, swap the source commit digest to match a different commit, resign locally. Without a CI or external timestamp binding, verifier checks 5-7 may pass on a fraudulently re-targeted trace.
- Evidence swap: replace a signed local test result artifact digest with one from a different run where tests actually passed. The digest check verifies integrity of the pointer, not whether the artifact belongs to *this* run's source.
- Redaction abuse: aggressively redact all command argv and stdout, pass verification with completeness=complete for the signed checkpoint, but destroy all forensic surface. The verifier needs a minimum redaction policy bound to the expected evidence contract.

Demo changes required:
- Demo 3 should also show the MissingEvidenceTable explicitly rendered, not just a `missing: [harness_observed, model_identity_observed]` list. Show how an investigator reads this two weeks later.
- Demo 5 needs a stronger post-hoc fabrication test: not just "no checkpoint was witnessed before gate," but "clock is clearly inconsistent with merge time" or "nonce/sequence range contradicts VCS timestamp."
- Add Demo 6: Redaction scenario. Show a trace where sensitive argv/stdout were redacted, the redaction manifest is present, and the forensics lead can still reconstruct which commands ran even if they cannot see the arguments.
- Add a forensic query exercise to every demo: given the trace, answer "what actually ran at sequence 42 and does it match the claimed source commit?" This surfaces whether the report shape is queryable or only declarative.
```
