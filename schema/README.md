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
| `flight-recorder-event.schema.json` | Records one ordered Block 09 recorder event with canonical hash fields, provenance, evidence, redaction, and optional witness reference. |
| `flight-recorder-run.schema.json` | Records run-level recorder metadata, source/task locks, event-chain closure, gaps, and profile state. |
| `flight-recorder-witness.schema.json` | Records a witness anchor that binds run id, source baseline, task hash, recorder version, and chain head. |
| `ci-witness.schema.json` | Records a CI witness that binds local report/run artifacts to repository, commit, workflow, job, and CI run identity. |
| `consumer-schema-version-declaration.schema.json` | Shows how downstream consumers declare supported schema versions. |
| `trace.schema.json` | Links specs, tasks, changes, evidence, observations, metric streams, external verdicts, accountability, and contract verification records. |
| `self-attestation-case.schema.json` | Defines local self-attestation verifier cases and expected proof states. |
| `evidence-bundle.schema.json` | Compatibility schema for reviewable proof bundles. |
| `gate-verdict.schema.json` | Compatibility schema for externally recorded gate results with explicit `origin: "external"`. |
| `decision-record.schema.json` | Compatibility schema for final human or external automated decisions. |

## Validation

Basic JSON syntax check:

```bash
jq empty schema/*.json
```

Schema validation target:

- JSON Schema Draft 2020-12
- stable `$id` per schema
- semver schema versions for artifacts once full examples are validated

Canonical validation commands for this path:

```bash
go test ./...
go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc
```

Validation commands exclude `.git/`, `.beads/`, `.sdp-trace-runs/`, `benchmarks/repos/`, temporary directories, editor caches, and generated dependency directories.

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

Flight-recorder schemas add Block 09 run evidence, not trust authority. A schema-valid local recorder chain can support local reconstruction only. Accountability, audit-grade, or external-trust claims require a verifier profile that checks witness evidence outside the mutable run artifact. Late-attach gaps remain `not_assessed`; requirement changes are represented by supersession events; unresolved redaction remains visible to verifier profiles and must not be hidden by summaries or query output.

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
