# Slice 1 Review: tools/osscompat probe harness

**Reviewer:** GLM adversarial plane
**Date:** 2026-05-20
**Files reviewed:** `main.go`, `main_test.go`, `probe.go`, `probe_test.go`, `runner.go`
**Build status (local checkout only, not source-bound proof):** `go vet` clean, all tests pass, 54.9% statement coverage. Reproduce with `go test -count=1 ./tools/osscompat/...` and `go vet ./tools/osscompat/...`.

---

## Findings

### F1 — `jsonschema-fixtures` probe returns `pass` without evidence

| Axis      | Quality / Trust |
|-----------|-----------------|
| Severity  | **Important**   |
| Location  | `probe.go:87–89` |

`runJSONSchemaFixtures` unconditionally returns `statePass` with a prose-only reason. When `check-jsonschema` is on PATH, the probe skips to `Run()` and reports `pass` without executing any validation. This violates the AGENTS.md trust rule: *"Every claim about a gate or verdict must be evidence-backed or marked `not_assessed`."*

The Description string compounds this: *"Validate flight-recorder fixtures against local schema refs"* — but no fixtures are validated.

**Fix:** Either implement actual fixture validation or change state to `stateNotAssessed` with a concrete reason (e.g., `"fixture validation not yet implemented; tool present"`).

---

### F2 — Five probe Descriptions overclaim what they actually do

| Axis      | Quality         |
|-----------|-----------------|
| Severity  | **Important**   |
| Location  | `probe.go:28–65` (registry definitions) |

Descriptions vs. actual behavior:

| Probe | Description says | Actually does |
|-------|-----------------|---------------|
| `jsonschema-fixtures` | "Validate flight-recorder fixtures against local schema refs" | Returns hardcoded `pass` |
| `jsonschema-wrap-drift` | "Validate live sdp-trace wrap output vs … schema.json" | Returns hardcoded `cannot_verify` |
| `opa-policy` | "Evaluate simplified adapter-capture pass/fail rule" | Runs `opa version` only |
| `intoto-wrap` | "Wrap command and sign link metadata" | Runs `in-toto-run --version` only |
| `cosign-local-sign` | "Sign and verify local blob with local key" | Runs `cosign version` only |

**Fix:** Align descriptions with actual behavior, e.g. *"Verify cosign is present and responds to `version`."*

---

### F3 — Silent error suppression in single-probe path

| Axis      | Quality         |
|-----------|-----------------|
| Severity  | **Important**   |
| Location  | `main.go:38`    |

```go
_ = printResults(stdout, []probeResult{r}, *asJSON)
```

The all-probes path (line 44–46) checks the error from `printResults` and returns exit 2 on failure. The single-probe path silently discards it. This is an inconsistency that could mask I/O failures.

**Fix:** Handle the error identically to the all-probes path:
```go
if err := printResults(stdout, []probeResult{r}, *asJSON); err != nil {
    fmt.Fprintf(stderr, "print results: %v\n", err)
    return 2
}
```

---

### F4 — `repoRoot()` is dead code

| Axis      | DX              |
|-----------|-----------------|
| Severity  | **Advisory**    |
| Location  | `probe.go:156–160` |

`repoRoot()` is defined but never called. Meanwhile `runCUEImport` (line 114–122) uses a hardcoded relative path `schema/flight-recorder-run.schema.json`. Either use `repoRoot()` to construct the path or remove the dead function.

---

### F5 — All registry probe implementations have 0% test coverage

| Axis      | Quality         |
|-----------|-----------------|
| Severity  | **Important**   |
| Location  | `probe.go:87–155` |

Coverage report confirms every `runXxx` function in the registry has 0% coverage. Tests exercise only synthetic probes (`go-version`, `missing-tool-test`). Real probes like `runOPAPolicy`, `runCUEImport`, etc. are completely untested.

**Fix:** Add at least one test per probe function. For external-tool probes, test the "tool missing" path via the existing `runProbe` guard, and add isolated unit tests for the probe logic (possibly skipping if the tool is absent, using `testing.Short()` or a build-tag gate).

---

### F6 — `jsonschema-wrap-drift` has `NeedsTool` but ignores it

| Axis      | Quality         |
|-----------|-----------------|
| Severity  | **Advisory**    |
| Location  | `probe.go:80–84` |

`jsonschema-wrap-drift` has `NeedsTool: "check-jsonschema"`, but its `Run` function unconditionally returns `stateCannotVerify` — it never uses the tool even when present. When `check-jsonschema` IS available, the probe still returns `cannot_verify`. The `NeedsTool` guard gives the false impression that installing the tool would enable real checking.

**Fix:** Remove `NeedsTool` and always return `cannot_verify`, or implement actual drift checking when the tool is present.

---

### F7 — No summary line in output

| Axis      | UX              |
|-----------|-----------------|
| Severity  | **Advisory**    |
| Location  | `runner.go:37–49` (`printResults`) |

Users must manually scan all lines to determine overall status. A summary like `3 pass, 2 fail, 2 not_assessed` at the end would improve scanability.

---

### F8 — No `-list` flag to enumerate probes without running them

| Axis      | UX              |
|-----------|-----------------|
| Severity  | **Advisory**    |
| Location  | `main.go:17–20` |

There is no way to discover available probes without running them. A `-list` flag would improve discoverability.

---

### F9 — No test for `-json` + `-probe` combination

| Axis      | Quality         |
|-----------|-----------------|
| Severity  | **Advisory**    |
| Location  | `main_test.go`  |

The test suite covers `-json` alone and `-probe` alone but never both together. This is a reachable code path (flag parsing accepts both simultaneously) that should have a test case.

---

## Summary

| Severity  | Count |
|-----------|-------|
| Critical  | 0     |
| Important | 4     |
| Advisory  | 5     |

**Positive observations:**
- Clean file decomposition (main / probe / runner).
- `go vet` clean, all tests pass (local checkout evidence; reproduce with `go test -count=1 ./tools/osscompat/...`).
- No command injection risk — all `exec.Command` args are fixed strings, no user input flows into commands.
- Honest `cannot_verify` and `not_assessed` states in the runner infrastructure.
- Good separation of `run()` for testability with dependency injection for stdout/stderr.
- Text output is well-aligned and readable.

**Verdict:** Not LGTM. Four Important findings require resolution — the most critical being F1 (evidence-free `pass` verdict violating trust rules) and F3 (silent error suppression).
