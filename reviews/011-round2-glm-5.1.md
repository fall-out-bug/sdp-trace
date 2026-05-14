Now I have a comprehensive picture. Here is my review:

---

## Code Review Findings — 011 Schema Documentation Validation (Round 2)

### F-01: `readIndex` error messages use hardcoded constant instead of actual `path` parameter
**Severity**: Low | **Category**: Correctness

`tools/schemadoc/main.go:99-104`

```go
func readIndex(path string) (*Index, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read %s: %w", indexPath, err)  // ← uses constant, not path
    }
    var idx Index
    if err := json.Unmarshal(b, &idx); err != nil {
        return nil, fmt.Errorf("parse %s: %w", indexPath, err)  // ← same
    }
```

`readIndex` accepts a `path` parameter but error messages always report `indexPath` (`"schema/index.json"`). In production the two always match, but in tests like `TestCheckFailsUnsupportedVersion` the actual path is a temp file — so error messages will misleadingly reference `schema/index.json`. Fix: use the `path` parameter in error format strings.

### F-02: `generateTable` does not canonicalize sort order
**Severity**: Low | **Category**: Maintainability

`tools/schemadoc/main.go:194-202`

`generateTable` emits rows in the order they appear in `idx.Schemas`. If someone reorders entries in `index.json` (no semantic change), the generated table changes, and `schemadoc -verify-readme` fails. The README table is effectively order-sensitive with no canonical sort. Fix: sort entries by name before emitting, or sort in `checkReadmeAt` before comparison. This would make the tool resilient to benign index reordering.

### F-03: No validation that `example_coverage: "not_assessed"` should not list examples
**Severity**: Low | **Category**: Correctness / Spec alignment

Spec US-003 says: *"Where examples exist, the index names representative examples or marks example coverage `not_assessed`"* — the "or" implies these are alternatives. The checker validates `example_coverage: "present"` requires a non-empty examples list, but does not validate the reverse: `example_coverage: "not_assessed"` with a populated examples list passes silently. An entry with both `"example_coverage": "not_assessed"` and a populated `"examples"` array is semantically ambiguous. Fix: emit a warning or error when `example_coverage` is `"not_assessed"` or empty but examples are listed.

### F-04: Missing test for README with absent schemadoc markers
**Severity**: Low | **Category**: Testability

The test suite covers README sync pass (`TestCheckReadmePassesWhenSynchronized`) and README drift (`TestCheckReadmeFailsWhenDrifted`), but does not test the case where both `<!-- schemadoc-start -->` and `<!-- schemadoc-end -->` markers are entirely absent from the README. The `checkReadmeAt` function handles this case (`startIdx == -1 || endIdx == -1`), but it is untested. A simple test with a bare README would close this gap.

### F-05: Missing test for invalid `example_coverage` value
**Severity**: Low | **Category**: Testability

The checker validates `example_coverage` against `validExampleCoverages` (accepts `"present"`, `"not_assessed"`, `""`), but no test exercises a truly invalid value like `"unknown"` or `"partial"`. A test asserting `invalid example_coverage` in the error output would close this gap.

### F-06: `TestRunAcceptsCurrentIndex` is a non-isolated integration canary
**Severity**: Low | **Category**: Testability

`tools/schemadoc/main_test.go:236-238`

```go
func TestRunAcceptsCurrentIndex(t *testing.T) {
    if err := run(false, false); err != nil {
        t.Fatalf("run: %v", err)
    }
}
```

This test calls `run(false, false)` which hits the real repo filesystem via `repoRoot()`. It acts as a canary — any schema drift in the real repo fails this test. This is intentional but fragile: during active development of the index, this test breaks and blocks all other test runs. Consider:
- Naming it `TestIntegrationCurrentIndex` or gating it behind a build tag
- Adding a comment explaining the canary purpose

### F-07: `repoRoot()` relies on `runtime.Caller(0)` — breaks under `go install`
**Severity**: Informational | **Category**: Maintainability

`tools/schemadoc/main.go:58-63`

```go
func repoRoot() string {
    _, file, _, ok := runtime.Caller(0)
    if !ok {
        return "."
    }
    return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
```

