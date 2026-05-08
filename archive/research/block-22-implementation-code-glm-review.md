**APPROVE** — None of the seven prior findings remain.

| Prior finding | Current status |
|---|---|
| **IdentityState missing from `ProfileStates`** | ✅ `IdentityState string \`json:"identity_state"\`` is the first field in the struct. Populated by `defaultProfileStates()` and all builder paths. |
| **Reason-code constants vs. string literals** | ✅ All 23 spec reason codes declared as `const` in `witness.go`. `profiles.go` uses constants exclusively (e.g. `ReasonMissingIdentity`, `ReasonSignerMismatch`, `ReasonRunMismatch`, etc.). |
| **Output safety asserted without scan** | ✅ `finalizeRecordForWrite` marshals the record *without* `OutputSafety`, runs `forbiddenOutputPresent(raw)` on the serialized bytes, and only sets `statePass` when the scan confirms absence. Unsafe candidates are downgraded to `fail`/`witness_unsafe_output_candidate` before write. |
| **Customer PKI signature/digest/freshness** | ✅ `verifyFreshnessSignature()` checks Ed25519 sig on canonical payload; `strongDigest()` + equality check on payload digest; `freshnessCurrent()` validates `IssuedAt`/`ValidUntil` window. |
| **Run-id replay check** | ✅ `runIDMatches(runsRoot, envelope.CI.RunID)` in `BuildCIEnvelopeProfile`; `runIDMatches(runsRoot, freshness.RunID)` in `BuildCustomerPKI`. Both emit `fail`/`witness_run_mismatch` on mismatch via `runIDsFromRoot`. |
| **Fixture matrix breadth** | ✅ 23 fixtures in `fixture-matrix.json` covering all spec-required rows: valid/malformed/mismatch/stale/env-only/unsupported/unsafe/cert-revoked/key-custody variants across all four profiles including `air-gapped-v1`. |
| **CLI reason mismatch** | ✅ `TestWitnessCommandBuildkiteRequiresExplicitEnvelope` clears all CI env vars and correctly asserts `"reason": "witness_identity_missing"`, matching the `ReasonMissingIdentity` branch when `ambientCIEnvPresent` returns false. |
