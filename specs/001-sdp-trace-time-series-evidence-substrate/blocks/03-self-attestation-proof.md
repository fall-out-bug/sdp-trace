# Block 03: Self-Attestation Proof

Status: implemented; final review findings fixed
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Beads mirror: `sdp-trace-cdn.13`
Audience: technical executive, CIO, CEO, release owners, future external policy consumers

## Purpose

Self-Attestation Proof verifies that the `sdp-trace` contract used by observers reports source reference, digest, local signature, identity policy, freshness, external attestation, and production trust as separate proof states.

The block answers the CIO question:

> Could a person or model silently change the contract and still make the repository look trusted?

## Proof States

The verifier must report these states separately:

- `schema_valid`
- `digest_verified`
- `locally_attested`
- `externally_attested`
- `production_release_verified`

Missing proof remains `not_assessed`; it must not be collapsed into trusted.

## In Scope

- immutable source reference in contract manifest
- self-attestation verification command
- self-attestation evidence record under `examples/self-trace/`
- negative fixtures for wrong source commit, wrong signer, wrong trusted identity policy, stale manifest, missing external attestation, and modified verification artifact
- summary in retired research artifact

## Out of Scope

- full TUF metadata
- enterprise PKI implementation
- external policy consumer pass/fail policy
- customer pilot execution

## Acceptance

Block 03 passes only when the verification command reports each proof state independently and refuses to call the release trusted when immutable source identity, signer identity, freshness, or external attestation is missing.

In the current shared dirty working tree, source-content verification is explicitly `not_assessed` because the manifest `source_commit` does not contain the final artifact set. That means `locally_attested` is assessed as `false`; this is the correct product state until final artifacts are committed and the manifest, DSSE envelope, and result are regenerated against that commit.

Production trust remains blocked until `externally_attested` and `production_release_verified` are backed by an accepted trust anchor such as GitHub OIDC plus Sigstore/Rekor or customer PKI.
