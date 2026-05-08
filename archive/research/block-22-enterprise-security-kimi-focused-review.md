**APPROVE**

No remaining critical or major enterprise/security findings. All prior findings are fixed in the reviewed spec:

| Finding ID | Fix location in spec | Verification |
|---|---|---|
| **S22-ES-01** | **Profile Requirements / GitLab CI** | Explicitly names malicious runner, compromised `CI_JOB_TOKEN`, and injected env threats; caps `same-job` or `same-runner` topologies below `external_witnessed` and below `ci_witnessed` without separate job isolation. |
| **S22-ES-02** | **CLI Boundary** | Requires profile normalizers to read only from the declared run-directory witness envelope or explicit input flags and to ignore verifier-process-inherited environment variables; env-only presence emits `cannot_verify` / `witness_environment_only_insufficient`. |
| **S22-ES-03** | **Witness Profile Contract**, **Closed Reason Codes**, **Fixture Matrix** | `run_binding` is a required contract field; `witness_run_binding_missing` and `witness_run_mismatch` reason codes exist; fixture `ci-run-id-mismatch` requires `fail` / `witness_run_mismatch`. |
| **S22-ES-04** | **Customer PKI**, **Normalized witness result**, **Closed Reason Codes** | Customer PKI profile distinguishes revocation state (`pass`/`fail`/`not_assessed`) and key custody state (`hsm`, `kms`, `software`, `unknown`, `not_assessed`, `cannot_verify`); reason codes `witness_revocation_not_assessed`, `witness_certificate_revoked`, and `witness_key_custody_not_assessed` are present. |
| **S22-ES-05** | **Air-Gapped Profile** | Explicitly forbids network calls in air-gapped validators; requires manually imported public keys, timestamps, and revocation snapshots to carry integrity digests verified out-of-band, with `fail` on mismatch. |
| **S22-ES-06** | **Safety Requirements** | Requires profile normalizers to validate or sanitize raw provider inputs against closed safety classes **before** structural parsing; rejects fields containing secret-like patterns; forbids relying solely on final-output redaction. |
| **S22-ES-07** | **Non-Goals** | Explicitly rules out composite, chained, or layered witness profiles for Block 22 and states that running multiple profiles must not upgrade trust scope beyond the strongest single established scope. |
