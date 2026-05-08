# Data Model: sdp-trace Time-Series Evidence Substrate

## Entity: Accountability Record

Represents human accountability for an artifact, release, assessment package, or external decision record.

Fields:

- `dri`: accountable identity object
- `approver`: approving identity object
- `escalation`: accountable identity object or escalation channel object
- `authority_scope`: idea, spec, plan, task, evidence, assessment_input, contract_release, external_verdict, custom
- `accountability_claim`: recording_only, content_approval, risk_acceptance, release_approval
- `approval_ref`: inspectable approval reference such as PR review, release approval, meeting decision, or signed-off record
- `risk_owner`: accountable identity object that owns residual risk
- `line_of_defense`: first, second, third

Rules:

- Accountable identity objects contain `identity_ref` and `actor_type`.
- Accountable actor types are human_user, human_role, or human_group.
- AI actors may appear in provenance as producers, reviewers, critics, or judges.
- AI actors must not be the sole `dri`, `approver`, `risk_owner`, or escalation owner.
- Assessment inputs must provide effective accountability for every referenced evidence item, either directly on the evidence event or inherited from an evidence package.
- External verdict accountability must distinguish recording accountability from correctness or risk-acceptance accountability through `accountability_claim`.
- Public examples may use synthetic human-held roles. Customer pilots must map these fields to an accepted identity or approval system.
- Accountability records do not decide policy outcomes; they identify who owns follow-up and escalation.

## Entity: Risk Classification

Represents the governance context for AI-assisted work.

Fields:

- `observed_autonomy_level`: assistive, collaborative, delegated, autonomous
- `observed_impact_level`: low, medium, high, critical
- `classification_source`: human_dri, customer_policy, external_governance_policy, policy_engine, not_assessed
- `classification_ref`
- `declared_oversight`: optional external assertion with `origin: external`, `policy_ref`, `required_oversight`, and `review_independence`
- `risk_notes`: optional concise context

Rules:

- Risk classification records observed autonomy and impact plus externally declared oversight obligations. It does not decide pass/fail.
- `sdp-trace` must not derive `required_oversight` from the observed classification. That derivation belongs to a policy consumer or external governance source.
- Downstream policy engines such as external policy consumer may use this classification to enforce review or approval policies.
- Higher autonomy and higher impact should preserve more evidence and stronger human approval references.

## Entity: Contract Manifest

Represents a versioned contract release candidate.

Fields:

- `contract_version`
- `signing_profile`: target value `sdp-trace-signature/sigstore-dsse-keyless-v1`
- `schema_versions`
- `artifacts`: path, kind, SHA-256 digest, schema `$id` when applicable
- `source_commit`
- `previous_manifest_digest`
- `issued_at`
- exactly one of `valid_until` or `freshness_policy`
- `trusted_identity_policy_ref`
- `approval_refs`
- `compatibility_notes`
- `accountability`

Rules:

- The manifest signs the contract release as a whole, not each schema independently.
- Listed artifacts include schemas, contract docs, validation scripts, fixtures, source commit metadata, and compatibility notes.
- A checkout can be schema-valid but not a trusted contract release if manifest digest verification fails or signature verification is missing for a trusted-release claim.
- A manifest with neither `valid_until` nor `freshness_policy` is invalid.
- `previous_manifest_digest` supports rollback detection but is not a full TUF implementation.

## Entity: Trusted Identity Policy

Represents who is allowed to issue a trusted contract release.

Fields:

- `id`
- `oidc_issuer`
- `source_uri`
- `protected_ref`
- `workflow_identity`
- `release_captain`: accountable identity object
- `required_approval_refs`
- `allowed_private_equivalent_profile`
- `created_at`
- `accountability`

Rules:

- Any signer identity that does not match this policy fails trusted release verification.
- A generic OIDC-authenticated actor is not enough; the signer must match issuer, repository or source URI, workflow identity, protected ref, and approval policy.
- Private equivalent profiles must specify envelope binding, trusted root or identity source, timestamp or freshness evidence, and audit-log or compensating-control status.

