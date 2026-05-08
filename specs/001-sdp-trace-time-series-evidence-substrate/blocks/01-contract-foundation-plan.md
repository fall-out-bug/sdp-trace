# Block 01: Contract Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking. Do not start implementation until the CTO approves this plan.

**Goal:** Build the first machine-checkable `sdp-trace` contracts for evidence, provenance, accountability, contract integrity, observations, metric movement, assessment input, and external verdict inputs without making native policy decisions.

**Architecture:** Keep `sdp-trace` as a portable contract, evidence, accountability, and integrity substrate. Schemas define structure, examples prove positive and negative behavior, scripts prove repeatable validation, and docs explain ownership/versioning/signing so `sdp-gate` can inherit contracts without becoming a dependency.

**Tech Stack:** JSON Schema Draft 2020-12, pinned local `ajv@8.20.0`, `jq`, SHA-256 digest verification, in-toto Statement, DSSE, Sigstore/Cosign keyless signing profile, shell validation scripts, Markdown SpecKit artifacts.

---

## Executive Plan

### CTO

Contract Foundation answers: "Do we have trustworthy movement data, or are we silently making policy calls?" The delivered answer is not a green/red score. It is a validated assessment input containing current/prior metric values, deltas, evidence coverage, provenance, and explicit `not_assessed` gaps.

### CIO

The block creates governed contracts: stable schema IDs, semver rules, offline-capable validation, compatibility notes, artifact redaction/integrity rules, accountability records, signed-release metadata, and a consumer schema-version declaration shape for future `sdp-gate` adoption.

### CEO

The block makes pilot claims defensible and governable. Before we run OpenCode, model, harness, or Kotlin+Bazel pilots, the repository can prove what counts as acceptable evidence, who is accountable for it, whether the contract release is intact, and what must be marked `not_assessed`.

## Architecture

```text
schema layer
  common definitions
  accountability
  risk classification
  contract manifest
  contract release verification
  trusted identity policy
  evidence event
  provenance record
  observation
  metric stream
  trace snapshot compatibility
  assessment input
  external verdict input
  consumer schema-version declaration

fixture layer
  positive assessment package
  not_assessed package
  negative native-policy-field package
  negative AI-as-sole-accountable-owner package
  negative modified-contract-manifest package
  negative unauthorized-signer package
  sdp-gate consumer declaration example
  contract manifest example
  contract release verification example
  trusted identity policy example

validation layer
  JSON syntax check
  Draft 2020-12 schema validation
  expected negative validation failure
  artifact-safety scan
  manifest digest verification
  signing-profile verification evidence

governance layer
  human accountability
  autonomy/impact oversight classification
  schema ownership and versioning
  sdp-trace / sdp-gate boundary
  contract release signing profile
  trace migration notes
```

Data flow remains:

```text
evidence event
  -> provenance record
  -> accountability record
  -> observation
  -> metric stream
  -> trace snapshot
  -> assessment input
  -> external policy consumer
```

External decisions are recorded only as:

```text
external policy decision
  -> external verdict input with origin=external
  -> evidence context
```

## Scope Lock

In scope:

- Schema contracts and examples required by `01-contract-foundation.md`.
- Validation and safety commands that a fresh checkout can run.
- Documentation needed for downstream consumers to inherit contract versions.
- Accountability and risk classification contracts required to identify human DRI, approval, escalation, and oversight.
- Contract manifest and release verification contracts using the target signing profile `sdp-trace-signature/sigstore-dsse-keyless-v1`.

Out of scope:

- Running OpenCode, MiniMax, Kimi, GLM, harness, or Kotlin+Bazel pilots.
- Implementing `sdp-gate` policy evaluation.
- Creating dashboards, ingestion daemons, external signing services, or compatibility claims.
- Operating enterprise PKI, private transparency logs, or full TUF-style repository metadata.
- Declaring any native pass/fail/readiness/degradation verdict in `sdp-trace`.

## Files

Create:

- `package.json`
- `package-lock.json`
- `scripts/validate-contracts.sh`
- `scripts/check-artifact-safety.sh`
- `scripts/verify-contract-manifest.sh`
- `schema/common.schema.json`
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
- `schema/consumer-schema-version-declaration.schema.json`
- `examples/contract-foundation/positive-assessment-input.json`
- `examples/contract-foundation/not-assessed-assessment-input.json`
- `examples/contract-foundation/negative-native-policy-field.json`
- `examples/contract-foundation/negative-ai-sole-accountable-owner.json`
- `examples/contract-foundation/negative-modified-contract-manifest.json`
- `examples/contract-foundation/negative-unauthorized-signer.json`
- `examples/contract-foundation/sdp-gate-consumer-declaration.example.json`
- `examples/contract-foundation/contract-manifest.example.json`
- `examples/contract-foundation/contract-release-verification.example.json`
- `examples/contract-foundation/trusted-identity-policy.example.json`
- `docs/accountability-model.md`
- `docs/contract-release-signing.md`

