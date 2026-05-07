# Block 22 Review Ledger

Status: Spec and implementation review assessed. Initial product-boundary,
tracing/evidence, and enterprise/security spec planes returned `REVISE`.
Valid critical and major spec findings were accepted and fixed. Focused spec
re-reviews returned `APPROVE` with no remaining critical or major findings.
Implementation code/correctness, tracing/evidence, and security
requirements-vs-implementation planes returned `APPROVE` with no remaining
critical or major findings after accepted fixes. PR-level code/correctness,
tracing/evidence, and requirements-vs-implementation security planes were run
for PR #15. Accepted PR-level security blockers were fixed and focused
re-review returned `APPROVE`; GitHub Actions `verify` passed for the first PR
head and must be re-checked after the final evidence commit.

## Spec Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| S22-PB-01 | critical | product-boundary | Customer PKI freshness without a live PKI service was underspecified and could accept self-claimed freshness. | Accepted and fixed. Added shared Freshness Evaluation and required signed customer PKI freshness evidence bound to payload digest, run id, policy digest, signer identity, issued time, optional valid-until time, and nonce or sequence. | MiniMax-M2.7 product-boundary review; `22-additional-ci-enterprise-witness-profiles.md` Freshness Evaluation |
| S22-PB-02 | major | product-boundary | Air-gapped profile left reviewers unable to answer what command or validation path can be run. | Accepted and fixed. `air-gapped-v1` is now a documentation and fixture profile id, not a witness `--kind`; reviewer validation uses repository fixture validation, and network calls are forbidden. | MiniMax-M2.7 product-boundary review; Air-Gapped Profile |
| S22-PB-03 | major | product-boundary | Independence thresholds were abstract and profile implementers could invent trust ceilings. | Accepted and fixed. Added closed `independence_state` topology enum and Trust-Scope Determination matrix; GitLab and Buildkite profile sections now state topology caps. | MiniMax-M2.7 product-boundary review; Witness Profile Contract; Trust-Scope Determination |
| S22-PB-04 | major | product-boundary | Buildkite trust ceiling did not state whether `external_witnessed` was achievable. | Accepted and fixed. `buildkite-v1` now caps at `ci_witnessed`; external witnessing requires a later reviewed profile with an independent anchor. | MiniMax-M2.7 product-boundary review; Buildkite |
| S22-PB-05 | major | product-boundary | CLI error handling for unknown kind, missing profile inputs, and exit states was undefined. | Accepted and fixed. CLI Boundary now defines closed kinds, single-kind behavior, usage error exit `2`, fail exit `1`, cannot-verify exit `3`, and success exit `0`. | MiniMax-M2.7 product-boundary review; CLI Boundary |
| S22-PB-06 | major | product-boundary | Customer PKI public input flags and safe path handling were unnamed. | Accepted and fixed. CLI Boundary now names authority policy, public cert/public key, payload digest, and freshness evidence inputs and rejects implicit scanning, private-key inputs, provider tokens, and customer directories. | MiniMax-M2.7 product-boundary review; CLI Boundary |
| S22-TE-01 | critical | tracing/evidence | Shared closed reason-code registry was missing. | Accepted and fixed. Added Closed Reason Codes registry with default verifier states for identity, signer, freshness, artifact, source, run, policy, environment-only, unsupported, malformed, unsafe, private-key, revocation, and key-custody cases. | OpenRouter Qwen tracing/evidence review; Closed Reason Codes |
| S22-TE-02 | critical | tracing/evidence | Witness CLI exit codes were undefined. | Accepted and fixed with S22-PB-05. | OpenRouter Qwen tracing/evidence review; CLI Boundary |
| S22-TE-03 | critical | tracing/evidence | Trust-scope determination rules were missing. | Accepted and fixed. Added Trust-Scope Determination matrix mapping required facts and independence states to `external_witnessed`, `ci_witnessed`, `local_observed`, `cannot_verify`, `not_assessed`, and `fail`. | OpenRouter Qwen tracing/evidence review; Trust-Scope Determination |
| S22-TE-04 | major | tracing/evidence | Source, run, and policy binding states appeared in normalized result but not in profile contract. | Accepted and fixed. Added `source_binding`, `run_binding`, and `policy_binding` to the witness profile contract. | OpenRouter Qwen tracing/evidence review; Witness Profile Contract |
| S22-TE-05 | major | tracing/evidence | Fixture matrix did not enumerate expected outputs per fixture. | Accepted and fixed. Added minimum expected fixture rows with profile id, verifier state, established scope, and required non-pass reason code. | OpenRouter Qwen tracing/evidence review; Fixture Matrix |
| S22-TE-06 | major | tracing/evidence | Cross-surface consumption mapping was underspecified. | Accepted and fixed. Added mapping from normalized witness result fields to gate/protected, managed harness, and cross-repository posture export consumption. | OpenRouter Qwen tracing/evidence review; Cross-Surface Consumption |
| S22-ES-01 | critical | enterprise/security | GitLab CI omitted co-located agent/witness topology caps. | Accepted and fixed. GitLab profile now names runner/job-token/env-injection threats and caps same-job or same-runner topologies below external witness, and below CI witness without separate job isolation. | Kimi K2P6 enterprise/security review; GitLab CI |
| S22-ES-02 | critical | enterprise/security | Profile normalizers could inherit verifier process environment as CI evidence. | Accepted and fixed. CLI Boundary now requires profile normalizers to read only declared witness envelopes or explicit flags and ignore process-inherited CI env vars. | Kimi K2P6 enterprise/security review; CLI Boundary |
| S22-ES-03 | major | enterprise/security | CI run-id replay resistance was not explicit. | Accepted and fixed. Closed reason codes and fixtures now include `witness_run_mismatch`; contract requires run/build/pipeline/job binding. | Kimi K2P6 enterprise/security review; Witness Profile Contract; Fixture Matrix |
| S22-ES-04 | major | enterprise/security | Customer PKI lacked revocation and key-custody verifier states. | Accepted and fixed. Added revocation and key-custody states to customer PKI profile, normalized result, and reason-code registry. | Kimi K2P6 enterprise/security review; Customer PKI |
| S22-ES-05 | major | enterprise/security | Air-gapped profile lacked network-call prohibition and import-integrity rules. | Accepted and fixed. Air-gapped profile now forbids network calls and requires integrity digests for manually imported public keys, timestamps, and revocation snapshots. | Kimi K2P6 enterprise/security review; Air-Gapped Profile |
| S22-ES-06 | major | enterprise/security | Input sanitization boundary before profile normalization was missing. | Accepted and fixed. Safety Requirements now require validation or sanitization before structural parsing and forbid relying only on final-output redaction. | Kimi K2P6 enterprise/security review; Safety Requirements |
| S22-ES-07 | major | enterprise/security | Composite witness profiles were not explicitly ruled out. | Accepted and fixed. Non-Goals now rule out composite, chained, or layered witness profiles for Block 22. | Kimi K2P6 enterprise/security review; Non-Goals |

