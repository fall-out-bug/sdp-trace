# Block 04: Release Finalization and External Trust Design

Status: implemented artifacts present; current closure stale until verifier evidence is repaired
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Beads mirror: `sdp-trace-cdn.21`
Audience: CTO, CIO, CEO, release owners, future `sdp-gate` consumers

## Purpose

Block 04 turns the Block 03 self-attestation learning into a release-finalization contract.

It has two deliberately separate goals:

1. Close the local source-content gap by defining how a final committed artifact set becomes a source-bound local release.
2. Define the external trust anchor shape for Sigstore/Rekor or customer PKI without claiming production trust until real external evidence exists.

The block answers the CIO question:

> Can this repository distinguish "the files are locally finalized against a commit" from "the release is externally trusted"?

## Product Boundary

`sdp-trace` records proof states. It does not accept risk, approve a release, or decide customer readiness.

`trusted_contract_release` must remain false unless every required proof state is assessed and successful under the selected trust profile. Missing external evidence is `not_assessed`, not a warning, score, or implicit approval.

## Proof State Model

Block 04 keeps the Block 03 proof states and adds more explicit source and external trust details:

- `schema_valid`
- `digest_verified`
- `source_commit_artifacts_verified`
- `locally_attested`
- `externally_attested`
- `production_release_verified`

For a local source-bound release:

- `source_commit_artifacts_verified`: assessed true
- `digest_verified`: assessed true
- `locally_attested`: assessed true
- `externally_attested`: `not_assessed`
- `production_release_verified`: `not_assessed`
- `trusted_contract_release`: false

For an externally trusted production release:

- `source_commit_artifacts_verified`: assessed true
- `digest_verified`: assessed true
- `locally_attested`: assessed true
- `externally_attested`: assessed true
- `production_release_verified`: assessed true
- `trusted_contract_release`: true

## In Scope

- A release-finalization spec that requires an immutable committed source reference before source-content proof can be assessed.
- A local finalization command design that refuses to claim source-content proof from a dirty working tree.
- External trust anchor design for:
  - Sigstore/Cosign keyless DSSE with GitHub OIDC and Rekor transparency evidence.
  - Customer PKI or private equivalent with an explicit policy, certificate/key identity, timestamp or audit evidence, and accepted transparency substitute.
- Negative fixture requirements for:
  - dirty worktree finalization
  - source commit mismatch
  - missing Rekor or customer audit evidence
  - OIDC issuer mismatch
  - source URI mismatch
  - protected ref mismatch
  - workflow identity mismatch
  - approval ref mismatch
  - expired release freshness
- Documentation updates that say which proof states are assessed and which remain `not_assessed`.

## Out of Scope

- Implementing a real public GitHub Actions release workflow in this design pass.
- Requiring Sigstore/Rekor for private or air-gapped customers.
- Storing private keys, raw certificates with secrets, credentials, raw customer data, or private prompt contents in the repository.
- Adding `sdp-gate` policy thresholds.
- Claiming pilot/customer readiness from local finalization alone.

## Recommended Approach

Use a two-layer release model.

Layer 1 is **source-bound local finalization**. It binds the manifest subject artifact set to `source_commit` and verifies that the selected commit contains every manifest artifact path with matching SHA-256 digest.

`source_commit` is the immutable 40-character git commit that contains the source artifacts named by the manifest. Generated proof artifacts that embed the manifest digest, DSSE envelope, public key fingerprint, or verification result are outside the manifest subject set and may be produced after this source commit. This avoids a self-referential commit hash and prevents circular digest churn.

Layer 2 is **external trust attestation**. It records either a public Sigstore/Rekor bundle or a customer-approved private equivalent. This layer is optional for local engineering proof and mandatory for `production_release_verified`.

This design is stricter than treating local DSSE as trusted, but it is the only option that preserves the UX promise: a CTO or CIO can inspect exactly what is proven, what is missing, and who is accountable.

## Release Finalization Flow

The implemented source-bound guard preserves this sequence:

