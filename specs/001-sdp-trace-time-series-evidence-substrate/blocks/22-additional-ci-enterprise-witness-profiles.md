# Block 22: Additional CI And Enterprise Witness Profiles

Status: implemented and PR-reviewed in PR #15. The reviewed implementation was
merged after code/correctness, tracing/evidence, and requirements-vs-
implementation review planes. External production trust remains outside Block
22.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13-product-gap-closure-roadmap.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/15-signed-checkpoint-replay-resistance.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/16-protected-gate-enforcement-profile.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/17-managed-harness-enforcement-profile.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/21-cross-repository-degradation-export.md`

## Goal

Define and implement additional witness profiles without making GitHub Actions
the hidden trust model for every CI or enterprise environment.

The product outcome is a provider-neutral witness profile contract plus a small
set of concrete profiles that share verifier semantics:

- GitLab CI;
- Buildkite as the second hosted CI profile for Block 22;
- customer-private PKI;
- air-gapped witness documentation and fixture semantics.

Buildkite is selected over Jenkins for this block because it has a clearer
hosted CI identity and artifact boundary for a first portable demo profile.
Jenkins remains a documented follow-up profile candidate because real
enterprise deployments often require it, but its plugin, controller, agent, and
credential topologies create more ways to overclaim witness independence.

## Problem

The current witness path is useful, but it is too easy to treat "CI witness" as
synonymous with "GitHub Actions shaped environment variables." That weakens the
product in two ways.

First, it creates UX confusion for platform teams. A GitLab or Buildkite user
should not have to reverse-engineer GitHub fields to understand what evidence
`sdp-trace` needs.

Second, it creates a trust bug. Environment variables alone do not prove a
witness boundary. Every supported profile must state where identity comes from,
what is signed or bound, how freshness is established, which artifacts are
covered, and which states remain unsupported, `not_assessed`, or
`cannot_verify`.

## Non-Goals

- No native policy decision, merge gate, release approval, compliance verdict,
  audit verdict, legal verdict, trust score, health score, risk score, badge,
  rank, or grade.
- No claim that GitLab, Buildkite, Jenkins, customer PKI, air-gapped operation,
  or any enterprise CI is fully supported beyond the profiles implemented in
  this block.
- No dependency on GitLab, Buildkite, Jenkins, Sigstore, Rekor, Vault, HSM,
  KMS, LDAP, SAML, OIDC provider SDKs, cloud SDKs, or customer infrastructure.
- No profile may upgrade trust scope from environment variables alone.
- No profile may treat local developer signatures, committed JSON, job logs, or
  unchecked artifact paths as external witness evidence.
- No network calls in the first implementation unless a later reviewed profile
  explicitly introduces them with deterministic `cannot_verify` behavior.
- No raw token, OIDC JWT, CI secret, private key, certificate private material,
  job log body, command body, provider URL with credentials, private filesystem
  path, or customer identity payload may be persisted or printed.
- No composite, chained, or layered witness profiles. Block 22 profiles are
  single-source only. Running multiple witness profiles and combining their
  results must not upgrade trust scope beyond the strongest single established
  scope; composite witnessing requires a later reviewed profile.
- Jenkins CI support is explicitly out of scope for Block 22. Jenkins follows
  as a separate reviewed profile candidate because its plugin, controller,
  agent, and credential topologies create more overclaim risks.
- No generic "enterprise witness" catch-all that bypasses profile-specific
  identity, signing, freshness, artifact-binding, and unsupported-state rules.

## Product Boundary

Block 22 may add:

- a provider-neutral witness profile schema or equivalent Go validation
  contract;
- a provider-neutral witness summary shape consumed by existing gate,
  protected, managed, and export surfaces without changing their policy
  ownership;
- profile-specific normalizers for GitLab CI and Buildkite;
- a customer-private PKI profile that validates declared public certificate or
  public key identity, signer authority policy, payload digest binding, and
  freshness metadata without requiring a live customer PKI;
- air-gapped profile documentation and fixtures that show what can and cannot
  be verified offline;
- deterministic fixture cases proving valid, missing, stale, mismatched,
  environment-only, unsupported, and malformed witness inputs.

Block 22 must preserve this vocabulary:

- `local_observed`: local structural evidence only;
- `ci_witnessed`: CI witness evidence bound to source, run, policy, and
  artifacts under an allowed profile;
- `external_witnessed`: external or enterprise witness evidence only when the
  selected profile has an independent timestamp, append-only, private PKI, or
  approved customer-equivalent anchor;
- `not_assessed`: profile or field was not assessed;
- `cannot_verify`: required profile evidence could not be verified;
- `fail`: evidence conflicts with the selected profile.

## Witness Profile Contract

Every profile must declare:

| Field | Purpose |
| --- | --- |
| `profile_id` | Stable closed profile id, for example `gitlab-ci-v1`, `buildkite-v1`, or `customer-pki-v1`. |
| `profile_version` | Profile contract version. Semantic changes require a new version. |
| `identity_source` | Verifier-readable source of witness identity, such as CI OIDC claims, CI job metadata plus trusted signer policy, or customer certificate subject. |
| `signing_boundary` | The process, key, or authority boundary that produced or authorized the witness. |
| `source_binding` | Source repository/ref/commit or tree digest facts covered by the witness, plus the non-pass state when they are absent or conflict. |
| `run_binding` | Run id, build id, pipeline id, job id, nonce, or sequence facts covered by the witness, plus the non-pass state when they are absent or conflict. |
| `policy_binding` | Authority policy, protected gate policy, managed harness policy, or checkpoint policy digest facts covered by the witness, plus the non-pass state when they are absent or conflict. |
| `freshness_boundary` | Timestamp, monotonic build/run identity, nonce, sequence, or signed checkpoint freshness fact required by the profile. |
| `artifact_binding` | Exact source, run, policy, checkpoint, report, and artifact digests covered by the witness. |
| `independence_state` | Closed witness topology state: `external_independent`, `ci_isolated_job`, `ci_shared_pipeline`, `ci_same_job`, `local_only`, `not_assessed`, or `cannot_verify`. This constrains trust scope but is not itself a trust scope. |
| `unsupported_states` | Closed reasons the profile refuses to assess or cannot support. |
| `safe_output_classes` | Sensitive classes verified absent from JSON and explain output, using the closed classes in Safety Requirements. |

The normalized witness result must expose:

- witness id, profile id/version, provider kind, producer, generated time, and
  schema version;
- trust scope requested by the input and trust scope actually established by
  verifier facts;
- identity state, signer authority state, freshness state, artifact binding
  state, source binding state, run binding state, policy binding state, and
  independence state;
- key custody state for enterprise signer profiles: `hsm`, `kms`, `software`,
  `unknown`, `not_assessed`, or `cannot_verify`;
- closed reason codes for every `fail`, `not_assessed`, or `cannot_verify`
  condition;
- path-redacted artifact ids and SHA-256-or-stronger digests when available;
- output-safety verification state.

### Closed Reason Codes

Profiles must use this initial closed reason-code registry. Adding a reason
code requires a reviewed profile version change.

| Reason code | Default verifier state | Meaning |
| --- | --- | --- |
| `witness_profile_verified` | `pass` | Selected witness profile established the requested trust scope from verifier facts. |
| `witness_identity_missing` | `cannot_verify` | Required witness identity was absent. |
| `witness_identity_mismatch` | `fail` | Witness identity conflicts with selected source, run, policy, or authority. |
| `witness_signer_authority_missing` | `cannot_verify` | Required signer authority policy was absent. |
| `witness_signer_mismatch` | `fail` | Signer identity is not allowed by policy. |
| `witness_freshness_missing` | `cannot_verify` | Required timestamp, nonce, sequence, or freshness evidence was absent. |
| `witness_freshness_stale` | `fail` | Freshness evidence is outside policy or superseded by later bound evidence. |
| `witness_artifact_digest_missing` | `cannot_verify` | Required artifact digest was absent. |
| `witness_artifact_digest_mismatch` | `fail` | Artifact digest conflicts with selected artifact bytes or manifest. |
| `witness_source_binding_missing` | `cannot_verify` | Required source ref, commit, tree, or baseline binding was absent. |
| `witness_source_mismatch` | `fail` | Source binding conflicts with the run under inspection. |
| `witness_run_binding_missing` | `cannot_verify` | Required run id, build id, pipeline id, job id, nonce, or sequence binding was absent. |
| `witness_run_mismatch` | `fail` | Embedded run/build/pipeline/job identity does not match the run under inspection. |
| `witness_policy_binding_missing` | `cannot_verify` | Required authority, checkpoint, protected, or managed policy digest binding was absent. |
| `witness_policy_mismatch` | `fail` | Policy binding conflicts with selected policy input. |
| `witness_environment_only_insufficient` | `cannot_verify` | Environment variables or agent-reported metadata were present but no authority-bound witness envelope existed. |
| `witness_unsupported_profile` | `cannot_verify` | Profile id or version is not supported by this verifier. |
| `witness_malformed_input` | `cannot_verify` | Input cannot be parsed as the selected profile. |
| `witness_unsafe_output_candidate` | `fail` | Serialized JSON or explain output would expose a forbidden safety class. |
| `witness_private_key_input_rejected` | `fail` | Input path or payload appears to contain private key material. |
| `witness_revocation_not_assessed` | `not_assessed` | Certificate revocation could not be checked in the selected offline profile. |
| `witness_certificate_revoked` | `fail` | Customer PKI certificate is revoked according to provided revocation evidence. |
| `witness_key_custody_not_assessed` | `not_assessed` | Key custody is absent from authority policy or cannot be verified. |

### Trust-Scope Determination

The established trust scope is derived from verifier facts. It is not selected
by profile name, CI provider name, or environment shape.

| Established scope | Required facts |
| --- | --- |
| `external_witnessed` | Identity, signer authority, source binding, run binding, policy binding, artifact binding, and freshness are `pass`; `independence_state` is `external_independent`; selected profile explicitly permits external witness scope. |
| `ci_witnessed` | Identity, signer authority or CI authority, source binding, run binding, policy binding, artifact binding, and freshness are `pass`; `independence_state` is `ci_isolated_job` or stronger; selected profile permits CI witness scope. |
| `local_observed` | Local structural facts are present, but CI or external authority, freshness, or independence requirements are absent. |
| `cannot_verify` | A required fact for the selected profile is missing, malformed, unsupported, or inaccessible and no conflicting evidence was observed. |
| `not_assessed` | The profile or field was intentionally outside the selected assessment scope. |
| `fail` | Any required binding, identity, signer, digest, freshness, revocation, policy, run, source, or safety fact conflicts with selected inputs. |

Environment variables, job logs, local signatures, committed JSON, unchecked
artifact paths, process-inherited CI variables, and agent-reported metadata can
never establish `ci_witnessed` or `external_witnessed` by themselves.

### Freshness Evaluation

Freshness semantics are shared across profiles:

- `pass`: timestamp, nonce, sequence, run/build identity, or signed checkpoint
  freshness evidence is bound to the selected run and policy, falls within the
  selected freshness policy window, and is not superseded by later bound
  evidence;
- `fail`: freshness evidence is expired, stale, superseded, contradicts the
  run/build identity, or uses a weak or mismatched digest;
- `cannot_verify`: required freshness evidence is missing, malformed, unsafe,
  or not readable under the selected profile;
- `not_assessed`: freshness is explicitly out of scope for a documentation-only
  or offline fixture and the resulting trust scope does not rely on it.

Customer PKI freshness must come from an explicit signed freshness evidence
artifact that binds payload digest, run id, policy digest, signer identity,
issued time, optional valid-until time, and nonce or sequence. A self-claimed
timestamp in unsigned witness JSON is not authority and must not upgrade trust.

### Cross-Surface Consumption

Downstream surfaces consume normalized witness facts; they do not own policy
decisions for Block 22.

| Normalized field | Gate / protected gate | Managed harness | Cross-repository posture export |
| --- | --- | --- | --- |
| `established_trust_scope` | Input fact for witness gate state and trust cap. | Input fact for managed witness binding state. | Witness-scope posture signal. |
| `profile_id` / `profile_version` | Recorded for deterministic reason output. | Recorded for managed profile diagnostics. | Grouping or safe metadata dimension when selected. |
| `identity_state` / `signer_authority_state` | Used to determine witness binding pass/fail/cannot_verify facts. | Used when managed policy requires signer or adapter authority. | Aggregated only as safe closed state counts. |
| `source_binding_state` / `run_binding_state` / `policy_binding_state` | Used as verifier facts for protected profile non-pass states. | Used to bind managed witness to selected run and policy. | Aggregated as closed state counts or refusal reasons. |
| `freshness_state` / `artifact_binding_state` | Used for stale, missing, or mismatched witness facts. | Used for managed witness freshness and artifact binding. | Aggregated as stale, untrusted, or not-assessed input posture. |
| `independence_state` | Constrains maximum trust scope. | Constrains managed witness binding confidence. | Exported as safe closed metadata only. |
| `reason_codes` | Rendered as deterministic explain/preview reasons. | Rendered as deterministic assessment reasons. | Counted only as closed reason-code facts. |

## Profile Requirements

### GitLab CI

GitLab CI profile `gitlab-ci-v1` must distinguish:

- malicious runner, compromised `CI_JOB_TOKEN`, or injected `.gitlab-ci.yml`
  environment-variable threats; environment-only inputs cannot establish trust;
- trusted CI identity available through reviewed profile inputs;
- project/ref/commit/pipeline/job identity bound to the run and source
  baseline;
- artifact digest binding to the selected `sdp-trace` report or checkpoint;
- same-job or same-runner co-located witness and build topologies that cap
  trust below `external_witnessed`, and below `ci_witnessed` when the witness
  lacks a separate job isolation boundary;
- missing or unverified OIDC/token identity as `cannot_verify` or
  `not_assessed`, not pass;
- environment-only metadata as insufficient for a trust upgrade.

### Buildkite

Buildkite profile `buildkite-v1` must distinguish:

- pipeline `env` blocks and agent hooks that can inject secrets into
  agent-reported metadata;
- organization/pipeline/build/job identity;
- agent-reported metadata versus signed or authority-bound witness facts;
- artifact digest binding;
- co-located agent and witness topologies that cap trust below
  `external_witnessed`;
- absent independent signer or missing artifact digest as `cannot_verify`.

The first `buildkite-v1` implementation caps at `ci_witnessed`; it cannot
establish `external_witnessed`. A future Buildkite external profile would need a
separate external witness anchor or customer PKI attestation bound to the
Buildkite run identity.

### Customer PKI

Customer-private PKI profile `customer-pki-v1` must distinguish:

- declared public certificate or public key identity from private key material;
- authority policy that allows the signer and profile;
- payload digest and checkpoint binding;
- signed freshness evidence declared through the explicit freshness artifact;
- certificate revocation state: provided revocation evidence can produce
  `pass` or `fail`; unavailable CRL/OCSP or equivalent revocation evidence is
  `not_assessed` unless the selected profile requires it;
- key custody state declared by authority policy as `hsm`, `kms`, `software`,
  `unknown`, `not_assessed`, or `cannot_verify`;
- unverifiable chain, missing authority, expired validity, weak digest, or
  missing freshness as deterministic non-pass states.

This profile may support `external_witnessed` only when the verifier can bind
the witness to an allowed signer authority and a freshness source outside the
mutable run artifact. Otherwise it must cap at `ci_witnessed`,
`local_observed`, `cannot_verify`, or `not_assessed` according to the evidence.

### Air-Gapped Profile

The air-gapped profile is documentation plus fixtures in Block 22.
`air-gapped-v1` is a fixture and documentation profile id, not a
`witness --kind` command. A reviewer validates committed air-gapped fixtures
through the repository fixture validation command, not through live network or
CI discovery.

It must explain:

- which evidence can be verified offline;
- which external checks are impossible in the air-gapped environment;
- how public-key, timestamp, and artifact-digest evidence can be carried into
  the environment;
- which states remain `not_assessed` or `cannot_verify`;
- why committed witness JSON is not authority by itself.
- network calls are explicitly forbidden in any air-gapped validator; any
  verification requiring external network access must emit `cannot_verify`;
- manually imported public keys, timestamps, and revocation snapshots must
  carry integrity digests verified against an out-of-band source, with
  verification failure emitting `fail`;
- freshness checks must use carried evidence internal monotonic sequence or
  signed timestamp chain, not comparison to an unsynchronized local wall clock,
  unless wall-clock synchronization is itself verified and declared.

## CLI Boundary

The existing witness command should remain the user-facing entrypoint unless
implementation review proves a new command is necessary:

```bash
go run ./cmd/sdp-trace witness --kind <profile-kind> --out <file> [--report-dir <dir>] <runs-root-or-run-dir>
```

Allowed Block 22 `--kind` values should be explicit closed identifiers. The
spec direction is:

- `github-actions` remains supported;
- `gitlab-ci` is added;
- `buildkite` is added;
- `customer-pki` is added only when required public inputs are supplied;
- air-gapped guidance is documented and fixture-backed, not a magic catch-all
  command kind.

Unknown `--kind` values are usage errors with exit code `2` and a deterministic
list of allowed values. Multiple `--kind` values are not allowed in Block 22.
Profile inputs that parse but cannot establish required profile facts emit
`cannot_verify` with exit code `3`; conflicting evidence emits `fail` with exit
code `1`; established witness facts emit exit code `0`.

Profile normalizers must read CI witness inputs only from the declared
run-directory witness envelope or explicit input flags. They must ignore
environment variables inherited from the verifier process. If matching CI
environment variables are present in the verifier process but no input envelope
exists, the result is `cannot_verify` with
`witness_environment_only_insufficient`.

Customer PKI requires explicit public input flags:

- `--customer-pki-authority-policy <path>` for allowed signer, profile,
  key-custody, revocation requirement, and freshness policy;
- `--customer-pki-public-cert <path>` or `--customer-pki-public-key <path>` for
  declared public identity;
- `--customer-pki-payload-digest <sha256-or-stronger-digest>` for the signed
  payload;
- `--customer-pki-freshness-evidence <path>` for signed freshness evidence.

The command must refuse implicit directory scanning, runtime-relative discovery,
private key paths, private key payloads, provider tokens, or customer
directories.

## Fixture Matrix

Block 22 fixtures must cover at least:

- valid GitLab CI witness;
- GitLab environment-only input that cannot upgrade trust;
- GitLab missing identity;
- GitLab source commit mismatch;
- GitLab stale or missing freshness evidence;
- GitLab artifact digest mismatch;
- GitLab same-job or same-runner topology cap;
- valid Buildkite witness;
- Buildkite same-job or agent-reported-only topology capped below
  `external_witnessed`;
- Buildkite stale freshness;
- Buildkite missing independent signer;
- Buildkite artifact digest mismatch;
- CI witness run/build/pipeline id mismatch;
- valid customer PKI witness with public certificate/key identity and allowed
  signer policy;
- customer PKI signer mismatch;
- customer PKI expired or stale freshness evidence;
- customer PKI revocation unavailable, revoked certificate, and key custody
  not assessed;
- customer PKI weak or missing digest;
- air-gapped offline package with explicit `not_assessed` external checks;
- malformed profile id, unsupported profile version, unsafe artifact ref, and
  unsafe output candidate.

Each fixture must state expected verifier result, trust scope, reason codes,
identity state, signing boundary state, freshness state, artifact binding state,
and output-safety state.

Minimum expected fixture rows:

| Fixture | Profile | Expected verifier state | Established scope | Required reason code when non-pass |
| --- | --- | --- | --- | --- |
| `gitlab-valid-ci-witness` | `gitlab-ci-v1` | `pass` | `ci_witnessed` | `witness_profile_verified` |
| `gitlab-env-only-cannot-verify` | `gitlab-ci-v1` | `cannot_verify` | `cannot_verify` | `witness_environment_only_insufficient` |
| `gitlab-missing-identity` | `gitlab-ci-v1` | `cannot_verify` | `cannot_verify` | `witness_identity_missing` |
| `gitlab-source-mismatch` | `gitlab-ci-v1` | `fail` | `fail` | `witness_source_mismatch` |
| `gitlab-stale-freshness` | `gitlab-ci-v1` | `fail` | `fail` | `witness_freshness_stale` |
| `gitlab-artifact-digest-mismatch` | `gitlab-ci-v1` | `fail` | `fail` | `witness_artifact_digest_mismatch` |
| `gitlab-same-job-topology-cap` | `gitlab-ci-v1` | `cannot_verify` | `cannot_verify` | `witness_environment_only_insufficient` |
| `buildkite-valid-ci-witness` | `buildkite-v1` | `pass` | `ci_witnessed` | `witness_profile_verified` |
| `buildkite-same-job-cap` | `buildkite-v1` | `cannot_verify` | `cannot_verify` | `witness_environment_only_insufficient` |
| `buildkite-missing-independent-signer` | `buildkite-v1` | `cannot_verify` | `cannot_verify` | `witness_signer_authority_missing` |
| `buildkite-stale-freshness` | `buildkite-v1` | `fail` | `fail` | `witness_freshness_stale` |
| `buildkite-artifact-digest-mismatch` | `buildkite-v1` | `fail` | `fail` | `witness_artifact_digest_mismatch` |
| `ci-run-id-mismatch` | `gitlab-ci-v1` or `buildkite-v1` | `fail` | `fail` | `witness_run_mismatch` |
| `customer-pki-valid-external-witness` | `customer-pki-v1` | `pass` | `external_witnessed` | `witness_profile_verified` |
| `customer-pki-signer-mismatch` | `customer-pki-v1` | `fail` | `fail` | `witness_signer_mismatch` |
| `customer-pki-expired-freshness` | `customer-pki-v1` | `fail` | `fail` | `witness_freshness_stale` |
| `customer-pki-revocation-not-assessed` | `customer-pki-v1` | `not_assessed` | `not_assessed` | `witness_revocation_not_assessed` |
| `customer-pki-revoked-certificate` | `customer-pki-v1` | `fail` | `fail` | `witness_certificate_revoked` |
| `customer-pki-private-key-rejected` | `customer-pki-v1` | `fail` | `fail` | `witness_private_key_input_rejected` |
| `air-gapped-offline-not-assessed` | `air-gapped-v1` | `not_assessed` | `not_assessed` | `witness_revocation_not_assessed` or profile-specific missing external check reason |
| `unsupported-profile-version` | any | `cannot_verify` | `cannot_verify` | `witness_unsupported_profile` |
| `unsafe-output-candidate` | any | `fail` | `fail` | `witness_unsafe_output_candidate` |

## Safety Requirements

Block 22 is safety-sensitive because CI and enterprise witness inputs can carry
tokens, private keys, provider URLs, job logs, and customer identity data.

JSON and explain output must avoid:

- raw command arguments;
- command names or executable paths when unsafe;
- stdout/stderr bodies;
- prompt or model response bodies;
- OIDC request tokens or JWT bodies;
- CI secrets and masked secret values;
- private keys and certificate private material;
- provider API tokens;
- authenticated provider URLs;
- private filesystem paths;
- raw job logs;
- unsafe personal identifiers;
- free-text parser errors containing input content;
- customer directory, LDAP, SAML, cloud, Vault, HSM, KMS, or PKI payloads.

Negative leak assertions must use synthetic sentinel values and must not echo
candidate secrets in failure output.

Profile normalizers must validate or sanitize raw provider inputs against the
closed safety classes before structural parsing. If a secret-like pattern is
detected in a field, the normalizer must reject the field or replace it with a
redaction digest; final-output redaction alone is insufficient.

OIDC JWTs or CI tokens may be parsed only to extract reviewed identity claims
required by the selected profile. Raw JWT bodies, token bodies, and secret
values must not be persisted in normalized results, intermediate artifacts,
debug logs, panic messages, or explain output.

## Acceptance Criteria

- A reviewer can inspect a single provider-neutral witness profile contract and
  see the same identity, signing, freshness, artifact-binding, independence,
  unsupported-state, and output-safety semantics used by GitHub Actions,
  GitLab CI, Buildkite, customer PKI, and air-gapped guidance.
- `gitlab-ci-v1`, `buildkite-v1`, and `customer-pki-v1` profile fixtures share
  verifier states and reason-code semantics instead of each inventing a local
  pass/fail vocabulary.
- No fixture or command can upgrade to `ci_witnessed` or `external_witnessed`
  from environment variables alone.
- Customer PKI uses public identity and authority-policy inputs only; private
  key material is never required, read, printed, or persisted.
- Air-gapped documentation clearly marks external checks that remain
  `not_assessed` or `cannot_verify`.
- Gate, protected gate, managed harness, and cross-repository posture surfaces
  can consume the normalized witness result without becoming policy owners.
- Safety-sensitive tests prove witness JSON and explain output do not leak CI
  tokens, OIDC tokens, private key material, provider URLs with credentials,
  raw job logs, private paths, or customer identity payloads.
- Go-first verification, schema validation, fixture validation, drift checks,
  and separate code, tracing/evidence, and requirements-vs-implementation
  reviews pass before implementation closure.

## Implementation Slices After Approval

Implementation must not start until this reviewed spec direction is approved.
After approval, likely slices are:

1. Provider-neutral witness profile/result contract and schema alignment.
2. GitLab CI profile normalization and fixtures.
3. Buildkite profile normalization and fixtures.
4. Customer PKI profile validation with public identity inputs and fixtures.
5. Air-gapped docs, fixture matrix, CLI docs, and cross-surface drift updates.
6. Safety-sensitive negative tests and output-safety checks.

Each slice needs focused verification, review, and scoped commit. PR-level
review must repeat code/correctness, tracing/evidence, and
requirements-vs-implementation planes.