Modify:

- `schema/README.md`
- `schema/trace.schema.json`
- `specs/001-sdp-trace-time-series-evidence-substrate/contracts/sdp-trace-sdp-gate-boundary.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md` only to mark completed tasks after implementation
- `docs/concepts.md` if boundary wording still implies native policy decisions
- `docs/cto-brief.en.md`
- `docs/cto-brief.ru.md`

Do not modify in this block:

- `archive/research/opencode-model-run-card.md`
- `archive/research/harness-run-card.md`
- `archive/research/kotlin-bazel-fixture-plan.md`
- compatibility matrix rows for real pilots
- `sdp_gate` repository files

## Work Packages

### WP0: Approval Gate

Trace: `sdp-trace-cdn` -> Block 01 plan approval.

- [ ] CTO reviews this plan.
- [ ] Implementation starts only after explicit approval.

Proof:

- Approved plan remains in `specs/001-sdp-trace-time-series-evidence-substrate/blocks/01-contract-foundation-plan.md`.

### WP1: Validator and Safety Foundation

Trace: `sdp-trace-cdn.8` -> T034, T036; `sdp-trace-cdn.2` -> T040.

- [ ] Add `package.json` and `package-lock.json` with `ajv@8.20.0` pinned as a dev validation dependency.
- [ ] Add `scripts/validate-contracts.sh` using local `scripts/validate-json-schema.mjs` and refusing to fetch network dependencies during validation.
- [ ] Add `scripts/check-artifact-safety.sh` with repository exclusions for `.git/`, `.beads/`, `.sdp-trace-runs/`, `benchmarks/repos/`, generated dependencies, editor caches, and temp directories.
- [ ] Document connected setup separately from offline validation: dependency installation may happen during setup, but validation must not call live network.

Proof:

- `rtk jq empty schema/*.json`
- `rtk npm ci`
- `rtk scripts/validate-contracts.sh`
- `rtk scripts/check-artifact-safety.sh`

### WP2: Common Contract Definitions

Trace: `sdp-trace-cdn.8` -> T034, T043.

- [ ] Create `schema/common.schema.json` with reusable definitions for schema version, IDs, timestamps, artifact refs, SHA-256 digests, redaction status, integrity status, provenance refs, evidence refs, unavailable fields, and `not_assessed` state.
- [ ] Use JSON Schema Draft 2020-12 `$schema` and stable `$id`.
- [ ] Encode `not_assessed` as state plus reason; do not use `"not_assessed"` as a typed scalar sentinel.

Proof:

- `rtk jq empty schema/common.schema.json`
- `rtk scripts/validate-contracts.sh`

### WP2A: Accountability and Risk Classification

Trace: `sdp-trace-cdn.2` -> T044, T045.

- [ ] Create `schema/accountability.schema.json` with `dri`, `approver`, `escalation`, `authority_scope`, `approval_ref`, `risk_owner`, and `line_of_defense`.
- [ ] Encode accountable identities as objects with `identity_ref` and `actor_type`; allowed accountable actor types are `human_user`, `human_role`, and `human_group`.
- [ ] Encode the rule that AI actors may be producers/reviewers/critics/judges in provenance but cannot be the sole `dri`, `approver`, `risk_owner`, or escalation owner.
- [ ] Require assessment inputs to provide effective accountability for every referenced evidence item, either directly or through a containing evidence package.
- [ ] Create `schema/risk-classification.schema.json` with `observed_autonomy_level`, `observed_impact_level`, `classification_source`, `classification_ref`, and optional external `declared_oversight`.
- [ ] Keep `required_oversight` and `review_independence` only inside external `declared_oversight`; `sdp-trace` must not derive them from observed autonomy/impact.
- [ ] Add `examples/contract-foundation/negative-ai-sole-accountable-owner.json` proving the forbidden accountability shape fails.
- [ ] Document the operating model in `docs/accountability-model.md`: first line owns execution/evidence, second line owns controls/policy challenge, third line owns independent assurance.

Proof:

