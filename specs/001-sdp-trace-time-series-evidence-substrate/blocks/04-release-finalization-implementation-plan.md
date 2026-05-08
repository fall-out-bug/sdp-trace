# Block 04 Implementation Plan: Source-Bound Release Finalization

Status: accepted for implementation
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Beads mirror: `sdp-trace-cdn.21`

## Goal

Implement Block 04 as a staged release-finalization capability:

1. Add a source-bound local finalization command that refuses dirty-tree source-content proof.
2. Extend release verification data to expose external trust proof states explicitly.
3. Add negative fixtures for unsupported external production trust claims.
4. Update docs so local source proof and external production trust cannot be confused.

## Consensus

Consensus is reached on the "local plus external trust design" scope. The user selected option 2 and rejected an additional confirmation gate as unnecessary. Implementation can proceed under the claim boundary in this plan.

## File Responsibilities

- `scripts/finalize-source-bound-release.sh`: command entrypoint for source-bound local finalization. It validates the repository cleanliness boundary before any source-content proof can be generated.
- `schema/contract-release-verification.schema.json`: adjacent release-verification schema that records source proof, external trust profile, protected ref, workflow identity, approval, transparency, and production verification states.
- `examples/contract-foundation/contract-release-verification.example.json`: local development release verification example. It must remain `trusted_contract_release: false`.
- `examples/contract-foundation/negative-trusted-release-*.json`: negative fixtures proving that production trust cannot be claimed when external proof states are missing or contradictory.
- `scripts/validate-contracts.sh`: validation entrypoint for schema examples and negative fixtures.
- `docs/contract-release-signing.md`: user-facing release signing and finalization guidance.
- retired research artifact: evidence note for residual `not_assessed` states.
- `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md`: SpecKit task status sync.

## Task 1: Source-Bound Finalization Guard

- Add `scripts/finalize-source-bound-release.sh`.
- The command must inspect the selected git repository before generating source-content proof.
- If `git status --porcelain --untracked-files=normal` is non-empty, it must exit non-zero and say that source-bound finalization is refused from a dirty working tree.
- If the repository is clean, it must resolve the requested source ref to a 40-character commit SHA, validate that the manifest exists, and verify every manifest artifact path exists at that commit with a matching SHA-256 digest.
- The implemented guard stops after source-content verification and emits a JSON finalization result. It must not claim external or production trust.
- Regenerating the manifest, DSSE envelope, and self-attestation result is a separate release-proof generation step; it must not be implied by this guard.

Verification:

```bash
rtk scripts/finalize-source-bound-release.sh --repo /path/to/dirty/git/repo --manifest manifest.json --source-ref HEAD
```

Expected result: non-zero exit with `Refusing source-bound finalization from dirty working tree`.

## Task 2: Release Verification Schema Extension

- Extend `schema/contract-release-verification.schema.json` with explicit fields:
  - `source_commit`
  - `source_commit_status`
  - `source_commit_artifact_status`
  - `source_commit_artifact_counts`
  - `external_trust_profile`
  - `external_attestation_ref`
  - `transparency_evidence_ref`
  - `source_uri_status`
  - `protected_ref_status`
  - `workflow_identity_status`
  - `approval_status`
  - `production_release_verified`
- For `trusted_contract_release: true`, require all local source, digest, identity, transparency, protected ref, workflow, approval, freshness, and production verification states to be successful.
- For `trusted_contract_release: true`, reject the local development external trust profile.

Verification:

```bash
rtk node scripts/validate-json-schema.mjs schema/contract-release-verification.schema.json examples/contract-foundation/contract-release-verification.example.json
```

Expected result: valid, while `trusted_contract_release` remains false.

## Task 3: External Trust Negative Fixtures

- Add negative fixtures for:
  - missing Rekor or customer audit evidence
  - OIDC issuer or identity policy mismatch
  - source URI mismatch
  - protected ref mismatch
  - workflow identity mismatch
  - approval mismatch
  - expired freshness
  - local profile used for production trust
- Wire every fixture into `scripts/validate-contracts.sh` with `expect_fail`.

Verification:

```bash
rtk npm run validate
```

Expected result: all negative fixtures fail validation for the intended reason and the full validation suite exits zero.

## Task 3A: Non-Circular Artifact Boundary

- Add validation that the contract manifest subject set excludes generated proof artifacts:
  - `examples/contract-foundation/contract-release-verification.example.json`
  - `examples/contract-foundation/contract-release.dsse.json`
  - `examples/contract-foundation/local-dev-signing-public.pem`
  - `examples/self-trace/self-attestation-verification-result.json`
- The check must fail if any of these paths appear in `examples/contract-foundation/contract-manifest.example.json`.

Verification:

```bash
rtk npm run validate
```

Expected result: validation fails if generated proof artifacts become manifest subjects.

## Task 4: Documentation Sync

- Update release docs to use these terms consistently:
  - `source-bound local release`
  - `externally trusted production release`
- State that local DSSE, schema validity, and private key possession cannot produce `trusted_contract_release: true`.
- Update self-attestation summary with exact remaining `not_assessed` states.
- Update SpecKit tasks T064-T068 only after implementation and verification complete.

Verification:

```bash
rtk rg -n "source-bound local release|externally trusted production release|trusted_contract_release" docs specs/001-sdp-trace-time-series-evidence-substrate
```

Expected result: wording consistently preserves the proof boundary.

## Task 5: Review and Fix Loop

- Run local self-review after implementation.
- Run validation:

```bash
rtk npm run validate
rtk git diff --check
```

- Run pi/subagent review after the block implementation.
- Register every review finding in Beads, including minor/P3 findings.
- Fix every valid finding before marking Block 04 complete.
