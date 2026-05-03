# Schema Reference

These schemas define the portable `sdp-trace` contract.

## Schemas

| Schema | Purpose |
|---|---|
| `common.schema.json` | Reusable IDs, timestamps, digests, accountability, and `not_assessed` definitions. |
| `accountability.schema.json` | Records human-held DRI, approver, escalation, risk owner, approval ref, and line of defense. |
| `risk-classification.schema.json` | Records observed autonomy/impact and externally declared oversight assertions. |
| `contract-manifest.schema.json` | Lists contract artifacts and SHA-256 digests for a contract release. |
| `contract-release-verification.schema.json` | Records manifest, source-bound artifact, signature, identity-policy, external trust, approval, freshness, and production verification status. |
| `proof-summary.schema.json` | Records live verifier output for a selected proof profile. Persisted proof summaries are audit artifacts, not trust authority until re-verified or externally signed. |
| `trusted-identity-policy.schema.json` | Declares which signer identity may issue a trusted contract release. |
| `evidence-event.schema.json` | Records one inspectable proof item. |
| `provenance-record.schema.json` | Records actor/model/harness/tool provenance and payload digests. |
| `observation.schema.json` | Records evidence-backed observations without policy verdicts. |
| `metric-stream.schema.json` | Records process movement across windows without interpretation labels. |
| `external-verdict-input.schema.json` | Records externally produced verdicts or quality assertions as external evidence. |
| `assessment-input.schema.json` | Packages trace evidence for a policy consumer such as `sdp-gate`. |
| `consumer-schema-version-declaration.schema.json` | Shows how downstream consumers declare supported schema versions. |
| `trace.schema.json` | Links specs, tasks, changes, evidence, observations, metric streams, external verdicts, accountability, and contract verification records. |
| `self-attestation-case.schema.json` | Defines local self-attestation verifier cases and expected proof states. |
| `evidence-bundle.schema.json` | Legacy compatibility schema for reviewable proof bundles. |
| `gate-verdict.schema.json` | Legacy compatibility schema for externally recorded gate results with explicit `origin: "external"`. |
| `decision-record.schema.json` | Legacy compatibility schema for final human or external automated decisions. |

## Validation

Basic JSON syntax check:

```bash
jq empty schema/*.json
```

Schema validation target:

- JSON Schema Draft 2020-12
- stable `$id` per schema
- semver schema versions for artifacts once full examples are validated

Pinned validation command for current examples:

```bash
npm ci
npm run verify:baseline
```

`npm ci` may use network during setup. `npm run verify:baseline` uses the local pinned `ajv@8.20.0` dependency through `scripts/validate-json-schema.mjs` and must not fetch dependencies during validation.

`npm run verify` and `npm run validate` are compatibility aliases for `npm run verify:baseline`. The canonical baseline command remains `npm run verify:baseline`.

Dirty checkouts fail closed by default: `npm run verify:baseline` exits `3` with `result: "cannot_verify"`. `scripts/verify.sh --profile baseline --allow-dirty` is only for local structural development and emits `trust_scope: "local_dirty_structural_only"`. Checked-in proof-summary examples use `artifact_role: "untrusted_shape_example"` and `trust_scope: "untrusted_shape_only"` with fixture digests, not live proof digests.

The historical full contract validation script remains available as:

```bash
npm run validate:contracts
```

It may fail while optional proof slices or stale release artifacts are intentionally incomplete. The baseline verifier is the default fresh-checkout structural proof command.

Validation commands exclude `.git/`, `.beads/`, `.sdp-trace-runs/`, `benchmarks/repos/`, temporary directories, editor caches, and generated dependency directories.

Artifact safety scan:

```bash
scripts/check-artifact-safety.sh
```

Verified hash recomputation:

```bash
scripts/verify-artifact-hashes.sh
```

Any JSON example that claims `integrity_status: "verified_hash"` for a local artifact must match the current SHA-256 digest. Provenance records whose `command` is a local file path and `digest_algorithm: "sha256"` are checked the same way.

Contract manifest digest verification:

```bash
scripts/verify-contract-manifest.sh examples/contract-foundation/contract-manifest.example.json
```

Contract release verification evidence:

```bash
scripts/verify-release-signature.sh examples/contract-foundation/contract-release.dsse.json examples/contract-foundation/local-dev-signing-public.pem
```

Source-bound local release finalization guard:

```bash
scripts/finalize-source-bound-release.sh --manifest examples/contract-foundation/contract-manifest.example.json --source-ref HEAD
```

This command must refuse dirty working trees. It can assess local source-content proof only when the selected source reference contains every manifest artifact path with a matching SHA-256 digest.

## Compatibility

Before `sdp-trace` v1.0, schema changes may be breaking only when examples and compatibility notes are updated in the same change.

After v1.0:

- additive optional fields are minor-version changes
- required field removals, enum semantic changes, or ownership-boundary changes are major-version changes
- downstream consumers such as `sdp-gate` must declare supported schema versions

`schema/trace.schema.json` remains a compatibility path until a replacement path and migration note are committed.

## Ownership Boundary

`sdp-trace` records evidence, provenance, observations, metric movement, accountability, manifest integrity, and external verdict inputs.

`sdp-trace` does not decide pass/fail, readiness, degradation, threshold sufficiency, or override outcomes. Those policy decisions belong to `sdp-gate` or another external policy consumer.

External verdicts may be recorded only through `external-verdict-input.schema.json` with explicit `origin: "external"`.

## Accountability and Signing

AI actors may produce or review artifacts, but accountable identities must be `human_user`, `human_role`, or `human_group`.

Trusted contract release status requires:

- manifest and artifact digests matched
- signature or approved private-equivalent verification
- trusted identity policy matched
- freshness status current
- source-content verification against the selected source reference
- source URI, protected ref, workflow identity, and approval status matched
- external attestation and transparency evidence, or a customer-approved audit/timestamp substitute, recorded as assessed
- `production_release_verified` assessed true

Schema-valid artifacts are not automatically trusted contract releases.

Profile naming is split deliberately:

- `signature_profile` records the envelope/signature shape, currently `sdp-trace-signature/sigstore-dsse-keyless-v1`.
- `external_trust_profile` records the release trust anchor. Production-capable values are `sigstore-rekor-keyless-v1` and `customer-pki-private-equivalent-v1`; `local-dev-private-equivalent-v1` is only a local development placeholder and cannot support `trusted_contract_release: true`.
- `signer_identity` may include a signer implementation label such as `local-dev-dsse-openssl-v1`.

The local development placeholder `local-dev-private-equivalent-v1` is valid evidence for source-bound engineering checks, but it is rejected for `trusted_contract_release: true`.
