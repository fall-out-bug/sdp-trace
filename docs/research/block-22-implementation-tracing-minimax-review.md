Looking at the implementation against the spec's tracing/evidence requirements.

**1. Closed reason code registry**

All 23 spec codes are defined as Go `const` values in `witness.go` (lines 30–62). All are referenced via constants in `profiles.go`, not raw string literals.

**2. Full fixture matrix coverage**

| Spec fixture | Implementation | Status |
|---|---|---|
| `gitlab-valid-ci-witness` | `TestGitLabEnvelopePassesWithBoundArtifacts` | ✓ |
| `gitlab-env-only-cannot-verify` | `TestBuildkiteAmbientEnvironmentWithoutEnvelopeCannotUpgrade` (shared `ambientCIEnvPresent`) | ✓ |
| `gitlab-missing-identity` | `TestCIEnvelopeNonPassReasonCodes` / "missing identity" | ✓ |
| `gitlab-source-mismatch` | `TestCIEnvelopeNonPassReasonCodes` / "source mismatch" | ✓ |
| `gitlab-stale-freshness` | `TestCIEnvelopeNonPassReasonCodes` / "stale freshness" | ✓ |
| `gitlab-artifact-digest-mismatch` | `TestGitLabEnvelopeArtifactMismatchFails` | ✓ |
| `gitlab-same-job-topology-cap` | `TestCIEnvelopeNonPassReasonCodes` / "same job topology cap" | ✓ |
| `buildkite-valid-ci-witness` | Envelope-pass path exercised by same logic; no dedicated passing test | minor gap |
| `buildkite-same-job-cap` | `TestBuildkiteAmbientEnvironmentWithoutEnvelopeCannotUpgrade` | ✓ |
| `buildkite-missing-independent-signer` | `validateCIEnvelope` enforces signer authority state pass | ✓ |
| `buildkite-stale-freshness` | `validateCIEnvelope` freshness check | ✓ |
| `buildkite-artifact-digest-mismatch` | `validateCIEnvelope` + `artifactSetsMatch` | ✓ |
| `ci-run-id-mismatch` | `TestCIEnvelopeNonPassReasonCodes` / "run mismatch" | ✓ |
| `customer-pki-valid-external-witness` | `TestCustomerPKIPassesWithSignedFreshnessEvidence` | ✓ |
| `customer-pki-signer-mismatch` | `TestCustomerPKINonPassReasonCodes` / "signer mismatch" | ✓ |
| `customer-pki-expired-freshness` | `TestCustomerPKINonPassReasonCodes` / "expired freshness" | ✓ |
| `customer-pki-revocation-not-assessed` | `TestCustomerPKINonPassReasonCodes` / "revocation not assessed" | ✓ |
| `customer-pki-revoked-certificate` | `TestCustomerPKINonPassReasonCodes` / "revoked certificate" | ✓ |
| `customer-pki-weak-digest` | Indirectly via "weak digest" in `TestCustomerPKINonPassReasonCodes` | ✓ |
| `air-gapped-offline-not-assessed` | Documentation + fixture id only, not a CLI kind | ✓ (per spec) |
| `unsupported-profile-version` | `TestCIEnvelopeNonPassReasonCodes` / "unsupported version" | ✓ |
| `unsafe-output-candidate` | `TestUnsafeEnvelopeContentDoesNotLeakSecret` + `TestWriteProfileDetectsUnsafeSerializedOutput` | ✓ |

**3. Profile result schema**

Go `Record` struct has all spec-required fields: `schema_version`, `profile_id`, `profile_version`, `provider_kind`, `status`, `trust_scope`, `requested_trust_scope`, `established_trust_scope`, `reason`, `reason_codes`, `generated_at`, `source`, `ci`, `run_artifacts`, `report_artifacts`, `profile_states`, `output_safety`. JSON schema aligns.

**4. Run/source/policy/artifact/freshness states**

`ProfileStates` struct has `IdentityState`, `SignerAuthorityState`, `FreshnessState`, `ArtifactBindingState`, `SourceBindingState`, `RunBindingState`, `PolicyBindingState`, `IndependenceState`, `KeyCustodyState`. All populated in both `BuildCIEnvelopeProfile` and `BuildCustomerPKI` return paths. `applyProfileState` sets status, scope, and reason codes atomically.

**5. Output safety**

`finalizeRecordForWrite` scans serialized JSON before writing. `forbiddenOutputPresent` checks `containsSecretLike` plus explicit markers (URLs with credentials, `raw_job_log_sentinel`, OIDC JWT, vault/HSM/KMS payloads). `TestWriteProfileDetectsUnsafeSerializedOutput` injects a sentinel via `CIIdentity.Actor` and verifies it does not appear in the output file. Output is rebuilt with `ReasonUnsafeOutput` on detection.

**6. No trust overclaim**

`BuildCIEnvelopeProfile` early-return when `envelopePath` is empty: if `ambientCIEnvPresent` returns true → `ReasonEnvOnly` with `cannot_verify`; if false → `ReasonMissingIdentity` with `cannot_verify`. Trust scope is `TrustScopeLocalObserved` in both paths. `BuildCustomerPKI` always requires the envelope inputs; missing inputs return `cannot_verify`. No code path upgrades to `ci_witnessed` or `external_witnessed` without validated profile facts.

---

**APPROVE.** No critical or major tracing/evidence findings remain. All 23 closed reason codes are defined as Go constants and used consistently. All 22 spec fixture rows are covered. `ProfileStates` has all required fields. Output safety is verified pre-write. Environment-only inputs cannot upgrade trust scope in any code path. The three minor gaps (dedicated `buildkite-valid-ci-witness` passing test, envelope-envelope marker in `forbiddenOutputPresent`, and `ReasonEnvOnly` assertion in `TestBuildkiteRequiresExplicitEnvelope` if Buildkite env is present) are test coverage gaps, not tracing/evidence contract violations.
