# Block 18: Redaction, Retention, And Forensic Profiles

Status: spec delta and implementation plan; awaiting explicit approval before
implementation.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/09-flight-recorder-trust-kernel.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-implementation-plan.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/16-protected-gate-enforcement-profile.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/17-managed-harness-enforcement-profile.md`
- `docs/flight-recorder.md`

## Goal

Make redaction and retention verifier-significant instead of decorative. A
reviewer should be able to tell whether a selected run is safe to inspect,
whether critical forensic evidence is reconstructable, and whether redaction
gaps cap or fail the selected profile.

The product outcome is a profile family that keeps the default path safe for
secrets while giving incident, compliance, and platform owners a stricter
forensic mode when they deliberately retain sanitized, encrypted, or external raw
evidence.

## Problem

Block 09 introduced retention modes and redaction states, but the current
implementation is still too shallow for high-trust use:

- digest-only evidence can prove that something existed, but often cannot
  support incident reconstruction;
- unresolved redaction already appears as a fixture, but redaction policy,
  authority, and profile impact are not first-class enough;
- protected and managed profiles avoid leaking raw output, but they do not yet
  express how redaction and retention cap downstream forensic claims;
- default safe retention is useful for DX, but it can be mistaken for a
  forensic record if the profile boundary is not explicit;
- encrypted or external raw references need digest binding and key/access state,
  otherwise they become unverifiable prose.

The weak framing would be "store more logs." That is wrong. The correct framing
is "prove what was retained, prove what was redacted before persistence, and
fail or cap forensic claims when reconstruction evidence is insufficient."

## Non-Goals

- No default raw stdout, stderr, prompt, source snippet, model response,
  credential, token, adapter secret, gateway token, OIDC token, or checkpoint
  key material persistence.
- No dependency on a named KMS, vault, SIEM, cloud object store, transparency
  log, or customer audit platform.
- No organization policy decision about acceptable privacy risk, incident
  severity, legal hold, or retention duration.
- No native merge, release, readiness, degradation, override approval, or
  risk-acceptance decision.
- No backfilling forensic proof for old runs that only retained digests.
- No arbitrary free-form evidence labels outside the current schema and claim
  tag rules.

## Product Boundary

Block 18 may emit verifier-derived redaction and retention facts:

- selected retention/redaction profile;
- effective retention mode per event and artifact;
- redaction policy digest and redaction rule identifiers applied before
  persistence;
- redaction authority and unresolved-redaction reason codes;
- raw payload digest, redacted payload digest, and sanitized excerpt digest
  where applicable;
- encrypted or external raw reference metadata, digest binding, access state, and
  key custody state without storing the key or raw payload;
- forensic completeness state for critical event families;
- deterministic condition rows explaining whether a profile passes, fails,
  caps to lower trust, or cannot verify.

Block 18 must not decide whether an organization accepts the residual risk.
External consumers may use these facts for legal, compliance, incident, or gate
policy.

## Profile Model

Block 18 separates recorder retention modes, recording policy profiles, and
verifier assessment profiles. Retention modes are embedded in run artifacts and
use the Block 09 / FR-054 identifiers exactly. Recording policy profiles map
event families to those modes. The assessment profile evaluates retained facts
against a selected redaction policy.

Retention modes:

- `digest_only`: only hash and metadata retained; safe by default but low
  forensic value.
- `sanitized_excerpt`: selected safe excerpts retained plus full raw digest.
- `encrypted_raw_ref`: raw artifact encrypted or sealed outside the committed
  run artifact; key held outside the run artifact.
- `external_artifact_ref`: raw artifact stored in an external artifact system,
  SIEM, CI artifact, customer audit log, or equivalent store with digest and
  access metadata.
- `not_assessed`: evidence unavailable; reason and next required evidence are
  mandatory.

Recording policy profiles:

- `safe_default`: pre-write redaction and digest-first retention. It maps raw
  payload classes to `digest_only` or `not_assessed` and does not persist raw
  payloads by default.
- `reviewable_sanitized`: maps declared event families to
  `sanitized_excerpt` where policy permits safe excerpts.
- `encrypted_raw_reference`: maps declared event families to
  `encrypted_raw_ref` with key custody and access state.
- `external_forensic_reference`: maps declared event families to
  `external_artifact_ref` with stable reference, digest, retention lifecycle,
  and access state.

Assessment profile:

- `forensic_retention` (`--profile forensic-retention`): selected verifier
  profile that requires critical event families to use `sanitized_excerpt`,
  `encrypted_raw_ref`, or `external_artifact_ref` according to policy.
  Digest-only critical evidence fails unless the policy explicitly classifies
  that event family as non-critical.

The default CLI behavior remains safe. Forensic profiles must be selected
explicitly and must explain which evidence cannot be reconstructed.

Selecting `forensic_retention` must be recorded as a trace event with actor
identity, selected profile, selected redaction policy digest, and declared
justification. Missing profile-selection trace evidence is `not_assessed` for
profile-selection accountability and `cannot_verify` when the selected policy
requires accountable profile selection.

## Redaction Policy Contract

Block 18 should add a portable redaction policy contract with:

- policy id, schema version, policy digest, and provenance reference;
- rule ids, detector family, and rule version;
- allowed retention modes: `digest_only`, `sanitized_excerpt`,
  `encrypted_raw_ref`, `external_artifact_ref`, `not_assessed`;
- redaction actions: `apply_rule`, `withhold`, `mark_unavailable`;
- forbidden committed-artifact persistence classes: credentials, tokens, raw
  prompts, raw model responses, source snippets, stdout/stderr bodies, OIDC
  tokens, adapter secrets, gateway tokens, and checkpoint key material;
- authority identity for policy selection and emergency withholding, expressed
  as a provenance actor or accountability reference rather than a prose string;
- profile mapping from event families to required retention modes;
- explicit critical event family classifications;
- unresolved-redaction handling and profile impact.

Policy absence is acceptable for `safe_default` only if the built-in policy id
and digest are emitted. For stricter profiles, missing policy is
`cannot_verify`.

`withhold` is a redaction action, not a retention mode. A withheld payload maps
to retention mode `not_assessed` unless the policy also records an accepted
`encrypted_raw_ref` or `external_artifact_ref` for the raw evidence.

Forbidden persistence classes apply to committed run artifacts. Encrypted or
external raw references are opt-in exceptions permitted only for explicitly
selected forensic retention policy paths and must not place raw payloads in the
committed artifact.

The built-in safe default policy must be versioned, for example
`builtin-safe-default-v1`, and evidence must record the version and digest used.
New `sdp-trace` versions must preserve verification of older built-in policy
versions or introduce an explicit schema/policy migration note. A new built-in
policy digest must not silently reinterpret old evidence.

Self-asserted redaction authority is insufficient for `forensic_retention`.
When authority identity cannot be verified against the selected provenance or
accountability reference, the result is `cannot_verify`. Withholding actions
must record withholding authority, requestor identity when different, reason
code, and declared justification. Anonymous withholding cannot satisfy forensic
retention.

## Event And Artifact Delta

Flight-recorder events and run manifests should expose:

- `redaction_policy_ref`;
- `redaction_rule_refs`;
- `redaction_authority`;
- `redaction_input_digest`;
- `redacted_payload_digest`;
- `retention_mode`;
- `retention_lifecycle`;
- `forensic_importance`;
- `raw_reference` with reference type, digest, access state, key custody state,
  retention expiry or lifecycle, and unavailable reason when applicable.

The verifier must not infer raw evidence availability from a URL, file path, or
prose note. Digest binding and access state must be machine-readable.

`redaction_input_digest` and `redacted_payload_digest` are persisted digests.
The verifier must not need raw content to verify that the event row records both
the original payload commitment and the retained redacted payload commitment.
`redaction_input_digest` means the pre-redaction payload digest.

Recorder-emitted fields:

| Field | Emitted by recorder | Used by verifier |
| --- | --- | --- |
| `redaction_policy_ref` | yes | policy digest and version binding |
| `redaction_rule_refs` | yes | rule coverage and unresolved-redaction checks |
| `redaction_authority` | yes | authority verification |
| `redaction_input_digest` | yes | pre-redaction payload commitment |
| `redacted_payload_digest` | yes | retained redacted payload commitment |
| `retention_mode` | yes | accepted evidence mode per event/artifact |
| `retention_lifecycle` | yes | raw-reference lifecycle checks |
| `forensic_importance` | yes | criticality evaluation |
| `raw_reference` | yes when applicable | digest, access, and key custody binding |

Verifier-computed fields stay in assessment condition rows, including
`critical_evidence_reconstructable`, `raw_reference_bound`,
`capped_to_retention_mode`, and final profile state.

For forensic profiles, digest algorithms must be SHA-256 or a stronger
approved algorithm declared by schema. Weak, legacy, unknown, or truncated
digests are `cannot_verify` for forensic retention.

Default critical event families for `forensic_retention` are:

- `command_started`;
- `command_finished`;
- `test_output_observed`;
- `file_mutation_observed`;
- `artifact_captured`;
- `model_identity_observed`;
- `harness_identity_observed`;
- `requirement_superseded`;
- `redaction_applied`;
- `run_closed`.

The selected redaction policy may classify additional event families as
critical. It may classify a default critical family as non-critical only with a
machine-readable reason and authority reference; that downgrade must be visible
in condition rows.

`raw_reference` schema shape:

```json
{
  "reference_type": "encrypted_raw_ref | external_artifact_ref",
  "reference_uri": "string",
  "digest": {
    "algorithm": "sha256",
    "value": "hex-encoded digest"
  },
  "access_state": "verified_available | restricted | unavailable | revoked | not_assessed",
  "access_state_last_verified": "RFC3339 timestamp or not_assessed",
  "key_custody_state": "not_applicable | holder_known | escrowed | destroyed | compromised | unknown | not_assessed",
  "retention_lifecycle": {
    "state": "active | expired | revoked | deleted | rotated | not_assessed",
    "policy_ref": "string or not_assessed",
    "expires_at": "RFC3339 timestamp or not_assessed"
  },
  "unavailable_reason": "missing_reference | access_denied | expired | key_unavailable | store_unreachable | digest_unverifiable | not_assessed"
}
```

External references with prose-only digest, URL-only evidence, missing access
state, missing required key custody state, or unverifiable access checks cannot
pass `forensic_retention`. Access state is an observation at assessment time.
Post-assessment revocation, key compromise, or access withdrawal must be
recorded as a superseding event, not by mutating closed evidence.

## Verifier Semantics

Required condition groups:

| Condition | Required behavior |
| --- | --- |
| `redaction_policy_bound` | `pass` when built-in or supplied policy digest is bound to the run and event rows; `cannot_verify` when required policy is missing; `fail` when event rows contradict the selected policy digest. |
| `redaction_prewrite_applied` | `pass` when persisted safe artifacts contain only allowed retained values and digests; `fail` when secret-like values are persisted; `cannot_verify` when the verifier cannot inspect required retained metadata. |
| `redaction_unresolved_visible` | `pass` when unresolved, unverifiable, or withheld redaction is visible with reason codes and impact; `fail` when unresolved redaction is hidden or contradicted; `cannot_verify` when reason evidence is inaccessible. |
| `retention_mode_declared` | `pass` when each relevant event/artifact declares a valid FR-054 retention mode; `fail` for invalid or contradictory modes; `cannot_verify` for missing required mode metadata. |
| `critical_evidence_reconstructable` | `pass` when critical event families have accepted retention for the selected profile; `fail` with `capped_to_retention_mode` when critical evidence is digest-only or unavailable; `cannot_verify` when raw-reference access cannot be verified. |
| `raw_reference_bound` | `pass` when encrypted or external raw references include digest binding and access/key custody state; `fail` for digest mismatch or contradictory metadata; `cannot_verify` for inaccessible stores, missing access verification, or unavailable keys. |
| `forensic_profile_not_overclaimed` | `pass` when forensic output does not exceed retained evidence; `fail` when digest-only or unavailable critical evidence is claimed as reconstructable; `cannot_verify` when profile-selection or policy evidence is inaccessible. |
| `profile_selection_accountable` | `pass` when forensic profile selection is recorded with actor, policy digest, and justification where policy requires it; `not_assessed` when accountability is out of scope for safe default; `cannot_verify` when required actor or policy evidence is missing. |

Exit behavior follows existing profile conventions:

- `0`: selected profile passes.
- `1`: verifier ran and required redaction/retention evidence failed or
  contradicted the policy.
- `2`: invalid invocation or unsupported profile.
- `3`: verifier could not verify required policy, raw reference, access state,
  or artifact binding.

When `forensic_retention` is selected and critical evidence is insufficient,
the assessment fails with a deterministic condition row and
`capped_to_retention_mode`, such as `digest_only` or `sanitized_excerpt`.
This cap is explanatory verifier output, not a separate pass state and not an
upgrade path. A later upgrade requires a new assessment over evidence that was
actually retained before or during the original run.

## CLI And UX

Preferred command surface:

```text
sdp-trace assess --profile forensic-retention ...
sdp-trace assess preview --profile forensic-retention ...
sdp-trace assess explain <assessment-result.json>
```

Preview must show what would be retained and redacted, never raw secret-like
values. Explain must list missing reconstruction evidence and next actions in
plain terms. Output must stay deterministic and safe for logs.

Preview is read-only and must not claim runtime equivalence unless it actually
executes the same redaction engine and policy resolver as assessment. The
current preview surface may report input status and safe policy/run metadata,
including rule ids, detector classes, action names, retention modes, and counts
or digests, not matched values. Preview over an existing run artifact evaluates
retained metadata. A future dry-run mode over a command descriptor may simulate
policy behavior, but it must label itself as simulation and must not claim
future proof.

The UX must make the trade-off visible:

- safe defaults reduce leakage but cap forensic claims;
- forensic reconstruction requires explicit retention choices and external key
  or artifact governance;
- withholding may be correct for privacy, but it keeps affected proof
  `not_assessed` or `cannot_verify`.

## Test And Fixture Expectations

Implementation requires tests and committed fixtures for:

- default safe profile redacts secret-like argv/stdout/stderr markers before
  persistence;
- preview output does not leak raw command args, stdout/stderr bodies, prompts,
  source snippets, credentials, tokens, model responses, adapter secrets,
  gateway tokens, or key material;
- digest-only critical evidence fails `forensic_retention`;
- sanitized excerpts pass only for event families allowed by policy;
- encrypted raw references pass only when digest binding, key custody state,
  and access state are present;
- external raw references pass only when digest, stable reference, retention
  lifecycle, and access state are present;
- external raw reference with present but unverifiable access state produces
  `cannot_verify` with reason code `access_unverifiable`;
- unresolved redaction fails or caps profiles according to policy;
- withhold action maps to `not_assessed` unless accepted encrypted or external
  raw reference evidence exists;
- missing redaction policy is `cannot_verify` outside built-in safe default;
- malformed or contradictory raw references fail.

## Review Plan

Spec review must use separate planes:

- requirements-vs-product-boundary review for privacy, forensic, and
  no-policy-decision scope;
- tracing/evidence review for redaction policy, retention lifecycle, and raw
  reference binding;
- code/correctness review after implementation for verifier behavior, safety
  tests, schemas, fixtures, and CLI output.

PR-level review repeats the same planes. Absent CI remains `not_assessed`.

Review disposition must be recorded in
`specs/001-sdp-trace-time-series-evidence-substrate/blocks/18-redaction-retention-forensic-profiles-review-ledger.md`
with finding id, severity, reviewer plane, finding, disposition, and evidence
reference. Spec-level findings must be closed or explicitly blocked before
implementation approval.

## Acceptance Criteria

- AC1: Spec, tasks, and schema deltas define explicit retention/redaction
  profiles without enabling raw persistence by default.
- AC2: Built-in safe default remains safe for committed examples and ordinary
  local use.
- AC3: Forensic profile cannot pass with digest-only critical evidence unless
  policy explicitly classifies the event family as non-critical.
- AC4: Encrypted and external raw references are verifier-visible and
  digest-bound, with access/key custody states recorded or `cannot_verify`.
- AC5: Redaction failures and unresolved redactions produce deterministic
  condition rows and never become pass through prose.
- AC6: CLI preview and explain surfaces are deterministic, useful to humans,
  and safety-sensitive negative tests prove secret-like values do not leak.
- AC7: Output remains verifier facts only; no native merge, release,
  readiness, degradation, legal, incident, or risk-acceptance decision is
  introduced.

## Implementation Slices

1. Schema and contract slice: redaction policy schema, retention/raw-reference
   fields, assessment-result schema version, fixture contracts, and schema
   parse checks.
2. Safe default slice: built-in policy id/digest, pre-write redaction metadata,
   and preview safety behavior.
3. Forensic verifier slice: profile selection, condition rows, critical event
   family evaluation, and deterministic exit codes.
4. Raw reference slice: encrypted/external reference binding, access/key
   custody states, and negative fixtures.
5. Explain and documentation slice: human-readable explanation, docs, review
   ledger, and evidence notes.

Block 18 must extend or version `schema/assessment-result.schema.json` for
forensic retention facts while preserving `assess explain` read compatibility
with Block 17 assessment results and Block 14/16 gate-result artifacts.

Implementation does not start until this spec direction is explicitly approved.
