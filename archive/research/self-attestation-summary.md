# Self-Attestation Summary

Date: 2026-05-01

Scope: Block 03 Self-Attestation Proof for the local development contract release.

## Purpose

Self-attestation separates contract shape, digest integrity, local DSSE attestation, external attestation, and production release verification. It must not collapse missing external trust anchors into a trusted product claim.

## Commands

```bash
scripts/verify-self-attestation.sh --case examples/self-trace/self-attestation-verification.json
scripts/verify-self-attestation.sh --all
npm run validate
```

Observed result on 2026-05-01: the local structural case reported assessed `schema_valid`, assessed `digest_verified`, assessed `source_commit_artifacts_verified: false`, assessed `locally_attested: false`, `externally_attested: not_assessed`, and `production_release_verified: not_assessed`. `locally_attested` was false because the shared working tree contained final artifacts that were not present in the manifest `source_commit`.

Regenerated on 2026-05-04 after Block 06 manifest subjects changed. The manifest digests, DSSE envelope, public key, verification artifact, and self-attestation case were refreshed. Current observed result: assessed `schema_valid: true`, assessed `digest_verified: true`, assessed `source_commit_artifacts_verified: true`, assessed `locally_attested: true`, `externally_attested: not_assessed`, `production_release_verified: not_assessed`. The negative suite continues to detect wrong source commit, wrong signer, wrong trusted identity policy, stale manifest, missing external attestation, modified verification artifact, and unsupported production-trust claims.

Block 04 adds the next release-finalization boundary: a source-bound local release can assess source-content proof only from a clean git tree and a source commit containing every manifest artifact with matching digest. It still cannot claim an externally trusted production release without accepted Sigstore/Rekor or customer PKI evidence.

## Evidence

- Positive case: `examples/self-trace/self-attestation-verification.json`
- Positive verification result: `examples/self-trace/self-attestation-verification-result.json`
- Verifier: `scripts/verify-self-attestation.sh`
- Source-bound finalization guard: `scripts/finalize-source-bound-release.sh`
- Contract manifest: `examples/contract-foundation/contract-manifest.example.json`
- Local DSSE envelope: `examples/contract-foundation/contract-release.dsse.json`
- Public key: `examples/contract-foundation/local-dev-signing-public.pem`

## Negative Fixtures

- `examples/self-trace/self-attestation-negative-wrong-source-commit.json`
- `examples/self-trace/self-attestation-negative-wrong-signer.json`
- `examples/self-trace/self-attestation-negative-wrong-policy.json`
- `examples/self-trace/self-attestation-negative-stale-manifest.json`
- `examples/self-trace/self-attestation-negative-missing-external-attestation.json`
- `examples/self-trace/self-attestation-negative-modified-verification-artifact.json`
- `examples/contract-foundation/negative-trusted-release-missing-external-evidence.json`
- `examples/contract-foundation/negative-trusted-release-oidc-issuer-mismatch.json`
- `examples/contract-foundation/negative-trusted-release-source-uri-mismatch.json`
- `examples/contract-foundation/negative-trusted-release-protected-ref-mismatch.json`
- `examples/contract-foundation/negative-trusted-release-workflow-identity-mismatch.json`
- `examples/contract-foundation/negative-trusted-release-approval-mismatch.json`
- `examples/contract-foundation/negative-trusted-release-expired-freshness.json`
- `examples/contract-foundation/negative-trusted-release-local-profile.json`

## Residual Not Assessed

- External Sigstore/Rekor or customer PKI attestation is not committed.
- Production release verification is not claimed.
- Source-content verification is assessed true after the 2026-05-04 manifest refresh. The source-bound finalization guard was used to verify that the selected `source_commit` contains every manifest artifact path with matching digest before regenerating the release proof.
- Release verification example fields `source_uri_status`, `protected_ref_status`, `workflow_identity_status`, and `approval_status` remain `not_assessed` because no external Sigstore/Rekor or accepted customer PKI evidence is committed.
- Fresh checkout reproducibility is not claimed in this dirty working tree session.
- The manifest `source_commit` is an immutable reference, but this repository still needs an external trust artifact before `trusted_contract_release: true`.

## Product Claim Boundary

Block 03 proves local structure, digest continuity, DSSE envelope binding, trusted identity policy enforcement for the local private-equivalent profile, freshness evaluation, source-content downgrade behavior, and negative-case detection.

Block 04 adds the source-bound local release guard and external production trust schema. It does not prove production trust in this dirty working tree. `trusted_contract_release: true` remains blocked until local source proof, external attestation, transparency or customer audit evidence, protected ref, workflow identity, approval, freshness, and production release verification are all assessed and successful.