This embeds the source file's path at compile time. It works for `go run ./tools/schemadoc` (which is how CI uses it) but will produce wrong paths if the binary is `go install`-ed and run from elsewhere. The fallback `"."` silently produces incorrect behavior rather than an explicit error. Since CI exclusively uses `go run`, this is acceptable, but worth documenting in a comment.

---

### Verified — No Issues Found

| Aspect | Assessment |
|---|---|
| Index ↔ schema file 1:1 correspondence | ✅ All 62 `*.schema.json` files matched in `index.json` |
| Example reference integrity | ✅ All 6 example paths resolve to existing files |
| README table ↔ index.json sync | ✅ Table between markers matches `generateTable()` output |
| Status values | ✅ All entries use `current`, `historical`, or `not_assessed` |
| Purpose fields | ✅ All entries have non-empty purpose strings |
| CI wiring | ✅ Both `schemadoc` and `schemadoc -verify-readme` run in CI |
| No Node.js tooling | ✅ Go-only, per spec FR-004 |
| Deterministic generation | ✅ Tested and confirmed |
| Pipe escaping in table | ✅ Tested and confirmed |
| `jq empty schema/*.json` also validates `index.json` syntax | ✅ Confirmed from CI step |

---

```json
{
  "schemaVersion": "subagent-result/v1",
  "role": "code-reviewer",
  "plane": "code",
  "subject": "011 Schema Documentation Validation PR review (round 2)",
  "disposition": "advisory",
  "findings": [
    {
      "id": "F-01",
      "severity": "low",
      "category": "correctness",
      "location": "tools/schemadoc/main.go:99-104",
      "summary": "readIndex error messages use hardcoded indexPath constant instead of path parameter, producing misleading messages in test contexts",
      "recommendation": "Replace `indexPath` with `path` in the fmt.Errorf format strings inside readIndex"
    },
    {
      "id": "F-02",
      "severity": "low",
      "category": "maintainability",
      "location": "tools/schemadoc/main.go:194-202",
      "summary": "generateTable emits rows in index.json order without sorting, causing unnecessary README drift on benign entry reordering",
      "recommendation": "Sort entries by name before emitting rows, or sort both sides in checkReadmeAt comparison"
    },
    {
      "id": "F-03",
      "severity": "low",
      "category": "correctness",
      "location": "tools/schemadoc/main.go:176-183",
      "summary": "No validation that example_coverage 'not_assessed' entries should not list examples; spec US-003 implies mutual exclusivity",
      "recommendation": "Add a check: if example_coverage is not 'present' but examples list is non-empty, emit a warning or error"
    },
    {
      "id": "F-04",
      "severity": "low",
      "category": "testability",
      "location": "tools/schemadoc/main_test.go",
      "summary": "No test for README with absent schemadoc markers (startIdx/endIdx == -1 path)",
      "recommendation": "Add TestCheckReadmeFailsWithoutMarkers exercising a README without schemadoc-start/end comments"
    },
    {
      "id": "F-05",
      "severity": "low",
      "category": "testability",
      "location": "tools/schemadoc/main_test.go",
      "summary": "No test for invalid example_coverage value (e.g. 'unknown', 'partial')",
      "recommendation": "Add TestCheckFailsInvalidExampleCoverage asserting 'invalid example_coverage' in error output"
    },
    {
      "id": "F-06",
      "severity": "low",
      "category": "testability",
      "location": "tools/schemadoc/main_test.go:236-238",
      "summary": "TestRunAcceptsCurrentIndex is a non-isolated integration canary that fails on any real-repo schema drift, blocking all test runs during active development",
      "recommendation": "Rename to TestIntegrationCurrentIndex or gate behind a build tag; add comment explaining canary purpose"
    },
    {
      "id": "F-07",
      "severity": "informational",
      "category": "maintainability",
      "location": "tools/schemadoc/main.go:58-63",
      "summary": "repoRoot() uses runtime.Caller(0) which embeds source paths at compile time; breaks under go install but works for go run (CI path)",
      "recommendation": "Add a comment documenting the go run-only constraint"
    }
  ],
  "veto": null,
  "evidence": {
    "schemaCountMatch": true,
    "exampleRefsValid": true,
    "readmeTableSync": true,
    "ciWired": true,
    "goOnly": true
  }
}
```
