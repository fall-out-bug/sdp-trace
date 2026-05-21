# Slice 3 Review: WS-017-B Wrap Output / Schema Drift

Reviewer: Kimi review plane  
Date: 2026-05-20  
Scope: `examples/flight-recorder/wrap-output-drift/*`, `tools/osscompat/probe_test.go`

## Verification Commands Run

```text
go test -count=1 ./tools/osscompat      → PASS
go vet ./tools/osscompat                → PASS
go vet ./...                            → PASS
gofmt -d tools/osscompat/probe_test.go  → PASS (no diff)
git diff --check                        → PASS
```

## Review Findings

### Advisory — Quality — `tools/osscompat/probe_test.go:141`

`TestWrapOutputIsNotValidJSON` is misnamed. It checks that `json.Unmarshal(data, &obj)` fails for a `map[string]any`, which proves the fixture is **not a JSON object** — not that it is "not valid JSON".

A JSON string (e.g., `"run_dir: ..."`) or JSON array would also fail to unmarshal into `map[string]any`, yet would be valid JSON. The schema requires an object, so the test behavior is semantically correct for its purpose; the name just overclaims.

**Evidence:**
```go
func TestWrapOutputIsNotValidJSON(t *testing.T) {
    ...
    var obj map[string]any
    if err := json.Unmarshal(data, &obj); err == nil {
        t.Fatal("expected wrap output to be non-JSON, but it parsed as JSON")
    }
}
```

**Suggested fix:** Rename to `TestWrapOutputIsNotJSONObject` to accurately describe the structural assertion.

---

## Correct

- **Drift evidence is honest:** `run.json` is a single plain-text line (`run_dir: .sdp-trace-runs/run-3068560305`), not a JSON object. The fixture is byte-identical to what the README claims and contains no secrets.
- **README is complete:** It lists the schema requirements, the drift dimensions, the blocker status, and the two resolution paths. It explicitly forbids deleting the directory or marking resolved without a source-bound proof commit.
- **Tests pass and are structurally sound:** `TestWrapOutputIsNotValidJSON` acts as a guard against accidental regression of the drift. If the fixture becomes a JSON object, the test fails, forcing an explicit update.
- **No unsafe operations:** `os.ReadFile` uses a relative path resolved from the package directory (standard Go test practice). No shell injection, no temp-file races, no secret leakage.
- **Probe architecture is clean:** `runJSONSchemaWrapDrift` returns `stateCannotVerify` with a documented reason, matching the repository rule that unverified claims must be marked `cannot_verify` or `not_assessed`.
- **File organization is clear:** The example directory is self-contained (`run.json` + `README.md`) and referenced from the `osscompat` test via a stable relative path.
