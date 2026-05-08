# Block 01: Contract Foundation

Status: contract scaffold implemented; product proof blocked until self-trace and self-attestation pass
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Audience: technical executive, CIO, CEO, implementation agents, future external policy consumers

## Purpose

Contract Foundation turns the product promise of `sdp-trace` into machine-checkable contract scaffolding. It defines the first stable set of trace substrate schemas, examples, validation commands, and negative tests before self-trace or pilot evidence begins.

This block does not prove that `sdp-trace` is a viable product. Product viability starts only when the repository traces its own development with these contracts and verifies its own contract release under an immutable source reference and accepted signing or attestation profile.

The block must answer a simple executive question:

> Can an observer validate the contract shapes that will be used to prove `sdp-trace` development itself?

The stronger technical executive question, "Can `sdp-trace` prove itself before proving anyone else?", is not answered by this block. It is the mandatory next block.

## Executive Outcomes

### technical executive

The technical executive gets a data substrate that can show process movement over time: prior/current values, deltas, dimensions, evidence coverage, and `not_assessed` gaps. The substrate does not hide policy choices inside opaque health scores.

### CIO

The CIO gets portable contracts with explicit schema versions, validation commands, redaction rules, and compatibility boundaries. These contracts can be governed across repositories and inherited by external policy consumer.

### CEO

The CEO gets a hard stop: no pilot credibility is claimed until `sdp-trace` traces itself. This block only defines the evidence and accountability contracts that the self-trace proof must use.

## In Scope

Schemas:

- `schema/accountability.schema.json`
- `schema/risk-classification.schema.json`
- `schema/contract-manifest.schema.json`
- `schema/contract-release-verification.schema.json`
- `schema/trusted-identity-policy.schema.json`
- `schema/evidence-event.schema.json`
- `schema/provenance-record.schema.json`
- `schema/observation.schema.json`
- `schema/metric-stream.schema.json`
- `schema/external-verdict-input.schema.json`
- `schema/assessment-input.schema.json`
- updates to `schema/trace.schema.json` only when needed for compatibility with the new entities

Examples:

- one positive assessment input example
- one `not_assessed` example
- one negative example proving native `sdp-trace` artifacts cannot contain forbidden native policy fields
- one negative example proving an AI actor cannot be the sole accountable owner or approver
- one consumer schema-version declaration example for a future external policy consumer consumer
- one contract manifest example with digests for schemas, docs, validation scripts, fixtures, source commit, and compatibility notes
- one contract release verification example for the target signing profile
- one trusted identity policy example naming the authorized release signer identity, protected source ref, workflow identity, required approval refs, and release captain

Validation:

- `jq empty schema/*.json`
- pinned Draft 2020-12 validation through local `ajv@8.20.0` and `scripts/validate-json-schema.mjs`
- dependency installation may use network during setup, but the accepted repository validation command must not require live network access when run in CI or an air-gapped pilot checkout
- one documented automated artifact-safety scan that fails committed examples containing obvious secrets, credentials, raw customer data markers, or private prompt contents
- one documented manifest verification command that recomputes listed artifact digests and fails when the checkout differs from the manifest
- documented target signature verification profile: `sdp-trace-signature/sigstore-dsse-keyless-v1`
- documented signer identity-policy verification; any OIDC signer that does not match the trusted identity policy fails verification
- documented exclusions for local/raw outputs

Documentation:

- schema ownership and versioning
- artifact safety and integrity rules
- `sdp-trace` / external policy consumer inheritance boundary
- consumer schema-version declaration rules
- accountability and escalation rules
- contract release signing and verification rules
- trusted signer identity policy rules
- one-minute technical executive decision narrative in Russian and English, mapped to SpecKit evidence and free of native policy verdict claims

## Out of Scope

- external policy consumer policy implementation
- dashboard or UI
- ingestion daemon
- pilot execution for OpenCode, MiniMax, Kimi, GLM, or Kotlin+Bazel
- compatibility matrix upgrades beyond `not_assessed`
- full TUF implementation, enterprise PKI, private transparency log operation, or long-lived signing-key management

## Scope Boundary: Contract Fixtures vs Pilot Execution