1. Require a clean committed source tree for source-bound finalization, or fail closed without emitting a source proof artifact.
2. Resolve the source reference to a 40-character commit SHA.
3. Verify:
   - every manifest artifact exists at `source_commit`
   - every source artifact digest matches the manifest
   - generated proof artifacts are not part of the manifest subject set
4. Emit a source-bound local result with external proof states as `not_assessed`.

A later release-proof generation step may regenerate the local DSSE envelope and self-attestation result over the manifest digest. That generation step does not change which source artifacts are assessed by `source_commit_artifacts_verified`.

## External Trust Anchor Profiles

### Public Sigstore/Rekor Profile

Signature profile: `sdp-trace-signature/sigstore-dsse-keyless-v1`.

External trust profile: `sigstore-rekor-keyless-v1`.

Required evidence:

- DSSE envelope over an in-toto Statement whose subject is the contract manifest digest.
- OIDC issuer matching trusted identity policy.
- Source URI matching trusted identity policy.
- Protected ref or release tag matching trusted identity policy.
- Workflow identity matching trusted identity policy.
- Rekor inclusion proof.
- Approval refs matching trusted identity policy.
- Freshness and rollback status assessed.

If any required evidence is missing, `externally_attested` and `production_release_verified` are `not_assessed` or assessed false depending on whether evidence is missing or contradictory.

### Customer PKI / Private Equivalent Profile

Required evidence:

- DSSE envelope or customer-approved equivalent binding the manifest digest.
- Certificate, public key, or identity fingerprint accepted by the customer's trusted identity policy.
- Timestamp, audit log, or offline transparency substitute.
- Customer-approved source reference and protected release channel.
- Approval refs accepted by the customer's governance process.
- Explicit customer policy acceptance for the private trust profile.
- Freshness and rollback status assessed.

Private equivalent trust is not "we skipped external trust"; it is a different external trust anchor. The artifact must name the profile and evidence used.

## Data Model Impact

Block 04 should extend or complement the current self-attestation case/result shape with:

- `source_commit_artifacts_verified`
- `source_commit_artifact_counts`
- `external_trust_profile`
- `external_attestation_ref`
- `transparency_evidence_ref`
- `protected_ref_status`
- `workflow_identity_status`
- `approval_status`, including release captain and required approval refs
- `production_release_verified`

The implementation should avoid a generic `trusted: true` shortcut. Every trusted state must be derivable from explicit proof states.

## UX Requirements

- A repository observer must be able to read one result artifact and understand:
  - what was locally proven
  - what was externally proven
  - what remains `not_assessed`
  - why `trusted_contract_release` is true or false
- Missing external evidence must be obvious in the first screen of the result artifact, not buried in prose.
- The docs must use "source-bound local release" and "externally trusted production release" as distinct terms.

## Acceptance

Block 04 design is accepted when:

- The spec explicitly separates source-bound local finalization from external production trust.
- The Socratic dialogue records the main objections and resolutions.
- No artifact claims external or production trust before real external evidence exists.

Block 04 implementation started after consensus was accepted for option 2: local source-bound finalization plus external trust schema/fixtures.

## Implemented Artifacts

- `scripts/finalize-source-bound-release.sh` verifies clean source-bound artifact proof and refuses dirty working trees.
- `schema/contract-release-verification.schema.json` records source commit artifact status, external trust profile, transparency evidence, source URI, protected ref, workflow identity, approval, and production release verification states.
- `examples/contract-foundation/negative-trusted-release-*.json` fixtures reject unsupported `trusted_contract_release: true` claims.
- `docs/contract-release-signing.md` and `archive/research/self-attestation-summary.md` distinguish a source-bound local release from an externally trusted production release.

## Delivery State

The selected scope is "local plus external trust design." Delivery remains staged:

1. Source-bound local finalization guard is implemented.
2. External trust evidence schema/fixtures are implemented.
3. Real Sigstore/Rekor or customer PKI verification remains `not_assessed` unless real evidence is available in this repository.
4. Historical Block 04 validation closure is stale in the current checkout until live verifier output supports it again.
