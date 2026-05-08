## Critical

**C1 — GitLab CI missing co-located agent/witness topology cap**
*Section: Profile Requirements / GitLab CI*
The Buildkite profile explicitly caps co-located topologies, but the GitLab profile omits the same cap. A witness executing inside the same GitLab runner job or pod as the build has no independence boundary and can be compromised by a malicious runner or injected job token.
*Fix required:* Add bullet: “same-job or same-runner co-located witness and build topologies that cap trust below `external_witnessed`, and below `ci_witnessed` when the witness lacks a separate job isolation boundary, because a compromised runner or malicious job definition can forge environment variables and job tokens.”

**C2 — Witness profiles must not inherit the verifier process environment as CI evidence**
*Section: CLI Boundary; Profile Requirements / GitLab CI and Buildkite*
The CLI Boundary does not require that profile normalizers read CI witness inputs only from the declared `<runs-root-or-run-dir>` witness envelope or explicit flags. Without this isolation, a local verifier invocation with exported CI environment variables can falsely establish a CI witness.
*Fix required:* Add CLI Boundary bullet: “Profile normalizers must read CI witness inputs only from the declared run-directory witness envelope or explicit input flags; they must ignore environment variables inherited from the verifier process. The command must detect and treat as `cannot_verify` any CI witness invocation where the input envelope is missing but matching CI env vars are present in the verifier process.”

---

## Major

**M1 — Missing explicit run-id replay resistance in CI profile requirements**
*Section: Witness Profile Contract; Profile Requirements / GitLab CI and Buildkite*
The contract mentions `run_binding_state` but neither CI profile explicitly requires that a witness result produced for a different run/build/pipeline id must be rejected. Block 22 must close this gap independently before gate-level binding (FR-059) is reached.
*Fix required:* Add Witness Profile Contract requirement: “The verifier must reject a witness result whose embedded run id, build id, or pipeline id does not match the run under inspection, emitting `fail` with reason `witness_run_mismatch`.” Add corresponding fixtures to the Fixture Matrix.

**M2 — Customer PKI lacks revocation and key-custody verifier states**
*Section: Profile Requirements / Customer PKI; Witness Profile Contract*
The profile lists `expired validity` but omits certificate revocation (CRL/OCSP unavailable) and key custody (HSM, KMS, software, unknown). Enterprise auditors require these as explicit `not_assessed` or `cannot_verify` states, not silent omissions.
*Fix required:* Add Customer PKI bullets: “certificate revocation state (CRL/OCSP unavailable as `not_assessed`, revoked as `fail`);” and “key custody state as declared by the authority policy or `not_assessed`.” Add `key_custody_state` to the normalized witness result fields.

**M3 — Air-gapped profile lacks mandatory network-call prohibition and import-integrity rules**
*Section: Air-Gapped Profile; Non-Goals*
The section describes offline documentation but does not explicitly forbid network calls in any air-gapped validator, nor does it require integrity verification for manually imported public keys and timestamps.
*Fix required:* Add Air-Gapped Profile bullets: “network calls are explicitly forbidden in any air-gapped validator; any verification requiring external network access must emit `cannot_verify`;” and “manually imported public keys and timestamps must carry an integrity digest verified against an out-of-band source, with verification failure emitting `fail`.”

**M4 — Safety requirements lack input-sanitization boundary before profile normalization**
*Section: Safety Requirements; Witness Profile Contract*
Output safety is required, but the spec does not mandate that profile normalizers scan raw CI or PKI inputs for secret-like values before parsing. A secret in an env var could leak via debug logs, panics, or intermediate structs before final-output redaction.
*Fix required:* Add Safety Requirements bullet: “Profile normalizers must validate or sanitize raw provider inputs against `safe_output_classes` before structural parsing; if a secret-like pattern is detected in an input field, the normalizer must reject the field or replace it with a redaction digest, and must not rely solely on final-output redaction.”

**M5 — Composite witness profiles are not explicitly ruled out**
*Section: Goal; CLI Boundary; Non-Goals*
Enterprise hybrid environments (e.g., GitLab CI + customer PKI signing) might attempt to combine profiles. The CLI shows single `--kind`, but the spec must explicitly state that Block 22 does not support composite or chained witness profiles.
*Fix required:* Add Non-Goals bullet: “No composite, chained, or layered witness profiles; Block 22 profiles are single-source only. Running multiple witness profiles and combining their results must not upgrade trust scope beyond the strongest single established scope, and composite witnessing is explicitly unsupported in this block.”

---

## Minor

**m1 — GitLab CI should name runner env-injection and job-token forgery as explicit threats**
*Section: Profile Requirements / GitLab CI; Problem*
The spec downgrades env-only evidence, but implementers need the threat model in the profile text to understand why the boundary exists.
*Fix required:* Add opening note to GitLab CI profile: “A malicious runner, compromised `CI_JOB_TOKEN`, or injected `.gitlab-ci.yml` environment variable can forge CI metadata; therefore, environment-only inputs cannot establish trust.”

**m2 — Buildkite should name pipeline `env` secret exposure as an explicit threat**
*Section: Profile Requirements / Buildkite*
Pipeline-level `env` blocks and agent hooks can inject secrets that agents then report as metadata.
*Fix required:* Add opening note to Buildkite profile: “Buildkite pipeline `env` blocks and agent hooks can inject secrets into agent-reported metadata; therefore, agent-reported inputs without signed or authority-bound facts cannot establish trust.”

**m3 — Normalized witness result should include `key_custody_state`**
*Section: Witness Profile Contract*
The `signing_boundary` field is too vague for enterprise consumers who must audit whether keys are HSM-backed.
*Fix required:* Add `key_custody_state` (allowed values: `hsm`, `kms`, `software`, `unknown`, `not_assessed`) to the normalized witness result fields and update `signing_boundary` description to reference it.

**m4 — Air-gapped timestamp semantics need wall-clock divergence guidance**
*Section: Air-Gapped Profile*
In air-gapped environments, carried-in timestamps should not be compared to unsynchronized local wall clocks.
*Fix required:* Add bullet: “Freshness checks in air-gapped environments must use the carried evidence's internal monotonic sequence or signed timestamp chain, not comparison to an unsynchronized local wall clock, unless wall clock synchronization is itself verified and declared.”

---

**Verdict: REVISE**
