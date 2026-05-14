# Schema Reference

These schemas define the portable `sdp-trace` contract.

## Schemas

The canonical schema list and metadata live in `schema/index.json`. The table below is generated from that index and must stay synchronized.

| Schema | Status | Purpose |
|---|---|---|
| `accountability.schema.json` | current | Records human-held DRI, approver, escalation, risk owner, approval ref, and line of defense. |
| `adapter-capture-run.schema.json` | current | Records adapter capture run metadata and event bindings. |
| `adapter-event.schema.json` | current | Records a single adapter event with provenance and digest bindings. |
| `adapter-registry.schema.json` | current | Declares adapter capabilities, versions, and event mappings. |
| `assessment-input.schema.json` | current | Packages trace evidence for an external policy consumer. |
| `assessment-result.schema.json` | current | Records assessment outcome, profile state, and gap list for a verifier run. |
| `authority-envelope.schema.json` | current | Declares selected actor/task authority without embedding downstream policy consequences. |
| `authority-evaluation.schema.json` | current | Records authority evaluation facts, attribution state, source coverage, and replay evidence refs. |
| `change-evidence-packet.v0.schema.json` | current | Packages change evidence for reviewable proof bundles (v0). |
| `checkpoint-verification.schema.json` | current | Records live signed-checkpoint verification states for signature, payload digest, replay binding, sequence, signer authority, and freshness checks. |
| `ci-artifact-observation.schema.json` | current | Records CI artifact observation results with digest bindings. |
| `ci-witness.schema.json` | current | Records a CI witness that binds local report/run artifacts to repository, commit, workflow, job, and CI run identity. |
| `common.schema.json` | current | Reusable IDs, timestamps, digests, accountability, and not_assessed definitions. |
| `consumer-schema-version-declaration.schema.json` | current | Shows how downstream consumers declare supported schema versions. |
| `contract-manifest.schema.json` | current | Lists contract artifacts and SHA-256 digests for a contract release. |
| `contract-release-verification.schema.json` | current | Records manifest, source-bound artifact, signature, identity-policy, external trust, approval, freshness, and production verification status. |
| `cross-repo-posture-export.schema.json` | current | Exports posture evidence across repository boundaries. |
| `decision-record.schema.json` | historical | Compatibility schema for final human or external automated decisions. |
| `delivery-trace-envelope.schema.json` | current | Wraps delivery trace records with envelope metadata. |
| `evidence-binding.schema.json` | current | Records cross-source binding state for git, harness, gateway, artifact, and custom evidence links. |
| `evidence-bundle-manifest.v0.schema.json` | current | Lists artifacts and digests for an evidence bundle (v0). |
| `evidence-bundle.schema.json` | historical | Compatibility schema for reviewable proof bundles. |
| `evidence-event.schema.json` | current | Records one inspectable proof item. |
| `external-verdict-input.schema.json` | current | Records externally produced gate, readiness, override, or custom verdicts as external evidence. |
| `flight-recorder-event.schema.json` | current | Records one ordered recorder event with canonical hash fields, provenance, evidence, redaction, and optional witness reference. |
| `flight-recorder-run.schema.json` | current | Records run-level recorder metadata, source/task locks, event-chain closure, gaps, and profile state. |
| `flight-recorder-witness.schema.json` | current | Records a witness anchor that binds run id, source baseline, task hash, recorder version, and chain head. |
| `forensics-query-pack-result.schema.json` | current | Records read-only forensic query-pack rows, input artifact digests, source references, row evidence states, and output-safety assertions without a policy verdict. |
| `gate-result.schema.json` | current | Records version-separated advisory and protected gate facts, including selected profile, protected gate state, checkpoint verification summary, protected conditions, and next-action hints without native policy ownership. |
| `gate-verdict.schema.json` | historical | Compatibility schema for portable gate results, including cannot_verify/not_assessed rationale, evidence requirements, and external policy references. |
| `github-pr-evidence-input.v0.schema.json` | current | Packages GitHub PR evidence for assessment input (v0). |
| `harness-event.schema.json` | current | Records one imported harness lifecycle event with digest-bound source identity and content treatment state. |
| `harness-observation-profile.schema.json` | current | Declares required and optional external harness event families, retention policy, degradation rules, and profile limits. |
| `harness-observation-run.schema.json` | current | Records an observed harness export run produced by harness observe. |
| `harness-observation-validation.schema.json` | current | Records harness validate state, per-family dimensions, and non-authority boundary. |
| `harness-session-profile.schema.json` | current | Declares harness session requirements and event capture scope. |
| `interaction-event.schema.json` | current | Records one interaction event with actor, intent, and payload refs. |
| `interaction-trace.schema.json` | current | Links interaction events into a traceable sequence. |
| `managed-harness-policy.schema.json` | current | Declares managed harness policy constraints and authority. |
| `metric-stream.schema.json` | current | Records process movement across windows without interpretation labels. |
| `observation.schema.json` | current | Records evidence-backed observations without policy verdicts. |
| `observed-action.schema.json` | current | Records observed review, mutation, harness, gateway, and custom events with source evidence refs. |
| `operation-record.schema.json` | current | Records one operation with accountability and evidence refs. |
| `promise-record.schema.json` | current | Records a promise or commitment with expected evidence. |
| `proof-summary.schema.json` | current | Records live verifier output for a selected proof profile. Persisted proof summaries are audit artifacts, not trust authority until re-verified or externally signed. |
| `provenance-record.schema.json` | current | Records actor/model/harness/tool provenance and payload digests. |
| `pr-review-common.schema.json` | current | Shared types for PR review packets and ledgers. |
| `pr-review-ledger.schema.json` | current | Records PR review ledger entries with verdict state. |
| `pr-review-packet.schema.json` | current | Packages PR review facts for external policy consumers. |
| `pr-review-profile.schema.json` | current | Declares PR review profile requirements and checks. |
| `pr-review-result.schema.json` | current | Records PR review run set results and gap list. |
| `pr-review-validation.schema.json` | current | Records PR review validation state and schema conformance. |
| `redaction-policy.schema.json` | current | Declares redaction rules and suppressed field policies. |
| `repo-observer-status.schema.json` | current | Records repository observer health and coverage state. |
| `review-ledger.schema.json` | current | Records review ledger entries with accountability and verdict. |
| `risk-classification.schema.json` | current | Records observed autonomy/impact and externally declared oversight assertions. |
| `self-attestation-case.schema.json` | current | Defines local self-attestation verifier cases and expected proof states. |
| `signed-checkpoint.schema.json` | current | Records a detached-signature checkpoint binding run id, nonce, source snapshot, task hash, contract digest, event count, and chain head for replay-resistant verification. |
| `trace.schema.json` | historical | Links specs, tasks, changes, evidence, observations, metric streams, external verdicts, accountability, and contract verification records. |
| `trusted-checkpoint-policy.schema.json` | current | Declares allowed checkpoint signer identities and authority boundaries for local signed, CI signed, or external witnessed checkpoint evidence. |
| `trusted-identity-policy.schema.json` | current | Declares which signer identity may issue a trusted contract release. |
| `witness-profile-result.schema.json` | current | Records witness profile evaluation results and trust boundary. |