- Accountability examples validate for human-held roles.
- AI-as-sole-accountable-owner fixture fails for the intended reason.

### WP2B: Contract Manifest and Release Signing

Trace: `sdp-trace-cdn.8` -> T046, T047, T048, T049, T050.

- [ ] Create `schema/contract-manifest.schema.json` for contract version, signing profile, schema versions, artifact digests, source commit, previous manifest digest, issued timestamp, exactly one of `valid_until` or `freshness_policy`, identity policy ref, approval refs, compatibility notes, and accountability.
- [ ] Create `schema/contract-release-verification.schema.json` for manifest digest status, artifact digest status, signing profile, signature status, signer identity, OIDC issuer, identity policy ref, transparency log status, freshness status, provenance refs, and accountability.
- [ ] Create `schema/trusted-identity-policy.schema.json` for expected OIDC issuer, source URI, protected ref, workflow identity, release captain, required approval refs, and allowed private-equivalent verifier profile.
- [ ] Add `scripts/verify-contract-manifest.sh` to recompute listed SHA-256 digests and fail if any checked-out artifact differs from the manifest.
- [ ] Add `docs/contract-release-signing.md` documenting target profile `sdp-trace-signature/sigstore-dsse-keyless-v1`: in-toto Statement subject, DSSE envelope, Sigstore/Cosign keyless signature, OIDC identity, workflow/ref policy, and public or private transparency-log handling.
- [ ] Add `examples/contract-foundation/contract-manifest.example.json`, `examples/contract-foundation/contract-release-verification.example.json`, and `examples/contract-foundation/trusted-identity-policy.example.json`.
- [ ] Add `examples/contract-foundation/negative-modified-contract-manifest.json` or equivalent fixture proving manifest digest mismatch is not trusted.
- [ ] Add `examples/contract-foundation/negative-unauthorized-signer.json` or equivalent fixture proving signer identity mismatch is not trusted.
- [ ] Produce one contract release verification evidence record for the selected signing profile shape or approved private equivalent before claiming contract scaffolding complete. This does not establish product trust without self-trace and self-attestation.

Proof:

- Contract manifest example validates.
- Contract release verification example validates.
- Trusted identity policy example validates.
- Manifest verification passes for the positive manifest and fails for the modified-contract negative fixture.
- Unauthorized signer fixture fails for the intended identity-policy reason.
- A real release verification evidence record exists; otherwise the block remains schema-complete but not trusted-release-complete.

### WP3: Evidence and Provenance Schemas

Trace: `sdp-trace-cdn.4` -> T015, T016.

- [ ] Create `schema/evidence-event.schema.json` with source, timestamps, actor, event type, status, artifact metadata, redaction state, integrity state, `dedupe_key`, `conflict_refs`, and external assertions.
- [ ] Create `schema/provenance-record.schema.json` with actor/model/harness/tool metadata, unavailable fields, payload digest, digest algorithm, optional `hash_prev`, and `chain_scope`.
- [ ] Reference accountability records where human sign-off or evidence ownership is required.
- [ ] Ensure neither schema assigns evidence strength as a native `sdp-trace` conclusion.

Proof:

- Positive fixture validates evidence and provenance records.
- `not_assessed` fixture validates unavailable model/tool fields through `unavailable_fields[]`.

### WP4: Observation and Movement Schemas

Trace: `sdp-trace-cdn.3` -> T008, T009, T010.

- [ ] Create `schema/observation.schema.json` for evidence-backed statements with evidence refs, provenance refs, scope, observed timestamp, and assessment status.
- [ ] Create `schema/metric-stream.schema.json` for samples and comparisons: previous value, current value, delta, unit, dimensions, evidence coverage, and missing evidence reasons.
- [ ] Require stream-level `assessment_state` when any sample or comparison is partial or `not_assessed`.
- [ ] Add fixture content showing movement over two windows without `degrading`, `improving`, `pass`, `fail`, `ready`, or `blocked` as native labels.

Proof:

- `examples/contract-foundation/positive-assessment-input.json` validates.
- The validated metric stream contains structural movement data only.

### WP5: Assessment and External Verdict Schemas

Trace: `sdp-trace-cdn.4` -> T017, T018, T041.

