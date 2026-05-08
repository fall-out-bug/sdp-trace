# Block 01: Contract Foundation Socratic Review

Status: author-revised
Block Spec: `01-contract-foundation.md`
Critic Provider: `kimi-coding/k2p6`
Judge Provider: `minimax/MiniMax-M2.7`

## Critic Rubrics

- problem and goal
- system boundaries
- roles and actors
- core scenarios
- non-goals
- assumptions and dependencies
- edge cases and resilience
- security and access
- observability and metrics
- testability and acceptance
- rollout, migration, and backward compatibility
- open questions and risks

## Critic Questions and Author Resolutions

### Q1: Pilot Scope Ambiguity

Severity: blocking

Question: The block excludes pilot execution, while the parent tasks include OpenCode, harness, model, and Kotlin+Bazel pilot run-cards. Are T027-T033 in this block or future work?

Resolution: The block now distinguishes contract fixtures from pilot execution. T027-T033 are later pilot-matrix work and consume this block's schemas. Contract Foundation may define generic fixture shapes and run-card contract skeletons, but it does not execute pilots or close compatibility matrix rows.

### Q2: Negative Fixture Enforceability

Severity: major

Question: JSON Schema cannot safely reject arbitrary prose values like `pass` or `fail` without false positives.

Resolution: The negative fixture is now field-path based. It rejects forbidden native policy fields such as `verdict`, `gate_result`, `readiness_verdict`, `degradation_status`, `policy_threshold`, `evidence_strength`, and related names outside `external-verdict-input`. Policy words inside prose summaries or external logs are not rejected by schema.

### Q3: Offline Validator Dependency

Severity: major

Question: `npx --yes ajv-cli@5.0.0` requires network access. Is that acceptable for CI or air-gapped pilots?

Resolution: The default developer path must not require live network access in CI or pilot checkout. The initial `ajv-cli@5.0.0` answer was superseded on 2026-05-01 after dependency review; current accepted validation uses exact local `ajv@8.20.0` through `scripts/validate-json-schema.mjs`.

### Q4: Typed `not_assessed` Representation

Severity: major

Question: How can `not_assessed` appear in typed fields like numeric metric values or string model versions without schema violations?

Resolution: `not_assessed` is now defined as state, not a string sentinel. Measured values use `assessment_state: "not_assessed"`, `value: null`, and `not_assessed_reason`. Optional unavailable metadata is omitted and listed in `unavailable_fields[]`. Required identity and boundary fields fail validation when missing.

### Q5: Structural Conformance vs Trustworthiness

Severity: major

Question: Schema-valid data can still be fraudulent or low quality. What kind of trust does the block prove?

Resolution: The block now states that it proves structural trust only: schema conformance, digest fields, redaction status, provenance completeness, and explicit `not_assessed` reasons. Producer honesty, artifact authenticity, fraud detection, and signed attestations remain external or future work.

### Q6: Secret and Raw Data Safety

Severity: major

Question: AC08 bans secrets and raw customer data, but is the check automated or manual?

Resolution: AC08 now requires a documented automated artifact-safety scan. The scan must fail committed examples containing obvious secrets, credentials, raw customer data markers, or private prompt contents.

### Q7: external policy consumer Schema Version Consumption

Severity: major

Question: Where does external policy consumer declare consumed schema versions, who maintains it, and how are breaking changes visible?

Resolution: The block now requires a consumer schema-version declaration example. `sdp-trace` defines and validates the portable declaration shape. The actual external policy consumer repository owns its real declaration. This makes contract changes visible without making external policy consumer a runtime dependency.

### Q8: Dedupe and Conflict Metadata

Severity: minor

Question: Edge cases mention dedupe and conflicts, but the evidence event entity lacks concrete fields.

Resolution: The positive fixture now requires `dedupe_key` and `conflict_refs` coverage. The schema task must represent duplicates and conflicts without collapsing them into a policy verdict.

### Q9: Missing Field Defaults

Severity: minor

Question: Should missing fields default to `not_assessed`, fail validation, or be optional?

Resolution: Required fields that establish identity, scope, timestamps, schema version, or artifact boundaries fail validation when missing. Optional fields may be omitted. Known unavailable assessment content must be explicit through `not_assessed` state and reason.

### Q10: Trace Schema Migration Trigger

Severity: minor

Question: When is `trace.schema.json` considered obsolete or in need of a migration note?

Resolution: T012 now has a decision rule. The existing trace schema remains the compatibility path unless it cannot represent observations, metric samples, evidence refs, external verdict refs, and schema versions without lossy or policy-owning fields. Migration notes are required for breaking or ownership-boundary changes.

## Author Assessment

