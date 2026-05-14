package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func schemaDir(t *testing.T, root string) string {
	d := filepath.Join(root, "schema")
	if err := os.MkdirAll(d, 0755); err != nil {
		t.Fatal(err)
	}
	return d
}

func TestExitCode(t *testing.T) {
	if got := exitCode(nil, os.Stderr); got != 0 {
		t.Fatalf("exitCode(nil) = %d, want 0", got)
	}
	err := errors.New("drift")
	if got := exitCode(err, os.Stderr); got != 1 {
		t.Fatalf("exitCode(err) = %d, want 1", got)
	}
}

func TestCheckPassesWhenIndexMatchesFiles(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sd, "b.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A purpose"},
			{Name: "b.schema.json", Status: "historical", Purpose: "B purpose"},
		},
	}
	if err := check(root, idx); err != nil {
		t.Fatalf("check: %v", err)
	}
}

func TestCheckFailsMissingIndexEntry(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{Version: indexVersion, Schemas: []SchemaEntry{}}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for missing index entry")
	}
	if !strings.Contains(err.Error(), "missing index entry: a.schema.json") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsExtraIndexEntry(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A"},
			{Name: "missing.schema.json", Status: "current", Purpose: "M"},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for extra index entry")
	}
	if !strings.Contains(err.Error(), "extra index entry") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsDuplicateName(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A"},
			{Name: "a.schema.json", Status: "historical", Purpose: "B"},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for duplicate name")
	}
	if !strings.Contains(err.Error(), "duplicate index entry") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsInvalidNameSuffix(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.json", Status: "current", Purpose: "A"},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for invalid name suffix")
	}
	if !strings.Contains(err.Error(), "name must end with .schema.json") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsMissingStatus(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "", Purpose: "A"},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for missing status")
	}
	if !strings.Contains(err.Error(), "missing status") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsInvalidStatus(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "draft", Purpose: "A"},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
	if !strings.Contains(err.Error(), "invalid status") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckPassesNotAssessedStatus(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: statusNotAssessed, Purpose: "A"},
		},
	}
	if err := check(root, idx); err != nil {
		t.Fatalf("check: %v", err)
	}
}

func TestCheckFailsMissingPurpose(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: ""},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for missing purpose")
	}
	if !strings.Contains(err.Error(), "missing purpose") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsInvalidExampleCoverage(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A", ExampleCoverage: "unknown"},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for invalid example_coverage")
	}
	if !strings.Contains(err.Error(), "invalid example_coverage") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsBrokenExampleRef(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A", ExampleCoverage: examplePresent, Examples: []string{"examples/nope.json"}},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for broken example ref")
	}
	if !strings.Contains(err.Error(), "broken example ref") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckFailsPresentWithoutExamples(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A", ExampleCoverage: examplePresent},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for present coverage without examples")
	}
	if !strings.Contains(err.Error(), "example_coverage is present but examples list is empty") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckPassesValidExampleRef(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "examples"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "examples", "ok.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A", ExampleCoverage: examplePresent, Examples: []string{"examples/ok.json"}},
		},
	}
	if err := check(root, idx); err != nil {
		t.Fatalf("check: %v", err)
	}
}

func TestCheckFailsUnsupportedVersion(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "index.json")
	if err := os.WriteFile(path, []byte(`{"version":"2","schemas":[]}`), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := readIndex(path)
	if err == nil {
		t.Fatal("expected error for unsupported version")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestGenerateTableIsDeterministic(t *testing.T) {
	idx := &Index{
		Schemas: []SchemaEntry{
			{Name: "b.schema.json", Status: "current", Purpose: "B"},
			{Name: "a.schema.json", Status: "current", Purpose: "A"},
		},
	}
	out1 := generateTable(idx)
	out2 := generateTable(idx)
	if out1 != out2 {
		t.Fatalf("generateTable not deterministic:\n%s\n!=\n%s", out1, out2)
	}
	lines := strings.Split(out1, "\n")
	if len(lines) < 4 {
		t.Fatalf("expected header + rows, got:\n%s", out1)
	}
}

func TestGenerateTableEscapesPipes(t *testing.T) {
	idx := &Index{
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A | B"},
		},
	}
	out := generateTable(idx)
	if strings.Contains(out, "| A | B |") {
		t.Fatalf("pipe not escaped:\n%s", out)
	}
	if !strings.Contains(out, "A \\| B") {
		t.Fatalf("pipe escape missing:\n%s", out)
	}
}

func TestCheckReadmePassesWhenSynchronized(t *testing.T) {
	root := t.TempDir()
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A"},
		},
	}
	readme := "# Schema Reference\n\n## Schemas\n\n<!-- schemadoc-start -->\n" + generateTable(idx) + "<!-- schemadoc-end -->\n\n## Validation\n"
	path := filepath.Join(root, "README.md")
	if err := os.WriteFile(path, []byte(readme), 0644); err != nil {
		t.Fatal(err)
	}
	if err := checkReadmeAt(path, idx); err != nil {
		t.Fatalf("checkReadmeAt: %v", err)
	}
}