## Validation

Basic JSON syntax check:

```bash
jq empty schema/*.json
```

Schema documentation freshness check (detects missing, stale, or extra schema entries):

```bash
go run ./tools/schemadoc
```

Schema validation target:

- JSON Schema Draft 2020-12
- stable `$id` per schema
- semver schema versions for artifacts once full examples are validated

Canonical validation commands for this path:

```bash
go test ./...
sdp-trace validate-fixtures examples/agentic-sdlc
```

Validation commands exclude `.git/`, `.beads/`, `.sdp-trace-runs/`, `benchmarks/repos/`, temporary directories, editor caches, and generated dependency directories.

## Compatibility

Before `sdp-trace` v1.0, schema changes may be breaking only when examples and compatibility notes are updated in the same change.

`gate-result.schema.json` accepts both `block14-gate-result-v1` and
`block16-gate-result-v1`. Advisory artifacts do not require protected profile
fields; protected-profile artifacts require `selected_profile`,
`protected_gate`, `checkpoint_verification`, and `protected_conditions` so
readers can avoid inferring protected conclusions from older advisory output.
Those block-numbered values are legacy compatibility tokens, not the naming
pattern for future public contracts. New compatibility tokens should use
semantic profile names; remove the block-numbered aliases before v1.0 unless a
retained migration note says otherwise.

After v1.0:

- additive optional fields are minor-version changes
- required field removals, enum semantic changes, or ownership-boundary changes are major-version changes
- downstream policy consumers must declare supported schema versions

`schema/trace.schema.json` remains a compatibility path until a replacement path and migration note are committed.

## Ownership Boundary

`sdp-trace` records evidence, provenance, observations, metric movement, accountability, manifest integrity, and external verdict inputs.

`sdp-trace` does not decide pass/fail, readiness, degradation, threshold sufficiency, or override outcomes. Those policy decisions belong to CI, release governance, customer governance, or another external policy consumer.

External verdicts may be recorded only through `external-verdict-input.schema.json` with explicit `origin: "external"`.

Flight-recorder schemas add run evidence, not trust authority. A schema-valid local recorder chain can support local reconstruction only. Accountability, audit-grade, or external-trust claims require a verifier profile that checks witness evidence outside the mutable run artifact. Late-attach gaps remain `not_assessed`; requirement changes are represented by supersession events; unresolved redaction remains visible to verifier profiles and must not be hidden by summaries or query output.

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