- [ ] Create `schema/external-verdict-input.schema.json` with `origin: "external"`, producer, producer type, verdict kind, source value, optional policy reference, artifact metadata, redaction status, and provenance refs.
- [ ] Create `schema/assessment-input.schema.json` packaging trace ref, evidence bundle refs, metric stream refs, observations, external verdict inputs, and `not_assessed` entries.
- [ ] Require accountability and risk classification refs in assessment inputs that claim readiness for policy-engine handoff.
- [ ] Reject forbidden native policy fields outside external verdict input records: `verdict`, `decision`, `gate_result`, `gate_verdict`, `readiness`, `readiness_verdict`, `degradation_status`, `policy_result`, `policy_threshold`, `evidence_strength`, `quality_score`, and `override_result`.
- [ ] Add `examples/contract-foundation/negative-native-policy-field.json` that fails for the intended forbidden-field reason.

Proof:

- Positive and `not_assessed` examples validate.
- Negative example fails validation for the documented forbidden native field.

### WP6: Trace Compatibility and Consumer Declaration

Trace: `sdp-trace-cdn.3` -> T012; `sdp-trace-cdn.8` -> T043.

- [ ] Update `schema/trace.schema.json` to represent `observation`, `metric_sample`, `metric_stream`, `external_verdict`, and schema-version metadata if the current trace schema can support them cleanly.
- [ ] If the current trace schema cannot support the new entities without lossy or policy-owning fields, document a replacement path and migration note before completing the block.
- [ ] Create `schema/consumer-schema-version-declaration.schema.json`.
- [ ] Add `examples/contract-foundation/sdp-gate-consumer-declaration.example.json` as a portable example, not a runtime `sdp-gate` dependency.

Proof:

- Current `examples/github-speckit/trace.json` still validates or has a documented migration path.
- Consumer declaration example validates.

### WP7: Documentation and Boundary Alignment

Trace: `sdp-trace-cdn.2` -> T005, T040; `sdp-trace-cdn.4` -> T042; `sdp-trace-cdn.11` -> T006, T007, T051 if touched.

- [ ] Update `schema/README.md` with schema IDs, versioning, migration rules, validator commands, offline validation rule, `not_assessed` encoding, and native vs external verdict ownership.
- [ ] Update `contracts/sdp-trace-sdp-gate-boundary.md` if implementation details require clearer contract inheritance language.
- [ ] Audit any touched docs for wording that implies `sdp-trace` owns policy thresholds or pass/fail decisions.
- [ ] Classify existing `gate-verdict` artifacts as compatibility/external decision records unless a later approved block replaces them.
- [ ] Rewrite `docs/cto-brief.en.md` and `docs/cto-brief.ru.md` as one-minute CTO decision narratives: problem, why it matters, what `sdp-trace` records, what it refuses to decide, where `sdp-gate` starts, and how Block 01 proves the foundation.
- [ ] Ensure every CTO narrative claim maps to a SpecKit artifact or avoids the claim.

Proof:

- No touched doc claims `sdp-trace` decides pass/fail, readiness, or degradation.
- Schema docs explain how a future `sdp-gate` consumer declares supported versions.
- CTO brief can be read in under one minute and answers "what is happening and why do I need it?" without marketing claims.

### WP8: Local Verification Evidence

Trace: `sdp-trace-cdn.8` -> T035, T036, T043.

- [ ] Run JSON syntax validation.
- [ ] Run positive fixture validation.
- [ ] Run `not_assessed` fixture validation.
- [ ] Run negative fixture validation and record the expected failure reason.
- [ ] Run negative AI-accountability fixture validation and record the expected failure reason.
- [ ] Run contract manifest digest verification.
- [ ] Run negative modified-manifest verification and record the expected failure reason.
- [ ] Run trusted identity policy validation and unauthorized signer negative verification.
- [ ] Record real release verification evidence for the selected signing profile or approved private equivalent.
- [ ] Run artifact-safety scan.
- [ ] Run `git diff --check`.

Proof:

- Commands and results are summarized in the implementation completion note.
- No pilot compatibility claim is made from this block alone.

### WP9: Post-Implementation Socratic Review

Trace: `sdp-trace-cdn` -> implementation proof review.

- [ ] Run clean-context PI critic on implemented artifacts with a provider different from the judge used here.
- [ ] Critic prompt asks only questions about contract gaps, policy leakage, portability, security, schema compatibility, and proof quality; it must not propose solutions.
- [ ] Author fixes implementation or docs if critic finds unresolved blocking/major issues.
- [ ] Run clean-context PI judge on a different provider from the critic.
- [ ] Store the review artifacts under `specs/001-sdp-trace-time-series-evidence-substrate/blocks/`.

Proof:

- PI judge returns PASS before the block is called implemented.

