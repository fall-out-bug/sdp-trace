# Slice 3 Review: WS-017-B Wrap Output / Schema Drift

**Reviewer:** GLM adversarial review plane
**Date:** 2026-05-20
**Files reviewed:**
- `examples/flight-recorder/wrap-output-drift/wrap-output.txt`
- `examples/flight-recorder/wrap-output-drift/README.md`
- `tools/osscompat/probe_test.go`

**Companion implementation files read for context:**
- `tools/osscompat/probe.go`
- `tools/osscompat/runner.go`
- `tools/osscompat/main.go`
- `tools/osscompat/main_test.go`
- `schema/flight-recorder-run.schema.json`

**Verification:**
- `go test -count=1 -v ./tools/osscompat/` — all 18 tests pass (6 skipped for absent tools), 0 failures
- `go vet ./tools/osscompat/` — clean, no diagnostics
- No secrets found in any reviewed file
- No trailing whitespace, merge-conflict markers, or other `git diff --check` issues

---

## Findings

### 1. Advisory — Quality — `probe_test.go:99` and `probe_test.go:109`

**`TestRunJSONSchemaFixtures` and `TestRunJSONSchemaWrapDrift` test only the hard-coded return values of `runJSONSchemaFixtures()` / `runJSONSchemaWrapDrift()`, not any real drift detection.**

Both probes in `probe.go:75-86` return `stateCannotVerify` with a static string. The tests at `probe_test.go:96-113` verify only that these functions return `stateCannotVerify` with a non-empty reason. This is technically correct — the probes are intentionally stub — but the tests add no behavioral coverage beyond what a compiler already enforces (the function signature guarantees a return). These tests document intent rather than verify behavior.

**Justification for Advisory (not higher):** The actual drift structural evidence lives in `TestWrapOutputIsNotValidJSON` (line 136), which does perform a real check. The two `cannot_verify` tests are harmless; they are just low-value. No incorrect assertion is present.

---

### 2. Advisory — UX — `README.md` drift summary table

**The README.md drift summary table compares "Live Output" vs "Schema Requirement" but does not mention the `anyOf` constraint requiring either `event_refs` or `event_chain_head`, nor the conditional `if/then` rules for `witnessed_run_recorder` and `externally_witnessed_run` profiles.**

The table at lines 34–41 lists 6 absent fields. The schema also requires, via `anyOf`, that at least one of `event_refs` or `event_chain_head` be present, and via `allOf`, that `witness_ref` + `event_chain_head` be present when profile is `witnessed_run_recorder` or `externally_witnessed_run`. These additional constraints are not called out in the drift documentation.

**Justification for Advisory:** The 6 listed fields are sufficient to demonstrate the fundamental mismatch (plain text vs. JSON). A reader already understands the gap is total. The omitted `anyOf`/conditional constraints add no actionable information for resolving the blocker.

---

### 3. Advisory — DX — `run.json` contains a nondeterministic run ID

**`run.json` contains `run-3068560305`, a verbatim captured run ID. The README says this is the "verbatim output of `sdp-trace wrap /bin/true`" but does not explain whether future runs should replace this fixture or whether the fixture is frozen.**

If someone re-runs `sdp-trace wrap /bin/true`, they get a different run ID. The README's instructions at line 13 imply the fixture should be re-created, but the blocker status text at line 43 says "do not delete this directory." This creates ambiguity about whether the fixture is a living sample or frozen evidence.

**Justification for Advisory:** The `TestWrapOutputIsNotValidJSON` test only checks that the content is non-JSON, so the specific run ID doesn't matter. No test would break if the fixture were regenerated. The ambiguity is purely editorial.

---

## Summary

| Severity | Axis | Count |
|---|---|---|
| Critical | any | 0 |
| Important | any | 0 |
| Advisory | Quality | 1 |
| Advisory | UX | 1 |
| Advisory | DX | 1 |

**Disposition:** Approved with advisory notes. Three Advisory findings remain (low-value stub tests, missing `anyOf`/`allOf` detail in drift table, fixture regeneration ambiguity). No blockers, no Important findings. The drift evidence is honest and structurally verified by `TestWrapOutputIsNotValidJSON`. The `cannot_verify` state for the two schema-drift probes is the correct honest call. Tests compile, pass, and `go vet` is clean.