func TestCheckReadmeFailsWithoutMarkers(t *testing.T) {
	root := t.TempDir()
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A"},
		},
	}
	path := filepath.Join(root, "README.md")
	if err := os.WriteFile(path, []byte("# No markers here\n"), 0644); err != nil {
		t.Fatal(err)
	}
	err := checkReadmeAt(path, idx)
	if err == nil {
		t.Fatal("expected error for missing markers")
	}
	if !strings.Contains(err.Error(), "README missing schemadoc-start marker") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckReadmeFailsMarkerOrder(t *testing.T) {
	root := t.TempDir()
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A"},
		},
	}
	readme := "\n<!-- schemadoc-end -->\n\n<!-- schemadoc-start -->\n"
	path := filepath.Join(root, "README.md")
	if err := os.WriteFile(path, []byte(readme), 0644); err != nil {
		t.Fatal(err)
	}
	err := checkReadmeAt(path, idx)
	if err == nil {
		t.Fatal("expected error for invalid marker order")
	}
	if !strings.Contains(err.Error(), "invalid schemadoc marker order") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckReadmeFailsWhenDrifted(t *testing.T) {
	root := t.TempDir()
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A"},
		},
	}
	readme := "# Schema Reference\n\n## Schemas\n\n<!-- schemadoc-start -->\n| Schema | Status | Purpose |\n|---|---|---|\n| `b.schema.json` | current | B |\n<!-- schemadoc-end -->\n\n## Validation\n"
	path := filepath.Join(root, "README.md")
	if err := os.WriteFile(path, []byte(readme), 0644); err != nil {
		t.Fatal(err)
	}
	err := checkReadmeAt(path, idx)
	if err == nil {
		t.Fatal("expected error for drifted README")
	}
	if !strings.Contains(err.Error(), "README schema table drift") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestRunAcceptsCurrentIndex(t *testing.T) {
	if err := run(false, false); err != nil {
		t.Fatalf("run: %v", err)
	}
}

func TestRunGenerate(t *testing.T) {
	if err := run(true, false); err != nil {
		t.Fatalf("run generate: %v", err)
	}
}

func TestRunVerifyReadme(t *testing.T) {
	if err := run(false, true); err != nil {
		t.Fatalf("run verify-readme: %v", err)
	}
}

func TestRunVerifyReadmeFailsWhenDrifted(t *testing.T) {
	root := repoRoot()
	readme := "# Schema Reference\n\n## Schemas\n\n<!-- schemadoc-start -->\n| Schema | Status | Purpose |\n|---|---|---|\n| `b.schema.json` | current | B |\n<!-- schemadoc-end -->\n\n## Validation\n"
	path := filepath.Join(root, "schema", "README.md")
	orig, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read README: %v", err)
	}
	if err := os.WriteFile(path, []byte(readme), 0644); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.WriteFile(path, orig, 0644); err != nil {
			t.Errorf("restore README: %v", err)
		}
	}()
	if err := run(false, true); err == nil {
		t.Fatal("expected error for drifted README")
	}
}

func TestCheckFailsNotAssessedWithExamples(t *testing.T) {
	root := t.TempDir()
	sd := schemaDir(t, root)
	if err := os.WriteFile(filepath.Join(sd, "a.schema.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A", ExampleCoverage: exampleNotAssessed, Examples: []string{"examples/nope.json"}},
		},
	}
	err := check(root, idx)
	if err == nil {
		t.Fatal("expected error for not_assessed with examples")
	}
	if !strings.Contains(err.Error(), "example_coverage is \"not_assessed\" but examples list is non-empty") {
		t.Fatalf("error missing expected text: %v", err)
	}
}

func TestCheckReadmeFailsMissingFile(t *testing.T) {
	idx := &Index{
		Version: indexVersion,
		Schemas: []SchemaEntry{
			{Name: "a.schema.json", Status: "current", Purpose: "A"},
		},
	}
	err := checkReadmeAt(filepath.Join(t.TempDir(), "README.md"), idx)
	if err == nil {
		t.Fatal("expected error for missing README")
	}
	if !strings.Contains(err.Error(), "read") {
		t.Fatalf("error missing expected text: %v", err)
	}
}