## Epic to Task Trace

| Executive Outcome | Epic / Beads Area | SpecKit Task | Implementation Output | Proof |
|---|---|---|---|---|
| Governed contracts | `sdp-trace-cdn.8` | T034 | `schema/README.md`, `package.json`, `package-lock.json` | Draft 2020-12 and validator version documented |
| Artifact safety | `sdp-trace-cdn.2` | T040 | `scripts/check-artifact-safety.sh`, docs | Safety scan passes |
| Human accountability | `sdp-trace-cdn.2` | T044 | `schema/accountability.schema.json`, `docs/accountability-model.md` | AI-as-sole-accountable-owner rejected |
| Risk-based oversight | `sdp-trace-cdn.2` | T045 | `schema/risk-classification.schema.json` | Autonomy/impact/oversight recorded without pass/fail |
| Contract integrity | `sdp-trace-cdn.8` | T046, T048 | `schema/contract-manifest.schema.json`, `scripts/verify-contract-manifest.sh` | Modified artifact digest mismatch fails |
| Contract release signing | `sdp-trace-cdn.8` | T047, T049, T050 | `schema/contract-release-verification.schema.json`, `schema/trusted-identity-policy.schema.json`, `docs/contract-release-signing.md` | Target profile, identity policy, and real verification status recorded |
| Evidence substrate | `sdp-trace-cdn.4` | T015 | `schema/evidence-event.schema.json` | Positive fixture validates |
| Provenance substrate | `sdp-trace-cdn.4` | T016 | `schema/provenance-record.schema.json` | Unavailable model/tool fields validate as `not_assessed` |
| CTO movement data | `sdp-trace-cdn.3` | T008, T009, T010 | observation and metric-stream schemas plus positive fixture | Prior/current/delta data validates without native verdicts |
| Trace compatibility | `sdp-trace-cdn.3` | T012 | updated trace schema or migration note | Existing trace example remains valid or replacement path exists |
| Policy boundary | `sdp-trace-cdn.4` | T017, T018, T041 | assessment-input and external-verdict schemas plus negative fixture | Native policy field rejected; external verdict accepted |
| Consumer inheritance | `sdp-trace-cdn.8` | T043 | consumer declaration schema and example | Example validates without `sdp-gate` runtime dependency |
| Boundary docs | `sdp-trace-cdn.2`, `sdp-trace-cdn.4` | T005, T042 | boundary contract and schema docs | Docs state ownership and migration rules |
| CTO narrative | `sdp-trace-cdn.11` | T051 | `docs/cto-brief.en.md`, `docs/cto-brief.ru.md` | One-minute RU/EN decision narrative with no native verdict claims and explicit self-proof blocker |
| Self-proof blocker | `sdp-trace-cdn.12` | T020-T026, T052-T061 | no implementation in Block 01 | Tasks remain future/pending and block customer pilot claims |
| Later pilots consume contract | `sdp-trace-cdn.6`, `.7`, `.10` | T027-T033 | no implementation in Block 01 | Tasks remain future/pending |

## Self-Trace Alignment

We should not self-trace implementation work before these contracts exist, because that would produce evidence in an unstable format. That bootstrapping choice is acceptable only if self-trace becomes the immediate next product proof.

The first real self-trace block starts immediately after Contract Foundation scaffolding validates. It consumes these schemas to record this block's own plan, validation commands, PI review, crisis review, and implementation evidence as a partial assessment input.

Until that self-trace validates, Block 01 is not evidence that the product works. It is evidence that the contract scaffolding is ready to be tested on the repository itself.

## Completion Criteria

Block 01 contract scaffolding is complete only when:

- Positive assessment input validates.
- `not_assessed` assessment input validates.
- Negative native-policy-field fixture fails for the intended reason.
- Negative AI-as-sole-accountable-owner fixture fails for the intended reason.
- Negative modified-contract-manifest fixture fails for the intended reason.
- Consumer schema-version declaration validates.
- Contract manifest example validates and digest verification passes.
- Contract release verification example records the target signing profile and verification status.
- Trusted identity policy example validates and unauthorized signer fixture fails.
- One real release verification evidence record exists for the selected profile or approved private equivalent.
- Artifact-safety scan passes.
- Trace schema compatibility is preserved or a migration note is committed.
- PI post-implementation judge returns PASS.
- No native `sdp-trace` artifact claims pass/fail/readiness/degradation.

Block 01 does not complete product proof. Product proof requires Phase 5 self-trace and Phase 5A self-attestation.