This block defines the contract that later pilot blocks must use. It may create generic fixture shapes and run-card contract skeletons, but it does not execute OpenCode, MiniMax, Kimi, GLM, harness, or JVM+Bazel pilots.

Tasks T027-T033 belong to a later pilot-matrix block. They consume the schemas and examples produced here. They are not complete just because Contract Foundation validates.

## Contract Principles

1. Evidence first: every claim is linked to evidence, provenance, or `not_assessed`.
2. No native policy verdicts: `sdp-trace` does not decide pass/fail, readiness, degradation, override, or evidence sufficiency.
3. External assertions are explicit: external verdicts and evidence-strength claims carry producer, origin, policy reference when available, artifact reference, and provenance.
4. Movement is structural: prior value, current value, delta, unit, dimensions, and evidence coverage are allowed; interpretation labels are external.
5. Safety is required before commit: no secrets, credentials, raw customer data, or private prompt contents in committed examples.
6. Versioning is part of the contract: schemas use Draft 2020-12, stable `$id`, and semver-compatible versioning.
7. Accountability is human: AI actors may produce, review, critique, or judge artifacts, but a human identity or named human-held role owns approval and escalation.
8. Contract integrity is verified: a trusted contract release is a signed manifest plus successful digest verification, not just the current files in a checkout.

## Target Data Flow

```text
evidence event
  -> provenance record
  -> accountable artifact
  -> observation
  -> metric sample / metric stream
  -> trace snapshot
  -> assessment input
  -> external policy consumer such as external policy consumer
```

External verdicts can be recorded as evidence:

```text
external policy decision -> external verdict input -> trace snapshot / assessment input
```

The external verdict path never becomes a native `sdp-trace` decision.

## Accountability Contract

Contract Foundation must make accountability explicit enough for a CEO to ask who owns the next step when the process lies, stalls, or becomes disputed.

Every contract release, assessment input, and evidence package must carry an `accountability` object with:

- `dri`: accountable identity object for a human identity or named human-held role
- `approver`: identity object for the human identity or named human-held role that approved the artifact or release
- `escalation`: identity object or channel object for unresolved failures
- `authority_scope`: idea, spec, plan, task, evidence, assessment_input, contract_release, or external_verdict
- `accountability_claim`: recording_only, content_approval, risk_acceptance, release_approval
- `approval_ref`: inspectable reference such as PR review, release approval, meeting decision, or signed-off record
- `risk_owner`: identity object for the human identity or named role owning residual risk
- `line_of_defense`: first, second, or third

Identity objects must carry `identity_ref` and `actor_type`. Allowed accountable actor types are `human_user`, `human_role`, and `human_group`. AI and system actor types are allowed in provenance, but not as sole accountable actors.

AI identities may appear in provenance as producers, reviewers, critics, or judges. They must not be the sole `dri`, `approver`, `risk_owner`, or escalation owner.

Evidence events may omit a direct accountability object only when they are contained in an evidence package or assessment input that provides effective accountability for every referenced evidence item. An assessment input must not claim completeness or trusted-release readiness if any referenced evidence lacks direct or inherited effective accountability.

If review independence is required by an external policy assertion, the accountability record must preserve enough actor identity and line-of-defense metadata for a policy consumer to detect self-review or same-line review. `sdp-trace` records the facts; external policy consumer decides whether the separation satisfies policy.

For public examples, synthetic human-held roles are allowed. For customer pilots, these fields must map to the customer's accepted identity system or approval process.

## Risk and Oversight Classification

Contract Foundation uses a minimal risk classification aligned with current AI governance practice:

- `observed_autonomy_level`: assistive, collaborative, delegated, or autonomous
- `observed_impact_level`: low, medium, high, or critical
- `classification_source`: human_dri, customer_policy, external_governance_policy, policy_engine, or not_assessed
- `classification_ref`: inspectable reference for the source that supplied the classification
- `declared_oversight`: optional external assertion containing `origin: "external"`, `policy_ref`, `required_oversight`, and `review_independence`

The classification does not decide pass/fail and does not compute oversight obligations inside `sdp-trace`. `sdp-trace` records observed autonomy and impact plus externally declared oversight requirements when supplied. A policy consumer such as external policy consumer decides whether the recorded classification satisfies a policy.