## Entity: Contract Release Verification

Represents the result of checking a checkout against a contract manifest and signing profile.

Fields:

- `manifest_ref`
- `manifest_digest`
- `manifest_digest_status`: matched, mismatch, missing, not_assessed
- `artifact_digest_status`: matched, mismatch, missing, not_assessed
- `signature_profile`
- `signature_status`: valid, invalid, missing, stale, not_assessed
- `signer_identity`
- `oidc_issuer`
- `identity_policy_ref`
- `transparency_log_status`: included, unavailable, not_required, not_assessed
- `freshness_status`: current, expired, rollback_suspected, not_assessed
- `verified_at`
- `provenance_refs`
- `accountability`

Rules:

- Trusted contract release requires matched manifest and artifact digests plus valid signature or approved equivalent identity-policy verification.
- A signature is valid only when the signer matches the trusted identity policy.
- Private or air-gapped environments may record a non-public equivalent, but cannot silently treat missing verification as trusted.
- `not_assessed` signature status is allowed as evidence but cannot support a `trusted_contract_release` claim.
- Verification evidence is structural. It does not prove model quality or business acceptance of residual risk.

## Entity: Evidence Event

Represents one inspectable proof item.

Fields:

- `id`: stable event identifier
- `source`: producer such as `github-actions`, `opencode`, `manual`, `test-log`, or `file-inspection`
- `external_ref`: source-local reference such as run URL, command ID, file path, or review ID
- `observed_at`: when the source observed the event
- `collected_at`: when `sdp-trace` collected or recorded it
- `actor`: human, agent, system, or tool identity
- `event_type`: command, file, test, ci_run, review, scan, deployment, harness_log, model_output, manual, custom
- `status`: success, failure, warning, skipped, pending, not_assessed
- `summary`: short human-readable description
- `artifact_uri`: optional inspectable artifact reference
- `artifact_hash`: optional hash for artifact integrity
- `artifact_hash_algorithm`: sha256 unless a future schema version explicitly adds another algorithm
- `redaction_status`: none, redacted, withheld, not_assessed
- `redaction_note`: required when `redaction_status` is not none
- `integrity_status`: verified_hash, external_asserted, unverified, not_assessed
- `external_assertions`: optional externally produced verdicts, scores, strength labels, or quality claims
- `accountability`: optional accountability record for human sign-off or evidence ownership

Rules:

- Missing evidence must not be converted into success.
- `sdp-trace` does not assign evidence strength. If a source system asserts strength, quality, or verdict, record it in `external_assertions[]` with producer and origin.
- A committed artifact reference must not contain secrets, credentials, raw customer data, or private prompt contents.
- If raw evidence cannot be committed, record a sanitized summary, hash when available, and redaction note.
- Pending evidence is not completed evidence. Metric samples that depend on pending evidence must either remain linked to pending status or record `not_assessed_reason`.
- Duplicate or conflicting events are preserved with source identity and conflict/dedupe metadata; `sdp-trace` does not collapse them into a verdict.

## Entity: Provenance Record

Represents the origin chain for an evidence event, observation, or trace snapshot.

Fields:

- `actor_id`
- `actor_type`
- `harness`
- `model_family`
- `model_version`
- `tool_name`
- `command`
- `prompt_ref`
- `context_refs`
- `captured_at`
- `hash_prev`
- `payload_digest`
- `digest_algorithm`: sha256
- `chain_scope`: optional identifier for the provenance chain

Rules:

- Model and tool identity may be `not_assessed` when unavailable.
- Provenance records origin. It does not imply quality.
- `payload_digest` is a SHA-256 digest of the canonicalized recorded payload.
- `hash_prev` is optional and only links records inside the same `chain_scope`.
- A digest proves content continuity for the committed record. It is not an authentication signature.
- Signed producer assertions may be recorded as external assertions. Contract release signing uses the contract manifest and release verification entities.

## Entity: Observation

Represents a dated, evidence-backed statement about process state.

Fields:

