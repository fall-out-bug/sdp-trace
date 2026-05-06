# Block 16 Protected Gate Fixtures

These fixtures are committed examples for the protected gate output shape.

| Scenario | Artifact | Coverage |
|---|---|---|
| Missing checkpoint | `missing-checkpoint-cannot-verify.gate-result.json` | Committed JSON fixture. |
| Local signed checkpoint cannot pass protected profile | `local-signed-fail.gate-result.json` | Committed JSON fixture plus Go tests. |
| Local signed checkpoint with invalid run binding | `local-signed-invalid-run-binding-fail.gate-result.json` | Committed JSON fixture. |
| Missing signer policy | `missing-signer-policy-cannot-verify.gate-result.json` | Committed JSON fixture. |
| Signer mismatch | `signer-mismatch-fail.gate-result.json` | Committed JSON fixture. |
| Missing CI witness | `missing-ci-witness-cannot-verify.gate-result.json` | Committed JSON fixture plus Go tests. |
| Absent witness freshness | `absent-freshness-cannot-verify.gate-result.json` | Committed JSON fixture plus Go tests. |
| Stale witness freshness | `stale-ci-witness-fail.gate-result.json` | Committed JSON fixture plus Go tests. |
| CI source mismatch | `ci-source-mismatch-fail.gate-result.json` | Committed JSON fixture. |
| CI artifact mismatch | `ci-artifact-mismatch-fail.gate-result.json` | Committed JSON fixture plus Go tests. |
| Malformed override with trust-scope failure | `malformed-override-trust-scope-fail.gate-result.json` | Committed JSON fixture plus Go tests. |
| CI signed checkpoint with bound witness can satisfy protected profile | `ci-signed-pass.gate-result.json` | Committed JSON fixture plus Go tests. |
| Override present in protected profile | `override-present-protected-profile.gate-result.json` | Committed JSON fixture. |
| Missing protected input flags | Go test | `TestProtectedGateRequiresCheckpointPolicyAndWitnessFlags`. |

The JSON fixtures are static output-shape examples. The Go tests remain the
authoritative executable behavior checks for signatures, digests, and run
bindings that are generated during test execution.