## Contract Release Signing Profile

The target signing profile is `sdp-trace-signature/sigstore-dsse-keyless-v1`.

Contract Foundation signs the contract release, not each schema independently:

```text
contract-manifest.json
  -> SHA-256 digest
  -> in-toto Statement subject
  -> DSSE envelope
  -> Sigstore/Cosign keyless signature bundle
```

The contract manifest includes:

- contract version
- signing profile
- schema `$id`s and versions
- SHA-256 digests for schemas, contract docs, validation scripts, fixtures, and compatibility notes
- source commit
- previous manifest digest when available
- issued timestamp
- exactly one of `valid_until` or explicit `freshness_policy`
- required signer identity policy reference
- accountability object for the release
- approval refs

The trusted identity policy includes:

- expected OIDC issuer
- expected repository or source URI
- expected protected branch, tag, or ref pattern
- expected workflow identity or build system identity
- required release captain identity or role
- required approval refs or CODEOWNERS-style approval evidence
- allowed private-equivalent verifier profile when public Sigstore/Rekor is unavailable

Signature verification must check:

- current manifest digest
- every listed artifact digest
- DSSE envelope binding
- Sigstore certificate identity
- expected OIDC issuer
- expected repository, workflow, and ref policy when available
- transparency log inclusion when the selected environment supports Rekor or an equivalent log
- rollback or stale-contract indicators

Air-gapped or private customer environments may replace public Sigstore/Rekor with an approved equivalent, but the artifact contract remains in-toto statement plus DSSE envelope plus explicit identity policy. The equivalent must specify envelope binding, trusted root or identity source, signing identity, timestamp or freshness evidence, and audit-log or compensating-control status.

A checkout may be schema-valid but not a trusted contract release until manifest digest verification and signature or approved-equivalent verification succeed. Missing, stale, invalid, or `not_assessed` signature verification can be recorded as evidence, but it must not support a `trusted_contract_release` claim.

Block 01 is not product proof even when local release verification evidence exists. A schema-only verification record proves shape; a local private-equivalent signing record proves envelope and digest mechanics. Product trust requires self-trace plus self-attestation against an immutable source reference and accepted trust anchor.

This block does not implement TUF. Rollback protection is limited to single-chain `previous_manifest_digest` continuity plus explicit freshness. Delegated roles, threshold metadata, snapshot metadata, timestamp metadata, target metadata, and key rotation are out of scope and require a future block.

## `not_assessed` Encoding

`not_assessed` is a state, not a magic replacement string. Typed values must not be overwritten with `"not_assessed"`.

Schemas must use one of these patterns:

- measured values: `assessment_state: "not_assessed"`, `value: null`, and a required `not_assessed_reason`
- unavailable optional metadata: omit the typed field and record the missing field in an `unavailable_fields[]` entry with `field`, `state: "not_assessed"`, and `reason`
- partially assessed records: `assessment_state: "partial"` plus explicit `not_assessed` entries for the missing parts

Required fields that establish identity, scope, timestamps, schema version, or artifact boundaries must not silently default to `not_assessed`; missing required fields fail validation.

## Structural Trust Boundary

Contract Foundation proves structural trust only:

- JSON Schema conformance
- stable schema IDs and versions
- required evidence/provenance references
- SHA-256 digest fields where artifacts are referenced
- redaction and integrity status fields
- explicit `not_assessed` reasons
- accountability fields with human-owned approval and escalation
- manifest digest continuity and signature verification status

It does not prove producer honesty, model quality, or policy sufficiency. It does prove whether the current checkout matches a declared contract manifest and whether that manifest was verified under the selected signing profile. Fraud detection, business acceptance of residual risk, and trust-anchor governance are future work or external policy concerns.

## Required Positive Fixture

The positive fixture must prove:

- an evidence event has source, timestamps, actor, status, redaction state, and inspectable artifact metadata
- an evidence event can carry `dedupe_key` and `conflict_refs` without collapsing conflicts into a verdict
- a provenance record captures actor/model/harness/tool metadata with unavailable fields marked `not_assessed`
- an observation references evidence and provenance
- a metric stream compares two windows through prior/current/delta values
- an assessment input packages evidence, observations, metric streams, and `not_assessed` gaps without native policy verdicts
- accountability fields identify human-held DRI, approver, escalation, risk owner, and line of defense
- contract release verification records manifest digest status and signature verification status

