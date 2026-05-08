## Second-Pass Block 13B Review

### Prior Finding 1 — Doctor missing required acceptance checks

| Required check | Prior state | Current state |
|---|---|---|
| Wrapper availability | Hardcoded `"pass"` | Still hardcoded `"pass"` but semantically correct for a single-binary — if `sdp-trace doctor` runs, the wrapper binary exists. Acceptable. |
| Output directory writeability | Absent | ✅ `writablePathCheck("output_directory", ...)` probes filesystem with temp-file write test |
| Expected-evidence references | Absent | ✅ `expectedEvidenceReferenceCheck(contract)` validates all event types against known set |
| Report directory policy | Absent | ✅ `writablePathCheck("report_directory", ...)` probes filesystem |

`ciWitnessPrerequisiteCheck` enumerates all required GitHub Actions OIDC fields. Tests (`TestDoctorReportsOfflineDevAndCannotVerifyCI`, `TestDoctorReportsUnwritableOutputDirectoryCannotVerify`, `TestDoctorReportsUnsupportedExpectedEvidenceCannotVerify`, `TestDoctorReportsContractLoadFailureCannotVerify`) cover all branches. Result is `offline_dev` not `pass`. **RESOLVED.**

### Prior Finding 2 — Preview missing required surfaces

`runPreviewCommand` now emits:

- **`boundaries`** via `previewBoundaries()` — all six boundaries with explicit states: `process_wrapper`→`pass`, `adapter_socket`→`not_integrated`, `tool_wrapper`→`unsupported`, `vcs_pr_observer`→`not_integrated`, `ci_observer`→`offline_dev`, `external_witness`→`not_integrated` ✅
- **`offline_implications`** via `previewOfflineImplications()` — `ci_witnessed`→`offline_dev`, `external_witnessed`→`not_integrated` ✅

`TestPreviewOutputsNoWritePlan` asserts boundary states and offline implications are present and correct. No raw argv leakage confirmed. **RESOLVED.**

### New-issue scan

I checked the full spec acceptance criteria against current code:

- **State taxonomy**: All 15 spec states present as stable Go constants. ✅
- **Boundary spec**: All 6 boundaries present as Go constants. ✅
- **Retention modes**: All 5 modes present, emitted by both doctor and preview. ✅
- **Determinism**: No `time.Now()`, random, or hostname fields in doctor or preview output. ✅
- **No-overclaim**: Doctor reports `offline_dev`; preview includes `"warning": "no run artifacts were written"`; CI witness is `cannot_verify` outside GitHub Actions; `audit_grade_gate` remains `cannot_verify`. ✅
- **Safety**: `NewCommandDescriptor` stores basename + SHA-256 digest only; tests confirm no raw-argv leak in serialized output. ✅
- **Overhead measurement**: Correctly `not_assessed` per spec — protocol is documented, measurement is a follow-up execution step. ✅
- **External trust**: Correctly `not_integrated` per spec. ✅

No new critical or major findings.

---

## NO_CRITICAL_OR_MAJOR

Both prior major findings are fully resolved. **VERDICT: ACCEPT.**
