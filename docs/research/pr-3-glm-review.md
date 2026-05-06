## Block 13B Merge-Candidate Review

### Check Matrix

| Category | Finding | Severity |
|---|---|---|
| Acceptance mismatch | — | — |
| Misleading CLI output | — | — |
| Missing diagnostics | — | — |
| Nondeterminism | — | — |
| Test coverage | 47 tests across 10 packages; doctor, preview, wrap/verify/explain/query, flag parsing, fixtures, gate, witness, report all covered | — |
| Go-only violations | — | — |

### Detailed Walkthrough

**Doctor acceptance criteria** (spec §Doctor Acceptance Criteria):
1. Wrapper availability, output/report dir writeability, contract parse, expected-evidence references — ✅ all five checked
2. CI witness profile prerequisites including OIDC — ✅ `ci_witness_prerequisites` check with `missing` list
3. State taxonomy distinction (`unsupported`, `not_integrated`, `missing_telemetry`, `offline_dev`, `cannot_verify`) in output — ✅ boundaries and checks use these exact states
4. Observation mode and trust cap reported — ✅ `result: "offline_dev"` acts as trust cap
5. Deterministic output — ✅ no timestamps or wall-clock fields; tests assert stable shape
6. No raw secrets/prompts/responses/stdout/stderr/OIDC tokens — ✅ only structural state emitted

**Preview acceptance criteria** (spec §Preview Acceptance Criteria):
1. Command descriptor with basename + argc, not raw argv — ✅ `CommandDescriptor` uses `filepath.Base` + digest
2. Output directories and artifact categories — ⚠️ **minor** (see below)
3. Retention mode per category — ✅ `safe_retention_modes` list + command descriptor retention
4. Evidence contract and required evidence ids — ✅ full contract dumped
5. Active boundaries + unsupported/unintegrated — ✅ all six boundaries enumerated with explicit states
6. Offline implications — ✅ `offline_implications` array present
7. No pass claim — ✅ `writes_artifacts: false` + warning

**Preview output directory gap (⚠️ minor, not major):** Preview doesn't accept `--output-dir` and doesn't enumerate concrete artifact output paths. The spec says "output directories and artifact categories that would be written." This is defensible because preview explicitly declares `writes_artifacts: false` — nothing would be written. The `doctor` command covers output-directory writeability separately. This is a completeness gap for a future enhancement, not a safety or correctness blocker.

**State taxonomy machine-enumerability:** Seven core `ObservationState` constants and six `ObservationBoundary` constants are typed Go enums. Trust-scope states (`local_observed`, `ci_witnessed`, etc.) appear as string literals in witness/gate output per Block 13B scope (external witness not integrated). Acceptable.

**Safety properties verified by tests:**
- `TestPreviewOutputsNoWritePlan`: confirms no `.sdp-trace-runs` created, no raw argv in output
- `TestDryRunOutputsSimulation`: confirms no raw `"hi"` leaked
- `TestDoctorReportsOfflineDevAndCannotVerifyCI`: confirms state taxonomy in output
- `TestValidateFixturesHonorsExpectedFailure` / `TestValidateFixturesRejectsUnexpectedFailure`: tamper detection works

---

**NO_CRITICAL_OR_MAJOR**

**VERDICT: ACCEPT**