## Required `not_assessed` Fixture

The `not_assessed` fixture must prove:

- missing raw logs do not become success
- unavailable model version does not make the whole sample invalid
- missing evidence explains what cannot be assessed and why
- pending evidence remains distinct from completed evidence
- missing signature verification does not become trusted contract release; it is recorded as `not_assessed` or invalid depending on the claim being made
- partially assessed metric streams carry stream-level assessment state so technical executive-facing movement does not hide unassessed samples inside an apparently complete comparison

## Required Negative Fixture

The negative fixture must fail validation when a native `sdp-trace` artifact includes forbidden native policy fields outside `external-verdict-input`.

Forbidden native field names include:

- `verdict`
- `decision`
- `gate_result`
- `gate_verdict`
- `readiness`
- `readiness_verdict`
- `degradation_status`
- `policy_result`
- `policy_threshold`
- `evidence_strength`
- `quality_score`
- `override_result`

Values such as `pass`, `fail`, `ready`, `blocked`, `degrading`, and `improving` are not rejected merely because they appear in prose summaries or external logs. They are rejected when `sdp-trace` exposes them through native policy fields. The same policy values may appear inside `external-verdict-input` with `origin: "external"`.

The negative fixtures must also fail validation or verification when:

- an AI actor is the sole `dri`, `approver`, `risk_owner`, or escalation owner
- an artifact claims `trusted_contract_release` while manifest digest verification fails
- an artifact claims `trusted_contract_release` while signature verification is missing, stale, or invalid
- a contract manifest omits required approval refs for contract release
- a contract manifest omits both `valid_until` and `freshness_policy`
- an OIDC signer identity does not match the trusted identity policy
- a self-review is represented as independent review without distinct actor identity or line-of-defense metadata

## Consumer Schema-Version Declaration

Contract Foundation exports a consumer declaration example for future external policy consumer integration. The example records:

- consumer name
- consumed schema `$id`s
- supported semver ranges
- unsupported or pending schema versions
- last validated fixture refs
- compatibility notes

The real external policy consumer repository owns its own consumer declaration. `sdp-trace` only defines the portable contract shape and example so breaking changes are visible before downstream consumers adopt them.

## Trace Schema Compatibility Decision

Task T012 must run before the first trace snapshot fixture is accepted. `schema/trace.schema.json` remains the compatibility path unless it cannot represent observations, metric samples, evidence refs, external verdict refs, and schema versions without lossy or policy-owning fields.

If the existing trace schema cannot represent those entities cleanly, the block must add a replacement path and migration note before schemas are marked complete. A migration note is required for any change that removes a required field, changes enum semantics, changes ownership boundaries, or makes a prior valid trace snapshot invalid.

## Acceptance Criteria

- AC01: All new schemas are valid JSON and declare Draft 2020-12.
- AC02: Schema docs define `$id`, schema versioning, compatibility, and migration rules.
- AC03: Positive fixture validates with the pinned local `ajv@8.20.0` command without requiring live network access.
- AC04: `not_assessed` fixture validates and preserves missing-evidence reasons.
- AC05: Negative fixture fails validation for native policy verdicts.
- AC06: `schema/trace.schema.json` remains usable or has a documented replacement path and migration note.
- AC07: `schema/README.md` explains `sdp-trace` ownership vs external policy ownership.
- AC08: The documented artifact-safety scan passes and no committed example contains raw secrets, credentials, raw customer data, or private prompt contents.
- AC09: Consumer schema-version declaration example validates and does not make external policy consumer a runtime dependency.
- AC10: Accountability schema and fixtures prove AI actors cannot be sole accountable owners or approvers.
- AC11: Contract manifest schema and example validate and include digests for schemas, docs, validation scripts, fixtures, source commit, approval refs, and compatibility notes.
- AC12: Contract release verification schema and example record the target signing profile, manifest digest status, signature status, signer identity policy, and rollback or freshness status.
- AC13: Manifest verification fails for an intentionally modified artifact that no longer matches its declared digest.
- AC14: Trusted identity policy example validates and rejects a signer identity outside the allowed issuer, repository, workflow, protected ref, or release-captain policy.
- AC15: Metric stream examples expose stream-level `assessment_state` when any sample or comparison is partial or `not_assessed`.
- AC16: Signing proof is not theoretical: completion evidence includes one real release verification result for the selected signing profile or approved private equivalent.
- AC17: Freshness is explicit: manifests without `valid_until` or `freshness_policy` fail validation.
- AC18: technical executive brief documents answer "what is happening and why do I need it?" in under one minute in RU/EN without marketing claims or native `sdp-trace` policy verdicts.

