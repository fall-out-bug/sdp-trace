Now I have a thorough understanding of the implementation. Let me compile findings.

---

## Code Review Findings — 011 Schema Documentation Validation

### F-001: No duplicate-name detection in index entries
**Severity**: medium (correctness)
**Location**: `tools/schemadoc/main.go:122-127`

The `check` function builds a `map[string]*SchemaEntry` from `idx.Schemas`, but if two entries share the same `name`, the second silently overwrites the first. The schema file would not appear as "missing from index" because the key exists, and one duplicate entry would pass validation invisibly.

**Impact**: A copy-paste error creating duplicate entries in `index.json` would go undetected. The spec (FR-002) says "detect missing schema entries, extra schema entries" — duplicate names are a form of extra entry that is invisible to the current logic.

**Fix**: Before building the map, iterate `idx.Schemas` and reject duplicate `Name` values.

---

### F-002: README ↔ index synchronization is not verified by the checker
**Severity**: medium (maintainability/integration)
**Location**: `tools/schemadoc/main.go:18` — `readmePath` is declared but never read or compared

The spec says "The table below is generated from that index and must stay synchronized" (README.md line 8), but the `check` function never verifies that `schema/README.md`'s table matches `schema/index.json`. The `--generate` flag produces a table to stdout, but there is no CI step or test that compares the generated output against the committed README.

**Impact**: A contributor could edit `index.json` without updating the README (or vice versa) and the checker would still pass. This directly undermines US-001's guarantee that "CI catches drift."

**Fix**: Either (a) add a `--verify-readme` mode that extracts and compares the table, or (b) have CI run `go run ./tools/schemadoc --generate > /tmp/table.md && diff` against the README table section, or (c) remove the README table and direct readers to `index.json` only.

---

### F-003: `os.Stat` errors not distinguished from permission failures
**Severity**: low (correctness)
**Location**: `tools/schemadoc/main.go:132`, `tools/schemadoc/main.go:149`

Both `os.Stat` calls treat any error as "file missing." If a file exists but is unreadable due to permissions, the checker reports "extra index entry (file missing)" or "broken example ref" — a misleading diagnosis. Similarly, `os.IsNotExist` should be checked to avoid masking filesystem errors.

**Impact**: Low in CI (permissions are typically uniform), but could confuse local debugging.

**Fix**: Use `os.IsNotExist(err)` to distinguish missing files from permission/access errors, and surface the latter differently.

---

### F-004: `not_assessed` status is valid but no test covers it
**Severity**: low (testability)
**Location**: `tools/schemadoc/main_test.go`

Tests cover `current`, `historical`, empty, and invalid statuses. There is no test for `statusNotAssessed` ("not_assessed"), which is a valid status. Given that the spec (US-002) explicitly mentions "active, historical, or not assessed" and the string `not_assessed` appears in the valid status map, a passing test for this status is a gap in coverage.

**Fix**: Add a `TestCheckPassesNotAssessedStatus` test case.

---

### F-005: `generateTable` does not escape pipe characters in purpose text
**Severity**: medium (correctness)
**Location**: `tools/schemadoc/main.go:165-168`

The `generateTable` function writes Markdown table rows using `fmt.Sprintf` without escaping `|` characters. If any `Purpose` string contains `|`, the Markdown table will be malformed. Currently no purpose string contains a pipe, but the index is hand-maintained and this is a latent formatting bug.

**Impact**: If a purpose containing `|` is added, the README table breaks silently.

**Fix**: Escape `|` → `\|` in purpose and name fields before writing table rows, or use a Markdown table writer that handles escaping.

---

### F-006: `readmePath` constant is unused — dead code
**Severity**: low (maintainability)
**Location**: `tools/schemadoc/main.go:18`

The constant `readmePath` is declared but never referenced. This is dead code that suggests an incomplete integration (see F-002).

**Fix**: Either use it (F-002) or remove it.

---

### F-007: No validation that index entry names match the `*.schema.json` glob
**Severity**: low (correctness)
**Location**: `tools/schemadoc/main.go:check()`

