**Verdict: APPROVE**

All three prior blockers are fixed by the provided fixes diff. No remaining critical or major findings.

**Blocker verification**

| ID | Fix | Evidence |
|---|---|---|
| **M1** `strongDigest` rejected SHA-384/512 | Changed length check from `!= 64` to `< 64`; added `TestStrongDigestAcceptsSHA256OrStrongerHex` covering 64-, 96-, and 128-char hex. | `internal/witness/profiles.go` `strongDigest` lines 753–758; `profiles_test.go` `TestStrongDigestAcceptsSHA256OrStrongerHex` |
| **M2** `finalizeRecordForWrite` copied tainted `ProfileID` / `ProfileVersion` / `ProviderKind` into safe record | Removed the three copy lines; safe record now derives those fields only from `baseRecord(record.Kind)`. | `internal/witness/profiles.go` `finalizeRecordForWrite` lines 412–417 |
| **M3** GitHub Actions path omitted schema-required normalized fields | `BuildGitHubActionsWithFetcher` now uses `baseRecord`, populates `ProfileStates`, `EstablishedTrustScope`, `ReasonCodes`, and calls `finalizeRecordForWrite` before write. | `internal/witness/witness.go` `BuildGitHubActionsWithFetcher`; `WriteGitHubActions` finalization call |

**Minor note (non-blocking)**
- `cmd/sdp-trace/main_test.go` references `writeJSONFileForTest`, which is not defined in the provided diffs; this is a test-helper compile-time issue, not a security or correctness defect.