## Traceability

| Epic / Area | Task | Purpose |
|---|---|---|
| `sdp-trace-cdn.8` | T034 | Pin Draft 2020-12, validator, schema IDs, versioning |
| `sdp-trace-cdn.2` | T040 | Artifact safety and integrity rules |
| `sdp-trace-cdn.2` | T044 | Accountability schema and human approval rules |
| `sdp-trace-cdn.2` | T045 | Risk classification and oversight metadata |
| `sdp-trace-cdn.8` | T046 | Contract manifest schema and digest verification |
| `sdp-trace-cdn.8` | T047 | Target signing profile and release verification evidence |
| `sdp-trace-cdn.8` | T048 | Negative AI-accountability and modified-contract-manifest fixtures |
| `sdp-trace-cdn.8` | T049 | Trusted signer identity policy and mismatch fixture |
| `sdp-trace-cdn.8` | T050 | Real release signature verification evidence |
| `sdp-trace-cdn.11` | T051 | One-minute technical executive decision narrative in RU/EN |
| `sdp-trace-cdn.3` | T008 | Observation schema |
| `sdp-trace-cdn.3` | T009 | Metric stream schema |
| `sdp-trace-cdn.3` | T010 | Movement-data example |
| `sdp-trace-cdn.3` | T012 | Trace schema compatibility path |
| `sdp-trace-cdn.4` | T015 | Evidence event schema |
| `sdp-trace-cdn.4` | T016 | Provenance record schema |
| `sdp-trace-cdn.4` | T017 | Assessment input schema |
| `sdp-trace-cdn.4` | T018 | External verdict input schema |
| `sdp-trace-cdn.8` | T035 | Positive and `not_assessed` fixtures |
| `sdp-trace-cdn.8` | T036 | Repository validation command |
| `sdp-trace-cdn.9` | T037 | Pilot evidence package outline consumes, but does not complete, this block's contract |
| `sdp-trace-cdn.4` | T041 | Negative policy-verdict fixture |
| `sdp-trace-cdn.4` | T042 | Schema ownership/versioning docs |
| `sdp-trace-cdn.8` | T043 | Schema-version and compatibility verification |
| `sdp-trace-cdn.6` | T027-T028 | Later pilot run-cards consume this block's contract |
| `sdp-trace-cdn.7` | T029-T031 | Later Kotlin+Bazel pilot fixtures consume this block's contract |
| `sdp-trace-cdn.10` | T032-T033 | Later compatibility matrices consume evidence produced under this block's contract |

## Exit Condition

Contract Foundation is complete when a fresh clone can run the documented syntax and schema validation commands and see:

- positive fixture passes
- `not_assessed` fixture passes
- negative policy-verdict fixture fails for the intended reason
- assessment input contains no native gate/degradation/readiness decision
- schema docs state how external policy consumer inherits supported contract versions
- consumer schema-version declaration example validates
- artifact-safety scan passes for committed examples
- accountability examples validate and AI-as-sole-accountable-owner fails
- contract manifest example validates and digest verification passes
- modified-contract negative fixture fails manifest verification
- release verification example records `sdp-trace-signature/sigstore-dsse-keyless-v1` status without making public Rekor mandatory for private environments
- trusted identity policy example validates and signer mismatch fails
- local release verification evidence exists for contract scaffolding, while product trust remains blocked until self-trace and self-attestation pass
- technical executive brief explains the tool in RU/EN without saying `sdp-trace` owns gate/degradation decisions