The blocking ambiguity is resolved. The remaining major questions are converted into explicit contract rules and acceptance criteria. No implementation has started; these changes only tighten the block specification and review record.

## Judge Result

Verdict: PASS

The judge found no unresolved blocking, major, or minor issues. It also found no scope creep or contradictions. A new Socratic cycle is not required before preparing the implementation plan.

## Governance Expansion Review

Status: author-revised
Critic Provider: `zai/glm-5.1`
Judge Provider: `kimi-coding/k2p6`

### GQ1: Optional Accountability on Evidence

Severity: blocking

Question: If evidence events can omit accountability and later become the sole basis for a critical assessment input, who is accountable?

Resolution: Evidence events may omit direct accountability only when a containing evidence package or assessment input provides effective accountability for every referenced evidence item. Assessment inputs cannot claim completeness or trusted-release readiness if referenced evidence lacks direct or inherited accountability.

### GQ2: Authorized Release Signer

Severity: blocking

Question: Without a required signer identity policy, can any OIDC-authenticated actor produce a signed manifest that passes verification?

Resolution: The block now requires a trusted identity policy with expected OIDC issuer, source URI, protected ref, workflow identity, release captain, required approval refs, and private-equivalent verifier profile. Any signer mismatch fails trusted release verification.

### GQ3: Policy Leakage Through Risk Classification

Severity: blocking

Question: Does `required_oversight` inside `sdp-trace` prescribe policy obligations instead of recording observed state?

Resolution: Risk classification now records observed autonomy and impact plus externally declared oversight assertions. `sdp-trace` does not derive required oversight. Oversight requirements appear only as external declarations with origin and policy reference.

### GQ4: Three Lines Independence

Severity: major

Question: Can the same person appear as both first-line DRI and second-line reviewer under different labels?

Resolution: Accountability identities are now structured with `identity_ref`, `actor_type`, and line-of-defense metadata. `sdp-trace` records enough identity and line data for a policy consumer to detect self-review or same-line review; external policy consumer decides whether separation satisfies policy.

### GQ5: Stream-Level `not_assessed`

Severity: major

Question: Can a metric stream silently mix assessed and `not_assessed` samples while appearing fully assessed?

Resolution: Metric streams now carry stream-level `assessment_state`. If any sample or comparison is partial or `not_assessed`, the stream state must be `partial` or `not_assessed`.

### GQ6: Safety Scan Overclaim

Severity: major

Question: Does the artifact safety scan prove committed artifacts contain no proprietary data, or only that common patterns were checked?

Resolution: The safety scan is treated as lower-bound automated evidence. Customer-facing examples still require accountable evidence-owner approval; scan pass alone is not a proof of absence for all proprietary data.

### GQ7: Freshness Definition

Severity: major

Question: What defines an expired manifest, and can a manifest without freshness rules remain current forever?

Resolution: Contract manifests must include exactly one of `valid_until` or `freshness_policy`. Omission fails validation.

### GQ8: Machine-Readable AI Actor Detection

Severity: major

Question: How can schema validation reject AI-as-sole-accountable-owner without machine-readable actor type?

Resolution: Accountable identities are objects with `identity_ref` and `actor_type`. Accountable actor types are limited to `human_user`, `human_role`, and `human_group`.

### GQ9: Network and Air-Gapped Verification

Severity: major

Question: How does no-network validation coexist with Sigstore/Rekor/OIDC verification?

Resolution: The block distinguishes offline structural validation from release verification. Private equivalents must preserve DSSE envelope binding, trusted identity source, timestamp or freshness evidence, and audit-log or compensating-control status.

### GQ10: TUF Scope Creep

Severity: minor

Question: When does rollback/freshness become a partial TUF implementation?

Resolution: The block now limits rollback protection to linear `previous_manifest_digest` continuity plus explicit freshness. Delegated roles, threshold metadata, snapshot metadata, timestamp metadata, target metadata, and key rotation are out of scope and require a future block.

### GQ11: Authority Scope Ambiguity

Severity: minor

Question: Does `authority_scope: external_verdict` mean owning the verdict's correctness or just recording it?

Resolution: Accountability records now include `accountability_claim`: `recording_only`, `content_approval`, `risk_acceptance`, or `release_approval`.

### GQ12: Theoretical Signing Profile

Severity: minor

Question: Can the block be called complete without ever producing a real signature or verification evidence?

Resolution: The block is not complete as a trusted contract foundation until one real release verification evidence record exists for the target signing profile or an approved private equivalent. Schema-only examples prove shape, not that signing works.

### Governance Judge Result

Verdict: PASS

The judge found no unresolved blocking, major, or minor governance issues. It also found no scope creep or contradictions. A new governance cycle is not required before presenting the updated implementation plan.