The checker verifies that each indexed file exists on disk, but does not validate that the `Name` field ends in `.schema.json`. An entry with `name: "index.json"` or `name: "random.txt"` would pass validation as long as the file exists in `schema/`. This could let non-schema files into the index.

**Impact**: Low (a file existing in `schema/` that isn't a schema is unlikely), but the checker's glob (`schema/*.schema.json`) specifically targets `.schema.json` files while the index validation has no corresponding constraint.

**Fix**: Add a check that `entry.Name` matches `*.schema.json`.

---

### F-008: `TestGenerateTableIsDeterministic` does not verify determinism
**Severity**: low (testability)
**Location**: `tools/schemadoc/main_test.go:194-208`

The test name asserts determinism, but the test only checks that the output contains both schema names — it never calls `generateTable` twice and compares outputs, nor does it verify row order matches index order. The test comment says "Rows should follow index order" but does not assert it.

**Fix**: Call `generateTable` twice, assert `==`. Optionally assert the specific order of lines.

---

### F-009: `Index.Version` is never validated
**Severity**: informational
**Location**: `tools/schemadoc/main.go:68-71`, `schema/index.json:2`

The index JSON has a `"version": "1"` field that is parsed but never validated. If the version changes to `"2"` with a different schema shape, the checker would silently accept or fail unpredictably.

**Impact**: Low for now (single version), but a version check would prevent forward-compatibility surprises.

**Fix**: Validate that `idx.Version` equals `"1"` and reject unknown versions.

---

### F-010: `jq empty schema/*.json` in CI also checks `index.json`
**Severity**: informational
**Location**: `.github/workflows/ci.yml:40`

The glob `schema/*.json` includes `schema/index.json`, so `jq empty` validates its JSON syntax. This is benign (valid JSON) but means `index.json` is checked as though it were a schema file by the JSON syntax step. The `schemadoc` checker handles its structural validation separately. Not a defect, just worth noting for awareness.

---

### Requirements Coverage Assessment

| Requirement | Status | Notes |
|---|---|---|
| FR-001 Go checker | ✅ met | `tools/schemadoc/main.go` |
| FR-002 Detect missing/extra/broken | ⚠️ partial | Duplicate names not detected (F-001); README drift not detected (F-002) |
| FR-003 Status + purpose required | ✅ met | Checked for all entries |
| FR-004 Deterministic generation | ⚠️ partial | Table generation is deterministic but not verified by test (F-008) |
| FR-005 Distinguish syntax vs semantic | ✅ met | README section explains this |
| US-001 CI catches drift | ⚠️ partial | `index.json` ↔ files is checked; README ↔ `index.json` is not (F-002) |
| US-002 Purpose/status metadata | ✅ met | Validated by checker |
| US-003 Example coverage | ⚠️ partial | No test for `not_assessed` example coverage state; spec says "distinguish present example refs from `not_assessed` refs" but the checker treats missing `examples` field as valid (no examples = no coverage claimed). No `not_assessed` sentinel value is defined or tested. |
| US-004 No Node.js tooling | ✅ met | Go-only |

---

```json
{
  "schemaVersion": "subagent-result/v1",
  "role": "code-reviewer",
  "plane": "code",
  "subject": "011 Schema Documentation Validation PR review",
  "disposition": "conditional_pass",
  "summary": "The schemadoc checker is clean Go code with good test coverage for the happy path and most error paths. The core index↔file validation works. However, two behavioral gaps should be addressed before merge: (1) no duplicate-name detection in index entries (F-001), and (2) README↔index synchronization is not verified despite the README claiming it must stay synchronized (F-002). F-005 (pipe escaping) is a latent formatting bug. The remaining findings are testability and maintainability improvements.",
  "findings": [
    {
      "id": "F-001",
      "severity": "medium",
      "category": "correctness",
      "location": "tools/schemadoc/main.go:122-127",
      "title": "No duplicate-name detection in index entries",
      "detail": "The map-building loop at line 124 silently overwrites entries with duplicate names. A duplicated entry would not be detected as 'missing' or 'extra', violating FR-002.",
      "fix": "Iterate idx.Schemas before building the map and reject duplicate Name values."
    },
    {
      "id": "F-002",
      "severity": "medium",
      "category": "integration",
      "location": "tools/schemadoc/main.go:18, schema/README.md:8",
      "title": "README↔index synchronization is not verified by the checker",
      "detail": "readmePath is declared but never used. The README says 'must stay synchronized' but no CI step or test enforces this. A contributor could edit index.json without updating the README and the checker would pass.",
      "fix": "Add --verify-readme mode or CI diff step, or remove the README table."
    },
    {
      "id": "F-003",
      "severity": "low",
      "category": "correctness",
      "location": "tools/schemadoc/main.go:132,149",
      "title": "os.Stat errors not distinguished from permission failures",
      "detail": "Any os.Stat error is treated as 'file missing'. Permission errors would produce misleading diagnostics.",
      "fix": "Use os.IsNotExist(err) to distinguish missing files from access errors."
    },
    {
      "id": "F-004",
      "severity": "low",
      "category": "testability",
      "location": "tools/schemadoc/main_test.go",
      "title": "not_assessed status has no test coverage",
      "detail": "Tests cover current, historical, empty, and invalid statuses but not not_assessed, despite it being a valid status and explicitly mentioned in US-002.",
      "fix": "Add TestCheckPassesNotAssessedStatus."
    },
    {
      "id": "F-005",
      "severity": "medium",
      "category": "correctness",
      "location": "tools/schemadoc/main.go:165-168",
      "title": "generateTable does not escape pipe characters in purpose text",
      "detail": "Markdown table rows use fmt.Sprintf without escaping |. Any purpose containing | would break the table layout.",
      "fix": "Escape | to \\| in name and purpose fields before writing table rows."
    },
    {
      "id": "F-006",
      "severity": "low",
      "category": "maintainability",
      "location": "tools/schemadoc/main.go:18",
      "title": "readmePath constant is unused dead code",
      "detail": "The constant readmePath is declared but never referenced, suggesting incomplete README verification integration.",
      "fix": "Use it (F-002) or remove it."
    },
    {
      "id": "F-007",
      "severity": "low",
      "category": "correctness",
      "location": "tools/schemadoc/main.go:check()",
      "title": "No validation that index entry names match *.schema.json glob",
      "detail": "The checker's file glob targets *.schema.json but index entry names are not validated for the same suffix. A non-schema file could be accepted.",
      "fix": "Add suffix check: strings.HasSuffix(entry.Name, \".schema.json\")."
    },
    {
      "id": "F-008",
      "severity": "low",
      "category": "testability",
      "location": "tools/schemadoc/main_test.go:194-208",
      "title": "TestGenerateTableIsDeterministic does not actually verify determinism",
      "detail": "The test only checks output contains entries, never calls generateTable twice or asserts equality. Row order is not asserted despite the comment.",
      "fix": "Call generateTable twice and assert ==."
    },
    {
      "id": "F-009",
      "severity": "informational",
      "category": "maintainability",
      "location": "tools/schemadoc/main.go:68-71, schema/index.json:2",
      "title": "Index.Version is parsed but never validated",
      "detail": "If version changes to 2 with a different schema shape, the checker would silently accept or fail unpredictably.",
      "fix": "Validate idx.Version equals \"1\" and reject unknown versions."
    }
  ],
  "requirementsCoverage": {
    "FR-001": "met",
    "FR-002": "partial — duplicate names and README drift not detected",
    "FR-003": "met",
    "FR-004": "partial — not verified by test",
    "FR-005": "met",
    "US-001": "partial — index↔files checked; README↔index not checked",
    "US-002": "met",
    "US-003": "partial — not_assessed example coverage sentinel not defined or tested",
    "US-004": "met"
  },
  "blockers": ["F-001", "F-002"]
}
```