- `id`
- `scope`
- `observed_at`
- `statement`
- `evidence_refs`
- `provenance_refs`
- `assessment_status`: assessed, partial, not_assessed
- `accountability`: optional accountability record when the observation is manually asserted or approved

Rules:

- Observations are not policy verdicts.
- Observations may be consumed by external policy consumer.

## Entity: Metric Sample

Represents one measured value for a process metric.

Fields:

- `metric_name`
- `value`
- `unit`
- `window_start`
- `window_end`
- `dimensions`
- `evidence_refs`
- `provenance_refs`
- `not_assessed_reason`
- `assessment_state`: assessed, partial, not_assessed

Common dimensions:

- repository
- team
- scope
- harness
- model_family
- model_version
- language
- build_system
- stack
- phase

Rules:

- Each sample must have evidence or `not_assessed_reason`.
- Missing dimensions are not the same as `not_assessed`. A dimension may be absent or `not_applicable` while the sample itself remains assessed.
- Thresholds and traffic-light ratings are external policy.

## Entity: Metric Stream

Represents ordered samples for the same metric and comparable dimensions.

Fields:

- `metric_name`
- `dimensions`
- `samples`
- `comparisons`
- `created_at`
- `updated_at`
- `assessment_state`: assessed, partial, not_assessed
- `not_assessed`: optional list of stream-level gaps

Rules:

- Stream comparison shows movement by referencing a previous sample/window and recording previous value, current value, delta, unit, and evidence coverage.
- `sdp-trace` may record a raw numeric or categorical delta. It must not label that movement as degradation, improvement, pass, fail, ready, or blocked.
- Policy labels may appear only as `External Verdict Input` records.
- If any sample or comparison is partial or `not_assessed`, the stream-level `assessment_state` must be `partial` or `not_assessed`.

## Entity: External Verdict Input

Represents a policy or quality assertion produced outside `sdp-trace`.

Fields:

- `id`
- `producer`: system, person, or organization that produced the verdict
- `producer_type`: policy_engine, reviewer, ci_system, customer, other
- `origin`: always external
- `verdict_kind`: gate, degradation, readiness, evidence_strength, quality_score, override, custom
- `verdict_value`: source-local value
- `policy_ref`: optional policy, rubric, or rule identifier used by the producer
- `generated_at`
- `artifact_uri`
- `artifact_hash`
- `artifact_hash_algorithm`
- `redaction_status`
- `provenance_refs`

Rules:

- External verdict inputs are evidence about another system's decision. They are not native `sdp-trace` decisions.
- Assessment inputs may include external verdict inputs as context, but must not convert them into a new `sdp-trace` verdict.
- Schemas must make `origin: external` explicit for this entity.

## Entity: Trace Snapshot

Represents a point-in-time graph across delivery artifacts.

Node kinds:

- spec
- plan
- task
- change
- evidence
- observation
- metric_sample
- external_verdict
- decision
- actor

Relation examples:

- task implements spec
- change fulfills task
- evidence supports observation
- metric_sample derived_from evidence
- external_verdict consumes assessment_input
- decision references external_verdict

Rules:

- A trace snapshot can link to an external verdict input, but native `decision` nodes are only recorded decisions from outside `sdp-trace`.
- Schema version used by each node or package must be recorded once full validation is enabled.

## Entity: Assessment Input

Represents the package handed to external policy consumer or another policy engine.

Fields:

- `id`
- `scope`
- `trace_snapshot_ref`
- `evidence_bundle_refs`
- `metric_stream_refs`
- `observations`
- `not_assessed`
- `generated_at`
- `schema_version`
- `producer`
- `accountability`
- `risk_classification`
- `contract_release_verification_ref`

Rules:

- It must be usable without Beads.
- It must not contain policy decisions owned by external policy consumer.
- It may include external verdict inputs as observed evidence, clearly marked as external.
- `schema_version` follows semver. Additive optional fields are minor changes; required field removals or semantic changes are major changes.
- It must identify human accountability for the package. AI actors may be producers in provenance, but not sole accountable owners.
