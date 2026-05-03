# sdp-trace Contract Release Signing

Signature profile: `sdp-trace-signature/sigstore-dsse-keyless-v1`.

External trust profiles:

- `sigstore-rekor-keyless-v1`
- `customer-pki-private-equivalent-v1`

Local signer implementation label:

- `local-dev-dsse-openssl-v1`

The target public profile is:

```text
contract-manifest.json
  -> SHA-256 digest
  -> in-toto Statement subject
  -> DSSE envelope
  -> Sigstore/Cosign keyless signature bundle
```

## Trusted Identity Policy

A signature is not enough. Verification must also match a trusted identity policy:

- expected OIDC issuer
- source URI
- protected branch, tag, or ref
- workflow identity
- release captain
- required approval refs

Any signer that does not match this policy fails trusted-release verification.

## Private Equivalent

Private or air-gapped environments may use an approved equivalent. The equivalent must still record:

- DSSE envelope binding
- trusted identity source
- timestamp or freshness evidence
- audit-log or compensating-control status
- explicit customer policy acceptance for the private profile
- customer-approved source reference and protected release channel
- release captain and required approval refs
- accepted transparency substitute, such as customer audit log or offline timestamp evidence

The local development fixture uses `local-dev-dsse-openssl-v1` as a signer implementation label to prove that the envelope binding and digest verification path are executable. It is not a replacement for production Sigstore identity or customer PKI trust.

The manifest digest is verified by the release verification record and DSSE envelope. Files that embed the manifest digest, such as release verification examples or negative signer fixtures, are intentionally outside the manifest artifact list to avoid circular digest dependencies.

## Source-Bound Local Release

A source-bound local release proves that the selected `source_commit` contains every manifest artifact path with the exact SHA-256 digest recorded in the manifest.

Run:

```bash
scripts/finalize-source-bound-release.sh --manifest examples/contract-foundation/contract-manifest.example.json --source-ref HEAD
```

The command refuses to run from a dirty working tree because uncommitted files cannot be honestly attributed to an immutable source reference. A successful source-bound local release may assess `source_commit_artifacts_verified: true`, but it still records external trust and production release verification as `not_assessed`.

This is intentionally narrower than production trust. A source-bound local release does not prove protected GitHub OIDC identity, Rekor inclusion, customer PKI audit evidence, protected ref status, workflow identity, or release approval.

## Externally Trusted Production Release

An externally trusted production release requires the local source-bound proof plus an accepted external trust profile:

- `sigstore-rekor-keyless-v1`
- `customer-pki-private-equivalent-v1`

For `trusted_contract_release: true`, the release verification record must show successful values for source commit artifacts, manifest digest, artifact digest, signature, identity policy, source URI, protected ref, workflow identity, approval, transparency or customer audit evidence, freshness, and `production_release_verified`. `approval_status` covers the release captain and required approval refs from the trusted identity policy.

`local-dev-private-equivalent-v1` is deliberately rejected for `trusted_contract_release: true`. Local DSSE, schema validity, and private key possession are useful engineering evidence, but none of them is production release trust.

## Freshness and Rollback

Block 01 does not implement TUF. Rollback protection is limited to:

- `previous_manifest_digest`
- exactly one of `valid_until` or `freshness_policy`
- release verification freshness status

Delegated roles, threshold metadata, snapshot metadata, timestamp metadata, target metadata, and key rotation require a future block.
