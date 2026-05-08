**Verdict: APPROVE**

No critical or major tracing/evidence findings remain in the PR-level diff. The implementation faithfully applies the closed reason-code registry, trust-scope rules, fixture expectations, environment-only non-upgrade, and review-ledger statements from the spec. The remaining gaps are minor test-coverage omissions that do not violate the contract.

---

### Minor Findings

**1. Missing dedicated test for `buildkite-missing-independent-signer` fixture**

- **File/function**: `internal/witness/profiles_test.go` / test suite
- **Spec evidence**: Fixture matrix row `buildkite-missing-independent-signer` expects `cannot_verify` / `witness_signer_authority_missing` with `buildkite-v1` profile. The shared `validateCIEnvelope` function contains the correct logic, but no test calls `BuildCIEnvelopeProfile(KindBuildkite, ...)` with an envelope whose `ProfileStates.SignerAuthorityState != statePass`.
- **Required fix**: Add a test case in `TestCIEnvelopeNonPassReasonCodes` (or a parallel Buildkite version) that constructs a Buildkite envelope with `SignerAuthorityState: stateCannotVerify` and asserts `StatusCannotVerify`, `ReasonMissingSigner`. (The existing GitLab-only table test doesn’t exercise the Buildkite kind.)

**2. Missing dedicated test for `ci-run-id-mismatch` with Buildkite kind**

- **File/function**: `internal/witness/profiles_test.go` / `TestCIEnvelopeNonPassReasonCodes`
- **Spec evidence**: Fixture matrix row `ci-run-id-mismatch` applies to `gitlab-ci-v1` or `buildkite-v1`. The current test mutates `RunBindingState = stateFail` only for GitLab (`KindGitLabCI`).
- **Required fix**: Add a parallel sub-test for `KindBuildkite` with the same mutation, asserting `StatusFail` / `ReasonRunMismatch`, to prove the replay check works identically for both profiles.

**3. Trust-scope field values in `cannot_verify` paths**

- **Files**: `internal/witness/profiles.go` (`BuildCIEnvelopeProfile` no-envelope branches), schema `trust_scope` enum.
- **Spec evidence**: The record’s `TrustScope` field is set to `local_observed` when the envelope is missing, while `EstablishedTrustScope` is `cannot_verify`. The schema expects `trust_scope` values to be one of the closed vocabulary. This is semantically correct (local structural evidence is present), but downstream consumers looking only at `trust_scope` may see `local_observed` instead of `cannot_verify`. The spec’s normalized result description leaves the exact mapping for these edge cases implicit.
- **Required fix**: No behavioural change is required—the witness result still correctly exposes `established_trust_scope: cannot_verify`. Optionally, document that `trust_scope` reflects the local evidence floor while `established_trust_scope` is the verifier-derived scope, so reviewers using the schema don’t misinterpret the field.

---

All three findings are minor (test coverage and documentation). They do not block the PR from an tracing/evidence correctness standpoint.