## Implementation Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| S22-IMPL-01 | major | security requirements-vs-implementation | Customer PKI signature verification, run-id replay resistance, unsafe input path handling, and fixture coverage needed hardening before the implementation could support the reviewed spec. | Accepted and fixed. Customer PKI verifies signed freshness evidence against public authority, payload digest, policy digest, signer id, freshness window, and current run id; unsafe paths reject traversal, URLs, private-key names, and symlink inputs; fixture/test coverage includes replay, signer, revocation, weak digest, private-key, and unsafe-output cases. | Kimi K2P6 implementation security review; `internal/witness/profiles.go`; `internal/witness/profiles_test.go` |
| S22-IMPL-02 | major | security requirements-vs-implementation | Input/output sanitization did not cover all closed safety classes, JWT-shaped bodies, and common provider-token markers. | Accepted and fixed. `containsSecretLike` and `forbiddenOutputPresent` now cover JWT-like three-part tokens, bearer/provider token markers, `cloud_payload`, `pki_payload`, and `free_text_parser_error_with_input`; tests prove JWT-shaped input is rejected before parsing and JWT-shaped serialized output is replaced with a safe `witness_unsafe_output` record. Focused re-review returned `APPROVE`. | Kimi K2P6 focused implementation security re-review; `internal/witness/profiles.go`; `internal/witness/profiles_test.go` |
| S22-IMPL-03 | minor | code/correctness | Initial implementation review packets omitted or truncated untracked files, so they were insufficient evidence. | Replaced. Full-patch code/correctness and tracing/evidence re-reviews were run against tracked and untracked Block 22 files and returned `APPROVE` with no remaining critical or major findings. | Qwen code/correctness re-review; MiniMax tracing/evidence re-review; `/tmp/block22-implementation-full.patch` generation pattern |

## PR-Level Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| S22-PR-01 | major | requirements-vs-implementation security | `strongDigest` accepted only 64-character SHA-256 digests, contradicting the reviewed `sha256-or-stronger` digest contract. | Accepted and fixed. `strongDigest` now accepts valid hex digests of at least 64 characters and tests cover 64-, 96-, and 128-character hex. Focused Kimi K2P6 PR re-review returned `APPROVE`. | Kimi K2P6 PR security review and focused re-review; `internal/witness/profiles.go`; `internal/witness/profiles_test.go` |
| S22-PR-02 | major | requirements-vs-implementation security | `finalizeRecordForWrite` copied potentially tainted `profile_id`, `profile_version`, and `provider_kind` into the safe fallback record after detecting unsafe output. | Accepted and fixed. The safe fallback now derives these fields only from `baseRecord(record.Kind)`, and the JWT-shaped output test poisons `ProfileID` to prove the sentinel is not persisted. Focused Kimi K2P6 PR re-review returned `APPROVE`. | Kimi K2P6 PR security review and focused re-review; `internal/witness/profiles.go`; `internal/witness/profiles_test.go` |
| S22-PR-03 | major | requirements-vs-implementation security | The existing GitHub Actions witness path did not populate the new schema-required normalized fields. | Accepted and fixed. `BuildGitHubActionsWithFetcher` now uses the normalized `baseRecord`, populates `profile_states`, `established_trust_scope`, `reason_codes`, and `output_safety`, and `WriteGitHubActions` finalizes output before writing. Focused Kimi K2P6 PR re-review returned `APPROVE`. | Kimi K2P6 PR security review and focused re-review; `internal/witness/witness.go`; `internal/witness/witness_test.go` |
| S22-PR-04 | minor | code/correctness | PR code review found minor consistency and coverage gaps: `omitempty` conflicted with schema-required fields, `ci_same_job` was a bare string, and CLI envelope success was not tested end to end. | Accepted and fixed. Required normalized fields no longer use `omitempty`, `independenceSameJob` is a named constant, and the CLI has an explicit Buildkite `--witness-envelope` pass test. | Qwen PR code/correctness review; `internal/witness/witness.go`; `internal/witness/profiles.go`; `cmd/sdp-trace/main_test.go` |
| S22-PR-05 | minor | tracing/evidence | PR tracing review requested dedicated Buildkite signer-authority and run-binding negative coverage. | Accepted and fixed. Added Buildkite-specific tests for `witness_signer_authority_missing` and `witness_run_mismatch`. | DeepSeek PR tracing/evidence review; `internal/witness/profiles_test.go` |
| S22-PR-06 | not_assessed | tracing/evidence | Initial MiniMax PR tracing reviewer returned `404 page not found`. | Replaced, not counted as evidence. The plane was rerun with OpenRouter DeepSeek and returned `APPROVE` with no remaining critical or major findings. | MiniMax failed attempt not committed; DeepSeek PR tracing/evidence review |

## Review Evidence State

- Socratic spec review: assessed with MiniMax-M2.7 product-boundary,
  OpenRouter Qwen tracing/evidence, and Kimi K2P6 enterprise/security planes.
  Initial reviews returned `REVISE`; accepted critical and major findings were
  fixed. Focused product-boundary, tracing/evidence, and enterprise/security
  re-reviews returned `APPROVE` with no remaining critical or major findings.
- Implementation review: assessed with separate code/correctness,
  tracing/evidence, and security requirements-vs-implementation planes. Initial
  incomplete review packets were treated as non-evidence and replaced with
  full-patch review packets. Accepted critical/major security findings were
  fixed and focused Kimi K2P6 re-review returned `APPROVE`; no remaining
  critical or major implementation findings are recorded.
- PR-level review: assessed for PR #15 with Qwen code/correctness, DeepSeek
  tracing/evidence after replacing a failed MiniMax attempt, and Kimi K2P6
  requirements-vs-implementation security. Initial Kimi security review returned
  `BLOCK`; accepted major findings were fixed and focused Kimi re-review returned
  `APPROVE`. No remaining critical or major PR-level findings are recorded.
- CI: GitHub Actions `verify` passed during the PR-level review cycle; final
  PR-head CI must be re-checked outside this ledger immediately before
  ready/merge because committing ledger evidence creates a new head.
